# sim — measurements

> The results of `sim`'s measurements: what was measured, how, on what, and what numbers came out. [architecture.md](architecture.md) refers here and holds only what was decided from them. A measurement is cited by its number: `measurements 1.1`, and from another repository `measurements 1.1 sim`.
>
> A fact no run checks that a measurement rests on is a record of [sources.md](sources.md).

## 1. The time of runs

The time of the gate and of a decisive run: by it the budget of the job "Gate" and the sizes of runs were chosen .

### 1.1. "Gate" on runners

- A measurement on a runner, pull request #33: building the stage — 20–25 seconds, the run of the gate — about four minutes on `ubuntu-24.04` and two and a half on `ubuntu-24.04-arm`, almost the same with the Go cache and without it.
- A measurement on a runner — the pull request of the pin of `olympiad-math` 0.1.1. "Gate" took: pull request #40, that pin — 11 min 52 s on `ubuntu-24.04` and 7 min 3 s on `ubuntu-24.04-arm`; #41, the move to `spec` v0.12.0 — 6 min 54 s and 6 min 49 s; #42, runs of no more than an hour — 8 min 14 s and 4 min 36 s; #43, the decisive run on the candidate's commit — 8 min 23 s and 4 min 42 s. On arm64 the gate is faster after #42, and on x64 the runner's time wanders from run to run more than the change gave — the measurement shows no gain there.
- On `engine` v0.1.6, pull request #49, "Gate" on 100 learners took 5 min 39 s on `ubuntu-24.04` with a cold Go cache and 4 min 27 s with a warm one, and 3 min 13 s and 3 min 7 s on `ubuntu-24.04-arm`.
- On 150, "Gate" took 8 min 21 s with a cold cache and 6 min 5 s with a warm one on `ubuntu-24.04`, and 4 min 16 s on `ubuntu-24.04-arm`. This is a measurement on `modules` v0.2.0 and `learner-model` v0.1.4.
- On `modules` v0.2.5 and `learner-model` v0.1.6, pull request #52 and the push of its merge, "Gate" on 150 learners took 6 min 17 s and 5 min 0 s on `ubuntu-24.04` and 3 min 36 s and 3 min 32 s on `ubuntu-24.04-arm`, both times with a cold Go cache. Recomputing the slowest measurement gave 300 learners about 11 min 59 s, right up against 12, and 250 about 10 min 5 s.
- On 250, pull request #53, "Gate" took 7 min 41 s on `ubuntu-24.04` and 5 min 28 s on `ubuntu-24.04-arm` with a cold cache, and on `ubuntu-24.04` with a warm one 5 min 34 s and, on a slow runner, 9 min 38 s: for callers without a cache on such a runner the gate comes close to 12 minutes, and their CI will show their time on the new tag.
- The gate as jobs, 2026-11-16, the times of the jobs from the GitHub API. On `main` before it, run 43805864156: "Gate / amd64" 502 s, its step `gate` 458 s, and "Gate / arm64" 362 s and 328 s; "Go / targets" 450 s, "Go / arm64" 462 s. On pull request #80, four parts, run 43812443762: on amd64 a part 135–168 s, its step `gate-part` 94–123 s, and the check 47 s, its `gate-check` 13 s; on arm64 121–129 s and 46 s; "Go / targets" 429 s, "Go / arm64" 456 s. Two parts, run 43813867051: on amd64 the parts 351 s and 232 s, `gate-part` 319 s and 206 s on two runners, and the check 53 s; on arm64 205 s and 202 s and 49 s, and "Go / targets" 702 s, the time of an x64 runner wandering from run to run as before. The image of `toolchain` came from the cache in every job, 12–31 s, and the Go cache was cold: a pull request saves none, and the key was new. The part jobs started up to 5 minutes after the run in the run of four, the runners of the organization being busy. `internal/measure`, whose `TestEngineOnRuns` runs by processes, took in "Go / targets" 21.4 s, 22.0 s and 29.8 s at v3, v1 and 386, against 27.0 s, 27.2 s and 38.8 s on `main`, 19.2 s against 28.5 s on arm64 in "Go / arm64", and the same in js/wasm and wasip1, 103–138 s. In the container, at a load average of 40–70 from neighbouring containers, `just gate exact` took 4 min 1 s, four `just gate-part` 84, 62, 85 and 77 s, and `just gate-check exact` 10 s, with the table of `just gate exact` to the byte.
- `.github` v0.1.30, 2026-11-16, the times of the jobs from the GitHub API, each from the start of the run where so said. On `main` before it, run 43816538542: the run 8 min 4 s; "Go / targets" 443 s and over at 8 min 4 s, "Go / arm64" 456 s, "Go / fma-diff" 143 s, "CodeQL / analyze (go)" 81 s; the gate's slowest part 200 s on amd64 and 127 s on arm64, the check 35 s and 43 s, "Gate / amd64" over at 5 min 5 s. On pull request #81, run 43855954103, the Go caches of `go.yml` cold — a pull request saves none, and `main` has none yet —: the run 7 min 38 s; the jobs "Go / target (…)" amd64-v3 73 s, cpu-off 71 s, 386 111 s, js-wasm 245 s and wasip1 166 s, started 1.7–2.8 min after the run, the runners of the organization busy, and "Go / targets" over at 6 min 42 s; "Go / arm64" 452 s, the longest job and the end of the run; "Go / fma-diff" 85 s; "CodeQL / analyze (go)" 61 s, run 43855954129. The gate in four parts: on amd64 the parts 130–187 s and the check 48 s, over at 4 min 13 s; on arm64 115–117 s and 40 s, over at 3 min 52 s. On x64 every job took the image from `main`'s entry of the new file, which `documents.yml` had saved already; on arm64 there was none, and each job built the image. The run of the last commit, 43858561575, waited for the runners of the organization: jobs started up to 7 min after the run, "Go / target (js-wasm)" at 7 min 19 s, and the run took 11 min 52 s, though no job of it was longer than those above — "Go / arm64" 438 s, the targets 77–165 s, the gate's parts 114–196 s.
- `spec` v0.18.4 and `.github` v0.1.31, 2026-11-16, the times of the jobs from the GitHub API, each end from the start of the run. On `main` before it, run 43909340309: the run 8 min 18 s; "Go / toolchain" 6 s, the image from the cache; "Go / lint" 45 s, "Go / vuln" 46 s, "Go / licenses" 46 s, "Go / test" 91 s; the targets 69–203 s, "Go / targets" over at 497 s; "Go / arm64" 429 s; "Go / fma-diff" 73 s; "Documents / docs-check-repo" 41 s; the gate's parts 127–194 s on amd64 and 119–124 s on arm64, the checks 50 s and 39 s, "Gate / amd64" over at 302 s; the release of that commit, run 43910863766, its step "Work out the version" 22 s. On pull request #83, run 43913407672, the Dockerfile new, so every entry of the image and of the Go caches missed: the run 9 min 52 s, its jobs started 1–4 min late, the runners of the organization busy; "Go / toolchain" 81 s, its step "Toolchain image" 45 s building `toolchain`, and the core no entry of its own — the step "Toolchain stages" named none, `sim`'s `toolchain` adding nothing to it; "Go / lint" 124 s, "Go / vuln" 82 s, "Go / licenses" 81 s, "Go / test" 85 s, the cold Go caches; the targets 77–211 s, "Go / targets" over at 592 s; "Go / arm64" 452 s; "Go / fma-diff" 138 s; "Documents / docs-check-repo" 47 s, its build of the image 25 s; the gate's parts 131–156 s on amd64 and 117–121 s on arm64, each building the image, 27 s, the checks 58 s and 40 s; "CodeQL / analyze (go)" 58 s. The release of `toolchain-core` is measured by the merge's push to `main`, not here.

### 1.2. The gate in the development container

- The first code of the gate, on 300 learners: in the development container on 20 cores the gate takes 45 seconds, 11 minutes of CPU time with the repeat for `I1`; on the 4 cores of a runner this is about three minutes of run beyond building the stage.
- On `olympiad-math` 0.1.1 the gate on 300 learners took 8.5 minutes and 135 CPU minutes on 20 cores — about 35–40 minutes on a runner.
- In the development container the gate on 150 learners took 1 min 28 s and 23.6 CPU minutes on `engine` v0.1.6, and on the pin of `modules` v0.2.5 7 min and 28.4 CPU minutes on a host that neighbouring containers loaded, on `learner-model` v0.1.6 1 min 16 s and 18.7 CPU minutes, and on 250 learners 1 min 52 s and 29.6 CPU minutes at a load average of about 40–60 from neighbouring containers.

### 1.3. The decisive run

The recipe `just decisive`: processes by parts of the learners, on all the cores of the machine, then `sim merge`.
- On the 20 cores of the development container it took about 40 minutes for the base and the ceiling together on `olympiad-math` 0.1.0, and on 0.1.1 and 4,000 learners the base took 6 h 43 min and the ceiling 1 h 31 min.
- On 1,000 learners, with the choice made anew for the first 10 learners of a cell and with a process on one CPU and `GOGC` 400 — about 37 CPU minutes a process for the base and about 14 for the ceiling, so on a free machine the run takes less than an hour; on a machine that other containers loaded, 1 h 42 min and 41 min.
- On `engine` v0.1.6 and 1,500 learners — 38 min by the clock for the base, about 33 CPU minutes a process, and about 19 CPU minutes for the ceiling.
- On the packages of `modules` v0.2.5 and `learner-model` v0.1.5 — 47 min 48 s for the base, about 38.7 CPU minutes a process, and about 18.8 for the ceiling, on a host shared with neighbouring containers.
- On `learner-model` v0.1.6 — 62 min 24 s for the base at a load average of up to 46, about 25.9 CPU minutes a process, a third less, and about 18.0 for the ceiling.
- On 2,000 learners a subsample on a free machine promises about 38 minutes for the base and about 27 for the ceiling, and on a host shared with neighbouring containers the base took 62 min 37 s, about 31.3 CPU minutes a process, and the ceiling 30 min 37 s, about 24.1.
- In the judgement of `theta`, on 2,000 learners a cell: the candidate `theta` on the working learners took 8 h 10 min by the clock for 687.6 CPU minutes, about 34.4 a process — the host was busy or asleep, from a load average of 1.4 at the start to 44.6 at the end, and the CPU is that of a run of about 40 minutes. On a host that neighbouring containers loaded — a load average of 11 to 42 — the base on `held-out-model-1` took 54 min 11 s, 673.5 CPU minutes, and `theta` on it 38 min 5 s, 641.0. The base on `theta` by `go.mod` took 38 min 15 s, 639.8 CPU minutes, about 32.0 a process, and the ceiling 31 min 18 s, 483.3, about 24.2.
- On the release of `theta`, `learner-model` v0.1.8, on a host that neighbouring containers loaded — a load average of 7 to 28 —: the base took 38 min 27 s by the clock, 638.9 CPU minutes, about 31.9 a process, and the ceiling 30 min 17 s, 482.3, about 24.1.
- With the call of the engine at the end of each task segment of a lesson with a verdict, the same bits as on the release of `theta`: the base took 48 min 32 s by the clock and the ceiling, which calls nothing there, 38 min 40 s, on a host with a load average of 12 at the start and 31 at the end. The ceiling is slower than on the release of `theta` by as much as the base, about a quarter, so the host, not the step, accounts for most of it; what the step costs is not measured apart.
- With the sums of the note `E6` in each line of the base, the same bits as with the call of the engine past them: the base took 46 min 32 s by the clock, 727.2 CPU minutes of user time, about 36.4 a process, on a host with a load average of 2.3 at the start, about 38 a quarter of an hour in and 9.7 at the end. The ceiling was not run: it writes no `E6`, and its file is that of the run with the call of the engine.
- In the judgement of the candidate of consolidation in a lesson, `engine` 0e4680f, on 2,000 working learners a cell, with the call of the engine at the end of each task segment: 43 min 50 s by the clock, 725.6 CPU minutes of user time, about 36.3 a process, on a machine free at the start — a load average of 2.0; neighbours raised it to as much as 61.6 between the fifth and the twentieth minute, and the rest of the run held at 22–23, its own 20 processes. It is under the 45 minutes a decisive run is held to with the step, so no cut of the run is needed. On a free machine the base before the step was promised about 38 minutes by the subsample on 2,000 learners above; the candidate consolidates about once a pass, so this is not the step's cost apart.
- The cost of the three cells of fairness on olympiad maths, measured without them, on cells of olympiad maths alone: `sim run` of the base on one process, `GOGC` 400 and `GOMAXPROCS` 1, on a free machine — a load average of 1.4–2.3. `G4-band` took 5.8 s of user time on 20 learners, 29.8 s on 200 and 59.2 s on 400, and `S-adult` 30.5 s on 200: about 0.14 s a learner, beyond about 3 s of start. So one cell of olympiad maths on 2,000 learners is about 4.7 CPU minutes, and three about 14, 1.9 % of the 727.2 of the base with the sums of `E6`. The base with the call of the engine has not been run on a wholly free machine. Its 727.2 CPU minutes are 36.4 a process, 1.14 times the 31.9 of the release of `theta`, which took 38 min 27 s by the clock; scaled so and by 1.019, 68 cells take about 44–45 minutes, and the candidate of consolidation, 43 min 50 s, would take about 44 min 40 s. So the run with the cells is at the 45 minutes a decisive run is held to, not under them with room to spare. This is an estimate, not a run of the new populations; the base with the three cells measures it.
- The base with the three cells of fairness, 68 cells and 2,000 working learners a cell, from a clean clone at `v0.1.19`: 43 min 28 s by the clock, 747.4 CPU minutes of user time, about 37.4 a process, and 17.8 of system time. The run waited until the host's load average fell below 5 and started at 4.8; its own 20 processes then held it at about 20–22, neighbours raised it to as much as 27.9 between the third and the twentieth minute, and it was 16.0 at the end. Against the estimate above, the processor time is 1.028 times the 727.2 of the base with the sums of `E6`, about 0.9 % beyond the 1.019 the three cells were to add, and the clock is under 45 minutes, as for the candidate of consolidation on a machine as free at its start, so no cut is called for; on a wholly free machine the run is still not measured. This is the first base with the call of the engine on a machine free at its start; part 7 of 20 run again by the same build on one process took 12 min 50 s, and its 6,800 lines matched the file line for line.
- The files of its results — about 0.8–1.0 GB a rule.

## 2. Claude Code sessions

How `sim` keeps the tokens that the move to `spec` v0.16.3 and the translation of `docs/` saved, by the pattern of measurements 1.9 spec: budgets of size in `docs-check` on every pull request, and a measurement of the sessions at each epic's archive.

### 2.1. Budgets

`docs-check` measures in bytes, in KB of 1,024 (measurements 1.9 spec). The limits are the block `repository` of `spec/tools/docs-check/budgets.json`, with what [docs-budgets.json](../docs-budgets.json) at the root names in its place. What is measured, the largest of each in `sim` on 2026-11-13, after the translation of `docs/`, KB:

| What | Limit | The largest |
|---|---|---|
| CLAUDE.md / RUN.md / DEBT.md | 12 / 96 / 128 | 7.4 / 38.9 / 4.1 |
| a rule of `.claude/rules/` | 6 | 1.9, `gate` |
| an entry of the decision log, with its amendments | 90, raised from 40 | 78.6, the entry of the two criteria; the next is the entry of the published `olympiad-math`, 13.1, and the median 4.6 |
| the own text under a heading of a document | 24 | 22.0, the gate in CI in `docs/architecture.md` |
| a row of a table of RUN.md / DEBT.md | 4 / 4 | 1.4, В13 / 1.0, Д3 |

- The one raised limit is the entry's: the entry of the two criteria held the bands and the paths of the criteria, and every new pin added an amendment with the table of its paths. In Russian it was 92.6 KB; the translation took it to 78.6. 90 is the largest plus about 15%, as measurements 1.9 spec sets limits. Raising it again is an edit of `docs-budgets.json`, and the third raise since the last analysis is a threshold of a new analysis of the sessions.
- Before the translation of the second half of the log the check, with `docs-budgets.json` at `{}`, found only the entry of the two criteria over its limit: the second half of the log, still Russian, was within 40 KB an entry.
- Since the judgement of `theta` the table of the paths is a section of this file, and the entry of the two criteria cites it; the test of the paths reads the table here. That entry then stopped growing by a table at every new pin.
- The section of `docs/architecture.md` on the gate in CI is 22.0 of 24 KB: the next growth of that section moves a part of it into a section of its own rather than raising the limit.

### 2.2. The measurement of sessions

`just -f spec/justfile tokens-report --since … --until … --row` prints one row of the table below. It runs at the archive of each epic, for the epic's window, and once a month from `spec` (measurements 1.9 spec).

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-sim` — before, the row of measurements 1.9 spec | 2026-10-13T00:00 … 2026-11-12T14:27 | 10 | 3,275 / 1,195 | 325.5 / 79.5 | 6 | $397.59 |
| `-workspaces-sim` — the move to `spec` v0.16.3 and the translation of `docs/` | 2026-11-12T20:41 … 2026-11-13T05:18 | 7 | 564 / 454 | 180.3 / 105.0 | 2 | $60.56 |
| `-workspaces-sim` — epic E5, and the move and the translation outside it | 2026-11-11T23:07 … 2026-11-14T02:34 | 16 | 1,557 / 634 | 179.4 / 97.4 | 3 | $142.28 |
| `-workspaces-sim` — epic E7, and the end of E5's last task with the laying out of E7 | 2026-11-14T02:34 … 2026-11-14T15:27 | 6 | 421 / 50 | 172.5 / 75.0 | 1 | $30.30 |
| `-workspaces-sim` — epic E9, and the end of E7's last task with the laying out of E9 | 2026-11-14T15:27 … 2026-11-15T00:31 | 7 | 685 / 227 | 182.2 / 81.8 | 1 | $57.76 |
| `-workspaces-sim` — epic E11, and the end of E9's last task | 2026-11-15T00:31 … 2026-11-16T01:37 | 5 | 255 / 87 | 153.3 / 70.3 | 0 | $18.71 |
| `-workspaces-sim` — epic E12: its first task up to its archive, and the end of E11's last task with the laying out of E12 | 2026-11-16T01:37 … 2026-11-16T02:11 | 2 | 75 / 0 | 112.0 / — | 0 | $3.69 |

- The context of a main call fell from 325.5 to 180.3 thousand, −45%; by session it was 178.0, 153.4, 128.8 and 109.5 in its four tasks in turn, while the documents a session reads became English, one task at a time. A subagent's rose from 79.5 to 105.0: here the subagents were translators and blind readers, each of which reads the glossary of spec's `docs/concepts.md` and the rules of translation. Three of the seven sessions, of one, one and five calls, ran no task.
- The windows differ in work: E2–E5 with code, runs and judgements against the move to a tag and the translation of `docs/`. The comparison by sessions after the task is spec's work.
- The window of E5, measured in its last task, opens after the commit of the base on `learner-model` v0.1.6 and before 4113692, which laid E5 out, and ends at 02:34, while that task still ran: its last calls are not in the row. It holds the window of the move and the translation above: the same report run again on that window gives 7 sessions and $60.32, so the tasks of the epic itself took 9 sessions and about $82. The context of a main call, 179.4 thousand, is the same as there, 180.3, against 325.5 before them.
- The window of E7, measured in its last task, opens where that of E5 ended, so it holds the end of E5's last task and the laying out of E7, which ran in that task's session; it ends at 15:27, while E7's last task still ran, and its last calls are not in it. The context of a main call, 172.5 thousand, is near that of E5, 179.4; the cost is about $5 a session, against about $9 in E5. E7 had less to read and write: one change of code, the set of confirmation by the judgement, one decisive run, that of the judgement of `engine-lesson-1` (measurements 1.3), and the base on the release of the chosen not done, `[-]`.
- The window of E9, measured in its last task, opens where that of E7 ended, so it holds the end of E7's last task, $5.77, and the laying out of E9, a session of its own, $6.28; it ends at 00:31, while E9's last task still ran, and its last calls are not in it. By session: the criterion of fairness $9.55, the populations and the measure $18.13, the constraint of fairness in the criteria and the tag $7.76, the run of the base $8.20, so the four tasks took $43.64. The populations, the measure and the evaluation in code with their tests took the most — 201 calls and 77 of subagents —, and the context of a main call, 182.2 thousand, is near that of E5 and E7.
- The window of E11, measured in its last task, opens where that of E9 ended, so it holds the end of E9's last task with the discussion of judgement rounds by plan, $7.29, and two short sessions of `spec`'s orchestrator, $0.45; it ends at 01:37, while E11's last task still ran, and its last calls are not in it. By session: the move to `spec` v0.18.3 $6.62, the tables of `run/` $4.34 up to then, its review's subagents among them. No compaction; the context of a main call, 153.3 thousand, is the lowest of the epics, since neither task ran code or long runs.

## 3. The paths of the targets

The values of the base and of the ceiling that the path of each measure is made of, read by `sim paths` on the decisive runs of the base and of the engine's ceiling and pinned in `internal/criterion/paths.json`. A new pin of what the paths are read on gives a new subsection here.

### 3.1. The base on `theta`

What: the paths of "the base on `theta`", the candidate the judgement `model-2` chose. On: `learner-model` at commit a063835 — `v0.1.7-0.20261113040121-a063835087c7` in `go.mod` —, `engine` v0.1.6 (57e0ab8), the packages of `modules` v0.2.5 and 2,000 working learners a cell, seed 20261107. How: `just decisive engine runs/base.jsonl` and `just decisive ceiling runs/ceiling.jsonl` by `go.mod`, then `sim paths -table`.
- The base, SHA-256 `3b4e4d28206227dcec01442bb55cd3fe07c563651b694874706045f3b5b299d3`, is to the byte past its header the decisive run of the candidate `theta`, `runs/model-2-theta.jsonl` (judgements/model-2/runs.txt): the header names the libraries by versions, and the candidate's by commits.
- The ceiling, SHA-256 `8b3675601c9d9cd8b31f836231bf5ee955e9e68f2aaa595502306e16ecbd92ec`: its lines moved against the ceiling on `learner-model` v0.1.6, since the engine's ceiling chooses with the model's readout. One twentieth of it, part 7 of 20, was run anew and matched it line for line, 6,500 lines.
- Of 85 paths, 34 moved against those on `learner-model` v0.1.6, all of the populations `F-`: `M3.settle` of `F-break` 161.08 → 148.40 on olympiad maths and 336.39 → 293.84 on school maths, `M5.caught` of `F-exp` and `F-growth5` of school maths 3.34 → 2.79 and 2.96 → 2.43, `M4.false.break` of `F-` of olympiad maths lower by 0.004–0.012; the rest by less than 0.1 of their value.

Values in the direction of the criterion, less is better:

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.7448 | 0.0192 | 0.7256 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.5216 | 0.0000 | 0.5216 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0289 | 0.0000 | 0.0289 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.6559 | 0.0000 | 4.6559 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8226 | 0.7471 | 0.0755 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2348 | 0.0000 | 0.2348 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0076 | 0.0000 | 0.0076 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.8154 | 0.0000 | 15.8154 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1852 | 0.0000 | 1.1852 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.3496 | 0.3496 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0436 | 0.0000 | 1.0436 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3755 | 0.3755 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.7402 | 0.0129 | 0.7273 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 1.0026 | 0.0000 | 1.0026 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8701 | 0.7003 | 0.1698 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4269 | 0.0000 | 2.4269 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.7332 | 0.0132 | 0.7200 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.4269 | 0.0000 | 0.4269 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.5182 | 0.0000 | 0.5182 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8118 | 0.6890 | 0.1228 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3629 | 0.0000 | 0.3629 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0917 | 0.0000 | 1.0917 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.7304 | 0.0107 | 0.7198 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.4237 | 0.0000 | 0.4237 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8148 | 0.6746 | 0.1401 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7551 | 0.0000 | 0.7551 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 57.5000 | 0.0000 | 57.5000 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 41.2255 | 0.0000 | 41.2255 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 47.2078 | 0.0000 | 47.2078 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.4096 | 0.0000 | 54.4096 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0554 | 0.0000 | 0.0554 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.7019 | 0.0000 | 4.7019 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0287 | 0.0000 | 0.0287 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 16.0443 | 0.0000 | 16.0443 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.5051 | 0.0000 | 1.5051 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1110 | 0.1110 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0473 | 0.0000 | 0.0473 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.4712 | 0.0000 | 0.4712 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1651 | 0.1599 | 0.0051 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | 0.0306 | 0.0000 | 0.0306 |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 3.7208 | 0.0000 | 3.7208 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.6935 | 0.0000 | 0.6935 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1547 | 0.1544 | 0.0004 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 4.0114 | 0.0000 | 4.0114 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0582 | 0.0000 | 0.0582 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.5369 | 0.0000 | 0.5369 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1567 | 0.1475 | 0.0092 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | 0.0304 | 0.0000 | 0.0304 |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 3.7170 | 0.0000 | 3.7170 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.7865 | 0.0000 | 0.7865 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1108 | 0.0906 | 0.0202 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 2.7890 | 0.0000 | 2.7890 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1591 | 0.1499 | 0.0092 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | 0.0263 | 0.0000 | 0.0263 |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 3.6518 | 0.0000 | 3.6518 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1199 | 0.1026 | 0.0173 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 1.6252 | 0.0000 | 1.6252 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0900 | 0.0000 | 0.0900 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.6345 | 0.0000 | 0.6345 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1928 | 0.1813 | 0.0115 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | 0.0300 | 0.0000 | 0.0300 |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 3.6416 | 0.0000 | 3.6416 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0003 | 0.0000 | 0.0003 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.9174 | 0.0000 | 0.9174 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1581 | 0.1565 | 0.0016 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 2.4267 | 0.0000 | 2.4267 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 148.3999 | 0.0000 | 148.3999 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 293.8350 | 0.0000 | 293.8350 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1735 | 0.1417 | 0.0317 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1972 | 0.1195 | 0.0777 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1550 | 0.1484 | 0.0066 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0036 | 0.0000 | 0.0036 |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.7762 | 0.0000 | 0.7762 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0795 | 0.0000 | 0.0795 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.1936 | 0.0000 | 0.1936 |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.5218 | 0.0000 | 0.5218 |

### 3.2. The release of `theta`

What: the paths of "the release of `theta`". On: `learner-model` v0.1.8 (ab688b2), the release of `theta` with the code of a063835 (sources 3.8), `engine` v0.1.6 (57e0ab8), the packages of `modules` v0.2.5 and 2,000 working learners a cell, seed 20261107. How: `just decisive engine runs/base-v018.jsonl` and `just decisive ceiling runs/ceiling-v018.jsonl` by `go.mod`, then `sim paths -table`.
- The base, SHA-256 `7d3cb52e72cf4c55d28dc11a16e8ed2330a9f06e78ccd9ab3ae01027fae7bf16`, and the ceiling, SHA-256 `00ad8bb7e82fb45542b7c2807f1541ce603663459f3629b24a65832f64043218`: past the header, which names the library's version, each is byte for byte its run of 3.1.
- So all 85 paths are those of 3.1 to the last digit.
- For the criterion of the judgement of consolidation in a lesson the base was run anew, by `go.mod` on the same libraries, because a line under the rule `engine` carries the sums of the note `E6`: SHA-256 `2e53e18c08df5016a70d2759fa97cbc8df6a4b72422e952eb6db5aef495f9aad`. With every key `"E6.…"` and its value removed, the file is byte for byte the base above, as the base of the bench's step for consolidation in a lesson is (`cmp`); the ceiling is the file above. So `paths.json` changes only the base's SHA-256, and the table is the same.
- For the fairness of the base by the grade the base was run anew, from a clean clone at `v0.1.19` on the same libraries, because the decisive run holds the three cells of fairness by the grade, `D-grade1`, `D-grade3` and `D-grade5`, 68 cells: SHA-256 `9e89cb7cfd0f1f0c00bb69c057a273aff0071e11b97450859d7b268e48c07763`. Past the header and without the 6,000 lines of the populations `D-`, the file is byte for byte the base run anew for `E6` above (`cmp`), and the header differs only by the three populations; the ceiling, which runs no `D-`, is the file above. `sim paths -table` on the two gives the table below to the last digit, so `paths.json` changes only the base's SHA-256 again.

Values in the direction of the criterion, less is better:

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.7448 | 0.0192 | 0.7256 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.5216 | 0.0000 | 0.5216 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0289 | 0.0000 | 0.0289 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.6559 | 0.0000 | 4.6559 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8226 | 0.7471 | 0.0755 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2348 | 0.0000 | 0.2348 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0076 | 0.0000 | 0.0076 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.8154 | 0.0000 | 15.8154 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1852 | 0.0000 | 1.1852 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.3496 | 0.3496 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0436 | 0.0000 | 1.0436 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3755 | 0.3755 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.7402 | 0.0129 | 0.7273 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 1.0026 | 0.0000 | 1.0026 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8701 | 0.7003 | 0.1698 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4269 | 0.0000 | 2.4269 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.7332 | 0.0132 | 0.7200 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.4269 | 0.0000 | 0.4269 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.5182 | 0.0000 | 0.5182 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8118 | 0.6890 | 0.1228 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3629 | 0.0000 | 0.3629 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0917 | 0.0000 | 1.0917 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.7304 | 0.0107 | 0.7198 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.4237 | 0.0000 | 0.4237 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8148 | 0.6746 | 0.1401 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7551 | 0.0000 | 0.7551 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 57.5000 | 0.0000 | 57.5000 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 41.2255 | 0.0000 | 41.2255 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 47.2078 | 0.0000 | 47.2078 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.4096 | 0.0000 | 54.4096 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0554 | 0.0000 | 0.0554 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.7019 | 0.0000 | 4.7019 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0287 | 0.0000 | 0.0287 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 16.0443 | 0.0000 | 16.0443 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.5051 | 0.0000 | 1.5051 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1110 | 0.1110 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0473 | 0.0000 | 0.0473 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.4712 | 0.0000 | 0.4712 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1651 | 0.1599 | 0.0051 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | 0.0306 | 0.0000 | 0.0306 |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 3.7208 | 0.0000 | 3.7208 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.6935 | 0.0000 | 0.6935 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1547 | 0.1544 | 0.0004 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 4.0114 | 0.0000 | 4.0114 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0582 | 0.0000 | 0.0582 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.5369 | 0.0000 | 0.5369 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1567 | 0.1475 | 0.0092 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | 0.0304 | 0.0000 | 0.0304 |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 3.7170 | 0.0000 | 3.7170 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.7865 | 0.0000 | 0.7865 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1108 | 0.0906 | 0.0202 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 2.7890 | 0.0000 | 2.7890 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1591 | 0.1499 | 0.0092 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | 0.0263 | 0.0000 | 0.0263 |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 3.6518 | 0.0000 | 3.6518 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1199 | 0.1026 | 0.0173 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 1.6252 | 0.0000 | 1.6252 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0900 | 0.0000 | 0.0900 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.6345 | 0.0000 | 0.6345 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1928 | 0.1813 | 0.0115 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | 0.0300 | 0.0000 | 0.0300 |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 3.6416 | 0.0000 | 3.6416 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0003 | 0.0000 | 0.0003 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.9174 | 0.0000 | 0.9174 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1581 | 0.1565 | 0.0016 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 2.4267 | 0.0000 | 2.4267 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 148.3999 | 0.0000 | 148.3999 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 293.8350 | 0.0000 | 293.8350 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1735 | 0.1417 | 0.0317 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1972 | 0.1195 | 0.0777 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1550 | 0.1484 | 0.0066 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0036 | 0.0000 | 0.0036 |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.7762 | 0.0000 | 0.7762 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0795 | 0.0000 | 0.0795 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.1936 | 0.0000 | 0.1936 |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.5218 | 0.0000 | 0.5218 |
