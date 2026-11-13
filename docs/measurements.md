# sim — measurements

> The results of `sim`'s measurements: what was measured, how, on what, and what numbers came out. [architecture.md](architecture.md) refers here and holds only what was decided from them (RULES.md, The bar; R143 spec). A measurement is cited by section, like any document: `measurements 1.1`, and from another repository `measurements 1.1 sim`.
>
> An external fact a measurement rests on is a record of [sources.md](sources.md). The measurements of section 1 stood in architecture 3, 3.8 and 3.9 until T51 ([docs/plan/E5.md](plan/E5.md)); what was decided from them stays there.

## 1. The time of runs

The time of the gate and of a decisive run: by it the budget of the job "Gate" and the sizes of runs were chosen (architecture 3.8 and 3.9; R26, R31, R33).

### 1.1. "Gate" on runners

- A measurement on a runner, pull request #33 (RUN.md, T24): building the stage — 20–25 seconds, the run of the gate — about four minutes on `ubuntu-24.04` and two and a half on `ubuntu-24.04-arm`, almost the same with the Go cache and without it.
- A measurement on a runner — the pull request of T30 (RUN.md). "Gate" took: pull request #40 (T30) — 11 min 52 s on `ubuntu-24.04` and 7 min 3 s on `ubuntu-24.04-arm`; #41 (T37) — 6 min 54 s and 6 min 49 s; #42 (T36) — 8 min 14 s and 4 min 36 s; #43 (T31) — 8 min 23 s and 4 min 42 s. On arm64 the gate is faster after T36, and on x64 the runner's time wanders from run to run more than the change gave — the measurement shows no gain there.
- On `engine` v0.1.6, pull request #49 (T41), "Gate" on 100 learners took 5 min 39 s on `ubuntu-24.04` with a cold Go cache and 4 min 27 s with a warm one, and 3 min 13 s and 3 min 7 s on `ubuntu-24.04-arm`.
- On 150, "Gate" took 8 min 21 s with a cold cache and 6 min 5 s with a warm one on `ubuntu-24.04`, and 4 min 16 s on `ubuntu-24.04-arm`. This is a measurement on `modules` v0.2.0 and `learner-model` v0.1.4.
- On `modules` v0.2.5 and `learner-model` v0.1.6, pull request #52 and the push of its merge, "Gate" on 150 learners took 6 min 17 s and 5 min 0 s on `ubuntu-24.04` and 3 min 36 s and 3 min 32 s on `ubuntu-24.04-arm`, both times with a cold Go cache. Recomputing the slowest measurement gave 300 learners about 11 min 59 s, right up against 12, and 250 about 10 min 5 s (the 2026-11-12 amendment to R31).
- On 250, pull request #53, "Gate" took 7 min 41 s on `ubuntu-24.04` and 5 min 28 s on `ubuntu-24.04-arm` with a cold cache, and on `ubuntu-24.04` with a warm one 5 min 34 s and, on a slow runner, 9 min 38 s: for callers without a cache on such a runner the gate comes close to 12 minutes, and their CI will show their time on the new tag (the 2026-11-12 amendment to R31).

### 1.2. The gate in the development container

- The first code of the gate, on 300 learners: in the development container on 20 cores the gate takes 45 seconds, 11 minutes of CPU time with the repeat for `I1`; on the 4 cores of a runner this is about three minutes of run beyond building the stage.
- On `olympiad-math` 0.1.1 the gate on 300 learners took 8.5 minutes and 135 CPU minutes on 20 cores — about 35–40 minutes on a runner (R31).
- In the development container the gate on 150 learners took 1 min 28 s and 23.6 CPU minutes in T41, and in T43 7 min and 28.4 CPU minutes on a host that neighbouring containers loaded (R42), on `learner-model` v0.1.6 1 min 16 s and 18.7 CPU minutes (R43), and on 250 learners 1 min 52 s and 29.6 CPU minutes at a load average of about 40–60 from neighbouring containers (the 2026-11-12 amendment to R31).

### 1.3. The decisive run

The recipe `just decisive`: processes by parts of the learners, on all the cores of the machine, then `sim merge` (architecture 3.8).
- On the 20 cores of the development container it took about 40 minutes for the base and the ceiling together on `olympiad-math` 0.1.0, and on 0.1.1 and 4,000 learners the base took 6 h 43 min and the ceiling 1 h 31 min (the 2026-11-09 amendment to R04, "a new pin").
- On 1,000 learners, with the choice made anew for the first 10 learners of a cell and with a process on one CPU and `GOGC` 400 — about 37 CPU minutes a process for the base and about 14 for the ceiling, so on a free machine the run takes less than an hour; on a machine that other containers loaded, 1 h 42 min and 41 min (R33; the 2026-11-10 amendment to R04).
- On `engine` v0.1.6 and 1,500 learners — 38 min by the clock for the base, about 33 CPU minutes a process, and about 19 CPU minutes for the ceiling (the 2026-11-11 amendment to R33 and to R04).
- On the packages of `modules` v0.2.5 and `learner-model` v0.1.5 — 47 min 48 s for the base, about 38.7 CPU minutes a process, and about 18.8 for the ceiling, on a host shared with neighbouring containers (the 2026-11-11 amendment to R04, "packages with `time` and leading questions"; R42).
- On `learner-model` v0.1.6 — 62 min 24 s for the base at a load average of up to 46, about 25.9 CPU minutes a process, a third less, and about 18.0 for the ceiling (the 2026-11-11 amendment to R04, "the base on `learner-model` v0.1.6"; R43).
- On 2,000 learners a subsample on a free machine promises about 38 minutes for the base and about 27 for the ceiling, and on a host shared with neighbouring containers the base took 62 min 37 s, about 31.3 CPU minutes a process, and the ceiling 30 min 37 s, about 24.1 (the 2026-11-12 amendment to R33 and to R04).
- The files of its results — about 0.8–1.0 GB a rule (R28).
