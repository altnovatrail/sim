package measure

import (
	"math"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/run"
)

// taskStep is a step that gave a task on a skill, with the model's mastery of it and a pool.
func taskStep(at int64, skill string, mastered bool, pool *run.Pool) run.Step {
	task := learnermodel.Task{ID: "m/task/" + skill, Primary: []string{skill}}
	return run.Step{
		Input: engine.Input{Time: at, Modules: []engine.Module{{Pool: []engine.Candidate{{Task: task}}}},
			Readout: learnermodel.Readout{Modules: []learnermodel.ModuleReadout{{Skills: []learnermodel.SkillReadout{{Skill: skill, Mastered: mastered}}}}}},
		Choice: engine.Choice{Kind: engine.ChoiceTask, ID: task.ID, Reason: engine.ReasonQueue},
		Pool:   pool,
	}
}

// E1 counts the items given by P* when given, and the ceiling over the choices that gave something.
func TestE1(t *testing.T) {
	st := taskStep(learner.Start, "a", false, &run.Pool{Given: true, Ceiling: 0.5})
	st.Obs = []*run.Obs{{Items: []run.ItemObs{{PGiven: 0.7}, {PGiven: 0.85}, {PGiven: 0.3}, {PGiven: 0.99}, {PGiven: 0.6}}}}
	nothing := run.Step{Choice: engine.Choice{Kind: engine.ChoiceNothing}, Pool: &run.Pool{Given: true, Ceiling: 1}}
	vs := Read("e1", nil, []Sums{Engine(&run.Record{Steps: []run.Step{st, nothing}})}, EngineNumbers())
	for name, want := range map[string]float64{"E1.corridor": 0.4, "E1.hard": 0.2, "E1.easy": 0.2, "E1.ceiling": 0.5, "E2.nothing": 0.5} {
		if v := value(t, vs, name); !v.Read || math.Abs(v.Value-want) > 1e-12 {
			t.Errorf("%s = %+v, want %v", name, v, want)
		}
	}
}

// E2: 50 choices in a row on one skill with no change of its mastery, while another skill is
// within reach, is stuck; a change of mastery, another skill between, or no other skill within
// reach starts the count again.
func TestE2(t *testing.T) {
	reach := &run.Pool{Given: true, Reachable: []string{"a", "b"}}
	alone := &run.Pool{Given: true, Reachable: []string{"a"}}
	steps := func(n int, pool *run.Pool, flip int, between int) []run.Step {
		var out []run.Step
		for k := range n {
			if k == between {
				out = append(out, taskStep(0, "b", false, reach))
			}
			out = append(out, taskStep(0, "a", k >= flip && flip >= 0, pool))
		}
		return out
	}
	pair := func(n int) []run.Step {
		var out []run.Step
		for range n {
			st := taskStep(0, "a", false, reach)
			st.Input.Modules[0].Pool[0].Task.Primary = []string{"a", "b"}
			out = append(out, st)
		}
		return out
	}
	for _, c := range []struct {
		name  string
		steps []run.Step
		stuck float64
	}{
		{"fifty", steps(50, reach, -1, -1), 1},
		{"forty-nine", steps(49, reach, -1, -1), 0},
		{"mastered between", steps(50, reach, 25, -1), 0},
		{"another between", steps(50, reach, -1, 25), 0},
		{"nothing else within reach", steps(50, alone, -1, -1), 0},
		{"the other skill in the same task", pair(50), 0},
	} {
		if s := Engine(&run.Record{Steps: c.steps}); s["E2.stuck"] != c.stuck {
			t.Errorf("%s: stuck %v, want %v", c.name, s["E2.stuck"], c.stuck)
		}
	}
}

// E3 reads the repetitions: R* of the skill due and whether it is late, and the share of
// repetitions by week.
func TestE3(t *testing.T) {
	review := func(at int64, rstar, pref float64) run.Step {
		st := taskStep(at, "a", true, nil)
		st.Choice.Reason = engine.ReasonRepetition
		st.Input.Readout.Modules[0].Skills[0].Review = learnermodel.ReviewDate{At: at}
		// b is due too, but the choice is not on it.
		st.Input.Readout.Modules[0].Skills = append(st.Input.Readout.Modules[0].Skills, learnermodel.SkillReadout{Skill: "b", Mastered: true, Review: learnermodel.ReviewDate{Already: true}})
		st.Obs = []*run.Obs{{Before: []run.SkillObs{{Skill: "a", Retention: rstar, PTrueRef: pref}, {Skill: "b", Retention: 0.1}}}}
		return st
	}
	notDue := review(learner.Start+week, 0.2, 0.1)
	notDue.Input.Readout.Modules[0].Skills[0].Review = learnermodel.ReviewDate{At: learner.Start + 2*week}
	rec := &run.Record{Steps: []run.Step{
		review(learner.Start, 0.9, 0.8), review(learner.Start+day, 0.7, 0.6), taskStep(learner.Start+2*day, "a", true, nil),
		notDue, taskStep(learner.Start+week+day, "a", true, nil),
	}}
	vs := Read("e3", nil, []Sums{Engine(rec)}, EngineNumbers())
	for name, want := range map[string]float64{"E3.retention": 0.8, "E3.late": 0.5, "E3.week1": 2.0 / 3, "E3.week2": 0.5} {
		if v := value(t, vs, name); !v.Read || math.Abs(v.Value-want) > 1e-12 {
			t.Errorf("%s = %+v, want %v", name, v, want)
		}
	}
	if v := value(t, vs, "E3.week3"); v.Read {
		t.Error("a week with no choice is read")
	}
}

// E4 and E5 on hand-built steps: the skills within reach given, the tasks with a rubric written,
// each once, and a task closed by "Done" given again, as a task or in a lesson.
func TestE4E5(t *testing.T) {
	pool := &run.Pool{Given: true, Reachable: []string{"a", "b"}}
	written := taskStep(0, "e", false, nil)
	written.Attempts = []learner.Attempt{{Task: "m/task/e", Closure: learner.ClosedDone}}
	again := taskStep(0, "e", false, nil)
	again.Input.RubricDone = []engine.RubricDone{{Task: "m/task/e"}}
	again.Attempts = []learner.Attempt{{Task: "m/task/e", Closure: learner.ClosedDone}}
	lesson := run.Step{Input: engine.Input{RubricDone: again.Input.RubricDone}, Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/x"},
		Attempts: []learner.Attempt{{Task: "m/task/e", Closure: learner.ClosedLeft}}}
	s := Engine(&run.Record{Steps: []run.Step{taskStep(0, "a", false, pool), written, again, lesson}})
	if s["E4.reach"] != 2 || s["E4.given"] != 1 || s["E5.written"] != 1 || s["E5.again"] != 2 {
		t.Errorf("E4 %v of %v, E5 %v written and %v again", s["E4.given"], s["E4.reach"], s["E5.written"], s["E5.again"])
	}
	if s["E4.skills"] != 0 || s["E3.end.n"] != 0 {
		t.Error("a run with no learner reads the truth at its end")
	}
}

// The measures of the engine read off real runs: G0 on school maths, F-model on olympiad maths,
// E-base.
func TestEngineOnRuns(t *testing.T) {
	for _, c := range []struct {
		population string
		graph, n   int
		numbers    []string
	}{
		{"G0", 1, 12, []string{"E1.corridor", "E1.ceiling", "E2.nothing", "E3.end", "E4.given", "E4.mastered"}},
		// F-model repeats on olympiad maths, its 19 skills, enough to read every resample.
		{"F-model", 0, 24, []string{"E3.retention", "E3.late", "E3.week1"}},
		{"E-base", 0, 12, []string{"E5.written", "E5.again", "E4.given"}},
	} {
		recs, graph := recordsOn(t, c.population, c.graph, c.n)
		var sums []Sums
		for _, r := range recs {
			sums = append(sums, Engine(r))
		}
		vs := Read(c.population, graph, sums, EngineNumbers())
		for _, name := range c.numbers {
			if v := value(t, vs, name); !v.Read || v.Low > v.Value || v.High < v.Value || v.Value < 0 {
				t.Errorf("%s %s: %+v", c.population, name, v)
			}
		}
		if c.population == "G0" {
			corridor, ceiling := value(t, vs, "E1.corridor"), value(t, vs, "E1.ceiling")
			if ceiling.Value <= 0 || ceiling.Value > 1 || corridor.Value > 1 {
				t.Errorf("E1: corridor %v, ceiling %v", corridor.Value, ceiling.Value)
			}
		}
		if c.population == "E-base" && value(t, vs, "E5.written").Value <= 0 {
			t.Error("E-base writes no task with a rubric")
		}
	}
}
