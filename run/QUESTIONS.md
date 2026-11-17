# QUESTIONS.md — questions to the author of `sim`

Table 3 of [RUN.md](../RUN.md): the author's open questions from tasks' «Открыто» (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| В25 ' run/QUESTIONS.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row comes from a task's «Открыто» (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main` (R81 spec; RULES.md, Across repositories). It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В33.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В29 | Epic E10 with its judgement deferred until there are 100 children (R69): T72 and T73 are `[?]`, and RUN.md, «Epics», archives an epic only once table 1 is all checked, so E10 stays current and no next epic is laid out while the judgement waits — for months, by the author's threshold. Archive it, or keep it? | T77 | (a) once T76 is done, the task after it archives E10 with T72 and T73 `[?]` in the archive, by the author's word as an exception to «Epics»; the judgement, when a candidate comes, is laid out as an epic anew by R59, with R65 and `held-out-model-fair-1` as they are. (b) E10 stays current until the judgement returns; every new task stands outside it in its table 1 (R61) and no epic of table 2 is taken | (a): R59 lays out a judgement only when a library brings a candidate, and none is coming; (b) holds the plan on a judgement with no date |
| В30 | The epic's rules «одно закрепление» and «биты 68 ячеек не двигаются до выбора суда» were lifted by the executor while the judgement is deferred (R69): the judgement, when it returns, pins the tag — T71 or later — and the base that an amendment to R65 names before its first candidate. Confirm, or keep the pin? | T77 | (a) confirm R69: releases of `spec`, `engine` and `modules` and the rules in place of calibration move the base through the gate meanwhile, and T76 takes the tag of `spec` even if the bits move. (b) keep the pin: `go.mod`, `gate.json`, the packages and the paths stay until the judgement returns; T76 only with the same bits, otherwise after T73 | (a): the pin held the numbers of a running judgement; with none it would hold the gate against every release for an unknown time |
