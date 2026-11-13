package run

import (
	"errors"
	"slices"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// consolidation is a run on the synthetic sheet whose rule is laid by a test: the first choice is
// the lesson of a skill, whose tasks d1 and d3 are its links, and a task of the same skill outside
// the lesson is there to consolidate (R19; R48).
type consolidation struct {
	r      *Run
	m      *pack.Module
	lesson *pack.Lesson
	skill  string
	inside []string // the lesson's tasks, in its order
	out    string   // a task on the skill outside the lesson
}

func newConsolidation(t *testing.T) *consolidation {
	t.Helper()
	ms, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	m := ms[synthetic.Sheet]
	mods := []*pack.Module{m}
	l, err := learner.New(seed.Working, "G0", 0, mods, learner.Base())
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(l, mods, Maths)
	if err != nil {
		t.Fatal(err)
	}
	c := &consolidation{r: r, m: m}
	for i := range m.Lessons {
		if len(m.Lessons[i].Lesson.Tasks) == 2 {
			c.lesson = &m.Lessons[i]
			break
		}
	}
	if c.lesson == nil {
		t.Fatal("the sheet has no lesson of two tasks")
	}
	c.skill = c.lesson.Lesson.Skills[0]
	for _, lt := range c.lesson.Lesson.Tasks {
		c.inside = append(c.inside, lt.ID)
	}
	pool, _ := m.ForAge(l.Age)
	for _, cand := range pool {
		if slices.Contains(cand.Task.Primary, c.skill) && !slices.Contains(c.inside, cand.Task.ID) {
			c.out = cand.Task.ID
			break
		}
	}
	if c.out == "" {
		t.Fatalf("no task of %s outside its lesson", c.skill)
	}
	return c
}

// rule chooses the lesson first, then answers at the end of the link of the task after with the
// choice at, and otherwise as the engine does.
func (c *consolidation) rule(after string, at engine.Choice) *consolidation {
	c.r.choose = func(in engine.Input) (engine.Choice, error) {
		switch {
		case !in.HasLast:
			return engine.Choice{Kind: engine.ChoiceLesson, ID: c.lesson.Lesson.ID, Reason: engine.ReasonNewSkill}, nil
		case in.Last.InLesson && in.Last.Task == after:
			return at, nil
		}
		return engine.Choose(in)
	}
	return c
}

func (c *consolidation) task() engine.Choice {
	return engine.Choice{Kind: engine.ChoiceTask, ID: c.out, Reason: engine.ReasonConsolidation}
}

// A consolidation after a link is a choice of the module and an issue of the skill, given at once:
// the lesson before it, the consolidation, and the lesson resumed after it are steps in the order of
// time, and the lesson ends passed, one place of the session.
func TestConsolidation(t *testing.T) {
	for _, link := range []int{0, 1} {
		c := newConsolidation(t)
		c.rule(c.inside[link], c.task())
		if err := c.r.Step(); err != nil {
			t.Fatal(err)
		}
		steps := c.r.rec.Steps
		if len(steps) != 3 {
			t.Fatalf("link %d: %d steps, want the lesson, the consolidation and the lesson resumed", link, len(steps))
		}
		lesson, task, resumed := steps[0], steps[1], steps[2]
		if lesson.Choice.Kind != engine.ChoiceLesson || lesson.Resumed || lesson.Pass != nil || len(lesson.Attempts) != link+1 || len(lesson.Further) != link {
			t.Errorf("link %d: the lesson before %+v, %d attempts, %d further", link, lesson.Choice, len(lesson.Attempts), len(lesson.Further))
		}
		in := task.Input
		if task.Choice != c.task() || task.Pool == nil || len(task.Attempts) != 1 || task.Attempts[0].Task != c.out || len(task.Obs) != 1 ||
			!in.Last.InLesson || in.Last.Task != c.inside[link] || in.Time != lesson.Attempts[link].End {
			t.Errorf("link %d: the consolidation %+v on the input at %d after %+v", link, task.Choice, in.Time, in.Last)
		}
		if !slices.ContainsFunc(in.Runs, func(x engine.Run) bool { return x.Lesson == c.lesson.Lesson.ID && x.End == engine.RunOpen }) {
			t.Errorf("link %d: the run is not open at the end of the link: %v", link, in.Runs)
		}
		end := task.Attempts[0].End
		if !resumed.Resumed || !further(resumed.Choice, c.lesson) || resumed.Pool != nil || resumed.Input.Time != end ||
			resumed.Pass == nil || resumed.Pass.Abandoned || len(resumed.Attempts) != 1-link || len(resumed.Further) != 1-link {
			t.Errorf("link %d: the lesson resumed %+v at %d", link, resumed.Choice, resumed.Input.Time)
		}
		if link == 0 && resumed.Attempts[0].Start != end {
			t.Errorf("the next link given at %d, the consolidation closed at %d", resumed.Attempts[0].Start, end)
		}
		if got := c.r.choices[c.m.ID]; got != 2 {
			t.Errorf("link %d: %d choices of the module, want the lesson and the consolidation", link, got)
		}
		is := c.r.issued[c.skill]
		if is.LastIssued != in.Time {
			t.Errorf("link %d: the skill last issued at %d, the consolidation at %d", link, is.LastIssued, in.Time)
		}
		wantCode := lesson.Choice.Reason
		if link == 1 {
			wantCode = engine.ReasonConsolidation
		}
		if is.LastCode != wantCode {
			t.Errorf("link %d: the skill's last code %v, want %v", link, is.LastCode, wantCode)
		}
		last := resumed.Pass.End
		if link == 1 && last != end {
			t.Errorf("the lesson ends at %d, the consolidation after its last link at %d", last, end)
		}
		if run := c.r.runs[0]; len(c.r.runs) != 1 || run.End != engine.RunPassed || run.EndedAt != last {
			t.Errorf("link %d: the runs %+v, the pass ended at %d", link, c.r.runs, last)
		}
		if next := c.r.cal.Next(); next != last {
			t.Errorf("link %d: the next choice at %d, the lesson ended at %d: the lesson is one place", link, next, last)
		}
	}
}

// Further through the lesson moves nothing of the summary: each call at the end of a link has the
// lesson's choice counted once and its skills issued at its time; the lesson is one step.
func TestFurther(t *testing.T) {
	mods := modules(t, olympiad)
	rec := runOf(t, "G0", 1, Maths, mods)
	lessons := 0
	for k, s := range rec.Steps {
		if s.Choice.Kind != engine.ChoiceLesson {
			if len(s.Further) > 0 || s.Resumed {
				t.Fatalf("step %d: %v went further", k, s.Choice)
			}
			continue
		}
		lessons++
		if s.Resumed || len(s.Further) != len(s.Attempts) {
			t.Fatalf("step %d: %d calls further on %d links", k, len(s.Further), len(s.Attempts))
		}
		for i, f := range s.Further {
			in := f.Input
			skill := s.Skills()[0]
			is := in.Summary.Skills[slices.IndexFunc(in.Summary.Skills, func(x engine.SkillIssue) bool { return x.Skill == skill })]
			if !further(f.Choice, &pack.Lesson{Lesson: engine.Lesson{ID: s.Choice.ID}}) || in.Time != s.Attempts[i].End ||
				moduleChoices(in.Summary) != moduleChoices(s.Input.Summary)+1 || is.LastIssued != s.Input.Time || in.Last.Task != s.Attempts[i].Task {
				t.Fatalf("step %d, link %d: %+v at %d, %d choices", k, i, f.Choice, in.Time, moduleChoices(in.Summary))
			}
		}
	}
	if lessons == 0 {
		t.Fatal("no lesson chosen")
	}
	// No call within a lesson of skills with no verdict, and none under the ceiling.
	for _, rec := range []*Record{runOf(t, "E-base", 0, Essays, modules(t, essays)), ceilingOf(t, mods)} {
		for k, s := range rec.Steps {
			if len(s.Further) > 0 || s.Resumed {
				t.Fatalf("%s, step %d: a call within the lesson %s", rec.Unit, k, s.Choice.ID)
			}
		}
	}
}

func ceilingOf(t *testing.T, mods []*pack.Module) *Record {
	t.Helper()
	l, err := learner.New(seed.Working, "G0", 1, mods, learner.Base())
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(l, mods, Maths)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := r.Use(Ceiling).All()
	if err != nil {
		t.Fatal(err)
	}
	return rec
}

// At the end of a link the engine gives the lesson of the run or a task outside the lesson, and
// after a consolidation the lesson again; any other answer, or an input it refuses, stops the run.
func TestConsolidationRefused(t *testing.T) {
	c := newConsolidation(t)
	refuse := errors.New("refused")
	for name, at := range map[string]engine.Choice{
		"nothing":        {Kind: engine.ChoiceNothing},
		"another lesson": {Kind: engine.ChoiceLesson, ID: c.m.Graph.Module + "/lesson/none", Reason: engine.ReasonOpenRun},
		"its own lesson": {Kind: engine.ChoiceLesson, ID: c.lesson.Lesson.ID, Reason: engine.ReasonNewSkill},
		"a task inside":  {Kind: engine.ChoiceTask, ID: c.inside[1], Reason: engine.ReasonConsolidation},
		"no task":        {Kind: engine.ChoiceTask, ID: c.m.Graph.Module + "/task/none", Reason: engine.ReasonConsolidation},
	} {
		c := newConsolidation(t)
		c.rule(c.inside[0], at)
		if err := c.r.Step(); !errors.Is(err, ErrRun) {
			t.Errorf("%s at the end of a link: %v", name, err)
		}
	}
	// The engine refuses the input at the end of a link.
	c.rule(c.inside[0], engine.Choice{})
	choose := c.r.choose
	c.r.choose = func(in engine.Input) (engine.Choice, error) {
		if in.HasLast {
			return engine.Choice{}, refuse
		}
		return choose(in)
	}
	if err := c.r.Step(); !errors.Is(err, ErrRun) || !errors.Is(err, refuse) {
		t.Errorf("an input refused at the end of a link: %v", err)
	}
	// After the consolidation, another task, or an input refused.
	for _, back := range []error{nil, refuse} {
		c := newConsolidation(t).rule("", engine.Choice{})
		c.r.choose = func(in engine.Input) (engine.Choice, error) {
			switch {
			case !in.HasLast:
				return engine.Choice{Kind: engine.ChoiceLesson, ID: c.lesson.Lesson.ID, Reason: engine.ReasonNewSkill}, nil
			case in.Last.InLesson:
				return c.task(), nil
			case back != nil:
				return engine.Choice{}, back
			}
			return c.task(), nil
		}
		if err := c.r.Step(); !errors.Is(err, ErrRun) {
			t.Errorf("after the consolidation, %v: %v", back, err)
		}
	}
}

// The steps of a lesson with a consolidation read as a choice of a task between two parts of the
// lesson: the skills of the consolidation are its task's.
func TestConsolidationSkills(t *testing.T) {
	c := newConsolidation(t)
	c.rule(c.inside[0], c.task())
	if err := c.r.Step(); err != nil {
		t.Fatal(err)
	}
	for k, want := range [][]string{{c.skill}, {c.skill}, {c.skill}} {
		if got := c.r.rec.Steps[k].Skills(); !slices.Equal(got, want) {
			t.Errorf("step %d on %v, want %v", k, got, want)
		}
	}
	if e := c.r.rec.Steps[1].Evidence[0]; e == nil || e.Task.ID != c.out || e.Kind != learnermodel.EvidenceKind {
		t.Errorf("the consolidation's evidence %+v", e)
	}
}

// What the run refuses around a consolidation, laid by the rule in the run itself, as
// TestObserveRefuses does: a lesson no module has; a pool it cannot read at the end of the link; a
// consolidation the model cannot read before it, or after it.
func TestConsolidationFaults(t *testing.T) {
	ms, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	other := ms[synthetic.Workbook] // a module the learner does not learn and the model has no graph of
	for name, lay := range map[string]func(c *consolidation) engine.Choice{
		"the pool": func(c *consolidation) engine.Choice {
			pool, _ := other.ForAge(c.r.l.Age)
			c.r.modules, c.r.suited = append(c.r.modules, other), append(c.r.suited, engine.Module{Graph: other.Graph, Pool: pool})
			return c.task()
		},
		"before": func(c *consolidation) engine.Choice {
			c.r.modules, c.r.suited = append(c.r.modules, other), append(c.r.suited, engine.Module{Graph: other.Graph})
			return engine.Choice{Kind: engine.ChoiceTask, ID: other.Tasks[0].Candidate.Task.ID, Reason: engine.ReasonConsolidation}
		},
		"after": func(c *consolidation) engine.Choice {
			c.r.modules, c.r.suited = append(c.r.modules, other), append(c.r.suited, engine.Module{Graph: other.Graph})
			return c.task()
		},
	} {
		c := newConsolidation(t)
		c.r.choose = func(in engine.Input) (engine.Choice, error) {
			switch {
			case !in.HasLast:
				return engine.Choice{Kind: engine.ChoiceLesson, ID: c.lesson.Lesson.ID, Reason: engine.ReasonNewSkill}, nil
			case in.Last.InLesson:
				return lay(c), nil
			}
			return engine.Choose(in)
		}
		if err := c.r.Step(); !errors.Is(err, ErrRun) {
			t.Errorf("%s: %v", name, err)
		}
	}
	c := newConsolidation(t)
	c.r.choose = func(engine.Input) (engine.Choice, error) {
		return engine.Choice{Kind: engine.ChoiceLesson, ID: c.m.Graph.Module + "/lesson/none", Reason: engine.ReasonNewSkill}, nil
	}
	if err := c.r.Step(); !errors.Is(err, ErrRun) {
		t.Errorf("a lesson no module has: %v", err)
	}
}

// A lesson with no skills has none with a verdict: the engine is not called within it.
func TestVerdictsOfNoSkills(t *testing.T) {
	c := newConsolidation(t)
	if verdicts(c.m, &pack.Lesson{}) || !verdicts(c.m, c.lesson) {
		t.Error("a lesson with no skills is one of skills with a verdict, or the sheet's is not")
	}
}
