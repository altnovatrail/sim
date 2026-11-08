package learner

import (
	"errors"
	"fmt"
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
}

// Module is the learner's truth in one module.
type Module struct {
	Pack      *pack.Module
	Grade     learnermodel.Grade
	StartBand int32   // the place of the band the learner starts in
	Level     float64 // θ*, the module's level; it is not forgotten
}

// Skill is the learner's truth of one skill with a verdict.
type Skill struct {
	Module    *Module
	Offset    float64 // δ*
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
	return nil
}

// New draws the n-th learner of a population of the set, in the given modules, from its places
// (docs/architecture.md, sections 3.1 and 3.4). A learner has one grade and one age category, as a
// profile has (R107 spec). Its truth comes from the place params, in a fixed order: the grade,
// among those every module with grades names, when they are more than one; for each module by ID,
// the level and the offset of each skill with a verdict in the order of the graph; then, when no
// module gives a grade, the age category, among those every module names. Its calendar, answers,
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
		Unit: unit, a: a, set: set, modules: map[string]*Module{}, skills: map[string]*Skill{},
		answers: place(seed.Answers), retries: place(seed.Retries), attempts: place(seed.Attempts), seeds: place(seed.Attempt),
	}
	sorted := slices.SortedFunc(slices.Values(modules), func(x, y *pack.Module) int { return strings.Compare(x.ID, y.ID) })
	var grade learnermodel.Grade
	var grades []int32
	gradedModules := 0
	for _, pm := range sorted {
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
	for _, pm := range sorted {
		m := &Module{Pack: pm, StartBand: 1}
		if len(gradesOf(pm)) > 0 {
			m.Grade = grade
			band, err := learnermodel.BandByGrade(m.Grade, pm.Graph.Bands)
			if err != nil {
				return nil, fmt.Errorf("learner: %s: %w", pm.ID, err)
			}
			m.StartBand = band.Place
		}
		m.Level = float64(a.BandStep*float64(m.StartBand-1)) + float64(a.StartSpread*params.Normal())
		for _, sk := range pm.Graph.Skills {
			if sk.Assessability != learnermodel.AssessabilityVerdict {
				continue
			}
			l.skills[sk.ID] = &Skill{Module: m, Offset: float64(a.SkillSpread * params.Normal()), Stability: stabilityFloor}
		}
		l.modules[pm.ID] = m
	}
	if grade.Known {
		l.Age = "under_13"
		return l, nil
	}
	ages := slices.Clone(sorted[0].Ages)
	for _, pm := range sorted[1:] {
		ages = slices.DeleteFunc(ages, func(x string) bool { return !slices.Contains(pm.Ages, x) })
	}
	if len(ages) == 0 {
		return nil, fmt.Errorf("%w: its modules name no age category in common", ErrAssumptions)
	}
	l.Age = ages[params.Below(uint64(len(ages)))]
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
	return detmath.Pow(1+float64(float64(l.a.LawF*days)/s.Stability), -l.a.LawC)
}

// Difficulty gives β* of a task: the true difficulty of its band and its difficulty, with the
// author's miss, z of the task's place times the population's Miss (docs/architecture.md, section
// 3.3). z is drawn once a task and a run: the place of the task's unit and the purpose miss.
func (l *Learner) Difficulty(t *pack.Task) float64 {
	task := t.Candidate.Task
	z := missOf(l.set, moduleOf(task.ID, "task"), t.Local)
	return float64(l.a.BandStep*float64(task.Band.Place-1)) + float64(task.Difficulty-3) + float64(l.a.Miss*z)
}

// Chance gives P* of an item of a task at a time and a rung of the ladder: c* + (1 − c*) · R̄* ·
// σ(a · (θ* + δ̄* − β*) + lift · rung), c* the product over the item's cells of their guess, δ̄*
// and R̄* the means over the item's primary skills with a verdict (docs/architecture.md, section
// 3.1).
func (l *Learner) Chance(t *pack.Task, item int, at int64, rung int) (float64, error) {
	task := t.Candidate.Task
	m, err := l.Module(moduleOf(task.ID, "task"))
	if err != nil {
		return 0, err
	}
	it := task.Items[item]
	primary := task.Primary
	if it.OwnSkills {
		primary = it.Primary
	}
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
	logit := float64(l.a.Slope*float64(m.Level+offset-l.Difficulty(t))) + float64(l.a.RungLift*float64(rung))
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
