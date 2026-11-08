// Package seed gives the bench its numbers: every draw of sim comes from a place of randomness of
// the module format, a seed and names, by the platform's one generator (docs/architecture.md,
// sections 3.4 and 3.8; R27, R59 spec; the text of the format, section 22). spec has the
// generator's test vectors but no Go code of it, so the place is built here, by Go's own PCG, and
// the vectors hold it to the format (R12 in docs/decisions.md).
package seed

import (
	"errors"
	"hash/fnv"
	"math/rand/v2"
	"strconv"
	"strings"
)

// Stream is the stream of draws of one place of randomness. The format's five functions are
// its methods, each Go's own method of math/rand/v2 on the same PCG, drawing the same integers in
// the same order (the text of the format, section 22).
type Stream struct {
	r *rand.Rand
}

// ErrName is a name of a place the format does not take: an empty one, or one with "|", which
// separates the names in the text the high word is hashed from.
var ErrName = errors.New("seed: a name of a place is empty or holds |")

// High gives the high word of the state of the place of seed and names: FNV-1a 64 over the UTF-8
// of the seed in decimal followed, for each name, by "|" and the name. The low word is the seed.
func High(seed uint64, names ...string) (uint64, error) {
	if len(names) == 0 {
		return 0, ErrName
	}
	h := fnv.New64a()
	h.Write([]byte(strconv.FormatUint(seed, 10)))
	for _, name := range names {
		if name == "" || strings.Contains(name, "|") {
			return 0, ErrName
		}
		h.Write([]byte("|" + name))
	}
	return h.Sum64(), nil
}

// Place gives the stream of the place of seed and names, at least one.
func Place(seed uint64, names ...string) (*Stream, error) {
	high, err := High(seed, names...)
	if err != nil {
		return nil, err
	}
	return New(high, seed), nil
}

// New gives the stream of the generator whose state is high·2^64 + low.
func New(high, low uint64) *Stream {
	return &Stream{r: rand.New(rand.NewPCG(high, low))}
}

// Integer draws an integer of 64 bits: the format's integer, Go's Uint64.
func (s *Stream) Integer() uint64 {
	return s.r.Uint64()
}

// Below draws an integer from 0 to bound − 1: the format's below, Go's Uint64N. bound is above 0.
func (s *Stream) Below(bound uint64) uint64 {
	return s.r.Uint64N(bound)
}

// Between draws an integer from a to c, both included: the format's between, a + Int64N(c − a + 1).
// a is not above c, and c − a + 1 fits in an int64.
func (s *Stream) Between(a, c int64) int64 {
	return a + s.r.Int64N(c-a+1)
}

// Fraction draws a number in [0, 1): the format's fraction, Go's Float64, the low 53 bits of an
// integer divided by 2^53. Go's Float64 is inlined, and its product could fuse with what its caller
// does next — 1 − u, or u − 0.5, into one step on arm64 that rounds once —, so the result goes
// through float64(), which the Go specification makes round ("An explicit floating-point type
// conversion rounds to the precision of the target type, preventing fusion"). "Go / fma-diff"
// holds it.
func (s *Stream) Fraction() float64 {
	return float64(s.r.Float64())
}

// Permutation draws an order of 0 … n − 1: the format's permutation, Go's Perm.
func (s *Stream) Permutation(n int) []int {
	return s.r.Perm(n)
}
