// Package bench runs cells of virtual learners into a file of results, one line a learner, and
// reads two such files in pairs by learner (docs/architecture.md, section 3.6; R24 in
// docs/decisions.md). Two versions of one library do not go into one Go build, so each build of sim
// writes its own file, and the comparison reads them side by side.
package bench

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"sync"

	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// The libraries a file of results names, by module path.
var Libraries = []string{"github.com/altnovatrail/learner-model", "github.com/altnovatrail/engine"}

// ErrBench is a run or a comparison the bench refuses.
var ErrBench = errors.New("bench")

// Header is the first line of a file of results: the set of learners, the rule that chose, and the
// version or commit of each library the build ran on.
type Header struct {
	Set       seed.Set          `json:"set"`
	Rule      string            `json:"rule"`
	Libraries map[string]string `json:"libraries"`
}

// Line is one learner of a cell: its population, graph and number, and the sums of its run — of the
// measures of the model and of the engine — and the oracle's on the same answers.
type Line struct {
	Population string       `json:"population"`
	Graph      []string     `json:"graph"`
	N          int          `json:"n"`
	Sums       measure.Sums `json:"sums"`
	Oracle     measure.Sums `json:"oracle"`
}

// Config is what to run: the populations by name, the learners of each graph, the set and the rule;
// Commits names a library the build does not, by module path, a full SHA.
type Config struct {
	Populations []string
	Learners    int
	Set         seed.Set
	Rule        run.Rule
	Commits     map[string]string
}

// rules names the rules in a file of results.
var rules = map[run.Rule]string{run.Engine: "engine", run.Ceiling: "ceiling"}

// Names gives the version or commit of each library: the version the build took it at — a tag, or a
// pseudo-version, which holds its commit —, its replacement's, with the replacement's path when it
// is another module, or the commit the caller names. A library the build took from a workspace has
// no version, and the caller has to name its commit, a full SHA; a commit named for a library the
// build names a version of, or for a module that is no library, is refused
// (docs/architecture.md, section 3.9).
func Names(info *debug.BuildInfo, commits map[string]string) (map[string]string, error) {
	for _, path := range slices.Sorted(maps.Keys(commits)) {
		c := commits[path]
		if !slices.Contains(Libraries, path) {
			return nil, fmt.Errorf("%w: a commit named for %s, which is not one of %v", ErrBench, path, Libraries)
		}
		if len(c) != 40 || strings.Trim(c, "0123456789abcdef") != "" {
			return nil, fmt.Errorf("%w: the commit of %s is %q, not a full SHA", ErrBench, path, c)
		}
	}
	out := map[string]string{}
	for _, lib := range Libraries {
		var dep *debug.Module
		if info != nil {
			if i := slices.IndexFunc(info.Deps, func(d *debug.Module) bool { return d.Path == lib }); i >= 0 {
				dep = info.Deps[i]
			}
		}
		name := ""
		if dep != nil && dep.Replace != nil {
			name, dep = dep.Replace.Path+"@", dep.Replace
		}
		c, named := commits[lib]
		switch known := dep != nil && dep.Version != "" && dep.Version != "(devel)"; {
		case known && named:
			return nil, fmt.Errorf("%w: the build took %s at %s, and a commit is named for it too", ErrBench, lib, dep.Version)
		case known:
			if name == lib+"@" {
				name = ""
			}
			out[lib] = name + dep.Version
		case named:
			out[lib] = c
		default:
			return nil, fmt.Errorf("%w: the build names no version of %s: name its commit", ErrBench, lib)
		}
	}
	return out, nil
}

// modules gives the modules of a graph by full ID, published or synthetic.
func modules(graph []string) ([]*pack.Module, error) {
	pub, err := published.Modules()
	if err != nil {
		return nil, err
	}
	syn, err := synthetic.Modules()
	if err != nil {
		return nil, err
	}
	var out []*pack.Module
	for _, id := range graph {
		m := pub[id]
		if m == nil {
			m = syn[id]
		}
		if m == nil {
			return nil, fmt.Errorf("%w: no module %s", ErrBench, id)
		}
		out = append(out, m)
	}
	return out, nil
}

// Run runs the cells of a config, every graph of every population, and writes the file of results
// to w: the header, then one line a learner, in the order of the populations, their graphs and the
// learners. The learners run in parallel, and the lines come in that order all the same.
func Run(c Config, w io.Writer) error {
	info, _ := debug.ReadBuildInfo()
	names, err := Names(info, c.Commits)
	if err != nil {
		return err
	}
	rule, ok := rules[c.Rule]
	if !ok {
		return fmt.Errorf("%w: no rule %d", ErrBench, c.Rule)
	}
	var pops []population.Population
	for _, name := range c.Populations {
		p, ok := population.Named(name)
		if !ok {
			return fmt.Errorf("%w: no population %s", ErrBench, name)
		}
		pops = append(pops, p)
	}
	buf := bufio.NewWriter(w)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(Header{Set: c.Set, Rule: rule, Libraries: names}); err != nil {
		return err
	}
	for _, p := range pops {
		for _, graph := range p.Graphs {
			mods, err := modules(graph)
			if err != nil {
				return err
			}
			lines := make([]Line, c.Learners)
			errs := make([]error, c.Learners)
			var wg sync.WaitGroup
			next := make(chan int)
			for range max(1, runtime.GOMAXPROCS(0)) {
				wg.Go(func() {
					for n := range next {
						lines[n], errs[n] = one(c, p, graph, mods, n)
					}
				})
			}
			for n := range c.Learners {
				next <- n
			}
			close(next)
			wg.Wait()
			for n := range c.Learners {
				if errs[n] != nil {
					return errs[n]
				}
				if err := enc.Encode(lines[n]); err != nil {
					return err
				}
			}
		}
	}
	return buf.Flush()
}

// one runs the n-th learner of a population on a graph.
func one(c Config, p population.Population, graph []string, mods []*pack.Module, n int) (Line, error) {
	l, err := learner.New(c.Set, p.Name, n, mods, p.Assumptions(n))
	if err != nil {
		return Line{}, err
	}
	r, err := run.New(l, mods, p.Limit)
	if err != nil {
		return Line{}, err
	}
	rec, err := r.Use(c.Rule).All()
	if err != nil {
		return Line{}, err
	}
	sums := measure.Model(rec)
	maps.Copy(sums, measure.Engine(rec)) // the keys of the two are apart, but for learners, 1 in both
	return Line{Population: p.Name, Graph: graph, N: n, Sums: sums, Oracle: measure.Oracle(rec)}, nil
}

// File is a file of results read back.
type File struct {
	Header Header
	Lines  []Line
}

// Read reads a file of results.
func Read(r io.Reader) (*File, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(nil, 1<<26)
	f := &File{}
	for k := 0; sc.Scan(); k++ {
		var err error
		if k == 0 {
			err = json.Unmarshal(sc.Bytes(), &f.Header)
		} else {
			var l Line
			err = json.Unmarshal(sc.Bytes(), &l)
			f.Lines = append(f.Lines, l)
		}
		if err != nil {
			return nil, fmt.Errorf("%w: line %d: %w", ErrBench, k+1, err)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if f.Header.Rule == "" {
		return nil, fmt.Errorf("%w: a file of results with no header", ErrBench)
	}
	return f, nil
}

// Side is one side of a comparison: a file of results, and whether it is read as its rule's sums or
// as the oracle's.
type Side struct {
	File   *File
	Oracle bool
}

// cellOf names the cell of a line: its population and graph.
func cellOf(l Line) string {
	return l.Population + "/" + strings.Join(l.Graph, "+")
}

// Compare reads two sides in pairs by learner — the same set, the same cells and the same learners,
// or none — and writes, for each cell and each number of the measures, a's value, b's, and the
// difference a − b with its 95 % interval from the same resamples in both, tab-separated; a number
// a cell cannot read is "-", and a cell on a graph with no published package says so.
func Compare(a, b Side, w io.Writer) error {
	if a.File.Header.Set != b.File.Header.Set {
		return fmt.Errorf("%w: the sides ran on different sets of learners, %v and %v", ErrBench, a.File.Header.Set, b.File.Header.Set)
	}
	if len(a.File.Lines) != len(b.File.Lines) {
		return fmt.Errorf("%w: %d learners on one side and %d on the other", ErrBench, len(a.File.Lines), len(b.File.Lines))
	}
	sums := func(s Side, l Line) measure.Sums {
		if s.Oracle {
			return l.Oracle
		}
		return l.Sums
	}
	numbers := numbersOf(a, b)
	if _, err := fmt.Fprintln(w, "cell\tnumber\ta\tb\ta-b\tlow\thigh"); err != nil {
		return err
	}
	for start := 0; start < len(a.File.Lines); {
		first := a.File.Lines[start]
		cell := cellOf(first)
		var as, bs []measure.Sums
		end := start
		for ; end < len(a.File.Lines) && cellOf(a.File.Lines[end]) == cell; end++ {
			la, lb := a.File.Lines[end], b.File.Lines[end]
			if cellOf(lb) != cell || la.N != lb.N {
				return fmt.Errorf("%w: learner %s/%d against %s/%d: the sides are not the same learners", ErrBench, cell, la.N, cellOf(lb), lb.N)
			}
			as, bs = append(as, sums(a, la)), append(bs, sums(b, lb))
		}
		pairs, err := measure.ReadPair(cell, first.Graph, as, bs, numbers)
		if err != nil {
			return err
		}
		for _, p := range pairs {
			name := cell
			if p.Diff.NoGraph {
				name += " (no graph)"
			}
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", name, p.Diff.Name, show(p.A, p.A.Value), show(p.B, p.B.Value), show(p.Diff, p.Diff.Value), show(p.Diff, p.Diff.Low), show(p.Diff, p.Diff.High)); err != nil {
				return err
			}
		}
		start = end
	}
	return nil
}

// numbersOf gives the numbers two sides compare: the measures of the model, and of the engine unless
// a side is the oracle, which reads the model's alone; E3 is not read off the ceiling of the engine,
// whose path goes to perfection (docs/architecture.md, section 3.6; R24).
func numbersOf(a, b Side) []measure.Number {
	if a.Oracle || b.Oracle {
		return measure.ModelNumbers()
	}
	out := measure.ModelNumbers()
	ceiling := a.File.Header.Rule == rules[run.Ceiling] || b.File.Header.Rule == rules[run.Ceiling]
	for _, n := range measure.EngineNumbers() {
		if !ceiling || !strings.HasPrefix(n.Name, "E3.") {
			out = append(out, n)
		}
	}
	return out
}

// show writes a number of a value, or "-" when the value is not read.
func show(v measure.Value, x float64) string {
	if !v.Read {
		return "-"
	}
	return fmt.Sprintf("%.6g", x)
}
