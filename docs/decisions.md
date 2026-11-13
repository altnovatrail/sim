# sim decision log

> Why sim is built the way it is. Each entry says what was decided, why, what was rejected and where it is described.
>
> **Rule:** do not re-propose what was rejected without new facts. If a decision has to change, ask the author and update this file. An entry is never rewritten: an amendment with its date at the end of the entry corrects it, or a new entry replaces it, and the old one stays where it is, marked "Superseded".
>
> The repository was set up 2026-10-19 from the template of `spec` v0.1.0 (R21 spec), moved to `spec` v0.2.9 2026-10-20 (R01), and 2026-11-07 to v0.9.4 (R02) and the same day to v0.9.5 (R06), and 2026-11-08 to v0.9.8 (R07) and the same day to v0.9.9 (R08); the code epic E2 was laid out 2026-11-08 (R09), and its frame — the Go module, CI and gopls in the container — is recorded by R10 and R11, and the seeds and distributions by R12; the executor leads T11–T23 by themselves, up to the merge (R13); the simulator reads and translates the packages itself (R14); the simulator's learner — R15, the run — R16, the invariants — R17, the populations — R18, the synthetic packages — R19, the populations `S-`, `F-`, `W-` — R20, `A-`, `L-`, `E-` — R21, the model's measures — R22, the engine's measures — R23, the ceilings and two versions side by side — R24, the criteria and the report — R25, the gate — R26, the gate in CI of `modules` — R27, the paths of the targets and the decisive run — R28, the release — R29. Epic E2 was done 2026-11-09, and epic E3, the first candidates, was laid out the same day (R30). The packages are pinned anew — on `modules` v0.1.4 (R31). 2026-11-10 `sim` moved to `spec` v0.12.0 — format 0.4, the packages of `modules` v0.2.0 (R32). The simulator's runs are fitted into an hour: 1,000 learners a cell and the choice anew among the first 10 (R33). A candidate's decisive run is run by the same recipe, on its commit (R34). The first judgement is the one whose candidates are ready earlier, and the judgement's report lies in the repository (R35). The first judgement of the engine's candidates ended in the exit, and the base stayed (R36). The first judgement of the model's candidates did too (R37). Epic E4, «Стенд за библиотеками», was laid out the same day (R38). 2026-11-11 the base moved to `engine` v0.1.6, and the gate and the decisive run grew to 150 and 1,500 learners (amendments to R04, R31 and R33); the same day `sim` moved to `spec` v0.15.0 (R39), and a judgement began to hold the candidate's second library the same as the base's (R40). The simulator reads the check kind `time` and the leading question of format 0.4 on `learner-model` v0.1.5 (R41) and is pinned on the packages of `modules` v0.2.5 with them, with the paths and bands anew (R42); the base moved to `learner-model` v0.1.6 with the same bits, and the paths were read anew (R43); epic E5, the sizes and the model's second round, was laid out 2026-11-11 (R44). 2026-11-12 the gate grew to 250 learners, and the decisive run to 2,000 (amendments to R04, R31 and R33); the same day `sim` moved to `spec` v0.16.3, the line of epic E10 of `spec`, with the way of working by its template, and the translation of `docs/` was split into the tasks T51–T53 (R45). How it comes to differ from the template over time — by tools in the devcontainer, by rules in RUN.md — is recorded here as a decision, as R46 canvas-api recorded the differences from its model. Next number — R48.
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

— *Amendment 2026-11-13, the base on `theta` (T47; R46):* the paths are read anew on the base after the choice of the judgement `model-2`: `learner-model` at commit a063835 (`v0.1.7-0.20261113040121-a063835087c7` in `go.mod`), the same `engine` v0.1.6 (57e0ab8), the packages of `modules` v0.2.5 and 2,000 working learners a cell. The rule, the size and the resolution are those of the amendment "2,000 learners". This is calibration after a choice, not a choice: the base is what was chosen ([architecture.md](architecture.md), section 3.6). The values are measurements 3.1, pinned in `internal/criterion/paths.json` (R28) in place of the values of the amendment "2,000 learners"; from this amendment on, the table stands there and not here (the 2026-11-13 amendment to R28). 34 of the 85 paths moved, all on the populations `F-`, most of all `M3.settle` of `F-break`: 161.08 → 148.40 on olympiad maths and 336.39 → 293.84 on school maths. `held-out-model-1` is looked at: the model's next round confirms on `held-out-model-2`, whose seed an amendment to this entry writes down before the first run of its candidate.

— *Amendment 2026-11-13, the release of `theta` (T48; R46):* the paths are read anew on the release of `theta`, `learner-model` v0.1.8 (ab688b2), with the code of a063835, the same `engine` v0.1.6 (57e0ab8), the packages of `modules` v0.2.5 and 2,000 working learners a cell: the judgement holds the base to the libraries of the paths (R40). Both decisive runs are, past the header, byte for byte those of the amendment "the base on `theta`", so the values are the same to the last digit: measurements 3.2, in `internal/criterion/paths.json` (R28) with the libraries and the SHA-256 of the new runs.

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

## Packages

**R14. The simulator carries the published packages in itself — the index and the archives of the release of `modules`, embedded by `go:embed` in the package `internal/published`, with the digest of the index in the code, — and reads them itself, with the package `internal/pack`, through the chain of digests, by the types of the canonical view, and translates them into the inputs of `learner-model` and `engine` as `api` would translate them.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T12 (RUN.md, T12; R13).
- **Reading** — by the types `formatv0.ModuleIndex`, `formatv0.Entry` and `formatv0.Inventory`: the index — by the pinned SHA-256 and by the places of its members, whatever format version it is of; an entry of a format version outside the window — for now that is only 0.3 (R44, R101 spec) — the simulator skips, rather than rejecting the index ([format text](../spec/format/README.md), section 2); a version that the index marks `withdrawn` is a refusal: `api` does not start an attempt on it (R63 spec); the package — the SHA-256 of `inventory.json` equals the entry's `package_digest`, the digest and the size of `server.json`, the server part, are those in `inventory.json`, and the inventory and the module manifest name the same module, version and format version as the entry ([format text](../spec/format/README.md), sections 25 and 26). An archive lies in the folder of the release that its entry names. The packages are read once per process: the bytes are pinned, and so is what they give. The archive is read by a tar reader of the simulator's own: `archive/tar` pulls in `os/user`, whose cgo the rules of code do not check on every target.
- **Translation:** the full IDs — the module, the kind and the local ID joined by "/" ([format text](../spec/format/README.md), section 10), the anchor skills `core.` — as they are; the bands — by their place in the module manifest; the tiers — from `skill_tiers`; the graph — `learnermodel.NewGraph` (R28 spec); an answer slot — the check kind and, for a choice, the number of options, wherever they stand — beside the check or in it; a check by a rubric gives no answer slot, and a task with it is "with a rubric"; a rubric beside a verdict in one task is an error; an item's traps — from `trap_values`, as full IDs; the provenance of a key — "from the author", of a template — "checked by the solver", of tags — "from the author" (concepts 9.4 spec). Beyond the inputs of the libraries the simulator's module holds a task's canvas, its age categories and the number of its instances, and a lesson's number of segments: the learner's truth reads them ([architecture.md](architecture.md), section 3.1).
- **The age filter** — as in `api` (R107 spec): a task and a lesson suit by their own categories or, without them, by the module's categories; a lesson — only if each of its tasks suits too (architecture 3.7 api; a test of `engine`).
- **A template's instance** — `below(count)` of the place of the attempt's seed, where the task's name is its local ID, and the numbers are from `Instances` of the canonical view in the same order (R12; R60 spec).

— *Why:* the run (T14), the gate (T24) and the tests need the same bytes, and without the network: the embedded release gives them in `sim` itself, and the chain of digests from the pinned index holds them as it would hold downloaded ones (RULES.md, Hard rules; R77 spec). The translation of `engine` lives in its test, and its vectors do not catch a divergence from the translation of `sim` (T21 engine), so `sim` repeats its rules and holds the same facts of the release by a test of its own.

— *Rejected:* the packages in `testdata/` — only a test would read them; downloading them by tag in every run; `module-kit` (R62 spec; R06); `archive/tar`.

— *Consequence:* a new pin of the packages is a new folder in `internal/published` after `release-verify`, a new digest of the index and new bands ([architecture.md](architecture.md), section 3.7); how the gate at the callers takes the packages — T24. — RUN.md, T12; `internal/pack/`, `internal/published/`; [architecture.md](architecture.md), sections 3.3 and 4.

— *Amendment 2026-11-09:* the release `modules` v0.1.4 is pinned — the index SHA-256 `a244c0ba…edb2a`. From it, `olympiad-math` 0.1.1, whose archive lies in the folder `v0.1.1/`, where the index's entry names it, and `grade-3-math` and `english-essays` — 0.1.0, from the folder `v0.1.0/`. The index and the archive of `olympiad-math` 0.1.0 are removed from the embedded files. Why — R31.

— *Amendment 2026-11-10:* the release `modules` v0.2.0 is pinned — the index SHA-256 `a45c1e3d…f71f` and all three modules of version 0.2.0 in format 0.4, from the folder `v0.2.0/`. The folders `v0.1.0/`, `v0.1.1/` and `v0.1.4/` are removed: the simulator no longer reads packages of 0.3. Why — R32.

— *Amendment 2026-11-11:* the translation takes the check kind `time` — `CheckKindTime` of `learner-model` v0.1.5 — and, for a task, whether it has a leading question `leading_question`, not its words; the learner's truth of `G8` reads them. Why — R41.

— *Amendment 2026-11-11 (T43):* the release `modules` v0.2.5 is pinned — the index SHA-256 `cebf23cf…aec6a`, byte for byte the same as v0.2.1's. From it, `olympiad-math` and `grade-3-math` 0.3.0, whose archives lie in the folder `v0.2.1/`, where the index's entries name them, and `english-essays` 0.2.0 from the folder `v0.2.0/`. The index of v0.2.0 and the archives of the maths 0.2.0 are removed from the embedded files. Why — R42.

## The virtual learner in code

**R15. The virtual learner's truth is the package `internal/learner`: a population is a set of assumptions `Assumptions`, where the base population is `Base()`, and each other one changes one field. The learner draws its truth from the place `params` in a pinned order, the first answers — one `fraction` of the place `answers` per item, the answers after rungs — from `retries`. What section 3.1 did not name is decided by the simulator's assumptions: day 0 is 2026-09-07, a lesson segment without a task takes two minutes, the answers after a rung take no time.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T13 (RUN.md, T13; R13).
- **The order of `params`:** a learner has one grade, as a profile does, — evenly among the grades that each of its modules with grades names, when there is more than one (for a learner of both maths it is 3); then for each module in the order of IDs the module level — the shift of the starting band plus `StartSpread` · N(0, 1) — and the offset of each skill with a verdict in the order of the graph — `SkillSpread` · N(0, 1); then, if there is no grade, the age category — evenly among the categories that each of its modules names. With a grade the category is `under_13` ([architecture.md](architecture.md), section 3.1). Modules without a common grade or category, and assumptions with which the calendar does not go, are an error.
- **P*** — c* + (1 − c*) · R̄* · σ(a · (θ* + δ̄* − β*) + 0.5 · rung), c* is the product over the item's answer slots: a choice of k options — `ChoiceGuess`/k, no more than 1, `predicate` — `PredicateGuess`, the others — 0; δ̄* and R̄* are the means over the item's primary skills with a verdict — the item's own skills or, without them, the task's skills; β* = band step · (p − 1) + (d − 3) + `Miss` · z, z is a normal number of the place `task/<module>/<task>` with the purpose `miss`, one per task and set. R* — (1 + F · Δt / τ*)^(−C), Δt and τ* in days; before the skill's first attempt — 1.
- **The attempt:** the first answers — item by item in the order of the items; after a wrong one — up to three rungs of the ladder, on each an answer to each item still wrong; at the fourth wrong answer of an item the solution is shown. A task with a rubric is one submission and "done", 15 minutes. The time is two minutes per item of a sheet and a minute per line of a workbook. The attempt's seed is the next integer of the place `attempt`. The stability of each primary skill of the items changes once per attempt, by the share of correct first verdicts, and the last practice is the start of the attempt.
- **The calendar:** a day is a session when its `fraction` of the place `calendar` is below `SessionChance`; a session is `SessionAttempts` attempts in a row from 15:00 UTC, and the next one is issued when the previous one has closed.
- **A lesson** — segment by segment to the end: the segments without a task — two minutes each, before the lesson's tasks, and the tasks — one after another, like any. A task's miss is one per task and set for all the learners, and it is drawn once.

— *Why:* section 3.1 describes the truth by formulas, and the populations of section 3.2 change one assumption each; the set of assumptions holds each of them as one field, and a test shows that a population differs from the base one by exactly that field. The pinned order of `params` and the separate places for answers and repeats hold "the k-th first answer takes the same numbers under any rule" (section 3.4). The start day, the time of a lesson segment and of repeats section 3.1 did not name; these are the simulator's assumptions, like the speeds of the populations: the pauses depend on them, but not what the simulator compares — the base, the candidate and the ceiling go on the same learners and the same calendar.

— *Rejected:* populations as code over the learner; the truth on the code of `learner-model` (R03); a time for an answer after a rung — the repeats would then shift the calendar, and with it the R* of the next tasks, and a rule with another ladder would change the pauses.

— *Consequence:* T14 gives the learner the tasks of the engine's choice and builds evidence from the attempts; T16–T19 — the populations as fields of `Assumptions`. — RUN.md, T13; `internal/learner/`; [architecture.md](architecture.md), section 3.1.

— *Amendment 2026-11-08:* a lesson's segments go in the author's order, not "explanations first": the module holds the order of the segments (`pack.Lesson.Order`), and the run segment by segment goes through them in it. — R21.

## The run

**R16. A learner's run is the package `internal/run`: at each time of the calendar the simulator builds the engine's input, as `api` does, takes the choice, answers with the learner's truth and gives the model the evidence; each step is recorded with the input, the choice, the attempts and the evidence. The limits of a run — 200 first answers in maths, but no longer than a year, 180 days for the populations of forgetting, for the essays — until the end of the tasks with a rubric, but no longer than 60 days; "nothing to give" ends the session.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T14 (RUN.md, T14; R13).
- **The engine's input:** `Read` at the time of the choice over the learner's modules and with its grade; the pool and the lessons — by the age category ([architecture.md](architecture.md), section 3.3; R14); the lesson runs — "passed" with the time of the end; the attempts closed "Done"; the last closed attempt — the task, the lesson if the attempt is a segment of a run, and the result on each primary skill as facts; the summary of issues — for a skill, the time of its last issue by a task or by a lesson that teaches it, the code of the choice and the result of its last closed attempt, for a module the number of choices, a lesson as one, for a task the time when its solution went to the device (R104, R107 spec). The solution goes to the device when the attempt is closed by a correct answer or by the solution shown (R82 spec). The simulator gives the summary in the order of IDs.
- **Evidence** — only for a task with a verdict: the outcome — "from the first submission", "after a hint" with the rung, or "solution shown"; for each item the first and the last verdict and the wrong answers before the correct one; "the solution went to the device earlier" — except for a template, which has no such fact; the time of the last action — the end of the attempt; the starting band — the learner's starting band. The truth knows no traps, and the learner has no technique and no requests (section 3.1).
- **The step:** a choice of a task is an attempt; a choice of a lesson is the run segment by segment, the evidence of each of its tasks and the run "passed", and each skill of the lesson and each primary skill of its tasks is issued at the time of the lesson's choice; a choice is one place of the session, a lesson with tasks too; "nothing to give" ends the session, and the next choice is on the next day with a session. The step holds beside each attempt its evidence — an attempt with a rubric has none — and the number of the template's instance, chosen by the attempt's seed (R12); the "days" of a run are the day of the last choice; the record holds the model's state after each step.
- **The limits** — each has days, otherwise a run might not end: maths — 200 first answers of items with a verdict, like the 200 answers of standalone, and no longer than 365 days, so that a rule that gives nothing still ends the run; the populations of forgetting — 180 days; the essays — until all the tasks with a rubric that suit the learner's category are closed "Done", and no longer than 60 days ([architecture.md](architecture.md), section 3.1).

— *Why:* the simulator gives the engine its input as `api` would give it ([architecture.md](architecture.md), section 1): otherwise it would measure the engine on an input that does not exist in the platform's work. The record of a step with its input is needed by the invariants — `I1` repeats the choice by the input, `I2` reads the flags of the trial series and the summary (section 3.5), — and the attempts and the evidence by the measures. A year in maths is a limit that a rule that gives tasks does not reach: `G0` gathers 200 answers in 6–31 days.

— *Rejected:* a lesson as as many places of the session as it has tasks: the number of choices in the summary counts a lesson as one, and the session counts it the same way; "nothing to give" — a choice again a minute later: a choice on the same input would give the same, and time in a session without an attempt means nothing.

— *Consequence:* T15 reads the invariants from the record of the steps, T20 and T21 — the measures. — RUN.md, T14; `internal/run/`; [architecture.md](architecture.md), section 3.

## Invariants

**R17. The invariants `I1`–`I4` are the package `internal/invariant`: `Check` gives each violation in the record of a run with the invariant, the learner, the step and the choice, and `Same` — a violation of `I1` between two runs on the same seeds. The tolerance is zero.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T15 (RUN.md, T15; R13).
- **`I1`:** the choice by the input that the step stores is the same choice; two runs on the same seeds give the same choices, the same bytes of each piece of evidence (`EncodeEvidence`) and the same bytes of the model's state at the end (`EncodeState`).
- **`I2`:** a choice of a task or a lesson has one code from `ReasonCode`, "nothing to give" has none, and the code is true to its meaning by the facts that the simulator knows — the 2026-11-08 amendment to R04.
- **`I3`:** a skill of "the author's order" goes when its prerequisites with a rubric or without assessment are done — a lesson that teaches them is passed, or a task on them is closed "Done".
- **`I4`:** a task, not a template, that is the task of the last closed attempt outside a lesson, and whose solution went to the device, — not as the next task (R82, R107 spec).

— *Why:* the invariants are not statistics, one violation is enough (R04), so they read each step, not the numbers of a cell. The step's input, recorded by the run (R16), lets `I1` repeat the choice without running anew.

— *Rejected:* recomputing the reachability and the order of the queue in the simulator — the 2026-11-08 amendment to R04.

— *Consequence:* the gate (T24) holds the invariants on each learner, and `I1` also by a repeat of the run. — RUN.md, T15; `internal/invariant/`; [architecture.md](architecture.md), section 3.5.

— *Amendment 2026-11-09:* for a run with a repeat — the gate's — `invariant.CheckRepeated` checks `I2`–`I4` and does not take the choice anew at each step: `I1` is held to the bit by the repeat (`Same`), and the engine's refusal of its own input is caught by the run itself. But the first 10 learners of each cell (`bench.Rechosen`) go through `Check` with a repeat too: a repeat does not see a rule that changes its input in place in the same way in both runs. Without a repeat, in a decisive run, `Check` still takes each choice anew. Why — R31.

— *Amendment 2026-11-10:* the choice anew at each step is taken only for the first 10 learners of each cell (`bench.Rechosen`) — in every run under the engine's rule, the decisive one too, not only with a repeat. The others go through `invariant.CheckWithoutI1` — the new name of `CheckRepeated`: a decisive run has no repeat, and the name says that the check does not take the choice anew. In a decisive run the choice anew cost about 40 % of the CPU. Why — R33.

## Populations in code

**R18. The populations are the package `internal/population`: each has a name, graphs, a run limit and its difference from the base population as fields of `learner.Assumptions`; the two sides of one assumption — by the parity of the learner's number. The nineteen populations of standalone are `Standalone`, in the order of the table of section 3.2, and `G7` goes on grade 3 maths (the 2026-11-08 amendment to R03).** Decided by the executor on the author's instruction 2026-11-08 in the planning of T16 (RUN.md, T16; R13).
- **The fields of the assumptions:** `LevelShift` — the level 2.5 below or above the starting band (`G1`); `LearnLevel`, `LearnSkill`, `LearnReturn` — the growth of the level and of the offsets of the answer's primary skills, and the return of the others by a share of the way to the initial one, at each first answer of an item with a verdict, `LearnFade` — the slowing of the growth F/(F + k), k from zero, and the return of the other skills does not slow down, as in standalone (`G2`, `G2-half`, `G2-fading`); the answers are each module's own; `JumpBy`, `JumpFrom`, `JumpTo` — a jump of the level at an answer from the 50th to the 150th (`G3`, `G3-drop`); `BandBias` — the tasks of one band, drawn for the learner, are 0.75 harder (`G4-band`); `SkillHalves` — the first half of the graph's skills is stronger, the second weaker, or the other way round (`G5`); `TechniqueChance`, `TechniqueLift` — a technique on request before the first answer in a fifth of the attempts with a verdict and +0.5 to the logit of its items (`G8`); the window of the jump, the share of the technique and the fading are checked in `New`; the other populations change fields that were already there — `Slope`, `ChoiceGuess`, `Miss`, `StartSpread`, `SkillSpread`.
- **The order of the numbers:** what only some populations draw — the answer of the jump and the hard band — goes in `params` after what every learner draws, the age category too; the technique on request — by one `fraction` of the place `attempts` per attempt; the growth and the jump — by the first answers of items, after the attempt. The evidence of an attempt with a technique carries `TechniqueOnRequest`.

— *Why:* the table of section 3.2 sets each population by one assumption, and a test holds that a population differs from the base one by exactly the fields of this assumption. The parity of the number gives half of the learners to each side without an extra number, which would shift the truth. A technique on request is in NovaTrail a fact of the evidence, and the model reads it as help (R06 learner-model), so the simulator gives it in the evidence, and does not only raise the logit.

— *Rejected:* the side from `params`; populations as code over the learner (R15); the graph of `G7` "no graph".

— *Consequence:* T18 and T19 add the new populations the same way. — RUN.md, T16; `internal/population/`, `internal/learner/`, `internal/run/`; [architecture.md](architecture.md), section 3.2.

— *Amendment 2026-11-11:* the technique on request of `G8` on a task with a leading question gives the question too — a hint on request before the first error gives both the techniques and the question ([format text](../spec/format/README.md), section 19), — and the evidence carries `LeadingQuestionOnRequest`. The lift is the same, +0.5, and there is no new number from the place `attempts`. Why — R41.

## Synthetic packages

**R19. The synthetic packages are the package `internal/synthetic`: three canonical views in the form of the modules of the first slice, assembled by Go code, valid by the format's schema of the canonical view and translated by the same `pack.Translate`, with the IDs `synthetic/sheet`, `synthetic/workbook` and `synthetic/essays`, which the criterion and the gate reject, and the anchor skill of the shared taxonomy `core.multiplication-table`, which the workbook teaches and on which a skill of the sheet rests.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T17 (RUN.md, T17; R13).
- **The sheet** — four skills with a verdict on three bands of grades 1–6: a choice of four options, a number, `predicate`, a skill on the first one and on the anchor skill; tasks of difficulty 1, 3 and 5 on each band; lessons.
- **The workbook** — one band of grade 3, a chain of three skills, the first teaches the anchor skill; pages of six to ten lines of numbers and of a choice of three options; lessons.
- **The essays** — without grades, all four age categories: a skill of choice, a skill without assessment, which a reading lesson teaches, two skills with a rubric in a chain on it; a task only for adults; lessons.

— *Why:* the synthetics are for the code and the tests of the simulator, not for a judgement (R03, R04): they must give what the published packages lack — a skill on an anchor skill of another module and an answer slot `predicate` ([architecture.md](architecture.md), section 4), — and the form of the three modules, so that the simulator's code goes on them the same way. The types `formatv0` hold no required fields or enumerations, so the views are held by the schema of the canonical view from the module `spec` — a test through `github.com/santhosh-tekuri/jsonschema/v6` v6.0.3, as in `module-kit`, only in tests; the license is Apache-2.0 (the module's `LICENSE`, read 2026-11-08). The anchor skill is from the shared taxonomy: the format knows only its skills, and the synthetics do not meet the published packages. Their own IDs and `IsSynthetic` let the criterion and the gate reject the synthetics.

— *Rejected:* files `server.json` in `testdata/` — they would be kept by hand; synthetics with the IDs of the published modules — they would be easy to confuse with them.

— *Consequence:* `S-predicate` (T18) and the `A-` populations (T19) go on the synthetics, and on the published graphs their checks are "no graph", and `M9` of `A-linked` in the gate too ([architecture.md](architecture.md), section 3.7); the criterion and the gate (T22–T24) reject a module for which `IsSynthetic` holds. — RUN.md, T17; `internal/synthetic/`; [architecture.md](architecture.md), section 3.3.

— *Amendment 2026-11-09:* only `S-predicate` stays on the synthetics: there are no answer slots `predicate` in `olympiad-math` 0.1.1 either. The populations `A-` go on the published pair of maths, and the synthetic packages hold them in the tests of the code (R31).

— *Amendment 2026-11-11:* the sheet's skill on the anchor skill is answered by a time — `time`, the key `7:30`, — not by a number, and the sheet's tasks of difficulty 5 have a leading question. The guessing for a number and for a time is 0, in the truth and in the model, so the numbers of the synthetics did not shift; a new skill would have shifted the engine's choice for a learner of both synthetics, and the anchor test would no longer see a mastered anchor skill. Why — R41.

## Populations S-, F-, W-

**R20. The fifteen populations `S-`, `F-` and `W-` are `population.Own`, each by the fields of its assumption in `learner.Assumptions`; the break and the bad day are drawn from the place `calendar`, the shift and the error of an attempt from `attempts`, and "unreadable" is a count of the answers on the sheet.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T18 (RUN.md, T18; R13).
- **The fields:** `Adult`, `AdultFrom`, `AdultTo` — an adult without a grade, the first band and the level 2.5–5 higher, evenly, the category `adult`; `UnreadableEvery` — every twentieth answer of an item on the sheet is first "unreadable", and the learner writes it again; `BreakStartFrom`, `BreakStartTo`, `BreakMin`, `BreakMax`, `BreakEndBy` — a break of 14–60 days beginning from the 30th to the 60th day, ending by the 120th, and `BreakDrop` — the fall of each module's level after it — the truth reads it by the time, like the bad day, and does not change the level; `ExponentialLaw`, `LawBase` — R* = 0.9^(Δt/τ*); `AttemptShock` — a shift N(0, 1) to the logit of all the items of an attempt; `ErrorChance` — the share of attempts where each item with the primary skill of the task is "not yet" on the first answer; `BadDayChance`, `BadDayDrop` — the share of bad days with a session and the fall of the level on them; `RepeatChance`, `RepeatDays` — P* = 0.95 on a repeat of a task, not a template, whose solution went to the device within 30 days.
- **The numbers:** the break — first from `calendar`, the start, then the length; the bad day — a `fraction` from `calendar` after the decision on the session; the shift, then the error — from `attempts`, after the technique of `G8`; the adult's level — last in `params`. "Unreadable" takes no numbers: the count goes over each answer of an item on the sheet — the first one and those after the rungs; the first verdict is given by the answer written again, by its own `fraction` of `answers`, and the evidence carries the reason `UnverdictedUnreadableOther`; an answer after a rung written again changes nothing. A learner has one calendar, and "the solution went to the device" is one source too: the learner, and the run takes from it both the evidence and the summary. Assumptions with which the break does not fit, the law is not defined, the shares are outside [0, 1] or an adult lands in a module without the category `adult`, `New` rejects.
- **Graphs and limits** — section 3.2: `S-scale`, `S-adult` and `S-unreadable` — olympiad maths; `S-predicate` — the synthetic sheet (R19); `F-` — both maths, 180 days; `W-shock`, `W-error`, `W-badday` — grade 3 maths; `W-repeat` — both maths.

— *Why:* section 3.4 draws the days, the sessions, the breaks and the bad day from `calendar`, and what an attempt holds from `attempts`, "the k-th attempt takes the same numbers in the same order"; the pinned order draws the same under any rule. "Unreadable" is not an error and not a first verdict (R04 spec; R06 learner-model), so it changes nothing in the truth and takes no numbers, and it reaches the model as a reason.

— *Rejected:* "every tenth day" and "every fifth attempt" by a count — section 3.4 names a stream; "unreadable" as a share from a stream — section 3.2 says "every twentieth".

— *Consequence:* T19 adds `A-`, `L-` and `E-` the same way. In R18 the register's "`New`" is now `Own`. — RUN.md, T18; `internal/learner/`, `internal/population/`, `internal/run/`; [architecture.md](architecture.md), section 3.2.

## Populations A-, L-, E-

**R21. The seven populations `A-`, `L-` and `E-` are in `population.Own`, by the fields of their assumptions; `A-` go on the synthetic sheet and workbook together; the shares and the segment are drawn from `attempts`; true mastery is `learner.TrulyMastered`; a run abandoned at a segment is open until it is continued or 30 days end it "abandoned".** Decided by the executor on the author's instruction 2026-11-08 in the planning of T19 (RUN.md, T19; R13).
- **The fields:** `LevelCorrelation` — the z of each next module is ρ·z₁ + √(1 − ρ²)·z, ρ = 0.7 (`A-linked`); `TeacherBoost` — the skills that teach the anchor skill are 2.0 stronger (`A-unpractised`); `TeacherFades` — g = 0 for such a skill as soon as the truth has mastered it at the first band where it is taught (`A-fading`); `AbandonChance` — a fifth of the runs is abandoned at a segment drawn evenly among those still ahead (`L-abandon`); `ExplainLift` — the first explanation or example that the lesson run reaches raises the offsets of its skills by 0.25, once (`L-explain`); `LeaveChance` — a third of the attempts of tasks with a rubric is closed "left" after the same 15 minutes (`E-leave`).
- **True mastery** — the P* of the reference task, of difficulty 3 at the band, without a miss — neither the author's nor a whole band's, — a choice of five options with the population's guessing, R* = 1, no lower than 0.775 ([architecture.md](architecture.md), section 3.5): the learner's `Reference` and `TrulyMastered`. The truth does not read the model's numbers. The teaching skill and the first band where it is taught are held by the truth's skill, from `New`.
- **A lesson's segments** go in the author's order: the module holds the order of a lesson's segments (`pack.Lesson.Order`), and the run goes through them in it, not "explanations first", as R15 said.
- **Lesson runs in the run:** an abandoned one is open, without a time of the end; the next run of the same lesson continues it from the segment where it was abandoned — "passed" or open again; an open one ends "abandoned" 30 days after its start, and the time of the end is those 30 days. The engine gives an open run again at the very next choice (branch 1, architecture 3.3 engine), so it reaches 30 days only if the engine does not give it; the simulator gives the engine its input honestly, and a test holds the end "abandoned". The run holds the time of the start: the engine's `Run` has none. An attempt with a rubric "left" is a last attempt "left" and not "Done".

— *Why:* section 3.2 leads `A-` to where an anchor skill is a prerequisite of a skill of another module, and the published packages have no such thing (R19); section 3.4 draws the segment of `L-abandon` and the leaving of `E-leave` from `attempts`. True mastery is needed both by `A-fading` and by the measures of T20, so it is a method of the learner, not code of a measure.

— *Rejected:* `A-` on the published maths — there is no skill on an anchor skill there, and `M9` has nothing to be read by; "every fifth run" and "every third task" by a count — section 3.4 names a stream.

— *Consequence:* all forty-one populations of section 3.2 are in the register, and T20 reads the measures on them. — RUN.md, T19; `internal/learner/`, `internal/run/`, `internal/population/`; [architecture.md](architecture.md), section 3.2.

— *Amendment 2026-11-09:* the populations `A-` go on the pair of published maths — `olympiad-math` 0.1.1 and `grade-3-math` 0.1.0, — not on the synthetic sheet and workbook: two olympiad skills rest on skills that school maths teaches. Why — R31.

## The model's measures

**R22. The model's measures `M1`–`M10` are the package `internal/measure`: a learner's run gives sums, a cell's number is read by the sums of all its learners, and the 95 % interval from 2,000 resamples of the learners from the place `bootstrap` with the cell's name; what a measure reads about an attempt the run records as an observation `run.Obs`.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T20 (RUN.md, T20; R13).
- **Sums and numbers.** A learner's run gives named sums — squares, counts, shares by bins; a cell's number is the ratio of the sums of all its learners, not the mean of the learners' numbers, so a learner with more answers weighs more. The interval is the 51st value from below and the 51st from above of the number over 2,000 resamples of the learners, 2.5 % on each side; the resamples are drawn from the place of the working seed with the names `bootstrap` and the cell's name — the population and the graph, without the rule, — so the base, the candidate and the ceiling on the same learners take the same resamples, and a pair can be read by learner ([architecture.md](architecture.md), section 3.6). A number that cannot be read on even one resample is unread: its interval would stand only on the resamples where the case came up; on a graph with a synthetic module the numbers are marked "no graph" (R19).
- **The observation of an attempt.** Before each attempt with a verdict the run reads the model at the time of its start: the P of each item with retention and without — P without retention is c + (P − c)/R̄, where the model computes R̄ and c by its own formula ([architecture.md](architecture.md), section 4), — the pause since the last practice of the item's primary skills and the module's skills; after the attempt — the P* of the first answer, the help and the module's skills again. A skill in the observation is the model's θ + δ and the truth's θ* + δ*, the "mastered" of the model and of the truth at the model's mastery band or at the first band where the skill is taught, and the P and P* of the reference task of this band. In a lesson each task is closed as soon as the learner has answered it — the lesson run calls the run around each task (`learner.Each`), — so the model and the truth are read before and after it as they are then, with the lesson's tasks before it and without the tasks after. An attempt with a rubric is not observed.
- **The truth of the measures** is the one that answers: θ* after the break of `F-break` is lower by the fall, and with it are read θ* + δ*, true mastery, the P* of the reference task and command of the anchor skill. The bad day of `W-badday` does not enter the truth of the measures: it lowers the answers of that day, not what the learner knows.
- **The count.** "An answer" is the first answer of an item with a verdict; "the skill's answers" are the items of attempts from which the model took evidence of the skill. A pause of a skill of no less than 14 days is an attempt that answered a skill whose practice was 14 days ago or earlier; with it the window of the break of `M4`, of 7 days, and the catching up of `M5` open. Help in `M2` is a technique on request and a repeat of a task with a seen solution: the rungs of the ladder come after the first verdict and do not change it. A "mastered" of `M6` declared by the same attempt that gave true mastery waited zero answers. The jump of `M3` is the module's `JumpAt` for `G3` and `G3-drop`, and the first attempt after the end of the break for `F-break`, and not for the other `F-`, where the break does not drop the truth; a run that ended before the jump has no jump. The answers until "holds" in `M3` are counted until the first answer of a series of ten, so an error that holds from the first attempt after the jump gives zero. `M7` is the root of the mean square of the module's θ − θ*: a signed measure goes into the criterion by its absolute value (section 3.6). `M9` is read at the first choice of each day. "How many declared" in `M4` — per learner.

— *Why:* the numbers of a cell from sums add up over the learners, and a resample of the learners is a resample of their sums, so the interval is computed without a second run. The model itself reads P without retention only inside a step and does not give it out, and the simulator repeats its formula by what `Read` gives — the retention of the skills and the flag of the trial series.

— *Rejected:* a cell's number as the mean of the learners' numbers — a learner with three answers would weigh as much as a learner with two hundred; resamples by the name of the rule — the pair of the base and the candidate would take different resamples; the interval over the resamples where the number is read — it would be narrower and biased; the truth without the fall after a break — a model that correctly follows the fall would show an error; "the skill's answers" as all the items of the module's attempts — a skill that was not answered would accumulate answers.

— *Consequence:* T21 adds `E1`–`E5` the same way, and T23 reads the difference of a pair on the same resamples. — RUN.md, T20; `internal/measure/`, `internal/run/`; [architecture.md](architecture.md), section 3.5.

*Amendment 2026-11-08, in T22 (R24):* the answers until "holds" in `M5` are counted the same way as in `M3` — until the first answer of a series of five, so an error that holds from the first answer after a pause gives zero. Found on the oracle: it catches up at once, and the measure gave it five. The difference of a pair of cells on the same resamples is read already by T22 (`measure.ReadPair`), not by T23.

## The engine's measures

**R23. The engine's measures `E1`–`E5` are in the package `internal/measure`, the same way as the model's measures (R22); what the truth says about the pool at a choice — the pool's ceiling, whether there is anything to give, the reachable skills — the run records before the engine's choice, `run.Pool`; a skill's reachability is read by the truth.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T21 (RUN.md, T21; R13).
- **The pool at a choice.** Before calling the engine the run reads each task of the pool that the rule could give: without a rubric, with an item with a verdict, not a task whose solution went to the device less than 24 hours ago, — and if there are none, a second pass, where they are included, except the task of the last attempt (R82 spec; R03 engine and its 2026-11-07 amendment), — and all of whose primary skills with a verdict are reachable. The ceiling is the largest share of items with P* in the corridor 0.70–0.85 in one such task; "the pool gives" — such a task exists; the reachable skills are the primary skills with a verdict of such tasks.
- **Reachability by the truth.** A skill is reachable when the truth holds its prerequisites: a skill with a verdict the truth has mastered at the first band where it is taught; a skill without a verdict or of a module that the learner does not study holds — there is nothing to read it by; an anchor skill the truth holds as in `M9` — a skill that teaches it is mastered with R* at the time, — and if no module of the learner teaches it, it holds.
- **P\* of the issue** — the truth's chance at the time of an item's issue, without what only the attempt holds: the technique, the shift, the error and the repeat (`ItemObs.PGiven`).
- **The count.** `E1` — the shares of the issued items by the P* of the issue, the ceiling — the mean over the choices of a task or a lesson. The skills of a choice are the lesson's skills or the task's primary skills, as `learner-model` and `engine` count them: the items' own skills, and the task's skills when an item takes them. `E2` "being stuck" — 50 choices in a row, each on the skill, without a change of its "mastered" of the model, and at each there is a reachable skill on which the choice does not stand; "nothing to give" — among the choices where the pool gives. `E3` — a repetition is a choice with the code "repetition", the skill is the skill of the choice, for which the model reads "mastered" and whose review date has come, R* and "lateness" — by the observation of the attempt (the truth's R* is added to `run.SkillObs`); the share of repetitions — by weeks of the calendar from the first to the 53rd; R* at the end of the run — at the time of the last choice or of the close of its last attempt — for the skills mastered by the truth at the first band where they are taught. `E4` — the share of the skills reachable at least once that were issued, and the share of the skills with a verdict mastered by the truth by the end. `E5` — per learner: the tasks with a rubric written, each once, and the attempts of a task already closed "Done" — as a task or in a lesson.

— *Why:* reachability is what the engine is judged by, and reading it by the model would mean judging the rule through its own eyes; the truth is what the simulator knows. The ceiling and reachability require the truth at the time of the choice, and by the end of the run the truth is already different, so the run writes them, not the measure.

— *Rejected:* reachability as `engine` reads it — by the model's P no lower than 0.70 (`skills.go` of `engine` v0.1.0): this is a rule of the engine, a candidate may change it; the P* of the issue with the noise of the attempt — no rule foresees it; the ceiling only by the choices of a task — a lesson is an issue too.

— *Consequence:* T22 builds the engine's ceiling on the same pool — by the P* of the task with the largest share in the corridor. — RUN.md, T21; `internal/measure/`, `internal/run/`; [architecture.md](architecture.md), section 3.5.

## The ceilings and two versions side by side

**R24. The model's oracle is read on the answers of the same run as the model; the engine's ceiling is a rule of the run's own on the same pool, and where the pool gives no task with P\*, `engine` chooses; each build of `sim` writes a file of results — JSON Lines, a line per learner, — and the comparison reads two files in pairs by learner; the file names a library by the version from the build, and one built from the workspace by the commit that the caller passes.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T22 (RUN.md, T22; R13).
- **The model's oracle** — the same observations, where the estimate is the truth θ* + δ*, P is the P* of the answer, P without retention is P* too, "mastered" is true mastery, the P of the reference task is P*, the level is θ*, and it holds the anchor skill when the truth holds it (`measure.Oracle`). The model is not replaced for the engine: the `Readout` of `learner-model` holds the state inside and computes the P of tasks itself, so the oracle cannot be given to the engine; and on the same answers the oracle and the model are compared by learner.
- **The engine's ceiling** — the rule `run.Ceiling` instead of `engine.Choose`, on the pool `run.Pool` (R23): the task with the largest share of items with P* in the corridor, at a tie — the one with the mean P* closer to 0.775, then by ID, with the pause of 24 hours and the rule's second pass; a repetition — when for a skill mastered by the truth the P* of the reference task of the first band where it is taught falls below 0.775 from forgetting, — the best task on it, with the code "repetition". Where the pool gives no task with a verdict — tasks with a rubric, lessons, "nothing to give", — `engine` chooses: the ceiling is the ceiling of P*, and where there is no P*, the rule stays. The invariant `I2` does not apply to the ceiling: its codes are not branches of `engine`.
- **The file of results** — JSON Lines: the first line is the set of learners, the rule (`engine` or `ceiling`) and the version or the commit of each library; then a line per learner — the population, the graph, the number, the sums of the measures of the model and of the engine and the sums of the oracle on the same answers. JSON writes float64 numbers in the shortest form that reads back to the same bits, so a pair is compared to the last bit. The commands are `sim run` (the populations, the learners per graph, the set `working` or `held-out-<criterion>-<n>`, the rule, `-commit`) and `sim compare` (two sides, each the sums of a rule or of the oracle): for each cell and number — a, b and the difference a − b with a 95 % interval from the same resamples in both (`measure.ReadPair`); the oracle's side reads only the model's measures — the oracle has no measures of the engine, — and for the engine's ceiling `E3` is not read. The file `-out` is written beside and takes its place only when the run has ended, so a failed run does not touch the previous file. The learners go in parallel, and the lines in the order of the populations, the graphs and the learners, so the file of one configuration is the same bytes.
- **The name of a library** — from the build information (`runtime/debug`): the version by which the build took it — a tag or a pseudo-version, where there is a commit, — or the version of its replacement, with the path of the replacement, when it is another module. A library built from the workspace `go.work` the build does not name by a version, and its commit is passed by the caller, `-commit <module path>=<full SHA>`, — the gate's recipe knows it from the gate file ([architecture.md](architecture.md), section 3.9); without it `sim` writes no file. The commit is 40 hexadecimal characters and only for one of the two libraries; a commit of a library whose version the build knows is rejected — the file would have two names of one build.

— *Why:* one Go build does not hold two versions of one library ([architecture.md](architecture.md), section 3.6), so a pair is two files; JSON Lines is read line by line, and a learner is one line. The commit in the file is what the gate checks the numbers' recording against.

— *Rejected:* the oracle as a separate run where the engine reads the truth — the model cannot be replaced for the engine, and a run with other choices would compare the model and the trajectory at once; the engine's ceiling for tasks with a rubric and for lessons — they have no P*; CSV — the learners have different sets of sums; the name of a library through `-ldflags -X` — every caller would pass it even where the build knows the version.

— *Consequence:* T23 judges the candidates by `ReadPair` of the base, the candidate and the ceilings; T24 writes `go.work` and passes `-commit`. The values of `E3` for the engine's ceiling are not read — `E3` goes to perfection ([architecture.md](architecture.md), section 3.6). — RUN.md, T22; `internal/bench/`, `internal/measure/`, `internal/run/`, `main.go`; [architecture.md](architecture.md), section 3.6.

*Amendment 2026-11-08, in T24 (R26), by the author's decision:* "the learners go in parallel" was an error of T22. `bench.Run` handed the learners out to goroutines through `sync.WaitGroup.Go`, and R09 rejected goroutines in the code: a run is split into processes by learners, and the recipe splits it. The rule `goroutine` catches the statement `go` and did not catch `WaitGroup.Go`. The goroutines are removed. A run takes a part of the learners (`bench.Part`, `sim run -part k/n`), the recipe launches the parts as processes, and `sim merge` joins their files into the file of the whole run, the same to the byte.

## The criteria and the report

**R25. The criteria are the package `internal/criterion` and the command `sim judge`: the judgement file names the criterion, the runs of the base, of the engine's ceiling, of the candidates and of the confirmation; a check stands where the base reads the number; the path of `E1` in the model's criterion goes to the engine's ceiling; the score adds up the k-th resamples of the cells.** Decided by the executor on the author's instruction 2026-11-08 in the planning of T23 (RUN.md, T23; R13).
- **The judgement file** — JSON: the criterion (`model`, `engine` or `pair`); the files of results of the base (the rule `engine`) and of the engine's ceiling (the rule `ceiling`) on the same learners; the candidates — the name, the file, its own score group, the model's target (`lag`, `corridor`, `settle`), the number of changed constants, a new field, the distance from the base in the author's order, and in a pair — the library whose target it is; the confirmation — three runs on held-out learners. The model's oracle is the sums of the oracle in the file of the base (R24). The report is Markdown, a tool's output, so it is in English (R14 spec).
- **Direction and tolerances.** Signed measures — by their absolute value (`M3.lag`, `M8.shift`); shares where more is better — by their shortfall to one (`E1`, `E3` R*, `E4`); the count of written tasks `E5` — with the sign changed. The tolerances are the author's ([architecture.md](architecture.md), section 3.6); for the measures that the author's list does not name — the tolerance of the nearest in kind: the shares of `M5` and `M7` — 1 point, like the shares of `M3`; the shift of `M8` — 0.01 logit, like `M1`; the Brier of `M10` — 0.002, like `M2`; the share of hard first items of `M7` — 1 point. The simulator's resolution is 4.3 standard errors of the difference, and the standard error is the deviation of its values over 2,000 resamples.
- **Checks.** "No worse" — the worse end of the interval of the difference "candidate minus base" is no greater than the tolerance. The floor — the worse end of the candidate's interval: `M4` no greater than 0.20, `M6` no greater than 1.5 times the base's value. The targets — `lag` and `corridor` by the lower end of the interval of the share of the path, `settle` by the upper end of the share of those not caught up. A check and a floor stand where the base reads the number: where even the base does not read it, there is nothing to hold, and the report names such cells rather than failing the candidate. A number that the base reads and the candidate's run does not is a check not met. The invariants `I1`–`I4` the run writes into the learner's line under the rule `engine` (`bench.Line.Breaches`), and the file's header says that they are recorded; a file without this mark does not hold the invariants. `I1` between two runs is the gate's (T24). A floor and a target whose population is not in the runs are not met: a check that was not measured does not hold. In a pair a check common to both criteria stands once. A model's target for a candidate of the engine is rejected: the engine has no number for a target.
- **The score.** A measure's share of the path on a cell is (base − candidate)/(base − ceiling) in its direction, the path — by the whole cell. The ceiling is perfection (zero) for the measures where the path goes to perfection; the engine's ceiling for the measures `E`, in the model's criterion too — the oracle cannot be run through the engine (R24), and the engine's ceiling is the ceiling of the corridor of any pair of a model and an engine; the oracle for the other measures of the model. A measure without a path — the path is no greater than zero, or the base or the ceiling does not read it, — leaves the score, and the report names it; a measure that the candidate's run does not read leaves the score unread. A population's value is the mean over its cells and measures, so the graphs of a population share its weight equally. The interval is from 2,000 resamples: the k-th resamples of all the cells are added with weights, and each cell has its own resamples — its own places.
- **The choice and the exit.** "Better" — the lower end of the interval of the score is above zero, and of its own group too. Equal to the best — those whose difference of scores with it on the same resamples holds zero. The simplest — the fewest constants, then without a new field, then the nearest to the base, then by name. The confirmation is the same judgement on the runs of held-out learners: the base and the candidate under `engine`, the ceiling under `ceiling`, the set held-out and not the one on which the choice was made; without it the choice waits for confirmation. The report calls a run decisive only if the confirmation is decisive too. A pair — the constraints of both criteria, "better" by the score of its own library, and the upper end of the interval of the other's score no lower than zero. A decisive run — when each cell of each run has no fewer than 4,000 learners.
- **The slice of `M2` after pauses from 7 days** (`M2.brier.after7`) is added for the group "Forgets": section 3.6 names it, and the slices of `M2` divided the pauses into intervals.

— *Why:* a check that even the base does not read speaks of the graph, not of the candidate: on school maths for 40 learners "mastered" is declared so rarely that not every resample sees it, and any candidate would fail. Adding the k-th resamples gives the interval of the sum, because the resamples of the cells are independent.

— *Rejected:* a check that fails when even the base does not read it — it fails the base against itself too; the path of `E1` to the oracle — the oracle does not lead the engine; a population's weight on each of its graphs in full — a population with two graphs would weigh twice as much.

— *Consequence:* the values of the base and of the ceiling of which the path is made up the judgement for now reads from its own runs: they can be pinned by an amendment to R04 only after the first decisive run of the base ([architecture.md](architecture.md), section 3.6), and that run is for T24. — RUN.md, T23; `internal/criterion/`, `internal/bench/`, `internal/measure/`, `main.go`; [architecture.md](architecture.md), section 3.6.

— *Amendment 2026-11-09:* the judgement no longer reads the path from its own runs. It takes the path from the file of paths that the first decisive run of the base pinned (R28; the 2026-11-09 amendment to R04), — at each choice and in each confirmation. So the judgement file does not name a run of the ceiling, neither for the choice nor for the confirmation, and `Input` and `Confirm` do not carry it. A candidate's share is still its paired difference with the judgement's base, on the judgement's learners, but divided by the pinned path. The decisive run went in T26, not in T24, as "Consequence" said. — RUN.md, T26; `internal/criterion/`, `main.go`; [architecture.md](architecture.md), section 3.6.

— *Amendment 2026-11-11 (T40; R40):* a judgement refuses a run that is not of its own judgement. The base and the base of the confirmation are on the libraries on which the paths were read, and the base is on the versions of the embedded `gate.json`. The second library of a candidate of one library is the base's version or the commit of its pin.

## The gate

**R26. The gate is the file `gate.json` at the root of `sim`, the recipe `just gate`, the command `sim gate` and the job "Gate" on amd64 and arm64. go.mod names the libraries by the same point as the file. Until T26 the gate holds the bands of its own run on the base. `I1` is a repeat of each learner of the gate.** Decided by the author 2026-11-08 in the planning of T24 (RUN.md, T24) — the goroutines, go.mod, the bands until T26 and `I1`, — and the form of the file and of the checks by the executor.
- **The gate file** — JSON:
  - `libraries` — by module path: the repository, the full SHA of the commit and the tag of the release on it, empty while there is no release;
  - `modules` — the release, the digest of the index and the packages (module, version, package digest);
  - `dependencies` — the chosen versions of the shared dependencies, now the Go module `spec`;
  - `go` — the line `go` of each of the three modules;
  - `bands` — the reason of the recording and the band of each number: the cell, the number, whether it was read, the value and the half-width.

  The bands lie in the same file, as [architecture.md](architecture.md), section 3.9, says. A number that the recording did not read is recorded as "not read", and the gate holds exactly that: read later, it is outside the band. Section 3.7 calls an unread number "outside the band" for a number that was read at the recording.
- **The recipe `just gate [bands|exact|record]`:**
  - takes each library by the commit from the file into a temporary folder (`git fetch` by SHA), and the caller substitutes its code by a folder — `GATE_LEARNER_MODEL` or `GATE_ENGINE`;
  - checks the tag from the file against the commit from the file (`git ls-remote`): the job "Go" builds the tag of `go.mod`, and the gate builds the commit;
  - writes `go.work` with `use` of the three modules, sets `GOWORK` and builds `sim` once;
  - runs the gate as processes, a part of the learners per core (`sim gate run -part k/n`), with the commits that it sees in the folders;
  - then `sim gate check` or `sim gate record` joins the parts. A part that failed stops the others. The reason of the recording goes into the command in the shell's quotes.
- **The file of results** also names the populations, the learners per graph and the repeat of the run, so `sim merge` joins the parts of one run in its order and rejects the parts of another. Recording the bands rejects a run not on the commits of the file: bands of the caller's code are not a recording.
- **Checks** — each divergence gives a red gate with a reason:
  - the libraries are built from the workspace;
  - the chosen version of `spec` (`debug.ReadBuildInfo`) equals the recorded one;
  - the lines `go` of the three modules equal the recorded ones and are no newer than the Go that the gate is built with;
  - the go.mod of `sim` requires each library by the tag from the file, and while there is no release — by the pseudo-version of its commit;
  - the packages that the simulator embeds are those that the file names;
  - no learner violated `I1`–`I4`;
  - each number is in its band. `sim`'s own CI holds the numbers as recorded to the bit (`exact`), the callers by the bands.
- **Cells and numbers** — by section 3.7, in code (`internal/gate`): the first 300 working learners of each cell of the gate's populations, numbers without a direction. The half-width is half the 95 % interval from the same 2,000 resamples. `A-linked` goes only on the synthetic packages (R19), and the gate does not take a synthetic module ([architecture.md](architecture.md), section 3.3). So the gate does not run `A-linked`, and its report calls its numbers "no graph". A run with a synthetic module is red.
- **The bands until T26** are recorded by a run of the gate on `learner-model` v0.1.0 and `engine` v0.1.0, with the reason "the first record", on the base. T26 records them anew with its own reason.
- **`I1`** — each learner of the gate is run twice (`bench.Config.Repeat`), and `invariant.Same` checks the two runs against each other to the bit. If a measurement on the runner does not fit into the budget, the repeat is cut by an amendment to this entry.

— *Why:*
- The processes keep R09 and split a run over the cores without goroutines.
- go.mod at the same point as the file is the only way to build the job "Go" of the shared `go.yml` on a commit that is not in a release yet. The Go proxy serves a commit by SHA, and the job "Go" builds exactly the code on which the bands are recorded.
- The bands of its own run give the gate something to hold already in T24, not only from T26.

— *Rejected:*
- the bands in a separate file — section 3.9 puts them into the gate file;
- goroutines with an amendment to R09 — the author kept R09;
- go.mod only as a lower bound — the job "Go" would be red on steps 1–4 of the order of releases;
- a repeat of a part of the learners at once — the measurement has not said yet that a full repeat does not fit;
- `jq` in the recipe — it is not in the stage `toolchain`, and `sim gate libraries` reads the file with the same code as the gate.

— *Consequence:* the callers take the same recipe with `GATE_*` — these are tasks of their repositories. The threshold of the CI of `modules` — T25. The decisive run is split into parts by the same `-part` and joined by `sim merge` — T26. — RUN.md, T24; `gate.json`, `internal/gate/`, `internal/bench/`, `internal/published/`, `gate.go`, `main.go`, `justfile`, `.github/workflows/ci.yml`; [architecture.md](architecture.md), sections 3.7 and 3.9.

— *Amendment 2026-11-09:* the gate has 100 learners a cell, not 300, and on the repeat of the run for `I1` the choice anew at each step is taken only for the first 10 learners of a cell. On `olympiad-math` 0.1.1 the engine's choice takes tens of times longer, and the gate on 300 learners did not fit into `timeout-minutes: 15`. The gate is cut by an edit of `sim`, not by time, as [architecture.md](architecture.md), section 3.9, says. `A-linked` is no longer "no graph": it goes on the pair of published maths. Why — R31.

— *Amendment 2026-11-10:* the gate file names the libraries v0.1.4 — `engine` at ddfed3a, `learner-model` at acfc4f2, — the packages of `modules` v0.2.0 and `spec` v0.12.0. `just gate exact` matched on them to the bit, so the bands are the former ones. Why — R32.

— *Amendment 2026-11-10:* the gate's workspace is written by the hidden recipe `workspace`, one for the gate and for the decisive run of a candidate (R34). The gate takes from `GATE_LEARNER_MODEL` and `GATE_ENGINE` not only the folder but also the commit by full SHA. The check of the tag, the gate itself and its numbers are as before: `just gate exact` on the new recipe matches to the bit.

— *Amendment 2026-11-11 (T41):* the gate file names `engine` v0.1.6 at 57e0ab8; `learner-model` — still v0.1.4 at acfc4f2, the packages — `modules` v0.2.0, `spec` — v0.12.0. On 100 learners `just gate exact` matched on it to the bit, then the gate grew to 150 learners, and the bands were recorded anew by `just gate record` with a reason (the 2026-11-11 amendment to R31).

— *Amendment 2026-11-11 (T43):* the bands are recorded anew on `modules` v0.2.5 — `just gate record` with the reason T43, on 150 learners, and `just gate exact` is green. All 69 numbers are read where they were read before; 39 did not change in a single bit — 26 of school maths, 6 of olympiad maths, 3 of the pair of maths and all 4 of the essays. 23 numbers of olympiad maths shifted, 6 of school maths — of `G0`, `G2`, `G2-half` and `F-model`, in the third to fourth digit: a segment of a worked example was added to the time of the lessons — and `M9.notFalse` of `A-linked`. The share of hard first items `E1.hard` on olympiad maths shifted most of all — by 1.9 of the former half-width for `G0` (0.320 → 0.374) and by 1.5 for `G2-half`, and `E2.stuck` became non-zero for `G0`, `G1` and `G2` of olympiad maths (up to 0.067): 30 choices of one of five became checks of time, and they no longer have guessing. Why — R42.

— *Amendment 2026-11-11 (T45):* the gate file names `learner-model` v0.1.6 at 7d796e2. `just gate exact` matched on it to the bit, all 69 numbers, — and the bands were recorded anew by `just gate record` only for the sake of the reason, as in T42: the numbers and the half-widths are the same. Why — R43.

— *Amendment 2026-11-12 (T46):* the gate goes on 250 learners a cell, not on 150, and the bands are recorded anew by `just gate record` with the reason T46 on the same libraries and packages; `just gate exact` is green — 1 min 52 s by the clock and 29.6 CPU minutes in the container, at a load average of about 40–60 from neighbouring containers. Of the 69 numbers, 25 — with a zero band — are the same. Three numbers that the recording did not read on 150 learners are now read: false "mastered" `M4.false.all` and the lateness `M6.late` of `G0` of school maths, and the share of late repetitions `E3.late` of `F-model` of school maths; two stayed unread — `M4.false.pages` of `W-shock` and `M9.hasFalse` of `A-linked`. Of the 42 numbers with a non-zero band, three — those read for the first time — have no former half-width; for 38 of the other 39 the new value is less than one former half-width away from the former one — most of all `E5.written` of `E-base`, by 0.60, — and the half-width is on average 0.79 of the former one, about √(150/250) ≈ 0.77. The thirty-ninth is the share of pauses after which the model never catches up with forgetting, `M5.never` of `F-model` of school maths: on 150 it was 0 with a zero band, and it became 0.0036 ± 0.0055: among 250 learners there was a pause after which the model did not catch up with forgetting, and among the first 150 there was none. This is a rare event, and the former zero was a zero of the sample, not of the measure. Why the gate grew — the 2026-11-12 amendment to R31.

— *Amendment 2026-11-13 (T47):* the gate file names `learner-model` at commit a063835 without a version — `theta`, chosen and confirmed by the judgement `model-2` (R46) —, and the bands are recorded anew by `just gate record` with the reason T47 on 250 learners; `just gate exact` is green. Of the 69 numbers 61 are the same to the bit. Eight moved, all of `F-model`, each by less than one former half-width: most of all `E2.stuck` of olympiad maths, 0.240 → 0.208 (0.62 of the former half-width), and `M5.caught` of olympiad maths, 3.53 → 3.76 (0.54); the half-widths are within a tenth of the former ones. The release of `theta` puts its tag in place of the commit (RUN.md, T48).

— *Amendment 2026-11-13 (T48):* the gate file names `learner-model` v0.1.8 at ab688b2, the release of `theta` with the code of a063835 (R46). `just gate exact` matched on it to the bit, all 69 numbers, — and the bands were recorded anew by `just gate record` only for the sake of the reason, as in T45: the numbers and the half-widths are the same.

**R27. The gate of the CI of `modules` is the mode `just gate modules` and two commands: `sim gate changed` decides whether what the simulator reads has changed, by the canonical text of the translated module, and `sim gate thresholds` holds the invariants and the thresholds. The gate's numbers are extended by `E1.hard` and `E2.nothing`, from which the thresholds are taken.** Decided by the executor in T25 (RUN.md, T25), except for the thresholds of `E1`: the author decided them — the 2026-11-09 amendment to R05.
- **What is compared.** The canonical text is the JSON of the translated module, `pack.Module.Text`: the bands and the audience of the module manifest, the graph and the tiers, the tasks with their items, answer slots, traps, canvas and age, the numbers of the instances of templates, the lessons. This is exactly what `pack.Translate` takes, so the list of what the simulator reads is one place in the code, as [architecture.md](architecture.md), section 3.9, says. The texts for the learner and the module's version do not enter the text.
- **What it is compared with** — the newest published version of the module that is not withdrawn: `api` does not serve a withdrawn one (R63 spec). The states:
  - `new` — there is no published version, or no release at all (`-first`); the whole module has changed;
  - `changed` and `same` — the texts differ or are the same;
  - `transition` — the published version has another format version; the simulator does not run the module, even with `-all` (the 2026-11-07 amendment to R103 spec);
  - a package of a format version that the simulator does not read, without a transition — a refusal: a new tag of `sim` is needed;
  - a broken package — the server part is not the one in the inventory, or the module manifest is not the right one, — a refusal, in a transition too: the inventory holds `server_part` in all format versions ([format text](../spec/format/README.md), section 2).
- **The check of the packages.** A published one — by the chain of digests from the index that `release-verify` accepted: `inventory.json` equals the package digest in the entry, `server.json` — the inventory. A package of a pull request has no entry: `server.json` — against the inventory, the module manifest — against the inventory. `main` reads the files, and `pack` gets the bytes (`pack.ReadFiles`).
- **What is run** — the cells of the gate each module of whose graph is among the packages that run (`bench.Config.Packages`). The header of the file of results names these packages: module, version and digest. `A-linked` is still "no graph". A module on which not one population of the gate runs the report names, but the gate does not turn red from it: populations on a new graph are an edit of `sim`.
- **The thresholds** are point values: the gate's learners are always the same. The edges are inclusive, an unread number is red. The values for the thresholds of `E1` and `E2` are taken from the bands of `gate.json` of the same tag:
  - `E2.stuck` and `E2.nothing` — the largest over all the cells plus 1 point;
  - `E1.hard` and `E1.corridor` — by the cells of `G0` and `G2-half`, by the 2026-11-09 amendment to R05;
  - `M4.false.all` on `G0` — no greater than 0.20, `E5.again` on `E-base` — zero.
- **The check of the run.** The libraries are the commits of the file. All the cells of the packages are there, with all the learners and the repeat for `I1`. The checks of the build are the same as in `sim gate check`. `sim gate check` and `sim gate record` reject a run on substituted packages.
- **The gate's numbers** are extended: `E2.nothing` — on all populations, `E1.hard` — on `G0` and `G2-half`. `gate.json` is recorded anew with the reason T25. The 47 former bands did not change in a single bit, 18 were added.
- **The recipe** — `just gate modules <arguments of sim gate changed>`. First `sim gate changed`: if there is nothing to run, the recipe exits with 0 without touching the network. Otherwise — the same build through `go.work` from the file, the parts `sim gate run -package …` and `sim gate thresholds`. Relative paths are from the caller's folder. The form of the call for the job of `modules` — README, "Gate in the CI of modules".

— *Why:*
- The text from `Translate` cannot diverge from what the simulator really reads: a field that `Translate` begins to read enters the comparison by itself.
- The arguments are as in `module-kit diff` (R11 modules): the CI of `modules` already holds the index accepted by `release-verify` and the published packages as folders.
- A point, not an interval: the gate has no randomness between runs, and an interval would shift the threshold by noise that does not exist.

— *Rejected:*
- a projection of the canonical view by a list of fields — a second list beside `Translate`, which will diverge from it;
- a comparison with a withdrawn version — it is not served, and a comparison with it would say "same" when the learners get something else;
- the gate red on new modules until there is a population on them — `modules` could not add a module until `sim` pins it, and it pins only what is released.

— *Consequence:* the gate's job in the CI of `modules` will be connected by the session of `modules` by the README — Д5 in table 4 of RUN.md `modules`, after the first tag of `sim` (T27). — RUN.md, T25; `internal/gate/modules.go`, `internal/gate/gate.go`, `internal/pack/`, `internal/bench/`, `gate.go`, `justfile`, `gate.json`, README; [architecture.md](architecture.md), sections 3.7 and 3.9.

— *Amendment 2026-11-09:* `sim gate changed` gives, together with the changed modules, an unchanged one too, when the graph of a population of the gate holds it together with a changed one: so `A-linked` goes on both maths from the pull request if at least one of them has changed ([architecture.md](architecture.md), section 3.9). On stdout its folder goes as the line `along <folder>`, the recipe passes it to `sim gate run -along`, and such a package goes only in the cells that it shares with the changed one (`published.Package.Along`), not in all of its own. The line of such a module on stderr names the module that brought it. A cell of the gate that lacked the package of the second module of its graph the report names without turning red. In code — `gate.Runs` and `partial`. Why — R31.

— *Amendment 2026-11-11:* the translation reads whether a task has a leading question, and `time` as a check kind, so the canonical text holds both: `sim gate changed` sees an added or removed question and a check `number` that became `time`, but not an edit of the words of the question — the simulator does not read them. Why — R41.

**R28. The paths of the targets are pinned by the file `internal/criterion/paths.json`, embedded in the simulator: `sim paths` writes it from the decisive runs of the base and of the engine's ceiling, and a judgement takes the path from it. The decisive run is run by the recipe `just decisive` on the development machine, and the files of its results do not lie in the repository. The first file of bands is recorded on the commits of the releases of the base.** Decided by the author 2026-11-09 in the planning of T26 (RUN.md, T26) — the file of paths, the files of the run, where to run it and the full table in the amendment to R04, — the rest by the executor.
- **The file of paths** — JSON:
  - `libraries`, `set`, `learners` — what the paths were read on;
  - `base` and `ceiling` — the SHA-256 of the files of results of the two runs;
  - `paths` — the cell, the measure, whether the path was read, the value of the base and the value of the ceiling, both in the criterion's direction, so the path is the base minus the ceiling.

  The pairs are each measure of the score groups of both criteria on each cell of a published graph where the population has a weight. The path of `E1` is common to the two criteria. The ceiling is zero for the measures "to perfection", the engine's ceiling for the measures `E…`, the base's oracle for the measures `M…`, as the judgement computed it before.
- **`sim paths -base <file> -ceiling <file> [-out …] [-table]`** writes the file of paths and prints the rows of the table for the amendment to R04. A refusal if:
  - the runs are not under the rules `engine` and `ceiling`;
  - the set is not the working one;
  - a part of a run is given;
  - the run went on packages instead of the pinned ones;
  - the base lacks some population of the score groups;
  - the runs have different libraries or learners;
  - the base has a violation of an invariant;
  - a cell has fewer than 4,000 learners.
- **The judgement** takes the path from the embedded file: a candidate's share is its paired difference with the judgement's base, divided by the pinned path. A measure that is not in the file, that the file did not read or whose path is no greater than zero is not in the score, and the report names it. Without pinned paths the judgement refuses. The judgement no longer needs a run of the ceiling — the 2026-11-09 amendment to R25.
- **A test** holds the table of the 2026-11-09 amendment to R04 equal to `Table()` of the embedded file: the value in the log and the value the code judges by are one.
- **The recipe `just decisive engine|ceiling <file>`** builds the simulator by `go.mod`, where both libraries are v0.1.0. It runs `sim run -populations all -learners 4000 -set working` as processes by parts of the learners on all the cores (R09), then `sim merge` and the SHA-256 of the file.
- **The first file of bands** — `just gate record` with the reason T26, on the same commits v0.1.0: the numbers are the same to the bit as T24 and T25 recorded, the reason changes.

— *Why:*
- A path that the judgement reads from its own runs would go after the first choice from the new base, and in the confirmation from the held-out learners, while section 3.6 judges "at a pinned path".
- The file is in the code, not only in the log: the judgement cannot diverge from the record, and a test holds the record equal to the file.
- The files of a run are about 2 GB per rule. The seeds and the same bits give them anew, and the SHA-256 in the file of paths says whether it is that run.

— *Rejected:*
- the path from the judgement's runs, as it was — it would drift;
- keeping the files of a run as an artifact or a release file — space, retention and a signature check for the sake of what is reproduced;
- a run on a GitHub runner — 4 cores, about 2.5 hours, an artifact of gigabytes.

— *Consequence:* the first candidate is judged by the pinned paths. A new pinning run — only by an amendment to R04, like a new held-out set. — RUN.md, T26; `internal/criterion/`, `main.go`, `justfile`, `gate.json`; [architecture.md](architecture.md), sections 3.6–3.8.

— *Amendment 2026-11-09:* the paths are pinned anew, before the first candidate: the graph changed — `sim` pinned `olympiad-math` 0.1.1, and the populations `A-` went on the published pair of maths. This is not a new choice but the same calibration on the graphs that the learners will get. The runs and the new table of paths — the 2026-11-09 amendment to R04, "a new pin"; why — R31.

— *Amendment 2026-11-10:* a decisive run is 1,000 learners a cell (R33), so `sim paths` refuses a cell with fewer than 1,000 learners, not 4,000. The recipe `just decisive` runs `-learners 1000` and creates the folder of the file if there is none. The files of a run are about 0.5 GB per rule. Between judgements they lie in `runs/` of the working copy: `git` does not see it, and Rebuild Container does not erase it, as it erased the files of T30 in `/tmp`. The paths are pinned anew — the 2026-11-10 amendment to R04, "runs no longer than an hour". Reproducibility is checked by a full repeat of the ceiling: on 1,000 learners this is cheap.

— *Amendment 2026-11-10, a candidate's run (T31; R34):* the recipe `just decisive engine|ceiling <file> [set]` takes a set of learners — `working` or `held-out-<criterion>-<n>` — and a candidate: `CANDIDATE_LEARNER_MODEL`, `CANDIDATE_ENGINE` or both. Without a candidate it builds the simulator by `go.mod`, as before, so the file of the base is the one whose SHA-256 `paths.json` names.

— *Amendment 2026-11-11 (T41; the 2026-11-11 amendment to R33):* a decisive run is 1,500 learners a cell, so `sim paths` refuses a cell with fewer than 1,500 learners, and a judgement calls such a run a rough look. The recipe `just decisive` runs `-learners 1500`. The files of a run are about 0.6–0.8 GB per rule. The paths are pinned anew — the 2026-11-11 amendment to R04. Reproducibility is checked by a full repeat of the ceiling. The former files on 1,000 learners lie in `runs/` as `base-1000.jsonl` and `ceiling-1000.jsonl`, and the judgement files of `engine-1` and `model-1` name the base so.

— *Amendment 2026-11-11 (T43):* the paths are pinned anew after a change of the packages and of a library, not after a candidate: the pin of `modules` v0.2.5 changes the graph and the tasks, and the base since T42 is on `learner-model` v0.1.5. `paths.json` names `learner-model` v0.1.5 and `engine` v0.1.6, the table — the 2026-11-11 amendment to R04, "packages with `time` and leading questions". Reproducibility is checked by a full repeat of the ceiling. Why — R42.

— *Amendment 2026-11-11 (T45):* the paths are pinned anew after a change of the library, though the numbers did not shift: a judgement holds the base to the libraries of the paths (R40), and the base is on `learner-model` v0.1.6. `paths.json` names `learner-model` v0.1.6 and `engine` v0.1.6 and the new SHA-256 of the runs, and the values are the same as T43's, so the table of the 2026-11-11 amendment to R04, "packages with `time` and leading questions", is still the table of paths. Why — R43.

— *Amendment 2026-11-12 (T46; the 2026-11-12 amendment to R33):* a decisive run is 2,000 learners a cell, so `sim paths` refuses a cell with fewer than 2,000 learners, and a judgement calls such a run a rough look. The recipe `just decisive` runs `-learners 2000`. The files of a run are about 0.8–1.0 GB per rule. The paths are pinned anew — the 2026-11-12 amendment to R04. Reproducibility is checked by a full repeat of the ceiling. The former files on 1,500 learners lie in `runs/` as `base-1500.jsonl` and `ceiling-1500.jsonl`; no judgement refers to them.

— *Amendment 2026-11-13 (T47; R46):* the table of the paths is a section of [measurements.md](measurements.md) — measurements 3.1 for the base on `theta` —, and the amendment to R04 that pins new paths cites it. The test holds that section's table equal to `Table()` of the embedded file, not the log's, and `sim paths -table` prints its rows. Why: the results of a measurement go to measurements (RULES.md, How we work; R143 spec), and every new pin added a table of about 7 KB to R04, which stood at 78.6 KB under its limit of 90 KB (measurements 2.1). The tables of the earlier amendments to R04 stay where they are: the log is not rewritten. `paths.json` names `learner-model` `v0.1.7-0.20261113040121-a063835087c7` and the SHA-256 of the new runs. Rejected: raising the limit of an entry again — R04 would grow by a table at every pin.

— *Amendment 2026-11-13 (T48; R46):* the paths are pinned anew on the release of `theta`, though the numbers did not shift, as in T45; the table is a subsection of its own, measurements 3.2, as the amendment of T47 above has it. What and why — the 2026-11-13 amendment to R04, "the release of `theta`".

**R29. The release of `sim` is the shared `release.yml` of `.github` v0.1.21 on the line `0.1`, starting with `v0.1.0`, when a commit of `main` with green CI changes what the gate is built from. The only release file is the archive `sim-<version>.tar.gz` with what builds the simulator and runs the gate, signed without a key. Immutable releases and a ruleset of the tags `v*` are turned on before the first release.** Decided by the author 2026-11-09 in the planning of T27 (RUN.md, T27): the shared release with an archive (the 2026-11-09 amendment to R05), the GitHub settings — by the executor through `gh api`, the merge that releases the tag — by the author. The rest — the executor.
- **`paths`** — what the release gives:
  - the Go code without tests and `testdata/`;
  - the embedded packages `internal/published/v*` and the paths `internal/criterion/paths.json`;
  - `go.mod`, `go.sum`, `LICENSE`, `.gitattributes`;
  - the gate file `gate.json`, `justfile` and `.devcontainer/Dockerfile`, on which the gate runs;
  - `release.yml` itself.

  Tests, documents and the other workflows the release does not give.
- **The archive** — the recipe `release-files <version>`: the folder `sim-<version>` with the same paths, except `release.yml`. The tar is written the same way on any build: the names in order, without times and owners, mode 0644, `gzip -9n`. The list is taken from the tree, `find`: the shared release runs the recipe in the tree of `git archive`, without git.
- **`ci-release`** holds the archive, as R15 engine does:
  - two builds in a row — the same bytes;
  - the archive has exactly what git tracks at these paths, and in the same form;
  - the unpacked archive builds.

  It is called by `ci-lint`, so the archive is held by "Go / lint".
- **GitHub settings** (2026-11-09, `gh api`):
  - immutable releases are turned on;
  - the ruleset `v-tags` (id 26826001) forbids update, deletion and non_fast_forward of the tags `refs/tags/v*`, without bypasses, — the same as `engine`'s.

— *Why:*
- A shared release is one release code for all the repositories with code, and Dependabot raises it together with the other workflows of `.github`.
- The archive is exactly what the gate's numbers are recorded on, except the libraries that `gate.json` names. The signature binds it to the tag (R77 spec).
- A change of documents or tests gives the callers no new tag.

— *Rejected:*
- a workflow of its own without files, as `spec` has — release code of its own that only `sim` keeps;
- an archive of the Go code only, as `engine` has — the gate runs on `gate.json`, `justfile` and the stage `toolchain` too, and without them the archive does not hold that the numbers are recorded on this tag.

— *Consequence:* the first tag `v0.1.0` comes out with the merge of the pull request of T27. The callers take its tag and SHA — README, "Taking a release" — by their own tasks: Д5 RUN.md `modules`, the gate's job of `learner-model` and `engine`. — RUN.md, T27; `.github/workflows/release.yml`, `justfile`, `.gitignore`, README; [architecture.md](architecture.md), section 3.9.

## Epic E3

**R30. Epic E3, the first candidates, runs up to the first choices. In it: the small path in the judgement, a new pin of the packages with a new calibration of the paths before the first candidate, the decisive run on a candidate's commit, the judgement of the model's first candidate and of the engine's first candidate, the move to the releases of the chosen ones and the outcome.** These are the author's answers 2026-11-09 when planning E3 (RUN.md, table 1). The order of the tasks T29–T35 is a proposal of the planning.
- **What was found in planning** (read 2026-11-09):
  - `modules` v0.1.1 published `olympiad-math` 0.1.1 on 2026-11-08 — with new skills and hundreds of tasks. The index of the newest release, v0.1.4, names both versions. `sim` pins `modules` v0.1.0, so the base, the paths and the bands are taken on `olympiad-math` 0.1.0;
  - `just decisive` builds the simulator by `go.mod` and runs only the working learners, so there is nowhere to substitute a candidate. `sim run -set` already takes the held-out learners;
  - the libraries connect the gate — T22 RUN.md `learner-model` and T23 RUN.md `engine`, — but have not proposed candidates yet. Their basis is Д5–Д7 RUN.md `learner-model`, Д3 and Д4 RUN.md `engine`;
  - `spec` has come out up to v0.9.16, but the libraries do not require the Go module `spec`, so the gate needs no move.
- **Decided by the author:**
  - the small path — option (a) of В7: the 2026-11-09 amendment to R04, the code — T29;
  - a new pin — before the candidates: the index of `modules` v0.1.4 with `olympiad-math` 0.1.1, the decisive run of the base and of the ceiling anew, new paths and bands, the next tag of `sim` (T30);
  - the epic runs up to the first choices: the judgement of the model's first candidate (T32) and of the engine's (T33), the move to their releases (T34) and the outcome (T35).
- **A proposal of the planning:**
  - the order — the small path (T29), the pin (T30), the recipe of a candidate's run (T31, after T30: the run of the base by the recipe matches to the byte the one whose SHA-256 the new file of paths names), the judgements (T32, T33), the releases (T34), the outcome (T35);
  - if the candidates come as a pair, one task judges them, T32, by the criterion of the pair. The seed of its held-out learners is written down by an amendment to R04 before the run.
  - if one library has not proposed a candidate when the other's judgement is done, the author defers its judgement, and the move to the releases and the outcome go without it.

— *Why:*
- **The pin before the candidates.** The paths are calibration, not a choice, and are read before the first candidate (R04). A candidate chosen on a graph that the learners will not see is chosen on the wrong thing. Pinning after the first choices would mean calibrating anew when judgements had already been made on the paths.
- **Up to the first choices.** The simulator was built for them: the first version of each library is the base, and what improves it passes through a judgement (R103 spec). The judgements wait for candidates, but this is the work of `sim`, and in the epic it is seen together with what prepares it.

— *Rejected:*
- judging the first candidates on `olympiad-math` 0.1.0, and pinning 0.1.1 after — see the "Why";
- pinning 0.1.1 and leaving the paths of the olympiad cells from 0.1.0 — those would be the paths of another graph;
- ending the epic with readiness — the small path, the pin and the recipe, — and leading the judgement of each candidate by a task outside an epic. Then the first choices would enter no epic, and nobody would see their order in advance.

— *Consequence:*
- E2 is moved to the archive [docs/plan/E2.md](plan/E2.md);
- table 1 of RUN.md is E3 in the tasks T29–T35;
- the row E3 left table 2, and В7 table 3;
- В9 «в реестр spec» tells `spec` that the epic is laid out and the candidates are judged after the tag with the new pin.

Where it is written: RUN.md, tables 1–3 and the bodies of T29–T35; [architecture.md](architecture.md), sections 2.4 and 3.

## A new pin of the packages

**R31. The simulator is pinned to the release of `modules` v0.1.4 — `olympiad-math` 0.1.1, `grade-3-math` 0.1.0 and `english-essays` 0.1.0. The populations `A-` run on the pair of published maths. The gate holds 100 learners a cell, and on the repeat of the run takes the choice anew at every step only for the first 10 learners of a cell.** These are the author's answers 2026-11-09 when planning T30 and E3 (R30).
- **What was found in planning** (2026-11-09, measurements on a copy of the repository; the sources — [architecture.md](architecture.md), section 4):
  - `olympiad-math` 0.1.1 has 19 skills instead of 3 and 772 tasks instead of 27, 19 of them with instances of a template. There are now four kinds of answer, not only `number`. There are still no slots `predicate`;
  - two olympiad skills rest on skills that school maths teaches: `fractions-parts` on `core.fractions-of-a-whole`, `number-divisibility` on `core.division-with-remainder`. So the populations `A-` got a published graph;
  - the simulator reads the package without edits of the code; `G0` and `A-` on the pair run without violations of the invariants, and `M9` is read;
  - a learner of the olympiad graph runs for about 0.65 s instead of about 0.007 s. In the profile:
    - `engine.Choose` on a pool of 772 tasks — about 60 % of the CPU;
    - of that, about 28 % is the choice that `invariant.Check` takes anew at every step;
    - garbage collection — about 23 %;
  - the gate on 300 learners ran for 8.5 minutes and 135 CPU minutes on 20 cores — against 45 seconds and 11 minutes before. On the 4 cores of the runner this is about 35–40 minutes, at a budget of `timeout-minutes: 15`;
  - `sim gate changed` compares with the index that the CI of `modules` passes, so the pin does not touch it. But `A-linked` runs on both maths, and the command gave only the changed modules.
- **Decided by the author:**
  - the populations `A-` — on the published pair, in this same calibration;
  - the gate is cut by an edit of `sim`, and the budget stays 15 minutes: 100 learners a cell instead of 300, and on the repeat of the run the choice is not taken anew at every step. The executor, after review, kept the choice anew for the first 10 learners of each cell — about 3 % of the work: the repeat does not see a rule that changes its input in place the same way in both runs;
  - the index of `modules` v0.1.4 — already when planning E3 (R30).
- **A proposal of the execution:**
  - the archive of `olympiad-math` 0.1.1 lies in the folder `v0.1.1/`, where the index's record names it, and the unused index and archive of `olympiad-math` 0.1.0 are removed from the embedded files;
  - `gate.Runs` gives, together with a changed module, the unchanged one that the graph of a population of the gate holds together with it;
  - the tests that run the olympiad graph run on fewer learners, so that `ci-test` fits into the budget of the builds of `same-bits`;
  - the test of the ceiling lets the same task in again after a pause of 24 hours, as the rule does (R82 spec): on 27 tasks this did not happen.

— *Why:*
- **`A-` in this calibration.** Otherwise the paths of the group "Anchor" and the bands of `A-linked` would have to be pinned once more — after the first candidates.
- **100 learners.** [architecture.md](architecture.md), section 3.9, says to cut a gate that does not fit by an edit of `sim`, not by time. The bands become wider by about √3, about 1.7 times. Code that shifted nothing still reads the same numbers to the bit, and the candidates are judged by the decisive run. The gate loses only its sensitivity to a small unintended shift.
- **The choice anew on the repeat — only for the first 10.** `I1` is held to the bit by the repeat (`Same`), and the engine's refusal of its own input is caught by the run itself. The choice anew at every step for all learners would cost a third of the work. But the repeat does not see a rule that changes its input in place the same way in both runs, while the choice anew does see it, and on the first 10 learners of each cell it costs about 3 %.

— *Rejected:*
- raise the budget to 45 minutes and keep 300 learners — then every pull request of three repositories would wait for the gate about 40 minutes;
- wait until `engine` speeds up the choice on a large pool — the first candidates would wait for someone else's task. The question went as row В10 of table 3 of RUN.md;
- keep `A-` on synthetics — see the "Why".

— *Consequence:* the runs and the new table of the paths — the 2026-11-09 amendment to R04, "a new pin"; amendments to R14, R17, R19, R21 and R26–R28; a new tag of `sim`, which the gates of the libraries take (В9). Where it is written: RUN.md, T30; [architecture.md](architecture.md), sections 3, 3.2, 3.3, 3.7–3.9 and 4; README.

— *Amendment 2026-11-10:* the process of the gate's run, like that of the decisive one, runs on one processor and with `GOGC` 400 (R33): the bits are the same, and `just gate exact` matches to the bit, while the gate's time on the runner falls. The gate's size is as before — 100 learners a cell and the choice anew for the first 10; giving it back 300 is still for later, by an amendment to this entry.

— *Amendment 2026-11-10, a measurement (T35):* that the gate's time on the runner falls was confirmed only on `ubuntu-24.04-arm`: 7 min 3 s and 6 min 49 s on pull requests #40 and #41, 4 min 36 s and 4 min 42 s on #42 and #43. On `ubuntu-24.04` — 11 min 52 s and 6 min 54 s before T36, 8 min 14 s and 8 min 23 s after: the runner's time wanders more than the edit gave, and the measurement shows no gain there ([architecture.md](architecture.md), section 3.9).

— *Amendment 2026-11-11, after the speed-up of `engine` (T41; R38):* the gate's base is `engine` v0.1.6, and the gate has 150 learners a cell, not 100. The choice anew on the repeat is still for the first 10 learners of a cell, and the budget is still `timeout-minutes: 15`. The author's answers 2026-11-10 and 2026-11-11 when planning and executing T41 (RUN.md, T41).
- **What was found** (2026-11-10 and 2026-11-11):
  - on `engine` v0.1.6 the gate on 100 learners reads all 69 numbers the same to the bit as on v0.1.4, and the invariants are intact. The bands on 100 were not written anew;
  - in the development container, 100 learners: 1 min 44 s by the clock and 25.9 CPU minutes on v0.1.4, 1 min 5 s and 16.5 min on v0.1.6 — 1.57 times less;
  - "Gate" of `sim` on 100 learners, pull request #49: `ubuntu-24.04` — 5 min 39 s with a cold Go cache and 4 min 27 s with a warm one, `ubuntu-24.04-arm` — 3 min 13 s and 3 min 7 s. The constant part of the job is about 35 seconds, the run with the build up to 5 min on 100 learners;
  - the callers keep no Go cache, and the gate runs longer at theirs. "Gate" of `engine` on the code of v0.1.6 and 100 learners ran for 5 min 3 s and 7 min 44 s, the gate's step up to 7 min 18 s, at `timeout-minutes: 25` (В27 RUN.md `engine`). "Gate" of `learner-model` on the previous tag of `sim` and the slow choice ran for 8–12 min at `timeout-minutes: 15`; on v0.1.6 this is about 5–7.6 min on 100 learners;
  - recalculated from the slowest measurement to 150 learners: at `sim` — about 8 min, at `engine` — about 11.5, at `learner-model` — about 11.2. To 200: at `sim` — about 10.6, but at `learner-model` — up to 15, right at the timeout. To 250 at `sim` — about 13.1;
  - the gate on 150 learners in the development container — 1 min 28 s and 23.6 CPU minutes. All 69 numbers are read in the same places as on 100. Each new value is less than one previous half-width from the previous one, the most of all `M1` for `G0` of school maths, by 0.94. The half-width is on average 0.85 of the previous one, about 1/√1.5.
  - "Gate" on 150 learners, pull request #49, the second commit: `ubuntu-24.04` — 8 min 21 s with a cold Go cache and 6 min 5 s with a warm one, `ubuntu-24.04-arm` — 4 min 16 s. The recalculation gave about 8 min on a cold one — the measurement matched it.
- **Decided by the author:**
  - the gate's size is by the slowest measurement: the largest, in steps of 50, from 100 to 300, at which the slowest measurement, recalculated to this size, is no more than 12 minutes — 3 minutes of margin before the timeout;
  - 150, not 200: on 200 the gate of `learner-model` without a Go cache would come close to its 15 minutes.

— *Why:* the rule of R38 — the gate gets as many learners as fit into the budget on the runner. The runner's time wanders: on the same code "Gate" ran both 6 min 54 s and 11 min 52 s. A gate that fails by timeout holds the merge in three repositories, so the slowest measurement is taken as the measure, with a margin. It must fit at every caller, not only in its own CI with a cache. 150 learners narrow the bands by about √1.5 ≈ 1.22 times.

— *Rejected:*
- 200 with `timeout-minutes: 25` at `learner-model`, as at `engine` — every pull request of `learner-model` would wait for the gate longer, and this is an edit of someone else's workflow;
- 200 as it is — "Gate" of `learner-model` would fail by timeout on a slow runner;
- 300 — about 16 min even in its own CI;
- the size by a typical measurement, not by the slowest — a slow runner would bring the gate down.

— *Consequence:* `gate.Learners` is 150; the bands are recorded anew, `just gate record` with a reason; the thresholds of the CI of `modules` are taken from the new record (R27). 300 learners of the gate wait for a speed-up of `learner-model`: a third of the choice is now `Readout.Task` (row В32 RUN.md `engine`, «в реестр spec»). Where it is written: RUN.md, T41; [architecture.md](architecture.md), sections 3.7–3.9; README; `internal/gate/gate.go`, `gate.json`.

— *Amendment 2026-11-12, after the speed-up of `learner-model` (T46; R43, R44):* the gate has 250 learners a cell, not 150. The choice anew on the repeat is still for the first 10 learners of a cell, the budget is still `timeout-minutes: 15`, and the rule of the size is the same: the largest, in steps of 50, from 100 to 300, at which the slowest measurement, recalculated to this size, is no more than 12 minutes — with one refinement of the author: a size that passes the rule right at the edge, with less than a minute to 12, is not taken, and the gate takes the step below. By the rule alone it would come out as 300. The author's answers 2026-11-12 when planning T46 (RUN.md, T46).
- **What was found** (2026-11-12):
  - "Gate" of `sim` on the tag v0.1.9 — `learner-model` v0.1.6, `modules` v0.2.5, 150 learners —, both times with a cold Go cache: pull request #52 — 6 min 17 s on `ubuntu-24.04` and 3 min 36 s on `ubuntu-24.04-arm`, the push to main — 5 min 0 s and 3 min 32 s. In T41 on 150 learners with a cold cache it was 8 min 21 s, about 0.75 of this;
  - the callers have not moved to v0.1.9 yet. "Gate" of `learner-model` on `sim` v0.1.8 runs its already fast code — up to 6 min 12 s. "Gate" of `engine` on `sim` v0.1.6 runs with `learner-model` v0.1.4 — up to 8 min 21 s, about 6 min 16 s recalculated by the ratio 0.75;
  - recalculating the slowest measurement, about 6 min 17 s, with the constant part of the job about 35 s, as in T41: 200 learners — about 8 min 11 s, 250 — about 10 min 5 s, 300 — about 11 min 59 s;
  - "Gate" on 250 learners, pull request #53, with a cold Go cache: `ubuntu-24.04` — 7 min 41 s, the gate's step 7 min 3 s; `ubuntu-24.04-arm` — 5 min 28 s; the whole CI is green. A repeat of the job on `ubuntu-24.04` with a warm cache — 5 min 34 s, the gate's step 5 min 6 s. On the second commit of the pull request, also with a warm cache, — 9 min 38 s, the gate's step 9 min 9 s: the runner happened to be slow. The callers have no Go cache, and on such a runner their gate would come close to 12 minutes or a little beyond — about 11 min 45 s – 13 min 17 s, depending on how the cold cache is recalculated —, at a timeout of 15 at `learner-model` and 25 at `engine`. The author decided to keep 250: the slowest measurement is in fact 9 min 38 s, and the callers will show their time on the new tag;
  - in the development container the gate on 250 learners — 1 min 52 s and 29.6 CPU minutes, at a load average of about 40–60 from neighbouring containers.
- **Decided by the author:** 250 learners.

— *Why:* by the rule 300 learners pass right at the edge — about 11 min 59 s at a threshold of 12 —, while the runner's time on the same code wandered up to 1.7 times (the 2026-11-10 amendment "a measurement"), and the callers on the new tag have not been measured yet. 250 leave about two minutes beyond the rule's three minutes of margin. The bands on 250 are narrower than on 150 by √(5/3) ≈ 1.29 times.

— *Rejected:*
- 300 by the letter of the rule — there is no margin for the runner's noise;
- 300 with a rollback to 250 if the measurement of the pull request goes over 12 minutes — one lucky measurement does not protect the callers from a timeout.

— *Consequence:* `gate.Learners` is 250; the bands are recorded anew, `just gate record` with a reason (the 2026-11-12 amendment to R26); the thresholds of the CI of `modules` are taken from the new record (R27). Д7 is closed together with the 2026-11-12 amendment to R33. Where it is written: RUN.md, T46; [architecture.md](architecture.md), sections 3.7–3.9; README; `internal/gate/gate.go`, `gate.json`.

## The tag `spec` v0.12.0

**R32. `sim` stands on `spec` v0.12.0, the latest release of `spec`. With it the packages are pinned to `modules` v0.2.0 — all three modules in format 0.4 —, and the libraries to `learner-model` v0.1.4 and `engine` v0.1.4. The step of gopls in the image is the template's of `spec`, and R11 is superseded.** The author's instruction 2026-11-10 — «перейди на spec последней версии» — and their answers to three questions of the move.
- **What the tag decided for `sim`** (`git -C spec diff v0.9.9 v0.12.0`, read 2026-11-09 and 2026-11-10; the sources — [architecture.md](architecture.md), section 4):
  - format 0.4 came out with v0.10.0 (R102, R108, R111 spec). While the format is v0, the platform reads only its current version ([format text](../spec/format/README.md), section 2), so the simulator needs the window 0.4 and packages in 0.4. All the additions of 0.4 are optional, and a module of 0.3 moved to 0.4 gives the same verdicts and numbers. In the types `formatv0` a lesson's explanation is now `ExplanationPart`;
  - Д137 of table 4 of RUN.md `spec` is the answer to row В24 of `engine`: `sim` takes `engine` and `learner-model` v0.1.3 in the new pin, before the decisive run of the base. They came out through the gate with the same numbers;
  - Д123 and Д136 are the answers to the rows of `sim` В6, В8 and В9, by their recommendations: rows 4.5 and 4.8 of [COORDINATION.md](../spec/COORDINATION.md) are closed, and the first candidates are row 4.9;
  - the 2026-11-09 amendment to R103 spec: `sim` judges the candidates of one library in a round — a decisive run on each, one judgement with all of them, a confirmation only of the chosen one (row В31 of `learner-model`). This concerns T32 and T33;
  - the template's kind `go` gives the step of gopls with a lock file (R102 spec), of the same version 0.24.0 as `sim`'s;
  - the rest — the data region, the legal questions and the learners' data (R109, R112–R123 spec) — does not concern `sim`: it has no data of live learners.
- **What is changed here:**
  - `spec/` — to v0.12.0, commit 5b4640a; `gh release verify` confirmed the release;
  - `go.mod` — `spec` v0.12.0, `engine` and `learner-model` v0.1.4; `release-verify` accepted the archives of both libraries;
  - the format window — 0.4; the synthetic packages and the tests — in 0.4; a lesson's explanation — `ExplanationPart`;
  - the pin — `modules` v0.2.0: `release-verify` accepted the index and the three archives, all the modules of version 0.2.0;
  - `gate.json` — the libraries v0.1.4 by the commits of their tags, `modules` v0.2.0, `spec` v0.12.0. `just gate exact` matched to the bit: the bands are as before;
  - the step of gopls and `.devcontainer/gopls/go.mod` — from the template; `go.sum` is the same;
  - the tag — in CLAUDE.md, README.md and the opening of RUN.md.
- **Decided by the author:**
  - the paths are written anew in T36, by a fast run, not here: T36 changes the design of the run anyway, and no candidates are judged before it;
  - the libraries — v0.1.4, the latest releases: the code is the same as in v0.1.3, with only documents and the move to `spec` v0.12.0 on top;
  - the step of gopls — the template's.
- **Found in the move:** the first 100 learners of each of the 65 cells on the new simulator give the same bytes as the run of the base of T30 — 6,500 lines. Neither the fix of the libraries v0.1.3 nor format 0.4 shifted the base's numbers, so the paths of T30 stay in force until T36. But two parts of 0.4 the simulator will not read when a module uses them: the check kind `time` and the leading question `leading_question`. The published packages do not have them; this is an open question — [architecture.md](architecture.md), subsection 2.5, row В11 of table 3 of RUN.md.

— *Why now:* the author's instruction. v0.12.0 is the latest release of `spec` on 2026-11-10, and `modules`, `engine` and `learner-model` already stand on it. Without it the simulator will not read packages in 0.4, and the gate of the CI of `modules` on `v0.2.0` will not run.

— *Rejected:*
- stay on v0.9.9 — see the "Why now";
- the libraries v0.1.3 — exactly what Д137 named, but v0.1.4 is the same code on the same tag of `spec` as `sim`'s;
- the decisive run now — about 8 hours, and after T36 the paths are written anew anyway;
- a step of gopls of its own — it would diverge from the template, and gopls would have to be raised by hand.

— *Consequence:* table 3 of RUN.md — В6, В8 and В9 closed, В10 stays; R11 is superseded; the 2026-11-10 amendments to R04, R14 and R26. The merge releases the next tag of `sim`, and the gates of the libraries take it. Where it is written: RUN.md, T37; [architecture.md](architecture.md), the status and sections 3, 3.3 and 4; CLAUDE.md, README.md.

## Runs no longer than an hour

**R33. The decisive run and the confirmation go on 1,000 learners a cell instead of 4,000. The choice anew at every step for `I1` is taken only by the first 10 learners of a cell — in every run, as in the gate. The process of a run goes with `GOGC` 400, and the process of a part of a run also on one processor. The files of the decisive runs are kept in `runs/` of the working copy, which `git` does not see, and a judgement runs only the candidate.** The author's decision 2026-11-09, in T30: a run longer than an hour is not acceptable — the author is prototyping, and the functionality changes fast. The author's answers 2026-11-10 when planning T36 (RUN.md, T36); the runtime environment is a proposal of the execution.
- **What was found in planning** (2026-11-10, measurements on a copy of the repository; the sources — [architecture.md](architecture.md), section 4):
  - 10 learners a cell in all 65 cells — 1,064 s of CPU:
    - half of it — the six populations `F-`, which have 180 days;
    - school maths, the essays and `S-predicate` together — less than 1 %;
    - the rest — the olympiad graph, 772 tasks;
  - in the profile of `G0` and `F-model`:
    - `engine.Choose` — 64 % of the CPU, of that the choice anew in `invariant.Check` — 31 % of the whole;
    - garbage collection — about a quarter;
  - on one sample, the bytes of the results are the same in all the variants:
    - the choice anew only for the first 10 — 121 → 70 s of CPU;
    - `GOGC=400` — 121 → 87 s;
    - all together with `GOMAXPROCS=1` — 43 s, 2.8 times less;
    - the ceiling with `GOGC=400` and `GOMAXPROCS=1` — 22 → 16 s;
  - the memory at `GOGC=400` — about 100 MB a process, 20 processes — about 2 GB;
  - after these edits `engine.Choose` is 73 % of the CPU, and the simulator's own code about 17 %, mostly the P* of the pool at every step for the measures. Beyond that only `engine` speeds up the run (row В10 of table 3 of RUN.md). The run of the base on 4,000 learners would take about 2.4 hours;
  - the files of the decisive runs of T30 lay in `/tmp` and were lost after Rebuild Container.
- **Decided by the author:**
  - 1,000 learners a cell for the decisive run and the confirmation;
  - the choice anew for `I1` — for the first 10 learners of a cell, `bench.Rechosen`;
  - the run of the base is kept in `runs/`.
- **A proposal of the execution:**
  - `sim run` and `sim gate run` set themselves `GOGC` 400, and the run by parts also `GOMAXPROCS` 1, in the code, not in the recipe, so that every call gets them. R09 keeps goroutines out of the code anyway, and the parts of a run are processes, one per processor. A whole run is one process on the machine, and the processors are left to its garbage collection;
  - the populations of the decisive run — all of them: "no worse" is read on every population with a verdict ([architecture.md](architecture.md), section 3.6), and without any of them the judgement would not read its constraint;
  - `invariant.CheckRepeated` is renamed `CheckWithoutI1`: the decisive run has no repeat;
  - the recipe `just decisive` creates the file's folder if there is none.
- **The run** (the 2026-11-10 amendment to R04, "runs no longer than an hour"): the base on 1,000 learners — about 37 CPU minutes per process, the ceiling — about 14. On a free machine each run fits into an hour; on a machine that other containers loaded, the base took 1 h 42 min by the clock. The paths are close to the paths of T30, and the measures with a small path are the same seven.

— *Why:*
- **1,000 learners.** By its own means the simulator speeds the run up 2.8 times, and about seven were needed. The rest only a smaller number of learners gives, until `engine` speeds up the choice. 1,000 leave a margin for a busy machine and for the growth of the pools. The resolution of "no worse" — 4.3 standard errors of the paired difference — is wider by √4 = 2 times, and the author's tolerance, where it is smaller, gives way to the resolution. The judgement still tells apart what is stronger than the noise, and a rough look — fewer than 1,000 — still decides nothing.
- **The choice anew for the first 10.** The repeat, `Same`, holds `I1` to the bit at the gate, and the first 10 learners of each cell catch a rule that changes its input in place (R31). In the decisive run these are 650 learners and about 130 thousand choices. The choice anew for all cost about 40 % of the CPU. A rule that does not choose the same thing on the same input — for example, by the order of ranging over a map — the decisive run without a repeat will not see for the other learners. But a candidate, before it is judged, passes the gate in the CI of its library, and the gate repeats each of its 100 learners of each cell and holds `I1` to the bit (R31). And its decisive run, repeated on the same seeds, must give the same bytes.
- **`GOMAXPROCS` 1 and `GOGC` 400.** Neither changes a bit, and garbage collection cost about a quarter of the CPU. There are as many processes as processors, and the extra threads of each process only take CPU away from the others.
- **`runs/`.** A folder of the working copy survives Rebuild Container, and `/tmp` does not. `git` does not see it, and the files still do not lie in the repository (R28). A judgement runs the candidate, and the base's file is checked by SHA-256 against what `paths.json` names; how the recipe of a candidate's run does this, T31 decides.

— *Rejected:*
- keep 4,000 and wait until `engine` speeds up the choice — the judgements would wait for someone else's task, and the run would take about 2.4 hours even after the edits of the simulator;
- 1,500 learners — about 54 minutes on a free machine, right at the hour, and on a busy one more;
- 1,000 only on the olympiad cells and 4,000 on the others — the size of the run would depend on the graph, and with the growth of the packages the rule would have to be revisited;
- the choice anew for all learners — then about 750 learners would fit into an hour, and the resolution would be 2.3 times wider;
- cut the populations down to the score groups — see the "A proposal of the execution";
- run the base anew in every judgement — two runs instead of one.

— *Consequence:* the 2026-11-10 amendments to R04, R17, R28 and R31; a new `paths.json`. Returning to 4,000 learners — by an amendment to R04, when `engine` speeds up the choice (В10). Where it is written: RUN.md, T36; [architecture.md](architecture.md), sections 3, 3.6–3.8; README; `main.go`, `gate.go`, `internal/bench/`, `internal/invariant/`, `internal/criterion/`, `justfile`, `.gitignore`.

— *Amendment 2026-11-11, after the speed-up of `engine` (T41; R38):* the decisive run and the confirmation go on 1,500 learners a cell, not on 1,000. The choice anew for `I1` is still for the first 10 learners of a cell; `GOGC` 400 and one processor a part are as they were. The author's answers 2026-11-10 when planning T41 (RUN.md, T41).
- **What was found** (2026-11-10, the development machine, 20 threads; the head of pull request #57 of `engine`, `74afdf8`, by its files the same as the tag v0.1.6):
  - the sample — 100 learners a cell in all 65 cells, by parts on 20 processes: the base on `engine` v0.1.4 — 256 s, on v0.1.6 — 166 s, 1.54 times less; the ceiling — 83 s, it does not call `engine`;
  - recalculated to 1,000 learners the base is about 28 min instead of 43, to 1,500 about 41 min, to 1,600 about 44;
  - the run of the base on 1,500 learners — 38 min by the clock on an almost free machine, the ceiling — about 19 CPU minutes per process (the 2026-11-11 amendment to R04);
  - the measures with a path no greater than the author's tolerance — the same seven as on 1,000.
- **Decided by the author:** 1,500 learners; the choice anew for `I1` — for the first 10.

— *Why:* the rule of R38 — the decisive run gets as many learners as fit into an hour with a margin, and "with a margin" is no longer than 45 minutes on a free machine, a quarter of an hour for a busy machine and the growth of the pools. 1,500 is about 41 minutes. The resolution of "no worse" is narrower by √1.5 ≈ 1.22 times. R33's rejection of 1,500 stood on 54 minutes at `engine` v0.1.4 — a new fact removed it.

— *Rejected:*
- 1,600 — about 44 minutes, right at 45: no margin is left for the growth of the pools;
- 1,250 — about 35 minutes: more margin than R38 asks, and the resolution is narrower only by 1.12 times;
- the choice anew for `I1` for all learners — then about 1,000 would fit into 45 minutes, and `I1` is held to the bit by the gate's repeat anyway (R31).

— *Consequence:* `criterion.Decisive` and the recipe `just decisive` — 1,500; the 2026-11-11 amendment to R04 — a new `paths.json`; the judgements `engine-1` and `model-1` stay on their runs of 1,000, and the next judgements and confirmations go on 1,500. 4,000 still do not fit: a third of the choice is now `Readout.Task` of `learner-model` (row В32 RUN.md `engine`, «в реестр spec»), and beyond that `learner-model` speeds up the run. Where it is written: RUN.md, T41; [architecture.md](architecture.md), sections 3.6 and 3.8; README; `internal/criterion/criterion.go`, `justfile`.

— *Amendment 2026-11-12, after the speed-up of `learner-model` (T46; R43, R44):* the decisive run and the confirmation go on 2,000 learners a cell, not on 1,500. The choice anew for `I1` is still for the first 10 learners of a cell; `GOGC` 400 and one processor a part are as they were. The author's answers 2026-11-12 when planning and executing T46 (RUN.md, T46).
- **What was found** (2026-11-12, the development machine, 20 threads; `learner-model` v0.1.6, `engine` v0.1.6, the packages of `modules` v0.2.5):
  - the sample — 100 learners a cell in all 65 cells, by parts on 20 processes, on a free machine: the load average was 1.8 before the build of the simulator and 2.0 after it, at the start of the run. The base — 1 min 55 s against 166 s in T41, the ceiling — 1 min 20 s against 83 s: the base has the speed-up of `Readout.Task`, and the ceiling does not call `engine`;
  - recalculated, the base on 1,500 learners is about 28.8 minutes, on 1,750 about 33.5, on 2,000 about 38.3, on 2,250 about 43.1; the ceiling on 2,000 about 26.7. This matches the CPU measurement of T45 — about 25.9 minutes per process for the base on 1,500 (R43). In T41 such a sample promised 41 minutes on 1,500, and the run took 38: the recalculation overstates a little;
  - the runs on 2,000 learners went on a host that neighbouring containers shared: the base — 62 min 37 s by the clock and about 31.3 CPU minutes per process, the ceiling — 30 min 37 s and about 24.1 (the 2026-11-12 amendment to R04). The CPU per process of the ceiling is 4/3 more than on 1,500 in T45 (18.0), as expected, and that of the base only 1.21 times more (25.9 → 31.3): CPU time on an overloaded host wanders by itself, so the size is chosen by the sample, not by it;
  - the file of a run — about 1.0 GB for the base and 0.8 GB for the ceiling; `sim paths` reads both in 5 min 23 s.
- **Decided by the author:** the size is the largest on a grid in steps of 250, from 1,500 to 4,000, at which the base, recalculated from the sample on a free machine, is no longer than 41 minutes, as 1,500 was in T41. This is 2,000; the choice anew for `I1` — for the first 10.

— *Why:* the rule of R38 — the decisive run gets as many learners as fit into an hour with a margin, and "with a margin" is no longer than 45 minutes on a free machine (the 2026-11-11 amendment). 2,000 is about 38 minutes; 2,250 about 43, right at the edge, like the 1,600 at 44 minutes rejected in T41. The resolution of "no worse" is narrower by √(4/3) ≈ 1.15 times.

— *Rejected:*
- 2,250 — about 43 minutes: no margin is left for the growth of the pools;
- a step of 100 — 2,100 or 2,200 would give 3–5 % of resolution at the price of an uneven size and a smaller margin;
- no longer than 38 minutes — more margin than R38 asks: then 1,750, and the resolution is narrower only by 1.08 times;
- the size by the clock of a run on a busy host — it measures the neighbours, not the simulator.

— *Consequence:* `criterion.Decisive` and the recipe `just decisive` — 2,000; the 2026-11-12 amendment to R04 — a new `paths.json`; the judgements `engine-1` and `model-1` stay on their runs of 1,000, and the next judgements and confirmations, starting with `model-2`, go on 2,000. 4,000 still do not fit — about 77 minutes of the base; Д7 is closed together with the 2026-11-12 amendment to R31: the sizes are again as large as the time allows. Where it is written: RUN.md, T46; [architecture.md](architecture.md), sections 3.6 and 3.8; README; `internal/criterion/criterion.go`, `justfile`.

## A candidate's decisive run

**R34. A candidate's decisive run is run by the same recipe `just decisive` as the run of the base: with a set of learners and with the candidate from `CANDIDATE_LEARNER_MODEL`, `CANDIDATE_ENGINE` or both. A candidate is a commit by its full SHA, which the recipe fetches from GitHub, or a folder without uncommitted edits. The candidate's simulator is built through the gate's workspace: this piece of the recipes is one for both recipes, and the build before the run is held by `sim gate verify`. The base without a candidate is built by `go.mod`, as before.** The author's answers 2026-11-10 when planning T31 (RUN.md, T31); the rest is a proposal of the execution.
- **What was found in planning** (2026-11-10):
  - three candidates are open in `engine`, pull requests #46–48 — A, B and C (T31 `engine`). T33 will judge them in a round (the 2026-11-09 amendment to R103 spec);
  - GitHub gives a commit by its full SHA even when it is on no branch: the head of the open #48 of `engine` was fetched by `git fetch --depth 1`;
  - the gate's build is held by `gate.Verify`: the libraries from the workspace, the shared dependencies and the lines `go` of the gate file, the `go.mod` of `sim`, the pinned packages. It needs no parts of a run, but only `sim gate check` called it, together with the bands;
  - the header of the results file is written by `bench.Names`. It names a library from the workspace by the commit passed through `-commit`, and a library from `go.mod` by its version.
- **Decided by the author:**
  - the base without a candidate — by `go.mod`;
  - a candidate — a SHA or a clean folder, a pair — both variables.
- **A proposal of the execution:**
  - a hidden recipe `workspace dir prefix clean`. It takes each library from `<prefix>_<NAME>` or by the commit of the gate file, checks the tag, writes `go.work` and prints a line per library: the path, the folder, the commit. The gate calls it with `GATE` and without the check of cleanness, the decisive run with `CANDIDATE` and with it;
  - `sim gate verify -file <gate file> -go <folder>…` — `gate.Verify` without the bands: a candidate whose `go.mod` raised a shared dependency is not run, just as it would not pass the gate of its own CI;
  - the second library of a candidate is also from the workspace, by the commit of the gate file. So `gate.Verify` applies without edits, and the header names both libraries by commits;
  - a pair is both variables together with `CANDIDATE_PAIR=yes`. A variable left over from another run does not silently make a pair of a candidate, and a variable with a name that is not among the libraries of the gate file is a refusal;
  - a relative folder is from the caller's folder, as the run's file is;
  - the recipe checks the set of learners before the run, by a short run of one learner;
  - README — the section "Judging a candidate": from the commit of a pull request to the judgement, the confirmation and the recording of the bands.

— *Why:*
- **The base by `go.mod`.** Then its file is to the byte the one whose SHA-256 `paths.json` names (R28; the 2026-11-10 amendment to R04). Through the workspace the header would name the libraries by commits, the bytes would diverge, and the base with the paths would have to be written anew.
- **A SHA or a folder.** The head of a pull request lies on GitHub, and the SHA names it unambiguously. A folder is needed for a candidate that is not on GitHub yet. Without edits — because the header names a commit, and the file must be a run of exactly that commit.
- **One piece of the workspace.** The gate and a candidate's run build the simulator the same way; two copies would diverge.

— *Rejected:*
- the base through `go.work` — see the "Why";
- only a SHA — a candidate without a pull request could not be run;
- a folder with uncommitted edits — the header would name a commit, and other code would run;
- a copy of its own of the piece of `go.work` in the recipe `decisive`.

— *Consequence:* the 2026-11-10 amendments to R26 and R28; T32 and T33 run the candidates by this recipe. Where it is written: RUN.md, T31; [architecture.md](architecture.md), sections 3, 3.6, 3.8 and 3.9; README; `justfile`, `gate.go`.

## The order of the judgements and the judgement's report

**R35. Of the first two judgements — of the model and of the engine — the first is the one whose candidates are ready earlier, and the second goes on the base with what the first chose. The judgement's report lies in the repository, in `judgements/<criterion>-<n>/`. There is no row «В реестр spec» with a tag with a candidate.** The author's answers 2026-11-10 when planning T32 (RUN.md, T32).
- **What was found in planning** (2026-11-10):
  - the exploration of `learner-model` — its T30, pull request #42 — chose three candidates: `g10`, `halve035` and `forget-set`. All in the group "Forgets", without a target and without a new field, and their draft pull requests are T31 `learner-model`, not started yet;
  - the engine's candidates A, B and C are already ready, pull requests #46–48 of `engine`;
  - after the first choice the base is what was chosen ([architecture.md](architecture.md), section 3.6), so the order of the judgements decides on which base the second goes;
  - `sim judge` does not check that a candidate's second library is the same as the base's; from such an error the recipe's `CANDIDATE_PAIR` protects now (R34).
- **Decided by the author:**
  - **the order.** Now T33 goes first. If it chose the engine, T32 runs the base anew — the model v0.1.4 with the chosen engine, by the recipe on the new `go.mod` —, and the recipe builds the model's candidates with the engine at the commit of `gate.json` (R34). The path stays pinned (R28);
  - **the report.** The folder `judgements/<criterion>-<n>/`, where n is the number of the choice of this criterion, as for the held-out learners ([architecture.md](architecture.md), section 3.8): the judgement file, the report of `sim judge` as it is and the SHA-256 of each run. The outcome and its reasons are recorded by a decision of the log;
  - **`spec`.** The first candidates are led by row 4.9 of [COORDINATION.md](../spec/COORDINATION.md), and the order of releases by R103 spec: `spec` will see the tag of `sim` with the bands and the report on `main`, and the tag in "Gate" is raised by the library's session.

— *Why:*
- **Whoever is ready goes first.** The candidates of one library do not wait for someone else's, and the second judgement judges its candidates together with what is already chosen, as section 3.6 says.
- **The report in the repository.** The bands are recorded anew by a pull request with a reason — the report of the judgement and of the confirmation ([architecture.md](architecture.md), section 3.7) —, and review reads this report beside the bands. The runs are reproduced to the byte, but review should not have to take the report anew in order to read it.

— *Rejected:*
- judge as a pair — the engine's candidates would wait for the model's candidates, and the combinations 3 × 3 are nine runs;
- both judgements on the previous base — the chosen model and engine would never be judged together;
- the report only in the log — see the "Why";
- a row «В реестр spec» — what it would ask row 4.9 already leads.

— *Consequence:* T33 goes before T32; the reports are `judgements/engine-1/` and `judgements/model-1/`. Where it is written: RUN.md, T32 and T33; [architecture.md](architecture.md), section 3.6.

— *Amendment 2026-11-10, after the judgement `engine-1` (R36):* n in `judgements/<criterion>-<n>/` is the number of the round of the judgement of this criterion, not the number of the choice. A round may end in an exit without a confirmation, and then there is no choice, and the next round needs a folder of its own. The held-out set of a round is the first `held-out-<criterion>-<n>` not yet touched ([architecture.md](architecture.md), section 3.8): for `engine-2` it is `held-out-engine-1` again.

— *Amendment 2026-11-11 (T40; R40):* the judgement itself checks that a candidate's second library is the same as the base's, by the embedded `gate.json`; Д6 is closed. The recipe's protection — a pair requires `CANDIDATE_PAIR=yes` (R34) — stays.

## The judgement `engine-1`

**R36. The first judgement of the candidates of `engine` — A, B and C, pull requests #46–48 — ended in an exit: no candidate holds all the constraints of the engine's criterion, and the base stays as it is. There was no confirmation, and `held-out-engine-1` is still held out. `gate.json`, `go.mod` and the bands do not change.** The judgement of T33 (RUN.md), a round by the 2026-11-09 amendment to R103 spec; the report — `judgements/engine-1/` (R35).
- **The candidates** — all in the group `corridor`, without new constants and without a new field; the lines for the judgement file came from their descriptions:
  - A, `7ef8288c…`, `distance` 1 — a task in the corridor before an unseen one outside it;
  - B, `32aa835d…`, `distance` 2 — a seen solution after a pause does not change the order;
  - C, `97abff0d…`, `distance` 3 — the order of A and consolidation of an error only in the corridor.
- **The runs** — decisive, on 1,000 working learners a cell, by the recipe `just decisive` (R34), on a free machine:
  - the base — `runs/base.jsonl` of the run of T36, SHA-256 `b13a9264…4c92`, the one that `paths.json` names;
  - A — 48 min by the clock, SHA-256 `6ac032a5…e94d`;
  - B — 42 min, `21ad0cbd…f8e4`;
  - C — 41 min, `79353c9a…b8459`.

  So the hour of T36 by the clock is confirmed too (R33). One twentieth of the run of A was run anew and matched its file line for line, 3,250 lines.
- **The report** of `sim judge`:

  | candidate | score | group `corridor` | failed "no worse" |
  |---|---|---|---|
  | A | 0.124 [0.115; 0.133] | 0.133 [0.122; 0.144] | 3 |
  | B | 0.130 [0.118; 0.141] | 0.147 [0.134; 0.160] | 30 |
  | C | 0.127 [0.117; 0.138] | 0.137 [0.126; 0.148] | 3 |

  - All three are better than the base by the score.
  - A fails the lag `M3.lag` on `G2` of both maths (on school maths worse by 0.036 at a tolerance of 0.019) and `E4.given` for `F-growth5` (0.01015 at 0.01).
  - C fails the school `E1.corridor` for `G2` (0.01004 at 0.01) and `M3.lag` for `G2` (0.036 at 0.018), and also the olympiad `M3.lag` for `G2-fading`.
  - B fails `E4.given` on almost all the olympiad populations, `E3.retention` for `F-`, `M4.false.all` for `W-repeat` and `M3.lag`. The risk of repeating a seen solution came true (R06 learner-model).
  - For all of them three `E3.end` on school maths left the score — their path is within the tolerance (the 2026-11-09 amendment to R04).

— *Why an exit:* the choice takes only a candidate that holds all the constraints and is better than the base ([architecture.md](architecture.md), section 3.6). There are none such. The common reason of A and C is the model's lag `M3.lag` on `G2`: the new order of the tasks gives the model less evidence where the truth is learning. This is read as a signal for the next round of `engine`, not as an argument to loosen the tolerance: the tolerances were written down before the first candidate (R04).

— *Consequence:* pull requests #46–48 of `engine` are closed by its session, with a link to `judgements/engine-1/`. There is no row «В реестр spec» about branch 3 of C: C is not chosen. T32 judges the model's candidates on the previous base — `runs/base.jsonl` —, because there is no first choice (R35). The next judgement of the engine is `judgements/engine-2/`, and its confirmation is again on `held-out-engine-1`: it is untouched (the 2026-11-10 amendment to R35). Where it is written: RUN.md, T33; `judgements/engine-1/`.

## The judgement `model-1`

**R37. The first judgement of the candidates of `learner-model` — `g10`, `halve035` and `forget-set`, pull requests #43–45 — ended in an exit: no candidate holds all the constraints of the model's criterion, and the base stays as it is. There was no confirmation, and `held-out-model-1` is still held out. `gate.json`, `go.mod` and the bands do not change.** The judgement of T32 (RUN.md), one round by the 2026-11-09 amendment to R103 spec; the report — `judgements/model-1/` (R35 and its amendment). It went after the engine's judgement, which also ended in an exit (R36), so the base is the previous one.
- **The candidates** — all in the group `forgets`, without a target and without a new field; the fields for the judgement file came from the session of `learner-model`. All three drafts are made from `main` 9fb144f, where the model's code is the same as in the release v0.1.4:
  - `g10`, `7f735b20…` — `stabilityGrowth` 10, one constant;
  - `halve035`, `17e0a42f…` — `stabilityCut` 0.35, one constant;
  - `forget-set`, `75c6a908…` — both, two constants.
- **The runs** — decisive, on 1,000 working learners a cell, by the recipe `just decisive` with `CANDIDATE_LEARNER_MODEL` (R34). The engine is taken at the commit of `gate.json`, v0.1.4:
  - the base — `runs/base.jsonl`, SHA-256 `b13a9264…4c92`;
  - `g10` — 41 min by the clock, `e19ff45e…0434`;
  - `halve035` — 44 min, `400d9a6d…d66e`;
  - `forget-set` — 41 min, `cdf7755c…4c25`. Its first start broke off together with the session and was repeated whole.

  One twentieth of the run of `halve035` was run anew and matched its file line for line, 3,250 lines.
- **The report** of `sim judge`:

  | candidate | score | group `forgets` | failed "no worse" |
  |---|---|---|---|
  | `g10` | 0.025 [0.011; 0.037] | 0.195 [0.098; 0.294] | 4 |
  | `halve035` | 0.014 [0.002; 0.027] | 0.117 [0.031; 0.213] | 3 |
  | `forget-set` | 0.033 [0.020; 0.046] | 0.277 [0.184; 0.373] | 8 |

  - All three are better than the base by the score, and the group "Forgets" is above zero for all of them.
  - `halve035` fails three checks, all within a hair of the tolerance: `E3.late` for `F-growth5` (0.04967 at 0.049), `M3.settle` for `F-break` (19.42 at 18.11) and `M1.rms.200` for `A-unpractised` (0.03116 at 0.03076).
  - `g10` fails `M1.rms.200` for `G3-drop`, `G8` and `S-scale1.5` and `M3.settle` for `F-break`.
  - `forget-set` fails the same as `g10`, and also the olympiad `E1.corridor` for `G2`, `G2-fading` and `S-adult` and `M1.rms.200` for `F-none`.
  - Common to all is the catching up after a break, `M3.settle` for `F-break`: with the new numbers of forgetting the model catches up with the truth after a break more slowly. Why — the judgement does not read; this is a question for the next round.

— *Why an exit:* the choice takes only a candidate that holds all the constraints and is better than the base ([architecture.md](architecture.md), section 3.6). There are none such. The tolerances were written down before the first candidate (R04), and that `halve035` failed by a hair is no argument to shift them. This is a signal for the next round of the model: to change forgetting so as not to delay the catching up after a break.

— *Consequence:* pull requests #43–45 of `learner-model` are closed by its session, with a link to `judgements/model-1/`. There are no first choices for either the model or the engine, so T34 — the move to the releases of the chosen ones — the author does not do: its status is `[-]`. The base stays `learner-model` v0.1.4 and `engine` v0.1.4. The next round of the model is `judgements/model-2/`, its confirmation on `held-out-model-1`. Where it is written: RUN.md, T32; `judgements/model-1/`.

## Epic E4

**R38. Epic E4 is «Стенд за библиотеками»: the simulator catches up with `spec`, `engine` and `modules` before new candidates. In it: the move to `spec` v0.14.0, the base on `engine` v0.1.5, a judgement that sees a candidate's second library (Д6), the sizes of the runs after the speed-up of `engine`, the check kind `time` and the leading question on a release of `learner-model`, the pin of `modules` v0.2.1 and the summing up. The second rounds of candidates go as tasks of this epic, when the candidates come. The sizes after the speed-up are by measurement, and the rule of the hour stays.** The author's answers 2026-11-10 when laying out E4 (RUN.md, table 1).
- **What was found in planning** (read 2026-11-10):
  - `spec` has come out up to v0.14.0, and `sim` stands on v0.12.0 (R32);
  - `spec` answered В10 by the row's recommendation (Д167 RUN.md `spec`): `engine` speeds up the choice on a large pool by a task of its own. E4 of `engine` is the profile of `Choose` on `olympiad-math` 0.1.1, a speed-up without a change of the choice, about 10×, and "Gate" on a tag of `sim` that will give the gate back 300 learners (T37–T39 RUN.md `engine`);
  - `spec` answered В11 (Д173 RUN.md `spec`; row 4.10 of its COORDINATION.md). The answers to В10 and В11 came out with the tag v0.14.0, while the submodule stands on v0.12.0, so the rows of table 3 leave by the move T38. The order is this: `learner-model` releases `time` and the leading question in one release through the gate of `sim` (its Д14); then `sim` translates them in the learner's truth and in the evidence and releases a tag; then `modules` releases v0.2.1 (T56 RUN.md `modules`), and `sim` pins itself to it;
  - `engine` released v0.1.5 with the fix of R82 spec in the base: the task of the last attempt with a seen solution does not go next. It was found by the invariant `I4` of the gate of `sim` on the candidates of T31. According to `engine`, all 69 numbers of the gate are the same to the bit, and row В29 RUN.md `engine` asks `sim` and `api` to take v0.1.5;
  - there will be no new candidates soon: `engine` is speeding up the choice, and `learner-model` is recording the outcome of the first judgement.
- **Decided by the author:**
  - **the content** — «Стенд за библиотеками», the second rounds — when the candidates come;
  - **the sizes after the speed-up:** the gate gets 300 learners if it fits into 15 minutes on the runner; the decisive run gets as many as fit into an hour with a margin, but no more than 4,000. New sizes are new paths and new bands.
- **A proposal of the execution for the sizes**, T41 decides it when planning: if 300 learners do not fit into 15 minutes on the runner, the gate gets the largest size that fits, but no less than 100; "with a margin" is a decisive run no longer than 45 minutes on a free machine, a quarter of an hour for a busy machine and the growth of the pools.
- **A proposal of the execution:** the order of the tasks — T38–T44 (RUN.md, «Почему такой порядок»); the second rounds — as tasks with a number free at the moment of writing.

— *Why:*
- **The simulator behind the libraries.** Candidates can be judged only on what the libraries stand on and what the learners will get. `spec`, `engine` v0.1.5 and the packages 0.2.1 have already come out or are coming out, and there are no candidates yet.
- **Sizes by measurement.** R33 cut the runs down to an hour because the choice was expensive. When it becomes cheaper, it is worth giving back as much as the same rule allows, not the previous numbers blindly.

— *Rejected:*
- E4 as T35 proposed it — only the second rounds: the epic would wait for candidates, and the simulator would lag behind the libraries;
- give back 300 and 4,000 blindly — the decisive run could go over an hour;
- leave the sizes alone after the speed-up — the gate would stay coarser than it can be, and `engine` is waiting for a tag with 300 learners.

— *Consequence:*
- E3 is moved to the archive [docs/plan/E3.md](plan/E3.md);
- table 1 of RUN.md — E4 in the tasks T38–T44;
- the row E4 left table 2; В10 and В11 are closed by T38 with the answers that came with the tag `spec` v0.14.0;
- [architecture.md](architecture.md) — the status and subsection 2.5, "Decided"; README, "Status".

## The tag `spec` v0.15.0

**R39. `sim` stands on `spec` v0.15.0 — the latest tag, not v0.14.0, which the task named. The submodule, `go.mod` and `gate.json` are on it. The gate's numbers are the same to the bit, and the bands are as before. Rows В10 and В11 of table 3 of RUN.md are closed by the answers of `spec` that came with the tag.** The author's instruction 2026-11-11: «делай T38, но вышла более свежая версия spec» (RUN.md, T38).
- **What was found in the move** (read 2026-11-11; the sources — [architecture.md](architecture.md), section 4):
  - `spec` has come out up to v0.15.0, the tag is on commit 58d3f46, and `gh release verify` confirmed it;
  - for `sim` the tags v0.13.0–v0.15.0 change nothing that the simulator reads: `RULES.md`, the text and schemas of format 0.4 except the example of traps, the types `formatv0`, the template and the container are the same as on v0.12.0, and card architecture 3.9 spec is untouched;
  - the Go module `spec` added the packages `consentagev1`, `exportv1` and `legalv2` and the parts `Legal`, `ConsentAge` and `Export`, and holds the embedded files in parts. Its `go.mod` is the same, so there are no new dependencies;
  - for `sim` the tags decided four answers of table 4 of RUN.md `spec`:
    - Д167 — the answer to В10: `engine` speeds up the choice by a task of its own, and the gate and the decisive run grow after it. T41 did this: `engine` v0.1.6, the gate on 150 learners, the decisive run on 1,500 (the 2026-11-11 amendments to R31 and R33);
    - Д199 — the same move to `engine` v0.1.6: T41 too;
    - Д198 — the answer to В32 of `engine`: `learner-model` speeds up `Readout.Task` by a task of its own, and `sim` gives the gate back 300 learners and the decisive run 4,000 when both edits come out. This is the debt Д7 of table 4 of RUN.md;
    - Д173 — the answer to В11: row 4.10 of COORDINATION.md `spec`. First `learner-model` releases `time` and the leading question in one release through the gate of `sim`, then `sim` translates them and releases a tag, then `modules` releases v0.2.1. These are T42 and T43;
  - the libraries do not require the Go module `spec`: the `go.mod` of `learner-model` v0.1.4 is without `require`, and `engine` v0.1.6 requires only `learner-model`. `spec` is a dependency of the simulator itself, and its version is the tag of the submodule (R64 spec);
  - `just gate exact` on v0.15.0 — all 69 numbers are the same to the bit, and the bands were not written anew.
- **Decided by the author:** the move straight to v0.15.0.

— *Why:* the repository stands on the latest tag of `spec` (RULES.md, Across repositories), and between v0.14.0 and v0.15.0 nothing changes for the simulator. The move to v0.14.0 would only put off the next one.

— *Rejected:* v0.14.0, as the task named it — one more move right after.

— *Consequence:* В10 and В11 left table 3 of RUN.md; `gate.json` names `spec` v0.15.0; the tag — in CLAUDE.md, README and the opening of RUN.md. The merge releases a tag of `sim`, and the gates of the libraries take it. Where it is written: RUN.md, T38; [architecture.md](architecture.md), the status, section 2 and section 4; README.

— *Amendment 2026-11-12:* `sim` stands on `spec` v0.16.3 (R45): the submodule, `go.mod` and `gate.json` are on it, Go in the container is 1.27.4, and the gate's numbers are the same to the bit.

## A candidate's second library

**R40. A judgement holds the base and the candidates to the libraries of the base. `sim` embeds `gate.json`, as it embeds `paths.json`, and `sim judge` refuses the whole judgement if: the base or the base of the confirmation runs on other libraries than those on which the paths were read; the base is not on the versions that the gate file names; the second library of a candidate of one library — or of the same candidate on the held-out learners — is neither the base's version nor the commit of its pin. A pair has no second library. The headers of the results files, the recipe and the bytes of the runs do not change.** The author's answers 2026-11-11 when planning T40 (RUN.md, T40); the rest is a proposal of the execution. It closes Д6 of table 4 of RUN.md.
- **What was found in planning** (2026-11-11):
  - the base's header names both libraries by versions — the recipe `just decisive` builds the base by `go.mod` —, and a candidate's header names both by commits: the recipe takes the second library by the commit from `gate.json` (R34). The judgement cannot compare them directly: only the gate file matches a version with the commit of its tag;
  - the judgement did not check the base against the one on which the paths were read either, though `paths.json` names its versions (a finding of the review of T41);
  - while a library has no release, `go.mod` names it by the pseudo-version of a commit, and `sim gate verify` already matches a pseudo-version with a commit by the suffix of 12 characters;
  - the packages `criterion` and `gate` do not import each other.
- **Decided by the author:**
  - the source is `gate.json`, embedded in `sim`;
  - a wrong second library is a refusal of the whole judgement, not a verdict;
  - along the way the judgement checks the base's libraries against the paths, but not the SHA-256 of the base's file.
- **A proposal of the execution:**
  - `main` embeds `gate.json` and passes the judgement the pins of the libraries — `criterion.Pin`, the commit and the tag, by the module path;
  - `criterion.Judge` calls the check of the libraries after the paths. Without the pins the judgement does not judge. The check against the gate file accepts the tag, and while there is none, the pseudo-version of the commit, as `sim gate verify` does; for a candidate's second library also the commit of the pin, but not an empty name;
  - the refusal names the candidate, the library and both points: where the candidate ran and where the base did.

— *Why:*
- **`gate.json`.** It alone names both the version and the commit of its tag, and the candidates' runs are built exactly by it (R34). The headers, the recipe and the bytes of the runs stay as before, and the pinned run of the base, whose SHA-256 `paths.json` names, is the same (R28). The embedded file judges by the pin with which `sim` is built — as the paths do.
- **A refusal of the whole judgement.** A run on someone else's second library is not a candidate of this judgement: it is an error of the input, like a run on the wrong learners, not a verdict on the candidate.
- **The base against the paths.** A candidate's share is its difference from the base divided by the pinned path. On a base of other libraries the path is not the right one, and the share would be someone else's.

— *Rejected:*
- a field in the candidate's header "taken from the caller" — it changes the recipe and the header, and the previous candidates do not have it;
- a path to the gate file in the judgement file — the wrong one can be given;
- a check "wrong library" with "not met" for one candidate — the judgement would give a verdict on a run that is not of it;
- checking the SHA-256 of the base's file too — rough looks and tests would need a way around, and on someone else's base the check of the libraries already refuses;
- only Д6, without the check of the base against the paths — see the "Why".

— *Consequence:*
- the judgements `engine-1` and `model-1` on the present tree get a refusal: their base is on `engine` v0.1.4, and the paths on v0.1.6. Their reports are records on their own commits, and there they are reproduced;
- the recipe's protection stays: a pair requires `CANDIDATE_PAIR=yes` (R34);
- the 2026-11-11 amendments to R25 and R35; Д6 is closed.

Where it is written: RUN.md, T40; [architecture.md](architecture.md), sections 3 and 3.6; README; `main.go`, `internal/criterion/criterion.go`.

## Format 0.4 on the simulator

**R41. The simulator reads the check kind `time` and the leading question of format 0.4 on `learner-model` v0.1.5. `time` is a check kind of its own, and the truth guesses it with a probability of 0, as a number. The simulator reads the leading question only by whether the task has one: a request of `G8` before the first answer on such a task gives both the technique and the question, the evidence carries `LeadingQuestionOnRequest`, and the rise of the logit is as before. There is no decisive run here, and `paths.json` stays on `learner-model` v0.1.4: until T43 the judgement refuses.** The author's answers 2026-11-11 when planning T42 (RUN.md, T42); the rest is a proposal of the execution. It carries out row 4.10 of COORDINATION.md `spec` in the part of `sim` (Д173 RUN.md `spec`; R39).
- **What was found in planning** (2026-11-11; the sources — [architecture.md](architecture.md), section 4):
  - `learner-model` v0.1.5 adds `CheckKindTime` without guessing and the field of evidence `LeadingQuestionOnRequest` — "a hint the learner asked for gave the task's leading question". The model counts it as help, as a technique on request: the answer goes into the step, but does not set "mastered". The kind of evidence is 2; by R28 learner-model the model's numbers do not change;
  - the leading question is given only by a hint on request before the first error ([format text](../spec/format/README.md), sections 14 and 19). On the simulator only `G8` asks for it, and the base learner never does ([architecture.md](architecture.md), section 3.1);
  - `modules` v0.2.1 with `time` in 30 tasks and leading questions in 311 tasks came out 2026-11-10, before row 4.10 (row В38 RUN.md `modules`), so T43 waits only for this record.
- **Decided by the author:**
  - the leading question — one request, one rise, without a new assumption;
  - there is no decisive run in T42: the paths are read once by T43, already on `modules` v0.2.1.
- **A proposal of the execution:**
  - `release-verify` accepted `learner-model-proto-v0.1.5.tar.gz` — "learner-model v0.1.5, signed by the shared release.yml of .github v0.1.21"; `go.mod` and `gate.json` are on v0.1.5, the commit of the tag `3a2cd29`. `just gate exact` — all 69 numbers are the same to the bit. The bands are recorded anew by `just gate record` only for the sake of the reason: the values in them are the same to the bit, and the gate file now names the commit on which they are recorded;
  - the translation: `time` — `CheckKindTime`; `pack.Task.Leading` — the task has a `leading_question`. The canonical text holds both (the 2026-11-11 amendments to R14 and R27);
  - the truth: `Attempt.LeadingQuestion` — the technique of `G8` on a task with a question, without a new number from the place `attempts`; the evidence carries it in `LeadingQuestionOnRequest`, of kind 2 (the 2026-11-11 amendment to R18). The truth does not know traps, so there is no 12-hour miss for `time`;
  - the synthetics: the skill of the sheet on the anchor skill is answered by a time, and the tasks of the sheet of difficulty 5 have a leading question (the 2026-11-11 amendment to R19);
  - the judgement's test in `main_test.go` holds its runs and the embedded gate file to the libraries of the paths: between the move of a library and the decisive run that reads the paths anew, they diverge.

— *Why:*
- **`time` with a guessing of 0.** An answer of time is typed, not chosen, like a number, and the model guesses it the same way — 0 (R19 learner-model). In the rest the truth holds the model's assumptions ([architecture.md](architecture.md), section 3.1).
- **One request, one rise.** By the format text the leading question comes with the same hint as the techniques, and the model counts them as one help. How much the question helps beyond the technique there are no data — neither at the simulator nor at `learner-model`, where R06 and R21 learner-model rejected a weight without data. Without a new assumption the numbers of `G8` do not shift in T43 either, and the simulator judges the task as the learner will get it: the evidence carries what the server will record.
- **Only "there is a question".** The question's words are for the learner, and the simulator reads them no more than the text of the statement. So an edit of the words does not wake the gate of the CI of `modules`, while an added or removed question does.
- **Without a decisive run.** T43 pins the packages with `time` and questions, and these are new paths in any case. A run here on `modules` v0.2.0 would cost about an hour of the machine for paths that T43 would rewrite at once. There are no candidates now, so the judgement's refusal until T43 holds nobody up.

— *Rejected:*
- a rise of its own for the leading question beyond the technique — a value without data;
- a new population that asks for the question — one more cell of the decisive run, and also without data;
- translate `time` as `number` — then the simulator would measure another task (option (b) of [architecture.md](architecture.md), subsection 2.5);
- a new skill of the synthetics with `time` — it would shift the engine's choice for the learner of both synthetics, and the anchor's test would stop seeing the mastered anchor;
- the decisive run in T42, and T43 inside T42 — see the "Why" and RUN.md, T42.

— *Consequence:*
- until T43 reads the paths anew, `sim judge` refuses any judgement: the base is on `learner-model` v0.1.5, and `paths.json` was read on v0.1.4 (R40);
- the 2026-11-11 amendments to R14, R18, R19 and R27; [architecture.md](architecture.md), subsection 2.5, is done;
- the merge releases a tag of `sim`, and the gates of the libraries take it by tasks of their own. There is no row «в реестр spec»: the order that the release of `modules` v0.2.1 broke has already been raised by row В38 RUN.md `modules`.

Where it is written: RUN.md, T42; [architecture.md](architecture.md), subsection 2.5 and sections 3, 3.1–3.3; README; `internal/pack/translate.go`, `internal/learner/attempt.go`, `internal/run/run.go`, `internal/synthetic/synthetic.go`.

— *Amendment 2026-11-11 (T43):* the judgement no longer refuses. The paths are read anew on `learner-model` v0.1.5 and the packages of `modules` v0.2.5 (the 2026-11-11 amendment to R04), and the judgement's test holds its runs to the libraries of the paths, as before — now these are the same as in `gate.json`. Why — R42.

— *Amendment 2026-11-11 (T45):* the translation and the truth are the same, and the base is on `learner-model` v0.1.6, where `Readout.Task` is faster with the same bits: both decisive runs and the gate matched the runs on v0.1.5 to the bit. The paths are read anew on it, so the judgement does not refuse. Why — R43.

## The pin of `modules` v0.2.5

**R42. The simulator is pinned to the release of `modules` v0.2.5 — `olympiad-math` and `grade-3-math` 0.3.0 with the checks `time`, leading questions and worked examples, and `english-essays` 0.2.0. The paths are read anew on it and on `learner-model` v0.1.5, and so are the gate's bands. The sizes of the runs are as before: the gate — 150 learners a cell, the decisive run — 1,500.** The author's answers 2026-11-11 when planning and executing T43 (RUN.md, T43); the rest is a proposal of the execution. Together with R41 it carries out the part of `sim` of row 4.10 of COORDINATION.md `spec` (Д173 RUN.md `spec`).
- **What was found in planning** (2026-11-11; the sources — [architecture.md](architecture.md), section 4):
  - `modules` v0.2.1 came out 2026-11-10 with both maths at 0.3.0. The index of v0.2.5 is byte for byte the same as that of v0.2.1 — SHA-256 `cebf23cf…aec6a`: in v0.2.2–v0.2.5 only the synthetic samples and the move to `spec` v0.15.0 changed;
  - `olympiad-math` 0.3.0 has 775 tasks instead of 772 and 22 templates instead of 19, a leading question in 322 tasks; checks `number` — 572, `choice` — 166 instead of 196, `fraction` — 7, `time` — 30;
  - `grade-3-math` 0.3.0 has the same tasks and checks, and in each of the four lessons a new segment — a worked example, two minutes, as any segment without a task (R15);
  - the simulator reads the new packages without edits of the code, and a sample of the decisive run on a copy of the repository gave all 65 cells.
- **Decided by the author:**
  - the pin — the latest release, v0.2.5, not v0.2.1, which the task named;
  - reproducibility — by a full repeat of the ceiling;
  - in the execution: the decisive run gets 1,500 learners as before, though the base took 47 min 48 s. Neighbouring containers shared the host at that time, and the run holds the rule of the hour.
- **A proposal of the execution:**
  - `release-verify` accepted the index — "module-index.json: modules v0.2.5, signed by the shared release.yml of .github v0.1.15" — and both archives from v0.2.1 — "modules v0.2.1, signed by the shared release.yml of .github v0.1.15". The archive of the essays 0.2.0 is the same as the one accepted in R32;
  - the index lies in `internal/published/v0.2.5/`, the archives in the folders of the releases that their records name: `v0.2.1/` and `v0.2.0/` (the 2026-11-11 amendment to R14);
  - the tests of the facts of the pin — on the new packages, with the facts of 0.4: the checks `time`, the tasks with a leading question, the segment of a worked example. The test of a run's input holds the instance and the fact "solution seen" by whether the task is a template: the new template `squares-on-lines` came into the choice of its learner;
  - the decisive run, the new paths — the 2026-11-11 amendment to R04 — and the gate's bands (the 2026-11-11 amendments to R26 and R28). The ceiling was repeated whole and matched to the byte. The gate in the container ran for 7 min by the clock on an overloaded host, 28.4 CPU minutes — against 23.6 in T41; the CI of the pull request will show the time on the runners.

— *Why:*
- **v0.2.5.** The index is the same, and the simulator names the latest release of `modules`, as it stands on the latest tag of `spec`.
- **Now, not later.** Since T42 the judgement refuses until the paths are read anew. The pin requires new paths anyway, so one run gives both the packages that the learners will get and the paths on the new library.
- **1,500 at 47 minutes.** The 45 minutes of R38 are the margin on a free machine, and the run went on a busy host: the CPU per process of the base is about 38.7 minutes against 33 in T41, and part of the difference is the neighbours. The run holds the rule of the hour, while 1,400 learners would make the resolution of "no worse" coarser and would cost one more run of the base.

— *Rejected:*
- v0.2.1 — the same index, but not the latest release;
- one twentieth of the base and of the ceiling instead of a full repeat — about 17 minutes saved;
- 1,400 learners or a new measurement on a free host — see the "Why".

— *Consequence:*
- the judgement judges again (the 2026-11-11 amendment to R41);
- the thresholds of the CI of `modules` are taken from the new record of the bands (R27);
- the merge of the epic's branch releases a tag of `sim` with T42 and T43, and the gates of `learner-model` and `engine` take it by tasks of their own; row 4.10 of COORDINATION.md `spec` — by row В12 of table 3 of RUN.md, «в реестр spec».

Where it is written: RUN.md, T43; [architecture.md](architecture.md), the status, subsection 2.5 and sections 3, 3.3, 3.8 and 4; README; `internal/published/`, `gate.json`, `internal/criterion/paths.json`.

## The base on `learner-model` v0.1.6

**R43. The base stands on `learner-model` v0.1.6, where `Readout.Task` is faster with the same bits. The paths are read anew on it, and so are the gate's bands, and the numbers did not shift by a single bit. The sizes of the runs are as before: the gate — 150 learners a cell, the decisive run — 1,500; the measurement of epic E5 raises them (Д7).** The author's answer 2026-11-11, after the summing up of E4: «сделай что сейчас можешь из таблицы 4» and the move to v0.1.6 with the paths, not only a preliminary measurement (RUN.md, T45); the rest is a proposal of the execution.
- **What was found** (2026-11-11; the sources — [architecture.md](architecture.md), section 4):
  - `learner-model` v0.1.6 came out 2026-11-11, the tag on commit 7d796e2 — "Make Readout.Task faster with the same bits (T40)". Row В39 of its RUN.md promises the same API and the same P to the bit. The module's `go.mod` at the tag has no requirements — only `go 1.27.3` —, so it does not move the tag of `spec`;
  - `release-verify` accepted `learner-model-proto-v0.1.6.tar.gz` — "learner-model v0.1.6, signed by the shared release.yml of .github v0.1.21";
  - of the two debts of table 4, only a part of Д7 is done now: Д3 waits for stage 2, and the sizes wait for a measurement on a free machine and on the runners. Neighbouring containers shared the host that day, the load average 25–46 on 20 cores.
- **A proposal of the execution:**
  - `go.mod`, `go.sum` and `gate.json` — on v0.1.6 (7d796e2); the tests that name the pin — too;
  - `just gate exact` — all 69 numbers are the same to the bit; the bands are recorded anew only for the sake of the reason (the 2026-11-11 amendment to R26);
  - the decisive runs of the base and of the ceiling on 1,500 learners — the same bytes of the lines as in T43, and new paths with the same values (the 2026-11-11 amendments to R04 and R28).
- **A measurement of the speed-up, preliminary.** The CPU per process of the base is about 25.9 minutes against 38.7 in T43, a third less; of the ceiling 18.0 against 18.8: the ceiling hardly calls `Readout.Task`. By the clock the base took 62 min 24 s — longer than an hour — at a load average of up to 46. The gate in the container — 1 min 16 s by the clock and 18.7 CPU minutes against 28.4 in T43. This is not the measurement of Д7: that one goes on a free machine and on the runners, and by it E5 chooses the sizes (the 2026-11-11 amendments to R31 and R33).

— *Why:*
- **Now, not in E5.** The move costs the numbers nothing and takes one change of the pin off E5: the second round of the model is judged on one pin of `sim` (row В42 RUN.md `learner-model`), and the measurement of the sizes must already go on v0.1.6. The decisive runs also gave the first estimate of the speed-up on the previous sizes.
- **The paths anew, though the numbers are the same.** The judgement holds the base to the libraries of the paths (R40): without new paths it would refuse every base on v0.1.6, as from T42 to T43.
- **Without a second run of the ceiling.** The new ceiling matched the ceiling of T43 to the byte of the lines, and that one was itself repeated whole (R42).
- **62 minutes under the rule of the hour.** The time by the clock is the neighbours': the base took a third less CPU than in T43, where the run fitted into 47 minutes. The size of the decisive run does not change by this measurement, neither up nor down.

— *Rejected:*
- only a preliminary measurement on a copy of the repository, without the move — the pin would stay on v0.1.5, and the change would fall on E5;
- the move without new paths — the judgement would refuse;
- raise the sizes by this measurement — the host is busy, and the rule of the 2026-11-11 amendments to R31 and R33 requires a free machine and the runners.

— *Consequence:*
- Д7 stays: E5 raises the sizes, and its first task is now the measurement and the sizes, without the move to v0.1.6;
- the merge of the branch of epic E4 releases a tag of `sim` with this move too.

Where it is written: RUN.md, T45; [architecture.md](architecture.md), the status and sections 3, 3.6–3.9 and 4; README; `go.mod`, `gate.json`, `internal/criterion/paths.json`.

## Epic E5

**R44. Epic E5 is «Размеры и второй круг модели»: the sizes of the runs on `learner-model` v0.1.6 by measurement (Д7), the judgement `model-2` with a confirmation on `held-out-model-1`, the move of the base to the release of the chosen one and the summing up — T46–T49. The engine's judgement `engine-2` is a separate epic E6 of table 2 of RUN.md, on the base after `model-2`.** The author's answers 2026-11-11, after the summing up of E4 and T45: «разверни эпик E5», and the content — the sizes and `model-2`, with `engine-2` separately (RUN.md, table 1, «Решено автором»); the rest is a proposal of the execution.
- **What was found** (2026-11-11; RUN.md of `learner-model` and of `engine` on `main`, read through `gh api`):
  - E5 of `learner-model` — «Второй круг кандидатов модели» — is laid out: its T41 waits for a tag of `sim` with the sizes after v0.1.6, T43 is the exploration after the analysis of `M3.settle` for `F-break`, T44 up to three draft pull requests of candidates, T45 the judgement of `sim` and the release of the chosen one;
  - `engine` is in its E7, «Закрепление в уроке», and its E5 — «Второй круг коридора» — stands in table 2 not laid out: no candidates of the engine are in sight;
  - `spec` has not recorded rows В42 and В43 RUN.md `learner-model`, «в реестр spec», yet. Their recommendations are one pin of `sim` per round of the model and the round of the model before the round of the engine;
  - the base is already on `learner-model` v0.1.6 (R43), so the sizes are measured without a move.
- **A proposal of the execution:** the order of the tasks — RUN.md, «Почему такой порядок»; one pin from the tag of T46 to the choice of the judgement `model-2` — the epic's «Общее для задач»; the gate's size — from 100 to 300, as in the 2026-11-11 amendment to R31, not from 150, as in the row E5 of table 2 that T44 proposed.

— *Why:*
- **The sizes first.** The exploration and the judgement of `learner-model` go on the sizes of the decisive run of `sim`, and "Gate" of `learner-model` waits for a tag with them. The sizes move the paths and the bands, and in the middle of a round the pin does not change (В42 RUN.md `learner-model`).
- **`engine-2` separately.** The round of the model first — by the recommendation of В43 RUN.md `learner-model`, and the round of the engine is judged on the base after it. And there are no candidates of the engine even in the plan: an epic that waits for them would stand open without a date, and with it the summing up of the sizes and of the round of the model.

— *Rejected:*
- all three rounds in one epic, as T44 proposed (row E5 of table 2 of RUN.md as of commit 77828ca) — see the "Why";
- the judgement `model-2` before the sizes — the round would have to be judged anew or on rough sizes, and "Gate" of `learner-model` waits for a tag with them anyway.

— *Consequence:* E4 is in the archive [docs/plan/E4.md](plan/E4.md); the row E5 of table 2 of RUN.md went to table 1, and in its place stood E6 «Второй круг движка»; T46 closes Д7.

Where it is written: RUN.md, tables 1, 2 and 4 and «Эпик E5»; [architecture.md](architecture.md), the status and section 3; README.

— *Amendment 2026-11-12:* "one pin per round" has an exception, R45: the version of `spec` in `go.mod` and `gate.json` is raised to v0.16.3 in the middle of the round `model-2`, because the bits are the same; the packages, the libraries, the paths and the bands are as before. The move and the rebuilding of the documents by row 7.10 of COORDINATION.md `spec` are the tasks T50–T53, outside an epic.

— *Amendment 2026-11-13 (R47):* RUN.md, «Почему такой порядок» and table 1, «Решено автором», moved word for word to [docs/plan/E5.md](plan/E5.md), «Из RUN.md», when RUN.md was shortened; what they say is in force.

## The tag `spec` v0.16.3

**R45. `sim` stands on `spec` v0.16.3 — the latest tag, of the line 0.16 of epic E10 of `spec`. The submodule, `go.mod` and `gate.json` are on it, and Go in the container is 1.27.4; the gate's numbers are the same to the bit, so the paths, the bands and the runs of the base stay. Row 7.10 of COORDINATION.md `spec` — the move, the rebuilding of the documents and the translation of `docs/` into English — `sim` takes before the pilot of `module-kit` and splits into four tasks outside an epic, T50–T53; the way of working by the template v0.16.3 applies already from T50.** The author's instruction 2026-11-12 — «переходи на spec последней версии» — and their answers to two questions of the move (RUN.md, T50).
- **What was found in the move** (`git -C spec diff v0.15.0 v0.16.3`, read 2026-11-12; the sources — [architecture.md](architecture.md), section 4):
  - `spec` has come out up to v0.16.3, the tag is on commit 664398e, and `gh release verify` confirmed it;
  - the line 0.16 is epic E10 of `spec`, the saving of the sessions' tokens (R143 spec, R145 spec). The move to it is, for each repository, a row of item 7 of COORDINATION.md `spec`, for `sim` 7.10: one task outside an epic by the shared text T123 of `spec`, with the move to the tag, the rebuilding of the documents on the model of `spec` and the translation of `docs/`. Rows 7.2–7.12 wait for the pilot — row 7.1, `module-kit`. At about 20:30 UTC 2026-11-12 `module-kit` stood on v0.15.0 without a pull request of the move, and at 20:59 UTC it moved to v0.16.3 only by the tag — Go 1.27.4 and Д215 —, leaving the rest of row 7.1 to its task T43 (R54 module-kit): the pilot is not passed yet. A large repository, by the text of T123, proposes before starting to split the task by rule 11 of its RUN.md;
  - the tag changes nothing that the simulator reads: `format/`, `types/formatv0` and the root package of `spec` are the same, and what changed are the contracts `api/v1` and `admissionv0`, which the simulator does not take. The `go.mod` of the module `spec` is the same — there are no new dependencies. The header of a run's file names only the libraries, so the runs of the base and of the ceiling on 2,000 learners are in force;
  - the Go of the stage `toolchain` is 1.27.4: 1.27.3 has the vulnerability GO-2026-6859 of the standard library, and the move to the tag brings 1.27.4 to every Go repository together with `go 1.27.4` of the template's lock file of gopls. The directives `go` in `go.mod` do not change — they name the lowest version, not the container's version (the 2026-11-11 amendment to R64 spec);
  - RULES.md: one task is one session, and the hook `UserPromptSubmit` stops a phrase that starts another; a decision that changes another entry puts an amendment on it; external facts — as records of `docs/sources.md`, the results of measurements — in `docs/measurements.md`; `docs/` is in English, and a document in Russian stays so until its task of translation; `spec` triages its inbox on the author's word or at the end of an epic (R143 spec);
  - Д215 of `spec` — the line `**/sim/spec/.claude/rules/**` in `claudeMdExcludes`: the rules of the parts of `spec` are not loaded into a session of `sim`. Rows Д140, Д196 and Д198 of `spec`, whose "when" is a move or a task of `sim`, were carried out earlier: R32, R39–R43, T46;
  - the answer to В12 of table 3 of RUN.md is by the recommendation, (a): row 4.10 of COORDINATION.md `spec` is closed by the tag `sim` v0.1.9, and row 3.9 no longer waits for it (Д173 in DEBT.md `spec` and table 5 of RUN.md `spec`).
- **Decided by the author (2026-11-12):**
  - the whole row 7.10 — now, before the pilot, and in four tasks outside an epic by rule 11: T50 — the tag and the way of working, T51 — `sources.md`, `measurements.md` and the translation of architecture, T52 and T53 — the translation of the log, the budgets and the measurement "after". About 650 KB of Russian text — the log 424 KB, architecture 225 KB — do not fit into one check;
  - an exception to the rule of epic E5 "one pin per round": the version of `spec` in `go.mod` and `gate.json` changes in the middle of the round `model-2` if the bits are the same. The condition is `just gate exact`, all the numbers to the bit, and the first 100 learners of each cell on the new simulator byte for byte the same as their lines in the decisive run of the base on 2,000, as in R32.
- **What is changed here:**
  - `spec/` — to v0.16.3; `go.mod` and `go.sum` — the Go module `spec` v0.16.3; `gate.json` — `spec` v0.16.3, the bands were not written anew; `gate_test.go` — the version of `spec` in the build that the tests substitute is taken by `specRecorded` from the embedded `gate.json`, not as a literal, and the next move will not touch it;
  - `.devcontainer/Dockerfile` — `GO_VERSION` 1.27.4 and the sums, as in `spec`; `.devcontainer/gopls/go.mod` — `go 1.27.4`, as in the template, and `go.sum` is the same;
  - `.claude/settings.json` — by the template and the kind `go`: compaction at 400 thousand tokens, no claude.ai connectors, the hook `spec/tools/session/`, the plugin `gopls-lsp` and the line of Д215;
  - CLAUDE.md — by the template: the documents with `DEBT.md` and `sources.md`, the budgets and the rules of the parts `.claude/rules/` — `go`, `gate` and `ci`; AGENTS.md — the template's, with a line about the rules by paths;
  - RUN.md — «Как пользоваться» and «Как устроен план» by the template: one task is one session, the body of a done task goes at once to the epic's archive, table 4 is in a file of its own, `DEBT.md`; the body of T46 — in the new archive [docs/plan/E5.md](plan/E5.md); the section «Задачи вне эпика» with T51–T53; В12 is closed, В13 is in table 3;
  - the tag — in CLAUDE.md, README.md, the opening of RUN.md and the status of [architecture.md](architecture.md).
- **The bits are the same:** `just gate exact` on v0.16.3 — all 69 numbers to the bit; the first 100 learners of each of the 65 cells on the new simulator — 6,500 lines — byte for byte the same as their lines in the decisive run of the base of T46 on 2,000, `runs/base.jsonl`. So it is both on the container's Go 1.27.3 and on Go 1.27.4, on which CI will build the simulator: its toolchain was taken from the Go proxy through `GOTOOLCHAIN=go1.27.4`. The condition of the exception is met, and the bands were not written anew.

— *Why:* the author's instruction; the repository stands on the latest tag of `spec` (RULES.md, Across repositories), and the saving of E10 for `sim` is the sessions of every next task. In the middle of the round `model-2` — because the rule "one pin" holds the numbers on which the round is judged, and they did not shift by a bit, and there are no candidates of the second round yet. Before the pilot — the author's word; the pilot catches errors of the shared text T123, and what T50–T53 find, they write down in «Сделано».

— *Rejected:*
- wait for the pilot 7.1 — the author's decision: `module-kit` has so far taken only the tag (R54 module-kit), and its T43 is not done;
- wait for the end of the judgement `model-2`, as the rule of epic E5 says — it holds the numbers of the round, and they are the same;
- the whole row 7.10 as one task, as in the text of T123 — more than one check (rule 11);
- only the tag, and the rebuilding and the translation after the pilot, as `module-kit` did (R54 module-kit): the author's decision is the whole row now.

— *Consequence:* the 2026-11-12 amendments to R02, R39 and R44; В12 left table 3, and В13 «в реестр spec» names to `spec` the tasks of row 7.10. The merge releases a tag of `sim` with the same numbers: "Gate" of `learner-model` may stay on v0.1.10, and the judgement `model-2` goes on the same runs of the base. T51–T53 — each in a session of its own. Where it is written: RUN.md, «Задачи вне эпика», «Общее для задач» of E5 and table 3; [docs/plan/E5.md](plan/E5.md); CLAUDE.md; `.claude/`; DEBT.md; README.md.

— *Amendment 2026-11-13 (R47):* RUN.md no longer has the section «Задачи вне эпика»: the phrases of T50–T53 moved word for word to [docs/plan/E5.md](plan/E5.md), «Из RUN.md», and a task outside an epic now stands in table 1 of the current epic. The rest of the way of working by the template v0.16.3 is in force, in a shorter RUN.md.

## The judgement `model-2`

**R46. The second judgement of the candidates of `learner-model` — one candidate, `theta`, pull request #62 — chose `theta`, and the held-out learners `held-out-model-1` confirmed it. The base is now `learner-model` at the commit of `theta`, a063835: by its pseudo-version in `go.mod` and by the commit, without a version, in `gate.json`. The paths and the bands are read anew on it, and `held-out-model-1` is no longer held out.** The judgement of T47 (RUN.md), one round by the 2026-11-09 amendment to R103 spec; the report — `judgements/model-2/` (R35). It went on the sizes of T46 and on one pin of `sim` from `v0.1.10` (R44; the exception of R45).
- **What was found in planning** (2026-11-13; RUN.md of `learner-model` and its pull request #62, read through `gh`):
  - T44 of `learner-model` made one draft, `theta`, commit `a063835087c7df7f2eecd1912c44afda7d96bfe8`, from its `main` 45ae4d6: after a pause of a module of at least 14 days its weight of answers n falls to 20 before the step — the level θ of the module forgets. Group `jumps`, no target, two constants, no new field (R32 learner-model). Its "Gate" on `sim` v0.1.10 is green. So the questions of the task — how many candidates, in what order, what if they come apart — have one answer: one candidate;
  - the session of `learner-model` had already run the decisive run of the same commit on a clone of `sim` v0.1.10: the file SHA-256 `35051542…7c19`, its body past the header `c918fb7e…428c`;
  - a judgement holds the base to the versions of the gate file (R40), and a candidate's run names both libraries by commits (R34), so it cannot be the base of the new paths: after the move the base is run anew by `go.mod`.
- **The runs** — decisive, 2,000 learners a cell, by the recipe `just decisive` (R34), on `sim` after v0.1.12; their time is measurements 1.3:
  - the base — `runs/base-v016.jsonl`, the run of T46, SHA-256 `7db81ec7…c92e`, the one that `paths.json` named;
  - `theta` — SHA-256 `35051542…7c19`, to the byte the run of the session of `learner-model`: the run is reproduced on another clone and another tag of `sim`;
  - the confirmation on `held-out-model-1`, seed 20261108 — the base `a2b806dd…ced7`, `theta` `640f33d2…5152`.
- **The report** of `sim judge`:

  | run | score | group `jumps` | failed "no worse" |
  |---|---|---|---|
  | `theta`, working learners | 0.0103 [0.0052; 0.0152] | 0.0205 [0.0178; 0.0233] | 0 |
  | `theta`, `held-out-model-1` | 0.0422 [0.0328; 0.0522] | 0.0203 [0.0177; 0.0229] | 0 |

  - `theta` holds every constraint on both sets and is better than the base on both.
  - The catching up after a break, `M3.settle` of `F-break`, the common failure of `model-1` (R37), is 12.7 answers faster on olympiad maths and 42.6 on school maths: these are the new paths, measurements 3.1.

— *Why chosen:* the choice takes a candidate that holds all the constraints and is better than the base, and the confirmation asks the same of it on learners never looked at ([architecture.md](architecture.md), section 3.6). `theta` is the only candidate and passed both. The confirmation does not ask the score to repeat: on the held-out learners it is four times higher, and the candidate's own group, `jumps`, is the same on both sets.

— *Rejected:*
- the base of the new paths from the candidate's run — the judgement would refuse it (R40), see the planning;
- a full repeat of the ceiling — the new base matched the candidate's run to the byte past the header, and one twentieth of the new ceiling matched it line for line, as in R36 and R37.

— *Consequence:*
- `go.mod` — `learner-model` `v0.1.7-0.20261113040121-a063835087c7`; `gate.json` — commit a063835 without a version; the tests that name the pin;
- the paths anew — the 2026-11-13 amendments to R04 and R28, measurements 3.1; the bands anew — the 2026-11-13 amendment to R26;
- the merge releases a tag of `sim` with the bands of `theta`: T45 of `learner-model` raises "Gate" of #62 to it, merges it and releases `v0.1.7`, and T48 moves the base to that release. There is no row «В реестр spec» (R35);
- the next round of the model is `judgements/model-3/`, with its confirmation on `held-out-model-2`; the judgement `engine-2`, epic E6, goes on the base with `theta` (R44);
- on the present tree `sim judge` refuses the judgement file of `model-2`: its base is on `learner-model` v0.1.6, and the paths on a063835 (R40). The report is a record on 9acc4af, `main` before T47, where the judgement file gives it again, as for `engine-1` and `model-1` on their own commits; `runs.txt` names it.

Where it is written: RUN.md, T47 ([docs/plan/E5.md](plan/E5.md)); `judgements/model-2/`; [architecture.md](architecture.md), the status and sections 3, 3.6–3.8; [measurements.md](measurements.md), 1.3 and 3.1; README; `go.mod`, `gate.json`, `internal/criterion/paths.json`, the tests.

— *Amendment 2026-11-13 (T48):* the release of `theta` is `learner-model` `v0.1.8`, not `v0.1.7` as the consequence above says: `v0.1.7` went to the move of `learner-model` to `spec` v0.16.3, with the code of v0.1.6, and T45 of `learner-model` merged #62 after it, as the commit ab688b2 (sources 3.8). The code of the release is that of a063835 — no Go file and no `go.mod` changed between them. The base now stands on the release: `go.mod` names `learner-model` v0.1.8, and `gate.json` its commit ab688b2 with the version; `release-verify` accepted the archive. `just gate exact` is green on it, all 69 numbers to the bit, and the bands were recorded anew only for the sake of the reason (the 2026-11-13 amendment to R26, T48). The paths were read anew on the release, to the same values (the 2026-11-13 amendment to R04, "the release of `theta`"). The judgement `engine-2`, epic E6, goes on this base.

## A shorter RUN.md

**R47. RUN.md keeps only what is live and what no other file says: the executor's rules, the current epic, the tables 2 and 3, and the bodies of the tasks not yet done. The section «Задачи вне эпика» is gone: a task outside an epic stands in table 1 of the current epic and goes into its archive like the others. Rules that repeat RULES.md or CLAUDE.md are a line with a pointer, and their numbers stay. The description of the current epic and the history of the tables go to the epic's archive word for word, not lost.** The author's instruction 2026-11-13, given with T47: «удали секцию Задачи вне эпика - она больше не нужна» and «почисти RUN.md от ненужной информации чтобы сэкономить токены».
- **What changed** (RUN.md, 37.2 KB → 19.7 KB):
  - the opening — without the paragraph on the first epic and without the repeat of "one task is one session";
  - rules 3, 5 and 10 — a pointer to RULES.md, Hard rules and How we work; rule 6 — without the list of what `docs-check` checks, which is in CLAUDE.md, Environment; the block «Ссылки» — gone, as RULES.md, Citations says it;
  - «Как устроен план: эпики» — the same rules, shorter; the list of done epics — one line each;
  - the opening paragraphs of tables 2 and 3 — the history of epics proposed and of the closed rows В1–В14 — and the epic's «Почему такой порядок» and «Решено автором (2026-11-11)» — to [docs/plan/E5.md](plan/E5.md), «Из RUN.md», with the phrases of T50–T53. Table 3 keeps its next number, В15, which `docs-check` reads.

— *Why:* every session reads RUN.md whole, at its start, and half of it was history or a copy of the shared rules that CLAUDE.md already loads. The template of `spec` keeps the section «Задачи вне эпика» (R19 spec, R145 spec), but here it held only the finished T50–T53, whose bodies were in the archive already; a task outside an epic loses nothing by standing in table 1. How this repository differs from the template is recorded here (the opening of this log).

— *Rejected:*
- keep the section empty for the analysis task of R145 spec — such a task can stand in table 1 like any other;
- delete the history without moving it — the archive of the epic is where it would have gone at the epic's end, and the decisions cite it;
- renumber the rules — the log and CLAUDE.md cite them by number.

— *Consequence:* the 2026-11-13 amendments to R44 and R45. A move to a newer tag of `spec` that brings a new RUN.md of the template reads this entry first. Where it is written: RUN.md; [docs/plan/E5.md](plan/E5.md).
