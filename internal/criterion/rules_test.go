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

// rulesRun gives a run of the rules: the cells of fairness with a grade at the reference's estimate
// from the answer from(grade, n), the cell C-grade1 likewise, corrected after ten answers, and the
// reference's learners on whose run the guard fired when fired(n); learner-model at RulesCandidate.
func rulesRun(from func(string, int) int, fired func(int) bool) *bench.File {
	f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Libraries: maps.Clone(versions), Invariants: true}}
	f.Header.Libraries["github.com/altnovatrail/learner-model"] = "v0.1.13-0.20261116235821-" + RulesCandidate[:12]
	withFairness(f, olympiad, from, even)
	for i := range f.Lines {
		if l := &f.Lines[i]; l.Population == population.Reference && fired(l.N) {
			l.Fair.Guard = 1
		}
	}
	for n := range learners {
		tr := &measure.Trace{Bands: even("C-grade1", n), Corrected: population.CorrectAfter}
		for k := range measure.FairLength {
			x := 0.0
			if k < from("C-grade1", n) {
				x = 2.5
			}
			tr.Levels = append(tr.Levels, x)
		}
		f.Lines = append(f.Lines, bench.Line{Population: "C-grade1", Graph: olympiad, N: n, Sums: measure.Sums{"learners": 1}, Oracle: measure.Sums{"learners": 1}, Fair: tr})
	}
	return f
}

// evaluateRules evaluates a run of the rules against a base, the decisive run lowered to the test's
// learners and the base's SHA-256 the paths'.
func evaluateRules(t *testing.T, run, base *bench.File) (*RulesResult, error) {
	t.Helper()
	prior := decisive
	decisive = learners
	t.Cleanup(func() { decisive = prior })
	p, err := fixed()
	if err != nil {
		t.Fatal(err)
	}
	return EvaluateRules(Rules{Name: "rules-1", Date: "2026-11-17", Version: "v0.1.28", Run: run, Base: base, RunDigest: "abc", BaseDigest: p.Base, Pins: maps.Clone(pinned)})
}

// The evaluation of the rules holds each of the four conditions of the criterion on its own: grade
// 5 faster than the base, fair, the guard's false firings, and the correction against grade 1
// without it.
func TestEvaluateRules(t *testing.T) {
	// The base: grade 5 at the reference's estimate from answer 60 for every learner, grade 1 from 0.
	base := fairBase(func(g string, _ int) int { return map[string]int{"D-grade5": 60}[g] }, even)
	// The rules: every grade from 0, and the corrected twin of grade 1 as well.
	never := func(int) bool { return false }
	r, err := evaluateRules(t, rulesRun(soon, never), base)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Holds[0] || !r.Holds[1] || !r.Holds[2] || r.Holds[3] || r.All {
		t.Errorf("holds %v, all %v: grade 5 better, fair, no firing, and a correction that changes nothing", r.Holds, r.All)
	}
	if r.Better.Value != 1 || r.BaseSame.Value != 0 || r.Guard.Value != 0 || r.Learners != learners || r.Corrected != learners {
		t.Errorf("M11.same %v against %v, M12.false %v of %d, corrected %d", r.Better.Value, r.BaseSame.Value, r.Guard.Value, r.Learners, r.Corrected)
	}
	// The correction brings grade 1 to the reference sooner than grade 1 without it does.
	late := func(g string, _ int) int { return map[string]int{"D-grade1": 60}[g] }
	if r, err = evaluateRules(t, rulesRun(late, never), fairBase(late, even)); err != nil {
		t.Fatal(err)
	}
	if !r.Holds[3] || r.Holds[0] || !r.Holds[1] {
		t.Errorf("holds %v: the correction faster than grade 1, grade 5 as the base", r.Holds)
	}
	// Grade 5 slower than the base fails fair as well; a guard that fires on every learner of the
	// reference fails the third.
	slow := func(g string, _ int) int { return map[string]int{"D-grade5": 60}[g] }
	if r, err = evaluateRules(t, rulesRun(slow, func(int) bool { return true }), fairBase(soon, even)); err != nil {
		t.Fatal(err)
	}
	if r.Holds[0] || r.Holds[1] || r.Holds[2] || r.Guard.Value != 1 {
		t.Errorf("holds %v, M12.false %v: grade 5 slower, and the guard on every learner", r.Holds, r.Guard.Value)
	}
	var b bytes.Buffer
	if err := r.Write(&b); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"# The short run of the learner model's three rules: rules-1", "in docs/decisions.md", "| 3. M12.false of D-grade3/novatrail/olympiad-math, 60 learners | 1 [1; 1] | at most 0.05 | no |", "The rules are not checked", "## For spec's DPIA"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("the report has no %q:\n%s", want, b.String())
		}
	}
}

// The evaluation of the rules refuses what it does not read: a run with no name or digests, a base
// other than the paths', a run under another rule or not whole, with a breach, on another
// learner-model, without a cell, short of learners, or with cells of other learners.
func TestEvaluateRulesRefuses(t *testing.T) {
	p, err := fixed()
	if err != nil {
		t.Fatal(err)
	}
	never := func(int) bool { return false }
	base := fairBase(soon, even)
	other := func(f *bench.File) *bench.File {
		f.Header.Libraries["github.com/altnovatrail/learner-model"] = "v0.1.8"
		return f
	}
	ceiling := func(f *bench.File) *bench.File { f.Header.Rule = "ceiling"; return f }
	part := func(f *bench.File) *bench.File { f.Header.Part = bench.Part{K: 0, N: 2}; return f }
	unchecked := func(f *bench.File) *bench.File { f.Header.Invariants = false; return f }
	breach := func(f *bench.File) *bench.File { f.Lines[0].Breaches = []string{"I4"}; return f }
	without := func(name string) func(*bench.File) *bench.File {
		return func(f *bench.File) *bench.File {
			var keep []bench.Line
			for _, l := range f.Lines {
				if l.Population != name {
					keep = append(keep, l)
				}
			}
			f.Lines = keep
			return f
		}
	}
	renumber := func(f *bench.File) *bench.File {
		for i := range f.Lines {
			if f.Lines[i].Population == "C-grade1" || f.Lines[i].Population == population.Reference {
				continue
			}
			f.Lines[i].N += 1000
		}
		return f
	}
	shifted := func(f *bench.File) *bench.File {
		for i := range f.Lines {
			if f.Lines[i].Population == "D-grade1" || f.Lines[i].Population == population.Reference || f.Lines[i].Population == "D-grade5" {
				continue
			}
			f.Lines[i].N += 1000
		}
		return f
	}
	untraced := func(f *bench.File) *bench.File { f.Lines[len(f.Lines)-1].Fair = nil; return f }
	for _, c := range []struct {
		name       string
		run        func(*bench.File) *bench.File
		base       func(*bench.File) *bench.File
		digest     string
		decisiveTo int
	}{
		{"another base", nil, nil, "abc", learners},
		{"another learner-model", other, nil, p.Base, learners},
		{"the ceiling", ceiling, nil, p.Base, learners},
		{"a part", part, nil, p.Base, learners},
		{"no invariants", unchecked, nil, p.Base, learners},
		{"a breach in the base", nil, breach, p.Base, learners},
		{"no cell of the correction", without("C-grade1"), nil, p.Base, learners},
		{"short of learners", nil, nil, p.Base, learners + 1},
		{"another base's learners", renumber, nil, p.Base, learners},
		{"the correction's other learners", shifted, nil, p.Base, learners},
		{"a run with no trace", untraced, nil, p.Base, learners},
		{"a base with no trace", nil, untraced, p.Base, learners},
	} {
		run, b := rulesRun(soon, never), fairBase(soon, even)
		if c.run != nil {
			run = c.run(run)
		}
		if c.base != nil {
			b = c.base(b)
		}
		prior := decisive
		decisive = c.decisiveTo
		_, err := EvaluateRules(Rules{Name: "rules-1", Date: "2026-11-17", Run: run, Base: b, RunDigest: "abc", BaseDigest: c.digest, Pins: maps.Clone(pinned)})
		decisive = prior
		if !errors.Is(err, ErrCriterion) {
			t.Errorf("%s: %v", c.name, err)
		}
	}
	if _, err := EvaluateRules(Rules{Run: base, Base: base}); !errors.Is(err, ErrCriterion) {
		t.Errorf("no name: %v", err)
	}
}
