package pack

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/spec/types/formatv0"
)

// anchorPrefix opens the ID of an anchor skill, a skill of the shared taxonomy that another
// module's skills teach (R50 spec): it stays as it is, with no module before it.
const anchorPrefix = "core."

// Module is a module as the bench reads it: the graph and the inputs of engine, translated once
// from the canonical view, and what the virtual learners need beside them.
type Module struct {
	ID    string
	Graph learnermodel.Graph // by learnermodel.NewGraph, for both libraries (R28 spec)
	Ages  []string           // the module's age categories, its audience
	Tasks []Task             // ordered by full ID
	// Lessons are the module's lessons, ordered by full ID.
	Lessons []Lesson
}

// Task is a task of the module: what engine takes of it, and what the bench reads beside.
type Task struct {
	Candidate engine.Candidate // the task as engine takes it, with the traps of its items
	Local     string           // its ID in the module
	Canvas    formatv0.CanvasKind
	Rubric    bool     // a task with a rubric: no verdict and no evidence (architecture 3.8 spec)
	Ages      []string // its own age categories; none, the module's
	Instances int      // the number of instances of a template; 0 for a task

	numbers []formatv0.Numbers // the instances of a template, in the order the format numbers them
}

// Lesson is a lesson of the module: what engine takes of it, and its links beside.
type Lesson struct {
	Lesson engine.Lesson
	Links  int // its segments: explanations, worked examples and tasks, in the author's order
	// Order gives each link in the author's order: the index of its task in Lesson.Tasks, or -1 for
	// an explanation or a worked example.
	Order []int
	Ages  []string // its own age categories; none, the module's
}

// ErrTranslate is a canonical view the bench cannot translate.
var ErrTranslate = errors.New("pack: the canonical view cannot be translated")

// Translate reads a canonical view into the bench's module: full IDs — the module, the kind and
// the local ID joined by / (the text of the format, section 10) —, anchor skills as they are;
// bands by place from the manifest; tiers from skill_tiers; a slot's kind of check and a choice's
// number of options, whether they stand beside the check or in it; a check by a rubric gives no
// slot with a verdict. The graph is learner-model's own NewGraph, so that sim and api build one
// graph of one module alike (R28 spec).
func Translate(cv *formatv0.CanonicalView) (*Module, error) {
	module := string(cv.Manifest.ID)
	full := func(kind, id string) string {
		if strings.HasPrefix(id, anchorPrefix) {
			return id
		}
		return module + "/" + kind + "/" + id
	}
	fulls := func(kind string, ids []string) []string {
		var out []string
		for _, id := range ids {
			out = append(out, full(kind, id))
		}
		return out
	}

	var bands []learnermodel.Band
	place := map[formatv0.BandID]int32{}
	for k, b := range cv.Manifest.Bands {
		var grades []int32
		for _, g := range b.Grades {
			grades = append(grades, int32(g))
		}
		bands = append(bands, learnermodel.Band{Place: int32(k + 1), ID: string(b.ID), Grades: grades})
		place[b.ID] = int32(k + 1)
	}

	if cv.SkillGraph == nil {
		return nil, fmt.Errorf("%w: %s has no skill graph", ErrTranslate, module)
	}
	var skills []learnermodel.Skill
	for _, id := range sortedKeys(cv.SkillGraph.Skills) {
		sk := cv.SkillGraph.Skills[id]
		a, ok := assessability[sk.Assessability]
		if !ok {
			return nil, fmt.Errorf("%w: %s, skill %s has the assessability %q", ErrTranslate, module, id, sk.Assessability)
		}
		var teaches []string
		for _, t := range sk.Teaches {
			teaches = append(teaches, string(t))
		}
		skills = append(skills, learnermodel.Skill{
			ID: full("skill", string(id)), Assessability: a, Prerequisites: fulls("skill", sk.Prerequisites),
			Teaches: teaches, Tier: int32(cv.SkillTiers[id]),
		})
	}
	var retired []learnermodel.RetiredSkill
	for _, id := range sortedKeys(cv.SkillGraph.RetiredSkills) {
		var by []string
		for _, r := range cv.SkillGraph.RetiredSkills[id].ReplacedBy {
			by = append(by, full("skill", string(r)))
		}
		retired = append(retired, learnermodel.RetiredSkill{ID: full("skill", string(id)), ReplacedBy: by})
	}

	m := &Module{ID: module}
	for _, a := range cv.Manifest.Audience.AgeCategories {
		m.Ages = append(m.Ages, string(a))
	}
	var tasks []learnermodel.Task
	byLocal := map[formatv0.TaskID]learnermodel.Task{}
	for _, id := range sortedKeys(cv.Tasks) {
		st := cv.Tasks[id]
		p, ok := place[st.Band]
		if !ok {
			return nil, fmt.Errorf("%w: %s, task %s names the band %q, which the manifest has not", ErrTranslate, module, id, st.Band)
		}
		task := learnermodel.Task{
			ID: full("task", string(id)), Template: st.Template != nil,
			Band: learnermodel.BandRef{Place: p, ID: string(st.Band)}, Difficulty: int32(st.Difficulty),
			Primary: fulls("skill", skillIDs(st.Skills.Primary)), Supporting: fulls("skill", st.Skills.Supporting),
			// A key in a module is the author's; a template's instance's is computed by its
			// expression or its solver, checked by a solver (concepts 9.4 spec). The skill tags
			// are the author's.
			KeyProvenance: learnermodel.ProvenanceAuthor, SkillTagsProvenance: learnermodel.ProvenanceAuthor,
		}
		if st.Template != nil {
			task.KeyProvenance = learnermodel.ProvenanceSolver
		}
		var traps []engine.ItemTraps
		rubric, verdict := false, false
		for k, it := range st.Items {
			item := learnermodel.TaskItem{Place: int32(k + 1)}
			if it.Skills != nil {
				item.OwnSkills, item.Primary = true, fulls("skill", skillIDs(it.Skills.Primary))
				item.Supporting = fulls("skill", it.Skills.Supporting)
			}
			var caught []string
			for _, sl := range slots(it) {
				if sl.kind == formatv0.CheckKindRubric {
					rubric = true
					continue
				}
				kind, ok := checkKind[sl.kind]
				if !ok {
					return nil, fmt.Errorf("%w: %s, task %s has a check of the kind %q", ErrTranslate, module, id, sl.kind)
				}
				verdict = true
				slot := learnermodel.Slot{Check: kind}
				if kind == learnermodel.CheckKindChoice {
					slot.Options, slot.HasOptions = int32(sl.options), true
				}
				item.Slots = append(item.Slots, slot)
				for _, tv := range sl.traps {
					caught = append(caught, full("trap", string(tv)))
				}
			}
			task.Items = append(task.Items, item)
			if len(caught) > 0 {
				slices.Sort(caught)
				traps = append(traps, engine.ItemTraps{Place: item.Place, Traps: slices.Compact(caught)})
			}
		}
		if rubric && verdict {
			return nil, fmt.Errorf("%w: %s, task %s mixes a rubric with checks that give a verdict", ErrTranslate, module, id)
		}
		t := Task{
			Candidate: engine.Candidate{Task: task, Traps: traps}, Local: string(id), Canvas: st.Canvas.Kind,
			Rubric: rubric, Instances: len(cv.Instances[id]), numbers: cv.Instances[id],
		}
		for _, a := range st.AgeCategories {
			t.Ages = append(t.Ages, string(a))
		}
		if st.Template != nil && t.Instances == 0 {
			return nil, fmt.Errorf("%w: %s, template %s has no instance in the canonical view", ErrTranslate, module, id)
		}
		tasks = append(tasks, task)
		byLocal[id] = task
		m.Tasks = append(m.Tasks, t)
	}

	g, err := learnermodel.NewGraph(module, bands, skills, retired, tasks)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrTranslate, module, err)
	}
	m.Graph = g

	for _, id := range sortedKeys(cv.Lessons) {
		sl := cv.Lessons[id]
		l := Lesson{Lesson: engine.Lesson{ID: full("lesson", string(id)), Skills: fulls("skill", skillIDs(sl.Skills))}, Links: len(sl.Segments)}
		for _, seg := range sl.Segments {
			if seg.Task == nil {
				l.Order = append(l.Order, -1)
				continue
			}
			l.Order = append(l.Order, len(l.Lesson.Tasks))
			task, ok := byLocal[*seg.Task]
			if !ok {
				return nil, fmt.Errorf("%w: %s, lesson %s names the task %s, which the module has not", ErrTranslate, module, id, *seg.Task)
			}
			l.Lesson.Tasks = append(l.Lesson.Tasks, task)
		}
		for _, a := range sl.AgeCategories {
			l.Ages = append(l.Ages, string(a))
		}
		m.Lessons = append(m.Lessons, l)
	}
	return m, nil
}

var assessability = map[formatv0.SkillAssessability]learnermodel.Assessability{
	formatv0.SkillAssessabilityVerdict: learnermodel.AssessabilityVerdict,
	formatv0.SkillAssessabilityRubric:  learnermodel.AssessabilityRubric,
	formatv0.SkillAssessabilityNone:    learnermodel.AssessabilityNone,
}

var checkKind = map[formatv0.CheckKind]learnermodel.CheckKind{
	formatv0.CheckKindNumber:    learnermodel.CheckKindNumber,
	formatv0.CheckKindFraction:  learnermodel.CheckKindFraction,
	formatv0.CheckKindChoice:    learnermodel.CheckKindChoice,
	formatv0.CheckKindPredicate: learnermodel.CheckKindPredicate,
}

// slot is a place of an answer of an item, whichever way the item holds it: an answer slot of a
// sheet, or a part of a line of a workbook that has a check.
type slot struct {
	kind    formatv0.CheckKind
	options int // a choice's options: beside the check, or in it
	traps   []formatv0.TrapID
}

func slots(it formatv0.Item) []slot {
	var out []slot
	for _, sl := range it.AnswerSlots {
		s := slot{kind: sl.Check.Kind, options: len(sl.Options) + len(sl.Check.Options)}
		for _, tv := range sl.Check.TrapValues {
			s.traps = append(s.traps, tv.Trap)
		}
		out = append(out, s)
	}
	for _, part := range it.Line {
		if part.Check == nil {
			continue
		}
		s := slot{kind: part.Check.Kind, options: len(part.Check.Options)}
		for _, tv := range part.Check.TrapValues {
			s.traps = append(s.traps, tv.Trap)
		}
		out = append(out, s)
	}
	return out
}

func skillIDs(ids []formatv0.SkillID) []string {
	var out []string
	for _, id := range ids {
		out = append(out, string(id))
	}
	return out
}

// sortedKeys gives the keys of a map in order: the canonical view keys its entries by ID, and the
// bench reads them in the order of their IDs, never the map's.
func sortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	return slices.Sorted(maps.Keys(m))
}

// suits tells whether a task or a lesson with the age categories ages suits a learner of the age
// category age: its own age categories, or else the module's (R107 spec).
func (m *Module) suits(ages []string, age string) bool {
	if len(ages) == 0 {
		ages = m.Ages
	}
	return slices.Contains(ages, age)
}

// ForAge gives the pool and the lessons as api gives them to a learner of an age category: the
// tasks that suit them, and the lessons that suit them whose every task does too (R107 spec;
// architecture 3.7 api). The age category is api's filter, not an input of engine.
func (m *Module) ForAge(age string) ([]engine.Candidate, []engine.Lesson) {
	var pool []engine.Candidate
	suitsTask := map[string]bool{}
	for _, t := range m.Tasks {
		ok := m.suits(t.Ages, age)
		suitsTask[t.Candidate.Task.ID] = ok
		if ok {
			pool = append(pool, t.Candidate)
		}
	}
	var lessons []engine.Lesson
	for _, l := range m.Lessons {
		if !m.suits(l.Ages, age) || slices.ContainsFunc(l.Lesson.Tasks, func(t learnermodel.Task) bool { return !suitsTask[t.ID] }) {
			continue
		}
		lessons = append(lessons, l.Lesson)
	}
	return pool, lessons
}

// ErrNoTask is a task the module has not.
var ErrNoTask = errors.New("pack: no such task in the module")

// Task gives a task of the module by its full ID.
func (m *Module) Task(id string) (*Task, error) {
	i, ok := slices.BinarySearchFunc(m.Tasks, id, func(t Task, id string) int { return cmp.Compare(t.Candidate.Task.ID, id) })
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoTask, id)
	}
	return &m.Tasks[i], nil
}

// Instance gives the instance of a template an attempt with seed attemptSeed takes: its number,
// chosen as the platform chooses it — the place names the task's ID in the module, as the format's
// task has it (the text of the format, section 23; R60 spec) —, and its numbers, read from the canonical view, whose instances are in that order. A task that is no
// template has no instance.
func (m *Module) Instance(id string, attemptSeed uint64) (int, formatv0.Numbers, error) {
	t, err := m.Task(id)
	if err != nil {
		return 0, nil, err
	}
	n, err := seed.Instance(attemptSeed, t.Local, t.Instances)
	if err != nil {
		return 0, nil, err
	}
	return n, t.numbers[n], nil
}

// ItemPrimary gives the primary skills of the item at index i of the task: the item's own, when it
// has skills of its own, and the task's otherwise (the text of the format, section 12). The truth,
// the engine's results and the measures read an item's skills by this one rule.
func (t *Task) ItemPrimary(i int) []string {
	it := t.Candidate.Task.Items[i]
	if it.OwnSkills {
		return it.Primary
	}
	return t.Candidate.Task.Primary
}

// Text gives the canonical text of what the bench reads of the module: the module translated —
// its graph with the bands and the tiers, its audience, its tasks with their items, slots, traps,
// canvases and ages and the numbers of each template's instances, its lessons —, as JSON, whose
// fields keep their order and whose maps are written by their keys. What Translate does not read —
// the texts for the learner, the version — is not in it, so two packages give the same text when
// the bench reads the same of them, and Translate is the one list of what the bench reads
// (docs/architecture.md, section 3.9; R27 in docs/decisions.md).
func (m *Module) Text() ([]byte, error) {
	type task struct {
		Task
		Numbers []formatv0.Numbers
	}
	v := struct {
		Graph   learnermodel.Graph
		Ages    []string
		Tasks   []task
		Lessons []Lesson
	}{Graph: m.Graph, Ages: m.Ages, Lessons: m.Lessons}
	for _, t := range m.Tasks {
		v.Tasks = append(v.Tasks, task{Task: t, Numbers: t.numbers})
	}
	return json.Marshal(v)
}
