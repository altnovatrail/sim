package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/gate"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/seed"
)

// "sim gate" lists the libraries of the gate file, runs a part of the gate, and checks or records
// the parts merged; the check holds the build to the file first. A part alone is no whole run: the
// check is red and the record refuses.
func TestGate(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if got := run([]string{"gate", "libraries", "gate.json"}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim gate libraries = %d: %s", got, stderr.String())
	}
	want := "github.com/altnovatrail/engine altnovatrail/engine c6c5bb9bc2c6fc486f629da82ee834bd89fc2b35 v0.1.0\n" +
		"github.com/altnovatrail/learner-model altnovatrail/learner-model d2620a18d865c813ec2556c93fca92a1445e3a74 v0.1.0\n"
	if stdout.String() != want {
		t.Errorf("sim gate libraries wrote\n%s", stdout.String())
	}
	dir := t.TempDir()
	part := filepath.Join(dir, "part.jsonl")
	if got := run([]string{"gate", "run", "-part", "0/150", "-out", part}, &stdout, &stderr); got != 0 {
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
		t.Errorf("sim gate check of one part of 150 = %d: %s", got, stderr.String())
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
	if got := run([]string{"run", "-populations", "G0", "-learners", "4", "-out", whole}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim run = %d: %s", got, stderr.String())
	}
	var parts []string
	for _, p := range []string{"1/2", "0/2"} {
		path := filepath.Join(dir, strings.ReplaceAll(p, "/", "-")+".jsonl")
		if got := run([]string{"run", "-populations", "G0", "-learners", "4", "-part", p, "-out", path}, &stdout, &stderr); got != 0 {
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
			{Path: "github.com/altnovatrail/spec", Version: "v0.9.9"},
		}}, true
	}
	dir := t.TempDir()
	var parts []string
	for k := range 2 {
		f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Invariants: true, Part: bench.Part{K: k, N: 2},
			Populations: gate.Populations, Learners: gate.Learners, Repeat: true, Libraries: map[string]string{
				"github.com/altnovatrail/engine":        "c6c5bb9bc2c6fc486f629da82ee834bd89fc2b35",
				"github.com/altnovatrail/learner-model": "d2620a18d865c813ec2556c93fca92a1445e3a74",
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
}
