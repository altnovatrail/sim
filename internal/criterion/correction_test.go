package criterion

import (
	"bytes"
	"errors"
	"maps"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/seed"
)

// correctionRun gives a run of C-grade1 and D-grade1 on learner-model at RulesCandidate, each
// learner at its truth from the answer from(cell, n).
func correctionRun(from func(string, int) int) *bench.File {
	f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Libraries: maps.Clone(versions), Invariants: true, Learners: learners}}
	f.Header.Libraries["github.com/altnovatrail/learner-model"] = "v0.1.13-0.20261116235821-" + RulesCandidate[:12]
	for _, p := range []string{"D-grade1", "C-grade1"} {
		for n := range learners {
			tr := &measure.Trace{Bands: even(p, n)}
			if p == "C-grade1" {
				tr.Corrected = population.CorrectAfter
			}
			for k := range measure.FairLength {
				x := 0.0
				if k < from(p, n) {
					x = 2.5
				}
				tr.Levels, tr.Truth = append(tr.Levels, x), append(tr.Truth, 0)
			}
			f.Lines = append(f.Lines, bench.Line{Population: p, Graph: olympiad, N: n, Sums: measure.Sums{"learners": 1}, Fair: tr})
		}
	}
	return f
}

func evaluateCorrection(t *testing.T, f *bench.File) (*CorrectionResult, error) {
	t.Helper()
	prior := decisive
	decisive = learners
	t.Cleanup(func() { decisive = prior })
	return EvaluateCorrection(Correction{Name: "rules-2", Date: "2026-11-17", Version: "v0.1.28", Digest: "abc", Run: f, Pins: maps.Clone(pinned)})
}

// The correction holds when its learners come to the truth by the 50th answer more often than the
// same learners without it, and not when they come no sooner.
func TestEvaluateCorrection(t *testing.T) {
	faster := func(p string, n int) int {
		if p == "C-grade1" {
			return 0
		}
		return 60 + n%2
	}
	r, err := evaluateCorrection(t, correctionRun(faster))
	if err != nil {
		t.Fatal(err)
	}
	if !r.Holds || r.Numbers[1].A.Value != 1 || r.Numbers[1].B.Value != 0 || r.Learners != learners || r.Corrected != learners {
		t.Errorf("holds %v: M13.same %v against %v, %d learners, %d corrected", r.Holds, r.Numbers[1].A.Value, r.Numbers[1].B.Value, r.Learners, r.Corrected)
	}
	var b bytes.Buffer
	if err := r.Write(&b); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"# The adult's correction against the truth: rules-2", "in docs/decisions.md", "The correction holds the criterion.", "## For spec's DPIA"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("the report has no %q:\n%s", want, b.String())
		}
	}
	if r, err = evaluateCorrection(t, correctionRun(soon)); err != nil || r.Holds {
		t.Errorf("a correction no faster holds %v: %v", r != nil && r.Holds, err)
	}
	b.Reset()
	if err := r.Write(&b); err != nil || !strings.Contains(b.String(), "The correction does not hold the criterion.") {
		t.Errorf("the report of a correction that does not hold: %v\n%s", err, b.String())
	}
}

// The evaluation of the correction refuses a run with no name, not whole, on another learner-model,
// without the truth, short of learners, or with cells of other learners.
func TestEvaluateCorrectionRefuses(t *testing.T) {
	for name, change := range map[string]func(*bench.File){
		"the ceiling":           func(f *bench.File) { f.Header.Rule = "ceiling" },
		"another learner-model": func(f *bench.File) { f.Header.Libraries["github.com/altnovatrail/learner-model"] = "v0.1.8" },
		"no truth":              func(f *bench.File) { f.Lines[0].Fair.Truth = nil },
		"no trace":              func(f *bench.File) { f.Lines[0].Fair = nil },
		"short of learners":     func(f *bench.File) { f.Lines = f.Lines[1:] },
		"other learners":        func(f *bench.File) { f.Lines[len(f.Lines)-1].N += 1000 },
	} {
		f := correctionRun(soon)
		change(f)
		if _, err := evaluateCorrection(t, f); !errors.Is(err, ErrCriterion) {
			t.Errorf("%s: %v", name, err)
		}
	}
	if _, err := EvaluateCorrection(Correction{Run: correctionRun(soon)}); !errors.Is(err, ErrCriterion) {
		t.Errorf("no name: %v", err)
	}
}
