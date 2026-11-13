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

## 2. Claude Code sessions

How `sim` keeps what row 7.10 of `spec`'s COORDINATION.md saved, by the pattern of measurements 1.9 spec (R145 spec; T50–T53 of RUN.md, R45): budgets of size in `docs-check` on every pull request, and a measurement of the sessions at each epic's archive.

### 2.1. Budgets

`docs-check` measures in bytes, in KB of 1,024 (measurements 1.9 spec). The limits are the block `repository` of `spec/tools/docs-check/budgets.json`, with what [docs-budgets.json](../docs-budgets.json) at the root names in its place. What is measured, the largest of each in `sim` on 2026-11-13, after T53, KB:

| What | Limit | The largest |
|---|---|---|
| CLAUDE.md / RUN.md / DEBT.md | 12 / 96 / 128 | 7.4 / 38.9 / 4.1 |
| a rule of `.claude/rules/` | 6 | 1.9, `gate` |
| an entry of the decision log, with its amendments | 90, raised from 40 | 78.6, R04; the next is R31, 13.1, and the median 4.6 |
| the own text under a heading of a document | 24 | 22.0, architecture 3.9 |
| a row of a table of RUN.md / DEBT.md | 4 / 4 | 1.4, В13 / 1.0, Д3 |

- The one raised limit is the entry's: R04 holds the bands and the paths of the criteria, and every new pin adds an amendment with the table of its paths. In Russian it was 92.6 KB; T52 translated it to 78.6. 90 is the largest plus about 15%, as measurements 1.9 spec sets limits. Raising it again is an edit of `docs-budgets.json`, and the third raise since the last analysis is a threshold of R145 spec.
- Before T53 the check, with `docs-budgets.json` at `{}`, found only R04 over its limit: the second half of the log, still Russian, was within 40 KB an entry.
- Architecture 3.9 is 22.0 of 24 KB: the next growth of that section moves a part of it into a section of its own rather than raising the limit.

### 2.2. The measurement of sessions

`just -f spec/justfile tokens-report --since … --until … --row` prints one row of the table below. It runs at the archive of each epic, for the epic's window (RUN.md, «Как устроен план: эпики»), and once a month from `spec` (measurements 1.9 spec).

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-sim` — before, the row of measurements 1.9 spec | 2026-10-13T00:00 … 2026-11-12T14:27 | 10 | 3,275 / 1,195 | 325.5 / 79.5 | 6 | $397.59 |
| `-workspaces-sim` — T50–T53, row 7.10 of COORDINATION.md `spec` | 2026-11-12T20:41 … 2026-11-13T05:18 | 7 | 564 / 454 | 180.3 / 105.0 | 2 | $60.56 |

- The context of a main call fell from 325.5 to 180.3 thousand, −45%; by session it was 178.0 in T50, 153.4 in T51, 128.8 in T52 and 109.5 in T53, while the documents a session reads became English, one task at a time. A subagent's rose from 79.5 to 105.0: here the subagents were translators and blind readers, each of which reads the glossary, concepts 15 spec, and the rules of translation. Three of the seven sessions, of one, one and five calls, ran no task.
- The windows differ in work: E2–E5 with code, runs and judgements against the move to a tag and the translation of `docs/`. The comparison by sessions after the task is T124 spec.
