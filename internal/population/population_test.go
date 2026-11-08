package population

import (
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
