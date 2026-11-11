package synthetic

import (
	"slices"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/invariant"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
)

func modules(t *testing.T) map[string]*pack.Module {
	t.Helper()
	ms, err := Modules()
	if err != nil {
		t.Fatal(err)
	}
	if again, _ := Modules(); again[Sheet] != ms[Sheet] {
		t.Fatal("Modules translated the views again")
	}
	return ms
}

// kinds counts the kinds of check of a module's cells, and its tasks with a rubric apart.
func kinds(m *pack.Module) (map[learnermodel.CheckKind]int, int) {
	out, rubrics := map[learnermodel.CheckKind]int{}, 0
	for _, task := range m.Tasks {
		if task.Rubric {
			rubrics++
		}
		for _, it := range task.Candidate.Task.Items {
			for _, s := range it.Slots {
				out[s.Check]++
			}
		}
	}
	return out, rubrics
}

// The synthetic packages have the shape docs/architecture.md, section 3.3, gives them: a sheet with
// a choice, numbers, a predicate and a time, and leading questions, over three bands of grades 1–6;
// a workbook of one band of grade 3 with pages of up to ten lines; essays with a rubric and no
// grades; and a skill of the sheet that rests on the anchor skill a skill of the workbook teaches.
func TestShape(t *testing.T) {
	ms := modules(t)
	sheet, workbook, essays := ms[Sheet], ms[Workbook], ms[Essays]
	if k, _ := kinds(sheet); k[learnermodel.CheckKindChoice] == 0 || k[learnermodel.CheckKindNumber] == 0 || k[learnermodel.CheckKindPredicate] == 0 || k[learnermodel.CheckKindTime] == 0 {
		t.Errorf("the sheet has the kinds %v", k)
	}
	if !slices.ContainsFunc(sheet.Tasks, func(t pack.Task) bool { return t.Leading }) {
		t.Error("no task of the sheet has a leading question")
	}
	if len(sheet.Graph.Bands) != 3 || !slices.Equal(sheet.Graph.Bands[2].Grades, []int32{5, 6}) {
		t.Errorf("the sheet's bands %v", sheet.Graph.Bands)
	}
	if len(workbook.Graph.Bands) != 1 || !slices.Equal(workbook.Graph.Bands[0].Grades, []int32{3}) {
		t.Errorf("the workbook's bands %v", workbook.Graph.Bands)
	}
	longest := 0
	for _, task := range workbook.Tasks {
		if task.Canvas != "workbook" {
			t.Errorf("%s is on %q", task.Candidate.Task.ID, task.Canvas)
		}
		longest = max(longest, len(task.Candidate.Task.Items))
	}
	if longest != 10 {
		t.Errorf("the workbook's longest page has %d lines, want 10", longest)
	}
	if k, rubrics := kinds(essays); rubrics == 0 || k[learnermodel.CheckKindChoice] == 0 || k[learnermodel.CheckKindUnspecified] != 0 {
		t.Errorf("the essays have the kinds %v and %d tasks with a rubric", k, rubrics)
	}
	for id := range ms {
		if !IsSynthetic(id) {
			t.Errorf("%s is not named synthetic", id)
		}
	}
	if IsSynthetic("novatrail/olympiad-math") {
		t.Error("a published module is named synthetic")
	}
	for _, b := range essays.Graph.Bands {
		if len(b.Grades) != 0 {
			t.Errorf("the essays' band %s names grades", b.ID)
		}
	}
	rests := slices.ContainsFunc(sheet.Graph.Skills, func(s learnermodel.Skill) bool { return slices.Contains(s.Prerequisites, Anchor) })
	teaches := slices.ContainsFunc(workbook.Graph.Skills, func(s learnermodel.Skill) bool { return slices.Contains(s.Teaches, Anchor) })
	if !rests || !teaches {
		t.Errorf("the anchor: the sheet rests on it %v, the workbook teaches it %v", rests, teaches)
	}
	if pool, _ := essays.ForAge("under_13"); len(pool) != len(essays.Tasks)-1 {
		t.Errorf("a child gets %d tasks of the synthetic essays, want all but the one for adults", len(pool))
	}
}

func runOf(t *testing.T, population string, n int, limit run.Limit, mods ...*pack.Module) *run.Record {
	t.Helper()
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
		t.Fatalf("%s: %v", l.Unit, err)
	}
	return rec
}

// The base population runs on every synthetic package and on the sheet with the workbook, with the
// invariants held; a learner of both reads the anchor skill, whose status the workbook decides.
func TestRuns(t *testing.T) {
	ms := modules(t)
	anchors := map[learnermodel.AnchorStatus]int{}
	for n := range 8 {
		for _, c := range []struct {
			population string
			limit      run.Limit
			mods       []*pack.Module
		}{
			{"G0", run.Maths, []*pack.Module{ms[Sheet]}},
			{"G0", run.Maths, []*pack.Module{ms[Workbook]}},
			{"A-linked", run.Maths, []*pack.Module{ms[Sheet], ms[Workbook]}},
			{"E-base", run.Essays, []*pack.Module{ms[Essays]}},
		} {
			rec := runOf(t, c.population, n, c.limit, c.mods...)
			for _, b := range invariant.Check(rec) {
				t.Error(b)
			}
			if c.population != "A-linked" {
				continue
			}
			for _, s := range rec.Steps {
				for _, a := range s.Input.Readout.Anchors {
					if a.Anchor == Anchor {
						anchors[a.Status]++
					}
				}
			}
		}
	}
	if anchors[learnermodel.AnchorHas] == 0 || anchors[learnermodel.AnchorHasNot] == 0 {
		t.Errorf("the anchor is read %v: a learner of both should not hold it, then hold it", anchors)
	}
}

// A learner who asks for the technique, as G8 does, gets the leading question with it on the tasks
// of the sheet that have one: the evidence carries it, of the kind that records it, and only there
// (R41 in docs/decisions.md).
func TestLeadingQuestion(t *testing.T) {
	ms := modules(t)
	a := learner.Base()
	a.TechniqueChance, a.TechniqueLift = 0.2, 0.5
	asked := 0
	for n := range 4 {
		l, err := learner.New(seed.Working, "G8", n, []*pack.Module{ms[Sheet]}, a)
		if err != nil {
			t.Fatal(err)
		}
		r, err := run.New(l, []*pack.Module{ms[Sheet]}, run.Maths)
		if err != nil {
			t.Fatal(err)
		}
		rec, err := r.All()
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range rec.Steps {
			for i, at := range s.Attempts {
				e := s.Evidence[i]
				if e == nil {
					continue
				}
				task, err := ms[Sheet].Task(at.Task)
				if err != nil {
					t.Fatal(err)
				}
				if e.LeadingQuestionOnRequest != at.LeadingQuestion || at.LeadingQuestion != (at.Technique && task.Leading) {
					t.Fatalf("%s: a leading question %v in the evidence, %v in the attempt, the technique %v", at.Task, e.LeadingQuestionOnRequest, at.LeadingQuestion, at.Technique)
				}
				if !at.LeadingQuestion {
					continue
				}
				asked++
				data, err := learnermodel.EncodeEvidence(*e)
				if err != nil {
					t.Fatal(err)
				}
				if back, err := learnermodel.DecodeEvidence(data); err != nil || back.Kind != learnermodel.EvidenceKind || !back.LeadingQuestionOnRequest {
					t.Fatalf("%s: the evidence reads back as kind %d with a leading question %v: %v", at.Task, back.Kind, back.LeadingQuestionOnRequest, err)
				}
			}
		}
	}
	if asked == 0 {
		t.Error("no learner got a leading question")
	}
}
