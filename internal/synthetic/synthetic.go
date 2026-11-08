// Package synthetic builds the bench's own synthetic packages: canonical views in the shape of the
// three modules of the first slice, for the code and the tests of the bench and not to judge
// candidates by (docs/architecture.md, section 3.3; R03; R19 in docs/decisions.md). A sheet with a
// choice, numbers and a predicate over three bands of grades 1–6, one of whose skills rests on an
// anchor skill; a workbook of one band of grade 3 with pages of up to ten lines, one of whose skills
// teaches that anchor; essays with a rubric and no grades. The published packages have neither a
// skill that rests on another module's anchor nor a cell of predicate (docs/architecture.md,
// section 4), so the populations A- and S-predicate run here only. No learner's data is in them:
// every word is the bench's own (RULES.md, Hard rules).
package synthetic

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/spec/types/formatv0"
)

// The synthetic modules, by full ID, and the anchor skill the workbook teaches and the sheet rests
// on.
const (
	Sheet    = "synthetic/sheet"
	Workbook = "synthetic/workbook"
	Essays   = "synthetic/essays"
	// Anchor is an anchor skill of the shared taxonomy: the format names no other (the text of the
	// format; formatv0.SharedTaxonomy), and the published school maths for grade 3 teaches it too,
	// which matters not, since the synthetic modules never meet the published ones.
	Anchor = string(formatv0.SharedTaxonomyCoreMultiplicationTable)
)

func text(s string) *formatv0.Text { t := formatv0.Text(s); return &t }

func ptr[T any](v T) *T { return &v }

func manifest(id string, subject formatv0.Subject, ages []formatv0.AgeCategory, grades []formatv0.Grade, bands []formatv0.Band) formatv0.Manifest {
	for k := range bands {
		bands[k].Title = formatv0.Text("Band " + string(bands[k].ID))
	}
	return formatv0.Manifest{
		ModuleFormatVersion: "0.3", ID: formatv0.ModuleID(id), Version: "0.0.1", Title: formatv0.Text("Synthetic " + id),
		Subject: subject, Audience: formatv0.Audience{AgeCategories: ages, Grades: grades}, Bands: bands,
		Languages: formatv0.Languages{Explanation: "en"}, Authors: []formatv0.Author{{Name: "The sim bench"}}, License: "Apache-2.0",
	}
}

func skill(title string, a formatv0.SkillAssessability, prerequisites ...string) formatv0.Skill {
	return formatv0.Skill{Title: formatv0.Text(title), Assessability: a, Prerequisites: prerequisites}
}

func statement(s string) []formatv0.StatementPart { return []formatv0.StatementPart{{Text: text(s)}} }

// slotOf gives a slot of a sheet: a choice of n options by its options, a number by its key, a
// predicate by its predicate on the value.
func slotOf(kind formatv0.CheckKind, options int) formatv0.AnswerSlot {
	s := formatv0.AnswerSlot{Binding: formatv0.AnswerSlotBindingAnywhere, Check: formatv0.Check{Kind: kind}}
	switch kind {
	case formatv0.CheckKindChoice:
		s.Binding = formatv0.AnswerSlotBindingChoice
		for k := range options {
			s.Options = append(s.Options, formatv0.Text(fmt.Sprintf("option %d", k+1)))
		}
		s.Check.Key = ptr(string(s.Options[0]))
	case formatv0.CheckKindPredicate:
		s.Check.Predicate = ptr("value % 2 == 1 and 10 < value and value < 30")
	default:
		s.Check.Key = ptr("15")
	}
	return s
}

// SheetView gives the canonical view of the synthetic sheet: four skills with a verdict over three
// bands of grades 1–6 — one answered by a choice of four options, one by a number, one by a
// predicate, and one that rests on the anchor skill and on the first —, a task of each skill at
// difficulties 1, 3 and 5 on each band, and a lesson of each skill on the middle band.
func SheetView() *formatv0.CanonicalView {
	bands := []formatv0.Band{
		{ID: "grades-1-2", Grades: []formatv0.Grade{1, 2}},
		{ID: "grades-3-4", Grades: []formatv0.Grade{3, 4}},
		{ID: "grades-5-6", Grades: []formatv0.Grade{5, 6}},
	}
	cv := &formatv0.CanonicalView{
		Manifest: manifest(Sheet, formatv0.SubjectMathematics, []formatv0.AgeCategory{"under_13", "adult"}, []formatv0.Grade{1, 2, 3, 4, 5, 6}, bands),
		SkillGraph: &formatv0.SkillGraph{Skills: map[formatv0.SkillID]formatv0.Skill{
			"choosing":  skill("Choosing", formatv0.SkillAssessabilityVerdict),
			"counting":  skill("Counting", formatv0.SkillAssessabilityVerdict),
			"proving":   skill("Proving", formatv0.SkillAssessabilityVerdict),
			"on-anchor": skill("On the anchor", formatv0.SkillAssessabilityVerdict, "choosing", Anchor),
		}},
		SkillTiers: formatv0.SkillTiers{"choosing": 1, "counting": 1, "proving": 1, "on-anchor": 2},
		Tasks:      map[formatv0.TaskID]formatv0.Task{},
		Lessons:    map[formatv0.LessonID]formatv0.Lesson{},
	}
	cv.Manifest.Dependencies = []formatv0.SharedTaxonomy{formatv0.SharedTaxonomy(Anchor)}
	slots := map[formatv0.SkillID]formatv0.AnswerSlot{
		"choosing":  slotOf(formatv0.CheckKindChoice, 4),
		"counting":  slotOf(formatv0.CheckKindNumber, 0),
		"proving":   slotOf(formatv0.CheckKindPredicate, 0),
		"on-anchor": slotOf(formatv0.CheckKindNumber, 0),
	}
	for _, skill := range []formatv0.SkillID{"choosing", "counting", "proving", "on-anchor"} {
		for _, band := range bands {
			for _, d := range []int64{1, 3, 5} {
				id := formatv0.TaskID(fmt.Sprintf("%s-%s-d%d", skill, band.ID, d))
				cv.Tasks[id] = formatv0.Task{
					ID: id, Canvas: formatv0.Canvas{Kind: formatv0.CanvasKindHandwritingSheet}, Band: band.ID, Difficulty: d,
					Statement: statement("A synthetic task."), Solution: statement("A synthetic solution."),
					Skills: formatv0.Skills{Primary: []formatv0.SkillID{skill}},
					Items:  []formatv0.Item{{AnswerSlots: []formatv0.AnswerSlot{slots[skill]}}},
				}
			}
		}
		lesson := formatv0.LessonID(string(skill) + "-grades-3-4")
		cv.Lessons[lesson] = formatv0.Lesson{ID: lesson, Title: "A lesson", Skills: []formatv0.SkillID{skill}, Segments: []formatv0.Segment{
			{Explanation: statement("An explanation.")},
			{Task: ptr(formatv0.TaskID(fmt.Sprintf("%s-grades-3-4-d1", skill)))},
			{Task: ptr(formatv0.TaskID(fmt.Sprintf("%s-grades-3-4-d3", skill)))},
		}}
	}
	return cv
}

// WorkbookView gives the canonical view of the synthetic workbook: one band of grade 3, three skills
// with a verdict in a chain — the first teaches the anchor skill —, pages of six, eight and ten
// lines of numbers and choices of three options at difficulties 1, 3 and 5 and another page at 3,
// and a lesson of each skill.
func WorkbookView() *formatv0.CanonicalView {
	cv := &formatv0.CanonicalView{
		Manifest: manifest(Workbook, formatv0.SubjectMathematics, []formatv0.AgeCategory{"under_13"}, []formatv0.Grade{3}, []formatv0.Band{{ID: "grade-3", Grades: []formatv0.Grade{3}}}),
		SkillGraph: &formatv0.SkillGraph{Skills: map[formatv0.SkillID]formatv0.Skill{
			"tables":    {Title: "Tables", Assessability: formatv0.SkillAssessabilityVerdict, Teaches: []formatv0.SharedTaxonomy{formatv0.SharedTaxonomy(Anchor)}},
			"remainder": skill("Remainder", formatv0.SkillAssessabilityVerdict, "tables"),
			"compare":   skill("Compare", formatv0.SkillAssessabilityVerdict, "remainder"),
		}},
		SkillTiers: formatv0.SkillTiers{"tables": 1, "remainder": 2, "compare": 3},
		Tasks:      map[formatv0.TaskID]formatv0.Task{},
		Lessons:    map[formatv0.LessonID]formatv0.Lesson{},
	}
	for _, skill := range []formatv0.SkillID{"tables", "remainder", "compare"} {
		for k, page := range []struct {
			d     int64
			lines int
		}{{1, 6}, {3, 8}, {5, 10}, {3, 7}} {
			id := formatv0.TaskID(fmt.Sprintf("%s-page-%d", skill, k+1))
			var items []formatv0.Item
			for line := range page.lines {
				check := &formatv0.WrittenCheck{Kind: formatv0.CheckKindNumber, Key: ptr("56")}
				if line%3 == 2 {
					check = &formatv0.WrittenCheck{Kind: formatv0.CheckKindChoice, Key: ptr("<"), Options: []formatv0.Text{"<", ">", "="}}
				}
				items = append(items, formatv0.Item{Line: []formatv0.LinePart{{Text: text("7 × 8")}, {Binding: ptr(formatv0.LinePartBindingBox), Check: check}}})
			}
			cv.Tasks[id] = formatv0.Task{
				ID: id, Canvas: formatv0.Canvas{Kind: formatv0.CanvasKindWorkbook}, Band: "grade-3", Difficulty: page.d,
				Statement: statement("Fill in the page."), Solution: statement("A synthetic solution."),
				Skills: formatv0.Skills{Primary: []formatv0.SkillID{skill}}, Items: items,
			}
		}
		lesson := formatv0.LessonID(string(skill) + "-lesson")
		cv.Lessons[lesson] = formatv0.Lesson{ID: lesson, Title: "A lesson", Skills: []formatv0.SkillID{skill}, Segments: []formatv0.Segment{
			{Explanation: statement("An explanation.")},
			{Task: ptr(formatv0.TaskID(string(skill) + "-page-1"))},
		}}
	}
	return cv
}

// EssaysView gives the canonical view of the synthetic essays: no grades and every age category; a
// skill with a verdict answered by a choice of three options, a skill without assessment taught by a
// lesson of reading alone, and two skills with a rubric in a chain on it; tasks with a rubric of
// three criteria — one for adults alone —, choices, and a lesson of each skill but the one with a
// verdict.
func EssaysView() *formatv0.CanonicalView {
	cv := &formatv0.CanonicalView{
		Manifest: manifest(Essays, formatv0.SubjectWriting, []formatv0.AgeCategory{"under_13", "from_13_to_consent", "from_consent_to_18", "adult"}, nil,
			[]formatv0.Band{{ID: "first-steps"}, {ID: "whole-essays"}}),
		SkillGraph: &formatv0.SkillGraph{Skills: map[formatv0.SkillID]formatv0.Skill{
			"choosing-a-sentence": skill("Choosing a sentence", formatv0.SkillAssessabilityVerdict),
			"reading":             skill("Reading", formatv0.SkillAssessabilityNone),
			"writing":             skill("Writing", formatv0.SkillAssessabilityRubric, "reading"),
			"improving":           skill("Improving", formatv0.SkillAssessabilityRubric, "writing"),
		}},
		SkillTiers: formatv0.SkillTiers{"choosing-a-sentence": 1, "reading": 1, "writing": 2, "improving": 3},
		Tasks:      map[formatv0.TaskID]formatv0.Task{},
		Lessons:    map[formatv0.LessonID]formatv0.Lesson{},
	}
	sheet := formatv0.Canvas{Kind: formatv0.CanvasKindHandwritingSheet}
	rubric := formatv0.AnswerSlot{Binding: formatv0.AnswerSlotBindingSheet, Check: formatv0.Check{Kind: formatv0.CheckKindRubric, Criteria: []formatv0.Criterion{{Text: "One"}, {Text: "Two"}, {Text: "Three"}}}}
	for k := range 2 {
		id := formatv0.TaskID(fmt.Sprintf("choosing-%d", k+1))
		cv.Tasks[id] = formatv0.Task{ID: id, Canvas: sheet, Band: "first-steps", Difficulty: int64(2*k + 1),
			Statement: statement("Choose the best first sentence."), Solution: statement("A synthetic solution."),
			Skills: formatv0.Skills{Primary: []formatv0.SkillID{"choosing-a-sentence"}},
			Items:  []formatv0.Item{{AnswerSlots: []formatv0.AnswerSlot{slotOf(formatv0.CheckKindChoice, 3)}}}}
	}
	for _, w := range []struct {
		id    formatv0.TaskID
		skill formatv0.SkillID
		band  formatv0.BandID
		ages  []formatv0.AgeCategory
	}{
		{"a-first-essay", "writing", "first-steps", nil},
		{"a-second-essay", "writing", "first-steps", nil},
		{"an-essay-for-adults", "writing", "whole-essays", []formatv0.AgeCategory{"adult"}},
		{"an-essay-improved", "improving", "whole-essays", nil},
	} {
		cv.Tasks[w.id] = formatv0.Task{ID: w.id, Canvas: sheet, Band: w.band, Difficulty: 3, AgeCategories: w.ages, Statement: statement("Write a short essay."),
			Skills: formatv0.Skills{Primary: []formatv0.SkillID{w.skill}}, Items: []formatv0.Item{{AnswerSlots: []formatv0.AnswerSlot{rubric}}}}
	}
	cv.Lessons["reading-lesson"] = formatv0.Lesson{ID: "reading-lesson", Title: "Reading", Skills: []formatv0.SkillID{"reading"}, Segments: []formatv0.Segment{
		{Explanation: statement("A model essay.")}, {Explanation: statement("Another model essay.")},
	}}
	cv.Lessons["writing-lesson"] = formatv0.Lesson{ID: "writing-lesson", Title: "Writing", Skills: []formatv0.SkillID{"writing"}, Segments: []formatv0.Segment{
		{Explanation: statement("How to begin.")}, {Task: ptr(formatv0.TaskID("a-first-essay"))},
	}}
	cv.Lessons["improving-lesson"] = formatv0.Lesson{ID: "improving-lesson", Title: "Improving", Skills: []formatv0.SkillID{"improving"}, Segments: []formatv0.Segment{
		{Explanation: statement("How to improve.")}, {Task: ptr(formatv0.TaskID("an-essay-improved"))},
	}}
	return cv
}

// views are the synthetic canonical views, by module.
var views = map[string]func() *formatv0.CanonicalView{Sheet: SheetView, Workbook: WorkbookView, Essays: EssaysView}

// read translates the synthetic views once: they are fixed, and so is what they give.
var read = sync.OnceValues(func() (map[string]*pack.Module, error) {
	out := map[string]*pack.Module{}
	for _, id := range slices.Sorted(maps.Keys(views)) {
		m, err := pack.Translate(views[id]())
		if err != nil {
			return nil, fmt.Errorf("synthetic: %s: %w", id, err)
		}
		out[id] = m
	}
	return out, nil
})

// IsSynthetic tells whether a module is a synthetic one: the criterion and the gate refuse it, since
// they judge on published packages alone (R03, R04).
func IsSynthetic(module string) bool {
	return strings.HasPrefix(module, "synthetic/")
}

// Modules gives the synthetic modules, by full ID, translated as the published ones are. They are
// shared, and read only.
func Modules() (map[string]*pack.Module, error) {
	return read()
}
