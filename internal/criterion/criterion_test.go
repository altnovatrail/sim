package criterion

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
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
	return Input{Criterion: criterion, Base: file("engine", seed.Working, baseSide, false), Ceiling: file("ceiling", seed.Working, ceilSide, false), Candidates: cands}
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
	held, err := seed.HeldOut(seed.CriterionModel, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		s    side
		want string
	}{{betterSide, "confirmed"}, {baseSide, "not confirmed"}} {
		in := input(Model, Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false)})
		in.Confirm = &Confirm{Base: file("engine", held, baseSide, false), Ceiling: file("ceiling", held, ceilSide, false), Candidate: file("engine", held, c.s, false)}
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
	in.Base, in.Ceiling = only(in.Base), only(in.Ceiling)
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
	in.Base, in.Ceiling = noG0(in.Base), noG0(in.Ceiling)
	if r, err = Judge(in); err != nil || verdict(t, r, "floorless").Holds {
		t.Errorf("runs with no G0 hold the floor: %v", err)
	}
}

// The criterion of the engine scores by its groups; a pair holds the constraints of both and is
// better by its library's score, the other's interval not wholly below zero.
func TestEngineAndPair(t *testing.T) {
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
	good := file("engine", seed.Working, betterSide, false)
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	other := file("engine", held, betterSide, false)
	for name, in := range map[string]Input{
		"no criterion":      input("taste"),
		"no ceiling":        {Criterion: Model, Base: good},
		"ceiling as base":   {Criterion: Model, Base: file("ceiling", seed.Working, baseSide, false), Ceiling: file("ceiling", seed.Working, ceilSide, false)},
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
	in.Confirm = &Confirm{Base: good, Ceiling: file("ceiling", seed.Working, ceilSide, false), Candidate: good}
	if _, err := Judge(in); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation on the working learners: %v", err)
	}
	in.Confirm = &Confirm{Base: other, Ceiling: other, Candidate: other}
	if _, err := Judge(in); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation with its ceiling under the engine's rule: %v", err)
	}
	same := input(Model, Candidate{Name: "x", File: other})
	same.Base, same.Ceiling = file("engine", held, baseSide, false), file("ceiling", held, ceilSide, false)
	same.Confirm = &Confirm{Base: same.Base, Ceiling: same.Ceiling, Candidate: other}
	if _, err := Judge(same); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation on the learners of the choice: %v", err)
	}
	bad := input(Model, Candidate{Name: "x", File: good})
	bad.Ceiling = file("ceiling", held, ceilSide, false)
	if _, err := Judge(bad); !errors.Is(err, ErrCriterion) {
		t.Errorf("a ceiling of other learners: %v", err)
	}
	conf := input(Model, Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false)})
	conf.Confirm = &Confirm{Base: other, Ceiling: file("ceiling", seed.Working, ceilSide, false), Candidate: other}
	conf.Confirm.Base = file("engine", held, baseSide, false)
	if _, err := Judge(conf); !errors.Is(err, ErrCriterion) {
		t.Errorf("a confirmation whose ceiling ran on other learners: %v", err)
	}
}

// The report says whether the run decides, names the graphs, the candidates and the outcome.
func TestReport(t *testing.T) {
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	in := input(Model,
		Candidate{Name: "better", File: file("engine", seed.Working, betterSide, false), Group: "still"},
		Candidate{Name: "same", File: file("engine", seed.Working, betterSide, false), Constants: 3},
		Candidate{Name: "worse", File: file("engine", seed.Working, worseSide, false)},
	)
	in.Confirm = &Confirm{Base: file("engine", held, baseSide, false), Ceiling: file("ceiling", held, ceilSide, false), Candidate: file("engine", held, betterSide, false)}
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
