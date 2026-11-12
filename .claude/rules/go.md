---
paths:
  - "*.go"
  - "internal/**"
  - "go.mod"
  - "go.sum"
  - ".github/coderules.json"
  - ".github/codecov.yml"
---

# The Go code of sim

sim is one Go module, `github.com/altnovatrail/sim`: the root is `package main`, the program `sim`, and the parts of the bench are packages under `internal/` (R10; [docs/architecture.md](../../docs/architecture.md), section 3). `go.mod` requires spec's Go module at the tag of the `spec/` submodule (R64 spec) and the two libraries at the versions the base stands on, which `gate.json` names too ([gate](gate.md)).

A task that touches Go code ends with these green:

- `just ci-lint` — gofmt, `go vet`, the gopls lock against `GOPLS_VERSION`, and the release's file (R29);
- `just ci-test` — every package, with coverage in `coverage/go.out`, and no race detector: the code has no goroutine, and the rule `goroutine` keeps it out (the 2026-11-09 amendment to R10);
- `just rules` — the rules of code of the pinned `.github`, with this repository's lists, `.github/coderules.json`;
- `just same-bits` when the change touches arithmetic — `test-targets`, `fma-diff` and `rules`: the same bits on every build, as "Go / targets" and "Go / fma-diff" hold them (R09);
- `just ci-vuln` and, with `ALLOWED_LICENSES` set as `.github`'s `go.yml` sets it, `just ci-licenses` when a dependency changes.

Code comes with its tests, and the merge waits for `codecov/patch`, 90 % of the changed lines (R48 spec). A file of a NovaTrail release enters the code only once `just -f spec/justfile release-verify` has accepted it (R77 spec). A change that moves a number of the base is not code alone: the [gate](gate.md) rule says what it runs.
