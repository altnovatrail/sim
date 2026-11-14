package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/criterion"
	"github.com/altnovatrail/sim/internal/gate"
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
	// The judge is held here to the libraries the paths were read on, whichever this tree builds
	// with. Mostly they are the gate file's, and then nothing below changes; but between the move of
	// a library and the decisive run that reads the paths anew they differ, and the judge refuses
	// every base (R40, R41 in docs/decisions.md), as it did from T42 to T43. Then the runs are named
	// the paths' libraries, and the embedded gate file names their versions, with no commit, since
	// it knows only its own: a candidate's other library is held to the base's version alone.
	f, err := os.Open(filepath.Join("internal", "criterion", "paths.json"))
	if err != nil {
		t.Fatal(err)
	}
	p, err := criterion.ReadPaths(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"base.jsonl", "held.jsonl"} {
		relabel(t, filepath.Join(dir, name), filepath.Join(dir, name), p.Libraries)
	}
	g, err := gate.Read(bytes.NewReader(gateJSON))
	if err != nil {
		t.Fatal(err)
	}
	for path, version := range p.Libraries {
		if lib := g.Libraries[path]; lib.Version != version {
			lib.Version, lib.Commit = version, ""
			g.Libraries[path] = lib
		}
	}
	var held bytes.Buffer
	if err := g.Write(&held); err != nil {
		t.Fatal(err)
	}
	embedded := gateJSON
	t.Cleanup(func() { gateJSON = embedded })
	gateJSON = held.Bytes()
	good := write("j.json", `{"judgement":"model-2","criterion":"model","base":"base.jsonl",
		"candidates":[{"name":"same","file":"base.jsonl","group":"still","constants":1,"equal_on":["novatrail/grade-3-math"]}],
		"confirm":{"base":"held.jsonl","candidate":"held.jsonl"}}`)
	stdout.Reset()
	if got := run([]string{"judge", good}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim judge = %d: %s", got, stderr.String())
	}
	if !strings.Contains(stdout.String(), "# The criterion of the model") || !strings.Contains(stdout.String(), "the base stays") || strings.Contains(stdout.String(), "fails neutral") {
		t.Errorf("the report:\n%s", stdout.String())
	}
	// A graph the candidate names in equal_on and the runs have no cell of fails the check "neutral"
	// (R49).
	stdout.Reset()
	absent := write("absent.json", `{"criterion":"model","base":"base.jsonl","candidates":[{"name":"same","file":"base.jsonl","equal_on":["novatrail/english-essays"]}]}`)
	if got := run([]string{"judge", absent}, &stdout, &stderr); got != 0 || !strings.Contains(stdout.String(), "fails neutral novatrail/english-essays all") {
		t.Errorf("sim judge of a graph the runs do not have = %d:\n%s", got, stdout.String())
	}
	// A candidate of the engine built beside another commit of learner-model is no candidate of this
	// judgement: the judge holds it to the base's learner-model, by the embedded gate file (R40).
	beside := map[string]string{bench.Libraries[0]: strings.Repeat("f", 40), bench.Libraries[1]: p.Libraries[bench.Libraries[1]]}
	relabel(t, filepath.Join(dir, "base.jsonl"), filepath.Join(dir, "foreign.jsonl"), beside)
	foreign := write("foreign.json", `{"criterion":"engine","base":"base.jsonl","candidates":[{"name":"beside","file":"foreign.jsonl"}]}`)
	stderr.Reset()
	if got := run([]string{"judge", foreign}, &stdout, &stderr); got != 1 || !strings.Contains(stderr.String(), "candidate beside ran on "+bench.Libraries[0]) {
		t.Errorf("sim judge of a candidate beside another learner-model = %d: %s", got, stderr.String())
	}
	// A judgement file in the folder of one judgement that names another — a copy's leftover — is
	// refused before its runs are read; in its own folder it goes on (R52).
	for _, c := range []struct {
		folder string
		status int
	}{{"model-3", 1}, {"model-2", 0}} {
		if err := os.MkdirAll(filepath.Join(dir, "judgements", c.folder), 0o700); err != nil {
			t.Fatal(err)
		}
		copied := write(filepath.Join("judgements", c.folder, "judgement.json"), `{"judgement":"model-2","criterion":"model","base":"../../base.jsonl",
			"candidates":[{"name":"same","file":"../../base.jsonl"}],"confirm":{"base":"../../held.jsonl","candidate":"../../held.jsonl"}}`)
		stderr.Reset()
		if got := run([]string{"judge", copied}, &stdout, &stderr); got != c.status || c.status == 1 && !strings.Contains(stderr.String(), "lies in the folder of model-3") {
			t.Errorf("sim judge of a file naming model-2 in judgements/%s = %d: %s", c.folder, got, stderr.String())
		}
	}
	prior := gateJSON
	t.Cleanup(func() { gateJSON = prior })
	gateJSON = []byte(`{"extra":1}`)
	stderr.Reset()
	if got := run([]string{"judge", good}, &stdout, &stderr); got != 1 || !strings.Contains(stderr.String(), "the embedded gate file") {
		t.Errorf("sim judge with a gate file that does not read = %d: %s", got, stderr.String())
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
		{[]string{"judge", write("unnamed.json", `{"criterion":"model","base":"base.jsonl","confirm":{"base":"held.jsonl","candidate":"held.jsonl"}}`)}, 1},
		{[]string{"judge", write("refused.json", `{"criterion":"taste","base":"base.jsonl"}`)}, 1},
	} {
		stderr.Reset()
		if got := run(c.args, &stdout, &stderr); got != c.status {
			t.Errorf("sim %q = %d, want %d: %s", c.args, got, c.status, stderr.String())
		}
	}
}

// The judgement files in judgements/ still read as sim judge reads them, a judgement that names
// itself names its folder, and a confirmation ran on its judgement's own held-out set, as runs.txt
// names the set of each run (R35, R52 in docs/decisions.md): model-2 on held-out-model-1. The
// records are not edited: each gives its report again on its own commit.
func TestJudgements(t *testing.T) {
	dirs, err := filepath.Glob(filepath.Join("judgements", "*"))
	if err != nil || len(dirs) == 0 {
		t.Fatalf("judgements: %v, %v", dirs, err)
	}
	confirmed := 0
	for _, dir := range dirs {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			continue
		}
		name := filepath.Base(dir)
		raw, err := os.ReadFile(filepath.Join(dir, "judgement.json"))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		var j judgement
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&j); err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if j.Judgement != "" && j.Judgement != name {
			t.Errorf("%s names the judgement %s", name, j.Judgement)
		}
		if j.Confirm == nil {
			continue
		}
		runs, err := os.ReadFile(filepath.Join(dir, "runs.txt"))
		if err != nil {
			t.Fatal(err)
		}
		sets := map[string]seed.Set{} // by the file's path from the root, from the header on the line after its SHA-256
		lines := strings.Split(string(runs), "\n")
		for i := 0; i+1 < len(lines); i++ {
			if fields := strings.Fields(lines[i]); len(fields) == 2 && !strings.HasPrefix(lines[i], "#") && isSHA256(fields[0]) {
				var h struct {
					Set seed.Set `json:"set"`
				}
				if err := json.Unmarshal([]byte(lines[i+1]), &h); err != nil {
					t.Fatalf("%s, runs.txt, %s: %v", name, fields[1], err)
				}
				sets[filepath.Clean(fields[1])] = h.Set
			}
		}
		want, err := seed.Confirmation(name)
		if err != nil {
			t.Errorf("%s has a confirmation and no held-out set: %v", name, err)
			continue
		}
		for _, file := range []string{j.Confirm.Base, j.Confirm.Candidate} {
			if got, ok := sets[filepath.Clean(filepath.Join(dir, file))]; !ok || got != want {
				t.Errorf("%s: the confirmation's %s ran on %v, and the judgement's set is %v", name, file, got, want)
			}
		}
		confirmed++
	}
	if confirmed == 0 {
		t.Error("no judgement in judgements/ has a confirmation")
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

// TestMain runs the tests of the commands with an alone that leaves their process as it was, so that
// a test that runs a part does not leave the tests after it on one processor.
func TestMain(m *testing.M) {
	tuned = alone
	alone = func(bench.Part) {}
	os.Exit(m.Run())
}

// tuned is the alone of the program, which TestMain takes out of the commands' way.
var tuned func(bench.Part)

// A process that runs a part goes on one processor, and every process of a run has the collector at
// gcPercent; a run of the whole keeps its processors. The test gives the process back what it had.
func TestAlone(t *testing.T) {
	procs := runtime.GOMAXPROCS(0)
	gc := debug.SetGCPercent(100)
	debug.SetGCPercent(gc)
	t.Cleanup(func() {
		runtime.GOMAXPROCS(procs)
		debug.SetGCPercent(gc)
	})
	tuned(bench.Part{})
	if got := runtime.GOMAXPROCS(0); got != procs {
		t.Errorf("a run of the whole: GOMAXPROCS is %d, want %d", got, procs)
	}
	if got := debug.SetGCPercent(gc); got != gcPercent {
		t.Errorf("a run of the whole: the collector's percent is %d, want %d", got, gcPercent)
	}
	tuned(bench.Part{K: 1, N: 4})
	if got := runtime.GOMAXPROCS(0); got != 1 {
		t.Errorf("a part: GOMAXPROCS is %d, want 1", got)
	}
	if got := debug.SetGCPercent(gc); got != gcPercent {
		t.Errorf("a part: the collector's percent is %d, want %d", got, gcPercent)
	}
}

// relabel writes the file of results from to the file to with the libraries of its header put to
// libraries, its lines as they are.
func relabel(t *testing.T, from, to string, libraries map[string]string) {
	t.Helper()
	raw, err := os.ReadFile(from)
	if err != nil {
		t.Fatal(err)
	}
	first, rest, _ := bytes.Cut(raw, []byte("\n"))
	var h bench.Header
	if err := json.Unmarshal(first, &h); err != nil {
		t.Fatal(err)
	}
	h.Libraries = maps.Clone(libraries)
	if first, err = json.Marshal(h); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(to, append(append(first, '\n'), rest...), 0o600); err != nil {
		t.Fatal(err)
	}
}

// isSHA256 tells whether a word is a SHA-256 in hex, as sha256sum writes it.
func isSHA256(s string) bool {
	b, err := hex.DecodeString(s)
	return err == nil && len(b) == sha256.Size
}
