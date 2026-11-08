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
	// The break, days from BreakStart to BreakEnd, the last not included; none when they are equal.
	BreakStart, BreakEnd int
	bad                  map[int]bool // the bad days with a session
}

// NewCalendar gives the calendar of the learner's place calendar, made once and kept for the
// learner's bad days and its break. A population with a break draws it first: its start, then its length,
// evenly, so that it ends by BreakEndBy (docs/architecture.md, section 3.2).
func (l *Learner) NewCalendar() *Calendar {
	if l.cal != nil {
		return l.cal // a learner has one calendar: its bad days and its break are drawn once
	}
	s, err := l.set.Place(l.Unit, seed.Calendar)
	if err != nil {
		panic(err) // the unit is the bench's own
	}
	c := &Calendar{s: s, a: l.a, day: -1, bad: map[int]bool{}}
	if l.a.BreakMax > 0 {
		c.BreakStart = int(s.Between(int64(l.a.BreakStartFrom), int64(l.a.BreakStartTo)))
		c.BreakEnd = c.BreakStart + int(s.Between(int64(l.a.BreakMin), int64(min(l.a.BreakMax, l.a.BreakEndBy-c.BreakStart))))
	}
	l.cal = c
	return c
}

// After tells whether a time is past the end of the break (F-break).
func (c *Calendar) After(at int64) bool {
	return c.BreakEnd > c.BreakStart && at >= Start+int64(c.BreakEnd)*secondsPerDay
}

// Bad tells whether the day of a time is a bad day with a session (W-badday).
func (c *Calendar) Bad(at int64) bool {
	return c.bad[int((at-Start)/secondsPerDay)]
}

// Next gives when the next attempt is given: in the session under way while it has attempts
// left, else at the start of the next day that has a session.
func (c *Calendar) Next() int64 {
	if c.left == 0 {
		for {
			c.day++
			if c.day >= c.BreakStart && c.day < c.BreakEnd {
				continue
			}
			if c.s.Fraction() < c.a.SessionChance {
				break
			}
		}
		if c.a.BadDayChance > 0 {
			c.bad[c.day] = c.s.Fraction() < c.a.BadDayChance
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

// EndSession ends the session under way: the next attempt is given on the next day that has a
// session. The run calls it when the engine has nothing to give (R16 in docs/decisions.md).
func (c *Calendar) EndSession() {
	c.left = 0
}

// DropAt gives the time the learner's break ends and its level falls, and whether it falls: F-break
// (M3 after a break). The break of the other F- populations drops nothing.
func (l *Learner) DropAt() (int64, bool) {
	if l.cal == nil || l.cal.BreakEnd <= l.cal.BreakStart || l.a.BreakDrop == 0 {
		return 0, false
	}
	return Start + int64(l.cal.BreakEnd)*secondsPerDay, true
}
