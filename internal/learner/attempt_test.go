package learner

import (
	"errors"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"

	"github.com/altnovatrail/learner-model/detmath"
	"github.com/altnovatrail/sim/internal/seed"
)

// An attempt draws one fraction of the place answers for each item's first answer, whatever came
// before, and the answers after a rung from retries: the k-th first answer of a learner meets the
// same numbers under any rule (docs/architecture.md, section 3.4). The first answer of each item is
// right exactly when its fraction is below P*.
func TestAnswerFirst(t *testing.T) {
	l := learnerOf(t, "G0", 5, olympiad)
	shadow := learnerOf(t, "G0", 5, olympiad)
	answers, _ := seed.Working.Place(l.Unit, seed.Answers)
	at := Start
	for _, tk := range tasks(t, olympiad) {
		var ps []float64
		for i := range tk.Candidate.Task.Items {
			p, err := shadow.Chance(tk, i, at, 0)
			if err != nil {
				t.Fatal(err)
			}
			ps = append(ps, p)
		}
		a, err := l.Answer(tk, at)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := shadow.Answer(tk, at); err != nil {
			t.Fatal(err)
		}
		for _, res := range a.Items {
			if f := answers.Fraction(); res.FirstRight != (f < ps[res.Place-1]) {
				t.Fatalf("%s item %d: first right %v with fraction %v and P* %v", a.Task, res.Place, res.FirstRight, f, ps[res.Place-1])
			}
		}
		if a.End != a.Start+SheetItemSeconds*int64(len(a.Items)) {
			t.Fatalf("%s: from %d to %d", a.Task, a.Start, a.End)
		}
		at = a.End + 3*secondsPerDay
	}
}

// How an attempt closes: right at once with no rung, right after rungs, or the solution shown with
// the fourth wrong answer; the ladder lifts P*.
func TestAnswerClosure(t *testing.T) {
	seen := map[Closure]bool{}
	rungsSeen := map[int]bool{}
	for n := range 200 {
		l := learnerOf(t, "G0", n, olympiad)
		tk := task(t, olympiad, "pig-56-d3-2")
		a, err := l.Answer(tk, Start)
		if err != nil {
			t.Fatal(err)
		}
		seen[a.Closure] = true
		rungsSeen[a.Rungs] = true
		res := a.Items[0]
		switch {
		case res.FirstRight && (a.Rungs != 0 || a.Closure != ClosedRight || res.Wrong != 0):
			t.Fatalf("right at once: %+v", a)
		case !res.FirstRight && a.Closure == ClosedRight && (res.Wrong != a.Rungs || !res.Right):
			t.Fatalf("right after %d rungs: %+v", a.Rungs, a)
		case a.Closure == ClosedSolutionShown && (res.Wrong != 4 || a.Rungs != 3 || res.Right):
			t.Fatalf("the solution shown: %+v", a)
		}
	}
	if !seen[ClosedRight] || !seen[ClosedSolutionShown] || !rungsSeen[0] || !rungsSeen[1] {
		t.Errorf("closures %v, rungs %v", seen, rungsSeen)
	}
	l := learnerOf(t, "G0", 1, olympiad)
	tk := task(t, olympiad, "enum-34-d1-1")
	p0, _ := l.Chance(tk, 0, Start, 0)
	p3, _ := l.Chance(tk, 0, Start, 3)
	if p3 <= p0 {
		t.Errorf("the ladder does not lift P*: %v, %v", p0, p3)
	}
}

// A task with a rubric is written, sent once and done, in fifteen minutes, with no verdict; a line
// of a workbook takes a minute.
func TestAnswerRubricAndWorkbook(t *testing.T) {
	l := learnerOf(t, "E-base", 1, essays)
	for _, tk := range tasks(t, essays) {
		if !tk.Rubric {
			continue
		}
		a, err := l.Answer(tk, Start)
		if err != nil {
			t.Fatal(err)
		}
		if a.Closure != ClosedDone || len(a.Items) != 0 || a.End != Start+EssaySeconds {
			t.Fatalf("a task with a rubric: %+v", a)
		}
	}
	w := learnerOf(t, "G0", 1, grade3)
	tk := tasks(t, grade3)[0]
	a, err := w.Answer(tk, Start)
	if err != nil {
		t.Fatal(err)
	}
	if a.End != Start+WorkbookLineSeconds*int64(len(a.Items)) || len(a.Items) != len(tk.Candidate.Task.Items) {
		t.Fatalf("a workbook: %d items from %d to %d", len(a.Items), a.Start, a.End)
	}
	if _, err := learnerOf(t, "G0", 1, olympiad).Answer(tk, Start); err == nil {
		t.Error("an answer to a module the learner does not learn: no error")
	}
}

// The stability moves once an attempt, by the share of right first verdicts: it grows after right
// answers once retention has fallen, halves after wrong ones, and never goes below a day.
func TestPractise(t *testing.T) {
	l := learnerOf(t, "G0", 6, olympiad)
	tk := task(t, olympiad, "enum-34-d1-1")
	s := l.skills[tk.Candidate.Task.Primary[0]]
	s.Stability, s.Last, s.Practised = 4, Start, true
	at := Start + 10*secondsPerDay
	r := l.Retention(s, at)
	a, err := l.Answer(tk, at)
	if err != nil {
		t.Fatal(err)
	}
	share := 0.0
	if a.Items[0].FirstRight {
		share = 1
	}
	want := max(1, float64(float64(4*detmath.Pow(1+float64(15*(1-r)), share))*detmath.Pow(0.5, 1-share)))
	if s.Stability != want || s.Last != at {
		t.Errorf("τ* %v, want %v; last %d", s.Stability, want, s.Last)
	}
	s.Stability = 1
	for range 5 {
		if _, err := l.Answer(tk, at); err != nil {
			t.Fatal(err)
		}
		if s.Stability < 1 {
			t.Fatalf("τ* %v below a day", s.Stability)
		}
	}
}

// A lesson is passed link by link in the author's order: its links with no task take two minutes
// each, and each task is answered as any task is.
func TestPassLesson(t *testing.T) {
	m := module(t, olympiad)
	l := learnerOf(t, "L-explain", 1, olympiad)
	for i := range m.Lessons {
		lesson := &m.Lessons[i]
		p, err := l.PassLesson(m, lesson, Start)
		if err != nil {
			t.Fatal(err)
		}
		if len(p.Attempts) != len(lesson.Lesson.Tasks) || p.Lesson != lesson.Lesson.ID {
			t.Fatalf("%s: %d attempts for %d tasks", p.Lesson, len(p.Attempts), len(lesson.Lesson.Tasks))
		}
		before := 0 // the links before the first task, in the author's order
		for _, i := range lesson.Order {
			if i >= 0 {
				break
			}
			before++
		}
		first := Start + int64(before)*LinkSeconds
		after := 0 // the links after the last task
		for k := len(lesson.Order) - 1; k >= 0 && lesson.Order[k] < 0; k-- {
			after++
		}
		if len(p.Attempts) > 0 && (p.Attempts[0].Start != first || p.End != p.Attempts[len(p.Attempts)-1].End+int64(after)*LinkSeconds) {
			t.Fatalf("%s: the first task at %d, want %d; the end %d", p.Lesson, p.Attempts[0].Start, first, p.End)
		}
	}
	other := learnerOf(t, "G0", 1, grade3)
	if _, err := other.PassLesson(m, &m.Lessons[0], Start); err == nil {
		t.Error("a lesson of a module the learner does not learn: no error")
	}
}

// What Answer refuses before it draws anything: a task of a module the learner does not learn, a
// rubric one too, and a task with nothing to answer; the attempt seeds stay where they were.
func TestAnswerRefuses(t *testing.T) {
	l := learnerOf(t, "G0", 1, olympiad)
	seeds, _ := seed.Working.Place(l.Unit, seed.Attempt)
	for _, tk := range tasks(t, essays) {
		if tk.Rubric {
			if _, err := l.Answer(tk, Start); err == nil {
				t.Error("a rubric task of a module the learner does not learn: no error")
			}
			break
		}
	}
	empty := *task(t, olympiad, "enum-34-d1-1")
	empty.Candidate.Task.Items = []learnermodel.TaskItem{{Place: 1}}
	if _, err := l.Answer(&empty, Start); !errors.Is(err, ErrNoAnswer) {
		t.Errorf("a task with nothing to answer: %v", err)
	}
	a, err := l.Answer(task(t, olympiad, "enum-34-d1-1"), Start)
	if err != nil || a.Seed != seeds.Integer() {
		t.Errorf("the refused attempts moved the seeds: %v", err)
	}
}
