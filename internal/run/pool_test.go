package run

import (
	"slices"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// A task whose solution went less than 24 hours ago is not given in the first pass; in the second
// it is, but for the task of the last attempt; a template is given always.
func TestGiven(t *testing.T) {
	solved := func(id string) (int64, bool) { return 100, id != "fresh" }
	task := func(id string, template bool) learnermodel.Task { return learnermodel.Task{ID: id, Template: template} }
	for _, c := range []struct {
		task         learnermodel.Task
		at           int64
		second, want bool
	}{
		{task("fresh", false), 200, false, true},
		{task("seen", false), 200, false, false},
		{task("seen", false), 200, true, true},
		{task("last", false), 200, true, false},
		{task("seen", false), 100 + pause, false, true},
		{task("last", true), 200, false, true},
	} {
		if got := given(c.task, c.at, solved, "last", c.second); got != c.want {
			t.Errorf("%s at %d, second %v: %v", c.task.ID, c.at, c.second, got)
		}
	}
}

// Each step reads the pool before the choice: a ceiling among shares, skills within reach that
// have a verdict, and a pool that gives while the engine gives something.
func TestPool(t *testing.T) {
	for _, id := range []string{olympiad, grade3} {
		mods := modules(t, id)
		rec := runOf(t, "G0", 1, Maths, mods)
		reached := 0
		for _, s := range rec.Steps {
			p := s.Pool
			if p == nil || p.Ceiling < 0 || p.Ceiling > 1 {
				t.Fatalf("%s: the pool %+v", rec.Unit, p)
			}
			if s.Choice.Kind == engine.ChoiceTask && !p.Given {
				t.Errorf("%s: a task given from a pool that gives none", rec.Unit)
			}
			for _, sk := range p.Reachable {
				if !slices.ContainsFunc(mods[0].Graph.Skills, func(g learnermodel.Skill) bool {
					return g.ID == sk && g.Assessability == learnermodel.AssessabilityVerdict
				}) {
					t.Errorf("%s: %s within reach has no verdict", rec.Unit, sk)
				}
			}
			reached = max(reached, len(p.Reachable))
		}
		if reached == 0 {
			t.Errorf("%s: no skill ever within reach", id)
		}
	}
}

// The truth holds a prerequisite skill once it masters it, and an anchor once a skill that teaches
// it is mastered; a skill of the sheet that rests on the anchor is out of reach before.
func TestHolds(t *testing.T) {
	ms, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	sheet, workbook, grade := ms[synthetic.Sheet], ms[synthetic.Workbook], modules(t, grade3)[0]
	onAnchor, withPre := "", ""
	for _, s := range sheet.Graph.Skills {
		if slices.Contains(s.Prerequisites, synthetic.Anchor) {
			onAnchor = s.ID
		}
	}
	for _, s := range grade.Graph.Skills {
		if len(s.Prerequisites) > 0 {
			withPre = s.ID
		}
	}
	if onAnchor == "" || withPre == "" {
		t.Fatalf("no skill on the anchor, %q, or on another skill, %q", onAnchor, withPre)
	}
	strong, weak := learner.Base(), learner.Base()
	strong.LevelShift, weak.LevelShift = 6, -6
	check := func(a learner.Assumptions, want bool) {
		for _, c := range []struct {
			mods  []*pack.Module
			skill string
		}{{[]*pack.Module{sheet, workbook}, onAnchor}, {[]*pack.Module{grade}, withPre}} {
			l, err := learner.New(seed.Working, "A-linked", 0, c.mods, a)
			if err != nil {
				t.Fatal(err)
			}
			r, err := New(l, c.mods, Maths)
			if err != nil {
				t.Fatal(err)
			}
			if r.holds(c.mods[0], c.skill, learner.Start) != want {
				t.Errorf("level shift %v: %s held %v", a.LevelShift, c.skill, !want)
			}
			if !r.holds(c.mods[0], "synthetic/none/skill/none", learner.Start) {
				t.Error("a skill the graph does not have is not held")
			}
		}
	}
	check(strong, true)
	check(weak, false)
}

// The skills of a choice are its task's primary skills as learner-model counts them: the items'
// own, and the task's only when an item takes them; a lesson's are its skills.
func TestStepSkills(t *testing.T) {
	own := learnermodel.Task{ID: "m/task/own", Primary: []string{"z"}, Items: []learnermodel.TaskItem{
		{OwnSkills: true, Primary: []string{"y", "x"}}, {OwnSkills: true, Primary: []string{"x"}}}}
	takes := own
	takes.ID, takes.Items = "m/task/takes", append(slices.Clone(own.Items), learnermodel.TaskItem{})
	bare := learnermodel.Task{ID: "m/task/bare", Primary: []string{"z"}}
	in := engine.Input{Modules: []engine.Module{{
		Pool:    []engine.Candidate{{Task: own}, {Task: takes}, {Task: bare}},
		Lessons: []engine.Lesson{{ID: "m/lesson/l", Skills: []string{"w"}}},
	}}}
	for id, want := range map[string][]string{"m/task/own": {"x", "y"}, "m/task/takes": {"x", "y", "z"}, "m/task/bare": {"z"}} {
		if got := (Step{Input: in, Choice: engine.Choice{Kind: engine.ChoiceTask, ID: id}}).Skills(); !slices.Equal(got, want) {
			t.Errorf("%s: %v, want %v", id, got, want)
		}
	}
	if got := (Step{Input: in, Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/l"}}).Skills(); !slices.Equal(got, []string{"w"}) {
		t.Errorf("the lesson: %v", got)
	}
	if got := (Step{Input: in, Choice: engine.Choice{Kind: engine.ChoiceNothing}}).Skills(); got != nil {
		t.Errorf("nothing: %v", got)
	}
}
