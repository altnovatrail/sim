package bench

import (
	"bytes"
	"errors"
	"maps"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/engine"

	"github.com/altnovatrail/sim/internal/invariant"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
)

func runTo(t *testing.T, c Config) *File {
	t.Helper()
	var buf bytes.Buffer
	if err := Run(c, &buf); err != nil {
		t.Fatal(err)
	}
	f, err := Read(&buf)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

// A file of results names its set, rule and libraries, and holds a line a learner of every graph of
// every population, in order; run again, it is the same bytes.
func TestRun(t *testing.T) {
	c := Config{Populations: []string{"G7", "S-predicate"}, Learners: 3, Set: seed.Working}
	var a, b bytes.Buffer
	if err := Run(c, &a); err != nil {
		t.Fatal(err)
	}
	if err := Run(c, &b); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a.Bytes(), b.Bytes()) {
		t.Error("two runs of one config give two files")
	}
	f, err := Read(&a)
	if err != nil {
		t.Fatal(err)
	}
	if f.Header.Set != seed.Working || f.Header.Rule != "engine" || f.Header.Libraries[Libraries[0]] != "v0.1.8" || f.Header.Libraries[Libraries[1]] != "v0.1.6" {
		t.Errorf("the header %+v", f.Header)
	}
	var cells []string
	for _, l := range f.Lines {
		if len(cells) == 0 || cells[len(cells)-1] != cellOf(l) {
			cells = append(cells, cellOf(l))
		}
		if l.Sums["learners"] != 1 || l.Oracle["learners"] != 1 || len(l.Breaches) != 0 {
			t.Errorf("%s/%d has no sums, or breaches %v", cellOf(l), l.N, l.Breaches)
		}
	}
	want := []string{"G7/novatrail/grade-3-math", "S-predicate/synthetic/sheet"}
	if strings.Join(cells, " ") != strings.Join(want, " ") || len(f.Lines) != 6 {
		t.Errorf("cells %v of %d lines", cells, len(f.Lines))
	}
	// The note of consolidation in a lesson is read under the engine alone: the lessons of school maths
	// give its sums, and the ceiling's file has none (R49).
	if l := f.Lines[0]; l.Sums["E6.passes"] <= 0 || l.Sums["E6.lesson.items"] <= 0 {
		t.Errorf("no note on %s/%d: %v", cellOf(l), l.N, l.Sums)
	}
	ceil := runTo(t, Config{Populations: []string{"G7"}, Learners: 1, Set: seed.Working, Rule: run.Ceiling})
	if ceil.Header.Rule != "ceiling" {
		t.Errorf("the ceiling's rule %q", ceil.Header.Rule)
	}
	for key := range ceil.Lines[0].Sums {
		if strings.HasPrefix(key, "E6.") {
			t.Errorf("the ceiling's line has %s", key)
		}
	}
}

// Compare reads two files by learner: a file against itself differs by nothing, to the last bit; a
// graph with a synthetic module says so; files of other learners are refused.
func TestCompare(t *testing.T) {
	c := Config{Populations: []string{"G7", "S-predicate"}, Learners: 4, Set: seed.Working}
	f := runTo(t, c)
	var out bytes.Buffer
	if err := Compare(Side{File: f}, Side{File: f}, &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if lines[0] != "cell\tnumber\ta\tb\ta-b\tlow\thigh" {
		t.Errorf("the head %q", lines[0])
	}
	read, noGraph := 0, 0
	for _, l := range lines[1:] {
		col := strings.Split(l, "\t")
		if len(col) != 7 {
			t.Fatalf("line %q", l)
		}
		if col[4] != "-" {
			read++
			if col[4] != "0" || col[5] != "0" || col[6] != "0" || col[2] != col[3] {
				t.Errorf("a file against itself: %q", l)
			}
		}
		if strings.Contains(col[0], "(no graph)") {
			noGraph++
		}
	}
	if read == 0 || noGraph == 0 {
		t.Errorf("%d numbers read, %d of no graph", read, noGraph)
	}
	out.Reset()
	if err := Compare(Side{File: f, Oracle: true}, Side{File: f}, &out); err != nil || !strings.Contains(out.String(), "M1.rms.200\t0\t") || strings.Contains(out.String(), "\tE") {
		t.Errorf("the oracle against the engine reads the measures of the model alone: %v", err)
	}
	out.Reset()
	ceil := runTo(t, Config{Populations: []string{"G7", "S-predicate"}, Learners: 4, Set: seed.Working, Rule: run.Ceiling})
	if err := Compare(Side{File: ceil}, Side{File: f}, &out); err != nil || strings.Contains(out.String(), "\tE3.") || !strings.Contains(out.String(), "\tE1.corridor\t") {
		t.Errorf("the ceiling of the engine against the engine reads no E3: %v", err)
	}

	other := runTo(t, Config{Populations: []string{"G7"}, Learners: 4, Set: seed.Working})
	short := runTo(t, Config{Populations: []string{"G7", "S-predicate"}, Learners: 3, Set: seed.Working})
	held, err := seed.HeldOut(seed.CriterionModel, 1)
	if err != nil {
		t.Fatal(err)
	}
	moved := *f
	moved.Header.Set = held
	swapped := *f
	swapped.Lines = append(append([]Line{}, f.Lines[1:]...), f.Lines[0])
	for name, b := range map[string]*File{"other cells": other, "fewer learners": short, "another set": &moved, "another order": &swapped} {
		if err := Compare(Side{File: f}, Side{File: b}, &out); !errors.Is(err, ErrBench) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// Names takes a library's version from the build, its replacement's, or the commit the caller
// names, and refuses a library the build took from a workspace with no commit named.
func TestNames(t *testing.T) {
	info := &debug.BuildInfo{Deps: []*debug.Module{
		{Path: Libraries[0], Version: "v0.1.0"},
		{Path: Libraries[1], Version: "v0.1.0", Replace: &debug.Module{Path: "../engine", Version: "v0.0.0-20261108-abcdef"}},
	}}
	got, err := Names(info, nil)
	if err != nil || got[Libraries[0]] != "v0.1.0" || got[Libraries[1]] != "../engine@v0.0.0-20261108-abcdef" {
		t.Errorf("Names = %v, %v", got, err)
	}
	same := &debug.BuildInfo{Deps: []*debug.Module{
		{Path: Libraries[0], Version: "v0.1.0", Replace: &debug.Module{Path: Libraries[0], Version: "v0.1.1"}},
		{Path: Libraries[1], Version: "v0.1.0"},
	}}
	if got, err := Names(same, nil); err != nil || got[Libraries[0]] != "v0.1.1" {
		t.Errorf("a replacement by another version of the same module: %v, %v", got, err)
	}
	devel := &debug.BuildInfo{Deps: []*debug.Module{{Path: Libraries[0], Version: "(devel)"}}}
	if _, err := Names(devel, nil); !errors.Is(err, ErrBench) {
		t.Errorf("a library of a workspace with no commit: %v", err)
	}
	commits := map[string]string{Libraries[0]: strings.Repeat("a", 40), Libraries[1]: strings.Repeat("b", 40)}
	if got, err := Names(nil, commits); err != nil || got[Libraries[1]] != commits[Libraries[1]] {
		t.Errorf("named commits: %v, %v", got, err)
	}
	for name, c := range map[string]map[string]string{
		"a module that is no library":        {"github.com/altnovatrail/spec": strings.Repeat("a", 40)},
		"a short commit":                     {Libraries[0]: "abc123"},
		"a branch":                           {Libraries[0]: "main"},
		"a commit beside a version it knows": {Libraries[0]: strings.Repeat("a", 40)},
	} {
		if _, err := Names(info, c); !errors.Is(err, ErrBench) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// What Run and Read refuse.
func TestRefuses(t *testing.T) {
	var buf bytes.Buffer
	for name, c := range map[string]Config{
		"no population": {Populations: []string{"Z9"}, Learners: 1, Set: seed.Working},
		"no rule":       {Populations: []string{"G7"}, Learners: 1, Set: seed.Working, Rule: run.Rule(9)},
	} {
		if err := Run(c, &buf); !errors.Is(err, ErrBench) {
			t.Errorf("%s: %v", name, err)
		}
	}
	if _, err := modules([]string{"novatrail/none"}, nil); !errors.Is(err, ErrBench) {
		t.Errorf("a module that is not there: %v", err)
	}
	for name, in := range map[string]string{"empty": "", "no header": "{}\n", "not json": "{\"rule\":\"engine\"}\n{\n"} {
		if _, err := Read(strings.NewReader(in)); !errors.Is(err, ErrBench) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// failing is a writer that fails once it has taken some bytes.
type failing struct{ left int }

func (f *failing) Write(p []byte) (int, error) {
	if len(p) > f.left {
		return 0, errors.New("full")
	}
	f.left -= len(p)
	return len(p), nil
}

// A writer that fails stops Run and Compare with its error, at the header or at a line.
func TestWriteFails(t *testing.T) {
	c := Config{Populations: []string{"G7"}, Learners: 1, Set: seed.Working}
	f := runTo(t, c)
	for _, left := range []int{0, 300} {
		if err := Run(c, &failing{left: left}); err == nil {
			t.Errorf("Run into a writer full after %d bytes: no error", left)
		}
	}
	for _, left := range []int{0, 60} {
		if err := Compare(Side{File: f}, Side{File: f}, &failing{left: left}); err == nil {
			t.Errorf("Compare into a writer full after %d bytes: no error", left)
		}
	}
}

// The parts of a run, merged, are the run of the whole to the byte; Merge refuses parts of other
// runs, a part twice and a part missing; a part out of range is refused.
func TestParts(t *testing.T) {
	// Past Rechosen learners a cell, so that a repeated run checks both ways.
	c := Config{Populations: []string{"G7", "S-predicate"}, Learners: Rechosen + 2, Set: seed.Working, Repeat: true}
	var whole bytes.Buffer
	if err := Run(c, &whole); err != nil {
		t.Fatal(err)
	}
	var parts []*File
	for k := range 3 {
		p := c
		p.Part = Part{K: k, N: 3}
		parts = append(parts, runTo(t, p))
	}
	slices.Reverse(parts)
	merged, err := Merge(parts)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := Write(merged, &out); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out.Bytes(), whole.Bytes()) {
		t.Errorf("the merged parts are not the whole:\n%.300s\nwant\n%.300s", out.String(), whole.String())
	}
	for _, l := range merged.Lines {
		if len(l.Breaches) != 0 {
			t.Errorf("%s/%d: %v", cellOf(l), l.N, l.Breaches)
		}
	}
	other := runTo(t, Config{Populations: []string{"G7"}, Learners: 2, Set: seed.Working, Rule: run.Ceiling, Part: Part{K: 1, N: 3}})
	fewer := c
	fewer.Learners, fewer.Part = 4, Part{K: 2, N: 3}
	outside := *parts[1]
	outside.Header.Part = Part{K: 5, N: 3}
	for name, ps := range map[string][]*File{
		"none":           nil,
		"twice":          {parts[0], parts[0], parts[1]},
		"missing":        {parts[0], parts[1]},
		"another run":    {parts[0], parts[1], other},
		"fewer learners": {parts[0], parts[1], runTo(t, fewer)},
		"a part outside": {parts[0], &outside, parts[2]},
	} {
		if _, err := Merge(ps); !errors.Is(err, ErrBench) {
			t.Errorf("%s: %v", name, err)
		}
	}
	twice := c
	twice.Populations = []string{"G7", "G7"}
	if err := Run(twice, &out); !errors.Is(err, ErrBench) {
		t.Errorf("a population twice: %v", err)
	}
	for _, p := range []Part{{K: 3, N: 3}, {K: -1, N: 2}, {K: 0, N: -1}, {K: 2, N: 0}} {
		bad := c
		bad.Part = p
		if err := Run(bad, &out); !errors.Is(err, ErrBench) {
			t.Errorf("part %+v: %v", p, err)
		}
	}
	if err := Write(merged, &failing{left: 10}); err == nil {
		t.Error("Write into a full writer: no error")
	}
	if err := Write(merged, &failing{left: 400}); err == nil {
		t.Error("Write into a writer full at a line: no error")
	}
}

// A run on a package in place of the pinned one runs only the graphs the package covers, names the
// package in its header, reads the same as the pinned run when the package is the pinned one, and
// its parts merge in its order.
func TestPackages(t *testing.T) {
	pub, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	pins, err := published.Pinned()
	if err != nil {
		t.Fatal(err)
	}
	pin := pins[slices.IndexFunc(pins, func(p published.Package) bool { return p.Module == "novatrail/grade-3-math" })]
	taken := []Package{{Package: pin, Module: pub[pin.Module]}}
	c := Config{Populations: []string{"G0", "E-base"}, Learners: 3, Set: seed.Working, Packages: taken}
	got := runTo(t, c)
	if !slices.Equal(got.Header.Packages, []published.Package{pin}) {
		t.Errorf("the header names the packages %v", got.Header.Packages)
	}
	if want := []string{"G0/novatrail/grade-3-math"}; !slices.Equal(CellsOn(c.Populations, got.Header.Packages), want) {
		t.Errorf("the cells on the package are %v", CellsOn(c.Populations, got.Header.Packages))
	}
	// A package that came along runs only on the graphs it shares with one that did not: olympiad
	// maths along with school maths is A-linked on both, and none of its own.
	olympiad := pins[slices.IndexFunc(pins, func(p published.Package) bool { return p.Module == "novatrail/olympiad-math" })]
	olympiad.Along = true
	both := []published.Package{pin, olympiad}
	if want := []string{"G0/novatrail/grade-3-math", "A-linked/novatrail/olympiad-math+novatrail/grade-3-math"}; !slices.Equal(CellsOn([]string{"G0", "A-linked"}, both), want) {
		t.Errorf("the cells with a package along are %v", CellsOn([]string{"G0", "A-linked"}, both))
	}
	if got := CellsOn([]string{"G0", "A-linked"}, []published.Package{olympiad}); got != nil {
		t.Errorf("a package along alone runs %v", got)
	}
	plain := runTo(t, Config{Populations: []string{"G0"}, Learners: 3, Set: seed.Working})
	var same []Line
	for _, l := range plain.Lines {
		if cellOf(l) == "G0/novatrail/grade-3-math" {
			same = append(same, l)
		}
	}
	if len(got.Lines) != 3 || !slices.EqualFunc(got.Lines, same, func(a, b Line) bool { return a.N == b.N && maps.Equal(a.Sums, b.Sums) }) {
		t.Errorf("the run on the pinned package is not the pinned run: %d lines", len(got.Lines))
	}
	var parts []*File
	for k := range 2 {
		p := c
		p.Part = Part{K: k, N: 2}
		parts = append(parts, runTo(t, p))
	}
	merged, err := Merge(parts)
	if err != nil || len(merged.Lines) != 3 {
		t.Errorf("the parts on the package merge: %v", err)
	}
	pinned := c
	pinned.Packages, pinned.Part = nil, Part{K: 1, N: 2}
	if _, err := Merge([]*File{parts[0], runTo(t, pinned)}); !errors.Is(err, ErrBench) {
		t.Errorf("a part on the package merges with a part on the pin: %v", err)
	}
	var buf bytes.Buffer
	for name, ps := range map[string][]Package{
		"no module read":    {{Package: pin}},
		"another module":    {{Package: pin, Module: pub["novatrail/olympiad-math"]}},
		"the package twice": {taken[0], taken[0]},
	} {
		bad := c
		bad.Packages = ps
		if err := Run(bad, &buf); !errors.Is(err, ErrBench) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// Every run with invariants takes each choice again, I1, for the first Rechosen learners of a cell
// alone: a choice changed in the record is a breach of I1 for the last of them, and none past them.
func TestRechosen(t *testing.T) {
	p, _ := population.Named("G7")
	c := Config{Populations: []string{p.Name}, Learners: 1, Set: seed.Working}
	mods, err := modules(p.Graphs[0], nil)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := runOne(c, p, mods, 0)
	if err != nil {
		t.Fatal(err)
	}
	bad := *rec
	bad.Steps = slices.Clone(rec.Steps)
	bad.Steps[0].Choice.Kind = engine.ChoiceNothing
	i1 := func(bs []invariant.Breach) bool {
		return slices.ContainsFunc(bs, func(b invariant.Breach) bool { return b.Invariant == "I1" })
	}
	if !i1(checkOf(Rechosen - 1)(&bad)) {
		t.Errorf("learner %d: a choice changed is no breach of I1", Rechosen-1)
	}
	if i1(checkOf(Rechosen)(&bad)) || len(checkOf(Rechosen)(rec)) != 0 || len(checkOf(0)(rec)) != 0 {
		t.Errorf("learner %d takes the choices again, or the run breaks an invariant", Rechosen)
	}
}
