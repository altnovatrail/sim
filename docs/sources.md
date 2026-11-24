# sim — sources

> The facts of the outside world the repository's documents under `docs/` rest on and no run of ours checks: a law, a regulator, a standard's normative text, a provider's terms and policies, a price, a license, a page that changes in place with no version, a text outside the repositories. Here are their records — the address, the date it was read and a quote in the source's language — and [architecture.md](architecture.md) and the decision log hold the claim and a citation of its one record by its number, `sources 1.1`, not a chain. A record's heading names its topic, not a task.
>
> No record is needed for a library's documentation or code at a version or commit — a link beside the claim —, for a behaviour a test or a measurement checks — they are its evidence —, or for spec's tags and the files of NovaTrail's and MathTrail's repositories — git, at a commit.
>
> A record is a subsection "N.M" of a topic: one source, or one reading of several sources for one decision. A record's number never changes and is never given again: a new one goes at the end of its topic, a new topic at the end of the file, and a record removed leaves a gap. The next topic is 8: the numbers of topics 1 and 3–7 are not given again. A source read again with the same text changes the record's date; with another text, the claims and decisions on it are rewritten in the same change and the record holds the new reading alone — where a decision waits on the author, an open question, and both readings stand until it is closed. The model is `spec/docs/sources.md`. The results of measurements are not here but in [measurements.md](measurements.md).

## 2. GitHub and other external sources

### 2.3. GitHub Docs: runners and limits

Read 2026-11-07.

GitHub Docs: [GitHub-hosted runners](https://docs.github.com/en/actions/reference/runners/github-hosted-runners) — for public repositories `ubuntu-24.04` is Linux x64 and `ubuntu-24.04-arm` is arm64, each with 4 CPU and 16 GB, and "Use of the standard GitHub-hosted runners is free and unlimited on public repositories"; [Limits](https://docs.github.com/en/actions/reference/limits) — "Each job in a workflow can run for up to 6 hours of execution time".

### 2.4. The Go module `go.yaml.in/yaml/v4` v4.0.0

Read 2026-11-24 for the seeds and the distributions.

The Go module `go.yaml.in/yaml/v4` v4.0.0, its repository at the tag: [`LICENSE`](https://github.com/yaml/go-yaml/blob/v4.0.0/LICENSE) — "Apache License, Version 2.0, January 2004"; [`NOTICE`](https://github.com/yaml/go-yaml/blob/v4.0.0/NOTICE) — some of the files are "still covered by their original MIT license", and "All the remaining project files are covered by the Apache license". `module-kit` takes the same version in its `go.mod`.

### 2.5. AS 241

Read 2026-11-08 for the seeds and the distributions.

Wichura, M. J., "Algorithm AS 241: The percentage points of the normal distribution", Applied Statistics 37(3), 1988, pp. 477–484 — the coefficients PPND16. The article was not read here: the coefficients stand as its well-known publication gives them, and a test holds them — Φ(x) through `math.Erfc` returns p from 1e-300 to 1 − 1e-12 within an error of 1e-14 in x, and `inverseNormal(0.975)` = 1.959963984540054.
