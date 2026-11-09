package measure

import (
	"reflect"
	"testing"

	"github.com/altnovatrail/sim/internal/run"
)

// The oracle is the ceiling of the model where it is one: no error of the estimate, no false
// mastery, no wait for one, and calibration no worse than the model's; on the same answers.
func TestOracle(t *testing.T) {
	recs, graph := records(t, "G0", 60)
	var model, oracle []Sums
	for _, r := range recs {
		model, oracle = append(model, Model(r)), append(oracle, Oracle(r))
	}
	pairs, err := ReadPair("G0", graph, oracle, model, ModelNumbers())
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]Pair{}
	for _, p := range pairs {
		byName[p.Diff.Name] = p
	}
	for _, name := range []string{"M1.rms.200", "M1.mean.200", "M4.false.all", "M6.late", "M7.error.10", "M3.lag"} {
		if p := byName[name]; !p.A.Read || p.A.Value != 0 {
			t.Errorf("the oracle's %s = %+v, want 0", name, p.A)
		}
	}
	if p := byName["M2.brier.all"]; !p.Diff.Read || p.Diff.Value >= 0 || p.Diff.High >= 0 {
		t.Errorf("the oracle's Brier is no better than the model's: %+v", p.Diff)
	}
	a := readOracle(t, "A-linked")
	if v := value(t, a, "M9.hasFalse"); v.Read && v.Value != 0 {
		t.Errorf("the oracle reads an anchor falsely: %+v", v)
	}
}

func readOracle(t *testing.T, name string) []Value {
	t.Helper()
	recs, graph := records(t, name, 6)
	var sums []Sums
	for _, r := range recs {
		sums = append(sums, Oracle(r))
	}
	return Read(name, graph, sums, ModelNumbers())
}

// A cell against itself differs by nothing, to the last bit; two cells of different learners are
// refused.
func TestReadPair(t *testing.T) {
	sums, graph := cell(t, "G0", 8)
	pairs, err := ReadPair("G0", graph, sums, sums, ModelNumbers())
	if err != nil {
		t.Fatal(err)
	}
	read := 0
	for _, p := range pairs {
		if !reflect.DeepEqual(p.A, p.B) || (p.Diff.Read && (p.Diff.Value != 0 || p.Diff.Low != 0 || p.Diff.High != 0)) {
			t.Errorf("%s: %+v", p.Diff.Name, p)
		}
		read += b2i(p.Diff.Read)
	}
	if read == 0 {
		t.Error("no number read")
	}
	if _, err := ReadPair("G0", graph, sums, sums[1:], ModelNumbers()); err == nil {
		t.Error("cells of different learners are compared")
	}
}

// The standard error of a value is the deviation of its values over the resamples, and its ends
// the 51st from each end of them; M2 cuts the answers after pauses of at least 7 days apart.
func TestSE(t *testing.T) {
	learners := []Sums{{"x": 1, "n": 1}, {"x": 3, "n": 1}, {"x": 5, "n": 1}, {"x": 7, "n": 1}}
	v := Read("se", nil, learners, []Number{{"mean", ratio("x", "n")}})[0]
	if !v.Read || len(v.Resampled) != Resamples || v.SE <= 0.5 || v.SE >= 2 {
		t.Errorf("SE %v of the mean of 1, 3, 5, 7 over %d resamples", v.SE, len(v.Resampled))
	}
	vals := make([]float64, Resamples)
	for i := range vals {
		vals[i] = float64(Resamples - i)
	}
	if low, high := Ends(vals); low != 51 || high != Resamples-50 {
		t.Errorf("Ends = %v, %v", low, high)
	}
	s := Sums{}
	rec := []*run.Obs{{Items: []run.ItemObs{{PModel: 0.5, Pause: 7 * day}, {PModel: 0.5, Pause: 6 * day}}}}
	m2(s, rec, [][]bool{{true, false}})
	if s["M2.after7.n"] != 1 || s["M2.all.n"] != 2 {
		t.Errorf("after 7 days %v of %v answers", s["M2.after7.n"], s["M2.all.n"])
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
