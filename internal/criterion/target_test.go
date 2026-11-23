package criterion

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// late gives a grade's learners who never come to the reference's estimate: of grade 1 those whose n
// is a multiple of every1, of grade 5 those whose n is a multiple of every5, none where it is 0.
func late(every1, every5 int) func(string, int) int {
	return func(grade string, n int) int {
		every := map[string]int{"D-grade1": every1, "D-grade5": every5}[grade]
		if every > 0 && n%every == 0 {
			return measure.FairLength
		}
		return 0
	}
}

// fairInput is a judgement of model-fair-1 by its target: the base less fair than the reference, a
// quarter of grade 1 and half of grade 5 never coming to its estimate.
func fairInput(cands ...Candidate) Input {
	in := input(Model, cands...)
	in.Base = fairBase(late(4, 2), even)
	in.Judgement, in.Target = "model-fair-1", Same
	return in
}

// fairCandidate is a candidate whose 65 cells go as a side says and whose grades go as from says.
func fairCandidate(set seed.Set, s side, from func(string, int) int) *bench.File {
	return withFairness(bare("engine", set, s, false), olympiad, from, even)
}

// The target M11.same pooled alone decides "better" and the choice: of the candidates that hold
// every constraint, the one fairest on both grades is chosen over one that mends grade 5 alone, one
// better by the score alone is not better, and one that mends grade 5 at grade 1's cost fails fair.
// Δ of a grade is the candidate's share less the base's, and the shares are those the evaluation reads.
func TestTarget(t *testing.T) {
	usual(t)
	r, err := Judge(fairInput(
		Candidate{Name: "fairer", File: fairCandidate(seed.Working, baseSide, soon)},
		Candidate{Name: "grade5", File: fairCandidate(seed.Working, baseSide, late(4, 0))},
		Candidate{Name: "scored", File: fairCandidate(seed.Working, betterSide, late(4, 2))},
		Candidate{Name: "traded", File: fairCandidate(seed.Working, baseSide, late(2, 0))},
	))
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]struct {
		target        float64
		holds, better bool
	}{
		"fairer": {0.375, true, true},
		"grade5": {0.25, true, true},
		"scored": {0, true, false},
		"traded": {0.125, false, true},
	} {
		v := verdict(t, r, name)
		if v.Target == nil || !v.Target.Score.Read || v.Target.Score.Value != want.target || v.Holds != want.holds || v.Better != want.better {
			t.Errorf("%s: target %+v, holds %v, better %v; want %+v", name, v.Target, v.Holds, v.Better, want)
		}
	}
	if s := verdict(t, r, "scored").Score; !s.Read || s.Low <= 0 {
		t.Errorf("the candidate better by the score: %+v", s)
	}
	if r.Chosen != "fairer" || len(r.Equals) != 1 {
		t.Errorf("chosen %q among %v", r.Chosen, r.Equals)
	}
	fairer := verdict(t, r, "fairer").Target
	if fairer.Name != Same || len(fairer.Grades) != 2 || fairer.Grades[0].Delta.Value != 0.25 || fairer.Grades[1].Delta.Value != 0.5 {
		t.Fatalf("the grades of the target: %+v", fairer.Grades)
	}
	if s := fairer.Score; s.Low > s.Value || s.High < s.Value || len(s.Resampled) != measure.Resamples {
		t.Errorf("the target's interval: %v [%v, %v], %d resamples", s.Value, s.Low, s.High, len(s.Resampled))
	}
	// The candidate's numbers of M11 are those its evaluation would read, interval and all.
	e, err := evaluation(t, fairCandidate(seed.Working, baseSide, late(4, 0)))
	if err != nil {
		t.Fatal(err)
	}
	for i, g := range verdict(t, r, "grade5").Target.Grades {
		got, want := g.Numbers, e.Grades[i]
		if g.Cell != want.Cell || got[0].B.Value != want.Answers.Value || got[1].B.Value != want.Same.Value || got[1].B.Low != want.Same.Low || got[1].B.High != want.Same.High || got[2].B.High != want.Band.High || g.SameHolds != want.SameHolds {
			t.Errorf("%s: the target reads %+v, the evaluation %+v", g.Cell, got[1].B, want.Same)
		}
	}
	var b bytes.Buffer
	if err := r.Write(&b); err != nil {
		t.Fatal(err)
	}
	g5 := "D-grade5/" + olympiad[0]
	for _, want := range []string{
		"The judgement model-fair-1 goes by its target, M11.same of D-grade1 and D-grade5 pooled",
		"- the target, M11.same pooled: 0.375 [",
		"- score, a reference that decides nothing: ",
		"- the score of each group, a reference too: still ",
		"#### The target: M11.same of the grades",
		"| `" + g5 + "` | `M11.same` | 0.5 [",
		"| `" + g5 + "` | `M11.answers` | 100 [",
		"- fails fair D-grade1/" + olympiad[0] + " M11.same",
		"chosen: fairer",
	} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("the report does not say %q:\n%s", want, b.String())
		}
	}
	if strings.Contains(b.String(), "bounds of R55 on the held-out") {
		t.Errorf("the bounds stand in a report with no confirmation:\n%s", b.String())
	}
}

// The chosen one is confirmed by the target on the held-out learners of model-fair-1, and the report
// writes whether it holds the bounds of fairness there, which decide nothing.
func TestTargetConfirm(t *testing.T) {
	usual(t)
	held, err := seed.Confirmation("model-fair-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		from          func(string, int) int
		outcome, says string
	}{
		{soon, "chosen and confirmed: fairer", "The chosen one on its held-out learners holds the bounds of R55."},
		{late(0, 4), "chosen and confirmed: fairer", "| `D-grade5/" + olympiad[0] + "` | no: 0.75 ["},
		{late(4, 2), "not confirmed", "The chosen one on its held-out learners does not hold the bounds of R55."},
	} {
		in := fairInput(Candidate{Name: "fairer", File: fairCandidate(seed.Working, baseSide, soon)})
		in.Confirm = &Confirm{Base: fairCandidate(held, baseSide, late(4, 2)), Candidate: fairCandidate(held, baseSide, c.from)}
		r, err := Judge(in)
		if err != nil {
			t.Fatal(err)
		}
		var b bytes.Buffer
		if err := r.Write(&b); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(r.Outcome, c.outcome) || !strings.Contains(b.String(), c.says) || !strings.Contains(b.String(), "#### The bounds of R55 on the held-out learners") {
			t.Errorf("outcome %q, want %q and %q:\n%s", r.Outcome, c.outcome, c.says, b.String())
		}
	}
}

// Runs without the cells of fairness leave the target unread, and no candidate is better.
func TestTargetUnread(t *testing.T) {
	usual(t)
	in := fairInput(Candidate{Name: "bare", File: bare("engine", seed.Working, betterSide, false)})
	in.Base = bare("engine", seed.Working, baseSide, false)
	r, err := Judge(in)
	if err != nil {
		t.Fatal(err)
	}
	v := verdict(t, r, "bare")
	if v.Target == nil || v.Target.Score.Read || v.Better || r.Chosen != "" {
		t.Fatalf("a target with no cells: %+v, chosen %q", v.Target, r.Chosen)
	}
	var b bytes.Buffer
	if err := r.Write(&b); err != nil {
		t.Fatal(err)
	}
	if want := "| `D-grade1/" + olympiad[0] + "` | — | the runs have no such cell | — | — |"; !strings.Contains(b.String(), want) || !strings.Contains(b.String(), "- the target, M11.same pooled: not read") {
		t.Errorf("the report does not say %q:\n%s", want, b.String())
	}
}

// A judgement goes by the target its criterion wrote down and by no other: model-fair-1 names
// M11.same, no other judgement names a target, and a candidate of a judgement with a target names
// no group or goal.
func TestTargetRefuses(t *testing.T) {
	usual(t)
	cand := func(group, goal string) Candidate {
		return Candidate{Name: "c", File: fairCandidate(seed.Working, baseSide, soon), Group: group, Goal: goal}
	}
	for name, c := range map[string]struct {
		change func(*Input)
		want   string
	}{
		"another judgement":      {func(in *Input) { in.Judgement = "model-2" }, `the judgement "model-2" names the target "M11.same", which only model-fair-1 wrote down`},
		"no judgement":           {func(in *Input) { in.Judgement = "" }, `the judgement "" names the target`},
		"no target":              {func(in *Input) { in.Target = "" }, `the judgement model-fair-1 goes by the target M11.same, which its file is to name, and it names ""`},
		"another target":         {func(in *Input) { in.Target = "M11.band" }, `and it names "M11.band"`},
		"a group":                {func(in *Input) { in.Candidates[0] = cand("still", "") }, "candidate c names a group or a goal"},
		"a goal":                 {func(in *Input) { in.Candidates[0] = cand("", "lag") }, "candidate c names a group or a goal"},
		"the engine's criterion": {func(in *Input) { in.Criterion = Engine }, "the judgement model-fair-1 is not of the criterion engine"},
	} {
		in := fairInput(cand("", ""))
		c.change(&in)
		if _, err := Judge(in); !errors.Is(err, ErrCriterion) || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// The target is the mean over grades, so a grade of fairness on a second graph is refused rather
// than weighed twice; the names of the grades are listed as the report writes them.
func TestTargetOneCellAGrade(t *testing.T) {
	usual(t)
	k := slices.IndexFunc(population.Fairness, func(p population.Population) bool { return p.Name == "D-grade5" })
	prior := population.Fairness[k].Graphs
	t.Cleanup(func() { population.Fairness[k].Graphs = prior })
	population.Fairness[k].Graphs = append(slices.Clone(prior), []string{synthetic.Sheet})
	if _, err := Judge(fairInput(Candidate{Name: "fairer", File: fairCandidate(seed.Working, baseSide, soon)})); !errors.Is(err, ErrCriterion) || !strings.Contains(err.Error(), "the target reads one cell a grade, and D-grade5 has more than one") {
		t.Errorf("a grade on two graphs: %v", err)
	}
	for want, names := range map[string][]string{"": nil, "a": {"a"}, "a and b": {"a", "b"}, "a, b and c": {"a", "b", "c"}} {
		if got := joinAnd(names); got != want {
			t.Errorf("joinAnd(%q) = %q", names, got)
		}
	}
}
