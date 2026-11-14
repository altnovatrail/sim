package measure

import (
	"slices"

	"github.com/altnovatrail/engine"
	"github.com/altnovatrail/sim/internal/invariant"
	"github.com/altnovatrail/sim/internal/run"
)

// Note gives the sums of one learner's run for the note of consolidation in a lesson, E6, which
// decides nothing: no check, no path and no weight, only the note of a judgement beside its
// criterion (docs/architecture.md, sections 3.5 and 3.6; R49 in docs/decisions.md). It reads the
// lessons whose skills all have a verdict, the ones the run calls the engine in (R48):
//   - the consolidations in a lesson — a task right before the lesson resumed — and the passes;
//   - the items with a verdict of the lesson's task links and of the consolidations, by P* when
//     given, in the corridor 0.70–0.85, below 0.50 and above 0.95, as E1 reads them;
//   - P* when given of the items of the next task link of the same pass after a link that failed
//     on a skill of the lesson, as the engine reads a failure (invariant.Failed).
func Note(rec *run.Record) Sums {
	s := Sums{}
	var skills []string // the skills of the lesson in pass, while all of them have a verdict
	failed := false     // the last task link of the pass failed on a skill of its lesson
	for k, st := range rec.Steps {
		switch {
		case st.Choice.Kind == engine.ChoiceLesson:
			if !st.Resumed {
				skills, failed = verdictLesson(st.Input, st.Choice.ID), false
			}
			if skills == nil {
				continue
			}
			for i, o := range st.Obs {
				if o != nil { // a task with a rubric has no items to read, and is a link all the same
					if failed {
						for _, it := range o.Items {
							s.add("E6.next.sum", it.PGiven)
							s.add("E6.next.items", 1)
						}
					}
					shares(s, "E6.lesson", o.Items)
				}
				failed = i < len(st.Results) && slices.ContainsFunc(st.Results[i], func(r engine.SkillResult) bool {
					return slices.Contains(skills, r.Skill) && invariant.Failed(r)
				})
			}
			if st.Pass != nil {
				s.add("E6.passes", 1)
				skills = nil
			}
		case skills != nil && st.Choice.Kind == engine.ChoiceTask && k+1 < len(rec.Steps) && rec.Steps[k+1].Resumed:
			s.add("E6.consolidations", 1)
			for _, o := range st.Obs {
				if o != nil {
					shares(s, "E6.consolidation", o.Items)
				}
			}
		}
	}
	return s
}

// verdictLesson gives the skills of a lesson of an input when it has skills and every one has a
// verdict, as the run calls the engine within it (run.Verdicts; R48); nil otherwise.
func verdictLesson(in engine.Input, id string) []string {
	for _, m := range in.Modules {
		i := slices.IndexFunc(m.Lessons, func(l engine.Lesson) bool { return l.ID == id })
		if i < 0 {
			continue
		}
		if skills := m.Lessons[i].Skills; run.Verdicts(m.Graph, skills) {
			return skills
		}
		return nil
	}
	return nil
}

// NoteNumbers are the numbers of the note of consolidation in a lesson a cell is read for: the
// consolidations a pass, the shares of the items of the lessons' links and of the consolidations,
// and the mean P* of the next link after a failure. The criteria hold none of them (R49).
func NoteNumbers() []Number {
	out := []Number{{"E6.consolidations", ratio("E6.consolidations", "E6.passes")}}
	for _, of := range []string{"E6.lesson", "E6.consolidation"} {
		for _, share := range []string{"corridor", "hard", "easy"} {
			out = append(out, Number{of + "." + share, ratio(of+"."+share, of+".items")})
		}
	}
	return append(out, Number{"E6.next", ratio("E6.next.sum", "E6.next.items")})
}
