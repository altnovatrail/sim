package run

import (
	"bytes"
	"errors"
	"slices"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/spec/types/formatv0"
)

const (
	olympiad = "novatrail/olympiad-math"
	grade3   = "novatrail/grade-3-math"
	essays   = "novatrail/english-essays"
)

func modules(t *testing.T, ids ...string) []*pack.Module {
	t.Helper()
	ms, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	var out []*pack.Module
	for _, id := range ids {
		out = append(out, ms[id])
	}
	return out
}

func runOf(t *testing.T, population string, n int, limit Limit, mods []*pack.Module) *Record {
	t.Helper()
	l, err := learner.New(seed.Working, population, n, mods, learner.Base())
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(l, mods, limit)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := r.All()
	if err != nil {
		t.Fatalf("%s: %v", l.Unit, err)
	}
	return rec
}

// The base population runs on each published package to the end of its limit: 200 answers of
// items with a verdict in each of the maths, and in essays until the tasks with a rubric run out
// or 60 days pass; and a learner of both maths runs on both.
func TestRunPublished(t *testing.T) {
	for n := range 5 {
		for _, id := range []string{olympiad, grade3} {
			rec := runOf(t, "G0", n, Maths, modules(t, id))
			if rec.Answers < 200 || rec.Days > 365 || len(rec.Steps) == 0 {
				t.Errorf("%s in %s: %d answers in %d days", rec.Unit, id, rec.Answers, rec.Days)
			}
		}
		rec := runOf(t, "E-base", n, Essays, modules(t, essays))
		done := map[string]bool{}
		for _, s := range rec.Steps {
			for _, a := range s.Attempts {
				if a.Closure == learner.ClosedDone {
					done[a.Task] = true
				}
			}
		}
		pool, _ := modules(t, essays)[0].ForAge(rec.Learner.Age)
		rubrics := 0
		for _, c := range pool {
			if t2, _ := modules(t, essays)[0].Task(c.Task.ID); t2.Rubric {
				rubrics++
			}
		}
		if len(done) != rubrics && rec.Days <= 60 {
			t.Errorf("%s: %d of %d tasks with a rubric done in %d days", rec.Unit, len(done), rubrics, rec.Days)
		}
		both := runOf(t, "A-linked", n, Maths, modules(t, olympiad, grade3))
		if both.Answers < 200 {
			t.Errorf("%s in both maths: %d answers", both.Unit, both.Answers)
		}
	}
}

// A run repeated on the same seeds gives the same bits: every choice, and the model's state at the
// end (docs/architecture.md, section 3.5, I1).
func TestRunRepeats(t *testing.T) {
	for _, id := range []string{olympiad, grade3} {
		a := runOf(t, "G0", 11, Maths, modules(t, id))
		b := runOf(t, "G0", 11, Maths, modules(t, id))
		if len(a.Steps) != len(b.Steps) {
			t.Fatalf("%s: %d steps, then %d", id, len(a.Steps), len(b.Steps))
		}
		for k := range a.Steps {
			if a.Steps[k].Choice != b.Steps[k].Choice {
				t.Fatalf("%s: step %d chose %v, then %v", id, k, a.Steps[k].Choice, b.Steps[k].Choice)
			}
		}
		sa, err := learnermodel.EncodeState(a.State)
		if err != nil {
			t.Fatal(err)
		}
		sb, _ := learnermodel.EncodeState(b.State)
		if !bytes.Equal(sa, sb) {
			t.Fatalf("%s: two states at the end", id)
		}
	}
}

// What the bench gives the engine and the model is what api would give: the pool and the lessons
// the age category suits, the summary of what was handed out before, the last closed attempt, and
// evidence whose outcome is the attempt's.
func TestRunInput(t *testing.T) {
	rec := runOf(t, "G0", 2, Maths, modules(t, olympiad))
	choices := int64(0)
	solutions := map[string]int64{}
	for k, s := range rec.Steps {
		in := s.Input
		if k > 0 {
			prev := rec.Steps[k-1]
			if !in.HasLast || len(prev.Attempts) == 0 || in.Last.Task != prev.Attempts[len(prev.Attempts)-1].Task {
				t.Fatalf("step %d: the last attempt %+v", k, in.Last)
			}
		}
		if got := moduleChoices(in.Summary); got != choices {
			t.Fatalf("step %d: the summary counts %d choices, %d were made", k, got, choices)
		}
		for _, ts := range in.Summary.Tasks {
			if solutions[ts.Task] != ts.SolutionSentAt {
				t.Fatalf("step %d: the solution of %s went at %d, the summary says %d", k, ts.Task, solutions[ts.Task], ts.SolutionSentAt)
			}
		}
		if s.Choice.Kind != engine.ChoiceNothing {
			choices++
		}
		if len(s.Evidence) != len(s.Attempts) || len(s.Instances) != len(s.Attempts) {
			t.Fatalf("step %d: %d attempts, %d evidence, %d instances", k, len(s.Attempts), len(s.Evidence), len(s.Instances))
		}
		for i, a := range s.Attempts {
			e := s.Evidence[i]
			if e == nil || s.Instances[i] != -1 {
				t.Fatalf("step %d: evidence %v, instance %d", k, e, s.Instances[i])
			}
			want := learnermodel.OutcomeFirstSubmission
			switch {
			case a.Closure == learner.ClosedSolutionShown:
				want = learnermodel.OutcomeSolutionShown
			case a.Rungs > 0:
				want = learnermodel.OutcomeAfterHint
			}
			if e.Outcome.Kind != want || e.PracticeTime != a.End || len(e.Items) != len(a.Items) {
				t.Fatalf("step %d: evidence %+v of attempt %+v", k, e.Outcome, a)
			}
			if _, seen := solutions[a.Task]; e.SolutionSeenBefore != seen {
				t.Fatalf("step %d: %s seen before %v, the run says %v", k, a.Task, e.SolutionSeenBefore, seen)
			}
			solutions[a.Task] = a.End
		}
	}
	e := runOf(t, "E-base", 0, Essays, modules(t, essays))
	pool, lessons := modules(t, essays)[0].ForAge(e.Learner.Age)
	if in := e.Steps[0].Input; len(in.Modules) != 1 || len(in.Modules[0].Pool) != len(pool) || len(in.Modules[0].Lessons) != len(lessons) {
		t.Errorf("essays: the engine gets %d tasks and %d lessons, api would give %d and %d",
			len(in.Modules[0].Pool), len(in.Modules[0].Lessons), len(pool), len(lessons))
	}
	for _, s := range e.Steps {
		if s.Choice.Kind == engine.ChoiceLesson && (s.Pass == nil || s.Pass.Lesson != s.Choice.ID) {
			t.Fatalf("a lesson chosen and %+v passed", s.Pass)
		}
		for i, a := range s.Attempts {
			if (a.Closure == learner.ClosedDone) != (s.Evidence[i] == nil) {
				t.Fatalf("%s closed %v with evidence %v", a.Task, a.Closure, s.Evidence[i])
			}
		}
	}
	last := e.Steps[len(e.Steps)-1]
	if in := last.Input; len(in.Runs) == 0 || in.Runs[0].End != engine.RunPassed || len(in.RubricDone) == 0 && rubricsIn(e) > 0 {
		t.Errorf("essays: runs %v, done %v", in.Runs, in.RubricDone)
	}
}

func moduleChoices(s engine.Summary) int64 {
	var n int64
	for _, m := range s.Modules {
		n += m.Choices
	}
	return n
}

func rubricsIn(rec *Record) int {
	n := 0
	for _, s := range rec.Steps[:len(rec.Steps)-1] {
		for _, a := range s.Attempts {
			if a.Closure == learner.ClosedDone {
				n++
			}
		}
	}
	return n
}

// When the engine has nothing to give, the session ends, and the run goes on the next day with a
// session until its limit: here a module whose pool and lessons are empty.
func TestRunNothing(t *testing.T) {
	m := *modules(t, olympiad)[0]
	m.Tasks, m.Lessons = nil, nil
	rec := runOf(t, "G0", 1, Limit{Answers: 200, Days: 10}, []*pack.Module{&m})
	if rec.Answers != 0 || rec.Days < 1 || rec.Days > 10 || len(rec.Steps) == 0 {
		t.Fatalf("%d answers, %d days, %d steps", rec.Answers, rec.Days, len(rec.Steps))
	}
	days := map[int64]bool{}
	for _, s := range rec.Steps {
		if s.Choice.Kind != engine.ChoiceNothing {
			t.Fatalf("the engine chose %v from nothing", s.Choice)
		}
		days[s.Input.Time] = true
	}
	if len(days) != len(rec.Steps) {
		t.Error("two choices of nothing in one session")
	}
}

func TestNewRefuses(t *testing.T) {
	mods := modules(t, olympiad)
	l, err := learner.New(seed.Working, "G0", 0, mods, learner.Base())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := New(l, nil, Maths); !errors.Is(err, ErrRun) {
		t.Errorf("no module: %v", err)
	}
	if _, err := New(l, mods, Limit{Answers: 200}); !errors.Is(err, ErrRun) {
		t.Errorf("a limit with no days: %v", err)
	}
	if _, err := New(l, modules(t, grade3), Maths); !errors.Is(err, ErrRun) {
		t.Errorf("a module the learner does not learn: %v", err)
	}
	r, err := New(l, mods, Maths)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := r.task("no/task/x"); !errors.Is(err, ErrRun) {
		t.Errorf("a task no module has: %v", err)
	}
	if _, _, err := r.lesson("no/lesson/x"); !errors.Is(err, ErrRun) {
		t.Errorf("a lesson no module has: %v", err)
	}
}

// The result of an attempt on each primary skill, in facts, as the engine reads a failure.
func TestResults(t *testing.T) {
	task := &pack.Task{Candidate: engine.Candidate{Task: learnermodel.Task{
		ID: "m/task/t", Primary: []string{"m/skill/a"},
		Items: []learnermodel.TaskItem{{Place: 1}, {Place: 2, OwnSkills: true, Primary: []string{"m/skill/b"}}},
	}}}
	a := learner.Attempt{Items: []learner.ItemResult{{Place: 1, FirstRight: true}, {Place: 2}}, Closure: learner.ClosedSolutionShown}
	got := results(task, a)
	want := []engine.SkillResult{
		{Skill: "m/skill/a", Answered: 1, FirstRight: 1, SolutionShown: true},
		{Skill: "m/skill/b", Answered: 1, FirstRight: 0, SolutionShown: true},
	}
	if !slices.Equal(got, want) {
		t.Errorf("results %v, want %v", got, want)
	}
}

// An input the engine refuses — here a pool with a task twice — stops the run with the reason.
func TestRunRefused(t *testing.T) {
	m := *modules(t, olympiad)[0]
	m.Tasks = append(slices.Clone(m.Tasks), m.Tasks[0])
	l, err := learner.New(seed.Working, "G0", 0, []*pack.Module{&m}, learner.Base())
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(l, []*pack.Module{&m}, Maths)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.All(); !errors.Is(err, ErrRun) {
		t.Errorf("a pool with a task twice: %v", err)
	}
}

// A template is given as one task, and each attempt takes the instance its seed names, as the
// platform chooses it (R12; R14): on a synthetic module of one skill and one template.
func TestRunTemplate(t *testing.T) {
	cv := &formatv0.CanonicalView{
		Manifest: formatv0.Manifest{
			ModuleFormatVersion: "0.3", ID: "example/template", Version: "0.1.0",
			Audience: formatv0.Audience{AgeCategories: []formatv0.AgeCategory{"under_13"}},
			Bands:    []formatv0.Band{{ID: "only", Grades: []formatv0.Grade{3}}},
		},
		SkillGraph: &formatv0.SkillGraph{Skills: map[formatv0.SkillID]formatv0.Skill{"s": {Assessability: formatv0.SkillAssessabilityVerdict}}},
		SkillTiers: formatv0.SkillTiers{"s": 1},
		Tasks: map[formatv0.TaskID]formatv0.Task{"tpl": {
			ID: "tpl", Canvas: formatv0.Canvas{Kind: formatv0.CanvasKindHandwritingSheet}, Band: "only", Difficulty: 3,
			Skills:   formatv0.Skills{Primary: []formatv0.SkillID{"s"}},
			Template: &formatv0.Template{Parameters: map[formatv0.TemplateName]formatv0.Domain{"n": {From: 1, To: 4}}},
			Items:    []formatv0.Item{{AnswerSlots: []formatv0.AnswerSlot{{Check: formatv0.Check{Kind: formatv0.CheckKindNumber}}}}},
		}},
		Instances: map[formatv0.TaskID][]formatv0.Numbers{"tpl": {{"n": 1}, {"n": 2}, {"n": 3}, {"n": 4}}},
	}
	m, err := pack.Translate(cv)
	if err != nil {
		t.Fatal(err)
	}
	rec := runOf(t, "G0", 0, Limit{Answers: 20, Days: 30}, []*pack.Module{m})
	if len(rec.Steps) == 0 {
		t.Fatal("no step")
	}
	for k, s := range rec.Steps {
		for i, a := range s.Attempts {
			want, err := seed.Instance(a.Seed, "tpl", 4)
			if err != nil || s.Instances[i] != want {
				t.Fatalf("step %d: instance %d, the seed names %d", k, s.Instances[i], want)
			}
			if s.Evidence[i].HasSolutionSeenFact {
				t.Fatalf("step %d: a template's evidence has the fact of a solution seen", k)
			}
		}
	}
}
