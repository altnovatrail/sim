package criterion

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// cellData is what a judgement reads of one cell: the pair of the base and the candidate, every
// number turned by its direction, and the paths fixed on the cell.
type cellData struct {
	cell, population string
	graph            []string
	noGraph          bool
	pair             map[string]measure.Pair // a is the base, b the candidate: Diff is base − candidate
	path             map[string]Path         // the fixed paths of the cell, by number
	base, candidate  []measure.Sums          // the sums of each learner, the n-th of both the same learner
	learners         []int                   // the number of the n-th learner in the cell
}

// noGraph tells whether a graph has a module with no published package (R19).
func noGraph(graph []string) bool {
	return slices.ContainsFunc(graph, synthetic.IsSynthetic)
}

// names are the numbers a judgement reads off each cell, in a fixed order: not those of fairness by
// the grade, which no cell of its own gives and the constraint of fairness reads in pairs of grades
// (R55).
func names() []string {
	return slices.DeleteFunc(slices.Sorted(maps.Keys(defs)), func(name string) bool {
		return slices.ContainsFunc(measure.FairNumbers(), func(n measure.Number) bool { return n.Name == name })
	})
}

// read reads the cells of a candidate against the base, each beside its fixed paths. The cells of
// fairness by the grade are left out: the constraint of fairness alone reads them, in pairs of
// grades, and no other check, path or note does (R04, R55).
func read(base, candidate *bench.File, paths *Paths) ([]cellData, error) {
	pairs, err := bench.Cells(bench.Side{File: base}, bench.Side{File: candidate})
	if err != nil {
		return nil, fmt.Errorf("the candidate against the base: %w", err)
	}
	ns := names()
	numbers := directed(ns)
	var out []cellData
	for _, c := range pairs {
		if p, _ := population.Named(c.Population); p.Fair() {
			continue
		}
		ps, err := measure.ReadPair(c.Cell, c.Graph, c.A, c.B, numbers)
		if err != nil {
			return nil, err
		}
		d := cellData{cell: c.Cell, population: c.Population, graph: c.Graph, noGraph: noGraph(c.Graph),
			pair: map[string]measure.Pair{}, path: paths.on(c.Cell), base: c.A, candidate: c.B, learners: c.N}
		for i, name := range ns {
			d.pair[name] = ps[i]
		}
		out = append(out, d)
	}
	return out, nil
}

// judge judges one candidate against the base on its criterion, beside the base's cells of fairness
// by the grade; a pair holds the constraints of both, once each, and is scored by the library whose
// goal it names, the other's score beside. A judgement with a target reads it too, and the target
// alone decides whether the candidate is better, the score standing beside as a reference (R65).
func judge(criterion, tgt string, base *bench.File, baseFair []bench.FairCell, paths *Paths, c Candidate) (Verdict, error) {
	candFair, err := bench.Fair(c.File)
	if err != nil {
		return Verdict{}, fmt.Errorf("%w: candidate %s: %w", ErrCriterion, c.Name, err)
	}
	if a, b := fairCells(baseFair), fairCells(candFair); !slices.Equal(a, b) {
		return Verdict{}, fmt.Errorf("%w: candidate %s has the cells of fairness by the grade [%s], and the base [%s]: both runs need the populations D-, which just decisive and just rough run (R55)", ErrCriterion, c.Name, strings.Join(b, ", "), strings.Join(a, ", "))
	}
	cells, err := read(base, c.File, paths)
	if err != nil {
		return Verdict{}, fmt.Errorf("%w: candidate %s: %w", ErrCriterion, c.Name, err)
	}
	v := Verdict{Candidate: c.Name}
	v.Checks = append(v.Checks, invariants(c.File))
	crits := []string{criterion}
	if criterion == Pair {
		crits = []string{Model, Engine}
	}
	seen := map[string]bool{}
	for _, crit := range crits {
		nws := modelNotWorse
		if crit == Engine {
			nws = engineNotWorse
		}
		for _, k := range notWorseChecks(cells, nws) {
			if key := k.Cell + " " + k.Number; !seen[key] {
				seen[key] = true
				v.Checks = append(v.Checks, k)
			}
		}
		if crit == Model {
			v.Checks = append(v.Checks, floors(cells)...)
			if c.Goal != "" {
				v.Checks = append(v.Checks, goals(cells, c.Goal)...)
			}
		}
	}
	v.Checks = append(v.Checks, neutral(cells, c.EqualOn)...)
	fc, err := fair(baseFair, candFair)
	if err != nil {
		return Verdict{}, fmt.Errorf("%w: candidate %s: %w", ErrCriterion, c.Name, err)
	}
	v.Checks = append(v.Checks, fc...)
	v.Holds = !slices.ContainsFunc(v.Checks, func(k Check) bool { return !k.Met && !k.Skipped && !k.NoGraph })
	if v.Note, err = note(cells); err != nil {
		return Verdict{}, fmt.Errorf("%w: candidate %s: %w", ErrCriterion, c.Name, err)
	}
	own, other := criterion, ""
	if criterion == Pair {
		own, other = c.Library, Model
		if c.Library == Model {
			other = Engine
		}
	}
	groups := groupsOf(own, own)
	v.Score = score(cells, groups, "")
	if c.Group != "" {
		g := score(cells, groups, c.Group)
		v.Group = &g
	}
	v.Better = v.Score.Read && v.Score.Low > 0 && (v.Group == nil || v.Group.Read && v.Group.Low > 0)
	if other != "" {
		o := score(cells, groupsOf(other, other), "")
		v.Other = &o
		v.Better = v.Better && o.Read && o.High >= 0
	}
	if tgt != "" {
		if v.Target, err = target(baseFair, candFair); err != nil {
			return Verdict{}, fmt.Errorf("%w: candidate %s: %w", ErrCriterion, c.Name, err)
		}
		for _, g := range groups {
			v.Target.Groups = append(v.Target.Groups, GroupScore{g.id, score(cells, groups, g.id)})
		}
		v.Better = v.Target.Score.Read && v.Target.Score.Low > 0
	}
	return v, nil
}

// invariants is the check of I1–I4 on the candidate's runs: they were checked, and no learner has a
// breach.
func invariants(f *bench.File) Check {
	k := Check{Kind: "invariants", Number: "I1–I4", Met: f.Header.Invariants}
	if !k.Met {
		k.Detail = "the file of results did not record the invariants"
		return k
	}
	for _, l := range f.Lines {
		if len(l.Breaches) > 0 {
			k.Met = false
			k.Cell = l.Population + "/" + strings.Join(l.Graph, "+")
			k.Detail = fmt.Sprintf("%d breaches on learner %d, the first: %s", len(l.Breaches), l.N, l.Breaches[0])
			return k
		}
	}
	return k
}

// notWorseChecks gives the checks of "no worse": on each cell a check holds on, the worst end of the
// 95 % interval of the paired difference, the candidate's less the base's, is within the tolerance —
// the larger of the author's and 4.3 standard errors of the difference (R04). Where the base reads no
// such number, the check has nothing to hold.
func notWorseChecks(cells []cellData, nws []notWorse) []Check {
	var out []Check
	for _, d := range cells {
		for _, nw := range nws {
			if !nw.applies(d.population) {
				continue
			}
			for _, name := range nw.numbers {
				out = append(out, notWorseCheck(Check{Kind: "no worse", Cell: d.cell, Number: name, NoGraph: d.noGraph}, d.pair[name]))
			}
		}
	}
	return out
}

// notWorseCheck completes a check of "no worse" of one number on one cell by its pair: a is the base
// and b the candidate, both turned by the number's direction.
func notWorseCheck(k Check, p measure.Pair) Check {
	switch {
	case !p.A.Read:
		k.Skipped, k.Detail = true, "the base reads none"
	case !p.B.Read || !p.Diff.Read:
		k.Detail = "the candidate's run does not read it"
	default:
		worst, tol := -p.Diff.Low, max(defs[k.Number].tolerance, float64(resolution*p.Diff.SE))
		k.Met = worst <= tol
		k.Detail = fmt.Sprintf("worse by at most %.4g, the tolerance %.4g", worst, tol)
	}
	return k
}

// fair gives the constraint of fairness by the grade, one in the criteria of the model and of the
// engine and in a pair (R55, R57): on the cell of each grade but the reference's, the pairs of the
// grade and the reference of the candidate against those of the base, learner by learner, and each
// number of M11 — M11.same by its shortfall, M11.answers and M11.band as they are — "no worse", as
// on any cell. The pairs read as a cell of the reference's name, so that both sides and every grade
// resample the same learners, as the evaluation does. A grade whose cell the runs do not have fails
// each number: olympiad maths is published, and a check never measured does not hold. The two sides
// have the same cells of fairness, and read has held their learners to one another.
func fair(base, candidate []bench.FairCell) ([]Check, error) {
	var ns []string
	for _, n := range measure.FairNumbers() {
		ns = append(ns, n.Name)
	}
	numbers := directed(ns)
	var out []Check
	for _, g := range grades(base) {
		cell, a := g.cell, g.pairs
		if a == nil {
			for _, name := range ns {
				out = append(out, Check{Kind: "fair", Cell: cell, Number: name, Detail: "the runs have no cell " + cell + " beside its reference " + population.Reference})
			}
			continue
		}
		b := fairCell(candidate, cell)
		ps, err := measure.ReadPair(a.Reference, a.Graph, a.Pairs, b.Pairs, numbers)
		if err != nil {
			return nil, err
		}
		for k, name := range ns {
			out = append(out, notWorseCheck(Check{Kind: "fair", Cell: cell, Number: name}, ps[k]))
		}
	}
	return out, nil
}

// fairCells names the cells of fairness by the grade read in pairs, sorted.
func fairCells(cells []bench.FairCell) []string {
	out := []string{}
	for _, c := range cells {
		out = append(out, c.Cell)
	}
	slices.Sort(out)
	return out
}

// floors gives the floor of the first slice (R04): on G0 and G0-skills1, false "mastered" M4 at
// most 20 % and the answers until "mastered" M6 at most 1.5 times the base's, by the worst end of
// the candidate's interval. Where the base reads no such number, the floor has nothing to hold, as
// "no worse" has not.
func floors(cells []cellData) []Check {
	out := missing(cells, "floor", "G0", "G0-skills1")
	for _, d := range cells {
		if d.population != "G0" && d.population != "G0-skills1" {
			continue
		}
		m4, m6 := d.pair["M4.false.all"], d.pair["M6.late"]
		k4 := Check{Kind: "floor", Cell: d.cell, Number: "M4.false.all", NoGraph: d.noGraph, Skipped: !m4.A.Read,
			Met: m4.B.Read && m4.B.High <= 0.20, Detail: fmt.Sprintf("at most %.4g, the floor 0.2", m4.B.High)}
		k6 := Check{Kind: "floor", Cell: d.cell, Number: "M6.late", NoGraph: d.noGraph, Skipped: !m6.A.Read,
			Met: m6.B.Read && m6.B.High <= float64(1.5*m6.A.Value), Detail: fmt.Sprintf("at most %.4g, the floor 1.5 × %.4g", m6.B.High, m6.A.Value)}
		out = append(out, k4, k6)
	}
	return out
}

// goals gives the checks of the goal a candidate of the model names (R04): on G2-half, M3 closes at
// least 0.55 of its fixed path to the oracle (lag) or E1 at least 0.27 of its fixed path to the
// ceiling (corridor), by the low end of the interval of the share; on G3, at most 25 % of the
// learners never settle after the jump (settle), by the high end.
func goals(cells []cellData, goal string) []Check {
	var out []Check
	if goal == "settle" {
		out = missing(cells, "goal", "G3")
	} else {
		out = missing(cells, "goal", "G2-half")
	}
	for _, d := range cells {
		switch {
		case goal == "lag" && d.population == "G2-half":
			out = append(out, share(d, "M3.lag", 0.55))
		case goal == "corridor" && d.population == "G2-half":
			out = append(out, share(d, "E1.corridor", 0.27))
		case goal == "settle" && d.population == "G3":
			p := d.pair["M3.unsettled"]
			out = append(out, Check{Kind: "goal", Cell: d.cell, Number: "M3.unsettled", NoGraph: d.noGraph, Met: p.B.Read && p.B.High <= 0.25,
				Detail: fmt.Sprintf("at most %.4g, the goal 0.25", p.B.High)})
		}
	}
	return out
}

// neutral gives the checks that a candidate gives the base's runs on the graphs it names: on each
// cell of each graph, every learner's sums are the base's to the bit, so every paired difference is
// zero — on the graphs where no skill of a lesson has a task outside it, a variant of consolidation
// in a lesson has nothing to give (R49). A graph the runs have no cell of is not met: a check that
// is never measured does not hold.
func neutral(cells []cellData, graphs []string) []Check {
	var out []Check
	for _, g := range graphs {
		found := false
		for _, d := range cells {
			if strings.Join(d.graph, "+") != g {
				continue
			}
			found = true
			k := Check{Kind: "neutral", Cell: d.cell, Number: "all", NoGraph: d.noGraph, Met: true,
				Detail: "every learner's sums are the base's"}
			for n := range d.base {
				if key, ok := differs(d.base[n], d.candidate[n]); ok {
					k.Met = false
					k.Detail = fmt.Sprintf("learner %d differs first in %s: the base %v, the candidate %v", d.learners[n], key, d.base[n][key], d.candidate[n][key])
					break
				}
			}
			out = append(out, k)
		}
		if !found {
			out = append(out, Check{Kind: "neutral", Cell: g, Number: "all", Detail: "the runs have no cell on the graph " + g})
		}
	}
	return out
}

// differs gives the first key, in their order, whose sums two learners do not share, a key one of
// them lacks reading as zero, as the measures read it.
func differs(a, b measure.Sums) (string, bool) {
	keys := append(slices.Collect(maps.Keys(a)), slices.Collect(maps.Keys(b))...)
	slices.Sort(keys)
	for _, k := range slices.Compact(keys) {
		if a[k] != b[k] {
			return k, true
		}
	}
	return "", false
}

// note gives the note of consolidation in a lesson: the numbers of E6 on each cell of a published
// graph where the base or the candidate consolidated at least once, base and candidate paired by
// learner. It decides nothing (R49).
func note(cells []cellData) ([]Note, error) {
	consolidated := func(ls []measure.Sums) bool {
		return slices.ContainsFunc(ls, func(s measure.Sums) bool { return s["E6.consolidations"] > 0 })
	}
	var out []Note
	for _, d := range cells {
		if d.noGraph || !consolidated(d.base) && !consolidated(d.candidate) {
			continue
		}
		ps, err := measure.ReadPair(d.cell, d.graph, d.base, d.candidate, measure.NoteNumbers())
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			out = append(out, Note{Cell: d.cell, Pair: p})
		}
	}
	return out, nil
}

// missing gives a failed check for each population a check needs that the runs do not have: a check
// that is never measured is not met.
func missing(cells []cellData, kind string, populations ...string) []Check {
	var out []Check
	for _, p := range populations {
		if !slices.ContainsFunc(cells, func(d cellData) bool { return d.population == p }) {
			out = append(out, Check{Kind: kind, Cell: p, Detail: "the runs have no learner of " + p})
		}
	}
	return out
}

// share is the check that a candidate closes at least a share of the path of a number on a cell.
func share(d cellData, name string, least float64) Check {
	k := Check{Kind: "goal", Cell: d.cell, Number: name, NoGraph: d.noGraph}
	f := fraction(d, name)
	switch {
	case f.path == smallPath:
		k.Detail = "the fixed path is within the author's tolerance, so the share is not read"
		return k
	case f.path == noPath || !f.read:
		k.Detail = "the path is not read, or the candidate's run does not read the number"
		return k
	}
	low, _ := measure.Ends(f.resampled)
	k.Met = low >= least
	k.Detail = fmt.Sprintf("closes at least %.4g of the path, the goal %.4g", low, least)
	return k
}

// pathKind is how the fixed path of a number on a cell stands in the score.
type pathKind int

const (
	noPath    pathKind = iota // not fixed, not read, or not above zero: the number leaves the score
	smallPath                 // above zero but within the author's tolerance: it leaves the score too
	onPath                    // above the tolerance: the share is read
)

// slack is how far above the author's tolerance a path may be and still be within it: the
// difference of two numbers turned by their direction (1 − v for a shortfall) is off by a few ulps,
// and a path of 0.010000000000000009 is the tolerance 0.01 written in decimal.
const slack = 1e-12

// frac is the share of the path of a number on a cell the candidate closes: (base − candidate)/path,
// the difference paired on the learners of the judgement and the path the fixed one, on the whole
// cell and over the resamples.
type frac struct {
	path      pathKind
	read      bool // the candidate's run reads the number
	value     float64
	resampled []float64
}

// fraction gives the share of the path of a number on a cell the candidate closes, by the path
// fixed on the first decisive run of the base (Paths): the same path at every choice and in every
// confirmation, whatever the base and the learners of the judgement. A number with no fixed path,
// one not read, or one of zero or less has no share. Nor has one whose path is above zero but
// within the author's tolerance of the number, which the author has called "no change": a ceiling
// no further from the base than that sets no path — nor, where the path goes to perfection, a base
// no further from it —, and the noise of the difference would decide the share (R04). Either leaves
// the score whether the candidate's run reads the number or not, as its share would not count.
func fraction(d cellData, name string) frac {
	p := d.pair[name]
	x, ok := d.path[name]
	path := float64(x.Base - x.Ceiling)
	switch {
	case !ok || !x.Read || path <= 0:
		return frac{}
	case path-defs[name].tolerance <= slack:
		return frac{path: smallPath}
	}
	f := frac{path: onPath, read: p.Diff.Read}
	if !f.read {
		return f
	}
	f.value = float64(p.Diff.Value / path)
	f.resampled = make([]float64, len(p.Diff.Resampled))
	for i, x := range p.Diff.Resampled {
		f.resampled[i] = float64(x / path)
	}
	return f
}

// score gives the score over groups — all of them, or the one named: the share of the path each
// number closes, the mean over the numbers and cells of a population, weighted by the populations'
// weights. A number with no path, or with a path within the author's tolerance, leaves the score
// with its weight, and the score names it; a number the candidate's run does not read leaves the
// score unread. Its interval comes of the resamples of each cell apart, the k-th of every cell
// added up.
func score(cells []cellData, groups []group, only string) Score {
	var s Score
	var total []float64
	value, weights, unread := 0.0, 0.0, false
	for _, g := range groups {
		if only != "" && g.id != only {
			continue
		}
		for _, wp := range g.populations {
			var sum []float64
			one, n := 0.0, 0
			for _, d := range cells {
				if d.population != wp.population || d.noGraph {
					continue
				}
				for _, name := range g.numbers {
					f := fraction(d, name)
					switch {
					case f.path == smallPath:
						s.Small = append(s.Small, d.cell+" "+name)
						continue
					case f.path == noPath:
						s.Left = append(s.Left, d.cell+" "+name)
						continue
					case !f.read:
						unread = true
						continue
					}
					if sum == nil {
						sum = make([]float64, len(f.resampled))
					}
					for i, x := range f.resampled {
						sum[i] += x
					}
					one += f.value
					n++
				}
			}
			if n == 0 {
				continue
			}
			if total == nil {
				total = make([]float64, len(sum))
			}
			for i := range sum {
				total[i] += float64(wp.weight * float64(sum[i]/float64(n)))
			}
			value += float64(wp.weight * float64(one/float64(n)))
			weights += wp.weight
		}
	}
	if total == nil || unread {
		return s
	}
	for i := range total {
		total[i] = float64(total[i] / weights)
	}
	s.Read, s.Value, s.Resampled = true, float64(value/weights), total
	s.Low, s.High = measure.Ends(total)
	return s
}
