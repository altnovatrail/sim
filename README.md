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

A component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring.

It stands beside learner-model and engine, the two libraries it imports at build time and whose CI takes it as a gate at an exact tag, so that a change that worsens the criterion is not merged; the CI of modules runs it when what it reads of a package changes; and it derives its learners' seeds with spec's random generator, one for the whole platform.

Its card is in the architecture of [spec](https://github.com/altnovatrail/spec); its own documents, its decisions and its plan are in [docs/](docs/README.md).
