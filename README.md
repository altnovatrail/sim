# sim

[![CI](https://github.com/altnovatrail/sim/actions/workflows/ci.yml/badge.svg)](https://github.com/altnovatrail/sim/actions/workflows/ci.yml)
[![CodeQL](https://github.com/altnovatrail/sim/actions/workflows/codeql.yml/badge.svg)](https://github.com/altnovatrail/sim/actions/workflows/codeql.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/altnovatrail/sim)](go.mod)
[![Coverage](https://codecov.io/gh/altnovatrail/sim/branch/main/graph/badge.svg)](https://codecov.io/gh/altnovatrail/sim)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/altnovatrail/sim/badge)](https://scorecard.dev/viewer/?uri=github.com/altnovatrail/sim)
[![License](https://img.shields.io/github/license/altnovatrail/sim)](LICENSE)

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_sim&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=altnovatrail_sim)

The NovaTrail learner simulator: thousands of virtual learners of different populations answer tasks through learner-model and engine, so that a change to either is measured before it is merged. It holds the criterion a new rule is chosen by, written down before the first candidate runs, and the tolerance bands that CI holds every change to; it measures the algorithms and never changes them. It reads the module packages as pinned data, through spec's types and without module-kit, and holds no data of live learners.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.15.0), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, learner-model and engine; it comes together with learner-model and engine, since the simulator needs both the model and the rule of choice.

## Status

Detailing is done, and so is the code epic, E2 ([RUN.md](RUN.md)). The first epic decided what the component's card left open: the virtual learners, the criteria with their bands, and how the CI of learner-model, engine and modules takes the gate. spec took its recommendation on that last point and closed its open question with it (R103 spec). E2 built the bench on the first versions of learner-model, engine and modules. The decisive runs of the base and of the ceiling of the engine fixed the paths the criteria score by, and the first tag, `v0.1.0`, came out on 2026-11-09: it is the gate the CI of those three repositories takes. Its own CI is .github's shared checks, with the same bits held on every build, the gate on amd64 and arm64, CodeQL for Go and the coverage in Codecov. The bench runs: `sim run`, `sim merge`, `sim compare`, `sim judge`, `sim paths` and `sim gate`; [docs/architecture.md](docs/architecture.md), section 3, is what it is built of. The bench is pinned to the packages learners get, modules v0.2.5 with both maths at 0.3.0, with checks of time and leading questions, and every module in format 0.4, on spec v0.15.0, learner-model v0.1.6 and engine v0.1.6; the gate takes 250 learners a cell, and a decisive run 2,000 and at most an hour, keeping its file of results in `runs/`, which git ignores (R31, R33 in [docs/decisions.md](docs/decisions.md)). E3 judged the first candidates of both libraries, and both judgements ended in an exit, so the base stays (R36, R37; [docs/plan/E3.md](docs/plan/E3.md)). E4 kept the bench up with the libraries: spec v0.15.0 and engine v0.1.6 with the sizes of runs after its faster choice, the check kind `time` and the leading question of format 0.4 on learner-model v0.1.5 (R41), and the pin of modules that has them, with the paths read anew (R42). The base then moved to learner-model v0.1.6, whose faster readout gives the same bits, with the paths read anew (R43). The current epic, E5, has raised the sizes of runs by a measurement on it — the gate from 150 learners to 250, a decisive run from 1,500 to 2,000, from the tag after v0.1.9 — and is to judge the second round of the model's candidates on them; the engine's second round is the epic after it, E6 (R44).

## Gate

`just gate` builds the bench on the commits of the libraries that [gate.json](gate.json) names, through a `go.work` it writes in a temporary folder, runs the gate and holds its numbers to their bands; `just gate exact` holds them to the recorded numbers to the last bit, as this repository's CI does on amd64 and arm64, and `just gate record "<reason>"` records the bands anew. A caller puts its own code in place of a library with `GATE_LEARNER_MODEL` or `GATE_ENGINE`: a folder, or a commit by its full SHA.

### Gate in the CI of modules

The CI of modules judges content, not a library, so its gate holds invariants and thresholds rather than bands, and runs only where what the bench reads of a package changed ([docs/architecture.md](docs/architecture.md), section 3.9). Its job checks out sim by the tag and the SHA of its commit and runs, from the root of that checkout:

```
just gate modules (-index <module-index.json> | -first) [-published <folder>]... [-all] <package folder>...
```

- `-index` is `module-index.json` of the newest release of modules, which `release-verify` has accepted; `-first` says there is no release yet, and every package is new. Give exactly one of the two.
- Each `-published` is a published package of that index as a folder, with its `inventory.json` and `server.json`: the packages of every version of the modules, withdrawn ones too, as the CI of modules already unpacks them for `module-kit diff`. sim finds each by the digest of its `inventory.json` and checks it by the chain of digests from the index.
- Each `<package folder>` is a package that `module-kit build` wrote for the pull request. Give every module the pull request builds, not only the changed ones: a cell on two modules runs on both of them.
- `-all` runs every package: give it when the pull request moves the tag of sim.
- Relative paths are taken from the folder `just` was called in.

The recipe first asks `sim gate changed` what each package is, and prints a line a package:

| State | Meaning | Runs |
|---|---|---|
| `new` | the module has no published version that is not withdrawn | yes |
| `changed` | what the bench reads differs from the newest published version that is not withdrawn | yes |
| `same` | what the bench reads is the same: only texts for the learner, or the version, changed | with `-all`, or when a cell of the gate runs it with a changed module: A-linked runs on both maths, so a change to either brings the other along, and its line names the module that brought it |
| `transition` | the newest published version is of another version of the module format | never |

What the bench reads is what it translates the server part into: the bands and the audience of the manifest, the graph and its tiers, the tasks, the lessons and the instances of templates, and none of the texts for the learner. If nothing runs, the recipe says so and exits with 0 before it fetches anything. Otherwise it builds the bench on both libraries at the commits that gate.json names, runs the gate's populations on the graphs of the packages that run, and checks the numbers against the thresholds:

- `E2`: the share of learners stuck, and of the choices of nothing while the pool has a task to give, on every cell. Each is at most its largest value in gate.json plus one percentage point.
- `E1` on `G0` and `G2-half`. The share of items with P* below 0.50 is at most its largest value in gate.json plus one point. The share in the corridor is at least its smallest value less one point.
- `M4` on `G0`: at most 20 % false "mastered".
- `E5` on `E-base`: no task with a rubric given again.

The recipe exits with 0 when the gate is green, 1 when it is red or refuses its input, and 2 on a call it does not know. It writes a table of every number against its threshold to stdout, for the job's summary. Run it in the toolchain stage of sim's devcontainer at the same tag, as this repository's job "Gate" does: `docker build --target toolchain - < .devcontainer/Dockerfile`, then `docker run` with both the checkout of sim and the package folders mounted, so that every path the recipe is given is inside the container.

## Judging a candidate

A candidate of learner-model or engine is judged by a criterion on decisive runs of 2,000 learners a cell ([docs/architecture.md](docs/architecture.md), sections 3.6 and 3.8; R34 in [docs/decisions.md](docs/decisions.md)). From the root of this repository:

1. **The base.** `just decisive engine runs/base.jsonl` builds the bench on the libraries go.mod names and runs the base on the working learners. Its file is the one [paths.json](internal/criterion/paths.json) names: `sha256sum runs/base.jsonl` gives its `base`. A file kept in `runs/` from an earlier run needs no run again.
2. **The candidate.** `CANDIDATE_ENGINE=<full SHA> just decisive engine runs/engine-a.jsonl` builds the bench on that commit, through the same workspace as the gate, with the other library at the commit of gate.json. The SHA may be the head of a pull request that is on no branch. A folder works in place of the SHA, if it has nothing uncommitted. `CANDIDATE_LEARNER_MODEL` names a candidate of learner-model, and both together with `CANDIDATE_PAIR=yes` a pair. A relative folder is of the caller's folder. The recipe holds the build to gate.json before it runs, and the file names both libraries by commit.
3. **The judgement.** A judgement file names the criterion — `model`, `engine` or `pair` — the base's file and each candidate's, with its name, its group and what makes it simple:

   ```json
   {
     "criterion": "engine",
     "base": "runs/base.jsonl",
     "candidates": [
       {"name": "A", "file": "runs/engine-a.jsonl", "group": "<group>", "constants": 1}
     ]
   }
   ```

   `sim judge <judgement file>` writes the report: the constraints, the score and the choice. It refuses runs that are not of the judgement: a base on other libraries than those the paths were read on, or than the embedded gate.json names, and a candidate of one library on another commit of the other (R40).
4. **The confirmation.** The candidate chosen and the base run once more, on the held-out learners of the choice: `just decisive engine runs/base-held-out.jsonl held-out-engine-1`, and the same with `CANDIDATE_ENGINE`. The judgement file then names them under `"confirm": {"base": …, "candidate": …}`.
5. **The bands.** A candidate confirmed goes into gate.json and go.mod, and `just gate record "<reason>"` records the bands anew in a pull request to this repository.

## Taking a release

A release comes out by itself once CI passes on main and something the gate is built from has changed: the next patch of the line `0.1`, from `v0.1.0`. Releases are immutable, and the tags `v*` are under a ruleset, so a tag stays on its commit for good.

A caller of the gate takes sim from git, not as a file:

- write both the tag and the SHA of its commit in its workflow: the tag for the reader, and the SHA for what holds the content, as actions are pinned;
- check out sim at that SHA, then run `just gate` or `just gate modules` from it, in the toolchain stage of its devcontainer (see [Gate](#gate)).

A release's one file is `sim-<version>.tar.gz`, the archive of what builds the bench and runs the gate, signed without a key. `just -f spec/justfile release-verify sim <tag> <file>` checks it, and `gh release verify <tag> -R altnovatrail/sim` checks the release.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, what was decided, and what the code builds.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.15.0), checked out here as a git submodule at `v0.15.0`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).
