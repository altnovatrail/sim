package gate

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/spec/types/formatv0"
)

// The gate of the CI of modules judges content, not a library, and a change of content moves the
// numbers by design: it holds the invariants and thresholds, not bands, on the packages a pull
// request built, and only on the modules where what the bench reads changed (docs/architecture.md,
// section 3.9; the amendment of 2026-11-07 to R04, the amendments to R05 and R27 in
// docs/decisions.md).

// State is what a package of a pull request is against the published versions of its module.
type State string

const (
	// New is a module with no published version, or with none to compare when there is no release
	// yet: it changed whole.
	New State = "new"
	// Changed is a package whose canonical text differs from that of the newest published version.
	Changed State = "changed"
	// Same is a package the bench reads the same of as the newest published version.
	Same State = "same"
	// Transition is a package of another version of the format than the newest published version:
	// a transition of the format, where the bench does not run (the amendment of 2026-11-07 to R103
	// spec).
	Transition State = "transition"
)

// Package is a package as a folder holds it: the folder, and its inventory.json and server.json.
type Package struct {
	Folder            string
	Inventory, Server []byte
}

// Change is a package of a pull request and what it is: its module and version, and the published
// version it was compared with, if any. Read is the package read, when the bench reads its format.
type Change struct {
	State           State
	Module, Version string
	Against         string
	Folder          string
	Read            *bench.Package
}

// ErrModules is a comparison or a run of the gate of modules the bench refuses.
var ErrModules = errors.New("gate: modules")

// Changes tells, for each package of a pull request, whether what the bench reads of it changed
// against the newest published version of its module that is not withdrawn — the one api gives
// (R63 spec) —, by their canonical texts (pack.Module.Text). idx is the index of the newest release
// of modules, which release-verify has accepted, and pub the packages of its versions, each
// checked by the chain of digests from it; with no index, there is no release yet, and every
// package is new. A package of another version of the format than that version is a transition,
// whose content the bench does not compare across the formats; a package of a version of the format
// the bench does not read, which is no transition, is refused: it wants a newer tag of sim. A
// package that is broken — a server part apart from its inventory, a manifest of another module —
// is refused, a transition too.
func Changes(idx *formatv0.ModuleIndex, pub, packages []Package) ([]Change, error) {
	byDigest := map[formatv0.Digest]Package{}
	for _, p := range pub {
		byDigest[pack.Digest(p.Inventory)] = p
	}
	var out []Change
	for _, p := range packages {
		read, inv, err := readOne(p)
		if inv == nil || err != nil && !errors.Is(err, pack.ErrFormat) {
			return nil, err
		}
		if slices.ContainsFunc(out, func(c Change) bool { return c.Module == string(inv.ID) }) {
			return nil, fmt.Errorf("%w: %s: a second package of %s", ErrModules, p.Folder, inv.ID)
		}
		c := Change{State: New, Module: string(inv.ID), Version: string(inv.Version), Folder: p.Folder}
		version, entry, ok := newest(idx, inv.ID)
		if ok {
			c.Against = string(version)
			if entry.ModuleFormatVersion != inv.ModuleFormatVersion {
				c.State = Transition
				out = append(out, c)
				continue
			}
		}
		if err != nil {
			return nil, err
		}
		c.Read = read
		if ok {
			was, err := publishedText(byDigest, inv.ID, version, entry)
			if err != nil {
				return nil, err
			}
			now, err := read.Module.Text()
			if err != nil {
				return nil, fmt.Errorf("%w: %s: %w", ErrModules, p.Folder, err)
			}
			c.State = Same
			if !bytes.Equal(was, now) {
				c.State = Changed
			}
		}
		out = append(out, c)
	}
	return out, nil
}

// ReadPackages reads the packages of a pull request as a run takes them, in place of the pinned ones: each
// checked against its inventory and translated.
func ReadPackages(packages []Package) ([]bench.Package, error) {
	var out []bench.Package
	for _, p := range packages {
		read, _, err := readOne(p)
		if err != nil {
			return nil, err
		}
		out = append(out, *read)
	}
	return out, nil
}

// readOne reads a package of a pull request: the inventory, whenever it reads, and the package
// translated, when the bench reads its version of the format.
func readOne(p Package) (*bench.Package, *formatv0.Inventory, error) {
	cv, inv, err := pack.ReadFiles(p.Inventory, p.Server, nil)
	if err != nil {
		return nil, inv, fmt.Errorf("%w: %s: %w", ErrModules, p.Folder, err)
	}
	m, err := pack.Translate(cv)
	if err != nil {
		return nil, inv, fmt.Errorf("%w: %s: %w", ErrModules, p.Folder, err)
	}
	id := published.Package{Module: string(inv.ID), Version: string(inv.Version), Digest: string(pack.Digest(p.Inventory))}
	return &bench.Package{Package: id, Module: m}, inv, nil
}

// newest gives the newest version of a module in an index that is not withdrawn, by Semantic
// Versioning — X.Y.Z, with no pre-release, as the format writes a version (the text of the format,
// section 11).
func newest(idx *formatv0.ModuleIndex, id formatv0.ModuleID) (formatv0.ModuleVersion, formatv0.Entry, bool) {
	var best formatv0.ModuleVersion
	var entry formatv0.Entry
	if idx == nil {
		return best, entry, false
	}
	for _, v := range slices.Sorted(maps.Keys(idx.Modules[id])) {
		e := idx.Modules[id][v]
		if e.Withdrawn == nil && (best == "" || semver(v, best) > 0) {
			best, entry = v, e
		}
	}
	return best, entry, best != ""
}

// semver compares two versions X.Y.Z part by part, as numbers; a part that is no number counts as
// 0, which the index's form does not let in.
func semver(a, b formatv0.ModuleVersion) int {
	pa, pb := strings.Split(string(a), "."), strings.Split(string(b), ".")
	for k := range max(len(pa), len(pb)) {
		var x, y int
		if k < len(pa) {
			x, _ = strconv.Atoi(pa[k])
		}
		if k < len(pb) {
			y, _ = strconv.Atoi(pb[k])
		}
		if c := cmp.Compare(x, y); c != 0 {
			return c
		}
	}
	return 0
}

// publishedText gives the canonical text of a published version, from the published package whose
// inventory.json is its entry's package digest.
func publishedText(byDigest map[formatv0.Digest]Package, id formatv0.ModuleID, version formatv0.ModuleVersion, e formatv0.Entry) ([]byte, error) {
	p, ok := byDigest[e.PackageDigest]
	if !ok {
		return nil, fmt.Errorf("%w: no published package is %s %s, %s", ErrModules, id, version, e.PackageDigest)
	}
	cv, inv, err := pack.ReadFiles(p.Inventory, p.Server, &e)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrModules, p.Folder, err)
	}
	if inv.ID != id || inv.Version != version {
		return nil, fmt.Errorf("%w: %s is %s %s, and the index names %s %s by its digest", ErrModules, p.Folder, inv.ID, inv.Version, id, version)
	}
	m, err := pack.Translate(cv)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrModules, p.Folder, err)
	}
	text, err := m.Text()
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrModules, p.Folder, err)
	}
	return text, nil
}

// ModulesConfig gives the run of the gate of modules, or a part of it: the gate's run, on the
// packages of a pull request in place of the pinned ones, so on the graphs they cover alone.
func ModulesConfig(part bench.Part, commits map[string]string, packages []bench.Package) bench.Config {
	c := Config(part, commits)
	c.Packages = packages
	return c
}

// tolerance is the author's tolerance of the shares of E1 and E2, one percentage point
// (docs/architecture.md, section 3.6).
const tolerance = 0.01

// threshold is a threshold of the gate of modules: a number, on the cells of populations — every
// population when none —, held at most or at least to a limit. The limit is a value of its own, or
// is taken from the bands of the gate file at the same tag of sim, on the cells of the same
// populations: the largest recorded value plus the tolerance, or the smallest less it.
type threshold struct {
	number      string
	populations []string
	most        bool
	limit       float64
	from        string // "largest", "smallest", or "" for a limit of its own
}

// thresholds are the thresholds of the gate of modules (docs/architecture.md, section 3.9): the
// learners stuck and the choices of nothing no more than on the published graphs; on G0 and
// G2-half, the share of items with P* below 0.50 no more, and the share in the corridor no less,
// than on the published graphs (the amendment of 2026-11-09 to R05); no more than 20 % false
// "mastered" on G0, the floor of the model's criterion; no task with a rubric given again.
var thresholds = []threshold{
	{number: "E2.stuck", most: true, from: "largest"},
	{number: "E2.nothing", most: true, from: "largest"},
	{number: "E1.hard", populations: []string{"G0", "G2-half"}, most: true, from: "largest"},
	{number: "E1.corridor", populations: []string{"G0", "G2-half"}, from: "smallest"},
	{number: "M4.false.all", populations: []string{"G0"}, most: true, limit: 0.20},
	{number: "E5.again", populations: []string{"E-base"}, most: true, limit: 0},
}

// on tells whether a threshold stands on a cell.
func (t threshold) on(cell string) bool {
	pop, _, _ := strings.Cut(cell, "/")
	return len(t.populations) == 0 || slices.Contains(t.populations, pop)
}

// limitIn gives the limit of a threshold from a gate file, or why there is none.
func (t threshold) limitIn(f *File) (float64, error) {
	if t.from == "" {
		return t.limit, nil
	}
	var values []float64
	for _, b := range f.Bands.Numbers {
		if b.Number == t.number && b.Read && t.on(b.Cell) {
			values = append(values, b.Value)
		}
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("the gate file has no %s read to take its threshold from", t.number)
	}
	if t.from == "largest" {
		return slices.Max(values) + tolerance, nil
	}
	return slices.Min(values) - tolerance, nil
}

// Held is one number of the gate of modules against its threshold.
type Held struct {
	Cell, Number string
	Most         bool
	Limit        float64
	Now          Value
	OK           bool
}

// ThresholdReport is the verdict of the gate of modules: every number against its threshold, why it
// is red, if it is, and the modules on which no cell of the gate runs.
type ThresholdReport struct {
	Packages  []published.Package
	Rows      []Held
	Reasons   []string
	Uncovered []string
}

// OK tells whether the gate of modules is green.
func (r *ThresholdReport) OK() bool { return len(r.Reasons) == 0 }

// Thresholds checks a run of the gate of modules against the thresholds, their limits from the gate
// file: the run is the gate's on packages, on the commits of the libraries the file names, with
// every cell the packages cover and all its learners; no learner breaks an invariant; every number
// with a threshold is within it, the edges included, and a number the run could not read is not. A
// package on which no cell of the gate runs — a module the bench has no population on — is named,
// not red: populations on a new graph are a change of sim.
func Thresholds(f *File, results *bench.File) (*ThresholdReport, error) {
	h := results.Header
	r := &ThresholdReport{Packages: h.Packages}
	if len(h.Packages) == 0 {
		return nil, fmt.Errorf("%w: the run took no package in place of the pinned ones", ErrModules)
	}
	if h.Set != seed.Working || h.Rule != "engine" || !h.Invariants || !h.Repeat || h.Part != (bench.Part{}) || !slices.Equal(h.Populations, Populations) || h.Learners != Learners {
		r.Reasons = append(r.Reasons, fmt.Sprintf("the run is not the gate's: set %s, rule %s, invariants %v, repeat %v, part %d of %d, populations %v, learners %d",
			h.Set.Experiment, h.Rule, h.Invariants, h.Repeat, h.Part.K, h.Part.N, h.Populations, h.Learners))
	}
	for _, path := range slices.Sorted(maps.Keys(f.Libraries)) {
		if got, want := h.Libraries[path], f.Libraries[path].Commit; got != want {
			r.Reasons = append(r.Reasons, fmt.Sprintf("the run took %s at %q, and the gate file names %s", path, got, want))
		}
	}
	count := map[string]int{}
	for _, l := range results.Lines {
		count[l.Population+"/"+strings.Join(l.Graph, "+")]++
	}
	cells := bench.CellsOn(Populations, h.Packages)
	for _, cell := range cells {
		if count[cell] != learners {
			r.Reasons = append(r.Reasons, fmt.Sprintf("cell %s has %d learners, not %d", cell, count[cell], learners))
		}
	}
	for _, cell := range slices.Sorted(maps.Keys(count)) {
		if !slices.Contains(cells, cell) {
			r.Reasons = append(r.Reasons, fmt.Sprintf("cell %s is in the run, and the packages do not cover it", cell))
		}
	}
	for _, p := range h.Packages {
		if bench.CellsOn(Populations, []published.Package{p}) == nil {
			r.Uncovered = append(r.Uncovered, p.Module)
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
	for _, t := range thresholds {
		limit, err := t.limitIn(f)
		if err != nil {
			return nil, err
		}
		for _, v := range values {
			if v.Number != t.number || !t.on(v.Cell) {
				continue
			}
			row := Held{Cell: v.Cell, Number: v.Number, Most: t.most, Limit: limit, Now: v}
			row.OK = v.Read && (t.most && v.Value.Value <= limit || !t.most && v.Value.Value >= limit)
			if !row.OK {
				r.Reasons = append(r.Reasons, fmt.Sprintf("%s %s is %s, and its threshold is %s %s", v.Cell, v.Number, show(v.Value.Value, v.Read), bound(t.most), show(limit, true)))
			}
			r.Rows = append(r.Rows, row)
		}
	}
	slices.SortStableFunc(r.Rows, func(a, b Held) int {
		return cmp.Compare(slices.Index(cells, a.Cell), slices.Index(cells, b.Cell))
	})
	return r, nil
}

// bound writes which way a threshold holds.
func bound(most bool) string {
	if most {
		return "at most"
	}
	return "at least"
}

// Write writes the report in Markdown: the packages, the table of every number against its
// threshold, the modules no cell runs on, then the reasons the gate is red, or that it is green.
func (r *ThresholdReport) Write(w io.Writer) error {
	var b strings.Builder
	b.WriteString("# The gate of modules\n\nOn the packages:\n\n")
	for _, p := range r.Packages {
		fmt.Fprintf(&b, "- %s %s, %s\n", p.Module, p.Version, p.Digest)
	}
	b.WriteString("\n| cell | number | threshold | now | held |\n|---|---|---|---|---|\n")
	for _, row := range r.Rows {
		fmt.Fprintf(&b, "| %s | %s | %s %s | %s | %s |\n", row.Cell, row.Number, bound(row.Most), show(row.Limit, true), show(row.Now.Value.Value, row.Now.Read), yes(row.OK))
	}
	for _, p := range NoGraph {
		var held []string
		for _, t := range thresholds {
			if t.on(p + "/") {
				held = append(held, t.number)
			}
		}
		fmt.Fprintf(&b, "| %s | %s | no graph | | |\n", p, strings.Join(held, ", "))
	}
	for _, m := range r.Uncovered {
		fmt.Fprintf(&b, "\nNo cell of the gate runs on %s: populations on a new graph are a change of sim.\n", m)
	}
	if r.OK() {
		b.WriteString("\nThe gate of modules is green: every number is within its threshold, and no learner breaks an invariant.\n")
	} else {
		b.WriteString("\nThe gate of modules is red:\n\n")
		for _, reason := range r.Reasons {
			fmt.Fprintf(&b, "- %s\n", reason)
		}
	}
	_, err := io.WriteString(w, b.String())
	return err
}
