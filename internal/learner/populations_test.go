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
