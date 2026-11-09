package learner

import (
	"errors"
	"math"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/learner-model/detmath"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/seed"
)

const (
	olympiad = "novatrail/olympiad-math"
	grade3   = "novatrail/grade-3-math"
	essays   = "novatrail/english-essays"
)

func module(t *testing.T, id string) *pack.Module {
	t.Helper()
	ms, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	return ms[id]
}

// tasks gives the tasks of a published module, each by its place in the module.
func tasks(t *testing.T, id string) []*pack.Task {
	t.Helper()
	m := module(t, id)
	var out []*pack.Task
	for i := range m.Tasks {
		out = append(out, &m.Tasks[i])
	}
	return out
}

func learnerOf(t *testing.T, population string, n int, ids ...string) *Learner {
	t.Helper()
	var ms []*pack.Module
	for _, id := range ids {
		ms = append(ms, module(t, id))
	}
	l, err := New(seed.Working, population, n, ms, Base())
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func task(t *testing.T, mod, local string) *pack.Task {
	t.Helper()
	task, err := module(t, mod).Task(mod + "/task/" + local)
	if err != nil {
		t.Fatal(err)
	}
	return task
}

// The truth of a learner is its places': the same unit gives the same truth, and every grade of
// olympiad maths comes, each with the start band whose grades name it, the level around it.
func TestNew(t *testing.T) {
	a, b := learnerOf(t, "G0", 7, olympiad), learnerOf(t, "G0", 7, olympiad)
	ma, _ := a.Module(olympiad)
	mb, _ := b.Module(olympiad)
	if *ma != *mb || a.Age != "under_13" {
		t.Fatalf("one unit, two truths: %+v, %+v; age %q", ma, mb, a.Age)
	}
	grades := map[int32]int{}
	var sum, sum2 float64
	const n = 600
	for i := range n {
		l := learnerOf(t, "G0", i, olympiad)
		m, _ := l.Module(olympiad)
		if !m.Grade.Known || m.Grade.Value < 1 || m.Grade.Value > 6 {
			t.Fatalf("grade %+v", m.Grade)
		}
		grades[m.Grade.Value]++
		if want := (m.Grade.Value + 1) / 2; m.StartBand != want {
			t.Fatalf("grade %d starts in band %d, want %d", m.Grade.Value, m.StartBand, want)
		}
		dev := m.Level - float64(2.5*float64(m.StartBand-1))
		sum += dev
		sum2 += float64(dev * dev)
		if len(l.skills) != 19 {
			t.Fatalf("%d skills with a verdict, want 19", len(l.skills))
		}
	}
	for g := int32(1); g <= 6; g++ {
		if grades[g] < n/6-40 || grades[g] > n/6+40 {
			t.Errorf("grade %d came %d times of %d", g, grades[g], n)
		}
	}
	if mean, sd := float64(sum/n), math.Sqrt(float64(sum2/n)); math.Abs(mean) > 0.15 || math.Abs(sd-1) > 0.1 {
		t.Errorf("the level around the start band: mean %v, spread %v", mean, sd)
	}
	if l := learnerOf(t, "G0", 1, grade3); l.modules[grade3].Grade != (learnermodel.Grade{Value: 3, Known: true}) || l.modules[grade3].StartBand != 1 {
		t.Errorf("school maths for grade 3: %+v", l.modules[grade3])
	}
	ages := map[string]int{}
	for i := range 200 {
		l := learnerOf(t, "E-base", i, essays)
		if l.modules[essays].Grade.Known || l.modules[essays].StartBand != 1 {
			t.Fatalf("essays: %+v", l.modules[essays])
		}
		ages[l.Age]++
		if len(l.skills) != 1 {
			t.Fatalf("essays have %d skills with a verdict, want 1", len(l.skills))
		}
	}
	if len(ages) != 4 {
		t.Errorf("essays give the ages %v; want each of the module's four", ages)
	}
	if _, err := New(seed.Working, "G0", 0, nil, Base()); err != ErrNoModule {
		t.Errorf("New with no module: %v", err)
	}
	if _, err := a.Module(essays); err == nil {
		t.Error("a module the learner does not learn: no error")
	}
	if a.Skill("no/skill/x") != nil {
		t.Error("a skill the learner has not")
	}
}

// The author's miss is the task's: two learners of a set see it alike, and another set draws its
// own.
func TestDifficulty(t *testing.T) {
	tk := task(t, olympiad, "enum-34-d1-1")
	a, b := learnerOf(t, "G0", 1, olympiad), learnerOf(t, "G2-half", 9, olympiad)
	if a.Difficulty(tk) != b.Difficulty(tk) {
		t.Errorf("two learners see two misses: %v, %v", a.Difficulty(tk), b.Difficulty(tk))
	}
	s, _ := seed.Working.Place(seed.Task(olympiad, "enum-34-d1-1"), seed.Miss)
	z := s.Normal()
	if want := float64(2.5*float64(tk.Candidate.Task.Band.Place-1)) + float64(tk.Candidate.Task.Difficulty-3) + float64(0.5*z); a.Difficulty(tk) != want {
		t.Errorf("β* = %v, want %v", a.Difficulty(tk), want)
	}
	held, _ := seed.HeldOut(seed.CriterionModel, 1)
	h, err := New(held, "G0", 1, []*pack.Module{module(t, olympiad)}, Base())
	if err != nil {
		t.Fatal(err)
	}
	if h.Difficulty(tk) == a.Difficulty(tk) {
		t.Error("another set draws the same miss")
	}
}

// P* by its formula, on an item of a number, a choice of three options and with a rung.
func TestChance(t *testing.T) {
	l := learnerOf(t, "G0", 3, grade3)
	m := l.modules[grade3]
	for _, tk := range tasks(t, grade3) {
		for i, it := range tk.Candidate.Task.Items {
			primary := tk.Candidate.Task.Primary
			if it.OwnSkills {
				primary = it.Primary
			}
			var off float64
			for _, id := range primary {
				off += l.skills[id].Offset
			}
			off = float64(off / float64(len(primary)))
			c := 1.0
			for _, s := range it.Slots {
				if s.Check == learnermodel.CheckKindChoice {
					c = float64(c / float64(s.Options))
				} else {
					c = 0
				}
			}
			for rung := range 4 {
				got, err := l.Chance(tk, i, Start, rung)
				if err != nil {
					t.Fatal(err)
				}
				want := c + float64((1-c)*detmath.Sigmoid(float64(m.Level+off-l.Difficulty(tk))+float64(0.5*float64(rung))))
				if math.Abs(got-want) > 1e-15 {
					t.Fatalf("%s item %d rung %d: P* %v, want %v", tk.Candidate.Task.ID, i+1, rung, got, want)
				}
			}
		}
	}
	other := task(t, olympiad, "enum-34-d1-1")
	if _, err := l.Chance(other, 0, Start, 0); err == nil {
		t.Error("the chance of a task of a module the learner does not learn: no error")
	}
}

func TestGuess(t *testing.T) {
	l := learnerOf(t, "G0", 1, grade3)
	l.a.ChoiceGuess, l.a.PredicateGuess = 1.5, 0.1
	choice := learnermodel.Slot{Check: learnermodel.CheckKindChoice, Options: 3, HasOptions: true}
	cases := []struct {
		slots []learnermodel.Slot
		want  float64
	}{
		{[]learnermodel.Slot{choice}, 0.5},
		{[]learnermodel.Slot{choice, choice}, 0.25},
		{[]learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 1, HasOptions: true}}, 1},
		{[]learnermodel.Slot{{Check: learnermodel.CheckKindPredicate}}, 0.1},
		{[]learnermodel.Slot{choice, {Check: learnermodel.CheckKindNumber}}, 0},
	}
	for _, c := range cases {
		if got := l.guess(c.slots); math.Abs(got-c.want) > 1e-15 {
			t.Errorf("guess(%v) = %v, want %v", c.slots, got, c.want)
		}
	}
}

// R* after an attempt: 1 at once, the law after a pause, and 1 always where the population does
// not forget.
func TestRetention(t *testing.T) {
	l := learnerOf(t, "G0", 2, olympiad)
	tk := task(t, olympiad, "enum-34-d1-1")
	id := tk.Candidate.Task.Primary[0]
	s := l.skills[id]
	if l.Retention(s, Start) != 1 {
		t.Error("R* before the first attempt is not 1")
	}
	if _, err := l.Answer(tk, Start); err != nil {
		t.Fatal(err)
	}
	if !s.Practised || s.Last != Start || l.Retention(s, Start) != 1 {
		t.Fatalf("after an attempt: %+v", s)
	}
	at := Start + 7*secondsPerDay
	want := detmath.Pow(1+float64(float64(19.0/81*7)/s.Stability), -0.5)
	if got := l.Retention(s, at); math.Abs(got-want) > 1e-15 || got >= 1 {
		t.Errorf("R* after 7 days = %v, want %v", got, want)
	}
	l.a.Forgets = false
	if l.Retention(s, at) != 1 {
		t.Error("a learner who does not forget forgets")
	}
}

// A learner of both maths has one grade, the one both name, and the start band of each module by
// it; assumptions no learner can live by are refused.
func TestNewShared(t *testing.T) {
	for n := range 20 {
		l := learnerOf(t, "A-linked", n, olympiad, grade3)
		o, g := l.modules[olympiad], l.modules[grade3]
		if o.Grade != (learnermodel.Grade{Value: 3, Known: true}) || g.Grade != o.Grade || o.StartBand != 2 || g.StartBand != 1 || l.Age != "under_13" {
			t.Fatalf("both maths: %+v, %+v, %q", o, g, l.Age)
		}
	}
	for _, a := range []Assumptions{
		func() Assumptions { a := Base(); a.SessionChance = 0; return a }(),
		func() Assumptions { a := Base(); a.SessionAttempts = 0; return a }(),
	} {
		if _, err := New(seed.Working, "G0", 0, []*pack.Module{module(t, olympiad)}, a); !errors.Is(err, ErrAssumptions) {
			t.Errorf("assumptions %+v: %v", a, err)
		}
	}
	if _, err := New(seed.Working, "E-base", 0, []*pack.Module{module(t, essays), module(t, grade3)}, Base()); err != nil {
		t.Errorf("essays with school maths: %v", err)
	}
	g3 := *module(t, grade3)
	g3.Graph.Bands = []learnermodel.Band{{Place: 1, ID: "grade-7", Grades: []int32{7}}}
	if _, err := New(seed.Working, "G0", 0, []*pack.Module{module(t, olympiad), &g3}, Base()); !errors.Is(err, ErrAssumptions) {
		t.Errorf("modules with no grade in common: %v", err)
	}
	e := *module(t, essays)
	e.Ages = []string{"no-such-age"}
	if _, err := New(seed.Working, "E-base", 0, []*pack.Module{module(t, essays), &e}, Base()); !errors.Is(err, ErrAssumptions) {
		t.Errorf("modules with no age category in common: %v", err)
	}
}
