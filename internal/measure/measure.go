// Package measure reads the measures of the bench off the runs of virtual learners, by the truth
// the bench knows and what the learner model and the engine gave (docs/architecture.md, section
// 3.5; R04; R22 in docs/decisions.md). Each learner's run gives its sums; a number of a cell is
// read from the sums of all its learners, and its 95 % interval from 2,000 resamples of the
// learners, the same resamples for every number of the cell, so that two cells of the same learners
// compare by learner (docs/architecture.md, section 3.6).
package measure

import (
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
	var keys []string
	for _, l := range learners {
		keys = append(keys, slices.Collect(maps.Keys(l))...)
	}
	slices.Sort(keys)
	keys = slices.Compact(keys)
	rows := make([][]float64, len(learners))
	all := make([]int, len(learners))
	for i, l := range learners {
		rows[i] = make([]float64, len(keys))
		for k, key := range keys {
			rows[i][k] = l[key]
		}
		all[i] = i
	}
	whole := total(keys, rows, all)
	vals := make([][]float64, len(numbers))
	if len(learners) > 0 {
		s, err := seed.Place(seed.Working.Seed, "bootstrap", cell)
		if err != nil {
			panic(err) // a cell's name is the bench's own
		}
		picks := make([]int, len(learners))
		for range Resamples {
			for i := range picks {
				picks[i] = int(s.Below(uint64(len(learners))))
			}
			t := total(keys, rows, picks)
			for k, n := range numbers {
				if v, ok := n.Of(t); ok {
					vals[k] = append(vals[k], v)
				}
			}
		}
	}
	out := make([]Value, len(numbers))
	for k, n := range numbers {
		v, ok := n.Of(whole)
		out[k] = Value{Name: n.Name, Value: v, Read: ok && len(vals[k]) == Resamples, NoGraph: noGraph}
		if out[k].Read {
			sorted := slices.Sorted(slices.Values(vals[k]))
			out[k].Low, out[k].High = sorted[tail], sorted[Resamples-1-tail]
		}
	}
	return out
}

func itoa(n int) string { return strconv.Itoa(n) }
