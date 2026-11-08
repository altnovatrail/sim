package run

import (
	"errors"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/seed"
)

// Each attempt with a verdict is observed: an item beside each item answered, in order, with P and
// P* among probabilities and P without retention no lower than with it; the skills after it, and
// before it once the model has a skill of the module;
// and the answers before it counted — in a lesson too, task by task. An attempt with a rubric is not observed; the essays have tasks
// of a choice too.
func TestObserve(t *testing.T) {
	for _, id := range []string{olympiad, grade3, essays} {
		limit := Maths
		if id == essays {
			limit = Essays
		}
		rec := runOf(t, "G0", 3, limit, modules(t, id))
		answers, observed, lessons := 0, 0, 0
		for _, s := range rec.Steps {
			if s.Pass != nil && len(s.Obs) > 1 {
				lessons++
			}
			if len(s.Obs) != len(s.Attempts) {
				t.Fatalf("%s: %d observations of %d attempts", rec.Unit, len(s.Obs), len(s.Attempts))
			}
			for i, o := range s.Obs {
				a := s.Attempts[i]
				if o == nil {
					if len(a.Items) > 0 {
						t.Errorf("%s: an attempt of %d items is not observed", rec.Unit, len(a.Items))
					}
					continue
				}
				observed++
				if len(o.Items) != len(a.Items) || o.Answers != answers || len(o.After) == 0 || (answers > 0 && len(o.Before) == 0) {
					t.Fatalf("%s: observed %d items of %d, %d answers before of %d", rec.Unit, len(o.Items), len(a.Items), o.Answers, answers)
				}
				for k, it := range o.Items {
					if it.PModel <= 0 || it.PModel >= 1 || it.PTrue < 0 || it.PTrue > 1 || it.PFresh < it.PModel-1e-12 || it.PGiven <= 0 || it.PGiven > 1 || it.Many != (len(o.Items) > 1) {
						t.Errorf("%s: item %d observed %+v", rec.Unit, k, it)
					}
					if it.PTrue != a.Items[k].Chance {
						t.Errorf("%s: P* %v beside the chance %v", rec.Unit, it.PTrue, a.Items[k].Chance)
					}
				}
				answers += len(a.Items)
			}
		}
		if observed == 0 || id == olympiad && lessons == 0 {
			t.Errorf("%s: %d attempts observed, %d lessons of several tasks", id, observed, lessons)
		}
	}
}

// The c of an item: 1/k for a choice of k options, the product over its cells, and 0 for any cell
// that is not a choice.
func TestGuessOfItem(t *testing.T) {
	choice := func(k int32) learnermodel.Slot {
		return learnermodel.Slot{Check: learnermodel.CheckKindChoice, Options: k, HasOptions: true}
	}
	if c := guess([]learnermodel.Slot{choice(4), choice(5)}); c != float64(1.0/20) {
		t.Errorf("two choices: c = %v", c)
	}
	if c := guess([]learnermodel.Slot{choice(4), {Check: learnermodel.CheckKindNumber}}); c != 0 {
		t.Errorf("a choice and a number: c = %v", c)
	}
}

// A task of a module the learner does not have is read by no P, and its module has no skills to
// read; the model refuses to read a module it has no graph of.
func TestObserveRefuses(t *testing.T) {
	mods := modules(t, olympiad)
	l, err := learner.New(seed.Working, "G0", 0, mods, learner.Base())
	if err != nil {
		t.Fatal(err)
	}
	r, err := New(l, mods, Maths)
	if err != nil {
		t.Fatal(err)
	}
	for range 3 {
		if err := r.Step(); err != nil {
			t.Fatal(err)
		}
	}
	other := modules(t, grade3)[0]
	if _, err := r.observeBefore(other, &other.Tasks[0], int64(r.cal.Day())*86400, false, nil); !errors.Is(err, ErrRun) {
		t.Errorf("a task of another module: %v", err)
	}
	if ss, err := r.skillsOf(learnermodel.Readout{}, other, 0); ss != nil || err != nil {
		t.Errorf("the skills of a module with no readout: %v, %v", ss, err)
	}
	r.modules = append(r.modules, other)
	if _, err := r.readAt(0); !errors.Is(err, ErrRun) {
		t.Errorf("a read of a module with no graph: %v", err)
	}
}
