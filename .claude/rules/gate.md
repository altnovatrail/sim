---
paths:
  - "gate.json"
  - "gate.go"
  - "gate_test.go"
  - "internal/gate/**"
  - "internal/criterion/**"
  - "judgements/**"
---

# The gate, the paths and the judgements

[gate.json](../../gate.json) names the commits of learner-model and engine the bands were recorded on, the pin of modules, spec's version and the bands themselves (R05, R26). [internal/criterion/paths.json](../../internal/criterion/paths.json), embedded in the simulator, holds the paths the criteria score by, read from the decisive runs of the base and of the ceiling, whose SHA-256 it names (R28). How they are used is [README.md](../../README.md), "Gate" and "Judging a candidate", and [docs/architecture.md](../../docs/architecture.md), sections 3.6–3.9.

- `just gate exact` — every number of the gate the same to the bit as the bands recorded: a change that should move no number ends with it green (R26). `just gate` holds the numbers to their bands, as the callers' "Gate" does.
- `just gate record "<reason>"` records the bands anew, with the reason, when a change moves the base on purpose (R26).
- `just decisive <rule> <file> [<set>]` — a decisive run, 2,000 learners a cell (R33); a candidate's by `CANDIDATE_ENGINE` or `CANDIDATE_LEARNER_MODEL` at its commit (R34). The files go to `runs/`, which git ignores; a judgement's file, report and SHA-256 go to `judgements/<round>/` (R35).
- A task that moves the base's numbers runs the decisive runs of the base and of the ceiling anew, writes the paths with `sim paths` — new values by an amendment to R04 that cites their table, `sim paths -table`, as a new subsection of section 3 of [docs/measurements.md](../../docs/measurements.md), which the test of the paths reads (the 2026-11-13 amendment to R28) — and the bands with `just gate record` (section 3.7 of the architecture).
- `sim judge` holds the base and the candidates to the libraries of the paths and of the base (R40), and a judgement is recorded as a decision of the log, as R36 and R37.
- The time of a decisive run is written with the host's load average; a measurement that chooses a size runs on a free machine (RUN.md, epic E5).
