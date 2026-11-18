package criterion

import (
	"bytes"
	"errors"
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// withFairness adds to a file of results the cells of fairness by the grade on a graph, each learner
// with its trace: the learner n of a grade apart from the reference, which stays at 0, by 2.5 until
// the answer from(grade, n), and its attempts after 50 answers at the bands bands(grade, n). Their sums are of
// a learner who stays put, so that no other check could tell them apart.
func withFairness(f *bench.File, graph []string, from func(grade string, n int) int, bands func(grade string, n int) []int) *bench.File {
	for _, p := range population.Fairness {
		for n := range learners {
			tr := &measure.Trace{Bands: bands(p.Name, n)}
			for k := range measure.FairLength {
				x := 0.0
				if p.Name != population.Reference && k < from(p.Name, n) {
					x = 2.5
				}
				tr.Levels = append(tr.Levels, x)
			}
			sums := measure.Sums{"learners": 1, "M1.200.sq": 4, "M1.200.n": 1, "E1.corridor": float64(n % 2), "E1.items": 1}
			f.Lines = append(f.Lines, bench.Line{Population: p.Name, Graph: graph, N: n, Sums: sums, Oracle: measure.Sums{"learners": 1}, Fair: tr})
		}
	}
	f.Header.Learners = learners
	return f
}

// soon has every grade at the reference's estimate from the first answer, and even has the same
// attempts at every band.
func soon(string, int) int   { return 0 }
func even(string, int) []int { return []int{1, 1, 1} }

func fairBase(from func(string, int) int, bands func(string, int) []int) *bench.File {
	return withFairness(bare("engine", seed.Working, baseSide, false), olympiad, from, bands)
}

// evaluation evaluates a file of results with the decisive run lowered to the test's learners.
func evaluation(t *testing.T, f *bench.File) (*Fairness, error) {
	t.Helper()
	prior := decisive
	decisive = learners
	t.Cleanup(func() { decisive = prior })
	return Evaluate(Evaluation{Name: "fairness-1", Digest: "abc", Date: "2026-11-14", Version: "v0.1.18", Run: f, Pins: maps.Clone(pinned)})
}

// The base holds when every grade comes to the reference's estimate by 50 answers for at least 0.90
// of the learners at the low end and the bands differ by at most 0.05 at the high end; it does not
// when a grade comes late or its bands part, and the report says which and gives the paragraph for
// dpia 10.2 spec.
func TestEvaluate(t *testing.T) {
	f, err := evaluation(t, fairBase(func(_ string, n int) int { return 10 * (n % 6) }, even))
	if err != nil {
		t.Fatal(err)
	}
	if !f.Holds || len(f.Grades) != 2 || f.Grades[0].Cell != "D-grade1/"+olympiad[0] || f.Grades[1].Cell != "D-grade5/"+olympiad[0] {
		t.Fatalf("a base that holds: %+v", f.Grades)
	}
	for _, g := range f.Grades {
		if g.Answers.Value != 25 || g.Same.Value != 1 || g.Band.Value != 0 || g.Learners != learners || g.Reference != population.Reference+"/"+olympiad[0] {
			t.Errorf("%s: answers %v, same %v, band %v", g.Cell, g.Answers.Value, g.Same.Value, g.Band.Value)
		}
	}
	var b bytes.Buffer
	if err := f.Write(&b); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"# The fairness of the learner model by the grade: fairness-1", "R55", "D-grade1, D-grade3, D-grade5", "- sim v0.1.18", "SHA-256 abc", "| D-grade1/novatrail/olympiad-math | 60 | 25 [", "| yes |", "The base holds the criterion.", "## For dpia 10.2 spec", "On 2026-11-14", "Grade 5: a share of 1 [1; 1]", "github.com/altnovatrail/learner-model v1.2.0 and github.com/altnovatrail/engine v3.4.0"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("the report does not say %q:\n%s", want, b.String())
		}
	}
	late := func(grade string, n int) int {
		if grade == "D-grade5" && n%5 == 0 {
			return measure.FairLength
		}
		return 0
	}
	apart := func(grade string, n int) []int {
		if grade == "D-grade1" {
			return []int{2, 1, 0}
		}
		return even(grade, n)
	}
	f, err = evaluation(t, fairBase(late, apart))
	if err != nil {
		t.Fatal(err)
	}
	if f.Holds || !f.Grades[0].SameHolds || f.Grades[0].BandHolds || f.Grades[1].SameHolds || !f.Grades[1].BandHolds {
		t.Fatalf("a base that does not hold: %+v", f.Grades)
	}
	b.Reset()
	if err := f.Write(&b); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"| no: M11.band |", "| no: M11.same |", "The base does not hold the criterion."} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("the report does not say %q:\n%s", want, b.String())
		}
	}
	both := Grade{}
	if holds(both) != "no: M11.same and M11.band" || showInterval(measure.Value{}) != "unread" {
		t.Errorf("an unread grade: %q", holds(both))
	}
}

// An evaluation reads only a whole decisive run of the base under the engine, on the working learners,
// the pinned packages and the libraries of the gate file, with no breach and every cell of fairness on
// a published graph; it needs its name, date and digest.
func TestEvaluateRefuses(t *testing.T) {
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	change := func(f func(*bench.File)) *bench.File {
		g := fairBase(soon, even)
		f(g)
		return g
	}
	for name, f := range map[string]*bench.File{
		"the ceiling":       change(func(g *bench.File) { g.Header.Rule = "ceiling" }),
		"held-out learners": change(func(g *bench.File) { g.Header.Set = held }),
		"a part":            change(func(g *bench.File) { g.Header.Part = bench.Part{K: 1, N: 2} }),
		"packages in place": change(func(g *bench.File) { g.Header.Packages = []published.Package{{Module: olympiad[0]}} }),
		"other libraries":   change(func(g *bench.File) { g.Header.Libraries = map[string]string{bench.Libraries[0]: "v9.9.9"} }),
		"a breach":          withFairness(bare("engine", seed.Working, baseSide, true), olympiad, soon, even),
		"no cell":           bare("engine", seed.Working, baseSide, false),
		"a synthetic graph": withFairness(bare("engine", seed.Working, baseSide, false), []string{synthetic.Sheet}, soon, even),
		"no trace":          change(func(g *bench.File) { g.Lines[len(g.Lines)-1].Fair = nil }),
		"no invariants":     change(func(g *bench.File) { g.Header.Invariants = false }),
	} {
		if _, err := evaluation(t, f); !errors.Is(err, ErrCriterion) {
			t.Errorf("%s: %v", name, err)
		}
	}
	prior := decisive
	decisive = learners + 1
	t.Cleanup(func() { decisive = prior })
	for name, e := range map[string]Evaluation{
		"a rough look": {Name: "x", Digest: "d", Date: "2026-11-14", Run: fairBase(soon, even), Pins: maps.Clone(pinned)},
		"no date":      {Name: "x", Digest: "d", Run: fairBase(soon, even), Pins: maps.Clone(pinned)},
		"no pins":      {Name: "x", Digest: "d", Date: "2026-11-14", Run: fairBase(soon, even)},
	} {
		if _, err := Evaluate(e); !errors.Is(err, ErrCriterion) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// The cells of fairness change nothing a judgement reads but the constraint of fairness: the report
// of a candidate against a base with them is that without them but for the lines of fair, and the
// paths read off a base with them and a ceiling without them are those of a base without them (R04,
// R55).
func TestFairnessLeftOut(t *testing.T) {
	base, ceiling := bare("engine", seed.Working, baseSide, false), file("ceiling", seed.Working, ceilSide, false)
	p := usePaths(t, base, ceiling)
	with, err := RecordPaths(fairBase(soon, even), ceiling, "base", "ceiling")
	if err != nil || !reflect.DeepEqual(with, p) {
		t.Fatalf("the paths of a base with the cells of fairness: %v", err)
	}
	// The cells of fairness first in the file: each cell still reads its own oracle.
	first := fairBase(soon, even)
	slices.SortStableFunc(first.Lines, func(a, b bench.Line) int {
		pa, _ := population.Named(a.Population)
		pb, _ := population.Named(b.Population)
		return b2i(pb.Fair()) - b2i(pa.Fair())
	})
	if got, err := RecordPaths(first, ceiling, "base", "ceiling"); err != nil || !reflect.DeepEqual(got, p) {
		t.Fatalf("the paths of a base with the cells of fairness first: %v", err)
	}
	// Both candidates fail fair, one by its absent cells and the other by its numbers, and the rest
	// of the report is the same.
	report := func(base, cand *bench.File) (string, int) {
		t.Helper()
		in := input(Model, Candidate{Name: "better", File: cand, Group: "still"})
		in.Base = base
		r, err := Judge(in)
		if err != nil {
			t.Fatal(err)
		}
		var b bytes.Buffer
		if err := r.Write(&b); err != nil {
			t.Fatal(err)
		}
		lines := strings.SplitAfter(b.String(), "\n")
		n := len(lines)
		lines = slices.DeleteFunc(lines, func(l string) bool { return strings.HasPrefix(l, "- fails fair ") })
		return strings.Join(lines, ""), n - len(lines)
	}
	plain, n := report(bare("engine", seed.Working, baseSide, false), bare("engine", seed.Working, betterSide, false))
	late := func(string, int) int { return measure.FairLength }
	fair, m := report(fairBase(soon, even), withFairness(bare("engine", seed.Working, betterSide, false), olympiad, late, even))
	if fair != plain || n == 0 || m == 0 {
		t.Errorf("the cells of fairness change the report beyond fair (%d and %d lines of it):\n%s\nagainst\n%s", n, m, fair, plain)
	}
	// A candidate that holds fair is chosen alike whether the other sums of its cells of fairness are
	// the base's or far worse: no check, score or choice reads them.
	worse := withFairness(bare("engine", seed.Working, betterSide, false), olympiad, soon, even)
	for i, l := range worse.Lines {
		if p, _ := population.Named(l.Population); p.Fair() {
			worse.Lines[i].Sums = measure.Sums{"learners": 1, "M1.200.sq": 100, "M1.200.n": 1, "E1.corridor": 0, "E1.items": 1}
		}
	}
	same, _ := report(fairBase(soon, even), withFairness(bare("engine", seed.Working, betterSide, false), olympiad, soon, even))
	other, k := report(fairBase(soon, even), worse)
	if same != other || k != 0 || !strings.Contains(same, "chosen: better") {
		t.Errorf("the other sums of the cells of fairness change the report:\n%s\nagainst\n%s", other, same)
	}
}

// The constraint fair holds a candidate as close to the reference as the base on every grade, and
// fails the grade and the number where it is not: in the criteria of the model and of the engine, in a
// pair once, and in the confirmation; runs without the cells fail it, and a line of fairness with
// no trace is refused (R55).
func TestFair(t *testing.T) {
	usual(t)
	half := func(grade string, n int) int {
		if grade == "D-grade5" && n%2 == 0 {
			return measure.FairLength
		}
		return 0
	}
	apart := func(grade string, n int) []int {
		if grade == "D-grade1" {
			return []int{2, 1, 0}
		}
		return even(grade, n)
	}
	better := func(from func(string, int) int, bands func(string, int) []int) *bench.File {
		return withFairness(bare("engine", seed.Working, betterSide, false), olympiad, from, bands)
	}
	fairs := func(v Verdict) (all int, fails []string) {
		for _, k := range v.Checks {
			if k.Kind == "fair" {
				all++
				if !k.Met {
					fails = append(fails, k.Cell+" "+k.Number)
				}
			}
		}
		return
	}
	g1, g5 := "D-grade1/"+olympiad[0], "D-grade5/"+olympiad[0]
	for _, crit := range []string{Model, Engine, Pair} {
		lib := ""
		if crit == Pair {
			lib = Model
		}
		r, err := Judge(input(crit,
			Candidate{Name: "fair", File: file("engine", seed.Working, betterSide, false), Library: lib},
			Candidate{Name: "late", File: better(half, even), Library: lib},
			Candidate{Name: "apart", File: better(soon, apart), Library: lib},
		))
		if err != nil {
			t.Fatal(err)
		}
		for name, want := range map[string][]string{
			"fair":  nil,
			"late":  {g5 + " M11.answers", g5 + " M11.same"},
			"apart": {g1 + " M11.band"},
		} {
			v := verdict(t, r, name)
			all, fails := fairs(v)
			slices.Sort(fails)
			if all != 6 || !slices.Equal(fails, want) || v.Holds != (want == nil) {
				t.Errorf("%s, %s: %d checks of fair, fails %v, holds %v; want %v", crit, name, all, fails, v.Holds, want)
			}
		}
		if r.Chosen != "fair" {
			t.Errorf("%s: chosen %q", crit, r.Chosen)
		}
		var b bytes.Buffer
		if err := r.Write(&b); err != nil {
			t.Fatal(err)
		}
		if want := "- fails fair " + g1 + " M11.band: worse by at most 0.3333, the tolerance 0.01"; !strings.Contains(b.String(), want) {
			t.Errorf("%s: the report does not say %q:\n%s", crit, want, b.String())
		}
	}
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	in := input(Model, Candidate{Name: "fair", File: file("engine", seed.Working, betterSide, false)})
	in.Confirm = &Confirm{Base: file("engine", held, baseSide, false), Candidate: withFairness(bare("engine", held, betterSide, false), olympiad, half, even)}
	in.Judgement = "model-2"
	r, err := Judge(in)
	if err != nil || r.Confirmed == nil || r.Confirmed.Holds || !strings.Contains(r.Outcome, "not confirmed") {
		t.Errorf("a confirmation less fair than the base: %v, %+v", err, r)
	}
	in = input(Model, Candidate{Name: "bare", File: bare("engine", seed.Working, betterSide, false)})
	in.Base = bare("engine", seed.Working, baseSide, false)
	if r, err = Judge(in); err != nil {
		t.Fatal(err)
	}
	all, fails := fairs(verdict(t, r, "bare"))
	slices.Sort(fails)
	if want := []string{g1 + " M11.answers", g1 + " M11.band", g1 + " M11.same", g5 + " M11.answers", g5 + " M11.band", g5 + " M11.same"}; all != 6 || !slices.Equal(fails, want) || verdict(t, r, "bare").Holds {
		t.Errorf("runs without the cells of fairness: %d checks, fails %v", all, fails)
	}
	// A side without its trace, or without the cells the other has, is refused, and the error names
	// the side.
	for name, c := range map[string]struct {
		change func(*Input)
		want   string
	}{
		"the base with no trace":                      {func(in *Input) { in.Base.Lines[len(in.Base.Lines)-1].Fair = nil }, "the base: bench: learner D-grade5"},
		"the candidate with no trace":                 {func(in *Input) { in.Candidates[0].File.Lines[len(in.Base.Lines)-1].Fair = nil }, "candidate fair: bench: learner D-grade5"},
		"the confirmation's base with no trace":       {func(in *Input) { in.Confirm.Base.Lines[len(in.Base.Lines)-1].Fair = nil }, "the confirmation's base: bench: learner D-grade5"},
		"a candidate without the cells":               {func(in *Input) { in.Candidates[0].File = bare("engine", seed.Working, betterSide, false) }, "candidate fair has the cells of fairness by the grade [], and the base [" + g1 + ", " + g5 + "]"},
		"a candidate with the cells the base has not": {func(in *Input) { in.Base = bare("engine", seed.Working, baseSide, false) }, "candidate fair has the cells of fairness by the grade [" + g1 + ", " + g5 + "], and the base []"},
	} {
		in := input(Model, Candidate{Name: "fair", File: file("engine", seed.Working, betterSide, false)})
		in.Confirm = &Confirm{Base: file("engine", held, baseSide, false), Candidate: file("engine", held, betterSide, false)}
		in.Judgement = "model-2"
		c.change(&in)
		if _, err := Judge(in); !errors.Is(err, ErrCriterion) || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v", name, err)
		}
	}
}
