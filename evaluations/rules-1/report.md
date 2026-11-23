# The short run of the learner model's three rules: rules-1

## The criterion

Written down before the run, in docs/decisions.md of sim. The rules — the start a band below the band of the grade, the adult's correction and the guard — are checked when all four hold, on the cells D-grade1/novatrail/olympiad-math, D-grade3/novatrail/olympiad-math, D-grade5/novatrail/olympiad-math and C-grade1/novatrail/olympiad-math of olympiad maths, the learner n one learner in all of them and in the base's run: 1. M11.same of D-grade5/novatrail/olympiad-math, by the 50th answer against D-grade3, above the base's at the low end of its interval; 2. the constraint fair on grades 1 and 5, against the base; 3. M12.false, the share of the learners of D-grade3/novatrail/olympiad-math whose module level the guard lowered at least once, at most 0.05 at the high end of its interval; 4. the adult's correction, C-grade1/novatrail/olympiad-math — D-grade1/novatrail/olympiad-math with the level set to the band of the truth after the 10th answer —, with M11.same above that of D-grade1/novatrail/olympiad-math at the low end of the interval of the paired difference. The intervals are 95 %, from 2000 resamples of the learners' numbers.

## The runs

- the rules: github.com/altnovatrail/learner-model v0.1.13-0.20261116235821-df338e58e35f and github.com/altnovatrail/engine v0.1.6, 2000 learners a cell, SHA-256 8b1d9dcdd33de9216f7dec883e9bad6a0ef5ed63827f4d5ed44a405f876b2c16
- the base: github.com/altnovatrail/learner-model v0.1.8 and github.com/altnovatrail/engine v0.1.6, 2000 learners a cell, SHA-256 9e89cb7cfd0f1f0c00bb69c057a273aff0071e11b97450859d7b268e48c07763
- sim v0.1.28-0.20261117002128-502c93a90aaa
- the date: 2026-11-17

## The numbers

| Condition | The rules | Against | Holds |
|---|---|---|---|
| 1. M11.same of D-grade5/novatrail/olympiad-math | 0.812 [0.795; 0.829] | the base 0.5525 [0.53; 0.5735] | yes |
| 2. fair, 6 checks | 1 failed | the base | no |
| 3. M12.false of D-grade3/novatrail/olympiad-math, 2000 learners | 0.0005 [0; 0.0015] | at most 0.05 | yes |
| 4. M11.same of C-grade1/novatrail/olympiad-math less D-grade1/novatrail/olympiad-math | -0.497 [-0.519; -0.475] | above 0 | no |

M11 of the grades against grade 3:

| Cell | Side | M11.answers | M11.same | M11.band |
|---|---|---|---|---|
| D-grade1/novatrail/olympiad-math | the rules | 17.1 [15.76; 18.44] | 0.909 [0.8965; 0.9215] | 0.01773 [0.0146; 0.021] |
| D-grade1/novatrail/olympiad-math | the base | 32.52 [30.62; 34.36] | 0.812 [0.795; 0.829] | 0.006742 [0.00411; 0.01023] |
| D-grade5/novatrail/olympiad-math | the rules | 32.52 [30.62; 34.36] | 0.812 [0.795; 0.829] | 0.006615 [0.003917; 0.01007] |
| D-grade5/novatrail/olympiad-math | the base | 83.33 [79.72; 87.23] | 0.5525 [0.53; 0.5735] | 0.0127 [0.007317; 0.01804] |

The checks of fair:

- D-grade1/novatrail/olympiad-math M11.answers: yes, worse by at most -13.25, the tolerance 4.96
- D-grade1/novatrail/olympiad-math M11.same: yes, worse by at most -0.0755, the tolerance 0.04789
- D-grade1/novatrail/olympiad-math M11.band: no, worse by at most 0.01423, the tolerance 0.01
- D-grade5/novatrail/olympiad-math M11.answers: yes, worse by at most -46.61, the tolerance 9.489
- D-grade5/novatrail/olympiad-math M11.same: yes, worse by at most -0.2315, the tolerance 0.06339
- D-grade5/novatrail/olympiad-math M11.band: yes, worse by at most 0.0004043, the tolerance 0.0131

The cell of the correction, 2000 of whose 2000 learners the adult corrected:

| Number | C-grade1/novatrail/olympiad-math | D-grade1/novatrail/olympiad-math | Difference |
|---|---|---|---|
| M11.answers | 102 [98.35; 105.6] | 17.1 [15.76; 18.44] | 84.86 [80.99; 88.74] |
| M11.same | 0.412 [0.3905; 0.4335] | 0.909 [0.8965; 0.9215] | -0.497 [-0.519; -0.475] |
| M11.band | 0.04753 [0.04231; 0.0529] | 0.01773 [0.0146; 0.021] | 0.0298 [0.02438; 0.0346] |

The rules are not checked: not every condition holds.

## For spec's DPIA

On 2026-11-17, sim ran the learner model's three rules — github.com/altnovatrail/learner-model v0.1.13-0.20261116235821-df338e58e35f and github.com/altnovatrail/engine v0.1.6 — on virtual learners of one skill under grades 1, 3 and 5 of olympiad maths, 2000 a grade, against the base — github.com/altnovatrail/learner-model v0.1.8 and github.com/altnovatrail/engine v0.1.6 —, by a criterion written down before the run. Grade 5 came within 0.5 of the estimate under grade 3 after at most 50 answers for a share of 0.812 [0.795; 0.829] of the learners, against 0.5525 [0.53; 0.5735] under the base, and grade 1 for 0.909 [0.8965; 0.9215], against 0.812 [0.795; 0.829]. The constraint fair held 5 of its 6 checks of grades 1 and 5. The guard lowered the level of a share of 0.0005 [0; 0.0015] of the learners of grade 3, whose truth is the band of their grade. An adult's correction to the band of the truth after the 10th answer changed the share of grade 1 that came to the estimate under grade 3 by -0.497 [-0.519; -0.475]. The rules are not checked: not every condition holds. The numbers are given with their 95 % intervals.
