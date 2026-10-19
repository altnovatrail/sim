# sim — архитектура

> Статус: каркас из заготовки `spec`, 2026-10-19.
>
> Что делает компонент, чего не делает и что о нём ещё предстоит решить. Это подробности его карточки — architecture 3.9 spec. Карточка остаётся сводкой: если подробности меняют то, что в ней сказано, она правится в `spec` тем же путём, что и вопросы в его реестры (RULES.md, Across repositories).
>
> Почему решено так, а не иначе, — в [decisions.md](decisions.md); план работ — в [RUN.md](../RUN.md). Понятия — concepts spec. Документы и решения `spec` лежат в `spec/docs/` и цитируются с суффиксом «spec».

## 1. Назначение и границы

**Делает:**
- гоняет `learner-model` и `engine` на тысячах виртуальных учеников разных профилей: виртуальный ученик отвечает на задачи, модель ученика оценивает его владение навыками, движок выбирает, что дать дальше (architecture 3.9 spec; concepts 11, 12 spec). В работе платформы граф и свидетельства библиотекам передаёт `api` (architecture 3.8 spec), а на стенде — сам `sim`;
- меряет то, что платформа обещает о модели ученика и движке:
  - коридор успеха — задача решается с вероятностью от 0,70 до 0,85 (concepts 12 spec);
  - осторожное «освоено» (concepts 11 spec);
  - тот же вход даёт тот же выход, и у каждого выбора есть код причины (concepts 12 spec).

  В первом срезе — на графах обоих модулей: движок держит коридор успеха, не зацикливается и не объявляет «освоено» раньше времени (product 10 spec);
- держит критерий, по которому выбирают новое правило, и полосы допуска (architecture 3.9 spec). Критерий записывается до первого прогона кандидата (architecture 3.8 spec; concepts 11 spec; R163 standalone);
- служит гейтом CI (architecture 3.9 spec; repositories 4 spec):
  - каждое изменение `learner-model` или `engine` прогоняется на стенде, и изменение, которое ухудшает критерий, не сливается (repositories 4 spec). Как CI библиотеки берёт стенд, решает открытый вопрос Т-16 (раздел 2.3);
  - CI `modules` гоняет `sim`, если изменился граф навыков (architecture 3.10 spec; Т-16).

**Не делает:**
- не меняет алгоритмы, а только меряет их (architecture 3.9 spec). Модель ученика — дело `learner-model`, выбор следующего урока или задачи — `engine` (architecture 3.8 spec);
- не держит данных живых учеников (architecture 3.9 spec): ни ответов, ни работ, ни журналов `api`. Виртуальные ученики синтетические, как любые примеры в репозиториях NovaTrail (RULES.md, Hard rules);
- не обращается к ИИ: всё, что обращается к модели, — в `tutor`, а модель ученика и выбор считает код (architecture 1 spec; R03 spec);
- не читает модули через `module-kit` и не проверяет их: пакеты он читает по типам из `spec` (architecture 3.9 spec; repositories 2.1 spec), а модуль проверяет `module-kit` (architecture 3.7 spec);
- не работает в платформе: это инструмент CI, а не часть сервиса (repositories 1, 2.1 spec).

**Данные и связи:**
- данные — пакеты модулей, закреплённые по версии: из них — графы навыков и задачи. Их `sim` читает по типам канонического вида из `spec`, без `module-kit` (architecture 3.9 spec). Что внутри пакета — открытый вопрос Т-15; его рекомендация — скомпилированный канонический вид, который читают `api`, `sim` и eval. На этапе 1 пакеты лежат в релизах `modules` (architecture 3.10 spec);
- данных живых учеников нет (architecture 3.9 spec);
- при сборке зависит от `learner-model` и `engine` (repositories 2.1, 3 spec). CI этих библиотек берёт `sim` по тегу — как инструмент CI, а не как зависимость при сборке, поэтому цикла нет (repositories 1 spec; Т-16). CI `modules` тоже берёт `sim` по точному тегу (repositories 2.1, 3 spec);
- время и зерно библиотеки получают снаружи (architecture 3.8 spec): на стенде их даёт `sim`. Генератор случайности и вывод зёрен задаёт `spec`, с тест-векторами (architecture 1 spec; Т-18);
- образец — `tools/learners` и `docs/learners.md` в standalone (architecture 3.9 spec; раздел 4 этого документа). Код MathTrail под MIT переносится свободно, уведомление MIT сохраняется (R12 spec);
- этап 1, вместе с `learner-model` и `engine`: стенду нужны и модель, и правило выбора, а стенд и критерий появляются до первого кандидата (repositories 6 spec; product 11 spec).

## 2. Решить при детализации

Пункты «Решить при детализации» из карточки — по подразделу на пункт. Каждый решает своя задача эпика E1 в [RUN.md](../RUN.md). Решённый пункт называет запись журнала, а само решение переезжает в раздел 3. Вопрос, который выходит за этот репозиторий, уходит «в реестр spec» (RULES.md, Across repositories).

У карточки два пункта — подразделы 2.1 и 2.2. Подраздел 2.3 — сверх её пунктов: карточка называет Т-16 среди того, что `sim` делает, — гейт CI, — но не среди того, что решить (architecture 3.9 spec). Его называет задача T08 `spec`, которая завела этот репозиторий. Т-16 — вопрос реестра `spec`, и его решение касается ещё `learner-model`, `engine` и `modules`; поэтому то, что его закроет, уходит «в реестр spec».

### 2.1. Профили виртуальных учеников

- **Суть:** какие виртуальные ученики есть на стенде — как каждый отвечает, как меняется его владение навыками, как он занимается во времени — и на каких графах навыков он идёт. Профиль — это допущение о живых учениках: стенд его не мерит, а принимает, и от набора профилей зависит, что критерий вообще может увидеть. Решается здесь: это пункт карточки (architecture 3.9 spec).
- **Что известно:**
  - в standalone виртуальные ученики — синтетические дети 1–6 класса, каждый из своего зерна, и одни и те же дети проходят каждое правило. Генераторов девятнадцать, и каждый меняет одно допущение базовой популяции: G0 стоит на месте, G1 размещён на уровень мимо, G2 быстро учится, G3 делает скачок и так далее. Основной — G2-half, ученик со скоростью, которую автор MathTrail считает правдоподобной. Каждый ребёнок даёт 200 ответов, а каждый десятый ответ начинает новый день (раздел 4);
  - там же: «How fast real children learn is not known: the speeds are the bench's, not a measurement» (раздел 4);
  - ученик NovaTrail — любого возраста, и в первом срезе два непохожих модуля: ребёнок пишет от руки на листе, взрослый правит текст на текстовой канве (product 10 spec; R06 spec, R08 spec). Второй модуль окончательно не выбран (О-1). Профили standalone — только дети и только математика;
  - модель ученика и движок NovaTrail делают то, чего не делали в standalone, и профили должны это задеть:
    - забывание: без практики оценка тает, и движок возвращает навык на повторение. В standalone забывания не было (concepts 11, 12 spec). Время библиотеки получают снаружи (architecture 3.8 spec), поэтому профилю нужен ритм занятий и перерывы;
    - пререквизиты: навык предлагается, когда освоены его пререквизиты. Правило standalone связей между темами не читало (concepts 12 spec);
    - уровень — на модуль или область, а не на всю модель (concepts 11 spec; architecture 3.8 spec): ученик может учить два модуля сразу;
    - угадывание зависит от вида ответа — 1/k при выборе из k вариантов, около нуля при свободном ответе; в standalone оно было зашито как 0,2 (concepts 11 spec);
    - свидетельство — исход попытки: с первой отправки, после подсказки номер k, с наброском-подсказкой, решение показано, без проверки. Серия однотипных пунктов на одном листе весит меньше суммы (concepts 11 spec);
  - задачи стенда — из пакетов модулей (architecture 3.9 spec). Сложность задачи задаёт модуль по своей шкале (concepts 6 spec). В standalone задачу писала модель чата, и генераторы задавали её промах — насколько задача на самом деле сложнее или легче заказанной (раздел 4);
  - виртуальные ученики синтетические, данных живых учеников на стенде нет (architecture 3.9 spec; RULES.md, Hard rules);
  - пакетов ещё нет: формат модуля v0 и пакет решает эпик E2 в RUN.md `spec` (Т-7, Т-15).
- **Варианты:** пока пусто — их приносит задача.
- **Задача:** T01.

### 2.2. Критерий и полосы допуска

- **Суть:** по какому правилу новое правило модели ученика или движка признаётся лучше прежнего или хотя бы не хуже — это критерий, — и насколько числа могут сдвинуться на каждом изменении, чтобы гейт CI его пропустил, — это полосы допуска. Записать их нужно до первого прогона кандидата (R163 standalone): иначе числа кандидата влияли бы на то, чем его судят. Решается здесь: это пункт карточки (architecture 3.9 spec).
- **Что известно:**
  - каждое изменение `learner-model` и `engine` проходит стенд по критерию, записанному до первого прогона кандидата, а изменение, которое ухудшает критерий, не сливается (architecture 3.8 spec; repositories 4 spec; concepts 11 spec);
  - что стенд должен показать в первом срезе: на графах обоих модулей движок держит коридор успеха, не зацикливается и не объявляет «освоено» раньше времени (product 10 spec). Коридор — от 0,70 до 0,85, с целью 0,775 (concepts 12 spec); «освоено» — осторожная оценка (concepts 11 spec); воспроизводимость и код причины у каждого выбора проверяет стенд (concepts 12 spec);
  - в standalone это два разных механизма (раздел 4):
    - **критерий** выбирает новое правило. В нём ограничения, оценка, выбор среди кандидатов, подтверждение на отложенных зёрнах и выход, если ни один кандидат не прошёл. Оценка шага рейтинга — доля пути от действующего правила к «оракулу», который знает, где ребёнок на самом деле; у правила «освоено» и у «долга» — доля пути от базы к совершенству, потому что для нового правила «освоено» оракул не потолок. Решающий прогон — 4 000 детей на ячейку (R163 standalone);
    - **страж** (guard) — гейт CI на каждом pull request. Он гоняет действующее правило на первых 300 детях генераторов G0–G3 и держит каждое число в пределах двух полуширин 95 % интервала, записанного рядом с ним. Полосы записываются заново, когда изменение намеренно двигает модель, и новый файл полос читается на ревью;
  - допуск проверки «не хуже» в standalone — больший из допуска автора и разрешения стенда, 4,3 стандартной ошибки: допуск тоньше разрешения проваливал бы равного кандидата случайно (R163 standalone);
  - критерии standalone выбирали шаг рейтинга, правило «освоено» и то, сколько неверный ответ стоит ребёнку на экране, — «долг». Задачу у всех правил выбирал один и тот же код сервиса, так что правило выбора они не судили (раздел 4). У NovaTrail это два репозитория, и `engine` нужен свой критерий;
  - отправная модель NovaTrail — модель standalone; отличия — уровень на модуль и забывание (concepts 11 spec; architecture 3.8 spec). С чем сравнивается первый кандидат и есть ли у стенда свой «оракул», ещё не решено;
  - стенд и критерий появляются до первого кандидата (repositories 6 spec).
- **Варианты:** пока пусто.
- **Задача:** T02.

### 2.3. Как CI библиотек берёт стенд (Т-16)

- **Суть:** как CI `learner-model`, `engine` и `modules` получает `sim`, подставляет свой код и читает гейт; сколько гейт идёт; в каком порядке выходят релизы, когда изменение ломает API библиотеки, на котором собран `sim`, или намеренно двигает числа и полосы надо записать заново в другом репозитории. Карточка этого пункта среди того, что решить, не называет, но называет Т-16 среди того, что `sim` делает (architecture 3.9 spec). Решается вместе с `spec`: Т-16 — вопрос его реестра.
- **Что известно:**
  - варианты Т-16: CI библиотеки берёт `sim` по точному тегу и подставляет код pull request через `go.work`; гейт только на кандидатах релиза; три библиотеки в одном репозитории. Рекомендация — `go.work`: репозитории остаются раздельными (R10 spec), а каждый pull request проходит стенд. CI `modules` тоже гоняет `sim`, если меняется граф (Т-16);
  - сборка, сеть и CI — разные связи: гейт `sim` в CI библиотек — не цикл, хотя `sim` их импортирует (repositories 1, 3 spec);
  - `engine` зависит от `learner-model` при сборке (repositories 2.1, 3 spec), так что `sim` собирается с обоими, а pull request подменяет одну из них;
  - в standalone стенд и продукт жили в одном репозитории: стенд — свой модуль Go, а продукт в нём подключён через `replace ../..` (R155 standalone). Страж шёл отдельной задачей CI на каждом pull request, с таймаутом 10 минут, а сам на машине разработки занимал несколько секунд (раздел 4). В NovaTrail стенд и библиотеки — в разных репозиториях (R10 spec, R18 spec);
  - полосы стража в standalone записаны на amd64 при `GOAMD64` v1 и на процессоре с FMA. На arm64, на amd64 при `GOAMD64` v3 и выше и на процессоре без FMA последние цифры чисел могут разойтись, хотя и в пределах полос (раздел 4);
  - общие workflow — Go, релиз с подписью, проверка лицензий — лежат в `.github` с actions, закреплёнными по SHA (architecture 3.11 spec). Какие из них общие, а какие свои, решает детализация `.github`.
- **Варианты:** пока пусто.
- **Задача:** T03.

## 3. Устройство

Пока пусто. Сюда переезжает то, что решено по разделу 2: из каких частей компонент состоит, его типы и данные, как он связан с другими репозиториями.

## 4. Источники

Внешний факт — условия поставщика, цена, закон, документация библиотеки — встаёт сюда с источником и датой прочтения (RULES.md, The bar).

Репозиторий standalone, ветка `main`, коммит c6cbeb2, прочитано 2026-10-19 — тот же коммит, что читала architecture 8 spec:
- [`docs/learners.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/learners.md) — образец стенда (разделы 2.1–2.3):
  - «What it is, and where it comes from»: «every tenth answer starts a new day. A child gives 200 answers»;
  - «The children»: дети синтетические, каждый из своего зерна, и одни и те же дети проходят каждое правило; таблица девятнадцати генераторов, G2-half — «The main learner, the speed the author takes as plausible (2026-10-03)»; промах модели чата — задача «truly harder or easier than the difficulty asked for»; «How fast real children learn is not known: the speeds are the bench's, not a measurement»;
  - «The rules»: «the choice of task, the trial series and mastery are the service's for every rule, and only the estimate differs»; «The oracle … the ceiling of every rule that learns of a child from answers»; «It is not a ceiling for a new rule of mastery»;
  - «Choosing mastery»: «The score: the share of the way from the baseline to perfect»;
  - «The criterion»: ограничения, оценка, выбор, подтверждение на отложенных зёрнах и выход; «A candidate is chosen on the working seeds, in a decision run of 4,000 children a cell, and confirmed once on the held-out seeds»;
  - «Choosing how a wrong answer is charged»: «Its aim is the screen, not the estimate: a single wrong answer should cost the child little on the card»;
  - «The guard»: «It runs the service's own path on the first 300 children of each of G0, G1, G2 and G3 … and holds it within two half-widths of the 95 % interval recorded with it»; «CI runs it on every pull request, as a job of its own»; полосы записываются заново, когда изменение намеренно двигает модель, и «the new `guard.csv` is read in the change's review»; «The bands are recorded on amd64, built at `GOAMD64` v1, Go's default, where Go fuses no multiplication with an addition, and on a processor with FMA … elsewhere the last digits of the numbers may part from them, well within any band»;
  - «Reading the live numbers»: «Two tables of the service's report hold them to live ones»; стенд читает своих детей так же, «and this is how the two are laid side by side»;
  - «The paper's numbers»: «On arm64 the last digits may differ, since Go fuses a multiplication and an addition there into one step that rounds once, as it does on amd64 at `GOAMD64` v3 and above; so may they on a processor without FMA»;
  - «Running it»: «`just ci-learners` runs the guard, a few seconds on the development machine».
- [`docs/decisions.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/decisions.md):
  - R163 standalone: «A new rule of the student model is chosen by a criterion written down before any candidate ran, on a bench that knows its ceiling». Там же — допуск «never under 4.3 standard errors of the paired difference» и почему: допуски автора «are finer than the bench resolves on 1,000 children»;
  - R155 standalone: «`tools/learners` is a Go module of its own … the product through `replace ../..` and no dependency of its own».
- [`.github/workflows/ci.yml`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/.github/workflows/ci.yml), задача `learners-guard` — «Student model guard», `timeout-minutes: 10`, шаг `just ci-learners | tee -a "$GITHUB_STEP_SUMMARY"`.
- [`LICENSE`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/LICENSE): «MIT License», «Copyright (c) 2026 MathTrail».
