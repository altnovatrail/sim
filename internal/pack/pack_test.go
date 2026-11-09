package pack

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/spec/types/formatv0"
)

const mod = "example/synthetic"

func ptr[T any](v T) *T { return &v }

// view gives a canonical view of a small synthetic module: two bands, skills with a verdict, with
// a rubric and retired, a prerequisite from the module and an anchor skill; a task on a sheet with
// a trap and items of their own skills, a task on a workbook with a choice whose options stand in
// the check and a predicate, a choice whose options stand beside the check, a task with a rubric,
// a template with three instances, and lessons. Synthetic, as every example in NovaTrail is.
func view() *formatv0.CanonicalView {
	sheet := formatv0.Canvas{Kind: formatv0.CanvasKindHandwritingSheet}
	workbook := formatv0.Canvas{Kind: formatv0.CanvasKindWorkbook}
	return &formatv0.CanonicalView{
		Manifest: formatv0.Manifest{
			ModuleFormatVersion: "0.3", ID: mod, Version: "0.1.0",
			Audience: formatv0.Audience{AgeCategories: []formatv0.AgeCategory{"under_13", "adult"}},
			Bands:    []formatv0.Band{{ID: "low", Grades: []formatv0.Grade{1, 2}}, {ID: "high", Grades: []formatv0.Grade{3, 4}}},
		},
		SkillGraph: &formatv0.SkillGraph{
			Skills: map[formatv0.SkillID]formatv0.Skill{
				"a": {Assessability: formatv0.SkillAssessabilityVerdict, Teaches: []formatv0.SharedTaxonomy{"core.x"}},
				"b": {Assessability: formatv0.SkillAssessabilityVerdict, Prerequisites: []string{"a", "core.y"}},
				"r": {Assessability: formatv0.SkillAssessabilityRubric},
			},
			RetiredSkills: map[formatv0.SkillID]formatv0.RetiredSkill{"old": {ReplacedBy: []formatv0.SkillID{"a"}}},
		},
		SkillTiers: formatv0.SkillTiers{"a": 1, "b": 2, "r": 1},
		Tasks: map[formatv0.TaskID]formatv0.Task{
			"t-sheet": {
				ID: "t-sheet", Canvas: sheet, Band: "low", Difficulty: 3, Skills: formatv0.Skills{Primary: []formatv0.SkillID{"a"}},
				Items: []formatv0.Item{
					{AnswerSlots: []formatv0.AnswerSlot{{Check: formatv0.Check{Kind: formatv0.CheckKindNumber, TrapValues: []formatv0.TrapValue{{Value: "1", Trap: "off-by-one"}, {Value: "2", Trap: "off-by-one"}}}}}},
					{Skills: &formatv0.Skills{Primary: []formatv0.SkillID{"b"}, Supporting: []string{"core.x"}},
						AnswerSlots: []formatv0.AnswerSlot{{Check: formatv0.Check{Kind: formatv0.CheckKindNumber}}}},
				},
			},
			"t-workbook": {
				ID: "t-workbook", Canvas: workbook, Band: "high", Difficulty: 1, AgeCategories: []formatv0.AgeCategory{"under_13"},
				Skills: formatv0.Skills{Primary: []formatv0.SkillID{"b"}},
				Items: []formatv0.Item{{Line: []formatv0.LinePart{
					{Text: ptr(formatv0.Text("3"))},
					{Check: &formatv0.WrittenCheck{Kind: formatv0.CheckKindChoice, Options: []formatv0.Text{"<", ">", "="}}},
					{Check: &formatv0.WrittenCheck{Kind: formatv0.CheckKindPredicate}},
				}}},
			},
			"t-beside": {
				ID: "t-beside", Canvas: sheet, Band: "low", Difficulty: 2, Skills: formatv0.Skills{Primary: []formatv0.SkillID{"a"}},
				Items: []formatv0.Item{{AnswerSlots: []formatv0.AnswerSlot{{Options: []formatv0.Text{"yes", "no"}, Check: formatv0.Check{Kind: formatv0.CheckKindChoice}}}}},
			},
			"t-rubric": {
				ID: "t-rubric", Canvas: sheet, Band: "low", Difficulty: 3, Skills: formatv0.Skills{Primary: []formatv0.SkillID{"r"}},
				Items: []formatv0.Item{{AnswerSlots: []formatv0.AnswerSlot{{Check: formatv0.Check{Kind: formatv0.CheckKindRubric}}}}},
			},
			"t-template": {
				ID: "t-template", Canvas: sheet, Band: "high", Difficulty: 4, Skills: formatv0.Skills{Primary: []formatv0.SkillID{"a"}},
				Template: &formatv0.Template{Parameters: map[formatv0.TemplateName]formatv0.Domain{"n": {From: 1, To: 3}}},
				Items:    []formatv0.Item{{AnswerSlots: []formatv0.AnswerSlot{{Check: formatv0.Check{Kind: formatv0.CheckKindNumber}}}}},
			},
		},
		Instances: map[formatv0.TaskID][]formatv0.Numbers{"t-template": {{"n": 1}, {"n": 2}, {"n": 3}}},
		Lessons: map[formatv0.LessonID]formatv0.Lesson{
			"l-child": {ID: "l-child", Skills: []formatv0.SkillID{"b"}, Segments: []formatv0.Segment{
				{Explanation: []formatv0.StatementPart{{}}}, {Task: ptr(formatv0.TaskID("t-workbook"))},
			}},
			"l-adult": {ID: "l-adult", Skills: []formatv0.SkillID{"a"}, AgeCategories: []formatv0.AgeCategory{"adult"},
				Segments: []formatv0.Segment{{Task: ptr(formatv0.TaskID("t-sheet"))}}},
		},
	}
}

func TestTranslate(t *testing.T) {
	m, err := Translate(view())
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Graph.Bands) != 2 || m.Graph.Bands[1].Place != 2 || !slices.Equal(m.Graph.Bands[1].Grades, []int32{3, 4}) {
		t.Errorf("bands %v", m.Graph.Bands)
	}
	var ids []string
	for _, s := range m.Graph.Skills {
		ids = append(ids, s.ID+"@"+strconv.Itoa(int(s.Tier)))
	}
	if want := []string{mod + "/skill/a@1", mod + "/skill/b@2", mod + "/skill/r@1"}; !slices.Equal(ids, want) {
		t.Errorf("skills %v, want %v", ids, want)
	}
	if b := m.Graph.Skills[1]; !slices.Equal(b.Prerequisites, []string{"core.y", mod + "/skill/a"}) {
		t.Errorf("the prerequisites of b are %v: the module's by full ID, the anchor as it is", b.Prerequisites)
	}
	if !slices.Equal(m.Graph.Skills[0].Teaches, []string{"core.x"}) || len(m.Graph.Retired) != 1 || m.Graph.Retired[0].ReplacedBy[0] != mod+"/skill/a" {
		t.Errorf("teaches %v, retired %v", m.Graph.Skills[0].Teaches, m.Graph.Retired)
	}

	sheet, err := m.Task(mod + "/task/t-sheet")
	if err != nil {
		t.Fatal(err)
	}
	st := sheet.Candidate.Task
	if st.Band != (learnermodel.BandRef{Place: 1, ID: "low"}) || st.Difficulty != 3 || st.KeyProvenance != learnermodel.ProvenanceAuthor || st.Template {
		t.Errorf("t-sheet: %+v", st)
	}
	if len(sheet.Candidate.Traps) != 1 || !slices.Equal(sheet.Candidate.Traps[0].Traps, []string{mod + "/trap/off-by-one"}) {
		t.Errorf("t-sheet traps %v: the item's, by full ID, once", sheet.Candidate.Traps)
	}
	if it := st.Items[1]; !it.OwnSkills || !slices.Equal(it.Primary, []string{mod + "/skill/b"}) || !slices.Equal(it.Supporting, []string{"core.x"}) {
		t.Errorf("t-sheet item 2: %+v", it)
	}

	wb, _ := m.Task(mod + "/task/t-workbook")
	slots := wb.Candidate.Task.Items[0].Slots
	want := []learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 3, HasOptions: true}, {Check: learnermodel.CheckKindPredicate}}
	if !slices.Equal(slots, want) || wb.Canvas != formatv0.CanvasKindWorkbook {
		t.Errorf("t-workbook: slots %v, canvas %q", slots, wb.Canvas)
	}
	beside, _ := m.Task(mod + "/task/t-beside")
	if s := beside.Candidate.Task.Items[0].Slots[0]; s.Options != 2 || !s.HasOptions {
		t.Errorf("t-beside: a choice whose options stand beside the check has %d options", s.Options)
	}
	rubric, _ := m.Task(mod + "/task/t-rubric")
	if !rubric.Rubric || len(rubric.Candidate.Task.Items[0].Slots) != 0 {
		t.Errorf("t-rubric: rubric %v, slots %v", rubric.Rubric, rubric.Candidate.Task.Items[0].Slots)
	}

	tpl, _ := m.Task(mod + "/task/t-template")
	if !tpl.Candidate.Task.Template || tpl.Instances != 3 || tpl.Candidate.Task.KeyProvenance != learnermodel.ProvenanceSolver {
		t.Errorf("t-template: %+v", tpl)
	}
	// The instance an attempt takes is the one the platform's choice names, by the task's ID in
	// the module, and its numbers are the canonical view's in that order.
	for _, s := range []uint64{0, 1, 20261107} {
		n, numbers, err := m.Instance(tpl.Candidate.Task.ID, s)
		want, _ := seed.Instance(s, "t-template", 3)
		if err != nil || n != want || numbers["n"] != int64(want+1) {
			t.Errorf("seed %d: instance %d with %v, %v; want %d", s, n, numbers, err, want)
		}
	}

	if len(m.Lessons) != 2 || m.Lessons[1].Links != 2 || len(m.Lessons[1].Lesson.Tasks) != 1 {
		t.Errorf("lessons %+v", m.Lessons)
	}
	pool, lessons := m.ForAge("adult")
	if len(pool) != 4 || len(lessons) != 1 || lessons[0].ID != mod+"/lesson/l-adult" {
		t.Errorf("an adult gets %d tasks and %v", len(pool), lessons)
	}
	pool, lessons = m.ForAge("under_13")
	if len(pool) != 5 || len(lessons) != 1 || lessons[0].ID != mod+"/lesson/l-child" {
		t.Errorf("a child gets %d tasks and %v", len(pool), lessons)
	}
}

// What a canonical view may hold that the bench cannot translate: each fails with ErrTranslate.
func TestTranslateRefuses(t *testing.T) {
	cases := map[string]func(*formatv0.CanonicalView){
		"no skill graph": func(cv *formatv0.CanonicalView) { cv.SkillGraph = nil },
		"an assessability": func(cv *formatv0.CanonicalView) {
			cv.SkillGraph.Skills["a"] = formatv0.Skill{Assessability: "sometimes"}
		},
		"a band the manifest has not": func(cv *formatv0.CanonicalView) {
			task := cv.Tasks["t-sheet"]
			task.Band = "middle"
			cv.Tasks["t-sheet"] = task
		},
		"a kind of check": func(cv *formatv0.CanonicalView) {
			task := cv.Tasks["t-beside"]
			task.Items[0].AnswerSlots[0].Check.Kind = "spoken"
			cv.Tasks["t-beside"] = task
		},
		"a rubric beside a verdict": func(cv *formatv0.CanonicalView) {
			task := cv.Tasks["t-rubric"]
			task.Items[0].AnswerSlots = append(task.Items[0].AnswerSlots, formatv0.AnswerSlot{Check: formatv0.Check{Kind: formatv0.CheckKindNumber}})
			cv.Tasks["t-rubric"] = task
		},
		"a template with no instance": func(cv *formatv0.CanonicalView) { cv.Instances = nil },
		"a lesson of a missing task": func(cv *formatv0.CanonicalView) {
			l := cv.Lessons["l-adult"]
			l.Segments = []formatv0.Segment{{Task: ptr(formatv0.TaskID("t-missing"))}}
			cv.Lessons["l-adult"] = l
		},
		"a graph NewGraph refuses": func(cv *formatv0.CanonicalView) {
			task := cv.Tasks["t-sheet"]
			task.Skills.Primary = []formatv0.SkillID{"missing"}
			cv.Tasks["t-sheet"] = task
		},
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			cv := view()
			change(cv)
			if _, err := Translate(cv); !errors.Is(err, ErrTranslate) {
				t.Errorf("Translate: %v, want ErrTranslate", err)
			}
		})
	}
}

// archive builds a package's archive as a release has it: a gzip of a tar of the files, in order.
func archive(t *testing.T, files ...[2]string) []byte {
	t.Helper()
	var tarball bytes.Buffer
	for _, f := range files {
		h := make([]byte, 512)
		copy(h, f[0])
		copy(h[124:], strconv.FormatInt(int64(len(f[1])), 8)+"\x00")
		h[156] = '0'
		tarball.Write(h)
		tarball.WriteString(f[1])
		tarball.Write(make([]byte, (512-len(f[1])%512)%512))
	}
	tarball.Write(make([]byte, 1024))
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	zw.Write(tarball.Bytes())
	zw.Close()
	return gz.Bytes()
}

func digest(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// pkg builds the archive of the synthetic module as a release has it, with an inventory that names
// its server part, and the entry an index gives it; change edits the inventory before it is
// written.
func pkg(t *testing.T, cv *formatv0.CanonicalView, change func(*formatv0.Inventory)) ([]byte, formatv0.Entry) {
	t.Helper()
	server, err := json.Marshal(cv)
	if err != nil {
		t.Fatal(err)
	}
	inv := formatv0.Inventory{
		ModuleFormatVersion: cv.Manifest.ModuleFormatVersion, ID: cv.Manifest.ID, Version: cv.Manifest.Version,
		ServerPart: formatv0.Descriptor{Digest: formatv0.Digest("sha256:" + digest(string(server))), Size: int64(len(server))},
	}
	if change != nil {
		change(&inv)
	}
	inventory, err := json.Marshal(inv)
	if err != nil {
		t.Fatal(err)
	}
	e := formatv0.Entry{ModuleFormatVersion: "0.3", PackageDigest: formatv0.Digest("sha256:" + digest(string(inventory))), Release: "v0.1.0"}
	return archive(t, [2]string{"inventory.json", string(inventory)}, [2]string{"server.json", string(server)}), e
}

// A package whose chain of digests holds, and every way it can break.
func TestReadPackage(t *testing.T) {
	good, entry := pkg(t, view(), nil)
	cv, err := ReadPackage(good, mod, "0.1.0", entry)
	if err != nil || cv.Manifest.ID != mod {
		t.Fatalf("ReadPackage: %v", err)
	}
	bad := func(change func(*formatv0.Inventory)) ([]byte, formatv0.Entry) { return pkg(t, view(), change) }
	otherServer, otherServerEntry := bad(func(inv *formatv0.Inventory) { inv.ServerPart.Digest = "sha256:" + formatv0.Digest(digest("other")) })
	otherSize, otherSizeEntry := bad(func(inv *formatv0.Inventory) { inv.ServerPart.Size++ })
	otherID, otherIDEntry := bad(func(inv *formatv0.Inventory) { inv.ID = "example/other" })
	otherManifest, otherManifestEntry := pkg(t, func() *formatv0.CanonicalView { v := view(); v.Manifest.Version = "0.2.0"; return v }(),
		func(inv *formatv0.Inventory) { inv.Version = "0.1.0" })
	for name, c := range map[string]struct {
		data  []byte
		entry formatv0.Entry
		want  error
	}{
		"another inventory":   {good, formatv0.Entry{ModuleFormatVersion: "0.3", PackageDigest: "sha256:" + formatv0.Digest(digest("x"))}, ErrDigest},
		"another server part": {otherServer, otherServerEntry, ErrDigest},
		"another size":        {otherSize, otherSizeEntry, ErrDigest},
		"another module":      {otherID, otherIDEntry, ErrIdentity},
		"another manifest":    {otherManifest, otherManifestEntry, ErrIdentity},
		"no server part":      {archive(t, [2]string{"inventory.json", "{}"}), entry, nil},
		"no gzip":             {[]byte("plain"), entry, nil},
		"a broken header":     {gzipped(t, header("x", "zz")), entry, nil},
		"a short block":       {gzipped(t, header("x", "1000\x00")), entry, nil},
		"a truncated gzip":    {good[:len(good)/2], entry, nil},
		"inventory not JSON": {archive(t, [2]string{"inventory.json", "{"}, [2]string{"server.json", "{}"}),
			formatv0.Entry{ModuleFormatVersion: "0.3", PackageDigest: "sha256:" + formatv0.Digest(digest("{"))}, nil},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := ReadPackage(c.data, mod, "0.1.0", c.entry)
			if err == nil || (c.want != nil && !errors.Is(err, c.want)) {
				t.Errorf("ReadPackage: %v, want %v", err, c.want)
			}
		})
	}
	notJSON := "{"
	inv := `{"module_format_version":"0.3","id":"example/synthetic","version":"0.1.0","server_part":{"digest":"sha256:` + digest(notJSON) + `","size":1}}`
	if _, err := ReadPackage(archive(t, [2]string{"inventory.json", inv}, [2]string{"server.json", notJSON}), mod, "0.1.0",
		formatv0.Entry{ModuleFormatVersion: "0.3", PackageDigest: formatv0.Digest("sha256:" + digest(inv))}); err == nil {
		t.Error("a server part that is not JSON: no error")
	}
}

func header(name, size string) []byte {
	h := make([]byte, 512)
	copy(h, name)
	copy(h[124:], size)
	return h
}

func gzipped(t *testing.T, data []byte) []byte {
	t.Helper()
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	zw.Write(data)
	zw.Close()
	return gz.Bytes()
}

// An index is read by the places of its members whatever its version; an entry is given only if
// the bench may read it.
func TestIndex(t *testing.T) {
	data := `{"module_format_version":"0.4","modules":{"a/b":{"0.1.0":{"module_format_version":"0.3","package_digest":"sha256:00","release":"v0.1.0","added":"later"},"0.2.0":{"module_format_version":"0.4","package_digest":"sha256:01","release":"v0.2.0"},"0.0.9":{"module_format_version":"0.3","package_digest":"sha256:02","release":"v0.0.9","withdrawn":"wrong_key"}}}}`
	idx, err := ReadIndex([]byte(data), digest(data))
	if err != nil {
		t.Fatalf("ReadIndex of an index of a later version: %v", err)
	}
	if e, err := Entry(idx, "a/b", "0.1.0"); err != nil || e.Release != "v0.1.0" {
		t.Errorf("Entry of 0.1.0: %v, %v", e, err)
	}
	if _, err := Entry(idx, "a/b", "0.2.0"); !errors.Is(err, ErrFormat) {
		t.Errorf("Entry of a later format: %v", err)
	}
	if _, err := Entry(idx, "a/b", "0.0.9"); !errors.Is(err, ErrWithdrawn) {
		t.Errorf("Entry of a withdrawn version: %v", err)
	}
	if _, err := Entry(idx, "a/c", "0.1.0"); err == nil {
		t.Error("Entry of a module the index has not: no error")
	}
	if _, err := ReadIndex([]byte(data), digest("x")); !errors.Is(err, ErrDigest) {
		t.Errorf("another digest: %v", err)
	}
	if _, err := ReadIndex([]byte("{"), digest("{")); err == nil {
		t.Error("an index that is not JSON: no error")
	}
}

// files gives the inventory.json and server.json of the synthetic module as a folder holds them;
// change edits the inventory before it is written.
func files(t *testing.T, cv *formatv0.CanonicalView, change func(*formatv0.Inventory)) (inventory, server []byte) {
	t.Helper()
	got, err := Untar(func() []byte { a, _ := pkg(t, cv, change); return a }())
	if err != nil {
		t.Fatal(err)
	}
	return got["inventory.json"], got["server.json"]
}

// A package as a folder holds it reads with an entry of an index and without one, and keeps the
// inventory when it cannot be read further.
func TestReadFiles(t *testing.T) {
	inventory, server := files(t, view(), nil)
	e := formatv0.Entry{ModuleFormatVersion: "0.3", PackageDigest: Digest(inventory)}
	for _, entry := range []*formatv0.Entry{nil, &e} {
		cv, inv, err := ReadFiles(inventory, server, entry)
		if err != nil || cv.Manifest.ID != mod || inv.Version != "0.1.0" {
			t.Fatalf("ReadFiles with the entry %v: %v", entry, err)
		}
	}
	later := func() *formatv0.CanonicalView { v := view(); v.Manifest.ModuleFormatVersion = "0.4"; return v }()
	laterInventory, laterServer := files(t, later, nil)
	otherServer := []byte(string(server) + " ")
	for name, c := range map[string]struct {
		inventory, server []byte
		entry             *formatv0.Entry
		want              error
		inv               bool
	}{
		"another digest":              {inventory, server, &formatv0.Entry{ModuleFormatVersion: "0.3", PackageDigest: Digest([]byte("x"))}, ErrDigest, false},
		"another format in the entry": {inventory, server, &formatv0.Entry{ModuleFormatVersion: "0.4", PackageDigest: Digest(inventory)}, ErrIdentity, true},
		"a format outside the window": {laterInventory, laterServer, nil, ErrFormat, true},
		"another server part":         {inventory, otherServer, nil, ErrDigest, true},
		"inventory not JSON":          {[]byte("{"), server, nil, nil, false},
	} {
		t.Run(name, func(t *testing.T) {
			_, inv, err := ReadFiles(c.inventory, c.server, c.entry)
			if err == nil || (c.want != nil && !errors.Is(err, c.want)) || (inv != nil) != c.inv {
				t.Errorf("ReadFiles: %v, the inventory %v; want %v", err, inv, c.want)
			}
		})
	}
	// An archive of another module than its entry names, whose inventory and manifest agree.
	other := func() *formatv0.CanonicalView { v := view(); v.Manifest.ID = "example/other"; return v }()
	data, entry := pkg(t, other, nil)
	if _, err := ReadPackage(data, mod, "0.1.0", entry); !errors.Is(err, ErrIdentity) {
		t.Errorf("an archive of another module: %v", err)
	}
	if !Readable("0.3") || Readable("0.4") {
		t.Error("the window of the format is not 0.3 alone")
	}
}

// The canonical text changes with what the bench reads — a task tagged to another skill, the
// numbers of an instance — and not with a text for the learner or the version.
func TestText(t *testing.T) {
	text := func(change func(*formatv0.CanonicalView)) string {
		t.Helper()
		v := view()
		change(v)
		m, err := Translate(v)
		if err != nil {
			t.Fatal(err)
		}
		data, err := m.Text()
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	base := text(func(*formatv0.CanonicalView) {})
	if base == "" || !json.Valid([]byte(base)) {
		t.Fatalf("the text is no JSON: %q", base)
	}
	for name, c := range map[string]struct {
		change func(*formatv0.CanonicalView)
		same   bool
	}{
		"the version": {func(v *formatv0.CanonicalView) { v.Manifest.Version = "0.2.0" }, true},
		"an explanation of a lesson": {func(v *formatv0.CanonicalView) {
			v.Lessons["l-child"].Segments[0].Explanation[0] = formatv0.StatementPart{Text: ptr(formatv0.Text("another"))}
		}, true},
		"an option's text": {func(v *formatv0.CanonicalView) { v.Tasks["t-beside"].Items[0].AnswerSlots[0].Options[0] = "yes!" }, true},
		"a task tagged to another skill": {func(v *formatv0.CanonicalView) {
			task := v.Tasks["t-beside"]
			task.Skills.Primary = []formatv0.SkillID{"b"}
			v.Tasks["t-beside"] = task
		}, false},
		"the numbers of an instance": {func(v *formatv0.CanonicalView) { v.Instances["t-template"][0] = formatv0.Numbers{"n": 4} }, false},
		"a difficulty": {func(v *formatv0.CanonicalView) {
			task := v.Tasks["t-sheet"]
			task.Difficulty = 5
			v.Tasks["t-sheet"] = task
		}, false},
	} {
		if got := text(c.change) == base; got != c.same {
			t.Errorf("%s: the same text %v, want %v", name, got, c.same)
		}
	}
}
