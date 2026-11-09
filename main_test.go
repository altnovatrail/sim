package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/criterion"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/seed"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		status int
		stdout string
		stderr string
	}{
		{name: "no command", args: nil, status: 0, stdout: usage},
		{name: "help", args: []string{"help"}, status: 0, stdout: usage},
		{name: "-h", args: []string{"-h"}, status: 0, stdout: usage},
		{name: "-help", args: []string{"-help"}, status: 0, stdout: usage},
		{name: "--help", args: []string{"--help"}, status: 0, stdout: usage},
		{
			name:   "unknown command",
			args:   []string{"nope", "-update"},
			status: 2,
			stderr: "sim: unknown command \"nope\"\n\n" + usage,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if got := run(c.args, &stdout, &stderr); got != c.status {
				t.Errorf("run(%q) = %d, want %d", c.args, got, c.status)
			}
			if got := stdout.String(); got != c.stdout {
				t.Errorf("run(%q) wrote to stdout\n%s\nwant\n%s", c.args, got, c.stdout)
			}
			if got := stderr.String(); got != c.stderr {
				t.Errorf("run(%q) wrote to stderr\n%s\nwant\n%s", c.args, got, c.stderr)
			}
		})
	}
}

// The usage names the command form and its commands.
func TestUsage(t *testing.T) {
	for _, want := range []string{"usage: sim <command>", "\n  run ", "\n  compare ", "\n  judge "} {
		if !strings.Contains(usage, want) {
			t.Errorf("usage does not say %q:\n%s", want, usage)
		}
	}
}

// "sim run" writes a file of results, "sim compare" reads two by learner: a file against itself
// differs by nothing; flags sim does not take give status 2, a run it cannot do status 1.
func TestCommands(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.jsonl")
	var stdout, stderr bytes.Buffer
	if got := run([]string{"run", "-populations", "G7", "-learners", "2", "-out", a}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim run = %d: %s", got, stderr.String())
	}
	if got := run([]string{"run", "-populations", "G7", "-learners", "1", "-rule", "ceiling", "-set", "held-out-model-1"}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim run -rule ceiling = %d: %s", got, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"rule":"ceiling"`) || !strings.Contains(stdout.String(), "held-out-model-1") {
		t.Errorf("sim run to stdout: %.200s", stdout.String())
	}
	// A run that fails leaves the file that was there as it was.
	before, err := os.ReadFile(a)
	if err != nil {
		t.Fatal(err)
	}
	if got := run([]string{"run", "-populations", "Z9", "-learners", "1", "-out", a}, &stdout, &stderr); got != 1 {
		t.Errorf("sim run of no population = %d", got)
	}
	if after, err := os.ReadFile(a); err != nil || !bytes.Equal(before, after) {
		t.Errorf("a failed run changed the file: %v", err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Errorf("a failed run left %d files", len(entries))
	}
	stdout.Reset()
	if got := run([]string{"compare", "-a-oracle", a, a}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim compare = %d: %s", got, stderr.String())
	}
	if !strings.Contains(stdout.String(), "M1.rms.200\t0\t") {
		t.Errorf("the oracle against the engine: %.300s", stdout.String())
	}
	stdout.Reset()
	if got := run([]string{"run", "-populations", "G7", "-learners", "1"}, &stdout, &stderr); got != 0 || !strings.Contains(stdout.String(), `"rule":"engine"`) {
		t.Errorf("sim run with the default rule = %d", got)
	}
	stdout.Reset()
	if got := run([]string{"run", "-learners", "1", "-populations", "all", "-out", filepath.Join(dir, "missing", "x")}, &stdout, &stderr); got != 1 {
		t.Errorf("sim run into a folder that is not there = %d", got)
	}
	for _, c := range []struct {
		args   []string
		status int
	}{
		{[]string{"run", "-nope"}, 2},
		{[]string{"run", "-h"}, 0},
		{[]string{"run", "-commit", "nothing"}, 2},
		{[]string{"run", "-learners", "1", "-populations", "G7", "-commit", "github.com/altnovatrail/engine=" + strings.Repeat("e", 40)}, 1},
		{[]string{"run", "-rule", "oracle"}, 1},
		{[]string{"run", "-set", "nope"}, 1},
		{[]string{"run", "-set", "held-out-model-x"}, 1},
		{[]string{"run", "-set", "held-out-model-99"}, 1},
		{[]string{"run", "-learners", "0"}, 1},
		{[]string{"run", "-populations", "Z9", "-learners", "1"}, 1},
		{[]string{"compare", a}, 2},
		{[]string{"compare", "-nope", a, a}, 2},
		{[]string{"compare", a, filepath.Join(dir, "none")}, 1},
		{[]string{"compare", a, filepath.Join(dir, "missing.jsonl")}, 1},
	} {
		stderr.Reset()
		if got := run(c.args, &stdout, &stderr); got != c.status {
			t.Errorf("sim %q = %d, want %d: %s", c.args, got, c.status, stderr.String())
		}
	}
	bad := filepath.Join(dir, "bad.jsonl")
	if err := os.WriteFile(bad, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := run([]string{"compare", a, bad}, &stdout, &stderr); got != 1 {
		t.Errorf("sim compare with a file of no header = %d", got)
	}
}

// "sim judge" reads a judgement file — its runs relative to it — and writes the report; a file it
// cannot read or a judgement it refuses gives status 1, flags it does not take status 2.
func TestJudge(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	for _, c := range [][]string{
		{"run", "-populations", "G7", "-learners", "3", "-out", filepath.Join(dir, "base.jsonl")},
		{"run", "-populations", "G7", "-learners", "3", "-set", "held-out-model-1", "-out", filepath.Join(dir, "held.jsonl")},
	} {
		if got := run(c, &stdout, &stderr); got != 0 {
			t.Fatalf("sim %q = %d: %s", c, got, stderr.String())
		}
	}
	write := func(name, body string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	good := write("j.json", `{"criterion":"model","base":"base.jsonl",
		"candidates":[{"name":"same","file":"base.jsonl","group":"still","constants":1}],
		"confirm":{"base":"held.jsonl","candidate":"held.jsonl"}}`)
	stdout.Reset()
	if got := run([]string{"judge", good}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim judge = %d: %s", got, stderr.String())
	}
	if !strings.Contains(stdout.String(), "# The criterion of the model") || !strings.Contains(stdout.String(), "the base stays") {
		t.Errorf("the report:\n%s", stdout.String())
	}
	for _, c := range []struct {
		args   []string
		status int
	}{
		{[]string{"judge"}, 2},
		{[]string{"judge", "-nope", good}, 2},
		{[]string{"judge", filepath.Join(dir, "none.json")}, 1},
		{[]string{"judge", write("unknown.json", `{"criterion":"model","extra":1}`)}, 1},
		{[]string{"judge", write("nobase.json", `{"criterion":"model","base":"none.jsonl"}`)}, 1},
		{[]string{"judge", write("ceiling.json", `{"criterion":"model","base":"base.jsonl","ceiling":"base.jsonl"}`)}, 1},
		{[]string{"judge", write("nocand.json", `{"criterion":"model","base":"base.jsonl","candidates":[{"name":"x","file":"none.jsonl"}]}`)}, 1},
		{[]string{"judge", write("noconf.json", `{"criterion":"model","base":"base.jsonl","confirm":{"base":"none.jsonl"}}`)}, 1},
		{[]string{"judge", write("bad.json", `{"criterion":"model","base":"j.json"}`)}, 1},
		{[]string{"judge", write("refused.json", `{"criterion":"taste","base":"base.jsonl"}`)}, 1},
	} {
		stderr.Reset()
		if got := run(c.args, &stdout, &stderr); got != c.status {
			t.Errorf("sim %q = %d, want %d: %s", c.args, got, c.status, stderr.String())
		}
	}
}

// main passes the arguments after the program's name to run and exits with its status.
func TestMainExits(t *testing.T) {
	args, prior := os.Args, exit
	t.Cleanup(func() { os.Args, exit = args, prior })
	status := -1
	exit = func(code int) { status = code }
	os.Args = []string{"sim", "help"}
	main()
	if status != 0 {
		t.Errorf("main with help exited %d, want 0", status)
	}
}

// "sim paths" reads the paths off the decisive runs of the base and of the ceiling — here built by
// hand, every learner of each population the score needs, on one graph and with no sums, so that every path is not read — writes them to
// -out and their table to stdout; a rough look, a file it cannot read and flags it does not take
// are refused.
func TestPaths(t *testing.T) {
	dir := t.TempDir()
	write := func(name, rule string, n int) string {
		f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: rule, Invariants: rule == "engine"}}
		for _, pop := range []string{"G0", "G0-skills1", "G1", "G2", "G2-half", "G2-fading", "G3", "G3-drop", "S-adult", "F-model", "F-exp", "F-C1", "F-growth5", "F-break", "W-shock", "W-error", "W-badday", "L-explain"} {
			for k := range n {
				f.Lines = append(f.Lines, bench.Line{Population: pop, Graph: []string{"novatrail/olympiad-math"}, N: k, Sums: measure.Sums{"learners": 1}, Oracle: measure.Sums{"learners": 1}})
			}
		}
		var buf bytes.Buffer
		if err := bench.Write(f, &buf); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	base, ceiling := write("base.jsonl", "engine", criterion.Decisive), write("ceiling.jsonl", "ceiling", criterion.Decisive)
	out := filepath.Join(dir, "paths.json")
	var stdout, stderr bytes.Buffer
	if got := run([]string{"paths", "-base", base, "-ceiling", ceiling, "-out", out, "-table"}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim paths = %d: %s", got, stderr.String())
	}
	if !strings.Contains(stdout.String(), "| `G0/novatrail/olympiad-math` | `M1.rms.200` | — | — | — |") {
		t.Errorf("the table:\n%s", stdout.String())
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	p, err := criterion.ReadPaths(bytes.NewReader(data))
	sum := sha256.Sum256(func() []byte { b, _ := os.ReadFile(base); return b }())
	if err != nil || p.Learners != criterion.Decisive || p.Base != hex.EncodeToString(sum[:]) {
		t.Errorf("the paths: %v, %d learners, base %s", err, p.Learners, p.Base)
	}
	rough := write("rough.jsonl", "engine", 3)
	for _, c := range []struct {
		args   []string
		status int
	}{
		{[]string{"paths"}, 2},
		{[]string{"paths", "-base", base}, 2},
		{[]string{"paths", "-nope"}, 2},
		{[]string{"paths", "-base", base, "-ceiling", ceiling, "-out", out, "extra"}, 2},
		{[]string{"paths", "-base", base, "-ceiling", ceiling, "-table"}, 2},
		{[]string{"paths", "-base", filepath.Join(dir, "none.jsonl"), "-ceiling", ceiling, "-out", out}, 1},
		{[]string{"paths", "-base", base, "-ceiling", filepath.Join(dir, "none.jsonl"), "-out", out}, 1},
		{[]string{"paths", "-base", out, "-ceiling", ceiling, "-out", out}, 1},
		{[]string{"paths", "-base", rough, "-ceiling", ceiling, "-out", out}, 1},
	} {
		stderr.Reset()
		if got := run(c.args, &stdout, &stderr); got != c.status {
			t.Errorf("sim %q = %d, want %d: %s", c.args, got, c.status, stderr.String())
		}
	}
}
