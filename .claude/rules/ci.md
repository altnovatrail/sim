---
paths:
  - ".github/workflows/**"
  - ".github/dependabot.yml"
  - ".devcontainer/**"
  - "justfile"
---

# CI, the release and the devcontainer

- **CI**, `.github/workflows/ci.yml`, calls `.github`'s shared workflows by the SHA of a release: secrets, DCO, the documents and `go.yml` with `same-bits`, which runs the `ci-` recipes of the [justfile](../../justfile) in the `toolchain` stage of `.devcontainer/Dockerfile` (R30 spec; R10). Its "Gate" is the workflow `.github/workflows/gate.yml`, called on `ubuntu-24.04` and `ubuntu-24.04-arm` with mode `exact`: four jobs of `just gate-part <k> 4` and the check `just gate-check`, named "Gate / amd64" and "Gate / arm64", which the ruleset of `main` requires, in the same stage, under `timeout-minutes: 15` (R26, R31, R67). It takes the stage's image from the cache by three steps copied from the pinned `go.yml`, the same file, image, path and key, so that one entry serves both, and a push to `main` alone saves it (R63, R68): Dependabot moves the SHA, not the copy, so a pull request that moves `.github` holds the steps to its `go.yml` by hand, and the recipes `test-targets`, `fma-diff` and `rules` to the jobs of its `go.yml` whose scripts they run. Its Go cache is keyed by `go.sum` and the Dockerfile, never `gate.json`, and a pull request does not save it (R67). A caller takes `gate.yml` by a tag: a change of its inputs is a change for the callers. CodeQL and Scorecard are workflows of their own.
- **The release**, `.github/workflows/release.yml`, is `.github`'s shared one: a commit of `main` with CI green that changes what the gate is built from gives the next tag of line `0.1`, with `sim-<version>.tar.gz`, which `just release-files` writes (R29). The callers take sim by the tag and the SHA of its commit.
- **The devcontainer** is spec's, with gopls of the template's kind `go` (R32): every tool at an exact version, and an archive checked against its sha256 (RULES.md, Hard rules). `GO_VERSION` and its sums, `GOPLS_VERSION` and `.devcontainer/gopls/` come from spec with a move to its tag (R102 spec; the 2026-11-11 amendment to R64 spec), and `just ci-lint` holds the lock to the `ARG`. The directive `go` of `go.mod` names the least version, not the container's.
- **Dependabot** raises the actions and `.github`'s workflows; the tools of the justfile, `GOVULNCHECK` and `GO_LICENSES`, are raised by hand.

A task that changes CI or the justfile ends with the recipes it touches green in the container, and its pull request's CI green; one that changes the release runs `just release-files` too.
