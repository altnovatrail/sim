package measure

import (
	"math"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/run"
)

// lessonInput is an input whose module has the lesson m/lesson/x on the skill s, with the
// assessability given.
func lessonInput(a learnermodel.Assessability) engine.Input {
	return engine.Input{Modules: []engine.Module{{
		Graph:   learnermodel.Graph{Skills: []learnermodel.Skill{{ID: "s", Assessability: a}}},
		Lessons: []engine.Lesson{{ID: "m/lesson/x", Skills: []string{"s"}}},
	}}}
}

// link is an attempt of a lesson's link: its items by P* when given, and its result on a skill.
func link(skill string, failed bool, ps ...float64) (*run.Obs, []engine.SkillResult) {
	o := &run.Obs{}
	for _, p := range ps {
		o.Items = append(o.Items, run.ItemObs{PGiven: p})
	}
	r := engine.SkillResult{Skill: skill, Answered: 2, FirstRight: 2}
	if failed {
		r.FirstRight = 0
	}
	return o, []engine.SkillResult{r}
}

// The note reads a lesson with a consolidation after its first link: the consolidation a pass, the
// shares of the links' items and of the consolidation's, and P* of the next link after the failure,
// which comes after the consolidation (R48, R49).
func TestNote(t *testing.T) {
	in := lessonInput(learnermodel.AssessabilityVerdict)
	lesson := run.Step{Input: in, Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/x", Reason: engine.ReasonNewSkill}}
	o1, r1 := link("s", true, 0.75, 0.3)
	lesson.Obs, lesson.Results = []*run.Obs{o1}, [][]engine.SkillResult{r1}
	consolidation := taskStep(0, "s", false, &run.Pool{Given: true})
	consolidation.Choice.Reason = engine.ReasonConsolidation
	oc, rc := link("s", true, 0.8)
	consolidation.Obs, consolidation.Results = []*run.Obs{oc}, [][]engine.SkillResult{rc}
	back := run.Step{Input: in, Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/x", Reason: engine.ReasonOpenRun}, Resumed: true, Pass: &learner.Pass{}}
	o2, r2 := link("s", false, 0.97, 0.6)
	o3, r3 := link("s", false, 0.72)
	back.Obs, back.Results = []*run.Obs{o2, nil, o3}, [][]engine.SkillResult{r2, nil, r3}
	after := taskStep(0, "s", false, nil) // a task after the pass is no consolidation in a lesson
	after.Obs = []*run.Obs{{Items: []run.ItemObs{{PGiven: 0.1}}}}
	s := Note(&run.Record{Steps: []run.Step{lesson, consolidation, back, after}})
	vs := Read("note", nil, []Sums{s}, NoteNumbers())
	for name, want := range map[string]float64{
		"E6.consolidations": 1, "E6.lesson.corridor": 0.4, "E6.lesson.hard": 0.2, "E6.lesson.easy": 0.2,
		"E6.consolidation.corridor": 1, "E6.consolidation.hard": 0, "E6.next": (0.97 + 0.6) / 2,
	} {
		if v := value(t, vs, name); !v.Read || math.Abs(v.Value-want) > 1e-12 {
			t.Errorf("%s = %+v, want %v", name, v, want)
		}
	}
	if s["E6.lesson.items"] != 5 || s["E6.passes"] != 1 || s["E6.next.items"] != 2 {
		t.Errorf("the sums %v", s)
	}
}

// A lesson with a skill without a verdict is not read, the engine is not called within it; a link
// that does not fail, or fails on a skill not of the lesson, gives no next link; and a run with no
// lesson reads nothing.
func TestNoteNothing(t *testing.T) {
	rubric := run.Step{Input: lessonInput(learnermodel.AssessabilityRubric), Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/x"}, Pass: &learner.Pass{}}
	o, r := link("s", false, 0.7)
	rubric.Obs, rubric.Results = []*run.Obs{o}, [][]engine.SkillResult{r}
	if s := Note(&run.Record{Steps: []run.Step{rubric}}); len(s) != 0 {
		t.Errorf("a lesson without a verdict: %v", s)
	}
	for _, failedOn := range []string{"", "other"} {
		st := run.Step{Input: lessonInput(learnermodel.AssessabilityVerdict), Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/x"}, Pass: &learner.Pass{}}
		o1, r1 := link("s", false, 0.7)
		if failedOn != "" {
			o1, r1 = link(failedOn, true, 0.7)
		}
		o2, r2 := link("s", false, 0.8)
		st.Obs, st.Results = []*run.Obs{o1, o2}, [][]engine.SkillResult{r1, r2}
		s := Note(&run.Record{Steps: []run.Step{st}})
		if s["E6.next.items"] != 0 || s["E6.lesson.items"] != 2 || s["E6.passes"] != 1 {
			t.Errorf("failed on %q: %v", failedOn, s)
		}
	}
	// A task with a rubric after a failed link is the next link, though it has no items to read; the
	// link after it is not.
	st := run.Step{Input: lessonInput(learnermodel.AssessabilityVerdict), Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/x"}, Pass: &learner.Pass{}}
	o1, r1 := link("s", true, 0.7)
	_, r2 := link("s", false)
	o3, r3 := link("s", false, 0.8)
	st.Obs, st.Results = []*run.Obs{o1, nil, o3}, [][]engine.SkillResult{r1, r2, r3}
	if s := Note(&run.Record{Steps: []run.Step{st}}); s["E6.next.items"] != 0 || s["E6.lesson.items"] != 2 {
		t.Errorf("a link with a rubric after the failure: %v", s)
	}
	if s := Note(&run.Record{Steps: []run.Step{taskStep(0, "s", false, nil)}}); len(s) != 0 {
		t.Errorf("a run with no lesson: %v", s)
	}
	unknown := run.Step{Input: lessonInput(learnermodel.AssessabilityVerdict), Choice: engine.Choice{Kind: engine.ChoiceLesson, ID: "m/lesson/y"}}
	if s := Note(&run.Record{Steps: []run.Step{unknown}}); len(s) != 0 {
		t.Errorf("a lesson the input does not have: %v", s)
	}
}

// On real runs of the base the note reads the lessons of school maths, and no consolidation: the
// base answers "further" at the end of every link (R48).
func TestNoteOnRuns(t *testing.T) {
	recs, graph := recordsOn(t, "G2", 1, 12)
	var sums []Sums
	for _, r := range recs {
		sums = append(sums, Note(r))
	}
	vs := Read("G2", graph, sums, NoteNumbers())
	if v := value(t, vs, "E6.lesson.corridor"); !v.Read || v.Value < 0 || v.Value > 1 {
		t.Errorf("E6.lesson.corridor %+v", v)
	}
	if v := value(t, vs, "E6.consolidations"); !v.Read || v.Value != 0 {
		t.Errorf("E6.consolidations %+v", v)
	}
	if v := value(t, vs, "E6.consolidation.corridor"); v.Read {
		t.Errorf("the base consolidated: %+v", v)
	}
}
