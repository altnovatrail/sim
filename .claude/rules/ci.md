---
paths:
  - ".github/workflows/**"
  - ".github/dependabot.yml"
  - ".devcontainer/**"
  - "justfile"
---

# CI, the release and the devcontainer

- **CI**, `.github/workflows/ci.yml`, calls `.github`'s shared workflows by the SHA of a release: secrets, DCO, the documents and `go.yml` with `same-bits`, which runs the `ci-` recipes of the [justfile](../../justfile) in the `toolchain` stage of `.devcontainer/Dockerfile` (R30 spec; R10). Its own job "Gate" runs `just gate exact` on `ubuntu-24.04` and `ubuntu-24.04-arm`, in the same stage, under `timeout-minutes: 15` (R26, R31). CodeQL and Scorecard are workflows of their own.
- **The release**, `.github/workflows/release.yml`, is `.github`'s shared one: a commit of `main` with CI green that changes what the gate is built from gives the next tag of line `0.1`, with `sim-<version>.tar.gz`, which `just release-files` writes (R29). The callers take sim by the tag and the SHA of its commit.
- **The devcontainer** is spec's, with gopls of the template's kind `go` (R32): every tool at an exact version, and an archive checked against its sha256 (RULES.md, Hard rules). `GO_VERSION` and its sums, `GOPLS_VERSION` and `.devcontainer/gopls/` come from spec with a move to its tag (R102 spec; the 2026-11-11 amendment to R64 spec), and `just ci-lint` holds the lock to the `ARG`. The directive `go` of `go.mod` names the least version, not the container's.
- **Dependabot** raises the actions and `.github`'s workflows; the tools of the justfile, `GOVULNCHECK` and `GO_LICENSES`, are raised by hand.

A task that changes CI or the justfile ends with the recipes it touches green in the container, and its pull request's CI green; one that changes the release runs `just release-files` too.
