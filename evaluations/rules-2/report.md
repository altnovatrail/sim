# The adult's correction against the truth: rules-2

## The criterion

Written down before the run, by R73 of docs/decisions.md of sim. C-grade1/novatrail/olympiad-math is D-grade1/novatrail/olympiad-math — grade 1, the truth of grade 3 and the unit of D-grade3 — with an adult who sets the module level to the band of the truth after the 10th answer. M13.answers is the number of answers before the first of 10 in a row with |θ − θ*| below 0.5, θ the model's module level and θ* the truth's after each answer, a learner for whom it never holds counting as 200; M13.same is the share of learners for whom that number is at most 50. The correction holds when the paired difference of M13.same, C-grade1/novatrail/olympiad-math's less D-grade1/novatrail/olympiad-math's, is above zero at the low end of its 95 % interval, from 2000 resamples of the learners' numbers.

## The run

- github.com/altnovatrail/learner-model v0.1.13-0.20261116235821-df338e58e35f and github.com/altnovatrail/engine v0.1.6, 2000 learners a cell, SHA-256 552eaff7f98c8c792e3a49e3c899288e240cfebf4222278ce4baf3a6e92723c0
- sim v0.1.28-0.20261117003236-2095e69d2658
- the date: 2026-11-17

## The numbers

2000 of the 2000 learners of C-grade1/novatrail/olympiad-math were corrected.

| Number | C-grade1/novatrail/olympiad-math | D-grade1/novatrail/olympiad-math | Difference |
|---|---|---|---|
| M13.answers | 76.49 [73.16; 79.77] | 119.7 [116.1; 123.5] | -43.24 [-47.32; -38.87] |
| M13.same | 0.56 [0.54; 0.581] | 0.351 [0.3295; 0.372] | 0.209 [0.185; 0.233] |

The correction holds the criterion.

## For dpia 10.2 spec

On 2026-11-17, sim checked the adult's correction of R157 spec — github.com/altnovatrail/learner-model v0.1.13-0.20261116235821-df338e58e35f and github.com/altnovatrail/engine v0.1.6 — on 2000 virtual learners of grade 1 of olympiad maths whose truth is the band of grade 3, by a criterion written down before the run (R73 sim). With an adult who set the level to the band of the truth after the 10th answer, the model's estimate came within 0.5 of the truth for 10 answers in a row by the 50th answer for a share of 0.56 [0.54; 0.581] of the learners, against 0.351 [0.3295; 0.372] of the same learners without the correction, a difference of 0.209 [0.185; 0.233]. The correction holds the criterion. The numbers are given with their 95 % intervals.
