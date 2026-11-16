package criterion

import (
	"fmt"
	"io"
	"strings"

	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
)

// Write writes the report in Markdown: whether the run decides or is a rough look, the graphs of
// the choice and the cells of no graph, each candidate's score, group, the checks it fails and the
// note of consolidation in a lesson, the choice, and the confirmation.
func (r *Report) Write(w io.Writer) error {
	var b strings.Builder
	kind := "a rough look, which decides nothing"
	if r.Decisive {
		kind = "a decisive run"
	}
	fmt.Fprintf(&b, "# The criterion of the %s\n\n", r.Criterion)
	fmt.Fprintf(&b, "This is %s: the fewest learners of a cell is %d, and a decisive run takes %d.\n\n", kind, r.Learners, Decisive)
	if r.Target != "" {
		fmt.Fprintf(&b, "The judgement %s goes by its target, %s of %s pooled: the mean of the grades' Δ, each the candidate's %s against %s less the base's, with its 95 %% interval. The target alone decides \"better\" and the choice; the score of the criterion is a reference, which decides nothing; the bounds of R55 are written for the confirmation and decide nothing either (R65).\n\n",
			r.Judgement, r.Target, gradeNames(), r.Target, population.Reference)
	}
	fmt.Fprintf(&b, "The choice is made on the graphs %s.", strings.Join(r.Graphs, ", "))
	if len(r.NoGraph) > 0 {
		fmt.Fprintf(&b, " No graph — a module with no published package, its checks neither met nor failed: %s.", strings.Join(r.NoGraph, ", "))
	}
	b.WriteString("\n\n## Candidates\n\n")
	for _, v := range r.Verdicts {
		writeVerdict(&b, v)
	}
	b.WriteString("## Outcome\n\n")
	if len(r.Equals) > 1 {
		chosen := "is chosen"
		if !r.Decisive {
			chosen = "would be chosen by a decisive run"
		}
		fmt.Fprintf(&b, "Equal to the best: %s; the simplest of them %s.\n\n", strings.Join(r.Equals, ", "), chosen)
	}
	if r.Confirmed != nil {
		fmt.Fprintf(&b, "The confirmation on %s, the held-out learners of the judgement %s:\n\n", r.HeldOut, r.Judgement)
		writeVerdict(&b, *r.Confirmed)
		if t := r.Confirmed.Target; t != nil {
			writeBounds(&b, *t)
		}
	}
	fmt.Fprintf(&b, "%s.\n", r.Outcome)
	_, err := io.WriteString(w, b.String())
	return err
}

// writeVerdict writes one candidate's verdict.
func writeVerdict(b *strings.Builder, v Verdict) {
	fmt.Fprintf(b, "### %s\n\n", v.Candidate)
	fmt.Fprintf(b, "- holds every constraint: %v; better than the base: %v\n", v.Holds, v.Better)
	if t := v.Target; t != nil {
		fmt.Fprintf(b, "- the target, %s pooled: %s\n", t.Name, showScore(t.Score))
		fmt.Fprintf(b, "- score, a reference that decides nothing: %s\n", showScore(v.Score))
		var gs []string
		for _, g := range t.Groups {
			gs = append(gs, g.Group+" "+showScore(g.Score))
		}
		fmt.Fprintf(b, "- the score of each group, a reference too: %s\n", strings.Join(gs, "; "))
	} else {
		fmt.Fprintf(b, "- score: %s\n", showScore(v.Score))
	}
	if v.Group != nil {
		fmt.Fprintf(b, "- its group: %s\n", showScore(*v.Group))
	}
	if v.Other != nil {
		fmt.Fprintf(b, "- the other library's score: %s\n", showScore(*v.Other))
	}
	writeLeft(b, "", v.Score)
	if v.Other != nil {
		writeLeft(b, "the other library's score: ", *v.Other)
	}
	var failed, skipped []string
	for _, k := range v.Checks {
		switch {
		case k.NoGraph:
		case k.Skipped:
			skipped = append(skipped, k.Cell+" "+k.Number)
		case !k.Met:
			failed = append(failed, fmt.Sprintf("%s %s %s: %s", k.Kind, k.Cell, k.Number, k.Detail))
		}
	}
	for _, f := range failed {
		fmt.Fprintf(b, "- fails %s\n", f)
	}
	if len(skipped) > 0 {
		fmt.Fprintf(b, "- nothing to hold, the base reads none: %s\n", strings.Join(skipped, "; "))
	}
	b.WriteString("\n")
	if v.Target != nil {
		writeTarget(b, *v.Target)
	}
	writeNote(b, v.Note)
}

// writeTarget writes the numbers of M11 of each grade of the target, the base's and the candidate's
// as they are, with their intervals, and their difference, the candidate's less the base's: that of
// M11.same is the grade's Δ (R65).
func writeTarget(b *strings.Builder, t Target) {
	fmt.Fprintf(b, "#### The target: %s of the grades\n\n", t.Name)
	b.WriteString("| cell | number | base | candidate | candidate − base |\n|---|---|---|---|---|\n")
	for _, g := range t.Grades {
		if g.Numbers == nil {
			fmt.Fprintf(b, "| `%s` | — | the runs have no such cell | — | — |\n", g.Cell)
			continue
		}
		for _, p := range g.Numbers {
			fmt.Fprintf(b, "| `%s` | `%s` | %s | %s | %s |\n", g.Cell, p.Diff.Name, showWith(p.A), showWith(p.B), showWith(negated(p.Diff)))
		}
	}
	b.WriteString("\n")
}

// writeBounds writes whether the chosen candidate holds the bounds of R55 on its held-out learners:
// on every grade, M11.same at least sameFloor at the low end of its interval and M11.band at most
// bandCeiling at the high end. They decide nothing in a judgement with a target; the outcome goes to
// dpia 10.2 spec (R65).
func writeBounds(b *strings.Builder, t Target) {
	fmt.Fprintf(b, "#### The bounds of R55 on the held-out learners, which decide nothing here\n\n")
	fmt.Fprintf(b, "| cell | M11.same at least %.2f at its low end | M11.band at most %.2f at its high end |\n|---|---|---|\n", sameFloor, bandCeiling)
	holds := true
	for _, g := range t.Grades {
		fmt.Fprintf(b, "| `%s` | %s: %s | %s: %s |\n", g.Cell, yesNo(g.SameHolds), showWith(g.Same.B), yesNo(g.BandHolds), showWith(g.Band.B))
		holds = holds && g.SameHolds && g.BandHolds
	}
	verdict := "holds"
	if !holds {
		verdict = "does not hold"
	}
	fmt.Fprintf(b, "\nThe chosen one on its held-out learners %s the bounds of R55.\n\n", verdict)
}

// yesNo writes whether a bound holds.
func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}

// writeNote writes the note of consolidation in a lesson, E6, on the cells where the base or the
// candidate consolidated: it explains where a shift of E1 comes from and decides nothing (R49).
func writeNote(b *strings.Builder, note []Note) {
	if len(note) == 0 {
		return
	}
	b.WriteString("#### Note, which decides nothing: consolidation in a lesson\n\n")
	b.WriteString("| cell | number | base | candidate | candidate − base |\n|---|---|---|---|---|\n")
	for _, n := range note {
		diff := "—"
		if n.Diff.Read {
			diff = showWith(negated(n.Diff))
		}
		fmt.Fprintf(b, "| `%s` | `%s` | %s | %s | %s |\n", n.Cell, n.Diff.Name, showValue(n.A), showValue(n.B), diff)
	}
	b.WriteString("\n")
}

// showWith writes a value with its interval as a score is written, or says it is not read.
func showWith(v measure.Value) string {
	return showScore(Score{Value: v.Value, Low: v.Low, High: v.High, Read: v.Read})
}

// showValue writes a value of the note, or a dash when it is not read.
func showValue(v measure.Value) string {
	if !v.Read {
		return "—"
	}
	return fmt.Sprintf("%.4g", v.Value)
}

// writeLeft writes the numbers that left a score, and why. The group's score is of the same numbers
// as the whole one, so only the whole score and, in a pair, the other library's are written.
func writeLeft(b *strings.Builder, of string, s Score) {
	if len(s.Left) > 0 {
		fmt.Fprintf(b, "- %sleft the score, its path not read or not above zero: %s\n", of, strings.Join(s.Left, "; "))
	}
	if len(s.Small) > 0 {
		fmt.Fprintf(b, "- %sleft the score, its path within the author's tolerance: %s\n", of, strings.Join(s.Small, "; "))
	}
}

// showScore writes a score with its interval, or says it is not read.
func showScore(s Score) string {
	if !s.Read {
		return "not read"
	}
	return fmt.Sprintf("%.4g [%.4g, %.4g]", s.Value, s.Low, s.High)
}
