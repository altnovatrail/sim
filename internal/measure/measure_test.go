package measure

import (
	"math"
	"reflect"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

func value(t *testing.T, vs []Value, name string) Value {
	t.Helper()
	for _, v := range vs {
		if v.Name == name {
			return v
		}
	}
	t.Fatalf("no number %s", name)
	return Value{}
}

// A number of a cell is read from the sums of its learners, its interval from 2,000 resamples drawn
// the same way every time; a number with nothing to read is not read; a graph with a synthetic
// module is no graph.
func TestRead(t *testing.T) {
	learners := []Sums{{"x": 1, "n": 1}, {"x": 3, "n": 1}, {"x": 5, "n": 1}, {"x": 7, "n": 1}}
	numbers := []Number{{"mean", ratio("x", "n")}, {"none", ratio("x", "missing")}, {"rms", root("x", "n")}}
	a, b := Read("cell", []string{"novatrail/olympiad-math"}, learners, numbers), Read("cell", nil, learners, numbers)
	mean := value(t, a, "mean")
	if !mean.Read || mean.Value != 4 || mean.Low > 4 || mean.High < 4 || mean.Low < 1 || mean.High > 7 || mean.Low == mean.High {
		t.Errorf("mean %+v", mean)
	}
	if !reflect.DeepEqual(a[0], b[0]) || !reflect.DeepEqual(a[2], b[2]) {
		t.Error("two reads of one cell give two intervals")
	}
	if value(t, a, "none").Read {
		t.Error("a number with nothing to read is read")
	}
	if got := value(t, a, "rms").Value; math.Abs(got-2) > 1e-12 {
		t.Errorf("rms = %v", got)
	}
	if mean.NoGraph {
		t.Error("a published graph is no graph")
	}
	if v := Read("cell", []string{"novatrail/olympiad-math", synthetic.Sheet}, learners, numbers); !v[0].NoGraph || v[0].Value != mean.Value || v[0].Low != mean.Low || v[0].High != mean.High || !v[0].Read {
		t.Errorf("a graph with a synthetic module: %+v", v[0])
	}
	if v := Read("empty", nil, nil, numbers); v[0].Read {
		t.Error("an empty cell reads a number")
	}
}

func skill(id string, evidence int64, est, truth float64, mastered, truly bool) run.SkillObs {
	return run.SkillObs{Skill: id, Evidence: evidence, Estimate: est, Truth: truth, Mastered: mastered, TrulyMastered: truly, Band: 1}
}

// The measures of the model on a record built by hand: what each reads and counts.
func TestModelSums(t *testing.T) {
	const a = "m/skill/a"
	obs := []*run.Obs{
		{At: 0, Answers: 0, Items: []run.ItemObs{{PModel: 0.6, PFresh: 0.7, PTrue: 0.4, Pause: -1}},
			Before: []run.SkillObs{skill(a, 0, 0, 0, false, false)}, After: []run.SkillObs{skill(a, 1, 0.5, 0, false, false)}, Level: 1, TrueLevel: 0.5},
		{At: 20 * day, Answers: 1, Repeat: true, Items: []run.ItemObs{{PModel: 0.8, PFresh: 0.9, PTrue: 0.9, Pause: 20 * day, Help: true}, {PModel: 0.2, PFresh: 0.25, PTrue: 0.3, Pause: 20 * day, Many: true}},
			Before: []run.SkillObs{{Skill: a, Evidence: 1, Estimate: 0.5, Last: 100}}, After: []run.SkillObs{skill(a, 2, 0.7, 0.2, true, false)}, Level: 1, TrueLevel: 0.5},
	}
	rec := &run.Record{Steps: []run.Step{
		{Attempts: []learner.Attempt{{Items: []learner.ItemResult{{FirstRight: true}}}}, Obs: obs[:1]},
		{Attempts: []learner.Attempt{{Items: []learner.ItemResult{{FirstRight: true}, {FirstRight: false}}}}, Obs: obs[1:]},
	}}
	s := Model(rec)
	vs := Read("hand", nil, []Sums{s}, ModelNumbers())
	checks := map[string]float64{
		"M2.brier.all":    (0.16 + 0.04 + 0.04) / 3,
		"M2.bias.all":     (0.6+0.8+0.2)/3 - 2.0/3,
		"M2.brier.pause7": (0.04 + 0.04) / 2,
		"M2.brier.pages":  0.04,
		"M2.brier.help":   0.04,
		"M4.false.all":    1,
		"M4.false.break":  1,
		"M4.false.pages":  1,
		"M7.hard":         2.0 / 3,
		"M7.streak":       1,
		"M8.shift":        0.2,
		"M10.retention":   (0.04 + 0.04) / 2,
		"M10.fresh":       (0.01 + 0.0625) / 2,
		"M1.mean.5":       0, // not reached: three answers
	}
	for name, want := range checks {
		v := value(t, vs, name)
		if name == "M1.mean.5" {
			if v.Read {
				t.Errorf("M1 after 5 answers is read after 3")
			}
			continue
		}
		if !v.Read || math.Abs(v.Value-want) > 1e-12 {
			t.Errorf("%s = %v (read %v), want %v", name, v.Value, v.Read, want)
		}
	}
	if v := value(t, vs, "M7.error.5"); v.Read {
		t.Error("M7 after 5 answers read after 3")
	}
}

// The late mastery, M6, and the catch-up after a pause, M5, on records built by hand.
func TestLateAndCatch(t *testing.T) {
	const a = "m/skill/a"
	step := func(o *run.Obs, rights ...bool) run.Step {
		var items []learner.ItemResult
		for _, r := range rights {
			items = append(items, learner.ItemResult{FirstRight: r})
		}
		return run.Step{Attempts: []learner.Attempt{{Items: items}}, Obs: []*run.Obs{o}}
	}
	one := []run.ItemObs{{PModel: 0.5, Pause: -1}}
	late := &run.Record{Steps: []run.Step{
		step(&run.Obs{Items: one, After: []run.SkillObs{skill(a, 1, 0, 0, false, true)}}, true),
		step(&run.Obs{Items: one, After: []run.SkillObs{skill(a, 2, 0, 0, false, true)}}, true),
		step(&run.Obs{Items: one, After: []run.SkillObs{skill(a, 3, 0, 0, true, true)}}, true),
	}}
	vs := Read("late", nil, []Sums{Model(late)}, ModelNumbers())
	if v := value(t, vs, "M6.late"); !v.Read || v.Value != 2 {
		t.Errorf("M6.late = %+v, want 2", v)
	}
	if v := value(t, vs, "M6.never"); v.Value != 0 {
		t.Errorf("M6.never = %+v", v)
	}
	ref := func(p, pt float64, evidence, last int64) []run.SkillObs {
		return []run.SkillObs{{Skill: a, Evidence: evidence, PRef: p, PTrueRef: pt, Last: last}}
	}
	var steps []run.Step
	steps = append(steps, step(&run.Obs{At: 30 * day, Items: one, Before: ref(0.9, 0.5, 1, 1), After: ref(0.9, 0.5, 2, 30*day)}, true))
	for k := range 6 {
		steps = append(steps, step(&run.Obs{At: int64(31+k) * day, Items: one, Before: ref(0.6, 0.58, int64(2+k), int64(30+k)*day), After: ref(0.6, 0.58, int64(3+k), int64(31+k)*day)}, true))
	}
	vs = Read("catch", nil, []Sums{Model(&run.Record{Steps: steps})}, ModelNumbers())
	if v := value(t, vs, "M5.caught"); !v.Read || v.Value != 1 {
		t.Errorf("M5.caught = %+v, want 1: one answer far, then five near", v)
	}
	if v := value(t, vs, "M5.never"); v.Value != 0 {
		t.Errorf("M5.never = %+v", v)
	}
}

// cell runs n learners of a population on its first graph and gives their sums and the graph.
func cell(t *testing.T, name string, n int) ([]Sums, []string) {
	t.Helper()
	recs, graph := records(t, name, n)
	var sums []Sums
	for _, rec := range recs {
		sums = append(sums, Model(rec))
	}
	return sums, graph
}

// records runs n learners of a population on its first graph and gives their records and the graph.
func records(t *testing.T, name string, n int) ([]*run.Record, []string) {
	t.Helper()
	ms, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	syn, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	p, ok := population.Named(name)
	if !ok {
		t.Fatalf("no population %s", name)
	}
	var mods []*pack.Module
	for _, id := range p.Graphs[0] {
		if m := ms[id]; m != nil {
			mods = append(mods, m)
		} else {
			mods = append(mods, syn[id])
		}
	}
	var recs []*run.Record
	for k := range n {
		l, err := learner.New(seed.Working, p.Name, k, mods, p.Assumptions(k))
		if err != nil {
			t.Fatal(err)
		}
		r, err := run.New(l, mods, p.Limit)
		if err != nil {
			t.Fatal(err)
		}
		rec, err := r.All()
		if err != nil {
			t.Fatal(err)
		}
		recs = append(recs, rec)
	}
	return recs, p.Graphs[0]
}

// readCell reads the measures of the model off n learners of a population.
func readCell(t *testing.T, name string, n int) []Value {
	t.Helper()
	sums, graph := cell(t, name, n)
	return Read(name, graph, sums, ModelNumbers())
}

// The measures read off real runs: G0 on olympiad maths; the settling after the jump of G3 and
// after the break of F-break; the anchor of a learner of both synthetic maths, A-linked.
func TestModelOnRuns(t *testing.T) {
	vs := readCell(t, "G0", 8)
	for _, name := range []string{"M1.rms.200", "M2.brier.all", "M2.ece.all", "M7.error.10", "M7.streak"} {
		if v := value(t, vs, name); !v.Read || v.Low > v.Value || v.High < v.Value {
			t.Errorf("%s: %+v", name, v)
		}
	}
	if v := value(t, vs, "M2.brier.all"); v.Value <= 0 || v.Value >= 0.25 {
		t.Errorf("Brier %v is no better than a coin", v.Value)
	}
	for _, name := range []string{"G3", "F-break"} {
		s := readCell(t, name, 6)
		settled, never := value(t, s, "M3.settle"), value(t, s, "M3.unsettled")
		if !never.Read || (settled.Read && settled.Value < 0) {
			t.Errorf("%s: settle %+v, never %+v", name, settled, never)
		}
	}
	a := readCell(t, "A-linked", 6)
	if has, not := value(t, a, "M9.hasFalse"), value(t, a, "M9.notFalse"); !has.Read && !not.Read || has.NoGraph {
		t.Error("A-linked reads no anchor, or reads it on no published graph")
	}
}

// M9 reads the anchors at the start of each session against the truth.
func TestAnchorReadings(t *testing.T) {
	in := engine.Input{Time: 0, Readout: learnermodel.Readout{Anchors: []learnermodel.AnchorReadout{{Anchor: "core.x", Status: learnermodel.AnchorHas}}}}
	rec := &run.Record{Steps: []run.Step{{Input: in}}}
	if s := Model(rec); s["M9.has"] != 0 {
		t.Error("M9 read without a learner")
	}
}

// Only the attempts that answer a skill count its answers: a skill of the module left alone is not
// waited on by M5 or M6, and a mastery declared with the attempt that the truth mastered the skill
// waits no answers.
func TestPractised(t *testing.T) {
	const a, b = "m/skill/a", "m/skill/b"
	one := []run.ItemObs{{PModel: 0.5, Pause: -1}}
	step := func(o *run.Obs) run.Step {
		return run.Step{Attempts: []learner.Attempt{{Items: []learner.ItemResult{{FirstRight: true}}}}, Obs: []*run.Obs{o}}
	}
	// a is answered, mastered by the truth and the model at once; b, answered long ago, waits.
	last := func(o run.SkillObs, at int64) run.SkillObs { o.Last = at; return o }
	rec := &run.Record{Steps: []run.Step{
		step(&run.Obs{At: 30 * day,
			Before: []run.SkillObs{{Skill: a}, {Skill: b, Evidence: 1, Last: 1}},
			After:  []run.SkillObs{skill(a, 1, 0, 0, true, true), {Skill: b, Evidence: 1, Last: 1, TrulyMastered: true}}, Items: one}),
		step(&run.Obs{At: 31 * day,
			Before: []run.SkillObs{last(skill(a, 1, 0, 0, true, true), 30*day), {Skill: b, Evidence: 1, Last: 1, TrulyMastered: true}},
			After:  []run.SkillObs{last(skill(a, 2, 0, 0, true, true), 31*day), {Skill: b, Evidence: 1, Last: 1, TrulyMastered: true}}, Items: one}),
	}}
	s := Model(rec)
	if s["M5.episodes"] != 0 {
		t.Errorf("a skill not answered after its pause opens %v episodes", s["M5.episodes"])
	}
	if s["M6.late.n"] != 1 || s["M6.late.sum"] != 0 || s["M6.cases"] != 2 || s["M6.never"] != 1 {
		t.Errorf("M6: late %v of %v, %v cases, %v never", s["M6.late.sum"], s["M6.late.n"], s["M6.cases"], s["M6.never"])
	}
	if s["M4.all.declared"] != 1 || s["M4.all.false"] != 0 {
		t.Errorf("M4: %v declared, %v false", s["M4.all.declared"], s["M4.all.false"])
	}
}

// A number some resample cannot read is not read: one learner of five has a case.
func TestReadNeedsEveryResample(t *testing.T) {
	learners := []Sums{{"x": 1, "n": 1}, {}, {}, {}, {}}
	if v := Read("rare", nil, learners, []Number{{"mean", ratio("x", "n")}}); v[0].Read || v[0].Value != 1 {
		t.Errorf("a number of one learner in five: %+v", v[0])
	}
}

// M3 after the jump of G3: an error under 0.5 from the first attempt after it settles at once, and a
// run that ends before the jump has none.
func TestSettle(t *testing.T) {
	ms, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	p, _ := population.Named("G3")
	mods := []*pack.Module{ms[population.Olympiad]}
	l, err := learner.New(seed.Working, p.Name, 0, mods, p.Assumptions(0))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := l.Module(population.Olympiad)
	id := m.Pack.Graph.Skills[0].ID
	rec := func(answers int) *run.Record {
		r := &run.Record{Learner: l}
		for k := range answers {
			o := &run.Obs{Answers: k, Items: []run.ItemObs{{Pause: -1}}, After: []run.SkillObs{skill(id, 1, 0.1, 0, false, false)}}
			r.Steps = append(r.Steps, run.Step{Attempts: []learner.Attempt{{Items: []learner.ItemResult{{FirstRight: true}}}}, Obs: []*run.Obs{o}})
		}
		return r
	}
	s := Model(rec(m.JumpAt + 20))
	if s["M3.jumped"] != 1 || s["M3.settle.n"] != 1 || s["M3.settle.sum"] != 0 {
		t.Errorf("settled at once: jumped %v, settle %v of %v", s["M3.jumped"], s["M3.settle.sum"], s["M3.settle.n"])
	}
	if s := Model(rec(m.JumpAt)); s["M3.jumped"] != 0 || s["M3.unsettled"] != 0 {
		t.Errorf("a run that ends at the jump: jumped %v, unsettled %v", s["M3.jumped"], s["M3.unsettled"])
	}
}
