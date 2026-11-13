package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/gate"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/spec/types/formatv0"
)

// "sim gate" lists the libraries of the gate file, runs a part of the gate, and checks or records
// the parts merged; the check holds the build to the file first. A part alone is no whole run: the
// check is red and the record refuses.
func TestGate(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if got := run([]string{"gate", "libraries", "gate.json"}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim gate libraries = %d: %s", got, stderr.String())
	}
	want := "github.com/altnovatrail/engine altnovatrail/engine 57e0ab8852e65bb56724d08fc2e9cdb6a839d959 v0.1.6\n" +
		"github.com/altnovatrail/learner-model altnovatrail/learner-model a063835087c7df7f2eecd1912c44afda7d96bfe8\n"
	if stdout.String() != want {
		t.Errorf("sim gate libraries wrote\n%s", stdout.String())
	}
	dir := t.TempDir()
	part := filepath.Join(dir, "part.jsonl")
	// The part Learners of Learners+1 holds no learner of a cell: the gate on the learners of olympiad
	// maths 0.1.1 takes minutes, and what is checked here is a part that is not the whole.
	if got := run([]string{"gate", "run", "-part", fmt.Sprintf("%d/%d", gate.Learners, gate.Learners+1), "-out", part}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim gate run = %d: %s", got, stderr.String())
	}
	other := filepath.Join(dir, "other.jsonl") // never written
	gateFile := filepath.Join(dir, "gate.json")
	data, err := os.ReadFile("gate.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(gateFile, data, 0o600); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if got := run([]string{"gate", "check", "-file", gateFile, "-go", ".", part}, &stdout, &stderr); got != 1 || !strings.Contains(stderr.String(), "is missing") {
		t.Errorf("sim gate check of one part = %d: %s", got, stderr.String())
	}
	if got := run([]string{"gate", "record", "-reason", "a test", "-file", gateFile, "-go", ".", part}, &stdout, &stderr); got != 1 {
		t.Errorf("sim gate record of a part = %d", got)
	}
	if after, _ := os.ReadFile(gateFile); !bytes.Equal(after, data) {
		t.Error("a refused record changed the gate file")
	}
	for _, c := range []struct {
		args   []string
		status int
	}{
		{[]string{"gate"}, 2},
		{[]string{"gate", "nope"}, 2},
		{[]string{"gate", "libraries"}, 2},
		{[]string{"gate", "libraries", "-nope", "gate.json"}, 2},
		{[]string{"gate", "libraries", filepath.Join(dir, "none.json")}, 1},
		{[]string{"gate", "libraries", part}, 1},
		{[]string{"gate", "run", "-part", "2/1"}, 2},
		{[]string{"gate", "run", "-part", "x"}, 2},
		{[]string{"gate", "check", "-file", gateFile}, 2},
		{[]string{"gate", "check", "-nope", part}, 2},
		{[]string{"gate", "check", "-file", filepath.Join(dir, "none.json"), part}, 1},
		{[]string{"gate", "check", "-file", gateFile, "-go", dir, part}, 1},
		{[]string{"gate", "check", "-file", gateFile, "-go", ".", other}, 1},
		{[]string{"gate", "check", "-file", gateFile, "-go", ".", gateFile}, 1},
	} {
		stderr.Reset()
		if got := run(c.args, &stdout, &stderr); got != c.status {
			t.Errorf("sim %q = %d, want %d: %s", c.args, got, c.status, stderr.String())
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("go 1.27.3\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := run([]string{"gate", "check", "-file", gateFile, "-go", dir, part}, &stdout, &stderr); got != 1 {
		t.Errorf("a go.mod with no module path = %d", got)
	}
}

// "sim run -part" runs a part, and "sim merge" merges the parts into the whole, to the byte.
func TestMerge(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	whole := filepath.Join(dir, "whole.jsonl")
	if got := run([]string{"run", "-populations", "G7", "-learners", "4", "-out", whole}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim run = %d: %s", got, stderr.String())
	}
	var parts []string
	for _, p := range []string{"1/2", "0/2"} {
		path := filepath.Join(dir, strings.ReplaceAll(p, "/", "-")+".jsonl")
		if got := run([]string{"run", "-populations", "G7", "-learners", "4", "-part", p, "-out", path}, &stdout, &stderr); got != 0 {
			t.Fatalf("sim run -part %s = %d: %s", p, got, stderr.String())
		}
		parts = append(parts, path)
	}
	merged := filepath.Join(dir, "merged.jsonl")
	if got := run(append([]string{"merge", "-out", merged}, parts...), &stdout, &stderr); got != 0 {
		t.Fatalf("sim merge = %d: %s", got, stderr.String())
	}
	a, _ := os.ReadFile(whole)
	b, _ := os.ReadFile(merged)
	if !bytes.Equal(a, b) {
		t.Error("the merged parts are not the whole")
	}
	if st, err := os.Stat(merged); err != nil || st.Mode().Perm() != 0o644 {
		t.Errorf("a new file of results has the mode %v: %v", st.Mode().Perm(), err)
	}
	for _, c := range []struct {
		args   []string
		status int
	}{
		{[]string{"merge"}, 2},
		{[]string{"merge", "-nope", parts[0]}, 2},
		{[]string{"merge", filepath.Join(dir, "none.jsonl")}, 1},
		{[]string{"merge", parts[0]}, 1},
	} {
		if got := run(c.args, &stdout, &stderr); got != c.status {
			t.Errorf("sim %q = %d, want %d", c.args, got, c.status)
		}
	}
}

// A whole run of the gate, its parts merged, checks green against its own record and red against a
// gate file with a band taken out; the record writes the bands into the gate file. The run is built
// by hand, every learner of every cell with no sums, so that every number is held as not read; the
// build is the workspace's.
func TestGateWhole(t *testing.T) {
	prior := buildInfo
	t.Cleanup(func() { buildInfo = prior })
	buildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Deps: []*debug.Module{
			{Path: "github.com/altnovatrail/engine", Version: "(devel)"},
			{Path: "github.com/altnovatrail/learner-model", Version: "(devel)"},
			{Path: "github.com/altnovatrail/spec", Version: specRecorded(t)},
		}}, true
	}
	dir := t.TempDir()
	var parts []string
	for k := range 2 {
		f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Invariants: true, Part: bench.Part{K: k, N: 2},
			Populations: gate.Populations, Learners: gate.Learners, Repeat: true, Libraries: map[string]string{
				"github.com/altnovatrail/engine":        "57e0ab8852e65bb56724d08fc2e9cdb6a839d959",
				"github.com/altnovatrail/learner-model": "a063835087c7df7f2eecd1912c44afda7d96bfe8",
			}}}
		for _, cell := range gate.Cells() {
			pop, graph, _ := strings.Cut(cell, "/")
			for n := k; n < gate.Learners; n += 2 {
				f.Lines = append(f.Lines, bench.Line{Population: pop, Graph: strings.Split(graph, "+"), N: n, Sums: measure.Sums{"learners": 1}})
			}
		}
		path := filepath.Join(dir, fmt.Sprintf("part-%d.jsonl", k))
		var buf bytes.Buffer
		if err := bench.Write(f, &buf); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
			t.Fatal(err)
		}
		parts = append(parts, path)
	}
	gateFile := filepath.Join(dir, "gate.json")
	data, err := os.ReadFile("gate.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(gateFile, data, 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	args := append([]string{"-file", gateFile, "-go", "."}, parts...)
	if got := run(append([]string{"gate", "record"}, args...), &stdout, &stderr); got != 1 {
		t.Errorf("a record for no reason = %d", got)
	}
	// The go.mod of a library: the gate reads the go line of every module of the workspace.
	for _, name := range []string{"engine", "learner-model"} {
		lib := filepath.Join(dir, name)
		if err := os.MkdirAll(lib, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(lib, "go.mod"), []byte("module github.com/altnovatrail/"+name+"\n\ngo 1.27.3\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		args = append([]string{"-go", lib}, args...)
	}
	stderr.Reset()
	if got := run(append([]string{"gate", "record", "-reason", "a test"}, args...), &stdout, &stderr); got != 0 {
		t.Fatalf("sim gate record = %d: %s", got, stderr.String())
	}
	if st, err := os.Stat(gateFile); err != nil || st.Mode().Perm() != 0o600 {
		t.Errorf("the record changed the mode of the gate file: %v", err)
	}
	stdout.Reset()
	if got := run(append([]string{"gate", "check", "-exact"}, args...), &stdout, &stderr); got != 0 || !strings.Contains(stdout.String(), "The gate is green") {
		t.Errorf("sim gate check -exact = %d:\n%.400s", got, stdout.String())
	}
	f, err := readGate(gateFile)
	if err != nil {
		t.Fatal(err)
	}
	f.Bands.Numbers = f.Bands.Numbers[1:]
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(gateFile, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if got := run(append([]string{"gate", "check"}, args...), &stdout, &stderr); got != 1 || !strings.Contains(stdout.String(), "has no band") {
		t.Errorf("sim gate check with a band out = %d:\n%.400s", got, stdout.String())
	}
	if got := run([]string{"gate", "record", "-reason", "a test", "-file", gateFile, "-go", ".", parts[0], parts[1]}, &stdout, &stderr); got != 1 {
		t.Errorf("a record with the build apart from the file = %d", got)
	}
	// A whole run of the gate of modules on the pinned packages, no number read: against the bands of
	// gate.json, the thresholds are red.
	pins, err := published.Pinned()
	if err != nil {
		t.Fatal(err)
	}
	modules := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Invariants: true, Populations: gate.Populations, Learners: gate.Learners, Repeat: true, Packages: pins,
		Libraries: map[string]string{
			"github.com/altnovatrail/engine":        "57e0ab8852e65bb56724d08fc2e9cdb6a839d959",
			"github.com/altnovatrail/learner-model": "a063835087c7df7f2eecd1912c44afda7d96bfe8",
		}}}
	for _, cell := range bench.CellsOn(gate.Populations, pins) {
		pop, graph, _ := strings.Cut(cell, "/")
		for n := range gate.Learners {
			modules.Lines = append(modules.Lines, bench.Line{Population: pop, Graph: strings.Split(graph, "+"), N: n, Sums: measure.Sums{"learners": 1}})
		}
	}
	var buf2 bytes.Buffer
	if err := bench.Write(modules, &buf2); err != nil {
		t.Fatal(err)
	}
	whole := filepath.Join(dir, "modules.jsonl")
	if err := os.WriteFile(whole, buf2.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if got := run([]string{"gate", "thresholds", "-go", filepath.Join(dir, "engine"), "-go", filepath.Join(dir, "learner-model"), "-go", ".", "-file", "gate.json", whole}, &stdout, &stderr); got != 1 || !strings.Contains(stdout.String(), "The gate of modules is red") || !strings.Contains(stdout.String(), "E1.hard is not read") {
		t.Errorf("sim gate thresholds on a whole run with no number read = %d:\n%.600s%s", got, stdout.String(), stderr.String())
	}
}

// folders writes the pinned packages into folders, as the CI of modules unpacks the published ones,
// each named by its package digest, and gives them by module.
func folders(t *testing.T, dir string) map[string]string {
	t.Helper()
	pins, err := published.Pinned()
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, p := range pins {
		name := strings.Replace(p.Digest, ":", "-", 1)
		archive, err := published.Archive(p)
		if err != nil {
			t.Fatal(err)
		}
		files, err := pack.Untar(archive)
		if err != nil {
			t.Fatal(err)
		}
		folder := filepath.Join(dir, name)
		if err := os.MkdirAll(folder, 0o700); err != nil {
			t.Fatal(err)
		}
		for _, f := range []string{"inventory.json", "server.json"} {
			if err := os.WriteFile(filepath.Join(folder, f), files[f], 0o600); err != nil {
				t.Fatal(err)
			}
		}
		out[p.Module] = folder
	}
	return out
}

// "sim gate changed" writes the folders the gate of modules runs on, and what each package is;
// "sim gate run -package" runs a part on a package; "sim gate thresholds" refuses a part alone and
// holds a whole run of the gate of modules, here with no number read, red.
func TestGateModules(t *testing.T) {
	dir := t.TempDir()
	pkgs := folders(t, dir)
	index := filepath.Join("internal", "published", published.Release, "module-index.json")
	grade3 := pkgs["novatrail/grade-3-math"]
	olympiad := pkgs["novatrail/olympiad-math"]
	// An index with no school maths: it is new, and A-linked brings olympiad maths, which did not
	// change, into the run with it.
	var idx formatv0.ModuleIndex
	data, err := os.ReadFile(index)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatal(err)
	}
	delete(idx.Modules, "novatrail/grade-3-math")
	data, err = json.Marshal(idx)
	if err != nil {
		t.Fatal(err)
	}
	noGrade3 := filepath.Join(dir, "no-grade-3.json")
	if err := os.WriteFile(noGrade3, data, 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	for _, c := range []struct {
		args   []string
		stdout string
		stderr string
	}{
		{[]string{"-first", grade3}, grade3 + "\n", "novatrail/grade-3-math 0.3.0: new\n"},
		{[]string{"-index", index, "-published", grade3, grade3}, "", "novatrail/grade-3-math 0.3.0: same against 0.3.0\n"},
		{[]string{"-index", index, "-published", grade3, "-all", grade3}, grade3 + "\n", "novatrail/grade-3-math 0.3.0: same against 0.3.0\n"},
		{[]string{"-index", noGrade3, "-published", olympiad, grade3, olympiad}, grade3 + "\nalong " + olympiad + "\n",
			"novatrail/grade-3-math 0.3.0: new\nnovatrail/olympiad-math 0.3.0: same against 0.3.0, run with novatrail/grade-3-math, which shares a graph of the gate with it\n"},
	} {
		stdout.Reset()
		stderr.Reset()
		if got := run(append([]string{"gate", "changed"}, c.args...), &stdout, &stderr); got != 0 || stdout.String() != c.stdout || stderr.String() != c.stderr {
			t.Errorf("sim gate changed %q = %d:\n%s%s", c.args, got, stdout.String(), stderr.String())
		}
	}
	part := filepath.Join(dir, "part.jsonl")
	if got := run([]string{"gate", "run", "-part", "0/300", "-package", grade3, "-out", part}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim gate run -package = %d: %s", got, stderr.String())
	}
	f, err := readResults([]string{part})
	if err != nil {
		t.Fatal(err)
	}
	if len(f[0].Header.Packages) != 1 || slices.ContainsFunc(f[0].Lines, func(l bench.Line) bool { return l.Graph[0] != "novatrail/grade-3-math" }) {
		t.Errorf("the part on grade-3-math: %+v", f[0].Header)
	}
	notJSON := filepath.Join(dir, "index.json")
	if err := os.WriteFile(notJSON, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args   []string
		status int
	}{
		{[]string{"gate", "changed"}, 2},
		{[]string{"gate", "changed", grade3}, 2},
		{[]string{"gate", "changed", "-first", "-index", index, grade3}, 2},
		{[]string{"gate", "changed", "-first", "-published", grade3, grade3}, 2},
		{[]string{"gate", "changed", "-nope", grade3}, 2},
		{[]string{"gate", "changed", "-index", filepath.Join(dir, "none.json"), grade3}, 1},
		{[]string{"gate", "changed", "-index", notJSON, grade3}, 1},
		{[]string{"gate", "changed", "-index", index, "-published", dir, grade3}, 1},
		{[]string{"gate", "changed", "-first", dir}, 1},
		{[]string{"gate", "changed", "-index", index, grade3}, 1},
		{[]string{"gate", "run", "-package", dir}, 1},
		{[]string{"gate", "run", "-package", grade3, "-package", grade3}, 1},
		{[]string{"gate", "thresholds", "-go", ".", part}, 1},
		{[]string{"gate", "thresholds", "-file", "gate.json"}, 2},
	} {
		stderr.Reset()
		if got := run(c.args, &stdout, &stderr); got != c.status {
			t.Errorf("sim %q = %d, want %d: %s", c.args, got, c.status, stderr.String())
		}
	}
}

// "sim gate verify" holds the build alone to the gate file: a build of the workspace with the
// go.mod of every module passes, a build of the module cache or a module missing fails, and so do
// flags it does not take.
func TestGateVerify(t *testing.T) {
	prior := buildInfo
	t.Cleanup(func() { buildInfo = prior })
	workspace := &debug.BuildInfo{Deps: []*debug.Module{
		{Path: "github.com/altnovatrail/engine", Version: "(devel)"},
		{Path: "github.com/altnovatrail/learner-model", Version: "(devel)"},
		{Path: "github.com/altnovatrail/spec", Version: specRecorded(t)},
	}}
	buildInfo = func() (*debug.BuildInfo, bool) { return workspace, true }
	dir := t.TempDir()
	args := []string{"gate", "verify", "-file", "gate.json", "-go", "."}
	for _, name := range []string{"engine", "learner-model"} {
		lib := filepath.Join(dir, name)
		if err := os.MkdirAll(lib, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(lib, "go.mod"), []byte("module github.com/altnovatrail/"+name+"\n\ngo 1.27.3\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		args = append(args, "-go", lib)
	}
	var stdout, stderr bytes.Buffer
	if got := run(args, &stdout, &stderr); got != 0 || !strings.Contains(stdout.String(), "The build holds to the gate file") {
		t.Fatalf("sim gate verify on the workspace = %d: %s%s", got, stdout.String(), stderr.String())
	}
	stdout.Reset()
	if got := run(args[:len(args)-2], &stdout, &stderr); got != 1 || !strings.Contains(stdout.String(), "the gate was given no go.mod of github.com/altnovatrail/learner-model") {
		t.Errorf("sim gate verify without a module = %d: %s", got, stdout.String())
	}
	workspace = &debug.BuildInfo{Deps: []*debug.Module{
		{Path: "github.com/altnovatrail/engine", Version: "v0.1.4"},
		{Path: "github.com/altnovatrail/learner-model", Version: "v0.1.4"},
		{Path: "github.com/altnovatrail/spec", Version: "v0.11.0"},
	}}
	stdout.Reset()
	if got := run(args, &stdout, &stderr); got != 1 || !strings.Contains(stdout.String(), "github.com/altnovatrail/engine is built at v0.1.4") || !strings.Contains(stdout.String(), "chose github.com/altnovatrail/spec at v0.11.0") {
		t.Errorf("sim gate verify on the module cache = %d: %s", got, stdout.String())
	}
	for _, bad := range [][]string{
		{"gate", "verify", "-file", "gate.json", "extra"},
		{"gate", "verify", "-nope"},
	} {
		if got := run(bad, &stdout, &stderr); got != 2 {
			t.Errorf("run(%q) = %d, want 2", bad, got)
		}
	}
	if got := run([]string{"gate", "verify", "-file", filepath.Join(dir, "none.json")}, &stdout, &stderr); got != 1 {
		t.Errorf("sim gate verify of no gate file = %d, want 1", got)
	}
	if got := run([]string{"gate", "verify", "-go", dir}, &stdout, &stderr); got != 1 {
		t.Errorf("sim gate verify with a folder that has no go.mod = %d, want 1", got)
	}
}

// specRecorded is the version of spec's module that the gate file records: a build of the
// workspace that the tests stand in for chose it, and it moves with each move to a tag of spec.
func specRecorded(t *testing.T) string {
	t.Helper()
	f, err := gate.Read(bytes.NewReader(gateJSON))
	if err != nil {
		t.Fatal(err)
	}
	v := f.Dependencies["github.com/altnovatrail/spec"]
	if v == "" {
		t.Fatal("the gate file records no version of github.com/altnovatrail/spec")
	}
	return v
}
