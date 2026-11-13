package criterion

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/seed"
)

// Path is the path of a number on a cell, fixed: the value of the base and of the ceiling, both
// turned by the number's direction, so that the path is base − ceiling. Read tells that both were
// read; the ceiling is zero for a number whose path goes to perfection, the ceiling of the engine
// for the engine's numbers and the oracle of the base for the model's.
type Path struct {
	Cell    string  `json:"cell"`
	Number  string  `json:"number"`
	Read    bool    `json:"read"`
	Base    float64 `json:"base"`
	Ceiling float64 `json:"ceiling"`
}

// Paths are the paths the criteria hold, read once on the first decisive run of the base, before
// the first candidate, and fixed (docs/architecture.md, section 3.6; the amendment of 2026-11-09 to
// R04, R28 in docs/decisions.md): the libraries, the set and the learners of the runs they were
// read on, the SHA-256 of the files of results of the base and of the ceiling of the engine, and the
// path of every number of the groups of the score on every cell of a published graph.
type Paths struct {
	Libraries map[string]string `json:"libraries"`
	Set       seed.Set          `json:"set"`
	Learners  int               `json:"learners"`
	Base      string            `json:"base"`
	Ceiling   string            `json:"ceiling"`
	Paths     []Path            `json:"paths"`
}

//go:embed paths.json
var pathsJSON []byte

// fixed gives the paths the criteria hold, those of paths.json, read once; a test puts its own in
// their place.
var fixed = sync.OnceValues(func() (*Paths, error) { return ReadPaths(bytes.NewReader(pathsJSON)) })

// ReadPaths reads a file of paths.
func ReadPaths(r io.Reader) (*Paths, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var p Paths
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("%w: the paths: %w", ErrCriterion, err)
	}
	return &p, nil
}

// Write writes a file of paths: JSON indented by two spaces, a number in its shortest form that
// reads back to the same bits.
func (p *Paths) Write(w io.Writer) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

// on gives the fixed paths of a cell, by number.
func (p *Paths) on(cell string) map[string]Path {
	out := map[string]Path{}
	for _, x := range p.Paths {
		if x.Cell == cell {
			out[x.Number] = x
		}
	}
	return out
}

// decisive is the learners of every cell a run the paths are read on must have; a test lowers it.
var decisive = Decisive

// covered are the populations a run the paths are read on must have: every population with a weight
// in a group of the score, but those on synthetic graphs alone; a test puts its own in their place.
var covered = func() []string {
	var out []string
	for _, g := range slices.Concat(modelGroups, engineGroups) {
		for _, w := range g.populations {
			if !strings.HasPrefix(w.population, "A-") {
				out = append(out, w.population)
			}
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}()

// pathNumbers gives the numbers with a path on a population: those of the groups of the score of
// both criteria where the population has a weight, in the order of their names.
func pathNumbers(population string) []string {
	var out []string
	for _, g := range slices.Concat(modelGroups, engineGroups) {
		if slices.ContainsFunc(g.populations, func(w weighted) bool { return w.population == population }) {
			out = append(out, g.numbers...)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// RecordPaths reads the paths off the decisive runs of the base and of the ceiling of the engine,
// whose files of results have the SHA-256 given: on the working learners, the whole runs and not a
// part, on the pinned packages, with every population of the groups of the score, at least a
// decisive run's learners on every cell, on the same libraries and the same learners, the base with
// no breach of an invariant. A cell on a graph with no published package
// has no path (R19).
func RecordPaths(base, ceiling *bench.File, baseDigest, ceilingDigest string) (*Paths, error) {
	b, c := base.Header, ceiling.Header
	switch {
	case b.Rule != "engine" || c.Rule != "ceiling":
		return nil, fmt.Errorf("%w: the base ran under %q and the ceiling under %q", ErrCriterion, b.Rule, c.Rule)
	case b.Set != seed.Working || c.Set != seed.Working:
		return nil, fmt.Errorf("%w: the paths are read on the working learners, not on %s and %s", ErrCriterion, b.Set.Experiment, c.Set.Experiment)
	case b.Part != (bench.Part{}) || c.Part != (bench.Part{}):
		return nil, fmt.Errorf("%w: the paths are read on whole runs, not on parts", ErrCriterion)
	case len(b.Packages) > 0 || len(c.Packages) > 0:
		return nil, fmt.Errorf("%w: the paths are read on the pinned packages, not on packages in their place", ErrCriterion)
	case !maps.Equal(b.Libraries, c.Libraries):
		return nil, fmt.Errorf("%w: the base ran on %v and the ceiling on %v", ErrCriterion, b.Libraries, c.Libraries)
	}
	for _, l := range base.Lines {
		if len(l.Breaches) > 0 {
			return nil, fmt.Errorf("%w: the base breaks an invariant: %s", ErrCriterion, l.Breaches[0])
		}
	}
	for _, pop := range covered {
		if !slices.ContainsFunc(base.Lines, func(l bench.Line) bool { return l.Population == pop }) {
			return nil, fmt.Errorf("%w: the base has no learner of %s, whose paths the score needs", ErrCriterion, pop)
		}
	}
	againstCeiling, err := bench.Cells(bench.Side{File: base}, bench.Side{File: ceiling})
	if err != nil {
		return nil, fmt.Errorf("%w: the ceiling against the base: %w", ErrCriterion, err)
	}
	againstOracle, err := bench.Cells(bench.Side{File: base}, bench.Side{File: base, Oracle: true})
	if err != nil {
		return nil, fmt.Errorf("%w: the oracle against the base: %w", ErrCriterion, err)
	}
	p := &Paths{Libraries: b.Libraries, Set: b.Set, Base: baseDigest, Ceiling: ceilingDigest}
	for k, cell := range againstCeiling {
		if p.Learners == 0 || len(cell.A) < p.Learners {
			p.Learners = len(cell.A)
		}
		if len(cell.A) < decisive {
			return nil, fmt.Errorf("%w: cell %s has %d learners, and the paths are read on a decisive run of %d", ErrCriterion, cell.Cell, len(cell.A), decisive)
		}
		ns := pathNumbers(cell.Population)
		if noGraph(cell.Graph) || len(ns) == 0 {
			continue
		}
		numbers := directed(ns)
		bases := measure.Read(cell.Cell, cell.Graph, cell.A, numbers)
		ceilings := measure.Read(cell.Cell, cell.Graph, cell.B, numbers)
		oracles := measure.Read(cell.Cell, cell.Graph, againstOracle[k].B, numbers)
		for i, name := range ns {
			ceil := oracles[i]
			if strings.HasPrefix(name, "E") {
				ceil = ceilings[i]
			}
			x := Path{Cell: cell.Cell, Number: name, Read: bases[i].Read && ceil.Read}
			if defs[name].perfect {
				x.Read = bases[i].Read
				ceil = measure.Value{}
			}
			if x.Read {
				x.Base, x.Ceiling = bases[i].Value, ceil.Value
			}
			p.Paths = append(p.Paths, x)
		}
	}
	return p, nil
}

// Table gives the rows of the table of the paths in Markdown, as the amendment to R04 that fixes
// them holds them: the cell, the number, the base, the ceiling and the path, to four decimals with
// a decimal point, as the documents in English write numbers; a path not read is a dash.
func (p *Paths) Table() string {
	show := func(x float64) string {
		return strconv.FormatFloat(x, 'f', 4, 64)
	}
	var b strings.Builder
	for _, x := range p.Paths {
		if !x.Read {
			fmt.Fprintf(&b, "| `%s` | `%s` | — | — | — |\n", x.Cell, x.Number)
			continue
		}
		fmt.Fprintf(&b, "| `%s` | `%s` | %s | %s | %s |\n", x.Cell, x.Number, show(x.Base), show(x.Ceiling), show(float64(x.Base-x.Ceiling)))
	}
	return b.String()
}
