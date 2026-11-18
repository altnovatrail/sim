package criterion

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/seed"
)

// The criterion of the short run of the three rules of R157 spec, written down before the run
// (R72): the rules are learner-model at RulesCandidate; grade 5's M11.same holds above the base's
// at the low end of its interval; M12.false of the reference's cell holds at most guardCeiling at
// the high end of its interval; and the adult's correction, the cell of population.Rules, has the
// paired difference of M11.same against grade 1's above zero at its low end. fair is the constraint
// of R57 as it is.
const (
	RulesCandidate = "df338e58e35f97c51eca49be0064f19d0fcfc0a2"
	guardCeiling   = 0.05
	// betterCell is the cell whose M11.same the rules must raise, and correctionCell and
	// uncorrectedCell the cells of the adult's correction and of the same learners without it.
	betterCell      = "D-grade5/" + population.Olympiad
	correctionCell  = "C-grade1/" + population.Olympiad
	uncorrectedCell = "D-grade1/" + population.Olympiad
	referenceCell   = population.Reference + "/" + population.Olympiad
)

// Rules is an evaluation of the three rules: its name, as its folder evaluations/<name>/, the run of
// the rules and the base's run with their SHA-256, the date, the version of sim that reads them and
// the libraries of the gate file, which the engine of the run is held to.
type Rules struct {
	Name, Date, Version   string
	Run, Base             *bench.File
	RunDigest, BaseDigest string
	Pins                  map[string]Pin
}

// RulesResult is what the evaluation of the three rules reads: each of the four conditions of R72,
// its numbers and whether it holds, and whether all hold.
type RulesResult struct {
	Rules
	// Better is M11.same of grade 5 under the rules, BaseSame the same of the base's run.
	Better, BaseSame measure.Value
	// Grades are M11 of grades 1 and 5 under the rules and under the base, each number by name.
	Grades []RulesGrade
	// Fair are the checks of the constraint fair, the rules against the base.
	Fair []Check
	// Guard is M12.false of the reference's cell, of Learners learners.
	Guard    measure.Value
	Learners int
	// Correction is M11 of the cell of the correction against the same learners without it, each
	// number by name; Corrected is how many of its CorrectionLearners the adult corrected.
	Correction                    []measure.Pair
	Corrected, CorrectionLearners int
	// Holds are the four conditions in their order, and All whether every one holds.
	Holds [4]bool
	All   bool
}

// RulesGrade is M11 of one grade's cell under the rules and under the base.
type RulesGrade struct {
	Cell        string
	Rules, Base []measure.Value
}

// EvaluateRules reads the evaluation of the three rules (R72): the run of the rules — under the
// engine, on the working learners, a whole run on the pinned packages with its invariants and no
// breach, learner-model at RulesCandidate and the engine at the gate file's pin, the cells of
// population.Fairness and of population.Rules of at least a decisive run's learners each — against
// the base's run, the file paths.json names.
func EvaluateRules(r Rules) (*RulesResult, error) {
	if r.Name == "" || r.Date == "" || r.RunDigest == "" || r.BaseDigest == "" {
		return nil, fmt.Errorf("%w: an evaluation of the rules needs its name, its date and the SHA-256 of both runs", ErrCriterion)
	}
	paths, err := fixed()
	if err != nil {
		return nil, err
	}
	if r.BaseDigest != paths.Base {
		return nil, fmt.Errorf("%w: the base's run has the SHA-256 %s, and the paths name the base %s", ErrCriterion, r.BaseDigest, paths.Base)
	}
	for _, side := range []struct {
		name string
		f    *bench.File
	}{{"the run of the rules", r.Run}, {"the base's run", r.Base}} {
		if err := whole(side.name, side.f); err != nil {
			return nil, err
		}
	}
	if err := onCandidate("the run of the rules", r.Run, r.Pins); err != nil {
		return nil, err
	}
	run, err := bench.Fair(r.Run)
	if err != nil {
		return nil, fmt.Errorf("%w: the run of the rules: %w", ErrCriterion, err)
	}
	base, err := bench.Fair(r.Base)
	if err != nil {
		return nil, fmt.Errorf("%w: the base's run: %w", ErrCriterion, err)
	}
	cells := map[string]*bench.FairCell{}
	for _, name := range []string{betterCell, uncorrectedCell, correctionCell} {
		c := fairCell(run, name)
		switch {
		case c == nil:
			return nil, fmt.Errorf("%w: the run of the rules has no cell %s beside %s", ErrCriterion, name, population.Reference)
		case len(c.Pairs) < decisive:
			return nil, fmt.Errorf("%w: cell %s has %d learners, and the evaluation reads %d", ErrCriterion, name, len(c.Pairs), decisive)
		}
		cells[name] = c
	}
	for _, name := range []string{betterCell, uncorrectedCell} {
		if b := fairCell(base, name); b == nil || !slices.Equal(b.N, cells[name].N) {
			return nil, fmt.Errorf("%w: the base's cell %s is not of the learners of the run of the rules", ErrCriterion, name)
		}
	}
	if !slices.Equal(cells[correctionCell].N, cells[uncorrectedCell].N) {
		return nil, fmt.Errorf("%w: cell %s is not of the learners of %s", ErrCriterion, correctionCell, uncorrectedCell)
	}
	out := &RulesResult{Rules: r}
	for _, name := range []string{uncorrectedCell, betterCell} {
		c, cb := cells[name], fairCell(base, name)
		out.Grades = append(out.Grades, RulesGrade{Cell: name, Rules: measure.Read(c.Reference, c.Graph, c.Pairs, measure.FairNumbers()), Base: measure.Read(cb.Reference, cb.Graph, cb.Pairs, measure.FairNumbers())})
	}
	out.Better, out.BaseSame = out.Grades[1].Rules[1], out.Grades[1].Base[1]
	out.Holds[0] = out.Better.Read && out.BaseSame.Read && out.Better.Low > out.BaseSame.Value
	if out.Fair, err = fair(base, run); err != nil {
		return nil, err
	}
	out.Holds[1] = !slices.ContainsFunc(out.Fair, func(k Check) bool { return !k.Met })
	var guard []measure.Sums
	for _, l := range r.Run.Lines {
		if l.Population == population.Reference {
			guard = append(guard, measure.Guarded(l.Fair))
		}
	}
	out.Guard, out.Learners = measure.Read(referenceCell, cells[betterCell].Graph, guard, measure.GuardNumbers())[0], len(guard)
	out.Holds[2] = out.Guard.Read && out.Guard.High <= guardCeiling
	c, u := cells[correctionCell], cells[uncorrectedCell]
	out.CorrectionLearners = len(c.Pairs)
	if out.Correction, err = measure.ReadPair(c.Reference, c.Graph, c.Pairs, u.Pairs, measure.FairNumbers()); err != nil {
		return nil, err
	}
	for _, l := range r.Run.Lines {
		if l.Population+"/"+strings.Join(l.Graph, "+") == correctionCell && l.Fair.Corrected > 0 {
			out.Corrected++
		}
	}
	same := out.Correction[1].Diff
	out.Holds[3] = same.Read && same.Low > 0
	out.All = !slices.Contains(out.Holds[:], false)
	return out, nil
}

// whole holds a run the evaluation of the rules reads to what an evaluation reads: under the engine,
// on the working learners, a whole run on the pinned packages, its invariants recorded and none
// broken.
func whole(name string, f *bench.File) error {
	h := f.Header
	switch {
	case h.Rule != "engine" || h.Set != seed.Working:
		return fmt.Errorf("%w: %s is under %q on %s, and the evaluation reads the engine on the working learners", ErrCriterion, name, h.Rule, h.Set.Experiment)
	case h.Part != (bench.Part{}) || len(h.Packages) > 0:
		return fmt.Errorf("%w: %s is not a whole run on the pinned packages", ErrCriterion, name)
	case !h.Invariants:
		return fmt.Errorf("%w: %s did not record the invariants", ErrCriterion, name)
	}
	for _, l := range f.Lines {
		if len(l.Breaches) > 0 {
			return fmt.Errorf("%w: %s breaks an invariant: %s", ErrCriterion, name, l.Breaches[0])
		}
	}
	return nil
}

// Write writes the report of the evaluation of the rules in Markdown: the criterion, the two runs,
// the four conditions with their numbers and intervals, the outcome, and the paragraph for dpia 10.2
// spec (R72).
func (r *RulesResult) Write(w io.Writer) error {
	var b strings.Builder
	yes := func(ok bool) string {
		if ok {
			return "yes"
		}
		return "no"
	}
	fmt.Fprintf(&b, "# The short run of the three rules of R157 spec: %s\n\n", r.Name)
	b.WriteString("## The criterion\n\n")
	fmt.Fprintf(&b, "Written down before the run, by R72 of docs/decisions.md of sim. The rules — the start a band below the band of the grade, the adult's correction and the guard — are checked when all four hold, on the cells %s, %s, %s and %s of olympiad maths, the learner n one learner in all of them and in the base's run: 1. M11.same of %s, by the %dth answer against %s, above the base's at the low end of its interval; 2. the constraint fair of R57 on grades 1 and 5, against the base; 3. M12.false, the share of the learners of %s whose module level the guard lowered at least once, at most %.2f at the high end of its interval; 4. the adult's correction, %s — %s with the level set to the band of the truth after the %dth answer —, with M11.same above that of %s at the low end of the interval of the paired difference. The intervals are 95 %%, from %d resamples of the learners' numbers.\n\n",
		uncorrectedCell, referenceCell, betterCell, correctionCell, betterCell, measure.FairBy, population.Reference, referenceCell, guardCeiling, correctionCell, uncorrectedCell, population.CorrectAfter, uncorrectedCell, measure.Resamples)
	b.WriteString("## The runs\n\n")
	for _, side := range []struct {
		name, digest string
		f            *bench.File
	}{{"the rules", r.RunDigest, r.Run}, {"the base", r.BaseDigest, r.Base}} {
		fmt.Fprintf(&b, "- %s: %s, %d learners a cell, SHA-256 %s\n", side.name, libraries(side.f.Header), side.f.Header.Learners, side.digest)
	}
	fmt.Fprintf(&b, "- sim %s\n- the date: %s\n\n", r.Version, r.Date)
	b.WriteString("## The numbers\n\n| Condition | The rules | Against | Holds |\n|---|---|---|---|\n")
	fmt.Fprintf(&b, "| 1. M11.same of %s | %s | the base %s | %s |\n", betterCell, showInterval(r.Better), showInterval(r.BaseSame), yes(r.Holds[0]))
	failed := 0
	for _, k := range r.Fair {
		failed += b2i(!k.Met)
	}
	fmt.Fprintf(&b, "| 2. fair, %d checks | %d failed | the base | %s |\n", len(r.Fair), failed, yes(r.Holds[1]))
	fmt.Fprintf(&b, "| 3. M12.false of %s, %d learners | %s | at most %.2f | %s |\n", referenceCell, r.Learners, showInterval(r.Guard), guardCeiling, yes(r.Holds[2]))
	fmt.Fprintf(&b, "| 4. M11.same of %s less %s | %s | above 0 | %s |\n\n", correctionCell, uncorrectedCell, showInterval(r.Correction[1].Diff), yes(r.Holds[3]))
	b.WriteString("M11 of the grades against grade 3:\n\n| Cell | Side | M11.answers | M11.same | M11.band |\n|---|---|---|---|---|\n")
	for _, g := range r.Grades {
		for _, side := range []struct {
			name string
			vs   []measure.Value
		}{{"the rules", g.Rules}, {"the base", g.Base}} {
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", g.Cell, side.name, showInterval(side.vs[0]), showInterval(side.vs[1]), showInterval(side.vs[2]))
		}
	}
	b.WriteString("\n")
	b.WriteString("The checks of fair:\n\n")
	for _, k := range r.Fair {
		fmt.Fprintf(&b, "- %s %s: %s, %s\n", k.Cell, k.Number, yes(k.Met), k.Detail)
	}
	fmt.Fprintf(&b, "\nThe cell of the correction, %d of whose %d learners the adult corrected:\n\n| Number | %s | %s | Difference |\n|---|---|---|---|\n", r.Corrected, r.CorrectionLearners, correctionCell, uncorrectedCell)
	for _, p := range r.Correction {
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", p.A.Name, showInterval(p.A), showInterval(p.B), showInterval(p.Diff))
	}
	verdict := "are checked: all four conditions hold"
	if !r.All {
		verdict = "are not checked: not every condition holds"
	}
	fmt.Fprintf(&b, "\nThe rules %s.\n\n## For dpia 10.2 spec\n\n", verdict)
	fmt.Fprintf(&b, "On %s, sim ran the three rules of R157 spec — %s — on virtual learners of one skill under grades 1, 3 and 5 of olympiad maths, %d a grade, against the base — %s —, by a criterion written down before the run (R72 sim). Grade 5 came within %.1f of the estimate under grade 3 after at most %d answers for a share of %s of the learners, against %s under the base, and grade 1 for %s, against %s. The constraint fair held %d of its %d checks of grades 1 and 5. The guard lowered the level of a share of %s of the learners of grade 3, whose truth is the band of their grade. An adult's correction to the band of the truth after the %dth answer changed the share of grade 1 that came to the estimate under grade 3 by %s. The rules %s. The numbers are given with their 95 %% intervals.\n",
		r.Date, libraries(r.Run.Header), r.Learners, libraries(r.Base.Header), measure.FairGap, measure.FairBy, showInterval(r.Better), showInterval(r.BaseSame), showInterval(r.Grades[0].Rules[1]), showInterval(r.Grades[0].Base[1]), len(r.Fair)-failed, len(r.Fair), showInterval(r.Guard), population.CorrectAfter, showInterval(r.Correction[1].Diff), verdict)
	_, err := io.WriteString(w, b.String())
	return err
}
