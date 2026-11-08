package seed

import (
	"strconv"
	"testing"
)

func TestHeldOut(t *testing.T) {
	cases := []struct {
		criterion string
		n         int
		want      Set
	}{
		{CriterionModel, 1, Set{Seed: 20261108, Experiment: "held-out-model-1"}},
		{CriterionEngine, 1, Set{Seed: 20261109, Experiment: "held-out-engine-1"}},
	}
	for _, c := range cases {
		got, err := HeldOut(c.criterion, c.n)
		if err != nil || got != c.want {
			t.Errorf("HeldOut(%q, %d) = %v, %v; want %v", c.criterion, c.n, got, err, c.want)
		}
	}
	for _, c := range []struct {
		criterion string
		n         int
	}{{CriterionModel, 2}, {CriterionPair, 1}} {
		if _, err := HeldOut(c.criterion, c.n); err != ErrNotRecorded {
			t.Errorf("HeldOut(%q, %d): %v, want ErrNotRecorded", c.criterion, c.n, err)
		}
	}
}

func TestUnits(t *testing.T) {
	if got := Learner("G2-half", 17); got != "G2-half/17" {
		t.Errorf("Learner = %q", got)
	}
	if got := Task("novatrail/olympiad-math", "enum-34-d1-1"); got != "task/novatrail/olympiad-math/enum-34-d1-1" {
		t.Errorf("Task = %q", got)
	}
}

// No learner is drawn alike in two sets (docs/architecture.md, section 3.8): the low word of a
// place is its set's seed, so places of different sets never share a state, and their first draws
// part. Checked across the recorded sets for the first 300 learners of one population of each
// family of section 3.2, every purpose of a learner, and the miss of tasks of each module; the
// test of every population comes with the populations (T16–T19 in RUN.md).
func TestSetsApart(t *testing.T) {
	sets := []Set{Working}
	for _, c := range []string{CriterionModel, CriterionEngine} {
		s, err := HeldOut(c, 1)
		if err != nil {
			t.Fatal(err)
		}
		sets = append(sets, s)
	}
	type place struct{ unit, purpose string }
	var places []place
	purposes := []string{Params, Calendar, Answers, Retries, Attempts, Attempt, Choice}
	for _, population := range []string{"G0", "G2-half", "S-adult", "F-model", "W-shock", "A-linked", "L-explain", "E-base"} {
		for n := range 300 {
			for _, p := range purposes {
				places = append(places, place{Learner(population, n), p})
			}
		}
	}
	for _, module := range []string{"novatrail/olympiad-math", "novatrail/grade-3-math", "novatrail/english-essays"} {
		for n := range 30 {
			places = append(places, place{Task(module, "task-"+strconv.Itoa(n)), Miss})
		}
	}
	seen := map[uint64]string{}
	for _, s := range sets {
		for _, pl := range places {
			st, err := s.Place(pl.unit, pl.purpose)
			if err != nil {
				t.Fatal(err)
			}
			first := st.Integer()
			where := s.Experiment + " " + pl.unit + " " + pl.purpose
			if prior, ok := seen[first]; ok {
				t.Fatalf("%s draws first what %s does", where, prior)
			}
			seen[first] = where
		}
	}
}

func TestPlaceRefuses(t *testing.T) {
	for _, names := range [][]string{nil, {""}, {"a|b"}, {"ok", ""}} {
		if _, err := Place(1, names...); err != ErrName {
			t.Errorf("Place(1, %q): %v, want ErrName", names, err)
		}
	}
	if _, err := Working.Place("a|b", Params); err != ErrName {
		t.Errorf("Set.Place with |: %v, want ErrName", err)
	}
}

func TestBetweenAndPermutation(t *testing.T) {
	s := New(0, 2)
	for range 100 {
		if x := s.Between(-3, 3); x < -3 || x > 3 {
			t.Fatalf("Between(-3, 3) = %d", x)
		}
	}
	p := s.Permutation(5)
	seen := map[int]bool{}
	for _, x := range p {
		seen[x] = true
	}
	if len(p) != 5 || len(seen) != 5 {
		t.Fatalf("Permutation(5) = %v", p)
	}
}

// Named gives a set by its experiment: working, or a held-out set R04 has recorded.
func TestNamed(t *testing.T) {
	if s, err := Named("working"); err != nil || s != Working {
		t.Errorf("Named(working) = %v, %v", s, err)
	}
	if s, err := Named("held-out-engine-1"); err != nil || s != (Set{Seed: 20261109, Experiment: "held-out-engine-1"}) {
		t.Errorf("Named(held-out-engine-1) = %v, %v", s, err)
	}
	if _, err := Named("held-out-pair-1"); err != ErrNotRecorded {
		t.Errorf("Named(held-out-pair-1): %v", err)
	}
}
