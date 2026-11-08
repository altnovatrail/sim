package seed

import (
	"errors"
	"strconv"
)

// Set is a set of virtual learners: a seed and the name of its experiment, the first name of every
// place it draws from (docs/architecture.md, sections 3.4 and 3.8). Sets that must stay apart take
// seeds of their own, as the format bids ("a use that needs streams apart from another's takes a
// seed of its own", the text of the format, section 22).
type Set struct {
	Seed       uint64
	Experiment string
}

// Working is the set of working learners: the gate and the decisive runs take the same ones.
var Working = Set{Seed: 20261107, Experiment: "working"}

// heldOut are the held-out sets R04 has recorded, by the name of their experiment. The seed of
// each next choice of a criterion, and of each pair, is recorded by an amendment to R04 before
// its candidate first runs (docs/architecture.md, section 3.8), and only then comes in here.
var heldOut = map[string]uint64{
	"held-out-model-1":  20261108,
	"held-out-engine-1": 20261109,
}

// ErrNotRecorded is a held-out set R04 has not recorded yet.
var ErrNotRecorded = errors.New("seed: no such held-out set is recorded in R04")

// The criteria whose choices take held-out sets of their own (docs/architecture.md, section 3.8).
const (
	CriterionModel  = "model"
	CriterionEngine = "engine"
	CriterionPair   = "pair"
)

// HeldOut gives the held-out set of the n-th choice of a criterion: experiment
// held-out-<criterion>-<n>, at the seed R04 recorded for it.
func HeldOut(criterion string, n int) (Set, error) {
	experiment := "held-out-" + criterion + "-" + strconv.Itoa(n)
	s, ok := heldOut[experiment]
	if !ok {
		return Set{}, ErrNotRecorded
	}
	return Set{Seed: s, Experiment: experiment}, nil
}

// The purposes of the places of a virtual learner and of a task (docs/architecture.md, section
// 3.4).
const (
	Params   = "params"   // the learner's truth
	Calendar = "calendar" // days, sessions and breaks
	Answers  = "answers"  // the first answers of items
	Retries  = "retries"  // the answers after a rung of the ladder or after unreadable
	Attempts = "attempts" // what each attempt holds: shocks, errors, techniques, abandons
	Attempt  = "attempt"  // the seed of each attempt, for the instances of templates
	Choice   = "choice"   // the seed of a choice of engine, if its rule ever draws at random
	Miss     = "miss"     // the author's miss of a task's difficulty
)

// Learner gives the unit of the n-th virtual learner of a population: <population>/<n>.
func Learner(population string, n int) string {
	return population + "/" + strconv.Itoa(n)
}

// Task gives the unit of a task of a module: task/<module>/<task>.
func Task(module, task string) string {
	return "task/" + module + "/" + task
}

// Place gives the stream of the place of a unit and a purpose in the set: names experiment,
// unit and purpose, as standalone's seeded(unit, purpose) is (docs/architecture.md, section 3.4).
func (s Set) Place(unit, purpose string) (*Stream, error) {
	return Place(s.Seed, s.Experiment, unit, purpose)
}
