package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
			args:   []string{"gate", "-update"},
			status: 2,
			stderr: "sim: unknown command \"gate\"\n\n" + usage,
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
	for _, want := range []string{"usage: sim <command>", "\n  run ", "\n  compare "} {
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
	if got := run([]string{"run", "-populations", "G0", "-learners", "2", "-out", a}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim run = %d: %s", got, stderr.String())
	}
	if got := run([]string{"run", "-populations", "G0", "-learners", "1", "-rule", "ceiling", "-set", "held-out-model-1"}, &stdout, &stderr); got != 0 {
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
	if got := run([]string{"run", "-populations", "G0", "-learners", "1"}, &stdout, &stderr); got != 0 || !strings.Contains(stdout.String(), `"rule":"engine"`) {
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
		{[]string{"run", "-learners", "1", "-populations", "G0", "-commit", "github.com/altnovatrail/engine=" + strings.Repeat("e", 40)}, 1},
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
