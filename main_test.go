package main

import (
	"bytes"
	"os"
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

// The usage names the command form and says that no command is there yet, so a caller that runs
// sim before its commands come is told why it has nothing to run.
func TestUsage(t *testing.T) {
	for _, want := range []string{"usage: sim <command>", "no commands yet"} {
		if !strings.Contains(usage, want) {
			t.Errorf("usage does not say %q:\n%s", want, usage)
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
