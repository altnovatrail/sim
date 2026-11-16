# RUN.md — work plan of sim

What sim does, does not do and has still to decide — [docs/architecture.md](docs/architecture.md); why — [docs/decisions.md](docs/decisions.md). The platform is in `spec/`, a submodule at tag v0.18.3: the component's card — architecture 3.9 spec. Small tasks, one at a time, gathered in epics; only the current epic is in detail (R19 spec). Other tables, each a file opening with its rules (R153 spec): 2, epics for later — [run/EPICS.md](run/EPICS.md); 3, questions to the author — [run/QUESTIONS.md](run/QUESTIONS.md); 4, debt — [run/DEBT.md](run/DEBT.md). No `run/TASKS.md`: a task outside an epic stands in table 1 of the current epic (R47, R61).

Language (R154 spec): rules and tables in English, terse, for the session; the epic's description, table 1 and the bodies of tasks in Russian — the author writes and reviews them.

## Rules

A task starts with «Выполни задачу T01 из RUN.md» in a new session in the devcontainer. One task — one session (RULES.md, How we work).

1. Do only the named task; touch other tasks' files only where it requires.
2. Read first: [CLAUDE.md](CLAUDE.md) with RULES.md of `spec/`, the task, the sections it cites, [docs/decisions.md](docs/decisions.md).
   - Conflict with the documents or something missing — stop and ask; never guess.
   - A gap in the component itself — a subsection of [docs/architecture.md](docs/architecture.md), section 2. A gap across repositories — a row of [run/QUESTIONS.md](run/QUESTIONS.md) whose question opens with the mark «В реестр spec» (R81 spec; RULES.md, Across repositories).
   - From [run/DEBT.md](run/DEBT.md) and [run/QUESTIONS.md](run/QUESTIONS.md) take only the rows the task and its sections name and those `grep -n` finds by the paths and parts it changes (R143 spec).
   - The checks of a part the task touches — its rule in `.claude/rules/` (CLAUDE.md, Environment).
3. Devcontainer only, with `spec/` checked out (RULES.md, Hard rules; CLAUDE.md).
4. Outside this repository — editing `spec` or other repositories, their settings, anything that costs money — only after a warning (what, where, cost) and the author's confirmation. `spec/` is read only here; a task of its own moves it to a new tag. The exception of R146 spec, the task's own checks in a browser or with a real AI call, does not arise: `sim` has none (R51).
5. External facts — from the primary source, as a record of `docs/sources.md` (RULES.md, How we work).
6. End with the checks of what the task changed; documents — `just -f spec/justfile docs-check-repo` (R32 spec; what it checks — CLAUDE.md, Environment). Then `/code-review high`, one round: fix the findings, do not rerun. A small change may skip the review — the executor decides.
7. Do not commit: the author commits after review. Otherwise only on the author's word «доведи до мержа», and in a task of spec's orchestrator, «Сделай эпик E…»: its executor works by the contract `spec/tools/orchestrator/executor.md` — commits, pushes and opens the pull request, never merges (RULES.md, How we work; R156 spec; R60).
8. End with a short report in Russian: what was done, files, how to check, what the review found or why there was none, what is open and what goes «в реестр spec».
9. When done: table 1 `[ ]` → `[x]`; append to the body «Сделано (дата)» — what was done, how checked — and «Открыто». A question for the author from «Открыто» is also a row of [run/QUESTIONS.md](run/QUESTIONS.md); what is put off, a row of [run/DEBT.md](run/DEBT.md).
   - Then move the body verbatim, relative links fixed, to the epic archive `docs/plan/<epic>.md`, with a gist sentence — an unheaded list item at the start of the receiving section: «T01 — сделана <дата>: <итог одной фразой>». Neither stays here (R153 spec).
   - A task outside an epic stands in table 1 of the current epic and lives as its task. Between epics there is none: it waits for the next epic and becomes its first task (R47, R50, R61).
10. Language — RULES.md, Hard rules; above, «Language».
11. A task larger than one review (≈ 30–60 minutes of reading): propose a split before starting.
12. Before a new task, close the findings on the previous one (R30 spec); an allowed departure from rule 1.
    - SonarCloud sees merged code: [New Code](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim); CodeQL and Scorecard — `gh api 'repos/altnovatrail/sim/code-scanning/alerts?state=open'`.
    - Fix in code first; "false positive" is the last resort, with the reason in words — shared reasons in R33 spec, own ones in the log, like R51 canvas-api. Never skip silently.
    - Report these fixes in a separate paragraph: the author commits them apart from the task.

**Statuses:** `[ ]` not started; `[x]` done; `[-]` not done, reason in the task; `[?]` deferred, reason and when to return in the task. A task that fails the author's review returns to `[ ]`.

**Who:** «Claude» — the executor does it, the author reviews; «Claude, проверяю я» — texts the author reads line by line; «я + Claude» — the author does some steps: decisions, services, GitHub settings, confirmations by rule 4.

**Plan:** «да» — plan mode first: decisions missing from the documents, or larger than one review; «—» — run at once.

**Task body** (Russian): «Цель», «Сделать», «Как проверить», «Разделы»; as needed «Не делать», «Решить при планировании», «Мой шаг», «Стоимость», «Риски»; from planning «Почему», «Решено автором (дата)», «Выяснилось при планировании»; when done «Сделано (дата)», «Открыто».

## Epics

- **Table 1** is the current epic: its tasks in detail, run one at a time.
- **Next epic.** Once table 1 is all checked, the author picks the next epic from [run/EPICS.md](run/EPICS.md) and splits it in plan mode into tasks of one review each, into table 1 and this file. A judgement epic is laid out only when a library brings a candidate with a reason of its own; no rounds by plan (R59).
- **A done epic** goes to `docs/plan/<epic>.md` (R59 canvas-api): bodies and gist sentences are there already (rule 9); move its description and table of tasks verbatim; questions and debt stay in `run/`; until the next epic, table 1 reads «Current epic — none» and names the last done epic. Measure the epic's sessions — `just -f spec/justfile tokens-report --since <first day of the epic> --row` — as a row of measurements 2.2 (R145 spec). No task «итог эпика»: the epic's last task archives it, measures it and proposes the next (R59).
- Never read an archive whole: find a task by heading — `grep -n '^#### T03\.' RUN.md docs/plan/*.md`; from another repository — «T03 sim». Task numbers are global, never reused, and taken free at the moment of writing.

## Table 1. Current epic — none

Last done epic: E11, [docs/plan/E11.md](docs/plan/E11.md). Next by [run/EPICS.md](run/EPICS.md): E10, once `learner-model` takes its round (row 4.18 of COORDINATION.md spec). When Д10 of [run/DEBT.md](run/DEBT.md), the move of CI to the latest `.github`, runs — row В25 of [run/QUESTIONS.md](run/QUESTIONS.md) (R61).
