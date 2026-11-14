// Command sim is NovaTrail's learner simulator: virtual learners of different populations answer
// tasks through learner-model and engine, and the numbers they give hold every change to either
// before it is merged (docs/architecture.md). It is a tool of CI, not a library: the parts of the
// bench live in internal/ packages (R10 in docs/decisions.md), and a caller that builds sim with
// go.work runs it as `go run github.com/altnovatrail/sim <command>`.
//
// Its commands run cells of learners into a file of results, compare two such files by learner
// (R24) and judge candidates by a criterion (R25); the rest come with the tasks of epic E2 in RUN.md.
package main

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/criterion"
	"github.com/altnovatrail/sim/internal/gate"
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
  judge    judge candidates against the base by a criterion, as a judgement
           file names them, and write the report
  merge    merge the files of the parts of a run into the file of the whole
  paths    read the paths of the score off the decisive runs of the base and
           of the ceiling into the file the criterion holds, and write
           their table
  gate     the gate: its libraries, a part of its run, its check against
           the gate file, and the recording of its bands

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
	case "judge":
		err = judgeCommand(args[1:], stdout, stderr)
	case "merge":
		err = mergeCommand(args[1:], stdout, stderr)
	case "paths":
		err = pathsCommand(args[1:], stdout, stderr)
	case "gate":
		err = gateCommand(args[1:], stdout, stderr)
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

// gcPercent is the garbage collector's share of the live heap a run lets grow before it collects:
// four times Go's own. A process of a run then holds about 100 MB, and a run takes about 28 % less
// processor than at Go's own 100 (R33 in docs/decisions.md).
const gcPercent = 400

// alone sets the collector of a process of a run to gcPercent, and a process that runs a part to one
// processor: the parts of a run are processes, as many as the processors, with no goroutine of the
// bench's own (R09), and another processor for each would only take the others'. A run of the whole is
// the one process on the machine, and keeps the processors for its collector. Neither changes a bit
// of what the run writes. It is a variable so that the tests of the commands leave their own process
// as it was.
var alone = func(part bench.Part) {
	if part.N > 1 {
		runtime.GOMAXPROCS(1)
	}
	debug.SetGCPercent(gcPercent)
}

// runCommand is "sim run".
func runCommand(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	pops := fs.String("populations", "all", "the populations, by name and comma-separated, or all")
	learners := fs.Int("learners", 100, "the learners of each graph of each population")
	set := fs.String("set", "working", "the set of learners: working, or held-out-<criterion>-<n> or held-out-<criterion>-<variant>-<n>")
	rule := fs.String("rule", "engine", "the rule that chooses: engine, or ceiling")
	out := fs.String("out", "", "the file of results; stdout when empty")
	commits := commitFlag(fs)
	part := partFlag(fs)
	if err := parse(fs, args); err != nil {
		return err
	}
	c := bench.Config{Learners: *learners, Commits: commits, Part: *part}
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
	alone(*part)
	return writeTo(*out, stdout, func(w io.Writer) error { return bench.Run(c, w) })
}

// commitFlag adds the flag -commit to a command: a library the build names no version of, as
// <module path>=<full SHA>, repeatable.
func commitFlag(fs *flag.FlagSet) map[string]string {
	commits := map[string]string{}
	fs.Func("commit", "a library the build names no version of, as <module path>=<full SHA>; repeatable", func(v string) error {
		path, sha, ok := strings.Cut(v, "=")
		if !ok || path == "" || sha == "" {
			return fmt.Errorf("want <module path>=<full SHA>, got %q", v)
		}
		commits[path] = sha
		return nil
	})
	return commits
}

// partFlag adds the flag -part to a command: the part k of n of the learners, so that a recipe runs
// the parts as processes and merges their files (R09; the amendment of 2026-11-08 to R24).
func partFlag(fs *flag.FlagSet) *bench.Part {
	p := &bench.Part{}
	fs.Func("part", "the part k/n of the learners, those whose number gives k mod n; all of them when not given", func(v string) error {
		k, n, ok := strings.Cut(v, "/")
		var err error
		if p.K, err = strconv.Atoi(k); err == nil && ok {
			p.N, err = strconv.Atoi(n)
		}
		if err != nil || !ok || p.N < 1 || p.K < 0 || p.K >= p.N {
			return fmt.Errorf("want k/n with 0 ≤ k < n, got %q", v)
		}
		return nil
	})
	return p
}

// writeTo writes what write gives to stdout, or to a file: beside its place, and put there only once
// write is done, so that a run that fails leaves a file that was there as it was. The file keeps the
// mode of the one it replaces, or takes 0644, as git checks a file out.
func writeTo(path string, stdout io.Writer, write func(io.Writer) error) error {
	if path == "" {
		return write(stdout)
	}
	mode := os.FileMode(0o644)
	if st, err := os.Stat(path); err == nil {
		mode = st.Mode().Perm()
	}
	f, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	if err = f.Chmod(mode); err == nil {
		err = write(f)
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), path)
	}
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}

// readResults reads files of results.
func readResults(paths []string) ([]*bench.File, error) {
	var out []*bench.File
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		file, err := bench.Read(f)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		out = append(out, file)
	}
	return out, nil
}

// mergeCommand is "sim merge".
func mergeCommand(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	out := fs.String("out", "", "the file of the whole; stdout when empty")
	if err := parse(fs, args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, "usage: sim merge [-out <file>] <part>…")
		return errUsage
	}
	parts, err := readResults(fs.Args())
	if err != nil {
		return err
	}
	whole, err := bench.Merge(parts)
	if err != nil {
		return err
	}
	return writeTo(*out, stdout, func(w io.Writer) error { return bench.Write(whole, w) })
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
	files, err := readResults(fs.Args())
	if err != nil {
		return err
	}
	return bench.Compare(bench.Side{File: files[0], Oracle: *aOracle}, bench.Side{File: files[1], Oracle: *bOracle}, stdout)
}

// judgement is a judgement file: the criterion, the files of results of the base and of each
// candidate — with the graphs it claims to give the base's runs on, equal_on (R49) —, and of the
// confirmation when it has run; a path is relative to the judgement file. The paths of the score
// are those the criterion holds fixed, so no ceiling runs beside them (R28).
type judgement struct {
	Criterion  string `json:"criterion"`
	Base       string `json:"base"`
	Candidates []struct {
		Name      string   `json:"name"`
		File      string   `json:"file"`
		Group     string   `json:"group"`
		Goal      string   `json:"goal"`
		Constants int      `json:"constants"`
		NewField  bool     `json:"new_field"`
		Distance  int      `json:"distance"`
		Library   string   `json:"library"`
		EqualOn   []string `json:"equal_on"`
	} `json:"candidates"`
	Confirm *struct {
		Base      string `json:"base"`
		Candidate string `json:"candidate"`
	} `json:"confirm"`
}

// gateJSON is the gate file the bench is built with: the judge holds the base to its libraries, and
// a candidate of one library to the other of them, as it holds the score to the paths it embeds
// (R40 in docs/decisions.md).
//
//go:embed gate.json
var gateJSON []byte

// pins gives the libraries of a gate file as the points the criterion holds the base to.
func pins(data []byte) (map[string]criterion.Pin, error) {
	f, err := gate.Read(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	out := map[string]criterion.Pin{}
	for _, path := range slices.Sorted(maps.Keys(f.Libraries)) {
		out[path] = criterion.Pin{Commit: f.Libraries[path].Commit, Version: f.Libraries[path].Version}
	}
	return out, nil
}

// judgeCommand is "sim judge".
func judgeCommand(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("judge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := parse(fs, args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: sim judge <judgement file>")
		return errUsage
	}
	raw, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	var j judgement
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&j); err != nil {
		return fmt.Errorf("%s: %w", fs.Arg(0), err)
	}
	dir := filepath.Dir(fs.Arg(0))
	read := func(path string) (*bench.File, error) {
		if path == "" {
			return nil, nil
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(dir, path)
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		file, err := bench.Read(f)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		return file, nil
	}
	in := criterion.Input{Criterion: j.Criterion}
	if in.Pins, err = pins(gateJSON); err != nil {
		return fmt.Errorf("the embedded gate file: %w", err)
	}
	if in.Base, err = read(j.Base); err != nil {
		return err
	}
	for _, c := range j.Candidates {
		f, err := read(c.File)
		if err != nil {
			return err
		}
		in.Candidates = append(in.Candidates, criterion.Candidate{Name: c.Name, File: f, Group: c.Group, Goal: c.Goal,
			Constants: c.Constants, NewField: c.NewField, Distance: c.Distance, Library: c.Library, EqualOn: c.EqualOn})
	}
	if c := j.Confirm; c != nil {
		in.Confirm = &criterion.Confirm{}
		for _, x := range []struct {
			path string
			to   **bench.File
		}{{c.Base, &in.Confirm.Base}, {c.Candidate, &in.Confirm.Candidate}} {
			if *x.to, err = read(x.path); err != nil {
				return err
			}
		}
	}
	r, err := criterion.Judge(in)
	if err != nil {
		return err
	}
	return r.Write(stdout)
}

// pathsCommand is "sim paths": it reads the paths of the score off the decisive runs of the base
// and of the ceiling of the engine, writes them into the file the criterion holds, and writes the
// rows of their table, as the latest subsection of section 3 of docs/measurements.md holds them (R28).
func pathsCommand(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("paths", flag.ContinueOnError)
	fs.SetOutput(stderr)
	base := fs.String("base", "", "the file of results of the decisive run of the base")
	ceiling := fs.String("ceiling", "", "the file of results of the decisive run of the ceiling of the engine")
	out := fs.String("out", "", "the file of the paths, internal/criterion/paths.json of the repository")
	table := fs.Bool("table", false, "write the rows of the table of the paths to stdout")
	if err := parse(fs, args); err != nil {
		return err
	}
	if *base == "" || *ceiling == "" || *out == "" || fs.NArg() != 0 {
		fmt.Fprintln(stderr, "usage: sim paths -base <file> -ceiling <file> -out <file> [-table]")
		return errUsage
	}
	b, bd, err := readDigest(*base)
	if err != nil {
		return err
	}
	c, cd, err := readDigest(*ceiling)
	if err != nil {
		return err
	}
	p, err := criterion.RecordPaths(b, c, bd, cd)
	if err != nil {
		return err
	}
	if err := writeTo(*out, stdout, p.Write); err != nil {
		return err
	}
	if *table {
		_, err = io.WriteString(stdout, p.Table())
	}
	return err
}

// readDigest reads a file of results and gives its SHA-256 beside it, in hex, read in the same pass.
func readDigest(path string) (*bench.File, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	h := sha256.New()
	file, err := bench.Read(io.TeeReader(f, h))
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", path, err)
	}
	return file, hex.EncodeToString(h.Sum(nil)), nil
}
