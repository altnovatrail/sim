// Command sim is NovaTrail's learner simulator: virtual learners of different populations answer
// tasks through learner-model and engine, and the numbers they give hold every change to either
// before it is merged (docs/architecture.md). It is a tool of CI, not a library: the parts of the
// bench live in internal/ packages (R10 in docs/decisions.md), and a caller that builds sim with
// go.work runs it as `go run github.com/altnovatrail/sim <command>`.
//
// Its commands run cells of learners into a file of results and compare two such files by learner
// (R24); the rest come with the tasks of epic E2 in RUN.md.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/population"
	runpkg "github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
)

// usage is what sim prints when asked how it is used.
const usage = `usage: sim <command> [flags]

sim is NovaTrail's learner simulator. Its commands:

  run      run cells of virtual learners and write a file of results, one
           line a learner, to stdout or -out
  compare  read two files of results in pairs by learner and write each
           number of each cell, a, b and a-b with its 95 % interval

Run "sim <command> -h" for its flags.
`

// exit ends the program with a status; a test puts its own in its place, since a test cannot start
// sim as a process of its own on every target it runs on (js/wasm has none).
var exit = os.Exit

func main() {
	exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run carries out the command args name and gives the exit status: 0 when it is done, 1 when it
// fails, 2 when the command or its flags are not ones sim knows.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return 0
	}
	var err error
	switch args[0] {
	case "help", "-h", "-help", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	case "run":
		err = runCommand(args[1:], stdout, stderr)
	case "compare":
		err = compareCommand(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "sim: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
	switch {
	case err == nil, errors.Is(err, flag.ErrHelp):
		return 0
	case errors.Is(err, errUsage):
		return 2
	}
	fmt.Fprintf(stderr, "sim: %v\n", err)
	return 1
}

// errUsage is flags sim does not take; the flag package has said why.
var errUsage = errors.New("usage")

// parse parses a command's flags, and takes its usage error as errUsage.
func parse(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return err
		}
		return errUsage
	}
	return nil
}

// runCommand is "sim run".
func runCommand(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	pops := fs.String("populations", "all", "the populations, by name and comma-separated, or all")
	learners := fs.Int("learners", 100, "the learners of each graph of each population")
	set := fs.String("set", "working", "the set of learners: working, or held-out-<criterion>-<n>")
	rule := fs.String("rule", "engine", "the rule that chooses: engine, or ceiling")
	out := fs.String("out", "", "the file of results; stdout when empty")
	commits := map[string]string{}
	fs.Func("commit", "a library the build names no version of, as <module path>=<full SHA>; repeatable", func(v string) error {
		path, sha, ok := strings.Cut(v, "=")
		if !ok || path == "" || sha == "" {
			return fmt.Errorf("want <module path>=<full SHA>, got %q", v)
		}
		commits[path] = sha
		return nil
	})
	if err := parse(fs, args); err != nil {
		return err
	}
	c := bench.Config{Learners: *learners, Commits: commits}
	if *pops == "all" {
		for _, p := range population.All() {
			c.Populations = append(c.Populations, p.Name)
		}
	} else {
		c.Populations = strings.Split(*pops, ",")
	}
	var err error
	if c.Set, err = seed.Named(*set); err != nil {
		return fmt.Errorf("-set %s: %w", *set, err)
	}
	switch *rule {
	case "engine":
		c.Rule = runpkg.Engine
	case "ceiling":
		c.Rule = runpkg.Ceiling
	default:
		return fmt.Errorf("no rule %q: engine or ceiling", *rule)
	}
	if *learners <= 0 {
		return fmt.Errorf("-learners %d: want at least one", *learners)
	}
	if *out == "" {
		return bench.Run(c, stdout)
	}
	// The file is written beside its place and put there only once the run is done, so that a run
	// that fails leaves a file that was there as it was.
	f, err := os.CreateTemp(filepath.Dir(*out), filepath.Base(*out)+".*")
	if err != nil {
		return err
	}
	err = bench.Run(c, f)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), *out)
	}
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}

// compareCommand is "sim compare".
func compareCommand(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	fs.SetOutput(stderr)
	aOracle := fs.Bool("a-oracle", false, "read the first file as the oracle's sums, the ceiling of the model")
	bOracle := fs.Bool("b-oracle", false, "read the second file as the oracle's sums")
	if err := parse(fs, args); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		fmt.Fprintln(stderr, "usage: sim compare [-a-oracle] [-b-oracle] <a> <b>")
		return errUsage
	}
	var sides [2]bench.Side
	for k, path := range fs.Args() {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		file, err := bench.Read(f)
		f.Close()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		sides[k] = bench.Side{File: file}
	}
	sides[0].Oracle, sides[1].Oracle = *aOracle, *bOracle
	return bench.Compare(sides[0], sides[1], stdout)
}
