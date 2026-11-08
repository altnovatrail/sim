package run

import (
	"fmt"
	"slices"

	learnermodel "github.com/altnovatrail/learner-model"
	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
)

// Obs is what the measures read of one attempt with a verdict (docs/architecture.md, section 3.5;
// R22): the model and the truth of its items at the attempt, and of the skills of its module before
// and after it.
type Obs struct {
	At      int64     // when the attempt was given
	Repeat  bool      // the task, not a template, had its solution seen before: a repeat (M8)
	Answers int       // the learner's first answers of items with a verdict before the attempt
	Items   []ItemObs // beside the attempt's items
	Before  []SkillObs
	After   []SkillObs
	// Level and TrueLevel are θ of the model and θ* of the truth of the attempt's module after it.
	Level, TrueLevel float64
}

// ItemObs is one item of an attempt.
type ItemObs struct {
	PTrue  float64 // P* of the first answer
	PModel float64 // P of the model at the attempt, with retention
	// PFresh is P of the model by the same estimate with R̄ = 1, c + (P − c)/R̄, by its own formula
	// c + (1 − c) · R̄ · σ, R̄ the mean R of the item's recorded skills with a verdict in the order of
	// their IDs, 1 for a skill with no practice and in the trial series (learner-model v0.1.0,
	// chanceOf in step.go), for M10.
	PFresh float64
	Pause  int64 // seconds since the last practice of the item's primary skills, -1 for none
	Help   bool  // the answer had help: the technique on request, or a task whose solution was seen
	Many   bool  // the attempt has more than one item: a page
}

// SkillObs is one skill with a verdict of the attempt's module.
type SkillObs struct {
	Skill         string
	Evidence      int64
	Estimate      float64 // θ + δ of the model
	Truth         float64 // θ* + δ*
	Mastered      bool    // the model's mastery
	Band          int32   // the model's mastery band, or the first band the skill is taught on
	TrulyMastered bool    // the truth has mastered it on Band (docs/architecture.md, section 3.5)
	PRef          float64 // P of the model on the reference task of Band, with retention
	PTrueRef      float64 // P* of the truth on it, with R*
	Last          int64   // the model's last practice of it
}

// readAt reads the model at a time, over the learner's modules.
func (r *Run) readAt(at int64) (learnermodel.Readout, error) {
	ids := make([]string, len(r.modules))
	for i, m := range r.modules {
		ids[i] = m.ID
	}
	ro, err := learnermodel.Read(r.state, r.graphs, at, r.grade, ids)
	if err != nil {
		return ro, fmt.Errorf("%w: the model refuses to read at %d: %w", ErrRun, at, err)
	}
	return ro, nil
}

// skillsOf gives what the model and the truth say of the skills with a verdict of a module at a time.
func (r *Run) skillsOf(ro learnermodel.Readout, m *pack.Module, at int64) ([]SkillObs, error) {
	i := slices.IndexFunc(ro.Modules, func(mr learnermodel.ModuleReadout) bool { return mr.Module == m.ID })
	if i < 0 {
		return nil, nil
	}
	var out []SkillObs
	for _, sr := range ro.Modules[i].Skills {
		truth, _, ok := r.l.Truth(sr.Skill, at)
		if !ok {
			continue
		}
		band := r.l.TaughtOn(sr.Skill)
		if sr.Mastered && sr.MasteryBand.Place > 0 {
			band = sr.MasteryBand.Place
		}
		ref := learnermodel.Task{
			ID: m.ID + "/task/reference", Band: learnermodel.BandRef{Place: band, ID: m.Graph.Bands[band-1].ID}, Difficulty: 3, Primary: []string{sr.Skill},
			KeyProvenance: learnermodel.ProvenanceAuthor, SkillTagsProvenance: learnermodel.ProvenanceAuthor,
			Items: []learnermodel.TaskItem{{Place: 1, Slots: []learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 5, HasOptions: true}}}},
		}
		tr, err := ro.Task(ref)
		if err != nil {
			return nil, fmt.Errorf("%w: the reference task of %s: %w", ErrRun, sr.Skill, err)
		}
		out = append(out, SkillObs{
			Skill: sr.Skill, Evidence: sr.EvidenceCount, Estimate: sr.Estimate, Truth: truth,
			Mastered: sr.Mastered, Band: band, TrulyMastered: r.l.TrulyMasteredAt(sr.Skill, band, at),
			PRef: tr.Items[0].Probability, PTrueRef: r.l.ReferenceAt(sr.Skill, band, at), Last: sr.LastPractice,
		})
	}
	return out, nil
}

// observeBefore reads the model and the truth before an attempt on a task at a time: P of each item
// with and without retention, and the skills of its module. seen tells whether the task's solution
// went to the device before; read, when given, is what the model read at that time already.
func (r *Run) observeBefore(m *pack.Module, t *pack.Task, at int64, seen bool, read *learnermodel.Readout) (*Obs, error) {
	if t.Rubric {
		return nil, nil
	}
	var ro learnermodel.Readout
	var err error
	if read != nil {
		ro = *read
	} else if ro, err = r.readAt(at); err != nil {
		return nil, err
	}
	o := &Obs{At: at, Repeat: seen && !t.Candidate.Task.Template, Answers: r.rec.Answers}
	if o.Before, err = r.skillsOf(ro, m, at); err != nil {
		return nil, err
	}
	tr, err := ro.Task(t.Candidate.Task)
	if err != nil {
		return nil, fmt.Errorf("%w: the model reads no P of %s: %w", ErrRun, t.Candidate.Task.ID, err)
	}
	mr := slices.IndexFunc(ro.Modules, func(x learnermodel.ModuleReadout) bool { return x.Module == m.ID })
	many := 0
	for _, it := range t.Candidate.Task.Items {
		many += b2i(len(it.Slots) > 0)
	}
	for k, it := range t.Candidate.Task.Items {
		if len(it.Slots) == 0 {
			continue
		}
		io := ItemObs{PModel: tr.Items[k].Probability, Pause: -1, Help: seen && !t.Candidate.Task.Template, Many: many > 1}
		c := guess(it.Slots)
		rbar, n, last := 0.0, 0, int64(-1)
		for _, id := range slices.Sorted(slices.Values(t.ItemPrimary(k))) {
			if !verdict(m, id) {
				continue
			}
			n++
			j := -1
			if mr >= 0 {
				j = slices.IndexFunc(ro.Modules[mr].Skills, func(s learnermodel.SkillReadout) bool { return s.Skill == id })
			}
			if j < 0 {
				rbar++ // a skill with no practice forgets nothing
				continue
			}
			s := ro.Modules[mr].Skills[j]
			rbar += s.Retention
			if s.EvidenceCount > 0 && s.LastPractice > last {
				last = s.LastPractice
			}
		}
		io.PFresh = io.PModel
		if trial := mr < 0 || !ro.Modules[mr].TrialOver; !trial && n > 0 && rbar > 0 {
			rbar = float64(rbar / float64(n))
			io.PFresh = c + float64(float64(io.PModel-c)/rbar)
		}
		if last >= 0 {
			io.Pause = at - last
		}
		o.Items = append(o.Items, io)
	}
	return o, nil
}

// observeAfter completes an observation with the attempt — the truth's P* of each first answer and
// the help it had — and the skills and the module's level after it.
func (r *Run) observeAfter(m *pack.Module, t *pack.Task, a learner.Attempt, o *Obs) (*Obs, error) {
	if o == nil {
		return nil, nil
	}
	for k, res := range a.Items {
		if k < len(o.Items) {
			o.Items[k].PTrue = res.Chance
			o.Items[k].Help = o.Items[k].Help || a.Technique
		}
	}
	ro, err := r.readAt(a.End)
	if err != nil {
		return nil, err
	}
	if o.After, err = r.skillsOf(ro, m, a.End); err != nil {
		return nil, err
	}
	for _, ms := range r.state.Modules {
		if ms.Module == m.ID {
			o.Level = ms.Level
		}
	}
	if level, err := r.l.LevelAt(m.ID, a.End); err == nil {
		o.TrueLevel = level
	}
	return o, nil
}

// verdict tells whether a skill of a module has a verdict.
func verdict(m *pack.Module, id string) bool {
	return slices.ContainsFunc(m.Graph.Skills, func(s learnermodel.Skill) bool {
		return s.ID == id && s.Assessability == learnermodel.AssessabilityVerdict
	})
}

// guess gives the model's c of an item: the product over its cells of 1/k for a choice of k options
// and 0 for the rest (R03 learner-model).
func guess(slots []learnermodel.Slot) float64 {
	c := 1.0
	for _, s := range slots {
		if s.Check != learnermodel.CheckKindChoice || s.Options <= 0 {
			return 0
		}
		c = float64(c / float64(s.Options))
	}
	return c
}
