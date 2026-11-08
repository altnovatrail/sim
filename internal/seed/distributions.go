package seed

import (
	"math"

	"github.com/altnovatrail/learner-model/detmath"
)

// The format gives no normal or exponential distribution, and sim answers for theirs: each is
// built here on one fraction of the stream, with arithmetic whose every product goes straight into
// float64() and with detmath.Log and, for the normal, math.Sqrt, which give the same bits on every
// build (R09; R09 learner-model). So the bands of the gate, read to the last digit, are the same on any
// processor (docs/architecture.md, sections 3.4 and 3.7; R12 in docs/decisions.md).

// Normal draws a number of the standard normal distribution: the inverse of its distribution
// function at one fraction, by Wichura's algorithm AS 241, PPND16, which is good to about 1e-16.
// A fraction of exactly 0 — one draw in 2^53 — is taken as 2^−53, the smallest one above it, so
// the result is finite and the two tails end at the same distance: the largest 1 − u is 2^−53 too.
func (s *Stream) Normal() float64 {
	return normalAt(s.Fraction())
}

// normalAt gives the normal number of one fraction u, 0 ≤ u < 1.
func normalAt(u float64) float64 {
	if u == 0 {
		u = 0x1p-53
	}
	return inverseNormal(u)
}

// Exponential draws a number of the exponential distribution with mean 1: −ln(1 − u) at one
// fraction u. 1 − u is above 0, since u is below 1.
func (s *Stream) Exponential() float64 {
	return exponentialAt(s.Fraction())
}

// exponentialAt gives the exponential number of one fraction u, 0 ≤ u < 1. It is 0 − ln(1 − u)
// rather than −ln(1 − u), so that u = 0 gives +0 and not −0.
func exponentialAt(u float64) float64 {
	return 0 - detmath.Log(1-u)
}

// The coefficients of AS 241, PPND16 (Wichura, M. J., "Algorithm AS 241: The percentage points of
// the normal distribution", Applied Statistics 37(3), 1988, pp. 477–484), lowest power first.
var (
	centralNum = []float64{
		3.3871328727963666080e0, 1.3314166789178437745e+2, 1.9715909503065514427e+3,
		1.3731693765509461125e+4, 4.5921953931549871457e+4, 6.7265770927008700853e+4,
		3.3430575583588128105e+4, 2.5090809287301226727e+3,
	}
	centralDen = []float64{
		1, 4.2313330701600911252e+1, 6.8718700749205790830e+2, 5.3941960214247511077e+3,
		2.1213794301586595867e+4, 3.9307895800092710610e+4, 2.8729085735721942674e+4,
		5.2264952788528545610e+3,
	}
	nearNum = []float64{
		1.42343711074968357734e0, 4.63033784615654529590e0, 5.76949722146069140550e0,
		3.64784832476320460504e0, 1.27045825245236838258e0, 2.41780725177450611770e-1,
		2.27238449892691845833e-2, 7.74545014278341407640e-4,
	}
	nearDen = []float64{
		1, 2.05319162663775882187e0, 1.67638483018380384940e0, 6.89767334985100004550e-1,
		1.48103976427480074590e-1, 1.51986665636164571966e-2, 5.47593808499534494600e-4,
		1.05075007164441684324e-9,
	}
	tailNum = []float64{
		6.65790464350110377720e0, 5.46378491116411436990e0, 1.78482653991729133580e0,
		2.96560571828504891230e-1, 2.65321895265761230930e-2, 1.24266094738807843860e-3,
		2.71155556874348757815e-5, 2.01033439929228813265e-7,
	}
	tailDen = []float64{
		1, 5.99832206555887937690e-1, 1.36929880922735805310e-1, 1.48753612908506148525e-2,
		7.86869131145613259100e-4, 1.84631831751005468180e-5, 1.42151175831644588870e-7,
		2.04426310338993978564e-15,
	}
)

// inverseNormal gives the x at which the standard normal distribution function is p, 0 < p < 1.
func inverseNormal(p float64) float64 {
	q := p - 0.5
	if math.Abs(q) <= 0.425 {
		r := 0.180625 - float64(q*q)
		return float64(q * float64(horner(centralNum, r)/horner(centralDen, r)))
	}
	r := p
	if q > 0 {
		r = 1 - p
	}
	r = math.Sqrt(-detmath.Log(r))
	var x float64
	if r <= 5 {
		r -= 1.6
		x = float64(horner(nearNum, r) / horner(nearDen, r))
	} else {
		r -= 5
		x = float64(horner(tailNum, r) / horner(tailDen, r))
	}
	if q < 0 {
		return -x
	}
	return x
}

// horner gives the polynomial of coefficients c, lowest power first, at r.
func horner(c []float64, r float64) float64 {
	v := c[len(c)-1]
	for i := len(c) - 2; i >= 0; i-- {
		v = float64(v*r) + c[i]
	}
	return v
}
