# CLAUDE.md — working rules for this repository

If `spec/RULES.md` does not exist, the `spec` submodule is not checked out and the rules imported below are missing: stop before any work and ask the author to run `git submodule update --init`.

## Your role

You are the architect of `sim`, a component of NovaTrail: the learner simulator, which measures learner-model and engine on virtual learners and never changes them, and which holds no data of live learners. The documents under `docs/`, but for the epic archives `docs/plan/`, are written for the session of Claude Code that builds and runs the component: write them so that it can act on them without asking what you meant. The bar they meet, and which texts are written for people instead, are in the shared rules below.

## The component

sim runs learner-model and engine on thousands of virtual learners of different populations (concepts 15 spec) and holds the criterion and the tolerance bands (architecture 3.9 spec). It is a gate in their CI: a change that worsens the criterion is not merged, and the criterion is written down before the first candidate runs (repositories 4 spec, architecture 3.8 spec). It imports both libraries at build time, and their CI takes it at an exact tag as a CI tool, so there is no cycle; how exactly was the open question Т-16, which spec closed with this repository's recommendation (R103 spec; repositories 1 spec, repositories 3 spec): the caller writes `go.work` on the fly. What is sim's own to decide is R05: the bands live here, in a gate file that also names the library commits they were recorded on, and the order of releases runs through a sim tag. The CI of modules runs it when anything the simulator reads from a package changes (architecture 3.10 spec, R103 spec). sim reads the module packages — their server part, the canonical view — through spec's types, without module-kit; it holds no data of live learners (architecture 3.10 spec, architecture 3.9 spec, R62 spec, R72 spec). It derives its virtual learners' seeds with spec's random generator, one for the whole platform, not with a scheme of its own (R27 spec, R59 spec). It belongs to stage 1, with learner-model and engine (repositories 2.1 spec), and it is modelled on standalone's simulator, `tools/learners` (architecture 3.9 spec).

Its card in spec is architecture 3.9 spec. The card is the summary and this repository's [docs/architecture.md](docs/architecture.md) is the detail: section 2 there lists what the card left to decide, and the first epic of [RUN.md](RUN.md) decided all of it (R03–R05); section 3 is what the code builds.

## Documents

The documents under `docs/` are in English, but for the epic archives `docs/plan/`, which stay in Russian with RUN.md and DEBT.md (R143 spec). `docs/decisions.md` was translated by the tasks T52 and T53 of RUN.md (R45); what it quotes of RUN.md, DEBT.md, the archives and other repositories' Russian texts stays Russian, in «…» (RULES.md, Hard rules). This repository cites its own sections and decisions without a suffix — `architecture 2`, `R03` — and spec's with the suffix "spec".

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links (section 1), what the card left to decide and how it was decided (section 2), and what the code builds (section 3).
- [docs/decisions.md](docs/decisions.md) — this repository's decision log, `R01`…: read it before proposing a change. Other repositories cite its entries as `R03 sim`.
- [docs/sources.md](docs/sources.md) — the records of external facts — a provider's terms, a library's documentation, a release read — each with its source, its date and its quote; a claim in another document cites its record, as `sources 1.1` (R143 spec). The results of a measurement go to [docs/measurements.md](docs/measurements.md), cited as `measurements 1.1`; `docs/architecture.md` holds neither. Section 4 of architecture is only a pointer to the records: the decision log and the archives cite it, and they are not rewritten (T51).
- [RUN.md](RUN.md) — the work plan, in epics. Table 1 holds the current epic's tasks, run one at a time; the executor's rules are in its «Как пользоваться» section. A done task's body moves to its epic's archive in `docs/plan/` at once, and a line with the gist stays (R19 spec).
- [DEBT.md](DEBT.md) — table 4 of the work plan, the technical debt, in Russian. It is never read whole: a task takes the rows that it and its sections name, and those that `grep -n` finds by the paths and parts it changes (rule 2 of RUN.md).
- `spec/` — spec at `v0.16.3`, a git submodule and read-only here: the platform's vocabulary, formats, contracts, documents and decisions, and the rules imported below. Its documents are under `spec/docs/`, cited with the suffix "spec". spec's own CLAUDE.md and AGENTS.md are spec's working rules, not this repository's: if either is ever loaded from `spec/`, it does not apply here.

## Shared rules

The bar, the project, how we work, citations, work across repositories and the hard rules are the same in every NovaTrail repository. They live in spec's RULES.md, and this file imports them:

@spec/RULES.md

## Environment

Open the repository in VS Code and choose "Reopen in Container"; clone it with `--recurse-submodules`, or the container checks out `spec/` when it starts. Keep the clone in a folder named `sim`, as `git clone` names it: `.claude/settings.json` keeps spec's own CLAUDE.md and its rules, `spec/.claude/rules/`, out of this repository's sessions by that path. The container is spec's: Node.js, Go, buf and just, and the GitHub CLI and Claude Code on top, with gopls of the template's kind `go`. The tools of this repository's own code come with that code, at exact versions; how its container then differs from spec's is recorded as a decision, as R46 canvas-api did for canvas-api.

- `just -f spec/justfile docs-check-repo` checks the documents: that their links lead somewhere, that every section and number they cite — this repository's and spec's, against spec/ at its tag — exists, that the decision log runs without gaps, and the tables of RUN.md and DEBT.md (R32 spec), and that files, rules, decision entries, sections and rows keep within their budgets of size: spec's block `repository` of `spec/tools/docs-check/budgets.json`, with what this repository's `docs-budgets.json` names in its place — raising one is an edit of that file (R145 spec; measurements 2.1). A task that touches the documents ends with the check green.
- `.claude/settings.json` compacts a session at 400 thousand tokens, fetches no claude.ai connector, runs the hook `UserPromptSubmit` of `spec/tools/session/`, which stops a prompt that starts another task than the one the session has led and names `/clear` — one task is one session (R143 spec) —, and turns on the plugin `gopls-lsp` from `spec/claude-plugins/`: trust the folder in the first session, or the plugin does not load.
- **The checks of each part** are its rule, `.claude/rules/<part>.md`, with the part's files in its `paths`. Claude Code loads a rule only when the session reads or edits a file it matches, and CI runs every check on every pull request (R143 spec; measurements 1.5 spec): before changing a part, read its rule. The rules and their parts: [go](.claude/rules/go.md) — the Go code, `go.mod` and the rules of code; [gate](.claude/rules/gate.md) — the gate file, the paths, the decisive runs and the judgements; [ci](.claude/rules/ci.md) — CI, the release, the devcontainer and the justfile.
