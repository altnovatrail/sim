# RUN.md — work plan of sim

What sim does, does not do and has still to decide — [docs/architecture.md](docs/architecture.md); why — [docs/decisions.md](docs/decisions.md). The platform is in `spec/`, a submodule at tag v0.19.11: the component's card — architecture 3.9 spec. Small tasks, one at a time, gathered in epics; only the current epic is in detail (R19 spec). Other tables, each a file opening with its rules (R153 spec): 2, epics for later — [run/EPICS.md](run/EPICS.md); 3, questions to the author — [run/QUESTIONS.md](run/QUESTIONS.md); 4, debt — [run/DEBT.md](run/DEBT.md). No `run/TASKS.md`: a task outside an epic stands in table 1 of the current epic (R47, R61).

Language (R154 spec): rules and tables in English, terse, for the session; the epic's description, table 1 and the bodies of tasks in Russian — the author writes and reviews them. Commit messages, issues and a pull request's title and description in English, a person's too; a session's terse, for the session, never a copy of the report (the 2026-11-17 amendments to R14 spec; R75, R77).

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
9. When done: table 1 `[ ]` → `[x]`; append to the body «Сделано (дата)» — one line, what was done and its pull request, which says how it was checked; where the author commits after review, the report says it. What stays open is rows, not text: a question for the author, of [run/QUESTIONS.md](run/QUESTIONS.md); what is put off, of [run/DEBT.md](run/DEBT.md) (R160 spec).
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

**Task body** (Russian): «Цель», «Сделать», «Как проверить», «Разделы»; as needed «Не делать», «Решить при планировании», «Мой шаг», «Стоимость», «Риски»; from planning «Почему», «Решено автором (дата)», «Выяснилось при планировании»; when done «Сделано (дата)» (rule 9).

## Epics

- **Table 1** is the current epic: its tasks in detail, run one at a time.
- **Next epic.** Once table 1 is all checked, the author picks the next epic from [run/EPICS.md](run/EPICS.md) and splits it in plan mode into tasks of one review each, into table 1 and this file. A judgement epic is laid out only when a library brings a candidate with a reason of its own; no rounds by plan (R59).
- **A done epic** goes to `docs/plan/<epic>.md` (R59 canvas-api): bodies and gist sentences are there already (rule 9); move its description and table of tasks verbatim; questions and debt stay in `run/`; until the next epic, table 1 reads «Current epic — none» and names the last done epic. Measure the epic's sessions — `just -f spec/justfile tokens-report --since <first day of the epic> --row` — as a row of measurements 2.2 (R145 spec). No task «итог эпика»: the epic's last task archives it, measures it and proposes the next (R59).
- Never read an archive whole: find a task by heading — `grep -n '^#### T03\.' RUN.md docs/plan/*.md`; from another repository — «T03 sim». Task numbers are global, never reused, and taken free at the moment of writing.

## Table 1. Current epic — E10. Суд круга справедливости модели

Основание — круг строки 4.17 COORDINATION.md spec, который `learner-model` взял своим эпиком E8 (R64; sources 3.10). Последний сделанный эпик — E12, [docs/plan/E12.md](docs/plan/E12.md).

**Результат** (строка 4.18 COORDINATION.md spec):
- критерий суда круга справедливости записан до первого кандидата, как у `engine-lesson-1` (R49): своя папка и свой отложенный набор (R52), критерий модели как есть, цель — `M11.same` у `D-grade1` и `D-grade5` (R65, R66). `sim judge` судит по нему; тег суда — T71 или позже — и его базу называет поправка к R65 до первого кандидата (R69);
- суд кандидатов `learner-model` по коммиту, `CANDIDATE_LEARNER_MODEL` (R34), прошёл: выбранный подтверждён на своём отложенном наборе, его полосы записаны, и база стоит на его релизе, — или выход записан, и база остаётся;
- итог — строкой «в реестр spec» для dpia 10.2 spec; эпик — в архиве, следующий предложен (R59).

**Суд отложен** (R69): долгих прогонов подбора констант на `sim` нет, пока детей меньше 100, а разведка T61 `learner-model` кандидатов не дала. T72 и T73 — `[?]`; критерий и набор суда ждут его возвращения. Эпик остаётся текущим, пока таблица 1 не вся отмечена (правила RUN.md, «Epics»; строка В29 [run/QUESTIONS.md](run/QUESTIONS.md)).

| Статус | ID | Задача | Кто | План | Зависит от |
|---|---|---|---|---|---|
| [x] | T69 | Развернуть E10 «Суд круга справедливости модели» | Claude | — | — |
| [x] | T70 | Критерий суда круга справедливости — до первого кандидата | Claude, проверяю я | да | — |
| [x] | T71 | Цель `M11.same` и набор суда в коде, тег | Claude | да | T70 |
| [?] | T72 | Суд `model-fair-1` | я + Claude | да | T71; кандидаты — черновые pull request T62 `learner-model` (строка 4.17 COORDINATION.md spec), не их релиз |
| [?] | T73 | База — на релиз выбранного, архив эпика | я + Claude | да | T72; при выборе — релиз выбранного, T63 `learner-model` |
| [x] | T74 | Скорость CI: строка 7.19 | Claude | — | — |
| [x] | T75 | Переход на релиз `.github` строки 7.13 | Claude | да | T74; релиз `.github` строки 7.13 COORDINATION.md spec |
| [x] | T76 | Переход на тег `spec` строки 7.14: стадия `toolchain-core` | я + Claude | да | T75; тег `spec` строки 7.14 COORDINATION.md spec; релиз `.github` с T66 `.github` |
| [x] | T77 | Суд круга справедливости — до 100 детей | Claude | — | — |
| [x] | T78 | Переход на тег `spec` v0.18.6 | Claude | — | — |
| [x] | T79 | Короткий прогон трёх правил R157 spec | Claude | да | T78; pull request #84 `learner-model` (строка 4.20 COORDINATION.md spec) — его коммит, не релиз |
| [x] | T80 | README без «Status», «License», фразы о языках и версии spec | Claude | — | — |
| [x] | T81 | Переход на тег `spec` v0.18.14 | Claude | — | — |
| [x] | T82 | База гейта — релиз `learner-model` v0.1.13 (Д12) | Claude | — | релиз `learner-model` v0.1.13 с #84 |
| [x] | T83 | Переход на тег `spec` v0.18.27 и релиз `.github` v0.1.33 | Claude | — | — |
| [x] | T84 | Тексты без истории (R160 spec) и тег `spec` v0.18.34 | Claude | — | — |
| [x] | T85 | Тексты без истории: `docs/architecture.md`, разделы 1–3.4 | Claude | — | T84 |
| [x] | T86 | Тексты без истории: `docs/architecture.md`, разделы 3.5–3.9 | Claude | — | T84 |
| [x] | T87 | Тексты без истории: `docs/README.md`, README, CLAUDE.md, `.claude/rules/`, план | Claude | — | T84 |
| [x] | T88 | Комментарии кода по R160 spec: корень, `internal/bench`, `criterion`, `gate` | Claude | — | T84 |
| [x] | T89 | Комментарии кода по R160 spec: `internal/learner`, `measure`, `run`, `invariant` | Claude | — | T84 |
| [x] | T90 | Комментарии кода по R160 spec: остальные пакеты, justfile, `.github/`, `.devcontainer/` | Claude | — | T84 |
| [x] | T91 | Переход на тег `spec` v0.19.4 | Claude | — | — |
| [x] | T92 | Переход на тег `spec` v0.19.11 | Claude | — | — |

**Почему такой порядок.**
- Суд — по коммитам черновых pull request кандидатов, не по их релизу: `sim` гонит кандидата из его коммита (R34), а релиз выбранного ждёт полос `sim` (R103 spec).
- База на релиз — после суда; при выходе T73 только архивирует эпик. Отдельной задачи итога нет: последняя архивирует (R59).
- T74–T92 — вне эпика, в таблице 1, а не в `run/TASKS.md` (R61): каждая — по строке COORDINATION.md spec, слову автора или эпику `spec` в прогоне его оркестратора (R156 spec), её причина — в её теле; тела сделанных — в [docs/plan/E10.md](docs/plan/E10.md).
- T87–T90 — строка 7.31 COORDINATION.md spec, по одному ревью (правило 11): друг от друга не зависят и не двигают ни бит, ни закрепления.

## Epic E10. Суд круга справедливости модели

**Почему эпик такой.** Оценка `fairness-1` нашла: база — `learner-model` v0.1.8 и `engine` v0.1.6 — не держит критерий справедливости по классу: у `D-grade1` и `D-grade5` `M11.same` 0,812 [0,795; 0,829] и 0,5525 [0,53; 0,5735] при границе 0,90 (R58). `fair` держит кандидата «не хуже базы», но не награждает за лучшее, так что круг, который подтягивает базу к границе, судится со своей целью, записанной до первого кандидата (R58, R59). `learner-model` взял круг эпиком E8: разбор T59, «Gate» на теге `sim` с критерием E10 — T60, разведка T61, кандидаты — черновые pull request T62, суд и выпуск — T63, итог — T64 (sources 3.10).

**Общее для задач:**
- общее для задач E5, E7 и E9 действует и здесь ([docs/plan/E5.md](docs/plan/E5.md), «Эпик E5»; [docs/plan/E7.md](docs/plan/E7.md), «Эпик E7»; [docs/plan/E9.md](docs/plan/E9.md), «Эпик E9»): тесты и `codecov/patch`, `same-bits`, файл релиза — только после `release-verify`, язык; файлы решающих прогонов — в `runs/`; время — с load average хоста, а замер, по которому что-то выбирается, — на свободной машине;
- **критерий — до первого кандидата:** пока T70 не смержена, ни один кандидат круга не гоняется в `sim`, даже грубым взглядом; числа критерия после первого прогона кандидата не меняются (R103 spec; R163 standalone);
- **биты 68 ячеек не двигаются до выбора суда:** T71 меняет чтение результатов, а не прогон, так что файл базы — `runs/base-t64.jsonl`, чей SHA-256 называет `paths.json` (R58), — и 69 чисел гейта прежние. Сдвигает их только выбранный и подтверждённый кандидат — в T72 и T73;
- **одно закрепление:** от тега T71 до выбора суда `go.mod`, `gate.json`, пакеты и пути не меняются; релиз `spec`, `engine` или `modules`, который выйдет за это время, ждёт конца суда и берётся своей задачей. Исключение — тег `spec` с теми же битами (R45);
- **пока суд отложен** (R69): правила «биты 68 ячеек не двигаются до выбора суда» и «одно закрепление» не действуют, а «критерий — до первого кандидата» и «отложенный набор суда» — действуют — суда нет, и держать им нечего; вернувшийся суд называет свой тег и базу поправкой к R65 до первого кандидата (R69);
- **отложенный набор суда** не гоняет никто до подтверждения T72: ни разведка `learner-model`, ни грубый взгляд (R49, R52; раздел 3.8 [docs/architecture.md](docs/architecture.md));
- **в прогоне оркестратора `spec`** каждая задача — своя ветка от `main` и свой pull request, если промпт не называет общую — ветку эпика `spec` или пакета задач (R159 spec); мержит оркестратор на зелёных проверках (R156 spec; правило 7); вопросы плана — строки [run/QUESTIONS.md](run/QUESTIONS.md);
- **не входит:** кандидат `engine` — он идёт своей причиной (R59); числа забывания `learner-model` (sources 3.10).

## Задачи эпика

#### T72. Суд `model-fair-1`

- **Цель:** кандидаты круга справедливости `learner-model` осуждены по критерию T70, R65, на теге T71 или позже, который называет поправка к R65 (R69). Выбранный подтверждён на нетронутом отложенном наборе суда, и его полосы записаны, — или выход записан, и база остаётся ([docs/architecture.md](docs/architecture.md), разделы 3.6 и 3.8). Итог — строкой «в реестр spec» для dpia 10.2 spec.
- **Ждёт:** T71; черновых pull request кандидатов T62 `learner-model` (строка 4.17 COORDINATION.md spec) — их коммитов, не релиза. Если разведка T61 `learner-model` кончится выходом без кандидатов, исполнитель предлагает задаче `[-]` строкой [run/QUESTIONS.md](run/QUESTIONS.md), итог несёт T64 `learner-model` (sources 3.10), а T73 только архивирует эпик.
- **Сделать:** как T58 в [docs/plan/E7.md](docs/plan/E7.md) и T47 в [docs/plan/E5.md](docs/plan/E5.md):
  - тег — T71 или позже — и база — поправкой к R65 до первого кандидата (R69); кандидаты — коммиты черновых pull request `learner-model`, имена, цели и число констант — из их описаний (R25);
  - решающий прогон каждого — `CANDIDATE_LEARNER_MODEL=<sha> just decisive engine runs/<файл>` (R34) — на свободной машине, с load average хоста; база — `runs/base-t64.jsonl`, который называет `paths.json`;
  - `sim judge` по критерию T70, R65; файл суда, отчёт и SHA-256 прогонов — в `judgements/model-fair-1/` (R35);
  - если выбран — подтверждение, один раз, на отложенном наборе суда; если подтвердился — `gate.json` называет коммит кандидата, `go.mod` — его псевдоверсию; пути заново на нём — поправкой к R04 и подразделом раздела 3 [docs/measurements.md](docs/measurements.md); полосы — `just gate record` с причиной — отчётом суда; мерж выпускает тег, после которого `learner-model` выпускает выбранного (T63 `learner-model`; R103 spec);
  - если нет — выход: база остаётся, отчёт записан; набор больше не отложен, если подтверждение шло (R52);
  - итог — решением журнала, как R46 и R53; строка «в реестр spec» — исход для строки 4.18 и для item 10 dpia 10.2 spec: `M11` выбранного у `D-grade1` и `D-grade5` с интервалами подтверждения и держит ли он границу R55, или выход с числами суда.
- **Решить при планировании:** что делать, если решающий прогон на свободной машине дольше 45 минут (раздел 3 [docs/architecture.md](docs/architecture.md), «What stays open»); повторять ли потолок на коммите выбранного; кандидату с иным стартом по классу, чем полоса класса (R51 spec), — ждать ли ответа `spec` на его строку до суда или только до выпуска.
- **Риски:** σ₀ 3,0 в разведке T30 `learner-model` срывала 28 ограничений, хотя `off-band` была лучше (sources 3.10): кандидат, что поднимает `M11.same`, может сорвать «не хуже» на других ячейках. Вероятен выход; решает суд.
- **Мой шаг:** решения; чтение отчёта суда; свободная машина на время решающих прогонов.
- **Отложено (2026-11-16):** `[?]` — решение автора 2026-11-16 через координатора `spec` (T171 spec), долгие прогоны — после 100 детей (R69): разведка T61 `learner-model` кандидатов не дала — порог 0,90 подбором констант не берётся, его T62 и T63 — `[?]` (R42 `learner-model`; sources 3.11), а долгих прогонов подбора констант на `sim` нет, пока детей меньше 100. Не `[-]` строкой таблицы 3, как предвидел «Ждёт»: суд вернётся. Критерий R65, цель R66, имя папки `judgements/model-fair-1/` (её ещё нет) и нетронутый набор `held-out-model-fair-1` — как есть; тег и базу суда называет поправка к R65 до первого кандидата. Итог для dpia 10.2 spec несёт T64 `learner-model`. Вернуться, когда детей будет 100, или когда решение `spec` о правилах вместо калибровки вернёт кандидата в T62 `learner-model` — по R59.
- **Как проверить:** часть прогона кандидата воспроизводится до байта на том же коммите (R28); `sim judge` не отказывает (R40); `just gate exact`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** строка 4.18 COORDINATION.md spec; R28, R34, R35, R40, R46, R52, R53, R55, R58; R65; [docs/architecture.md](docs/architecture.md), разделы 3.6–3.9; [docs/README.md](docs/README.md), «Judging a candidate»; dpia 10.2 spec; sources 3.10.

#### T73. База — на релиз выбранного, архив эпика

- **Цель:** при выборе — `gate.json` и `go.mod` называют релиз выбранного, пути прочитаны на нём, и база оценена заново, `evaluations/fairness-2/` (R55); в обоих исходах — эпик в архиве и следующий предложен (R59).
- **Ждёт:** T72; при выборе — релиза выбранного, T63 `learner-model`.
- **Решить при планировании:** при выборе — что делать, если решающий прогон оценки на свободной машине дольше 45 минут, как в T72.
- **Мой шаг:** при выборе — свободная машина на время решающего прогона оценки.
- **Сделать:**
  - при выборе, как T48 в [docs/plan/E5.md](docs/plan/E5.md): `release-verify` принимает релиз (R77 spec); `gate.json` и `go.mod` — на нём; `just gate exact` — числа те же, что на коммите кандидата; пути заново (R40); оценка базы — решающий прогон, `sim evaluate`, `evaluations/fairness-2/` (R55), и её итог — строкой «в реестр spec» для item 10 dpia 10.2 spec; мерж выпускает тег (R29);
  - в обоих исходах: раздел 3 [docs/architecture.md](docs/architecture.md) сверить с построенным; что меняет карточку architecture 3.9 spec — строкой «в реестр spec»; таблицы 2–4; эпик — в архив `docs/plan/E10.md`, замер токенов его сессий — строкой measurements 2.2; следующий эпик предложить (R59).
- **Как проверить:** при выборе — `just gate exact`, `just ci-test`, `just ci-lint`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Отложено (2026-11-16):** `[?]` — вслед за T72, по тому же решению автора (R69): суда нет, и базе не на что переходить; эпик не уходит в архив, пока таблица 1 не вся отмечена (правила RUN.md, «Epics»). Вернуться вместе с T72.
- **Разделы:** R04, R26, R29, R40, R55, R58, R59; [docs/architecture.md](docs/architecture.md), разделы 3, 3.6 и 3.9; architecture 3.9 spec; dpia 10.2 spec.
