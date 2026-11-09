// Package gate holds the gate of the bench: the gate file, which names what the gate's numbers were
// recorded on — the commits of the libraries, the pinned packages, the versions the workspace chose
// and the go lines — and the bands; the gate's cells and numbers; the checks of a run against the
// file; and the recording of the bands anew (docs/architecture.md, sections 3.7 and 3.9; R26 in
// docs/decisions.md).
package gate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// Library is a library the gate builds the bench on: its repository, the full SHA of its commit,
// and the tag of the release on that commit, empty while it is in none.
type Library struct {
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
	Version    string `json:"version"`
}

// Modules are the packages the gate runs on: the release of modules and the digest of its index, and
// each package by its module, version and package digest.
type Modules struct {
	Release  string              `json:"release"`
	Index    string              `json:"index"`
	Packages []published.Package `json:"packages"`
}

// Band is the band of one number of the gate on one cell: its recorded value and the half-width of
// its 95 % interval; a number the recording could not read is recorded as not read.
type Band struct {
	Cell   string  `json:"cell"`
	Number string  `json:"number"`
	Read   bool    `json:"read"`
	Value  float64 `json:"value"`
	Half   float64 `json:"half"`
}

// Bands are the bands of the gate and why they were recorded.
type Bands struct {
	Reason  string `json:"reason"`
	Numbers []Band `json:"numbers"`
}

// File is the gate file: the libraries by module path, the packages, the versions of the shared
// dependencies the workspace chose, the go line of each module by its path, and the bands.
type File struct {
	Libraries    map[string]Library `json:"libraries"`
	Modules      Modules            `json:"modules"`
	Dependencies map[string]string  `json:"dependencies"`
	Go           map[string]string  `json:"go"`
	Bands        Bands              `json:"bands"`
}

// Read reads a gate file; a field it does not know is a mistake.
func Read(r io.Reader) (*File, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var f File
	if err := dec.Decode(&f); err != nil {
		return nil, fmt.Errorf("gate: the gate file: %w", err)
	}
	return &f, nil
}

// Write writes a gate file, indented, its keys in order; a number in the shortest form that reads
// back to the same bits.
func (f *File) Write(w io.Writer) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(f); err != nil {
		return err
	}
	_, err := w.Write(buf.Bytes())
	return err
}

// Learners are the learners of each cell of the gate: the first of the working set. They were 300,
// as standalone's guard had; on olympiad maths 0.1.1 a choice of the engine takes some ten times
// longer, and 300 would not fit the job's budget of fifteen minutes, which the gate keeps by a
// change of sim, not of the time (docs/architecture.md, section 3.9; R31 in docs/decisions.md).
const Learners = 100

// learners are the learners whole holds each cell to: Learners, but a test of the checks lowers it,
// since reading the numbers off 100 learners a cell under the race detector takes minutes.
var learners = Learners

// Populations are the populations of the gate, each on every graph it runs on, in the order of
// population (docs/architecture.md, section 3.7). A-linked runs on the two published maths together,
// since olympiad maths 0.1.1 rests skills on school maths (R31 in docs/decisions.md).
var Populations = []string{"G0", "G1", "G2", "G2-half", "G3", "F-model", "W-shock", "A-linked", "E-base"}

// numbers are the numbers of the gate of each population, besides the two of E2 on every one. The
// share of items with P* below 0.50, E1.hard, and the choices of nothing, E2.nothing, are what the
// thresholds of the CI of modules are taken from (R27 in docs/decisions.md).
var numbers = map[string][]string{
	"G0":       {"M1.rms.200", "E1.corridor", "E1.hard", "M4.false.all", "M6.late", "M6.never"},
	"G1":       {"M7.error.10"},
	"G2":       {"M3.lag", "E1.corridor"},
	"G2-half":  {"M3.lag", "E1.corridor", "E1.hard"},
	"G3":       {"M3.unsettled"},
	"A-linked": {"M9.hasFalse", "M9.notFalse"},
	"F-model":  {"M2.brier.after7", "M5.caught", "M5.never", "E3.late"},
	"W-shock":  {"M4.false.pages"},
	"E-base":   {"E5.written", "E5.again"},
}

// Numbers gives the numbers of the gate of a population: its own, then the share of learners stuck
// and the share of the choices of nothing while the pool has a task to give.
func Numbers(population string) []string {
	return append(slices.Clone(numbers[population]), "E2.stuck", "E2.nothing")
}

// Config gives the run of the gate, or a part of it: its learners under the engine, each run twice
// for I1, the libraries named by the commits the recipe gives.
func Config(part bench.Part, commits map[string]string) bench.Config {
	return bench.Config{Populations: Populations, Learners: Learners, Set: seed.Working, Rule: run.Engine, Commits: commits, Part: part, Repeat: true}
}

// Value is a number of the gate read off a run, on its cell.
type Value struct {
	Cell, Number string
	measure.Value
}

// Values reads the numbers of the gate off a run of it, in the order of its cells and numbers.
func Values(results *bench.File) ([]Value, error) {
	cells, err := bench.Cells(bench.Side{File: results}, bench.Side{File: results})
	if err != nil {
		return nil, err
	}
	all := slices.Concat(measure.ModelNumbers(), measure.EngineNumbers())
	var out []Value
	for _, c := range cells {
		var ns []measure.Number
		for _, name := range Numbers(c.Population) {
			ns = append(ns, all[slices.IndexFunc(all, func(n measure.Number) bool { return n.Name == name })])
		}
		for _, v := range measure.Read(c.Cell, c.Graph, c.A, ns) {
			v.Resampled = nil
			out = append(out, Value{Cell: c.Cell, Number: v.Name, Value: v})
		}
	}
	return out, nil
}

// whole tells whether a run is the whole run of the gate: the working learners under the engine,
// the invariants recorded, every cell of the gate with all its learners, and no part.
func whole(results *bench.File) []string {
	var out []string
	h := results.Header
	if h.Set != seed.Working || h.Rule != "engine" || !h.Invariants || h.Part != (bench.Part{}) {
		out = append(out, fmt.Sprintf("the run is not the gate's: set %s, rule %s, invariants %v, part %d of %d", h.Set.Experiment, h.Rule, h.Invariants, h.Part.K, h.Part.N))
	}
	if len(h.Packages) > 0 {
		out = append(out, "the run took packages in place of the pinned ones: it is the gate of the CI of modules, held to thresholds, not bands")
	}
	count := map[string]int{}
	for _, l := range results.Lines {
		count[l.Population+"/"+strings.Join(l.Graph, "+")]++
	}
	for _, cell := range Cells() {
		if count[cell] != learners {
			out = append(out, fmt.Sprintf("cell %s has %d learners, not %d", cell, count[cell], learners))
		}
	}
	return out
}

// Cells gives the cells of the gate: every population on each of its graphs.
func Cells() []string {
	return bench.CellsOf(Populations)
}

// Row is one number of the gate against its band.
type Row struct {
	Cell, Number string
	Band         *Band // nil: the gate file has no band for it
	Now          Value
	InBand       bool
	Exact        bool // the number is the recorded one to the last bit
}

// Report is the gate's verdict: every number against its band, and why it is red, if it is.
type Report struct {
	Rows    []Row
	Reasons []string
	Exact   bool // the numbers are held to the last bit, not to their bands
}

// OK tells whether the gate is green.
func (r *Report) OK() bool { return len(r.Reasons) == 0 }

// Check checks a run of the gate against the gate file: the run is the gate's whole run, no learner
// breaks an invariant, and every number is within its band — two half-widths of the recorded value,
// the edges included —, or, held exactly, the recorded value to the last bit. A number the recording
// could not read is held as not read. The populations with no published graph are not run, and the
// report names them (docs/architecture.md, section 3.7).
func Check(f *File, results *bench.File, exact bool) (*Report, error) {
	r := &Report{Exact: exact, Reasons: whole(results)}
	for _, l := range results.Lines {
		if slices.ContainsFunc(l.Graph, synthetic.IsSynthetic) {
			r.Reasons = append(r.Reasons, "the run has a synthetic module: "+strings.Join(l.Graph, "+"))
			break
		}
	}
	for _, l := range results.Lines {
		for _, b := range l.Breaches {
			r.Reasons = append(r.Reasons, "an invariant: "+b)
		}
	}
	values, err := Values(results)
	if err != nil {
		return nil, err
	}
	for _, v := range values {
		row := Row{Cell: v.Cell, Number: v.Number, Now: v}
		if i := slices.IndexFunc(f.Bands.Numbers, func(b Band) bool { return b.Cell == v.Cell && b.Number == v.Number }); i >= 0 {
			b := f.Bands.Numbers[i]
			row.Band = &b
			switch {
			case !b.Read:
				row.InBand, row.Exact = !v.Read, !v.Read
			case v.Read:
				row.InBand = math.Abs(v.Value.Value-b.Value) <= float64(2*b.Half)
				row.Exact = math.Float64bits(v.Value.Value) == math.Float64bits(b.Value)
			}
		}
		switch {
		case row.Band == nil:
			r.Reasons = append(r.Reasons, fmt.Sprintf("%s %s has no band in the gate file", v.Cell, v.Number))
		case exact && !row.Exact:
			r.Reasons = append(r.Reasons, fmt.Sprintf("%s %s is %s, not the recorded %s to the last bit", v.Cell, v.Number, show(v.Value.Value, v.Read), show(row.Band.Value, row.Band.Read)))
		case !row.InBand:
			r.Reasons = append(r.Reasons, fmt.Sprintf("%s %s is %s, out of its band %s ± 2 × %s", v.Cell, v.Number, show(v.Value.Value, v.Read), show(row.Band.Value, row.Band.Read), show(row.Band.Half, row.Band.Read)))
		}
		r.Rows = append(r.Rows, row)
	}
	return r, nil
}

// Record gives the gate file with the bands of a run of the gate, recorded for a reason: the value
// of every number and the half-width of its interval (docs/architecture.md, section 3.7). The run
// must be on the commits the file names: bands of a caller's code are no record.
func Record(f *File, results *bench.File, reason string) (*File, error) {
	if reason == "" {
		return nil, fmt.Errorf("gate: bands are recorded for a reason")
	}
	for _, path := range slices.Sorted(maps.Keys(f.Libraries)) {
		if got, want := results.Header.Libraries[path], f.Libraries[path].Commit; got != want {
			return nil, fmt.Errorf("gate: the run took %s at %q, and the gate file names %s", path, got, want)
		}
	}
	if rs := whole(results); len(rs) > 0 {
		return nil, fmt.Errorf("gate: %s", strings.Join(rs, "; "))
	}
	for _, l := range results.Lines {
		if len(l.Breaches) > 0 {
			return nil, fmt.Errorf("gate: an invariant: %s", l.Breaches[0])
		}
	}
	values, err := Values(results)
	if err != nil {
		return nil, err
	}
	out := *f
	out.Bands = Bands{Reason: reason}
	for _, v := range values {
		b := Band{Cell: v.Cell, Number: v.Number, Read: v.Read}
		if v.Read {
			b.Value, b.Half = v.Value.Value, float64((v.High-v.Low)/2)
		}
		out.Bands.Numbers = append(out.Bands.Numbers, b)
	}
	return &out, nil
}

// show writes a number in its shortest form that reads back to the same bits, or "not read".
func show(x float64, read bool) string {
	if !read {
		return "not read"
	}
	return strconv.FormatFloat(x, 'g', -1, 64)
}

// Write writes the report in Markdown: the table of every number against its band, then the
// reasons the gate is red, or that it is green.
func (r *Report) Write(w io.Writer) error {
	var b strings.Builder
	b.WriteString("# The gate\n\n| cell | number | recorded | half-width | now | in band | to the bit |\n|---|---|---|---|---|---|---|\n")
	for _, row := range r.Rows {
		rec, half := "no band", ""
		if row.Band != nil {
			rec, half = show(row.Band.Value, row.Band.Read), show(row.Band.Half, row.Band.Read)
		}
		in, bit := yes(row.InBand), yes(row.Exact)
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n", row.Cell, row.Number, rec, half, show(row.Now.Value.Value, row.Now.Read), in, bit)
	}
	mode := "within their bands"
	if r.Exact {
		mode = "the recorded ones to the last bit"
	}
	if r.OK() {
		fmt.Fprintf(&b, "\nThe gate is green: every number is %s, and no learner breaks an invariant.\n", mode)
	} else {
		fmt.Fprintf(&b, "\nThe gate is red, the numbers held %s:\n\n", mode)
		for _, reason := range r.Reasons {
			fmt.Fprintf(&b, "- %s\n", reason)
		}
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func yes(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// cell is a population of the gate on one of its graphs.
type cell struct {
	population string
	graph      []string
}

// cells gives every population of the gate on each of its graphs, in the order of Populations.
func cells() []cell {
	var out []cell
	for _, name := range Populations {
		p, ok := population.Named(name)
		if !ok {
			panic("gate: no population " + name) // Populations are names of the registry
		}
		for _, g := range p.Graphs {
			out = append(out, cell{name, g})
		}
	}
	return out
}
