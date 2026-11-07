# sim

[![CI](https://github.com/altnovatrail/sim/actions/workflows/ci.yml/badge.svg)](https://github.com/altnovatrail/sim/actions/workflows/ci.yml)
[![CodeQL](https://github.com/altnovatrail/sim/actions/workflows/codeql.yml/badge.svg)](https://github.com/altnovatrail/sim/actions/workflows/codeql.yml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/altnovatrail/sim/badge)](https://scorecard.dev/viewer/?uri=github.com/altnovatrail/sim)
[![License](https://img.shields.io/github/license/altnovatrail/sim)](LICENSE)

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)

The NovaTrail learner simulator: thousands of virtual learners of different populations answer tasks through learner-model and engine, so that a change to either is measured before it is merged. It holds the criterion a new rule is chosen by, written down before the first candidate runs, and the tolerance bands that CI holds every change to; it measures the algorithms and never changes them. It reads the module packages as pinned data, through spec's types and without module-kit, and holds no data of live learners.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.9.5), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, learner-model and engine; it comes together with learner-model and engine, since the simulator needs both the model and the rule of choice.

## Status

Detailing is done; the code comes next. The first epic of [RUN.md](RUN.md) decided what the component's card left open: the virtual learners, the criteria with their bands, and how the CI of learner-model, engine and modules takes the gate. spec took its recommendation on that last point and closed its open question with it (R103 spec). There is no simulator here yet: [docs/architecture.md](docs/architecture.md), section 3, is what the next epic builds.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, what was decided, and what the code builds.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.9.5), checked out here as a git submodule at `v0.9.5`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).
