package learner

import (
	"testing"

	"github.com/altnovatrail/sim/internal/seed"
)

// A day is a session exactly when its fraction of the place calendar is below 5/7; a session is ten
// attempts in a row from 15:00 UTC, the next given when the last closed.
func TestCalendar(t *testing.T) {
	l := learnerOf(t, "G0", 8, olympiad)
	c := l.NewCalendar()
	days, _ := seed.Working.Place(l.Unit, seed.Calendar)
	sessions := 0
	day := -1
	for range 300 {
		for {
			day++
			if days.Fraction() < 5.0/7 {
				break
			}
		}
		at := Start + int64(day)*secondsPerDay + 15*3600
		for k := range 10 {
			got := c.Next()
			if got != at || c.Day() != day {
				t.Fatalf("session %d, attempt %d: given at %d on day %d, want %d on day %d", sessions, k, got, c.Day(), at, day)
			}
			if again := c.Next(); again != got {
				t.Fatal("Next moved without Done")
			}
			at += 600
			c.Done(at)
		}
		sessions++
	}
	if share := float64(float64(sessions) / float64(day+1)); share < 0.66 || share > 0.77 {
		t.Errorf("%d sessions in %d days", sessions, day+1)
	}
	b := learnerOf(t, "G0", 8, olympiad).NewCalendar()
	if b.Next() != learnerOf(t, "G0", 8, olympiad).NewCalendar().Next() {
		t.Error("one unit, two calendars")
	}
}
