package criterion

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
)

// Same is the target M11.same: the share of the learners of each grade but the reference's who come
// to the reference's estimate, the candidate's less the base's, pooled over the grades.
const Same = "M11.same"

// targets are the judgements written down with a target of their own, by name, and their target: the
// judgement file names it, and no other judgement may, so that a judgement goes by the target its
// criterion wrote down before its first candidate and by no other.
var targets = map[string]string{
	"model-fair-1": Same,
}

// Target is a candidate's target against the base: each grade's numbers of M11 and Δ, the target's
// score, the mean of the grades' Δ with its 95 % interval, and the score of each group of the
// criterion, which stands beside its whole score as a reference and decides nothing.
type Target struct {
	Name   string
	Grades []TargetGrade
	Score  Score
	Groups []GroupScore
}

// GroupScore is the score of one group of a criterion.
type GroupScore struct {
	Group string
	Score Score
}

// TargetGrade is one grade of the target: its cell; M11.answers, M11.same and M11.band as they are,
// in the order of measure.FairNumbers, A the base's and B the candidate's, each with its interval,
// and Diff, base − candidate; Δ, the candidate's M11.same less the base's; and whether the candidate
// holds the bounds of fairness on the grade, which the report writes for the confirmation and which
// decide nothing.
type TargetGrade struct {
	Cell                 string
	Numbers              []measure.Pair
	Same, Band           measure.Pair
	Delta                measure.Value
	SameHolds, BandHolds bool
}

// target reads the target M11.same of a candidate against the base off the pairs of the grades and
// the reference that fair reads, resampled by the reference's cell, so that both grades and both
// sides resample the same learners n: Δ of a grade is the candidate's share less the base's, and the
// target's score is the mean of the grades' Δ, the k-th resample of each grade added up. A grade
// whose cell the runs do not have, or whose share a side does not read, leaves the score unread, and
// an unread score is not better than the base. The two sides have the same cells of fairness. The
// mean is over grades, (Δ1 + Δ5) / 2, so a grade on more than one graph is refused: the populations
// of fairness are on olympiad maths alone.
func target(base, candidate []bench.FairCell) (*Target, error) {
	t := &Target{Name: Same}
	var total []float64
	value, read := 0.0, true
	seen := map[string]bool{}
	for _, g := range grades(base) {
		grade, _, _ := strings.Cut(g.cell, "/")
		if seen[grade] {
			return nil, fmt.Errorf("the target reads one cell a grade, and %s has more than one", grade)
		}
		seen[grade] = true
		tg := TargetGrade{Cell: g.cell, Delta: measure.Value{Name: Same}}
		if g.pairs == nil {
			t.Grades, read = append(t.Grades, tg), false
			continue
		}
		ps, err := measure.ReadPair(g.pairs.Reference, g.pairs.Graph, g.pairs.Pairs, fairCell(candidate, g.cell).Pairs, measure.FairNumbers())
		if err != nil {
			return nil, err
		}
		tg.Numbers, tg.Same, tg.Band = ps, named(ps, Same), named(ps, "M11.band")
		tg.SameHolds, tg.BandHolds = bounds(tg.Same.B, tg.Band.B)
		if tg.Same.Diff.Read {
			tg.Delta = negated(tg.Same.Diff)
		}
		t.Grades = append(t.Grades, tg)
		if !tg.Delta.Read {
			read = false
			continue
		}
		if total == nil {
			total = make([]float64, len(tg.Delta.Resampled))
		}
		for i, x := range tg.Delta.Resampled {
			total[i] += x
		}
		value += tg.Delta.Value
	}
	if !read || total == nil {
		return t, nil
	}
	n := float64(len(t.Grades))
	for i := range total {
		total[i] = float64(total[i] / n)
	}
	t.Score = Score{Value: float64(value / n), Read: true, Resampled: total}
	t.Score.Low, t.Score.High = measure.Ends(total)
	return t, nil
}

// named gives the pair of a number by its name among pairs read together.
func named(ps []measure.Pair, name string) measure.Pair {
	return ps[slices.IndexFunc(ps, func(p measure.Pair) bool { return p.Diff.Name == name })]
}

// negated gives a value with its sign turned, its interval and resamples too: b − a of a difference
// a − b. It is 0 − x, not −x, so that no value is −0.
func negated(v measure.Value) measure.Value {
	out := v
	out.Value, out.Low, out.High = 0-v.Value, 0-v.High, 0-v.Low
	out.Resampled = nil
	for _, x := range v.Resampled {
		out.Resampled = append(out.Resampled, 0-x)
	}
	return out
}

// checkTarget refuses a judgement that names a target its criterion did not write down, or names none
// when it did, and a candidate that names a group or a goal where the judgement's target judges every
// candidate.
func (in Input) checkTarget() error {
	want := targets[in.Judgement]
	switch {
	case in.Target == want:
	case want == "":
		return fmt.Errorf("%w: the judgement %q names the target %q, which only %s wrote down before their first candidate", ErrCriterion, in.Judgement, in.Target, strings.Join(slices.Sorted(maps.Keys(targets)), ", "))
	default:
		return fmt.Errorf("%w: the judgement %s goes by the target %s, which its file is to name, and it names %q", ErrCriterion, in.Judgement, want, in.Target)
	}
	for _, c := range in.Candidates {
		if in.Target != "" && (c.Group != "" || c.Goal != "") {
			return fmt.Errorf("%w: candidate %s names a group or a goal, and the judgement %s judges every candidate by its target %s", ErrCriterion, c.Name, in.Judgement, in.Target)
		}
	}
	return nil
}

// deciding gives the score that decides whether a candidate is better than the base and which is
// chosen: the target's in a judgement that names one, where the criterion's score is a reference that
// decides nothing, and the criterion's otherwise.
func (v Verdict) deciding() Score {
	if v.Target != nil {
		return v.Target.Score
	}
	return v.Score
}

// gradeNames names the grades of fairness but the reference's, "D-grade1 and D-grade5".
func gradeNames() string {
	var out []string
	for _, p := range population.Fairness {
		if p.Name != p.Reference {
			out = append(out, p.Name)
		}
	}
	return joinAnd(out)
}
