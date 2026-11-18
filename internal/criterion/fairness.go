package criterion

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
)

// The bounds the base holds the criterion of fairness by the grade to, for every grade alike: M11.same
// at least sameFloor at the low end of its interval and M11.band at most bandCeiling at the high end
// (R55).
const (
	sameFloor   = 0.90
	bandCeiling = 0.05
)

// Evaluation is an evaluation of the base by the criterion of fairness by the grade: its name, as its
// folder evaluations/<name>/, the decisive run of the base and its SHA-256, the date the evaluation
// is made on, the libraries the base is fixed at — those of the gate file — and the version of sim
// that reads it.
type Evaluation struct {
	Name, Digest, Date, Version string
	Run                         *bench.File
	Pins                        map[string]Pin
}

// Fairness is what an evaluation reads: each grade against the reference, and whether the base holds.
type Fairness struct {
	Evaluation
	Grades []Grade
	Holds  bool
}

// Grade is one grade against the reference grade: its cell, the numbers of M11 with their 95 %
// intervals, and whether it holds the bounds.
type Grade struct {
	Cell, Reference      string
	Answers, Same, Band  measure.Value
	Learners             int
	SameHolds, BandHolds bool
}

// Evaluate reads the evaluation of fairness by the grade off a decisive run of the base: under the
// engine, on the working learners, the whole run on the pinned packages and the libraries of the gate
// file, with no breach of an invariant and at least a decisive run's learners on each cell of
// fairness, which are every population of fairness on a published graph. Each grade but the
// reference's is read against the reference in pairs by learner, the resamples drawn by the
// reference's cell, so that every grade resamples the same learners (R55).
func Evaluate(e Evaluation) (*Fairness, error) {
	switch {
	case e.Name == "" || e.Date == "" || e.Digest == "":
		return nil, fmt.Errorf("%w: an evaluation needs its name, its date and the SHA-256 of its run", ErrCriterion)
	case len(e.Pins) == 0:
		return nil, fmt.Errorf("%w: an evaluation needs the libraries the base is fixed at, those of the gate file", ErrCriterion)
	}
	if err := whole("the base", e.Run); err != nil {
		return nil, err
	}
	if err := onPins("the base", e.Run, e.Pins); err != nil {
		return nil, err
	}
	cells, err := bench.Fair(e.Run)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCriterion, err)
	}
	f := &Fairness{Evaluation: e, Holds: true}
	for _, g := range grades(cells) {
		cell, c := g.cell, g.pairs
		switch {
		case c == nil:
			return nil, fmt.Errorf("%w: the run has no cell %s, nor its pairs with %s: a run without the cells of fairness evaluates nothing", ErrCriterion, cell, population.Reference)
		case noGraph(c.Graph):
			return nil, fmt.Errorf("%w: cell %s is on a graph with no published package", ErrCriterion, cell)
		case len(c.Pairs) < decisive:
			return nil, fmt.Errorf("%w: cell %s has %d learners, and an evaluation reads a decisive run of %d", ErrCriterion, cell, len(c.Pairs), decisive)
		}
		vs := measure.Read(c.Reference, c.Graph, c.Pairs, measure.FairNumbers())
		d := Grade{Cell: c.Cell, Reference: c.Reference, Answers: vs[0], Same: vs[1], Band: vs[2], Learners: len(c.Pairs)}
		d.SameHolds, d.BandHolds = bounds(d.Same, d.Band)
		f.Holds = f.Holds && d.SameHolds && d.BandHolds
		f.Grades = append(f.Grades, d)
	}
	return f, nil
}

// bounds tells whether a grade's M11.same and M11.band hold the bounds of R55: M11.same at least
// sameFloor at the low end of its interval, M11.band at most bandCeiling at the high end.
func bounds(same, band measure.Value) (bool, bool) {
	return same.Read && same.Low >= sameFloor, band.Read && band.High <= bandCeiling
}

// gradeCell is the cell of a grade that the evaluation and the constraint fair read: its name, and its
// pairs, nil when a file has none.
type gradeCell struct {
	cell  string
	pairs *bench.FairCell
}

// grades gives the cells of fairness by the grade that the evaluation and the constraint fair read:
// the cell of each grade but the reference's on each of its graphs, with its pairs among cells (R55).
func grades(cells []bench.FairCell) []gradeCell {
	var out []gradeCell
	for _, p := range population.Fairness {
		if p.Name == p.Reference {
			continue
		}
		for _, g := range p.Graphs {
			cell := p.Name + "/" + strings.Join(g, "+")
			out = append(out, gradeCell{cell, fairCell(cells, cell)})
		}
	}
	return out
}

// fairCell gives the pairs of a cell of fairness by its name, nil when cells have none.
func fairCell(cells []bench.FairCell, name string) *bench.FairCell {
	if i := slices.IndexFunc(cells, func(c bench.FairCell) bool { return c.Cell == name }); i >= 0 {
		return &cells[i]
	}
	return nil
}

// Write writes the report of the evaluation in Markdown: the criterion, the base and the run, each
// grade's numbers with their intervals, whether the base holds, and the paragraph for dpia 10.2 spec
// (R55).
func (f *Fairness) Write(w io.Writer) error {
	var b strings.Builder
	h := f.Run.Header
	fmt.Fprintf(&b, "# The fairness of the learner model by the grade: %s\n\n", f.Name)
	b.WriteString("## The criterion\n\n")
	fmt.Fprintf(&b, "Written down before its first run, by R55 of docs/decisions.md of sim. The populations %s, on olympiad maths, differ by the grade alone: the assumptions of G0, one truth — that of grade 3 — and one unit for the learner n of every grade, so that the model starts on bands 1, 2 and 3. On the module level θ the model gives after each answer, against the learner n of the reference %s: M11.answers is the number of answers before the first of %d in a row with |θ − θ of the reference| below %.1f, a learner for whom it never holds counting as %d; M11.same is the share of learners for whom that number is at most %d; M11.band is the largest absolute difference, over the bands, between the grade's and the reference's shares of the attempts with a verdict at that band from answer %d on. The intervals are 95 %%, from %d resamples of the learners' numbers, the same in every cell. The base holds when, for every grade alike, M11.same is at least %.2f at the low end of its interval and M11.band at most %.2f at the high end; M11.answers has no bound.\n\n",
		strings.Join(fairnessNames(), ", "), population.Reference, measure.FairAnswers, measure.FairGap, measure.FairLength, measure.FairBy, measure.FairBy+1, measure.Resamples, sameFloor, bandCeiling)
	b.WriteString("## The base\n\n")
	for _, path := range bench.Libraries {
		fmt.Fprintf(&b, "- %s %s\n", path, h.Libraries[path])
	}
	fmt.Fprintf(&b, "- sim %s\n- the decisive run: %d learners a cell, SHA-256 %s\n- the date: %s\n\n", f.Version, h.Learners, f.Digest, f.Date)
	b.WriteString("## The numbers\n\n| Cell | Learners | M11.answers | M11.same | M11.band | Holds |\n|---|---|---|---|---|---|\n")
	for _, g := range f.Grades {
		fmt.Fprintf(&b, "| %s | %d | %s | %s | %s | %s |\n", g.Cell, g.Learners, showInterval(g.Answers), showInterval(g.Same), showInterval(g.Band), holds(g))
	}
	verdict := "holds"
	if !f.Holds {
		verdict = "does not hold"
	}
	fmt.Fprintf(&b, "\nThe base %s the criterion.\n\n", verdict)
	b.WriteString("## For dpia 10.2 spec\n\n")
	ref := gradeOf(population.Reference)
	fmt.Fprintf(&b, "On %s, sim evaluated the learner model by the grade (R147 spec; R55 sim): virtual learners of one skill under grades %s of olympiad maths, %d a grade. The base — %s — %s the criterion written down before the first run.", f.Date, gradeList(), f.minLearners(), libraries(h), verdict)
	for _, g := range f.Grades {
		name, _, _ := strings.Cut(g.Cell, "/")
		fmt.Fprintf(&b, " Grade %s: a share of %s of the learners held within %.1f of the estimate under grade %s after at most %d answers, and %s answers went by on average, a learner who never held counting as %d; the share of the tasks at a band after the %dth answer differed from grade %s's by at most %s.", gradeOf(name), showInterval(g.Same), measure.FairGap, ref, measure.FairBy, showInterval(g.Answers), measure.FairLength, measure.FairBy, ref, showInterval(g.Band))
	}
	b.WriteString(" The numbers are given with their 95 % intervals.\n")
	_, err := io.WriteString(w, b.String())
	return err
}

func fairnessNames() []string {
	var out []string
	for _, p := range population.Fairness {
		out = append(out, p.Name)
	}
	return out
}

func (f *Fairness) minLearners() int {
	least := 0
	for _, g := range f.Grades {
		if least == 0 || g.Learners < least {
			least = g.Learners
		}
	}
	return least
}

func libraries(h bench.Header) string {
	var out []string
	for _, path := range bench.Libraries {
		out = append(out, path+" "+h.Libraries[path])
	}
	return strings.Join(out, " and ")
}

func gradeOf(name string) string {
	return strings.TrimPrefix(name, "D-grade")
}

// gradeList lists the grades of the populations of fairness: "1, 3 and 5".
func gradeList() string {
	var out []string
	for _, name := range fairnessNames() {
		out = append(out, gradeOf(name))
	}
	return joinAnd(out)
}

// joinAnd lists names as "a, b and c": one alone as it is, none as nothing.
func joinAnd(names []string) string {
	if len(names) < 2 {
		return strings.Join(names, "")
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

// showInterval writes a value with its interval, or "unread".
func showInterval(v measure.Value) string {
	if !v.Read {
		return "unread"
	}
	return fmt.Sprintf("%.4g [%.4g; %.4g]", v.Value, v.Low, v.High)
}

func holds(g Grade) string {
	switch {
	case g.SameHolds && g.BandHolds:
		return "yes"
	case !g.SameHolds && !g.BandHolds:
		return "no: M11.same and M11.band"
	case !g.SameHolds:
		return "no: M11.same"
	}
	return "no: M11.band"
}
