package learner

import (
	"errors"
	"fmt"
	"slices"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/learner-model/detmath"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/spec/types/formatv0"
)

// The time a learner spends, in seconds: on an item of a sheet two minutes, as in standalone, on a
// line of a workbook a minute, on an essay fifteen minutes (docs/architecture.md, section 3.1); on
// a link of a lesson that is no task — an explanation or a worked example — two minutes, an
// assumption of the bench (R15 in docs/decisions.md). Answers after a rung take no time of their
// own: the next task comes at the time the attempt's items would.
const (
	SheetItemSeconds    = 120
	WorkbookLineSeconds = 60
	EssaySeconds        = 900
	LinkSeconds         = 120
)

// Closure is why an attempt closed.
type Closure int

const (
	ClosedRight         Closure = iota + 1 // every item right, at the first answer or after a rung
	ClosedSolutionShown                    // the fourth wrong answer: the solution is shown
	ClosedDone                             // a task with a rubric, sent once and done
	ClosedLeft                             // a task with a rubric, left for another (E-leave)
)

// ItemResult is how one item of an attempt went.
type ItemResult struct {
	Place      int32
	FirstRight bool    // the first verdict of the item, the one the model reads (R06 learner-model)
	Right      bool    // whether it was right by the attempt's end
	Wrong      int     // its wrong answers before it was right, or all of them
	Unreadable bool    // its first answer was unreadable and written again (S-unreadable)
	Chance     float64 // P* of its first answer, as the truth gave it
}

// Attempt is an attempt of a learner on a task.
type Attempt struct {
	Task    string // the task's full ID
	Start   int64  // when it was given, in seconds UTC
	End     int64  // when it closed
	Seed    uint64 // the attempt's seed, for the instance of a template
	Items   []ItemResult
	Rungs   int // the rungs of the ladder it reached, 0 to 3
	Closure Closure
	// Technique tells that the learner took the technique on request before its first answer.
	Technique bool
}

// rungs is the number of rungs of the ladder of hints (concepts 10 spec).
const rungs = 3

// ErrNoAnswer is a task with no rubric and no cell to answer.
var ErrNoAnswer = errors.New("learner: the task has nothing to answer")

// Answer has the learner answer a task given at a time: the first answers of the items, one
// fraction of the place answers each in the items' order, then, after a wrong answer, again on
// each rung of the ladder, one fraction of the place retries for each item not yet right; with the
// fourth wrong answer the solution is shown — to a child by the adult's act, as by default, to a
// learner of the category adult by its own button (concepts 10 spec; R82 spec). A task with a
// rubric is written, sent once and done. The attempt's seed is the next integer of the place
// attempt. Then the stability of each primary skill of the items moves once, by the share of its
// first verdicts that were right, and its last practice becomes the attempt's start
// (docs/architecture.md, section 3.1).
func (l *Learner) Answer(t *pack.Task, at int64) (Attempt, error) {
	task := t.Candidate.Task
	if _, err := l.Module(moduleOf(task.ID, "task")); err != nil {
		return Attempt{}, err
	}
	if !t.Rubric && !slices.ContainsFunc(task.Items, func(it learnermodel.TaskItem) bool { return len(it.Slots) > 0 }) {
		return Attempt{}, fmt.Errorf("%w: %s", ErrNoAnswer, task.ID)
	}
	a := Attempt{Task: task.ID, Start: at, Seed: l.seeds.Integer()}
	if l.a.TechniqueChance > 0 {
		// Drawn on every attempt, so that the k-th attempt meets the same numbers; a task with a rubric
		// takes no technique.
		a.Technique = l.attempts.Fraction() < l.a.TechniqueChance && !t.Rubric
	}
	lift := 0.0
	if a.Technique {
		lift = l.a.TechniqueLift
	}
	// What an attempt holds comes from the place attempts in a fixed order: the technique, the shock,
	// the error (docs/architecture.md, section 3.4).
	if l.a.AttemptShock != 0 {
		lift += float64(l.a.AttemptShock * l.attempts.Normal())
	}
	holdsError := l.a.ErrorChance > 0 && l.attempts.Fraction() < l.a.ErrorChance && !t.Rubric
	solvedAt, seen := l.solved[task.ID]
	repeat := l.a.RepeatChance > 0 && seen && !task.Template && at-solvedAt <= int64(l.a.RepeatDays)*secondsPerDay
	if t.Rubric {
		a.End, a.Closure = at+EssaySeconds, ClosedDone
		if l.a.LeaveChance > 0 && l.attempts.Fraction() < l.a.LeaveChance {
			a.Closure = ClosedLeft
		}
		return a, nil
	}
	for i, it := range task.Items {
		if len(it.Slots) == 0 {
			continue
		}
		p, err := l.chance(t, i, at, 0, lift)
		if err != nil {
			return Attempt{}, err
		}
		if repeat {
			p = l.a.RepeatChance
		}
		right := l.answers.Fraction() < p
		if holdsError && hasPrimary(t, i) {
			right = false
		}
		res := ItemResult{Place: it.Place, FirstRight: right, Right: right, Wrong: b2i(!right), Chance: p}
		if holdsError && hasPrimary(t, i) {
			res.Chance = 0
		}
		res.Unreadable = l.unreadable(t)
		a.Items = append(a.Items, res)
	}
	for r := 1; r <= rungs && slices.ContainsFunc(a.Items, func(x ItemResult) bool { return !x.Right }); r++ {
		a.Rungs = r
		for k := range a.Items {
			if a.Items[k].Right {
				continue
			}
			p, err := l.chance(t, int(a.Items[k].Place-1), at, r, lift)
			if err != nil {
				return Attempt{}, err
			}
			if repeat {
				p = l.a.RepeatChance
			}
			// An answer again that could not be read is written once more, and its fraction is the one
			// that counts: it takes no verdict and adds no mistake.
			l.unreadable(t)
			if l.retries.Fraction() < p {
				a.Items[k].Right = true
			} else {
				a.Items[k].Wrong++
			}
		}
	}
	a.Closure = ClosedRight
	if slices.ContainsFunc(a.Items, func(x ItemResult) bool { return !x.Right }) {
		a.Closure = ClosedSolutionShown
	}
	per := int64(SheetItemSeconds)
	if t.Canvas == formatv0.CanvasKindWorkbook {
		per = WorkbookLineSeconds
	}
	a.End = at + per*int64(len(a.Items))
	l.practise(t, a)
	l.learn(t, a)
	// The solution goes to the device when the attempt closes right or with it shown (R82 spec).
	l.solved[task.ID] = a.End
	return a, nil
}

// practise moves the stability of each primary skill of the attempt's items once, by the share S̄
// of its first verdicts that were right: τ* := max(τ₀; τ* · (1 + g · (1 − R*))^S̄ · (1/2)^(1 − S̄)),
// R* at the attempt's start; and sets its last practice to the attempt's start (R04, R06
// learner-model; docs/architecture.md, section 3.1). The skills move in the order of their IDs.
func (l *Learner) practise(t *pack.Task, a Attempt) {
	right, all := map[string]int{}, map[string]int{}
	var ids []string
	for _, res := range a.Items {
		for _, id := range t.ItemPrimary(int(res.Place - 1)) {
			if l.skills[id] == nil {
				continue
			}
			if all[id] == 0 {
				ids = append(ids, id)
			}
			all[id]++
			right[id] += b2i(res.FirstRight)
		}
	}
	slices.Sort(ids)
	for _, id := range ids {
		s := l.skills[id]
		r := l.Retention(s, a.Start)
		share := float64(float64(right[id]) / float64(all[id]))
		g := l.a.Growth
		if l.a.TeacherFades && s.Teacher && l.TrulyMastered(id, s.TaughtOn) {
			g = 0
		}
		grow := detmath.Pow(1+float64(g*(1-r)), share)
		halve := detmath.Pow(0.5, 1-share)
		s.Stability = max(stabilityFloor, float64(float64(s.Stability*grow)*halve))
		s.Last, s.Practised = a.Start, true
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Pass is a learner's pass of a lesson, link by link to its end, or to the link where it was left.
type Pass struct {
	Lesson    string // the lesson's full ID
	Start     int64
	End       int64
	Attempts  []Attempt // the attempts on its tasks, in the author's order
	Abandoned bool      // the pass was left at a link (L-abandon)
}

// Each is what a lesson's pass calls around each of its tasks, so that the caller reads the truth
// and the model between them, as they stand when the task is given and once it is answered.
type Each struct {
	Before func(t *pack.Task, at int64) error
	After  func(t *pack.Task, a Attempt) error
}

// PassLesson has the learner pass a lesson of a module given at a time, link by link in the author's
// order to its end (docs/architecture.md, section 3.1): an explanation or a worked example takes two
// minutes and changes nothing of the truth (R15) but for L-explain, where the first one reached
// teaches the lesson's skills once; a task is answered as any task is. L-abandon leaves a pass at a
// link drawn evenly, and the next pass of the lesson goes on from that link. each, when given, is
// called around each task: Before when it is given, After once it is answered.
func (l *Learner) PassLesson(m *pack.Module, lesson *pack.Lesson, at int64, each *Each) (Pass, error) {
	if _, err := l.Module(m.ID); err != nil {
		return Pass{}, err
	}
	p := Pass{Lesson: lesson.Lesson.ID, Start: at}
	from := l.resume[lesson.Lesson.ID]
	stop := len(lesson.Order)
	if l.a.AbandonChance > 0 {
		// Two numbers of the place attempts each pass: whether it is left, and at which link of
		// those still to go.
		leave := l.attempts.Fraction() < l.a.AbandonChance
		link := from + int(l.attempts.Below(uint64(max(1, len(lesson.Order)-from))))
		if leave {
			stop, p.Abandoned = link, true
		}
	}
	for k := from; k < stop; k++ {
		i := lesson.Order[k]
		if i < 0 {
			if l.a.ExplainLift != 0 && !l.taught[lesson.Lesson.ID] {
				l.taught[lesson.Lesson.ID] = true
				for _, id := range lesson.Lesson.Skills {
					if s := l.skills[id]; s != nil {
						s.Offset += l.a.ExplainLift
					}
				}
			}
			at += LinkSeconds
			continue
		}
		t, err := m.Task(lesson.Lesson.Tasks[i].ID)
		if err != nil {
			return Pass{}, err
		}
		if each != nil {
			if err := each.Before(t, at); err != nil {
				return Pass{}, err
			}
		}
		a, err := l.Answer(t, at)
		if err != nil {
			return Pass{}, err
		}
		if each != nil {
			if err := each.After(t, a); err != nil {
				return Pass{}, err
			}
		}
		p.Attempts = append(p.Attempts, a)
		at = a.End
	}
	if p.Abandoned {
		l.resume[lesson.Lesson.ID] = stop
	} else {
		delete(l.resume, lesson.Lesson.ID)
	}
	p.End = at
	return p, nil
}

// learn moves the truth by the attempt's first answers, one by one, when the population learns or
// jumps (G2, G2-half, G2-fading; G3, G3-drop), as standalone's child does after its k-th answer
// (docs/architecture.md, sections 3.2 and 4): the module's level grows by LearnLevel and the
// offsets of the item's primary skills by LearnSkill more, both slowed by F/(F + k) when LearnFade
// is F, k counted from 0; every other skill of the module goes back LearnReturn of the way to where
// it started, unslowed; and the level jumps by JumpBy at its answer of the jump. The answers are
// the module's own.
func (l *Learner) learn(t *pack.Task, a Attempt) {
	m := l.modules[moduleOf(t.Candidate.Task.ID, "task")]
	for _, res := range a.Items {
		k := m.Answered
		m.Answered++
		if l.a.LearnLevel != 0 || l.a.LearnSkill != 0 || l.a.LearnReturn != 0 {
			rate := 1.0
			if l.a.LearnFade > 0 {
				rate = float64(l.a.LearnFade / float64(l.a.LearnFade+float64(k)))
			}
			m.Level += float64(l.a.LearnLevel * rate)
			primary := t.ItemPrimary(int(res.Place - 1))
			for _, sk := range m.Pack.Graph.Skills {
				s := l.skills[sk.ID]
				if s == nil {
					continue
				}
				if slices.Contains(primary, sk.ID) {
					s.Offset += float64(l.a.LearnSkill * rate)
				} else {
					s.Offset += float64(l.a.LearnReturn * float64(s.Start-s.Offset))
				}
			}
		}
		if m.JumpAt > 0 && m.Answered == m.JumpAt {
			m.Level += l.a.JumpBy
		}
	}
}

// hasPrimary tells whether the item at index i has one of the task's primary skills (W-error).
func hasPrimary(t *pack.Task, i int) bool {
	return slices.ContainsFunc(t.ItemPrimary(i), func(s string) bool { return slices.Contains(t.Candidate.Task.Primary, s) })
}

// unreadable counts an answer of an item on a sheet, and tells whether it is the UnreadableEvery-th,
// unreadable first and written again (S-unreadable): every answer counts, the first and those
// after a rung (docs/architecture.md, section 3.2).
func (l *Learner) unreadable(t *pack.Task) bool {
	if l.a.UnreadableEvery <= 0 || t.Canvas != formatv0.CanvasKindHandwritingSheet {
		return false
	}
	l.sheetItems++
	return l.sheetItems%l.a.UnreadableEvery == 0
}
