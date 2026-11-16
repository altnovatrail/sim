# EPICS.md — epics for later of `sim`

Table 2 of [RUN.md](../RUN.md) (R153 spec, R154 spec).

- A row has a result and a basis, no tasks yet. Row order is the order epics are taken, until the author decides otherwise.
- Read by the choice and planning of the next epic, not by each task. Once an epic is split into tasks, they go to table 1 and the body of RUN.md, and its row leaves (RUN.md, Epics).
- A judgement is laid out only when a library brings a candidate with a reason of its own — a bug, a product need, a row of `spec`; no rounds by plan (R59).

## Table 2. Epics for later

| ID | Epic | Result | Basis |
|---|---|---|---|
| E10 | «Суд круга справедливости модели» — the judgement of `learner-model`'s round of fairness | A judgement of `learner-model`'s candidates aimed at fairness by the grade — `M11.same` of `D-grade1` and `D-grade5` —, modelled on `engine-lesson-1` (R49): its own folder and held-out set, the model's criterion as is, written down before the first candidate. The chosen one is released and the base stands on it — or the exit is recorded | R59 — option (a) of row В22; R58 — the base does not hold it, and `fair` does not reward better; R55. Waits for `learner-model` to take the round: rows 4.17 and 4.18 of COORDINATION.md spec. Tag — `v0.1.20` or later (amendment to R57) |
