package run

import (
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/seed"
)

// The guard fired where the count of a module goes from one below guardFailures to it, and nowhere
// else: not on the next failure of the same run, nor in another module.
func TestFired(t *testing.T) {
	state := func(counts ...int32) learnermodel.State {
		s := learnermodel.State{}
		for i, c := range counts {
			s.Modules = append(s.Modules, learnermodel.ModuleState{Module: []string{"a", "b"}[i], FailuresInARow: c})
		}
		return s
	}
	for _, c := range []struct {
		before, after learnermodel.State
		module        string
		want          bool
	}{
		{state(2), state(3), "a", true},
		{state(3), state(4), "a", false},
		{state(1), state(2), "a", false},
		{state(2), state(0), "a", false},
		{state(0, 2), state(0, 3), "b", true},
		{state(2, 0), state(3, 0), "b", false},
		{state(), state(3), "a", false},
		{state(2), state(3), "c", false},
	} {
		if got := fired(c.before, c.after, c.module); got != c.want {
			t.Errorf("fired(%v, %v, %s) = %v", c.before, c.after, c.module, got)
		}
	}
}

// An adult of C- corrects the module level once, right after the attempt that brings the first
// answers to CorrectAfter, to the band of the truth: the run records the answers before it, the
// model's level stands in the truth's band then, and the state at the end differs from the same
// learner's with no adult. A learner with no adult has no correction.
func TestCorrection(t *testing.T) {
	mods := modules(t, olympiad)
	runOf := func(after int) (*Record, *learner.Learner) {
		a := learner.Base()
		a.Grade, a.TruthGrade, a.CorrectAfter = 1, 3, after
		l, err := learner.New(seed.Working, "D-grade3", 0, mods, a)
		if err != nil {
			t.Fatal(err)
		}
		r, err := New(l, mods, Maths)
		if err != nil {
			t.Fatal(err)
		}
		var at learnermodel.State
		for !r.Done() {
			if err := r.Step(); err != nil {
				t.Fatal(err)
			}
			if r.corrected && at.Modules == nil {
				at = r.state
			}
		}
		if after > 0 {
			m, _ := l.Module(olympiad)
			ro, err := learnermodel.Read(at, r.graphs, 0, r.grade, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !ro.Modules[0].HasLevelBand || ro.Modules[0].LevelBand.Place != m.TrueBand {
				t.Errorf("after the correction the level stands in band %v, and the truth's is %d", ro.Modules[0].LevelBand, m.TrueBand)
			}
		}
		return r.Record(), l
	}
	plain, _ := runOf(0)
	corrected, _ := runOf(10)
	if plain.Corrected != 0 {
		t.Errorf("a learner with no adult corrected after %d answers", plain.Corrected)
	}
	if corrected.Corrected < 10 || corrected.Corrected >= 20 {
		t.Errorf("the adult corrected after %d answers, want the attempt that reaches 10", corrected.Corrected)
	}
	if plain.State.Modules[0].Level == corrected.State.Modules[0].Level {
		t.Errorf("the correction left the level at the end %v", plain.State.Modules[0].Level)
	}
	for _, c := range []struct {
		rec  *Record
		want int
	}{{plain, 0}, {corrected, 1}} {
		steps := 0
		for _, s := range c.rec.Steps {
			if s.Correction != nil {
				steps++
				if !s.Correction.LevelSetByAdult {
					t.Errorf("the step records evidence %+v as the correction", *s.Correction)
				}
			}
		}
		if steps != c.want {
			t.Errorf("%d steps record a correction, want %d", steps, c.want)
		}
	}
	if _, err := learner.New(seed.Working, "D-grade3", 0, mods, learner.Assumptions{CorrectAfter: -1}); err == nil {
		t.Error("a correction after -1 answers is taken")
	}
	if _, err := learner.New(seed.Working, "D-grade3", 0, modules(t, olympiad, grade3), learner.Assumptions{CorrectAfter: 10}); err == nil {
		t.Error("a correction of a learner of two modules is taken")
	}
}
