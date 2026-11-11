package learner

import (
	"errors"
	"math"
	"slices"
	"testing"

	"github.com/altnovatrail/learner-model/detmath"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
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
// in about one attempt in five, and with it the leading question of a task that has one.
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
	// The same hint gives the leading question on a task that has one, with no lift and no number of
	// its own: the attempt is the one on the task without it, the question beside (R41 in
	// docs/decisions.md). A learner who asks for nothing gets none.
	leading, plainTask := *tk, *tk
	leading.Leading, plainTask.Leading = true, false
	with, without, base := learnerWith(t, "G8", 1, a, olympiad), learnerWith(t, "G8", 1, a, olympiad), learnerOf(t, "G8", 1, olympiad)
	asked := 0
	for k := range 100 {
		at := Start + int64(k)*secondsPerDay
		got, err := with.Answer(&leading, at)
		if err != nil {
			t.Fatal(err)
		}
		plain, err := without.Answer(&plainTask, at)
		if err != nil {
			t.Fatal(err)
		}
		if got.LeadingQuestion != got.Technique || plain.LeadingQuestion || got.Technique != plain.Technique || !slices.Equal(got.Items, plain.Items) {
			t.Fatalf("attempt %d: %+v with a leading question, %+v without", k, got, plain)
		}
		if b, err := base.Answer(&leading, at); err != nil || b.LeadingQuestion {
			t.Fatalf("attempt %d of the base: a leading question %v, %v", k, b.LeadingQuestion, err)
		}
		asked += b2i(got.LeadingQuestion)
	}
	if asked == 0 {
		t.Error("G8 never gets the leading question")
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
		// The params of an adult: no grade, the level, an offset of each skill with a verdict, then its
		// shift above.
		p, _ := seed.Working.Place(l.Unit, seed.Params)
		level := p.Normal()
		for range len(l.skills) {
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
			id := tk.Candidate.Task.Primary[0]
			truth, _, _ := l.Truth(id, at)
			ref, mastered := l.ReferenceAt(id, 1, at), l.TrulyMasteredAt(id, 1, at)
			lv2, _ := l.LevelAt(olympiad, at)
			l.a.BreakDrop = 0
			noDrop, _ := l.Chance(tk, 0, at, 0)
			truth0, _, _ := l.Truth(id, at)
			ref0, mastered0 := l.ReferenceAt(id, 1, at), l.TrulyMasteredAt(id, 1, at)
			l.a.BreakDrop = 1
			if (p < noDrop) != c.After(at) || (truth < truth0) != c.After(at) || (ref < ref0) != c.After(at) || (mastered && !mastered0) {
				t.Fatalf("day %d: P* %v with the fall, %v without; truth %v and %v", c.Day(), p, noDrop, truth, truth0)
			}
			if want := level - float64(b2i(c.After(at))); lv2 != want {
				t.Fatalf("day %d: the level %v, want %v", c.Day(), lv2, want)
			}
			if l.NewCalendar() != c {
				t.Fatal("a second calendar")
			}
			c.Done(at + 120)
		}
		if end, ok := l.DropAt(); !ok || end != Start+int64(c.BreakEnd)*secondsPerDay {
			t.Fatalf("the fall at %d, %v", end, ok)
		}
		l.a.BreakDrop = 0
		if _, ok := l.DropAt(); ok {
			t.Fatal("a break with no fall drops the level")
		}
	}
	if _, err := learnerOf(t, "G0", 0, olympiad).LevelAt(grade3, Start); err == nil {
		t.Error("the level of a module the learner does not learn")
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

func syntheticModules(t *testing.T) []*pack.Module {
	t.Helper()
	ms, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	return []*pack.Module{ms[synthetic.Sheet], ms[synthetic.Workbook]}
}

// A-linked: the levels of a learner's two modules are tied as ρ = 0.7 says; without it, they are not.
func TestLevelCorrelation(t *testing.T) {
	corr := func(rho float64) float64 {
		a := Base()
		a.LevelCorrelation = rho
		var sx, sy, sxx, syy, sxy float64
		const n = 600
		for k := range n {
			l, err := New(seed.Working, "A-linked", k, syntheticModules(t), a)
			if err != nil {
				t.Fatal(err)
			}
			x := l.modules[synthetic.Sheet].Level - float64(a.BandStep*float64(l.modules[synthetic.Sheet].StartBand-1))
			y := l.modules[synthetic.Workbook].Level
			sx, sy, sxx, syy, sxy = sx+x, sy+y, sxx+float64(x*x), syy+float64(y*y), sxy+float64(x*y)
		}
		cov := float64(sxy/n) - float64(float64(sx/n)*float64(sy/n))
		vx := float64(sxx/n) - float64(float64(sx/n)*float64(sx/n))
		vy := float64(syy/n) - float64(float64(sy/n)*float64(sy/n))
		return float64(cov / math.Sqrt(float64(vx*vy)))
	}
	if c := corr(0.7); math.Abs(c-0.7) > 0.08 {
		t.Errorf("tied at 0.7, the levels correlate %v", c)
	}
	if c := corr(0); math.Abs(c) > 0.12 {
		t.Errorf("untied, the levels correlate %v", c)
	}
}

// A-unpractised and A-fading: the skill that teaches the anchor is stronger by the boost; once the
// truth has mastered it, it grows its stability no more.
func TestTeacher(t *testing.T) {
	a := Base()
	a.TeacherBoost = 2
	l, err := New(seed.Working, "A-unpractised", 0, syntheticModules(t), a)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := New(seed.Working, "A-unpractised", 0, syntheticModules(t), Base())
	tables := synthetic.Workbook + "/skill/tables"
	if got := l.skills[tables].Offset - b.skills[tables].Offset; math.Abs(got-2) > 1e-12 {
		t.Errorf("the teaching skill is stronger by %v", got)
	}
	if got := l.skills[synthetic.Workbook+"/skill/compare"].Offset - b.skills[synthetic.Workbook+"/skill/compare"].Offset; got != 0 {
		t.Errorf("a skill that teaches nothing is stronger by %v", got)
	}
	if !l.TrulyMastered(tables, 1) {
		t.Fatalf("a skill 2 above its level is not mastered: P* %v", l.Reference(tables, 1))
	}
	// A-fading against the same learner without it: the same answers, a mastered teaching skill, and
	// the stability grows only without fading.
	ms, _ := synthetic.Modules()
	page, _ := ms[synthetic.Workbook].Task(synthetic.Workbook + "/task/tables-page-1")
	stability := func(fades bool) (float64, float64) {
		f := Base()
		f.TeacherFades = fades
		fl, _ := New(seed.Working, "A-fading", 0, syntheticModules(t), f)
		s := fl.skills[tables]
		s.Offset += 10 // the truth masters it
		s.Stability, s.Last, s.Practised = 2, Start, true
		at, err := fl.Answer(page, Start+20*secondsPerDay)
		if err != nil {
			t.Fatal(err)
		}
		right := 0
		for _, res := range at.Items {
			right += b2i(res.FirstRight)
		}
		return s.Stability, float64(float64(right) / float64(len(at.Items)))
	}
	fading, share := stability(true)
	plain, again := stability(false)
	if share != again || share == 0 {
		t.Fatalf("the two learners answered apart: %v, %v", share, again)
	}
	if want := max(1, float64(2*detmath.Pow(0.5, 1-share))); fading != want || plain <= fading {
		t.Errorf("the stability of a mastered teaching skill: %v with fading, want %v; %v without", fading, want, plain)
	}
}

// The reference task: 1/5 + 4/5 · σ(θ* + δ* − β) on the band, and 0 for a skill with no truth.
func TestReference(t *testing.T) {
	l := learnerWith(t, "G0", 0, Base(), olympiad)
	id := task(t, olympiad, "enum-34-d1-1").Candidate.Task.Primary[0]
	s := l.skills[id]
	want := 0.2 + float64(0.8*detmath.Sigmoid(s.Module.Level+s.Offset-2.5))
	if got := l.Reference(id, 2); math.Abs(got-want) > 1e-15 {
		t.Errorf("Reference = %v, want %v", got, want)
	}
	if l.Reference("no/skill/x", 1) != 0 || l.TrulyMastered("no/skill/x", 1) {
		t.Error("a skill with no truth is mastered")
	}
}

// L-abandon, L-explain and E-leave.
func TestLessonsAndEssays(t *testing.T) {
	m := module(t, olympiad)
	a := Base()
	a.AbandonChance = 1
	l := learnerWith(t, "L-abandon", 0, a, olympiad)
	stops, _ := seed.Working.Place(l.Unit, seed.Attempts)
	for i := range m.Lessons {
		lesson := &m.Lessons[i]
		stops.Fraction()
		link := int(stops.Below(uint64(lesson.Links)))
		p, err := l.PassLesson(m, lesson, Start, nil)
		if err != nil {
			t.Fatal(err)
		}
		want := 0 // the tasks among the links before the one it is left at, in the author's order
		for _, i := range lesson.Order[:link] {
			if i >= 0 {
				want++
			}
		}
		if !p.Abandoned || len(p.Attempts) != want {
			t.Fatalf("%s left at link %d of %d: abandoned %v, %d attempts, want %d", p.Lesson, link, lesson.Links, p.Abandoned, len(p.Attempts), want)
		}
	}
	e := Base()
	e.ExplainLift = 0.25
	x := learnerWith(t, "L-explain", 0, e, olympiad)
	lesson := &m.Lessons[0]
	id := lesson.Lesson.Skills[0]
	before := x.skills[id].Offset
	if _, err := x.PassLesson(m, lesson, Start, nil); err != nil {
		t.Fatal(err)
	}
	once := x.skills[id].Offset
	if _, err := x.PassLesson(m, lesson, Start+secondsPerDay, nil); err != nil {
		t.Fatal(err)
	}
	if math.Abs(once-before-0.25) > 1e-12 || x.skills[id].Offset != once {
		t.Errorf("the first pass lifted the skill by %v, the second by %v", once-before, x.skills[id].Offset-once)
	}
	// The next pass of a lesson left goes on from the link where it was left.
	long := &m.Lessons[0]
	r := learnerWith(t, "L-abandon", 1, a, olympiad)
	first, err := r.PassLesson(m, long, Start, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.a.AbandonChance = 0
	second, err := r.PassLesson(m, long, Start+secondsPerDay, nil)
	if err != nil {
		t.Fatal(err)
	}
	if second.Abandoned || len(first.Attempts)+len(second.Attempts) != len(long.Lesson.Tasks) {
		t.Errorf("left with %d attempts, gone on with %d, of %d tasks", len(first.Attempts), len(second.Attempts), len(long.Lesson.Tasks))
	}
	// The explanations teach only when a pass reaches one: not a lesson of tasks alone.
	tasksOnly := *lesson
	tasksOnly.Lesson.ID, tasksOnly.Order = "x/lesson/tasks-only", []int{0}
	tasksOnly.Lesson.Tasks = lesson.Lesson.Tasks[:1]
	y := learnerWith(t, "L-explain", 1, e, olympiad)
	before = y.skills[id].Offset
	if _, err := y.PassLesson(m, &tasksOnly, Start, nil); err != nil {
		t.Fatal(err)
	}
	if y.skills[id].Offset != before {
		t.Errorf("a lesson of tasks alone taught by %v", y.skills[id].Offset-before)
	}
	lv := Base()
	lv.LeaveChance = 1
	w := learnerWith(t, "E-leave", 0, lv, essays)
	for _, tk := range tasks(t, essays) {
		at, err := w.Answer(tk, Start)
		if err != nil {
			t.Fatal(err)
		}
		if tk.Rubric != (at.Closure == ClosedLeft) {
			t.Fatalf("%s: rubric %v, closed %v", tk.Candidate.Task.ID, tk.Rubric, at.Closure)
		}
	}
}
