# DEBT.md — technical debt of `sim`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д3 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row comes from a task's «Открыто» (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed. What closed it is named by the «Сделано» of the task that closed it.
- «Д3» of the root `DEBT.md` before T67 is the same row. The records of the closed Д1, Д2 and Д4–Д9 that `DEBT.md` held are in [docs/plan/E11.md](../docs/plan/E11.md), «Из RUN.md и DEBT.md, при переходе на `run/`».

## Table 4. Technical debt

Numbers are not reused: next — Д14.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д3 | No check of the simulator's numbers against live learners: the populations are its assumptions, their speeds, laws of forgetting and shares not measured (R03). Standalone checked two tables of its service's report against live children, and its simulator computed the same numbers for its own (sources 1.1) | T01 | The simulator may choose a rule good on its assumptions and worse on live learners, unseen | Stage 2, when `api` holds live learners' data (architecture 1 spec; repositories 2.1 spec): read the path against privacy spec first, propose it as a row of [QUESTIONS.md](QUESTIONS.md) «В реестр spec» |
| Д12 | The base of the gate stands on a commit of `learner-model`, df338e5 of its pull request #84, not on a release; the paths and the base's decisive run are still those of T64, on v0.1.8 (R74) | T79 | `sim judge` refuses a base on v0.1.8 against the gate's pin, `sim evaluate` the run of `fairness-1`, and a judgement has no paths of the rules; none runs while the judgement is deferred (R69) | Once `learner-model` releases the rules after merging #84: a task moves `gate.json` and `go.mod` to the release, `release-verify` first, `just gate exact` the same numbers, as T48 (R46); the paths with the judgement's return (R65) |
| Д13 | The adult's correction of `C-grade1` is one fact of `LevelCorrection` to place 2 at the 10th answer, while `SetModuleLevel` of `spec` v0.18.14 takes one band at a time and appends nothing for the band the level stands in; the end state is the same but for the learners whose level already stands in place 2, where `sim`'s fact ends the trial series and the guard's count and the contract's would not. Their share is not measured: the trace holds θ, not the band of the level (R75) | T81 | `rules-1` and `rules-2` overstate or understate the correction for that share, unknown, in the numbers of row В33 for dpia 10.2 spec | Before the correction is measured again, or when `spec` asks for the numbers of В33 anew: a task counts the learners of `C-grade1` with the level in place 2 at the 10th answer, by the band `learner-model` gives, and has `run.correct` skip the fact for them, the criterion first (R73) |
