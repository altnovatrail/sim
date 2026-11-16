# DEBT.md — technical debt of `sim`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д3 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row comes from a task's «Открыто» (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed.

## Table 4. Technical debt

Numbers are not reused: next — Д11.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д3 | No check of the simulator's numbers against live learners: the populations are its assumptions, their speeds, laws of forgetting and shares not measured (R03). Standalone checked two tables of its service's report against live children, and its simulator computed the same numbers for its own (sources 1.1) | T01 | The simulator may choose a rule good on its assumptions and worse on live learners, unseen | Stage 2, when `api` holds live learners' data (architecture 1 spec; repositories 2.1 spec): read the path against privacy spec first, propose it as a row of [QUESTIONS.md](QUESTIONS.md) «В реестр spec» |
| Д10 | CI calls `.github`'s shared checks at `v0.1.21`, without the cached `toolchain` image — in `.github` since `v0.1.22`, latest `v0.1.24` (Д146 spec); `.github/dependabot.yml` has no group `shared-workflows`, `patterns: ["altnovatrail/.github/*"]` (Д88 spec). Neither waits for a tag of `spec` (sources 7.9; R60) | T66 | Each job builds the image from the network: an outage of nodejs.org, GitHub or the Ubuntu snapshot fails the checks; Dependabot opens a pull request per `.github` reference | A task of its own — the latest `.github` release and the Dependabot group —, before the next CI task |
