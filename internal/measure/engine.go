package measure

import (
	"maps"
	"slices"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/run"
)

// The thresholds of the measures of the engine (docs/architecture.md, section 3.5).
const (
	stuckChoices = 50   // E2: choices in a row on one skill
	hardP        = 0.50 // E1: P* below it is too hard
	easyP        = 0.95 // E1: P* above it is too easy
	week         = 7 * day
)

// weeks are the weeks E3 reads the share of repetitions in: those of the longest limit.
var weeks = (max(run.Maths.Days, run.Forgetting.Days, run.Essays.Days) + 6) / 7

// Engine gives the sums of one learner's run for the measures of the engine, E1–E5
// (docs/architecture.md, section 3.5; R23 in docs/decisions.md).
func Engine(rec *run.Record) Sums {
	s := Sums{"learners": 1}
	e1(s, rec)
	e2(s, rec)
	e3(s, rec)
	e4(s, rec)
	e5(s, rec)
	return s
}

// gives tells whether a step gave something: a task or a lesson, but for a lesson going on after a
// consolidation, which is no new choice (R48).
func gives(st run.Step) bool {
	return !st.Resumed && (st.Choice.Kind == engine.ChoiceTask || st.Choice.Kind == engine.ChoiceLesson)
}

// e1: of the items with a verdict given, the shares with P* at the time given in the corridor
// 0.70–0.85, below 0.50 and above 0.95; and the ceiling of the pool, the mean over the choices that
// gave something of the largest share in the corridor of one task the rule could give.
func e1(s Sums, rec *run.Record) {
	for _, st := range rec.Steps {
		for _, o := range st.Obs {
			if o == nil {
				continue
			}
			shares(s, "E1", o.Items)
		}
		if gives(st) && st.Pool != nil && st.Pool.Given {
			s.add("E1.ceiling.sum", st.Pool.Ceiling)
			s.add("E1.ceiling.n", 1)
		}
	}
}

// shares adds items with a verdict to the sums of a measure's shares by P* at the time given: all of
// them, and those in the corridor, below 0.50 and above 0.95. E1 reads all the items so, and the
// note E6 those of the lessons and of the consolidations (R49).
func shares(s Sums, of string, items []run.ItemObs) {
	for _, it := range items {
		s.add(of+".items", 1)
		switch {
		case it.PGiven >= run.CorridorLow && it.PGiven <= run.CorridorHigh:
			s.add(of+".corridor", 1)
		case it.PGiven < hardP:
			s.add(of+".hard", 1)
		case it.PGiven > easyP:
			s.add(of+".easy", 1)
		}
	}
}

// mastered tells whether the model reads a skill mastered in an input.
func mastered(in engine.Input, id string) bool {
	for _, m := range in.Readout.Modules {
		for _, sk := range m.Skills {
			if sk.Skill == id {
				return sk.Mastered
			}
		}
	}
	return false
}

// e2: whether the learner has 50 choices in a row on one skill with no change of its "mastered"
// while the pool holds another skill within reach; and the choices of nothing while the pool has a
// task the rule could give, of the choices with such a pool.
func e2(s Sums, rec *run.Record) {
	type streak struct {
		n        int
		mastered bool
	}
	streaks := map[string]streak{}
	stuck := false
	for _, st := range rec.Steps {
		if st.Resumed {
			continue // a lesson going on after a consolidation is no new choice (R48)
		}
		if st.Pool != nil && st.Pool.Given {
			s.add("E2.given", 1)
			if st.Choice.Kind == engine.ChoiceNothing {
				s.add("E2.nothing", 1)
			}
		}
		on := st.Skills()
		next := map[string]streak{}
		for _, id := range on {
			// Another skill within reach is one the choice is not on.
			other := st.Pool != nil && slices.ContainsFunc(st.Pool.Reachable, func(r string) bool { return !slices.Contains(on, r) })
			if !other {
				continue
			}
			m := mastered(st.Input, id)
			k, ok := streaks[id]
			if !ok || k.mastered != m {
				k = streak{mastered: m}
			}
			k.n++
			next[id] = k
			stuck = stuck || k.n >= stuckChoices
		}
		streaks = next
	}
	if stuck {
		s.add("E2.stuck", 1)
	}
}

// e3: of the choices by repetition, R* of the skill due for review when it is given, and whether
// P* of the reference task of its band of mastery is already below 0.775 — late; the share of
// repetitions among the choices of each week of the calendar; and the mean R* of the skills the
// truth has mastered on the last day of the run.
func e3(s Sums, rec *run.Record) {
	for _, st := range rec.Steps {
		if !gives(st) {
			continue
		}
		w := itoa(int((st.Input.Time-learner.Start)/week) + 1)
		s.add("E3.week"+w+".choices", 1)
		if st.Choice.Reason != engine.ReasonRepetition {
			continue
		}
		s.add("E3.week"+w+".reviews", 1)
		if len(st.Obs) == 0 || st.Obs[0] == nil {
			continue
		}
		on := st.Skills()
		for _, b := range st.Obs[0].Before {
			if !slices.Contains(on, b.Skill) || !due(st.Input, b.Skill) {
				continue
			}
			s.add("E3.retention.sum", b.Retention)
			s.add("E3.reviews", 1)
			if b.PTrueRef < learner.Mastery {
				s.add("E3.late", 1)
			}
		}
	}
	l, end, skills := endOf(rec)
	for _, id := range skills {
		if l.TrulyMasteredAt(id, l.TaughtOn(id), end) {
			_, r, _ := l.Truth(id, end)
			s.add("E3.end.sum", r)
			s.add("E3.end.n", 1)
		}
	}
}

// due tells whether the model reads a skill mastered with its review date come at an input.
func due(in engine.Input, id string) bool {
	for _, m := range in.Readout.Modules {
		for _, sk := range m.Skills {
			if sk.Skill == id {
				return sk.Mastered && (sk.Review.Already || sk.Review.At <= in.Time)
			}
		}
	}
	return false
}

// endOf gives the learner of a run, the time it ended — its last choice, or the close of the last
// attempt that choice gave — and the skills with a verdict of its modules, by full ID; no skills for
// a run with no learner or no step.
func endOf(rec *run.Record) (*learner.Learner, int64, []string) {
	if rec.Learner == nil || len(rec.Steps) == 0 {
		return nil, 0, nil
	}
	st := rec.Steps[len(rec.Steps)-1]
	last := st.Input
	end := last.Time
	for _, a := range st.Attempts {
		end = max(end, a.End)
	}
	var skills []string
	for _, m := range last.Modules {
		for _, sk := range m.Graph.Skills {
			if sk.Assessability == learnermodel.AssessabilityVerdict {
				skills = append(skills, sk.ID)
			}
		}
	}
	return rec.Learner, end, skills
}

// e4: of the skills within reach at some choice, the share given at least once; and of the skills
// with a verdict, the share the truth has mastered by the end of the run.
func e4(s Sums, rec *run.Record) {
	reach, given := map[string]bool{}, map[string]bool{}
	for _, st := range rec.Steps {
		if st.Pool != nil {
			for _, id := range st.Pool.Reachable {
				reach[id] = true
			}
		}
		if gives(st) {
			for _, id := range st.Skills() {
				given[id] = true
			}
		}
	}
	for _, id := range slices.Sorted(maps.Keys(reach)) {
		s.add("E4.reach", 1)
		if given[id] {
			s.add("E4.given", 1)
		}
	}
	l, end, skills := endOf(rec)
	for _, id := range skills {
		s.add("E4.skills", 1)
		if l.TrulyMasteredAt(id, l.TaughtOn(id), end) {
			s.add("E4.mastered", 1)
		}
	}
}

// e5: the tasks with a rubric the learner wrote by the end of the run, each once, and the times it
// was given a task it had closed by "Done" — as a task or in a lesson.
func e5(s Sums, rec *run.Record) {
	written := map[string]bool{}
	for _, st := range rec.Steps {
		for _, a := range st.Attempts {
			if slices.ContainsFunc(st.Input.RubricDone, func(d engine.RubricDone) bool { return d.Task == a.Task }) {
				s.add("E5.again", 1)
			}
			if a.Closure == learner.ClosedDone && !written[a.Task] {
				written[a.Task] = true
				s.add("E5.written", 1)
			}
		}
	}
}

// EngineNumbers are the numbers of the measures of the engine a cell is read for.
func EngineNumbers() []Number {
	out := []Number{
		{"E1.corridor", ratio("E1.corridor", "E1.items")},
		{"E1.hard", ratio("E1.hard", "E1.items")},
		{"E1.easy", ratio("E1.easy", "E1.items")},
		{"E1.ceiling", ratio("E1.ceiling.sum", "E1.ceiling.n")},
		{"E2.stuck", ratio("E2.stuck", "learners")},
		{"E2.nothing", ratio("E2.nothing", "E2.given")},
		{"E3.retention", ratio("E3.retention.sum", "E3.reviews")},
		{"E3.late", ratio("E3.late", "E3.reviews")},
		{"E3.end", ratio("E3.end.sum", "E3.end.n")},
	}
	for w := 1; w <= weeks; w++ {
		out = append(out, Number{"E3.week" + itoa(w), ratio("E3.week"+itoa(w)+".reviews", "E3.week"+itoa(w)+".choices")})
	}
	return append(out,
		Number{"E4.given", ratio("E4.given", "E4.reach")},
		Number{"E4.mastered", ratio("E4.mastered", "E4.skills")},
		Number{"E5.written", ratio("E5.written", "learners")},
		Number{"E5.again", ratio("E5.again", "learners")},
	)
}
