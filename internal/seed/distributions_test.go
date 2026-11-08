package seed

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// inverseNormal at two points known to 16 digits, and everywhere else held to the distribution
// function itself: Φ(x) = erfc(−x/√2)/2, by math.Erfc, which is not AS 241, gives back p from
// p = 1e-300 to 1 − 1e-12 within what an error of 1e-14 in x makes of it. In the tail dp/p is
// about |x|·dx, so the tolerance on p is 1e-14 · p · max(1, |x|).
func TestInverseNormal(t *testing.T) {
	if got := inverseNormal(0.5); got != 0 {
		t.Errorf("inverseNormal(0.5) = %v", got)
	}
	if got := inverseNormal(0.975); math.Abs(got-1.959963984540054) > 1e-15 {
		t.Errorf("inverseNormal(0.975) = %.17g", got)
	}
	ps := []float64{1e-300, 1e-100, 1e-20, 1e-10, 1e-5, 0.001, 0.02425, 0.07, 0.075, 0.08, 0.2, 0.3, 0.45}
	for _, p := range ps {
		for _, q := range []float64{p, 1 - p} {
			if q >= 1-1e-12 {
				continue
			}
			x := inverseNormal(q)
			phi := float64(math.Erfc(float64(-x/math.Sqrt2)) / 2)
			// Near 1 the distribution function is flat, and p is held from the other side.
			got, want := phi, q
			if q > 0.5 {
				got, want = float64(math.Erfc(float64(x/math.Sqrt2))/2), 1-q
			}
			tol := float64(1e-14 * want)
			if ax := math.Abs(x); ax > 1 {
				tol = float64(tol * ax)
			}
			if math.Abs(got-want) > tol {
				t.Errorf("inverseNormal(%v) = %.17g, and Φ of it is %.17g", q, x, phi)
			}
		}
	}
}

// A fraction of 0 is taken as 2^−53, so Normal stays finite, and the exponential of it is +0.
func TestNormalOfZero(t *testing.T) {
	if x := normalAt(0); x != inverseNormal(0x1p-53) || math.IsInf(x, 0) || math.IsNaN(x) || x > -8 {
		t.Errorf("normalAt(0) = %v", x)
	}
	if lo, hi := normalAt(0), inverseNormal(1-0x1p-53); lo != -hi {
		t.Errorf("the tails end at %v and %v", lo, hi)
	}
	if y := exponentialAt(0); math.Float64bits(y) != 0 {
		t.Errorf("exponentialAt(0) = %v, bits %x", y, math.Float64bits(y))
	}
}

// The draws of a stream, many of them, have the mean and the spread the distributions have.
func TestDistributionMoments(t *testing.T) {
	const n = 200000
	s, err := Working.Place("moments", "test")
	if err != nil {
		t.Fatal(err)
	}
	var sumN, sumN2, sumE, sumE2 float64
	for range n {
		x := s.Normal()
		sumN += x
		sumN2 += float64(x * x)
		y := s.Exponential()
		if y < 0 {
			t.Fatalf("Exponential gave %v", y)
		}
		sumE += y
		sumE2 += float64(y * y)
	}
	meanN, meanE := float64(sumN/n), float64(sumE/n)
	varN, varE := float64(sumN2/n)-float64(meanN*meanN), float64(sumE2/n)-float64(meanE*meanE)
	// Five standard errors: the mean of n draws, and the variance, whose error is about √(2/n) for
	// the normal and √(8/n) for the exponential.
	if math.Abs(meanN) > 5/math.Sqrt(n) || math.Abs(varN-1) > float64(5*math.Sqrt(2.0/n)) {
		t.Errorf("normal: mean %v, variance %v", meanN, varN)
	}
	if math.Abs(meanE-1) > 5/math.Sqrt(n) || math.Abs(varE-1) > float64(5*math.Sqrt(8.0/n)) {
		t.Errorf("exponential: mean %v, variance %v", meanE, varE)
	}
}

// The golden bits: the first draws of each distribution at one place, written as the hex of their
// bits in testdata/distributions.txt. The same file passes on every build the shared go.yml runs
// the tests on — GOAMD64 v1 and v3, 386, js/wasm and arm64 — so the distributions give the same
// bits everywhere (docs/architecture.md, section 3.4). SIM_UPDATE=1 writes the file again; a
// change to it is read on review.
func TestGoldenBits(t *testing.T) {
	s, err := Working.Place("golden", "distributions")
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	for i := range 32 {
		fmt.Fprintf(&b, "normal %d %016x\n", i, math.Float64bits(s.Normal()))
	}
	for i := range 32 {
		fmt.Fprintf(&b, "exponential %d %016x\n", i, math.Float64bits(s.Exponential()))
	}
	got := b.String()
	file := filepath.Join("testdata", "distributions.txt")
	if os.Getenv("SIM_UPDATE") == "1" {
		if err := os.WriteFile(file, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		gs, ws := bufio.NewScanner(strings.NewReader(got)), bufio.NewScanner(strings.NewReader(string(want)))
		for gs.Scan() && ws.Scan() {
			if gs.Text() != ws.Text() {
				t.Fatalf("the bits part from %s: got %q, want %q", file, gs.Text(), ws.Text())
			}
		}
		t.Fatalf("the bits part from %s", file)
	}
}
