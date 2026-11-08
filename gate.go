package main

import (
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
	case "check", "record":
		return gateCheck(args[0], args[1:], stdout, stderr)
	}
	fmt.Fprintf(stderr, "sim gate: unknown command %q\n%s\n", args[0], gateUsage)
	return errUsage
}

const gateUsage = `usage: sim gate libraries <gate file>
       sim gate run -part k/n -commit <path>=<sha>… [-out <file>]
       sim gate check [-exact] -file <gate file> -go <module dir>… <part>…
       sim gate record -reason <why> -file <gate file> -go <module dir>… <part>…`

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

// gateRun runs a part of the gate.
func gateRun(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("gate run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	commits := commitFlag(fs)
	part := partFlag(fs)
	out := fs.String("out", "", "the file of results of the part; stdout when empty")
	if err := parse(fs, args); err != nil {
		return err
	}
	return writeTo(*out, stdout, func(w io.Writer) error { return bench.Run(gate.Config(*part, commits), w) })
}

// gateCheck checks the run of the gate, its parts merged, against the gate file and writes the
// report; or records its bands into the gate file. Both first hold the build to the file.
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
	if cmd == "record" {
		if len(reasons) > 0 {
			return fmt.Errorf("the build departs from the gate file: %v", reasons)
		}
		recorded, err := gate.Record(f, results, *reason)
		if err != nil {
			return err
		}
		return writeTo(*file, stdout, recorded.Write)
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
