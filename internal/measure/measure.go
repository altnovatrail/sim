// Package measure reads the measures of the bench off the runs of virtual learners, by the truth
// the bench knows and what the learner model and the engine gave (docs/architecture.md, section
// 3.5; R04; R22 in docs/decisions.md). Each learner's run gives its sums; a number of a cell is
// read from the sums of all its learners, and its 95 % interval from 2,000 resamples of the
// learners, the same resamples for every number of the cell, so that two cells of the same learners
// compare by learner (docs/architecture.md, section 3.6).
package measure

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"

	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// Sums are what one learner's run gives the measures: named sums, added over learners.
type Sums map[string]float64

func (s Sums) add(key string, v float64) { s[key] += v }

// Number is a number of a cell: its name and how it is read from the sums of the cell's learners.
// It is read only when its sums allow it, for a count above zero.
type Number struct {
	Name string
	Of   func(Sums) (float64, bool)
}

// Value is a number read off a cell: its value and the ends of its 95 % interval. A number the cell
// cannot read is not Read, and in a check it is unmet; on a graph with a module whose package is not
// published it is NoGraph, and its checks are neither met nor failed (docs/architecture.md, section
// 3.6; R19).
type Value struct {
	Name             string
	Value, Low, High float64
	Read             bool
	NoGraph          bool
}

// Resamples is the number of resamples of the learners an interval is read from, and tail the
// resamples left out below its low end and above its high end, 2.5 % of them each.
const (
	Resamples = 2000
	tail      = Resamples / 40
)

// ratio is the number num/den of the sums, read when den is above zero.
func ratio(num, den string) func(Sums) (float64, bool) {
	return func(s Sums) (float64, bool) {
		if s[den] <= 0 {
			return 0, false
		}
		return float64(s[num] / s[den]), true
	}
}

// root is √(num/den), a root mean square.
func root(num, den string) func(Sums) (float64, bool) {
	return func(s Sums) (float64, bool) {
		if s[den] <= 0 {
			return 0, false
		}
		return math.Sqrt(float64(s[num] / s[den])), true
	}
}

// total adds the sums of some learners, each as its values under the cell's keys, picked by their
// places in the order they are given.
func total(keys []string, learners [][]float64, picks []int) Sums {
	sum := make([]float64, len(keys))
	for _, i := range picks {
		for k, v := range learners[i] {
			sum[k] += v
		}
	}
	out := make(Sums, len(keys))
	for k, key := range keys {
		out[key] = sum[k]
	}
	return out
}

// Read reads the numbers off a cell of learners on a graph, each with its 95 % interval: the 51st
// lowest and the 51st highest of its values over Resamples resamples of the learners, drawn from
// the place of the working seed with the names bootstrap and the cell's name, so that the intervals
// are the same on every run and every build. The cell's name is its population and graph, not its
// rule: the cells of the base, a candidate and a ceiling on the same learners draw the same
// resamples. A number that some resample cannot read is not read: its interval would rest only on
// the resamples that drew a case.
func Read(cell string, graph []string, learners []Sums, numbers []Number) []Value {
	noGraph := slices.ContainsFunc(graph, synthetic.IsSynthetic)
	keys, r := rows(learners)
	all := make([]int, len(learners))
	for i := range all {
		all[i] = i
	}
	whole := total(keys, r[0], all)
	vals := make([][]float64, len(numbers))
	resample(cell, len(learners), func(picks []int) {
		t := total(keys, r[0], picks)
		for k, n := range numbers {
			if v, ok := n.Of(t); ok {
				vals[k] = append(vals[k], v)
			}
		}
	})
	out := make([]Value, len(numbers))
	for k, n := range numbers {
		v, ok := n.Of(whole)
		out[k] = interval(Value{Name: n.Name, Value: v, NoGraph: noGraph}, ok, vals[k])
	}
	return out
}

// interval completes a value with the ends of its 95 % interval from its values over the
// resamples, and reads it only when the whole cell and every resample read it.
func interval(v Value, ok bool, vals []float64) Value {
	v.Read = ok && len(vals) == Resamples
	if v.Read {
		sorted := slices.Sorted(slices.Values(vals))
		v.Low, v.High = sorted[tail], sorted[Resamples-1-tail]
	}
	return v
}

// resample calls f with the places of Resamples resamples of n learners, drawn from the place of the
// working seed with the names bootstrap and the cell's name; none for no learners.
func resample(cell string, n int, f func(picks []int)) {
	if n == 0 {
		return
	}
	s, err := seed.Place(seed.Working.Seed, "bootstrap", cell)
	if err != nil {
		panic(err) // a cell's name is the bench's own
	}
	picks := make([]int, n)
	for range Resamples {
		for i := range picks {
			picks[i] = int(s.Below(uint64(n)))
		}
		f(picks)
	}
}

// rows gives the sums of each learner as its values under the keys of all of them, sorted.
func rows(learners ...[]Sums) ([]string, [][][]float64) {
	var keys []string
	for _, ls := range learners {
		for _, l := range ls {
			keys = append(keys, slices.Collect(maps.Keys(l))...)
		}
	}
	slices.Sort(keys)
	keys = slices.Compact(keys)
	out := make([][][]float64, len(learners))
	for j, ls := range learners {
		out[j] = make([][]float64, len(ls))
		for i, l := range ls {
			out[j][i] = make([]float64, len(keys))
			for k, key := range keys {
				out[j][i][k] = l[key]
			}
		}
	}
	return keys, out
}

// Pair is a number read off two cells of the same learners — two rules, or a rule and the oracle —
// and their difference a − b, with its 95 % interval from the same resamples in both
// (docs/architecture.md, section 3.6).
type Pair struct {
	A, B, Diff Value
}

// ReadPair reads the numbers off two cells of the same learners on a graph, the n-th learner of
// each the same learner, and the difference of each number, a − b, with its interval from the same
// resamples of the learners in both: the cells compare by learner.
func ReadPair(cell string, graph []string, a, b []Sums, numbers []Number) ([]Pair, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("measure: cell %s has %d learners on one side and %d on the other", cell, len(a), len(b))
	}
	noGraph := slices.ContainsFunc(graph, synthetic.IsSynthetic)
	keys, r := rows(a, b)
	all := make([]int, len(a))
	for i := range all {
		all[i] = i
	}
	wa, wb := total(keys, r[0], all), total(keys, r[1], all)
	xs, ys, diffs := make([][]float64, len(numbers)), make([][]float64, len(numbers)), make([][]float64, len(numbers))
	resample(cell, len(a), func(picks []int) {
		ta, tb := total(keys, r[0], picks), total(keys, r[1], picks)
		for k, n := range numbers {
			x, okx := n.Of(ta)
			y, oky := n.Of(tb)
			if okx {
				xs[k] = append(xs[k], x)
			}
			if oky {
				ys[k] = append(ys[k], y)
			}
			if okx && oky {
				diffs[k] = append(diffs[k], x-y)
			}
		}
	})
	out := make([]Pair, len(numbers))
	for k, n := range numbers {
		x, okx := n.Of(wa)
		y, oky := n.Of(wb)
		out[k] = Pair{
			A:    interval(Value{Name: n.Name, Value: x, NoGraph: noGraph}, okx, xs[k]),
			B:    interval(Value{Name: n.Name, Value: y, NoGraph: noGraph}, oky, ys[k]),
			Diff: interval(Value{Name: n.Name, Value: x - y, NoGraph: noGraph}, okx && oky, diffs[k]),
		}
	}
	return out, nil
}

func itoa(n int) string { return strconv.Itoa(n) }
