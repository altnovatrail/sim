// Command sim is NovaTrail's learner simulator: virtual learners of different populations answer
// tasks through learner-model and engine, and the numbers they give hold every change to either
// before it is merged (docs/architecture.md). It is a tool of CI, not a library: the parts of the
// bench live in internal/ packages (R10 in docs/decisions.md), and a caller that builds sim with
// go.work runs it as `go run github.com/altnovatrail/sim <command>`.
//
// The commands come with the tasks of epic E2 in RUN.md; until then sim only says how it is used.
package main

import (
	"fmt"
	"io"
	"os"
)

// usage is what sim prints when asked how it is used.
const usage = `usage: sim <command>

sim is NovaTrail's learner simulator. It has no commands yet: they come with
the tasks of epic E2 in RUN.md.
`

// exit ends the program with a status; a test puts its own in its place, since a test cannot start
// sim as a process of its own on every target it runs on (js/wasm has none).
var exit = os.Exit

func main() {
	exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run carries out the command args name and gives the exit status: 0 when it is done, 2 when the
// command is not one sim knows.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return 0
	}
	switch args[0] {
	case "help", "-h", "-help", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	}
	fmt.Fprintf(stderr, "sim: unknown command %q\n\n%s", args[0], usage)
	return 2
}
