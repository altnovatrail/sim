# The fairness of the learner model by the grade: fairness-1

## The criterion

Written down before its first run, by R55 of docs/decisions.md of sim. The populations D-grade1, D-grade3, D-grade5, on olympiad maths, differ by the grade alone: the assumptions of G0, one truth — that of grade 3 — and one unit for the learner n of every grade, so that the model starts on bands 1, 2 and 3. On the module level θ the model gives after each answer, against the learner n of the reference D-grade3: M11.answers is the number of answers before the first of 10 in a row with |θ − θ of the reference| below 0.5, a learner for whom it never holds counting as 200; M11.same is the share of learners for whom that number is at most 50; M11.band is the largest absolute difference, over the bands, between the grade's and the reference's shares of the attempts with a verdict at that band from answer 51 on. The intervals are 95 %, from 2000 resamples of the learners' numbers, the same in every cell. The base holds when, for every grade alike, M11.same is at least 0.90 at the low end of its interval and M11.band at most 0.05 at the high end; M11.answers has no bound.

## The base

- github.com/altnovatrail/learner-model v0.1.8
- github.com/altnovatrail/engine v0.1.6
- sim v0.1.19
- the decisive run: 2000 learners a cell, SHA-256 9e89cb7cfd0f1f0c00bb69c057a273aff0071e11b97450859d7b268e48c07763
- the date: 2026-11-14

## The numbers

| Cell | Learners | M11.answers | M11.same | M11.band | Holds |
|---|---|---|---|---|---|
| D-grade1/novatrail/olympiad-math | 2000 | 32.52 [30.62; 34.36] | 0.812 [0.795; 0.829] | 0.006742 [0.00411; 0.01023] | no: M11.same |
| D-grade5/novatrail/olympiad-math | 2000 | 83.33 [79.72; 87.23] | 0.5525 [0.53; 0.5735] | 0.0127 [0.007317; 0.01804] | no: M11.same |

The base does not hold the criterion.

## For dpia 10.2 spec

On 2026-11-14, sim evaluated the learner model by the grade (R147 spec; R55 sim): virtual learners of one skill under grades 1, 3 and 5 of olympiad maths, 2000 a grade. The base — github.com/altnovatrail/learner-model v0.1.8 and github.com/altnovatrail/engine v0.1.6 — does not hold the criterion written down before the first run. Grade 1: a share of 0.812 [0.795; 0.829] of the learners held within 0.5 of the estimate under grade 3 after at most 50 answers, and 32.52 [30.62; 34.36] answers went by on average, a learner who never held counting as 200; the share of the tasks at a band after the 50th answer differed from grade 3's by at most 0.006742 [0.00411; 0.01023]. Grade 5: a share of 0.5525 [0.53; 0.5735] of the learners held within 0.5 of the estimate under grade 3 after at most 50 answers, and 83.33 [79.72; 87.23] answers went by on average, a learner who never held counting as 200; the share of the tasks at a band after the 50th answer differed from grade 3's by at most 0.0127 [0.007317; 0.01804]. The numbers are given with their 95 % intervals.
