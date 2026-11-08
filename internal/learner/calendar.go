package learner

import "github.com/altnovatrail/sim/internal/seed"

// Start is day 0 of every calendar, 2026-09-07 00:00:00 UTC, a Monday: every learner starts on
// the same day (docs/architecture.md, section 3.1). The date is the bench's; only the days after
// it matter.
const Start int64 = 1788739200

// sessionStart is the time of day a session begins, 15:00 UTC.
const sessionStart = 15 * 3600

// Calendar is a learner's days, sessions and the times its attempts are given: a day is a session
// with the population's chance, one fraction of the place calendar each day in order, and a session
// is the population's number of attempts in a row from 15:00 UTC, the next given at once
// (docs/architecture.md, section 3.1). Times are whole seconds UTC, as in evidence (R08
// learner-model).
type Calendar struct {
	s    *seed.Stream
	a    Assumptions
	day  int   // the day of the session under way, or before the first, -1
	left int   // the attempts left in it
	next int64 // when its next attempt is given
}

// NewCalendar gives the calendar of the learner's place calendar.
func (l *Learner) NewCalendar() *Calendar {
	s, err := l.set.Place(l.Unit, seed.Calendar)
	if err != nil {
		panic(err) // the unit is the bench's own
	}
	return &Calendar{s: s, a: l.a, day: -1}
}

// Next gives when the next attempt is given: in the session under way while it has attempts
// left, else at the start of the next day that has a session.
func (c *Calendar) Next() int64 {
	if c.left == 0 {
		for {
			c.day++
			if c.s.Fraction() < c.a.SessionChance {
				break
			}
		}
		c.left = c.a.SessionAttempts
		c.next = Start + int64(c.day)*secondsPerDay + sessionStart
	}
	return c.next
}

// Done tells the calendar that the attempt given at Next closed at end: the next is given then.
func (c *Calendar) Done(end int64) {
	c.next = end
	c.left--
}

// Day gives the day of the session under way, from 0.
func (c *Calendar) Day() int {
	return c.day
}
