// Package learner is the virtual learner of the bench: its truth, which the bench knows and the
// learner model estimates from its attempts, and its calendar (docs/architecture.md, section 3.1;
// R03). The truth is the bench's own model, not learner-model's code: the bench measures the model
// on a learner it knows, not on a copy of itself. The base population holds the model's
// assumptions (R03, R04, R06 learner-model), and each population changes one of them
// (docs/architecture.md, section 3.2), by a field of Assumptions.
package learner

// Assumptions are what a population assumes of its learners. Base gives the base population's; a
// population is Base with one field changed (R15 in docs/decisions.md).
type Assumptions struct {
	// BandStep is the true step between neighbouring bands, in logits: the shift 2.5 · (p − 1)
	// of a start band and of a task's band (R03 learner-model).
	BandStep float64
	// StartSpread is the spread of the module's level around the start band, N(0, σ²).
	StartSpread float64
	// SkillSpread is the spread of a skill's offset from the module's level, N(0, σ²).
	SkillSpread float64
	// Miss multiplies a task's z, drawn once a task and a run: the author's miss of its difficulty
	// (docs/architecture.md, section 3.3).
	Miss float64
	// Slope is a of σ(a · (θ* + δ̄* − β*)).
	Slope float64
	// ChoiceGuess is the chance of guessing a choice of k options, times k: 1 gives 1/k, as the
	// model has it (R03 learner-model); it is capped at 1.
	ChoiceGuess float64
	// PredicateGuess is the chance of guessing a cell of predicate; the model has 0.
	PredicateGuess float64
	// RungLift is what each rung of the ladder of hints adds to the logit of the items not yet
	// right (concepts 10 spec).
	RungLift float64
	// Forgets tells whether retention falls with time; false keeps it at 1.
	Forgets bool
	// LawC and LawF are C and F of the law of retention R* = (1 + F · Δt / τ*)^(−C), Δt and τ* in
	// days (R04 learner-model).
	LawC, LawF float64
	// Growth is g, by which practice spaced out grows the stability (R06 learner-model).
	Growth float64
	// SessionChance is the chance that a day of the calendar has a session.
	SessionChance float64
	// SessionAttempts is the number of attempts of a session.
	SessionAttempts int
}

// Base gives the assumptions of the base population: the model's own, but the author's miss
// (docs/architecture.md, section 3.1).
func Base() Assumptions {
	return Assumptions{
		BandStep:        2.5,
		StartSpread:     1,
		SkillSpread:     0.5,
		Miss:            0.5,
		Slope:           1,
		ChoiceGuess:     1,
		PredicateGuess:  0,
		RungLift:        0.5,
		Forgets:         true,
		LawC:            0.5,
		LawF:            19.0 / 81,
		Growth:          15,
		SessionChance:   5.0 / 7,
		SessionAttempts: 10,
	}
}
