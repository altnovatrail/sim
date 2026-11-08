package measure

import (
	"maps"
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/run"
)

// The answer counts M1 reads the error after, the pauses M2 cuts by, in seconds, and the thresholds
// of M3, M4, M5 and M10 (docs/architecture.md, section 3.5).
var (
	errorAfter = []int{5, 10, 20, 50, 100, 200}
	pauses     = []struct {
		name     string
		from, to int64
	}{{"pause0", 0, day}, {"pause1", day, 7 * day}, {"pause7", 7 * day, 30 * day}, {"pause30", 30 * day, 1 << 62}}
)

const (
	day           = 86400
	lagFrom       = 101 // M3 reads the lag from this answer on
	settleError   = 0.5 // M3: the error of the skills under which an answer is settled
	settleAnswers = 10  // M3: in a row
	breakPause    = 14 * day
	breakWindow   = 7 * day
	pagesEvidence = 2    // M4: a mastery declared with at most this much evidence is "of pages"
	catchError    = 0.05 // M5: |P − P*| of the reference task under which forgetting is caught up
	catchAnswers  = 5    // M5: in a row
	startStreak   = 15   // M7: the first answers the longest wrong streak is read in
	startHard     = 10   // M7: the first items whose share with P* below 0.5 is read
	freshPause    = day  // M10: the answers after pauses of at least a day
)

// attempts gives the observed attempts of a record in order, with their results.
func attempts(rec *run.Record) ([]*run.Obs, [][]bool) {
	var obs []*run.Obs
	var rights [][]bool
	for _, s := range rec.Steps {
		for i, o := range s.Obs {
			if o == nil {
				continue
			}
			var r []bool
			for _, res := range s.Attempts[i].Items {
				r = append(r, res.FirstRight)
			}
			obs, rights = append(obs, o), append(rights, r)
		}
	}
	return obs, rights
}

// errors gives the errors θ + δ − θ* − δ* of the skills of a snapshot that have evidence.
func errors(skills []run.SkillObs) []float64 {
	var out []float64
	for _, s := range skills {
		if s.Evidence > 0 {
			out = append(out, s.Estimate-s.Truth)
		}
	}
	return out
}

func rms(xs []float64) (float64, bool) {
	if len(xs) == 0 {
		return 0, false
	}
	sq := 0.0
	for _, x := range xs {
		sq += float64(x * x)
	}
	return float64(sq / float64(len(xs))), true // the mean square; its root is taken at the cell
}

// Model gives the sums of one learner's run for the measures of the model, M1–M10
// (docs/architecture.md, section 3.5).
func Model(rec *run.Record) Sums {
	obs, rights := attempts(rec)
	return model(rec, obs, rights, false)
}

// Oracle gives the sums of the ceiling of the model on one learner's run, the oracle: its estimate
// is the truth θ* + δ*, its P is P*, it declares "mastered" by true mastery, and it holds an anchor
// when the truth does (docs/architecture.md, section 3.6). It reads the same answers as the model
// of the run, so that the two compare by learner (R24 in docs/decisions.md).
func Oracle(rec *run.Record) Sums {
	obs, rights := attempts(rec)
	out := make([]*run.Obs, len(obs))
	for k, o := range obs {
		c := *o
		c.Items = slices.Clone(o.Items)
		for i := range c.Items {
			c.Items[i].PModel, c.Items[i].PFresh = c.Items[i].PTrue, c.Items[i].PTrue
		}
		c.Before, c.After = oracleSkills(o.Before), oracleSkills(o.After)
		c.Level = o.TrueLevel
		out[k] = &c
	}
	return model(rec, out, rights, true)
}

// oracleSkills gives the oracle's snapshot of skills: the truth in place of the model's estimate,
// mastery and P of the reference task.
func oracleSkills(skills []run.SkillObs) []run.SkillObs {
	out := slices.Clone(skills)
	for i := range out {
		out[i].Estimate, out[i].Mastered, out[i].PRef = out[i].Truth, out[i].TrulyMastered, out[i].PTrueRef
	}
	return out
}

// model gives the sums of the measures of the model on the observations of a run, the model's or
// the oracle's.
func model(rec *run.Record, obs []*run.Obs, rights [][]bool, oracle bool) Sums {
	s := Sums{"learners": 1}
	m1(s, obs)
	m2(s, obs, rights)
	m3(s, rec, obs)
	m4(s, obs)
	m5(s, obs)
	m6(s, obs)
	m7(s, obs, rights)
	m8(s, obs)
	m9(s, rec, oracle)
	m10(s, obs, rights)
	return s
}

// m1: the error after 5, 10, 20, 50, 100 and 200 answers — the root mean square and the mean of θ +
// δ − θ* − δ* over the skills with evidence, read once each count is first reached.
func m1(s Sums, obs []*run.Obs) {
	next := 0
	for _, o := range obs {
		after := o.Answers + len(o.Items)
		for next < len(errorAfter) && after >= errorAfter[next] {
			k := itoa(errorAfter[next])
			for _, e := range errors(o.After) {
				s.add("M1."+k+".sq", float64(e*e))
				s.add("M1."+k+".sum", e)
				s.add("M1."+k+".n", 1)
			}
			next++
		}
	}
}

// m2: the calibration of P with retention at the first verdicts — Brier, the bias, and the error of
// calibration over ten bins —, overall and cut by the pause since the skills' last practice, by
// items of attempts of many items, and by answers with help.
func m2(s Sums, obs []*run.Obs, rights [][]bool) {
	for k, o := range obs {
		for i, it := range o.Items {
			y := 0.0
			if rights[k][i] {
				y = 1
			}
			cuts := []string{"all"}
			for _, p := range pauses {
				if it.Pause >= p.from && it.Pause < p.to {
					cuts = append(cuts, p.name)
				}
			}
			if it.Pause >= 7*day {
				cuts = append(cuts, "after7") // after pauses of at least 7 days, the cut of "forgets"
			}
			if it.Many {
				cuts = append(cuts, "pages")
			}
			if it.Help {
				cuts = append(cuts, "help")
			}
			bin := itoa(min(9, int(float64(10*it.PModel))))
			for _, c := range cuts {
				d := it.PModel - y
				s.add("M2."+c+".brier", float64(d*d))
				s.add("M2."+c+".p", it.PModel)
				s.add("M2."+c+".y", y)
				s.add("M2."+c+".n", 1)
				s.add("M2."+c+".bin"+bin+".p", it.PModel)
				s.add("M2."+c+".bin"+bin+".y", y)
				s.add("M2."+c+".bin"+bin+".n", 1)
			}
		}
	}
}

// m3: the lag — θ + δ − θ* − δ* of the skills with evidence from the 101st answer on, on average —;
// and after a jump, a drop or a break, the answers until the error of the skills falls under 0.5 to
// stay there for ten answers in a row — none when it is there from the first attempt after —, and
// whether it never does.
func m3(s Sums, rec *run.Record, obs []*run.Obs) {
	for _, o := range obs {
		if o.Answers+len(o.Items) < lagFrom {
			continue
		}
		for _, e := range errors(o.After) {
			s.add("M3.lag.sum", e)
			s.add("M3.lag.n", 1)
		}
	}
	jump, ok := jumpAnswer(rec, obs)
	if !ok {
		return
	}
	s.add("M3.jumped", 1)
	streak, from := 0, 0
	for _, o := range obs {
		after := o.Answers + len(o.Items)
		if after <= jump {
			continue
		}
		ms, ok := rms(errors(o.After))
		if !ok || ms >= float64(settleError*settleError) {
			streak = 0
			continue
		}
		if streak == 0 {
			from = max(o.Answers, jump)
		}
		streak += len(o.Items)
		if streak >= settleAnswers {
			s.add("M3.settle.sum", float64(from-jump))
			s.add("M3.settle.n", 1)
			return
		}
	}
	s.add("M3.unsettled", 1)
}

// jumpAnswer gives the answer after which the truth jumped: the jump of G3 or G3-drop, or the end of
// the break of F-break, counted in first answers of items with a verdict; none when the run did not
// get past it.
func jumpAnswer(rec *run.Record, obs []*run.Obs) (int, bool) {
	l := rec.Learner
	if l == nil || len(obs) == 0 {
		return 0, false
	}
	if end, ok := l.DropAt(); ok {
		for _, o := range obs {
			if o.At >= end {
				return o.Answers, true
			}
		}
		return 0, false
	}
	for _, o := range obs {
		if len(o.After) == 0 {
			continue
		}
		last := obs[len(obs)-1]
		if m, err := l.Module(moduleOf(o.After[0].Skill)); err == nil && m.JumpAt > 0 && last.Answers+len(last.Items) > m.JumpAt {
			return m.JumpAt, true
		}
	}
	return 0, false
}

// moduleOf gives the module of a skill by its full ID.
func moduleOf(skill string) string {
	module, _, _ := strings.Cut(skill, "/skill/")
	return module
}

// before gives a skill's snapshot before an attempt, and whether the module had it.
func before(o *run.Obs, skill string) (run.SkillObs, bool) {
	i := slices.IndexFunc(o.Before, func(b run.SkillObs) bool { return b.Skill == skill })
	if i < 0 {
		return run.SkillObs{}, false
	}
	return o.Before[i], true
}

// practised tells whether an attempt answered a skill: the model took evidence of it. The answers
// of a skill are the items of the attempts that answered it.
func practised(o *run.Obs, a run.SkillObs) bool {
	b, _ := before(o, a.Skill)
	return a.Evidence > b.Evidence
}

// paused tells whether an attempt answered a skill after a pause of it of at least 14 days.
func paused(o *run.Obs, a run.SkillObs) bool {
	b, ok := before(o, a.Skill)
	return ok && practised(o, a) && b.Evidence > 0 && o.At-b.Last >= breakPause
}

// declared gives the skills an attempt's model declared mastered: not before it, mastered after.
func declared(o *run.Obs) []run.SkillObs {
	var out []run.SkillObs
	for _, a := range o.After {
		if b, _ := before(o, a.Skill); a.Mastered && !b.Mastered {
			out = append(out, a)
		}
	}
	return out
}

// m4: the share of masteries declared whose skill the truth has not mastered on the band of mastery
// at the moment; and of them, those declared within 7 days after a pause of the skill of at least
// 14 days, and those declared on at most two pieces of evidence.
func m4(s Sums, obs []*run.Obs) {
	window := map[string]int64{}
	for _, o := range obs {
		for _, a := range o.After {
			if paused(o, a) {
				window[a.Skill] = o.At + breakWindow
			}
		}
		for _, d := range declared(o) {
			bad := 0.0
			if !d.TrulyMastered {
				bad = 1
			}
			cuts := []string{"all"}
			if end, ok := window[d.Skill]; ok && o.At <= end {
				cuts = append(cuts, "break")
			}
			if d.Evidence <= pagesEvidence {
				cuts = append(cuts, "pages")
			}
			for _, c := range cuts {
				s.add("M4."+c+".declared", 1)
				s.add("M4."+c+".false", bad)
			}
		}
	}
}

// m5: after a pause of a skill of at least 14 days, the answers of the skill until |P − P*| of the
// reference task of its band falls under 0.05 to stay there for five answers in a row — none when it
// is there from the first answer —, as M3 counts its answers, and whether it never does.
func m5(s Sums, obs []*run.Obs) {
	type episode struct{ answers, streak, from int }
	open := map[string]*episode{}
	for _, o := range obs {
		for _, a := range o.After {
			if paused(o, a) && open[a.Skill] == nil {
				open[a.Skill] = &episode{}
				s.add("M5.episodes", 1)
			}
			ep := open[a.Skill]
			if ep == nil || !practised(o, a) {
				continue
			}
			n := len(o.Items)
			if d := a.PRef - a.PTrueRef; d > -catchError && d < catchError {
				if ep.streak == 0 {
					ep.from = ep.answers
				}
				ep.streak += n
			} else {
				ep.streak = 0
			}
			ep.answers += n
			if ep.streak >= catchAnswers {
				s.add("M5.caught.sum", float64(ep.from))
				s.add("M5.caught.n", 1)
				delete(open, a.Skill)
			}
		}
	}
	s.add("M5.never", float64(len(open)))
}

// m6: for each skill the truth has mastered on its band, the answers of the skill from that moment
// to the model's mastery — none when both come with one attempt —, and whether it never came; a
// mastery declared before the true one is false, and M4's.
func m6(s Sums, obs []*run.Obs) {
	since := map[string]int{}
	done := map[string]bool{}
	for _, o := range obs {
		for _, a := range o.After {
			if done[a.Skill] {
				continue
			}
			n := 0
			if practised(o, a) {
				n = len(o.Items)
			}
			k, waiting := since[a.Skill]
			switch {
			case a.Mastered && !waiting && !a.TrulyMastered:
				done[a.Skill] = true // declared before the truth mastered it
			case a.Mastered:
				if !waiting {
					since[a.Skill], n = 0, 0 // declared with the attempt the truth mastered it
				}
				s.add("M6.late.sum", float64(k+n))
				s.add("M6.late.n", 1)
				done[a.Skill] = true
			case waiting:
				since[a.Skill] = k + n
			case a.TrulyMastered:
				since[a.Skill] = 0
			}
		}
	}
	for _, skill := range slices.Sorted(maps.Keys(since)) {
		if !done[skill] {
			s.add("M6.never", 1)
		}
	}
	s.add("M6.cases", float64(len(since)))
}

// m7: the start — θ − θ* of the module after 5 and 10 answers, the longest streak of wrong answers
// among the first 15, and the share of the first 10 items with P* below 0.5.
func m7(s Sums, obs []*run.Obs, rights [][]bool) {
	for _, at := range []int{5, 10} {
		for _, o := range obs {
			if o.Answers+len(o.Items) >= at {
				d := o.Level - o.TrueLevel
				s.add("M7."+itoa(at)+".sq", float64(d*d))
				s.add("M7."+itoa(at)+".n", 1)
				break
			}
		}
	}
	longest, streak, seen, hard := 0, 0, 0, 0
	for k, o := range obs {
		for i, it := range o.Items {
			if seen < startStreak {
				if rights[k][i] {
					streak = 0
				} else {
					streak++
					longest = max(longest, streak)
				}
			}
			if seen < startHard && it.PTrue < 0.5 {
				hard++
			}
			seen++
		}
	}
	if seen > 0 {
		s.add("M7.streak.sum", float64(longest))
		s.add("M7.streak.n", 1)
		s.add("M7.hard", float64(hard))
		s.add("M7.items", float64(min(seen, startHard)))
	}
}

// m8: the mean shift of θ + δ of the primary skills of a task, not a template, whose solution went
// to the device before, by the evidence of that repeat. The truth does not move on a repeat of
// W-repeat, so every shift is an error.
func m8(s Sums, obs []*run.Obs) {
	for _, o := range obs {
		if !o.Repeat {
			continue
		}
		for _, a := range o.After {
			if b, ok := before(o, a.Skill); ok && practised(o, a) {
				s.add("M8.sum", a.Estimate-b.Estimate)
				s.add("M8.n", 1)
			}
		}
	}
}

// m9: at the start of each session, the share of readings "holds" an anchor skill when the truth
// does not hold it, and of "does not hold" when it does. The truth holds an anchor when a skill that
// teaches it is truly mastered with its retention at the time (docs/architecture.md, section 3.5).
func m9(s Sums, rec *run.Record, oracle bool) {
	l := rec.Learner
	if l == nil {
		return
	}
	lastDay := int64(-1)
	for _, st := range rec.Steps {
		d := st.Input.Time / day
		if d == lastDay {
			continue
		}
		lastDay = d
		for _, a := range st.Input.Readout.Anchors {
			holds := false
			for _, m := range st.Input.Modules {
				for _, sk := range m.Graph.Skills {
					if slices.Contains(sk.Teaches, a.Anchor) && l.ReferenceAt(sk.ID, l.TaughtOn(sk.ID), st.Input.Time) >= learner.Mastery {
						holds = true
					}
				}
			}
			status := a.Status
			if oracle {
				status = learnermodel.AnchorHasNot
				if holds {
					status = learnermodel.AnchorHas
				}
			}
			switch status {
			case learnermodel.AnchorHas:
				s.add("M9.has", 1)
				if !holds {
					s.add("M9.hasFalse", 1)
				}
			case learnermodel.AnchorHasNot:
				s.add("M9.not", 1)
				if holds {
					s.add("M9.notFalse", 1)
				}
			}
		}
	}
}

// m10: on the answers after pauses of at least a day, Brier of P with retention and of P by the same
// estimate without it: the model with forgetting is no worse on learners who do not forget.
func m10(s Sums, obs []*run.Obs, rights [][]bool) {
	for k, o := range obs {
		for i, it := range o.Items {
			if it.Pause < freshPause {
				continue
			}
			y := 0.0
			if rights[k][i] {
				y = 1
			}
			d, f := it.PModel-y, it.PFresh-y
			s.add("M10.retention", float64(d*d))
			s.add("M10.fresh", float64(f*f))
			s.add("M10.n", 1)
		}
	}
}

// ece reads the error of calibration of a cut: over the bins, the share of answers in each times
// |mean P − mean y|.
func ece(cut string) func(Sums) (float64, bool) {
	return func(s Sums) (float64, bool) {
		n := s["M2."+cut+".n"]
		if n <= 0 {
			return 0, false
		}
		e := 0.0
		for b := range 10 {
			k := "M2." + cut + ".bin" + itoa(b)
			if nb := s[k+".n"]; nb > 0 {
				d := float64(s[k+".p"]/nb) - float64(s[k+".y"]/nb)
				if d < 0 {
					d = -d
				}
				e += float64(float64(nb/n) * d)
			}
		}
		return e, true
	}
}

// difference reads a − b of two ratios.
func difference(an, ad, bn, bd string) func(Sums) (float64, bool) {
	return func(s Sums) (float64, bool) {
		x, okx := ratio(an, ad)(s)
		y, oky := ratio(bn, bd)(s)
		return x - y, okx && oky
	}
}

// ModelNumbers are the numbers of the measures of the model a cell is read for.
func ModelNumbers() []Number {
	var out []Number
	for _, k := range errorAfter {
		out = append(out,
			Number{"M1.rms." + itoa(k), root("M1."+itoa(k)+".sq", "M1."+itoa(k)+".n")},
			Number{"M1.mean." + itoa(k), ratio("M1."+itoa(k)+".sum", "M1."+itoa(k)+".n")})
	}
	cuts := []string{"all", "pages", "help", "after7"}
	for _, p := range pauses {
		cuts = append(cuts, p.name)
	}
	for _, c := range cuts {
		out = append(out,
			Number{"M2.brier." + c, ratio("M2."+c+".brier", "M2."+c+".n")},
			Number{"M2.bias." + c, difference("M2."+c+".p", "M2."+c+".n", "M2."+c+".y", "M2."+c+".n")},
			Number{"M2.ece." + c, ece(c)})
	}
	out = append(out,
		Number{"M3.lag", ratio("M3.lag.sum", "M3.lag.n")},
		Number{"M3.settle", ratio("M3.settle.sum", "M3.settle.n")},
		Number{"M3.unsettled", ratio("M3.unsettled", "M3.jumped")},
	)
	for _, c := range []string{"all", "break", "pages"} {
		out = append(out, Number{"M4.false." + c, ratio("M4."+c+".false", "M4."+c+".declared")},
			Number{"M4.declared." + c, ratio("M4."+c+".declared", "learners")})
	}
	out = append(out,
		Number{"M5.caught", ratio("M5.caught.sum", "M5.caught.n")},
		Number{"M5.never", ratio("M5.never", "M5.episodes")},
		Number{"M6.late", ratio("M6.late.sum", "M6.late.n")},
		Number{"M6.never", ratio("M6.never", "M6.cases")},
		Number{"M7.error.5", root("M7.5.sq", "M7.5.n")},
		Number{"M7.error.10", root("M7.10.sq", "M7.10.n")},
		Number{"M7.streak", ratio("M7.streak.sum", "M7.streak.n")},
		Number{"M7.hard", ratio("M7.hard", "M7.items")},
		Number{"M8.shift", ratio("M8.sum", "M8.n")},
		Number{"M9.hasFalse", ratio("M9.hasFalse", "M9.has")},
		Number{"M9.notFalse", ratio("M9.notFalse", "M9.not")},
		Number{"M10.retention", ratio("M10.retention", "M10.n")},
		Number{"M10.fresh", ratio("M10.fresh", "M10.n")},
	)
	return out
}
