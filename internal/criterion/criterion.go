// Package criterion judges candidates against the base by the criterion of the model, of the
// engine, or of a pair of both: the constraints, the score as the share of the path from the base to
// the ceiling, "better", the choice, the confirmation on held-out learners and the way out
// (docs/architecture.md, section 3.6; R25 in docs/decisions.md). It reads files of results that
// sim run wrote, the candidates' beside the base's on the same learners, and holds the paths fixed
// on the first decisive run of the base (R28).
package criterion

import (
	"cmp"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
)

// The criteria.
const (
	Model  = "model"
	Engine = "engine"
	Pair   = "pair"
)

// Decisive is the learners of every cell a decisive run takes, and the recipe just decisive runs; a
// run of fewer is a rough look and decides nothing (docs/architecture.md, section 3.8). It was 4,000
// until a decisive run took hours on olympiad maths 0.1.1, and a run is to take an hour at most (R33
// in docs/decisions.md); then 1,000, until engine v0.1.6 made the choice faster and 1,500 fit in 45
// minutes (the amendment of 2026-11-11 to R33); then 1,500, until learner-model v0.1.6 made its
// readout faster and 2,000 fit, by a sample on a free machine, in about 38 minutes (the amendment of
// 2026-11-12 to R33).
const Decisive = 2000

// ErrCriterion is a judgement the criterion refuses.
var ErrCriterion = errors.New("criterion")

// Candidate is a candidate rule: its file of results on the learners of the base, the group of the
// score it names its own, its goal, if it names one of the model's, and what the choice among equals
// reads — the constants it changes, whether it adds a field to the model's state or the engine's
// input, and its distance from the base, in the author's order. In a pair it names the library whose
// goal it is. EqualOn are the graphs, by the IDs of their modules joined by "+", on which it claims to
// give the same runs as the base: there every learner's sums are the base's, or it fails the check
// "neutral" (R49 in docs/decisions.md).
type Candidate struct {
	Name      string
	File      *bench.File
	Group     string
	Goal      string
	Constants int
	NewField  bool
	Distance  int
	Library   string
	EqualOn   []string
}

// Confirm are the runs of the confirmation on the held-out learners of the choice: the base and the
// chosen candidate.
type Confirm struct {
	Base, Candidate *bench.File
}

// Input is a judgement: the criterion, the base's run, the candidates on the same learners, the
// confirmation when it has run, and the libraries the base is fixed at.
type Input struct {
	Criterion  string
	Base       *bench.File
	Candidates []Candidate
	Confirm    *Confirm
	// Pins are the libraries of the gate file, by module path: the base runs on them, and a
	// candidate of one library on the other of them (R40 in docs/decisions.md).
	Pins map[string]Pin
}

// Pin is the point a library of the base is fixed at: the full SHA of its commit, and the tag of
// the release on it, empty while it is in none.
type Pin struct {
	Commit, Version string
}

// names tells whether a file of results names the library at its pin: by the tag, or, while there
// is none, by the pseudo-version of the commit, as go.mod names it then (docs/architecture.md,
// section 3.9).
func (p Pin) names(got string) bool {
	if p.Version != "" {
		return got == p.Version
	}
	return len(p.Commit) >= 12 && !strings.Contains(got, "@") && strings.HasSuffix(got, "-"+p.Commit[:12])
}

// String names the point as the gate names it.
func (p Pin) String() string {
	switch {
	case p.Commit == "":
		return "no commit"
	case p.Version != "":
		return p.Version + " (" + p.Commit + ")"
	}
	return "the pseudo-version of commit " + p.Commit
}

// Check is one check of a candidate: a constraint on a cell and a number, or on its run.
type Check struct {
	Kind    string // invariants, floor, no worse, goal, neutral
	Cell    string
	Number  string
	Met     bool
	Skipped bool // the base reads no such number on the cell: there is nothing to hold
	NoGraph bool // a graph with no published package: neither met nor failed
	Detail  string
}

// Score is a score: its value, its 95 % interval, its values over the resamples, and the numbers
// that left it — a number with no fixed path: none recorded, none read, or one of zero or less.
type Score struct {
	Value, Low, High float64
	Read             bool
	Resampled        []float64
	Left             []string // left the score: the path not read, or not above zero
	Small            []string // left the score: the path above zero but within the author's tolerance
}

// Verdict is how a candidate stands: its checks, whether it holds them all, its score, its group's,
// and in a pair the other library's score; whether it is better than the base; and the note of
// consolidation in a lesson, which decides nothing.
type Verdict struct {
	Candidate string
	Checks    []Check
	Holds     bool
	Score     Score
	Group     *Score // the score of the group it names, if it names one
	Other     *Score
	Better    bool
	Note      []Note
}

// Note is a number of the note of consolidation in a lesson, E6, on a cell where the base or the
// candidate consolidated: the base's value, the candidate's and their difference, base − candidate,
// with the interval from the same resamples (docs/architecture.md, section 3.6; R49). No check,
// path or weight reads it.
type Note struct {
	Cell string
	measure.Pair
}

// Report is what a judgement gives.
type Report struct {
	Criterion string
	Decisive  bool
	Learners  int      // the fewest learners of a cell over the runs
	Graphs    []string // the graphs the choice is made on: those of published packages
	NoGraph   []string // the cells on a graph with no published package
	Verdicts  []Verdict
	Chosen    string   // the candidate chosen, if any; on a rough look, the one a decisive run would choose
	Equals    []string // the candidates equal to the best, the chosen among them
	Confirmed *Verdict // the chosen on its held-out learners, when the confirmation has run
	Outcome   string
}

// Judge judges the candidates of an input. A rough look — fewer learners a cell than a decisive run
// takes — decides nothing, and its outcome says so first (docs/architecture.md, section 3.6).
func Judge(in Input) (*Report, error) {
	r, err := decide(in)
	if err != nil {
		return nil, err
	}
	if !r.Decisive {
		o := []byte(r.Outcome) // the outcomes are the criterion's own, each opening with a capital letter of ASCII
		if len(o) > 0 && 'A' <= o[0] && o[0] <= 'Z' {
			o[0] += 'a' - 'A'
		}
		r.Outcome = "A rough look, which decides nothing: " + string(o)
	}
	return r, nil
}

// decide judges the candidates of an input as a decisive run would.
func decide(in Input) (*Report, error) {
	if err := in.check(); err != nil {
		return nil, err
	}
	paths, err := fixed()
	if err != nil {
		return nil, err
	}
	if len(paths.Paths) == 0 {
		return nil, fmt.Errorf("%w: no path is fixed yet: the first decisive run of the base fixes them (sim paths)", ErrCriterion)
	}
	if err := in.libraries(paths); err != nil {
		return nil, err
	}
	runs := slices.Concat([]*bench.File{in.Base}, files(in.Candidates))
	if in.Confirm != nil {
		runs = append(runs, in.Confirm.Base, in.Confirm.Candidate)
	}
	r := &Report{Criterion: in.Criterion, Learners: fewest(runs)}
	r.Decisive = r.Learners >= Decisive
	var scores [][]float64
	var eligible []int
	for k, c := range in.Candidates {
		v, err := judge(in.Criterion, in.Base, paths, c)
		if err != nil {
			return nil, err
		}
		r.Verdicts = append(r.Verdicts, v)
		scores = append(scores, v.Score.Resampled)
		if v.Holds && v.Better {
			eligible = append(eligible, k)
		}
	}
	cells, err := bench.Cells(bench.Side{File: in.Base}, bench.Side{File: in.Base})
	if err != nil {
		return nil, err
	}
	for _, c := range cells {
		if noGraph(c.Graph) {
			r.NoGraph = append(r.NoGraph, c.Cell)
		} else if g := strings.Join(c.Graph, "+"); !slices.Contains(r.Graphs, g) {
			r.Graphs = append(r.Graphs, g)
		}
	}
	slices.Sort(r.Graphs)
	if len(eligible) == 0 {
		r.Outcome = "No candidate holds every constraint and is better than the base: the base stays as it is"
		return r, nil
	}
	best := slices.MaxFunc(eligible, func(a, b int) int { return cmp.Compare(r.Verdicts[a].Score.Value, r.Verdicts[b].Score.Value) })
	var equals []int
	for _, k := range eligible {
		if k == best {
			equals = append(equals, k)
			continue
		}
		diff := make([]float64, len(scores[best]))
		for i := range diff {
			diff[i] = scores[best][i] - scores[k][i]
		}
		if low, high := measure.Ends(diff); low <= 0 && high >= 0 {
			equals = append(equals, k)
		}
	}
	chosen := slices.MinFunc(equals, func(a, b int) int {
		x, y := in.Candidates[a], in.Candidates[b]
		return cmp.Or(cmp.Compare(x.Constants, y.Constants), cmp.Compare(b2i(x.NewField), b2i(y.NewField)), cmp.Compare(x.Distance, y.Distance), strings.Compare(x.Name, y.Name))
	})
	for _, k := range equals {
		r.Equals = append(r.Equals, in.Candidates[k].Name)
	}
	r.Chosen = in.Candidates[chosen].Name
	if in.Confirm == nil {
		r.Outcome = "Chosen: " + r.Chosen + "; it awaits its confirmation on the held-out learners of the choice"
		return r, nil
	}
	c := in.Candidates[chosen]
	c.File = in.Confirm.Candidate
	v, err := judge(in.Criterion, in.Confirm.Base, paths, c)
	if err != nil {
		return nil, err
	}
	r.Confirmed = &v
	if v.Holds && v.Better {
		r.Outcome = "Chosen and confirmed: " + r.Chosen
	} else {
		r.Outcome = "Chosen, but not confirmed on the held-out learners: the base stays as it is; no second candidate is tried on them"
	}
	return r, nil
}

// check refuses an input whose runs do not go together.
func (in Input) check() error {
	switch in.Criterion {
	case Model, Engine, Pair:
	default:
		return fmt.Errorf("%w: no criterion %q", ErrCriterion, in.Criterion)
	}
	if in.Base == nil || in.Base.Header.Rule != "engine" {
		return fmt.Errorf("%w: a judgement needs the base's run under the engine's rule", ErrCriterion)
	}
	for _, c := range in.Candidates {
		if c.File == nil || c.File.Header.Rule != "engine" {
			return fmt.Errorf("%w: candidate %s has no run under the engine's rule", ErrCriterion, c.Name)
		}
		if c.Group != "" && !slices.ContainsFunc(groupsOf(in.Criterion, c.Library), func(g group) bool { return g.id == c.Group }) {
			return fmt.Errorf("%w: candidate %s names group %q, which its criterion does not have", ErrCriterion, c.Name, c.Group)
		}
		if c.Goal != "" && !slices.Contains([]string{"lag", "corridor", "settle"}, c.Goal) {
			return fmt.Errorf("%w: candidate %s names goal %q: lag, corridor or settle", ErrCriterion, c.Name, c.Goal)
		}
		if c.Goal != "" && (in.Criterion == Engine || in.Criterion == Pair && c.Library != Model) {
			return fmt.Errorf("%w: candidate %s names a goal of the model, and the engine has no goal number", ErrCriterion, c.Name)
		}
		if in.Criterion == Pair && c.Library != Model && c.Library != Engine {
			return fmt.Errorf("%w: candidate %s of a pair names library %q: model or engine", ErrCriterion, c.Name, c.Library)
		}
	}
	if f := in.Confirm; f != nil {
		if f.Base == nil || f.Candidate == nil {
			return fmt.Errorf("%w: a confirmation needs the base's run and the candidate's", ErrCriterion)
		}
		if f.Base.Header.Rule != "engine" || f.Candidate.Header.Rule != "engine" {
			return fmt.Errorf("%w: the confirmation's base and candidate ran under %q and %q", ErrCriterion, f.Base.Header.Rule, f.Candidate.Header.Rule)
		}
		// The held-out sets of a criterion are its own (docs/architecture.md, section 3.8); which of
		// them a round takes, its judgement file's folder does not say, and README names it.
		if !strings.HasPrefix(f.Base.Header.Set.Experiment, "held-out-"+in.Criterion+"-") || f.Base.Header.Set == in.Base.Header.Set {
			return fmt.Errorf("%w: the confirmation ran on %s, not on held-out learners of the criterion %s apart from the choice's", ErrCriterion, f.Base.Header.Set.Experiment, in.Criterion)
		}
	}
	return nil
}

// libraries refuses runs that are not on the libraries of the base: the base and the confirmation's
// base on those the paths were read on and the gate file fixes, and a candidate of one library, the
// chosen one on its held-out learners too, on the other library of the base — by its version, by
// the commit of the pin, as the recipe of the decisive run builds a candidate, or by the
// pseudo-version of that commit (R34, R40 in docs/decisions.md). A pair changes both libraries, and
// has no other to hold.
func (in Input) libraries(paths *Paths) error {
	if len(in.Pins) == 0 {
		return fmt.Errorf("%w: a judgement needs the libraries the base is fixed at, those of the gate file", ErrCriterion)
	}
	bases, runs := []*bench.File{in.Base}, slices.Clone(in.Candidates)
	if in.Confirm != nil {
		bases = append(bases, in.Confirm.Base)
		runs = append(runs, Candidate{Name: "chosen, on its held-out learners,", File: in.Confirm.Candidate})
	}
	for _, b := range bases {
		if !maps.Equal(b.Header.Libraries, paths.Libraries) {
			return fmt.Errorf("%w: the base ran on %v, and the paths were read on %v", ErrCriterion, b.Header.Libraries, paths.Libraries)
		}
		for _, path := range bench.Libraries {
			if pin, got := in.Pins[path], b.Header.Libraries[path]; !pin.names(got) {
				return fmt.Errorf("%w: the gate file names %s at %v, and the base ran on %q", ErrCriterion, path, pin, got)
			}
		}
	}
	other := map[string]string{Model: bench.Libraries[1], Engine: bench.Libraries[0]}[in.Criterion]
	if other == "" {
		return nil
	}
	want, pin := in.Base.Header.Libraries[other], in.Pins[other]
	for _, c := range runs {
		if got := c.File.Header.Libraries[other]; got == "" || got != want && got != pin.Commit && !pin.names(got) {
			return fmt.Errorf("%w: candidate %s ran on %s at %q, and the base on %s, %v: a candidate of one library is judged on the other library of the base", ErrCriterion, c.Name, other, got, want, pin)
		}
	}
	return nil
}

// groupsOf gives the groups of a criterion; a pair's are those of the library its candidate names.
func groupsOf(criterion, library string) []group {
	if criterion == Engine || criterion == Pair && library == Engine {
		return engineGroups
	}
	return modelGroups
}

// files gives the candidates' files.
func files(cs []Candidate) []*bench.File {
	var out []*bench.File
	for _, c := range cs {
		out = append(out, c.File)
	}
	return out
}

// fewest gives the fewest learners of a cell over files.
func fewest(fs []*bench.File) int {
	least := -1
	for _, f := range fs {
		n := map[string]int{}
		for _, l := range f.Lines {
			n[l.Population+"/"+strings.Join(l.Graph, "+")]++
		}
		for _, cell := range slices.Sorted(maps.Keys(n)) {
			if least < 0 || n[cell] < least {
				least = n[cell]
			}
		}
	}
	return max(least, 0)
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
