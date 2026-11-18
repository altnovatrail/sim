package bench

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// A learner of a population of fairness runs on its reference's unit and carries the trace of M11
// under the engine alone; a learner of any other population has none (R55).
func TestFairLine(t *testing.T) {
	ms, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	graph := []string{synthetic.Sheet}
	mods := []*pack.Module{ms[synthetic.Sheet]}
	grade1, _ := population.Named("D-grade1")
	grade1.Graphs = [][]string{graph}
	other, _ := population.Named("S-predicate")
	c := Config{Learners: 1, Set: seed.Working}
	rec, err := runOne(c, grade1, mods, 4)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Unit != population.Reference+"/4" {
		t.Errorf("the learner of D-grade1 runs on %s", rec.Unit)
	}
	l, err := one(c, grade1, graph, mods, 4)
	if err != nil {
		t.Fatal(err)
	}
	if l.Population != "D-grade1" || l.Fair == nil || len(l.Fair.Levels) != measure.FairLength {
		t.Errorf("the line of D-grade1: %s, trace %+v", l.Population, l.Fair)
	}
	if l, err := one(c, other, graph, mods, 4); err != nil || l.Fair != nil {
		t.Errorf("S-predicate has a trace: %+v, %v", l.Fair, err)
	}
	c.Rule = run.Ceiling
	if l, err := one(c, grade1, graph, mods, 4); err != nil || l.Fair != nil {
		t.Errorf("a trace under the ceiling: %+v, %v", l.Fair, err)
	}
	for _, name := range []string{"D-grade1", population.Reference} {
		if err := Run(Config{Populations: []string{name}, Learners: 1, Set: seed.Working, Rule: run.Ceiling}, &strings.Builder{}); !errors.Is(err, ErrBench) {
			t.Errorf("the ceiling of %s: %v", name, err)
		}
	}
	all, ceiling := Populations(run.Engine), Populations(run.Ceiling)
	if len(all) != len(population.All()) || len(ceiling) != len(all)-len(population.Fairness) || slices.Contains(ceiling, population.Reference) {
		t.Errorf("all the populations: %d under the engine, %d under the ceiling", len(all), len(ceiling))
	}
}

// fairFile builds a file of results under a rule: two learners of S-predicate, and under the engine
// two of each population of fairness on olympiad maths, their traces apart by the grade's offset
// until the answer from.
func fairFile(rule string, from map[string]int) *File {
	f := &File{Header: Header{Set: seed.Working, Rule: rule}}
	for n := range 2 {
		f.Lines = append(f.Lines, Line{Population: "S-predicate", Graph: []string{synthetic.Sheet}, N: n, Sums: measure.Sums{"learners": 1}})
	}
	if rule != "engine" {
		return f
	}
	for _, p := range population.Fairness {
		for n := range 2 {
			tr := &measure.Trace{Bands: []int{n, 1}}
			for k := range measure.FairLength {
				x := 0.0
				if k < from[p.Name] {
					x = 2.5
				}
				tr.Levels = append(tr.Levels, x)
			}
			f.Lines = append(f.Lines, Line{Population: p.Name, Graph: []string{population.Olympiad}, N: n, Sums: measure.Sums{"learners": 1}, Fair: tr})
		}
	}
	return f
}

// The cells of fairness read in pairs, each grade's against the reference's, the resamples by the
// reference's cell; against the ceiling, which does not run them, their lines are left out of both
// sides; a grade without its reference, a learner without a trace or out of step is refused.
func TestFair(t *testing.T) {
	f := fairFile("engine", map[string]int{"D-grade1": 30, "D-grade5": 60})
	cells, err := Fair(f)
	if err != nil {
		t.Fatal(err)
	}
	ref := population.Reference + "/" + population.Olympiad
	if len(cells) != 2 || cells[0].Cell != "D-grade1/"+population.Olympiad || cells[1].Population != "D-grade5" || cells[0].Reference != ref || cells[1].Reference != ref {
		t.Fatalf("the cells of fairness: %+v", cells)
	}
	for _, c := range cells {
		want := map[string]float64{"D-grade1": 30, "D-grade5": 60}[c.Population]
		if len(c.Pairs) != 2 || c.N[1] != 1 || c.Pairs[1]["M11.answers.sum"] != want || c.Pairs[0]["M11.band.grade.2"] != 1 {
			t.Errorf("%s: %v", c.Cell, c.Pairs)
		}
	}
	if cs, err := Fair(fairFile("ceiling", nil)); err != nil || len(cs) != 0 {
		t.Errorf("a file with no cell of fairness: %v, %v", cs, err)
	}
	pairs, err := Cells(Side{File: f}, Side{File: fairFile("ceiling", nil)})
	if err != nil || len(pairs) != 1 || pairs[0].Population != "S-predicate" {
		t.Errorf("the base against the ceiling: %+v, %v", pairs, err)
	}
	if pairs, err := Cells(Side{File: f}, Side{File: f}); err != nil || len(pairs) != 4 {
		t.Errorf("the base against itself: %d cells, %v", len(pairs), err)
	}
	drop := func(pop string, n int) *File {
		g := fairFile("engine", nil)
		for k, l := range g.Lines {
			if l.Population == pop && l.N == n {
				g.Lines = append(g.Lines[:k], g.Lines[k+1:]...)
				break
			}
		}
		return g
	}
	noReference := fairFile("engine", nil)
	noReference.Lines = slices.DeleteFunc(noReference.Lines, func(l Line) bool { return l.Population == population.Reference })
	untraced := fairFile("engine", nil)
	untraced.Lines[len(untraced.Lines)-1].Fair = nil
	swapped := fairFile("engine", nil)
	for k, l := range swapped.Lines {
		if l.Population == population.Reference {
			swapped.Lines[k].N = 1 - l.N
		}
	}
	for name, g := range map[string]*File{
		"no reference":        noReference,
		"a reference short":   drop(population.Reference, 0),
		"a learner short":     drop("D-grade5", 1),
		"a line no trace":     untraced,
		"learners not paired": swapped,
	} {
		if _, err := Fair(g); !errors.Is(err, ErrBench) {
			t.Errorf("%s: %v", name, err)
		}
	}
}
