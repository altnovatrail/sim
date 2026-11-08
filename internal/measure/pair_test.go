package measure

import (
	"testing"
)

// The oracle is the ceiling of the model where it is one: no error of the estimate, no false
// mastery, no wait for one, and calibration no worse than the model's; on the same answers.
func TestOracle(t *testing.T) {
	recs, graph := records(t, "G0", 20)
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
		if p.A != p.B || (p.Diff.Read && (p.Diff.Value != 0 || p.Diff.Low != 0 || p.Diff.High != 0)) {
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

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
