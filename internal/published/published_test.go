package published

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/spec/types/formatv0"
)

const (
	olympiad = "novatrail/olympiad-math"
	grade3   = "novatrail/grade-3-math"
	essays   = "novatrail/english-essays"
)

func modules(t *testing.T) map[string]*pack.Module {
	t.Helper()
	ms, err := Modules()
	if err != nil {
		t.Fatal(err)
	}
	if again, _ := Modules(); again[olympiad] != ms[olympiad] {
		t.Fatal("Modules read the packages again")
	}
	return ms
}

// The packages read through the chain of digests, and their translation holds what the release
// has: the skills with their tiers in the order of the graph, the bands, the tasks, the lessons —
// the same facts engine's test of the release holds (T21 engine), so that a translation of sim
// that parted from it shows here.
func TestReleasePackages(t *testing.T) {
	ms := modules(t)
	for _, tc := range []struct {
		module  string
		skills  []string // local ID and tier, ordered by ID
		bands   int
		tasks   int
		lessons int
	}{
		{olympiad, []string{"algorithms-weighing-pouring:3", "arithmetic-tricks:1", "combinatorics-enumeration:1", "counting-figures:1",
			"counting-gaps:1", "fractions-parts:2", "games-strategy:2", "geometry-cubes-paper:1", "geometry-grid:2", "logic-knights-liars:2",
			"logic-ordering:1", "logic-sets:3", "number-divisibility:3", "parity-alternation:1", "percent-basic:3", "pigeonhole-basic:2",
			"ratio-sharing:3", "time-calendar:2", "time-clocks:1"}, 3, 772, 41},
		{grade3, []string{"comparing-fractions:3", "division-with-remainder:2", "fractions-of-a-whole:2", "times-tables:1"}, 1, 12, 4},
		{essays, []string{"choosing-a-first-sentence:1", "endings:4", "improving-an-essay:4", "paragraphs:3",
			"reading-model-essays:1", "writing-a-first-sentence:2"}, 2, 8, 6},
	} {
		m := ms[tc.module]
		var got []string
		for _, s := range m.Graph.Skills {
			got = append(got, strings.TrimPrefix(s.ID, tc.module+"/skill/")+":"+itoa(s.Tier))
		}
		if !slices.Equal(got, tc.skills) {
			t.Errorf("%s: skills %v, want %v", tc.module, got, tc.skills)
		}
		if len(m.Graph.Bands) != tc.bands || len(m.Tasks) != tc.tasks || len(m.Lessons) != tc.lessons {
			t.Errorf("%s: %d bands, %d tasks, %d lessons; want %d, %d, %d", tc.module, len(m.Graph.Bands), len(m.Tasks), len(m.Lessons), tc.bands, tc.tasks, tc.lessons)
		}
	}
	g3 := ms[grade3]
	i := slices.IndexFunc(g3.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == grade3+"/skill/fractions-of-a-whole" })
	if !slices.Contains(g3.Graph.Skills[i].Teaches, "core.fractions-of-a-whole") {
		t.Errorf("fractions-of-a-whole teaches %v", g3.Graph.Skills[i].Teaches)
	}
	// A task or a lesson with no age categories of its own takes the module's: school maths for
	// grade 3 is under_13 alone; of essays, one task and one lesson are for adults alone.
	if pool, lessons := ms[essays].ForAge("under_13"); len(lessons) != 5 || len(pool) != 7 {
		t.Errorf("a child gets %d lessons and %d tasks of essays, want 5 and 7", len(lessons), len(pool))
	}
	if pool, lessons := g3.ForAge("adult"); len(lessons) != 0 || len(pool) != 0 {
		t.Errorf("an adult gets %d lessons and %d tasks of school maths for grade 3, want none", len(lessons), len(pool))
	}
	l := lesson(t, ms[olympiad], olympiad+"/lesson/counting-gaps-grades-3-4")
	var ids []string
	for _, task := range l.Lesson.Tasks {
		ids = append(ids, task.ID)
	}
	want := []string{olympiad + "/task/gaps-34-d1-1", olympiad + "/task/gaps-34-d3-1", olympiad + "/task/gaps-34-d5-1"}
	if !slices.Equal(ids, want) {
		t.Errorf("lesson counting-gaps-grades-3-4 holds %v, want %v", ids, want)
	}
}

// The kinds of check of the release: numbers, fractions and choices of five options in the check on
// the sheet of olympiad maths 0.1.1, numbers, fractions and choices of three options on the workbook
// of grade 3 — the options in the check —, choices of
// three options beside the check and rubrics in essays (R12 of this log; docs/architecture.md,
// section 4). No cell of predicate is published.
func TestReleaseChecks(t *testing.T) {
	ms := modules(t)
	count := func(m *pack.Module) map[string]int {
		c := map[string]int{}
		for _, task := range m.Tasks {
			if task.Rubric {
				c["rubric"]++
			}
			for _, it := range task.Candidate.Task.Items {
				for _, s := range it.Slots {
					c[itoa(int32(s.Check))+"/"+itoa(s.Options)]++
				}
			}
		}
		return c
	}
	number, fraction, choice := itoa(int32(learnermodel.CheckKindNumber)), itoa(int32(learnermodel.CheckKindFraction)), itoa(int32(learnermodel.CheckKindChoice))
	if c := count(ms[olympiad]); c[number+"/0"] != 569 || c[fraction+"/0"] != 7 || c[choice+"/5"] != 196 || len(c) != 3 {
		t.Errorf("olympiad maths has the checks %v", c)
	}
	if c := count(ms[grade3]); c[number+"/0"] != 78 || c[fraction+"/0"] != 17 || c[choice+"/3"] != 17 || len(c) != 3 {
		t.Errorf("grade 3 has the checks %v", c)
	}
	if c := count(ms[essays]); c[choice+"/3"] != 3 || c["rubric"] != 5 || len(c) != 2 {
		t.Errorf("essays have the checks %v", c)
	}
	for _, m := range ms {
		for _, task := range m.Tasks {
			if (task.Canvas == "") || (task.Rubric && len(slices.Concat(itemSlots(task)...)) > 0) {
				t.Errorf("%s: canvas %q, rubric %v", task.Candidate.Task.ID, task.Canvas, task.Rubric)
			}
		}
	}
}

// Of the release, only olympiad maths 0.1.1 has templates: 19 tasks, each with its instances, and an
// attempt's seed takes one of them; a task that is no template has no instance (docs/architecture.md,
// section 4). How the seed takes an instance is tested on a template of pack's own test.
func TestReleaseTemplates(t *testing.T) {
	ms := modules(t)
	for id, want := range map[string]int{olympiad: 19, grade3: 0, essays: 0} {
		m, n := ms[id], 0
		for _, task := range m.Tasks {
			k, _, err := m.Instance(task.Candidate.Task.ID, 1)
			switch {
			case task.Candidate.Task.Template:
				n++
				if task.Instances == 0 || err != nil || k < 0 || k >= task.Instances {
					t.Errorf("%s: %d instances, the instance %d, %v", task.Candidate.Task.ID, task.Instances, k, err)
				}
			case task.Instances != 0 || err == nil:
				t.Errorf("%s is no template, and still has an instance", task.Candidate.Task.ID)
			}
		}
		if n != want {
			t.Errorf("%s has %d templates, want %d", id, n, want)
		}
	}
	if _, _, err := ms[olympiad].Instance("no/such/task", 1); !errors.Is(err, pack.ErrNoTask) {
		t.Errorf("Instance of a missing task: %v", err)
	}
}

// A pin that does not hold: another index digest, a module or a version the index has not.
func TestPinsRefuse(t *testing.T) {
	if _, err := readPin(Release, strings.Repeat("0", 64), versions); !errors.Is(err, pack.ErrDigest) {
		t.Errorf("another index digest: %v", err)
	}
	if _, err := readPin(Release, IndexDigest, map[formatv0.ModuleID]formatv0.ModuleVersion{olympiad: "9.9.9"}); err == nil {
		t.Error("a version the index has not: no error")
	}
	if _, err := readPin("v9.9.9", IndexDigest, versions); err == nil {
		t.Error("a release that is not embedded: no error")
	}
}

func lesson(t *testing.T, m *pack.Module, id string) pack.Lesson {
	t.Helper()
	for _, l := range m.Lessons {
		if l.Lesson.ID == id {
			return l
		}
	}
	t.Fatalf("%s has no lesson %s", m.ID, id)
	return pack.Lesson{}
}

func itemSlots(task pack.Task) [][]learnermodel.Slot {
	var out [][]learnermodel.Slot
	for _, it := range task.Candidate.Task.Items {
		out = append(out, it.Slots)
	}
	return out
}

func itoa[N ~int32](n N) string {
	return strconv.Itoa(int(n))
}

// Pinned names the three packages of the first slice by the pinned index, each with its version and
// package digest: olympiad maths 0.1.1, the others 0.1.0.
func TestPinned(t *testing.T) {
	ps, err := Pinned()
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 3 || ps[0].Module != "novatrail/english-essays" || ps[2].Module != "novatrail/olympiad-math" {
		t.Fatalf("Pinned = %+v", ps)
	}
	for _, p := range ps {
		want := "0.1.0"
		if p.Module == olympiad {
			want = "0.1.1"
		}
		if p.Version != want || len(p.Digest) != len("sha256:")+64 {
			t.Errorf("%+v", p)
		}
	}
}

// Archive gives each pinned package from the folder of the release its entry names, and refuses a
// package the pinned index has not.
func TestArchive(t *testing.T) {
	ps, err := Pinned()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range ps {
		data, err := Archive(p)
		if err != nil {
			t.Fatal(err)
		}
		if files, err := pack.Untar(data); err != nil || files["inventory.json"] == nil {
			t.Errorf("%s: %v", p.Module, err)
		}
	}
	if _, err := Archive(Package{Module: olympiad, Version: "9.9.9"}); err == nil {
		t.Error("a package the index has not: no error")
	}
}
