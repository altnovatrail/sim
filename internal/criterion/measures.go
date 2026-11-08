package criterion

import (
	"math"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/measure"
)

// direction is how a number goes into a criterion, so that less is better for every one
// (docs/architecture.md, section 3.6).
type direction int

const (
	plain     direction = iota // less is better as it is
	absolute                   // a number with a sign goes by its modulus
	shortfall                  // a share where more is better goes by its shortfall from one
	more                       // a count where more is better goes with its sign turned
)

// def is a number of a criterion: its direction, the author's tolerance of "no worse", and whether
// its path goes to perfection, zero in its direction, rather than to a ceiling.
type def struct {
	dir       direction
	tolerance float64
	perfect   bool
}

// defs are the numbers the criteria read, by their names in measure. The tolerances are the
// author's of docs/architecture.md, section 3.6; those it does not list take the tolerance of the
// measure nearest them in kind (R25 in docs/decisions.md).
var defs = map[string]def{
	"M1.rms.200":      {plain, 0.01, false},
	"M2.brier.all":    {plain, 0.002, false},
	"M2.brier.after7": {plain, 0.002, false},
	"M2.brier.pages":  {plain, 0.002, false},
	"M3.lag":          {absolute, 0.01, false},
	"M3.settle":       {plain, 1, false},
	"M3.unsettled":    {plain, 0.01, false},
	"M4.false.all":    {plain, 0.02, true},
	"M4.false.break":  {plain, 0.02, true},
	"M4.false.pages":  {plain, 0.02, true},
	"M5.caught":       {plain, 1, false},
	"M5.never":        {plain, 0.01, false}, // a share, as M3's
	"M6.late":         {plain, 1, true},
	"M6.never":        {plain, 0.02, true},
	"M7.error.10":     {plain, 0.01, false},
	"M7.hard":         {plain, 0.01, false},    // a share, as M3's
	"M8.shift":        {absolute, 0.01, false}, // logits, as M1's
	"M9.hasFalse":     {plain, 0.02, true},
	"M9.notFalse":     {plain, 0.02, true},
	"M10.retention":   {plain, 0.002, false}, // Brier, as M2's
	"E1.corridor":     {shortfall, 0.01, false},
	"E2.stuck":        {plain, 0.01, true},
	"E2.nothing":      {plain, 0.01, true},
	"E3.late":         {plain, 0.02, true},
	"E3.retention":    {shortfall, 0.01, true},
	"E3.end":          {shortfall, 0.01, true},
	"E4.given":        {shortfall, 0.01, true},
	"E4.mastered":     {shortfall, 0.01, true},
	"E5.written":      {more, 1, true},
	"E5.again":        {plain, 1, true},
}

// resolution is the bench's resolution of "no worse": 4.3 standard errors of the paired difference
// (docs/architecture.md, section 3.6; R04).
const resolution = 4.3

// directed gives the numbers of names, each turned by its direction.
func directed(names []string) []measure.Number {
	all := slices.Concat(measure.ModelNumbers(), measure.EngineNumbers())
	var out []measure.Number
	for _, name := range names {
		i := slices.IndexFunc(all, func(n measure.Number) bool { return n.Name == name })
		of, d := all[i].Of, defs[name]
		out = append(out, measure.Number{Name: name, Of: func(s measure.Sums) (float64, bool) {
			v, ok := of(s)
			switch d.dir {
			case absolute:
				v = math.Abs(v)
			case shortfall:
				v = 1 - v
			case more:
				v = -v
			}
			return v, ok
		}})
	}
	return out
}

// group is a group of the score: its weight, its populations with theirs — they add up to the
// group's — and its numbers.
type group struct {
	id, name    string
	weight      float64
	populations []weighted
	numbers     []string
}

type weighted struct {
	population string
	weight     float64
}

// The groups of the score of each criterion (docs/architecture.md, section 3.6).
var (
	modelGroups = []group{
		{"still", "stands still", 1, []weighted{{"G0", 1}}, []string{"M1.rms.200", "E1.corridor"}},
		{"learns", "learns", 2, []weighted{{"G2-half", 1.25}, {"G2-fading", 0.5}, {"G2", 0.25}}, []string{"M3.lag", "E1.corridor"}},
		{"jumps", "changes at once", 1, []weighted{{"G3", 0.4}, {"G3-drop", 0.4}, {"F-break", 0.2}}, []string{"M3.settle"}},
		{"off-band", "not in its band", 0.5, []weighted{{"G1", 0.375}, {"S-adult", 0.125}}, []string{"M7.error.10", "M7.hard"}},
		{"mastered", "mastered", 1, []weighted{{"G0", 0.5}, {"G0-skills1", 0.5}}, []string{"M4.false.all", "M6.late"}},
		{"forgets", "forgets", 1, []weighted{{"F-model", 0.4}, {"F-exp", 0.2}, {"F-C1", 0.2}, {"F-growth5", 0.2}}, []string{"M2.brier.after7", "M5.caught", "M4.false.break"}},
		{"pages", "pages", 1, []weighted{{"W-shock", 0.5}, {"W-error", 0.25}, {"W-badday", 0.25}}, []string{"M2.brier.pages", "M4.false.pages"}},
		{"anchor", "anchor", 0.5, []weighted{{"A-fading", 0.25}, {"A-unpractised", 0.25}}, []string{"M9.hasFalse", "M9.notFalse"}},
	}
	engineGroups = []group{
		{"corridor", "corridor", 2, []weighted{{"G2-half", 1.25}, {"G0", 0.5}, {"G2", 0.25}}, []string{"E1.corridor"}},
		{"reviews", "repetitions", 1, []weighted{{"F-model", 0.5}, {"F-exp", 0.25}, {"F-growth5", 0.25}}, []string{"E3.late", "E3.end"}},
		{"coverage", "coverage", 1, []weighted{{"G2-half", 0.75}, {"L-explain", 0.25}}, []string{"E4.mastered"}},
	}
)

// notWorse is a check of "no worse": its numbers, and the populations it holds on — all, or those
// whose names begin with one of the prefixes, or are one of the names.
type notWorse struct {
	numbers []string
	on      []string // empty: every population
}

// applies tells whether a check of "no worse" holds on a population.
func (n notWorse) applies(population string) bool {
	if len(n.on) == 0 {
		return true
	}
	return slices.ContainsFunc(n.on, func(p string) bool {
		return p == population || strings.HasSuffix(p, "-") && strings.HasPrefix(population, p)
	})
}

// The checks of "no worse" of each criterion (docs/architecture.md, section 3.6): on every
// population, and on a measure's own populations (section 3.5).
var (
	modelNotWorse = []notWorse{
		{numbers: []string{"M1.rms.200", "M2.brier.all", "M4.false.all", "M6.never", "E1.corridor"}},
		{[]string{"M3.lag"}, []string{"G2", "G2-half", "G2-fading"}},
		{[]string{"M3.settle", "M3.unsettled"}, []string{"G3", "G3-drop", "F-break"}},
		{[]string{"M5.caught", "M5.never", "E3.late", "E3.retention", "E3.end"}, []string{"F-"}},
		{[]string{"M7.error.10", "M7.hard"}, []string{"G1", "S-adult", "G0-start0.5", "G0-start2"}},
		{[]string{"M8.shift"}, []string{"W-repeat"}},
		{[]string{"M9.hasFalse", "M9.notFalse"}, []string{"A-"}},
		{[]string{"M10.retention"}, []string{"F-none"}},
	}
	engineNotWorse = []notWorse{
		{numbers: []string{"E2.stuck", "E2.nothing", "E1.corridor", "M4.false.all", "E4.given", "E4.mastered"}},
		{[]string{"M3.lag"}, []string{"G2", "G2-half", "G2-fading"}},
		{[]string{"M3.settle", "M3.unsettled"}, []string{"G3", "G3-drop", "F-break"}},
		{[]string{"E3.late", "E3.retention", "E3.end"}, []string{"F-"}},
		{[]string{"E5.written", "E5.again"}, []string{"E-"}},
		{[]string{"M9.hasFalse", "M9.notFalse"}, []string{"A-"}},
	}
)
