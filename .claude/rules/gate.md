---
paths:
  - "gate.json"
  - "gate.go"
  - "gate_test.go"
  - "internal/gate/**"
  - "internal/criterion/**"
  - "judgements/**"
  - "evaluations/**"
---

# The gate, the paths and the judgements

[gate.json](../../gate.json) names the commits of learner-model and engine the bands were recorded on, the pin of modules, spec's version and the bands themselves. [internal/criterion/paths.json](../../internal/criterion/paths.json), embedded in the simulator, holds the paths the criteria score by, read from the decisive runs of the base and of the ceiling, whose SHA-256 it names. How the gate runs and a candidate is judged is in [docs/README.md](../../docs/README.md) and [docs/architecture.md](../../docs/architecture.md).

- `just gate exact` — every number of the gate the same to the bit as the bands recorded: a change that should move no number ends with it green. `just gate` holds the numbers to their bands, as the callers' "Gate" does. `just gate-part <k> <n> <folder>` and `just gate-check <mode> <folder>` are the same gate as jobs, as `gate.yml` runs it in CI.
- `just gate record "<reason>"` records the bands anew, with the reason, when a change moves the base on purpose.
- `just decisive <rule> <file> [<set>]` — a decisive run, 2,000 learners a cell; a candidate's by `CANDIDATE_ENGINE` or `CANDIDATE_LEARNER_MODEL` at its commit. `just rough <rule> <file> [<learners>]` — a rough look, the same run on fewer working learners, 300 unless named; it decides nothing. The files go to `runs/`, which git ignores; a judgement's file, report and SHA-256 go to `judgements/<round>/`.
- A task that moves the base's numbers runs the decisive runs of the base and of the ceiling anew, writes the paths with `sim paths` — their table, `sim paths -table`, as a new subsection of the paths of the targets in [docs/measurements.md](../../docs/measurements.md), which the test of the paths reads, and the log's entry of the two criteria rewritten in place with the new values, citing it — and the bands with `just gate record`. While the judgement is deferred, a move of the base records the bands alone and leaves the paths and the base's decisive run that `paths.json` names: the judgement that returns reads the paths on its base.
- `sim judge` holds the base and the candidates to the libraries of the paths and of the base, and a confirmation to the held-out set that `internal/seed/sets.go` records for the judgement its file names by `"judgement"`; a judgement is recorded as a decision of the log.
- The time of a decisive run is written with the host's load average; a measurement that chooses a size runs on a free machine, since a busy one skews the time the size is chosen by.
- An evaluation of fairness by the grade is `evaluations/fairness-<n>/`: its `evaluation.json` names the decisive run of the base by its SHA-256, and `sim evaluate` writes `report.md`, which names the version of sim from its build — so both are made from a clean clone at a tag. An evaluation of the learner model's three rules — the start a band below the grade's, the adult's correction and the guard —, `evaluations/rules-<n>/`, is the same by the criteria the log wrote down for them, made at the commit of its task, since the base it ran on has no tag.
