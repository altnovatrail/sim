# QUESTIONS.md — questions to the author of `sim`

Table 3 of [RUN.md](../RUN.md): the author's open questions from tasks' «Открыто» (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| В25 ' run/QUESTIONS.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row comes from a task's «Открыто» (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main` (R81 spec; RULES.md, Across repositories). It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В26.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В25 | When does Д10 of [DEBT.md](DEBT.md) run — the move of CI to the latest `.github` release with the Dependabot group? E10 waits for `learner-model` with no date, and between epics no task runs (R50) | T67 | (a) an epic of its own now, one task, as R60 laid out E11 for the move to a tag; (b) the first task of E10, whenever it is laid out; (c) wait for the next CI task | (a): Д10's «when» is before the next CI task, and while it waits every job builds its image from the network |
