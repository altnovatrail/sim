// Package bench runs cells of virtual learners into a file of results, one line a learner, and
// reads two such files in pairs by learner (docs/architecture.md, section 3.6; R24 in
// docs/decisions.md). Two versions of one library do not go into one Go build, so each build of sim
// writes its own file, and the comparison reads them side by side.
package bench

import (
	"bufio"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"

	"runtime/debug"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/invariant"
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
	// Populations, Learners and Repeat are what the run ran, so that the parts of one run merge in its
	// order and parts of another run do not; Part is the part of the learners the file holds, the zero
	// part all of them.
	Populations []string `json:"populations,omitempty"`
	Learners    int      `json:"learners,omitempty"`
	Repeat      bool     `json:"repeat,omitempty"`
	Part        Part     `json:"part,omitzero"`
	// Packages are the packages the run took in place of the pinned ones, and it ran only the graphs
	// whose every module is one of them: the gate of the CI of modules, on the packages of a pull
	// request (docs/architecture.md, section 3.9). None, the pinned packages and every graph.
	Packages []published.Package `json:"packages,omitempty"`
	// Invariants tells that each line records the breaches of I1–I4 on its run: under the engine.
	Invariants bool `json:"invariants"`
}

// Line is one learner of a cell: its population, graph and number, and the sums of its run — of the
// measures of the model and of the engine, and under the engine of the note E6 too (R49) — and the
// oracle's on the same answers.
type Line struct {
	Population string       `json:"population"`
	Graph      []string     `json:"graph"`
	N          int          `json:"n"`
	Sums       measure.Sums `json:"sums"`
	Oracle     measure.Sums `json:"oracle"`
	// Breaches are the breaches of the invariants I1–I4 on the run, under the engine; the ceiling of
	// the engine is no branch of engine, and its runs are not checked (R25 in docs/decisions.md).
	Breaches []string `json:"breaches,omitempty"`
}

// Config is what to run: the populations by name, the learners of each graph, the set and the rule;
// Commits names a library the build does not, by module path, a full SHA. Part takes a part of the
// learners, so that a recipe runs the parts as processes of their own and merges their files; Repeat
// runs each learner twice and holds the two runs to the same bits, I1.
type Config struct {
	Populations []string
	Learners    int
	Set         seed.Set
	Rule        run.Rule
	Commits     map[string]string
	Part        Part
	Repeat      bool
	// Packages take the place of the pinned packages of their modules, and only the graphs whose every
	// module is one of them run.
	Packages []Package
}

// Package is a package a run takes in place of the pinned one: what it is, and its module read.
type Package struct {
	published.Package
	Module *pack.Module
}

// Part is the K-th of N parts of the learners: those whose number n gives K by n mod N. The zero
// part is the whole. The bench draws a learner from its seeds alone, so a learner is the same in any
// part (docs/architecture.md, section 3.4), and the parts run as processes, not goroutines, as R09
// has it (the amendment of 2026-11-08 to R24 in docs/decisions.md).
type Part struct {
	K int `json:"k"`
	N int `json:"n"`
}

// valid tells whether a part is one: the zero part, or k of n with 0 ≤ k < n.
func (p Part) valid() bool {
	return p == Part{} || p.N > 0 && p.K >= 0 && p.K < p.N
}

// takes tells whether a part takes the n-th learner.
func (p Part) takes(n int) bool {
	return p.N <= 1 || n%p.N == p.K
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

// modules gives the modules of a graph by full ID: a package the run takes in place of the pinned
// one, or else the pinned, published or synthetic.
func modules(graph []string, packages []Package) ([]*pack.Module, error) {
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
		var m *pack.Module
		if i := slices.IndexFunc(packages, func(p Package) bool { return p.Package.Module == id }); i >= 0 {
			m = packages[i].Module
		} else if m = pub[id]; m == nil {
			m = syn[id]
		}
		if m == nil {
			return nil, fmt.Errorf("%w: no module %s", ErrBench, id)
		}
		out = append(out, m)
	}
	return out, nil
}

// Run runs the cells of a config, every graph of every population — with packages in place of the
// pinned ones, the graphs they cover —, and writes the file of results to w: the header, then one
// line a learner of its part, in the order of the populations, their graphs and the learners.
func Run(c Config, w io.Writer) error {
	if !c.Part.valid() {
		return fmt.Errorf("%w: no part %d of %d", ErrBench, c.Part.K, c.Part.N)
	}
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
	for k, name := range c.Populations {
		p, ok := population.Named(name)
		if !ok {
			return fmt.Errorf("%w: no population %s", ErrBench, name)
		}
		if slices.Contains(c.Populations[:k], name) {
			return fmt.Errorf("%w: population %s comes twice", ErrBench, name)
		}
		pops = append(pops, p)
	}
	var taken []published.Package
	for k, p := range c.Packages {
		if p.Module == nil || p.Module.ID != p.Package.Module || slices.ContainsFunc(c.Packages[:k], func(q Package) bool { return q.Package.Module == p.Package.Module }) {
			return fmt.Errorf("%w: the package of %s is not its module read, or comes twice", ErrBench, p.Package.Module)
		}
		taken = append(taken, p.Package)
	}
	buf := bufio.NewWriter(w)
	enc := json.NewEncoder(buf)
	h := Header{Set: c.Set, Rule: rule, Libraries: names, Invariants: c.Rule == run.Engine, Populations: c.Populations, Learners: c.Learners, Repeat: c.Repeat, Part: c.Part, Packages: taken}
	if err := enc.Encode(h); err != nil {
		return err
	}
	for _, p := range pops {
		for _, graph := range p.Graphs {
			if !covers(taken, graph) {
				continue
			}
			mods, err := modules(graph, c.Packages)
			if err != nil {
				return err
			}
			for n := range c.Learners {
				if !c.Part.takes(n) {
					continue
				}
				line, err := one(c, p, graph, mods, n)
				if err != nil {
					return err
				}
				if err := enc.Encode(line); err != nil {
					return err
				}
			}
		}
	}
	return buf.Flush()
}

// Rechosen are the first learners of each cell whose every choice is taken again on the input it
// keeps, I1, in every run with invariants: the gate's, repeated, and a decisive one alike. Taking
// each choice again on all the learners would double the engine's work, the most of a run. A
// repetition holds the bits of every learner, but not that the input a choice keeps is the one it was
// made on — a rule that changed its input in place would change it alike in both runs —, and the
// first learners of every cell catch that (R31, R33 in docs/decisions.md).
const Rechosen = 10

// covers tells whether a run on packages runs a graph: with none, every graph; with some, a graph
// whose every module is one of them, and one at least not a package that only came along.
func covers(packages []published.Package, graph []string) bool {
	if len(packages) == 0 {
		return true
	}
	of := func(id string) int {
		return slices.IndexFunc(packages, func(p published.Package) bool { return p.Module == id })
	}
	return !slices.ContainsFunc(graph, func(id string) bool { return of(id) < 0 }) &&
		slices.ContainsFunc(graph, func(id string) bool { return !packages[of(id)].Along })
}

// runOne runs the n-th learner of a population on the modules of a graph.
func runOne(c Config, p population.Population, mods []*pack.Module, n int) (*run.Record, error) {
	l, err := learner.New(c.Set, p.Name, n, mods, p.Assumptions(n))
	if err != nil {
		return nil, err
	}
	r, err := run.New(l, mods, p.Limit)
	if err != nil {
		return nil, err
	}
	return r.Use(c.Rule).All()
}

// checkOf gives the check of the invariants of the n-th learner of a cell: I1 by taking each choice
// again only for the first Rechosen.
func checkOf(n int) func(*run.Record) []invariant.Breach {
	if n < Rechosen {
		return invariant.Check
	}
	return invariant.CheckWithoutI1
}

// one runs the n-th learner of a population on a graph.
func one(c Config, p population.Population, graph []string, mods []*pack.Module, n int) (Line, error) {
	rec, err := runOne(c, p, mods, n)
	if err != nil {
		return Line{}, err
	}
	sums := measure.Model(rec)
	maps.Copy(sums, measure.Engine(rec)) // the keys of the two are apart, but for learners, 1 in both
	line := Line{Population: p.Name, Graph: graph, N: n, Sums: sums, Oracle: measure.Oracle(rec)}
	if c.Rule == run.Engine {
		// The note of consolidation in a lesson is read under the engine alone, as the breaches
		// are: the ceiling of the engine does not consolidate, and its file stays as it was (R49).
		maps.Copy(sums, measure.Note(rec))
		for _, b := range checkOf(n)(rec) {
			line.Breaches = append(line.Breaches, b.String())
		}
	}
	if c.Repeat {
		again, err := runOne(c, p, mods, n)
		if err != nil {
			return Line{}, err
		}
		for _, b := range invariant.Same(rec, again) {
			line.Breaches = append(line.Breaches, b.String())
		}
	}
	return line, nil
}

// Merge merges the files of the parts of one run into the file of the whole: the same set, rule,
// libraries, invariants, populations, learners and repetition, every part of N once, and the lines
// in the order of the run's populations, their graphs and the learners — the order of a run of the
// whole.
func Merge(parts []*File) (*File, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("%w: no part to merge", ErrBench)
	}
	first := parts[0].Header
	n := max(1, first.Part.N)
	seen := make([]bool, n)
	out := &File{Header: first}
	out.Header.Part = Part{}
	for _, f := range parts {
		h := f.Header
		if h.Set != first.Set || h.Rule != first.Rule || h.Invariants != first.Invariants || !maps.Equal(h.Libraries, first.Libraries) ||
			!slices.Equal(h.Populations, first.Populations) || h.Learners != first.Learners || h.Repeat != first.Repeat || !slices.Equal(h.Packages, first.Packages) || max(1, h.Part.N) != n {
			return nil, fmt.Errorf("%w: part %d of %d is not of the run of part %d of %d", ErrBench, h.Part.K, h.Part.N, first.Part.K, first.Part.N)
		}
		if !h.Part.valid() {
			return nil, fmt.Errorf("%w: no part %d of %d", ErrBench, h.Part.K, h.Part.N)
		}
		if seen[h.Part.K] {
			return nil, fmt.Errorf("%w: part %d of %d comes twice", ErrBench, h.Part.K, n)
		}
		seen[h.Part.K] = true
		out.Lines = append(out.Lines, f.Lines...)
	}
	if i := slices.Index(seen, false); i >= 0 {
		return nil, fmt.Errorf("%w: part %d of %d is missing", ErrBench, i, n)
	}
	order := map[string]int{}
	for k, cell := range CellsOn(first.Populations, first.Packages) {
		order[cell] = k
	}
	slices.SortStableFunc(out.Lines, func(a, b Line) int {
		return cmp.Or(cmp.Compare(order[cellOf(a)], order[cellOf(b)]), cmp.Compare(a.N, b.N))
	})
	return out, nil
}

// Write writes a file of results: its header, then a line a learner.
func Write(f *File, w io.Writer) error {
	buf := bufio.NewWriter(w)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(f.Header); err != nil {
		return err
	}
	for _, l := range f.Lines {
		if err := enc.Encode(l); err != nil {
			return err
		}
	}
	return buf.Flush()
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

// CellsOf names the cells of populations, by name: each on every graph it runs on, in order; a name
// that is no population has none.
func CellsOf(populations []string) []string {
	return CellsOn(populations, nil)
}

// CellsOn names the cells of populations that a run on packages runs: those on the graphs the
// packages cover; with none, every graph.
func CellsOn(populations []string, packages []published.Package) []string {
	var out []string
	for _, name := range populations {
		p, _ := population.Named(name)
		for _, graph := range p.Graphs {
			if !covers(packages, graph) {
				continue
			}
			out = append(out, p.Name+"/"+strings.Join(graph, "+"))
		}
	}
	return out
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
	cells, err := Cells(a, b)
	if err != nil {
		return err
	}
	numbers := numbersOf(a, b)
	if _, err := fmt.Fprintln(w, "cell\tnumber\ta\tb\ta-b\tlow\thigh"); err != nil {
		return err
	}
	for _, c := range cells {
		pairs, err := measure.ReadPair(c.Cell, c.Graph, c.A, c.B, numbers)
		if err != nil {
			return err
		}
		for _, p := range pairs {
			name := c.Cell
			if p.Diff.NoGraph {
				name += " (no graph)"
			}
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", name, p.Diff.Name, show(p.A, p.A.Value), show(p.B, p.B.Value), show(p.Diff, p.Diff.Value), show(p.Diff, p.Diff.Low), show(p.Diff, p.Diff.High)); err != nil {
				return err
			}
		}
	}
	return nil
}

// Cell is a cell of two sides, its learners in pairs: the n-th of A and of B the same learner, whose
// number in the cell is the n-th of N.
type Cell struct {
	Cell, Population string
	Graph            []string
	A, B             []measure.Sums
	N                []int
}

// Cells reads two sides in pairs by learner: the same set, the same cells and the same learners in
// the same order, or none.
func Cells(a, b Side) ([]Cell, error) {
	if a.File.Header.Set != b.File.Header.Set {
		return nil, fmt.Errorf("%w: the sides ran on different sets of learners, %v and %v", ErrBench, a.File.Header.Set, b.File.Header.Set)
	}
	if len(a.File.Lines) != len(b.File.Lines) {
		return nil, fmt.Errorf("%w: %d learners on one side and %d on the other", ErrBench, len(a.File.Lines), len(b.File.Lines))
	}
	sums := func(s Side, l Line) measure.Sums {
		if s.Oracle {
			return l.Oracle
		}
		return l.Sums
	}
	var out []Cell
	for k, la := range a.File.Lines {
		lb := b.File.Lines[k]
		cell := cellOf(la)
		if cellOf(lb) != cell || la.N != lb.N {
			return nil, fmt.Errorf("%w: learner %s/%d against %s/%d: the sides are not the same learners", ErrBench, cell, la.N, cellOf(lb), lb.N)
		}
		if len(out) == 0 || out[len(out)-1].Cell != cell {
			out = append(out, Cell{Cell: cell, Population: la.Population, Graph: la.Graph})
		}
		c := &out[len(out)-1]
		c.A, c.B, c.N = append(c.A, sums(a, la)), append(c.B, sums(b, lb)), append(c.N, la.N)
	}
	return out, nil
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
