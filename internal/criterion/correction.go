package criterion

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
)

// Correction is an evaluation of the adult's correction against the truth: its name, as its folder
// evaluations/<name>/, the run of C-grade1 and D-grade1 and its SHA-256, the date, the version of
// sim that reads it and the libraries of the gate file, which the engine of the run is held to.
type Correction struct {
	Name, Date, Version, Digest string
	Run                         *bench.File
	Pins                        map[string]Pin
}

// CorrectionResult is what the evaluation of the correction reads: M13 of the cell of the
// correction against the same learners without it, each number by name, the learners and how many
// of them the adult corrected, and whether the criterion holds.
type CorrectionResult struct {
	Correction
	Numbers             []measure.Pair
	Learners, Corrected int
	Holds               bool
}

// EvaluateCorrection reads the evaluation of the correction against the truth: a whole run of
// the rules, as EvaluateRules holds it, whose cells C-grade1 and D-grade1 of olympiad maths have the
// same learners in the same order, at least a decisive run's, each with the truth in its trace. It
// holds when the paired difference of M13.same, the correction's less grade 1's, is above zero at
// the low end of its interval.
func EvaluateCorrection(c Correction) (*CorrectionResult, error) {
	if c.Name == "" || c.Date == "" || c.Digest == "" {
		return nil, fmt.Errorf("%w: an evaluation of the correction needs its name, its date and the SHA-256 of its run", ErrCriterion)
	}
	if err := whole("the run of the correction", c.Run); err != nil {
		return nil, err
	}
	if err := onCandidate("the run of the correction", c.Run, c.Pins); err != nil {
		return nil, err
	}
	cells := map[string][]bench.Line{}
	for _, l := range c.Run.Lines {
		cell := l.Population + "/" + strings.Join(l.Graph, "+")
		if cell != correctionCell && cell != uncorrectedCell {
			continue
		}
		if l.Fair == nil || len(l.Fair.Truth) == 0 {
			return nil, fmt.Errorf("%w: learner %s/%d has no truth in its trace", ErrCriterion, cell, l.N)
		}
		cells[cell] = append(cells[cell], l)
	}
	corrected, plain := cells[correctionCell], cells[uncorrectedCell]
	switch {
	case len(corrected) < decisive || len(plain) < decisive:
		return nil, fmt.Errorf("%w: cells %s and %s have %d and %d learners, and the evaluation reads %d", ErrCriterion, correctionCell, uncorrectedCell, len(corrected), len(plain), decisive)
	case !slices.EqualFunc(corrected, plain, func(a, b bench.Line) bool { return a.N == b.N }):
		return nil, fmt.Errorf("%w: cell %s is not of the learners of %s", ErrCriterion, correctionCell, uncorrectedCell)
	}
	out := &CorrectionResult{Correction: c, Learners: len(corrected)}
	var a, b []measure.Sums
	for i := range corrected {
		a, b = append(a, measure.Caught(corrected[i].Fair)), append(b, measure.Caught(plain[i].Fair))
		out.Corrected += b2i(corrected[i].Fair.Corrected > 0)
	}
	var err error
	if out.Numbers, err = measure.ReadPair(referenceCell, corrected[0].Graph, a, b, measure.TruthNumbers()); err != nil {
		return nil, err
	}
	same := out.Numbers[1].Diff
	out.Holds = same.Read && same.Low > 0
	return out, nil
}

// onCandidate holds a run of the rules to the libraries the criterion of the rules names:
// learner-model at RulesCandidate, the engine at the gate file's pin.
func onCandidate(name string, f *bench.File, pins map[string]Pin) error {
	pins = maps.Clone(pins)
	pins["github.com/altnovatrail/learner-model"] = Pin{Commit: RulesCandidate}
	return onPins(name, f, pins)
}

// onPins holds a run to the points of its libraries: by the tag, or by the pseudo-version of the
// commit while there is none.
func onPins(name string, f *bench.File, pins map[string]Pin) error {
	for _, lib := range bench.Libraries {
		if pin, got := pins[lib], f.Header.Libraries[lib]; !pin.names(got) {
			return fmt.Errorf("%w: %s is on %s %q, and it is fixed at %v", ErrCriterion, name, lib, got, pin)
		}
	}
	return nil
}

// Write writes the report of the evaluation of the correction in Markdown: the criterion, the run,
// M13 of both cells and their difference, the outcome, and the paragraph for the DPIA.
func (r *CorrectionResult) Write(w io.Writer) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# The adult's correction against the truth: %s\n\n## The criterion\n\n", r.Name)
	fmt.Fprintf(&b, "Written down before the run, by R73 of docs/decisions.md of sim. %s is %s — grade 1, the truth of grade 3 and the unit of %s — with an adult who sets the module level to the band of the truth after the %dth answer. M13.answers is the number of answers before the first of %d in a row with |θ − θ*| below %.1f, θ the model's module level and θ* the truth's after each answer, a learner for whom it never holds counting as %d; M13.same is the share of learners for whom that number is at most %d. The correction holds when the paired difference of M13.same, %s's less %s's, is above zero at the low end of its 95 %% interval, from %d resamples of the learners' numbers.\n\n",
		correctionCell, uncorrectedCell, population.Reference, population.CorrectAfter, measure.FairAnswers, measure.FairGap, measure.FairLength, measure.FairBy, correctionCell, uncorrectedCell, measure.Resamples)
	fmt.Fprintf(&b, "## The run\n\n- %s, %d learners a cell, SHA-256 %s\n- sim %s\n- the date: %s\n\n", libraries(r.Run.Header), r.Run.Header.Learners, r.Digest, r.Version, r.Date)
	fmt.Fprintf(&b, "## The numbers\n\n%d of the %d learners of %s were corrected.\n\n| Number | %s | %s | Difference |\n|---|---|---|---|\n", r.Corrected, r.Learners, correctionCell, correctionCell, uncorrectedCell)
	for _, p := range r.Numbers {
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", p.A.Name, showInterval(p.A), showInterval(p.B), showInterval(p.Diff))
	}
	verdict := "holds"
	if !r.Holds {
		verdict = "does not hold"
	}
	fmt.Fprintf(&b, "\nThe correction %s the criterion.\n\n## For dpia 10.2 spec\n\n", verdict)
	fmt.Fprintf(&b, "On %s, sim checked the adult's correction of R157 spec — %s — on %d virtual learners of grade 1 of olympiad maths whose truth is the band of grade 3, by a criterion written down before the run (R73 sim). With an adult who set the level to the band of the truth after the %dth answer, the model's estimate came within %.1f of the truth for %d answers in a row by the %dth answer for a share of %s of the learners, against %s of the same learners without the correction, a difference of %s. The correction %s the criterion. The numbers are given with their 95 %% intervals.\n",
		r.Date, libraries(r.Run.Header), r.Learners, population.CorrectAfter, measure.FairGap, measure.FairAnswers, measure.FairBy, showInterval(r.Numbers[1].A), showInterval(r.Numbers[1].B), showInterval(r.Numbers[1].Diff), verdict)
	_, err := io.WriteString(w, b.String())
	return err
}
