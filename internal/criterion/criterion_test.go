package criterion

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

const learners = 60

var olympiad = []string{"novatrail/olympiad-math"}

// noise gives the k-th learner its own small deviation, the same on every side.
func noise(k int) float64 { return float64(0.01 * float64((k*7)%11-5)) }

// side is how one side's learners go: the rms of M1 after 200 answers, the share of items in the
// corridor, the lag M3, the share that never settles, and the oracle's rms.
type side struct {
	rms, corridor, lag, unsettled, oracleRMS float64
}

// file builds a file of results of a rule on a set: learners of each population on the olympiad
// maths, and of A-fading on the synthetic graph, their sums as a side says.
func file(rule string, set seed.Set, s side, breach bool) *bench.File {
	f := &bench.File{Header: bench.Header{Set: set, Rule: rule, Invariants: rule == "engine"}}
	for _, p := range []string{"G0", "G0-skills1", "G2-half", "G3", "A-fading"} {
		graph := olympiad
		if p == "A-fading" {
			graph = []string{synthetic.Sheet, synthetic.Workbook}
		}
		for k := range learners {
			e := noise(k)
			r := s.rms + e
			sums := measure.Sums{
				"learners": 1, "M1.200.sq": float64(r * r), "M1.200.n": 1,
				"E1.corridor": float64(100 * (s.corridor + e)), "E1.items": 100,
				"M3.lag.sum": s.lag + e, "M3.lag.n": 1,
				"M3.unsettled": float64(b2i(k%100 < int(float64(100*s.unsettled)))), "M3.jumped": 1,
			}
			o := float64(s.oracleRMS + e)
			oracle := measure.Sums{"learners": 1, "M1.200.sq": float64(o * o), "M1.200.n": 1, "M3.lag.sum": 0, "M3.lag.n": 1}
			l := bench.Line{Population: p, Graph: graph, N: k, Sums: sums, Oracle: oracle}
			if breach && k == 3 {
				l.Breaches = []string{"I4: a solved task again"}
			}
			f.Lines = append(f.Lines, l)
		}
	}
	return f
}

var (
	baseSide   = side{rms: 0.5, corridor: 0.1, lag: 0.5, unsettled: 0.4}
	ceilSide   = side{rms: 0.5, corridor: 0.5, lag: 0.5, unsettled: 0.4}
	betterSide = side{rms: 0.3, corridor: 0.3, lag: 0.1, unsettled: 0.05}
	worseSide  = side{rms: 0.6, corridor: 0.1, lag: 0.5, unsettled: 0.4}
)

func input(criterion string, cands ...Candidate) Input {
	return Input{Criterion: criterion, Base: file("engine", seed.Working, baseSide, false), Candidates: cands}
}

// usePaths fixes the paths read off a base and a ceiling for the rest of a test.
func usePaths(t *testing.T, base, ceiling *bench.File) *Paths {
	t.Helper()
	prior, priorDecisive, priorCovered := fixed, decisive, covered
	decisive, covered = learners, []string{"G0", "G0-skills1", "G2-half", "G3"}
	p, err := RecordPaths(base, ceiling, "base", "ceiling")
	if err != nil {
		t.Fatal(err)
	}
	fixed = func() (*Paths, error) { return p, nil }
	t.Cleanup(func() { fixed, decisive, covered = prior, priorDecisive, priorCovered })
	return p
}

// usual fixes the paths of the usual base and ceiling.
func usual(t *testing.T) {
	t.Helper()
	usePaths(t, file("engine", seed.Working, baseSide, false), file("ceiling", seed.Working, ceilSide, false))
}

func verdict(t *testing.T, r *Report, name string) Verdict {
	t.Helper()
	i := slices.IndexFunc(r.Verdicts, func(v Verdict) bool { return v.Candidate == name })
	if i < 0 {
		t.Fatalf("no verdict of %s", name)
	}
	return r.Verdicts[i]
}

func failed(v Verdict) []string {
	var out []string
	for _, k := range v.Checks {
		if !k.Met && !k.Skipped && !k.NoGraph {
			out = append(out, k.Kind+" "+k.Cell+" "+k.Number)
		}
	}
	return out
}

// A candidate equal to the base holds "no worse" and is not better; one worse beyond the tolerance
// fails "no worse"; one better holds it, scores above zero and is chosen; the base stays without it.
func TestNoWorseAndBetter(t *testing.T) {
	usual(t)
	r, err := Judge(input(Model,
		Candidate{Name: "equal", File: file("engine", seed.Working, baseSide, false)},
		Candidate{Name: "worse", File: file("engine", seed.Working, worseSide, false)},
		Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false), Group: "still"},
	))
	if err != nil {
		t.Fatal(err)
	}
	equal, worse, better := verdict(t, r, "equal"), verdict(t, r, "worse"), verdict(t, r, "better")
	if !equal.Holds || equal.Better || equal.Score.Value != 0 || equal.Score.Low != 0 || equal.Score.High != 0 {
		t.Errorf("equal: holds %v, better %v, score %+v; fails %v", equal.Holds, equal.Better, equal.Score.Value, failed(equal))
	}
	if worse.Holds || !slices.ContainsFunc(failed(worse), func(f string) bool { return strings.HasSuffix(f, "M1.rms.200") }) {
		t.Errorf("worse: holds %v, fails %v", worse.Holds, failed(worse))
	}
	if !better.Holds || !better.Better || better.Score.Low <= 0 || better.Group == nil || better.Group.Low <= 0 || equal.Group != nil {
		t.Errorf("better: holds %v, better %v, score %+v, group %+v; fails %v", better.Holds, better.Better, better.Score.Value, better.Group, failed(better))
	}
	if r.Chosen != "better" || r.Confirmed != nil || !strings.Contains(r.Outcome, "awaits") {
		t.Errorf("chosen %q, outcome %q", r.Chosen, r.Outcome)
	}
	if r.Decisive || r.Learners != learners || !slices.Equal(r.Graphs, []string{"novatrail/olympiad-math"}) || len(r.NoGraph) != 1 {
		t.Errorf("the report: decisive %v, %d learners, graphs %v, no graph %v", r.Decisive, r.Learners, r.Graphs, r.NoGraph)
	}
	r, err = Judge(input(Model, Candidate{Name: "equal", File: file("engine", seed.Working, baseSide, false)}))
	if err != nil || r.Chosen != "" || !strings.Contains(r.Outcome, "base stays") {
		t.Errorf("no candidate better: %v, %+v", err, r)
	}
}

// A number whose path is zero or less leaves the score with its weight, and the score names it.
func TestPathLeaves(t *testing.T) {
	in := input(Model, Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false)})
	s := baseSide
	s.oracleRMS = 0.9 // the oracle worse than the base: no path for M1
	in.Base = file("engine", seed.Working, s, false)
	usePaths(t, in.Base, file("ceiling", seed.Working, ceilSide, false))
	r, err := Judge(in)
	if err != nil {
		t.Fatal(err)
	}
	v := verdict(t, r, "better")
	if !slices.Contains(v.Score.Left, "G0/novatrail/olympiad-math M1.rms.200") || !v.Score.Read {
		t.Errorf("left %v, read %v", v.Score.Left, v.Score.Read)
	}
}

// Of candidates equal to the best, the simplest is chosen: the fewest constants, then no new field,
// then the nearest to the base.
func TestChoice(t *testing.T) {
	usual(t)
	r, err := Judge(input(Model,
		Candidate{Name: "b", File: file("engine", seed.Working, betterSide, false), Constants: 2},
		Candidate{Name: "c", File: file("engine", seed.Working, betterSide, false), Constants: 1, NewField: true},
		Candidate{Name: "a", File: file("engine", seed.Working, betterSide, false), Constants: 1, Distance: 2},
		Candidate{Name: "d", File: file("engine", seed.Working, betterSide, false), Constants: 1, Distance: 1},
	))
	if err != nil {
		t.Fatal(err)
	}
	if r.Chosen != "d" || len(r.Equals) != 4 {
		t.Errorf("chosen %q among %v", r.Chosen, r.Equals)
	}
}

// The chosen candidate is confirmed on held-out learners when it holds and is better there too;
// otherwise the base stays.
func TestConfirm(t *testing.T) {
	usual(t)
	held, err := seed.HeldOut(seed.CriterionModel, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		s    side
		want string
	}{{betterSide, "confirmed"}, {baseSide, "not confirmed"}} {
		in := input(Model, Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false)})
		in.Confirm = &Confirm{Base: file("engine", held, baseSide, false), Candidate: file("engine", held, c.s, false)}
		r, err := Judge(in)
		if err != nil {
			t.Fatal(err)
		}
		if r.Confirmed == nil || !strings.Contains(r.Outcome, c.want) {
			t.Errorf("outcome %q, want %q", r.Outcome, c.want)
		}
	}
}

// A breach of an invariant fails the candidate; the goals of the model hold or fail by the share of
// the path and the share that never settles.
func TestInvariantsAndGoals(t *testing.T) {
	usual(t)
	r, err := Judge(input(Model,
		Candidate{Name: "breaks", File: file("engine", seed.Working, betterSide, true)},
		Candidate{Name: "lag", File: file("engine", seed.Working, betterSide, false), Goal: "lag"},
		Candidate{Name: "corridor", File: file("engine", seed.Working, betterSide, false), Goal: "corridor"},
		Candidate{Name: "settle", File: file("engine", seed.Working, betterSide, false), Goal: "settle"},
		Candidate{Name: "unsettled", File: file("engine", seed.Working, baseSide, false), Goal: "settle"},
	))
	if err != nil {
		t.Fatal(err)
	}
	if v := verdict(t, r, "breaks"); v.Holds || v.Checks[0].Met || !strings.Contains(v.Checks[0].Detail, "I4") {
		t.Errorf("a breach: holds %v, %+v", v.Holds, v.Checks[0])
	}
	for _, name := range []string{"lag", "corridor", "settle"} {
		v := verdict(t, r, name)
		goals := slices.DeleteFunc(slices.Clone(v.Checks), func(k Check) bool { return k.Kind != "goal" })
		if len(goals) == 0 || !v.Holds {
			t.Errorf("%s: goals %+v, fails %v", name, goals, failed(v))
		}
	}
	if v := verdict(t, r, "unsettled"); v.Holds {
		t.Errorf("40 %% never settled holds the goal of 25 %%")
	}
	// A file that did not record the invariants, and runs with no learner of the populations a floor
	// or a goal needs, fail those checks.
	unchecked := file("engine", seed.Working, betterSide, false)
	unchecked.Header.Invariants = false
	r, err = Judge(input(Model, Candidate{Name: "unchecked", File: unchecked}))
	if err != nil {
		t.Fatal(err)
	}
	if v := verdict(t, r, "unchecked"); v.Holds || !strings.Contains(v.Checks[0].Detail, "did not record") {
		t.Errorf("unrecorded invariants: holds %v, %+v", v.Holds, v.Checks[0])
	}
	only := func(f *bench.File) *bench.File {
		g := *f
		g.Lines = slices.DeleteFunc(slices.Clone(f.Lines), func(l bench.Line) bool { return l.Population != "G0" })
		return &g
	}
	in := input(Model, Candidate{Name: "lag", File: only(file("engine", seed.Working, betterSide, false)), Goal: "lag"},
		Candidate{Name: "settle", File: only(file("engine", seed.Working, betterSide, false)), Goal: "settle"})
	in.Base = only(in.Base)
	r, err = Judge(in)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"lag", "settle"} {
		if v := verdict(t, r, name); v.Holds || !slices.ContainsFunc(v.Checks, func(k Check) bool { return strings.Contains(k.Detail, "no learner of") }) {
			t.Errorf("%s: a goal never measured holds", name)
		}
	}
	noG0 := func(f *bench.File) *bench.File {
		g := *f
		g.Lines = slices.DeleteFunc(slices.Clone(f.Lines), func(l bench.Line) bool { return l.Population == "G0" })
		return &g
	}
	in = input(Model, Candidate{Name: "floorless", File: noG0(file("engine", seed.Working, betterSide, false))})
	in.Base = noG0(in.Base)
	if r, err = Judge(in); err != nil || verdict(t, r, "floorless").Holds {
		t.Errorf("runs with no G0 hold the floor: %v", err)
	}
}

// The criterion of the engine scores by its groups; a pair holds the constraints of both and is
// better by its library's score, the other's interval not wholly below zero.
func TestEngineAndPair(t *testing.T) {
	usual(t)
	r, err := Judge(input(Engine, Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false), Group: "corridor"}))
	if err != nil {
		t.Fatal(err)
	}
	if v := verdict(t, r, "better"); !v.Better || !v.Holds || r.Chosen != "better" {
		t.Errorf("engine: better %v, holds %v, fails %v", v.Better, v.Holds, failed(v))
	}
	r, err = Judge(input(Pair,
		Candidate{Name: "pair", File: file("engine", seed.Working, betterSide, false), Library: Model, Group: "learns"},
		Candidate{Name: "worse", File: file("engine", seed.Working, worseSide, false), Library: Engine},
	))
	if err != nil {
		t.Fatal(err)
	}
	if v := verdict(t, r, "pair"); !v.Better || !v.Holds || v.Other == nil || !v.Other.Read {
		t.Errorf("pair: better %v, holds %v, other %+v, fails %v", v.Better, v.Holds, v.Other, failed(v))
	}
	v := verdict(t, r, "worse")
	if v.Better || v.Holds {
		t.Errorf("a worse pair: better %v, holds %v", v.Better, v.Holds)
	}
	seen := map[string]bool{}
	for _, k := range v.Checks {
		if key := k.Kind + " " + k.Cell + " " + k.Number; k.Kind == "no worse" && seen[key] {
			t.Errorf("a pair checks %s twice", key)
		} else {
			seen[key] = true
		}
	}
}

// What Judge refuses.
func TestRefuses(t *testing.T) {
	usual(t)
	good := file("engine", seed.Working, betterSide, false)
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	other := file("engine", held, betterSide, false)
	for name, in := range map[string]Input{
		"no criterion":      input("taste"),
		"no base":           {Criterion: Model},
		"ceiling as base":   {Criterion: Model, Base: file("ceiling", seed.Working, baseSide, false)},
		"candidate ceiling": input(Model, Candidate{Name: "x", File: file("ceiling", seed.Working, betterSide, false)}),
		"group":             input(Model, Candidate{Name: "x", File: good, Group: "corridor"}),
		"goal":              input(Model, Candidate{Name: "x", File: good, Goal: "speed"}),
		"pair library":      input(Pair, Candidate{Name: "x", File: good}),
		"engine goal":       input(Engine, Candidate{Name: "x", File: good, Goal: "corridor"}),
		"pair engine goal":  input(Pair, Candidate{Name: "x", File: good, Library: Engine, Goal: "lag"}),
		"other learners":    input(Model, Candidate{Name: "x", File: other}),
	} {
		if _, err := Judge(in); !errors.Is(err, ErrCriterion) {
			t.Errorf("%s: %v", name, err)
		}
	}
	in := input(Model, Candidate{Name: "x", File: good})
	in.Confirm = &Confirm{Base: good}
	if _, err := Judge(in); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation without its runs: %v", err)
	}
	in.Confirm = &Confirm{Base: good, Candidate: good}
	if _, err := Judge(in); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation on the working learners: %v", err)
	}
	in.Confirm = &Confirm{Base: other, Candidate: file("ceiling", held, betterSide, false)}
	if _, err := Judge(in); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation with its candidate under the ceiling's rule: %v", err)
	}
	same := input(Model, Candidate{Name: "x", File: other})
	same.Base = file("engine", held, baseSide, false)
	same.Confirm = &Confirm{Base: same.Base, Candidate: other}
	if _, err := Judge(same); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation on the learners of the choice: %v", err)
	}
	conf := input(Model, Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false)})
	conf.Confirm = &Confirm{Base: file("engine", held, baseSide, false), Candidate: file("engine", seed.Working, betterSide, false)}
	if _, err := Judge(conf); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation whose candidate ran on other learners: %v", err)
	}
	// With no path fixed, or none to read, there is nothing to score by.
	fixed = func() (*Paths, error) { return &Paths{}, nil }
	if _, err := Judge(input(Model, Candidate{Name: "x", File: good})); !errors.Is(err, ErrCriterion) || !strings.Contains(err.Error(), "no path is fixed") {
		t.Errorf("no path fixed: %v", err)
	}
	fixed = func() (*Paths, error) { return nil, ErrCriterion }
	if _, err := Judge(input(Model, Candidate{Name: "x", File: good})); !errors.Is(err, ErrCriterion) {
		t.Errorf("paths that do not read: %v", err)
	}
}

// The paths are read off the base and the ceiling: the base's value and the oracle's for the model's
// numbers, the ceiling of the engine's for the engine's, zero for those whose path goes to
// perfection, each turned by its direction; no path on a synthetic graph; and the runs they are
// read on are decisive runs of the base and the ceiling on the working learners.
func TestRecordPaths(t *testing.T) {
	base, ceiling := file("engine", seed.Working, baseSide, false), file("ceiling", seed.Working, ceilSide, false)
	p := usePaths(t, base, ceiling)
	if p.Learners != learners || p.Set != seed.Working || p.Base != "base" || p.Ceiling != "ceiling" {
		t.Errorf("the paths name %d learners, %+v, %s and %s", p.Learners, p.Set, p.Base, p.Ceiling)
	}
	at := func(cell, number string) Path {
		t.Helper()
		x, ok := p.on(cell)[number]
		if !ok {
			t.Fatalf("no path of %s on %s", number, cell)
		}
		return x
	}
	g0 := "G0/novatrail/olympiad-math"
	if x := at(g0, "E1.corridor"); !x.Read || math.Abs(x.Base-0.9) > 1e-3 || math.Abs(x.Ceiling-0.5) > 1e-3 {
		t.Errorf("E1 on G0, the shortfall of the corridor to the ceiling of the engine: %+v", x)
	}
	if x := at(g0, "M1.rms.200"); !x.Read || x.Base <= 0.45 || x.Ceiling >= 0.1 {
		t.Errorf("M1 on G0, to the oracle: %+v", x)
	}
	if x := at(g0, "M4.false.all"); x.Read {
		t.Errorf("M4 read on runs with no sum of it: %+v", x)
	}
	if x := at("G2-half/novatrail/olympiad-math", "M3.lag"); !x.Read || x.Ceiling != 0 {
		t.Errorf("M3 on G2-half, to the oracle with no lag: %+v", x)
	}
	if slices.ContainsFunc(p.Paths, func(x Path) bool { return strings.HasPrefix(x.Cell, "A-") }) {
		t.Error("a path on a synthetic graph")
	}
	rows := p.Table()
	if !strings.Contains(rows, "| `"+g0+"` | `E1.corridor` | 0,8998 | 0,4998 | 0,4000 |\n") || !strings.Contains(rows, "| `"+g0+"` | `M4.false.all` | — | — | — |\n") {
		t.Errorf("the table:\n%s", rows)
	}
	var buf bytes.Buffer
	if err := p.Write(&buf); err != nil {
		t.Fatal(err)
	}
	back, err := ReadPaths(bytes.NewReader(buf.Bytes()))
	if err != nil || back.Table() != rows || !reflect.DeepEqual(back, p) {
		t.Errorf("the paths do not read back as written: %v", err)
	}
	if _, err := ReadPaths(strings.NewReader(`{"paths":[],"extra":1}`)); !errors.Is(err, ErrCriterion) {
		t.Errorf("an unknown field: %v", err)
	}
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	part := file("engine", seed.Working, baseSide, false)
	part.Header.Part = bench.Part{K: 0, N: 2}
	broken := file("engine", seed.Working, baseSide, true)
	elsewhere := file("ceiling", seed.Working, ceilSide, false)
	elsewhere.Header.Libraries = map[string]string{"github.com/altnovatrail/engine": "v9.9.9"}
	short := file("ceiling", seed.Working, ceilSide, false)
	short.Lines = short.Lines[1:]
	noG3 := file("engine", seed.Working, baseSide, false)
	noG3.Lines = slices.DeleteFunc(noG3.Lines, func(l bench.Line) bool { return l.Population == "G3" })
	substituted := file("engine", seed.Working, baseSide, false)
	substituted.Header.Packages = []published.Package{{Module: "novatrail/olympiad-math", Version: "0.1.1"}}
	for name, c := range map[string][2]*bench.File{
		"the rules swapped":  {ceiling, base},
		"held-out learners":  {file("engine", held, baseSide, false), file("ceiling", held, ceilSide, false)},
		"a part":             {part, ceiling},
		"other libraries":    {base, elsewhere},
		"a breach":           {broken, ceiling},
		"a learner short":    {base, short},
		"a population short": {noG3, ceiling},
		"packages in place":  {substituted, ceiling},
	} {
		if _, err := RecordPaths(c[0], c[1], "", ""); !errors.Is(err, ErrCriterion) {
			t.Errorf("%s: %v", name, err)
		}
	}
	decisive = learners + 1
	if _, err := RecordPaths(base, ceiling, "", ""); !errors.Is(err, ErrCriterion) || !strings.Contains(err.Error(), "decisive run") {
		t.Errorf("a rough look: %v", err)
	}
}

// The paths the criterion holds read, and the amendment to R04 that fixes them holds the same table.
func TestPathsFixed(t *testing.T) {
	p, err := fixed()
	if err != nil {
		t.Fatal(err)
	}
	decisions, err := os.ReadFile(filepath.Join("..", "..", "docs", "decisions.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(decisions), p.Table()) {
		t.Error("the table of the paths in the amendment to R04 is not that of paths.json: write it anew with sim paths -table")
	}
}

// The report says whether the run decides, names the graphs, the candidates and the outcome.
func TestReport(t *testing.T) {
	usual(t)
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	in := input(Model,
		Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false), Group: "still"},
		Candidate{Name: "same", File: file("engine", seed.Working, betterSide, false), Constants: 3},
		Candidate{Name: "worse", File: file("engine", seed.Working, worseSide, false)},
	)
	in.Confirm = &Confirm{Base: file("engine", held, baseSide, false), Candidate: file("engine", held, betterSide, false)}
	r, err := Judge(in)
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err := r.Write(&b); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"a rough look", "novatrail/olympiad-math", "No graph", "### better", "### worse", "fails no worse", "Equal to the best", "confirmation", "Chosen and confirmed: better"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("the report does not say %q:\n%s", want, b.String())
		}
	}
	r.Decisive = true
	b.Reset()
	if err := r.Write(&b); err != nil || !strings.Contains(b.String(), "a decisive run") {
		t.Errorf("a decisive report: %v", err)
	}
	if showScore(Score{}) != "not read" {
		t.Error("an unread score is shown")
	}
}

// A number whose fixed path is above zero but within the author's tolerance has no share: at the
// tolerance or below it leaves the score, read by the candidate or not, and above it the share is
// read (the amendment of 2026-11-09 to R04). A path to perfection is the base's own shortfall, and a
// path the tolerance in decimal, off by a few ulps, is within it.
func TestPathWithinTolerance(t *testing.T) {
	for _, c := range []struct {
		number        string
		base, ceiling float64
		read          bool
		path          pathKind
	}{
		{"M3.settle", 3, 2.5, true, smallPath},
		{"M3.settle", 3, 2, true, smallPath},
		{"M3.settle", 3, 1.5, true, onPath},
		{"M3.settle", 3, 3, true, noPath},
		{"M3.settle", 3, 2, false, smallPath},
		{"E3.end", 0.0083, 0, true, smallPath},
		{"E3.end", 0.0116, 0, true, onPath},
		{"E1.corridor", 1 - 0.70, 1 - 0.71, true, smallPath},
		{"M2.brier.after7", 0.1559, 0.1555, true, smallPath},
		{"M2.brier.after7", 0.1537, 0.1502, true, onPath},
	} {
		d := cellData{
			pair: map[string]measure.Pair{c.number: {Diff: measure.Value{Read: c.read, Value: 0.5, Resampled: []float64{0.5, 0.25}}}},
			path: map[string]Path{c.number: {Read: true, Base: c.base, Ceiling: c.ceiling}},
		}
		if f := fraction(d, c.number); f.path != c.path {
			t.Errorf("%s, the path %v: %v, want %v", c.number, c.base-c.ceiling, f.path, c.path)
		}
	}
}

// A path within the author's tolerance leaves the score with its weight, the score names it apart,
// a goal on it is not met, and the report says why.
func TestPathSmall(t *testing.T) {
	ceil := baseSide
	ceil.corridor += 0.005 // the ceiling no further from the base than the tolerance of E1, 0.01
	usePaths(t, file("engine", seed.Working, baseSide, false), file("ceiling", seed.Working, ceil, false))
	r, err := Judge(input(Model, Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false), Goal: "corridor"}))
	if err != nil {
		t.Fatal(err)
	}
	v := verdict(t, r, "better")
	if !slices.Contains(v.Score.Small, "G0/novatrail/olympiad-math E1.corridor") || slices.Contains(v.Score.Left, "G0/novatrail/olympiad-math E1.corridor") || !v.Score.Read {
		t.Errorf("small %v, left %v, read %v", v.Score.Small, v.Score.Left, v.Score.Read)
	}
	if !slices.ContainsFunc(v.Checks, func(k Check) bool {
		return k.Kind == "goal" && k.Number == "E1.corridor" && !k.Met && strings.Contains(k.Detail, "tolerance")
	}) {
		t.Errorf("the goal on the corridor: fails %v", failed(v))
	}
	var b bytes.Buffer
	if err := r.Write(&b); err != nil || !strings.Contains(b.String(), "its path within the author's tolerance: G0/novatrail/olympiad-math E1.corridor") {
		t.Errorf("the report: %v\n%s", err, b.String())
	}
}

// The recipe of the decisive run, just decisive, runs Decisive learners a cell, of the working ones
// unless it is given another set: a decisive run of another size would be refused by sim paths, or
// read as a rough look, and a base of other learners would not be the one paths.json names.
func TestDecisiveRecipe(t *testing.T) {
	justfile, err := os.ReadFile(filepath.Join("..", "..", "justfile"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{fmt.Sprintf("-learners %d -set ", Decisive), `decisive rule out set="working":`} {
		if !strings.Contains(string(justfile), want) {
			t.Errorf("the recipe decisive has no %q", want)
		}
	}
}
