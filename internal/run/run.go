// Package run runs a virtual learner through learner-model and engine, as the platform runs a
// learner through them: at each time its calendar gives, the bench builds the engine's input as
// api would — the model's reading, the pool and the lessons the learner's age category suits, the
// runs of lessons, the attempts closed by "Done", the last closed attempt and the summary of what
// was handed out —, takes the engine's choice, has the learner answer by its truth, and gives the
// model the evidence (docs/architecture.md, sections 1 and 3; architecture 3.8 spec; R104, R107
// spec; R16 in docs/decisions.md). Within a lesson whose skills all have a verdict it calls the
// engine at the end of each task link, as api will, and plays its answer: further through the
// lesson, or a task of consolidation (R48). It records each step, for the invariants and the
// measures.
package run

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
)

// Limit is when a learner's run ends: after Answers first answers of items with a verdict, when
// Answers is above 0; after Days days of the calendar; and, when Rubrics, once every task with a
// rubric the learner's age category suits is done (docs/architecture.md, section 3.1).
type Limit struct {
	Answers int
	Days    int
	Rubrics bool
}

// Maths is the run of a learner in mathematics: 200 answers of items with a verdict, as 200
// answers in standalone, within a year of the calendar, so that a rule that gives nothing for
// good still ends (R16).
var Maths = Limit{Answers: 200, Days: 365}

// Forgetting is the run of the populations of forgetting: 180 days of the calendar.
var Forgetting = Limit{Days: 180}

// Essays is the run of a learner of essays: until the tasks with a rubric it is suited run out,
// within 60 days.
var Essays = Limit{Days: 60, Rubrics: true}

// Step is one choice of the run and what came of it. Evidence and Instances stand beside
// Attempts, one for each. A lesson with a task of consolidation is several steps, in the order of
// time: the lesson up to the consolidation, the consolidation — a choice of a task —, and the
// lesson after it, Resumed (R48).
type Step struct {
	Input    engine.Input // the input as the engine got it; the choice is the same on it again
	Choice   engine.Choice
	Attempts []learner.Attempt // the attempts the choice gave: one for a task, a lesson's tasks
	// Evidence is the evidence the model got of each attempt, as Append gave it back; nil for an
	// attempt of a task with a rubric, which gives none.
	Evidence []*learnermodel.Evidence
	// Instances is the number of the instance of each attempt's template, chosen by the attempt's
	// seed as the platform chooses it; -1 for a task that is no template.
	Instances []int
	// Obs is what the measures read of each attempt with a verdict: the model and the truth before
	// and after it; nil for an attempt of a task with a rubric (R22 in docs/decisions.md).
	Obs  []*Obs
	Pass *learner.Pass // the pass of a lesson, on the step that ends it, with all its attempts
	// Pool is what the truth says of the pool at the choice, before the engine chose (R23); nil on a
	// step Resumed.
	Pool *Pool
	// Further are the calls of the engine at the end of the step's task links answered with the
	// lesson of the same run: further through the lesson, no new choice (R48).
	Further []Call
	// Resumed tells that the step is a lesson going on after a task of consolidation: its choice is
	// the engine's answer after the consolidation, the lesson of the same run, and no new choice.
	Resumed bool
}

// Call is a call of the engine within a lesson: the input it got and its answer.
type Call struct {
	Input  engine.Input
	Choice engine.Choice
}

// Record is a learner's run.
type Record struct {
	Unit    string
	Learner *learner.Learner
	Steps   []Step
	State   learnermodel.State // the model's state at the end
	Answers int                // first answers of items with a verdict
	Days    int                // the day of the calendar of the last choice, from 1
}

// Run is a learner's run in progress: the model's state and what api keeps of the learner.
type Run struct {
	l        *learner.Learner
	cal      *learner.Calendar
	modules  []*pack.Module  // the modules of the engine: the learner's, by ID
	suited   []engine.Module // their graphs, and the pool and the lessons the age category suits
	rubrics  map[string]bool // the tasks with a rubric of the pools
	graphs   learnermodel.Graphs
	grade    learnermodel.Grade
	limit    Limit
	rule     Rule
	ended    bool     // the calendar went past the limit's last day
	verdicts []string // the skills with a verdict of the modules, by full ID, for the ceiling
	// choose is the engine's rule: engine.Choose, or a test's.
	choose func(engine.Input) (engine.Choice, error)

	state      learnermodel.State
	runs       []engine.Run
	runStarts  []int64 // when each run began, for a run left open
	rubricDone []engine.RubricDone
	last       engine.Last
	hasLast    bool
	issued     map[string]engine.SkillIssue
	choices    map[string]int64

	rec Record
}

// ErrRun is a run that cannot go on: an input the engine or the model refuses, or a choice the
// bench cannot follow.
var ErrRun = errors.New("run: the run cannot go on")

// New starts the run of a learner in the given modules, the same the learner was drawn in. Every
// limit has its days, so that every run ends.
func New(l *learner.Learner, modules []*pack.Module, limit Limit) (*Run, error) {
	if len(modules) == 0 {
		return nil, fmt.Errorf("%w: no module", ErrRun)
	}
	if limit.Days <= 0 {
		return nil, fmt.Errorf("%w: a limit with no days, %+v, may never end", ErrRun, limit)
	}
	r := &Run{
		l: l, cal: l.NewCalendar(), graphs: learnermodel.Graphs{}, limit: limit, rubrics: map[string]bool{}, choose: engine.Choose,
		issued: map[string]engine.SkillIssue{}, choices: map[string]int64{},
		rec: Record{Unit: l.Unit, Learner: l},
	}
	r.modules = slices.SortedFunc(slices.Values(modules), func(x, y *pack.Module) int { return strings.Compare(x.ID, y.ID) })
	for _, m := range r.modules {
		pool, lessons := m.ForAge(l.Age)
		r.suited = append(r.suited, engine.Module{Graph: m.Graph, Pool: pool, Lessons: lessons})
		for _, c := range pool {
			if t, err := m.Task(c.Task.ID); err == nil && t.Rubric {
				r.rubrics[c.Task.ID] = true
			}
		}
		lm, err := l.Module(m.ID)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRun, err)
		}
		if lm.Grade.Known {
			r.grade = lm.Grade
		}
		r.graphs[m.ID] = learnermodel.ModuleGraph{Graph: m.Graph, Current: true}
		for _, s := range m.Graph.Skills {
			if s.Assessability == learnermodel.AssessabilityVerdict {
				r.verdicts = append(r.verdicts, s.ID)
			}
		}
	}
	return r, nil
}

// Use sets the rule that chooses in the run; a new run chooses by the engine.
func (r *Run) Use(rule Rule) *Run {
	r.rule = rule
	return r
}

// All runs the learner to the end of its limit and gives the record; on an error, the record up
// to the step that failed.
func (r *Run) All() (*Record, error) {
	for !r.Done() {
		if err := r.Step(); err != nil {
			return &r.rec, err
		}
	}
	return &r.rec, nil
}

// Record gives the record so far: the steps taken and the model's state after them.
func (r *Run) Record() *Record {
	return &r.rec
}

// Done tells whether the run has reached its limit.
func (r *Run) Done() bool {
	if r.ended || r.limit.Answers > 0 && r.rec.Answers >= r.limit.Answers {
		return true
	}
	return r.limit.Rubrics && r.rubricsDone()
}

// rubricsDone tells whether every task with a rubric the learner is suited is done.
func (r *Run) rubricsDone() bool {
	done := 0
	for _, d := range r.rubricDone {
		if r.rubrics[d.Task] {
			done++
		}
	}
	return done >= len(r.rubrics)
}

// abandonAfter is how long a run of a lesson left open stays open: then it ends "abandoned"
// (docs/architecture.md, section 3.2, L-abandon).
const abandonAfter = 30 * 86400

// Step takes one choice at the time the calendar gives, and follows it. A time past the limit's
// last day ends the run instead. Before the choice, a run left open for 30 days ends abandoned.
func (r *Run) Step() error {
	if r.Done() {
		return nil
	}
	at := r.cal.Next()
	day := r.cal.Day() + 1
	if day > r.limit.Days {
		r.ended = true
		return nil
	}
	r.rec.Days = day
	for i := range r.runs {
		if r.runs[i].End == engine.RunOpen && at >= r.runStarts[i]+abandonAfter {
			r.runs[i].End, r.runs[i].EndedAt = engine.RunAbandoned, r.runStarts[i]+abandonAfter
		}
	}
	in, err := r.input(at)
	if err != nil {
		return err
	}
	pool, err := r.pool(at)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRun, err)
	}
	choose := r.choose
	if r.rule == Ceiling {
		choose = func(in engine.Input) (engine.Choice, error) { return r.ceiling(in, pool) }
	}
	choice, err := choose(in)
	if err != nil {
		return fmt.Errorf("%w: the engine refuses the input at %d: %w", ErrRun, at, err)
	}
	step := Step{Input: in, Choice: choice, Pool: pool}
	steps := []Step{step}
	switch choice.Kind {
	case engine.ChoiceNothing:
		r.cal.EndSession()
	case engine.ChoiceTask:
		end, err := r.give(&steps[0], at)
		if err != nil {
			return err
		}
		r.cal.Done(end)
	case engine.ChoiceLesson:
		var end int64
		if steps, end, err = r.passLesson(step, at); err != nil {
			return err
		}
		r.cal.Done(end)
	default:
		return fmt.Errorf("%w: the engine chose %v", ErrRun, choice.Kind)
	}
	r.rec.Steps = append(r.rec.Steps, steps...)
	r.rec.State = r.state
	return nil
}

// give follows the choice of a task at a time: the learner answers it, the attempt is closed, and
// the module has one choice more. It gives when the attempt closed.
func (r *Run) give(step *Step, at int64) (int64, error) {
	m, t, err := r.task(step.Choice.ID)
	if err != nil {
		return 0, err
	}
	_, seen := r.l.SolvedAt(t.Candidate.Task.ID)
	before, err := r.observeBefore(m, t, at, seen, &step.Input.Readout)
	if err != nil {
		return 0, err
	}
	a, err := r.l.Answer(t, at)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrRun, err)
	}
	e, err := r.close(m, t, a, step.Choice.Reason, "", at, seen)
	if err != nil {
		return 0, err
	}
	o, err := r.observeAfter(m, t, a, before)
	if err != nil {
		return 0, err
	}
	step.Obs = []*Obs{o}
	n, err := instance(m, t, a)
	if err != nil {
		return 0, err
	}
	step.Attempts, step.Evidence, step.Instances = []learner.Attempt{a}, []*learnermodel.Evidence{e}, []int{n}
	r.choices[m.ID]++
	return a.End, nil
}

// passLesson follows the choice of a lesson at a time: the module has one choice more and each skill
// of the lesson is handed out then, the run of the lesson is open, and the learner passes it link by
// link. At the end of each task link of a lesson whose skills all have a verdict the engine is
// called, as api will call it (R48): the lesson of the same run goes on, and a task outside the
// lesson is a consolidation, given at once. The ceiling calls nothing there and goes on. It gives the
// steps the choice became: the lesson, or, with consolidations, its parts and the consolidations
// between them, in the order of time, and when the pass ended.
func (r *Run) passLesson(step Step, at int64) ([]Step, int64, error) {
	m, lesson, err := r.lesson(step.Choice.ID)
	if err != nil {
		return nil, 0, err
	}
	seen := map[string]bool{}
	for _, lt := range lesson.Lesson.Tasks {
		_, seen[lt.ID] = r.l.SolvedAt(lt.ID)
	}
	r.choices[m.ID]++
	r.issue(lesson.Lesson.Skills, at)
	r.open(lesson, at)
	consolidates := r.rule == Engine && verdicts(m, lesson)
	var steps []Step
	cur := &step
	// Each task of the lesson is closed as it is answered, so that the model and the truth are read
	// before and after it as they stand then, with the tasks before it and not those after.
	var before *Obs
	each := &learner.Each{
		Before: func(t *pack.Task, at int64) (err error) {
			before, err = r.observeBefore(m, t, at, seen[t.Candidate.Task.ID], nil)
			return err
		},
		After: func(t *pack.Task, a learner.Attempt) error {
			e, err := r.close(m, t, a, step.Choice.Reason, lesson.Lesson.ID, at, seen[a.Task])
			if err != nil {
				return err
			}
			n, err := instance(m, t, a)
			if err != nil {
				return err
			}
			o, err := r.observeAfter(m, t, a, before)
			if err != nil {
				return err
			}
			cur.Attempts, cur.Evidence, cur.Instances, cur.Obs = append(cur.Attempts, a), append(cur.Evidence, e), append(cur.Instances, n), append(cur.Obs, o)
			return nil
		},
		End: func(_ *pack.Task, a learner.Attempt) (int64, error) {
			if !consolidates {
				return a.End, nil
			}
			in, c, err := r.call(a.End)
			if err != nil {
				return 0, err
			}
			if further(c, lesson) {
				cur.Further = append(cur.Further, Call{Input: in, Choice: c})
				return a.End, nil
			}
			task, back, err := r.consolidate(lesson, in, c)
			if err != nil {
				return 0, err
			}
			steps = append(steps, *cur, *task)
			cur = back
			return task.Attempts[0].End, nil
		},
	}
	pass, err := r.l.PassLesson(m, lesson, at, each)
	if errors.Is(err, ErrRun) {
		return nil, 0, err
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrRun, err)
	}
	r.pass(lesson, pass)
	cur.Pass = &pass
	return append(steps, *cur), pass.End, nil
}

// consolidate follows the engine's answer c at the end of a task link of a lesson, on the input in,
// when it is not further through the lesson: a task outside the lesson is a choice of its module and
// is given at once, with the pool read as at any choice; then the engine is called again, and the
// lesson of the same run must go on (R48). Any other answer is a run that cannot go on. It gives the
// step of the consolidation and the step of the lesson resumed.
func (r *Run) consolidate(lesson *pack.Lesson, in engine.Input, c engine.Choice) (*Step, *Step, error) {
	if c.Kind != engine.ChoiceTask || slices.ContainsFunc(lesson.Lesson.Tasks, func(t learnermodel.Task) bool { return t.ID == c.ID }) {
		return nil, nil, fmt.Errorf("%w: at the end of a link of %s at %d the engine chose %v %s (%v), neither the lesson nor a task outside it",
			ErrRun, lesson.Lesson.ID, in.Time, c.Kind, c.ID, c.Reason)
	}
	pool, err := r.pool(in.Time)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrRun, err)
	}
	task := &Step{Input: in, Choice: c, Pool: pool}
	end, err := r.give(task, in.Time)
	if err != nil {
		return nil, nil, err
	}
	back, b, err := r.call(end)
	if err != nil {
		return nil, nil, err
	}
	if !further(b, lesson) {
		return nil, nil, fmt.Errorf("%w: after the consolidation %s at %d the engine chose %v %s (%v), not the lesson %s",
			ErrRun, c.ID, end, b.Kind, b.ID, b.Reason, lesson.Lesson.ID)
	}
	return task, &Step{Input: back, Choice: b, Resumed: true}, nil
}

// call calls the engine within a lesson at a time, on the input api would build then.
func (r *Run) call(at int64) (engine.Input, engine.Choice, error) {
	in, err := r.input(at)
	if err != nil {
		return in, engine.Choice{}, err
	}
	c, err := r.choose(in)
	if err != nil {
		return in, c, fmt.Errorf("%w: the engine refuses the input at %d: %w", ErrRun, at, err)
	}
	return in, c, nil
}

// further tells whether an answer of the engine within a lesson is the lesson of the same run:
// further through it.
func further(c engine.Choice, lesson *pack.Lesson) bool {
	return c.Kind == engine.ChoiceLesson && c.ID == lesson.Lesson.ID && c.Reason == engine.ReasonOpenRun
}

// verdicts tells whether a lesson has skills and every one has a verdict: only then is the engine
// called within it (concepts 5 spec; R48).
func verdicts(m *pack.Module, lesson *pack.Lesson) bool {
	if len(lesson.Lesson.Skills) == 0 {
		return false
	}
	for _, id := range lesson.Lesson.Skills {
		if !verdict(m, id) {
			return false
		}
	}
	return true
}

// issue hands out skills at a time: the time of their last issue does not go back, so that a lesson
// going on after a consolidation leaves the consolidation's time (R48).
func (r *Run) issue(skills []string, at int64) {
	for _, s := range skills {
		is := r.issued[s]
		is.Skill, is.Issued, is.LastIssued = s, true, max(is.LastIssued, at)
		r.issued[s] = is
	}
}

// open records the run of a lesson as open while it is passed, as api holds it: an open run of the
// lesson goes on, and a new pass is a new run (concepts 5 spec).
func (r *Run) open(lesson *pack.Lesson, at int64) {
	if slices.ContainsFunc(r.runs, func(x engine.Run) bool { return x.Lesson == lesson.Lesson.ID && x.End == engine.RunOpen }) {
		return
	}
	r.runs = append(r.runs, engine.Run{Lesson: lesson.Lesson.ID, Skills: lesson.Lesson.Skills, End: engine.RunOpen})
	r.runStarts = append(r.runStarts, at)
}

// pass ends the open run of a lesson by its pass: passed, or open still when the pass was left.
func (r *Run) pass(lesson *pack.Lesson, p learner.Pass) {
	if p.Abandoned {
		return
	}
	for i := range r.runs {
		if r.runs[i].Lesson == lesson.Lesson.ID && r.runs[i].End == engine.RunOpen {
			r.runs[i].End, r.runs[i].EndedAt = engine.RunPassed, p.End
			return
		}
	}
}

// instance gives the number of the instance of a template an attempt takes, by its seed; -1 for a
// task that is no template (R12; R14).
func instance(m *pack.Module, t *pack.Task, a learner.Attempt) (int, error) {
	if !t.Candidate.Task.Template {
		return -1, nil
	}
	n, _, err := m.Instance(t.Candidate.Task.ID, a.Seed)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrRun, err)
	}
	return n, nil
}

// input builds the engine's input at a time as api would build it (architecture 3.6 engine).
func (r *Run) input(at int64) (engine.Input, error) {
	ro, err := r.readAt(at)
	if err != nil {
		return engine.Input{}, err
	}
	return engine.Input{
		Time: at, Readout: ro, Modules: r.suited, Runs: slices.Clone(r.runs), RubricDone: slices.Clone(r.rubricDone),
		Last: r.last, HasLast: r.hasLast, Summary: r.summary(),
	}, nil
}

// summary gives the summary of what was handed out, in the order of IDs (R104, R107 spec).
func (r *Run) summary() engine.Summary {
	var s engine.Summary
	for _, id := range slices.Sorted(maps.Keys(r.issued)) {
		s.Skills = append(s.Skills, r.issued[id])
	}
	for _, id := range slices.Sorted(maps.Keys(r.choices)) {
		s.Modules = append(s.Modules, engine.ModuleChoices{Module: id, Choices: r.choices[id]})
	}
	ids, ats := r.l.Solved()
	for i, id := range ids {
		s.Tasks = append(s.Tasks, engine.TaskSolution{Task: id, SolutionSentAt: ats[i]})
	}
	return s
}

// task finds a task of a choice in the learner's modules.
func (r *Run) task(id string) (*pack.Module, *pack.Task, error) {
	for _, m := range r.modules {
		if t, err := m.Task(id); err == nil {
			return m, t, nil
		}
	}
	return nil, nil, fmt.Errorf("%w: the engine chose the task %s, which no module of the learner has", ErrRun, id)
}

// lesson finds a lesson of a choice in the learner's modules.
func (r *Run) lesson(id string) (*pack.Module, *pack.Lesson, error) {
	for _, m := range r.modules {
		for i := range m.Lessons {
			if m.Lessons[i].Lesson.ID == id {
				return m, &m.Lessons[i], nil
			}
		}
	}
	return nil, nil, fmt.Errorf("%w: the engine chose the lesson %s, which no module of the learner has", ErrRun, id)
}

// close records a closed attempt as api would: the evidence of a task with a verdict, appended to
// the model's state; the attempt closed by "Done" of a task with a rubric; the solution sent to the
// device when the attempt closed right or with the solution shown (R82 spec); the summary of each
// primary skill, handed out at issued — the time of the choice of the task, or of its lesson, and
// never back (R48) —;
// and the last closed attempt. seen tells whether the task's solution went to the device before
// the attempt, as the learner keeps it.
func (r *Run) close(m *pack.Module, t *pack.Task, a learner.Attempt, code engine.ReasonCode, lesson string, issued int64, seen bool) (*learnermodel.Evidence, error) {
	task := t.Candidate.Task
	results := results(t, a)
	for _, res := range results {
		is := r.issued[res.Skill]
		is.Skill, is.Issued, is.LastIssued, is.LastCode, is.LastResult, is.HasLast = res.Skill, true, max(is.LastIssued, issued), code, res, true
		r.issued[res.Skill] = is
	}
	r.last = engine.Last{Task: task.ID, Left: a.Closure == learner.ClosedLeft, Lesson: lesson, InLesson: lesson != "", Results: results}
	r.hasLast = true
	if t.Rubric {
		if a.Closure == learner.ClosedDone {
			r.rubricDone = append(r.rubricDone, engine.RubricDone{Task: task.ID, At: a.End})
		}
		return nil, nil
	}
	e, err := r.evidence(m, t, a, seen)
	if err != nil {
		return nil, err
	}
	state, got, err := learnermodel.Append(r.state, e, r.graphs)
	if err != nil {
		return nil, fmt.Errorf("%w: the model refuses the evidence of %s: %w", ErrRun, task.ID, err)
	}
	r.state = state
	for _, res := range a.Items {
		r.rec.Answers += b2i(res.Place > 0)
	}
	return &got, nil
}

// evidence gives the evidence of an attempt on a task with a verdict (R82 spec; concepts 11 spec):
// its outcome by its wrong answers, the first and final verdict of each item, whether the task's
// solution went to the device before — no such fact for a template, whose instances differ —, and
// the time of its last action, and the technique taken on request (G8) with the leading question
// the same hint gave, when the task has one. The truth has no traps, and the learner asks for
// nothing else.
func (r *Run) evidence(m *pack.Module, t *pack.Task, a learner.Attempt, seen bool) (learnermodel.Evidence, error) {
	lm, err := r.l.Module(m.ID)
	if err != nil {
		return learnermodel.Evidence{}, fmt.Errorf("%w: %w", ErrRun, err)
	}
	task := t.Candidate.Task
	e := learnermodel.Evidence{
		Kind: learnermodel.EvidenceKind, Task: task, PracticeTime: a.End,
		StartBand:                learnermodel.BandRef{Place: lm.StartBand, ID: m.Graph.Bands[lm.StartBand-1].ID},
		SolutionSeenBefore:       seen && !task.Template,
		HasSolutionSeenFact:      !task.Template,
		TechniqueOnRequest:       a.Technique,
		LeadingQuestionOnRequest: a.LeadingQuestion,
	}
	switch {
	case a.Closure == learner.ClosedSolutionShown:
		e.Outcome = learnermodel.Outcome{Kind: learnermodel.OutcomeSolutionShown}
	case a.Rungs > 0:
		e.Outcome = learnermodel.Outcome{Kind: learnermodel.OutcomeAfterHint, Rung: int32(a.Rungs), HasRung: true}
	default:
		e.Outcome = learnermodel.Outcome{Kind: learnermodel.OutcomeFirstSubmission}
	}
	for _, res := range a.Items {
		first := learnermodel.VerdictNotYet
		if res.FirstRight {
			first = learnermodel.VerdictCorrect
		}
		final := learnermodel.VerdictNotYet
		if res.Right {
			final = learnermodel.VerdictCorrect
		}
		ia := learnermodel.ItemAnswer{
			Place: res.Place, Submitted: true,
			// An answer written again after it could not be read: no verdict and no mistake (R04 spec).
			Unverdicted:  unreadable(res.Unreadable),
			FirstVerdict: first, HasFirstVerdict: true,
			FinalVerdict: final, HasFinalVerdict: true,
		}
		if res.Right {
			ia.WrongBeforeFirstCorrect, ia.HasFirstCorrect = int32(res.Wrong), true
		}
		e.Items = append(e.Items, ia)
	}
	return e, nil
}

// results gives the attempt's result on each primary skill of its task, in the order of IDs: of
// the items the skill is primary in, how many got a first verdict, how many of them right, and
// whether the solution was shown (architecture 3.6 engine).
func results(t *pack.Task, a learner.Attempt) []engine.SkillResult {
	by := map[string]*engine.SkillResult{}
	for _, s := range t.Candidate.Task.Primary {
		by[s] = &engine.SkillResult{Skill: s}
	}
	for _, res := range a.Items {
		for _, s := range t.ItemPrimary(int(res.Place - 1)) {
			sr := by[s]
			if sr == nil {
				sr = &engine.SkillResult{Skill: s}
				by[s] = sr
			}
			sr.Answered++
			sr.FirstRight += int32(b2i(res.FirstRight))
		}
	}
	var out []engine.SkillResult
	for _, s := range slices.Sorted(maps.Keys(by)) {
		sr := by[s]
		sr.SolutionShown = a.Closure == learner.ClosedSolutionShown
		out = append(out, *sr)
	}
	return out
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// unreadable gives the reasons of the submissions of an item before its first verdict: one that
// could not be read, when the learner wrote it again (S-unreadable).
func unreadable(yes bool) []learnermodel.UnverdictedReason {
	if yes {
		return []learnermodel.UnverdictedReason{learnermodel.UnverdictedUnreadableOther}
	}
	return nil
}
