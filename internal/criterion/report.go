package criterion

import (
	"fmt"
	"io"
	"strings"
)

// Write writes the report in Markdown: whether the run decides or is a rough look, the graphs of
// the choice and the cells of no graph, each candidate's score, group and the checks it fails, the
// choice, and the confirmation.
func (r *Report) Write(w io.Writer) error {
	var b strings.Builder
	kind := "a rough look, which decides nothing"
	if r.Decisive {
		kind = "a decisive run"
	}
	fmt.Fprintf(&b, "# The criterion of the %s\n\n", r.Criterion)
	fmt.Fprintf(&b, "This is %s: the fewest learners of a cell is %d, and a decisive run takes %d.\n\n", kind, r.Learners, Decisive)
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
		fmt.Fprintf(&b, "Equal to the best: %s; the simplest of them is chosen.\n\n", strings.Join(r.Equals, ", "))
	}
	if r.Confirmed != nil {
		b.WriteString("The confirmation on the held-out learners:\n\n")
		writeVerdict(&b, *r.Confirmed)
	}
	fmt.Fprintf(&b, "%s.\n", r.Outcome)
	_, err := io.WriteString(w, b.String())
	return err
}

// writeVerdict writes one candidate's verdict.
func writeVerdict(b *strings.Builder, v Verdict) {
	fmt.Fprintf(b, "### %s\n\n", v.Candidate)
	fmt.Fprintf(b, "- holds every constraint: %v; better than the base: %v\n", v.Holds, v.Better)
	fmt.Fprintf(b, "- score: %s\n", showScore(v.Score))
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
