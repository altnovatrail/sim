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

## Table 1. Current epic — E10. Суд круга справедливости модели

Развёрнут 2026-11-16 задачей T69 (R64), когда `learner-model` взял круг строки 4.17 COORDINATION.md spec — эпик E8, T59–T64, R40 learner-model, `main` e3aaee6 (sources 3.10). Последний сделанный эпик — E12, [docs/plan/E12.md](docs/plan/E12.md).

**Результат** (строка 4.18 COORDINATION.md spec):
- критерий суда круга справедливости записан до первого кандидата, как у `engine-lesson-1` (R49): своя папка и свой отложенный набор (R52), критерий модели как есть, цель — `M11.same` у `D-grade1` и `D-grade5` (В22, вариант (а); R59). `sim judge` судит по нему, и тег T71 с ним — закрепление суда: позже `v0.1.20` и `v0.1.22`, с `fair` и файлом базы 68 ячеек (поправка 2026-11-14 к R57; R64);
- суд кандидатов `learner-model` по коммиту, `CANDIDATE_LEARNER_MODEL` (R34), прошёл: выбранный подтверждён на своём отложенном наборе, его полосы записаны, и база стоит на его релизе, — или выход записан, и база остаётся;
- итог — строкой «в реестр spec» для dpia 10.2 spec; эпик — в архиве, следующий предложен (R59).

| Статус | ID | Задача | Кто | План | Зависит от |
|---|---|---|---|---|---|
| [x] | T69 | Развернуть E10 «Суд круга справедливости модели» | Claude | — | — |
| [x] | T70 | Критерий суда круга справедливости — до первого кандидата | Claude, проверяю я | да | — |
| [x] | T71 | Цель `M11.same` и набор суда в коде, тег | Claude | да | T70 |
| [ ] | T72 | Суд `model-fair-1` | я + Claude | да | T71; кандидаты — черновые pull request T62 `learner-model` (строка 4.17 COORDINATION.md spec), не их релиз |
| [ ] | T73 | База — на релиз выбранного, архив эпика | я + Claude | да | T72; при выборе — релиз выбранного, T63 `learner-model` |
| [x] | T74 | Скорость CI: строка 7.19 | Claude | — | — |
| [x] | T75 | Переход на релиз `.github` строки 7.13 | Claude | да | T74; релиз `.github` строки 7.13 COORDINATION.md spec |
| [ ] | T76 | Переход на тег `spec` строки 7.14: стадия `toolchain-core` | я + Claude | да | T75; тег `spec` строки 7.14 COORDINATION.md spec; релиз `.github` с T66 `.github` |

**Почему такой порядок.**
- Критерий — первым и отдельно от кода, как T61 в [docs/plan/E9.md](docs/plan/E9.md): я вычитываю его, пока ничто его не гоняет, и он ничего не ждёт в `learner-model`. Документы тега не выпускают (R29): тег с критерием — тег T71, и его ждут T60 `learner-model` — «Gate» на теге суда — и разведка T61 `learner-model`, которая читает наборы `sim judge` критерием E10 (sources 3.10).
- Код и тег — своей задачей: `M11` нет ни в одной группе счёта, у него нет пути и цели (R55), так что `sim judge` цели `M11.same` сейчас не считает; строка набора суда — в `internal/seed/sets.go` (R52). Мерж выпускает тег (R29), и с него до выбора закрепление не меняется.
- Суд — по коммитам черновых pull request кандидатов, не по их релизу: `sim` гонит кандидата из его коммита (R34), а релиз выбранного ждёт полос `sim` (R103 spec).
- База на релиз — после суда, как T59 в [docs/plan/E7.md](docs/plan/E7.md) и T48 в [docs/plan/E5.md](docs/plan/E5.md); при выходе задача только архивирует эпик. Отдельной задачи итога нет: последняя архивирует (R59).
- T74–T76 — вне эпика: скорость CI по строке 7.19 COORDINATION.md spec, слово автора 2026-11-16 (R61, R67). Ни бит, ни закрепления они не двигают; T76 переводит `spec` на тег и меняет `.devcontainer/`, и до выбора суда — только с теми же битами (R45).

## Epic E10. Суд круга справедливости модели

**Почему эпик такой.** Оценка `fairness-1` нашла: база — `learner-model` v0.1.8 и `engine` v0.1.6 — не держит критерий справедливости по классу: у `D-grade1` и `D-grade5` `M11.same` 0,812 [0,795; 0,829] и 0,5525 [0,53; 0,5735] при границе 0,90 (R58). `fair` держит кандидата «не хуже базы», но не награждает за лучшее, так что круг, который подтягивает базу к границе, судится со своей целью, записанной до первого кандидата (R58, R59). `learner-model` взял круг эпиком E8: разбор T59, «Gate» на теге `sim` с критерием E10 — T60, разведка T61, кандидаты — черновые pull request T62, суд и выпуск — T63, итог — T64 (sources 3.10).

**Общее для задач:**
- общее для задач E5, E7 и E9 действует и здесь ([docs/plan/E5.md](docs/plan/E5.md), «Эпик E5»; [docs/plan/E7.md](docs/plan/E7.md), «Эпик E7»; [docs/plan/E9.md](docs/plan/E9.md), «Эпик E9»): тесты и `codecov/patch`, `same-bits`, файл релиза — только после `release-verify`, язык; файлы решающих прогонов — в `runs/`; время — с load average хоста, а замер, по которому что-то выбирается, — на свободной машине;
- **критерий — до первого кандидата:** пока T70 не смержена, ни один кандидат круга не гоняется в `sim`, даже грубым взглядом; числа критерия после первого прогона кандидата не меняются (R103 spec; R163 standalone);
- **биты 68 ячеек не двигаются до выбора суда:** T71 меняет чтение результатов, а не прогон, так что файл базы — `runs/base-t64.jsonl`, чей SHA-256 называет `paths.json` (R58), — и 69 чисел гейта прежние. Сдвигает их только выбранный и подтверждённый кандидат — в T72 и T73;
- **одно закрепление:** от тега T71 до выбора суда `go.mod`, `gate.json`, пакеты и пути не меняются; релиз `spec`, `engine` или `modules`, который выйдет за это время, ждёт конца суда и берётся своей задачей. Исключение — тег `spec` с теми же битами (R45);
- **отложенный набор суда** не гоняет никто до подтверждения T72: ни разведка `learner-model`, ни грубый взгляд (R49, R52; раздел 3.8 [docs/architecture.md](docs/architecture.md));
- **в прогоне оркестратора `spec`** каждая задача — своя ветка от `main` и свой pull request, мержит оркестратор на зелёных проверках (R156 spec; правило 7); вопросы плана — строки [run/QUESTIONS.md](run/QUESTIONS.md);
- **не входит:** кандидат `engine` — он идёт своей причиной (R59); числа забывания `learner-model` (sources 3.10).

## Задачи эпика

#### T72. Суд `model-fair-1`

- **Цель:** кандидаты круга справедливости `learner-model` осуждены по критерию T70, R65, на теге T71. Выбранный подтверждён на нетронутом отложенном наборе суда, и его полосы записаны, — или выход записан, и база остаётся ([docs/architecture.md](docs/architecture.md), разделы 3.6 и 3.8). Итог — строкой «в реестр spec» для dpia 10.2 spec.
- **Ждёт:** T71; черновых pull request кандидатов T62 `learner-model` (строка 4.17 COORDINATION.md spec) — их коммитов, не релиза. Если разведка T61 `learner-model` кончится выходом без кандидатов, исполнитель предлагает задаче `[-]` строкой [run/QUESTIONS.md](run/QUESTIONS.md), итог несёт T64 `learner-model` (sources 3.10), а T73 только архивирует эпик.
- **Сделать:** как T58 в [docs/plan/E7.md](docs/plan/E7.md) и T47 в [docs/plan/E5.md](docs/plan/E5.md):
  - тег T71 — поправкой к R65; кандидаты — коммиты черновых pull request `learner-model`, имена, цели и число констант — из их описаний (R25);
  - решающий прогон каждого — `CANDIDATE_LEARNER_MODEL=<sha> just decisive engine runs/<файл>` (R34) — на свободной машине, с load average хоста; база — `runs/base-t64.jsonl`, который называет `paths.json`;
  - `sim judge` по критерию T70, R65; файл суда, отчёт и SHA-256 прогонов — в `judgements/model-fair-1/` (R35);
  - если выбран — подтверждение, один раз, на отложенном наборе суда; если подтвердился — `gate.json` называет коммит кандидата, `go.mod` — его псевдоверсию; пути заново на нём — поправкой к R04 и подразделом раздела 3 [docs/measurements.md](docs/measurements.md); полосы — `just gate record` с причиной — отчётом суда; мерж выпускает тег, после которого `learner-model` выпускает выбранного (T63 `learner-model`; R103 spec);
  - если нет — выход: база остаётся, отчёт записан; набор больше не отложен, если подтверждение шло (R52);
  - итог — решением журнала, как R46 и R53; строка «в реестр spec» — исход для строки 4.18 и для item 10 dpia 10.2 spec: `M11` выбранного у `D-grade1` и `D-grade5` с интервалами подтверждения и держит ли он границу R55, или выход с числами суда.
- **Решить при планировании:** что делать, если решающий прогон на свободной машине дольше 45 минут (раздел 3 [docs/architecture.md](docs/architecture.md), «What stays open»); повторять ли потолок на коммите выбранного; кандидату с иным стартом по классу, чем полоса класса (R51 spec), — ждать ли ответа `spec` на его строку до суда или только до выпуска.
- **Риски:** σ₀ 3,0 в разведке T30 `learner-model` срывала 28 ограничений, хотя `off-band` была лучше (sources 3.10): кандидат, что поднимает `M11.same`, может сорвать «не хуже» на других ячейках. Вероятен выход; решает суд.
- **Мой шаг:** решения; чтение отчёта суда; свободная машина на время решающих прогонов.
- **Как проверить:** часть прогона кандидата воспроизводится до байта на том же коммите (R28); `sim judge` не отказывает (R40); `just gate exact`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** строка 4.18 COORDINATION.md spec; R28, R34, R35, R40, R46, R52, R53, R55, R58; R65; [docs/architecture.md](docs/architecture.md), разделы 3.6–3.9; README, «Judging a candidate»; dpia 10.2 spec; sources 3.10.

#### T73. База — на релиз выбранного, архив эпика

- **Цель:** при выборе — `gate.json` и `go.mod` называют релиз выбранного, пути прочитаны на нём, и база оценена заново, `evaluations/fairness-2/` (R55); в обоих исходах — эпик в архиве и следующий предложен (R59).
- **Ждёт:** T72; при выборе — релиза выбранного, T63 `learner-model`.
- **Решить при планировании:** при выборе — что делать, если решающий прогон оценки на свободной машине дольше 45 минут, как в T72.
- **Мой шаг:** при выборе — свободная машина на время решающего прогона оценки.
- **Сделать:**
  - при выборе, как T48 в [docs/plan/E5.md](docs/plan/E5.md): `release-verify` принимает релиз (R77 spec); `gate.json` и `go.mod` — на нём; `just gate exact` — числа те же, что на коммите кандидата; пути заново (R40); оценка базы — решающий прогон, `sim evaluate`, `evaluations/fairness-2/` (R55), и её итог — строкой «в реестр spec» для item 10 dpia 10.2 spec; мерж выпускает тег (R29);
  - в обоих исходах: раздел 3 [docs/architecture.md](docs/architecture.md) сверить с построенным; что меняет карточку architecture 3.9 spec — строкой «в реестр spec»; таблицы 2–4; эпик — в архив `docs/plan/E10.md`, замер токенов его сессий — строкой measurements 2.2; следующий эпик предложить (R59).
- **Как проверить:** при выборе — `just gate exact`, `just ci-test`, `just ci-lint`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R04, R26, R29, R40, R55, R58, R59; [docs/architecture.md](docs/architecture.md), разделы 3, 3.6 и 3.9; architecture 3.9 spec; dpia 10.2 spec.

#### T76. Переход на тег `spec` строки 7.14: стадия `toolchain-core`

- **Цель:** `sim` — на теге `spec` строки 7.14 COORDINATION.md spec, со стадией `toolchain-core` в `.devcontainer/Dockerfile` и apt шаблона только из источников снимка; проверки и релиз, которые общий `go.yml` релиза строки 7.13 гоняет в `toolchain-core`, идут в ней; те же биты. Последняя задача строки 7.19: она одна меняет `.devcontainer/`.
- **Почему:** строки 7.13, 7.14 и 7.19 COORDINATION.md spec; переход на тег `spec` — своей задачей и решением журнала (RULES.md, Across repositories), как E11 (R60).
- **Ждёт:** T75; тега `spec` строки 7.14; релиза `.github` с T66 `.github` — последней частью строки 7.13: стадия `toolchain-core` в релизе и проверках, «`release.yml` строит стадию `toolchain-core` Dockerfile вызывающего, без кеша, а проверки без браузера и Go-тестов — lint, vuln, licenses и документы — идут в ней; у Dockerfile без этой стадии — как прежде» (sources 5.4). T66 ждёт перехода `.github` на тот же тег `spec`, а он меняет `.devcontainer/` `.github`, и по решению координатора `spec` (T171 spec) ночью `.devcontainer/` не меняется, поэтому T75 взяла `v0.1.30` без T66 (R68). По правилу «одно закрепление» эпика E10 тег `spec` до выбора суда берётся, только если биты те же (R45): `just gate exact` зелёный и пути прежние; иначе — после T73.
- **Сделать:** как переход E11 (R60): подмодуль `spec/` — на тег; `.devcontainer/` — по шаблону вида `go` этого тега, со стадией `toolchain-core`; каждый `uses:` `.github` в `.github/workflows/` — на SHA релиза с T66, как T75 (R63, R68); шаги образа в `gate.yml` — по `go.yml` этого релиза: если T66 меняет путь файла образа или его стадию, «A change to how either builds the image — `toolchain-core` … — takes a new path in all of them at once» (R42 .github), — тот же путь и ключ; рецепты `test-targets`, `fma-diff` и `rules` — сверить с job его `go.yml`; модуль Go `spec` в `go.mod` и в `gate.json` — на тег, если он его выпускает, и `just gate exact` — те же 69 чисел; что тег меняет в правилах и документах — по его журналу; решение журнала о переходе.
- **Мой шаг:** пересобрать контейнер после мержа — «Rebuild Container»: в прогоне оркестратора ночью пересобрать его некому, поэтому `.devcontainer/` меняет только эта задача.
- **Как проверить:** `just gate exact`; `just ci-lint`, `just ci-test`, `just same-bits`; `just -f spec/justfile docs-check-repo`; сборка стадий `toolchain` и `toolchain-core` из Dockerfile; CI pull request зелёный; время изменённых job до и после — в отчёте; `/code-review high`, один круг.
- **Разделы:** строки 7.13, 7.14 и 7.19 COORDINATION.md spec; measurements 3.1 spec; R76 spec, R156 spec; R45, R60, R63, R67, R68; sources 5.4.
