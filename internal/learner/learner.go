package learner

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
	"strings"
	"sync"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/learner-model/detmath"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/seed"
)

// secondsPerDay turns the integer seconds of the calendar into the days of the law of retention.
const secondsPerDay = 86400

// stabilityFloor is τ₀, in days: the stability a skill starts with and never goes below, as the
// model's (R04 learner-model).
const stabilityFloor = 1.0

// Learner is a virtual learner: its truth in each of its modules and its streams of draws.
type Learner struct {
	Unit string // <population>/<n>
	Age  string // one age category, as a profile has one (R107 spec)

	a       Assumptions
	set     seed.Set
	modules map[string]*Module
	skills  map[string]*Skill

	answers, retries, attempts, seeds *seed.Stream

	cal        *Calendar        // the learner's one calendar, once made: its bad days and its break
	sheetItems int              // answers of items on a sheet so far, for S-unreadable
	solved     map[string]int64 // when the solution of each task last went to the device
	taught     map[string]bool  // the lessons whose explanations have taught once, for L-explain
	resume     map[string]int   // the link where the open pass of each lesson was left, for L-abandon
}

// Module is the learner's truth in one module.
type Module struct {
	Pack      *pack.Module
	Grade     learnermodel.Grade
	StartBand int32   // the place of the band the learner starts in
	Level     float64 // θ*, the module's level; it is not forgotten
	JumpAt    int     // the answer of the module at which its level jumps, when the population jumps
	Answered  int     // first answers of items with a verdict in the module so far
	HardBand  int32   // the place of the band whose tasks are truly harder, when the population says so
}

// Skill is the learner's truth of one skill with a verdict.
type Skill struct {
	Module    *Module
	Offset    float64 // δ*
	Start     float64 // δ* as it was drawn, which learning goes back to
	Teacher   bool    // the skill teaches an anchor skill
	TaughtOn  int32   // the first band the graph teaches it on
	Stability float64 // τ*, in days
	Last      int64   // the time of the last attempt on it, in seconds
	Practised bool
}

// ErrNoModule is a learner given no module, or a task of a module it does not learn.
var ErrNoModule = errors.New("learner: no such module")

// ErrAssumptions is a population whose assumptions no learner can live by.
var ErrAssumptions = errors.New("learner: the assumptions cannot be lived by")

// validate refuses assumptions the calendar or the truth cannot run on.
func (a Assumptions) validate() error {
	if !(a.SessionChance > 0 && a.SessionChance <= 1) || a.SessionAttempts < 1 {
		return fmt.Errorf("%w: a session comes with the chance %v, %d attempts long", ErrAssumptions, a.SessionChance, a.SessionAttempts)
	}
	if a.JumpBy != 0 && !(a.JumpFrom >= 1 && a.JumpFrom <= a.JumpTo) {
		return fmt.Errorf("%w: a jump between answers %d and %d", ErrAssumptions, a.JumpFrom, a.JumpTo)
	}
	if a.BreakMax > 0 && !(a.BreakMin >= 1 && a.BreakMin <= a.BreakMax && a.BreakStartFrom >= 0 && a.BreakStartFrom <= a.BreakStartTo && a.BreakStartTo+a.BreakMin <= a.BreakEndBy) {
		return fmt.Errorf("%w: a break from day %d–%d, %d–%d days long, by day %d", ErrAssumptions, a.BreakStartFrom, a.BreakStartTo, a.BreakMin, a.BreakMax, a.BreakEndBy)
	}
	if a.ExponentialLaw && !(a.LawBase > 0 && a.LawBase < 1) {
		return fmt.Errorf("%w: a law of base %v", ErrAssumptions, a.LawBase)
	}
	if a.Adult && a.AdultFrom > a.AdultTo {
		return fmt.Errorf("%w: an adult %v to %v above the first band", ErrAssumptions, a.AdultFrom, a.AdultTo)
	}
	if !(a.LevelCorrelation >= -1 && a.LevelCorrelation <= 1) {
		return fmt.Errorf("%w: a correlation of %v", ErrAssumptions, a.LevelCorrelation)
	}
	for _, c := range []float64{a.ErrorChance, a.BadDayChance, a.RepeatChance, a.AbandonChance, a.LeaveChance} {
		if !(c >= 0 && c <= 1) {
			return fmt.Errorf("%w: a chance of %v", ErrAssumptions, c)
		}
	}
	if a.UnreadableEvery < 0 || a.RepeatDays < 0 || a.AttemptShock < 0 {
		return fmt.Errorf("%w: unreadable every %d, a repeat within %d days, a shock of %v", ErrAssumptions, a.UnreadableEvery, a.RepeatDays, a.AttemptShock)
	}
	if !(a.TechniqueChance >= 0 && a.TechniqueChance <= 1) || a.LearnFade < 0 {
		return fmt.Errorf("%w: the technique with the chance %v, a fade of %v", ErrAssumptions, a.TechniqueChance, a.LearnFade)
	}
	return nil
}

// New draws the n-th learner of a population of the set, in the given modules, from its places
// (docs/architecture.md, sections 3.1 and 3.4). A learner has one grade and one age category, as a
// profile has (R107 spec). Its truth comes from the place params, in a fixed order: the grade,
// among those every module with grades names, when they are more than one; for each module by ID,
// the level and the offset of each skill with a verdict in the order of the graph; then, when no
// module gives a grade, the age category, among those every module names; and last what only some
// populations draw, for each module by ID — the answer of the jump, the band that is harder —, so
// that a population that draws more moves none of the others (R18). Its calendar, answers,
// retries, attempts and attempt seeds come from places of their own, so that one never moves
// another.
func New(set seed.Set, population string, n int, modules []*pack.Module, a Assumptions) (*Learner, error) {
	if len(modules) == 0 {
		return nil, ErrNoModule
	}
	if err := a.validate(); err != nil {
		return nil, err
	}
	unit := seed.Learner(population, n)
	place := func(purpose string) *seed.Stream {
		s, err := set.Place(unit, purpose)
		if err != nil {
			panic(err) // the names are the bench's own: a population with | in it is a bug
		}
		return s
	}
	params := place(seed.Params)
	l := &Learner{
		Unit: unit, a: a, set: set, modules: map[string]*Module{}, skills: map[string]*Skill{}, solved: map[string]int64{}, taught: map[string]bool{}, resume: map[string]int{},
		answers: place(seed.Answers), retries: place(seed.Retries), attempts: place(seed.Attempts), seeds: place(seed.Attempt),
	}
	sorted := slices.SortedFunc(slices.Values(modules), func(x, y *pack.Module) int { return strings.Compare(x.ID, y.ID) })
	var grade learnermodel.Grade
	var grades []int32
	gradedModules := 0
	for _, pm := range sorted {
		if a.Adult {
			break
		}
		g := gradesOf(pm)
		if len(g) == 0 {
			continue
		}
		if gradedModules == 0 {
			grades = g
		} else {
			grades = slices.DeleteFunc(grades, func(x int32) bool { return !slices.Contains(g, x) })
		}
		gradedModules++
	}
	if gradedModules > 0 {
		if len(grades) == 0 {
			return nil, fmt.Errorf("%w: its modules name no grade in common", ErrAssumptions)
		}
		g := grades[0]
		if len(grades) > 1 {
			g = grades[params.Below(uint64(len(grades)))]
		}
		grade = learnermodel.Grade{Value: g, Known: true}
	}
	var first *float64 // the first module's z, which A-linked ties the others to
	for _, pm := range sorted {
		m := &Module{Pack: pm, StartBand: 1}
		if len(gradesOf(pm)) > 0 && !a.Adult {
			m.Grade = grade
			band, err := learnermodel.BandByGrade(m.Grade, pm.Graph.Bands)
			if err != nil {
				return nil, fmt.Errorf("learner: %s: %w", pm.ID, err)
			}
			m.StartBand = band.Place
		}
		z := params.Normal()
		if a.LevelCorrelation != 0 {
			if first == nil {
				first = &z
			} else {
				z = float64(a.LevelCorrelation**first) + float64(math.Sqrt(1-float64(a.LevelCorrelation*a.LevelCorrelation))*z)
			}
		}
		m.Level = float64(a.BandStep*float64(m.StartBand-1)) + float64(a.StartSpread*z) + a.LevelShift
		var verdict []string
		for _, sk := range pm.Graph.Skills {
			if sk.Assessability == learnermodel.AssessabilityVerdict {
				verdict = append(verdict, sk.ID)
			}
		}
		for k, id := range verdict {
			off := float64(a.SkillSpread * params.Normal())
			teacher := teaches(pm, id)
			if a.TeacherBoost != 0 && teacher {
				off += a.TeacherBoost
			}
			if a.SkillHalves != 0 {
				if 2*k < len(verdict) {
					off += a.SkillHalves
				} else {
					off -= a.SkillHalves
				}
			}
			l.skills[id] = &Skill{Module: m, Offset: off, Start: off, Stability: stabilityFloor, Teacher: teacher, TaughtOn: taughtOn(pm, id)}
		}
		l.modules[pm.ID] = m
	}
	switch {
	case a.Adult:
		l.Age = "adult"
		for _, pm := range sorted {
			if !slices.Contains(pm.Ages, "adult") {
				return nil, fmt.Errorf("%w: an adult in %s, whose audience has no adult", ErrAssumptions, pm.ID)
			}
		}
	case grade.Known:
		l.Age = "under_13"
	default:
		ages := slices.Clone(sorted[0].Ages)
		for _, pm := range sorted[1:] {
			ages = slices.DeleteFunc(ages, func(x string) bool { return !slices.Contains(pm.Ages, x) })
		}
		if len(ages) == 0 {
			return nil, fmt.Errorf("%w: its modules name no age category in common", ErrAssumptions)
		}
		l.Age = ages[params.Below(uint64(len(ages)))]
	}
	// What only some populations draw comes after what every learner draws, so that a population
	// that adds a draw moves none of the others: the answer of the jump, the band that is harder.
	for _, pm := range sorted {
		m := l.modules[pm.ID]
		if a.JumpBy != 0 {
			m.JumpAt = int(params.Between(int64(a.JumpFrom), int64(a.JumpTo)))
		}
		if a.BandBias != 0 && len(pm.Graph.Bands) > 0 {
			m.HardBand = int32(params.Below(uint64(len(pm.Graph.Bands)))) + 1
		}
		if a.Adult {
			m.Level += a.AdultFrom + float64(float64(a.AdultTo-a.AdultFrom)*params.Fraction())
		}
	}
	return l, nil
}

// gradesOf gives the grades the bands of a module name, ascending and once each.
func gradesOf(m *pack.Module) []int32 {
	var out []int32
	for _, b := range m.Graph.Bands {
		out = append(out, b.Grades...)
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// Module gives the learner's truth in a module.
func (l *Learner) Module(id string) (*Module, error) {
	m, ok := l.modules[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoModule, id)
	}
	return m, nil
}

// Skill gives the learner's truth of a skill with a verdict, or nil.
func (l *Learner) Skill(id string) *Skill {
	return l.skills[id]
}

// moduleOf gives the module of an entry by its full ID, <module>/<kind>/<local>.
func moduleOf(full, kind string) string {
	if i := strings.Index(full, "/"+kind+"/"); i >= 0 {
		return full[:i]
	}
	return ""
}

// Retention gives R* of a skill at a time: the law of retention since its last attempt, at its
// stability; 1 before its first attempt, and always when the population does not forget.
func (l *Learner) Retention(s *Skill, at int64) float64 {
	if !l.a.Forgets || !s.Practised {
		return 1
	}
	days := float64(float64(at-s.Last) / secondsPerDay)
	if days <= 0 {
		return 1
	}
	if l.a.ExponentialLaw {
		return detmath.Pow(l.a.LawBase, float64(days/s.Stability))
	}
	return detmath.Pow(1+float64(float64(l.a.LawF*days)/s.Stability), -l.a.LawC)
}

// Difficulty gives β* of a task: the true difficulty of its band and its difficulty, with the
// author's miss, z of the task's place times the population's Miss (docs/architecture.md, section
// 3.3). z is drawn once a task and a run: the place of the task's unit and the purpose miss.
func (l *Learner) Difficulty(t *pack.Task) float64 {
	task := t.Candidate.Task
	z := missOf(l.set, moduleOf(task.ID, "task"), t.Local)
	beta := float64(l.a.BandStep*float64(task.Band.Place-1)) + float64(task.Difficulty-3) + float64(l.a.Miss*z)
	if m := l.modules[moduleOf(task.ID, "task")]; m != nil && m.HardBand == task.Band.Place {
		beta += l.a.BandBias
	}
	return beta
}

// Chance gives P* of an item of a task at a time and a rung of the ladder: c* + (1 − c*) · R̄* ·
// σ(a · (θ* + δ̄* − β*) + lift · rung), c* the product over the item's cells of their guess, δ̄*
// and R̄* the means over the item's primary skills with a verdict (docs/architecture.md, section
// 3.1).
func (l *Learner) Chance(t *pack.Task, item int, at int64, rung int) (float64, error) {
	return l.chance(t, item, at, rung, 0)
}

// chance is Chance with lift more on the logit: the technique taken on request (G8).
func (l *Learner) chance(t *pack.Task, item int, at int64, rung int, lift float64) (float64, error) {
	task := t.Candidate.Task
	m, err := l.Module(moduleOf(task.ID, "task"))
	if err != nil {
		return 0, err
	}
	it := task.Items[item]
	primary := t.ItemPrimary(item)
	var offset, retention float64
	n := 0
	for _, id := range primary {
		s := l.skills[id]
		if s == nil {
			continue
		}
		offset += s.Offset
		retention += l.Retention(s, at)
		n++
	}
	if n > 0 {
		offset, retention = float64(offset/float64(n)), float64(retention/float64(n))
	} else {
		retention = 1
	}
	level := l.level(m, at)
	if l.cal != nil && l.a.BadDayDrop != 0 && l.cal.Bad(at) {
		level -= l.a.BadDayDrop
	}
	logit := float64(l.a.Slope*float64(level+offset-l.Difficulty(t))) + float64(l.a.RungLift*float64(rung)) + lift
	c := l.guess(it.Slots)
	return c + float64(float64(float64(1-c)*retention)*detmath.Sigmoid(logit)), nil
}

// guess gives c* of an item: the product over its cells of the chance of guessing each — a choice
// of k options ChoiceGuess/k, capped at 1, a predicate PredicateGuess, and 0 for the rest.
func (l *Learner) guess(slots []learnermodel.Slot) float64 {
	c := 1.0
	for _, s := range slots {
		switch {
		case s.Check == learnermodel.CheckKindChoice && s.Options > 0:
			c = float64(c * min(1, float64(l.a.ChoiceGuess/float64(s.Options))))
		case s.Check == learnermodel.CheckKindPredicate:
			c = float64(c * l.a.PredicateGuess)
		default:
			return 0
		}
	}
	return c
}

// misses holds z of each task and set: the author's miss is the task's, one a run, and every
// learner of the set sees the same (docs/architecture.md, section 3.3), so it is drawn once.
var misses = struct {
	sync.Mutex
	z map[missKey]float64
}{z: map[missKey]float64{}}

type missKey struct {
	set          seed.Set
	module, task string
}

// missOf gives z of a task in a set: the normal number of the place of the task's unit and the
// purpose miss.
func missOf(set seed.Set, module, task string) float64 {
	misses.Lock()
	defer misses.Unlock()
	k := missKey{set, module, task}
	z, ok := misses.z[k]
	if !ok {
		s, err := set.Place(seed.Task(module, task), seed.Miss)
		if err != nil {
			panic(err) // a task ID with | is refused by the format
		}
		z = s.Normal()
		misses.z[k] = z
	}
	return z
}

// SolvedAt tells when the solution of a task last went to the device, if it did: when an attempt on
// it closed right or with the solution shown (R82 spec). The run reads it for the evidence and the
// summary, and W-repeat for its chance, so the fact has one source.
func (l *Learner) SolvedAt(task string) (int64, bool) {
	at, ok := l.solved[task]
	return at, ok
}

// Solved gives the tasks whose solution went to the device, and when it last did, in the order of
// their IDs.
func (l *Learner) Solved() ([]string, []int64) {
	ids := slices.Sorted(maps.Keys(l.solved))
	ats := make([]int64, len(ids))
	for i, id := range ids {
		ats[i] = l.solved[id]
	}
	return ids, ats
}

// teaches tells whether a skill of a module teaches an anchor skill (R50 spec).
func teaches(m *pack.Module, id string) bool {
	i := slices.IndexFunc(m.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == id })
	return i >= 0 && len(m.Graph.Skills[i].Teaches) > 0
}

// Mastery is the middle of the corridor of success, 0.775: the truth has mastered a skill on a band
// when P* of the reference task reaches it (docs/architecture.md, section 3.5).
const Mastery = 0.775

// TrulyMastered tells whether the truth has mastered a skill on a band: P* of the reference task —
// difficulty 3 on that band, no miss, a choice of five options, R* = 1 — is at least 0.775, as
// "truly mastered" in standalone (docs/architecture.md, section 3.5). It reads no number of the
// model: not its threshold, not its guess.
func (l *Learner) TrulyMastered(id string, band int32) bool {
	return l.Reference(id, band) >= Mastery
}

// TrulyMasteredAt is TrulyMastered at a time: with the level the truth has then, fallen after the
// break of F-break. A bad day is not counted: it lowers the answers of that day, not what the learner
// knows (docs/architecture.md, section 3.5; R22).
func (l *Learner) TrulyMasteredAt(id string, band int32, at int64) bool {
	s := l.skills[id]
	return s != nil && l.reference(s, band, l.level(s.Module, at)) >= Mastery
}

// level gives θ* of a module at a time: fallen by BreakDrop after the break (F-break).
func (l *Learner) level(m *Module, at int64) float64 {
	if l.cal != nil && l.a.BreakDrop != 0 && l.cal.After(at) {
		return m.Level - l.a.BreakDrop
	}
	return m.Level
}

// LevelAt gives θ* of a module at a time, as level does.
func (l *Learner) LevelAt(module string, at int64) (float64, error) {
	m, err := l.Module(module)
	if err != nil {
		return 0, err
	}
	return l.level(m, at), nil
}

// Reference gives P* of the reference task of a skill on a band: c + (1 − c) · σ(a · (θ* + δ* −
// β)), c the guess of a choice of five options, β the band's true difficulty at difficulty 3, with
// no miss — the author's nor that of a whole band — and R* = 1. A skill the learner has no truth of
// gives 0.
func (l *Learner) Reference(id string, band int32) float64 {
	s := l.skills[id]
	if s == nil {
		return 0
	}
	return l.reference(s, band, s.Module.Level)
}

// reference gives P* of the reference task of a skill on a band at a level of its module.
func (l *Learner) reference(s *Skill, band int32, level float64) float64 {
	beta := float64(l.a.BandStep * float64(band-1))
	c := l.guess([]learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 5, HasOptions: true}})
	return c + float64(float64(1-c)*detmath.Sigmoid(float64(l.a.Slope*float64(level+s.Offset-beta))))
}

// taughtOn gives the band a skill is first taught on in its module's graph, the first if none.
func taughtOn(m *pack.Module, id string) int32 {
	i := slices.IndexFunc(m.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == id })
	if i < 0 || len(m.Graph.Skills[i].BandsTaught) == 0 {
		return 1
	}
	return m.Graph.Skills[i].BandsTaught[0]
}

// Truth gives θ* + δ* of a skill with a verdict at a time, with the level fallen after the break of
// F-break, and its R* then; false for a skill the learner has no truth of.
func (l *Learner) Truth(id string, at int64) (float64, float64, bool) {
	s := l.skills[id]
	if s == nil {
		return 0, 0, false
	}
	return l.level(s.Module, at) + s.Offset, l.Retention(s, at), true
}

// TaughtOn gives the first band the graph teaches a skill on; 1 for a skill with no truth.
func (l *Learner) TaughtOn(id string) int32 {
	if s := l.skills[id]; s != nil {
		return s.TaughtOn
	}
	return 1
}

// ReferenceAt gives P* of the reference task of a skill on a band at a time: with the level then, as
// TrulyMasteredAt reads it, and the skill's R* then, c + (1 − c) · R* · σ (docs/architecture.md,
// section 3.5, M5).
func (l *Learner) ReferenceAt(id string, band int32, at int64) float64 {
	s := l.skills[id]
	if s == nil {
		return 0
	}
	c := l.guess([]learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 5, HasOptions: true}})
	return c + float64(float64(l.reference(s, band, l.level(s.Module, at))-c)*l.Retention(s, at))
}
