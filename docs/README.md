# The documents of sim

What the [README](../README.md) of the repository leaves out, which says only what the component is and where it stands in NovaTrail (the second 2026-11-17 amendment to R21 spec; R77 in [decisions.md](decisions.md)). How to run the gate, judge a candidate and take a release are the sections below, moved from the README as they were.

- [architecture.md](architecture.md) — the component: what it does and does not do, its data and links, what was decided, and what the code builds.
- [decisions.md](decisions.md) — the decision log: what was decided, why, and what was rejected.
- [sources.md](sources.md) — the records of external facts the documents rest on: each source, the date it was read and its quote.
- [measurements.md](measurements.md) — the results of measurements: the times of the gate and of decisive runs.
- [RUN.md](../RUN.md) — the work plan, in epics.
- `run/` — the plan's other tables: [the epics for later](../run/EPICS.md), [the open questions](../run/QUESTIONS.md) and [the technical debt](../run/DEBT.md).

The working rules, for people and coding agents alike, are in [CLAUDE.md](../CLAUDE.md) and [AGENTS.md](../AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec), checked out here as a git submodule: clone this repository with `git clone --recurse-submodules`.

## Gate

`just gate` builds the simulator on the commits of the libraries that [gate.json](../gate.json) names, through a `go.work` it writes in a temporary folder, runs the gate and holds its numbers to their bands; `just gate exact` holds them to the recorded numbers to the last bit, as this repository's CI does on amd64 and arm64, and `just gate record "<reason>"` records the bands anew. A caller puts its own code in place of a library with `GATE_LEARNER_MODEL` or `GATE_ENGINE`: a folder, or a commit by its full SHA.

The same gate runs as jobs (R67 in [decisions.md](decisions.md)): `just gate-part <k> <n> <folder>` runs part k of n of the learners, as one job of n, into the folder, and `just gate-check <bands|exact> <folder>` holds the build and the merged parts as `just gate` does, refusing a part missing, repeated or of another run. The jobs of one gate run on one kind of machine, since each splits its part by its processors. [.github/workflows/gate.yml](../.github/workflows/gate.yml) is that, four jobs and the check, and this repository's CI calls it on amd64 and arm64. A caller calls it by the SHA of a tag of sim:

```yaml
  gate:
    name: Gate
    uses: altnovatrail/sim/.github/workflows/gate.yml@<SHA of the tag> # <tag>
    with:
      name: amd64
      sim-tag: <tag>
      sim-commit: <SHA of the tag>
      library: GATE_ENGINE
```

It checks out sim at `sim-commit`, which `sim-tag` must be on, and the caller's own commit in place of the library `library` names, and holds the numbers to their bands; the check is "Gate / amd64".

### Gate in the CI of modules

The CI of modules judges content, not a library, so its gate holds invariants and thresholds rather than bands, and runs only where what the simulator reads of a package changed ([architecture.md](architecture.md), section 3.9). Its job checks out sim by the tag and the SHA of its commit and runs, from the root of that checkout:

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
| `changed` | what the simulator reads differs from the newest published version that is not withdrawn | yes |
| `same` | what the simulator reads is the same: only texts for the learner, or the version, changed | with `-all`, or when a cell of the gate runs it with a changed module: A-linked runs on both maths, so a change to either brings the other along, and its line names the module that brought it |
| `transition` | the newest published version is of another version of the module format | never |

What the simulator reads is what it translates the server part into: the bands and the audience of the manifest, the graph and its tiers, the tasks, the lessons and the instances of templates, and none of the texts for the learner. If nothing runs, the recipe says so and exits with 0 before it fetches anything. Otherwise it builds the simulator on both libraries at the commits that gate.json names, runs the gate's populations on the graphs of the packages that run, and checks the numbers against the thresholds:

- `E2`: the share of learners stuck, and of the choices of nothing while the pool has a task to give, on every cell. Each is at most its largest value in gate.json plus one percentage point.
- `E1` on `G0` and `G2-half`. The share of items with P* below 0.50 is at most its largest value in gate.json plus one point. The share in the corridor is at least its smallest value less one point.
- `M4` on `G0`: at most 20 % false "mastered".
- `E5` on `E-base`: no task with a rubric given again.

The recipe exits with 0 when the gate is green, 1 when it is red or refuses its input, and 2 on a call it does not know. It writes a table of every number against its threshold to stdout, for the job's summary. Run it in the toolchain stage of sim's devcontainer at the same tag, as this repository's job "Gate" does: `docker build --target toolchain - < .devcontainer/Dockerfile`, then `docker run` with both the checkout of sim and the package folders mounted, so that every path the recipe is given is inside the container.

## Judging a candidate

A candidate of learner-model or engine is judged by a criterion on decisive runs of 2,000 learners a cell ([architecture.md](architecture.md), sections 3.6 and 3.8; R34 in [decisions.md](decisions.md)). From the root of this repository:

1. **The base.** `just decisive engine runs/base.jsonl` builds the simulator on the libraries go.mod names and runs the base on the working learners. Its file is the one [paths.json](../internal/criterion/paths.json) names: `sha256sum runs/base.jsonl` gives its `base`. A file kept in `runs/` from an earlier run needs no run again.
2. **The candidate.** `CANDIDATE_ENGINE=<full SHA> just decisive engine runs/engine-a.jsonl` builds the simulator on that commit, through the same workspace as the gate, with the other library at the commit of gate.json. The SHA may be the head of a pull request that is on no branch. A folder works in place of the SHA, if it has nothing uncommitted. `CANDIDATE_LEARNER_MODEL` names a candidate of learner-model, and both together with `CANDIDATE_PAIR=yes` a pair. A relative folder is of the caller's folder. The recipe holds the build to gate.json before it runs, and the file names both libraries by commit.
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

   `sim judge <judgement file>` writes the report: the constraints, the score and the choice. Every criterion holds the constraint `fair` of fairness by the grade: on `D-grade1` and `D-grade5` each number of `M11` against grade 3 no worse than the base's, so both runs need the cells of the populations `D-`, which `just decisive` and `just rough` run; the report names the cell and the number where it does not hold (R55, R57). It refuses runs that are not of the judgement: a base on other libraries than those the paths were read on, or than the embedded gate.json names, and a candidate of one library on another commit of the other (R40).

   A candidate may name `"equal_on"`: the graphs, as the IDs of their modules joined by `+`, on which it gives the base's runs. On each cell of them every learner's sums must be the base's, or the candidate fails the constraint `neutral`. The judgement of consolidation in a lesson names `["novatrail/grade-3-math", "novatrail/english-essays"]` (R49). Beside each candidate the report writes the note of consolidation in a lesson, `E6`, on the cells where a run consolidated; it decides nothing.
4. **The confirmation.** The candidate chosen and the base run once more, on the held-out learners of the choice, its judgement's own — `held-out-engine-lesson-1` for `engine-lesson-1`, `held-out-model-fair-1` for `model-fair-1`, `held-out-engine-1` for `engine-2` ([architecture.md](architecture.md), section 3.8): `just decisive engine runs/base-held-out.jsonl held-out-engine-lesson-1`, and the same with `CANDIDATE_ENGINE`. The judgement file then names them under `"confirm": {"base": …, "candidate": …}`, and names its judgement, as its folder: `"judgement": "engine-lesson-1"`; a file in a folder of `judgements/` that names another judgement is refused. `sim judge` takes the confirmation only on the set that [internal/seed/sets.go](../internal/seed/sets.go) records for that judgement, and refuses a set kept for another round or variant, or one another judgement has looked at (R52).
5. **The bands.** A candidate confirmed goes into gate.json and go.mod, and `just gate record "<reason>"` records the bands anew in a pull request to this repository.

**A rough look** goes the same way, on fewer learners a cell, before a candidate is ready for its decisive run: `just rough engine runs/base-300.jsonl` for the base and `CANDIDATE_ENGINE=<full SHA> just rough engine runs/engine-a-300.jsonl` for the candidate — 300 learners a cell unless a third argument names from 1 to 1,999, and only the working learners, since a look at held-out learners leaves them held out no longer. A judgement file of the two, without a confirmation, gives a report named a rough look, whose outcome decides nothing (R49).

**The judgement of the round of fairness**, `model-fair-1`, goes by the target its criterion wrote down before its first candidate ([architecture.md](architecture.md), section 3.6; R65, R66 in [decisions.md](decisions.md)). Its file names the judgement and the target, and its candidates name no group or goal:

```json
{
  "judgement": "model-fair-1",
  "criterion": "model",
  "target": "M11.same",
  "base": "../../runs/base.jsonl",
  "candidates": [
    {"name": "A", "file": "../../runs/model-fair-1-a.jsonl", "constants": 1}
  ]
}
```

The target is the candidate's `M11.same` less the base's, on `D-grade1` and `D-grade5` pooled, with its interval: it alone decides "better" and the choice, and the model's score stands in the report as a reference. The report writes `M11` of the base and of the candidate on both grades, and, for the confirmation on `held-out-model-fair-1`, whether the chosen one holds the bounds of R55; they decide nothing here. `sim judge` refuses `model-fair-1` without the target, and any other judgement with one. A rough look of the round, as learner-model's exploration runs it, uses the same file without a confirmation.

**The evaluation of fairness by the grade** reads the decisive run of the base, whose cells of the populations `D-grade1`, `D-grade3` and `D-grade5` the engine's ceiling does not run ([architecture.md](architecture.md), section 3.6; R55, R56 in [decisions.md](decisions.md)). An evaluation file in its folder, `evaluations/fairness-<n>/evaluation.json`, names the evaluation, the run relative to the file, its SHA-256 and the date:

```json
{"evaluation": "fairness-1", "run": "../../runs/base.jsonl", "sha256": "<sha256sum of the run>", "date": "YYYY-MM-DD"}
```

`sim evaluate evaluations/fairness-1/evaluation.json > evaluations/fairness-1/report.md` writes the report: each grade against grade 3, whether the base holds the bounds of R55, and the paragraph for spec's DPIA. It refuses a run of another SHA-256, a rough look, held-out learners, a run that did not record the invariants, and a run on other libraries than the gate file names. The report names the version of sim from its build, so the run and the report are made from a clean clone at a tag, as `fairness-1` was at v0.1.19 (R58). Since the gate stands on the commit of the three rules, `sim evaluate` refuses the run of `fairness-1`, which is on v0.1.8: its report stays the record of its tag (R74). A judgement reads those cells by the constraint of fairness alone.

**The evaluation of the three rules** of R157 spec reads a run of the populations `D-grade1`, `D-grade3`, `D-grade5` and `C-grade1`, 2,000 working learners a cell, on the commit of `learner-model` that R72 names, against the base's decisive run that `paths.json` names ([architecture.md](architecture.md), section 3.6; R72 in [decisions.md](decisions.md)). The run is `sim run -populations D-grade1,D-grade3,D-grade5,C-grade1 -learners 2000`, in parts with `-part k/n` and `sim merge`, and its evaluation file names its criterion and both runs:

```json
{"evaluation": "rules-1", "criterion": "rules", "run": "../../runs/rules.jsonl", "sha256": "<sha256sum of the run>", "base": "../../runs/base-t64.jsonl", "base_sha256": "<sha256sum of the base>", "date": "YYYY-MM-DD"}
```

`sim evaluate evaluations/rules-1/evaluation.json > evaluations/rules-1/report.md` writes the report: the four conditions of R72 with their numbers, M11 of grades 1 and 5 under the rules and the base, and the paragraph for spec's DPIA. The adult's correction against the truth, R73, is an evaluation file of the criterion `correction`, which names its run alone — `sim run -populations D-grade1,C-grade1 -learners 2000` — and its report is `evaluations/rules-2/report.md`.

## Taking a release

A release comes out by itself once CI passes on main and something the gate is built from has changed: the next patch of the line `0.1`, from `v0.1.0`. Releases are immutable, and the tags `v*` are under a ruleset, so a tag stays on its commit for good.

A caller of the gate takes sim from git, not as a file:

- write both the tag and the SHA of its commit in its workflow: the tag for the reader, and the SHA for what holds the content, as actions are pinned;
- check out sim at that SHA, then run `just gate` or `just gate modules` from it, in the toolchain stage of its devcontainer (see [Gate](#gate)).

A release's one file is `sim-<version>.tar.gz`, the archive of what builds the simulator and runs the gate, signed without a key. `just -f spec/justfile release-verify sim <tag> <file>` checks it, and `gh release verify <tag> -R altnovatrail/sim` checks the release.
