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

	// LevelShift moves the level of every module from the start band's: G1's 2.5 above or below.
	LevelShift float64
	// LearnLevel, LearnSkill and LearnReturn are what each first answer of an item with a verdict
	// teaches (G2): the module's level grows by LearnLevel, the offsets of the item's primary skills
	// by LearnSkill more, and the offset of every other skill of the module goes back LearnReturn of
	// the way to where it started. LearnFade, when above 0, slows them as F/(F + k) at the k-th
	// answer (G2-fading).
	LearnLevel, LearnSkill, LearnReturn, LearnFade float64
	// JumpBy moves the level of the module once, at an answer drawn from JumpFrom to JumpTo (G3,
	// G3-drop).
	JumpBy           float64
	JumpFrom, JumpTo int
	// BandBias makes the tasks of one band of each module, drawn for the learner, truly harder by
	// that much than the author named, beyond the miss (G4-band).
	BandBias float64
	// SkillHalves makes the first half of the skills of a module's graph stronger by that much and
	// the second half weaker by as much (G5); a negative value turns it round.
	SkillHalves float64
	// TechniqueChance is the chance that an attempt takes the technique on request before its first
	// answer, one fraction of the place attempts each attempt, and TechniqueLift what that adds to
	// the logit of its items (G8).
	TechniqueChance, TechniqueLift float64

	// Adult makes a learner an adult without a grade (S-adult): the first band to start in, a level
	// AdultFrom to AdultTo above it, drawn evenly, and the age category adult.
	Adult              bool
	AdultFrom, AdultTo float64
	// UnreadableEvery makes every UnreadableEvery-th answer of an item on a sheet unreadable first,
	// written again (S-unreadable): no verdict and no mistake (R04 spec).
	UnreadableEvery int
	// The break of the calendar (F-): one break, starting on a day from BreakStartFrom to
	// BreakStartTo, BreakMin to BreakMax days long, drawn evenly, ending by day BreakEndBy; and
	// BreakDrop, by which the level of every module falls after it (F-break).
	BreakStartFrom, BreakStartTo, BreakMin, BreakMax, BreakEndBy int
	BreakDrop                                                    float64
	// ExponentialLaw turns the law of retention into R* = LawBase^(Δt/τ*) (F-exp).
	ExponentialLaw bool
	LawBase        float64
	// AttemptShock is the spread of a shift N(0, σ²) of the logit of all items of an attempt, one
	// for each attempt (W-shock).
	AttemptShock float64
	// ErrorChance is the chance that an attempt holds an error: every item with one of the task's
	// primary skills is not yet right at its first answer (W-error).
	ErrorChance float64
	// BadDayChance is the chance that a day with a session is bad, and BadDayDrop by how much the
	// level of every module falls on it (W-badday).
	BadDayChance, BadDayDrop float64
	// RepeatChance is P* of every item of a task, not a template, whose solution went to the device
	// within RepeatDays days (W-repeat).
	RepeatChance float64
	RepeatDays   int

	// LevelCorrelation ties the levels of a learner's modules: the second and later draw their
	// N(0, 1) as ρ·z₁ + √(1 − ρ²)·z (A-linked).
	LevelCorrelation float64
	// TeacherBoost makes the skills that teach an anchor skill stronger by that much: the learner
	// knows the anchor outside the platform (A-unpractised).
	TeacherBoost float64
	// TeacherFades stops a skill that teaches an anchor skill from growing its stability, g = 0, once
	// the truth has mastered it (A-fading).
	TeacherFades bool
	// AbandonChance is the chance that a pass of a lesson is abandoned at a link drawn evenly, one
	// fraction and one link of the place attempts each pass (L-abandon).
	AbandonChance float64
	// ExplainLift makes the explanations and the worked examples of a lesson teach: at its first
	// pass, the offsets of its skills grow by that much (L-explain).
	ExplainLift float64
	// LeaveChance is the chance that an attempt of a task with a rubric is left, closed "left", one
	// fraction of the place attempts each such attempt (E-leave).
	LeaveChance float64
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
