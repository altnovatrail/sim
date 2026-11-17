package measure

import (
	"math"
	"strconv"

	"github.com/altnovatrail/sim/internal/run"
)

// The numbers of M11, fairness by the grade (docs/architecture.md, section 3.5; R55): the run's
// length in answers, which a learner who never comes to the reference's estimate counts as; the
// answers by which the estimate is the same; the gap under which, and the answers for which in a row,
// it holds, as M3 settles.
const (
	FairLength  = 200
	FairBy      = 50
	FairGap     = settleError
	FairAnswers = settleAnswers
)

// Trace is what M11 reads of one learner's run, written into its line in the cells of fairness by
// the grade alone: the model's module level after each answer, and the attempts with a verdict at
// each band once FairBy answers are behind.
type Trace struct {
	// Levels are θ of the model after each of the first FairLength answers — the first answers of
	// items with a verdict —, in their order: each answer of an attempt has the level after it.
	Levels []float64 `json:"levels"`
	// Truth is θ* of the truth's module level after each of the same answers, as run.Obs gives it,
	// for M13 (R73 in docs/decisions.md). A file written before it has none.
	Truth []float64 `json:"truth,omitempty"`
	// Bands count the attempts with a verdict whose first answer is the 51st or later, to the end of
	// the run, those of lessons and consolidations among them, by the place of their task's band:
	// Bands[p-1] for the place p.
	Bands []int `json:"bands"`
	// Guard is how often the guard of learner-model lowered the module level on the whole run, as
	// run.Record counts it; read by M12 on the reference's cell alone (R72 in docs/decisions.md). A file
	// of a model with no guard has none.
	Guard int `json:"guard,omitempty"`
	// Corrected is the first answers before the adult corrected the module level, as run.Record
	// gives it, 0 when no adult did (C-).
	Corrected int `json:"corrected,omitempty"`
}

// Traced gives the trace of M11 of one learner's run.
func Traced(rec *run.Record) *Trace {
	obs, _ := attempts(rec)
	t := &Trace{Levels: []float64{}, Bands: []int{}, Guard: rec.Guard, Corrected: rec.Corrected}
	for _, o := range obs {
		for k := o.Answers; k < o.Answers+len(o.Items) && k < FairLength; k++ {
			t.Levels, t.Truth = append(t.Levels, o.Level), append(t.Truth, o.TrueLevel)
		}
		if o.Answers >= FairBy && o.Band > 0 {
			for len(t.Bands) < int(o.Band) {
				t.Bands = append(t.Bands, 0)
			}
			t.Bands[o.Band-1]++
		}
	}
	return t
}

// Fair gives the sums of M11 of the learner n of a grade against the learner n of the reference
// grade, a pair of one truth and one unit: the answers k before the first of FairAnswers in a row at
// which the two levels lie within FairGap — none before it gives 0, as M3 counts, and a learner for
// whom it never holds counts as FairLength —; whether k is at most FairBy; and the attempts of each
// at each band after FairBy answers, every band either has under both sides. So a cell of pairs
// reads as a cell of learners, its resamples of n those of both sides at once.
func Fair(grade, reference *Trace) Sums {
	s := Sums{"M11.pairs": 1}
	k := settled(grade.Levels, reference.Levels)
	s.add("M11.answers.sum", float64(k))
	if k <= FairBy {
		s.add("M11.same", 1)
	}
	for p := range max(len(grade.Bands), len(reference.Bands)) {
		for _, side := range []struct {
			name  string
			bands []int
		}{{"grade", grade.Bands}, {"reference", reference.Bands}} {
			c := 0.0
			if p < len(side.bands) {
				c = float64(side.bands[p])
			}
			s.add("M11.band."+side.name+"."+strconv.Itoa(p+1), c)
			s.add("M11.band."+side.name, c)
		}
	}
	return s
}

// settled gives the answers k before the first of FairAnswers in a row at which two series of
// levels lie within FairGap, as M3 counts catching up: 0 when it holds from the first answer, and
// FairLength when it never holds.
func settled(a, b []float64) int {
	streak, from := 0, 0
	for i := range min(len(a), len(b)) {
		if math.Abs(a[i]-b[i]) >= FairGap {
			streak = 0
			continue
		}
		if streak == 0 {
			from = i
		}
		if streak++; streak >= FairAnswers {
			return from
		}
	}
	return FairLength
}

// band is M11.band: the largest absolute difference, over the bands, between the grade's share of
// the attempts after FairBy answers at a band and the reference's, read when both sides have any.
func band(s Sums) (float64, bool) {
	g, r := s["M11.band.grade"], s["M11.band.reference"]
	if g <= 0 || r <= 0 {
		return 0, false
	}
	worst := 0.0
	for p := 1; ; p++ {
		place := strconv.Itoa(p)
		x, ok := s["M11.band.grade."+place]
		if !ok {
			return worst, true
		}
		worst = max(worst, math.Abs(float64(x/g)-float64(s["M11.band.reference."+place]/r)))
	}
}

// FairNumbers are the numbers of M11, read off a cell of pairs of a grade and the reference
// (docs/architecture.md, section 3.5; R55).
func FairNumbers() []Number {
	return []Number{
		{"M11.answers", ratio("M11.answers.sum", "M11.pairs")},
		{"M11.same", ratio("M11.same", "M11.pairs")},
		{"M11.band", band},
	}
}

// Guarded gives the sums of M12 of one learner, read on the reference's cell, whose truth is the
// band of its grade (R72 in docs/decisions.md): whether the guard of learner-model lowered its
// module level at least once on the run, a false firing there by construction.
func Guarded(t *Trace) Sums {
	s := Sums{"M12.learners": 1}
	if t.Guard > 0 {
		s.add("M12.false", 1)
	}
	return s
}

// GuardNumbers are the numbers of M12, read off a cell of learners by their sums of Guarded: the
// share of them whose level the guard lowered (R72).
func GuardNumbers() []Number {
	return []Number{{"M12.false", ratio("M12.false", "M12.learners")}}
}

// Caught gives the sums of M13 of one learner, against its own truth (R73 in docs/decisions.md): the
// answers k before the first of FairAnswers in a row at which the model's module level lies within
// FairGap of the truth's, counted as Fair counts them — FairLength when it never holds —, and whether
// k is at most FairBy.
func Caught(t *Trace) Sums {
	s := Sums{"M13.learners": 1}
	k := settled(t.Levels, t.Truth)
	s.add("M13.answers.sum", float64(k))
	if k <= FairBy {
		s.add("M13.same", 1)
	}
	return s
}

// TruthNumbers are the numbers of M13, read off a cell of learners by their sums of Caught: the mean
// answers until the estimate holds at the truth, and the share for whom it does by FairBy (R73).
func TruthNumbers() []Number {
	return []Number{
		{"M13.answers", ratio("M13.answers.sum", "M13.learners")},
		{"M13.same", ratio("M13.same", "M13.learners")},
	}
}
