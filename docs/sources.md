# sim — sources

> The external facts the repository's documents under `docs/` rest on: a provider's terms, a price, a law, a standard, a library's documentation. Here are their records — the address, the version or date of the page, the date it was read and a quote in the source's language — and [architecture.md](architecture.md) and the decision log hold the claim and a citation of the record by its number, as a section: `sources 1.1` (RULES.md, The bar; R143 spec).
>
> A record is a subsection "N.M" of a topic: one source, or one reading of several sources for one decision. A record's number never changes: a new one goes at the end of its topic, a new topic at the end of the file. A source read again is an item with its date in the same record, and the earlier text is not edited. The model is sources spec. The results of measurements are not here but in [measurements.md](measurements.md).

## 1. MathTrail standalone

### 1.1. standalone c6cbeb2: the example for the simulator

The repository standalone, branch `main`, commit c6cbeb2, read 2026-10-19 — the same commit that architecture 8 spec read.

- [`docs/learners.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/learners.md) — the example for the simulator (architecture 2.1–2.3 and 3):
  - "What it is, and where it comes from": "every tenth answer starts a new day. A child gives 200 answers";
  - "The children": the children are synthetic, each from a seed of its own, and the same children go through every rule; a table of nineteen generators, G2-half — "The main learner, the speed the author takes as plausible (2026-10-03)"; the chat model's miss — a task "truly harder or easier than the difficulty asked for"; "How fast real children learn is not known: the speeds are the bench's, not a measurement";
  - "The rules": "the choice of task, the trial series and mastery are the service's for every rule, and only the estimate differs"; "The oracle … the ceiling of every rule that learns of a child from answers"; "It is not a ceiling for a new rule of mastery";
  - "Choosing mastery": "The score: the share of the way from the baseline to perfect";
  - "The criterion": the constraints, the score, the choice, the confirmation on the held-out seeds and the exit; "A candidate is chosen on the working seeds, in a decision run of 4,000 children a cell, and confirmed once on the held-out seeds";
  - "Choosing how a wrong answer is charged": "Its aim is the screen, not the estimate: a single wrong answer should cost the child little on the card";
  - "The guard": "It runs the service's own path on the first 300 children of each of G0, G1, G2 and G3 … and holds it within two half-widths of the 95 % interval recorded with it"; "CI runs it on every pull request, as a job of its own"; the bands are recorded anew when a change moves the model on purpose, and "the new `guard.csv` is read in the change's review"; "The bands are recorded on amd64, built at `GOAMD64` v1, Go's default, where Go fuses no multiplication with an addition, and on a processor with FMA … elsewhere the last digits of the numbers may part from them, well within any band";
  - "Reading the live numbers": "Two tables of the service's report hold them to live ones"; the simulator reads its own children the same way, "and this is how the two are laid side by side";
  - "The paper's numbers": "On arm64 the last digits may differ, since Go fuses a multiplication and an addition there into one step that rounds once, as it does on amd64 at `GOAMD64` v3 and above; so may they on a processor without FMA";
  - "Running it": "`just ci-learners` runs the guard, a few seconds on the development machine".
- [`docs/decisions.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/decisions.md):
  - R163 standalone: "A new rule of the student model is chosen by a criterion written down before any candidate ran, on a bench that knows its ceiling". In the same place — the tolerance "never under 4.3 standard errors of the paired difference", and why: the author's tolerances "are finer than the bench resolves on 1,000 children";
  - R155 standalone: "`tools/learners` is a Go module of its own … the product through `replace ../..` and no dependency of its own".
- [`.github/workflows/ci.yml`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/.github/workflows/ci.yml), the job `learners-guard` — "Student model guard", `timeout-minutes: 10`, the step `just ci-learners | tee -a "$GITHUB_STEP_SUMMARY"`.
- [`LICENSE`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/LICENSE): "MIT License", "Copyright (c) 2026 MathTrail".

### 1.2. standalone 89986bc: the populations, the criterion and the guard

The repository standalone, branch `main`, commit 89986bc, read 2026-11-07 — the commit that `modules` carries the tasks over from (R05 modules). Between c6cbeb2 and 89986bc, in `tools/learners` and `docs/learners.md` only `go.mod` and `go.sum` changed — Go 1.27.3 and new versions of `golang.org/x` — and the numbers of `testdata`: the assumptions of the populations are the same (architecture 3.1–3.4).

- [`docs/learners.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/learners.md):
  - "The children": the base population — "a grade from 1 to 6, and the start the service gives that grade; a true level of the start plus N(0, 1); a true offset per topic of N(0, 0.5²); and a task that is truly harder or easier than the difficulty asked for, by N(0, 0.5²) … A child answers correctly with chance 0.2 + 0.8·σ(level + offset − difficulty), and never asks for the hint"; "The generators each change the base population in one respect"; "A generator's children are drawn from seeds named after it, so a generator added leaves every child of the others as they were"; a table of nineteen populations — G0, G1, G2, G3, G4, G5, G6, G7, G8, G0-exact, G0-miss0.25, G0-miss1, G2-half, G2-fading, G3-drop, G0-start0.5, G0-start2, G0-topics0.3, G0-topics1;
  - "What it is, and where it comes from": the simulator goes through the product's path — `profile.New`, `tutor.Next`, `Profile.Ask`, `Profile.Issue`, `Profile.Record`; "A task is answered two minutes after it is issued, the next is issued at once, and every tenth answer starts a new day";
  - "The held-out seeds": the held-out children — "seed 20261003, experiment `held-out`"; "A test proves that no child of any generator is drawn alike in any two of the four sets"; "A set that has been looked at is held out no longer".
- [`tools/learners/child.go`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/tools/learners/child.go): the constants of the populations — `spreadAroundStart = 1.0`, `spreadOfTopics = 0.5`, `writingError = 0.5`, `baseFloor = 0.2`, `misplacedBy = 2.5`, `learningPerAnswer = 0.01`, `topicLearning = 0.02`, `forgetting = 0.01`, `jumpBy = 1.0`, `jumpEarliest = 50`, `jumpLatest = 150`, `harderHostBias = 0.75`, `groupOffset = 0.75`, `spreadInGroup = 0.25`, `hintShare = 0.2`, `fadeAfter = 50.0`; `G2-half` — `c.learn(learningPerAnswer/2, topicLearning/2)`; the miss is drawn at every issue — `writing: c.draws.NormFloat64()` in `answerDraws`; `seeded(unit, purpose)` — FNV-64a of "`%d|%s|%s|%s`" of the main seed, the experiment, the unit and the purpose, and PCG.
- [`tools/learners/loop.go`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/tools/learners/loop.go): `answerTime = 2 * time.Minute`, `answersPerDay = 10`.

The repository standalone, commit 89986bc, read 2026-11-07 in the planning of T02 — the criterion and the guard that were followed in writing architecture 3.5–3.8.

- [`docs/learners.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/learners.md):
  - "The measures": "A topic is truly mastered at a level when the child's true chance on a task of that level and topic at difficulty 3, written as asked, is at least 0.775, the corridor's middle"; "Every number is read off all the children of its cell, with a 95 % interval from 2,000 resamples of the children";
  - "The criterion": the numbers of the targets and the tolerances "come from the service's rule and the rules already on the bench, never from a candidate: they are calibration, not a choice"; "A number the run cannot read leaves its constraint unmet and its score unread";
  - "Choosing the step": the targets — "on G2-half, the lag, `|r6_lag|`, is at most 45 % of the service's", "c* = (0.35 − 0.257) / (0.605 − 0.257) = 0.27", "on G3, the share of children not caught up after the jump, `r6_jump_unsettled`, is at most 25 %"; the tolerance — "4.3 standard errors of the paired difference between the slow constant step and the service there, in the same run", and "A candidate exactly as good as the service passes such a check 99 times in 100"; the score — the groups "Stays put" 1, "Learns" 2 (G2-half 1.25, G2-fading 0.5, G2 0.25), "Changes at once" 1 (G3 0.5, G3-drop 0.5), "Placed a level off" 0.5 (G1); "The chosen one is the simplest of A* and its equals"; "no second candidate is tried on the held-out seeds, since a set tried twice is held out no longer";
  - "Choosing mastery": "false masteries are at most 20 %, on G0 and on G0-topics1"; the number of answers until "mastered" — "at most 1.5 times the baseline's";
  - "The guard": "holds it within two half-widths of the 95 % interval recorded with it"; "the edges included"; "a band is room for a change, not for chance"; "fails naming each number outside its band with the band itself"; "the new `guard.csv` is read in the change's review, as a snapshot is, since recording turns a red guard green".

Read again 2026-11-12 for T51: "The measures" — "The names are those of `cells.csv`".

### 1.3. standalone 6b29dd9

The repository standalone, branch `main`, commit 6b29dd9, read 2026-11-07: since 89986bc, in `tools/learners` and `docs/learners.md` only `tools/learners/go.sum` has changed, so architecture 3.1–3.8 read 89986bc, sources 1.2.

## 2. Go, GitHub and other external sources

### 2.1. Go Modules Reference: workspaces

Read 2026-11-07 for architecture 3.9, in the planning and the execution of T03.

[Go Modules Reference](https://go.dev/ref/mod): "A workspace is a collection of modules on disk that are used as the main modules when running minimal version selection (MVS)"; "A `use` adds a module on disk to the set of main modules in a workspace"; "MVS starts at the main modules (special vertices in the graph that have no version) and traverses the graph, tracking the highest required version of each module"; "replace directives in go.work files override any replaces of the same module or module version in workspace modules"; "If GOWORK is set to off, the command will be in a single-module context"; "It is generally inadvisable to commit go.work files into version control systems … CI systems should generally not be allowed to use the go.work file so that they can test the behavior of the module as it would be used when required by other modules". That the build fails when the API of a module in `use` has parted from what the simulator requires, the documentation does not say in words — it is a conclusion of architecture 3.9.

### 2.2. Go: floating-point numbers and `GOAMD64`

Read 2026-11-07 for architecture 3.9, in the planning and the execution of T03.

[The Go Programming Language Specification](https://go.dev/ref/spec), "Floating-point operators": "An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements, and produce a result that differs from the value obtained by executing and rounding the instructions individually". [Go Minimum Requirements](https://go.dev/wiki/MinimumRequirements#amd64): "GOAMD64=v1 (default)", and FMA from "GOAMD64=v3".

### 2.3. GitHub Docs: runners, limits and Dependabot

Read 2026-11-07 for architecture 3.9, in the planning and the execution of T03.

GitHub Docs: [GitHub-hosted runners](https://docs.github.com/en/actions/reference/runners/github-hosted-runners) — for public repositories `ubuntu-24.04` is Linux x64 and `ubuntu-24.04-arm` is arm64, each with 4 CPU and 16 GB, and "Use of the standard GitHub-hosted runners is free and unlimited on public repositories"; [Limits](https://docs.github.com/en/actions/reference/limits) — "Each job in a workflow can run for up to 6 hours of execution time"; [Dependabot supported ecosystems](https://docs.github.com/en/code-security/dependabot/ecosystems-supported-by-dependabot/supported-ecosystems-and-repositories), "GitHub Actions" — "Dependabot only supports updates to GitHub Actions using the GitHub repository syntax, such as `actions/checkout@v6` or `actions/checkout@<commit>`".

### 2.4. The Go module `go.yaml.in/yaml/v4` v4.0.0

Read 2026-11-08 for the seeds and the distributions of T11 (R12).

The Go module `go.yaml.in/yaml/v4` v4.0.0 from the Go proxy: `LICENSE` — "Apache License, Version 2.0, January 2004"; `NOTICE` — some of the files are "still covered by their original MIT license", and "All the remaining project files are covered by the Apache license". `module-kit` takes the same version in its `go.mod`.

### 2.5. AS 241

Read 2026-11-08 for the seeds and the distributions of T11 (R12).

Wichura, M. J., "Algorithm AS 241: The percentage points of the normal distribution", Applied Statistics 37(3), 1988, pp. 477–484 — the coefficients PPND16. The article was not read here: the coefficients stand as its well-known publication gives them, and a test holds them — Φ(x) through `math.Erfc` returns p from 1e-300 to 1 − 1e-12 within an error of 1e-14 in x, and `inverseNormal(0.975)` = 1.959963984540054.

### 2.6. The Go proxy: the Origin of the libraries' tags

Read 2026-11-08 for the gate of T24 (R26).

The Go proxy, the Origin of the tags: `github.com/altnovatrail/learner-model@v0.1.0` — "Hash: d2620a18d865c813ec2556c93fca92a1445e3a74", "Ref: refs/tags/v0.1.0"; `github.com/altnovatrail/engine@v0.1.0` — "Hash: c6c5bb9bc2c6fc486f629da82ee834bd89fc2b35".

### 2.7. `actions/cache` v6.1.0

Read 2026-11-08 for the gate of T24 (R26).

[actions/cache](https://github.com/actions/cache), the release v6.1.0 of 2026-06-26, its tag on the commit 55cc8345863c7cc4c66a329aec7e433d2d1c52a9 (`gh api repos/actions/cache/releases/latest`).

## 3. `learner-model`

### 3.1. Commit 5db0c8c, T02

The repository [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model), branch `main`, commit 5db0c8c, read 2026-11-07: since 6ae7dec, which T01 read, only RUN.md and the entry R13 learner-model on the ruleset have changed. For architecture 3.5–3.8 these were read: R12 learner-model — the first version of the model as the base —, architecture 3.1 and 3.2 learner-model — the threshold 0.71875 and the measures of forgetting that `learner-model` proposes —, and the rows Д5–Д8 and В24 of the RUN.md of `learner-model`.

### 3.2. Commit a55d52e, T03

Read 2026-11-07 for architecture 3.9, in the planning and the execution of T03.

The repository [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model), branch `main`, commit a55d52e: the module `github.com/altnovatrail/learner-model` with the package `detmath`, no releases. R09 learner-model — «Библиотека даёт одни биты на любой сборке одной своей версии — компилятором Go, на любой архитектуре, уровне GOAMD64, процессоре и версии Go»; `.github/workflows/ci.yml`, the job `go-arm64` — "The shared go.yml runs on ubuntu-24.04 alone", `runs-on: ubuntu-24.04-arm`, and the steps `docker build --target toolchain` and `docker run … toolchain just ci-test`; T22 of the RUN.md of `learner-model` — «job в `ci.yml`: `sim` по точному тегу, код pull request — через `go.work`», and how the gate runs a change of both libraries at once — by T03 of `sim`.

### 3.3. Release v0.1.0

Read 2026-11-08 for the move to `spec` v0.9.8 (R07).

The repository [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model): the release v0.1.0, 2026-11-07 (`gh release list -R altnovatrail/learner-model`) — the first version of the model, the base of the model's criterion (architecture 3.6). The record of the commit 5db0c8c, sources 3.1, speaks of the time of T02, when there were no releases.

Read 2026-11-08 for the planning of the epic E2 (R09).

The repository [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model), the tag `v0.1.0`: `go.mod` without `require`; `func NewGraph(module string, bands []Band, skills []Skill, retired []RetiredSkill, tasks []Task) (Graph, error)`, `Append`, `Read`, `Readout.Task`; the package `detmath` — `Exp`, `Log`, `Pow`, `Sigmoid`, `Retention`, "The engine may take this package as it is".

Read 2026-11-08 for the seeds and the distributions of T11 (R12).

The release of `learner-model` v0.1.0: `just -f spec/justfile release-verify learner-model v0.1.0 learner-model-proto-v0.1.0.tar.gz …` — "learner-model-proto-v0.1.0.tar.gz: learner-model v0.1.0, signed by the shared release.yml of .github v0.1.15".

Read 2026-11-08 for the model's measures of T20 (R22).

The Go module `learner-model` v0.1.0, `model.go` — `chance`: "P = c + (1 − c) · R̄ · σ", and `guess`: c is the product of 1/k over the answer slots of a choice, 0 with an answer slot that is not a choice; `step.go` — `chanceOf`: "δ̄ and R̄ the mean δ and R of its J recorded skills, in the order of their IDs — 0 and 1 with none … In the trial series P is without retention, R̄ = 1", a skill without practice — R = 1 (`retention`).

### 3.4. Commit 50f1b51

Read 2026-11-08 for the move to `spec` v0.9.9 (R08).

The repository [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model), branch `main`, commit 50f1b51: its own job `go-arm64` in `.github/workflows/ci.yml` is still there, and `go.yml` is pinned at v0.1.15 — it has no input `same-bits` yet.

### 3.5. `ci.yml`: the job `go-arm64`, T24

Read 2026-11-08 for the gate of T24 (R26).

The repository [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model), branch `main`: `.github/workflows/ci.yml`, the job `go-arm64` — `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`, `docker build --target toolchain --tag toolchain - < .devcontainer/Dockerfile`, `docker run --rm --user "$(id -u):$(id -g)" --env HOME=/tmp --volume "$PWD:/work" --workdir /work toolchain just ci-test`.

### 3.6. Release v0.1.5, T42

Read 2026-11-11 for T42 (R41): the release [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model) `v0.1.5` — "Bring format 0.4 into the model: time and the leading question (T37, R28)", its tag on the commit `3a2cd293aaa6ee826350f0c8a0b14cc07c1a6930`; `just -f spec/justfile release-verify learner-model v0.1.5 learner-model-proto-v0.1.5.tar.gz` — "learner-model v0.1.5, signed by the shared release.yml of .github v0.1.21"; in the code of v0.1.5 — `CheckKindTime`, the field of evidence `LeadingQuestionOnRequest` — "whether a hint the learner asked for gave the task's leading question, which a hint gives only before the first wrong answer", and `EvidenceKind = 2`; [format text](../spec/format/README.md), section 19, — "a hint that the learner asks for gives … the task's `leading_question` when it has one"; the row В38 of the RUN.md of [altnovatrail/modules](https://github.com/altnovatrail/modules) — «`modules` v0.2.1 с `olympiad-math` 0.3.0 и `grade-3-math` 0.3.0 — с `time` у 30 задач и наводящими вопросами у 311 задач и 8 шаблонов — вышел в 06:23 того же дня, на `spec` v0.12.0».

### 3.7. Release v0.1.6, T44 and T45

Read 2026-11-11 for T44: the release [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model) `v0.1.6` of 2026-11-11, its tag on the commit `7d796e290599eeb15481f0098c9356bf6c31afe9` (`git ls-remote`), — "Make Readout.Task faster with the same bits (T40)"; the row В39 of its RUN.md — «`Readout.Task` ускорен релизом `learner-model` `v0.1.6`, коммит 7d796e2, … с тем же API и теми же P до бита … все 69 чисел гейта в полосах и до бита те же»; the rows В42 and В43 — a single pin of `sim` for the second round of the model, and the model's round before the engine's round, both with the recommendation (a).

Read 2026-11-11 for T45 (R43): the release [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model) `v0.1.6` — the archive `learner-model-proto-v0.1.6.tar.gz` (`gh release download`); `just -f spec/justfile release-verify learner-model v0.1.6 learner-model-proto-v0.1.6.tar.gz` accepted it — "learner-model v0.1.6, signed by the shared release.yml of .github v0.1.21". The module's `go.mod` at the tag holds only `go 1.27.3`, with no requirements.

### 3.8. Releases v0.1.7 and v0.1.8, T48

Read 2026-11-13 for T48 (R46): the releases [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model) (`gh release view`, `git ls-remote`).
- `v0.1.7` of 05:48 UTC, its tag on the commit `f85b84ee1d0a03f46e2b77ebff9778741fd1a2d9` — "Merge pull request #64 from altnovatrail/spec-v0.16.3 / Move to spec v0.16.3, with Go 1.27.4 (T47, R33)". Its pull requests touch no Go file: against v0.1.6 the release differs only in `.claude/settings.json`, `.devcontainer/`, `.github/workflows/ci.yml`, CLAUDE.md, README.md, RUN.md, `docs/` and `spec`, so its code is that of v0.1.6.
- `v0.1.8` of 18:08 UTC, immutable, its tag on the commit `ab688b243f4ae511eec0a5e18e166f64b8148240` — "Merge pull request #62 from altnovatrail/t44-theta"; the notes — "theta: let a module's level forget after a pause of 14 days (R04, T44, T45)". The merge's second parent, 1a4c478, — "Take theta to the gate of sim v0.1.13 and record its judgement (T45) … the code is that of a063835. The merge releases v0.1.8." The commit a063835 is an ancestor of the tag, and between them only `.claude/settings.json`, `.devcontainer/`, `.github/workflows/ci.yml`, CLAUDE.md, README.md, RUN.md, `docs/` and `spec` changed — no Go file and no `go.mod`; against v0.1.7 the release brings the code of `theta`, `model.go` and `step.go` with their tests; the module's `go.mod` at the tag holds only `go 1.27.3`, with no requirements.
- `just -f spec/justfile release-verify learner-model v0.1.8 learner-model-proto-v0.1.8.tar.gz` accepted it — "learner-model v0.1.8, signed by the shared release.yml of .github v0.1.21".

### 3.9. Release v0.1.9 and commit d2d7ccd, T49

Read 2026-11-14 for T49: the release [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model) `v0.1.9` of 2026-11-13, 19:26 UTC, immutable, its tag on the commit `fc104d1075005a471bb0243a1e1c8713422c80f6` (`gh release view`, `git ls-remote`). Its commits against v0.1.8 — "Record the release v0.1.8 of theta (T45)", "Sum up E5: section 3 checked against the code, E7 proposed (T46)", "Move the finished epic E5 to its archive, docs/plan/E5.md" and "Move the shared workflows to .github v0.1.24, which runs the tests in wasip1 (T48)"; the files they change (`gh api …/compare/v0.1.8...v0.1.9`) — `.github/workflows/` `ci.yml`, `codeql.yml`, `release.yml` and `scorecard.yml`, CLAUDE.md, README.md, RUN.md, `docs/architecture.md`, `docs/decisions.md` and `docs/plan/E5.md`: no Go file and no `go.mod`, so its code is that of v0.1.8.

The branch `main`, commit d2d7ccd of 2026-11-14, 00:07 UTC — "Merge pull request #73 from altnovatrail/t51-gate". `.github/workflows/ci.yml` — `SIM_TAG: v0.1.14`, `SIM_COMMIT: 9c8f782189a4d1d8c2fef59d496f324c03e3cb82`. RUN.md, table 1 — «Текущий эпик — E7. Третий круг кандидатов модели», T51 «Гейт на теге `sim` с базой на `v0.1.8`» done, T52 «Разведка третьего круга» not, nor «| [ ] | T53 | Кандидаты — черновые pull request | Claude, проверяю я | да | T52 |» and «| [ ] | T54 | Суд `model-3` и выпуск выбранного | я + Claude | — | T53; суд `model-3` RUN.md `sim` (В50) |»; table 3, В50 — «„В реестр spec“ — суд `model-3` и порядок кругов на `sim`. Задачи под суд `model-3` у `sim` нет», «тег с новым шагом посреди круга модели сдвинул бы все прогоны, и разведку пришлось бы повторить», with the recommendation «а: круг модели готов первым — разведке не нужно ничего, кроме тега, — а кандидатов `engine-2` ещё нет, и шага закрепления в уроке в `sim` тоже». No open pull requests.

Read again 2026-11-14, 15:30 UTC, for T60 (`git ls-remote`, `gh release list`, `gh pr list`, RUN.md at `main`): `main` is still d2d7ccd, and `v0.1.9` is still the latest release. RUN.md, table 1, the same — «| [ ] | T52 | Разведка третьего круга | я + Claude | да | T51 |», and T53 and T54 not done. No pull request is open; the branch `e7-deploy` of the closed pull request #71 is left, 3 commits ahead of `main` and 9 behind.

Read again 2026-11-14, 16:14 UTC, for the laying out of E9 (R54; `gh api` of the commit of `main` and of RUN.md at `main`, `gh pr list`): `main` is still d2d7ccd, and RUN.md, table 1, still has «| [ ] | T52 | Разведка третьего круга | я + Claude | да | T51 |»; no pull request is open.

Read again 2026-11-15, 00:30 UTC, for T65 (`gh release list`, `git ls-remote`, `gh pr list`, RUN.md and `.github/workflows/ci.yml` at `main` by `gh api`): `main` is still d2d7ccd, `v0.1.9` is still the latest release, and no pull request is open; the branch `e7-deploy` is still there. RUN.md, table 1, still has «| [ ] | T52 | Разведка третьего круга | я + Claude | да | T51 |», and T53–T55 not done; `ci.yml` — still `SIM_TAG: v0.1.14`.

### 3.10. Commit e3aaee6: epic E8, the round of fairness, T69

Read 2026-11-16 for the laying out of E10 (R64; `gh api` of the commit of `main`, `gh release list`, `gh pr list`, `gh api …/compare/v0.1.8...v0.1.10`, and RUN.md and `docs/decisions.md` in a clone in a temporary directory): the branch `main` of [altnovatrail/learner-model](https://github.com/altnovatrail/learner-model), commit `e3aaee69a7bd6a168f0d664b3888f1e8e89ce752` of 2026-11-16, 02:27 UTC — "Merge pull request #77 from altnovatrail/t58-fairness-round". The latest release is `v0.1.10` of 2026-11-16, 01:17 UTC; against `v0.1.8` it changes no Go file and no `go.mod`. No pull request is open.
- R40 learner-model: «Epic E8 is the round of fairness by the grade: why the model leaves a start above the truth more slowly than it rises from a start below it, then candidates that bring `M11.same` of `D-grade5` and `D-grade1` up to 0.90 under the model's criterion as is, judged by `sim`'s E10 on the candidate's commit, and the chosen one released or the exit recorded, and the outcome «в реестр spec» for dpia 10.2 spec and dpia 8 spec — T59–T64.»
- RUN.md, table 1 — «Current epic — E8. Круг справедливости по классу», no task done: T59 «Разбор: почему со старта выше правды оценка сходит дольше»; «| [ ] | T60 | Гейт на теге `sim` с критерием суда круга | Claude | — | критерий суда E10 записан в `sim` (строка 4.18 COORDINATION.md spec) |»; T61 «Разведка круга»; «| [ ] | T62 | Кандидаты — черновые pull request | Claude, проверяю я | да | T61 |»; «| [ ] | T63 | Суд и выпуск выбранного | я + Claude | — | T62; суд E10 `sim` по коммиту кандидата (строка 4.18 COORDINATION.md spec); кандидату с иным стартом по классу — ответ `spec` на его строку «в реестр spec» |»; T64 «Итог круга и архив эпика».
- «Почему такой порядок»: «Гейт — до разведки, как T51 в E7: разведка и суд идут на одном теге … Тег — тот, на котором `sim` записал критерий E10, до первого кандидата (R103 spec; R49 sim)»; «Суд ждёт не релиза `sim`, а его суда по коммиту кандидата: `sim` гонит черновик из его коммита (строка 4.18). Тег `sim` с полосами выбранного нужен только мержу выбранного (R103 spec).»
- «Выяснилось при записи»: «σ₀ 3,0 уже пробовала разведка T30, на 300 учениках и без `M11`: по общему счёту не лучше базы и нарушает 28 ограничений, но группа `off-band` — ученики, чей класс не их уровень, — лучше на +0,062 сверх шума»; «Старт в полосе класса — правило R51 spec, … и иной старт меняет `spec`». «Не входит»: «код `sim`: критерий и суд — E10 `sim`».
- T61, «Сделать»: «разбор — `sim judge` критерием E10 и `sim compare` против базы»; T63, «Сделать»: «прочитать отчёт суда и подтверждение на отложенном наборе E10».

## 4. `engine`

### 4.1. Commit 08bc48a, T02 and T03

The repository [altnovatrail/engine](https://github.com/altnovatrail/engine), branch `main`, commit 08bc48a, read 2026-11-07: in table 1 of the RUN.md of `engine` only the move to `spec` v0.2.9, T07, is done, and the tasks of detailing T01–T06 have not begun. So the engine's criterion (architecture 3.6) rests on none of its decisions, and the engine's base is its first version, by В24 of `learner-model`.

Read 2026-11-07 for architecture 3.9, in the planning and the execution of T03.

The repository [altnovatrail/engine](https://github.com/altnovatrail/engine), commit 08bc48a — the same that T02 read: there is no code and no `go.mod`, and the tasks of detailing have not begun.

### 4.2. Commit b962ff4

Read 2026-11-08 for the move to `spec` v0.9.8 (R07).

The repository [altnovatrail/engine](https://github.com/altnovatrail/engine), branch `main`, commit b962ff4: the detailing is decided, R03–R10 engine, and with it architecture 3.5, 3.6 and 3.9 rest on its decisions — the record of the commit 08bc48a, sources 4.1, speaks of the time of T02. R03 engine, the 2026-11-07 amendments — «задача с виденным решением — не кандидат 24 часа, кроме второго прохода: если выбор не дал ничего — ни урока, ни задачи, ни в одном модуле, — правило проходит ещё раз», and «пререквизит, которому во входе учить нечем, навык после себя не держит»; R07 engine — «Движок обходится без случайности: генератора и зерна выбора у него нет, и поле „зерно“ в записи выбора пусто»; R10 engine — «Движок даёт тот же выбор и тот же код на любой сборке одной своей версии и версии `learner-model` — компилятором Go, на любой архитектуре, в WebAssembly тоже».

### 4.3. No releases, 2026-11-08

Read 2026-11-08 for the move to `spec` v0.9.9 (R08).

There are no releases of `engine` (`gh release list -R altnovatrail/engine`, 2026-11-08); the row 4.7 of [COORDINATION.md](../spec/COORDINATION.md) at v0.9.9 — «E2: T13–T22; сделаны T13–T19».

### 4.4. Release v0.1.0

Read 2026-11-08 for the planning of the epic E2 (R09).

The repository [altnovatrail/engine](https://github.com/altnovatrail/engine), the release `v0.1.0` — immutable, published 2026-11-08, its tag on the commit c6c5bb9, the files `engine-v0.1.0.tar.gz` and its bundle (`gh release view v0.1.0 -R altnovatrail/engine`). `gh release verify v0.1.0 -R altnovatrail/engine` — "Release v0.1.0 verified!"; `just -f spec/justfile release-verify engine v0.1.0 …` on `spec` v0.9.9 — "engine-v0.1.0.tar.gz: engine v0.1.0, signed by the shared release.yml of .github v0.1.21". The Go proxy serves `github.com/altnovatrail/engine@v0.1.0` with the Origin "Hash: c6c5bb9bc2c6fc486f629da82ee834bd89fc2b35", "Ref: refs/tags/v0.1.0". The `go.mod` of the tag — `module github.com/altnovatrail/engine`, `go 1.27.3`, `require github.com/altnovatrail/learner-model v0.1.0`. The API — `func Choose(in Input) (Choice, error)`; T21 of the RUN.md of `engine`: «перевод `server.json` у теста свой: расхождение с переводом `api` и `sim` векторы не ловят — его поймает стенд `sim`». The record sources 4.3, "There are no releases of `engine`", speaks of the time of T09.

Read 2026-11-08 for the packages of T12 (R14).

The repository [altnovatrail/engine](https://github.com/altnovatrail/engine), the tag `v0.1.0`: `release_internal_test.go` — the translation of `server.json` into the engine's types and the filter `forAge`, with the comment "The translation of server.json into the engine's types is this test's, not api's or sim's: the golden vectors do not catch a difference between them" and "archive/tar would pull os/user, whose cgo the rules test cannot type-check for every target".

### 4.5. Commit 9751ba4, T10

Read 2026-11-08 for the frame of T10 (R10, R11).

The repository [altnovatrail/engine](https://github.com/altnovatrail/engine), branch `main`, commit 9751ba4: `justfile` — the recipes `ci-lint`, `ci-test`, `github`, `script`, `test-targets`, `fma-diff`, `rules`, `same-bits`, `ci-vuln` (`golang.org/x/vuln/cmd/govulncheck@v1.9.0`) and `ci-licenses` (`github.com/google/go-licenses/v2@v2.0.1`); `.github/workflows/ci.yml` — the job `go` with `same-bits: true` on `go.yml@706f4067a66c91b0442033f82a7924eeb8b60497 # v0.1.21`; `.github/codecov.yml`, `.github/coderules.json`, `dependabot.yml` with `gomod`; `.devcontainer/gopls/` — the lock file of gopls v0.24.0; R12 engine — «`ci-lint` — `gofmt`, `go vet` и сверка версии gopls …, без других линтеров»; R13 engine — gopls from the lock file.

### 4.6. Releases v0.1.7 and v0.1.8 and commit e0fdb3f, T49

Read 2026-11-14 for T49: the repository [altnovatrail/engine](https://github.com/altnovatrail/engine) (`gh release list`, `git ls-remote`, `gh api …/compare`).
- `v0.1.7` of 2026-11-11, 22:59 UTC, its tag on the commit `68218a9675f83e349a2a567d420fbfee8d402b33`, and `v0.1.8` of 2026-11-12, 00:36 UTC, on `d8a35a4baed9639b7cf60e95c09f22ec70371e0f`. Against v0.1.6, the commit 57e0ab8 that `gate.json` names, v0.1.7 changes `.claude/settings.json`, `.devcontainer/`, `.github/workflows/ci.yml`, CLAUDE.md, README.md, RUN.md, `docs/` and `spec`, and v0.1.8 then RUN.md, `docs/`, `go.mod` and `go.sum` alone: no Go file, so their code is that of v0.1.6.
- The branch `main`, commit e0fdb3f of 2026-11-13, 17:52 UTC — "Merge pull request #67 from altnovatrail/e7-consolidation-in-lesson", the pull request "Record the planning of T52: it waits for sim's tag with the step". Against v0.1.8 it changes, of the Go files, only `bench_internal_test.go` and `release_internal_test.go`. Its RUN.md, table 1 — «| [ ] | T52 | Вариант правила | Claude, проверяю я | да | T50; тег `sim` с шагом и критерием T50 |»; T52, «Решено автором (2026-11-13)» — «T52 ждёт тега `sim` строки 4.15 целиком: кода до тега нет», and «Выяснилось при планировании» — «„Gate“ `engine` стоит на `sim` `v0.1.6`».

### 4.7. Branch `e7-consolidation-in-lesson`, commit 0e4680f, R50

Read 2026-11-14 for the laying out of E7 (R50): the repository [altnovatrail/engine](https://github.com/altnovatrail/engine) (`git ls-remote`, `gh api …/compare/main...e7-consolidation-in-lesson`, and RUN.md and `docs/plan/E7.md` at that branch).
- The branch `e7-consolidation-in-lesson` is one commit ahead of `main` e0fdb3f: `0e4680f0631e113bf91f54846a0a4b1322c76219` of 2026-11-14, 03:33 UTC, "Put the variant of consolidation in a lesson into code (T52)", with `rule.go`, `lessons.go`, `skills.go`, their tests and documents. No pull request is open.
- RUN.md, table 1 — «| [x] | T52 | Вариант правила | …» and «| [ ] | T53 | Суд, выбор и релиз | я + Claude | да | T52; суд `sim` |»; «Решено автором (2026-11-12)» — «эпик идёт в одной ветке, `e7-consolidation-in-lesson`, без других веток и pull request; … вариант T52–T53 выходит вторым релизом».
- `docs/plan/E7.md`, T52, «Открыто» — «кандидат суда `sim` — полный SHA этого коммита. Строка файла суда — `{"name": "E7", "group": "corridor", "constants": 0, "new_field": false, "distance": 1, "equal_on": ["novatrail/grade-3-math", "novatrail/english-essays"]}`» and «риск T53: по грубому взгляду выигрыша нет, и вероятен выход»; «Сделано» — the rough look on `sim` `v0.1.15`, 300 working learners a cell: «вариант не держит всех ограничений и не лучше базы. Группа «Коридор» +0,00002 [−0,0019; +0,0020], около нуля; «не хуже» не держит на трёх ячейках олимпиадной …; `neutral` держит».

### 4.8. Commit 0e4680f again, T58

Read 2026-11-14 for T58 (R53): the repository [altnovatrail/engine](https://github.com/altnovatrail/engine) (`git ls-remote`, `gh api …/commits/0e4680f…/check-runs`, and `docs/plan/E7.md` at the branch `e7-consolidation-in-lesson`).
- `refs/heads/e7-consolidation-in-lesson` is still `0e4680f0631e113bf91f54846a0a4b1322c76219`, and `refs/heads/main` still e0fdb3f. The commit has no check runs, and no pull request is open from the branch.
- `docs/plan/E7.md`, T52, «Сделано» — «`just gate` на варианте красный только двумя полосами нулевой ширины — `E2.stuck` на `G2-half` олимпиадной, 0,016, и на `A-linked`, 0,008, — инварианты целы, `ErrRun` нет»; «грубый взгляд рецептом `sim` `v0.1.15`, 300 рабочих учеников на ячейку … «не хуже» не держит на трёх ячейках олимпиадной — `E2.stuck` на `G2-half` и `E4.given` на `G5` и `S-scale3.5`; `neutral` держит»; and of the review's fixes — «`just gate` после них дал те же 69 чисел до бита, а грубый взгляд шёл на коде до ревью».

Read again 2026-11-14, 15:30 UTC, for T60 (`git ls-remote`, `gh release list`, `gh pr list`, RUN.md at `main` and at the branch): `refs/heads/e7-consolidation-in-lesson` is still 0e4680f and `main` still e0fdb3f, `v0.1.8` is still the latest release, and no pull request is open. RUN.md of the branch, table 1 — «| [ ] | T53 | Суд, выбор и релиз | я + Claude | да | T52; суд `sim` |»; table 2 — «| E5 | Второй круг коридора | Второй круг суда движка на нетронутом `held-out-engine-1` …», not laid out.

### 4.9. Epic E5 of `engine`: the candidates D, E and F, T65

Read 2026-11-15, 00:35 UTC, for T65: the repository [altnovatrail/engine](https://github.com/altnovatrail/engine) (`git ls-remote`, `gh release list`, `gh pr list`, `gh api` of `compare/…` and of RUN.md, `docs/measurements.md` and `.github/workflows/ci.yml` at `main`).
- `main` is the commit 77ee605 of 2026-11-14, 22:23 UTC — "Merge pull request #71 from altnovatrail/e5-corridor-round-2" —, and it is the only branch; `v0.1.8` is still the latest release, and no pull request is open. Against v0.1.6, the commit 57e0ab8 that `gate.json` names, `main` changes, of the Go files, only `bench_internal_test.go` and `release_internal_test.go`. `ci.yml` — `SIM_TAG: v0.1.17`, `SIM_COMMIT: 6f5151f5eb241b5b44b11b71296a0e3cd2c201a1`.
- RUN.md, table 1 — «Текущий эпик — E5. Второй круг коридора»: «| [x] | T57 | Разбор коридора на новой базе | …», «| [x] | T58 | Варианты кандидата коридора | …», «| [ ] | T59 | Суд `engine-2`, выбор и релиз | я + Claude | да | T58; E6 `sim` («Второй круг движка») |»; its order — «Суд — когда до него дойдёт очередь `sim`: E6 `sim` идёт после его E9 и E8 и судит на базе после суда `model-3`».
- T59, «Из T58 (2026-11-14)» — «кандидаты суда — полные SHA коммитов ветки `e5-corridor-round-2`: D 29e205a15f5a05a443851a255f00cffbcbd3d5b2, E e9832fca943b7ee457ddc759396e225a5114033a, F 389798ba6bce0035356b73f42baed39ab16bf833; голова ветки — код базы. Строки файла суда — все `"group": "corridor"`, `"new_field": false`, `"equal_on": ["novatrail/grade-3-math", "novatrail/english-essays"]`: E — `"constants": 0, "distance": 1`, D — `0` и `2`, F — `1` и `3`». `compare/<commit>...main` gives all three as ancestors of `main` — 4, 3 and 2 commits behind it.
- Table 3, В46 — «Грубый взгляд T58: все три кандидата, D, E и F, не держат «не хуже» по олимпиадному `M3.settle` на `G3-drop` — на 0,3–2,8 сверх допуска», with the recommendation «(а): грубый взгляд на 300 учениках ничего не решает»; `docs/measurements.md`, section 6 — the rough look on «A clone of `sim` `v0.1.17` (6f5151f; `gh release verify` accepts it)».

## 5. `.github`

### 5.1. Commit af40f94, T03

Read 2026-11-07 for architecture 3.9, in the planning and the execution of T03.

The repository [altnovatrail/.github](https://github.com/altnovatrail/.github), branch `main`, commit af40f94, the latest release v0.1.15: architecture 3.1 .github — «Входов у обеих обёрток нет, а имена рецептов у них одни»; `go.yml` — `runs-on: ubuntu-24.04` in both jobs; architecture 3.2 .github — the row «`learner-model`, `engine` | гейт `sim` по тегу через `go.work`», the row `modules` — «проверка пакетов CLI `module-kit` и стенд `sim` по тегам», and `sim` has no CI row of its own, but for the test vectors of the generator; R03 .github — «Свои — CI и релиз `spec`, гейт `sim`, eval `modules`…».

### 5.2. Commit 706f406, release v0.1.21

Read 2026-11-08 for the move to `spec` v0.9.9 (R08).

The repository [altnovatrail/.github](https://github.com/altnovatrail/.github), branch `main`, commit 706f406, the latest release v0.1.21: `.github/workflows/go.yml` — the input `same-bits`, `type: boolean`, `default: false`, `Hold the bits the same on every build, by "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules" (R107 in spec).`; architecture 3.1 .github, the row `go.yml`, — what `targets`, `fma-diff`, `arm64` and `rules` do, and «с `same-bits` — `learner-model` и `engine` (R107 spec)»; the section on a repository in Go — «С `v0.1.18` … у обёртки есть вход `same-bits`», «"Go / arm64" — с `v0.1.19`, "Go / rules" — с `v0.1.20`»; `tools/coderules/check.go` — `"goroutine", // a go statement in code`, `"math",      // a function of math whose result is not defined exactly`, `"maprange",  // a range over a map in code`, `"fused",     // a float * or / whose result does not go straight into float64()`, and `mathFuncs` — 14 functions: `Abs`, `Ceil`, `Copysign`, `Float64bits`, `Float64frombits`, `Floor`, `Frexp`, `Inf`, `IsInf`, `IsNaN`, `Ldexp`, `Signbit`, `Sqrt`, `Trunc`; the rules do not touch the files `_test.go` where they say "in code"; `go.yml` — the job `targets` runs `run GOAMD64=v3 go test -count=1 ./...` and four more such lines, `timeout-minutes: 30`, and the job `arm64` runs `toolchain just ci-test`, then the same lines, `timeout-minutes: 45`. The record of the commit af40f94, sources 5.1, «Входов у обеих обёрток нет», speaks of the time of T03.

Read 2026-11-08 for the planning of the epic E2 (R09).

The repository [altnovatrail/.github](https://github.com/altnovatrail/.github), the release v0.1.21 (706f406): `go.yml` builds the stage `toolchain` without the submodule `spec` and runs `just ci-lint|ci-test|ci-vuln|ci-licenses`; `release.yml`, the step that builds the files, — "release-files wrote nothing into dist/, or dist/ is a link." — a release without a file in `dist/` fails.

Read 2026-11-08 for the frame of T10 (R10, R11).

The repository [altnovatrail/.github](https://github.com/altnovatrail/.github), the release v0.1.21 (706f406): `tools/coderules/lists.go` — the fields `code-imports`, `test-imports`, `test-math`, `requires`, `nested-modules`, "A field it does not know is a mistake"; `tools/coderules/check.go` — the rules `import`, `cgo`, `asm`, `linkname`, `buildtag`, `require`, `replace`, `nested`, `math`, `fused`, `assignop`, `complex`, `typeparam`, `maprange`, `unicode`, `goroutine`; `go.yml`, the job `targets` — `run GOOS=js GOARCH=wasm go test -count=1 -exec="$(go env GOROOT)/lib/wasm/go_js_wasm_exec" ./...`, and `ALLOWED_LICENSES` — "0BSD,Apache-2.0,BlueOak-1.0.0,BSD-2-Clause,BSD-3-Clause,CC0-1.0,ISC,MIT,MIT-0,Unlicense".

### 5.3. Release v0.1.24, T68

Read 2026-11-16 for T68 (R63), from a clone of [altnovatrail/.github](https://github.com/altnovatrail/.github) in a temporary directory.

`gh release list -R altnovatrail/.github` — `v0.1.24`, «Latest», of 2026-11-13, 05:34 UTC; then `v0.1.23` of 2026-11-12 and `v0.1.22` of 2026-11-10. `gh release verify v0.1.24 -R altnovatrail/.github` — "Resolved tag v0.1.24 to sha1:fab102d8bbb4659575e107ee488e072aaaa71821", "Release v0.1.24 verified!". The notes: v0.1.22 — "E17: cache of the toolchain image"; v0.1.23 — the archive of E17, the move to `spec` v0.16.3, the sources and measurements in their own files; v0.1.24 — the translation of its documents and "E19, T58: wasip1 tests in node:wasi in go.yml".

`git diff v0.1.21 v0.1.24 -- .github/workflows/`, of the workflows `sim` calls:
- `go.yml` — no input changes; each of the jobs `check`, `targets`, `fma-diff`, `arm64` and `rules` takes the stage's image from the cache: «Each job that runs in the toolchain stage takes its image from the caller's cache of GitHub Actions and builds it from the Dockerfile only on a miss, so that an outage of what the Dockerfile downloads from does not fail the checks while the Dockerfile stays the same (R36)»; the steps «Toolchain image from the cache» — `actions/cache/restore@55cc8345863c7cc4c66a329aec7e433d2d1c52a9 # v6.1.0`, `path: ${{ runner.temp }}/toolchain.tar`, `key: toolchain-${{ runner.os }}-${{ runner.arch }}-${{ hashFiles('.devcontainer/Dockerfile') }}` —, «Toolchain image», with `docker load --quiet --input "$tar" && docker image inspect toolchain`, the build on a miss and `docker save --output "$tar" toolchain`, and «Toolchain image into the cache» — `actions/cache/save`, `if: steps.toolchain.outputs.cache-hit != 'true'`, `key: ${{ steps.toolchain.outputs.cache-primary-key }}`; «The caches are the caller's: one a pull request saves is "restored by re-runs of the pull request" alone, and the one of main, which every pull request reads, its CI on a push to main saves». In `targets` and `arm64`, `run GOOS=wasip1 GOARCH=wasm go build ./...` became `run GOOS=wasip1 GOARCH=wasm go test -count=1 -exec="node --stack-size=8192 --disable-warning=ExperimentalWarning $wasi" ./...`, by an executor that runs the test binary in `node:wasi` — «wasip1 in node:wasi, by an executor written outside the checkout»;
- `documents.yml` — the image from the cache by the same key, in `${{ runner.temp }}/novatrail-toolchain.tar`, through `TOOLCHAIN_IMAGE_FILE` of `spec/tools/toolchain/toolchain.sh`, which «spec ships … from v0.13.0»;
- `release.yml` and `scorecard.yml` — comments alone; `secrets.yml`, `dco.yml` and `codeql.yml` — unchanged; `tools/coderules/` — unchanged.

`docs/architecture.md` of v0.1.24, the template of `.github/dependabot.yml`: «# Every reference to a shared workflow of .github in one pull request, to one release: / # Dependabot names each altnovatrail/.github/.github/workflows/<file>.yml. / shared-workflows: / patterns: ["altnovatrail/.github/*"]».

Go 1.27.4 of the container, `src/syscall/fs_wasip1.go`: `func Fchmod(fd int, mode uint32) error { var stat Stat_t; return Fstat(fd, &stat) }` — on wasip1 a change of a file's mode only checks that the file is there. `GOOS=wasip1 GOARCH=wasm go test` by the executor of `go.yml` v0.1.24, in Node v24.21.0 of the container, on `sim` at `main` c5e9a4c: every package `ok` but the root, `TestMerge` — "a new file of results has the mode -rw-------".

### 5.4. Release v0.1.30, T75

Read 2026-11-16 for T75 (R68), from a clone of [altnovatrail/.github](https://github.com/altnovatrail/.github) in a temporary directory.

`gh release list -R altnovatrail/.github` — `v0.1.30`, «Latest», of 2026-11-16, 11:09 UTC; then `v0.1.29`–`v0.1.25` the same day, from 06:21 UTC. `gh release verify v0.1.30 -R altnovatrail/.github` — "Resolved tag v0.1.30 to sha1:8311cbfc01082e9b196c0af8d7b49192115cc84c", "Release v0.1.30 verified!". Its notes: "CodeQL's Go from go.dev's archive the caller's Dockerfile pins, without the image (T65)". By `git diff` of each tag against the one before, of the workflows: `v0.1.25` — `go.yml`; `v0.1.26` — `go.yml`, `typescript.yml`, `documents.yml`, `self-ci.yml`; `v0.1.27` — `go.yml`; `v0.1.28` — `release.yml`; `v0.1.29` — `typescript.yml`; `v0.1.30` — `codeql.yml`.

`git diff v0.1.24 v0.1.30 -- .github/workflows/`, of the workflows `sim` calls:
- the inputs of `workflow_call` of `go.yml`, `release.yml`, `codeql.yml`, `documents.yml`, `secrets.yml`, `dco.yml` and `scorecard.yml` are the same, but for comments; `secrets.yml`, `dco.yml`, `scorecard.yml` and `tools/coderules/` do not change at all;
- `go.yml`, R41 `.github`: «The jobs of the shared `go.yml` that run Go — `check`, `targets`, `fma-diff`, `arm64` and `rules` — take the caches of Go from the caller's cache of GitHub Actions: one folder `${{ runner.temp }}/go-cache` … The key is `go-<OS>-<architecture>-<check>-<hash of .devcontainer/Dockerfile>-<hash of **/go.sum>`, with no restore key. … Only a push to `main` saves»;
- `go.yml` and `documents.yml`, R42 `.github`: «keeps it in one file, `${{ runner.temp }}/novatrail-toolchain.tar`, of the image `novatrail-toolchain`, the file and the name that `spec/tools/toolchain/toolchain.sh` gives, by the key of R31 spec. Only a push to `main` saves it. In `go.yml` and `typescript.yml` a job `toolchain` looks the entry up and, on a miss, alone builds the image and hands its file to the jobs after it by an artifact of the run … `arm64` builds its own»; «`documents.yml` cannot take the artifact of `go.yml`: two called workflows of one run of the caller see each other's jobs only through the caller's `needs`». The steps in `go.yml` v0.1.30: «Toolchain image from the cache» — `path: ${{ runner.temp }}/novatrail-toolchain.tar`, `key: toolchain-${{ runner.os }}-${{ runner.arch }}-${{ hashFiles('.devcontainer/Dockerfile') }}`; «Toolchain image» — one script with `SAVE`, `docker build --target toolchain --tag novatrail-toolchain - < .devcontainer/Dockerfile`, `docker save` only `if [ "$SAVE" = true ] && [ "$restored" = false ]`; in `arm64` `SAVE: ${{ github.event_name == 'push' && github.ref == 'refs/heads/main' }}` and «Toolchain image into the cache» with `if: steps.toolchain.outputs.cache-hit != 'true' && github.event_name == 'push' && github.ref == 'refs/heads/main'`;
- `go.yml`, R43 `.github`: «With `same-bits`, the targets of the shared `go.yml` run as a matrix of the job `target`, one job a target — `amd64-v3`, `cpu-off`, `386`, `js-wasm` and `wasip1`, the jobs "Go / target (<target>)" —, and the job `targets` after them is the one check "Go / targets" of the caller's ruleset … "Go / fma-diff" stays one job, and its six targets run at once … The script of the targets is one for the matrix and `arm64` and runs the targets `TARGETS` names; `arm64` names `cpu-off js-wasm wasip1`. The input and the checks of the ruleset keep their names». The matrix in `go.yml`: `target: [amd64-v3, cpu-off, "386", js-wasm, wasip1]`; in `arm64`: `TARGETS: cpu-off js-wasm wasip1`;
- `release.yml`, R44 `.github`: «signs and checks its files eight at once … The inputs, the identity of the signature and what a consumer checks (R77 spec) do not change»;
- `codeql.yml`, R45 `.github`: «The analysis of `go` … reads `ARG GO_VERSION` and `ARG GO_SHA256_AMD64` from `.devcontainer/Dockerfile` — from the stage `toolchain` and the stages it is built from … No step of `codeql.yml` runs `docker`; the input, the checks and the categories keep their names»; «All nine callers that name `go` — … `sim` … — have them so on `main`». `.devcontainer/Dockerfile` of `sim`: `ARG GO_VERSION=1.27.4` and `ARG GO_SHA256_AMD64=…` in the stage `toolchain`, whose `ENV` sets `PATH` and `GOTOOLCHAIN=local` alone.

RUN.md of `.github` at `v0.1.30`, table 1 of E21: «| [ ] | T66 | Стадия `toolchain-core` в релизе и проверках | Claude | да | T60–T65; строка 7.14 COORDINATION.md spec |»; its «Цель»: «`release.yml` строит стадию `toolchain-core` Dockerfile вызывающего, без кеша, а проверки без браузера и Go-тестов — lint, vuln, licenses и документы — идут в ней; у Dockerfile без этой стадии — как прежде.»

## 6. `modules`

### 6.1. Commit 1a28f4a, T03

Read 2026-11-07 for architecture 3.9, in the planning and the execution of T03.

The repository [altnovatrail/modules](https://github.com/altnovatrail/modules), branch `main`, commit 1a28f4a, no releases: R03 modules — «Стенд идёт на пакете, который собрал CI: `sim` читает часть для сервера по типам канонического вида, без `module-kit`», «Pull request, который меняет тег `sim`, гоняет стенд на графах всех модулей», «переход формата идёт без `sim`»; architecture 3.2 modules — until the answer of T03, the gate is the graph; the recipe `modules-check` builds the packages into a temporary folder (`mktemp -d`, `module-kit build`), and `ci.yml` does not upload them as an artifact.

### 6.2. Release v0.1.0

Read 2026-11-08 for the planning of the epic E2 (R09).

The release [altnovatrail/modules](https://github.com/altnovatrail/modules) `v0.1.0`, the packages by the release's `module-index.json`: `package_digest` is the sha256 of the package's file `inventory.json`, not of the archive; the prerequisites of the three modules are only their own skills, so there is no skill of another module that would rest on an anchor skill, though grade 3 maths teaches three (`teaches`: `core.division-with-remainder`, `core.fractions-of-a-whole`, `core.multiplication-table`); no package has answer slots `predicate`; all 27 answer slots of olympiad maths are `number`, and the school maths and the essays have a choice of three options, with `options` of the first inside `check` and of the second beside it.

Read 2026-11-08 for the packages of T12 (R14).

The release [altnovatrail/modules](https://github.com/altnovatrail/modules) `v0.1.0`: `just -f spec/justfile release-verify modules v0.1.0 …` for `module-index.json` and the three archives — "modules v0.1.0, signed by the shared release.yml of .github v0.1.15"; the SHA-256 of `module-index.json` — `ed4d8bb325fdcb3ea54da439e08c58d927f773ba0f097338f7767ad76be8c90e`. The packages have no templates: no task has `template`, and `instances` of the canonical view are empty.

### 6.3. Releases v0.1.4 and v0.1.1, T30

Read 2026-11-09 for T30 (R31).

The releases [altnovatrail/modules](https://github.com/altnovatrail/modules) `v0.1.4` and `v0.1.1`, read 2026-11-09: `just -f spec/justfile release-verify modules v0.1.4 module-index.json` — "module-index.json: modules v0.1.4, signed by the shared release.yml of .github v0.1.15", the SHA-256 of the index — `a244c0ba2832ba6360f45bc74d293572b84b491fc823c7b01bdf6004ff9edb2a`; `release-verify modules v0.1.1` for the archive `sha256-123cb0d3….tar.gz` — "modules v0.1.1, signed by the shared release.yml of .github v0.1.15". The index `v0.1.4` names `olympiad-math` 0.1.0 and 0.1.1, and `grade-3-math` and `english-essays` only 0.1.0. `olympiad-math` 0.1.1: 19 skills and tiers instead of 3, 772 tasks instead of 27 — all on the sheet `handwriting_sheet`, 19 of them with the instances of a template —, 41 lessons instead of 9; the answer slots — 569 `number`, 7 `fraction` and 196 choices of five options in `check`, no `predicate`; `fractions-parts` rests on `core.fractions-of-a-whole`, and `number-divisibility` on `core.division-with-remainder`.

### 6.4. Releases v0.2.1–v0.2.5, T43

Read 2026-11-11 for T43 (R42): the releases [altnovatrail/modules](https://github.com/altnovatrail/modules) `v0.2.1`–`v0.2.5`. The index of `v0.2.5` is byte for byte that of `v0.2.1`, SHA-256 `cebf23cf94456ef313ca00c97ca05db7fccb2575780f8dfd899623ea60aaec6a`; the notes of `v0.2.2`–`v0.2.5` — "T56: release v0.2.1 is out and verified", "T58: recipe of synthetic samples on canvas's generator", "T60: first synthetic samples of grade-3-math", "Move to spec v0.15.0 (R30)". `just -f spec/justfile release-verify modules v0.2.5 module-index.json` — "module-index.json: modules v0.2.5, signed by the shared release.yml of .github v0.1.15"; `release-verify modules v0.2.1` for `sha256-e4de172e….tar.gz` and `sha256-dd2c6eb2….tar.gz` — "modules v0.2.1, signed by the shared release.yml of .github v0.1.15". `olympiad-math` 0.3.0 has 775 tasks, 22 templates, a leading question in 322 tasks, and the checks `number` 572, `choice` 166, `fraction` 7, `time` 30; `grade-3-math` 0.3.0 has the same 12 tasks and a segment `worked_example` in each of the 4 lessons.

### 6.5. Commit 8e36fbe: Д5, T49

Read 2026-11-14 for T49: the repository [altnovatrail/modules](https://github.com/altnovatrail/modules), branch `main`, commit 8e36fbe of 02:22 UTC — "Merge pull request #102 from altnovatrail/t68-spec-v0.17.2". Its table 4 is in DEBT.md, and the row Д5 is open there — «Гейта `sim` в CI нет, и релизы выходят без него — от первого, `v0.1.0`, до `v0.2.5`», «Перевод `time` вышел в `sim` v0.1.9 …, новейший — v0.1.10 (2026-11-12): job ждёт только своей задачи».

Read again 2026-11-14 for T60: `main` is the commit 6e4391f of 04:12 UTC — "Merge pull request #103 from altnovatrail/t69-check-board" —, and the row Д5 of its DEBT.md is still open — «Гейта `sim` в CI нет, и релизы выходят без него — от первого, `v0.1.0`, до `v0.2.5`».

Read again 2026-11-15, 00:30 UTC, for T65 (`gh release list`, `gh api` of the commit of `main`, of `compare/v0.2.7...main` and of DEBT.md at `main`): `main` is the commit 7e99049 of 2026-11-14, 19:58 UTC, and `v0.2.7` is still the latest release; against it `main` changes README.md, RUN.md, `docs/architecture.md`, `docs/measurements.md` and `docs/plan/E9.md` alone, no module. The row Д5 of its DEBT.md is still open — «Гейта `sim` в CI нет, и релизы выходят без него — от первого, `v0.1.0`, до `v0.2.5`».

### 6.6. Releases v0.2.6 and v0.2.7, T60

Read 2026-11-14 for T60: the releases [altnovatrail/modules](https://github.com/altnovatrail/modules) `v0.2.6` of 02:30 UTC and `v0.2.7` of 04:21 UTC (`gh release view`, `gh release download … -p module-index.json`).
- The index of `v0.2.6` is byte for byte that of `v0.2.5`, SHA-256 `cebf23cf94456ef313ca00c97ca05db7fccb2575780f8dfd899623ea60aaec6a`; its notes — "T62: the samples page canvas v0.3.11 opens the eval set", "T62: two of the author's samples of grade-3-math", "T63: epic E8 done and archived", "Д5 and Д8: the gate of sim waits; the corridor of grade-3-math measured", "T64: move to spec v0.16.3 and rebuild the documents (R31)".
- The index of `v0.2.7`, SHA-256 `6c46a1a902d96580c13ea75a37759ad7a2d043896c79d68116edc3851e02083a`, differs from it by one module alone — `"novatrail/check-board": {"0.1.0": {"module_format_version": "0.4", "package_digest": "sha256:28893a0f82709a2526d13950826ef8cb442608ad85d9c824f53a235122846f56", "release": "v0.2.7"}}`; the entries of `novatrail/english-essays` 0.2.0, `novatrail/grade-3-math` 0.3.0 and `novatrail/olympiad-math` 0.3.0 are the same. Its notes — "T69: the board module novatrail/check-board 0.1.0 (R33)".
- R33 modules — "The module of the check board of `tutor` is `novatrail/check-board` 0.1.0 …: the eleven tasks of the board in its order, from the banks of `tutor` at `v0.1.35` …; three skills by the kinds of the banks and a fourth that is not assessed, 22 traps, one lesson of all eleven tasks and three synthetic samples. The explanation language is Russian, the audience `adult` alone, one band without grades".

## 7. `spec` and `module-kit`

### 7.1. The Go module `spec` v0.9.9: the generator

Read 2026-11-08 for the planning of the epic E2 (R09).

The Go module `spec` at the tag v0.9.9: there is no implementation of the generator in Go — the test vectors `format/vectors/random/*.yaml` (12 cases) and `format/vectors/templates/*.yaml` by the types `formatv0.RandomVector` and `formatv0.TemplateVector`, embedded in `spec.Format`; the second implementation is `tools/format-check/random.mjs`. The implementation of `module-kit` is its `internal/random`, which no other module can import.

### 7.2. The move to `spec` v0.12.0

Read 2026-11-10 for the move to `spec` v0.12.0 (R32): `gh release verify v0.12.0 -R altnovatrail/spec` — "Release v0.12.0 verified!", the commit 5b4640a; the format text, section 2, — "Version 0.4 comes out with v0.10.0" and "While the format is v0 (0.y) … the platform reads only the current version"; the release [altnovatrail/modules](https://github.com/altnovatrail/modules) `v0.2.0` — "Move to spec v0.12.0 and module-kit v0.3.0, with the modules in format 0.4 (R28)": the index of format 0.4 names all three modules at version 0.2.0 in format 0.4, and the earlier versions in 0.3; `release-verify modules v0.2.0` for `module-index.json` and the three archives — "modules v0.2.0, signed by the shared release.yml of .github v0.1.15", the SHA-256 of the index — `a45c1e3d9bdfae094d08c1f0e62f51f0e9432b6c5afbf4193ce98c386c89f71f`; `release-verify engine v0.1.4` and `learner-model v0.1.4` — "signed by the shared release.yml of .github v0.1.21", the commits of the tags — ddfed3a and acfc4f2; Д137 of table 4 of the RUN.md of `spec` — «`engine` `v0.1.3` (fed6cd2) и `learner-model` `v0.1.3` (990d6b1) вышли через гейт с теми же числами».

### 7.3. The move to `spec` v0.15.0

Read 2026-11-11 for the move to `spec` v0.15.0 (R39): `gh release verify v0.15.0 -R altnovatrail/spec` — "Release v0.15.0 verified!", its tag on the commit 58d3f46; `RULES.md`, `format/` but for the example of traps, `types/formatv0`, `repo-template/` and `.devcontainer/` are the same as at v0.12.0 (`git diff v0.12.0 v0.15.0`); the Go module `spec` added the packages `consentagev1`, `exportv1` and `legalv2` and the parts `Legal`, `ConsentAge` and `Export`, and its `go.mod` is the same; the RUN.md of `spec`, table 4 — Д167 and Д173 (the answers to В10 and В11), Д198 and Д199 (the faster choice: `learner-model` and `engine`).

### 7.4. The move to `spec` v0.16.3

Read 2026-11-12 for the move to `spec` v0.16.3 (R45): `gh release verify v0.16.3 -R altnovatrail/spec` — "Release v0.16.3 verified!", its tag on the commit 664398e; `git -C spec diff --stat v0.15.0 v0.16.3 -- format/ types/formatv0 spec.go go.mod` — empty, and of the Go code only `contracts/novatrail/api/v1` and `types/admissionv0` changed; `spec/.devcontainer/Dockerfile` — `ARG GO_VERSION=1.27.4`, `GO_SHA256_AMD64=7226736c2339719cda28fbf3a54d71aa63bcefe376e98cef00034268dab24840`, `GO_SHA256_ARM64=3392582988a273a87ab5043ee1b2c9c17a2c3c00580bf9326c47a2e0aaaeac92`, and `repo-template/kinds/go/gopls/go.mod` — `go 1.27.4`; the vulnerability GO-2026-6859 — by the 2026-11-11 amendment to R64 spec, `vuln.go.dev` itself was not read here; `GOTOOLCHAIN=go1.27.4 go version` — "go version go1.27.4 linux/amd64", the toolchain from the Go proxy. The repository [altnovatrail/module-kit](https://github.com/altnovatrail/module-kit), branch `main` (`gh api`): at about 20:30 UTC — the commit 660ea37, the submodule `spec` at 58d3f46 (v0.15.0), no open pull requests; at 20:59 UTC — the commit 2f81f3f, the pull request #72 `Move to spec v0.16.3, with Go 1.27.4 (R54)`: R54 module-kit — the move to the tag, Go 1.27.4 and Д215, and the rest of the row 7.1 is its T43.

### 7.5. `spec` v0.17.0–v0.17.2, T49

Read 2026-11-14 for T49, without a move (`gh release list`, `git ls-remote`, `gh api` of the files at each tag and of `compare/v0.16.3...v0.17.2`): the releases [altnovatrail/spec](https://github.com/altnovatrail/spec) `v0.17.0` of 2026-11-13, 19:32 UTC, its tag on the commit 77af949, `v0.17.1` of 23:37 UTC on 54c5f83 and `v0.17.2` of 2026-11-14, 00:29 UTC, on fb45a89; `gh release verify v0.17.2 -R altnovatrail/spec` — "Release v0.17.2 verified!". The numbers the quotes cite are `spec`'s at those tags, not at v0.16.3, which `spec/` holds, so they stand in code.
- From v0.17.0, the decision `R147` — "The criterion of fairness is the simulator's, in `sim`, which holds the criteria and the gate", "It compares populations that differ only by the grade, which sets the start band", "The numbers — how many learners, how many attempts, what "the same" means — `sim` writes down before the first run, as for each criterion"; the card of `sim` — "A third criterion comes from `spec`: the fairness of the learner model by the grade — virtual learners of the same skill and different grades come to the same estimate and the same kind of task within a number of attempts written down before the first run, and its result goes into `dpia 10.2`, before the first child from the US"; the row 4.16 of COORDINATION.md — «| [ ] | 4.16 | `sim` | задачи ещё нет: её записывает сессия `sim` по `R147` | критерий справедливости модели …».
- From v0.17.1, the row 7.10 of COORDINATION.md is among the closed — «экономия токенов …, раньше пилота по слову автора (строка `sim` В13) | 2026-11-13, pull request 54–57, релизы `v0.1.11` и `v0.1.12`»; table 5 of its RUN.md, the triage of the inbox, — «| `sim` | В13 | `e99dec9dd7a0` | закрытая строка 7.10 …»; DEBT.md, Д236 — «`.github` и `sim` держат предел `entry` для всех записей — 50 и 90 КБ ради одной», to be closed «своим переходом на тег с T149: `"entries": {"R18": …}` и `{"R04": …}` вместо `entry`».
- Against v0.16.3, `format/README.md` changes only the language of `LICENSE` and `NOTICE`, and `types/formatv0`, RULES.md and the Go module's `go.mod` are the same; the new Go packages are `types/exportv2` and `types/legalv3`. The template's RUN.md adds to rule 4 the checks of a task without the author, and to rule 9 where the body of a task outside an epic goes — «Когда текущего эпика нет, тело уходит в архив …».

### 7.6. The move to `spec` v0.17.3

Read 2026-11-14 for the move to `spec` v0.17.3 (R51; T56 in [docs/plan/E7.md](plan/E7.md)): the release [altnovatrail/spec](https://github.com/altnovatrail/spec) `v0.17.3` of 2026-11-14, 03:06 UTC, after the laying out of E7; `gh release verify v0.17.3 -R altnovatrail/spec` — "Release v0.17.3 verified!", its tag on the commit 5b766ee. What v0.17.0–v0.17.2 bring `sim` is record 7.5.
- `git -C spec diff --stat v0.17.2 v0.17.3` — `RUN.md` and `draft/README.md` alone: the epic E11 of `spec`'s table 2 — «Задача с моей бумаги» — как в макете.
- `git -C spec diff v0.16.3 v0.17.3`: `spec.go`, `go.mod`, `types/formatv0` and RULES.md the same, and of `format/` only the language of `LICENSE` and `NOTICE`; `tools/docs-check/`, `tools/session/`, `claude-plugins/`, the template's `.claude/` and its kind `go` the same. The template: `.devcontainer/devcontainer.json` gives the container its name in its network — `"--network-alias=spec"`, which `new-repo` makes `"--network-alias=<repository>"` after the repository's own network, the first (Д195 in DEBT.md `spec`); its CLAUDE.md adds that `docs-check-repo` "matches `paths` only against the files git lists, as CI has them, so a part whose files `.gitignore` keeps out of git … has no rule: what must hold for it stays in this file" (T131 tutor); its RUN.md adds to rule 4 «Исключение — проверки задачи по R146 spec: их исполнитель проводит без меня» and to rule 9 «Когда текущего эпика нет, тело уходит в архив последнего сделанного эпика, в его раздел «Задачи вне эпика, сделанные после него»» (the 2026-11-13 amendment to R19 spec).
- `just new-repo-draft sim go 3.9 v0.17.3`, on a clone of `spec` in a temporary directory: its `.devcontainer/` differs from `sim`'s by the name in the network alone, besides `sim`'s own comments (R11, R32); its `.claude/settings.json` and `AGENTS.md` are `sim`'s.
- The key `entries` of the budgets is in no tag up to v0.17.3: `tools/docs-check/budgets.mjs` knows only `files`, `rule`, `entry`, `section`, `sections` and `rows`, and `spec`'s RUN.md at v0.17.3 holds T149 not done — «`tools/docs-check/budgets.mjs` — ключ `entries` в блоках `budgets.json` и в `docs-budgets.json`: предел записи по её номеру, `{"R18": 56}`». DEBT.md `spec`, Д236 — «`.github` и `sim` — своим переходом на тег с T149».
- The rows of DEBT.md `spec` that name `sim` and are new since v0.16.3: Д236 and Д234 — «Догадка, не проверена: `~/.claude/plugins/known_marketplaces.json` на общем томе контейнеров держит один путь `novatrail` — `/workspaces/sim/spec/claude-plugins`, — а его видит только контейнер `sim`», a task of `spec`.

### 7.7. `spec` v0.17.4 and v0.17.5, T60

Read 2026-11-14 for T60, without a move (`gh release list`, `gh release verify`, `git -C spec diff v0.17.3 v0.17.5` after `git fetch --tags`): the releases [altnovatrail/spec](https://github.com/altnovatrail/spec) `v0.17.4` of 04:17 UTC, its tag on the commit 33149be — "Merge pull request #149 from altnovatrail/paper-mockup-home" —, and `v0.17.5` of 05:44 UTC on a14d383 — "Merge pull request #150 from altnovatrail/e12-tasks"; `gh release verify v0.17.5 -R altnovatrail/spec` — "Release v0.17.5 verified!".
- Against v0.17.3 they change CLAUDE.md, COORDINATION.md, DEBT.md, RUN.md, `draft/README.md` and, of `docs/`, architecture, concepts, decisions, product, repositories, sources and the new `plan/E12.md`: the epic E12 of `spec`, «the app's interface by the mockup», the decision `R148`, the term "interface mockup" and the question `Т-41`, sign-in through Apple. These numbers are `spec`'s at v0.17.5, not at v0.17.3, which `spec/` holds, so they stand in code. `RULES.md`, `format/`, `types/`, `go.mod`, `tools/`, `repo-template/`, `.devcontainer/`, `.claude/` and `justfile` are the same, and of architecture the card 3.9 is the same.
- T149 is still not done in RUN.md at v0.17.5 — «`tools/docs-check/budgets.mjs` — ключ `entries` в блоках `budgets.json` и в `docs-budgets.json`», and `tools/` is the same as at v0.17.3, so the key `entries` is in no tag up to v0.17.5.

Read again 2026-11-14, 16:14 UTC, for the laying out of E9 (R54; `gh release list`, `git ls-remote`, RUN.md at `main` by `gh api`): `v0.17.5` is still the latest release, and `main` is its commit a14d383. T149 stands among «Задачи вне эпика» of RUN.md, whose done tasks go to an archive, with no «Сделано», and its body still asks for «`tools/docs-check/budgets.mjs` — ключ `entries` в блоках `budgets.json` и в `docs-budgets.json`», so neither a tag nor `main` has the key `entries`.

### 7.8. `spec` v0.17.6–v0.17.11, T65

Read 2026-11-15, 00:30 UTC, for T65, without a move (`gh release list`, `gh release verify`, `git -C spec diff v0.17.3 v0.17.11` after `git fetch --tags`): the releases [altnovatrail/spec](https://github.com/altnovatrail/spec) `v0.17.6` of 2026-11-14, 16:34 UTC, its tag on the commit 9b94c0a — "Merge pull request #151 from altnovatrail/t153-mockup" —, up to `v0.17.11` of 2026-11-15, 00:07 UTC, on 3080b49 — "Merge pull request #156 from altnovatrail/t157-coordination-rows"; `main` is that commit. `gh release verify v0.17.11 -R altnovatrail/spec` — "Release v0.17.11 verified!".
- Against v0.17.5 they change `.claude/rules/tools.md`, CLAUDE.md, COORDINATION.md, DEBT.md, RUN.md, `justfile`, of `docs/` architecture, decisions, product and `plan/E12.md`, `draft/` with the new `draft/mockup/`, the new `tools/mockup/`, and `tools/docs-check/budgets.json` alone in `tools/docs-check/`, where `spec`'s own budget of its DEBT.md goes from 128 to 132: the interface mockup of `spec`'s epic E12. Against v0.17.3, `RULES.md`, `format/`, `types/`, `go.mod`, `spec.go`, `repo-template/`, `.devcontainer/`, `tools/session/`, `claude-plugins/` and `.claude/settings.json` are the same, and so is the card 3.9 of architecture.
- T149 is still not done in RUN.md at v0.17.11: its body has no «Сделано», and `tools/docs-check/budgets.mjs` is as at v0.17.3, so the key `entries` is in no tag up to v0.17.11.
- COORDINATION.md at v0.17.11 — «| [ ] | 4.15 | `sim` | задачи ещё нет: её записывает сессия `sim` по предложению T50 `engine` | …» and «| [ ] | 4.16 | `sim` | задачи ещё нет: её записывает сессия `sim` …»; table 5 of RUN.md, the rows of other repositories recorded in `spec`, holds one row of `sim` — «| `sim` | В13 | `e99dec9dd7a0` | закрытая строка 7.10 …».

### 7.9. The move to `spec` v0.18.3

Read 2026-11-16, 00:50 UTC, for the move to `spec` v0.18.3 (R60; T66 in [docs/plan/E11.md](plan/E11.md)), after `git fetch --tags` in `spec/`: the release [altnovatrail/spec](https://github.com/altnovatrail/spec) `v0.18.3` of 2026-11-16, 00:42 UTC, its tag on the commit 922886a — "Merge pull request #162 from altnovatrail/e14" —, which is `main` of `spec` too; `gh release verify v0.18.3 -R altnovatrail/spec` — "Release v0.18.3 verified!". `git -C spec log --oneline v0.17.3..v0.18.3` — 53 commits, of them T149, T162, T164, T165 (R153 spec), R154 spec, R155 spec, T166 and the epic E14 with its R156 spec, the orchestrator of «Сделай эпик E…».
- What the simulator reads is the same: `git -C spec diff v0.17.3 v0.18.3` changes neither `spec.go` nor `go.mod` of the module `spec`, nor `format/` nor `types/formatv0`; it renames `types/exportv2` to `types/exportv3` and changes `contracts/`, which the simulator does not import. The hash of `go.mod` of the module `spec` in `go.sum` is the same as at v0.17.3.
- `tools/docs-check/budgets.mjs` knows the key `entries` — «entries names a limit of its own for one entry by its number» —, and `spec`'s block `repository` of `budgets.json` keeps `"entry": 40`. DEBT.md `spec` moved to `run/DEBT.md`; its rows that name `sim` are Д34, Д88, Д146, Д170, Д236, Д248 and Д250. Д236 — «`.github`, `sim`, on move to a tag with T149: `"entries": {"R18": …}` / `{"R04": …}` instead of `entry`»; Д248 — «`gopls-lsp` loads via devcontainer `containerEnv` `"CLAUDE_CODE_PLUGIN_DIRS": "${containerWorkspaceFolder}/spec/claude-plugins/gopls-lsp"`, not marketplace `novatrail`. … same change drops `extraKnownMarketplaces`, `enabledPlugins` from `.claude/settings.json`»; Д250 — the plan's tables 2, 3 and 4 to `run/EPICS.md`, `run/QUESTIONS.md` and `run/DEBT.md`, and the plan in English, terse, «Each repo on its move to a `spec` tag with `R153`, `R154` and `R155`». Д146 — «Pending, move to cached `.github` release: `api`, `engine`, `infra`, `sim` (on `v0.1.21`)», and Д88 — the group `shared-workflows` in `dependabot.yml`: neither waits for a tag of `spec`. The releases of `.github` go up to `v0.1.24`, of 2026-11-13 (`gh release list -R altnovatrail/.github`).
- RULES.md: the rules and tables of the work plan are «terse, for the session» (R154 spec); table 3 is «`run/QUESTIONS.md`, or its RUN.md until the repository moves its tables to `run/` (R153 spec)»; «Do not commit: the author commits after review, unless the author's word, as «доведи до мержа», says otherwise. «Сделай эпик E…» says it for each task it runs» (R156 spec).
- The card 3.9 of architecture now ends its item "To decide when detailed" with «`sim` decided it before its first run (R55 sim): an evaluation of the base, whose result is in `dpia 10.2`, and the constraint `fair` of the criteria of the model and of the engine, which holds a candidate no worse than the base (R57 sim). The base of `learner-model` v0.1.8 and `engine` v0.1.6 does not hold it (R58 sim; COORDINATION.md, closed row 4.16)».
- Table 5 of `spec`, `run/INBOX.md`, records every row «в реестр spec» of table 3 of `sim`, В15, В16 and В18–В24, each with its fingerprint; `fingerprintOf` of `tools/inbox/inbox.mjs` gives the same nine fingerprints on RUN.md of `sim` at 59885be. COORDINATION.md — «| [ ] | 4.18 | `sim` | E10 «Суд круга справедливости модели» таблицы 2: его разворачивает сессия `sim`, когда `learner-model` возьмёт круг строки 4.17 (R59 sim) | …».
- `just new-repo-draft sim go 3.9 v0.18.3`, on a clone of `spec` in a temporary directory: its `.devcontainer/` differs from `sim`'s by the base image `mcr.microsoft.com/devcontainers/base:3.0.13-noble` at `sha256:48be4adec1c1ecbe030064700aef85382f5205786a8f82937fc44a7473717eae`, the comment on the stages, the comment on what the image copies and `CLAUDE_CODE_PLUGIN_DIRS`, besides `sim`'s own comments (R11, R32); its `.claude/settings.json` has no `extraKnownMarketplaces` and no `enabledPlugins`; `GO_VERSION`, `GOPLS_VERSION`, `.devcontainer/gopls/` and AGENTS.md are `sim`'s. Its RUN.md is in English, with the tables in `run/`, and its CLAUDE.md names `run/`.
