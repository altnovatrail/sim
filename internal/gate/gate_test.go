package gate

import (
	"bytes"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/seed"
)

// results builds a run of the gate: every cell with all its learners, each learner's sums for every
// number of the gate, the k-th a little apart; drop leaves a key out of one cell.
func results(shift float64, dropCell, dropKey string) *bench.File {
	f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Invariants: true, Libraries: map[string]string{
		"github.com/altnovatrail/engine": strings.Repeat("c", 40), "github.com/altnovatrail/learner-model": strings.Repeat("d", 40),
	}}}
	for _, cell := range Cells() {
		pop, graph, _ := strings.Cut(cell, "/")
		for k := range learners {
			e := float64(0.001*float64(k%7)) + shift
			s := measure.Sums{"learners": 1,
				"M1.200.sq": float64((0.3 + e) * (0.3 + e)), "M1.200.n": 1,
				"E1.corridor": float64(10 + k%3), "E1.items": 100,
				"M4.all.false": float64(k % 5 / 4), "M4.all.declared": 1,
				"M4.pages.false": 0, "M4.pages.declared": 1,
				"M6.late.sum": float64(4 + k%3), "M6.late.n": 1, "M6.never": float64(k % 9 / 8), "M6.cases": 1,
				"M7.10.sq": float64((0.5 + e) * (0.5 + e)), "M7.10.n": 1,
				"M3.lag.sum": e, "M3.lag.n": 1, "M3.unsettled": float64(k % 4 / 3), "M3.jumped": 1,
				"M2.after7.brier": 0.1 + e, "M2.after7.n": 1,
				"M5.caught.sum": float64(3 + k%2), "M5.caught.n": 1, "M5.never": 0, "M5.episodes": 1,
				"E3.late": float64(k % 2), "E3.reviews": 1,
				"M9.hasFalse": 0, "M9.has": 1, "M9.notFalse": 0, "M9.not": 1,
				"E5.written": 4, "E5.again": 0, "E2.stuck": 0,
			}
			if cell == dropCell {
				delete(s, dropKey)
			}
			f.Lines = append(f.Lines, bench.Line{Population: pop, Graph: strings.Split(graph, "+"), N: k, Sums: s})
		}
	}
	return f
}

func file(t *testing.T) *File {
	t.Helper()
	prior := learners
	learners = 30
	t.Cleanup(func() { learners = prior })
	ps, err := published.Pinned()
	if err != nil {
		t.Fatal(err)
	}
	return &File{
		Libraries: map[string]Library{
			"github.com/altnovatrail/engine":        {Repository: "altnovatrail/engine", Commit: strings.Repeat("c", 40), Version: "v0.1.0"},
			"github.com/altnovatrail/learner-model": {Repository: "altnovatrail/learner-model", Commit: strings.Repeat("d", 40)},
		},
		Modules:      Modules{Release: published.Release, Index: "sha256:" + published.IndexDigest, Packages: ps},
		Dependencies: map[string]string{"github.com/altnovatrail/spec": "v0.9.9"},
		Go:           map[string]string{"github.com/altnovatrail/sim": "1.27.3", "github.com/altnovatrail/engine": "1.27.3", "github.com/altnovatrail/learner-model": "1.27.3"},
	}
}

const g0 = "G0/novatrail/olympiad-math"

func band(f *File, cell, number string) *Band {
	i := slices.IndexFunc(f.Bands.Numbers, func(b Band) bool { return b.Cell == cell && b.Number == number })
	return &f.Bands.Numbers[i]
}

// A run checked against its own record is green to the last bit; a number moved within two
// half-widths is in its band but not the recorded one; moved further, it is out, and the report
// names it with its band; a number of no graph is never red.
func TestCheck(t *testing.T) {
	rec, err := Record(file(t), results(0, "", ""), "a test")
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Bands.Numbers) == 0 || rec.Bands.Reason != "a test" {
		t.Fatalf("bands %+v", rec.Bands)
	}
	for _, exact := range []bool{true, false} {
		r, err := Check(rec, results(0, "", ""), exact)
		if err != nil || !r.OK() {
			t.Errorf("exact %v: %v %v", exact, err, r.Reasons)
		}
	}
	b := band(rec, g0, "M1.rms.200")
	b.Value += b.Half // within two half-widths
	if r, _ := Check(rec, results(0, "", ""), false); !r.OK() {
		t.Errorf("within the band: %v", r.Reasons)
	}
	if r, _ := Check(rec, results(0, "", ""), true); r.OK() || !strings.Contains(r.Reasons[0], "to the last bit") {
		t.Errorf("held exactly, a number moved: %v", r.Reasons)
	}
	b.Value += float64(2 * b.Half) // three half-widths off
	r, _ := Check(rec, results(0, "", ""), false)
	if r.OK() || !strings.Contains(strings.Join(r.Reasons, "\n"), g0+" M1.rms.200 is") {
		t.Errorf("out of the band: %v", r.Reasons)
	}
	var out bytes.Buffer
	if err := r.Write(&out); err != nil || !strings.Contains(out.String(), "The gate is red") || !strings.Contains(out.String(), "| "+g0+" | M1.rms.200 |") {
		t.Errorf("the red report:\n%s", out.String())
	}
	b.Value -= float64(3 * b.Half)
	r, _ = Check(rec, results(0, "", ""), false)
	out.Reset()
	if err := r.Write(&out); err != nil || !strings.Contains(out.String(), "The gate is green") || !strings.Contains(out.String(), "no graph") {
		t.Errorf("the green report:\n%s", out.String())
	}
}

// A number the recording could not read is held as not read: read later, it is out of its band; a
// number with no band, a breach of an invariant and a run that is not the gate's whole make it red.
func TestCheckRed(t *testing.T) {
	rec, err := Record(file(t), results(0, g0, "M6.late.n"), "a test")
	if err != nil {
		t.Fatal(err)
	}
	if band(rec, g0, "M6.late").Read {
		t.Fatal("a number no learner reads is recorded as read")
	}
	if r, _ := Check(rec, results(0, g0, "M6.late.n"), true); !r.OK() {
		t.Errorf("not read, as recorded: %v", r.Reasons)
	}
	if r, _ := Check(rec, results(0, "", ""), false); r.OK() {
		t.Error("a number recorded as not read is read, and the gate is green")
	}
	none := *rec
	none.Bands.Numbers = slices.DeleteFunc(slices.Clone(rec.Bands.Numbers), func(b Band) bool { return b.Cell == g0 && b.Number == "E2.stuck" })
	if r, _ := Check(&none, results(0, g0, "M6.late.n"), false); r.OK() || !strings.Contains(strings.Join(r.Reasons, " "), "has no band") {
		t.Errorf("no band: %v", r.Reasons)
	}
	broken := results(0, g0, "M6.late.n")
	broken.Lines[5].Breaches = []string{"I1: the runs differ"}
	if r, _ := Check(rec, broken, false); r.OK() {
		t.Error("a breach of an invariant left the gate green")
	}
	if _, err := Record(rec, broken, "a test"); err == nil {
		t.Error("bands recorded on a run with a breach")
	}
	part := results(0, "", "")
	part.Header.Part = bench.Part{K: 0, N: 2}
	short := results(0, "", "")
	short.Lines = short.Lines[1:]
	for name, f := range map[string]*bench.File{"a part": part, "a learner short": short} {
		if r, _ := Check(rec, f, false); r.OK() {
			t.Errorf("%s: green", name)
		}
		if _, err := Record(rec, f, "a test"); err == nil {
			t.Errorf("%s: recorded", name)
		}
	}
	if _, err := Record(rec, results(0, "", ""), ""); err == nil {
		t.Error("bands recorded for no reason")
	}
	// Bands of a caller's code are no record; a run on a synthetic module is red.
	theirs := results(0, "", "")
	theirs.Header.Libraries["github.com/altnovatrail/engine"] = strings.Repeat("e", 40)
	if _, err := Record(rec, theirs, "a test"); err == nil || !strings.Contains(err.Error(), "the run took") {
		t.Errorf("bands recorded on a caller's code: %v", err)
	}
	synth := results(0, "", "")
	synth.Lines[0].Graph = []string{"synthetic/sheet"}
	if r, _ := Check(rec, synth, false); r.OK() || !strings.Contains(strings.Join(r.Reasons, " "), "synthetic module") {
		t.Errorf("a synthetic module: %v", r.Reasons)
	}
}

// A gate file reads back as it was written; a field it does not know is refused.
func TestReadWrite(t *testing.T) {
	rec, err := Record(file(t), results(0, "", ""), "a test")
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := rec.Write(&buf); err != nil {
		t.Fatal(err)
	}
	back, err := Read(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	var again bytes.Buffer
	if err := back.Write(&again); err != nil || !bytes.Equal(buf.Bytes(), again.Bytes()) {
		t.Errorf("the gate file does not read back as written: %v", err)
	}
	if _, err := Read(strings.NewReader(`{"bands":{},"extra":1}`)); err == nil {
		t.Error("an unknown field is read")
	}
	if got := Numbers("G1"); !slices.Equal(got, []string{"M7.error.10", "E2.stuck"}) {
		t.Errorf("Numbers(G1) = %v", got)
	}
	if len(Cells()) != 14 {
		t.Errorf("%d cells: %v", len(Cells()), Cells())
	}
	c := Config(bench.Part{K: 1, N: 4}, nil)
	if c.Learners != Learners || !c.Repeat || c.Set != seed.Working || c.Part.N != 4 {
		t.Errorf("Config = %+v", c)
	}
}

// The build is held to the gate file: the libraries from the workspace, the versions the workspace
// chose, the go lines, the point sim's go.mod names, and the packages.
func TestVerify(t *testing.T) {
	f := file(t)
	info := &debug.BuildInfo{Deps: []*debug.Module{
		{Path: "github.com/altnovatrail/engine", Version: "(devel)"},
		{Path: "github.com/altnovatrail/learner-model", Version: "(devel)"},
		{Path: "github.com/altnovatrail/spec", Version: "v0.9.8", Replace: &debug.Module{Path: "github.com/altnovatrail/spec", Version: "v0.9.9"}},
	}}
	mods := func(simEngine, simLM, goLine string) []Module {
		return []Module{
			{Path: "github.com/altnovatrail/sim", Go: goLine, Requires: map[string]string{"github.com/altnovatrail/engine": simEngine, "github.com/altnovatrail/learner-model": simLM}},
			{Path: "github.com/altnovatrail/engine", Go: "1.27.3"},
			{Path: "github.com/altnovatrail/learner-model", Go: "1.27.3"},
		}
	}
	pseudo := "v0.1.1-0.20261108000000-" + strings.Repeat("d", 12)
	if got := Verify(f, info, "go1.27.3", mods("v0.1.0", pseudo, "1.27.3")); len(got) != 0 {
		t.Errorf("a build as the file names it: %v", got)
	}
	built := *info
	built.Deps = []*debug.Module{{Path: "github.com/altnovatrail/engine", Version: "v0.1.0"}, info.Deps[1], {Path: "github.com/altnovatrail/spec", Version: "v0.9.8"}}
	for name, c := range map[string]struct {
		info *debug.BuildInfo
		goV  string
		mods []Module
		want string
	}{
		"a library from the cache": {&built, "go1.27.3", mods("v0.1.0", pseudo, "1.27.3"), "not from the workspace"},
		"spec at another version":  {&built, "go1.27.3", mods("v0.1.0", pseudo, "1.27.3"), "the workspace chose github.com/altnovatrail/spec at v0.9.8"},
		"no spec":                  {nil, "go1.27.3", mods("v0.1.0", pseudo, "1.27.3"), "at nothing"},
		"another go line":          {info, "go1.28.0", mods("v0.1.0", pseudo, "1.28.0"), "has the go line 1.28.0, and the gate file records"},
		"a go newer than the gate": {info, "go1.27.2", mods("v0.1.0", pseudo, "1.27.3"), "newer than the Go of the gate"},
		"sim at another tag":       {info, "go1.27.3", mods("v0.0.9", pseudo, "1.27.3"), "requires github.com/altnovatrail/engine at \"v0.0.9\""},
		"sim at another commit":    {info, "go1.27.3", mods("v0.1.0", "v0.1.1-0.20261108000000-aaaaaaaaaaaa", "1.27.3"), "the pseudo-version of commit"},
		"no go.mod of a library":   {info, "go1.27.3", mods("v0.1.0", pseudo, "1.27.3")[:2], "no go.mod of github.com/altnovatrail/learner-model"},
		"no go.mod of sim":         {info, "go1.27.3", mods("v0.1.0", pseudo, "1.27.3")[1:], "no go.mod of github.com/altnovatrail/sim"},
	} {
		if got := strings.Join(Verify(f, c.info, c.goV, c.mods), "\n"); !strings.Contains(got, c.want) {
			t.Errorf("%s: %q, want %q", name, got, c.want)
		}
	}
	moved := *f
	moved.Modules.Release = "v0.2.0"
	moved.Modules.Packages = moved.Modules.Packages[1:]
	if got := strings.Join(Verify(&moved, info, "go1.27.3", mods("v0.1.0", pseudo, "1.27.3")), "\n"); !strings.Contains(got, "v0.2.0") || !strings.Contains(got, "the packages") {
		t.Errorf("other packages: %q", got)
	}
}

// ParseGoMod reads the module path, the go line and both forms of require, comments aside.
func TestParseGoMod(t *testing.T) {
	m, err := ParseGoMod([]byte("module example.com/x // the module\n\ngo 1.27.3\n\nrequire a.com/b v1.0.0\n\nrequire (\n\tc.com/d v0.1.0 // indirect\n\n\te.com/f v2.0.0\n)\n"))
	if err != nil || m.Path != "example.com/x" || m.Go != "1.27.3" || m.Requires["a.com/b"] != "v1.0.0" || m.Requires["c.com/d"] != "v0.1.0" || m.Requires["e.com/f"] != "v2.0.0" {
		t.Errorf("ParseGoMod = %+v, %v", m, err)
	}
	if _, err := ParseGoMod([]byte("go 1.27.3\n")); err == nil {
		t.Error("a go.mod with no module path is read")
	}
}
