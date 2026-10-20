# sim

The NovaTrail learner simulator: thousands of virtual learners of different profiles answer tasks through learner-model and engine, so that a change to either is measured before it is merged. It holds the criterion a new rule is chosen by, written down before the first candidate runs, and the tolerance bands that CI holds every change to; it measures the algorithms and never changes them. It reads the module packages as pinned data, through spec's types and without module-kit, and holds no data of live learners.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.2.9), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, learner-model and engine; it comes together with learner-model and engine, since the simulator needs both the model and the rule of choice.

## Status

Detailing. The documents come first: the first epic of [RUN.md](RUN.md) decides what the component's card left open, and the code follows in the order RUN.md sets. There is no simulator here yet: its virtual learners, its criterion and its bands are what the first epic decides, along with what it recommends to spec on how the libraries' CI takes it, an open question of spec's.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, and what is still to decide.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.2.9), checked out here as a git submodule at `v0.2.9`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).
