package run

import (
	"maps"
	"slices"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
)

// The corridor of success E1 reads P* in, and the pause before a task with a seen solution
// (docs/architecture.md, section 3.5; R82 spec; R03 engine).
const (
	CorridorLow  = 0.70
	CorridorHigh = 0.85
	pause        = 24 * 3600
)

// Pool is what the truth says of the pool at a choice, before the engine chooses (R23): what any
// rule could have given, so that the measures of the engine read the choice against it.
type Pool struct {
	// Ceiling is the largest share of items with a verdict whose P* is in the corridor, of one task
	// of the pool the rule could give: no rule that gives whole tasks gets above it (R104 spec).
	Ceiling float64
	// Given tells that the pool has a task the rule could give.
	Given bool
	// Reachable are the skills within reach by the truth, by full ID: skills with a verdict whose
	// prerequisites the truth holds, with a task of the pool the rule could give on them.
	Reachable []string
}

// Skills gives the skills a choice is on: a lesson's skills, or a task's primary skills as
// learner-model and the engine count them — the items' own, and the task's when an item takes them
// or the task has no items.
func (s Step) Skills() []string {
	for _, m := range s.Input.Modules {
		switch s.Choice.Kind {
		case engine.ChoiceLesson:
			for _, l := range m.Lessons {
				if l.ID == s.Choice.ID {
					return l.Skills
				}
			}
		case engine.ChoiceTask:
			for _, c := range m.Pool {
				if c.Task.ID == s.Choice.ID {
					var out []string
					takesTask := len(c.Task.Items) == 0
					for _, it := range c.Task.Items {
						if it.OwnSkills {
							out = append(out, it.Primary...)
						} else {
							takesTask = true
						}
					}
					if takesTask {
						out = append(out, c.Task.Primary...)
					}
					slices.Sort(out)
					return slices.Compact(out)
				}
			}
		}
	}
	return nil
}

// given tells whether a task of the pool is one the rule may give at a time in a pass: a task, not
// a template, whose solution went to the device less than 24 hours ago is not in the first pass,
// and in the second it is, but for the task of the last closed attempt (R82 spec; R03 engine and
// its amendment of 2026-11-07).
func given(t learnermodel.Task, at int64, solved func(string) (int64, bool), last string, second bool) bool {
	if t.Template {
		return true
	}
	sent, seen := solved(t.ID)
	if !seen || at-sent >= pause {
		return true
	}
	return second && t.ID != last
}

// pool reads the pool at a time: P* of the items of each task the rule could give — in the first
// pass, or in the second when the first has none —, and the skills within reach.
func (r *Run) pool(at int64) (*Pool, error) {
	last := ""
	if r.hasLast {
		last = r.last.Task
	}
	open := map[string]bool{}
	isOpen := func(m *pack.Module, id string) bool {
		if v, ok := open[id]; ok {
			return v
		}
		open[id] = r.holds(m, id, at)
		return open[id]
	}
	// The tasks within reach, read once: the share of their items with a verdict whose P* is in the
	// corridor, and their primary skills with a verdict.
	type within struct {
		task  learnermodel.Task
		share float64
		on    []string
	}
	var tasks []within
	for i, m := range r.modules {
		for _, c := range r.suited[i].Pool {
			t, err := m.Task(c.Task.ID)
			if err != nil {
				return nil, err
			}
			if t.Rubric {
				continue
			}
			var on []string
			in, all := 0, 0
			for k, it := range c.Task.Items {
				if len(it.Slots) == 0 {
					continue
				}
				ps, err := r.l.Chance(t, k, at, 0)
				if err != nil {
					return nil, err
				}
				all++
				in += b2i(ps >= CorridorLow && ps <= CorridorHigh)
				for _, id := range t.ItemPrimary(k) {
					if verdict(m, id) {
						on = append(on, id)
					}
				}
			}
			if all == 0 || len(on) == 0 || slices.ContainsFunc(on, func(id string) bool { return !isOpen(m, id) }) {
				continue // no item with a verdict on a skill with one, or a skill out of reach by the truth
			}
			tasks = append(tasks, within{task: c.Task, share: float64(float64(in) / float64(all)), on: on})
		}
	}
	for _, second := range []bool{false, true} {
		p := &Pool{}
		reach := map[string]bool{}
		for _, w := range tasks {
			if !given(w.task, at, r.l.SolvedAt, last, second) {
				continue
			}
			p.Given = true
			p.Ceiling = max(p.Ceiling, w.share)
			for _, id := range w.on {
				reach[id] = true
			}
		}
		if p.Given || second {
			p.Reachable = slices.Sorted(maps.Keys(reach))
			return p, nil
		}
	}
	return nil, nil // not reached: the second pass returns
}

// holds tells whether the truth holds the prerequisites of a skill of a module at a time: each is
// a skill with a verdict of the learner the truth has mastered on its first band taught; a skill
// the truth has no verdict of, or of a module the learner does not learn, is held — nothing reads
// it; an anchor is held when a skill of the learner's modules that teaches it is truly mastered
// with its R* then, or when none teaches it (docs/architecture.md, section 3.5, M9).
func (r *Run) holds(m *pack.Module, id string, at int64) bool {
	i := slices.IndexFunc(m.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == id })
	if i < 0 {
		return true
	}
	for _, pre := range m.Graph.Skills[i].Prerequisites {
		if _, _, ok := r.l.Truth(pre, at); ok {
			if !r.l.TrulyMasteredAt(pre, r.l.TaughtOn(pre), at) {
				return false
			}
			continue
		}
		taught, held := false, false
		for _, om := range r.modules {
			for _, s := range om.Graph.Skills {
				if slices.Contains(s.Teaches, pre) {
					taught = true
					held = held || r.l.ReferenceAt(s.ID, r.l.TaughtOn(s.ID), at) >= learner.Mastery
				}
			}
		}
		if taught && !held {
			return false
		}
	}
	return true
}
