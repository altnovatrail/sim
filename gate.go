package main

import (
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

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/gate"
	"github.com/altnovatrail/spec/types/formatv0"
)

// errRed is a gate that ran and is red: its report says why.
var errRed = errors.New("the gate is red")

// buildInfo gives what the build says of itself; a test puts a build of the workspace in its place,
// since go test builds from the module cache.
var buildInfo = debug.ReadBuildInfo

// gateCommand is "sim gate": the recipe of the gate calls it, built once through the workspace it
// writes from the gate file (docs/architecture.md, section 3.9; R26).
func gateCommand(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, gateUsage)
		return errUsage
	}
	switch args[0] {
	case "libraries":
		return gateLibraries(args[1:], stdout, stderr)
	case "run":
		return gateRun(args[1:], stdout, stderr)
	case "check", "record", "thresholds":
		return gateCheck(args[0], args[1:], stdout, stderr)
	case "changed":
		return gateChanged(args[1:], stdout, stderr)
	}
	fmt.Fprintf(stderr, "sim gate: unknown command %q\n%s\n", args[0], gateUsage)
	return errUsage
}

const gateUsage = `usage: sim gate libraries <gate file>
       sim gate run -part k/n -commit <path>=<sha>… [-package <folder>…] [-out <file>]
       sim gate check [-exact] -file <gate file> -go <module dir>… <part>…
       sim gate record -reason <why> -file <gate file> -go <module dir>… <part>…
       sim gate changed (-index <module-index.json> | -first) [-published <folder>]… [-all] <package folder>…
       sim gate thresholds -file <gate file> -go <module dir>… <part>…`

// gateLibraries writes the libraries of the gate file, one a line: the module path, the repository,
// the commit and the tag of its release, if it has one, for the recipe to fetch and hold the tag to
// the commit.
func gateLibraries(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("gate libraries", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := parse(fs, args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, gateUsage)
		return errUsage
	}
	f, err := readGate(fs.Arg(0))
	if err != nil {
		return err
	}
	for _, path := range slices.Sorted(maps.Keys(f.Libraries)) {
		l := f.Libraries[path]
		line := path + " " + l.Repository + " " + l.Commit
		if l.Version != "" {
			line += " " + l.Version
		}
		if _, err := fmt.Fprintln(stdout, line); err != nil {
			return err
		}
	}
	return nil
}

// gateRun runs a part of the gate; with packages, a part of the gate of modules on them.
func gateRun(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("gate run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	commits := commitFlag(fs)
	part := partFlag(fs)
	out := fs.String("out", "", "the file of results of the part; stdout when empty")
	folders := listFlag(fs, "package", "the folder of a package of a pull request of modules, to run in place of the pinned one; repeatable")
	along := listFlag(fs, "along", "the folder of a package of the pull request that did not change and came along with one that did, run only on the graphs it shares with it; repeatable")
	if err := parse(fs, args); err != nil {
		return err
	}
	c := gate.Config(*part, commits)
	if len(*folders) > 0 {
		pkgs, err := readFolders(*folders)
		if err != nil {
			return err
		}
		taken, err := gate.ReadPackages(pkgs)
		if err != nil {
			return err
		}
		pkgs, err = readFolders(*along)
		if err != nil {
			return err
		}
		came, err := gate.ReadPackages(pkgs)
		if err != nil {
			return err
		}
		for _, p := range came {
			p.Along = true
			taken = append(taken, p)
		}
		c = gate.ModulesConfig(*part, commits, taken)
	}
	return writeTo(*out, stdout, func(w io.Writer) error { return bench.Run(c, w) })
}

// listFlag adds a repeatable flag of strings to a command.
func listFlag(fs *flag.FlagSet, name, usage string) *[]string {
	var out []string
	fs.Func(name, usage, func(v string) error {
		out = append(out, v)
		return nil
	})
	return &out
}

// readFolders reads the packages in folders, as module-kit build writes them and as the CI of
// modules unpacks the published ones: inventory.json and server.json, and nothing else read.
func readFolders(folders []string) ([]gate.Package, error) {
	var out []gate.Package
	for _, d := range folders {
		inventory, err := os.ReadFile(filepath.Join(d, "inventory.json"))
		if err != nil {
			return nil, err
		}
		server, err := os.ReadFile(filepath.Join(d, "server.json"))
		if err != nil {
			return nil, err
		}
		out = append(out, gate.Package{Folder: d, Inventory: inventory, Server: server})
	}
	return out, nil
}

// gateChanged tells whether what the bench reads of each package of a pull request of modules
// changed against the newest published version of its module: a line a package on stderr, and on
// stdout the folders the gate of modules runs on — the changed and new ones, or with -all, when the
// pull request moves the tag of sim, every one but a transition of the format —, a line each, and
// after "along " those that did not change but share a graph of the gate with one of them
// (gate.Runs; docs/architecture.md, section 3.9; R27, R31).
func gateChanged(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("gate changed", flag.ContinueOnError)
	fs.SetOutput(stderr)
	index := fs.String("index", "", "module-index.json of the newest release of modules, which release-verify has accepted")
	first := fs.Bool("first", false, "modules has no release yet, and every package is new")
	all := fs.Bool("all", false, "run every package but a transition of the format: the pull request moves the tag of sim")
	pubs := listFlag(fs, "published", "the folder of a published package of the index; repeatable")
	if err := parse(fs, args); err != nil {
		return err
	}
	if fs.NArg() == 0 || (*index == "") == !*first || *first && len(*pubs) > 0 {
		fmt.Fprintln(stderr, gateUsage)
		return errUsage
	}
	var idx *formatv0.ModuleIndex
	if *index != "" {
		data, err := os.ReadFile(*index)
		if err != nil {
			return err
		}
		idx = &formatv0.ModuleIndex{}
		if err := json.Unmarshal(data, idx); err != nil {
			return fmt.Errorf("%s: %w", *index, err)
		}
	}
	pub, err := readFolders(*pubs)
	if err != nil {
		return err
	}
	pkgs, err := readFolders(fs.Args())
	if err != nil {
		return err
	}
	changes, err := gate.Changes(idx, pub, pkgs)
	if err != nil {
		return err
	}
	runs := gate.Runs(changes, *all)
	for _, c := range changes {
		line := fmt.Sprintf("%s %s: %s", c.Module, c.Version, c.State)
		if c.Against != "" {
			line += " against " + c.Against
		}
		if i := slices.IndexFunc(runs, func(r gate.Change) bool { return r.Module == c.Module }); i >= 0 && runs[i].With != "" {
			line += ", run with " + runs[i].With + ", which shares a graph of the gate with it"
		}
		fmt.Fprintln(stderr, line)
	}
	for _, c := range runs {
		line := c.Folder
		if c.With != "" {
			line = "along " + line // run only on the graphs it shares with a changed one
		}
		if _, err := fmt.Fprintln(stdout, line); err != nil {
			return err
		}
	}
	return nil
}

// gateCheck checks the run of the gate, its parts merged, against the gate file and writes the
// report; or records its bands into the gate file; or checks the run of the gate of modules against
// its thresholds and writes that report. Each first holds the build to the file.
func gateCheck(cmd string, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("gate "+cmd, flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "gate.json", "the gate file")
	exact := fs.Bool("exact", false, "hold the numbers to the recorded ones to the last bit, not to their bands")
	reason := fs.String("reason", "", "why the bands are recorded anew")
	var dirs []string
	fs.Func("go", "the folder of a module of the workspace, whose go.mod the gate reads; repeatable", func(v string) error {
		dirs = append(dirs, v)
		return nil
	})
	if err := parse(fs, args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, gateUsage)
		return errUsage
	}
	f, err := readGate(*file)
	if err != nil {
		return err
	}
	var mods []gate.Module
	for _, d := range dirs {
		data, err := os.ReadFile(filepath.Join(d, "go.mod"))
		if err != nil {
			return err
		}
		m, err := gate.ParseGoMod(data)
		if err != nil {
			return fmt.Errorf("%s: %w", d, err)
		}
		mods = append(mods, m)
	}
	info, _ := buildInfo()
	reasons := gate.Verify(f, info, runtime.Version(), mods)
	parts, err := readResults(fs.Args())
	if err != nil {
		return err
	}
	results, err := bench.Merge(parts)
	if err != nil {
		return err
	}
	switch cmd {
	case "record":
		if len(reasons) > 0 {
			return fmt.Errorf("the build departs from the gate file: %v", reasons)
		}
		recorded, err := gate.Record(f, results, *reason)
		if err != nil {
			return err
		}
		return writeTo(*file, stdout, recorded.Write)
	case "thresholds":
		r, err := gate.Thresholds(f, results)
		if err != nil {
			return err
		}
		r.Reasons = append(reasons, r.Reasons...)
		if err := r.Write(stdout); err != nil {
			return err
		}
		if !r.OK() {
			return errRed
		}
		return nil
	}
	r, err := gate.Check(f, results, *exact)
	if err != nil {
		return err
	}
	r.Reasons = append(reasons, r.Reasons...)
	if err := r.Write(stdout); err != nil {
		return err
	}
	if !r.OK() {
		return errRed
	}
	return nil
}

// readGate reads a gate file.
func readGate(path string) (*gate.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	g, err := gate.Read(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return g, nil
}
