// Package invariant checks the invariants of the bench on every learner's run, with a tolerance
// of zero: one breach on one learner makes the run red (docs/architecture.md, section 3.5; R04;
// R17 in docs/decisions.md).
//
//   - I1, reproducibility: the choice taken again on the input it keeps is the same choice — each
//     call further through a lesson too (R48) —, and a run repeated on the same seeds gives the same
//     bits (Same);
//   - I2, the reason code: each choice of a task or a lesson has one code of the closed list, a
//     choice of nothing none, and the code is true to its meaning on the facts it names;
//   - I3, the author's order: a skill with a rubric or without assessment comes after the skills of
//     that kind it rests on are done;
//   - I4, a solved task again: a task, not a template, whose solution went to the device is not
//     given as the next task.
package invariant

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/run"
)

// Breach is one breach of an invariant on one step of one learner's run.
type Breach struct {
	Invariant string // I1 … I4
	Unit      string
	Step      int // from 0; -1 for the run as a whole
	Choice    engine.Choice
	Detail    string
}

func (b Breach) String() string {
	return fmt.Sprintf("%s: %s, step %d, choice %v %s (%v): %s", b.Invariant, b.Unit, b.Step, b.Choice.Kind, b.Choice.ID, b.Choice.Reason, b.Detail)
}

// Check gives every breach of I1–I4 in a run's record: I1 by taking each choice again on the input it
// keeps.
func Check(rec *run.Record) []Breach {
	return check(rec, true)
}

// CheckWithoutI1 gives every breach of I2–I4 in a run's record, and takes no choice again: past the
// first learners of a cell, taking each choice again would double the engine's work on every step,
// the most of a run. A repeated run holds I1 to the same bits by Same, and the first learners of
// every cell take it on each choice (R31, R33 in docs/decisions.md).
func CheckWithoutI1(rec *run.Record) []Breach {
	return check(rec, false)
}

func check(rec *run.Record, again bool) []Breach {
	var out []Breach
	for k, s := range rec.Steps {
		breach := func(inv, format string, args ...any) {
			out = append(out, Breach{Invariant: inv, Unit: rec.Unit, Step: k, Choice: s.Choice, Detail: fmt.Sprintf(format, args...)})
		}
		// The calls further through a lesson are choices of the engine too: taken again, and true to
		// their code (R48).
		calls := append([]run.Call{{Input: s.Input, Choice: s.Choice}}, s.Further...)
		for i, call := range calls {
			at := ""
			if i > 0 {
				at = fmt.Sprintf("call %d further through the lesson: ", i)
			}
			if again {
				c, err := engine.Choose(call.Input)
				switch {
				case err != nil:
					breach("I1", "%sthe engine refuses the input it took: %v", at, err)
				case c != call.Choice:
					breach("I1", "%sthe choice taken again is %v %s (%v)", at, c.Kind, c.ID, c.Reason)
				}
			}
			for _, d := range reasonCode(run.Step{Input: call.Input, Choice: call.Choice}) {
				breach("I2", "%s%s", at, d)
			}
		}
		for _, d := range authorOrder(s) {
			breach("I3", "%s", d)
		}
		if d := solvedAgain(s); d != "" {
			breach("I4", "%s", d)
		}
	}
	return out
}

// Same gives the breach of I1 between two runs on the same seeds: every choice, every piece of
// evidence and the model's state at the end are the same bits, or the first that differs is named.
func Same(a, b *run.Record) []Breach {
	breach := func(step int, format string, args ...any) []Breach {
		return []Breach{{Invariant: "I1", Unit: a.Unit, Step: step, Detail: fmt.Sprintf(format, args...)}}
	}
	if a.Unit != b.Unit || len(a.Steps) != len(b.Steps) {
		return breach(-1, "%s has %d steps, then %s %d", a.Unit, len(a.Steps), b.Unit, len(b.Steps))
	}
	for k := range a.Steps {
		sa, sb := a.Steps[k], b.Steps[k]
		if sa.Choice != sb.Choice || sa.Resumed != sb.Resumed {
			return breach(k, "chose %v (resumed %v), then %v (%v)", sa.Choice, sa.Resumed, sb.Choice, sb.Resumed)
		}
		if !slices.EqualFunc(sa.Further, sb.Further, func(x, y run.Call) bool { return x.Choice == y.Choice }) {
			return breach(k, "went further through the lesson %d times, then %d, or by other choices", len(sa.Further), len(sb.Further))
		}
		if len(sa.Evidence) != len(sb.Evidence) {
			return breach(k, "%d pieces of evidence, then %d", len(sa.Evidence), len(sb.Evidence))
		}
		for i := range sa.Evidence {
			ea, eb := encode(sa.Evidence[i]), encode(sb.Evidence[i])
			if !bytes.Equal(ea, eb) {
				return breach(k, "the evidence of attempt %d differs", i)
			}
		}
	}
	sa, errA := learnermodel.EncodeState(a.State)
	sb, errB := learnermodel.EncodeState(b.State)
	if errA != nil || errB != nil || !bytes.Equal(sa, sb) {
		return breach(-1, "the model's states at the end differ (%v, %v)", errA, errB)
	}
	return nil
}

// encode gives the bytes of a piece of evidence; nil for an attempt that gave none, and the error
// as bytes for one the codec refuses, so that two such differ only if their errors do.
func encode(e *learnermodel.Evidence) []byte {
	if e == nil {
		return nil
	}
	b, err := learnermodel.EncodeEvidence(*e)
	if err != nil {
		return []byte("error: " + err.Error())
	}
	return b
}

// view is what the input says of one skill, by its full ID.
type view struct {
	skill   learnermodel.Skill
	read    learnermodel.SkillReadout
	hasRead bool
	issue   engine.SkillIssue
	module  learnermodel.ModuleReadout
}

func viewOf(in engine.Input, id string) (view, bool) {
	var v view
	found := false
	for _, m := range in.Modules {
		if i := slices.IndexFunc(m.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == id }); i >= 0 {
			v.skill, found = m.Graph.Skills[i], true
			for _, mr := range in.Readout.Modules {
				if mr.Module == m.Graph.Module {
					v.module = mr
					if j := slices.IndexFunc(mr.Skills, func(s learnermodel.SkillReadout) bool { return s.Skill == id }); j >= 0 {
						v.read, v.hasRead = mr.Skills[j], true
					}
				}
			}
		}
	}
	if i := slices.IndexFunc(in.Summary.Skills, func(s engine.SkillIssue) bool { return s.Skill == id }); i >= 0 {
		v.issue = in.Summary.Skills[i]
	}
	return v, found
}

// failed is a failure on a skill by the facts of a result, as the engine and the model read it: of
// the items answered, fewer than half right at the first verdict, or the solution shown with none
// right (architecture 3.6 engine; R06 learner-model).
func failed(r engine.SkillResult) bool {
	return 2*r.FirstRight < r.Answered || r.SolutionShown && r.FirstRight == 0
}

// reasonCode gives the breaches of I2 on a step: a choice of a task or a lesson has one code of the
// closed list ReasonCode, a choice of nothing none; and the code holds, on one of the choice's
// skills, the facts its meaning names (the text of the contracts, section 15; docs/architecture.md,
// section 3.5; the amendment of 2026-11-08 to R04, R17). Of each meaning the bench checks what it
// can read without the engine's own rule: what a skill was given, what the model reads of it and
// of its module, the last attempt and the runs.
func reasonCode(s run.Step) []string {
	c, in := s.Choice, s.Input
	if c.Kind == engine.ChoiceNothing {
		if c.Reason != engine.ReasonUnspecified || c.ID != "" {
			return []string{"a choice of nothing with a code or an ID"}
		}
		return nil
	}
	if c.Kind != engine.ChoiceTask && c.Kind != engine.ChoiceLesson {
		return []string{"a choice of no kind the contracts name"}
	}
	if c.Reason < engine.ReasonConsolidation || c.Reason > engine.ReasonRotation {
		return []string{fmt.Sprintf("the code %d is not of the closed list", c.Reason)}
	}
	ids := s.Skills()
	if len(ids) == 0 {
		return []string{"the choice is on no skill of the input"}
	}
	// The codes of branches 3–9 of the engine are given on skills with a verdict, which the model
	// reads; the author's order on the others (architecture 3.5 engine).
	holds := func(meaning func(view) bool) bool {
		for _, id := range ids {
			if v, ok := viewOf(in, id); ok && v.skill.Assessability == learnermodel.AssessabilityVerdict && meaning(v) {
				return true
			}
		}
		return false
	}
	if (c.Reason == engine.ReasonConsolidation || c.Reason == engine.ReasonProbe) && c.Kind != engine.ChoiceTask {
		return []string{fmt.Sprintf("the code %d is given only to a task", c.Reason)}
	}
	var meaning string
	ok := false
	switch c.Reason {
	case engine.ReasonConsolidation:
		meaning = "a failure on the skill in the last attempt, not left, past the trial series of its module"
		ok = in.HasLast && !in.Last.Left && holds(func(v view) bool {
			if !v.module.TrialOver || v.module.LastInTrial {
				return false
			}
			return slices.ContainsFunc(in.Last.Results, func(r engine.SkillResult) bool { return r.Skill == v.skill.ID && failed(r) })
		})
	case engine.ReasonNewSkill:
		meaning = "a skill not given before"
		ok = holds(func(v view) bool { return !v.issue.Issued })
	case engine.ReasonQueue:
		meaning = "a skill given before and not mastered"
		ok = holds(func(v view) bool { return v.issue.Issued && !(v.hasRead && v.read.Mastered) })
	case engine.ReasonRepetition:
		meaning = "a mastered skill whose review date has come"
		ok = holds(func(v view) bool {
			return v.hasRead && v.read.Mastered && (v.read.Review.Already || v.read.Review.At <= in.Time)
		})
	case engine.ReasonProbe:
		meaning = "a skill given before, past the trial series of its module"
		ok = holds(func(v view) bool { return v.issue.Issued && v.module.TrialOver })
	case engine.ReasonSkillChosen:
		meaning = "the skill chosen on the map"
		ok = in.MapSkill != "" && slices.Contains(ids, in.MapSkill)
	case engine.ReasonAuthorOrder:
		meaning = "a skill with a rubric or without assessment"
		ok = slices.ContainsFunc(ids, func(id string) bool {
			v, ok := viewOf(in, id)
			return ok && v.skill.Assessability != learnermodel.AssessabilityVerdict
		})
	case engine.ReasonOpenRun:
		meaning = "the lesson of a run that is open"
		ok = c.Kind == engine.ChoiceLesson && slices.ContainsFunc(in.Runs, func(r engine.Run) bool { return r.Lesson == c.ID && r.End == engine.RunOpen })
	case engine.ReasonRotation:
		meaning = "a skill given before"
		ok = holds(func(v view) bool { return v.issue.Issued })
	}
	if !ok {
		return []string{fmt.Sprintf("the code %d means %s, and none of %s is", c.Reason, meaning, strings.Join(ids, ", "))}
	}
	return nil
}

// authorOrder gives the breach of I3 on a step: a choice by the author's order is on a skill whose
// prerequisites with a rubric or without assessment are done — a lesson that teaches them passed, a
// task on them closed by "Done" —, or cannot be done: no lesson suited to the learner teaches them
// and no task of the pool is on them (concepts 4 spec; product 10 spec). The prerequisites are
// followed through skills with a verdict too, and the choice's own skills — a lesson that teaches a
// skill with its prerequisite — wait for nothing, as the engine reads them (lessons.go engine).
// Which comes next by the graph, without the learner model, is the engine's rule.
func authorOrder(s run.Step) []string {
	c, in := s.Choice, s.Input
	if c.Reason != engine.ReasonAuthorOrder {
		return nil
	}
	own := s.Skills()
	primaries := func(t learnermodel.Task) []string {
		out := slices.Clone(t.Primary)
		for _, it := range t.Items {
			out = append(out, it.Primary...)
		}
		return out
	}
	done := func(id string) bool {
		if slices.ContainsFunc(in.Runs, func(r engine.Run) bool { return r.End == engine.RunPassed && slices.Contains(r.Skills, id) }) {
			return true
		}
		for _, d := range in.RubricDone {
			for _, m := range in.Modules {
				for _, cand := range m.Pool {
					if cand.Task.ID == d.Task && slices.Contains(primaries(cand.Task), id) {
						return true
					}
				}
			}
		}
		return false
	}
	doable := func(id string) bool {
		for _, m := range in.Modules {
			if slices.ContainsFunc(m.Lessons, func(l engine.Lesson) bool { return slices.Contains(l.Skills, id) }) ||
				slices.ContainsFunc(m.Pool, func(cand engine.Candidate) bool { return slices.Contains(primaries(cand.Task), id) }) {
				return true
			}
		}
		return false
	}
	var out []string
	for _, id := range own {
		v, ok := viewOf(in, id)
		if !ok || v.skill.Assessability == learnermodel.AssessabilityVerdict {
			continue
		}
		seen := map[string]bool{}
		next := slices.Clone(v.skill.Prerequisites)
		for len(next) > 0 {
			p := next[0]
			next = next[1:]
			if seen[p] || slices.Contains(own, p) {
				continue
			}
			seen[p] = true
			pv, ok := viewOf(in, p)
			if !ok {
				continue
			}
			if pv.skill.Assessability != learnermodel.AssessabilityVerdict && !done(p) && doable(p) {
				out = append(out, fmt.Sprintf("%s comes by the author's order before %s, which it rests on, is done", id, p))
			}
			next = append(next, pv.skill.Prerequisites...)
		}
	}
	return out
}

// solvedAgain gives the breach of I4 on a step: the choice is the task of the last closed attempt,
// not a template, and its solution went to the device then (R82 spec; R107 spec): "at once" is
// "not as the next task".
func solvedAgain(s run.Step) string {
	c, in := s.Choice, s.Input
	if c.Kind != engine.ChoiceTask || !in.HasLast || in.Last.InLesson || c.ID != in.Last.Task {
		return ""
	}
	for _, m := range in.Modules {
		for _, cand := range m.Pool {
			if cand.Task.ID == c.ID && cand.Task.Template {
				return ""
			}
		}
	}
	if slices.ContainsFunc(in.Summary.Tasks, func(t engine.TaskSolution) bool { return t.Task == c.ID }) {
		return "the task whose solution went to the device in the last attempt comes as the next task"
	}
	return ""
}
