package gate

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/spec/types/formatv0"
)

// release gives the index of the pinned release of modules and its packages as folders hold them,
// each named by its package digest, as the CI of modules unpacks them.
func release(t *testing.T) (*formatv0.ModuleIndex, []Package) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "published", published.Release, "module-index.json"))
	if err != nil {
		t.Fatal(err)
	}
	idx, err := pack.ReadIndex(data, published.IndexDigest)
	if err != nil {
		t.Fatal(err)
	}
	pins, err := published.Pinned()
	if err != nil {
		t.Fatal(err)
	}
	var out []Package
	for _, p := range pins {
		name := strings.Replace(p.Digest, ":", "-", 1)
		archive, err := published.Archive(p)
		if err != nil {
			t.Fatal(err)
		}
		files, err := pack.Untar(archive)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, Package{Folder: name, Inventory: files["inventory.json"], Server: files["server.json"]})
	}
	return idx, out
}

// of gives the package of a module among packages.
func of(t *testing.T, packages []Package, id string) Package {
	t.Helper()
	for _, p := range packages {
		var inv formatv0.Inventory
		if err := json.Unmarshal(p.Inventory, &inv); err != nil {
			t.Fatal(err)
		}
		if string(inv.ID) == id {
			return p
		}
	}
	t.Fatalf("no package of %s", id)
	return Package{}
}

// edited gives a package with its canonical view and its inventory changed, the inventory naming
// the server part anew, as module-kit build would write the package of a pull request.
func edited(t *testing.T, p Package, folder string, view func(*formatv0.CanonicalView), inventory func(*formatv0.Inventory)) Package {
	t.Helper()
	var cv formatv0.CanonicalView
	var inv formatv0.Inventory
	if err := json.Unmarshal(p.Server, &cv); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(p.Inventory, &inv); err != nil {
		t.Fatal(err)
	}
	if view != nil {
		view(&cv)
	}
	server, err := json.Marshal(cv)
	if err != nil {
		t.Fatal(err)
	}
	inv.ServerPart = formatv0.Descriptor{Digest: pack.Digest(server), Size: int64(len(server))}
	if inventory != nil {
		inventory(&inv)
	}
	data, err := json.Marshal(inv)
	if err != nil {
		t.Fatal(err)
	}
	return Package{Folder: folder, Inventory: data, Server: server}
}

const olympiad = "novatrail/olympiad-math"

// A task tagged to another skill, the graph untouched, is a change; a text for the learner is not;
// a module with no published version is new; another version of the format than the published is
// a transition; a withdrawn version is not the one compared.
func TestChanges(t *testing.T) {
	idx, pub := release(t)
	same, err := Changes(idx, pub, pub)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range same {
		if c.State != Same || c.Against != c.Version || c.Read == nil || c.Read.Digest != string(pack.Digest(of(t, pub, c.Module).Inventory)) {
			t.Errorf("a published package against itself: %+v", c)
		}
	}
	o := of(t, pub, olympiad)
	var task formatv0.TaskID
	retag := edited(t, o, "retag", func(cv *formatv0.CanonicalView) {
		task = slices.Sorted(func(yield func(formatv0.TaskID) bool) {
			for id := range cv.Tasks {
				if !yield(id) {
					return
				}
			}
		})[0]
		st := cv.Tasks[task]
		var other formatv0.SkillID
		for id := range cv.SkillGraph.Skills {
			if id != st.Skills.Primary[0] && (other == "" || id < other) {
				other = id
			}
		}
		st.Skills.Primary = []formatv0.SkillID{other}
		cv.Tasks[task] = st
		cv.Manifest.Version = "0.2.1"
	}, func(inv *formatv0.Inventory) { inv.Version = "0.2.1" })
	text := edited(t, o, "text", func(cv *formatv0.CanonicalView) {
		st := cv.Tasks[task]
		st.Statement = []formatv0.StatementPart{{Text: new(formatv0.Text("Another wording of the task."))}}
		cv.Tasks[task] = st
		cv.Manifest.Version = "0.2.1"
	}, func(inv *formatv0.Inventory) { inv.Version = "0.2.1" })
	for _, c := range []struct {
		p    Package
		want State
	}{{retag, Changed}, {text, Same}} {
		got, err := Changes(idx, pub, []Package{c.p})
		if err != nil || got[0].State != c.want || got[0].Version != "0.2.1" || got[0].Against != "0.2.0" {
			t.Errorf("%s: %+v, %v; want %s", c.p.Folder, got, err, c.want)
		}
	}
	if got, err := Changes(nil, nil, []Package{retag}); err != nil || got[0].State != New || got[0].Against != "" {
		t.Errorf("with no release: %+v, %v", got, err)
	}
	newer := edited(t, o, "newer", nil, func(inv *formatv0.Inventory) { inv.ModuleFormatVersion = "0.5" })
	if got, err := Changes(idx, pub, []Package{newer}); err != nil || got[0].State != Transition || got[0].Read != nil {
		t.Errorf("a transition of the format: %+v, %v", got, err)
	}
	if _, err := Changes(nil, nil, []Package{newer}); !errors.Is(err, pack.ErrFormat) {
		t.Errorf("a version of the format the bench does not read, with no transition: %v", err)
	}
	brokenNewer := newer
	brokenNewer.Server = append(slices.Clone(newer.Server), ' ')
	if _, err := Changes(idx, pub, []Package{brokenNewer}); !errors.Is(err, pack.ErrDigest) {
		t.Errorf("a broken package of another version of the format passes for a transition: %v", err)
	}
	if got, err := ReadPackages(pub); err != nil || len(got) != 3 || got[0].Module == nil {
		t.Errorf("ReadPackages of the published packages: %v", err)
	}
	if _, err := ReadPackages([]Package{newer}); !errors.Is(err, pack.ErrFormat) {
		t.Errorf("ReadPackages of a version of the format the bench does not read: %v", err)
	}
	// A withdrawn newer version is passed over: the compared one is the version api gives.
	withdrawn := *idx
	withdrawn.Modules = map[formatv0.ModuleID]map[formatv0.ModuleVersion]formatv0.Entry{}
	for id, vs := range idx.Modules {
		withdrawn.Modules[id] = map[formatv0.ModuleVersion]formatv0.Entry{}
		for v, e := range vs {
			withdrawn.Modules[id][v] = e
		}
	}
	reason := formatv0.WithdrawalReason("a test")
	withdrawn.Modules[olympiad]["0.10.0"] = formatv0.Entry{ModuleFormatVersion: "0.5", PackageDigest: "sha256:" + formatv0.Digest(strings.Repeat("0", 64)), Release: "v0.2.0", Withdrawn: &reason}
	withdrawn.Modules[olympiad]["0.0.9"] = formatv0.Entry{ModuleFormatVersion: "0.5", PackageDigest: "sha256:" + formatv0.Digest(strings.Repeat("0", 64)), Release: "v0.0.1"}
	if got, err := Changes(&withdrawn, pub, []Package{o}); err != nil || got[0].State != Same || got[0].Against != "0.2.0" {
		t.Errorf("against a withdrawn newer version: %+v, %v", got, err)
	}
	bare := withdrawn
	bare.Modules = map[formatv0.ModuleID]map[formatv0.ModuleVersion]formatv0.Entry{olympiad: {"0.10.0": withdrawn.Modules[olympiad]["0.10.0"]}}
	if got, err := Changes(&bare, nil, []Package{o}); err != nil || got[0].State != New {
		t.Errorf("every version withdrawn: %+v, %v", got, err)
	}
	broken := edited(t, o, "broken", nil, nil)
	broken.Server = append(slices.Clone(broken.Server), ' ')
	untranslatable := edited(t, o, "untranslatable", func(cv *formatv0.CanonicalView) { cv.SkillGraph = nil }, nil)
	other := of(t, pub, "novatrail/grade-3-math")
	for name, c := range map[string]struct {
		pub, packages []Package
	}{
		"inventory not JSON":         {pub, []Package{{Folder: "x", Inventory: []byte("{"), Server: o.Server}}},
		"a server part apart":        {pub, []Package{broken}},
		"a module twice":             {pub, []Package{o, retag}},
		"no published package":       {nil, []Package{o}},
		"an untranslatable package":  {pub, []Package{untranslatable}},
		"a published package apart":  {[]Package{{Folder: "apart", Inventory: of(t, pub, olympiad).Inventory, Server: other.Server}}, []Package{o}},
		"the index names another ID": {[]Package{other}, []Package{o}},
	} {
		in := idx
		if name == "the index names another ID" {
			moved := withdrawn
			moved.Modules = map[formatv0.ModuleID]map[formatv0.ModuleVersion]formatv0.Entry{olympiad: {"0.2.0": idx.Modules["novatrail/grade-3-math"]["0.2.0"]}}
			in = &moved
		}
		if _, err := Changes(in, c.pub, c.packages); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if semver("0.10.0", "0.9.1") <= 0 || semver("1.0", "1.0.0") != 0 || semver("0.1.0", "0.1.1") >= 0 {
		t.Error("semver does not order versions part by part")
	}
	c := ModulesConfig(bench.Part{K: 1, N: 4}, nil, []bench.Package{*same[0].Read})
	if len(c.Packages) != 1 || c.Part.N != 4 || !c.Repeat || !slices.Equal(c.Populations, Populations) {
		t.Errorf("the config of the gate of modules: %+v", c)
	}
}

// modulesRun is a run of the gate of modules on the pinned packages, as results builds it.
func modulesRun(t *testing.T) *bench.File {
	t.Helper()
	pins, err := published.Pinned()
	if err != nil {
		t.Fatal(err)
	}
	f := results(0, "", "")
	f.Header.Packages, f.Header.Repeat, f.Header.Populations, f.Header.Learners = pins, true, Populations, Learners
	return f
}

// set puts a sum of every learner of a cell.
func set(f *bench.File, cell, key string, v float64) {
	for _, l := range f.Lines {
		if l.Population+"/"+strings.Join(l.Graph, "+") == cell {
			l.Sums[key] = v
		}
	}
}

// The gate of modules is green within its thresholds, the edges included, and red past one, on a
// number not read, on a breach and on a run that is not its own; a module no cell runs on is named,
// not red.
func TestThresholds(t *testing.T) {
	rec, err := Record(file(t), results(0, "", ""), "a test")
	if err != nil {
		t.Fatal(err)
	}
	r, err := Thresholds(rec, modulesRun(t))
	if err != nil || !r.OK() {
		t.Fatalf("on the run the bands were recorded on: %v %v", err, r.Reasons)
	}
	edge := slices.IndexFunc(r.Rows, func(h Held) bool { return h.Cell == g0 && h.Number == "M4.false.all" })
	if edge < 0 || r.Rows[edge].Now.Value.Value != 0.2 || !r.Rows[edge].OK {
		t.Errorf("M4 at its threshold 0.20 is not held: %+v", r.Rows[edge])
	}
	if r.Rows[0].Cell != g0 {
		t.Errorf("the rows do not go in the order of the cells: %s first", r.Rows[0].Cell)
	}
	var out bytes.Buffer
	if err := r.Write(&out); err != nil || !strings.Contains(out.String(), "The gate of modules is green") || !strings.Contains(out.String(), "| "+g0+" | E1.hard | at most ") || strings.Contains(out.String(), "no graph") {
		t.Errorf("the green report:\n%s", out.String())
	}
	hard := band(rec, g0, "E1.hard").Value
	for name, c := range map[string]struct {
		change func(*bench.File)
		want   string
	}{
		"E1.hard past the largest": {func(f *bench.File) { set(f, g0, "E1.hard", float64(100*(hard+0.02))) }, g0 + " E1.hard is"},
		"E1 in the corridor short": {func(f *bench.File) { set(f, "G2-half/novatrail/grade-3-math", "E1.corridor", 0) }, "E1.corridor is 0, and its threshold is at least"},
		"E2 stuck":                 {func(f *bench.File) { set(f, "G3/novatrail/olympiad-math", "E2.stuck", 1) }, "E2.stuck is 1"},
		"E5 again":                 {func(f *bench.File) { set(f, "E-base/novatrail/english-essays", "E5.again", 1) }, "E5.again is 1, and its threshold is at most 0"},
		"M4 past 20 %":             {func(f *bench.File) { set(f, g0, "M4.all.false", 1) }, "M4.false.all is 1"},
		"not read":                 {func(f *bench.File) { set(f, g0, "E1.items", 0) }, g0 + " E1.hard is not read"},
		"a breach":                 {func(f *bench.File) { f.Lines[3].Breaches = []string{"I2: a code"} }, "an invariant: I2"},
		"a part":                   {func(f *bench.File) { f.Header.Part = bench.Part{K: 0, N: 2} }, "the run is not the gate's"},
		"no repeat":                {func(f *bench.File) { f.Header.Repeat = false }, "the run is not the gate's"},
		"a learner short":          {func(f *bench.File) { f.Lines = f.Lines[1:] }, "has 29 learners"},
		"a caller's library":       {func(f *bench.File) { f.Header.Libraries["github.com/altnovatrail/engine"] = strings.Repeat("e", 40) }, "the run took github.com/altnovatrail/engine"},
		"a cell not covered":       {func(f *bench.File) { f.Header.Packages = f.Header.Packages[:1] }, "and the packages do not cover it"},
	} {
		f := modulesRun(t)
		c.change(f)
		r, err := Thresholds(rec, f)
		if err != nil || r.OK() || !strings.Contains(strings.Join(r.Reasons, "\n"), c.want) {
			t.Errorf("%s: %v %v", name, err, r.Reasons)
		}
	}
	r, _ = Thresholds(rec, func() *bench.File {
		f := modulesRun(t)
		set(f, g0, "E1.corridor", 0)
		return f
	}())
	out.Reset()
	if err := r.Write(&out); err != nil || !strings.Contains(out.String(), "The gate of modules is red") || !strings.Contains(out.String(), "| no |") {
		t.Errorf("the red report:\n%s", out.String())
	}
	// A package of a module no population runs on is named and leaves the gate green.
	f := modulesRun(t)
	f.Header.Packages = append(f.Header.Packages, published.Package{Module: "novatrail/new", Version: "0.1.0", Digest: "sha256:" + strings.Repeat("1", 64)})
	r, err = Thresholds(rec, f)
	out.Reset()
	if err != nil || !r.OK() || !slices.Equal(r.Uncovered, []string{"novatrail/new"}) || r.Write(&out) != nil || !strings.Contains(out.String(), "No cell of the gate runs on novatrail/new") {
		t.Errorf("a module no cell runs on: %v %v %v\n%s", err, r.Reasons, r.Uncovered, out.String())
	}
	if _, err := Thresholds(rec, results(0, "", "")); !errors.Is(err, ErrModules) {
		t.Errorf("a run on the pinned packages: %v", err)
	}
	none := *rec
	none.Bands.Numbers = slices.DeleteFunc(slices.Clone(rec.Bands.Numbers), func(b Band) bool { return b.Number == "E1.hard" })
	if _, err := Thresholds(&none, modulesRun(t)); err == nil || !strings.Contains(err.Error(), "no E1.hard") {
		t.Errorf("a gate file with no E1.hard: %v", err)
	}
	// The gate and its record refuse a run of the gate of modules.
	if r, _ := Check(rec, modulesRun(t), false); r.OK() || !strings.Contains(strings.Join(r.Reasons, " "), "took packages") {
		t.Errorf("the gate on a run of modules: %v", r.Reasons)
	}
	if _, err := Record(rec, modulesRun(t), "a test"); err == nil {
		t.Error("bands recorded on a run of modules")
	}
}

// The gate of modules runs on the changed and new packages, and on a package that did not change
// when a graph of the gate holds it with one that did: A-linked runs on both maths when either
// changed. With all, every package but a transition runs, and none is brought.
func TestRuns(t *testing.T) {
	const grade3, essays = "novatrail/grade-3-math", "novatrail/english-essays"
	changes := func(states ...State) []Change {
		return []Change{{Module: essays, State: states[0]}, {Module: grade3, State: states[1]}, {Module: olympiad, State: states[2]}}
	}
	show := func(cs []Change) []string {
		var out []string
		for _, c := range cs {
			out = append(out, c.Module+"<"+c.With)
		}
		return out
	}
	for _, c := range []struct {
		in   []Change
		all  bool
		want []string
	}{
		{changes(Same, Same, Changed), false, []string{grade3 + "<" + olympiad, olympiad + "<"}},
		{changes(Same, New, Same), false, []string{grade3 + "<", olympiad + "<" + grade3}},
		{changes(Changed, Same, Same), false, []string{essays + "<"}},
		{changes(Same, Same, Same), false, nil},
		{changes(Same, Transition, Changed), false, []string{olympiad + "<"}},
		{changes(Same, Same, Transition), true, []string{essays + "<", grade3 + "<"}},
	} {
		if got := show(Runs(c.in, c.all)); !slices.Equal(got, c.want) {
			t.Errorf("Runs(%v, %v) = %v, want %v", c.in, c.all, got, c.want)
		}
	}
}

// A cell of the gate on a module of the run whose graph wants a module the run has no package of is
// named, and so is not a cell on a package that only came along; a package along is no module
// uncovered.
func TestPartial(t *testing.T) {
	g3 := published.Package{Module: "novatrail/grade-3-math"}
	om := published.Package{Module: olympiad}
	if got := partial([]published.Package{g3}); !slices.Equal(got, []string{"A-linked/novatrail/olympiad-math+novatrail/grade-3-math, for want of a package of novatrail/olympiad-math"}) {
		t.Errorf("partial = %v", got)
	}
	om.Along = true
	if got := partial([]published.Package{g3, om}); got != nil {
		t.Errorf("both maths: partial = %v", got)
	}
	if got := partial([]published.Package{om}); got != nil {
		t.Errorf("a package along alone: partial = %v", got)
	}
	var b bytes.Buffer
	r := &ThresholdReport{Partial: partial([]published.Package{g3})}
	if err := r.Write(&b); err != nil || !strings.Contains(b.String(), "The cell A-linked/novatrail/olympiad-math+novatrail/grade-3-math, for want of a package of novatrail/olympiad-math, does not run.") {
		t.Errorf("the report:\n%s", b.String())
	}
}
