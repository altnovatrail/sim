// Package population names the populations of virtual learners: each changes one assumption of
// the base population, on the graphs it says, and runs to its limit (docs/architecture.md, section
// 3.2; R03; R18 in docs/decisions.md). Two sides of one assumption take half the learners each:
// the even ones the first side, the odd ones the second.
package population

import (
	"slices"

	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/run"
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

// All gives the populations the bench has, a copy, so that no caller changes them for another.
func All() []Population {
	return slices.Clone(Standalone)
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
