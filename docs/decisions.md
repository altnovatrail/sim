# sim decision log

> Why sim is built the way it is. Each entry says what was decided, why, what was rejected and where it is described.
>
> **Rule:** do not re-propose what was rejected without new facts. If a decision has to change, ask the author and update this file. An entry is never rewritten: an amendment with its date at the end of the entry corrects it, or a new entry replaces it, and the old one stays where it is, marked "Superseded".
>
> The repository was set up 2026-10-19 from the template of `spec` v0.1.0 (R21 spec), moved to `spec` v0.2.9 2026-10-20 (R01), and 2026-11-07 to v0.9.4 (R02) and the same day to v0.9.5 (R06), and 2026-11-08 to v0.9.8 (R07) and the same day to v0.9.9 (R08); the code epic E2 was laid out 2026-11-08 (R09), and its frame — the Go module, CI and gopls in the container — is recorded by R10 and R11, and the seeds and distributions by R12; the executor leads T11–T23 by themselves, up to the merge (R13); the simulator reads and translates the packages itself (R14); the simulator's learner — R15, the run — R16, the invariants — R17, the populations — R18, the synthetic packages — R19, the populations `S-`, `F-`, `W-` — R20, `A-`, `L-`, `E-` — R21, the model's measures — R22, the engine's measures — R23, the ceilings and two versions side by side — R24, the criteria and the report — R25, the gate — R26, the gate in CI of `modules` — R27, the paths of the targets and the decisive run — R28, the release — R29. Epic E2 was done 2026-11-09, and epic E3, the first candidates, was laid out the same day (R30). The packages are pinned anew — on `modules` v0.1.4 (R31). 2026-11-10 `sim` moved to `spec` v0.12.0 — format 0.4, the packages of `modules` v0.2.0 (R32). The simulator's runs are fitted into an hour: 1,000 learners a cell and the choice anew among the first 10 (R33). A candidate's decisive run is run by the same recipe, on its commit (R34). The first judgement is the one whose candidates are ready earlier, and the judgement's report lies in the repository (R35). The first judgement of the engine's candidates ended in the exit, and the base stayed (R36). The first judgement of the model's candidates did too (R37). Epic E4, «Стенд за библиотеками», was laid out the same day (R38). 2026-11-11 the base moved to `engine` v0.1.6, and the gate and the decisive run grew to 150 and 1,500 learners (amendments to R04, R31 and R33); the same day `sim` moved to `spec` v0.15.0 (R39), and a judgement began to hold the candidate's second library the same as the base's (R40). The simulator reads the check kind `time` and the leading question of format 0.4 on `learner-model` v0.1.5 (R41) and is pinned on the packages of `modules` v0.2.5 with them, with the paths and bands anew (R42); the base moved to `learner-model` v0.1.6 with the same bits, and the paths were read anew (R43); epic E5, the sizes and the model's second round, was laid out 2026-11-11 (R44). 2026-11-12 the gate grew to 250 learners, and the decisive run to 2,000 (amendments to R04, R31 and R33); the same day `sim` moved to `spec` v0.16.3, the line of epic E10 of `spec`, with the way of working by its template, and the translation of `docs/` was split into the tasks T51–T53 (R45). How it comes to differ from the template over time — by tools in the devcontainer, by rules in RUN.md — is recorded here as a decision, as R46 canvas-api recorded the differences from its model. Next number — R46.
>
> Here entries are cited as `R03`, in other repositories as `R03 sim`. `spec`'s decisions are cited as "R18 spec", MathTrail's decisions as "R37 canvas-api", and in a list as "canvas-api: R37, R38". Each log has its own numbers.
>
> Open questions that concern several repositories are in `spec`'s registers: О- in product 13 spec, П- in privacy 9 spec, Т- in repositories 8 spec. What is still to be decided in the component itself is in [architecture.md](architecture.md), section 2.

## The tag `spec`

**R01. `sim` stands on `spec` v0.2.9: the submodule `spec/` is at this tag, the links to `spec` pinned to a tag name it, and what came from the template is as in the template v0.2.9. Since then the documents are checked by `just -f spec/justfile docs-check-repo`.**
- **What the tag brings `sim`** (`git diff v0.1.0 v0.2.9` in `spec/`, read 2026-10-20; Д6 in table 4 of RUN.md `spec`):
  - R27 spec — the generator of randomness and the derivation of seeds from `spec` — one for the whole platform: `module-kit`, widgets, `engine` and `sim`. The simulator derives the seeds of virtual learners with it, not with a scheme of its own — so card 3.9 spec now says too —, and passes the generator's test vectors in CI, as they do (Т-18). The scheme of standalone's simulator — FNV-64a and PCG, the function `seeded` in `tools/learners/child.go` — stays in `sim` only if epic E2 of `spec` chooses it for the whole platform: it chooses the algorithm by the recommendations of T02 module-kit and T03 engine and takes this scheme into account. This settles the question of T01 "the simulator's own generator or the generator of `spec`";
  - R28 spec — `learner-model` sets the types of the graph, and `engine` takes them from it; they are simple and not from `module-kit` (architecture 3.8 spec). On the simulator `sim` passes the graph to both libraries, and it translates the graph from a module package into one form, not two — this is the argument of R28 spec itself;
  - R23 spec — the success corridor is a rule of `engine`, and `learner-model` gives the probability of a correct answer. The "mastered" threshold, which in standalone is tied to the corridor's target, R23 spec left open: where it lives is decided by T01 engine and T01 learner-model, and a change to one goes through `sim` together with the other. This is an input of T02 — whether one criterion serves both libraries — and of T03 — what is put into the CI of each and in what order their releases come out;
  - R29 spec — the list of the engine's reason codes is a contract in `spec`: T03 engine proposes the codes, and epic E3 of `spec` records them. The simulator checks the reason code of each choice (concepts 12 spec), so this is an input of the measure of T02;
  - R30 spec — the CI checks are the same for all repositories, and `sim` will have them even before code: the secret scan, DCO, CodeQL, Scorecard and the documents check — shared checks from `.github`. Scorecard is shared if its server accepts the result of the shared workflow, and otherwise it is a copy in the repository (architecture 3.11 spec); for `tutor` the server accepted it (T16 `spec`). All the checks but Scorecard hold the merge, the SonarCloud quality gate among them, and a pull request goes into `main` only by a merge commit. T21 `spec` will connect them: its session writes the CI files here, together with `dependabot.yml` and the exclusion of `spec/**` for SonarCloud (Д7 in table 4 of RUN.md `spec`), and it also sets the ruleset of the branch `main`. Until then `sim` has no CI. The checks of the language will come with the first code of `sim` (repositories 1 spec);
  - architecture 3.11 spec — which workflows are shared and which are the repository's own was decided by `.github` (R03 .github). Shared are the checks before code and, with the first code of `module-kit`, the wrappers of Go and TypeScript around the repository's recipes, with the license check, and the release with a signature. The gate of `sim` is a workflow of its own in the CI of `learner-model`, `engine` and `modules`, not a shared one, and the test vectors of the generator of `spec` in `sim`, as in `module-kit` and `engine`, are a workflow of its own too (architecture 3.2 .github). The table of `.github`'s own workflows is written by the recommendation of Т-16 — `go.work`. [architecture.md](architecture.md), section 2.3, and T03 said that the detailing of `.github` still decides this;
  - the opening of architecture 3 spec — an item of the detailing beyond the card is entered into the card by the repository's own "outcome of detailing", with a row of its table 3 «в реестр spec», once the decision is taken. This answers the question that T08 `spec` left when it set up `sim`: whether to add Т-16 to "To decide when detailed" of card 3.9 spec. For `sim` this is section 2.3, and T04 writes the change of the card;
  - R31 spec — a release of `spec` comes out by itself, after green CI on `main`: the next patch of the line, immutable, with a GitHub attestation. A repository still moves to a new tag by a task of its own;
  - R32 spec — a repository's documents are checked by the `docs-check` of the `spec` its submodule stands on: the links, its own sections and log, the tables «В» and «Д», and the citations of `spec` by the submodule at its tag. An example of a citation is written in inline code, so that the check does not take it for a citation (RULES.md, Citations);
  - R33 spec and rule 12 of the template — before a new task the findings of SonarCloud, CodeQL and Scorecard from the last one are closed. A finding that a change cannot close stays open with a reason, and the reasons shared by all repositories are written in R33 spec;
  - RULES.md, but for the rule about examples of citations: a development or CI tool that ships nowhere and puts none of its own code into what ships may carry another OSI-approved license too, or free terms for open source, but only as a named exception: in the license check or in the decision that brings it in (Hard rules; R12 spec). This concerns the tools of the future code of `sim`;
  - R24, R25 and R26 spec — about other components: the key for eval, the page of samples, the account with the owner's profile and the children's profiles. They do not change the documents of `sim`;
  - devcontainer — the base image `devcontainers/base` 3.0.10 by digest; Claude Code is installed without the packages' scripts, its `install.cjs` is run by name, and `claude --version` stops the build on a stub; `git.alwaysSignOff` — a commit from Source Control is signed off by DCO (R22 spec, R30 spec).
- **What is changed here:**
  - `spec/` — at v0.2.9, commit 5eafb1a. The links to `spec` in README.md and the tag in CLAUDE.md and RUN.md — v0.2.9. That the repository was set up from the template v0.1.0 the opening of the log still says: this is history;
  - five examples of citations of the old template are in inline code: two in CLAUDE.md, Documents, two in «Ссылках» of RUN.md and one in the opening of the log. Rule 6 of RUN.md and the section Environment in CLAUDE.md name `docs-check-repo`, and rule 12 came into RUN.md. These places are as in the template v0.2.9;
  - `.devcontainer/` — the one that the recipe of `spec` v0.2.9 (`editDevcontainer`) builds for `sim` of the kind `go`. The name of the container and of the network, `postStartCommand` and the Go extension are its own, as on v0.1.0. The move brings in no differences from the template;
  - what the tag decided for `sim` is named in its documents:
    - T01 — R27 spec; T02 — R23 and R29 spec; T03 — R23 spec and R03 .github; T04 — the change of the card for section 2.3;
    - [architecture.md](architecture.md): section 1 — R27, R28 and R30 spec; the opening of section 2 — the change of the card; section 2.1 — R27 spec; section 2.2 — R23 and R29 spec; section 2.3 — R23 spec and R03 .github;
    - CLAUDE.md, The component, — R27 spec; «Почему такой порядок» and «Не входит» in RUN.md.

— *Why now:* T21 `spec` connects the shared checks from `.github` here (R30 spec). The shared `documents.yml` builds the stage `toolchain` from the `.devcontainer/Dockerfile` of its caller and runs `just -f spec/justfile docs-check-repo` in it ([`documents.yml`](https://github.com/altnovatrail/.github/blob/ff2a359/.github/workflows/documents.yml), the release v0.1.0 of `.github`, commit ff2a359, read 2026-10-20). And this recipe is in `spec` only from v0.2.3 (R32 spec): on v0.1.0 the required check "Documents / docs-check-repo" would be red and would hold any merge. v0.2.9 is the latest release of `spec` as of 2026-10-20, and immutable (`gh release view -R altnovatrail/spec`); `gh release verify v0.2.9 -R altnovatrail/spec` confirmed it by its attestation, and the tag is on commit 5eafb1a. Д6 in table 4 of RUN.md `spec` also says to take the latest release at the time of the task.

— *Rejected:*
- stay on v0.1.0 until the first code: the documents check from T21 `spec` would hold any merge, T01 would read `spec` without R27 spec, and a link to a section of `spec` that the tag does not have, or a gap in the log, would pass unnoticed (Д6 in table 4 of RUN.md `spec`);
- v0.2.3–v0.2.8: each has `docs-check-repo`, but Д6 says to take the latest release. Rule 12 of the template came in v0.2.7, R33 spec in v0.2.8, and T21 `spec`, which will connect the checks here, only in v0.2.9;
- move without the devcontainer, or carry over its changes selectively: the container of `sim` is the container of `spec` (CLAUDE.md, Environment), and a difference would have to be recorded as a decision, and there is no reason for one;
- examples of citations in guillemets, as before: `docs-check-repo` checks a number in guillemets against the log, and R32 spec rejected not checking them;
- rewrite T01–T04 for the new tag: the tag decided one fork for `sim` — the seed generator in T01 —, and the rest brings the tasks inputs that they sort out themselves. So the tasks have short inserts with a reference to the R.

— *Consequence:* T05 goes first in table 1 of RUN.md, and T01 depends on it. The devcontainer image was not built here — Д1 in table 4 of RUN.md: T01 waits for Rebuild Container too. The commit, push and merge of this move are the executor's, by a one-off instruction of the author recorded in T05; rule 7 of RUN.md stays as in the template. The CI files — the calls of the shared workflows, `dependabot.yml`, `.sonarcloud.properties` — are written here by T21 `spec`, not by a task of this repository. — RUN.md, T05, T01–T04; CLAUDE.md; README.md; `.devcontainer/`; [architecture.md](architecture.md), sections 1, 2, 2.1, 2.2 and 2.3.

**R02. `sim` stands on `spec` v0.9.4, the latest release of `spec`: the submodule `spec/` is at this tag, the links to `spec` pinned to a tag name it, and what came from the template is as in the template v0.9.4. The tag decided for `sim` much of what the tasks of E1 held as options — the modules of the first slice, the format and the package, the seed generator, the reason codes, the level per module, forgetting, the weights of evidence and the anchor skill —, and card 3.9 spec now names five items of the detailing instead of two. The tasks T01–T04 are narrowed to what the tag left to `sim`.** The author's instruction 2026-11-06: move to the latest release of `spec`, in one task, and the commit, pull request and merge are the executor's, as in T05. This is a task of its own, T06 (RULES.md, Across repositories).
- **What the tag decided for `sim`** (`git -C spec diff v0.2.9 v0.9.4`, read 2026-11-06; the parts of `sim` in Д9, Д12, Д15, Д18, Д20, Д21, Д30–Д32, Д40, Д42, Д44, Д46, Д53, Д59, Д80, Д82, Д83, Д86, Д106 and Д107 of table 4 of RUN.md `spec`; [COORDINATION.md](../spec/COORDINATION.md), rows 4.1, 4.2 and 4.5):
  - **card 3.9 spec** — five items of "To decide when detailed": the profiles, and forgetting among them; the criterion and the bands with their own measures of forgetting; the weights of evidence — learners with pages of a workbook; command of the anchor skill — a learner with both maths and the measure of a false "knows"; what the gate of `sim` in the CI of `modules` includes. The first four are sorted out by T01 and T02, the fifth by T03 ([architecture.md](architecture.md), section 2). The measures named there are proposals of `learner-model`, not decisions: T02 decides (Д82 `spec`);
  - **the three modules of the first slice** (R41 spec; О-1 decided): olympiad maths 1–6 on a sheet, grade 3 maths on a workbook and short essays in English for any age, with feedback by a rubric. The profiles are children with a sheet and with a workbook and a learner of any age with essays instead of an adult with a text canvas (Д18 `spec`). There are no language modules in the first slice, and the simulator checks forgetting on the two maths (architecture 3.9 spec). On the simulator the engine holds the corridor on all three graphs where there is a verdict, and goes by the author's order where there is none (product 10 spec; concepts 4 spec);
  - **the format and the package:** the format has its own version, and while it is v0, the platform reads only the current one — from v0.9.0 this is 0.3 (R44, R101 spec; Д20, Д107 `spec`). `sim` reads the server part of the package, `server.json`, — the canonical view, with the tiers, the anchor skills, the bands and difficulty of the tasks and the numbers of the instances of templates, — by the types `formatv0.CanonicalView` of the Go module `spec`, without `module-kit` (R45, R49–R51, R62, R72 spec; architecture 3.9 spec; Д21, Д30–Д32, Д44, Д46 `spec`). The Go code takes the same tag of `spec` as the version of the Go module: the module's version and the submodule move together, and the module `spec` is excluded from Dependabot `gomod` (R64 spec; Д46 `spec`). This comes with the first code of `sim`. Д59 and Д107 `spec` say that `sim` reads the packages through `module-kit` — this diverges from R62 spec and the card, and the question went out as row В1 of table 3 of RUN.md, «В реестр spec»;
  - **the seed generator** (R59 spec; Т-18 closed; Д40 `spec`): PCG, as in Go's `math/rand/v2`, and the stream of a place of randomness is derived by FNV-1a 64 from the seed and the names of the place. The function `seeded` of standalone's simulator is a particular case of this derivation: the place of the main seed with the names "experiment, unit, assignment", so the simulator's numbers are kept. The normal and exponential distributions are not in the format, and `sim` is responsible for their portability (architecture 3.9 spec). The simulator chooses the instance of a template by a seed with the same choice as the platform, and for the engine a template is one task (R60 spec; Д42 `spec`);
  - **the answers of `learner-model`** (card 3.8 spec; Д82, Д86 `spec`): the level is per module, and a learner of two modules has two levels, so the profiles have no option "per area" (R03 learner-model); guessing is per item, by its answer slots (R03 learner-model); forgetting is the retention and stability of a skill, and the profiles have a calendar and laws of forgetting other than the model's law (R04 learner-model); the weights of evidence are the item's first verdict, help and the share of the answer of one attempt, and the profiles have pages of a workbook with dependent answers, help and repeated tasks instead of «серий однотипных пунктов» (R06 learner-model); command of the anchor skill (R07 learner-model);
  - **the "mastered" threshold and the corridor's target** (the 2026-11-02 amendment to R23 spec; Т-16): the threshold is a constant of `learner-model`, and a new target of the corridor of `engine` is a new constant of the model, and both changes go through `sim` in one run. How the gate runs such a change when a pull request replaces one library is for T03 `sim` to decide;
  - **the gate in CI of `modules`** (the 2026-11-01 and 2026-11-02 additions to Т-16; R76 spec; Д80 `spec`): until the first tag of `sim` there is no gate of the graph in the CI of `modules`; "the graph changed" means the canonical text of `skill_graph` differs from the newest published version of the same format version, and a move of the format goes without the simulator. What the gate includes — the graph or everything the simulator reads from a package — T03 decides, and until its answer, the graph;
  - **the engine's inputs** (card 3.8 spec; Д83 `spec`): the lesson runs and the closed attempts of tasks with a rubric. On the simulator `sim` gives them, as it gives the graph;
  - **the reason codes** (R100 spec; Д106 `spec`): the closed list `ReasonCode` of eight values, without parameters ([contracts text](../spec/contracts/README.md), section 15). At each choice the simulator checks a code from this list;
  - **the cross-repository plan** (R76 spec; Д53 `spec`): row 4.1 of COORDINATION.md puts the move as the first task of the course of `sim`, row 4.2 epic E1, and 4.5 the code epic and the simulator's first tag;
  - **the rules and the template:** an archive of the Dockerfile is checked against its sha256 sum (R34 spec; Д12 `spec`); the packages of `apt` come from a snapshot of the Ubuntu archive (R40 spec; Д15 `spec`) — the Dockerfile of `sim` has none; code comes with its tests, and the coverage in Codecov holds the merge (R48 spec); a file of a NovaTrail release enters code or CI only after `release-verify` (R77 spec); a row «В реестр spec» opens a question in table 3, and the session of `spec` finds it by itself (R81 spec); the normative texts for authors are in English (R42 spec); a review is one round (R35 spec) — rule 6 of RUN.md was already so; the status badges in the README (R48 spec);
  - **devcontainer:** the base image `devcontainers/base` 3.0.12 by digest, Go 1.27.3 and buf 1.75.0 in the stage `toolchain`, a sha256 sum for each archive (R34, R64, R83 spec). The container in CLAUDE.md, Environment, is "Node.js, Go, buf and just";
  - **`.github`:** `documents.yml` will move to the script `toolchain.sh` once all its callers stand on `spec` v0.2.17 or newer, and it was waiting for four repositories, `sim` among them (Д5 `.github`, by [COORDINATION.md](../spec/COORDINATION.md)). Now `sim` leaves them;
  - the other decisions R36–R102 spec and the amendments to earlier ones — the interface packages, the path to AI, the rubric and eval, the contract `app` ↔ `api`, admission by age, the database of `api` — are about other components, and enter `sim` only by what is named above.
- **What is changed here:**
  - `spec/` — at v0.9.4, commit e803736. The tag in CLAUDE.md, README.md, the opening of RUN.md and the status of [architecture.md](architecture.md) — v0.9.4. That the repository was on v0.1.0 and v0.2.9 the opening of the log and R01 say: this is history;
  - the template: a three-way merge of the template v0.2.9 → v0.9.4 with CLAUDE.md, README.md and RUN.md. CLAUDE.md and RUN.md merged without conflicts, README.md with one: the status badges went above the description. What changed is the container in CLAUDE.md, Environment, the mark «В реестр spec» in rule 2 and the language in rule 10 of RUN.md. The template's AGENTS.md, LICENSE, `.claude/settings.json` and `docs/architecture.md` did not change between the tags;
  - `.devcontainer/` — the one that the recipe of `spec` v0.9.4 (`editDevcontainer`) builds for `sim` of the kind `go`; `devcontainer.json` did not change, and the Dockerfile is as in `spec`. The move brings in no differences from the template. The image was not built here — Д2 in table 4 of RUN.md;
  - [architecture.md](architecture.md) — the status, section 1, section 2 — the opening and subsections 2.1–2.3: what the tag decided is named decided, and the items of the card are laid out over the subsections;
  - RUN.md — table 1 and T06, «Почему такой порядок», the inserts in T01–T04; table 3 — В1 «В реестр spec»; table 4 — Д1 closed, Д2 opened;
  - CLAUDE.md, The component, — R59, R62 and R72 spec.

— *Why now:* the author's instruction 2026-11-06. Row 4.1 of [COORDINATION.md](../spec/COORDINATION.md) puts the move as the first task of the course of `sim`, and on v0.2.9 T01 would plan the profiles "per area" as options and an adult with a text canvas, T02 a criterion without measures of forgetting, and the simulator would derive the seeds not by R59 spec (Д18, Д40, Д82 `spec`). v0.9.4 is the latest release as of 2026-11-06, and immutable (`gh release view -R altnovatrail/spec`); `gh release verify v0.9.4 -R altnovatrail/spec` confirmed it by its attestation, and the tag is on commit e803736.

— *Rejected:*
- a tag of the line 0.6, as row 4.1 of COORDINATION.md names: format 0.3 came out only in the line 0.9, and on 0.6 `sim` would read format 0.2, and a second move would be needed (Д107 `spec`);
- v0.9.0–v0.9.3: the line is the same, but the author asked for the latest release;
- split the move into two tasks — the move itself and the inserts in T01–T04: the author chose one (rule 11 of RUN.md);
- rewrite T01–T04 anew: the tag decided the forks, but not the profiles, the measures and the gate — the tasks still decide them. So in the tasks what is decided is named decided, with a reference to the R, and the items that no longer exist are removed;
- keep Д1 until T01: the check it named passed in this task — the image 3.0.10, `claude --version` answers, `bin/claude.exe` is an ELF binary —, and the new image waits for its own rebuild anyway, Д2;
- carry over the changes of the devcontainer selectively — by the same arguments as in R01.

— *Consequence:* T06 stands first in table 1 of RUN.md, and T01 depends on it and on Rebuild Container (Д2). The commit, push and merge are the executor's, by a one-off instruction of the author recorded in T06; rule 7 of RUN.md stays as in the template. With the first code of `sim` come the tests and Codecov (R48 spec), the version of the Go module `spec` together with the submodule and its exclusion from Dependabot `gomod` (Д46 `spec`), and CodeQL for Go (Д9 `spec`). — RUN.md, T06, T01–T04, tables 1, 3 and 4; CLAUDE.md; README.md; `.devcontainer/`; [architecture.md](architecture.md), the status and sections 1, 2, 2.1–2.3.

— *Amendment 2026-11-12:* Go of the stage `toolchain` is 1.27.4, not 1.27.3: the move to `spec` v0.16.3 brought it, as to every Go repository (R45; the 2026-11-11 amendment to R64 spec).

## Virtual learners

**R03. The simulator's virtual learner is a model of the truth of its own, not the code of `learner-model`: the base population holds the model's assumptions, and each population changes one of them. Standalone's nineteen populations are carried over as assumptions, and the simulator's code is written anew; the main one is `G2-half` with a calendar. The new populations touch forgetting, the pages of a workbook, the anchor skill, lessons and essays. The criterion and the gate run on the published packages of the three modules, by version, and until their releases on the synthetic packages of `sim`, only for code and tests. The miss of difficulty is the task's, one per run. The seeds are places of randomness of the format, and the held-out learners take a seed of their own.** The author's answers 2026-11-07 in the planning of T01 (RUN.md); the rest is a proposal of the execution. The whole design is in [architecture.md](architecture.md), sections 3.1–3.4.
- **Decided by the author:**
  - standalone's populations are the example: their assumptions as a list "carried over, changed", and the simulator's code is written anew for `learner-model` and `engine`. The MIT notice (R12 spec) only if the code epic carries over a piece of code word for word;
  - the main population is `G2-half` with a calendar: growth of one logit per 200 answers, an assumption of the simulator, not a measurement; `G2` is a stress test. The main one is for both maths;
  - the graphs: a population is set over the canonical view of any package; the criterion and the gate run on the published packages of the three modules, pinned by version; until the releases of `modules`, on the synthetic packages of `sim` in the form of the three modules, for code and tests, not for the judgement of candidates;
  - there is no check of the simulator's numbers against live learners at stage 1 — Д3 in table 4 of RUN.md;
  - T01 is one task, not two (rule 11 of RUN.md).
- **A proposal of the execution, not the author's answer:**
  - the words: population and virtual learner, and "profile" only in the sense of the vocabulary, and "generator" is not said of populations. Card 3.9 spec calls the populations profiles — row В2 of table 3 of RUN.md, «В реестр spec»;
  - the truth of the base population: the grade and the starting band as in the model, the module level is the shift of the starting band plus N(0, 1), the skill offset is from N(0, 0.5²), the task's hardness with a miss from N(0, 0.5²), the answer of an item by the model's formula with guessing by answer slots and a slope of 1, retention by the model's law with a stability that the learner's own answers move (R03, R04, R06 learner-model). The truth's stability goes by the model's law, once per attempt and no less than a day. In an attempt there are three rungs of the hint ladder, +0.5 to the logit each, and from the fourth wrong answer the showing of the solution, for a child by an adult's action (concepts 10 spec; R82 spec); a lesson goes to its end; an essay is one submission and "Done";
  - the calendar: a day is a session with a probability of 5/7, a session is ten attempts from 15:00 UTC, time is in whole seconds UTC; a run is 200 answers of items with a verdict, for the populations of forgetting 180 days, and a break ends no later than the 120th, for essays no longer than 60 days;
  - standalone's populations with changes: G4 — the author's miss on the whole band, `G4-band`, instead of the second "host"; G7 — guessing for a choice of 0 or 1.5/k instead of a floor of 0 or 0.3; G8 — a technique on request instead of a hint; `G0-topics*` — `G0-skills*`;
  - the new populations — by the numbers of `learner-model` that wait for the simulator, and by the items of the card: `S-` — the step of bands, guessing of `predicate`, an adult without a grade, "unreadable"; `F-` — a calendar with breaks and laws of forgetting other than the model's law, a learner who does not forget, and a fall after a break; `W-` — a page with a common shift, an error on the whole page, a bad day, a repeated task; `A-` — a learner of both maths; `L-` and `E-` — lessons and essays;
  - the miss is the task's: z is one per run, and the populations of the miss change the multiplier at it; a template has one miss for all its instances; the instance is chosen by the attempt's seed with the same choice as the platform (R60 spec);
  - the seeds: places with the names "experiment, unit, purpose", the unit is `<population>/<number>` or `task/<module ID>/<task ID>`, the purposes are `params`, `calendar`, `answers`, `retries`, `attempts`, `attempt` and `choice`, and for a task `miss`; first answers and repeats are different streams, so that the repeats do not shift the first answers; the held-out learners have a seed of their own and an experiment name of their own; the normal and exponential distributions give the same bits on any build. The seeds themselves come with the criterion of T02.

— *Why:*
- **A truth of its own, not the model.** The simulator measures how far the estimate of `learner-model` and the choice of `engine` stray from the learner whom the simulator knows. Were the truth the model itself, the simulator would measure the model on itself; so the base population holds the model's assumptions in its own code, and the populations depart from them — as G0 and the others do in standalone ([architecture.md](architecture.md), section 4).
- **An example, not a carry-over of code.** Standalone's simulator runs the product's path — `profile.New`, `tutor.Next`, `Profile.Record` — with topics, guessing 0.2 and "every tenth answer starts a new day". In `sim` the model and the choice are two libraries, the level is per module, guessing is by answer slots, time is a calendar (R03, R04 learner-model), so little would be left of standalone's code, and its design would be dragged along with it. The populations' assumptions are carried over whole: on them standalone's simulator chose the model with which `learner-model` begins (R12 learner-model).
- **`G2-half`.** The criterion's targets are read on the main learner (R163 standalone), and the speed of live learners is unknown, as in standalone. One speed, named an assumption, is more honest than three, and `G2` and `G2-fading` hold it from both sides.
- **The calendar.** The model's forgetting goes by time, not by the number of the answer (R04 learner-model), and the card requires days, sessions and breaks (architecture 3.9 spec). The breaks of the `F-` populations give pauses from a day to two months, on which the calibration by the length of the pause is read.
- **The miss is the task's.** In NovaTrail the author writes a task, once, and all learners see the same task; a miss at every issue, as with the chat model in standalone, would average out over the learners and hide a task that is in fact harder.
- **Synthetic packages before the releases.** `modules` has no release, and grade 3 maths and the essays are not even in the repository ([COORDINATION.md](../spec/COORDINATION.md), row 3.6). The simulator's code and its tests can go before them, but candidates cannot be judged on packages that the platform does not have: the criterion is read on what the learners see.
- **A seed of their own for the held-out learners.** The places of one seed share the low word of the state, and FNV-1a does not guard against collisions; the format says outright to take a seed of their own for streams that need to be apart ([format text](../spec/format/README.md), section 22). Standalone did the same — the seed 20261003 and the experiment `held-out`.

— *Rejected:*
- carry over `tools/learners` and rework it: the product design of standalone would come along with the code, and the MIT notice would become required;
- no main population — the criterion's targets on all populations equally: there would be no "typical learner" on whom the targets are read; a main one of its own for each module — three arbitrary speeds instead of one;
- wait for the published packages: T02 and the simulator's first code would wait for the content of school maths and the essays;
- check against live learners now: the simulator holds no data of live learners (architecture 3.9 spec), and `api` stores them, stage 2 (architecture 1 spec; repositories 2.1 spec);
- a miss at every issue, as in standalone, and "every tenth answer starts a new day": see "Why";
- the truth is the code of `learner-model`: the simulator would measure the model on itself;
- populations "per area" and an adult with a text canvas: the tag v0.9.4 removed them (R02).

— *Consequence:* T02 writes down the measures, targets and bands over the populations of section 3.2, and with them the seeds of the runs. The code epic builds the virtual learner, the calendar, the populations and the synthetic packages. Д2 closed, Д3 — the check against live learners; В2 «В реестр spec» — the populations in card 3.9 spec. — RUN.md, T01, tables 3 and 4; [architecture.md](architecture.md), the status, section 2.1, sections 3.1–3.4 and 4.

— *Amendment 2026-11-07:* the divergence of the card from the vocabulary, row В2 of table 3 of RUN.md, `spec` decided by the recommendation: card 3.9 spec says "populations", and the vocabulary got "virtual learner" and "population" (concepts 15 spec; R103 spec). The same day: `modules` put out its first release, v0.1.0, with the packages of all three modules. The synthetic packages are still only for code and tests, and the criterion and the gate run on the published ones; "no graph" stays where a population has no published graph. — R06.

— *Amendment 2026-11-08:* the graph of `G7` is grade 3 maths, not olympiad maths: in the published package of olympiad maths all the answer slots are `number`, while a choice of three options exists in school maths and the essays ([architecture.md](architecture.md), section 4), and guessing by answer slots is read on maths. The two sides of one assumption go by the parity of the learner's number. — R18; [architecture.md](architecture.md), section 3.2.

## The criterion and the bands

**R04. The simulator judges changes by two criteria — of the learner model and of the engine — on a common basis: invariants with a tolerance of zero, measures by the learner's truth, the base — the first released versions of the libraries, the ceilings — the model's oracle and a choice by the true probability, the targets — a share of the path from the base to the ceiling, the tolerance — the larger of the author's tolerance and the simulator's resolution, a confirmation once on its own held-out learners, and the exit to the base. The CI gate is the guard: 300 working learners of the gate's populations, numbers in two-sided bands, which a pull request to `sim` writes anew. The criterion and the gate run on the published graphs, and the seeds of the runs are pinned.** Written down 2026-11-07 — before the first run of a candidate: there is no code of the simulator, of `learner-model` or of `engine` yet. The author's answers 2026-11-07 in the planning of T02 (RUN.md); the rest is a proposal of the execution. The whole design is in [architecture.md](architecture.md), sections 3.5–3.8.
- **Decided by the author:**
  - two criteria on a common basis: each has its own constraints and score, and the measures, the resolution and the tolerance, the choice, the confirmation and the exit are shared. A joint change of the "mastered" threshold and of the corridor's target (the 2026-11-02 amendment to R23 spec) is a candidate pair, and both criteria judge it on one run;
  - the base is the first released version of each library (R12 learner-model; В24 `learner-model` «в реестр spec»); the model's ceiling is the oracle, which knows the truth, and the engine's ceiling is a choice by P* from the same pool; the targets are a share of the path, and the values of the base and of the ceiling are read on the first run of the base and pinned;
  - the tolerance is the larger of the author's tolerance and 4.3 standard errors of the paired difference (R163 standalone);
  - the gate's bands are two-sided, and writing them anew is a pull request to `sim` with a reason and the criterion's report, read at review;
  - on a graph whose package is not published, the checks are "no graph"; a new package is a new pin and new bands, and earlier choices are not revisited by themselves;
  - T02 is one task.
- **A proposal of the execution, not the author's answer:**
  - true mastery is the P* of the reference task, of difficulty 3 at the band, without a miss, with five options and at R* = 1, no lower than 0.775, as standalone's "truly mastered"; the truth's command of the anchor skill is with retention at the time of reading;
  - the invariants `I1`–`I4` — reproducibility, a reason code from `ReasonCode` and true to its meaning, the author's order without a verdict, a task whose solution went to the device not as the next choice;
  - the measures `M1`–`M10` and `E1`–`E5` with their definitions — the thresholds of pauses of 7 and 14 days, the window of being stuck of 50 choices, "at most two pieces of evidence" for pages, 0.05 for the catching up of forgetting;
  - the resolution is 4.3 standard errors of the paired difference of the candidate and the base, in the same run;
  - the direction: in the criterion every measure goes so that less is better — a signed measure by its absolute value, a share where more is better by its shortfall; a path where the ceiling is no better than the base leaves the score, and for the measures where the ceiling need not be better — coverage, repetitions, false "mastered" and "knows" — the path goes to perfection;
  - the author's tolerances — 0.01 logit for errors and lag, 1 point for the corridor and shares, 2 points for false "mastered" and "knows", one answer for the counts of answers, and the others in section 3.6;
  - the model's floor — false "mastered" at most 20 % on `G0` and `G0-skills1`, the answers before "mastered" at most 1.5 times the base's; standalone's targets — 0.55 of the path for lag and 0.27 for the corridor on `G2-half`, at most 25 % not caught up on `G3` — only for a candidate of the model that names them as its target; the engine has no number for the target; a candidate names its score group, and it too is better than the base;
  - the groups and weights of the score — the model's first four groups as in standalone's step, "Mastered" as in its choice of "mastered", the new groups "Forgets", "Pages" and "Anchor", and the engine's groups "Corridor", "Repetitions" and "Coverage"; a population's weight is a share of the group's weight;
  - the populations and the gate's numbers; the sizes — 4,000 learners a cell for a decisive run, 300–1,000 for a rough look, 300 for the gate;
  - the seeds: the working ones — 20261107, the experiment `working`; the held-out ones — their own for each choice of each criterion, `held-out-<criterion>-<n>`: the model's first choice — 20261108, `held-out-model-1`, the engine's first choice — 20261109, `held-out-engine-1`; the seed of each next choice and of each pair — by an amendment to this entry before the first run of its candidate.

— *Why:*
- **Two criteria.** The model and the engine choose different things: the model an estimate, the engine what to give. One score for both would mix the estimate's error with the corridor, and an improvement of one would hide a worsening of the other. And five criteria, one for each part, as the step, "mastered" and the "debt" in standalone, would make every new kind of candidate wait for its own entry, while the model's first candidates are several of its numbers at once (Д5–Д7 in table 4 of RUN.md `learner-model`). So the model's criterion covers all the parts of the model by score groups, and a candidate names its own.
- **The base and the ceiling.** The criterion is written down before the first candidate, not before the first line of code (R12 learner-model): what to compare with must be given by the library itself, not by a copy in `sim`. The ceiling gives a scale: the share of the path does not depend on the units of the measure and carries standalone's targets over as they are — the lag at most 45 % of the service's and c* = 0.27 are shares of the path to the oracle. The values of the base and of the ceiling are calibration, read before the candidate, as c* ([architecture.md](architecture.md), section 4).
- **The resolution is by the pair "candidate and base".** In standalone the resolution was read on a pair of rules already standing on the simulator, "never from a candidate": a slow constant step beside the service. `sim` has no such close pair: there are no other versions of the libraries on the simulator, the constants are constants of the version, and changing them is a new version (architecture 3.1 learner-model). The standard error of the paired difference is the noise of the comparison on these learners, not the level of the candidate's numbers: with it a candidate equal to the base passes 99 times in 100, and no wider.
- **The direction of measures.** A signed measure cancels itself out in the mean — in `G1` the learners stand on both sides of the band —, and a share where more is better would give a path with a minus sign. The absolute value and the shortfall give all the measures one direction, as `|r6_lag|` and the shortfall of the corridor in standalone.
- **Invariants apart.** Reproducibility and the reason code are required by concepts 12 spec at every choice, and the author's order at every skill with a rubric (concepts 4 spec). This is not statistics: one violation is enough, and an interval would mean nothing here.
- **The engine by coverage, not by learnedness.** The populations' truth grows from the choice only by which skills the learner answers, by the stability τ* that its answers move, and by the lessons of `L-explain` (R03). So the engine's criterion measures the corridor, repetitions and coverage — what the choice moves the truth by.
- **Two-sided bands.** Bands on learners who are always the same are a place for a change, not for chance: code that shifted nothing reads the recorded numbers to the last digit. An improvement that passes by itself shifts the starting point silently; writing anew at review shows what shifted and why.
- **Published graphs.** R03 does not allow judging candidates on synthetics, and waiting for all three packages would mean that `learner-model` and `engine` stand still until the releases of school maths and the essays.
- **Held-out learners of its own for each choice** — as in standalone, where the "debt" had its own set: a set that has been looked at is no longer held out.

— *Rejected:*
- one criterion for both libraries; a criterion for each part;
- a base without a ceiling — the targets would have to be written as numbers without knowing the base's numbers; a copy of standalone's model in `sim` as the base — R12 learner-model rejected it;
- only the author's tolerance — an equal candidate would fail by chance; only the simulator's resolution — on a large run it would catch shifts that the author does not care about;
- a resolution by the pair "base and ceiling" — their difference is wider than a close candidate's, and the tolerance would let a real worsening through; and where there is no ceiling — for essays and being stuck — there would be no resolution at all;
- standalone's c* = 0.27 as the target of the engine's corridor: it is a share of the path to the model's oracle in the choice of the service, and the engine's path goes to a choice by P*;
- one count of held-out sets for both criteria, or a shared first set: the engine's first choice would take a set that the model's choice has already looked at;
- standalone's targets as a constraint for every candidate: even the base does not hold them — the choice of the step in standalone came to an exit —, and a candidate of the constants of forgetting would fail on someone else's target;
- one-sided bands; a paired check "no worse" on every pull request — two builds per run, and that would be for T03 to decide;
- wait for all three packages; synthetic packages until the releases — this contradicts R03;
- a model without forgetting as a second model on `F-none`: the library does not give one, and the simulator compares P with retention and without it on one estimate.

— *Consequence:* the code epic builds the measures, the criteria, the gate and the runs by sections 3.5–3.8 and before the first candidate reads the values of the base and of the ceiling — Д4 in table 4 of RUN.md; they are pinned by an amendment to this entry. T03 decides where the bands lie, the runner, the gate's time, the order of releases and what the gate includes in the CI of `modules`. — RUN.md, T02, table 4; [architecture.md](architecture.md), the status, sections 2.2, 3, 3.2, 3.4–3.8 and 4.

— *Amendment 2026-11-07:* in the CI of `modules` the gate holds not the bands but the invariants `I1`–`I4` and the absolute thresholds of the first slice, on the package that the pull request built, and it runs when what the simulator reads from a package has changed — R05. The bands of section 3.7 the gate holds on changes of `learner-model`, `engine` and `sim` itself, on the published packages, and a new release of `modules` is still a new pin and new bands. The rest of this entry stands. — [architecture.md](architecture.md), sections 3.7 and 3.9.

— *Amendment 2026-11-07:* the base — the first released version of each library — no longer rests on the single row В24 `learner-model`: `spec` recorded it by R103 spec — the first version of each library, released before the first tag of `sim`, and the gate holds every change after it. — R06.

— *Amendment 2026-11-07:* the engine chooses a lesson or a task as a whole, not an item (R104 spec). So the engine's ceiling is the task with the largest share of items whose P* is in the corridor, and on a tie, the one whose mean P* of the items is closer to 0.775, not the item with P* closer to 0.775; the ceiling of the pool for `E1` is the mean over issues of the largest share of items in the corridor of one reachable task. The invariant `I2` also reads the flags of the trial series from `Read` and the summary of issues: "consolidation" and "probe" are not in the module's trial series, "the queue" is the skill that has gone unissued for longest. — R06; [architecture.md](architecture.md), sections 3.5 and 3.6.

— *Amendment 2026-11-08:* the outcome of detailing of `engine` (R107 spec; R03–R10 engine) refines what the simulator checks and gives the engine. (1) `I2`: `ReasonCode` has a ninth value, `REASON_CODE_ROTATION`, "in rotation" — a skill that was issued before and has gone unissued for longest, but does not fit "the queue": it is mastered, or there are no reachable skills; "new skill" — the rule counts the prerequisites as met, and that includes a prerequisite that nothing in the input teaches (the 2026-11-07 amendment to R03 engine); "the queue" — the longest in its own module ([contracts text](../spec/contracts/README.md), section 15). (2) `I4`: "at once" of R82 spec is "not as the next task", and 24 hours with the second pass is the pause of `engine`, not an invariant. (3) The engine's ceiling and the ceiling of the pool of `E1` take tasks with the same exclusions by R82 spec as the rule — otherwise the path to the ceiling would be wider than the rule can go. (4) The summary of issues comes with the additions of R107 spec, and the simulator selects the pool and the lessons, as `api` does, by the age category — one per learner, not per module. — R07; [architecture.md](architecture.md), sections 1, 3, 3.1, 3.4, 3.5 and 3.6.

— *Amendment 2026-11-08:* `I2` checks the meaning of the code by the facts that it names, as necessary conditions on one of the skills of the choice — whether the skill was issued, whether it is mastered, whether its date has come, whether the module's trial series has ended, a failure in the last attempt, a rubric, an open run, a skill from the map —, while the order "for longest", reachability, "on pause" and "mastered for the queue" are the rule of `engine`, and its golden vectors hold them. "New skill" — the skill has not been issued: at "the bottom of the ladder" `engine` gives it even with unmet prerequisites (architecture 3.5 engine). Recomputing the rule in the simulator would mean measuring the engine on itself, like the truth on the model's code (R03). — R17; [architecture.md](architecture.md), section 3.5.

— *Amendment 2026-11-09:* the values of the base and of the ceiling, from which the path of the targets is made up, are read on the first decisive run of the base and pinned. How the code holds them — R28. Д4 of table 4 of RUN.md is closed.
- **Runs.**
  - The simulator — commit 6041e1d of `sim` and the recipe `just decisive` (RUN.md, T26).
  - The base — `learner-model` v0.1.0 (d2620a1) and `engine` v0.1.0 (c6c5bb9), the packages of `modules` v0.1.0.
  - The working learners — seed 20261107, 4,000 a cell: 41 populations, 65 cells.
  - The run of the base took 35 minutes on 20 cores of the development container, and its file is SHA-256 `bc6af579c88261c95816e70dd98d3f0d6af6f868095efe22de9a834eacff70d9`.
  - The run of the engine's ceiling — 9 minutes, SHA-256 `4fc3b92427353d57d7bca64cd12c1b34edf5dfab7e7f1d6af43bf5ce58311968`.
  - The files do not lie in the repository: the recipe gives them anew to the byte.
- **Values** — in the criterion's direction, less is better: the share in the corridor of `E1` by its shortfall from one, the R* of `E3` and the share of mastered skills of `E4` too. The ceiling is the oracle for the model's measures, the engine's ceiling for `E1` and zero for the measures "to perfection" ([architecture.md](architecture.md), section 3.6). The path is the base minus the ceiling; a dash is a measure that the run did not read.

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.8961 | 0.7541 | 0.1420 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.3913 | 0.0000 | 0.3913 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0094 | 0.0000 | 0.0094 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.3034 | 0.0000 | 4.3034 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8280 | 0.7510 | 0.0770 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2366 | 0.0000 | 0.2366 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0178 | 0.0000 | 0.0178 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.8116 | 0.0000 | 15.8116 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1702 | 0.0000 | 1.1702 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.4376 | 0.4376 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0513 | 0.0000 | 1.0513 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3708 | 0.3708 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.8812 | 0.7170 | 0.1641 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 1.2584 | 0.0000 | 1.2584 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8694 | 0.7001 | 0.1693 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4176 | 0.0000 | 2.4176 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.8852 | 0.7252 | 0.1600 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.1154 | 0.0000 | 0.1154 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.6722 | 0.0000 | 0.6722 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8113 | 0.6866 | 0.1247 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3595 | 0.0000 | 0.3595 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0914 | 0.0000 | 1.0914 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.8847 | 0.7210 | 0.1637 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.4966 | 0.0000 | 0.4966 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8152 | 0.6739 | 0.1413 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7609 | 0.0000 | 0.7609 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 53.1216 | 0.0000 | 53.1216 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 41.0204 | 0.0000 | 41.0204 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 42.4563 | 0.0000 | 42.4563 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.8831 | 0.0000 | 54.8831 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0196 | 0.0000 | 0.0196 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.3104 | 0.0000 | 4.3104 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0330 | 0.0000 | 0.0330 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 16.0446 | 0.0000 | 16.0446 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.6674 | 0.0000 | 1.6674 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1431 | 0.1431 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0083 | 0.0000 | 0.0083 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.7375 | 0.0000 | 0.7375 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1537 | 0.1502 | 0.0035 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 5.2860 | 0.0000 | 5.2860 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.7760 | 0.0000 | 0.7760 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1559 | 0.1555 | 0.0004 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 2.7931 | 0.0000 | 2.7931 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0104 | 0.0000 | 0.0104 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.8706 | 0.0000 | 0.8706 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1305 | 0.1131 | 0.0174 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 5.9726 | 0.0000 | 5.9726 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.7970 | 0.0000 | 0.7970 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1104 | 0.0906 | 0.0198 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 3.3137 | 0.0000 | 3.3137 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1334 | 0.1191 | 0.0142 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 6.3771 | 0.0000 | 6.3771 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1177 | 0.0999 | 0.0178 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 2.7257 | 0.0000 | 2.7257 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0116 | 0.0000 | 0.0116 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.8175 | 0.0000 | 0.8175 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1632 | 0.1553 | 0.0079 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 5.5010 | 0.0000 | 5.5010 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0003 | 0.0000 | 0.0003 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.9354 | 0.0000 | 0.9354 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1578 | 0.1564 | 0.0015 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 2.8993 | 0.0000 | 2.8993 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 94.8470 | 0.0000 | 94.8470 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 337.6990 | 0.0000 | 337.6990 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1739 | 0.1419 | 0.0320 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1975 | 0.1196 | 0.0779 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1555 | 0.1488 | 0.0067 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.2574 | 0.0000 | 0.2574 |

- **What follows from this.**
  - The model's targets are readable: on `G2-half` the path of the lag of `M3` is 0.6722 and 1.0914 logit, the path of the corridor of `E1` 0.1600 and 0.1247.
  - The share of hard first items `M7.hard` has a path of zero: the oracle sees the same chosen items as the base — the measure leaves the score of the group "Out of its band", and that group is judged by `M7` after 10 answers.
  - `M4` after a break and `M4` from one or two pages the run read in no cell: "mastered" was not declared there. These measures leave the score of their groups too.
  - The path of R* on the last day, `E3.end`, is from 0.0002 to 0.0116: the base loses almost no retention by the end, and the share of the path on so small a value is noisy.
- A new pinning run — only by a new amendment to this entry.

— *Amendment 2026-11-09, the author's answer to В7 when planning E3:* a measure whose pinned path is positive but no greater than the author's tolerance of that measure ([architecture.md](architecture.md), section 3.6) leaves the score the same way as a measure with a path no greater than zero: its share is not read, its weight goes with it, the report names it, and the target on it is not met. These measures still hold the constraints "no worse". On the paths of the amendment above these are `E3.end` on 4 cells of 6 — `F-model` on both maths, `F-exp` and `F-growth5` on school maths — and Brier of `M2` after pauses on 2 cells of 8: `F-model` and `F-growth5` on school maths. The table of paths above is not rewritten. The code — T29 RUN.md.
- *Why:* the author's tolerance is what the author has already called "not changed", and a ceiling indistinguishable from the base within the tolerance sets no path. For the measures whose path goes to perfection, as for `E3.end`, the path is the base's own shortfall, and a base no farther than the tolerance from perfection sets no path either. A difference of a few thousandths — within the noise on 4,000 learners — gives on such a path a share in units and decides the score of the whole group. The rule is one for all measures, and the group "Repetitions" stays with its lateness of `E3`.
- *Rejected:*
  - leave it as it is — noise would decide the score of the groups "Repetitions" and "Forgets";
  - remove `E3.end` from the group "Repetitions" — a fix of one measure: `M2` after pauses would stay with the same noise, and the next calibration may give a small path to another measure.

— *Amendment 2026-11-09, a new pin (T30; R31):* the values of the base and of the ceiling, from which the path of the targets is made up, are read anew, before the first candidate, on the decisive run of the base on the new pin: `modules` v0.1.4 with `olympiad-math` 0.1.1, and the populations `A-` on the pair of maths. They are pinned in `internal/criterion/paths.json` (R28; the 2026-11-09 amendment to R28). The judgement no longer holds the table of the first 2026-11-09 amendment: its values are of the graph of `olympiad-math` 0.1.0.
- **Runs.**
  - The simulator — the tree of the pull request of T30 and the recipe `just decisive` (RUN.md, T30). The run went on the tree before the fixes from the review; the fixes touch only the gate, and a small run of the engine and of the ceiling — 3 learners each of `G0`, `A-linked`, `F-model` and `E-base` — gives the same bytes on both trees.
  - The base — `learner-model` v0.1.0 (d2620a1) and `engine` v0.1.0 (c6c5bb9), the packages of `modules` v0.1.4.
  - The working learners — seed 20261107, 4,000 a cell: 41 populations, 65 cells.
  - The run of the base took 6 h 43 min on 20 cores of the development container, not 35 minutes as on 0.1.0. Its file is 2.1 GB, SHA-256 `a36dc5b222696690132c109d389c2b65800fe2329113b390023d4c7039bdd2e4`.
  - The run of the engine's ceiling — 1 h 31 min, SHA-256 `8a94b5530b85fcf119bb1da4462c88642acbd006eced82a72846abc603c9481e`.
  - Reproducibility was checked cheaply: one twentieth of the ceiling's run — the parts `k/400` at k = 0, 20, …, 380, on 20 cores in 4 minutes — was run anew and matched the file line for line, 13,000 lines. The full repeat of the ceiling was cancelled, so as not to spend another hour and a half on it.
  - The files do not lie in the repository: the recipe gives them anew to the byte.
- **Values** — in the criterion's direction, less is better, as in the first amendment:

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.7347 | 0.0163 | 0.7184 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.5151 | 0.0000 | 0.5151 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0300 | 0.0000 | 0.0300 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.6300 | 0.0000 | 4.6300 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8280 | 0.7510 | 0.0770 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2366 | 0.0000 | 0.2366 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0178 | 0.0000 | 0.0178 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.8116 | 0.0000 | 15.8116 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1797 | 0.0000 | 1.1797 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.3523 | 0.3523 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0513 | 0.0000 | 1.0513 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3708 | 0.3708 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.7313 | 0.0125 | 0.7189 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 0.9971 | 0.0000 | 0.9971 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8694 | 0.7001 | 0.1693 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4176 | 0.0000 | 2.4176 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.7210 | 0.0110 | 0.7100 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.4284 | 0.0000 | 0.4284 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.5019 | 0.0000 | 0.5019 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8113 | 0.6866 | 0.1247 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3595 | 0.0000 | 0.3595 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0914 | 0.0000 | 1.0914 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.7189 | 0.0092 | 0.7097 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.3883 | 0.0000 | 0.3883 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8152 | 0.6739 | 0.1413 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7609 | 0.0000 | 0.7609 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 55.1020 | 0.0000 | 55.1020 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 41.0204 | 0.0000 | 41.0204 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 52.0895 | 0.0000 | 52.0895 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.8831 | 0.0000 | 54.8831 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0592 | 0.0000 | 0.0592 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.6722 | 0.0000 | 4.6722 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0330 | 0.0000 | 0.0330 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 16.0446 | 0.0000 | 16.0446 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.4971 | 0.0000 | 1.4971 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1106 | 0.1106 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0446 | 0.0000 | 0.0446 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.4755 | 0.0000 | 0.4755 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1663 | 0.1612 | 0.0051 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | 0.0542 | 0.0000 | 0.0542 |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 3.5453 | 0.0000 | 3.5453 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.7760 | 0.0000 | 0.7760 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1559 | 0.1555 | 0.0004 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 2.7931 | 0.0000 | 2.7931 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0588 | 0.0000 | 0.0588 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.5361 | 0.0000 | 0.5361 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1569 | 0.1474 | 0.0095 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | 0.0476 | 0.0000 | 0.0476 |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 3.4966 | 0.0000 | 3.4966 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.7970 | 0.0000 | 0.7970 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1104 | 0.0906 | 0.0198 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 3.3137 | 0.0000 | 3.3137 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1606 | 0.1510 | 0.0096 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | 0.0483 | 0.0000 | 0.0483 |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 3.5636 | 0.0000 | 3.5636 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1177 | 0.0999 | 0.0178 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 2.7257 | 0.0000 | 2.7257 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0918 | 0.0000 | 0.0918 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.6435 | 0.0000 | 0.6435 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1943 | 0.1826 | 0.0117 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | 0.0449 | 0.0000 | 0.0449 |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 3.4742 | 0.0000 | 3.4742 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0003 | 0.0000 | 0.0003 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.9354 | 0.0000 | 0.9354 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1578 | 0.1564 | 0.0015 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 2.8993 | 0.0000 | 2.8993 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 172.7969 | 0.0000 | 172.7969 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 337.6990 | 0.0000 | 337.6990 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1739 | 0.1419 | 0.0320 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1975 | 0.1196 | 0.0779 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1555 | 0.1488 | 0.0067 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0039 | 0.0000 | 0.0039 |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.7703 | 0.0000 | 0.7703 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0634 | 0.0000 | 0.0634 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.1968 | 0.0000 | 0.1968 |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.5284 | 0.0000 | 0.5284 |

- **What is seen.**
  - On `olympiad-math` 0.1.1 the ceiling holds almost all the issued items in the corridor — a shortfall of 0.009–0.016 —, and the base only 26–28 %. The path of the corridor of `E1` there is about 0.71, and on 0.1.0 it was 0.14–0.16. The model's targets on `G2-half` are readable: the path of the lag of `M3` is 0.5019 and 1.0914 logit, the path of the corridor of `E1` 0.7100 and 0.1247.
  - `M4` after a break is now read on olympiad maths; on school maths there is still none.
  - `M9` is read for `A-unpractised` and `A-fading`, on the pair of maths.
  - A path no greater than the author's tolerance — the measure leaves the score (the 2026-11-09 amendment to this entry, the answer to В7) — for seven measures: `E3.end` on school maths for `F-model`, `F-exp` and `F-growth5` (0.0002–0.0003 at a tolerance of 0.01); Brier of `M2` after pauses on school maths for `F-model` and `F-growth5` (0.0004 and 0.0015 at 0.002); `M4` on `G0` of school maths (0.0178 at 0.02); `M9.hasFalse` for `A-unpractised` (0.0039 at 0.02). `E3.end` on olympiad maths is 0.045–0.092, and it stays in the score.
  - The share of hard first items `M7.hard` still has a path of zero.
- A new pinning run comes only by a new amendment to this entry.

— *Amendment 2026-11-10:* the paths of the amendment "a new pin" stay in force after the move to `spec` v0.12.0 too — the packages of `modules` v0.2.0 in format 0.4 and the libraries v0.1.4 (R32). The first 100 learners of each of the 65 cells on the new simulator give the same bytes as the run of the base of T30: 6,500 lines. T36 writes the paths anew, by a fast run (RUN.md).

— *Amendment 2026-11-10, runs no longer than an hour (T36; R33):* the decisive run and the confirmation go on 1,000 learners a cell, not on 4,000, and the choice repeated for `I1` takes only the first 10 learners of a cell. The resolution of "no worse" is still 4.3 standard errors of the paired difference in the same run, but on four times fewer learners it is 2 times wider. Where it is larger than the author's tolerance, it becomes the tolerance. A rough look is a run of fewer than 1,000 learners. Returning to 4,000 comes by a new amendment to this entry, when `engine` speeds up the choice (row В10 of table 3 of RUN.md). The values of the base and of the ceiling that the path of the targets is made of are read anew on a decisive run of this size. They are pinned in `internal/criterion/paths.json` (R28) in place of the values of the amendment "a new pin".
- **Runs.**
  - The simulator is the tree of the pull request of T36 and the recipe `just decisive` (RUN.md, T36).
  - The base is `learner-model` v0.1.4 (acfc4f2) and `engine` v0.1.4 (ddfed3a), the packages of `modules` v0.2.0.
  - The working learners are seed 20261107, 1,000 a cell: 41 populations, 65 cells. They are the first 1,000 of the same learners as in the run of T30.
  - Other containers loaded the host machine the whole time: load average 85–150 on 20 threads, and each of the 20 processes got from a third to a half of a CPU.
  - The run of the base took 1 h 42 min by the clock and about 37 CPU minutes per process, 748 minutes for all 20. On a free machine this is less than an hour. Its file is 525 MB, SHA-256 `b13a9264353936391f84e52c578805d460461354b4ca0017c22763dc9f9d4c92`.
  - The run of the engine's ceiling took 41 min by the clock and about 14 CPU minutes per process. Its file is 406 MB, SHA-256 `f579ecdf54bd11493f84da1fca5e687aba12ea0c9f442b0dbaecf8ddc1b4cc7d`.
  - The run of the ceiling was repeated whole — 46 min on the same busy machine — and matched the file to the byte: the same SHA-256.
  - The files lie in `runs/` of the working copy and are not in the repository: the recipe gives them anew to the byte.
- **Values** — in the direction of the criterion, less is better, as in the first amendment:

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.7331 | 0.0167 | 0.7164 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.5103 | 0.0000 | 0.5103 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0287 | 0.0000 | 0.0287 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.6323 | 0.0000 | 4.6323 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8200 | 0.7424 | 0.0777 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2397 | 0.0000 | 0.2397 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0135 | 0.0000 | 0.0135 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.8219 | 0.0000 | 15.8219 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1989 | 0.0000 | 1.1989 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.3500 | 0.3500 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0524 | 0.0000 | 1.0524 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3738 | 0.3738 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.7329 | 0.0131 | 0.7198 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 0.9860 | 0.0000 | 0.9860 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8700 | 0.7025 | 0.1676 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4492 | 0.0000 | 2.4492 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.7229 | 0.0119 | 0.7110 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.4289 | 0.0000 | 0.4289 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.4974 | 0.0000 | 0.4974 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8126 | 0.6884 | 0.1242 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3528 | 0.0000 | 0.3528 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0892 | 0.0000 | 1.0892 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.7199 | 0.0096 | 0.7103 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.4007 | 0.0000 | 0.4007 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8128 | 0.6726 | 0.1402 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7609 | 0.0000 | 0.7609 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 53.3750 | 0.0000 | 53.3750 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 40.8951 | 0.0000 | 40.8951 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 52.4346 | 0.0000 | 52.4346 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.2998 | 0.0000 | 54.2998 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0590 | 0.0000 | 0.0590 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.6284 | 0.0000 | 4.6284 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0360 | 0.0000 | 0.0360 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 15.9533 | 0.0000 | 15.9533 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.5312 | 0.0000 | 1.5312 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1089 | 0.1089 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0471 | 0.0000 | 0.0471 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.4833 | 0.0000 | 0.4833 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1675 | 0.1621 | 0.0054 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | 0.0602 | 0.0000 | 0.0602 |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 3.6381 | 0.0000 | 3.6381 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.6963 | 0.0000 | 0.6963 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1576 | 0.1572 | 0.0003 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 4.0788 | 0.0000 | 4.0788 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0578 | 0.0000 | 0.0578 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.5363 | 0.0000 | 0.5363 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1568 | 0.1470 | 0.0098 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | 0.0534 | 0.0000 | 0.0534 |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 3.6715 | 0.0000 | 3.6715 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.8342 | 0.0000 | 0.8342 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1098 | 0.0906 | 0.0193 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 3.8074 | 0.0000 | 3.8074 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1587 | 0.1487 | 0.0100 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | 0.0538 | 0.0000 | 0.0538 |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 3.7889 | 0.0000 | 3.7889 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1170 | 0.1001 | 0.0169 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 2.8574 | 0.0000 | 2.8574 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0890 | 0.0000 | 0.0890 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.6450 | 0.0000 | 0.6450 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1946 | 0.1831 | 0.0115 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | 0.0463 | 0.0000 | 0.0463 |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 3.4889 | 0.0000 | 3.4889 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.9128 | 0.0000 | 0.9128 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1548 | 0.1537 | 0.0012 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 3.0377 | 0.0000 | 3.0377 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 170.0686 | 0.0000 | 170.0686 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 336.6550 | 0.0000 | 336.6550 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1730 | 0.1411 | 0.0319 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1983 | 0.1194 | 0.0788 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1549 | 0.1484 | 0.0066 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0036 | 0.0000 | 0.0036 |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.7754 | 0.0000 | 0.7754 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0559 | 0.0000 | 0.0559 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.1978 | 0.0000 | 0.1978 |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.5087 | 0.0000 | 0.5087 |

- **What is seen.**
  - The paths are close to the paths of the amendment "a new pin": for the paths of the corridor `E1` the difference is at most 0.002, and for the lag `M3` at most 0.032. The count of answers until the forgotten is caught, `M5.caught`, moved the most, on school maths — by up to 1.3 answers for `F-model`: on 1,000 learners it is noisy.
  - The path is no larger than the author's tolerance for the same seven measures as before. A measure with such a path leaves the score (the 2026-11-09 amendment to this entry, the answer to В7). They are:
    - `E3.end` on school maths for `F-model`, `F-exp` and `F-growth5`;
    - Brier of `M2` after pauses on school maths for `F-model` and `F-growth5` (0.0003 and 0.0012 at a tolerance of 0.002);
    - `M4` on `G0` of school maths (0.0135 at 0.02);
    - `M9.hasFalse` for `A-unpractised` (0.0036 at 0.02).
  - The share of hard first items, `M7.hard`, still has a path of zero.
- A new pinning run comes only by a new amendment to this entry.

— *Amendment 2026-11-11, the base on `engine` v0.1.6 and 1,500 learners (T41; R31, R33, R38):* the base is `engine` v0.1.6, and the decisive run and the confirmation go on 1,500 learners a cell, not on 1,000. `engine` v0.1.6 is v0.1.5, the fix of R82 spec, plus a speed-up of the choice without a change of the choice (R20 engine). The gate on it reads all 69 numbers the same to the bit as on v0.1.4. The choice became cheaper, and the base on 1,500 learners now fits into 45 minutes on a free machine (the 2026-11-11 amendment to R33). The resolution of "no worse" is still 4.3 standard errors of the paired difference in the same run, but on 1,500 learners it is √1.5 ≈ 1.22 times narrower than on 1,000. A rough look is a run of fewer than 1,500 learners. The values of the base and of the ceiling that the path of the targets is made of are read anew on a decisive run of this size. They are pinned in `internal/criterion/paths.json` (R28) in place of the values of the amendment "runs no longer than an hour".
- **Runs.**
  - The simulator is the tree of the pull request of T41 and the recipe `just decisive` (RUN.md, T41).
  - The base is `learner-model` v0.1.4 (acfc4f2) and `engine` v0.1.6 (57e0ab8), the packages of `modules` v0.2.0. `release-verify` accepted the archive of `engine` v0.1.6.
  - The working learners are seed 20261107, 1,500 a cell: 41 populations, 65 cells. They are the same working learners as in the runs of T30 and T36: the first 1,000 of them are the learners of the run of T36, and all 1,500 are the first 1,500 of the 4,000 of the run of T30.
  - The run of the base took 38 min by the clock and about 33 CPU minutes per process, 655 minutes for all 20; the load average at the start was 5. Its file is 788 MB, SHA-256 `c09442e774dc1f9403e5d67fda72e881005241722ed4e464d4ac1f3603b32b4f`.
  - The run of the engine's ceiling took about 19 CPU minutes per process, 383 minutes for all 20. By the clock it took 1 h 30 min: the processes got less CPU than they asked for — going over an hour speaks of the machine, not of the run. Its file is 609 MB, SHA-256 `8f5382340ed611ddfc285ecab66c429dd349f0e4c3f471ca1c1fc79436c58469`.
  - The run of the ceiling was repeated whole — 19 min 23 s by the clock — and matched the file to the byte: the same SHA-256.
  - The files lie in `runs/` of the working copy and are not in the repository: the recipe gives them anew to the byte. The files of the runs on 1,000 learners, which the judgements `engine-1` and `model-1` refer to, lie there too, as `base-1000.jsonl` and `ceiling-1000.jsonl`.
- **Values** — in the direction of the criterion, less is better, as in the first amendment:

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.7323 | 0.0168 | 0.7155 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.5119 | 0.0000 | 0.5119 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0288 | 0.0000 | 0.0288 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.6196 | 0.0000 | 4.6196 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8240 | 0.7481 | 0.0760 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2302 | 0.0000 | 0.2302 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0104 | 0.0000 | 0.0104 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.7789 | 0.0000 | 15.7789 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1923 | 0.0000 | 1.1923 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.3462 | 0.3462 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0481 | 0.0000 | 1.0481 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3755 | 0.3755 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.7321 | 0.0123 | 0.7198 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 0.9969 | 0.0000 | 0.9969 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8698 | 0.7013 | 0.1685 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4447 | 0.0000 | 2.4447 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.7241 | 0.0114 | 0.7126 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.4315 | 0.0000 | 0.4315 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.5027 | 0.0000 | 0.5027 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8112 | 0.6882 | 0.1230 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3647 | 0.0000 | 0.3647 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0970 | 0.0000 | 1.0970 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.7181 | 0.0084 | 0.7096 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.3994 | 0.0000 | 0.3994 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8138 | 0.6749 | 0.1389 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7497 | 0.0000 | 0.7497 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 55.9157 | 0.0000 | 55.9157 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 41.2151 | 0.0000 | 41.2151 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 52.6909 | 0.0000 | 52.6909 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.3511 | 0.0000 | 54.3511 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0595 | 0.0000 | 0.0595 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.6632 | 0.0000 | 4.6632 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0296 | 0.0000 | 0.0296 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 16.0122 | 0.0000 | 16.0122 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.5272 | 0.0000 | 1.5272 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1103 | 0.1103 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0463 | 0.0000 | 0.0463 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.4779 | 0.0000 | 0.4779 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1677 | 0.1627 | 0.0051 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | 0.0567 | 0.0000 | 0.0567 |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 3.5371 | 0.0000 | 3.5371 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.7332 | 0.0000 | 0.7332 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1549 | 0.1545 | 0.0004 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 3.7969 | 0.0000 | 3.7969 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0564 | 0.0000 | 0.0564 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.5484 | 0.0000 | 0.5484 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1571 | 0.1475 | 0.0096 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | 0.0528 | 0.0000 | 0.0528 |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 3.6069 | 0.0000 | 3.6069 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.8469 | 0.0000 | 0.8469 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1106 | 0.0906 | 0.0200 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 2.9475 | 0.0000 | 2.9475 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1597 | 0.1502 | 0.0095 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | 0.0524 | 0.0000 | 0.0524 |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 3.6834 | 0.0000 | 3.6834 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1193 | 0.1021 | 0.0172 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 2.7043 | 0.0000 | 2.7043 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0915 | 0.0000 | 0.0915 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.6472 | 0.0000 | 0.6472 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1941 | 0.1827 | 0.0114 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | 0.0470 | 0.0000 | 0.0470 |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 3.5538 | 0.0000 | 3.5538 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0003 | 0.0000 | 0.0003 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.8907 | 0.0000 | 0.8907 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1561 | 0.1547 | 0.0014 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 2.4756 | 0.0000 | 2.4756 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 173.0952 | 0.0000 | 173.0952 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 338.1067 | 0.0000 | 338.1067 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1733 | 0.1415 | 0.0318 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1977 | 0.1195 | 0.0783 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1548 | 0.1482 | 0.0066 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0033 | 0.0000 | 0.0033 |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.7777 | 0.0000 | 0.7777 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0625 | 0.0000 | 0.0625 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.1999 | 0.0000 | 0.1999 |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.5211 | 0.0000 | 0.5211 |

- **What is seen.**
  - The paths are close to the paths of the amendment "runs no longer than an hour": for the paths of the corridor `E1` the difference is at most 0.002, for the lag `M3.lag` at most 0.012. Again the count of answers until the forgotten is caught, `M5.caught`, moved the most, on school maths — by up to 0.86 answers for `F-exp` —, and so did the share of false "mastered" `M4` on `G0` and `G0-skills1` of school maths — by 0.003 and 0.006: on the small pools of school maths these measures are noisy.
  - The path is no larger than the author's tolerance for the same seven measures as before, and a measure with such a path still leaves the score (the 2026-11-09 amendment to this entry, the answer to В7):
    - `E3.end` on school maths for `F-model`, `F-exp` and `F-growth5`;
    - Brier of `M2` after pauses on school maths for `F-model` and `F-growth5` (0.0004 and 0.0014 at a tolerance of 0.002);
    - `M4` on `G0` of school maths (0.0104 at 0.02);
    - `M9.hasFalse` for `A-unpractised` (0.0033 at 0.02).
  - The share of hard first items, `M7.hard`, still has a path of zero.
- A new pinning run comes only by a new amendment to this entry.

— *Amendment 2026-11-11, packages with `time` and leading questions (T43; R41, R42):* the paths are read anew on a new pin — `modules` v0.2.5 with `olympiad-math` and `grade-3-math` 0.3.0 — and on `learner-model` v0.1.5. The rule, the size and the resolution are the same as in the amendment "the base on `engine` v0.1.6 and 1,500 learners": 1,500 working learners a cell, 4.3 standard errors of the paired difference. The paths are pinned in `internal/criterion/paths.json` (R28) in place of the values of that amendment.
- **Runs.**
  - The simulator is the tree of the branch of epic E4 with T42 and T43 and the recipe `just decisive` (RUN.md, T43).
  - The base is `learner-model` v0.1.5 (3a2cd29) and `engine` v0.1.6 (57e0ab8), the packages of `modules` v0.2.5: `olympiad-math` and `grade-3-math` 0.3.0, `english-essays` 0.2.0.
  - The working learners are the same as in the previous amendment: seed 20261107, 1,500 a cell, 41 populations, 65 cells.
  - The run of the base took 47 min 48 s by the clock and about 38.7 CPU minutes per process, 774 minutes for all 20. Neighbouring containers shared the host: the 15-minute load average towards the end was 24.4 with 20 processes of its own. Its file is 780 MB, SHA-256 `4567a53214a4fca607cb0d35ba3d4e4969b6213a1f7a2d3206a248d793e61fd4`.
  - The run of the engine's ceiling took 34 min 45 s by the clock, about 18.8 CPU minutes per process, 377 minutes for all 20, at a load average of about 37–42. Its file is 610 MB, SHA-256 `ea7b9eecf621b920a500b821db07f391dd07425b1fcc2efd5efde0541cfd984b`.
  - The run of the ceiling was repeated whole — 61 min 21 s by the clock at a load average of up to 72 — and matched the file to the byte: the same SHA-256.
  - The files lie in `runs/` of the working copy and are not in the repository. The previous files on 1,500 learners are overwritten: only the previous `paths.json` referred to them.
- **Values** — in the direction of the criterion, less is better, as in the first amendment:

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.7438 | 0.0193 | 0.7245 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.5235 | 0.0000 | 0.5235 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0291 | 0.0000 | 0.0291 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.6284 | 0.0000 | 4.6284 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8240 | 0.7481 | 0.0760 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2300 | 0.0000 | 0.2300 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0104 | 0.0000 | 0.0104 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.7789 | 0.0000 | 15.7789 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1922 | 0.0000 | 1.1922 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.3465 | 0.3465 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0481 | 0.0000 | 1.0481 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3755 | 0.3755 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.7424 | 0.0130 | 0.7294 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 1.0088 | 0.0000 | 1.0088 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8698 | 0.7013 | 0.1685 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4449 | 0.0000 | 2.4449 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.7341 | 0.0124 | 0.7217 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.4307 | 0.0000 | 0.4307 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.5201 | 0.0000 | 0.5201 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8112 | 0.6882 | 0.1230 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3647 | 0.0000 | 0.3647 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0970 | 0.0000 | 1.0970 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.7284 | 0.0094 | 0.7189 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.4180 | 0.0000 | 0.4180 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8138 | 0.6749 | 0.1389 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7502 | 0.0000 | 0.7502 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 60.1562 | 0.0000 | 60.1562 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 41.2151 | 0.0000 | 41.2151 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 47.6854 | 0.0000 | 47.6854 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.3511 | 0.0000 | 54.3511 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0552 | 0.0000 | 0.0552 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.6811 | 0.0000 | 4.6811 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0296 | 0.0000 | 0.0296 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 16.0122 | 0.0000 | 16.0122 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.5272 | 0.0000 | 1.5272 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1103 | 0.1103 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0459 | 0.0000 | 0.0459 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.4752 | 0.0000 | 0.4752 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1658 | 0.1604 | 0.0054 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | 0.0416 | 0.0000 | 0.0416 |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 3.5640 | 0.0000 | 3.5640 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.7070 | 0.0000 | 0.7070 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1549 | 0.1545 | 0.0003 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 4.5364 | 0.0000 | 4.5364 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0580 | 0.0000 | 0.0580 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.5485 | 0.0000 | 0.5485 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1567 | 0.1474 | 0.0093 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | 0.0407 | 0.0000 | 0.0407 |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 3.7583 | 0.0000 | 3.7583 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.8464 | 0.0000 | 0.8464 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1106 | 0.0906 | 0.0200 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 2.9475 | 0.0000 | 2.9475 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1591 | 0.1498 | 0.0093 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | 0.0323 | 0.0000 | 0.0323 |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 3.6401 | 0.0000 | 3.6401 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1195 | 0.1024 | 0.0171 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 1.8584 | 0.0000 | 1.8584 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0908 | 0.0000 | 0.0908 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.6378 | 0.0000 | 0.6378 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1926 | 0.1815 | 0.0110 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | 0.0332 | 0.0000 | 0.0332 |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 3.6300 | 0.0000 | 3.6300 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0003 | 0.0000 | 0.0003 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.9054 | 0.0000 | 0.9054 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1559 | 0.1545 | 0.0014 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 3.0654 | 0.0000 | 3.0654 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 158.6395 | 0.0000 | 158.6395 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 338.0980 | 0.0000 | 338.0980 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1733 | 0.1415 | 0.0318 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1977 | 0.1195 | 0.0783 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1548 | 0.1482 | 0.0066 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0035 | 0.0000 | 0.0035 |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.7765 | 0.0000 | 0.7765 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0544 | 0.0000 | 0.0544 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.1995 | 0.0000 | 0.1995 |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.5216 | 0.0000 | 0.5216 |

- **What is seen.**
  - The paths are close to the paths of the previous amendment: for the corridor `E1` the difference is at most 0.010, for the lag `M3.lag` at most 0.019, for `M1` at most 0.012. The paths of the measures counted in answers moved the most: the answers until the error of the skills holds after a jump — `M3.settle` on olympiad maths — by 14.5 answers for `F-break` (173.1 → 158.6), by 5.0 for `G3-drop` and by 4.2 for `G3`, and the count of answers until the forgotten is caught, `M5.caught`, on school maths — by up to 0.85 answers. Of the shares, the most moved were the share of late repetitions `E3.late` for `F-model` on school maths, by 0.026, and the false "mastered" after a break `M4.false.break` for `F-` on olympiad maths, by up to 0.020. Olympiad maths 0.3.0 is 30 answers as a time instead of a choice and 3 new templates; school maths is a worked example in each lesson.
  - The path is no larger than the author's tolerance for the same seven measures as before, and a measure with such a path still leaves the score (the 2026-11-09 amendment to this entry, the answer to В7):
    - `E3.end` on school maths for `F-model`, `F-exp` and `F-growth5`;
    - Brier of `M2` after pauses on school maths for `F-model` and `F-growth5` (0.0003 and 0.0014 at a tolerance of 0.002);
    - `M4` on `G0` of school maths (0.0104 at 0.02);
    - `M9.hasFalse` for `A-unpractised` (0.0035 at 0.02).
  - The share of hard first items, `M7.hard`, still has a path of zero.
- A new pinning run comes only by a new amendment to this entry.

— *Amendment 2026-11-11, the base on `learner-model` v0.1.6 (T45; R43):* the paths are read anew on `learner-model` v0.1.6 (7d796e2) with the same `engine` v0.1.6 (57e0ab8), the same packages of `modules` v0.2.5 and the same 1,500 working learners a cell: the judgement holds the base to the libraries of the paths (R40). The values are the same to the last digit as in the amendment "packages with `time` and leading questions", and its table stays the table of the paths.
- **Runs.**
  - The simulator is the tree of the branch of epic E4 with T42–T44 and the recipe `just decisive` (RUN.md, T45).
  - The result lines of both runs are the same bytes as those of the runs of T43; only the header, which names the version of `learner-model`, differs. So the ceiling was not run a second time for reproducibility: its run of T43 was repeated whole and matched to the byte, and the new one matched it.
  - The run of the base took 62 min 24 s by the clock at a load average of up to 46 on 20 cores and about 25.9 CPU minutes per process against 38.7 in T43. The run of the engine's ceiling took 22 min 28 s and about 18.0 CPU minutes per process against 18.8.
  - SHA-256: the base `9e40ee6065848078a8d7850f9b03fb8f5b91d2e122db5d58a98169efc630c6e0`, the ceiling `06fcc680110cba673933d2af7fa95a7274a552ad698c1f7510b8b5cc10e538a0`. The files lie in `runs/` of the working copy and are not in the repository.

— *Amendment 2026-11-12, 2,000 learners (T46; the 2026-11-12 amendment to R33):* the decisive run and the confirmation go on 2,000 learners a cell, not on 1,500, on the same base — `learner-model` v0.1.6 (7d796e2), `engine` v0.1.6 (57e0ab8) and the packages of `modules` v0.2.5. The speed-up of `Readout.Task` in `learner-model` v0.1.6 made the base cheaper, and 2,000 learners now fit into 45 minutes on a free machine with room to spare (the 2026-11-12 amendment to R33). The resolution of "no worse" is still 4.3 standard errors of the paired difference in the same run, but on 2,000 learners it is √(4/3) ≈ 1.15 times narrower than on 1,500. A rough look is a run of fewer than 2,000 learners. The values of the base and of the ceiling are read anew on a decisive run of this size and pinned in `internal/criterion/paths.json` (R28) in place of the values of the amendment "packages with `time` and leading questions".
- **Runs.**
  - The simulator is the tree of the pull request of T46 and the recipe `just decisive` (RUN.md, T46).
  - The working learners are seed 20261107, 2,000 a cell: 41 populations, 65 cells. The first 1,500 of them are the learners of the runs of T41–T45, and all 2,000 are the first 2,000 of the 4,000 of the run of T30.
  - The run of the base took 62 min 37 s by the clock at a load average from 49 at the start to 10 at the end — neighbouring containers shared the host — and about 31.3 CPU minutes per process, 626 minutes for all 20. On a free machine the sample promised about 38 minutes (the 2026-11-12 amendment to R33). Its file is 1,041 MB, SHA-256 `7db81ec7bdf8c02af8826031b671871e23338f812747f2926412cf5547ddc92e`.
  - The run of the engine's ceiling took 30 min 37 s by the clock at a load average from 10 to 31, about 24.1 CPU minutes per process, 483 minutes for all 20. Its file is 814 MB, SHA-256 `29741eb6ce20013ee66ccc8145693a62b8f7efa2ededecbf41494e74a5338021`.
  - The run of the ceiling was repeated whole — 29 min 54 s by the clock at a load average from 14 to 31 — and matched the file to the byte: the same SHA-256.
  - The files lie in `runs/` of the working copy and are not in the repository. The previous files on 1,500 learners lie there too, as `base-1500.jsonl` and `ceiling-1500.jsonl`.
- **Values** — in the direction of the criterion, less is better, as in the first amendment:

| Cell | Measure | Base | Ceiling | Path |
|---|---|---|---|---|
| `G0/novatrail/olympiad-math` | `E1.corridor` | 0.7448 | 0.0192 | 0.7256 |
| `G0/novatrail/olympiad-math` | `M1.rms.200` | 0.5216 | 0.0000 | 0.5216 |
| `G0/novatrail/olympiad-math` | `M4.false.all` | 0.0289 | 0.0000 | 0.0289 |
| `G0/novatrail/olympiad-math` | `M6.late` | 4.6559 | 0.0000 | 4.6559 |
| `G0/novatrail/grade-3-math` | `E1.corridor` | 0.8226 | 0.7471 | 0.0755 |
| `G0/novatrail/grade-3-math` | `M1.rms.200` | 0.2348 | 0.0000 | 0.2348 |
| `G0/novatrail/grade-3-math` | `M4.false.all` | 0.0076 | 0.0000 | 0.0076 |
| `G0/novatrail/grade-3-math` | `M6.late` | 15.8154 | 0.0000 | 15.8154 |
| `G1/novatrail/olympiad-math` | `M7.error.10` | 1.1852 | 0.0000 | 1.1852 |
| `G1/novatrail/olympiad-math` | `M7.hard` | 0.3496 | 0.3496 | 0.0000 |
| `G1/novatrail/grade-3-math` | `M7.error.10` | 1.0436 | 0.0000 | 1.0436 |
| `G1/novatrail/grade-3-math` | `M7.hard` | 0.3755 | 0.3755 | 0.0000 |
| `G2/novatrail/olympiad-math` | `E1.corridor` | 0.7402 | 0.0129 | 0.7273 |
| `G2/novatrail/olympiad-math` | `M3.lag` | 1.0026 | 0.0000 | 1.0026 |
| `G2/novatrail/grade-3-math` | `E1.corridor` | 0.8701 | 0.7003 | 0.1698 |
| `G2/novatrail/grade-3-math` | `M3.lag` | 2.4269 | 0.0000 | 2.4269 |
| `G2-half/novatrail/olympiad-math` | `E1.corridor` | 0.7332 | 0.0132 | 0.7200 |
| `G2-half/novatrail/olympiad-math` | `E4.mastered` | 0.4269 | 0.0000 | 0.4269 |
| `G2-half/novatrail/olympiad-math` | `M3.lag` | 0.5182 | 0.0000 | 0.5182 |
| `G2-half/novatrail/grade-3-math` | `E1.corridor` | 0.8118 | 0.6890 | 0.1228 |
| `G2-half/novatrail/grade-3-math` | `E4.mastered` | 0.3629 | 0.0000 | 0.3629 |
| `G2-half/novatrail/grade-3-math` | `M3.lag` | 1.0917 | 0.0000 | 1.0917 |
| `G2-fading/novatrail/olympiad-math` | `E1.corridor` | 0.7304 | 0.0107 | 0.7198 |
| `G2-fading/novatrail/olympiad-math` | `M3.lag` | 0.4237 | 0.0000 | 0.4237 |
| `G2-fading/novatrail/grade-3-math` | `E1.corridor` | 0.8148 | 0.6746 | 0.1401 |
| `G2-fading/novatrail/grade-3-math` | `M3.lag` | 0.7551 | 0.0000 | 0.7551 |
| `G3/novatrail/olympiad-math` | `M3.settle` | 57.5000 | 0.0000 | 57.5000 |
| `G3/novatrail/grade-3-math` | `M3.settle` | 41.2255 | 0.0000 | 41.2255 |
| `G3-drop/novatrail/olympiad-math` | `M3.settle` | 47.2078 | 0.0000 | 47.2078 |
| `G3-drop/novatrail/grade-3-math` | `M3.settle` | 54.4096 | 0.0000 | 54.4096 |
| `G0-skills1/novatrail/olympiad-math` | `M4.false.all` | 0.0554 | 0.0000 | 0.0554 |
| `G0-skills1/novatrail/olympiad-math` | `M6.late` | 4.7019 | 0.0000 | 4.7019 |
| `G0-skills1/novatrail/grade-3-math` | `M4.false.all` | 0.0287 | 0.0000 | 0.0287 |
| `G0-skills1/novatrail/grade-3-math` | `M6.late` | 16.0443 | 0.0000 | 16.0443 |
| `S-adult/novatrail/olympiad-math` | `M7.error.10` | 1.5051 | 0.0000 | 1.5051 |
| `S-adult/novatrail/olympiad-math` | `M7.hard` | 0.1110 | 0.1110 | 0.0000 |
| `F-model/novatrail/olympiad-math` | `E3.end` | 0.0452 | 0.0000 | 0.0452 |
| `F-model/novatrail/olympiad-math` | `E3.late` | 0.4787 | 0.0000 | 0.4787 |
| `F-model/novatrail/olympiad-math` | `M2.brier.after7` | 0.1650 | 0.1598 | 0.0053 |
| `F-model/novatrail/olympiad-math` | `M4.false.break` | 0.0423 | 0.0000 | 0.0423 |
| `F-model/novatrail/olympiad-math` | `M5.caught` | 3.6322 | 0.0000 | 3.6322 |
| `F-model/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-model/novatrail/grade-3-math` | `E3.late` | 0.7043 | 0.0000 | 0.7043 |
| `F-model/novatrail/grade-3-math` | `M2.brier.after7` | 0.1549 | 0.1546 | 0.0003 |
| `F-model/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-model/novatrail/grade-3-math` | `M5.caught` | 3.9287 | 0.0000 | 3.9287 |
| `F-exp/novatrail/olympiad-math` | `E3.end` | 0.0573 | 0.0000 | 0.0573 |
| `F-exp/novatrail/olympiad-math` | `E3.late` | 0.5357 | 0.0000 | 0.5357 |
| `F-exp/novatrail/olympiad-math` | `M2.brier.after7` | 0.1568 | 0.1479 | 0.0090 |
| `F-exp/novatrail/olympiad-math` | `M4.false.break` | 0.0382 | 0.0000 | 0.0382 |
| `F-exp/novatrail/olympiad-math` | `M5.caught` | 3.7098 | 0.0000 | 3.7098 |
| `F-exp/novatrail/grade-3-math` | `E3.end` | 0.0002 | 0.0000 | 0.0002 |
| `F-exp/novatrail/grade-3-math` | `E3.late` | 0.8039 | 0.0000 | 0.8039 |
| `F-exp/novatrail/grade-3-math` | `M2.brier.after7` | 0.1108 | 0.0907 | 0.0202 |
| `F-exp/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-exp/novatrail/grade-3-math` | `M5.caught` | 3.3407 | 0.0000 | 3.3407 |
| `F-C1/novatrail/olympiad-math` | `M2.brier.after7` | 0.1590 | 0.1498 | 0.0092 |
| `F-C1/novatrail/olympiad-math` | `M4.false.break` | 0.0308 | 0.0000 | 0.0308 |
| `F-C1/novatrail/olympiad-math` | `M5.caught` | 3.6668 | 0.0000 | 3.6668 |
| `F-C1/novatrail/grade-3-math` | `M2.brier.after7` | 0.1198 | 0.1026 | 0.0173 |
| `F-C1/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-C1/novatrail/grade-3-math` | `M5.caught` | 1.7587 | 0.0000 | 1.7587 |
| `F-growth5/novatrail/olympiad-math` | `E3.end` | 0.0913 | 0.0000 | 0.0913 |
| `F-growth5/novatrail/olympiad-math` | `E3.late` | 0.6415 | 0.0000 | 0.6415 |
| `F-growth5/novatrail/olympiad-math` | `M2.brier.after7` | 0.1926 | 0.1816 | 0.0110 |
| `F-growth5/novatrail/olympiad-math` | `M4.false.break` | 0.0355 | 0.0000 | 0.0355 |
| `F-growth5/novatrail/olympiad-math` | `M5.caught` | 3.6129 | 0.0000 | 3.6129 |
| `F-growth5/novatrail/grade-3-math` | `E3.end` | 0.0003 | 0.0000 | 0.0003 |
| `F-growth5/novatrail/grade-3-math` | `E3.late` | 0.9162 | 0.0000 | 0.9162 |
| `F-growth5/novatrail/grade-3-math` | `M2.brier.after7` | 0.1580 | 0.1565 | 0.0016 |
| `F-growth5/novatrail/grade-3-math` | `M4.false.break` | — | — | — |
| `F-growth5/novatrail/grade-3-math` | `M5.caught` | 2.9646 | 0.0000 | 2.9646 |
| `F-break/novatrail/olympiad-math` | `M3.settle` | 161.0811 | 0.0000 | 161.0811 |
| `F-break/novatrail/grade-3-math` | `M3.settle` | 336.3860 | 0.0000 | 336.3860 |
| `W-shock/novatrail/grade-3-math` | `M2.brier.pages` | 0.1735 | 0.1417 | 0.0317 |
| `W-shock/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-error/novatrail/grade-3-math` | `M2.brier.pages` | 0.1972 | 0.1195 | 0.0777 |
| `W-error/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `W-badday/novatrail/grade-3-math` | `M2.brier.pages` | 0.1550 | 0.1484 | 0.0066 |
| `W-badday/novatrail/grade-3-math` | `M4.false.pages` | — | — | — |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0036 | 0.0000 | 0.0036 |
| `A-unpractised/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.7762 | 0.0000 | 0.7762 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.hasFalse` | 0.0795 | 0.0000 | 0.0795 |
| `A-fading/novatrail/olympiad-math+novatrail/grade-3-math` | `M9.notFalse` | 0.1936 | 0.0000 | 0.1936 |
| `L-explain/novatrail/olympiad-math` | `E4.mastered` | 0.5218 | 0.0000 | 0.5218 |

- **What is seen.**
  - The paths are close to the paths of the amendment "packages with `time` and leading questions": for the corridor `E1` the difference is at most 0.003, for the lag `M3.lag` at most 0.018, for `M1` at most 0.005. The paths of the measures counted in answers moved the most: `M3.settle` for `G3` on olympiad maths — by 2.7 answers (60.2 → 57.5), for `F-break` on olympiad maths — by 2.4 (158.6 → 161.1), and `M5.caught` on school maths — by up to 0.61 answers for `F-model` (4.54 → 3.93). Of the shares, the most moved were `E3.late` for `F-exp` on school maths, by 0.0425 (0.846 → 0.804), and `M9.hasFalse` for `A-fading`, by 0.025 (0.054 → 0.080). The learners are the same, the libraries are the same, and 500 are added to the previous 1,500, so the shift comes from the new learners. A measure of rare events, as `M9.hasFalse` for `A-fading`, moves more than the others: on 1,500 its path was 0.054, on 2,000 it is 0.080: the path of such a measure is unsteady, and the judgement reads its share with this caveat.
  - The path is no larger than the author's tolerance for the same seven measures as before, and a measure with such a path still leaves the score (the 2026-11-09 amendment to this entry, the answer to В7):
    - `E3.end` on school maths for `F-model`, `F-exp` and `F-growth5`;
    - Brier of `M2` after pauses on school maths for `F-model` and `F-growth5` (0.0003 and 0.0016 at a tolerance of 0.002);
    - `M4` on `G0` of school maths (0.0076 at 0.02);
    - `M9.hasFalse` for `A-unpractised` (0.0036 at 0.02).
  - The share of hard first items, `M7.hard`, still has a path of zero.
- A new pinning run comes only by a new amendment to this entry.

## The gate in CI

**R05. The gate's bands lie in `sim`, in the gate file, which also names the commits of `learner-model` and `engine` on which they are recorded; the gate of a caller puts in the caller's code and takes the other library at the commit from the file. A change that breaks the simulator's API or moves the numbers is first recorded in `sim` on the commit of its pull request, then a tag of `sim` comes out, and only with it is the library's pull request merged. A release of `sim` comes out by itself, after a green CI on `main`. In the CI of `modules` the gate holds invariants and thresholds on the package of the pull request, not bands.** The author's answers 2026-11-07 in the planning of T03 (RUN.md); the rest is a proposal of the execution. How it is built — [architecture.md](architecture.md), section 3.9. What closes Т-16 — the path of `go.work`, the gate on every pull request, the runner and the time — is a question of `spec`'s register, and `sim` does not decide it but passes it on as row В3 of table 3 of RUN.md, «В реестр spec».
- **Decided by the author:**
  - by the recommendation of Т-16, `go.work` on every pull request: a library's CI fetches `sim` by tag and commit SHA and writes `go.work` on the fly, without committing it;
  - the bands are in `sim`, as R04 decided, and the gate file names the commits of the libraries on which they are recorded; the order of releases on a breaking change, on new bands and for the pair "the 'mastered' threshold and the target of the corridor" runs through a tag of `sim` recorded on the commits of the pull requests;
  - the gate in the CI of `modules` runs on every change of what the simulator reads from a package, on the package of the pull request, with invariants and absolute thresholds, without bands;
  - a release of `sim` comes out by itself, after a green CI on `main`, and is immutable; the caller pins the tag and SHA of `sim` in its workflow and raises them itself.
- **A proposal of the execution, not the author's answer:**
  - the gate file names the commits of the libraries by full SHA, the release of `modules` by its tag and the digest of the index, the packages by ID, version and digest, the versions of the shared dependencies that the workspace chose — of the Go module of `spec` and of the line `go` —, and the bands. The simulator is built only at the commits from the file, in its own CI too: its recipes write `go.work` from the file, and the `go.mod` of `sim` names the libraries only as a lower bound. Another chosen version of a shared dependency is a red gate with a reason;
  - the gate is built in the stage `toolchain` of the devcontainer of `sim` of the same tag;
  - what has changed in a package is decided by a command of `sim`: the canonical text of what the simulator would read — the module manifest, the graph and the tiers, the tasks, the lessons and the instances without the texts for the learner — against the newest published version of the same format version;
  - the thresholds of the CI of `modules`: being stuck `E2` and "nothing to give" with a non-empty pool — no more than the largest value on the published graphs in the gate file, plus 1 percentage point; on `G0` and `G2-half`, items with P* below 0.50 — no more than 25 %, and the share in the corridor `E1` — no less than half of the pool's ceiling; false "mastered" `M4` on `G0` — no more than 20 %; not one closed task with a rubric given again, `E5`;
  - after a breaking change or new bands, once the library has come out, the file of `sim` names the commit of its release, and `go.mod` the release itself; a tag recorded on a pull request that was closed without a merge stays, and the next pull request in `sim` returns the file to the released version;
  - a change of `learner-model` that breaks the API of `engine` is a pair; the pair's P is merged only when Q has passed review and is green on the same tag, `api` moves to both versions together, and if Q is not merged after all, `learner-model` rolls forward;
  - the budget of the gate's job is `timeout-minutes: 15`, an estimate that the first code of the gate replaces with a measurement; the runner is `ubuntu-24.04`. This is part of the answer to Т-16 in row В3, not a decision of `sim`.

— *Why:*
- **Bands in `sim`.** The gate's numbers depend on the simulator no less than on the libraries: on the populations, the packages and the measures. In `sim` one place records them — with the criterion's report, which also runs there —, and a change of the simulator itself, a new pin of a package or a population, does not wait for a pull request in the libraries. A caller on the previous tag stays green: a tag is immutable.
- **Commits in the file.** Two versions of one library do not go into one Go build (R04), and an edit of the simulator for an API from a pull request is not released yet. The file names what the numbers are actually recorded on, and the main module of `go.work` is taken as it is, without a choice of versions: "MVS starts at the main modules (special vertices in the graph that have no version)" ([architecture.md](architecture.md), section 4). So a candidate pair passes with one tag: each library puts in its own code and takes the second at the commit on which the pair is recorded.
- **No bands in `modules`.** The gate of `modules` judges content, not a rule, and an edit of content moves the numbers by design. With bands, each batch of the port — 705 tasks (R05 modules) — would wait for a round "pull request in `sim` → tag". The promises of the first slice — the corridor, no being stuck, no early "mastered" (product 10 spec) — are absolute and are checked on any content. The libraries' bands stay on the published packages.
- **A release by itself.** Each round of the order of releases passes through a tag of `sim`; if the author set the tag, the round would wait for them. An immutable release locks the tag on its commit (R31 spec), and the SHA in the caller's workflow holds the same content as the tag.

— *Rejected:*
- the options of Т-16 that `sim` does not recommend: the gate only on release candidates — a worsening would already lie in `main`, while a change that worsens the criterion is not merged (repositories 4 spec); `go mod edit -replace` on the fly — edits the simulator's module file and gives the same as `go.work`; the three libraries in one repository — against R10 spec;
- bands in the library's repository, as `guard.csv` in standalone: this would change R04, the bands would spread over three repositories, and a change of the simulator itself would wait for a pull request in the libraries;
- the gate at the commit of the second library that MVS chooses by `go.mod`: a pull request of `engine` that raised `learner-model` would change the gate's numbers without a record;
- bands in the CI of `modules` too; the gate of `modules` on the graph only — it would not see reassigned tasks, nor a task's band and difficulty;
- a reusable `gate.yml` in `sim` that Dependabot would raise by itself: the gate would live half in `sim`, while R03 .github calls it a workflow of the callers' own; tags of `sim` by the author's hand.

— *Consequence:* the code epic builds the gate file, the command "has it changed" and the release of `sim` by section 3.9, and the first tag comes out after the first versions of both libraries (R12 learner-model; В24 `learner-model`). The gate on arm64 waits for the build of `engine` — Д5 in table 4 of RUN.md. The gate's jobs in the CI of `learner-model`, `engine` and `modules` are written by their sessions — by section 3.9 and by `spec`'s answer to В3. — RUN.md, T03, tables 3 and 4; [architecture.md](architecture.md), status, sections 1, 2.2, 2.3, 3, 3.4, 3.6, 3.7, 3.9 and 4.

— *Amendment 2026-11-07:* `spec` recorded the answer to Т-16 that this entry passed on as row В3 of table 3 of RUN.md — R103 spec, by the row's recommendation, and Т-16 is closed. The first tag after the first versions of both libraries is also R103 spec, not В24 `learner-model` alone, as in the "Consequence". The rest of this entry is in force. — R06.

— *Amendment 2026-11-08:* the gate on arm64 no longer waits for `engine`, which has decided on the same bits on any build too (R10 engine), and `sim`'s own CI runs the gate on `ubuntu-24.04-arm` as well, while `ubuntu-24.04` is enough for the callers. Д5 of RUN.md is closed, and the job is work of the code epic. — R07; [architecture.md](architecture.md), section 3.9.

— *Amendment 2026-11-09:* the thresholds of `E1` in the CI of `modules` are taken from the gate file, as those of `E2` are. The author's answer in the planning of T25 (RUN.md, T25).
- On `G0` and `G2-half` the share of items with P* below 0.50 is no more than its largest value in the gate file on these populations plus 1 percentage point.
- The share in the corridor is no less than the smallest value minus 1 point.
- The absolute 25 % and half of the pool's ceiling stay a target of the engine's criterion, not a threshold of the gate.

A measurement in planning, 2026-11-09: 300 working learners a cell, `engine` v0.1.0, the published packages.

| Cell | P* < 0.50 | Corridor | ½ of the pool's ceiling |
|---|---|---|---|
| `G0`, olympiad | 52.7 % | 10.9 % | 37.6 % |
| `G0`, school | 48.4 % | 19.0 % | 23.1 % |
| `G2-half`, olympiad | 42.5 % | 11.3 % | 39.4 % |
| `G2-half`, school | 23.9 % | 19.3 % | 34.3 % |

The previous thresholds were not met: the share of P* < 0.50 kept within 25 % only on `G2-half` of school maths, and the corridor fell short of half of the ceiling in every cell. The gate of `modules` would have been red on every pull request, even without an edit of content.

Why from the file: just as for `E2`, content must not make the engine worse than it already is on the published graphs, and the engine is judged by its criterion.

Rejected:
- the absolute thresholds as they were — the gate of `modules` is red until the accepted engine reaches them;
- `E1` only in the report, without red — thresholds that hold nothing.

The thresholds of `E2`, `M4` and `E5` are in force. The largest value of `E2` in the gate file is 20.3 % on `F-model` of school maths, so the threshold of `E2` is 21.3 %. How the code holds the thresholds is R27. — [architecture.md](architecture.md), section 3.9; README.

— *Amendment 2026-11-09:* the releases of `sim` are put out by the shared `release.yml` of `.github` v0.1.21, with one file — an archive of what builds the simulator and runs the gate —, not by a workflow of its own without files. The author's answer in the planning of T27 (RUN.md, T27). The shared release does not come out without a file, and a workflow of its own, like `spec`'s, would be release code that only `sim` keeps, and Dependabot would not raise it together with the other workflows of `.github`. The callers still take `sim` from git by tag and SHA; the archive signs the content of the tag, and `release-verify` accepts it. How it is built — R29. — [architecture.md](architecture.md), section 3.9.

## The tag `spec` v0.9.5

**R06. `sim` stands on `spec` v0.9.5, the latest release of `spec`: the tag brought `spec`'s answer to the detailing of `sim` — R103 spec, Т-16 closed —, and rows В1–В4 of table 3 of RUN.md are closed by it.** The author's instruction 2026-11-07, «обнови до последней версии spec». The task is T07 RUN.md.
- **What the tag decided for `sim`** (`git -C spec diff v0.9.4 v0.9.5`, read 2026-11-07; Д111 in table 4 of RUN.md `spec`): R103 spec — the answers of `spec`'s author to the rows of `sim` В1–В4 and of `learner-model` В24, all by the rows' recommendations:
  - В1 — `sim` reads `server.json` by the types `formatv0.CanonicalView`, without `module-kit`, and the window of format versions by the `module_format_version` of the module manifest in the canonical view; «через `module-kit`» for `sim` in Д59 and Д107 RUN.md `spec` is an error of those rows;
  - В2 — card 3.9 spec says "populations", and the vocabulary got "virtual learner" and "population" (concepts 15 spec);
  - В3 — Т-16 is closed by `sim`'s answer: `go.work` on the fly with the other library at the commit from the gate file, the gate on every pull request and commit of `main` on `ubuntu-24.04` with `timeout-minutes: 15`, the order of releases through a tag of `sim`, the gate of the CI of `modules` without bands. So repositories 1 and 4 spec and card 3.10 spec now say;
  - В4 — card 3.9 spec: "everything is decided", the gate's data, the synthetic packages, the release of `sim` and the first tag after the first versions of both libraries; card 3.8 spec — the first version of each library before the gate;
  - В24 `learner-model` — row 4.6 of [COORDINATION.md](../spec/COORDINATION.md) is split into 4.6–4.8: the first version of `learner-model`, the first version of `engine` and the changes of both through the gate. Row 4.2 — the detailing of `sim` — is closed, and 4.5 — the code epic of `sim` — says that the simulator can be started now on synthetic packages.
- **What else is in the tag:** R104 spec — the engine's inputs: the summary of issues from `api` and the flags of the trial series in `Read` of `learner-model`, and the engine chooses a lesson or a task, not an item. The simulator gives the engine its input as `api` would give it ([architecture.md](architecture.md), sections 1 and 3), so the simulator assembles the summary of issues from its own choices and attempts. And since the engine chooses a task as a whole, not an item, the engine's ceiling, the pool's ceiling of `E1` and the meanings of the codes in `I2` change: the 2026-11-07 amendment to R04. The 2026-11-07 amendment to R93 spec is about the eval sample file of `tutor` and `canvas` and changes no document of `sim`. The template, `.devcontainer/`, RULES.md and the documents check did not change between the tags (`git -C spec diff --stat v0.9.4 v0.9.5`).
- **What is changed here:**
  - `spec/` — to v0.9.5, commit 86f442a; `gh release verify v0.9.5 -R altnovatrail/spec` confirmed the release by its attestation. The tag in CLAUDE.md, README.md and the opening of RUN.md is v0.9.5;
  - RUN.md — table 3 without В1–В4, and its opening names what closed them; epic E2 in table 2 — by R103 spec and rows 4.5–4.8 of COORDINATION.md; the task T07;
  - [architecture.md](architecture.md) — the status, sections 1, 2, 2.1, 2.3, 3.1, 3.3, 3.6 and 3.9: the references to the rows of table 3 are replaced by `spec`'s answer; sections 1 and 3 — the summary of issues and the flags of the trial series at the engine's input, sections 3.5 and 3.6 — the ceilings and `I2` by the choice of a task as a whole (R104 spec); section 3.3 — the synthetic packages not "while `modules` has no releases";
  - CLAUDE.md and README.md — Т-16 closed, populations instead of profiles;
  - the 2026-11-07 amendments to R03, R04 and R05 — that `spec` recorded their rows, the first release of `modules`, and the second one to R04 — the choice of a task as a whole.

— *Why now:* the author's instruction 2026-11-07. v0.9.5 is the latest release of `spec` on 2026-11-07, and immutable (`gh release list -R altnovatrail/spec`). Д111 of `spec` has `sim` close В1–В4 by its move to the tag with R103 spec (RULES.md, Across repositories). — *Rejected:* waiting for the next release — rows В1–В4 would hang open, though the answer to them is out.

— *Consequence:* table 3 of RUN.md is empty, and epic E2 no longer needs `spec`'s answer about the gate. — RUN.md, T07; [architecture.md](architecture.md), status and sections 1, 2, 2.1, 2.3, 3.1, 3.3, 3.6, 3.9.

## The tag `spec` v0.9.8

**R07. `sim` stands on `spec` v0.9.8, the latest release of `spec`: the tag brings the outcome of the detailing of `engine` — R107 spec, the code "in rotation", the additions to the summary of issues, the age category as a filter of `api` —, and the 2026-11-07 amendment to R103 spec about the move of the format. The debt Д5 is closed: `engine` has decided the same bits on any build too (R10 engine).** The author's instruction 2026-11-08, «выпущена новая версия spec — перейди на неё». The task is T08 RUN.md.
- **What the tag decided for `sim`** (`git -C spec diff v0.9.5 v0.9.8`, read 2026-11-08; Д113 and Д114 in table 4 of RUN.md `spec`):
  - R107 spec — Д114 names `sim`: «код „по кругу“ среди кодов, которые стенд ждёт у выбора». `ReasonCode` has a ninth value, `REASON_CODE_ROTATION` — "a skill given before comes round again, the one not given for longest first, though it does not fit the queue of `REASON_CODE_QUEUE`: it is mastered, or no skill is within reach" ([contracts text](../spec/contracts/README.md), section 15). `I2` expects it — the 2026-11-08 amendment to R04;
  - R107 spec — in the same place: "at once" of R82 spec means "not as the next task", and the pause of `engine` is 24 hours; the engine takes no randomness, and the seed of the choice is empty; the summary of issues comes with additions: the issue of a lesson is the issue of each of its skills, a module has the number of choices, a skill the outcome of its last closed attempt as facts, if it is a probe or a repetition; the age category is a filter of `api`, not an input of the engine. The simulator gives the engine its input as `api` does, so it takes all of this ([architecture.md](architecture.md), sections 1, 3, 3.4 and 3.5);
  - the 2026-11-07 amendment to R103 spec — in a move of the format `modules` changes nothing of what the simulator reads from a package, except the record itself in the new format version, and the move goes without the simulator ([architecture.md](architecture.md), section 3.9);
  - the detailing of `engine` is decided — R03–R10 engine, row 4.4 of [COORDINATION.md](../spec/COORDINATION.md) is closed —, and the engine too has decided on the same bits on any build (R10 engine). This is the condition for returning to Д5: the gate on arm64 is work of the code epic, the 2026-11-08 amendment to R05;
  - `learner-model` released its first version, v0.1.0, 2026-11-07 (`gh release list -R altnovatrail/learner-model`), and row 4.6 of COORDINATION.md is closed: the model's base exists. The first version of `engine` — row 4.7 — not yet;
  - R105 spec and the new rule of RULES.md, Hard rules, "Production images" — about the images that run in production. `sim` is a CI tool, not a service (repositories 1 spec), and it has no image in production; R106 spec is about the deployment of `infra`. They change no document of `sim`. The template, `.devcontainer/` and the documents check did not change between the tags (`git -C spec diff --stat v0.9.5 v0.9.8`).
- **What is changed here:**
  - `spec/` — to v0.9.8, commit b1fbc61; `gh release verify v0.9.8 -R altnovatrail/spec` confirmed the release by its attestation. The tag in CLAUDE.md, README.md and the opening of RUN.md is v0.9.8;
  - [architecture.md](architecture.md) — the status; sections 1 and 3 — the summary of issues with additions and the pool by the age category; 3.1 — one category per learner; 3.4 — the seed of the choice is empty; 3.5 — `I2` with "in rotation", "new skill" and "the queue" by the contract, `I4` with "at once", the pool's ceiling of `E1`; 3.6 — the engine's ceiling with the exclusions of R82 spec; 3.9 — arm64 and the move of the format; section 4 — `engine` b962ff4 and `learner-model` v0.1.0;
  - the 2026-11-08 amendments to R04 and R05;
  - RUN.md — the task T08; table 3 — В5 «В реестр spec» about `docs-check`; table 4 — Д5 closed; epic E2 in table 2 — by the first version of `learner-model` and the detailing of `engine`, with a job on arm64.

— *Why now:* the author's instruction 2026-11-08. v0.9.8 is the latest release of `spec` on 2026-11-08, and immutable (`gh release list -R altnovatrail/spec`). The intermediate v0.9.6 and v0.9.7 are not taken separately: one move takes everything up to the latest, as T06 took v0.2.9 → v0.9.4. — *Rejected:* three moves, one tag each — the same, three times as long.

— *Consequence:* table 3 of RUN.md holds В5 «в реестр spec» about `docs-check`, and Д3 and Д4 stay in table 4. — RUN.md, T08; [architecture.md](architecture.md), status and sections 1, 3, 3.1, 3.4, 3.5, 3.6, 3.9 and 4.

## The tag `spec` v0.9.9

**R08. `sim` stands on `spec` v0.9.9, the latest release of `spec`: the tag answers row В5 — the name of a NovaTrail repository after a number or a section is the suffix of the citation wherever a word boundary follows it, and after О-, П- and Т- it is a word of the phrase (the 2026-11-08 amendment to R32 spec) —, and the shared `go.yml` got the input `same-bits` (architecture 3.11 spec). By which job `sim`'s own CI runs the gate on arm64 — one of its own or this input — the code epic decides.** The author's instruction 2026-11-08, «перейди на spec последней версии». The task is T09 RUN.md.
- **What the tag decided for `sim`** (`git -C spec diff v0.9.8 v0.9.9`, read 2026-11-08; Д116 in table 4 of RUN.md `spec`):
  - the 2026-11-08 amendment to R32 spec is the answer to В5, but not by the row's recommendation; it is its second option: the name of a NovaTrail repository after a number or a section is the suffix of the citation wherever a word boundary follows it, a name that can be a word of prose too; a repository's citation of its own does not put after it a repository's name as a word of prose; after О-, П- and Т- the name is a word of the phrase. Д116 has `sim` close В5. The documents of `sim` need no edit for the amendment: the documents check is green on v0.9.9 before the edits too, and the two places that В5 was about T08 has already rewritten so that the suffix ends the phrase;
  - the card architecture 3.11 spec — the Go wrapper has the input `same-bits` (R107 spec; the 2026-11-08 amendments to R16 .github): with it run the jobs `targets`, `fma-diff`, `arm64` and `rules`. It is one, boolean, and is not taken in parts. [architecture.md](architecture.md), section 3.9, went around the option that Д114 RUN.md `spec` promised by a job of its own on arm64, as `go-arm64` of `learner-model` does. Now the code epic has two ways, and section 3.9 names what is for and against. For `same-bits`: the job `arm64` runs `ci-test`, and with it the gate, by a recipe of `sim` through `go.work` from the gate file, and `fma-diff` and the rule `fused` hold the same bits of the simulator's distributions. Against: `targets` runs `go test ./...` without `go.work` from the file, so a test at `GOAMD64` v3 does not check the gate; the runs on 386 and in WebAssembly the simulator does not need, and within the job's 30 minutes they repeat all the tests four more times; the rule `math` lets only 14 exactly defined functions into the code, while the learner's truth needs `Exp`, `Pow` and `Log`; the rule `goroutine` lets goroutines only into tests. A job of its own carries none of these limits, but then the simulator holds arm64 and v3 itself. The conclusion of section 3.9 about the callers does not change: `same-bits` runs the same code on other builds, and the shared `go.yml` still has nothing to fetch `sim` with and build it through `go.work`;
  - row 4.7 of [COORDINATION.md](../spec/COORDINATION.md) — `engine` has done T13–T19 of its epic E2, there is no first version: the condition on which the author unfolds E2 has not come;
  - the rest — the 2026-11-08 amendments to R05 modules (the port of standalone from commit 42b4408), to R05 engine (the tier of a lesson in the author's order), to R94 and R95 spec (snapping and the eval report), the warning of `module-kit` about a number in a module's text (the 2026-11-08 amendment to R57 spec) and Т-38 — changes no document of `sim`: the simulator reads the published packages by the gate file, not the port, and `engine` holds the order of the lessons. The template, `.devcontainer/` and RULES.md did not change between the tags (`git -C spec diff --stat v0.9.8 v0.9.9`); `tools/docs-check` changed — that is the amendment to R32 spec itself.
- **What is changed here:**
  - `spec/` — to v0.9.9, commit fa5da0b; `gh release verify v0.9.9 -R altnovatrail/spec` confirmed the release by its attestation. The tag in CLAUDE.md, README.md and the opening of RUN.md is v0.9.9;
  - [architecture.md](architecture.md) — the status; section 2.3 — "Deferred" with a job on arm64; section 3.9 — why a caller takes `sim` with a job of its own, and the two ways of the gate on arm64; section 4 — `.github` 706f406, `learner-model` 50f1b51 and the releases of `engine`;
  - RUN.md — the task T09; table 3 — В5 closed; epic E2 in table 2 — the state of `engine` and the question of `same-bits` among what the planning decides.

— *Why now:* the author's instruction 2026-11-08. v0.9.9 is the latest release of `spec` on 2026-11-08, and immutable (`gh release list -R altnovatrail/spec`). — *Rejected:* deciding in the move by which job the gate runs on arm64: this is a choice of the simulator's code — whether it has goroutines, whether the gate runs in `ci-test` —, and there is no code yet, and a choice without it would be a guess. The planning of E2 decides it.

— *Consequence:* table 3 of RUN.md is empty, and Д3 and Д4 stay in table 4 — the tag did not touch them. — RUN.md, T09; [architecture.md](architecture.md), status and sections 2.3, 3.9 and 4.

## The code epic

**R09. Epic E2 runs up to and including the first tag of `sim`: the simulator's code, the decisive run of the base and of the ceilings, the first file of bands, the release and the first tag. The same bits of the simulator's code are held by the input `same-bits` of the shared `go.yml`, and `sim`'s own CI runs the gate by a job of its own on `ubuntu-24.04` and `ubuntu-24.04-arm` and compares the numbers to the last digit.** The author's answers 2026-11-08 in the planning of E2 (RUN.md, table 1); the order of the tasks T10–T28 is a proposal of the planning. This is the answer to the question that R08 left to the code epic, and to the question of table 2 about what the epic consists of.
- **What was found in planning** (read 2026-11-08; the sources — [architecture.md](architecture.md), section 4):
  - `engine` released its first version, `v0.1.0`, at c6c5bb9: the release is immutable, `gh release verify` and `just -f spec/justfile release-verify engine v0.1.0` on `spec` v0.9.9 accept it, and the Go proxy serves the module with an Origin at the tag's commit. Its `go.mod` requires only `learner-model` `v0.1.0`. So all three inputs of the simulator are out — with `learner-model` `v0.1.0` and `modules` `v0.1.0` (R07; the 2026-11-07 amendment to R03);
  - the Go module of `spec` has no generator — only the test vectors `format/vectors/random` and `format/vectors/templates` with the types `formatv0.RandomVector` and `formatv0.TemplateVector`, embedded in `spec.Format`; the implementation of `module-kit` is in its `internal/random`, not open to import. The simulator builds the place of randomness itself, by the [format text](../spec/format/README.md), section 22, and runs the vectors as Go tests;
  - the shared `go.yml` of `.github` `v0.1.21` runs without the submodule `spec`, and `targets` and the second half of `arm64` run `go test ./...` without `go.work` from the gate file (section 3.9). The shared `release.yml` puts out no release without a file in `dist/`;
  - the packages of `modules` `v0.1.0` have neither prerequisites across modules nor slots `predicate`, olympiad maths has only slots `number`, the index's `package_digest` is the sha256 of the package's file `inventory.json`, and the `options` of a choice lie sometimes inside `check`, sometimes beside it.
- **Decided by the author:**
  - the epic runs up to and including the first tag: the decisive run of the base and of the ceilings, their values by an amendment to R04, the first file of bands (Д4 in table 4 of RUN.md), the release and the first tag;
  - the gate on arm64 is the input `same-bits: true` at the job "Go" of the shared `go.yml`, as `engine` has it, and a gate job of its own in `ci.yml` on `ubuntu-24.04` and `ubuntu-24.04-arm`, comparing the numbers to the last digit.
- **A proposal of the planning, not the author's answer:**
  - the order: the frame and CI (T10); seeds and packages (T11, T12); the learner, the run and the invariants (T13–T15); the populations and the synthetic packages (T16–T19); the measures (T20, T21); the ceilings and the criteria (T22, T23); the gate and the command "has it changed" (T24, T25); the decisive run (T26); the release and the first tag (T27); the summing up (T28). The first run is `G0` on the published packages in T14, and synthetic packages are needed only by the populations without a published graph;
  - its own "workflow" of the generator's test vectors (section 1; architecture 3.2 .github) is Go tests in `ci-test` under "Go / test", as in `module-kit`: the vectors are embedded in the module `spec`, and they do not need the submodule;
  - what is open and is decided in the tasks: the layout of the packages and the tools (T10), the algorithm of the normal distribution (T11), the graph `G7` (T16), the form of the results file (T22), the form of the gate file (T24), where the decisive run goes (T26), by which workflow the release comes out (T27).

— *Why:*
- **Up to the first tag.** The gate jobs of the callers wait for the first tag — T22 `learner-model`, T23 `engine` and Д5 `modules` —, and through them every candidate (row 4.8 of [COORDINATION.md](../spec/COORDINATION.md)). Everything the decisive run and the tag need is out; a separate epic would wait for nothing and would only put off the gate of three repositories.
- **`same-bits`.** The gate's bands are read to the last digit (section 3.7), and the simulator's distributions must give the same bits on any build (section 3.4). `fma-diff` and the rule `fused` check this on all targets in advance, while a run of the gate on arm64 checks it only on one, and only by the numbers. The argument against, the rule `math`, is removed by `detmath` of `learner-model`: it gives `Exp`, `Log` and `Pow` with the same bits (R09 learner-model), and `sim` depends on it anyway. The argument about the runs on 386 and in WebAssembly — four extra runs of all the tests within 30 minutes — is accepted as the price: the gate and the decisive run do not go in `ci-test`, the code's tests are small, and a test that started processes would not pass on js/wasm, so the processes of the decisive run are started by a recipe, not by the code.
- **A gate job of its own.** The gate needs the commits of the libraries and the packages from the gate file and `go.work` from it (section 3.9), while the shared `go.yml` runs the recipes on the caller's code, without the submodule, and part of its jobs without `go.work`. So the gate does not go into `ci-test` but runs as a job of its own on both architectures.

— *Rejected:*
- the epic ends with the code, and the decisive run, the bands and the tag go to the next one: see the "Why";
- only jobs of its own, without `same-bits`: without the checks at `GOAMD64` v3 and `fma-diff`, nobody would see a divergence on a build without FMA until it moved the gate's numbers;
- the gate in `ci-test`: in `targets` and the second half of `arm64` it would run on the versions of the libraries from `go.mod`, not on the commits where the bands are recorded, and without the submodule;
- goroutines in the code for the decisive run: the rule `goroutine` does not let them in, and the run is split into processes by learners — the learners' seeds do not depend on the order (section 3.4).

— *Consequence:* table 1 of RUN.md is E2 in the tasks T10–T28, the row E2 left table 2, and В6 «в реестр spec» tells `spec` that the epic is laid out and the first tag is T27. The rules of code are in «Общее для задач кода» of RUN.md. T26 closes Д4. — RUN.md, tables 1–4 and the bodies of T10–T28; [architecture.md](architecture.md), status and sections 2.3, 3, 3.9 and 4; README.md.

**R10. `sim` has the Go module `github.com/altnovatrail/sim` with `go 1.27.3` and without `require`: the root is `package main`, the program `sim`, and the parts of the simulator are the packages `internal/<part>`. Its code is checked by the job "Go" of the shared `go.yml` of `.github` v0.1.21 with the input `same-bits` (R09): the recipes `ci-lint`, `ci-test`, `ci-vuln` and `ci-licenses` in the `justfile`, coverage in Codecov, the rules of code with the lists of `.github/coderules.json`; CodeQL reads Go. All references to the shared workflows are to v0.1.21. `ci-lint` is `gofmt`, `go vet` and a check of the version of gopls (R11), without other linters.** The author's answers 2026-11-08 in the planning of T10 (RUN.md, T10, «Решено автором»); the rest follows the frame of `engine` (R12 engine).
- **Decided by the author:**
  - `go.mod` without `require`. The Go module `spec` comes with T11 — the generator's test vectors —, at the same tag as the submodule (R64 spec), and `learner-model` and `engine` with T12, the translation of the packages, together with the check of the release of `learner-model` by `release-verify` and the lines of `.github/coderules.json`;
  - the layout: the root is `package main`, and the parts are `internal/<part>` by [architecture.md](architecture.md), section 3, "Design". A caller that builds the simulator through `go.work` calls `go run github.com/altnovatrail/sim <command>`;
  - the steps outside the repository — commits, pull requests, the merge, the ruleset of `main`, Codecov, the Dependabot pull requests and Rebuild Container — are done by the author (rule 7 of RUN.md).
- **A proposal of the execution, not the author's answer:**
  - the `justfile` follows `engine`: `ci-test` is `go test -race -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage/go.out ./...`; govulncheck v1.9.0 and go-licenses v2.0.1 run through `go run` at an exact version; `test-targets`, `fma-diff`, `rules` and `same-bits` run in the container the scripts of the commit of `.github` at which `ci.yml` pinned `go.yml`. There are no release recipes until T27, and no gate among the recipes `ci-` at all (R09);
  - in T10 the program only says how to use it: without a command and with `help` — the usage and code 0, with an unknown command — code 2. Without a package with a test there is no `coverage/go.out`, and "Go / test" fails; the test replaces `os.Exit`, because `sim` cannot be started as a process of its own on js/wasm;
  - `.github/coderules.json` lists only what the code imports today, and grows with the tasks; `.devcontainer/gopls` is a nested module that no code reads;
  - `gomod` in `dependabot.yml` with `ignore` for the modules `spec`, `learner-model` and `engine` in advance: they are raised by a move to a tag of `spec`, by a task with a check of the release or by the gate file (R64, R77 spec; [architecture.md](architecture.md), section 3.9).

— *Why:* the frame of `engine` stands on the same `go.yml` with `same-bits`, its recipes have already passed CI, and the shared names of the recipes are the same in all repositories in Go (R16 .github). A `require` that no line imports `go mod tidy` would remove, and the check of a release belongs to the task that takes the module. `cmd/sim` is an extra level: there is one program, and `sim` is never a library ([architecture.md](architecture.md), section 1).

— *Rejected:* lower bounds of `spec`, `learner-model` and `engine` in `go.mod` already in T10, as the body of the task said; `cmd/sim`; staticcheck and other linters in `ci-lint` — as in R12 engine; release recipes before T27, when it is decided by which workflow the release comes out.

— *Consequence:* the ruleset of `main` requires "Go / lint", "Go / test", "Go / vuln", "Go / licenses", "Go / targets", "Go / fma-diff", "Go / arm64", "Go / rules", "CodeQL / analyze (go)" and `codecov/patch`, once its status comes — a step of the author. The Dependabot pull requests #5–#9, which raised `.github` to v0.1.10, are overtaken. — RUN.md, T10, T11, T12; `go.mod`, `main.go`, `justfile`, `.github/`; [architecture.md](architecture.md), sections 1, 3 and 4; README.md.

— *Amendment 2026-11-08:* the author did Rebuild Container themselves, and gave the other steps outside the repository — commits, pull requests, the required checks of the ruleset of `main`, closing the Dependabot pull requests #5–#9 — to the executor; the merge stays with the author (RUN.md, T10). The rest of this entry is in force.

— *Amendment 2026-11-09, the author's answer in T30:* `ci-test` runs without the race detector — `go test -count=1 -covermode=count -coverpkg=./... -coverprofile=coverage/go.out ./...`. There are no goroutines in the code, where the rule `goroutine` does not let them in (R09), nor in the tests, so the detector can find no race. And it slowed the run of a learner of `olympiad-math` 0.1.1 about 10 times. On the pull request of T30 the packages of the root and `internal/measure` went over the 10 minutes that `go test` gives a package, and "Go / test" and "Go / arm64" failed by timeout. Rejected: keeping `-race` and running the tests of the mechanics on school maths instead of olympiad maths — the tests would read the numbers on the published olympiad graph less often, and the detector would still find no race. A test that starts a goroutine brings `-race` back in the same edit (R31).

**R11. gopls stands in the image of the devcontainer at an exact version, v0.24.0, and is built from its own lock file `.devcontainer/gopls/`, as in R13 engine and R14 learner-model. This is how the container of `sim` differs from the container of `spec`, until the template of `spec` gives gopls to the kind `go`.** The author's answer 2026-11-08 in the planning of T10 (RUN.md, T10, «Решено автором»).
- The step of the Dockerfile and the lock file are the same as in `engine`: `ARG GOPLS_VERSION=0.24.0`, `go install -mod=readonly tool` from `.devcontainer/gopls/`, and the build stops if gopls names another version. CI does not build the stage `devcontainer`, so `ci-lint` checks the lock file against the `ARG`. The Go extension does not update its tools and does not check for updates, and the linter on save is off (`devcontainer.json`).
- The image was not built here: there is no docker in the container. The same step, `go install -mod=readonly tool` from a copy of the lock file, built gopls, and `gopls version` named `golang.org/x/tools/gopls v0.24.0`.

— *Why:* without the step the Go extension installs gopls itself, without an exact version, and that is against "Exact versions only" (RULES.md, Hard rules). `engine` and `learner-model` already hold gopls this way, and a shift of the version in all three is one edit.

— *Rejected:* a container as `spec`'s, without gopls.

— *Consequence:* Rebuild Container is a step of the author; the version of gopls is raised by hand, together in the lock file and in the `ARG`. — `.devcontainer/`, `justfile`; CLAUDE.md, Environment — the difference is recorded here, as that section says.

— *Superseded 2026-11-10 by R32:* the step of gopls and its lock file are now the template's of `spec`, of the kind `go` (R102 spec), not its own. The text above is not edited.

## Seeds and distributions

**R12. The simulator builds the place of randomness itself — FNV-1a 64 from the seed and the names, and PCG of `math/rand/v2` — in the package `internal/seed`, and holds it by the test vectors of `spec`. The normal distribution is the inverse distribution function by AS 241 on one `fraction`, with `detmath.Log` of `learner-model` and `math.Sqrt`; the exponential one is `0 − ln(1 − u)` with `detmath.Log`. The YAML of the vectors is read by the test dependency `go.yaml.in/yaml/v4`.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T11 (RUN.md, table 1 and T11).
- **The place of randomness:** the high word is FNV-1a 64 (`hash/fnv`) over the UTF-8 of the seed in decimal and, for each name, "|" and the name; the low word is the seed itself; the stream is `rand.NewPCG(high, low)` ([format text](../spec/format/README.md), section 22). The five functions of the format are methods of the stream, each a Go method: `Uint64`, `Uint64N`, `a + Int64N(c − a + 1)`, `Float64`, `Perm`. A name that is empty or holds "|" is an error. The number of a template's instance is `below(count)` of the place of the attempt's seed with the names `template`, the task's ID and `instance` ([format text](../spec/format/README.md), section 23; R60 spec).
- **Sets and places:** the working learners are seed 20261107, the experiment `working`; the held-out ones are only those R04 recorded: `held-out-model-1` — 20261108, `held-out-engine-1` — 20261109, and a set not recorded is an error until an amendment to R04 records it; the unit of a learner is `<population>/<number>`, of a task `task/<module>/<task>`; the purposes are [architecture.md](architecture.md), section 3.4.
- **The normal one:** AS 241, PPND16 — Wichura, M. J., "Algorithm AS 241: The percentage points of the normal distribution", Applied Statistics 37(3), 1988, pp. 477–484 ([architecture.md](architecture.md), section 4) — rational approximations in three regions, with `detmath.Log` and `math.Sqrt`; a `fraction` equal to 0 (one in 2^53) is taken as 2^−53, the smallest above it, so that the number is finite and the tails end at the same distance: the largest 1 − u is also 2^−53. The exponential one is `0 − ln(1 − u)`, so that `u = 0` gives +0, not −0. Each product and quotient goes straight into `float64()`, and `Fraction` itself returns `float64(...)`: otherwise Go's built-in `Float64` would fuse on arm64 with `1 − u` into one step — "fma-diff" caught this.
- **Checks:** the vectors `vectors/random` and `vectors/templates` from `spec.Format`; Φ(x) through `math.Erfc` — it is on the list `test-math` — returns p from 1e-300 to 1 − 1e-12 within an error of 1e-14 in x; the moments of 200,000 numbers; the golden bits of the first 32 numbers of each distribution in `internal/seed/testdata/distributions.txt` — the same file on all the builds that `same-bits` runs; no place of one set draws the same first number as a place of another.

— *Why:* the format defines the generator through Go's PCG and its methods ([format text](../spec/format/README.md), section 22), while the Go module `spec` has no implementation in Go, and `internal/random` of `module-kit` cannot be imported from another module (R09). The inverse function takes exactly one `fraction` per number, so a place's stream does not depend on which distribution reads it, and it needs only the logarithm and the root — both with the same bits. `go.yaml.in/yaml/v4` already reads the same vectors in `module-kit`; its license is Apache-2.0, part of the files MIT (`LICENSE` and `NOTICE` of the module v4.0.0, read 2026-11-08), both on the list `ALLOWED_LICENSES`, and the module does not go into the code that `sim` builds. `learner-model` comes before T12 for `detmath`: `release-verify` of its release was accepted 2026-11-08 — "signed by the shared release.yml of .github v0.1.15".

— *Rejected:* an implementation of PCG of its own; Box–Muller — it needs the sine and the cosine, which are neither in `detmath` nor among the exactly defined functions of `math`; Go's `NormFloat64` and `ExpFloat64` — the ziggurat calls `math.Exp` and `math.Log`, and it takes more than one `fraction` per call; a YAML parser of its own in the test; the simulator evaluating the expressions of templates — it takes the instances from the canonical view ([format text](../spec/format/README.md), section 25), and the test enumerates them itself and understands only the constraint `p != q`.

— *Consequence:* T13–T19 take numbers only from the places of `internal/seed`; a new held-out set comes by an amendment to R04 and a line in `internal/seed`. — RUN.md, T11, T12; `internal/seed/`, `go.mod`, `.github/coderules.json`; [architecture.md](architecture.md), sections 3.4 and 4.

## How epic E2 is worked

**R13. The executor leads the tasks T11–T23 alone, one at a time and in a row: plans each, chooses the answers of the planning itself, executes it and brings it to a merge into `main` by a merge commit when all the checks of the pull request are green. Rule 7 of RUN.md and "Do not commit" of RULES.md, How we work, do not apply to these tasks; rule 4 does.** The author's instruction 2026-11-08, after the merge of T10: «Все Задачи T11 - T23 делай сам в автономном режиме. Каждую сначала планируй. Сам вибирай оптимальные ответы после планирования и выполняй задачи последовательно. Сам доводишь каждую задачу до мержа в main».
- The answers of the planning go into the task as «Решено исполнителем по поручению автора (дата)», and what decides for the documents goes into the log, as the author's answers did in the earlier tasks. Review is by rule 6, one round. Rule 12 comes before each task.
- Not included: an edit of `spec` and of other repositories, of their settings and everything that costs money — still with a warning and the author's confirmation (rule 4); T24–T28 — still with the author: they hold the author's steps — the gate's ruleset, reading the decisive run, the settings of the release.

— *Why:* the author's instruction. Like the one-off instructions of T04–T10 to bring a task to a merge, it cancels rule 7 and "Do not commit" only for the tasks named; the rules of RULES.md change only in `spec`, and what is recorded here is what the author instructed for their own tasks.

— *Rejected:* waiting for the author at each plan and merge — the author instructed the opposite.

— *Consequence:* RUN.md, table 1, «Решено автором (2026-11-08), после T10»; T11–T23.

## Пакеты

**R14. Опубликованные пакеты стенд несёт в себе — индекс и архивы релиза `modules`, встроенные `go:embed` в пакет `internal/published`, с дайджестом индекса в коде, — и читает их сам, пакетом `internal/pack`, через цепочку дайджестов, по типам канонического вида, а переводит во входы `learner-model` и `engine` так, как перевёл бы `api`.** Решено исполнителем по поручению автора 2026-11-08 при планировании T12 (RUN.md, T12; R13).
- **Чтение** — по типам `formatv0.ModuleIndex`, `formatv0.Entry` и `formatv0.Inventory`: индекс — по закреплённому SHA-256 и по местам его членов, какой бы версии формата он ни был; запись версии формата вне окна — пока это только 0.3 (R44, R101 spec) — стенд пропускает, а не отвергает индекс ([текст формата](../spec/format/README.md), раздел 2); версия, которую индекс помечает `withdrawn`, — отказ: `api` попытку на ней не начинает (R63 spec); пакет — SHA-256 `inventory.json` равен `package_digest` записи, дайджест и размер `server.json` — части для сервера в `inventory.json`, а инвентарь и паспорт называют тот модуль, ту версию и ту версию формата, что запись ([текст формата](../spec/format/README.md), разделы 25 и 26). Архив лежит в папке того релиза, который называет его запись. Пакеты читаются один раз за процесс: байты закреплены, и то, что они дают, — тоже. Архив читает своё чтение tar: `archive/tar` тянет `os/user`, чей cgo правила кода не проверяют на каждой цели.
- **Перевод:** полные ID — модуль, вид и локальный ID через «/» ([текст формата](../spec/format/README.md), раздел 10), навыки-опоры `core.` — как есть; полосы — местом в паспорте; ярусы — из `skill_tiers`; граф — `learnermodel.NewGraph` (R28 spec); ячейка — вид проверки и, у выбора, число вариантов, где бы они ни стояли — рядом с проверкой или в ней; проверка по рубрике ячейки не даёт, а задача с ней — «с рубрикой»; рубрика рядом с вердиктом в одной задаче — ошибка; ловушки пункта — из `trap_values`, полными ID; происхождение ключа — «от автора», у шаблона — «проверено решателем», тегов — «от автора» (concepts 9.4 spec). Сверх входов библиотек модуль стенда держит холст задачи, её возрастные категории и число экземпляров, а у урока — число звеньев: их читает правда ученика ([architecture.md](architecture.md), раздел 3.1).
- **Возрастной фильтр** — как у `api` (R107 spec): задача и урок подходят по своим категориям или, без них, по категориям модуля; урок — только если подходит и каждая его задача (architecture 3.7 api; тест `engine`).
- **Экземпляр шаблона** — `below(count)` места зерна попытки, где имя задачи — её локальный ID, а числа — из `Instances` канонического вида в том же порядке (R12; R60 spec).

— *Почему:* прогону (T14), гейту (T24) и тестам нужны одни и те же байты, и без сети: встроенный релиз даёт их в самом `sim`, а цепочка дайджестов от закреплённого индекса держит их так же, как держала бы скачанные (RULES.md, Hard rules; R77 spec). Перевод `engine` живёт в его тесте, и расхождение с переводом `sim` его векторы не ловят (T21 engine), поэтому `sim` повторяет его правила и держит те же факты релиза своим тестом.

— *Отвергнуто:* пакеты в `testdata/` — их читал бы только тест; скачивать их по тегу в каждом прогоне; `module-kit` (R62 spec; R06); `archive/tar`.

— *Следствие:* новое закрепление пакетов — новая папка в `internal/published` после `release-verify`, новый дайджест индекса и новые полосы (раздел 3.7 [architecture.md](architecture.md)); как гейт у вызывающих берёт пакеты — T24. — RUN.md, T12; `internal/pack/`, `internal/published/`; [architecture.md](architecture.md), разделы 3.3 и 4.

— *Поправка 2026-11-09:* закреплён релиз `modules` v0.1.4 — индекс SHA-256 `a244c0ba…edb2a`. Из него `olympiad-math` 0.1.1, чей архив лежит в папке `v0.1.1/`, где его называет запись индекса, а `grade-3-math` и `english-essays` — 0.1.0, из папки `v0.1.0/`. Индекс и архив `olympiad-math` 0.1.0 из встроенных файлов убраны. Почему — R31.

— *Поправка 2026-11-10:* закреплён релиз `modules` v0.2.0 — индекс SHA-256 `a45c1e3d…f71f` и все три модуля версии 0.2.0 в формате 0.4, из папки `v0.2.0/`. Папки `v0.1.0/`, `v0.1.1/` и `v0.1.4/` убраны: пакетов 0.3 стенд больше не читает. Почему — R32.

— *Поправка 2026-11-11:* перевод берёт вид проверки `time` — `CheckKindTime` `learner-model` v0.1.5 — и у задачи то, есть ли наводящий вопрос `leading_question`, а не его слова; их читает правда ученика `G8`. Почему — R41.

— *Поправка 2026-11-11 (T43):* закреплён релиз `modules` v0.2.5 — индекс SHA-256 `cebf23cf…aec6a`, тот же байт в байт, что у v0.2.1. Из него `olympiad-math` и `grade-3-math` 0.3.0, чьи архивы лежат в папке `v0.2.1/`, где их называют записи индекса, и `english-essays` 0.2.0 из папки `v0.2.0/`. Индекс v0.2.0 и архивы математик 0.2.0 из встроенных файлов убраны. Почему — R42.

## Виртуальный ученик в коде

**R15. Правда виртуального ученика — пакет `internal/learner`: популяция — набор допущений `Assumptions`, где базовая популяция — `Base()`, а каждая другая меняет одно поле. Правду ученик тянет из места `params` в закреплённом порядке, первые ответы — по одной `fraction` места `answers` на пункт, ответы после ступеней — из `retries`. То, чего раздел 3.1 не назвал, решено допущениями стенда: день 0 — 2026-09-07, звено урока без задачи — две минуты, ответы после ступени времени не берут.** Решено исполнителем по поручению автора 2026-11-08 при планировании T13 (RUN.md, T13; R13).
- **Порядок `params`:** класс у ученика один, как у профиля, — поровну среди классов, которые называет каждый его модуль с классами, когда их больше одного (у ученика обеих математик это 3); затем для каждого модуля в порядке ID уровень модуля — сдвиг полосы старта плюс `StartSpread` · N(0, 1) — и поправка каждого навыка с вердиктом в порядке графа — `SkillSpread` · N(0, 1); затем, если класса нет, возрастная категория — поровну среди категорий, которые называет каждый его модуль. С классом категория — `under_13` ([architecture.md](architecture.md), раздел 3.1). Модули без общего класса или категории и допущения, с которыми календарь не идёт, — ошибка.
- **P*** — c* + (1 − c*) · R̄* · σ(a · (θ* + δ̄* − β*) + 0,5 · ступень), c* — произведение по ячейкам пункта: выбор из k вариантов — `ChoiceGuess`/k, не больше 1, `predicate` — `PredicateGuess`, остальные — 0; δ̄* и R̄* — средние по основным навыкам пункта с вердиктом — своим навыкам пункта или, без них, навыкам задачи; β* = шаг полосы · (p − 1) + (d − 3) + `Miss` · z, z — нормальное число места `task/<модуль>/<задача>` с назначением `miss`, одно на задачу и набор. R* — (1 + F · Δt / τ*)^(−C), Δt и τ* в сутках; до первой попытки навыка — 1.
- **Попытка:** первые ответы — по пункту в порядке пунктов; после неверного — до трёх ступеней лестницы, на каждой — ответ каждого ещё неверного пункта; четвёртый неверный ответ пункта — показ решения. Задача с рубрикой — одна отправка и «готово», 15 минут. Время — две минуты на пункт листа и минута на строку рабочей тетради. Зерно попытки — следующее целое места `attempt`. Стабильность каждого основного навыка пунктов меняется раз за попытку, по доле верных первых вердиктов, а последняя практика — начало попытки.
- **Календарь:** день — занятие, когда его `fraction` места `calendar` ниже `SessionChance`; занятие — `SessionAttempts` попыток подряд с 15:00 UTC, следующая выдаётся, когда закрылась прошлая.
- **Урок** — по звеньям до конца: звенья без задачи — по две минуты, до задач урока, а задачи — одна за другой, как любые. Промах задачи — один на задачу и набор у всех учеников, и тянется он один раз.

— *Почему:* раздел 3.1 описывает правду формулами, а популяции раздела 3.2 меняют по одному допущению; набор допущений держит каждое из них одним полем, и тест показывает, что популяция отличается от базовой ровно им. Закреплённый порядок `params` и отдельные места для ответов и повторов держат «k-й первый ответ берёт одни и те же числа при любом правиле» (раздел 3.4). День начала, время звена урока и повторов раздел 3.1 не назвал; это допущения стенда, как скорости популяций: от них зависят паузы, но не то, что стенд сравнивает, — база, кандидат и потолок идут на одних учениках и одном календаре.

— *Отвергнуто:* популяции кодом поверх ученика; правда на коде `learner-model` (R03); время ответа после ступени — повторы тогда сдвигали бы календарь, а с ним R* следующих задач, и правило с другой лестницей меняло бы паузы.

— *Следствие:* T14 даёт ученику задачи выбора движка и строит из попыток свидетельства; T16–T19 — популяции полями `Assumptions`. — RUN.md, T13; `internal/learner/`; [architecture.md](architecture.md), раздел 3.1.

— *Поправка 2026-11-08:* звенья урока идут в порядке автора, а не «сначала объяснения»: модуль держит порядок звеньев (`pack.Lesson.Order`), и прохождение по звеньям проходит его так. — R21.

## Прогон

**R16. Прогон ученика — пакет `internal/run`: на каждом времени календаря стенд строит вход движка, как `api`, берёт выбор, отвечает правдой ученика и отдаёт модели свидетельство; каждый шаг записан с входом, выбором, попытками и свидетельствами. Пределы прогона — 200 первых ответов в математике, но не дольше года, 180 дней у популяций забывания, у сочинений — до конца задач с рубрикой, но не дольше 60 дней; «нечего дать» кончает занятие.** Решено исполнителем по поручению автора 2026-11-08 при планировании T14 (RUN.md, T14; R13).
- **Вход движка:** `Read` на время выбора по модулям ученика и с его классом; пул и уроки — по возрастной категории ([architecture.md](architecture.md), раздел 3.3; R14); прохождения уроков — «пройден» с временем конца; закрытые «Готово» попытки; последняя закрытая попытка — задача, урок, если она звено прохождения, и результат на каждом основном навыке фактами; сводка выдач — у навыка время последней выдачи задачей или уроком, который ему учит, код выбора и результат его последней закрытой попытки, у модуля число выборов, урок одним, у задачи время, когда её решение ушло на устройство (R104, R107 spec). Решение уходит на устройство, когда попытка закрыта верным ответом или показом решения (R82 spec). Сводку стенд отдаёт в порядке ID.
- **Свидетельство** — только у задачи с вердиктом: исход — «с первой отправки», «после подсказки» со ступенью или «решение показано»; у каждого пункта первый и последний вердикт и неверные ответы до верного; «решение уходило на устройство раньше» — кроме шаблона, у которого этого факта нет; время последнего действия — конец попытки; стартовая полоса — полоса старта ученика. Ловушек правда не знает, приёма и просьб у ученика нет (раздел 3.1).
- **Шаг:** выбор задачи — попытка; выбор урока — прохождение по звеньям, свидетельство каждой его задачи и прохождение «пройден», а каждый навык урока и каждый основной навык его задач выдан во время выбора урока; выбор — одно место занятия, урок с задачами тоже; «нечего дать» кончает занятие, и следующий выбор — в следующий день с занятием. Шаг держит рядом с каждой попыткой её свидетельство — у попытки с рубрикой его нет — и номер экземпляра шаблона, выбранный по зерну попытки (R12); «дни» прогона — день последнего выбора; состояние модели запись держит после каждого шага.
- **Пределы** — у каждого есть дни, иначе прогон мог бы не кончиться: математика — 200 первых ответов пунктов с вердиктом, как 200 ответов standalone, и не дольше 365 дней, чтобы правило, которое ничего не даёт, всё равно кончило прогон; популяции забывания — 180 дней; сочинения — пока не закрыты «Готово» все задачи с рубрикой, подходящие категории ученика, и не дольше 60 дней ([architecture.md](architecture.md), раздел 3.1).

— *Почему:* стенд даёт движку вход, как его дал бы `api` (раздел 1 [architecture.md](architecture.md)): иначе он мерил бы движок на входе, которого в работе платформы нет. Запись шага с входом нужна инвариантам — `I1` повторяет выбор по входу, `I2` читает признаки пробной серии и сводку (раздел 3.5), — а попытки и свидетельства — мерам. Год у математики — предел, которого правило, дающее задачи, не достигает: `G0` набирает 200 ответов за 6–31 день.

— *Отвергнуто:* урок — столько мест занятия, сколько в нём задач: число выборов в сводке считает урок одним, и занятие считает так же; «нечего дать» — выбор снова через минуту: выбор по тому же входу дал бы то же, а время в занятии без попытки ничего не значит.

— *Следствие:* T15 читает инварианты по записи шагов, T20 и T21 — меры. — RUN.md, T14; `internal/run/`; [architecture.md](architecture.md), раздел 3.

## Инварианты

**R17. Инварианты `I1`–`I4` — пакет `internal/invariant`: `Check` даёт каждое нарушение в записи прогона с инвариантом, учеником, шагом и выбором, а `Same` — нарушение `I1` между двумя прогонами на одних зёрнах. Допуск — ноль.** Решено исполнителем по поручению автора 2026-11-08 при планировании T15 (RUN.md, T15; R13).
- **`I1`:** выбор по входу, который хранит шаг, — тот же выбор; два прогона на одних зёрнах дают те же выборы, те же байты каждого свидетельства (`EncodeEvidence`) и те же байты состояния модели в конце (`EncodeState`).
- **`I2`:** у выбора задачи или урока — один код из `ReasonCode`, у «нечего дать» — никакого, и код верен своему смыслу по фактам, которые знает стенд, — поправка 2026-11-08 к R04.
- **`I3`:** навык «порядка автора» идёт, когда его пререквизиты с рубрикой или без оценки сделаны — урок, который им учит, пройден, или задача на них закрыта «Готово».
- **`I4`:** задача, не шаблон, той же, что последняя закрытая попытка вне урока, и её решение ушло на устройство, — не следующей задачей (R82, R107 spec).

— *Почему:* инварианты — не статистика, одного нарушения хватает (R04), поэтому они читают каждый шаг, а не числа ячейки. Вход шага, записанный прогоном (R16), даёт `I1` повторить выбор без прогона заново.

— *Отвергнуто:* пересчитывать досягаемость и порядок очереди в стенде — поправка 2026-11-08 к R04.

— *Следствие:* гейт (T24) держит инварианты на каждом ученике, а `I1` — ещё и повтором прогона. — RUN.md, T15; `internal/invariant/`; [architecture.md](architecture.md), раздел 3.5.

— *Поправка 2026-11-09:* у прогона с повтором — у гейта — `invariant.CheckRepeated` проверяет `I2`–`I4`, а выбор на каждом шаге заново не берёт: `I1` держит до бита повтор (`Same`), а отказ движка от своего входа ловит сам прогон. Но первые 10 учеников каждой ячейки (`bench.Rechosen`) идут через `Check` и при повторе: повтор не видит правило, которое меняет свой вход на месте одинаково в обоих прогонах. Без повтора, в решающем прогоне, `Check` по-прежнему берёт каждый выбор заново. Почему — R31.

— *Поправка 2026-11-10:* выбор на каждом шаге заново берут только первые 10 учеников каждой ячейки (`bench.Rechosen`) — в каждом прогоне под правилом движка, и в решающем тоже, а не только при повторе. Остальные идут через `invariant.CheckWithoutI1` — так переименована `CheckRepeated`: повтора у решающего прогона нет, и имя говорит, что проверка не берёт выбор заново. В решающем прогоне выбор заново стоил около 40 % процессора. Почему — R33.

## Популяции в коде

**R18. Популяции — пакет `internal/population`: у каждой имя, графы, предел прогона и отличие от базовой популяции полями `learner.Assumptions`; две стороны одного допущения — по чётности номера ученика. Девятнадцать популяций standalone — `Standalone`, в порядке таблицы раздела 3.2, а `G7` идёт на школьной математике 3 класса (поправка 2026-11-08 к R03).** Решено исполнителем по поручению автора 2026-11-08 при планировании T16 (RUN.md, T16; R13).
- **Поля допущений:** `LevelShift` — уровень на 2,5 ниже или выше полосы старта (`G1`); `LearnLevel`, `LearnSkill`, `LearnReturn` — рост уровня, поправок основных навыков ответа и возврат остальных на долю пути к начальной на каждом первом ответе пункта с вердиктом, `LearnFade` — замедление роста F/(F + k), k с нуля, а возврат остальных навыков не замедляется, как у standalone (`G2`, `G2-half`, `G2-fading`); ответы — свои у каждого модуля; `JumpBy`, `JumpFrom`, `JumpTo` — скачок уровня на ответе от 50-го до 150-го (`G3`, `G3-drop`); `BandBias` — задачи одной полосы, вытянутой ученику, труднее на 0,75 (`G4-band`); `SkillHalves` — первая половина навыков графа сильнее, вторая слабее, или наоборот (`G5`); `TechniqueChance`, `TechniqueLift` — приём по просьбе до первого ответа в пятой доле попыток с вердиктом и +0,5 к логиту её пунктов (`G8`); окно скачка, доля приёма и угасание проверяются в `New`; остальные популяции меняют уже бывшие поля — `Slope`, `ChoiceGuess`, `Miss`, `StartSpread`, `SkillSpread`.
- **Порядок чисел:** то, что тянут только некоторые популяции, — ответ скачка и трудная полоса, — идёт в `params` после того, что тянет каждый ученик, возрастную категорию тоже; приём по просьбе — по одной `fraction` места `attempts` на попытку; рост и скачок — по первым ответам пунктов, после попытки. Свидетельство попытки с приёмом несёт `TechniqueOnRequest`.

— *Почему:* таблица раздела 3.2 задаёт каждую популяцию одним допущением, и тест держит, что популяция отличается от базовой ровно полями этого допущения. Чётность номера даёт половину учеников на сторону без лишнего числа, которое сдвинуло бы правду. Приём по просьбе у NovaTrail — факт свидетельства, и модель читает его как помощь (R06 learner-model), поэтому стенд отдаёт его в свидетельстве, а не только поднимает логит.

— *Отвергнуто:* сторона из `params`; популяции кодом поверх ученика (R15); граф `G7` «нет графа».

— *Следствие:* T18 и T19 добавляют новые популяции тем же путём. — RUN.md, T16; `internal/population/`, `internal/learner/`, `internal/run/`; [architecture.md](architecture.md), раздел 3.2.

— *Поправка 2026-11-11:* приём по просьбе `G8` на задаче с наводящим вопросом даёт и его — подсказка по просьбе до первой ошибки даёт и приёмы, и вопрос ([текст формата](../spec/format/README.md), раздел 19), — и свидетельство несёт `LeadingQuestionOnRequest`. Подъём тот же, +0,5, и нового числа из места `attempts` нет. Почему — R41.

## Синтетические пакеты

**R19. Синтетические пакеты — пакет `internal/synthetic`: три канонических вида в форме модулей первого среза, собранные кодом Go, валидные по схеме канонического вида формата и переведённые тем же `pack.Translate`, с ID `synthetic/sheet`, `synthetic/workbook` и `synthetic/essays`, которые критерий и гейт отвергают, и навыком-опорой общей таксономии `core.multiplication-table`, который учит рабочая тетрадь и на который опирается навык листа.** Решено исполнителем по поручению автора 2026-11-08 при планировании T17 (RUN.md, T17; R13).
- **Лист** — четыре навыка с вердиктом на трёх полосах классов 1–6: выбор из четырёх вариантов, число, `predicate`, навык на первом и на навыке-опоре; задачи сложности 1, 3 и 5 на каждой полосе; уроки.
- **Рабочая тетрадь** — одна полоса 3 класса, цепочка трёх навыков, первый учит навыку-опоре; страницы из шести–десяти строк чисел и выбора из трёх вариантов; уроки.
- **Сочинения** — без классов, все четыре возрастные категории: навык выбора, навык без оценки, которому учит урок чтения, два навыка с рубрикой цепочкой на нём; задача только для взрослых; уроки.

— *Почему:* синтетика — для кода и тестов стенда, а не для суда (R03, R04): ей нужно дать то, чего нет в опубликованных пакетах, — навык на навыке-опоре другого модуля и ячейку `predicate` ([architecture.md](architecture.md), раздел 4), — и форму трёх модулей, чтобы код стенда шёл на ней тем же путём. Типы `formatv0` не держат обязательных полей и перечислений, поэтому виды держит схема канонического вида из модуля `spec` — тест через `github.com/santhosh-tekuri/jsonschema/v6` v6.0.3, как у `module-kit`, только в тестах; лицензия — Apache-2.0 (`LICENSE` модуля, прочитано 2026-11-08). Навык-опора — из общей таксономии: формат знает только её навыки, а синтетика с опубликованными пакетами не встречается. Свои ID и `IsSynthetic` дают критерию и гейту отвергнуть синтетику.

— *Отвергнуто:* файлы `server.json` в `testdata/` — их держали бы руками; синтетика с ID опубликованных модулей — её легко было бы спутать с ними.

— *Следствие:* `S-predicate` (T18) и `A-` популяции (T19) идут на синтетике, а на опубликованных графах их проверки — «нет графа», и `M9` `A-linked` в гейте тоже (раздел 3.7 [architecture.md](architecture.md)); критерий и гейт (T22–T24) отвергают модуль, для которого `IsSynthetic`. — RUN.md, T17; `internal/synthetic/`; [architecture.md](architecture.md), раздел 3.3.

— *Поправка 2026-11-09:* на синтетике остаётся только `S-predicate`: ячеек `predicate` нет и в `olympiad-math` 0.1.1. Популяции `A-` идут на опубликованной паре математик, а синтетические пакеты держат их в тестах кода (R31).

— *Поправка 2026-11-11:* навык листа на навыке-опоре отвечается временем — `time`, ключ `7:30`, — а не числом, и у задач листа сложности 5 есть наводящий вопрос. Угадывание у числа и у времени — 0, у правды и у модели, так что числа синтетики не сдвинулись; новый навык сдвинул бы выбор движка у ученика обеих синтетик, и тест якоря перестал бы видеть освоенную опору. Почему — R41.

## Популяции S-, F-, W-

**R20. Пятнадцать популяций `S-`, `F-` и `W-` — `population.Own`, каждая полями своего допущения в `learner.Assumptions`; перерыв и плохой день тянутся из места `calendar`, сдвиг и ошибка попытки — из `attempts`, а «не прочитать» — счёт ответов на листе.** Решено исполнителем по поручению автора 2026-11-08 при планировании T18 (RUN.md, T18; R13).
- **Поля:** `Adult`, `AdultFrom`, `AdultTo` — взрослый без класса, первая полоса и уровень на 2,5–5 выше, поровну, категория `adult`; `UnreadableEvery` — каждый двадцатый ответ пункта на листе сначала «не прочитать», и ученик пишет его заново; `BreakStartFrom`, `BreakStartTo`, `BreakMin`, `BreakMax`, `BreakEndBy` — перерыв в 14–60 дней с начала от 30-го до 60-го дня, конец к 120-му, а `BreakDrop` — падение уровня каждого модуля после него — правда читает его по времени, как плохой день, а не меняет уровень; `ExponentialLaw`, `LawBase` — R* = 0,9^(Δt/τ*); `AttemptShock` — сдвиг N(0, 1) к логиту всех пунктов попытки; `ErrorChance` — доля попыток, где каждый пункт с основным навыком задачи — «пока нет» на первом ответе; `BadDayChance`, `BadDayDrop` — доля плохих дней с занятием и падение уровня в них; `RepeatChance`, `RepeatDays` — P* = 0,95 на повторе задачи, не шаблона, чьё решение уходило на устройство за 30 дней.
- **Числа:** перерыв — первым из `calendar`, начало, затем длина; плохой день — `fraction` из `calendar` после решения о занятии; сдвиг, затем ошибка — из `attempts`, после приёма `G8`; уровень взрослого — последним в `params`. «Не прочитать» числа не берёт: счёт идёт по каждому ответу пункта на листе — первому и после ступеней; первый вердикт даёт ответ, написанный заново, по своей `fraction` `answers`, а свидетельство несёт причину `UnverdictedUnreadableOther`; ответ после ступени, написанный заново, ничего не меняет. Календарь у ученика один, и «решение уходило на устройство» — тоже один источник: ученик, а прогон берёт из него и свидетельство, и сводку. Допущения, с которыми перерыв не уложится, закон не определён, доли вне [0, 1] или взрослый попадает в модуль без категории `adult`, `New` отвергает.
- **Графы и пределы** — раздел 3.2: `S-scale`, `S-adult` и `S-unreadable` — олимпиадная математика; `S-predicate` — синтетический лист (R19); `F-` — обе математики, 180 дней; `W-shock`, `W-error`, `W-badday` — школьная математика 3 класса; `W-repeat` — обе математики.

— *Почему:* раздел 3.4 тянет дни, занятия, перерывы и плохой день из `calendar`, а то, что держит попытка, — из `attempts`, «k-я попытка берёт одни и те же числа в одном порядке»; закреплённый порядок тянет одно и то же при любом правиле. «Не прочитать» — не ошибка и не первый вердикт (R04 spec; R06 learner-model), поэтому оно ничего не меняет в правде и не берёт числа, а модели приходит причиной.

— *Отвергнуто:* «каждый десятый день» и «каждая пятая попытка» счётом — раздел 3.4 называет поток; «не прочитать» долей из потока — раздел 3.2 говорит «каждый двадцатый».

— *Следствие:* T19 добавляет `A-`, `L-` и `E-` тем же путём. В R18 «`New`» реестра — теперь `Own`. — RUN.md, T18; `internal/learner/`, `internal/population/`, `internal/run/`; [architecture.md](architecture.md), раздел 3.2.

## Популяции A-, L-, E-

**R21. Семь популяций `A-`, `L-` и `E-` — в `population.Own`, полями своих допущений; `A-` идут на синтетических листе и тетради вместе; доли и звено тянутся из `attempts`; настоящее освоение — `learner.TrulyMastered`; прохождение, брошенное на звене, открыто, пока его не продолжат или 30 дней не кончат его «брошен».** Решено исполнителем по поручению автора 2026-11-08 при планировании T19 (RUN.md, T19; R13).
- **Поля:** `LevelCorrelation` — z каждого следующего модуля ρ·z₁ + √(1 − ρ²)·z, ρ = 0,7 (`A-linked`); `TeacherBoost` — навыки, которые учат навыку-опоре, сильнее на 2,0 (`A-unpractised`); `TeacherFades` — g = 0 у такого навыка, как только правда освоила его на первой полосе, где ему учат (`A-fading`); `AbandonChance` — пятая доля прохождений брошена на звене, вытянутом поровну среди тех, что ещё впереди (`L-abandon`); `ExplainLift` — первое объяснение или пример, до которого дошло прохождение урока, поднимает поправки его навыков на 0,25, один раз (`L-explain`); `LeaveChance` — треть попыток задач с рубрикой закрыта «ушёл» после тех же 15 минут (`E-leave`).
- **Настоящее освоение** — P* эталонной задачи, сложности 3 на полосе, без промаха — ни автора, ни целой полосы, — выбор из пяти вариантов с угадыванием популяции, R* = 1, не ниже 0,775 ([architecture.md](architecture.md), раздел 3.5): `Reference` и `TrulyMastered` ученика. Числа модели правда не читает. Навык-учитель и первая полоса, где ему учат, — у навыка правды с `New`.
- **Звенья урока** идут в порядке автора: модуль держит порядок звеньев урока (`pack.Lesson.Order`), и прохождение проходит его так, а не «сначала объяснения», как говорил R15.
- **Прохождения в прогоне:** брошенное — открытое, без времени конца; следующее прохождение того же урока продолжает его со звена, где бросили, — «пройден» или снова открытое; открытое кончается «брошен» через 30 дней от начала, и время конца — эти 30 дней. Движок отдаёт открытое прохождение снова на следующем же выборе (ветвь 1, architecture 3.3 engine), так что до 30 дней оно доходит, только если движок его не отдаёт; стенд даёт ему вход честно, а конец «брошен» держит тест. Время начала держит прогон: в `Run` движка его нет. Попытка с рубрикой «ушёл» — последняя попытка «ушёл» и не «Готово».

— *Почему:* раздел 3.2 ведёт `A-` туда, где навык-опора — пререквизит навыка другого модуля, а в опубликованных пакетах такого нет (R19); раздел 3.4 тянет звено `L-abandon` и уход `E-leave` из `attempts`. Настоящее освоение нужно и `A-fading`, и мерам T20, поэтому оно — метод ученика, а не код меры.

— *Отвергнуто:* `A-` на опубликованных математиках — навыка на навыке-опоре там нет, и `M9` читать нечем; «каждое пятое прохождение» и «каждая третья задача» счётом — раздел 3.4 называет поток.

— *Следствие:* все сорок одна популяция раздела 3.2 — в реестре, и T20 читает меры на них. — RUN.md, T19; `internal/learner/`, `internal/run/`, `internal/population/`; [architecture.md](architecture.md), раздел 3.2.

— *Поправка 2026-11-09:* популяции `A-` идут на паре опубликованных математик — `olympiad-math` 0.1.1 и `grade-3-math` 0.1.0, — а не на синтетических листе и тетради: два олимпиадных навыка опираются на навыки, которым учит школьная математика. Почему — R31.

## Меры модели

**R22. Меры модели `M1`–`M10` — пакет `internal/measure`: прогон ученика даёт суммы, число ячейки читается по суммам всех её учеников, а 95 % интервал — из 2 000 повторных выборок учеников из места `bootstrap` с именем ячейки; что мера читает о попытке, прогон записывает наблюдением `run.Obs`.** Решено исполнителем по поручению автора 2026-11-08 при планировании T20 (RUN.md, T20; R13).
- **Суммы и числа.** Прогон ученика даёт именованные суммы — квадраты, счёты, доли по корзинам; число ячейки — отношение сумм всех её учеников, а не среднее чисел учеников, так что ученик с большим числом ответов весит больше. Интервал — 51-е снизу и 51-е сверху значения числа по 2 000 повторным выборкам учеников, по 2,5 % с каждой стороны; выборки тянутся из места рабочего зерна с именами `bootstrap` и имени ячейки — популяции и графа, без правила, — так что база, кандидат и потолок на тех же учениках берут те же выборки, и пару можно читать по ученику ([architecture.md](architecture.md), раздел 3.6). Число, которое не читается хоть на одной выборке, — непрочитанное: его интервал стоял бы только на выборках, где случай попался; на графе с синтетическим модулем числа помечены «нет графа» (R19).
- **Наблюдение попытки.** Перед каждой попыткой с вердиктом прогон читает модель на время её начала: P каждого пункта с удержанием и без — P без удержания c + (P − c)/R̄, где R̄ и c модель считает по своей формуле (раздел 4 [architecture.md](architecture.md)), — паузу с последней практики основных навыков пункта и навыки модуля; после попытки — P* первого ответа, помощь и навыки модуля снова. Навык в наблюдении — θ + δ модели и θ* + δ* правды, «освоено» модели и правды на полосе освоения модели или на первой полосе, где ему учат, и P и P* эталонной задачи этой полосы. В уроке каждая задача закрывается, как только ученик её ответил, — прохождение зовёт прогон вокруг каждой задачи (`learner.Each`), — так что модель и правда читаются до и после неё такими, какие они тогда, с задачами урока до неё и без задач после. Попытка с рубрикой не наблюдается.
- **Правда мер** — та, что отвечает: θ* после перерыва `F-break` ниже на падение, и с ним читаются θ* + δ*, настоящее освоение, P* эталонной задачи и владение навыком-опорой. Плохой день `W-badday` в правду мер не входит: он снижает ответы этого дня, а не то, что ученик знает.
- **Счёт.** «Ответ» — первый ответ пункта с вердиктом; «ответы навыка» — пункты попыток, по которым модель взяла свидетельство навыка. Пауза навыка не меньше 14 суток — попытка, ответившая навык, у которого была практика 14 суток назад и раньше; с неё открываются окно перерыва `M4` на 7 суток и догонка `M5`. Помощь `M2` — приём по просьбе и повтор задачи с виденным решением: ступени лестницы идут после первого вердикта и его не меняют. «Освоено» `M6`, объявленное той же попыткой, что дала настоящее освоение, ждало ноль ответов. Скачок `M3` — `JumpAt` модуля у `G3`, `G3-drop` и первая попытка после конца перерыва у `F-break`, а не у других `F-`, где перерыв правду не роняет; прогон, кончившийся до скачка, скачка не имеет. Ответы до «держится» `M3` — до первого ответа серии из десяти, так что ошибка, которая держится с первой попытки после скачка, даёт ноль. `M7` — корень среднего квадрата θ − θ* модуля: мера со знаком идёт в критерий модулем (раздел 3.6). `M9` читается на первом выборе каждого дня. «Сколько объявлено» `M4` — на ученика.

— *Почему:* числа ячейки из сумм складываются по ученикам, и повторная выборка учеников — это выборка их сумм, так что интервал считается без второго прогона. Модель сама читает P без удержания только внутри шага, а наружу не отдаёт, и стенд повторяет её формулу по тому, что `Read` отдаёт, — удержанию навыков и признаку пробной серии.

— *Отвергнуто:* число ячейки как среднее чисел учеников — ученик с тремя ответами весил бы как ученик с двумястами; выборки по имени правила — пара базы и кандидата брала бы разные выборки; интервал по выборкам, где число читается, — он был бы уже и смещён; правда без падения после перерыва — модель, которая верно идёт за падением, показывала бы ошибку; «ответы навыка» как все пункты попыток модуля — навык, которого не отвечали, копил бы ответы.

— *Следствие:* T21 добавляет `E1`–`E5` тем же путём, а T23 читает разность пары на тех же выборках. — RUN.md, T20; `internal/measure/`, `internal/run/`; [architecture.md](architecture.md), раздел 3.5.

*Поправка 2026-11-08, при T22 (R24):* ответы до «держится» `M5` считаются так же, как у `M3`, — до первого ответа серии из пяти, так что ошибка, которая держится с первого ответа после паузы, даёт ноль. Выяснилось на оракуле: он догоняет сразу, а мера давала ему пять. Разность пары ячеек на тех же выборках читает уже T22 (`measure.ReadPair`), а не T23.

## Меры движка

**R23. Меры движка `E1`–`E5` — в пакете `internal/measure`, тем же путём, что меры модели (R22); что правда говорит о пуле на выборе — потолок пула, есть ли что дать, досягаемые навыки, — прогон записывает до выбора движка, `run.Pool`; досягаемость навыка читается по правде.** Решено исполнителем по поручению автора 2026-11-08 при планировании T21 (RUN.md, T21; R13).
- **Пул на выборе.** До вызова движка прогон читает каждую задачу пула, которую правило могло дать: без рубрики, с пунктом с вердиктом, не задачу с решением, ушедшим на устройство меньше 24 часов назад, — а если таких нет, второй проход, где они есть, кроме задачи последней попытки (R82 spec; R03 engine и его поправка 2026-11-07), — и все её основные навыки с вердиктом досягаемы. Потолок — наибольшая доля пунктов с P* в коридоре 0,70–0,85 у одной такой задачи; «пул даёт» — такая задача есть; досягаемые навыки — основные навыки с вердиктом таких задач.
- **Досягаемость по правде.** Навык досягаем, когда правда держит его пререквизиты: навык с вердиктом правда освоила на первой полосе, где ему учат; навык без вердикта или модуля, которого ученик не учит, держится — читать его нечем; навык-опору правда держит, как в `M9`, — навык, который ему учит, освоен с R* на время, — а если ему не учит ни один модуль ученика, он держится.
- **P\* выдачи** — шанс правды на время выдачи пункта, без того, что держит только попытка: приёма, сдвига, ошибки и повтора (`ItemObs.PGiven`).
- **Счёт.** `E1` — доли выданных пунктов по P* выдачи, потолок — среднее по выборам задачи или урока. Навыки выбора — навыки урока или основные навыки задачи, как их считают `learner-model` и `engine`: свои навыки пунктов, а навыки задачи — когда пункт их берёт. `E2` «зацикливание» — 50 выборов подряд, каждый на навыке, без смены его «освоено» модели, и на каждом есть досягаемый навык, на котором выбор не стоит; «нечего дать» — среди выборов, где пул даёт. `E3` — повторение — выбор с кодом «повторение», навык — навык выбора, у которого модель читает «освоено» и срок повторения пришёл, R* и «опоздание» — по наблюдению попытки (R* правды добавлен в `run.SkillObs`); доля повторений — по неделям календаря от первой до 53-й; R* в конце прогона — на время последнего выбора или закрытия его последней попытки — у навыков, освоенных правдой на первой полосе, где им учат. `E4` — доля досягаемых хоть раз навыков, которые выдавались, и доля навыков с вердиктом, освоенных правдой к концу. `E5` — на ученика: написанные задачи с рубрикой, каждая один раз, и попытки задачи, уже закрытой «Готово», — задачей или в уроке.

— *Почему:* досягаемость — чем движок судится, и читать её по модели значило бы судить правило его же глазами; правда — то, что стенд знает. Потолок и досягаемость требуют правды на время выбора, а правда к концу прогона уже другая, поэтому их пишет прогон, а не мера.

— *Отвергнуто:* досягаемость, как её читает `engine`, — по P модели не ниже 0,70 (`skills.go` `engine` v0.1.0): это правило движка, кандидат может его менять; P* выдачи с шумом попытки — его не предвидит ни одно правило; потолок только по выборам задачи — урок тоже выдача.

— *Следствие:* T22 строит потолок движка на том же пуле — по P* задачи с наибольшей долей в коридоре. — RUN.md, T21; `internal/measure/`, `internal/run/`; [architecture.md](architecture.md), раздел 3.5.

## Потолки и две версии рядом

**R24. Оракул модели читается на ответах того же прогона, что модель; потолок движка — своё правило прогона на том же пуле, а где пул не даёт задачи с P\*, выбирает `engine`; каждая сборка `sim` пишет файл результатов — JSON Lines, строка на ученика, — а сравнение читает два файла парами по ученику; библиотеки файл называет версией из сборки, а собранную из рабочей области — коммитом, который передаёт вызывающий.** Решено исполнителем по поручению автора 2026-11-08 при планировании T22 (RUN.md, T22; R13).
- **Оракул модели** — те же наблюдения, где оценка — правда θ* + δ*, P — P* ответа, P без удержания — тоже P*, «освоено» — настоящее освоение, P эталонной задачи — P*, уровень — θ*, а навык-опору он держит, когда держит правда (`measure.Oracle`). Модель не подменяется для движка: `Readout` `learner-model` держит состояние внутри и сам считает P задач, так что движку оракула не дать; а на тех же ответах оракул и модель сравниваются по ученику.
- **Потолок движка** — правило `run.Ceiling` вместо `engine.Choose`, на пуле `run.Pool` (R23): задача с наибольшей долей пунктов с P* в коридоре, при равенстве — с средней P* ближе к 0,775, затем по ID, с паузой 24 часа и вторым проходом правила; повторение — когда у навыка, освоенного правдой, P* эталонной задачи первой полосы, где ему учат, падает ниже 0,775 от забывания, — лучшая задача на нём, с кодом «повторение». Где пул не даёт задачи с вердиктом — задачи с рубрикой, уроки, «нечего дать», — выбирает `engine`: потолок — потолок P*, а где P* нет, остаётся правило. Инвариант `I2` к потолку не относится: его коды — не ветви `engine`.
- **Файл результатов** — JSON Lines: первая строка — набор учеников, правило (`engine` или `ceiling`) и версия или коммит каждой библиотеки; дальше строка на ученика — популяция, граф, номер, суммы мер модели и движка и суммы оракула на тех же ответах. Числа float64 JSON пишет кратчайшей записью, которая читается назад теми же битами, так что пара сравнивается до последнего бита. Команды — `sim run` (популяции, ученики на граф, набор `working` или `held-out-<критерий>-<n>`, правило, `-commit`) и `sim compare` (две стороны, каждая — суммы правила или оракула): для каждой ячейки и числа — a, b и разность a − b с 95 % интервалом из тех же выборок в обеих (`measure.ReadPair`); сторона оракула читает только меры модели — мер движка у оракула нет, — а у потолка движка не читается `E3`. Файл `-out` пишется рядом и встаёт на место, только когда прогон кончился, так что неудачный прогон прежнего файла не трогает. Ученики идут параллельно, а строки — в порядке популяций, графов и учеников, так что файл одной конфигурации — одни байты.
- **Имя библиотеки** — из сведений о сборке (`runtime/debug`): версия, по которой сборка её взяла, — тег или псевдоверсия, где есть коммит, — или версия её замены, с путём замены, когда это другой модуль. Библиотеку, собранную из рабочей области `go.work`, сборка версией не называет, и её коммит передаёт вызывающий, `-commit <путь модуля>=<полный SHA>`, — рецепт гейта знает его из файла гейта ([architecture.md](architecture.md), раздел 3.9); без него `sim` файла не пишет. Коммит — 40 шестнадцатеричных знаков и только для одной из двух библиотек; коммит библиотеки, чью версию сборка знает, отвергается — в файле было бы два имени одной сборки.

— *Почему:* одна сборка Go не держит двух версий одной библиотеки ([architecture.md](architecture.md), раздел 3.6), так что пара — два файла; JSON Lines читается построчно, и ученик — одна строка. Коммит в файле — то, по чему гейт сверяет, на чём записаны числа.

— *Отвергнуто:* оракул как отдельный прогон, где движок читает правду, — модели для движка не подменить, а прогон с другими выборами сравнивал бы модель и траекторию сразу; потолок движка для задач с рубрикой и уроков — у них нет P*; CSV — у учеников разные наборы сумм; имя библиотеки через `-ldflags -X` — его передавал бы каждый вызывающий даже там, где сборка версию знает.

— *Следствие:* T23 судит кандидатов по `ReadPair` базы, кандидата и потолков; T24 пишет `go.work` и передаёт `-commit`. Значения `E3` у потолка движка не читаются — `E3` идёт к совершенству ([architecture.md](architecture.md), раздел 3.6). — RUN.md, T22; `internal/bench/`, `internal/measure/`, `internal/run/`, `main.go`; [architecture.md](architecture.md), раздел 3.6.

*Поправка 2026-11-08, при T24 (R26), по решению автора:* «ученики идут параллельно» было ошибкой T22. `bench.Run` раздавал учеников по горутинам через `sync.WaitGroup.Go`, а R09 горутины в коде отверг: прогон делится на процессы по ученикам, и делит его рецепт. Правило `goroutine` ловит оператор `go` и `WaitGroup.Go` не поймало. Горутины убраны. Прогон берёт часть учеников (`bench.Part`, `sim run -part k/n`), рецепт запускает части процессами, а `sim merge` сводит их файлы в файл целого прогона, до байта тот же.

## Критерии и отчёт

**R25. Критерии — пакет `internal/criterion` и команда `sim judge`: файл суда называет критерий, прогоны базы, потолка движка, кандидатов и подтверждения; проверка стоит там, где число читает база; путь `E1` в критерии модели идёт к потолку движка; оценка складывает k-е выборки ячеек.** Решено исполнителем по поручению автора 2026-11-08 при планировании T23 (RUN.md, T23; R13).
- **Файл суда** — JSON: критерий (`model`, `engine` или `pair`); файлы результатов базы (правило `engine`) и потолка движка (правило `ceiling`) на одних учениках; кандидаты — имя, файл, своя группа оценки, цель модели (`lag`, `corridor`, `settle`), число изменённых констант, новое поле, удалённость от базы в порядке автора, а в паре — библиотека, чья цель; подтверждение — три прогона на отложенных учениках. Оракул модели — суммы оракула в файле базы (R24). Отчёт — Markdown, вывод инструмента, поэтому по-английски (R14 spec).
- **Направление и допуски.** Меры со знаком — модулем (`M3.lag`, `M8.shift`); доли, где больше — лучше, — недобором до единицы (`E1`, `E3` R*, `E4`); счёт написанных задач `E5` — со сменой знака. Допуски — автора ([architecture.md](architecture.md), раздел 3.6); мерам, которых его список не называет, — допуск ближайшей по роду: доли `M5` и `M7` — 1 пункт, как доли `M3`; сдвиг `M8` — 0,01 логита, как `M1`; Brier `M10` — 0,002, как `M2`; доля трудных первых пунктов `M7` — 1 пункт. Разрешение стенда — 4,3 стандартной ошибки разности, а стандартная ошибка — отклонение её значений по 2 000 выборкам.
- **Проверки.** «Не хуже» — худший конец интервала разности «кандидат минус база» не больше допуска. Пол — худший конец интервала кандидата: `M4` не больше 0,20, `M6` не больше 1,5 значения базы. Цели — `lag` и `corridor` по нижнему концу интервала доли пути, `settle` по верхнему концу доли не догнавших. Проверка и пол стоят там, где число читает база: где не читает и она, держать нечего, и отчёт называет такие ячейки, а не проваливает кандидата. Число, которое база читает, а прогон кандидата нет, — проверка не выполнена. Инварианты `I1`–`I4` прогон пишет в строку ученика под правилом `engine` (`bench.Line.Breaches`), а заголовок файла говорит, что они записаны; файл без этой пометки инвариантов не держит. `I1` между двумя прогонами — гейта (T24). Пол и цель, чьей популяции в прогонах нет, не выполнены: проверка, которую не мерили, не держится. В паре проверка, общая обоим критериям, стоит один раз. Цель модели у кандидата движка отвергается: числа цели у движка нет.
- **Оценка.** Доля пути меры на ячейке — (база − кандидат)/(база − потолок) в её направлении, путь — по целой ячейке. Потолок — совершенство (ноль) у мер, где путь идёт к совершенству; потолок движка у мер `E`, в критерии модели тоже — оракул через движок не прогнать (R24), а потолок движка — потолок коридора всякой пары модели и движка; оракул у прочих мер модели. Мера без пути — путь не больше нуля, или база или потолок её не читают, — выходит из оценки, и отчёт её называет; мера, которую не читает прогон кандидата, оставляет оценку непрочитанной. Значение популяции — среднее по её ячейкам и мерам, так что графы популяции делят её вес поровну. Интервал — из 2 000 выборок: k-е выборки всех ячеек складываются с весами, и у каждой ячейки свои выборки — свои места.
- **Выбор и выход.** «Лучше» — нижний конец интервала оценки выше нуля, и своей группы тоже. Равные лучшему — те, чья разность оценок с ним по тем же выборкам держит ноль. Простейший — наименьшее число констант, затем без нового поля, затем ближайший к базе, затем по имени. Подтверждение — тот же суд на прогонах отложенных учеников: база и кандидат под `engine`, потолок под `ceiling`, набор отложенный и не тот, на котором выбирали; без него выбор ждёт подтверждения. Решающим отчёт зовёт прогон, только если и подтверждение решающее. Пара — ограничения обоих критериев, «лучше» по оценке своей библиотеки, а верхний конец интервала оценки другой не ниже нуля. Решающий прогон — когда в каждой ячейке каждого прогона не меньше 4 000 учеников.
- **Разрез `M2` после пауз от 7 суток** (`M2.brier.after7`) добавлен для группы «Забывает»: раздел 3.6 называет его, а разрезы `M2` делили паузы на промежутки.

— *Почему:* проверка, которую не читает и база, говорит о графе, а не о кандидате: на школьной математике у 40 учеников «освоено» объявляется так редко, что не каждая выборка его видит, и любой кандидат проваливался бы. Сложение k-х выборок даёт интервал суммы, потому что выборки ячеек независимы.

— *Отвергнуто:* проверка, которая проваливается, когда не читает и база, — она проваливает и базу против самой себя; путь `E1` к оракулу — оракул движка не ведёт; вес популяции на каждый её граф целиком — популяция с двумя графами весила бы вдвое.

— *Следствие:* значения базы и потолка, из которых складывается путь, суд пока читает из своих прогонов: закрепить их поправкой к R04 можно только после первого решающего прогона базы ([architecture.md](architecture.md), раздел 3.6), а он — за T24. — RUN.md, T23; `internal/criterion/`, `internal/bench/`, `internal/measure/`, `main.go`; [architecture.md](architecture.md), раздел 3.6.

— *Поправка 2026-11-09:* суд больше не читает путь из своих прогонов. Он берёт путь из файла путей, который закрепил первый решающий прогон базы (R28; поправка 2026-11-09 к R04), — при каждом выборе и в каждом подтверждении. Поэтому файл суда не называет прогона потолка, ни у выбора, ни у подтверждения, а `Input` и `Confirm` его не несут. Доля кандидата — по-прежнему его парная разность с базой суда, на учениках суда, но делённая на закреплённый путь. Решающий прогон шёл в T26, а не в T24, как говорило «Следствие». — RUN.md, T26; `internal/criterion/`, `main.go`; [architecture.md](architecture.md), раздел 3.6.

— *Поправка 2026-11-11 (T40; R40):* суд отказывает прогон не из своего суда. База и база подтверждения — на библиотеках, на которых прочитаны пути, а база — на версиях встроенного `gate.json`. Вторая библиотека кандидата одной библиотеки — версия базы или коммит её закрепления.

## Гейт

**R26. Гейт — файл `gate.json` в корне `sim`, рецепт `just gate`, команда `sim gate` и job «Gate» на amd64 и arm64. go.mod называет библиотеки той же точкой, что и файл. До T26 гейт держит полосы своего прогона на базе. `I1` — повтор каждого ученика гейта.** Решено автором 2026-11-08 при планировании T24 (RUN.md, T24) — горутины, go.mod, полосы до T26 и `I1`, — а вид файла и сверки — исполнителем.
- **Файл гейта** — JSON:
  - `libraries` — по пути модуля: репозиторий, полный SHA коммита и тег релиза на нём, пустой, пока релиза нет;
  - `modules` — релиз, дайджест индекса и пакеты (модуль, версия, дайджест пакета);
  - `dependencies` — выбранные версии общих зависимостей, сейчас модуль Go `spec`;
  - `go` — строка `go` каждого из трёх модулей;
  - `bands` — причина записи и полоса каждого числа: ячейка, число, прочитано ли, значение и полуширина.

  Полосы лежат в том же файле, как велит [architecture.md](architecture.md), раздел 3.9. Число, которое запись не прочла, записывается «не прочитано», и гейт держит именно это: прочитанное позже — вне полосы. Раздел 3.7 зовёт непрочитанное «вне полосы» для числа, которое при записи читалось.
- **Рецепт `just gate [bands|exact|record]`:**
  - берёт каждую библиотеку коммитом из файла во временную папку (`git fetch` по SHA), а вызывающий подставляет свой код папкой — `GATE_LEARNER_MODEL` или `GATE_ENGINE`;
  - сверяет тег из файла с коммитом из файла (`git ls-remote`): job «Go» собирает тег `go.mod`, а гейт — коммит;
  - пишет `go.work` с `use` трёх модулей, ставит `GOWORK` и собирает `sim` один раз;
  - гонит гейт процессами, по части учеников на ядро (`sim gate run -part k/n`), с коммитами, которые видит в папках;
  - затем `sim gate check` или `sim gate record` сводит части. Часть, которая упала, останавливает остальные. Причина записи идёт в команду в кавычках оболочки.
- **Файл результатов** называет и популяции, учеников на граф и повтор прогона, так что `sim merge` сводит части одного прогона в его порядке и отвергает части другого. Запись полос отвергает прогон не на коммитах файла: полосы кода вызывающего — не запись.
- **Сверки** — каждое расхождение даёт красный гейт с причиной:
  - библиотеки собраны из рабочей области;
  - выбранная версия `spec` (`debug.ReadBuildInfo`) равна записанной;
  - строки `go` трёх модулей равны записанным и не новее Go, которым собран гейт;
  - go.mod `sim` требует каждую библиотеку тегом из файла, а пока релиза нет — псевдоверсией его коммита;
  - пакеты, которые встраивает стенд, — те, что называет файл;
  - ни один ученик не нарушил `I1`–`I4`;
  - каждое число в полосе. Свой CI `sim` держит числа записанными до бита (`exact`), вызывающие — полосами.
- **Ячейки и числа** — по разделу 3.7, в коде (`internal/gate`): первые 300 рабочих учеников каждой ячейки популяций гейта, числа без направления. Полуширина — половина 95 % интервала из тех же 2 000 выборок. `A-linked` идёт только на синтетических пакетах (R19), а гейт синтетического модуля не берёт ([architecture.md](architecture.md), раздел 3.3). Поэтому `A-linked` гейт не гоняет, а его отчёт называет её числа «нет графа». Прогон с синтетическим модулем — красный.
- **Полосы до T26** записаны прогоном гейта на `learner-model` v0.1.0 и `engine` v0.1.0, с причиной «первая запись на базе». T26 записывает их заново со своей причиной.
- **`I1`** — каждый ученик гейта прогоняется дважды (`bench.Config.Repeat`), и `invariant.Same` сверяет два прогона до бита. Если замер на раннере не уложится в бюджет, повтор урезается поправкой к этой записи.

— *Почему:*
- Процессы держат R09 и делят прогон на ядра без горутин.
- go.mod на той же точке, что и файл, — единственный способ собрать job «Go» общего `go.yml` на коммите, которого ещё нет в релизе. Прокси Go отдаёт коммит по SHA, и job «Go» собирает ровно тот код, на котором записаны полосы.
- Полосы своего прогона дают гейту что держать уже в T24, а не только с T26.

— *Отвергнуто:*
- полосы в отдельном файле — раздел 3.9 кладёт их в файл гейта;
- горутины с поправкой к R09 — автор оставил R09;
- go.mod только нижней границей — job «Go» был бы красным на шагах 1–4 порядка релизов;
- повтор части учеников сразу — замер ещё не сказал, что полный повтор не укладывается;
- `jq` в рецепте — его нет в стадии `toolchain`, а `sim gate libraries` читает файл тем же кодом, что и гейт.

— *Следствие:* вызывающие берут тот же рецепт с `GATE_*` — это задачи их репозиториев. Порог CI `modules` — T25. Решающий прогон делится на части тем же `-part` и сводится `sim merge` — T26. — RUN.md, T24; `gate.json`, `internal/gate/`, `internal/bench/`, `internal/published/`, `gate.go`, `main.go`, `justfile`, `.github/workflows/ci.yml`; [architecture.md](architecture.md), разделы 3.7 и 3.9.

— *Поправка 2026-11-09:* в гейте 100 учеников на ячейку, а не 300, и при повторе прогона для `I1` выбор на каждом шаге заново берётся только у первых 10 учеников ячейки. На `olympiad-math` 0.1.1 выбор движка идёт в десятки раз дольше, и гейт на 300 учениках не укладывался в `timeout-minutes: 15`. Гейт урезан правкой `sim`, а не временем, как и велит раздел 3.9 [architecture.md](architecture.md). `A-linked` больше не «нет графа»: он идёт на паре опубликованных математик. Почему — R31.

— *Поправка 2026-11-10:* файл гейта называет библиотеки v0.1.4 — `engine` на ddfed3a, `learner-model` на acfc4f2, — пакеты `modules` v0.2.0 и `spec` v0.12.0. `just gate exact` сошёлся на них до бита, так что полосы прежние. Почему — R32.

— *Поправка 2026-11-10:* рабочее пространство гейта пишет скрытый рецепт `workspace`, один для гейта и для решающего прогона кандидата (R34). Гейт берёт из `GATE_LEARNER_MODEL` и `GATE_ENGINE` не только папку, но и коммит по полному SHA. Проверка тега, сам гейт и его числа прежние: `just gate exact` на новом рецепте сходится до бита.

— *Поправка 2026-11-11 (T41):* файл гейта называет `engine` v0.1.6 на 57e0ab8; `learner-model` — по-прежнему v0.1.4 на acfc4f2, пакеты — `modules` v0.2.0, `spec` — v0.12.0. На 100 учениках `just gate exact` сошёлся на нём до бита, затем гейт вырос до 150 учеников, и полосы записаны заново `just gate record` с причиной (поправка 2026-11-11 к R31).

— *Поправка 2026-11-11 (T43):* полосы записаны заново на `modules` v0.2.5 — `just gate record` с причиной T43, на 150 учениках, а `just gate exact` зелёный. Все 69 чисел прочитаны там же, где прежде; 39 не изменились ни в одном бите — 26 школьных, 6 олимпиадных, 3 пары математик и все 4 сочинений. Сдвинулись 23 олимпиадных, 6 школьных — у `G0`, `G2`, `G2-half` и `F-model`, в третьем–четвёртом знаке: на время уроков легло звено разобранного примера — и `M9.notFalse` у `A-linked`. Сильнее всех сдвинулась доля трудных первых пунктов `E1.hard` на олимпиадной — на 1,9 прежней полуширины у `G0` (0,320 → 0,374) и на 1,5 у `G2-half`, а `E2.stuck` стало ненулевым у `G0`, `G1` и `G2` олимпиадной (до 0,067): 30 выборов из пяти стали временем, и угадывания у них больше нет. Почему — R42.

— *Поправка 2026-11-11 (T45):* файл гейта называет `learner-model` v0.1.6 на 7d796e2. `just gate exact` сошёлся на нём до бита, все 69 чисел, — и полосы записаны заново `just gate record` только ради причины, как в T42: числа и полуширины — те же. Почему — R43.

— *Поправка 2026-11-12 (T46):* гейт идёт на 250 учениках на ячейку, а не на 150, и полосы записаны заново `just gate record` с причиной T46 на тех же библиотеках и пакетах; `just gate exact` зелёный — 1 мин 52 с по часам и 29,6 минуты процессора в контейнере, при load average около 40–60 от соседних контейнеров. Из 69 чисел 25 — с нулевой полосой — те же. Три числа, которые на 150 учениках запись не прочла, теперь прочитаны: ложные «освоено» `M4.false.all` и опоздание `M6.late` у `G0` школьной и доля опоздавших повторений `E3.late` у `F-model` школьной; непрочитанными остались два — `M4.false.pages` у `W-shock` и `M9.hasFalse` у `A-linked`. Из 42 чисел с ненулевой полосой три — те, что прочитаны впервые, — прежней полуширины не имеют; у 38 из остальных 39 новое значение отстоит от прежнего меньше чем на одну прежнюю полуширину — больше всех `E5.written` у `E-base`, на 0,60, — а полуширина — в среднем 0,79 прежней, около √(150/250) ≈ 0,77. Тридцать девятое — доля пауз, после которых модель так и не догоняет забывание, `M5.never` у `F-model` школьной: на 150 оно было 0 с нулевой полосой, а стало 0,0036 ± 0,0055: среди 250 учеников нашлась пауза, после которой модель не догнала забывание, а среди первых 150 такой не было. Это редкое событие, и прежний ноль был ноль выборки, а не меры. Почему гейт вырос — поправка 2026-11-12 к R31.

**R27. Гейт CI `modules` — режим `just gate modules` и две команды: `sim gate changed` решает, изменилось ли то, что стенд читает, по каноническому тексту переведённого модуля, а `sim gate thresholds` держит инварианты и пороги. Числа гейта пополнены `E1.hard` и `E2.nothing`, от которых берутся пороги.** Решено исполнителем при T25 (RUN.md, T25), кроме порогов `E1`: их решил автор — поправка 2026-11-09 к R05.
- **Что сравнивается.** Канонический текст — JSON переведённого модуля, `pack.Module.Text`: полосы и аудитория манифеста, граф и ярусы, задачи с пунктами, слотами, ловушками, канвой и возрастом, числа экземпляров шаблонов, уроки. Это ровно то, что берёт `pack.Translate`, так что список того, что стенд читает, — одно место в коде, как велит [architecture.md](architecture.md), раздел 3.9. Тексты для ученика и версия модуля в текст не входят.
- **С чем сравнивается** — с новейшей опубликованной версией модуля, которая не снята: снятую `api` не выдаёт (R63 spec). Состояния:
  - `new` — опубликованной версии нет, или релиза нет вовсе (`-first`); модуль изменился весь;
  - `changed` и `same` — тексты разные или одинаковые;
  - `transition` — у опубликованной версии другая версия формата; стенд модуль не гонит, даже с `-all` (поправка 2026-11-07 к R103 spec);
  - пакет версии формата, которую стенд не читает, без перехода — отказ: нужен новый тег `sim`;
  - сломанный пакет — часть для сервера не та, что в описи, или манифест не тот, — отказ, и в переходе тоже: опись держит `server_part` во всех версиях формата ([текст формата](../spec/format/README.md), раздел 2).
- **Проверка пакетов.** Опубликованный — по цепочке дайджестов от индекса, который принял `release-verify`: `inventory.json` равен дайджесту пакета в записи, `server.json` — описи. Пакет pull request записи не имеет: `server.json` — против описи, манифест — против описи. Файлы читает `main`, а `pack` получает байты (`pack.ReadFiles`).
- **Что гоняется** — ячейки гейта, у которых каждый модуль графа есть среди пакетов, которые идут (`bench.Config.Packages`). Заголовок файла результатов называет эти пакеты: модуль, версию и дайджест. `A-linked` по-прежнему «нет графа». Модуль, на котором не идёт ни одна популяция гейта, отчёт называет, но гейт от этого не краснеет: популяции на новом графе — правка `sim`.
- **Пороги** — точечные значения: ученики гейта всегда те же. Края включительно, непрочитанное число — красное. Значения для порогов `E1` и `E2` берутся из полос `gate.json` того же тега:
  - `E2.stuck` и `E2.nothing` — наибольшее по всем ячейкам плюс 1 пункт;
  - `E1.hard` и `E1.corridor` — по ячейкам `G0` и `G2-half`, по поправке 2026-11-09 к R05;
  - `M4.false.all` на `G0` — не больше 0,20, `E5.again` на `E-base` — ноль.
- **Сверка прогона.** Библиотеки — коммиты файла. Есть все ячейки пакетов, со всеми учениками и повтором для `I1`. Сверки сборки — те же, что у `sim gate check`. `sim gate check` и `sim gate record` прогон на подменённых пакетах отвергают.
- **Числа гейта** пополнены: `E2.nothing` — на всех популяциях, `E1.hard` — на `G0` и `G2-half`. `gate.json` записан заново с причиной T25. 47 прежних полос не изменились ни в одном бите, добавилось 18.
- **Рецепт** — `just gate modules <аргументы sim gate changed>`. Сначала `sim gate changed`: если гнать нечего, рецепт выходит с 0, не трогая сеть. Иначе — та же сборка через `go.work` из файла, части `sim gate run -package …` и `sim gate thresholds`. Относительные пути — от папки вызывающего. Форма вызова для job `modules` — README, «Gate in the CI of modules».

— *Почему:*
- Текст из `Translate` не может разойтись с тем, что стенд на самом деле читает: поле, которое `Translate` начнёт читать, само войдёт в сравнение.
- Аргументы — как у `module-kit diff` (R11 modules): CI `modules` уже держит индекс, принятый `release-verify`, и опубликованные пакеты папками.
- Точка, а не интервал: у гейта нет случайности между прогонами, а интервал сдвинул бы порог на шум, которого нет.

— *Отвергнуто:*
- проекция канонического вида по списку полей — второй список рядом с `Translate`, который разойдётся с ним;
- сравнение со снятой версией — её не выдают, и сравнение с ней сказало бы «same», когда ученики получают другое;
- гейт по новым модулям красным, пока на них нет популяции, — `modules` не мог бы добавить модуль, пока `sim` его не закрепит, а закрепляет он только выпущенное.

— *Следствие:* job гейта в CI `modules` подключит сессия `modules` по README — Д5 в таблице 4 RUN.md `modules`, после первого тега `sim` (T27). — RUN.md, T25; `internal/gate/modules.go`, `internal/gate/gate.go`, `internal/pack/`, `internal/bench/`, `gate.go`, `justfile`, `gate.json`, README; [architecture.md](architecture.md), разделы 3.7 и 3.9.

— *Поправка 2026-11-09:* `sim gate changed` отдаёт вместе с изменёнными модулями и неизменённый, когда граф популяции гейта держит его вместе с изменённым: так `A-linked` идёт на обеих математиках из pull request, если изменилась хотя бы одна (раздел 3.9 [architecture.md](architecture.md)). На stdout его папка идёт строкой `along <папка>`, рецепт отдаёт её `sim gate run -along`, и такой пакет идёт только в ячейках, которые он делит с изменённым (`published.Package.Along`), а не во всех своих. Строка такого модуля на stderr называет модуль, который его привёл. Ячейку гейта, которой не хватило пакета второго модуля её графа, отчёт называет, не краснея. В коде — `gate.Runs` и `partial`. Почему — R31.

— *Поправка 2026-11-11:* перевод читает, есть ли у задачи наводящий вопрос, и `time` — видом проверки, так что канонический текст держит и то, и другое: `sim gate changed` видит добавленный или снятый вопрос и проверку `number`, ставшую `time`, но не правку слов вопроса — их стенд не читает. Почему — R41.

**R28. Пути целей закреплены файлом `internal/criterion/paths.json`, встроенным в стенд: его пишет `sim paths` из решающих прогонов базы и потолка движка, и суд берёт путь из него. Решающий прогон гонит рецепт `just decisive` на машине разработки, а файлы его результатов в репозитории не лежат. Первый файл полос записан на коммитах релизов базы.** Решено автором 2026-11-09 при планировании T26 (RUN.md, T26) — файл путей, файлы прогона, где гонять и полная таблица в поправке к R04, — остальное исполнителем.
- **Файл путей** — JSON:
  - `libraries`, `set`, `learners` — на чём прочитаны пути;
  - `base` и `ceiling` — SHA-256 файлов результатов двух прогонов;
  - `paths` — ячейка, мера, прочитан ли путь, значение базы и значение потолка, оба в направлении критерия, так что путь — база минус потолок.

  Пары — каждая мера групп оценки обоих критериев на каждой ячейке опубликованного графа, где у популяции есть вес. Путь `E1` у двух критериев общий. Потолок — ноль у мер «до совершенства», потолок движка у мер `E…`, оракул базы у мер `M…`, как и раньше считал суд.
- **`sim paths -base <файл> -ceiling <файл> [-out …] [-table]`** пишет файл путей и печатает строки таблицы для поправки к R04. Отказ, если:
  - прогоны не под правилами `engine` и `ceiling`;
  - набор не рабочий;
  - дана часть прогона;
  - прогон шёл на пакетах вместо закреплённых;
  - в базе нет какой-то популяции групп оценки;
  - у прогонов разные библиотеки или ученики;
  - в базе есть нарушение инварианта;
  - в ячейке меньше 4 000 учеников.
- **Суд** берёт путь из встроенного файла: доля кандидата — его парная разность с базой суда, делённая на закреплённый путь. Меры, которой нет в файле, которую он не прочёл или у которой путь не больше нуля, в оценке нет, и отчёт её называет. Без закреплённых путей суд отказывает. Прогон потолка суду больше не нужен — поправка 2026-11-09 к R25.
- **Тест** держит таблицу поправки 2026-11-09 к R04 равной `Table()` встроенного файла: значение в журнале и значение, которым судит код, — одно.
- **Рецепт `just decisive engine|ceiling <файл>`** собирает стенд по `go.mod`, где обе библиотеки — v0.1.0. Он гонит `sim run -populations all -learners 4000 -set working` процессами по частям учеников на всех ядрах (R09), затем `sim merge` и SHA-256 файла.
- **Первый файл полос** — `just gate record` с причиной T26, на тех же коммитах v0.1.0: числа те же до бита, что записали T24 и T25, меняется причина.

— *Почему:*
- Путь, который суд читает из своих прогонов, после первого выбора шёл бы от новой базы, а в подтверждении — от отложенных учеников, а раздел 3.6 судит «при закреплённом пути».
- Файл в коде, а не только в журнале: суд не может разойтись с записью, а тест держит запись равной файлу.
- Файлы прогона — около 2 ГБ на правило. Зёрна и одни биты дают их заново, а SHA-256 в файле путей говорит, тот ли это прогон.

— *Отвергнуто:*
- путь из прогонов суда, как было, — он плыл бы;
- хранить файлы прогона артефактом или файлом релиза — место, срок и проверка подписи ради того, что воспроизводится;
- прогон на раннере GitHub — 4 ядра, около 2,5 часа, артефакт в гигабайты.

— *Следствие:* первый кандидат судится по закреплённым путям. Новый закрепляющий прогон — только поправкой к R04, как и новый отложенный набор. — RUN.md, T26; `internal/criterion/`, `main.go`, `justfile`, `gate.json`; [architecture.md](architecture.md), разделы 3.6–3.8.

— *Поправка 2026-11-09:* пути закреплены заново, до первого кандидата: изменился граф — `sim` закрепил `olympiad-math` 0.1.1, а популяции `A-` пошли на опубликованной паре математик. Это не новый выбор, а та же калибровка на графах, которые получат ученики. Прогоны и новая таблица путей — поправка 2026-11-09 к R04, «новое закрепление»; почему — R31.

— *Поправка 2026-11-10:* решающий прогон — 1 000 учеников на ячейку (R33), так что `sim paths` отказывает ячейке, в которой меньше 1 000 учеников, а не 4 000. Рецепт `just decisive` гонит `-learners 1000` и создаёт папку файла, если её нет. Файлы прогона — около 0,5 ГБ на правило. Между судами они лежат в `runs/` рабочей копии: `git` её не видит, а Rebuild Container её не стирает, как стёр файлы T30 в `/tmp`. Пути закреплены заново — поправка 2026-11-10 к R04, «прогоны не дольше часа». Воспроизводимость сверена полным повтором потолка: на 1 000 учениках это дёшево.

— *Поправка 2026-11-10, прогон кандидата (T31; R34):* рецепт `just decisive engine|ceiling <файл> [набор]` берёт набор учеников — `working` или `held-out-<критерий>-<n>` — и кандидата: `CANDIDATE_LEARNER_MODEL`, `CANDIDATE_ENGINE` или оба. Без кандидата он собирает стенд по `go.mod`, как прежде, так что файл базы — тот, чей SHA-256 называет `paths.json`.

— *Поправка 2026-11-11 (T41; поправка 2026-11-11 к R33):* решающий прогон — 1 500 учеников на ячейку, так что `sim paths` отказывает ячейке, в которой меньше 1 500 учеников, а суд называет такой прогон грубым взглядом. Рецепт `just decisive` гонит `-learners 1500`. Файлы прогона — около 0,6–0,8 ГБ на правило. Пути закреплены заново — поправка 2026-11-11 к R04. Воспроизводимость сверена полным повтором потолка. Прежние файлы на 1 000 учеников лежат в `runs/` как `base-1000.jsonl` и `ceiling-1000.jsonl`, и файлы судов `engine-1` и `model-1` называют базу так.

— *Поправка 2026-11-11 (T43):* пути закреплены заново после смены пакетов и библиотеки, а не после кандидата: закрепление `modules` v0.2.5 меняет граф и задачи, а база с T42 — на `learner-model` v0.1.5. `paths.json` называет `learner-model` v0.1.5 и `engine` v0.1.6, таблица — поправка 2026-11-11 к R04, «пакеты с `time` и наводящими вопросами». Воспроизводимость сверена полным повтором потолка. Почему — R42.

— *Поправка 2026-11-11 (T45):* пути закреплены заново после смены библиотеки, хотя числа не сдвинулись: суд держит базу на библиотеках путей (R40), а база — на `learner-model` v0.1.6. `paths.json` называет `learner-model` v0.1.6 и `engine` v0.1.6 и новые SHA-256 прогонов, а значения — те же, что у T43, так что таблица поправки 2026-11-11 к R04, «пакеты с `time` и наводящими вопросами», по-прежнему таблица путей. Почему — R43.

— *Поправка 2026-11-12 (T46; поправка 2026-11-12 к R33):* решающий прогон — 2 000 учеников на ячейку, так что `sim paths` отказывает ячейке, в которой меньше 2 000 учеников, а суд называет такой прогон грубым взглядом. Рецепт `just decisive` гонит `-learners 2000`. Файлы прогона — около 0,8–1,0 ГБ на правило. Пути закреплены заново — поправка 2026-11-12 к R04. Воспроизводимость сверена полным повтором потолка. Прежние файлы на 1 500 учеников лежат в `runs/` как `base-1500.jsonl` и `ceiling-1500.jsonl`; ни один суд на них не ссылается.

**R29. Релиз `sim` — общий `release.yml` `.github` v0.1.21 по линии `0.1`, начиная с `v0.1.0`, когда коммит `main` с зелёным CI меняет то, из чего собран гейт. Единственный файл релиза — архив `sim-<версия>.tar.gz` с тем, что собирает стенд и гоняет гейт, подписанный без ключа. Неизменяемые релизы и ruleset тегов `v*` включены до первого релиза.** Решено автором 2026-11-09 при планировании T27 (RUN.md, T27): общий релиз с архивом (поправка 2026-11-09 к R05), настройки GitHub — исполнителем через `gh api`, мерж, который выпускает тег, — автором. Остальное — исполнитель.
- **`paths`** — то, что даёт релиз:
  - код Go без тестов и `testdata/`;
  - встроенные пакеты `internal/published/v*` и пути `internal/criterion/paths.json`;
  - `go.mod`, `go.sum`, `LICENSE`, `.gitattributes`;
  - файл гейта `gate.json`, `justfile` и `.devcontainer/Dockerfile`, на которых гейт гонится;
  - сам `release.yml`.

  Тесты, документы и прочие workflow релиза не дают.
- **Архив** — рецепт `release-files <версия>`: папка `sim-<версия>` с теми же путями, кроме `release.yml`. Tar пишется одинаково на любой сборке: имена по порядку, без времён и владельцев, режим 0644, `gzip -9n`. Список берётся из дерева, `find`: общий релиз гонит рецепт в дереве `git archive`, без git.
- **`ci-release`** держит архив, как R15 engine:
  - две сборки подряд — те же байты;
  - в архиве ровно то, что git отслеживает по этим путям, и в том же виде;
  - распакованный архив собирается.

  Его зовёт `ci-lint`, так что архив держит «Go / lint».
- **Настройки GitHub** (2026-11-09, `gh api`):
  - неизменяемые релизы включены;
  - ruleset `v-tags` (id 26826001) запрещает update, deletion и non_fast_forward тегов `refs/tags/v*`, без обходов, — такой же, как у `engine`.

— *Почему:*
- Общий релиз — один код релиза на все репозитории с кодом, и Dependabot поднимает его вместе с остальными workflow `.github`.
- Архив — это и есть то, на чём записаны числа гейта, кроме библиотек, которые называет `gate.json`. Подпись привязывает его к тегу (R77 spec).
- Изменение документов или тестов не даёт вызывающим нового тега.

— *Отвергнуто:*
- свой workflow без файлов, как у `spec`, — свой код релиза, который держит один `sim`;
- архив только кода Go, как у `engine`, — гейт гонится и на `gate.json`, `justfile` и стадии `toolchain`, а без них архив не держит того, что числа записаны на этом теге.

— *Следствие:* первый тег `v0.1.0` выходит мержем pull request T27. Вызывающие берут его тег и SHA — README, «Taking a release» — своими задачами: Д5 RUN.md `modules`, job гейта `learner-model` и `engine`. — RUN.md, T27; `.github/workflows/release.yml`, `justfile`, `.gitignore`, README; [architecture.md](architecture.md), раздел 3.9.

## Эпик E3

**R30. Эпик E3, первые кандидаты, идёт до первых выборов. В нём: малый путь в суде, новое закрепление пакетов с новой калибровкой путей до первого кандидата, решающий прогон на коммите кандидата, суд первого кандидата модели и первого кандидата движка, перевод на релизы выбранных и итог.** Это ответы автора 2026-11-09 при планировании E3 (RUN.md, таблица 1). Порядок задач T29–T35 — предложение планирования.
- **Что выяснилось при планировании** (прочитано 2026-11-09):
  - `modules` v0.1.1 опубликовал 2026-11-08 `olympiad-math` 0.1.1 — с новыми навыками и сотнями задач. Индекс новейшего релиза, v0.1.4, называет обе версии. `sim` закрепляет `modules` v0.1.0, так что база, пути и полосы сняты на `olympiad-math` 0.1.0;
  - `just decisive` собирает стенд по `go.mod` и гонит только рабочих учеников, так что кандидата подставить некуда. Отложенных учеников `sim run -set` уже берёт;
  - библиотеки подключают гейт — T22 RUN.md `learner-model` и T23 RUN.md `engine`, — но кандидатов ещё не предложили. Их основа — Д5–Д7 RUN.md `learner-model`, Д3 и Д4 RUN.md `engine`;
  - `spec` вышел до v0.9.16, но модуль Go `spec` библиотеки не требуют, так что гейту переход не нужен.
- **Решено автором:**
  - малый путь — вариант (а) В7: поправка 2026-11-09 к R04, код — T29;
  - новое закрепление — до кандидатов: индекс `modules` v0.1.4 с `olympiad-math` 0.1.1, решающий прогон базы и потолка заново, новые пути и полосы, следующий тег `sim` (T30);
  - эпик идёт до первых выборов: суд первого кандидата модели (T32) и движка (T33), перевод на их релизы (T34) и итог (T35).
- **Предложение планирования:**
  - порядок — малый путь (T29), закрепление (T30), рецепт прогона кандидата (T31, после T30: прогон базы рецептом совпадает до байта с тем, чьи SHA-256 называет новый файл путей), суды (T32, T33), релизы (T34), итог (T35);
  - если кандидаты придут парой, их судит одна задача, T32, по критерию пары. Зерно её отложенных учеников записывается поправкой к R04 до прогона.
  - если одна библиотека кандидата не предложила, когда суд другой сделан, автор откладывает её суд, и перевод на релизы и итог идут без него.

— *Почему:*
- **Закрепление до кандидатов.** Пути — калибровка, а не выбор, и читаются до первого кандидата (R04). Кандидат, выбранный на графе, которого ученики не увидят, выбран не на том. Закрепить после первых выборов значило бы калибровать заново, когда на путях уже судили.
- **До первых выборов.** Ради них стенд и строился: первая версия каждой библиотеки — база, а то, что её улучшает, проходит через суд (R103 spec). Суды ждут кандидатов, но это работа `sim`, и в эпике она видна вместе с тем, что её готовит.

— *Отвергнуто:*
- судить первых кандидатов на `olympiad-math` 0.1.0, а 0.1.1 закрепить после — см. «Почему»;
- закрепить 0.1.1, а пути олимпиадных ячеек оставить от 0.1.0 — это были бы пути другого графа;
- кончить эпик готовностью — малым путём, закреплением и рецептом, — а суд каждого кандидата вести задачей вне эпика. Тогда первые выборы не вошли бы ни в один эпик, и их порядок никто не видел бы заранее.

— *Следствие:*
- E2 перенесён в архив [docs/plan/E2.md](plan/E2.md);
- таблица 1 RUN.md — E3 задачами T29–T35;
- строка E3 ушла из таблицы 2, В7 — из таблицы 3;
- В9 «в реестр spec» говорит `spec`, что эпик развёрнут и кандидатов судят после тега с новым закреплением.

Где записано: RUN.md, таблицы 1–3 и тела T29–T35; [architecture.md](architecture.md), разделы 2.4 и 3.

## Новое закрепление пакетов

**R31. Стенд закреплён на релизе `modules` v0.1.4 — `olympiad-math` 0.1.1, `grade-3-math` 0.1.0 и `english-essays` 0.1.0. Популяции `A-` идут на паре опубликованных математик. Гейт держит 100 учеников на ячейку и при повторе прогона берёт выбор на каждом шаге заново только у первых 10 учеников ячейки.** Это ответы автора 2026-11-09 при планировании T30 и E3 (R30).
- **Что выяснилось при планировании** (2026-11-09, замеры на копии репозитория; источники — [architecture.md](architecture.md), раздел 4):
  - `olympiad-math` 0.1.1 — 19 навыков вместо 3 и 772 задачи вместо 27, у 19 из них экземпляры шаблона. Видов ответа стало четыре, а не только `number`. Ячеек `predicate` по-прежнему нет;
  - два олимпиадных навыка опираются на навыки, которым учит школьная математика: `fractions-parts` — на `core.fractions-of-a-whole`, `number-divisibility` — на `core.division-with-remainder`. Так у популяций `A-` появился опубликованный граф;
  - стенд читает пакет без правок кода; `G0` и `A-` на паре идут без нарушений инвариантов, `M9` читается;
  - ученик олимпиадного графа идёт около 0,65 с вместо около 0,007 с. В профиле:
    - `engine.Choose` на пуле из 772 задач — около 60 % процессора;
    - из них около 28 % — выбор, который `invariant.Check` берёт заново на каждом шаге;
    - сборка мусора — около 23 %;
  - гейт на 300 учениках шёл 8,5 минуты и 135 минут процессора на 20 ядрах — против 45 секунд и 11 минут прежде. На 4 ядрах раннера это около 35–40 минут при бюджете `timeout-minutes: 15`;
  - `sim gate changed` сравнивает с индексом, который передаёт CI `modules`, так что закрепление его не касается. Но `A-linked` идёт на обеих математиках, а команда отдавала только изменённые модули.
- **Решено автором:**
  - популяции `A-` — на опубликованной паре, в этой же калибровке;
  - гейт урезан правкой `sim`, а бюджет остаётся 15 минут: 100 учеников на ячейку вместо 300, а при повторе прогона выбор на каждом шаге заново не берётся. Исполнитель по ревью оставил выбор заново у первых 10 учеников каждой ячейки — около 3 % работы: повтор не видит правило, которое меняет свой вход на месте одинаково в обоих прогонах;
  - индекс `modules` v0.1.4 — ещё при планировании E3 (R30).
- **Предложение исполнения:**
  - архив `olympiad-math` 0.1.1 лежит в папке `v0.1.1/`, где его называет запись индекса, а неиспользуемые индекс и архив `olympiad-math` 0.1.0 из встроенных файлов убраны;
  - `gate.Runs` отдаёт вместе с изменённым модулем и неизменённый, который граф популяции гейта держит вместе с ним;
  - тесты, которые гоняют олимпиадный граф, идут на меньшем числе учеников, чтобы `ci-test` укладывался в бюджет сборок `same-bits`;
  - тест потолка пускает ту же задачу снова после паузы в 24 часа, как и правило (R82 spec): на 27 задачах этого не случалось.

— *Почему:*
- **`A-` в этой калибровке.** Иначе пути группы «Опора» и полосы `A-linked` пришлось бы закреплять ещё раз — уже после первых кандидатов.
- **100 учеников.** Раздел 3.9 [architecture.md](architecture.md) велит урезать гейт, который не укладывается, правкой `sim`, а не временем. Полосы станут шире примерно в √3, около 1,7 раза. Код, который ничего не сдвинул, по-прежнему читает те же числа до бита, а кандидатов судит решающий прогон. Гейт теряет только чувствительность к малому ненамеренному сдвигу.
- **Выбор заново при повторе — только у первых 10.** `I1` до бита держит повтор (`Same`), а отказ движка от своего входа ловит сам прогон. Выбор заново на каждом шаге у всех учеников стоил бы трети работы. Но повтор не видит правило, которое меняет вход на месте одинаково в обоих прогонах, а выбор заново видит, и на первых 10 учениках каждой ячейки он стоит около 3 %.

— *Отвергнуто:*
- поднять бюджет до 45 минут и оставить 300 учеников — тогда каждый pull request трёх репозиториев ждал бы гейт около 40 минут;
- ждать, пока `engine` ускорит выбор на большом пуле, — первые кандидаты ждали бы чужой задачи. Вопрос ушёл строкой В10 таблицы 3 RUN.md;
- оставить `A-` на синтетике — см. «Почему».

— *Следствие:* прогоны и новая таблица путей — поправка 2026-11-09 к R04, «новое закрепление»; поправки к R14, R17, R19, R21 и R26–R28; новый тег `sim`, его берут гейты библиотек (В9). Где записано: RUN.md, T30; [architecture.md](architecture.md), разделы 3, 3.2, 3.3, 3.7–3.9 и 4; README.

— *Поправка 2026-11-10:* процесс прогона гейта, как и решающего, идёт на одном процессоре и с `GOGC` 400 (R33): биты те же, и `just gate exact` сходится до бита, а время гейта на раннере падает. Размер гейта прежний — 100 учеников на ячейку и выбор заново у первых 10; вернуть ему 300 — по-прежнему потом, поправкой к этой записи.

— *Поправка 2026-11-10, замер (T35):* то, что время гейта на раннере падает, подтвердилось только на `ubuntu-24.04-arm`: 7 мин 3 с и 6 мин 49 с на pull request #40 и #41, 4 мин 36 с и 4 мин 42 с на #42 и #43. На `ubuntu-24.04` — 11 мин 52 с и 6 мин 54 с до T36, 8 мин 14 с и 8 мин 23 с после: время раннера гуляет сильнее, чем дала правка, и выигрыша там замер не показывает ([architecture.md](architecture.md), раздел 3.9).

— *Поправка 2026-11-11, после ускорения `engine` (T41; R38):* база гейта — `engine` v0.1.6, и в гейте 150 учеников на ячейку, а не 100. Выбор заново при повторе — по-прежнему у первых 10 учеников ячейки, бюджет — по-прежнему `timeout-minutes: 15`. Ответы автора 2026-11-10 и 2026-11-11 при планировании и исполнении T41 (RUN.md, T41).
- **Что выяснилось** (2026-11-10 и 2026-11-11):
  - на `engine` v0.1.6 гейт на 100 учениках читает все 69 чисел те же до бита, что на v0.1.4, и инварианты целы. Полосы на 100 не переписывались;
  - в контейнере разработки, 100 учеников: 1 мин 44 с по часам и 25,9 мин процессора на v0.1.4, 1 мин 5 с и 16,5 мин на v0.1.6 — в 1,57 раза меньше;
  - «Gate» `sim` на 100 учениках, pull request #49: `ubuntu-24.04` — 5 мин 39 с с холодным кешем Go и 4 мин 27 с с тёплым, `ubuntu-24.04-arm` — 3 мин 13 с и 3 мин 7 с. Постоянная часть job — около 35 секунд, прогон со сборкой — до 5 мин на 100 учеников;
  - вызывающие кеша Go не держат, и у них гейт идёт дольше. «Gate» `engine` на коде v0.1.6 и 100 учениках шёл 5 мин 3 с и 7 мин 44 с, шаг гейта — до 7 мин 18 с, при `timeout-minutes: 25` (В27 RUN.md `engine`). «Gate» `learner-model` на прежнем теге `sim` и медленном выборе шёл 8–12 мин при `timeout-minutes: 15`; на v0.1.6 это около 5–7,6 мин на 100 учеников;
  - в пересчёте по самому медленному замеру на 150 учеников: у `sim` — около 8 мин, у `engine` — около 11,5, у `learner-model` — около 11,2. На 200: у `sim` — около 10,6, но у `learner-model` — до 15, впритык к тайм-ауту. На 250 у `sim` — около 13,1;
  - гейт на 150 учениках в контейнере разработки — 1 мин 28 с и 23,6 мин процессора. Все 69 чисел прочитаны там же, где на 100. Каждое новое значение отстоит от прежнего меньше чем на одну прежнюю полуширину, больше всех — `M1` у `G0` школьной, на 0,94. Полуширина — в среднем 0,85 прежней, около 1/√1,5.
  - «Gate» на 150 учениках, pull request #49, второй коммит: `ubuntu-24.04` — 8 мин 21 с с холодным кешем Go и 6 мин 5 с с тёплым, `ubuntu-24.04-arm` — 4 мин 16 с. Пересчёт давал около 8 мин на холодном — замер с ним сошёлся.
- **Решено автором:**
  - размер гейта — по самому медленному замеру: наибольший с шагом 50, от 100 до 300, при котором самый медленный замер, пересчитанный на этот размер, не больше 12 минут — 3 минуты запаса до тайм-аута;
  - 150, а не 200: на 200 гейт `learner-model` без кеша Go подходил бы к своим 15 минутам.

— *Почему:* правило R38 — гейту столько учеников, сколько укладывается в бюджет на раннере. Время раннера гуляет: на одном коде «Gate» шёл и 6 мин 54 с, и 11 мин 52 с. Гейт, упавший по тайм-ауту, держит мерж в трёх репозиториях, поэтому мерой взят самый медленный замер, с запасом. Укладываться он должен у каждого вызывающего, а не только в своём CI с кешем. 150 учеников сужают полосы примерно в √1,5 ≈ 1,22 раза.

— *Отвергнуто:*
- 200 с `timeout-minutes: 25` у `learner-model`, как у `engine`, — каждый pull request `learner-model` ждал бы гейт дольше, и это правка чужого workflow;
- 200 как есть — «Gate» `learner-model` падал бы по тайм-ауту на медленном раннере;
- 300 — около 16 мин даже в своём CI;
- размер по типичному замеру, а не по самому медленному, — медленный раннер ронял бы гейт.

— *Следствие:* `gate.Learners` — 150; полосы записаны заново, `just gate record` с причиной; пороги CI `modules` берутся из новой записи (R27). 300 учеников гейта ждут ускорения `learner-model`: треть выбора теперь — `Readout.Task` (строка В32 RUN.md `engine`, «в реестр spec»). Где записано: RUN.md, T41; [architecture.md](architecture.md), разделы 3.7–3.9; README; `internal/gate/gate.go`, `gate.json`.

— *Поправка 2026-11-12, после ускорения `learner-model` (T46; R43, R44):* в гейте 250 учеников на ячейку, а не 150. Выбор заново при повторе — по-прежнему у первых 10 учеников ячейки, бюджет — по-прежнему `timeout-minutes: 15`, правило размера — то же: наибольший с шагом 50, от 100 до 300, при котором самый медленный замер, пересчитанный на этот размер, не больше 12 минут, — с одним уточнением автора: размер, который проходит правило впритык, меньше чем с минутой до 12, не берётся, и гейт берёт шаг меньше. По одному правилу вышло бы 300. Ответы автора 2026-11-12 при планировании T46 (RUN.md, T46).
- **Что выяснилось** (2026-11-12):
  - «Gate» `sim` на теге v0.1.9 — `learner-model` v0.1.6, `modules` v0.2.5, 150 учеников, — оба раза с холодным кешем Go: pull request #52 — 6 мин 17 с на `ubuntu-24.04` и 3 мин 36 с на `ubuntu-24.04-arm`, push в main — 5 мин 0 с и 3 мин 32 с. В T41 на 150 учениках с холодным кешем было 8 мин 21 с, около 0,75 от этого;
  - вызывающие на v0.1.9 ещё не перешли. «Gate» `learner-model` на `sim` v0.1.8 гонит его уже быстрый код — до 6 мин 12 с. «Gate» `engine` на `sim` v0.1.6 идёт с `learner-model` v0.1.4 — до 8 мин 21 с, около 6 мин 16 с в пересчёте отношением 0,75;
  - в пересчёте самого медленного замера, около 6 мин 17 с, с постоянной частью job около 35 с, как в T41: 200 учеников — около 8 мин 11 с, 250 — около 10 мин 5 с, 300 — около 11 мин 59 с;
  - «Gate» на 250 учениках, pull request #53, с холодным кешем Go: `ubuntu-24.04` — 7 мин 41 с, шаг гейта — 7 мин 3 с; `ubuntu-24.04-arm` — 5 мин 28 с; весь CI зелёный. Повтор job на `ubuntu-24.04` с тёплым кешем — 5 мин 34 с, шаг гейта — 5 мин 6 с. На втором коммите pull request, тоже с тёплым кешем, — 9 мин 38 с, шаг гейта — 9 мин 9 с: раннер попался медленный. У вызывающих кеша Go нет, и на таком раннере их гейт подошёл бы к 12 минутам или чуть за них — около 11 мин 45 с – 13 мин 17 с, смотря как пересчитывать холодный кеш, — при тайм-ауте 15 у `learner-model` и 25 у `engine`. Автор решил оставить 250: самый медленный замер на деле — 9 мин 38 с, а вызывающие покажут своё время на новом теге;
  - в контейнере разработки гейт на 250 учениках — 1 мин 52 с и 29,6 минуты процессора при load average около 40–60 от соседних контейнеров.
- **Решено автором:** 250 учеников.

— *Почему:* по правилу 300 учеников проходят впритык — около 11 мин 59 с при пороге 12, — а время раннера на одном коде гуляло до 1,7 раза (поправка 2026-11-10 «замер»), и вызывающие на новом теге ещё не замерены. 250 оставляют около двух минут сверх трёх минут запаса правила. Полосы на 250 уже, чем на 150, в √(5/3) ≈ 1,29 раза.

— *Отвергнуто:*
- 300 по букве правила — запаса на шум раннера нет;
- 300 с откатом на 250, если замер pull request выйдет за 12 минут, — один удачный замер не защищает вызывающих от тайм-аута.

— *Следствие:* `gate.Learners` — 250; полосы записаны заново, `just gate record` с причиной (поправка 2026-11-12 к R26); пороги CI `modules` берутся из новой записи (R27). Д7 закрыт вместе с поправкой 2026-11-12 к R33. Где записано: RUN.md, T46; [architecture.md](architecture.md), разделы 3.7–3.9; README; `internal/gate/gate.go`, `gate.json`.

## Тег `spec` v0.12.0

**R32. `sim` стоит на `spec` v0.12.0, последнем релизе `spec`. С ним пакеты закреплены на `modules` v0.2.0 — все три модуля в формате 0.4 —, а библиотеки — на `learner-model` v0.1.4 и `engine` v0.1.4. Шаг gopls в образе — заготовки `spec`, R11 заменён.** Поручение автора 2026-11-10 — «перейди на spec последней версии» — и его ответы на три вопроса перехода.
- **Что тег решил за `sim`** (`git -C spec diff v0.9.9 v0.12.0`, прочитано 2026-11-09 и 2026-11-10; источники — [architecture.md](architecture.md), раздел 4):
  - формат 0.4 вышел с v0.10.0 (R102, R108, R111 spec). Пока формат v0, платформа читает только текущую его версию ([текст формата](../spec/format/README.md), раздел 2), так что стенду нужно окно 0.4 и пакеты 0.4. Все дополнения 0.4 необязательны, и модуль 0.3, переведённый на 0.4, даёт те же вердикты и числа. В типах `formatv0` пояснение урока — теперь `ExplanationPart`;
  - Д137 таблицы 4 RUN.md `spec` — ответ на строку `engine` В24: `sim` берёт `engine` и `learner-model` v0.1.3 в новом закреплении, до решающего прогона базы. Они вышли через гейт с теми же числами;
  - Д123 и Д136 — ответы на строки `sim` В6, В8 и В9 по их рекомендациям: строки 4.5 и 4.8 [COORDINATION.md](../spec/COORDINATION.md) закрыты, а первые кандидаты — строка 4.9;
  - поправка 2026-11-09 к R103 spec: кандидатов одной библиотеки `sim` судит кругом — решающий прогон на каждом, один суд со всеми, подтверждение только выбранного (строка `learner-model` В31). Это касается T32 и T33;
  - вид `go` заготовки даёт шаг gopls с lock-файлом (R102 spec), той же версии 0.24.0, что у `sim`;
  - остальное — регион данных, правовые вопросы и данные учеников (R109, R112–R123 spec) — `sim` не касается: данных живых учеников у него нет.
- **Что изменено здесь:**
  - `spec/` — на v0.12.0, коммит 5b4640a; `gh release verify` подтвердил релиз;
  - `go.mod` — `spec` v0.12.0, `engine` и `learner-model` v0.1.4; `release-verify` принял архивы обеих библиотек;
  - окно формата — 0.4; синтетические пакеты и тесты — в 0.4; пояснение урока — `ExplanationPart`;
  - закрепление — `modules` v0.2.0: `release-verify` принял индекс и три архива, все модули версии 0.2.0;
  - `gate.json` — библиотеки v0.1.4 коммитами тегов, `modules` v0.2.0, `spec` v0.12.0. `just gate exact` сошёлся до бита: полосы прежние;
  - шаг gopls и `.devcontainer/gopls/go.mod` — из заготовки; `go.sum` тот же;
  - тег — в CLAUDE.md, README.md и начале RUN.md.
- **Решено автором:**
  - пути пишутся заново в T36, быстрым прогоном, а не здесь: T36 и так меняет устройство прогона, а кандидатов до неё не судят;
  - библиотеки — v0.1.4, последние релизы: код тот же, что в v0.1.3, сверху только документы и переход на `spec` v0.12.0;
  - шаг gopls — заготовки.
- **Выяснилось при переходе:** первые 100 учеников каждой из 65 ячеек на новом стенде дают те же байты, что прогон базы T30, — 6 500 строк. Ни исправление библиотек v0.1.3, ни формат 0.4 чисел базы не сдвинули, так что пути T30 остаются в силе до T36. Но две части 0.4 стенд не прочтёт, когда модуль ими воспользуется: вид проверки `time` и наводящий вопрос `leading_question`. В опубликованных пакетах их нет; это открытый вопрос — [architecture.md](architecture.md), подраздел 2.5, строка В11 таблицы 3 RUN.md.

— *Почему сейчас:* поручение автора. v0.12.0 — последний релиз `spec` на 2026-11-10, на нём уже `modules`, `engine` и `learner-model`. Без него стенд не прочтёт пакетов 0.4, а гейт CI `modules` на `v0.2.0` не пойдёт.

— *Отвергнуто:*
- остаться на v0.9.9 — см. «Почему сейчас»;
- библиотеки v0.1.3 — ровно то, что назвал Д137, но v0.1.4 — тот же код на том же теге `spec`, что у `sim`;
- решающий прогон сейчас — около 8 часов, а после T36 пути всё равно пишутся заново;
- свой шаг gopls — расходился бы с заготовкой, и gopls пришлось бы поднимать руками.

— *Следствие:* таблица 3 RUN.md — В6, В8 и В9 закрыты, В10 остаётся; R11 заменён; поправки 2026-11-10 к R04, R14 и R26. Мерж выпускает следующий тег `sim`, и гейты библиотек берут его. Где записано: RUN.md, T37; [architecture.md](architecture.md), статус и разделы 3, 3.3 и 4; CLAUDE.md, README.md.

## Прогоны не дольше часа

**R33. Решающий прогон и подтверждение идут на 1 000 учениках на ячейку вместо 4 000. Выбор на каждом шаге заново для `I1` берут только первые 10 учеников ячейки — в каждом прогоне, как в гейте. Процесс прогона идёт с `GOGC` 400, а процесс части прогона — ещё и на одном процессоре. Файлы решающих прогонов хранятся в `runs/` рабочей копии, которую `git` не видит, и суд гонит только кандидата.** Решение автора 2026-11-09, при T30: прогон дольше часа недопустим — автор прототипирует, и функциональность меняется быстро. Ответы автора 2026-11-10 при планировании T36 (RUN.md, T36); среда выполнения — предложение исполнения.
- **Что выяснилось при планировании** (2026-11-10, замеры на копии репозитория; источники — [architecture.md](architecture.md), раздел 4):
  - 10 учеников на ячейку во всех 65 ячейках — 1 064 с процессора:
    - половина — шесть популяций `F-`, у которых 180 дней;
    - школьная математика, сочинения и `S-predicate` вместе — меньше 1 %;
    - остальное — олимпиадный граф, 772 задачи;
  - в профиле `G0` и `F-model`:
    - `engine.Choose` — 64 % процессора, из них выбор заново в `invariant.Check` — 31 % всего;
    - сборка мусора — около четверти;
  - на одной выборке, байты результатов во всех вариантах одни:
    - выбор заново только у первых 10 — 121 → 70 с процессора;
    - `GOGC=400` — 121 → 87 с;
    - всё вместе с `GOMAXPROCS=1` — 43 с, в 2,8 раза меньше;
    - потолок с `GOGC=400` и `GOMAXPROCS=1` — 22 → 16 с;
  - память при `GOGC=400` — около 100 МБ на процесс, 20 процессов — около 2 ГБ;
  - после этих правок `engine.Choose` — 73 % процессора, а собственный код стенда — около 17 %, в основном P* пула на каждом шаге для мер. Дальше прогон ускоряет только `engine` (строка В10 таблицы 3 RUN.md). Прогон базы на 4 000 учениках шёл бы около 2,4 часа;
  - файлы решающих прогонов T30 лежали в `/tmp` и пропали после Rebuild Container.
- **Решено автором:**
  - 1 000 учеников на ячейку у решающего прогона и подтверждения;
  - выбор заново для `I1` — у первых 10 учеников ячейки, `bench.Rechosen`;
  - прогон базы хранится в `runs/`.
- **Предложение исполнения:**
  - `sim run` и `sim gate run` ставят себе `GOGC` 400, а прогон по частям — ещё и `GOMAXPROCS` 1, в коде, а не в рецепте, так что их получает всякий вызов. R09 и так держит горутины вне кода, а части прогона — процессы, по одному на процессор. Прогон целиком — один процесс на машине, и процессоры остаются его сборке мусора;
  - популяции решающего прогона — все: «не хуже» читается на каждой популяции с вердиктом (раздел 3.6 [architecture.md](architecture.md)), и без какой-то из них суд не прочёл бы своего ограничения;
  - `invariant.CheckRepeated` переименована в `CheckWithoutI1`: у решающего прогона повтора нет;
  - рецепт `just decisive` создаёт папку файла, если её нет.
- **Прогон** (поправка 2026-11-10 к R04, «прогоны не дольше часа»): база на 1 000 учениках — около 37 минут процессора на процесс, потолок — около 14. На свободной машине каждый прогон укладывается в час; на машине, которую грузили чужие контейнеры, база шла 1 ч 42 мин по часам. Пути близки к путям T30, а мер с малым путём — те же семь.

— *Почему:*
- **1 000 учеников.** Стенд своими средствами ускоряет прогон в 2,8 раза, а нужно было около семи. Остальное даёт только меньшее число учеников, пока `engine` не ускорит выбор. 1 000 оставляют запас на занятую машину и на рост пулов. Разрешение «не хуже» — 4,3 стандартной ошибки парной разности — шире в √4 = 2 раза, и допуск автора там, где он меньше, уступает разрешению. Суд по-прежнему различает то, что сильнее шума, а грубый взгляд — меньше 1 000 — по-прежнему ничего не решает.
- **Выбор заново у первых 10.** Повтор, `Same`, держит `I1` до бита у гейта, а первые 10 учеников каждой ячейки ловят правило, которое меняет свой вход на месте (R31). В решающем прогоне это 650 учеников и около 130 тысяч выборов. Выбор заново у всех стоил около 40 % процессора. Правило, которое выбирает не одно и то же на одном входе, — например, по порядку обхода карты, — решающий прогон без повтора у остальных учеников не увидит. Но кандидат, прежде чем его судят, проходит гейт в CI своей библиотеки, а гейт повторяет каждого из своих 100 учеников каждой ячейки и держит `I1` до бита (R31). А его решающий прогон, повторённый на тех же зёрнах, должен дать те же байты.
- **`GOMAXPROCS` 1 и `GOGC` 400.** Ни то, ни другое не меняет бита, а сборка мусора стоила около четверти процессора. Процессов столько, сколько процессоров, и лишние потоки каждого процесса только отнимают процессор у других.
- **`runs/`.** Папка рабочей копии переживает Rebuild Container, а `/tmp` — нет. `git` её не видит, и файлы по-прежнему не лежат в репозитории (R28). Суд гонит кандидата, а файл базы сверяется по SHA-256 с тем, что называет `paths.json`; как это делает рецепт прогона кандидата, решает T31.

— *Отвергнуто:*
- оставить 4 000 и ждать, пока `engine` ускорит выбор, — суды ждали бы чужой задачи, а прогон шёл бы около 2,4 часа и после правок стенда;
- 1 500 учеников — около 54 минут на свободной машине, впритык к часу, а на занятой — больше;
- 1 000 только на олимпиадных ячейках и 4 000 на остальных — размер прогона зависел бы от графа, и с ростом пакетов правило пришлось бы пересматривать;
- выбор заново у всех учеников — тогда в час укладывались бы около 750 учеников, а разрешение было бы в 2,3 раза шире;
- урезать популяции до групп оценки — см. «Предложение исполнения»;
- гонять базу заново в каждом суде — два прогона вместо одного.

— *Следствие:* поправки 2026-11-10 к R04, R17, R28 и R31; новый `paths.json`. Вернуть 4 000 учеников — поправкой к R04, когда `engine` ускорит выбор (В10). Где записано: RUN.md, T36; [architecture.md](architecture.md), разделы 3, 3.6–3.8; README; `main.go`, `gate.go`, `internal/bench/`, `internal/invariant/`, `internal/criterion/`, `justfile`, `.gitignore`.

— *Поправка 2026-11-11, после ускорения `engine` (T41; R38):* решающий прогон и подтверждение идут на 1 500 учениках на ячейку, а не на 1 000. Выбор заново для `I1` — по-прежнему у первых 10 учеников ячейки; `GOGC` 400 и один процессор на часть — как были. Ответы автора 2026-11-10 при планировании T41 (RUN.md, T41).
- **Что выяснилось** (2026-11-10, машина разработки, 20 потоков; голова pull request #57 `engine`, `74afdf8`, по файлам та же, что тег v0.1.6):
  - выборка — 100 учеников на ячейку во всех 65 ячейках, по частям на 20 процессах: база на `engine` v0.1.4 — 256 с, на v0.1.6 — 166 с, в 1,54 раза меньше; потолок — 83 с, он `engine` не зовёт;
  - в пересчёте на 1 000 учеников база — около 28 мин вместо 43, на 1 500 — около 41 мин, на 1 600 — около 44;
  - прогон базы на 1 500 учениках — 38 мин по часам на почти свободной машине, потолок — около 19 мин процессора на процесс (поправка 2026-11-11 к R04);
  - мер с путём не больше допуска автора — те же семь, что на 1 000.
- **Решено автором:** 1 500 учеников; выбор заново для `I1` — у первых 10.

— *Почему:* правило R38 — решающему прогону столько учеников, сколько с запасом укладывается в час, а «с запасом» — не дольше 45 минут на свободной машине, четверть часа на занятую машину и рост пулов. 1 500 — около 41 минуты. Разрешение «не хуже» уже в √1,5 ≈ 1,22 раза. Отказ R33 от 1 500 стоял на 54 минутах при `engine` v0.1.4 — новый факт снял его.

— *Отвергнуто:*
- 1 600 — около 44 минут, впритык к 45: на рост пулов запаса не остаётся;
- 1 250 — около 35 минут: запаса больше, чем велит R38, а разрешение уже лишь в 1,12 раза;
- выбор заново для `I1` у всех учеников — тогда в 45 минут уложились бы около 1 000, а `I1` до бита и так держит повтор гейта (R31).

— *Следствие:* `criterion.Decisive` и рецепт `just decisive` — 1 500; поправка 2026-11-11 к R04 — новый `paths.json`; суды `engine-1` и `model-1` остаются на своих прогонах 1 000, а следующие суды и подтверждения идут на 1 500. 4 000 по-прежнему не укладываются: треть выбора теперь — `Readout.Task` `learner-model` (строка В32 RUN.md `engine`, «в реестр spec»), и дальше прогон ускоряет `learner-model`. Где записано: RUN.md, T41; [architecture.md](architecture.md), разделы 3.6 и 3.8; README; `internal/criterion/criterion.go`, `justfile`.

— *Поправка 2026-11-12, после ускорения `learner-model` (T46; R43, R44):* решающий прогон и подтверждение идут на 2 000 учениках на ячейку, а не на 1 500. Выбор заново для `I1` — по-прежнему у первых 10 учеников ячейки; `GOGC` 400 и один процессор на часть — как были. Ответы автора 2026-11-12 при планировании и исполнении T46 (RUN.md, T46).
- **Что выяснилось** (2026-11-12, машина разработки, 20 потоков; `learner-model` v0.1.6, `engine` v0.1.6, пакеты `modules` v0.2.5):
  - выборка — 100 учеников на ячейку во всех 65 ячейках, по частям на 20 процессах, на свободной машине: load average — 1,8 перед сборкой стенда и 2,0 после неё, у старта прогона. База — 1 мин 55 с против 166 с в T41, потолок — 1 мин 20 с против 83 с: у базы ускорение `Readout.Task`, а потолок `engine` не зовёт;
  - в пересчёте база на 1 500 учениках — около 28,8 минуты, на 1 750 — около 33,5, на 2 000 — около 38,3, на 2 250 — около 43,1; потолок на 2 000 — около 26,7. Это сходится с замером процессора T45 — около 25,9 минуты на процесс у базы на 1 500 (R43). В T41 такая выборка обещала 41 минуту на 1 500, а прогон шёл 38: пересчёт немного завышает;
  - прогоны на 2 000 учениках шли на хосте, который делили соседние контейнеры: база — 62 мин 37 с по часам и около 31,3 минуты процессора на процесс, потолок — 30 мин 37 с и около 24,1 (поправка 2026-11-12 к R04). Процессора на процесс у потолка — в 4/3 больше, чем на 1 500 в T45 (18,0), как и ждали, а у базы — лишь в 1,21 раза (25,9 → 31,3): процессорное время на перегруженном хосте само гуляет, так что размер выбран по выборке, а не по нему;
  - файл прогона — около 1,0 ГБ у базы и 0,8 ГБ у потолка; `sim paths` читает оба за 5 мин 23 с.
- **Решено автором:** размер — наибольший на сетке с шагом 250, от 1 500 до 4 000, при котором база в пересчёте выборки на свободной машине не дольше 41 минуты, как 1 500 в T41. Это 2 000; выбор заново для `I1` — у первых 10.

— *Почему:* правило R38 — решающему прогону столько учеников, сколько с запасом укладывается в час, а «с запасом» — не дольше 45 минут на свободной машине (поправка 2026-11-11). 2 000 — около 38 минут; 2 250 — около 43, впритык, как отвергнутые в T41 1 600 при 44 минутах. Разрешение «не хуже» уже в √(4/3) ≈ 1,15 раза.

— *Отвергнуто:*
- 2 250 — около 43 минут: на рост пулов запаса не остаётся;
- шаг 100 — 2 100 или 2 200 дали бы 3–5 % разрешения ценой некруглого размера и меньшего запаса;
- не дольше 38 минут — запаса больше, чем велит R38: тогда 1 750, и разрешение уже лишь в 1,08 раза;
- размер по часам прогона на занятом хосте — он мерит соседей, а не стенд.

— *Следствие:* `criterion.Decisive` и рецепт `just decisive` — 2 000; поправка 2026-11-12 к R04 — новый `paths.json`; суды `engine-1` и `model-1` остаются на своих прогонах 1 000, а следующие суды и подтверждения, начиная с `model-2`, идут на 2 000. 4 000 по-прежнему не укладываются — около 77 минут базы; Д7 закрыт вместе с поправкой 2026-11-12 к R31: размеры снова так велики, как позволяет время. Где записано: RUN.md, T46; [architecture.md](architecture.md), разделы 3.6 и 3.8; README; `internal/criterion/criterion.go`, `justfile`.

## Решающий прогон кандидата

**R34. Решающий прогон кандидата гонит тот же рецепт `just decisive`, что и прогон базы: с набором учеников и с кандидатом из `CANDIDATE_LEARNER_MODEL`, `CANDIDATE_ENGINE` или обеих. Кандидат — коммит по полному SHA, который рецепт скачивает с GitHub, или папка без незакоммиченных правок. Стенд кандидата собирается через рабочее пространство гейта: этот кусок рецептов один на оба рецепта, а сборку до прогона держит `sim gate verify`. База без кандидата собирается по `go.mod`, как прежде.** Ответы автора 2026-11-10 при планировании T31 (RUN.md, T31); остальное — предложение исполнения.
- **Что выяснилось при планировании** (2026-11-10):
  - в `engine` открыты три кандидата, pull request #46–48 — A, B и C (T31 `engine`). Судить их кругом будет T33 (поправка 2026-11-09 к R103 spec);
  - GitHub отдаёт коммит по полному SHA и тогда, когда его нет ни на одной ветке: голова открытого #48 `engine` скачалась `git fetch --depth 1`;
  - сборку гейта держит `gate.Verify`: библиотеки из рабочего пространства, общие зависимости и строки `go` файла гейта, `go.mod` `sim`, закреплённые пакеты. Частей прогона ей не нужно, но звала её только `sim gate check`, вместе с полосами;
  - шапку файла результатов пишет `bench.Names`. Библиотеку из рабочего пространства она называет коммитом, который передан через `-commit`, а библиотеку из `go.mod` — версией.
- **Решено автором:**
  - база без кандидата — по `go.mod`;
  - кандидат — SHA или чистая папка, пара — обе переменные.
- **Предложение исполнения:**
  - скрытый рецепт `workspace dir prefix clean`. Он берёт каждую библиотеку из `<prefix>_<ИМЯ>` или по коммиту файла гейта, проверяет тег, пишет `go.work` и печатает по строке на библиотеку: путь, папку, коммит. Гейт зовёт его с `GATE` и без проверки чистоты, решающий прогон — с `CANDIDATE` и с ней;
  - `sim gate verify -file <файл гейта> -go <папка>…` — `gate.Verify` без полос: кандидат, чей `go.mod` поднял общую зависимость, не прогоняется, как не прошёл бы и гейт своего CI;
  - вторая библиотека у кандидата — тоже из рабочего пространства, по коммиту файла гейта. Так `gate.Verify` применима без правок, а шапка называет обе библиотеки коммитами;
  - пара — обе переменные вместе с `CANDIDATE_PAIR=yes`. Переменная, оставшаяся от другого прогона, не делает из кандидата пару незаметно, а переменная с именем, которого нет среди библиотек файла гейта, — отказ;
  - относительная папка — от папки вызывающего, как файл прогона;
  - набор учеников рецепт проверяет до прогона, коротким прогоном одного ученика;
  - README — раздел «Judging a candidate»: от коммита pull request до суда, подтверждения и записи полос.

— *Почему:*
- **База по `go.mod`.** Тогда её файл до байта тот, чей SHA-256 называет `paths.json` (R28; поправка 2026-11-10 к R04). Через рабочее пространство шапка называла бы библиотеки коммитами, байты разошлись бы, и базу с путями пришлось бы писать заново.
- **SHA или папка.** Голова pull request лежит на GitHub, и SHA её однозначно называет. Папка нужна для кандидата, которого ещё нет на GitHub. Без правок — потому что шапка называет коммит, и файл должен быть прогоном именно его.
- **Один кусок рабочего пространства.** Гейт и прогон кандидата собирают стенд одинаково; две копии разошлись бы.

— *Отвергнуто:*
- база через `go.work` — см. «Почему»;
- только SHA — кандидата без pull request нельзя было бы прогнать;
- папка с незакоммиченными правками — шапка назвала бы коммит, а шёл бы другой код;
- своя копия куска `go.work` в рецепте `decisive`.

— *Следствие:* поправки 2026-11-10 к R26 и R28; T32 и T33 гонят кандидатов этим рецептом. Где записано: RUN.md, T31; [architecture.md](architecture.md), разделы 3, 3.6, 3.8 и 3.9; README; `justfile`, `gate.go`.

## Порядок судов и отчёт суда

**R35. Из двух первых судов — модели и движка — первым идёт тот, чьи кандидаты готовы раньше, а второй — на базе с выбранным первым. Отчёт суда лежит в репозитории, в `judgements/<критерий>-<n>/`. Строки «В реестр spec» при теге с кандидатом нет.** Ответы автора 2026-11-10 при планировании T32 (RUN.md, T32).
- **Что выяснилось при планировании** (2026-11-10):
  - разведка `learner-model` — его T30, pull request #42 — выбрала трёх кандидатов: `g10`, `halve035` и `forget-set`. Все в группе «Забывает», без цели и без нового поля, а их черновые pull request — T31 `learner-model`, ещё не начатая;
  - кандидаты движка A, B и C уже готовы, pull request #46–48 `engine`;
  - после первого выбора база — выбранное ([architecture.md](architecture.md), раздел 3.6), так что порядок судов решает, на какой базе идёт второй;
  - `sim judge` не сверяет, что вторая библиотека у кандидата та же, что у базы; от такой ошибки сейчас защищает `CANDIDATE_PAIR` рецепта (R34).
- **Решено автором:**
  - **порядок.** Сейчас первым идёт T33. Если он выбрал движок, T32 гонит базу заново — модель v0.1.4 с выбранным движком, рецептом по новому `go.mod`, — а кандидатов модели рецепт соберёт с движком по коммиту `gate.json` (R34). Путь остаётся закреплённым (R28);
  - **отчёт.** Папка `judgements/<критерий>-<n>/`, где n — номер выбора этого критерия, как у отложенных учеников (раздел 3.8 [architecture.md](architecture.md)): файл суда, отчёт `sim judge` как есть и SHA-256 каждого прогона. Итог и его причины записывает решение журнала;
  - **`spec`.** Первых кандидатов ведёт строка 4.9 [COORDINATION.md](../spec/COORDINATION.md), а порядок релизов — R103 spec: тег `sim` с полосами и отчёт `spec` увидит на `main`, а тег в «Gate» поднимает сессия библиотеки.

— *Почему:*
- **Кто готов — тот первый.** Кандидаты одной библиотеки не ждут чужих, а второй суд судит свои кандидаты вместе с уже выбранным, как раздел 3.6 и велит.
- **Отчёт в репозитории.** Полосы записываются заново pull request с причиной — отчётом суда и подтверждения (раздел 3.7 [architecture.md](architecture.md)), — и ревью читает этот отчёт рядом с полосами. Прогоны воспроизводятся до байта, но снимать отчёт заново, чтобы его прочесть, ревью не должно.

— *Отвергнуто:*
- судить парой — кандидаты движка ждали бы кандидатов модели, а сочетаний 3 × 3 — девять прогонов;
- оба суда на прежней базе — выбранные модель и движок вместе не судились бы ни разу;
- отчёт только в журнале — см. «Почему»;
- строка «В реестр spec» — то, что она спросила бы, уже ведёт строка 4.9.

— *Следствие:* T33 идёт раньше T32; отчёты — `judgements/engine-1/` и `judgements/model-1/`. Где записано: RUN.md, T32 и T33; [architecture.md](architecture.md), раздел 3.6.

— *Поправка 2026-11-10, после суда `engine-1` (R36):* n в `judgements/<критерий>-<n>/` — номер круга суда этого критерия, а не номер выбора. Круг может кончиться выходом без подтверждения, и тогда выбора нет, а следующему кругу нужна своя папка. Отложенный набор круга — первый ещё не тронутый `held-out-<критерий>-<n>` (раздел 3.8 [architecture.md](architecture.md)): у `engine-2` это снова `held-out-engine-1`.

— *Поправка 2026-11-11 (T40; R40):* суд сам сверяет, что вторая библиотека кандидата — та же, что у базы, по встроенному `gate.json`; Д6 закрыт. Защита рецепта — пара требует `CANDIDATE_PAIR=yes` (R34) — остаётся.

## Суд `engine-1`

**R36. Первый суд кандидатов `engine` — A, B и C, pull request #46–48 — кончился выходом: ни один кандидат не держит всех ограничений критерия движка, и база остаётся как есть. Подтверждения не было, и `held-out-engine-1` по-прежнему отложен. `gate.json`, `go.mod` и полосы не меняются.** Суд T33 (RUN.md), круг по поправке 2026-11-09 к R103 spec; отчёт — `judgements/engine-1/` (R35).
- **Кандидаты** — все в группе `corridor`, без новых констант и без нового поля; строки для файла суда дали их описания:
  - A, `7ef8288c…`, `distance` 1 — задача в коридоре раньше невиданной вне его;
  - B, `32aa835d…`, `distance` 2 — виденное решение после паузы не меняет порядок;
  - C, `97abff0d…`, `distance` 3 — порядок A и закрепление ошибки только в коридоре.
- **Прогоны** — решающие, на 1 000 рабочих учениках на ячейку, рецептом `just decisive` (R34), на свободной машине:
  - база — `runs/base.jsonl` прогона T36, SHA-256 `b13a9264…4c92`, тот, что называет `paths.json`;
  - A — 48 мин по часам, SHA-256 `6ac032a5…e94d`;
  - B — 42 мин, `21ad0cbd…f8e4`;
  - C — 41 мин, `79353c9a…b8459`.

  Так подтверждён и час T36 по часам (R33). Одна двадцатая прогона A прогнана заново и совпала с его файлом строка в строку, 3 250 строк.
- **Отчёт** `sim judge`:

  | кандидат | оценка | группа `corridor` | проваленных «не хуже» |
  |---|---|---|---|
  | A | 0,124 [0,115; 0,133] | 0,133 [0,122; 0,144] | 3 |
  | B | 0,130 [0,118; 0,141] | 0,147 [0,134; 0,160] | 30 |
  | C | 0,127 [0,117; 0,138] | 0,137 [0,126; 0,148] | 3 |

  - Все трое лучше базы по оценке.
  - A проваливает отставание `M3.lag` на `G2` обеих математик (на школьной — хуже на 0,036 при допуске 0,019) и `E4.given` у `F-growth5` (0,01015 при 0,01).
  - C проваливает школьные `E1.corridor` у `G2` (0,01004 при 0,01) и `M3.lag` у `G2` (0,036 при 0,018), а также олимпиадный `M3.lag` у `G2-fading`.
  - B проваливает `E4.given` почти на всех олимпиадных популяциях, `E3.retention` у `F-`, `M4.false.all` у `W-repeat` и `M3.lag`. Сбылся риск повторить виденное решение (R06 learner-model).
  - Из оценки у всех вышли три `E3.end` на школьной математике — их путь в пределах допуска (поправка 2026-11-09 к R04).

— *Почему выход:* выбор берёт только кандидата, который держит все ограничения и лучше базы ([architecture.md](architecture.md), раздел 3.6). Таких нет. Общая причина A и C — отставание модели `M3.lag` на `G2`: новый порядок задач даёт модели меньше свидетельств там, где правда учится. Это читается как сигнал для следующего круга `engine`, а не как довод ослабить допуск: допуски записаны до первого кандидата (R04).

— *Следствие:* pull request #46–48 `engine` закрывает его сессия, со ссылкой на `judgements/engine-1/`. Строки «В реестр spec» о ветке 3 у C нет: C не выбран. T32 судит кандидатов модели на прежней базе — `runs/base.jsonl`, — потому что первого выбора нет (R35). Следующий суд движка — `judgements/engine-2/`, а его подтверждение — снова на `held-out-engine-1`: тот не тронут (поправка 2026-11-10 к R35). Где записано: RUN.md, T33; `judgements/engine-1/`.

## Суд `model-1`

**R37. Первый суд кандидатов `learner-model` — `g10`, `halve035` и `forget-set`, pull request #43–45 — кончился выходом: ни один кандидат не держит всех ограничений критерия модели, и база остаётся как есть. Подтверждения не было, и `held-out-model-1` по-прежнему отложен. `gate.json`, `go.mod` и полосы не меняются.** Суд T32 (RUN.md), один круг по поправке 2026-11-09 к R103 spec; отчёт — `judgements/model-1/` (R35 и поправка к нему). Он шёл после суда движка, тоже кончившегося выходом (R36), так что база прежняя.
- **Кандидаты** — все в группе `forgets`, без цели и без нового поля; поля для файла суда дала сессия `learner-model`. Все три черновика сделаны от `main` 9fb144f, где код модели тот же, что в релизе v0.1.4:
  - `g10`, `7f735b20…` — `stabilityGrowth` 10, одна константа;
  - `halve035`, `17e0a42f…` — `stabilityCut` 0,35, одна константа;
  - `forget-set`, `75c6a908…` — обе, две константы.
- **Прогоны** — решающие, на 1 000 рабочих учениках на ячейку, рецептом `just decisive` с `CANDIDATE_LEARNER_MODEL` (R34). Движок взят по коммиту `gate.json`, v0.1.4:
  - база — `runs/base.jsonl`, SHA-256 `b13a9264…4c92`;
  - `g10` — 41 мин по часам, `e19ff45e…0434`;
  - `halve035` — 44 мин, `400d9a6d…d66e`;
  - `forget-set` — 41 мин, `cdf7755c…4c25`. Первый его запуск оборвался вместе с сессией и был повторён целиком.

  Одна двадцатая прогона `halve035` прогнана заново и совпала с его файлом строка в строку, 3 250 строк.
- **Отчёт** `sim judge`:

  | кандидат | оценка | группа `forgets` | проваленных «не хуже» |
  |---|---|---|---|
  | `g10` | 0,025 [0,011; 0,037] | 0,195 [0,098; 0,294] | 4 |
  | `halve035` | 0,014 [0,002; 0,027] | 0,117 [0,031; 0,213] | 3 |
  | `forget-set` | 0,033 [0,020; 0,046] | 0,277 [0,184; 0,373] | 8 |

  - Все трое лучше базы по оценке, и группа «Забывает» у всех выше нуля.
  - `halve035` проваливает три проверки, все на волоске от допуска: `E3.late` у `F-growth5` (0,04967 при 0,049), `M3.settle` у `F-break` (19,42 при 18,11) и `M1.rms.200` у `A-unpractised` (0,03116 при 0,03076).
  - `g10` проваливает `M1.rms.200` у `G3-drop`, `G8` и `S-scale1.5` и `M3.settle` у `F-break`.
  - `forget-set` проваливает то же, что `g10`, и ещё олимпиадный `E1.corridor` у `G2`, `G2-fading` и `S-adult` и `M1.rms.200` у `F-none`.
  - Общее у всех — догонка после перерыва `M3.settle` у `F-break`: с новыми числами забывания модель догоняет правду после перерыва дольше. Почему — суд не читает; это вопрос следующего круга.

— *Почему выход:* выбор берёт только кандидата, который держит все ограничения и лучше базы ([architecture.md](architecture.md), раздел 3.6). Таких нет. Допуски записаны до первого кандидата (R04), и то, что `halve035` не прошёл на волоске, — не довод их сдвинуть. Это сигнал для следующего круга модели: менять забывание так, чтобы не задерживать догонку после перерыва.

— *Следствие:* pull request #43–45 `learner-model` закрывает его сессия, со ссылкой на `judgements/model-1/`. Первых выборов нет ни у модели, ни у движка, так что T34 — перевод на релизы выбранных — автор не делает: её статус `[-]`. База остаётся `learner-model` v0.1.4 и `engine` v0.1.4. Следующий круг модели — `judgements/model-2/`, его подтверждение — на `held-out-model-1`. Где записано: RUN.md, T32; `judgements/model-1/`.

## Эпик E4

**R38. Эпик E4 — «Стенд за библиотеками»: стенд догоняет `spec`, `engine` и `modules` раньше новых кандидатов. В нём: переход на `spec` v0.14.0, база на `engine` v0.1.5, суд, который видит вторую библиотеку кандидата (Д6), размеры прогонов после ускорения `engine`, вид проверки `time` и наводящий вопрос на релизе `learner-model`, закрепление `modules` v0.2.1 и итог. Вторые круги кандидатов идут задачами этого эпика, когда придут кандидаты. Размеры после ускорения — по замеру, а правило часа остаётся.** Ответы автора 2026-11-10 при развёртывании E4 (RUN.md, таблица 1).
- **Что выяснилось при планировании** (прочитано 2026-11-10):
  - `spec` вышел до v0.14.0, а `sim` стоит на v0.12.0 (R32);
  - `spec` ответил на В10 по рекомендации строки (Д167 RUN.md `spec`): выбор на большом пуле ускоряет `engine` своей задачей. E4 `engine` — профиль `Choose` на `olympiad-math` 0.1.1, ускорение без изменения выбора, около 10×, и «Gate» на теге `sim`, который вернёт гейту 300 учеников (T37–T39 RUN.md `engine`);
  - `spec` ответил на В11 (Д173 RUN.md `spec`; строка 4.10 его COORDINATION.md). Ответы на В10 и В11 вышли с тегом v0.14.0, а подмодуль стоит на v0.12.0, так что строки таблицы 3 уходят переходом T38. Порядок такой: `learner-model` выпускает `time` и наводящий вопрос одним релизом через гейт `sim` (его Д14); затем `sim` переводит их в правде ученика и в свидетельстве и выпускает тег; затем `modules` выпускает v0.2.1 (T56 RUN.md `modules`), и `sim` закрепляется на нём;
  - `engine` выпустил v0.1.5 с исправлением R82 spec в базе: задача последней попытки с виденным решением не идёт следующей. Его нашёл инвариант `I4` гейта `sim` на кандидатах T31. По `engine`, все 69 чисел гейта те же до бита, а строка В29 RUN.md `engine` просит `sim` и `api` взять v0.1.5;
  - новых кандидатов скоро не будет: `engine` ускоряет выбор, а `learner-model` записывает итог первого суда.
- **Решено автором:**
  - **состав** — «Стенд за библиотеками», вторые круги — когда придут кандидаты;
  - **размеры после ускорения:** гейту — 300 учеников, если он укладывается на раннере в 15 минут; решающему прогону — сколько с запасом укладывается в час, но не больше 4 000. Новые размеры — это новые пути и новые полосы.
- **Предложение исполнения к размерам**, решает T41 при планировании: если 300 учеников на раннере не укладываются в 15 минут, гейту — наибольший размер, который укладывается, но не меньше 100; «с запасом» — решающий прогон не дольше 45 минут на свободной машине, четверть часа на занятую машину и рост пулов.
- **Предложение исполнения:** порядок задач — T38–T44 (RUN.md, «Почему такой порядок»); вторые круги — задачами с номером, свободным в момент записи.

— *Почему:*
- **Стенд за библиотеками.** Судить кандидатов можно только на том, на чём стоят библиотеки и что получат ученики. `spec`, `engine` v0.1.5 и пакеты 0.2.1 уже вышли или выходят, а кандидатов пока нет.
- **Размеры по замеру.** R33 урезал прогоны до часа, потому что выбор был дорог. Когда он подешевеет, вернуть стоит столько, сколько позволяет то же правило, а не прежние числа вслепую.

— *Отвергнуто:*
- E4 как его предложил T35 — только вторые круги: эпик ждал бы кандидатов, а стенд отставал бы от библиотек;
- вернуть 300 и 4 000 вслепую — решающий прогон мог бы выйти за час;
- не трогать размеры после ускорения — гейт оставался бы грубее, чем может, а `engine` ждёт тега с 300 учениками.

— *Следствие:*
- E3 перенесён в архив [docs/plan/E3.md](plan/E3.md);
- таблица 1 RUN.md — E4 задачами T38–T44;
- строка E4 ушла из таблицы 2; В10 и В11 закрывает T38 ответами, пришедшими с тегом `spec` v0.14.0;
- [architecture.md](architecture.md) — статус и подраздел 2.5, «Решено»; README, «Status».

## Тег `spec` v0.15.0

**R39. `sim` стоит на `spec` v0.15.0 — последнем теге, а не на v0.14.0, который назвала задача. Подмодуль, `go.mod` и `gate.json` — на нём. Числа гейта те же до бита, и полосы прежние. Строки В10 и В11 таблицы 3 RUN.md закрыты ответами `spec`, пришедшими с тегом.** Поручение автора 2026-11-11: «делай T38, но вышла более свежая версия spec» (RUN.md, T38).
- **Что выяснилось при переходе** (прочитано 2026-11-11; источники — [architecture.md](architecture.md), раздел 4):
  - `spec` вышел до v0.15.0, тег на коммите 58d3f46, и `gh release verify` его подтвердил;
  - для `sim` теги v0.13.0–v0.15.0 не меняют ничего, что стенд читает: `RULES.md`, текст и схемы формата 0.4 кроме примера ловушек, типы `formatv0`, заготовка и контейнер — те же, что на v0.12.0, а карточка architecture 3.9 spec не тронута;
  - модуль Go `spec` добавил пакеты `consentagev1`, `exportv1` и `legalv2` и части `Legal`, `ConsentAge` и `Export`, а встроенные файлы держит частями. Его `go.mod` — тот же, так что новых зависимостей нет;
  - за `sim` теги решили четыре ответа таблицы 4 RUN.md `spec`:
    - Д167 — ответ на В10: выбор ускоряет `engine` своей задачей, а гейт и решающий прогон растут после неё. Это сделала T41: `engine` v0.1.6, гейт на 150 учениках, решающий прогон на 1 500 (поправки 2026-11-11 к R31 и R33);
    - Д199 — тот же переход на `engine` v0.1.6: тоже T41;
    - Д198 — ответ на В32 `engine`: `Readout.Task` ускоряет `learner-model` своей задачей, а `sim` возвращает гейту 300 учеников и решающему прогону 4 000, когда выйдут обе правки. Это долг Д7 таблицы 4 RUN.md;
    - Д173 — ответ на В11: строка 4.10 COORDINATION.md `spec`. Сначала `learner-model` выпускает `time` и наводящий вопрос одним релизом через гейт `sim`, затем `sim` их переводит и выпускает тег, затем `modules` выпускает v0.2.1. Это T42 и T43;
  - библиотеки модуля Go `spec` не требуют: `go.mod` `learner-model` v0.1.4 — без `require`, а `engine` v0.1.6 требует только `learner-model`. `spec` — зависимость самого стенда, и её версия — тег подмодуля (R64 spec);
  - `just gate exact` на v0.15.0 — все 69 чисел до бита те же, полосы не переписывались.
- **Решено автором:** переход сразу на v0.15.0.

— *Почему:* репозиторий стоит на последнем теге `spec` (RULES.md, Across repositories), и между v0.14.0 и v0.15.0 для стенда ничего не меняется. Переход на v0.14.0 лишь отложил бы следующий.

— *Отвергнуто:* v0.14.0, как назвала задача, — ещё один переход сразу следом.

— *Следствие:* В10 и В11 ушли из таблицы 3 RUN.md; `gate.json` называет `spec` v0.15.0; тег — в CLAUDE.md, README и начале RUN.md. Мерж выпускает тег `sim`, и гейты библиотек берут его. Где записано: RUN.md, T38; [architecture.md](architecture.md), статус, раздел 2 и раздел 4; README.

— *Поправка 2026-11-12:* `sim` стоит на `spec` v0.16.3 (R45): подмодуль, `go.mod` и `gate.json` — на нём, Go в контейнере — 1.27.4, а числа гейта — те же до бита.

## Вторая библиотека кандидата

**R40. Суд держит базу и кандидатов на библиотеках базы. `sim` встраивает `gate.json`, как встраивает `paths.json`, и `sim judge` отказывает суд целиком, если: база или база подтверждения идут не на тех библиотеках, на которых прочитаны пути; база — не на версиях, которые называет файл гейта; вторая библиотека кандидата одной библиотеки — или его же на отложенных учениках — не версия базы и не коммит её закрепления. У пары второй библиотеки нет. Шапки файлов результатов, рецепт и байты прогонов не меняются.** Ответы автора 2026-11-11 при планировании T40 (RUN.md, T40); остальное — предложение исполнения. Закрывает Д6 таблицы 4 RUN.md.
- **Что выяснилось при планировании** (2026-11-11):
  - шапка базы называет обе библиотеки версиями — базу рецепт `just decisive` собирает по `go.mod`, — а шапка кандидата обе коммитами: вторую библиотеку рецепт берёт коммитом из `gate.json` (R34). Сравнить их напрямую суд не может: версию с коммитом её тега сопоставляет только файл гейта;
  - суд не сверял и базу с той, на которой прочитаны пути, хотя `paths.json` называет её версии (находка ревью T41);
  - пока у библиотеки нет релиза, `go.mod` называет её псевдоверсией коммита, и `sim gate verify` уже сопоставляет псевдоверсию с коммитом по суффиксу из 12 знаков;
  - пакеты `criterion` и `gate` друг друга не импортируют.
- **Решено автором:**
  - источник — `gate.json`, встроенный в `sim`;
  - не та вторая библиотека — отказ всему суду, а не вердикт;
  - заодно суд сверяет библиотеки базы с путями, а SHA-256 файла базы — нет.
- **Предложение исполнения:**
  - `main` встраивает `gate.json` и передаёт суду закрепления библиотек — `criterion.Pin`, коммит и тег, по пути модуля;
  - `criterion.Judge` после путей зовёт сверку библиотек. Без закреплений суд не судит. Сверка с файлом гейта принимает тег, а пока его нет — псевдоверсию коммита, как `sim gate verify`; вторую библиотеку кандидата — ещё и коммит закрепления, а пустое имя — нет;
  - отказ называет кандидата, библиотеку и обе точки: где шёл кандидат, а где база.

— *Почему:*
- **`gate.json`.** Он один называет и версию, и коммит её тега, а прогоны кандидатов собраны именно по нему (R34). Шапки, рецепт и байты прогонов остаются прежними, и закреплённый прогон базы, чей SHA-256 называет `paths.json`, тот же (R28). Встроенный файл судит тем закреплением, с которым собран `sim`, — как пути.
- **Отказ всему суду.** Прогон на чужой второй библиотеке — не кандидат этого суда: это ошибка входа, как прогон не на тех учениках, а не вердикт кандидату.
- **База против путей.** Доля кандидата — его разность с базой, делённая на закреплённый путь. На базе других библиотек путь не тот, и доля была бы чужой.

— *Отвергнуто:*
- поле в шапке кандидата «взято у вызывающего» — меняет рецепт и шапку, а у прежних кандидатов его нет;
- путь к файлу гейта в файле суда — можно указать не тот;
- проверка «не та библиотека» с «не выполнено» у одного кандидата — суд выносил бы вердикт по прогону, который не из него;
- сверять ещё и SHA-256 файла базы — грубым взглядам и тестам понадобился бы обход, а на чужой базе отказывает уже сверка библиотек;
- только Д6, без сверки базы с путями, — см. «Почему».

— *Следствие:*
- суды `engine-1` и `model-1` на нынешнем дереве получают отказ: их база — на `engine` v0.1.4, а пути — на v0.1.6. Их отчёты — записи на своих коммитах, и там они воспроизводятся;
- защита рецепта остаётся: пара требует `CANDIDATE_PAIR=yes` (R34);
- поправки 2026-11-11 к R25 и R35; Д6 закрыт.

Где записано: RUN.md, T40; [architecture.md](architecture.md), разделы 3 и 3.6; README; `main.go`, `internal/criterion/criterion.go`.

## Формат 0.4 на стенде

**R41. Стенд читает вид проверки `time` и наводящий вопрос формата 0.4 на `learner-model` v0.1.5. `time` — свой вид проверки, и правда угадывает его с вероятностью 0, как число. Наводящий вопрос стенд читает только тем, что он у задачи есть: просьба `G8` до первого ответа на такой задаче даёт и приём, и вопрос, свидетельство несёт `LeadingQuestionOnRequest`, а подъём логита — прежний. Решающего прогона здесь нет, и `paths.json` остаётся на `learner-model` v0.1.4: до T43 суд отказывает.** Ответы автора 2026-11-11 при планировании T42 (RUN.md, T42); остальное — предложение исполнения. Исполняет строку 4.10 COORDINATION.md `spec` в части `sim` (Д173 RUN.md `spec`; R39).
- **Что выяснилось при планировании** (2026-11-11; источники — [architecture.md](architecture.md), раздел 4):
  - `learner-model` v0.1.5 добавляет `CheckKindTime` без угадывания и поле свидетельства `LeadingQuestionOnRequest` — «подсказка по просьбе дала наводящий вопрос». Модель считает его помощью, как приём по просьбе: ответ идёт в шаг, но «освоено» не ставит. Вид свидетельства — 2; по R28 learner-model числа модели не меняются;
  - наводящий вопрос даёт только подсказка по просьбе до первой ошибки ([текст формата](../spec/format/README.md), разделы 14 и 19). На стенде её просит только `G8`, а базовый ученик — никогда ([architecture.md](architecture.md), раздел 3.1);
  - `modules` v0.2.1 с `time` у 30 задач и наводящими вопросами у 311 задач вышел 2026-11-10, раньше строки 4.10 (строка В38 RUN.md `modules`), так что T43 ждёт только этой записи.
- **Решено автором:**
  - наводящий вопрос — одна просьба, один подъём, без нового допущения;
  - решающего прогона в T42 нет: пути один раз читает T43, уже на `modules` v0.2.1.
- **Предложение исполнения:**
  - `release-verify` принял `learner-model-proto-v0.1.5.tar.gz` — «learner-model v0.1.5, signed by the shared release.yml of .github v0.1.21»; `go.mod` и `gate.json` — на v0.1.5, коммит тега `3a2cd29`. `just gate exact` — все 69 чисел до бита те же. Полосы записаны заново `just gate record` только ради причины: значения в них те же до бита, а файл гейта называет теперь коммит, на котором они записаны;
  - перевод: `time` — `CheckKindTime`; `pack.Task.Leading` — у задачи есть `leading_question`. Канонический текст держит оба (поправки 2026-11-11 к R14 и R27);
  - правда: `Attempt.LeadingQuestion` — приём `G8` на задаче с вопросом, без нового числа из места `attempts`; свидетельство несёт его в `LeadingQuestionOnRequest`, вида 2 (поправка 2026-11-11 к R18). Правда ловушек не знает, так что 12-часового промаха у `time` нет;
  - синтетика: навык листа на навыке-опоре отвечается временем, а у задач листа сложности 5 есть наводящий вопрос (поправка 2026-11-11 к R19);
  - тест суда в `main_test.go` держит свои прогоны и встроенный файл гейта на библиотеках путей: между переходом библиотеки и решающим прогоном, который читает пути заново, они расходятся.

— *Почему:*
- **`time` — с угадыванием 0.** Ответ времени набирают, а не выбирают, как число, и модель угадывает его так же — 0 (R19 learner-model). Правда в остальном держит допущения модели (раздел 3.1 [architecture.md](architecture.md)).
- **Одна просьба, один подъём.** По тексту формата наводящий вопрос приходит с той же подсказкой, что и приёмы, а модель считает их одной помощью. Насколько вопрос помогает сверх приёма, данных нет — ни у стенда, ни у `learner-model`, где R06 и R21 learner-model отвергли вес без данных. Без нового допущения числа `G8` не сдвигаются и в T43, а стенд судит задачу так, как её получит ученик: свидетельство несёт то, что запишет сервер.
- **Только «вопрос есть».** Слова вопроса — для ученика, а стенд читает их не больше, чем текст условия. Поэтому правка слов не будит гейт CI `modules`, а добавленный или снятый вопрос — будит.
- **Без решающего прогона.** T43 закрепляет пакеты с `time` и вопросами, а это новые пути в любом случае. Прогон здесь на `modules` v0.2.0 стоил бы около часа машины ради путей, которые T43 тут же перепишет. Кандидатов сейчас нет, так что отказ суда до T43 никого не держит.

— *Отвергнуто:*
- свой подъём наводящего вопроса сверх приёма — величина без данных;
- новая популяция, которая просит вопрос, — ещё одна ячейка решающего прогона, и тоже без данных;
- `time` переводить как `number` — тогда стенд мерил бы другую задачу (вариант (б) подраздела 2.5 [architecture.md](architecture.md));
- новый навык синтетики с `time` — он сдвинул бы выбор движка у ученика обеих синтетик, и тест якоря перестал бы видеть освоенную опору;
- решающий прогон в T42 и T43 внутри T42 — см. «Почему» и RUN.md, T42.

— *Следствие:*
- пока T43 не прочтёт пути заново, `sim judge` отказывает любому суду: база на `learner-model` v0.1.5, а `paths.json` прочитан на v0.1.4 (R40);
- поправки 2026-11-11 к R14, R18, R19 и R27; подраздел 2.5 [architecture.md](architecture.md) сделан;
- мерж выпускает тег `sim`, и гейты библиотек берут его своими задачами. Строки «в реестр spec» нет: порядок, который нарушил выход `modules` v0.2.1, уже подняла строка В38 RUN.md `modules`.

Где записано: RUN.md, T42; [architecture.md](architecture.md), подраздел 2.5 и разделы 3, 3.1–3.3; README; `internal/pack/translate.go`, `internal/learner/attempt.go`, `internal/run/run.go`, `internal/synthetic/synthetic.go`.

— *Поправка 2026-11-11 (T43):* суд больше не отказывает. Пути прочитаны заново на `learner-model` v0.1.5 и пакетах `modules` v0.2.5 (поправка 2026-11-11 к R04), а тест суда держит свои прогоны на библиотеках путей, как и прежде, — теперь это те же, что в `gate.json`. Почему — R42.

— *Поправка 2026-11-11 (T45):* перевод и правда — те же, а база — на `learner-model` v0.1.6, где `Readout.Task` быстрее при тех же битах: оба решающих прогона и гейт совпали с прогонами на v0.1.5 до бита. Пути прочитаны заново на ней, так что суд не отказывает. Почему — R43.

## Закрепление `modules` v0.2.5

**R42. Стенд закреплён на релизе `modules` v0.2.5 — `olympiad-math` и `grade-3-math` 0.3.0 с проверками `time`, наводящими вопросами и разобранными примерами и `english-essays` 0.2.0. Пути прочитаны заново на нём и на `learner-model` v0.1.5, полосы гейта — тоже. Размеры прогонов прежние: гейт — 150 учеников на ячейку, решающий прогон — 1 500.** Ответы автора 2026-11-11 при планировании и исполнении T43 (RUN.md, T43); остальное — предложение исполнения. Вместе с R41 исполняет часть `sim` строки 4.10 COORDINATION.md `spec` (Д173 RUN.md `spec`).
- **Что выяснилось при планировании** (2026-11-11; источники — [architecture.md](architecture.md), раздел 4):
  - `modules` v0.2.1 вышел 2026-11-10 с обеими математиками 0.3.0. Индекс v0.2.5 побайтно тот же, что у v0.2.1, — SHA-256 `cebf23cf…aec6a`: в v0.2.2–v0.2.5 менялись только синтетические образцы и переход на `spec` v0.15.0;
  - `olympiad-math` 0.3.0 — 775 задач вместо 772 и 22 шаблона вместо 19, наводящий вопрос у 322 задач; проверок `number` — 572, `choice` — 166 вместо 196, `fraction` — 7, `time` — 30;
  - `grade-3-math` 0.3.0 — те же задачи и проверки, а в каждом из четырёх уроков новое звено — разобранный пример, две минуты, как любое звено без задачи (R15);
  - стенд читает новые пакеты без правок кода, а выборка решающего прогона на копии репозитория дала все 65 ячеек.
- **Решено автором:**
  - закрепление — последний релиз, v0.2.5, а не v0.2.1, который назвала задача;
  - воспроизводимость — полным повтором потолка;
  - при исполнении: решающему прогону — по-прежнему 1 500 учеников, хотя база шла 47 мин 48 с. Хост в это время делили соседние контейнеры, а правило часа прогон держит.
- **Предложение исполнения:**
  - `release-verify` принял индекс — «module-index.json: modules v0.2.5, signed by the shared release.yml of .github v0.1.15» — и оба архива из v0.2.1 — «modules v0.2.1, signed by the shared release.yml of .github v0.1.15». Архив сочинений 0.2.0 — тот же, что принят в R32;
  - индекс лежит в `internal/published/v0.2.5/`, архивы — в папках релизов, которые называют их записи: `v0.2.1/` и `v0.2.0/` (поправка 2026-11-11 к R14);
  - тесты фактов закрепления — на новые пакеты, с фактами 0.4: проверки `time`, задачи с наводящим вопросом, звено разобранного примера. Тест входа прогона держит экземпляр и факт «решение видел» по тому, шаблон ли задача: новый шаблон `squares-on-lines` попал в выбор его ученика;
  - решающий прогон, новые пути — поправка 2026-11-11 к R04 — и полосы гейта (поправки 2026-11-11 к R26 и R28). Потолок повторён целиком и совпал до байта. Гейт в контейнере шёл 7 мин по часам на перегруженном хосте, 28,4 мин процессора — против 23,6 в T41; время на раннерах покажет CI pull request.

— *Почему:*
- **v0.2.5.** Индекс тот же, а стенд называет последний релиз `modules`, как стоит на последнем теге `spec`.
- **Сейчас, а не позже.** С T42 суд отказывает, пока пути не прочитаны заново. Закрепление и так требует новых путей, так что один прогон даёт и пакеты, которые получат ученики, и пути на новой библиотеке.
- **1 500 при 47 минутах.** 45 минут R38 — запас на свободной машине, а прогон шёл на занятом хосте: процессора на процесс у базы — около 38,7 минуты против 33 в T41, и часть разницы — соседи. Прогон держит правило часа, а 1 400 учеников сделали бы разрешение «не хуже» грубее и стоили бы ещё одного прогона базы.

— *Отвергнуто:*
- v0.2.1 — тот же индекс, но не последний релиз;
- одна двадцатая базы и потолка вместо полного повтора — около 17 минут экономии;
- 1 400 учеников или новый замер на свободном хосте — см. «Почему».

— *Следствие:*
- суд снова судит (поправка 2026-11-11 к R41);
- пороги CI `modules` берутся из новой записи полос (R27);
- мерж ветки эпика выпускает тег `sim` с T42 и T43, и гейты `learner-model` и `engine` берут его своими задачами; строка 4.10 COORDINATION.md `spec` — строкой В12 таблицы 3 RUN.md, «в реестр spec».

Где записано: RUN.md, T43; [architecture.md](architecture.md), статус, подраздел 2.5 и разделы 3, 3.3, 3.8 и 4; README; `internal/published/`, `gate.json`, `internal/criterion/paths.json`.

## База на `learner-model` v0.1.6

**R43. База стоит на `learner-model` v0.1.6, где `Readout.Task` быстрее при тех же битах. Пути прочитаны заново на ней, полосы гейта — тоже, а числа не сдвинулись ни в одном бите. Размеры прогонов прежние: гейт — 150 учеников на ячейку, решающий прогон — 1 500; поднимает их замер эпика E5 (Д7).** Ответ автора 2026-11-11, после итога E4: «сделай что сейчас можешь из таблицы 4» и переход на v0.1.6 с путями, а не только предварительный замер (RUN.md, T45); остальное — предложение исполнения.
- **Что выяснилось** (2026-11-11; источники — [architecture.md](architecture.md), раздел 4):
  - `learner-model` v0.1.6 вышел 2026-11-11, тег на коммите 7d796e2, — «Make Readout.Task faster with the same bits (T40)». Строка В39 его RUN.md обещает тот же API и те же P до бита. В `go.mod` модуля на теге требований нет — только `go 1.27.3`, — так что тег `spec` он не двигает;
  - `release-verify` принял `learner-model-proto-v0.1.6.tar.gz` — «learner-model v0.1.6, signed by the shared release.yml of .github v0.1.21»;
  - из двух долгов таблицы 4 сейчас делается только часть Д7: Д3 ждёт этапа 2, а размеры — замера на свободной машине и на раннерах. Хост в этот день делили соседние контейнеры, load average — 25–46 на 20 ядрах.
- **Предложение исполнения:**
  - `go.mod`, `go.sum` и `gate.json` — на v0.1.6 (7d796e2); тесты, которые называют закрепление, — тоже;
  - `just gate exact` — все 69 чисел до бита те же; полосы записаны заново только ради причины (поправка 2026-11-11 к R26);
  - решающие прогоны базы и потолка на 1 500 учениках — те же байты строк, что у T43, и новые пути с теми же значениями (поправки 2026-11-11 к R04 и R28).
- **Замер ускорения, предварительный.** Процессора на процесс у базы — около 25,9 минуты против 38,7 в T43, на треть меньше; у потолка — 18,0 против 18,8: потолок `Readout.Task` почти не зовёт. По часам база шла 62 мин 24 с — дольше часа — при load average до 46. Гейт в контейнере — 1 мин 16 с по часам и 18,7 минуты процессора против 28,4 в T43. Это не замер Д7: тот идёт на свободной машине и на раннерах, и по нему E5 выбирает размеры (поправки 2026-11-11 к R31 и R33).

— *Почему:*
- **Сейчас, а не в E5.** Переход ничего не стоит числам и снимает с E5 одну смену закрепления: второй круг модели судится на одном закреплении `sim` (строка В42 RUN.md `learner-model`), а замер размеров должен идти уже на v0.1.6. Решающие прогоны заодно дали первую оценку ускорения на прежних размерах.
- **Пути заново, хотя числа те же.** Суд держит базу на библиотеках путей (R40): без новых путей он отказывал бы каждой базе на v0.1.6, как от T42 до T43.
- **Без второго прогона потолка.** Новый потолок совпал до байта строк с потолком T43, а тот сам повторён целиком (R42).
- **62 минуты при правиле часа.** Время по часам — соседей: процессора база взяла на треть меньше, чем в T43, где прогон уложился в 47 минут. Размер решающего прогона по этому замеру не меняется ни вверх, ни вниз.

— *Отвергнуто:*
- только предварительный замер на копии репозитория, без перехода — закрепление осталось бы на v0.1.5, и смена легла бы на E5;
- переход без новых путей — суд отказывал бы;
- поднять размеры по этому замеру — хост занят, а правило поправок 2026-11-11 к R31 и R33 требует свободной машины и раннеров.

— *Следствие:*
- Д7 остаётся: размеры поднимает E5, а его первая задача теперь — замер и размеры, без перехода на v0.1.6;
- мерж ветки эпика E4 выпускает тег `sim` и с этим переходом.

Где записано: RUN.md, T45; [architecture.md](architecture.md), статус и разделы 3, 3.6–3.9 и 4; README; `go.mod`, `gate.json`, `internal/criterion/paths.json`.

## Эпик E5

**R44. Эпик E5 — «Размеры и второй круг модели»: размеры прогонов на `learner-model` v0.1.6 по замеру (Д7), суд `model-2` с подтверждением на `held-out-model-1`, перевод базы на релиз выбранного и итог — T46–T49. Суд движка `engine-2` — отдельный эпик E6 таблицы 2 RUN.md, на базе после `model-2`.** Ответы автора 2026-11-11, после итога E4 и T45: «разверни эпик E5», и состав — размеры и `model-2`, а `engine-2` — отдельно (RUN.md, таблица 1, «Решено автором»); остальное — предложение исполнения.
- **Что выяснилось** (2026-11-11; RUN.md `learner-model` и `engine` на `main`, прочитаны через `gh api`):
  - E5 `learner-model` — «Второй круг кандидатов модели» — развёрнут: его T41 ждёт тег `sim` с размерами после v0.1.6, T43 — разведка после разбора `M3.settle` у `F-break`, T44 — до трёх черновых pull request кандидатов, T45 — суд `sim` и выпуск выбранного;
  - `engine` — в своём E7, «Закрепление в уроке», а его E5 — «Второй круг коридора» — стоит в таблице 2 неразвёрнутым: кандидатов движка не видно;
  - строки В42 и В43 RUN.md `learner-model`, «в реестр spec», `spec` ещё не записал. Их рекомендации — одно закрепление `sim` на круг модели и круг модели раньше круга движка;
  - база уже на `learner-model` v0.1.6 (R43), так что размеры меряются без перехода.
- **Предложение исполнения:** порядок задач — RUN.md, «Почему такой порядок»; одно закрепление от тега T46 до выбора суда `model-2` — «Общее для задач» эпика; размер гейта — от 100 до 300, как в поправке 2026-11-11 к R31, а не от 150, как в строке E5 таблицы 2, которую предложила T44.

— *Почему:*
- **Размеры первыми.** Разведка и суд `learner-model` идут на размерах решающего прогона `sim`, а «Gate» `learner-model` ждёт тега с ними. Размеры двигают пути и полосы, а посреди круга закрепление не меняется (В42 RUN.md `learner-model`).
- **`engine-2` отдельно.** Круг модели первым — по рекомендации В43 RUN.md `learner-model`, и круг движка судится на базе после него. А кандидатов движка нет даже в плане: эпик, который ждёт их, стоял бы открытым без срока, а с ним — итог размеров и круга модели.

— *Отвергнуто:*
- все три круга в одном эпике, как предложила T44 (строка E5 таблицы 2 RUN.md в редакции коммита 77828ca) — см. «Почему»;
- суд `model-2` раньше размеров — круг пришлось бы судить заново или на грубых размерах, а «Gate» `learner-model` всё равно ждёт тега с ними.

— *Следствие:* E4 — в архиве [docs/plan/E4.md](plan/E4.md); строка E5 таблицы 2 RUN.md ушла в таблицу 1, а на её место встала E6 «Второй круг движка»; Д7 закрывает T46.

Где записано: RUN.md, таблицы 1, 2 и 4 и «Эпик E5»; [architecture.md](architecture.md), статус и раздел 3; README.

— *Поправка 2026-11-12:* у «одного закрепления на круг» — исключение R45: версия `spec` в `go.mod` и `gate.json` поднята до v0.16.3 посреди круга `model-2`, потому что биты те же; пакеты, библиотеки, пути и полосы прежние. Переход и перестройка документов по строке 7.10 COORDINATION.md `spec` — задачи T50–T53 вне эпика.

## Тег `spec` v0.16.3

**R45. `sim` стоит на `spec` v0.16.3 — последнем теге, из линии 0.16 эпика E10 `spec`. Подмодуль, `go.mod` и `gate.json` — на нём, Go в контейнере — 1.27.4; числа гейта те же до бита, так что пути, полосы и прогоны базы остаются. Строку 7.10 COORDINATION.md `spec` — переход, перестройку документов и перевод `docs/` на английский — `sim` берёт раньше пилота `module-kit` и делит на четыре задачи вне эпика, T50–T53; порядок работы по заготовке v0.16.3 — уже с T50.** Поручение автора 2026-11-12 — «переходи на spec последней версии» — и его ответы на два вопроса перехода (RUN.md, T50).
- **Что выяснилось при переходе** (`git -C spec diff v0.15.0 v0.16.3`, прочитано 2026-11-12; источники — [architecture.md](architecture.md), раздел 4):
  - `spec` вышел до v0.16.3, тег на коммите 664398e, и `gh release verify` его подтвердил;
  - линия 0.16 — эпик E10 `spec`, экономия токенов сессий (R143 spec, R145 spec). Переход на неё у каждого репозитория — строка пункта 7 COORDINATION.md `spec`, у `sim` — 7.10: одна задача вне эпика по общему тексту T123 `spec`, с переходом на тег, перестройкой документов по образцу `spec` и переводом `docs/`. Строки 7.2–7.12 ждут пилота — строки 7.1 `module-kit`. Около 20:30 UTC 2026-11-12 `module-kit` стоял на v0.15.0 без pull request перехода, а в 20:59 UTC перешёл на v0.16.3 только тегом — Go 1.27.4 и Д215, — оставив остальное строки 7.1 своей задаче T43 (R54 module-kit): пилот ещё не пройден. Большой репозиторий, по тексту T123, до начала предлагает разбить задачу по правилу 11 своего RUN.md;
  - тег не меняет ничего, что стенд читает: `format/`, `types/formatv0` и корневой пакет `spec` те же, а изменились контракты `api/v1` и `admissionv0`, которых стенд не берёт. `go.mod` модуля `spec` тот же — новых зависимостей нет. Шапка файла прогона называет только библиотеки, так что прогоны базы и потолка на 2 000 учеников в силе;
  - Go стадии `toolchain` — 1.27.4: в 1.27.3 — уязвимость GO-2026-6859 стандартной библиотеки, и переход на тег приносит 1.27.4 каждому репозиторию Go вместе с `go 1.27.4` lock-файла gopls заготовки. Директивы `go` в `go.mod` не меняются — они называют наименьшую версию, а не версию контейнера (поправка 2026-11-11 к R64 spec);
  - RULES.md: одна задача — одна сессия, и хук `UserPromptSubmit` останавливает фразу, которая начинает другую; решение, которое меняет другую запись, ставит к ней поправку; внешние факты — записями `docs/sources.md`, итоги замеров — в `docs/measurements.md`; `docs/` — на английском, а документ на русском остаётся таким до своей задачи перевода; входящие `spec` разбирает по слову автора или в конце эпика (R143 spec);
  - Д215 `spec` — строка `**/sim/spec/.claude/rules/**` в `claudeMdExcludes`: правила частей `spec` не грузятся в сессию `sim`. Строки Д140, Д196 и Д198 `spec`, у которых «когда» — переход или задача `sim`, исполнены раньше: R32, R39–R43, T46;
  - ответ на В12 таблицы 3 RUN.md — по рекомендации (а): строка 4.10 COORDINATION.md `spec` закрыта тегом `sim` v0.1.9, а строка 3.9 больше её не ждёт (Д173 в DEBT.md `spec` и таблица 5 RUN.md `spec`).
- **Решено автором (2026-11-12):**
  - вся строка 7.10 — сейчас, раньше пилота, и четырьмя задачами вне эпика по правилу 11: T50 — тег и порядок работы, T51 — `sources.md`, `measurements.md` и перевод architecture, T52 и T53 — перевод журнала, бюджеты и замер «после». Около 650 КБ русского текста — журнал 424 КБ, architecture 225 КБ — на одну проверку не ложатся;
  - исключение из правила эпика E5 «одно закрепление на круг»: версия `spec` в `go.mod` и `gate.json` меняется посреди круга `model-2`, если биты те же. Условие — `just gate exact`, все числа до бита, и первые 100 учеников каждой ячейки на новом стенде побайтно те же, что их строки в решающем прогоне базы на 2 000, как в R32.
- **Что изменено здесь:**
  - `spec/` — на v0.16.3; `go.mod` и `go.sum` — модуль Go `spec` v0.16.3; `gate.json` — `spec` v0.16.3, полосы не переписывались; `gate_test.go` — версию `spec` в сборке, которую подставляют тесты, берёт `specRecorded` из встроенного `gate.json`, а не литералом, и следующий переход её не тронет;
  - `.devcontainer/Dockerfile` — `GO_VERSION` 1.27.4 и суммы, как в `spec`; `.devcontainer/gopls/go.mod` — `go 1.27.4`, как в заготовке, а `go.sum` тот же;
  - `.claude/settings.json` — по заготовке и виду `go`: сжатие на 400 тыс. токенов, без коннекторов claude.ai, хук `spec/tools/session/`, плагин `gopls-lsp` и строка Д215;
  - CLAUDE.md — по заготовке: документы с `DEBT.md` и `sources.md`, бюджеты и правила частей `.claude/rules/` — `go`, `gate` и `ci`; AGENTS.md — заготовки, со строкой о правилах по путям;
  - RUN.md — «Как пользоваться» и «Как устроен план» по заготовке: одна задача — одна сессия, тело сделанной задачи — сразу в архив эпика, таблица 4 — в своём файле `DEBT.md`; тело T46 — в новом архиве [docs/plan/E5.md](plan/E5.md); раздел «Задачи вне эпика» с T51–T53; В12 закрыта, В13 — в таблице 3;
  - тег — в CLAUDE.md, README.md, начале RUN.md и статусе [architecture.md](architecture.md).
- **Биты те же:** `just gate exact` на v0.16.3 — все 69 чисел до бита; первые 100 учеников каждой из 65 ячеек на новом стенде — 6 500 строк — побайтно те же, что их строки в решающем прогоне базы T46 на 2 000, `runs/base.jsonl`. Так — и на Go 1.27.3 контейнера, и на Go 1.27.4, на котором стенд соберёт CI: его тулчейн взят из прокси Go через `GOTOOLCHAIN=go1.27.4`. Условие исключения выполнено, и полосы не переписывались.

— *Почему:* поручение автора; репозиторий стоит на последнем теге `spec` (RULES.md, Across repositories), а экономия E10 у `sim` — сессии каждой следующей задачи. Посреди круга `model-2` — потому что правило «одно закрепление» держит числа, на которых судится круг, а они не сдвинулись ни на бит, и кандидатов второго круга ещё нет. Раньше пилота — слово автора; пилот ловит ошибки общего текста T123, и что найдут T50–T53, они записывают в «Сделано».

— *Отвергнуто:*
- ждать пилота 7.1 — решение автора: `module-kit` взял пока только тег (R54 module-kit), а его T43 не сделана;
- ждать конца суда `model-2`, как велит правило эпика E5, — оно держит числа круга, а они те же;
- вся строка 7.10 одной задачей, как в тексте T123, — больше одной проверки (правило 11);
- только тег, а перестройку и перевод — после пилота, как сделал `module-kit` (R54 module-kit): решение автора — вся строка сейчас.

— *Следствие:* поправки 2026-11-12 к R02, R39 и R44; В12 ушла из таблицы 3, а В13 «в реестр spec» называет `spec` задачи строки 7.10. Мерж выпускает тег `sim` с теми же числами: «Gate» `learner-model` может остаться на v0.1.10, а суд `model-2` идёт на тех же прогонах базы. T51–T53 — каждая своей сессией. Где записано: RUN.md, «Задачи вне эпика», «Общее для задач» E5 и таблица 3; [docs/plan/E5.md](plan/E5.md); CLAUDE.md; `.claude/`; DEBT.md; README.md.
