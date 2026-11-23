package measure

import (
	"math"
	"slices"
	"testing"

	"github.com/altnovatrail/sim/internal/learner"
	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/published"
	"github.com/altnovatrail/sim/internal/run"
	"github.com/altnovatrail/sim/internal/seed"
	"github.com/altnovatrail/sim/internal/synthetic"
)

// attempt is an observed attempt of so many items, after so many answers, at a band, with the level
// after it.
type attempt struct {
	answers, items int
	band           int32
	level          float64
}

// recordOf gives a record of observed attempts, and a step of a task with a rubric between them, which
// has no observation.
func recordOf(as ...attempt) *run.Record {
	rec := &run.Record{}
	for _, a := range as {
		o := &run.Obs{Answers: a.answers, Band: a.band, Level: a.level, Items: make([]run.ItemObs, a.items)}
		items := make([]learner.ItemResult, a.items)
		rec.Steps = append(rec.Steps,
			run.Step{Attempts: []learner.Attempt{{Items: items}}, Obs: []*run.Obs{o}},
			run.Step{Attempts: []learner.Attempt{{}}, Obs: []*run.Obs{nil}})
	}
	return rec
}

// The trace gives each answer the level after its attempt, up to the 200th, and counts by band the
// attempts whose first answer is the 51st or later.
func TestTraced(t *testing.T) {
	var as []attempt
	for k := 0; k < 210; k += 3 {
		as = append(as, attempt{answers: k, items: 3, band: int32(1 + k%2), level: float64(k)})
	}
	tr := Traced(recordOf(as...))
	if len(tr.Levels) != FairLength {
		t.Fatalf("%d levels, want %d", len(tr.Levels), FairLength)
	}
	for k, l := range tr.Levels {
		if want := float64(k - k%3); l != want {
			t.Fatalf("answer %d has the level %v, want %v, that after its attempt", k+1, l, want)
		}
	}
	// The attempts from answer 51 on start at 51, 54, … 207: 53 of them, the odd starts at band 2.
	if !slices.Equal(tr.Bands, []int{26, 27}) {
		t.Errorf("the bands after 50 answers: %v", tr.Bands)
	}
	if empty := Traced(&run.Record{}); empty.Levels == nil || empty.Bands == nil || len(empty.Levels)+len(empty.Bands) != 0 {
		t.Errorf("a run with no answer: %+v", empty)
	}
}

// levels gives a trace of 200 levels: off by gap until the answer from, then the reference's, but
// the answers of flicker off again.
func levels(gap float64, from int, flicker ...int) *Trace {
	t := &Trace{}
	for k := range FairLength {
		x := 0.0
		if k < from || slices.Contains(flicker, k) {
			x = gap
		}
		t.Levels = append(t.Levels, x)
	}
	return t
}

// The pair counts the answers before the first of ten in a row within 0.5, as M3 counts: none
// before them gives 0, a series broken starts again, and a learner who never holds counts as 200;
// M11.same holds to 50 answers and not to 51.
func TestFair(t *testing.T) {
	ref := levels(0, 0)
	for _, c := range []struct {
		name    string
		grade   *Trace
		answers float64
		same    bool
	}{
		{"the same from the first answer", levels(0, 0), 0, true},
		{"within the gap", levels(0.49, FairLength), 0, true},
		{"at the gap", levels(0.5, FairLength), FairLength, false},
		{"from answer 31", levels(2.5, 30), 30, true},
		{"from answer 51", levels(2.5, 50), 50, true},
		{"from answer 52", levels(2.5, 51), 51, false},
		{"a series broken at its ninth", levels(2.5, 30, 38), 39, true},
		{"never", levels(-2.5, FairLength), FairLength, false},
		{"the last ten", levels(1, FairLength-FairAnswers), FairLength - FairAnswers, false},
		{"fewer than ten", levels(1, FairLength-FairAnswers+1), FairLength, false},
	} {
		s := Fair(c.grade, ref)
		if s["M11.pairs"] != 1 || s["M11.answers.sum"] != c.answers || (s["M11.same"] == 1) != c.same {
			t.Errorf("%s: %v answers, same %v; want %v, %v", c.name, s["M11.answers.sum"], s["M11.same"], c.answers, c.same)
		}
	}
	s := Fair(&Trace{Levels: ref.Levels, Bands: []int{1, 3}}, &Trace{Levels: ref.Levels, Bands: []int{2, 0, 2}})
	for key, want := range map[string]float64{
		"M11.band.grade.1": 1, "M11.band.grade.2": 3, "M11.band.grade.3": 0, "M11.band.grade": 4,
		"M11.band.reference.1": 2, "M11.band.reference.2": 0, "M11.band.reference.3": 2, "M11.band.reference": 4,
	} {
		if v, ok := s[key]; !ok || v != want {
			t.Errorf("%s = %v, %v; want %v", key, v, ok, want)
		}
	}
}

// The numbers of M11 read off a cell of pairs: the mean answers, the share the same, and the largest
// difference of the shares of a band; M11.band is not read where a side has no attempt after 50
// answers.
func TestFairNumbers(t *testing.T) {
	ref := levels(0, 0)
	var pairs []Sums
	for k := range 40 {
		g := levels(2.5, 10*(k%10))
		g.Bands, ref.Bands = []int{1, 3, 0}, []int{2, 2, 0}
		pairs = append(pairs, Fair(g, ref))
	}
	vs := Read(population.Reference+"/"+population.Olympiad, []string{population.Olympiad}, pairs, FairNumbers())
	want := map[string]float64{"M11.answers": 45, "M11.same": 0.6, "M11.band": 0.25}
	for _, v := range vs {
		if !v.Read || math.Abs(v.Value-want[v.Name]) > 1e-12 || v.Low > v.Value || v.High < v.Value {
			t.Errorf("%s = %v [%v; %v], read %v; want %v", v.Name, v.Value, v.Low, v.High, v.Read, want[v.Name])
		}
	}
	none := Read("cell", nil, []Sums{Fair(levels(0, 0), levels(0, 0))}, FairNumbers())
	if none[2].Read || none[0].Value != 0 || none[1].Value != 1 {
		t.Errorf("a pair with no attempt after 50 answers: band read %v, answers %v, same %v", none[2].Read, none[0].Value, none[1].Value)
	}
}

// On runs of the synthetic sheet, which has three bands of grades 1–6, the learner n of grade 1 and
// of grade 3 is one truth on one unit: the reference against itself holds from the first answer, and
// grade 1, whose model starts a band below, starts apart.
func TestFairOnRuns(t *testing.T) {
	ms, err := synthetic.Modules()
	if err != nil {
		t.Fatal(err)
	}
	mods := []*pack.Module{ms[synthetic.Sheet]}
	trace := func(name string, n int) (*Trace, *run.Record) {
		p, _ := population.Named(name)
		l, err := learner.New(seed.Working, p.Units(), n, mods, p.Assumptions(n))
		if err != nil {
			t.Fatal(err)
		}
		r, err := run.New(l, mods, p.Limit)
		if err != nil {
			t.Fatal(err)
		}
		rec, err := r.All()
		if err != nil {
			t.Fatal(err)
		}
		return Traced(rec), rec
	}
	for n := range 3 {
		ref, rr := trace(population.Reference, n)
		low, lr := trace("D-grade1", n)
		if rr.Unit != lr.Unit || len(ref.Levels) != FairLength || len(low.Levels) != FairLength {
			t.Fatalf("learner %d: units %s and %s, %d and %d levels", n, rr.Unit, lr.Unit, len(ref.Levels), len(low.Levels))
		}
		if s := Fair(ref, ref); s["M11.answers.sum"] != 0 || s["M11.same"] != 1 {
			t.Errorf("learner %d: the reference against itself %v", n, s)
		}
		if math.Abs(low.Levels[0]-ref.Levels[0]) < FairGap {
			t.Errorf("learner %d: grade 1 starts at %v and grade 3 at %v", n, low.Levels[0], ref.Levels[0])
		}
	}
}

// The trace of a learner of C-grade1 carries the answers before the adult's correction, and is
// that of the same learner of D-grade1 up to them; M12 reads the share of learners whose level the
// guard lowered, each counted once.
func TestRulesOnRuns(t *testing.T) {
	ms, err := published.Modules()
	if err != nil {
		t.Fatal(err)
	}
	mods := []*pack.Module{ms[population.Olympiad]}
	trace := func(name string) *Trace {
		p, _ := population.Named(name)
		l, err := learner.New(seed.Working, p.Units(), 0, mods, p.Assumptions(0))
		if err != nil {
			t.Fatal(err)
		}
		r, err := run.New(l, mods, p.Limit)
		if err != nil {
			t.Fatal(err)
		}
		rec, err := r.All()
		if err != nil {
			t.Fatal(err)
		}
		return Traced(rec)
	}
	low, corrected := trace("D-grade1"), trace("C-grade1")
	if low.Corrected != 0 || corrected.Corrected < population.CorrectAfter {
		t.Fatalf("corrected after %d and %d answers", low.Corrected, corrected.Corrected)
	}
	if len(corrected.Truth) != len(corrected.Levels) || corrected.Truth[0] == 0 {
		t.Errorf("the truth beside %d levels: %d, from %v", len(corrected.Levels), len(corrected.Truth), corrected.Truth[:1])
	}
	k := corrected.Corrected
	if !slices.Equal(low.Levels[:k], corrected.Levels[:k]) || slices.Equal(low.Levels[k:], corrected.Levels[k:]) {
		t.Errorf("the levels of grade 1 with the correction after %d answers: %v, without: %v", k, corrected.Levels[:k+1], low.Levels[:k+1])
	}
	var sums []Sums
	for _, g := range []int{0, 2, 0, 1} {
		sums = append(sums, Guarded(&Trace{Guard: g}))
	}
	if v := Read("D-grade3/x", nil, sums, GuardNumbers())[0]; !v.Read || v.Value != 0.5 || v.Low > v.Value || v.High < v.Value {
		t.Errorf("M12.false of two learners of four fired = %+v", v)
	}
}

// M13 counts the answers before the model's level holds within FairGap of the truth's for
// FairAnswers in a row, FairLength when it never does, and the share of learners by FairBy; the
// trace carries the truth beside each level.
func TestCaught(t *testing.T) {
	trace := func(from int) *Trace {
		tr := &Trace{}
		for k := range FairLength {
			x := 2.0
			if k >= from {
				x = 0.1
			}
			tr.Levels, tr.Truth = append(tr.Levels, x), append(tr.Truth, 0)
		}
		return tr
	}
	for _, c := range []struct {
		from    int
		k, same float64
	}{{0, 0, 1}, {FairBy, FairBy, 1}, {FairBy + 1, FairBy + 1, 0}, {FairLength, FairLength, 0}} {
		s := Caught(trace(c.from))
		if s["M13.answers.sum"] != c.k || s["M13.same"] != c.same || s["M13.learners"] != 1 {
			t.Errorf("caught from answer %d: %v", c.from, s)
		}
	}
	v := Read("D-grade3/x", nil, []Sums{Caught(trace(0)), Caught(trace(FairLength))}, TruthNumbers())
	if v[0].Value != FairLength/2 || v[1].Value != 0.5 {
		t.Errorf("M13 of two learners = %+v", v)
	}
}
