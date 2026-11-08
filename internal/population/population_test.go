package population

import (
	"maps"
	"reflect"
	"slices"
	"strconv"
	"testing"

	"github.com/altnovatrail/sim/internal/invariant"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// changed gives the names of the fields of a in which it differs from the base population's.
func changed(a learner.Assumptions) []string {
	base := reflect.ValueOf(learner.Base())
	v := reflect.ValueOf(a)
	var out []string
	for i := range v.NumField() {
		if !reflect.DeepEqual(v.Field(i).Interface(), base.Field(i).Interface()) {
			out = append(out, v.Type().Field(i).Name)
		}
	}
	return out
}

// The nineteen populations of standalone, each the base population with one assumption changed —
// the fields of that one assumption and no other — on the graphs docs/architecture.md, section
// 3.2, gives them.
func TestStandalone(t *testing.T) {
	want := []struct {
		name   string
		fields [][]string // by side
		graphs [][]string
	}{
		{"G0", [][]string{nil}, maths},
		{"G1", [][]string{{"LevelShift"}, {"LevelShift"}}, maths},
		{"G2", [][]string{{"LearnLevel", "LearnSkill", "LearnReturn"}}, maths},
		{"G2-half", [][]string{{"LearnLevel", "LearnSkill", "LearnReturn"}}, maths},
		{"G2-fading", [][]string{{"LearnLevel", "LearnSkill", "LearnReturn", "LearnFade"}}, maths},
		{"G3", [][]string{{"JumpBy", "JumpFrom", "JumpTo"}}, maths},
		{"G3-drop", [][]string{{"JumpBy", "JumpFrom", "JumpTo"}}, maths},
		{"G4-band", [][]string{{"BandBias"}}, olympiad},
		{"G5", [][]string{{"SkillSpread", "SkillHalves"}, {"SkillSpread", "SkillHalves"}}, maths},
		{"G6", [][]string{{"Slope"}, {"Slope"}}, maths},
		{"G7", [][]string{{"ChoiceGuess"}, {"ChoiceGuess"}}, grade3},
		{"G8", [][]string{{"TechniqueChance", "TechniqueLift"}}, maths},
		{"G0-exact", [][]string{{"Miss"}}, maths},
		{"G0-miss0.25", [][]string{{"Miss"}}, maths},
		{"G0-miss1", [][]string{{"Miss"}}, maths},
		{"G0-start0.5", [][]string{{"StartSpread"}}, maths},
		{"G0-start2", [][]string{{"StartSpread"}}, maths},
		{"G0-skills0.3", [][]string{{"SkillSpread"}}, maths},
		{"G0-skills1", [][]string{{"SkillSpread"}}, maths},
	}
	if len(Standalone) != len(want) {
		t.Fatalf("%d populations of standalone, want %d", len(Standalone), len(want))
	}
	for i, w := range want {
		p := Standalone[i]
		if p.Name != w.name || !reflect.DeepEqual(p.Graphs, w.graphs) || p.Limit != run.Maths {
			t.Errorf("population %d: %s on %v, limit %+v; want %s on %v", i, p.Name, p.Graphs, p.Limit, w.name, w.graphs)
		}
		for side, fields := range w.fields {
			if got := changed(p.Assumptions(side)); !slices.Equal(got, fields) {
				t.Errorf("%s, side %d, changes %v; want %v", p.Name, side, got, fields)
			}
			if p.Side(side) != side {
				t.Errorf("%s: learner %d takes side %d", p.Name, side, p.Side(side))
			}
		}
		if len(w.fields) == 2 && reflect.DeepEqual(p.Assumptions(0), p.Assumptions(1)) {
			t.Errorf("%s: its two sides are one", p.Name)
		}
		if got, ok := Named(w.name); !ok || got.Name != w.name {
			t.Errorf("Named(%q) = %v, %v", w.name, got.Name, ok)
		}
	}
	if _, ok := Named("G99"); ok {
		t.Error("a population that is not")
	}
}

// Every population runs on each of its graphs, and the invariants hold on its learners.
func TestPopulationsRun(t *testing.T) {
	ms, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	syn, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	ms = maps.Clone(ms)
	maps.Copy(ms, syn)
	techniques := 0
	for _, p := range All() {
		for _, g := range p.Graphs {
			var mods []*pack.Module
			for _, id := range g {
				mods = append(mods, ms[id])
			}
			for n := range 4 {
				l, err := learner.New(seed.Working, p.Name, n, mods, p.Assumptions(n))
				if err != nil {
					t.Fatal(err)
				}
				r, err := run.New(l, mods, p.Limit)
				if err != nil {
					t.Fatal(err)
				}
				rec, err := r.All()
				if err != nil {
					t.Fatalf("%s on %v: %v", p.Name, g, err)
				}
				for _, b := range invariant.Check(rec) {
					t.Error(b)
				}
				for _, s := range rec.Steps {
					for i, a := range s.Attempts {
						if s.Evidence[i] != nil && s.Evidence[i].TechniqueOnRequest != a.Technique {
							t.Fatalf("%s: the technique of %s is not in its evidence", rec.Unit, a.Task)
						}
						if a.Technique {
							if p.Name != "G8" {
								t.Fatalf("%s takes the technique", rec.Unit)
							}
							techniques++
						}
					}
				}
			}
		}
	}
	if techniques == 0 {
		t.Error("G8 never takes the technique")
	}
}

// No learner of any population is drawn alike in two sets (docs/architecture.md, section 3.8):
// the first draw of the params of the first 300 learners of every population, in each recorded set.
func TestSetsApart(t *testing.T) {
	sets := []seed.Set{seed.Working}
	for _, c := range []string{seed.CriterionModel, seed.CriterionEngine} {
		s, err := seed.HeldOut(c, 1)
		if err != nil {
			t.Fatal(err)
		}
		sets = append(sets, s)
	}
	seen := map[uint64]string{}
	for _, s := range sets {
		for _, p := range All() {
			for n := range 300 {
				st, err := s.Place(seed.Learner(p.Name, n), seed.Params)
				if err != nil {
					t.Fatal(err)
				}
				first := st.Integer()
				where := s.Experiment + " " + p.Name + "/" + strconv.Itoa(n)
				if prior, ok := seen[first]; ok {
					t.Fatalf("%s draws first what %s does", where, prior)
				}
				seen[first] = where
			}
		}
	}
}

// The populations of S-, F-, W-, A-, L- and E-: each the base population with the fields of its one assumption,
// on the graphs and to the limit docs/architecture.md, section 3.2, gives it. F- share the calendar
// with a break, and each but F-model changes one more thing on it.
func TestOwn(t *testing.T) {
	cal := []string{"SessionChance", "BreakStartFrom", "BreakStartTo", "BreakMin", "BreakMax", "BreakEndBy"}
	with := func(fields ...string) []string { return append(slices.Clone(cal), fields...) }
	want := []struct {
		name   string
		fields []string
		graphs [][]string
		limit  run.Limit
	}{
		{"S-scale1.5", []string{"BandStep"}, olympiad, run.Maths},
		{"S-scale3.5", []string{"BandStep"}, olympiad, run.Maths},
		{"S-predicate", []string{"PredicateGuess"}, [][]string{{synthetic.Sheet}}, run.Maths},
		{"S-adult", []string{"Adult", "AdultFrom", "AdultTo"}, olympiad, run.Maths},
		{"S-unreadable", []string{"UnreadableEvery"}, olympiad, run.Maths},
		{"F-model", with(), maths, run.Forgetting},
		{"F-none", with("Forgets"), maths, run.Forgetting},
		{"F-exp", with("ExponentialLaw", "LawBase"), maths, run.Forgetting},
		{"F-C1", with("LawC"), maths, run.Forgetting},
		{"F-growth5", with("Growth"), maths, run.Forgetting},
		{"F-break", with("BreakDrop"), maths, run.Forgetting},
		{"W-shock", []string{"AttemptShock"}, grade3, run.Maths},
		{"W-error", []string{"ErrorChance"}, grade3, run.Maths},
		{"W-badday", []string{"BadDayChance", "BadDayDrop"}, grade3, run.Maths},
		{"W-repeat", []string{"RepeatChance", "RepeatDays"}, maths, run.Maths},
		{"A-linked", []string{"LevelCorrelation"}, anchored, run.Maths},
		{"A-unpractised", []string{"TeacherBoost"}, anchored, run.Maths},
		{"A-fading", []string{"TeacherFades"}, anchored, run.Maths},
		{"L-abandon", []string{"AbandonChance"}, olympiad, run.Maths},
		{"L-explain", []string{"ExplainLift"}, olympiad, run.Maths},
		{"E-base", nil, essays, run.Essays},
		{"E-leave", []string{"LeaveChance"}, essays, run.Essays},
	}
	if len(Own) != len(want) {
		t.Fatalf("%d own populations, want %d", len(Own), len(want))
	}
	order := func(fields []string) []string {
		var out []string
		v := reflect.TypeOf(learner.Assumptions{})
		for i := range v.NumField() {
			if slices.Contains(fields, v.Field(i).Name) {
				out = append(out, v.Field(i).Name)
			}
		}
		return out
	}
	for i, w := range want {
		p := Own[i]
		if p.Name != w.name || !reflect.DeepEqual(p.Graphs, w.graphs) || p.Limit != w.limit {
			t.Errorf("population %d: %s on %v, limit %+v", i, p.Name, p.Graphs, p.Limit)
		}
		if got := changed(p.Assumptions(0)); !slices.Equal(got, order(w.fields)) {
			t.Errorf("%s changes %v; want %v", p.Name, got, order(w.fields))
		}
	}
	if len(All()) != len(Standalone)+len(Own) {
		t.Errorf("All has %d populations", len(All()))
	}
}
