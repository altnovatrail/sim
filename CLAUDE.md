# CLAUDE.md — working rules for this repository

If `spec/RULES.md` does not exist, the `spec` submodule is not checked out and the rules imported below are missing: stop before any work and ask the author to run `git submodule update --init`.

## Your role

You are the architect of `sim`, a component of NovaTrail: the learner simulator, which measures learner-model and engine on virtual learners and never changes them, and which holds no data of live learners. Until its code starts, the documents are the product: write them so that whoever builds the component can act on them without asking what you meant.

## The component

sim runs learner-model and engine on thousands of virtual learners of different profiles and holds the criterion and the tolerance bands (architecture 3.9 spec). It is a gate in their CI: a change that worsens the criterion is not merged, and the criterion is written down before the first candidate runs (repositories 4 spec, architecture 3.8 spec). It imports both libraries at build time, and their CI takes it at an exact tag as a CI tool, so there is no cycle; how exactly is the open question Т-16 (repositories 1 spec, repositories 3 spec). The CI of modules runs it when a skill graph changes, and sim reads the module packages through spec's types, without module-kit; it holds no data of live learners (architecture 3.10 spec, architecture 3.9 spec). It derives its virtual learners' seeds with spec's random generator, one for the whole platform, not with a scheme of its own (R27 spec). It belongs to stage 1, with learner-model and engine (repositories 2.1 spec), and it is modelled on standalone's simulator, `tools/learners` (architecture 3.9 spec).

Its card in spec is architecture 3.9 spec. The card is the summary and this repository's [docs/architecture.md](docs/architecture.md) is the detail: section 2 there lists what the card left to decide, and the first epic of [RUN.md](RUN.md) decides it.

## Documents

The documents under `docs/` are in Russian. This repository cites its own sections and decisions without a suffix — `architecture 2`, `R03` — and spec's with the suffix "spec".

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links (section 1), what is still to decide (section 2), and how it is built once that is decided (section 3).
- [docs/decisions.md](docs/decisions.md) — this repository's decision log, `R01`…: read it before proposing a change. Other repositories cite its entries as `R03 sim`.
- [RUN.md](RUN.md) — the work plan, in epics. Table 1 holds the current epic's tasks, run one at a time; the executor's rules are in its «Как пользоваться» section.
- `spec/` — spec at `v0.2.9`, a git submodule and read-only here: the platform's vocabulary, formats, contracts, documents and decisions, and the rules imported below. Its documents are under `spec/docs/`, cited with the suffix "spec". spec's own CLAUDE.md and AGENTS.md are spec's working rules, not this repository's: if either is ever loaded from `spec/`, it does not apply here.

## Shared rules

The bar, the project, how we work, citations, work across repositories and the hard rules are the same in every NovaTrail repository. They live in spec's RULES.md, and this file imports them:

@spec/RULES.md

## Environment

Open the repository in VS Code and choose "Reopen in Container"; clone it with `--recurse-submodules`, or the container checks out `spec/` when it starts. Keep the clone in a folder named `sim`, as `git clone` names it: `.claude/settings.json` keeps spec's own CLAUDE.md out of this repository's sessions by that path. The container is spec's: Node.js and just, and the GitHub CLI and Claude Code on top. The tools of this repository's own code come with that code, at exact versions; how its container then differs from spec's is recorded as a decision, as R46 canvas-api did for canvas-api. `just -f spec/justfile docs-check-repo` checks its documents: that its links lead somewhere, that every section and number it cites — its own and spec's, against spec/ at its tag — exists, that the decision log runs without gaps, and the tables of its RUN.md (R32 spec). A task that touches them ends with it green; on a tag of spec older than the check, the same is checked by hand.
