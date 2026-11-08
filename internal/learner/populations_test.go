package learner

import (
	"errors"
	"math"
	"testing"

	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/seed"
)

func learnerWith(t *testing.T, population string, n int, a Assumptions, ids ...string) *Learner {
	t.Helper()
	var ms []*pack.Module
	for _, id := range ids {
		ms = append(ms, module(t, id))
	}
	l, err := New(seed.Working, population, n, ms, a)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

// G1: the level of every module moves by the shift, and nothing else of the truth does.
func TestLevelShift(t *testing.T) {
	a := Base()
	a.LevelShift = 2.5
	l, b := learnerWith(t, "X", 1, a, olympiad), learnerWith(t, "X", 1, Base(), olympiad)
	if got := l.modules[olympiad].Level - b.modules[olympiad].Level; math.Abs(got-2.5) > 1e-12 {
		t.Errorf("the shift moved the level by %v", got)
	}
}

// G2 and G2-fading: each first answer moves the level by LearnLevel, the primary skills' offsets by
// LearnSkill more and the others back LearnReturn of the way; fading slows it as F/(F + k).
func TestLearn(t *testing.T) {
	for _, fade := range []float64{0, 50} {
		a := Base()
		a.LearnLevel, a.LearnSkill, a.LearnReturn, a.LearnFade = 0.01, 0.02, 0.01, fade
		l := learnerWith(t, "G2", 3, a, olympiad)
		m := l.modules[olympiad]
		tk := task(t, olympiad, "enum-34-d1-1")
		primary := tk.Candidate.Task.Primary[0]
		other := ""
		for id := range l.skills {
			if id != primary {
				other = id
				break
			}
		}
		l.skills[other].Offset = l.skills[other].Start + 1
		level, own, back := m.Level, l.skills[primary].Offset, l.skills[other].Offset
		if _, err := l.Answer(tk, Start); err != nil {
			t.Fatal(err)
		}
		if fade > 0 {
			level = m.Level
			if _, err := l.Answer(tk, Start+secondsPerDay); err != nil {
				t.Fatal(err)
			}
			if got, want := m.Level-level, float64(0.01*float64(fade/(fade+1))); math.Abs(got-want) > 1e-12 {
				t.Errorf("the second answer grew the level by %v, want %v", got, want)
			}
			continue
		}
		// The first answer, k = 0, is not slowed yet.
		if got := m.Level - level; math.Abs(got-0.01) > 1e-12 {
			t.Errorf("fade %v: the level grew by %v", fade, got)
		}
		if got := l.skills[primary].Offset - own; math.Abs(got-0.02) > 1e-12 {
			t.Errorf("fade %v: the primary skill grew by %v", fade, got)
		}
		if got := back - l.skills[other].Offset; math.Abs(got-0.01) > 1e-12 {
			t.Errorf("fade %v: another skill went back by %v", fade, got)
		}
	}
}

// G3 and G3-drop: the level jumps once, at the answer drawn between 50 and 150.
func TestJump(t *testing.T) {
	a := Base()
	a.JumpBy, a.JumpFrom, a.JumpTo = -1, 50, 150
	l := learnerWith(t, "G3-drop", 2, a, olympiad)
	m := l.modules[olympiad]
	if m.JumpAt < 50 || m.JumpAt > 150 {
		t.Fatalf("the jump at answer %d", m.JumpAt)
	}
	start := m.Level
	tk := task(t, olympiad, "enum-34-d1-1")
	for k := 1; m.Answered < 160; k++ {
		if _, err := l.Answer(tk, Start+int64(k)*secondsPerDay); err != nil {
			t.Fatal(err)
		}
		want := start
		if m.Answered >= m.JumpAt {
			want = start - 1
		}
		if m.Level != want {
			t.Fatalf("after answer %d the level is %v, want %v", m.Answered, m.Level, want)
		}
	}
}

// G4-band: the tasks of the drawn band are truly harder by the bias; the others are not.
func TestBandBias(t *testing.T) {
	a := Base()
	a.BandBias = 0.75
	l := learnerWith(t, "G4-band", 4, a, olympiad)
	b := learnerWith(t, "G4-band", 4, Base(), olympiad)
	hard := l.modules[olympiad].HardBand
	if hard < 1 || hard > 3 {
		t.Fatalf("the hard band is %d", hard)
	}
	for _, tk := range tasks(t, olympiad) {
		diff := l.Difficulty(tk) - b.Difficulty(tk)
		want := 0.0
		if tk.Candidate.Task.Band.Place == hard {
			want = 0.75
		}
		if math.Abs(diff-want) > 1e-12 {
			t.Fatalf("%s in band %d: harder by %v, want %v", tk.Candidate.Task.ID, tk.Candidate.Task.Band.Place, diff, want)
		}
	}
}

// G5: the first half of the skills of the graph is stronger by the halves, the second weaker.
func TestSkillHalves(t *testing.T) {
	a := Base()
	a.SkillHalves, a.SkillSpread = 0.75, 0
	l := learnerWith(t, "G5", 0, a, grade3)
	skills := module(t, grade3).Graph.Skills
	for k, s := range skills {
		want := 0.75
		if 2*k >= len(skills) {
			want = -0.75
		}
		if got := l.skills[s.ID].Offset; got != want {
			t.Errorf("%s, skill %d of %d: offset %v, want %v", s.ID, k, len(skills), got, want)
		}
	}
}

// G8: with the technique taken, every item's chance is the logit lifted by 0.5; the technique comes
// in about one attempt in five.
func TestTechnique(t *testing.T) {
	a := Base()
	a.TechniqueChance, a.TechniqueLift = 0.2, 0.5
	l := learnerWith(t, "G8", 0, a, olympiad)
	tk := task(t, olympiad, "enum-34-d1-1")
	p0, _ := l.chance(tk, 0, Start, 0, 0)
	p1, _ := l.chance(tk, 0, Start, 0, 0.5)
	if p1 <= p0 {
		t.Errorf("the technique does not lift the chance: %v, %v", p0, p1)
	}
	taken := 0
	for k := range 500 {
		at, err := l.Answer(tk, Start+int64(k)*secondsPerDay)
		if err != nil {
			t.Fatal(err)
		}
		if at.Technique {
			taken++
		}
	}
	if taken < 70 || taken > 130 {
		t.Errorf("the technique in %d attempts of 500", taken)
	}
}

// Assumptions of the new fields that no learner can live by are refused: a jump with no window, a
// window turned round, a chance of the technique above 1, a negative fade.
func TestNewFieldsRefused(t *testing.T) {
	for _, change := range []func(*Assumptions){
		func(a *Assumptions) { a.JumpBy = 1 },
		func(a *Assumptions) { a.JumpBy, a.JumpFrom, a.JumpTo = 1, 150, 50 },
		func(a *Assumptions) { a.TechniqueChance = 1.5 },
		func(a *Assumptions) { a.LearnFade = -1 },
	} {
		a := Base()
		change(&a)
		if _, err := New(seed.Working, "X", 0, []*pack.Module{module(t, olympiad)}, a); !errors.Is(err, ErrAssumptions) {
			t.Errorf("%+v: %v", a, err)
		}
	}
}

// S-adult: no grade, the first band, a level 2.5 to 5 above it, the age category adult.
func TestAdult(t *testing.T) {
	a := Base()
	a.Adult, a.AdultFrom, a.AdultTo = true, 2.5, 5
	for n := range 40 {
		l := learnerWith(t, "S-adult", n, a, olympiad)
		m := l.modules[olympiad]
		// The params of an adult: no grade, the level, the three offsets, then its shift above.
		p, _ := seed.Working.Place(l.Unit, seed.Params)
		level := p.Normal()
		for range 3 {
			p.Normal()
		}
		shift := 2.5 + float64(2.5*p.Fraction())
		if m.Grade.Known || m.StartBand != 1 || l.Age != "adult" || m.Level != level+shift || shift < 2.5 || shift > 5 {
			t.Fatalf("an adult: %+v, age %q, want the level %v", m, l.Age, level+shift)
		}
	}
}

// F-: the calendar has one break, from a day between the 30th and the 60th, 14 to 60 days long and
// over by the 120th, with no session in it; after it, F-break's level falls by 1 once.
func TestBreak(t *testing.T) {
	a := Base()
	a.SessionChance = 3.0 / 7
	a.BreakStartFrom, a.BreakStartTo, a.BreakMin, a.BreakMax, a.BreakEndBy, a.BreakDrop = 30, 60, 14, 60, 120, 1
	for n := range 30 {
		l := learnerWith(t, "F-break", n, a, olympiad)
		c := l.NewCalendar()
		if c.BreakStart < 30 || c.BreakStart > 60 || c.BreakEnd-c.BreakStart < 14 || c.BreakEnd-c.BreakStart > 60 || c.BreakEnd > 120 {
			t.Fatalf("a break from day %d to %d", c.BreakStart, c.BreakEnd)
		}
		level := l.modules[olympiad].Level
		tk := task(t, olympiad, "enum-34-d1-1")
		for range 200 {
			at := c.Next()
			if c.Day() >= c.BreakStart && c.Day() < c.BreakEnd {
				t.Fatalf("a session on day %d, in the break %d–%d", c.Day(), c.BreakStart, c.BreakEnd)
			}
			if _, err := l.Answer(tk, at); err != nil {
				t.Fatal(err)
			}
			if l.modules[olympiad].Level != level || c.After(at) != (c.Day() >= c.BreakEnd) {
				t.Fatalf("day %d, the break ending %d: level %v, after %v", c.Day(), c.BreakEnd, l.modules[olympiad].Level, c.After(at))
			}
			p, _ := l.Chance(tk, 0, at, 0)
			l.a.BreakDrop = 0
			noDrop, _ := l.Chance(tk, 0, at, 0)
			l.a.BreakDrop = 1
			if (p < noDrop) != c.After(at) {
				t.Fatalf("day %d: P* %v with the fall, %v without", c.Day(), p, noDrop)
			}
			if l.NewCalendar() != c {
				t.Fatal("a second calendar")
			}
			c.Done(at + 120)
		}
	}
}

// F-exp: the law R* = 0.9^(Δt/τ*).
func TestExponentialLaw(t *testing.T) {
	a := Base()
	a.ExponentialLaw, a.LawBase = true, 0.9
	l := learnerWith(t, "F-exp", 0, a, olympiad)
	s := l.skills[task(t, olympiad, "enum-34-d1-1").Candidate.Task.Primary[0]]
	s.Stability, s.Last, s.Practised = 2, Start, true
	if got, want := l.Retention(s, Start+6*secondsPerDay), 0.9*0.9*0.9; math.Abs(got-want) > 1e-15 {
		t.Errorf("R* = %v, want %v", got, want)
	}
}

// W-shock, W-error, W-badday and W-repeat change an attempt as their tables say.
func TestPageAssumptions(t *testing.T) {
	tk := tasks(t, grade3)[0]
	t.Run("W-error", func(t *testing.T) {
		a := Base()
		a.ErrorChance = 1
		l := learnerWith(t, "W-error", 0, a, grade3)
		at, err := l.Answer(tk, Start)
		if err != nil {
			t.Fatal(err)
		}
		for _, res := range at.Items {
			if res.FirstRight && hasPrimary(tk, int(res.Place-1)) {
				t.Fatalf("item %d right at first in an attempt that holds an error", res.Place)
			}
		}
	})
	t.Run("W-shock", func(t *testing.T) {
		a := Base()
		a.AttemptShock = 1
		l := learnerWith(t, "W-shock", 0, a, grade3)
		shocks, _ := seed.Working.Place(l.Unit, seed.Attempts)
		rights := map[bool]int{}
		for k := range 50 {
			shock := shocks.Normal()
			at, err := l.Answer(tk, Start+int64(k)*secondsPerDay)
			if err != nil {
				t.Fatal(err)
			}
			for _, res := range at.Items {
				if res.FirstRight {
					rights[shock > 0]++
				}
			}
		}
		if rights[true] <= rights[false] {
			t.Errorf("right first answers after a shock up %d, down %d", rights[true], rights[false])
		}
	})
	t.Run("W-badday", func(t *testing.T) {
		a := Base()
		a.BadDayChance, a.BadDayDrop = 1, 1
		l := learnerWith(t, "W-badday", 0, a, grade3)
		c := l.NewCalendar()
		at := c.Next()
		if !c.Bad(at) {
			t.Fatal("a day of a chance 1 of being bad is not")
		}
		bad, _ := l.Chance(tk, 0, at, 0)
		c.bad = map[int]bool{}
		good, _ := l.Chance(tk, 0, at, 0)
		if bad >= good {
			t.Errorf("a bad day does not lower P*: %v, %v", bad, good)
		}
	})
	t.Run("W-repeat", func(t *testing.T) {
		a := Base()
		a.RepeatChance, a.RepeatDays = 1, 30
		l := learnerWith(t, "W-repeat", 0, a, grade3)
		if _, err := l.Answer(tk, Start); err != nil {
			t.Fatal(err)
		}
		at, err := l.Answer(tk, Start+10*secondsPerDay)
		if err != nil {
			t.Fatal(err)
		}
		for _, res := range at.Items {
			if !res.FirstRight {
				t.Fatalf("item %d not right on a repeat of chance 1", res.Place)
			}
		}
	})
}

// S-unreadable: every twentieth answer of an item on a sheet is unreadable first; none on a
// workbook.
func TestUnreadable(t *testing.T) {
	a := Base()
	a.UnreadableEvery = 20
	l := learnerWith(t, "S-unreadable", 0, a, olympiad)
	tk := task(t, olympiad, "enum-34-d1-1")
	// Every answer on the sheet counts, the first and those after a rung: the task has one item, so
	// an attempt gives 1 + its rungs answers.
	counted, marked := 0, 0
	for k := 1; k <= 60; k++ {
		at, err := l.Answer(tk, Start+int64(k)*secondsPerDay)
		if err != nil {
			t.Fatal(err)
		}
		if got := at.Items[0].Unreadable; got != ((counted+1)%20 == 0) {
			t.Fatalf("attempt %d, answer %d: unreadable %v", k, counted+1, got)
		}
		if at.Items[0].Unreadable {
			marked++
		}
		counted += 1 + at.Rungs
	}
	if marked == 0 {
		t.Error("no first answer unreadable in 60 attempts")
	}
	w := learnerWith(t, "S-unreadable", 0, a, grade3)
	for k := range 30 {
		at, _ := w.Answer(tasks(t, grade3)[0], Start+int64(k)*secondsPerDay)
		for _, res := range at.Items {
			if res.Unreadable {
				t.Fatal("an unreadable answer on a workbook")
			}
		}
	}
}

// The bounds of the fields of S-, F- and W-, and an adult where no module has adults.
func TestOwnFieldsRefused(t *testing.T) {
	for _, change := range []func(*Assumptions){
		func(a *Assumptions) { a.ErrorChance = 2 },
		func(a *Assumptions) { a.BadDayChance = -0.1 },
		func(a *Assumptions) { a.RepeatChance = 95 },
		func(a *Assumptions) { a.RepeatDays = -1 },
		func(a *Assumptions) { a.UnreadableEvery = -20 },
		func(a *Assumptions) { a.AttemptShock = -1 },
		func(a *Assumptions) { a.ExponentialLaw, a.LawBase = true, 1.5 },
		func(a *Assumptions) { a.Adult, a.AdultFrom, a.AdultTo = true, 5, 2.5 },
		func(a *Assumptions) {
			a.BreakMax, a.BreakMin, a.BreakStartFrom, a.BreakStartTo, a.BreakEndBy = 60, 14, 30, 60, 50
		},
	} {
		a := Base()
		change(&a)
		if _, err := New(seed.Working, "X", 0, []*pack.Module{module(t, olympiad)}, a); !errors.Is(err, ErrAssumptions) {
			t.Errorf("%+v: %v", a, err)
		}
	}
	a := Base()
	a.Adult, a.AdultFrom, a.AdultTo = true, 2.5, 5
	if _, err := New(seed.Working, "S-adult", 0, []*pack.Module{module(t, grade3)}, a); !errors.Is(err, ErrAssumptions) {
		t.Errorf("an adult in school maths for grade 3: %v", err)
	}
}
