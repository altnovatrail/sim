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

// Decisive is the learners of every cell a decisive run takes; a run of fewer is a rough look and
// decides nothing (docs/architecture.md, section 3.8).
const Decisive = 4000

// ErrCriterion is a judgement the criterion refuses.
var ErrCriterion = errors.New("criterion")

// Candidate is a candidate rule: its file of results on the learners of the base, the group of the
// score it names its own, its goal, if it names one of the model's, and what the choice among equals
// reads — the constants it changes, whether it adds a field to the model's state or the engine's
// input, and its distance from the base, in the author's order. In a pair it names the library whose
// goal it is.
type Candidate struct {
	Name      string
	File      *bench.File
	Group     string
	Goal      string
	Constants int
	NewField  bool
	Distance  int
	Library   string
}

// Confirm are the runs of the confirmation on the held-out learners of the choice: the base and the
// chosen candidate.
type Confirm struct {
	Base, Candidate *bench.File
}

// Input is a judgement: the criterion, the base's run, the candidates on the same learners, and the
// confirmation when it has run.
type Input struct {
	Criterion  string
	Base       *bench.File
	Candidates []Candidate
	Confirm    *Confirm
}

// Check is one check of a candidate: a constraint on a cell and a number, or on its run.
type Check struct {
	Kind    string // invariants, floor, no worse, goal
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
	Left             []string
}

// Verdict is how a candidate stands: its checks, whether it holds them all, its score, its group's,
// and in a pair the other library's score; and whether it is better than the base.
type Verdict struct {
	Candidate string
	Checks    []Check
	Holds     bool
	Score     Score
	Group     *Score // the score of the group it names, if it names one
	Other     *Score
	Better    bool
}

// Report is what a judgement gives.
type Report struct {
	Criterion string
	Decisive  bool
	Learners  int      // the fewest learners of a cell over the runs
	Graphs    []string // the graphs the choice is made on: those of published packages
	NoGraph   []string // the cells on a graph with no published package
	Verdicts  []Verdict
	Chosen    string   // the candidate chosen, if any
	Equals    []string // the candidates equal to the best, the chosen among them
	Confirmed *Verdict // the chosen on its held-out learners, when the confirmation has run
	Outcome   string
}

// Judge judges the candidates of an input.
func Judge(in Input) (*Report, error) {
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
		if !strings.HasPrefix(f.Base.Header.Set.Experiment, "held-out-") || f.Base.Header.Set == in.Base.Header.Set {
			return fmt.Errorf("%w: the confirmation ran on %s, not on held-out learners apart from the choice's", ErrCriterion, f.Base.Header.Set.Experiment)
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
