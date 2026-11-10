package invariant

import (
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
)

const (
	olympiad = "novatrail/olympiad-math"
	grade3   = "novatrail/grade-3-math"
	essays   = "novatrail/english-essays"
)

func record(t *testing.T, population string, n int, limit run.Limit, ids ...string) *run.Record {
	t.Helper()
	ms, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	var mods []*pack.Module
	for _, id := range ids {
		mods = append(mods, ms[id])
	}
	l, err := learner.New(seed.Working, population, n, mods, learner.Base())
	if err != nil {
		t.Fatal(err)
	}
	r, err := run.New(l, mods, limit)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := r.All()
	if err != nil {
		t.Fatal(err)
	}
	return rec
}

// The invariants hold on every learner of the base population on the published packages: G0 in
// each of the maths and in both, a learner of essays (docs/architecture.md, section 3.5).
func TestInvariantsHold(t *testing.T) {
	cases := []struct {
		population string
		limit      run.Limit
		ids        []string
	}{
		{"G0", run.Maths, []string{olympiad}},
		{"G0", run.Maths, []string{grade3}},
		{"A-linked", run.Maths, []string{olympiad, grade3}},
		{"E-base", run.Essays, []string{essays}},
	}
	for _, c := range cases {
		// Olympiad maths 0.1.1 is 772 tasks, and a learner of it takes a hundred times longer: fewer of
		// them keep the suite within the budget of the builds same-bits runs it on.
		learners := 25
		if slices.Contains(c.ids, olympiad) {
			learners = 6
		}
		for n := range learners {
			rec := record(t, c.population, n, c.limit, c.ids...)
			for _, b := range Check(rec) {
				t.Error(b)
			}
			if n < 3 {
				for _, b := range Same(rec, record(t, c.population, n, c.limit, c.ids...)) {
					t.Error(b)
				}
			}
		}
	}
}

// first gives the first step of a record that matches, and its place.
func first(t *testing.T, rec *run.Record, match func(run.Step) bool) (int, run.Step) {
	t.Helper()
	for k, s := range rec.Steps {
		if match(s) {
			return k, s
		}
	}
	t.Fatalf("%s has no such step", rec.Unit)
	return 0, run.Step{}
}

func only(t *testing.T, bs []Breach, inv string) {
	t.Helper()
	if len(bs) == 0 {
		t.Fatalf("no breach of %s", inv)
	}
	for _, b := range bs {
		if b.Invariant != inv {
			t.Fatalf("a breach of %s where one of %s was made: %v", b.Invariant, inv, b)
		}
		if !strings.Contains(b.String(), inv) {
			t.Fatalf("the breach does not name its invariant: %v", b)
		}
	}
}

// Each invariant catches a breach laid in a record: a choice the input does not give again, a code
// against its meaning or outside the list, the author's order broken, a solved task at once.
func TestBreaches(t *testing.T) {
	// The first learner of G0 who gets a task by rotation: on the 19 skills of olympiad maths 0.1.1,
	// not every one does.
	rotation := func(s run.Step) bool {
		return s.Choice.Kind == engine.ChoiceTask && s.Choice.Reason == engine.ReasonRotation
	}
	var rec *run.Record
	for n := range 10 {
		if rec = record(t, "G0", n, run.Maths, olympiad); slices.ContainsFunc(rec.Steps, rotation) {
			break
		}
	}
	k, s := first(t, rec, rotation)

	t.Run("I1 another choice", func(t *testing.T) {
		bad := *rec
		bad.Steps = []run.Step{s}
		bad.Steps[0].Choice.ID = "novatrail/olympiad-math/task/enum-12-d1-3"
		var i1 []Breach
		for _, b := range Check(&bad) {
			if b.Invariant == "I1" {
				i1 = append(i1, b)
			}
		}
		only(t, i1, "I1")
		// CheckWithoutI1 takes no choice again: past the first learners of a cell, a repeated run
		// holds I1 by Same, and the gate repeats every candidate.
		if hasBreach(CheckWithoutI1(&bad), "I1") || len(CheckWithoutI1(rec)) != 0 {
			t.Errorf("CheckWithoutI1 gives %v", CheckWithoutI1(&bad))
		}
	})
	t.Run("I1 an input refused", func(t *testing.T) {
		bad := run.Record{Unit: rec.Unit, Steps: []run.Step{s}}
		bad.Steps[0].Input.Modules = append(bad.Steps[0].Input.Modules, bad.Steps[0].Input.Modules[0])
		var i1 []Breach
		for _, b := range Check(&bad) {
			if b.Invariant == "I1" {
				i1 = append(i1, b)
			}
		}
		only(t, i1, "I1")
	})
	t.Run("I1 another run", func(t *testing.T) {
		other := record(t, "G0", 2, run.Maths, olympiad)
		other.Unit = rec.Unit
		if len(Same(rec, other)) == 0 {
			t.Fatal("two learners' runs are the same")
		}
		cut := *rec
		cut.Steps = rec.Steps[:len(rec.Steps)-1]
		if len(Same(rec, &cut)) == 0 {
			t.Fatal("a run one step shorter is the same")
		}
		ev := *rec
		ev.Steps = append([]run.Step(nil), rec.Steps...)
		e := *ev.Steps[k].Evidence[0]
		e.PracticeTime++
		ev.Steps[k].Evidence = []*learnermodel.Evidence{&e}
		if len(Same(rec, &ev)) == 0 {
			t.Fatal("a run with other evidence is the same")
		}
		st := *rec
		st.State = learnermodel.State{}
		if len(Same(rec, &st)) == 0 {
			t.Fatal("a run with another state at the end is the same")
		}
	})

	codes := []struct {
		name   string
		reason engine.ReasonCode
		kind   engine.ChoiceKind
		id     string
	}{
		{"a code outside the list", 10, engine.ChoiceTask, s.Choice.ID},
		{"a choice of nothing with a code", engine.ReasonQueue, engine.ChoiceNothing, ""},
		{"a choice of no kind", engine.ReasonQueue, engine.ChoiceUnspecified, s.Choice.ID},
		{"new skill on a skill given", engine.ReasonNewSkill, engine.ChoiceTask, s.Choice.ID},
		{"consolidation with no failure", engine.ReasonConsolidation, engine.ChoiceTask, s.Choice.ID},
		{"author order on a skill with a verdict", engine.ReasonAuthorOrder, engine.ChoiceTask, s.Choice.ID},
		{"open run with no run", engine.ReasonOpenRun, engine.ChoiceTask, s.Choice.ID},
		{"skill chosen with no map", engine.ReasonSkillChosen, engine.ChoiceTask, s.Choice.ID},
		{"repetition before the review date", engine.ReasonRepetition, engine.ChoiceTask, s.Choice.ID},
		{"a task no pool has", engine.ReasonQueue, engine.ChoiceTask, "novatrail/olympiad-math/task/none"},
		{"probe on a lesson", engine.ReasonProbe, engine.ChoiceLesson, s.Input.Modules[0].Lessons[0].ID},
		{"consolidation on a lesson", engine.ReasonConsolidation, engine.ChoiceLesson, s.Input.Modules[0].Lessons[0].ID},
	}
	for _, c := range codes {
		t.Run("I2 "+c.name, func(t *testing.T) {
			bad := s
			bad.Choice = engine.Choice{Kind: c.kind, ID: c.id, Reason: c.reason}
			if got := reasonCode(bad); len(got) == 0 {
				t.Fatalf("no breach of I2 on %+v", bad.Choice)
			}
		})
	}
	t.Run("I2 holds on the step as it was", func(t *testing.T) {
		if got := reasonCode(s); len(got) != 0 {
			t.Fatal(got)
		}
	})

	t.Run("I3", func(t *testing.T) {
		e := record(t, "E-base", 0, run.Essays, essays)
		_, step := first(t, e, func(s run.Step) bool {
			if s.Choice.Reason != engine.ReasonAuthorOrder {
				return false
			}
			for _, id := range s.Skills() {
				v, _ := viewOf(s.Input, id)
				for _, p := range v.skill.Prerequisites {
					if pv, ok := viewOf(s.Input, p); ok && pv.skill.Assessability != learnermodel.AssessabilityVerdict {
						return true
					}
				}
			}
			return false
		})
		if got := authorOrder(step); len(got) != 0 {
			t.Fatalf("the author's order breached as it was: %v", got)
		}
		step.Input.Runs, step.Input.RubricDone = nil, nil
		if got := authorOrder(step); len(got) == 0 {
			t.Fatal("no breach of I3 with nothing done before")
		}
		if !hasBreach(Check(&run.Record{Unit: e.Unit, Steps: []run.Step{step}}), "I3") {
			t.Fatal("Check does not give the breach of I3")
		}
	})

	t.Run("I4", func(t *testing.T) {
		_, next := first(t, rec, func(s run.Step) bool {
			return s.Choice.Kind == engine.ChoiceTask && s.Input.HasLast && !s.Input.Last.InLesson
		})
		bad := next
		bad.Choice.ID = next.Input.Last.Task
		if solvedAgain(bad) == "" {
			t.Fatal("no breach of I4 on the last task again")
		}
		if !hasBreach(Check(&run.Record{Unit: rec.Unit, Steps: []run.Step{bad}}), "I4") {
			t.Fatal("Check does not give the breach of I4")
		}
		inLesson := bad
		inLesson.Input.Last.InLesson = true
		if solvedAgain(inLesson) != "" {
			t.Fatal("a breach of I4 after a task of a lesson")
		}
		other := bad
		other.Choice.ID = "novatrail/olympiad-math/task/enum-12-d1-3"
		if other.Choice.ID == other.Input.Last.Task {
			other.Choice.ID = "novatrail/olympiad-math/task/enum-12-d3-1"
		}
		if solvedAgain(other) != "" {
			t.Fatal("a breach of I4 on another task than the last")
		}
	})
}

func hasBreach(bs []Breach, inv string) bool {
	for _, b := range bs {
		if b.Invariant == inv {
			return true
		}
	}
	return false
}

// I3 on a graph built here: rubric skill A rests on verdict skill V, which rests on rubric skill C;
// a lesson teaches C. A by the author's order before C is done breaches I3 through V; the lesson's
// own skills wait for nothing; a prerequisite nothing can teach is not waited for.
func TestAuthorOrderGraph(t *testing.T) {
	rubric, verdict := learnermodel.AssessabilityRubric, learnermodel.AssessabilityVerdict
	graph := learnermodel.Graph{Module: "m", Skills: []learnermodel.Skill{
		{ID: "m/skill/a", Assessability: rubric, Prerequisites: []string{"m/skill/v"}},
		{ID: "m/skill/b", Assessability: rubric, Prerequisites: []string{"m/skill/c"}},
		{ID: "m/skill/c", Assessability: rubric},
		{ID: "m/skill/v", Assessability: verdict, Prerequisites: []string{"m/skill/c"}},
	}}
	task := func(id, skill string) engine.Candidate {
		return engine.Candidate{Task: learnermodel.Task{ID: "m/task/" + id, Primary: []string{skill}}}
	}
	step := run.Step{
		Choice: engine.Choice{Kind: engine.ChoiceTask, ID: "m/task/on-a", Reason: engine.ReasonAuthorOrder},
		Input: engine.Input{Modules: []engine.Module{{
			Graph:   graph,
			Pool:    []engine.Candidate{task("on-a", "m/skill/a"), task("on-c", "m/skill/c")},
			Lessons: []engine.Lesson{{ID: "m/lesson/bc", Skills: []string{"m/skill/b", "m/skill/c"}}},
		}}},
	}
	if got := authorOrder(step); len(got) != 1 {
		t.Fatalf("A before C is done, through V: %v", got)
	}
	doneByTask := step
	doneByTask.Input.RubricDone = []engine.RubricDone{{Task: "m/task/on-c"}}
	if got := authorOrder(doneByTask); len(got) != 0 {
		t.Fatalf("C done by a task: %v", got)
	}
	doneByRun := step
	doneByRun.Input.Runs = []engine.Run{{Lesson: "m/lesson/bc", Skills: []string{"m/skill/b", "m/skill/c"}, End: engine.RunPassed}}
	if got := authorOrder(doneByRun); len(got) != 0 {
		t.Fatalf("C done by a lesson: %v", got)
	}
	lesson := step
	lesson.Choice = engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/bc", Reason: engine.ReasonAuthorOrder}
	if got := authorOrder(lesson); len(got) != 0 {
		t.Fatalf("a lesson that teaches B with C, which B rests on: %v", got)
	}
	nothing := step
	nothing.Input.Modules = []engine.Module{{Graph: graph, Pool: []engine.Candidate{task("on-a", "m/skill/a")}}}
	if got := authorOrder(nothing); len(got) != 0 {
		t.Fatalf("C, which nothing can teach: %v", got)
	}
}
