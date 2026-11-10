// Package population names the populations of virtual learners: each changes one assumption of
// the base population, on the graphs it says, and runs to its limit (docs/architecture.md, section
// 3.2; R03; R18 in docs/decisions.md). Two sides of one assumption take half the learners each:
// the even ones the first side, the odd ones the second.
package population

import (
	"slices"

	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// The modules of the first slice, by full ID (R41 spec).
const (
	Olympiad = "novatrail/olympiad-math"
	Grade3   = "novatrail/grade-3-math"
	Essays   = "novatrail/english-essays"
)

// Population is a population of virtual learners.
type Population struct {
	Name string
	// Graphs are the graphs it runs on, each the modules of one learner: a learner of olympiad maths
	// and one of school maths are two graphs; a learner of both is one graph of two modules.
	Graphs [][]string
	Limit  run.Limit
	// sides change the base population's assumptions, the first for the even learners and the
	// second, when there is one, for the odd.
	sides []func(*learner.Assumptions)
}

// Assumptions gives the assumptions of the n-th learner of the population: the base population's,
// changed by its side.
func (p Population) Assumptions(n int) learner.Assumptions {
	a := learner.Base()
	if len(p.sides) > 0 {
		p.sides[p.Side(n)](&a)
	}
	return a
}

// Side gives the side of the n-th learner, from 0: the learners take the sides in turn.
func (p Population) Side(n int) int {
	return n % max(1, len(p.sides))
}

var (
	maths    = [][]string{{Olympiad}, {Grade3}}
	olympiad = [][]string{{Olympiad}}
	grade3   = [][]string{{Grade3}}
	essays   = [][]string{{Essays}}
	// anchored is the graph of A-: a learner of a module whose skill rests on an anchor skill and of
	// the module that teaches it. Olympiad maths, from 0.1.1 on, rests two skills on skills school maths
	// teaches — fractions-parts on core.fractions-of-a-whole, number-divisibility on
	// core.division-with-remainder —, so it is the two published maths together
	// (docs/architecture.md, section 3.2; R31 in docs/decisions.md).
	anchored = [][]string{{Olympiad, Grade3}}
)

// one gives the sides of a population of one side.
func one(f func(*learner.Assumptions)) []func(*learner.Assumptions) {
	return []func(*learner.Assumptions){f}
}

// two gives the sides of a population of two.
func two(f, g func(*learner.Assumptions)) []func(*learner.Assumptions) {
	return []func(*learner.Assumptions){f, g}
}

// Standalone are the nineteen populations of standalone, carried over as assumptions and not as
// code (docs/architecture.md, section 3.2; R03). Their graphs are both maths unless said; G7 runs
// on school maths for grade 3, the one published package with a choice (the amendment of
// 2026-11-08 to R03).
var Standalone = []Population{
	{Name: "G0", Graphs: maths, Limit: run.Maths},
	{Name: "G1", Graphs: maths, Limit: run.Maths, sides: two(
		func(a *learner.Assumptions) { a.LevelShift = -2.5 },
		func(a *learner.Assumptions) { a.LevelShift = 2.5 },
	)},
	{Name: "G2", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.LearnLevel, a.LearnSkill, a.LearnReturn = 0.01, 0.02, 0.01
	})},
	{Name: "G2-half", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.LearnLevel, a.LearnSkill, a.LearnReturn = 0.005, 0.01, 0.01
	})},
	{Name: "G2-fading", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.LearnLevel, a.LearnSkill, a.LearnReturn, a.LearnFade = 0.01, 0.02, 0.01, 50
	})},
	{Name: "G3", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.JumpBy, a.JumpFrom, a.JumpTo = 1, 50, 150
	})},
	{Name: "G3-drop", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.JumpBy, a.JumpFrom, a.JumpTo = -1, 50, 150
	})},
	{Name: "G4-band", Graphs: olympiad, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.BandBias = 0.75 })},
	{Name: "G5", Graphs: maths, Limit: run.Maths, sides: two(
		func(a *learner.Assumptions) { a.SkillHalves, a.SkillSpread = 0.75, 0.25 },
		func(a *learner.Assumptions) { a.SkillHalves, a.SkillSpread = -0.75, 0.25 },
	)},
	{Name: "G6", Graphs: maths, Limit: run.Maths, sides: two(
		func(a *learner.Assumptions) { a.Slope = 0.5 },
		func(a *learner.Assumptions) { a.Slope = 2 },
	)},
	{Name: "G7", Graphs: grade3, Limit: run.Maths, sides: two(
		func(a *learner.Assumptions) { a.ChoiceGuess = 0 },
		func(a *learner.Assumptions) { a.ChoiceGuess = 1.5 },
	)},
	{Name: "G8", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.TechniqueChance, a.TechniqueLift = 0.2, 0.5
	})},
	{Name: "G0-exact", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.Miss = 0 })},
	{Name: "G0-miss0.25", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.Miss = 0.25 })},
	{Name: "G0-miss1", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.Miss = 1 })},
	{Name: "G0-start0.5", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.StartSpread = 0.5 })},
	{Name: "G0-start2", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.StartSpread = 2 })},
	{Name: "G0-skills0.3", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.SkillSpread = 0.3 })},
	{Name: "G0-skills1", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.SkillSpread = 1 })},
}

// breaks gives the calendar of F-model: a day is a session with the chance 3/7, and one break of
// 14–60 days begins on a day from the 30th to the 60th and ends by the 120th of 180
// (docs/architecture.md, section 3.2).
func breaks(a *learner.Assumptions) {
	a.SessionChance = 3.0 / 7
	a.BreakStartFrom, a.BreakStartTo, a.BreakMin, a.BreakMax, a.BreakEndBy = 30, 60, 14, 60, 120
}

// Own are the populations of NovaTrail's own model and engine: the scale of bands, guessing of a
// predicate, an adult without a grade, an unreadable answer (S-); forgetting by other laws, on a
// calendar with a break (F-); pages of a workbook (W-); the anchor skill (A-); lessons (L-);
// essays (E-) (docs/architecture.md, section 3.2). S-predicate runs on the synthetic sheet alone: no published
// package has a predicate, and in the criterion and the gate its checks are "no graph" (R04).
var Own = []Population{
	{Name: "S-scale1.5", Graphs: olympiad, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.BandStep = 1.5 })},
	{Name: "S-scale3.5", Graphs: olympiad, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.BandStep = 3.5 })},
	{Name: "S-predicate", Graphs: [][]string{{synthetic.Sheet}}, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.PredicateGuess = 0.1 })},
	{Name: "S-adult", Graphs: olympiad, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.Adult, a.AdultFrom, a.AdultTo = true, 2.5, 5
	})},
	{Name: "S-unreadable", Graphs: olympiad, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.UnreadableEvery = 20 })},
	{Name: "F-model", Graphs: maths, Limit: run.Forgetting, sides: one(breaks)},
	{Name: "F-none", Graphs: maths, Limit: run.Forgetting, sides: one(func(a *learner.Assumptions) { breaks(a); a.Forgets = false })},
	{Name: "F-exp", Graphs: maths, Limit: run.Forgetting, sides: one(func(a *learner.Assumptions) {
		breaks(a)
		a.ExponentialLaw, a.LawBase = true, 0.9
	})},
	{Name: "F-C1", Graphs: maths, Limit: run.Forgetting, sides: one(func(a *learner.Assumptions) { breaks(a); a.LawC = 1 })},
	{Name: "F-growth5", Graphs: maths, Limit: run.Forgetting, sides: one(func(a *learner.Assumptions) { breaks(a); a.Growth = 5 })},
	{Name: "F-break", Graphs: maths, Limit: run.Forgetting, sides: one(func(a *learner.Assumptions) { breaks(a); a.BreakDrop = 1 })},
	{Name: "W-shock", Graphs: grade3, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.AttemptShock = 1 })},
	{Name: "W-error", Graphs: grade3, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.ErrorChance = 0.2 })},
	{Name: "W-badday", Graphs: grade3, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.BadDayChance, a.BadDayDrop = 0.1, 1
	})},
	{Name: "W-repeat", Graphs: maths, Limit: run.Maths, sides: one(func(a *learner.Assumptions) {
		a.RepeatChance, a.RepeatDays = 0.95, 30
	})},
	{Name: "A-linked", Graphs: anchored, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.LevelCorrelation = 0.7 })},
	{Name: "A-unpractised", Graphs: anchored, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.TeacherBoost = 2 })},
	{Name: "A-fading", Graphs: anchored, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.TeacherFades = true })},
	{Name: "L-abandon", Graphs: olympiad, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.AbandonChance = 0.2 })},
	{Name: "L-explain", Graphs: olympiad, Limit: run.Maths, sides: one(func(a *learner.Assumptions) { a.ExplainLift = 0.25 })},
	{Name: "E-base", Graphs: essays, Limit: run.Essays},
	{Name: "E-leave", Graphs: essays, Limit: run.Essays, sides: one(func(a *learner.Assumptions) { a.LeaveChance = 1.0 / 3 })},
}

// All gives the populations the bench has, a copy, so that no caller changes them for another.
func All() []Population {
	return slices.Concat(Standalone, Own)
}

// Named gives a population by its name.
func Named(name string) (Population, bool) {
	for _, p := range All() {
		if p.Name == name {
			return p, true
		}
	}
	return Population{}, false
}
