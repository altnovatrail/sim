# sim's own recipes. The ci- ones are what .github's shared go.yml runs, one job each, in the
# toolchain stage of .devcontainer/Dockerfile (R16 .github; R10 in docs/decisions.md): their names
# are the same in every Go repository of NovaTrail, and what they run is this repository's own.
# The gate, which builds the bench on the commits of its gate file, is no ci- recipe: it comes with
# its own job (R09).

# The tools only the checks need, at exact versions. Each is a Go program run straight from its
# module, as engine's justfile runs them, so nothing is installed for it and go.mod gains no
# require: the Go module proxy serves that version and the checksum database vouches for it.
# Dependabot does not see them, so they are raised by hand.
GOVULNCHECK := "golang.org/x/vuln/cmd/govulncheck@v1.9.0"
GO_LICENSES := "github.com/google/go-licenses/v2@v2.0.1"

# List the recipes.
default:
    @just --list

# Fails on code gofmt would rewrite, on what go vet reports, and on a gopls lock file that names
# another version than GOPLS_VERSION in the Dockerfile; then ci-release (R29).
# Check formatting, go vet, the version of gopls and the release's file.
ci-lint:
    #!/usr/bin/env bash
    set -euo pipefail
    # The Go files git tracks or would, outside every testdata/, the root's too, which glob magic
    # takes and a bare **/ would not: spec/ is a submodule, another repository's, and git lists
    # none of its files. A file deleted from the working tree but still in the index is skipped,
    # so that gofmt does not fail on it with no word of formatting.
    unformatted="$(git ls-files -z --cached --others --exclude-standard -- '*.go' ':(exclude,glob)**/testdata/**' \
        | while IFS= read -r -d '' file; do if [ -e "$file" ]; then printf '%s\0' "$file"; fi; done \
        | xargs -0 -r gofmt -l)"
    if [ -n "$unformatted" ]; then
        echo "ci-lint: gofmt would rewrite these files; run gofmt -w on them:" >&2
        echo "$unformatted" >&2
        exit 1
    fi
    go vet ./...
    # gopls comes into the devcontainer image from its lock file, .devcontainer/gopls/, and the
    # image's build stops when the gopls it built does not name GOPLS_VERSION (R32). CI never
    # builds that stage, so the two are held to each other here: a pull request that raises one
    # alone, or a move to a new tag of spec that drops the step from the Dockerfile, fails here.
    lock="$(sed -n 's/^[[:space:]]*\(require[[:space:]]\{1,\}\)\{0,1\}golang\.org\/x\/tools\/gopls[[:space:]]\{1,\}v\([^[:space:]]\{1,\}\).*$/\2/p' \
        .devcontainer/gopls/go.mod)"
    arg="$(sed -n 's/^ARG GOPLS_VERSION=\(.*\)$/\1/p' .devcontainer/Dockerfile)"
    if [ -z "$lock" ] || [ "$arg" != "$lock" ]; then
        echo "ci-lint: .devcontainer/gopls/go.mod names gopls ${lock:-at no version} and" \
            ".devcontainer/Dockerfile has GOPLS_VERSION ${arg:-unset}; take both from the go kind of spec's repo-template (R32)" >&2
        exit 1
    fi
    {{ just_executable() }} ci-release

# The coverage goes into coverage/go.out, one Go profile, which the job coverage of .github's shared
# go.yml sends to Codecov (R48 spec; R10). Go counts statements and leaves the tests themselves out;
# a package that no test runs is in it too, its statements at zero. Every test counts for every
# package of the module, -coverpkg=./...: a test of one part runs the code of others, and its lines
# count where they are.
# The checks of the same bits on every build (R09) are the shared go.yml's, with same-bits:
# "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules". "Go / arm64" runs this recipe on
# arm64 hardware, then the targets. The recipes test-targets, fma-diff and rules below run the same
# in the container, and same-bits runs all three; CI calls none of them.
# No race detector: there is no goroutine in the code, which the rule goroutine keeps out, nor in
# the tests, and it made a learner of olympiad maths, from 0.1.1 on, some ten times slower, past the ten
# minutes go test gives a package (the amendment of 2026-11-09 to R10 in docs/decisions.md). A test
# that starts a goroutine brings it back.
# Run the tests and write their coverage.
ci-test:
    mkdir -p coverage
    go test -count=1 -covermode=count -coverpkg=./... -coverprofile=coverage/go.out ./...

# The checks of the same bits on every build (R09; docs/architecture.md, section 3.9) in the
# container, as the shared go.yml of .github runs them with same-bits, by its own scripts: .github
# at the commit ci.yml pins its go.yml to, so that the container runs what CI runs and nothing here
# is kept in step by hand, as engine's recipes do. CI calls none of these recipes. ci.yml must call
# go.yml once, by a commit's SHA. The checkout of that commit is kept in the user's cache, by the
# SHA and the version of Go, with the tool of the rules built in it, and reused: the first run of a
# commit reaches the network. The path of the tool's module holds .github, which go run does not
# take.
# Give the folder of the checkout of .github that ci.yml pins, fetching it once.
[private]
github:
    #!/usr/bin/env bash
    set -euo pipefail
    shas="$(sed -n 's|^[[:space:]]*uses: altnovatrail/\.github/\.github/workflows/go\.yml@\([0-9a-f]\{40\}\).*$|\1|p' .github/workflows/ci.yml)"
    if [ "$(printf '%s' "$shas" | grep -c .)" != 1 ]; then
        echo "same-bits: .github/workflows/ci.yml must call the go.yml of .github once, by a commit's SHA; it calls it ${shas:+by }$(printf '%s' "$shas" | tr '\n' ' ')" >&2
        exit 1
    fi
    # The tool is built by the Go of the container, so a move of Go builds it again.
    dir="${XDG_CACHE_HOME:-$HOME/.cache}/sim-github/$shas-$(go env GOVERSION)"
    if [ ! -f "$dir/.complete" ]; then
        rm -rf "$dir.partial"
        mkdir -p "$(dirname "$dir")"
        git clone --quiet --filter=blob:none --no-checkout https://github.com/altnovatrail/.github.git "$dir.partial"
        git -C "$dir.partial" checkout --quiet "$shas"
        go -C "$dir.partial/tools/coderules" build -o "$dir.partial/coderules" .
        touch "$dir.partial/.complete"
        rm -rf "$dir"
        mv "$dir.partial" "$dir"
    fi
    echo "$dir"

# The script a job of the pinned go.yml runs in the toolchain stage: the heredoc SCRIPT of the job,
# without the indentation of the YAML, which its first line that is not blank gives. A job with no
# such script stops the recipe with the reason.
[private]
script job:
    #!/usr/bin/env bash
    set -euo pipefail
    file="$({{ just_executable() }} github)/.github/workflows/go.yml"
    awk -v job="{{ job }}" '
        $0 == "  " job ":" { injob = 1; next }
        injob && /^  [a-z0-9-]+:$/ { injob = 0 }
        injob && !inscript && /<<'"'"'SCRIPT'"'"'$/ { inscript = 1; next }
        inscript && /^[[:space:]]*SCRIPT$/ { done = 1; exit }
        inscript {
            if (indent == "" && $0 ~ /[^ ]/) { match($0, /^ */); indent = RLENGTH }
            print substr($0, indent + 1)
        }
        END { if (!done) exit 1 }
    ' "$file" || { echo "same-bits: the job {{ job }} of $file has no script SCRIPT" >&2; exit 1; }

# The tests on the builds where the bits are most at risk — GOAMD64 v3, v1 with
# GODEBUG=cpu.all=off, 386 and js/wasm in the container's Node, and wasip1 built and vetted —, by
# the script of the job targets.
# Run the tests on the builds the same bits are held on, as "Go / targets" does.
test-targets:
    #!/usr/bin/env bash
    set -euo pipefail
    script="$({{ just_executable() }} script targets)"
    bash -c "$script"

# The assembler of each package with fusion and without it on six targets, with the job's own
# control, by the script of the job fma-diff.
# Compare the assembler with fusion and without it, as "Go / fma-diff" does.
fma-diff:
    #!/usr/bin/env bash
    set -euo pipefail
    script="$({{ just_executable() }} script fma-diff)"
    bash -c "$script"

# The rules of code, by the tool tools/coderules of the pinned .github and this repository's
# lists, .github/coderules.json, as the job rules runs it.
# Check the rules of code, as "Go / rules" does.
rules:
    #!/usr/bin/env bash
    set -euo pipefail
    dir="$({{ just_executable() }} github)"
    "$dir/coderules" -lists .github/coderules.json .

# What a task's report names as checked on every target.
# Run test-targets, fma-diff and rules.
same-bits: test-targets fma-diff rules

# The release's one file (R29): sim-<version>.tar.gz, a folder of that name with what builds the
# bench and runs the gate — the Go code but its tests and testdata/, the packages and the paths it
# embeds, go.mod, go.sum, LICENSE, .gitattributes, gate.json, this justfile and the Dockerfile of
# its toolchain —, the paths that give a release but the release workflow itself. The tar is
# written the same by every build: names in order, no times or owners, mode 0644, and gzip with no
# name or time. .github's shared release runs it in a tree git archive made, with no git to ask, so
# the files come from the tree. The callers of the gate take sim from git by the tag and the SHA of
# its commit; this file binds the tag to the signed release (R77 spec).
# Build the release's file into dist/.
release-files version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{ quote(version) }}
    if ! [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
        echo "release-files: the version is v<major>.<minor>.<patch>, as v0.1.0, not ${version}" >&2
        exit 1
    fi
    name="sim-${version}"
    rm -rf dist
    mkdir dist
    (printf '%s\n' LICENSE go.mod go.sum .gitattributes gate.json justfile .devcontainer/Dockerfile internal/criterion/paths.json
        find internal/published -path 'internal/published/v*' -type f ! -name '*.go' -print
        find . -name testdata -prune -o -path './.*' -prune -o -path ./spec -prune -o \
            -type f -name '*.go' ! -name '*_test.go' -print | sed 's|^\./||') |
        LC_ALL=C sort -u |
        tar --create --verbatim-files-from --files-from=- --format=posix \
            --pax-option=exthdr.name=%d/PaxHeaders/%f,delete=atime,delete=ctime \
            --transform="s,^,${name}/," --mtime=@0 --owner=0 --group=0 --numeric-owner --mode=0644 |
        gzip -9n > "dist/${name}.tar.gz"

# Builds the release's file twice, a second apart, and fails unless the two are the same bytes;
# unless the archive holds what git tracks or would at the paths release-files takes, as it is in
# the tree, and nothing else; and unless the unpacked archive builds. So the file a release signs is
# what builds and runs the gate at its commit, and anyone can build it again with the same tar and
# gzip and compare. It refuses to run beside a dist/ of one's own, which release-files would
# replace. ci-lint runs it, so "Go / lint" holds it (R29).
# Check that the release's file is the code, and the same bytes on every build.
ci-release:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -e dist ]; then
        echo "ci-release: dist/ is there already, and release-files would replace it: move it away first" >&2
        exit 1
    fi
    version=v0.0.0-ci
    name="sim-${version}"
    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp" dist' EXIT
    {{ just_executable() }} release-files "$version"
    mv dist "$tmp/first"
    sleep 1
    {{ just_executable() }} release-files "$version"
    if [ "$(ls dist)" != "${name}.tar.gz" ] || ! cmp -s "$tmp/first/${name}.tar.gz" "dist/${name}.tar.gz"; then
        echo "ci-release: two builds of the release's files differ, or dist/ holds another file:" >&2
        ls -l "$tmp/first" dist >&2
        exit 1
    fi
    mkdir "$tmp/unpacked"
    tar -xzf "dist/${name}.tar.gz" -C "$tmp/unpacked"
    got="$(cd "$tmp/unpacked/$name" && find . -type f | sed 's|^\./||' | LC_ALL=C sort)"
    # The files git tracks or would, as the gofmt step of ci-lint takes them, but those deleted from
    # the tree: a new file before git add, or one removed before git rm, is the tree's as the release
    # sees it.
    want="$(git ls-files --cached --others --exclude-standard -- LICENSE go.mod go.sum .gitattributes gate.json justfile \
        .devcontainer/Dockerfile internal/criterion/paths.json 'internal/published/v*' '*.go' \
        ':(exclude)*_test.go' ':(exclude,glob)**/testdata/**' |
        while IFS= read -r file; do if [ -e "$file" ]; then printf '%s\n' "$file"; fi; done | LC_ALL=C sort -u)"
    if [ "$got" != "$want" ]; then
        echo "ci-release: the archive holds" >&2
        echo "$got" >&2
        echo "and git tracks" >&2
        echo "$want" >&2
        exit 1
    fi
    while IFS= read -r file; do
        if ! cmp -s "$file" "$tmp/unpacked/$name/$file"; then
            echo "ci-release: ${file} in the archive is not the tree's" >&2
            exit 1
        fi
    done <<< "$want"
    go -C "$tmp/unpacked/$name" build ./...
    echo "ci-release: ${name}.tar.gz is what builds and runs the gate, the same bytes on two builds" >&2

# Fail on a known vulnerability in code this module reaches.
ci-vuln:
    go run {{ GOVULNCHECK }} ./...

# Fails if a package this module is built from carries a license outside ALLOWED_LICENSES, SPDX
# identifiers separated by commas, which .github's go.yml sets.
# Check the licenses of what the module is built from against ALLOWED_LICENSES.
ci-licenses:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "${ALLOWED_LICENSES:-}" ]; then
        echo "ci-licenses: ALLOWED_LICENSES is not set. CI takes it from .github's go.yml; to run this here, set it to the list there." >&2
        exit 1
    fi
    go run {{ GO_LICENSES }} check ./... --allowed_licenses="$ALLOWED_LICENSES"

# The workspace a build on the libraries of gate.json goes through (docs/architecture.md, section
# 3.9; R26, R34 in docs/decisions.md), one for the gate and for the decisive run of a candidate. Each
# library comes from <prefix>_<NAME> — GATE_ENGINE, CANDIDATE_LEARNER_MODEL — when the caller sets
# it: a commit by its full SHA, fetched from GitHub even when it is on no branch, as the head of a
# pull request is, or a folder, of the caller's folder from when relative, which with clean=yes must
# have nothing uncommitted, so that the commit a file of results names is the code that ran. A
# library the caller does not name comes at the commit gate.json names, which is the commit of its
# tag. It writes dir/go.work with sim and the two libraries, and a line a library, its fields apart
# by tabs: its module path, its folder, the commit it is at, and the variable that named it or -.
# Write the workspace of a build on the libraries of gate.json.
[private]
[positional-arguments]
workspace dir prefix clean from:
    #!/usr/bin/env bash
    set -euo pipefail
    dir="$1"
    prefix="$2"
    clean="$3"
    from="$4"
    [ -x "$dir/plain" ] || GOWORK=off go build -o "$dir/plain" .
    libraries="$("$dir/plain" gate libraries gate.json)"
    uses=("$PWD")
    lines=()
    while read -r path repo commit version; do
        name="${path##*/}"
        var="${prefix}_$(printf '%s' "$name" | tr 'a-z-' 'A-Z_')"
        given="${!var:-}"
        lib="$given"
        case "$lib" in /*|'') ;; *) lib="$from/$lib" ;; esac
        if [ -z "$given" ] || [[ "$given" =~ ^[0-9a-f]{40}$ ]]; then
            lib="$dir/$name"
            git init -q "$lib"
            git -C "$lib" fetch -q --depth 1 "https://github.com/$repo" "${given:-$commit}"
            git -C "$lib" -c advice.detachedHead=false checkout -q FETCH_HEAD
        fi
        if [ -z "$given" ] && [ -n "${version:-}" ]; then
            # The tag gate.json names is on the commit it names: the job "Go" builds go.mod's tag,
            # and the gate the commit (R26). An annotated tag points at its commit through ^{}.
            tags="$(git ls-remote "https://github.com/$repo" "refs/tags/$version" "refs/tags/$version^{}")"
            tagged="$(printf '%s\n' "$tags" | awk -v ref="refs/tags/$version^{}" '$2 == ref { print $1 }')"
            if [ -z "$tagged" ]; then
                tagged="$(printf '%s\n' "$tags" | awk '{ print $1 }')"
            fi
            if [ "$tagged" != "$commit" ]; then
                echo "workspace: $path $version is on ${tagged:-no commit}, and gate.json names $commit" >&2
                exit 1
            fi
        fi
        if [ ! -d "$lib" ]; then
            echo "workspace: $var is $given: a commit by its full SHA, or a folder" >&2
            exit 1
        fi
        lib="$(cd "$lib" && pwd)"
        if [ -n "$given" ] && [ "$clean" = yes ] && [ -n "$(git -C "$lib" status --porcelain)" ]; then
            echo "workspace: $var, $lib, has changes not committed, and a file of results names a commit" >&2
            exit 1
        fi
        head="$(git -C "$lib" rev-parse HEAD)"
        echo "workspace: $path at $head${given:+, from $var}" >&2
        uses+=("$lib")
        named=-
        [ -z "$given" ] || named="$var"
        lines+=("$path	$lib	$head	$named")
    done <<< "$libraries"
    {
        echo "go $(sed -n 's/^go //p' go.mod)"
        echo "use ("
        for u in "${uses[@]}"; do echo "	$u"; done
        echo ")"
    } > "$dir/go.work"
    printf '%s\n' "${lines[@]}"

# The decisive run of a rule (docs/architecture.md, sections 3.6, 3.8 and 3.9; R09, R28, R33, R34 in
# docs/decisions.md): 2,000 learners of a set — working, or held-out-<criterion>-<n> or
# held-out-<criterion>-<variant>-<n> for a confirmation — of every cell of every population, as
# processes, one part of the learners each, as many as the processors, merged into one file of
# results, whose SHA-256 it prints. Without a candidate it builds the bench on the libraries go.mod
# names, the base while it names the releases gate.json does, so that the base's file is the one
# paths.json names to the byte. A candidate is CANDIDATE_LEARNER_MODEL or CANDIDATE_ENGINE, a commit
# by its full SHA or a folder with nothing uncommitted; a pair is both with CANDIDATE_PAIR=yes, so
# that a variable left set from another run does not make a pair unseen. The bench is then built
# through the workspace of the gate, with the other library at the commit of gate.json, the build is
# held to gate.json as the gate holds it before anything runs, and the file names both libraries by
# commit. A set is checked before the run. The file is no part of the repository: its seeds and the
# same bits on every build give it again to the byte, and runs/, which git ignores, keeps it between
# judgements. rule is engine, for the base or a candidate and the oracle on its answers, or ceiling,
# for the ceiling of the engine; a relative out or folder is of the caller's folder:
# `CANDIDATE_ENGINE=<sha> just decisive engine runs/engine-a.jsonl`.
# Run the decisive run of a rule into a file of results: rule engine or ceiling, set working or held-out-….
decisive rule out set="working": (cells "decisive" rule out set "2000")

# A rough look (docs/architecture.md, sections 3.6 and 3.8; R49 in docs/decisions.md): the run of
# just decisive, the same build and the same candidates, on fewer learners a cell — from 1 to 1,999,
# 300 unless named, as engine's rough looks ran — and on the working learners alone: a look at
# held-out learners would leave them held out no longer. sim judge names its report a rough look,
# and it decides nothing. The base's rough look runs the same way, without a candidate, on as many
# learners: `just rough engine runs/base-300.jsonl` and
# `CANDIDATE_ENGINE=<sha> just rough engine runs/engine-a-300.jsonl`.
# Run a rough look of a rule on fewer working learners than a decisive run: rule engine or ceiling.
rough rule out learners="300": (cells "rough" rule out "working" learners)

# The body of just decisive and just rough: the runs of a rule on so many learners a cell of a set,
# a rough look only of fewer than a decisive run's 2,000 and only on the working learners.
[private]
cells kind rule out set learners:
    #!/usr/bin/env bash
    set -euo pipefail
    kind={{ quote(kind) }}
    case {{ quote(rule) }} in engine|ceiling) ;; *) echo "$kind: rule {{ rule }}: engine or ceiling" >&2; exit 2 ;; esac
    learners={{ quote(learners) }}
    case "$learners" in ''|*[!0-9]*|0*) echo "$kind: learners $learners: a whole number above zero" >&2; exit 2 ;; esac
    if [ "$kind" = rough ] && { [ "$learners" -ge 2000 ] || [ {{ quote(set) }} != working ]; }; then
        echo "rough: $learners learners of {{ set }}: a rough look is fewer than 2,000 working learners a cell, and 2,000 a decisive run, just decisive" >&2
        exit 2
    fi
    out={{ quote(out) }}
    case "$out" in /*) ;; *) out={{ quote(invocation_directory()) }}/"$out" ;; esac
    mkdir -p "$(dirname "$out")"
    tmp="$(mktemp -d)"
    pids=()
    trap 'for p in "${pids[@]}"; do kill "$p" 2>/dev/null || true; done; wait; rm -rf "$tmp"' EXIT
    named=()
    while IFS='=' read -r var _; do
        [ "$var" = CANDIDATE_PAIR ] || named+=("$var")
    done < <(env | grep -E '^CANDIDATE_[A-Z_]+=.' || true)
    if [ "${#named[@]}" -gt 1 ] && [ "${CANDIDATE_PAIR:-}" != yes ]; then
        echo "$kind: ${named[*]} name a pair; set CANDIDATE_PAIR=yes to run it" >&2
        exit 2
    fi
    commits=()
    if [ "${#named[@]}" -gt 0 ]; then
        libraries="$({{ just_executable() }} -f {{ quote(justfile()) }} workspace "$tmp" CANDIDATE yes {{ quote(invocation_directory()) }})"
        gomods=(-go "$PWD")
        taken=0
        while IFS=$'\t' read -r path dir head from; do
            gomods+=(-go "$dir")
            commits+=(-commit "$path=$head")
            [ "$from" = - ] || taken=$((taken + 1))
        done <<< "$libraries"
        if [ "$taken" != "${#named[@]}" ]; then
            echo "$kind: ${named[*]}: a variable names no library of gate.json" >&2
            exit 2
        fi
        GOWORK="$tmp/go.work" go build -o "$tmp/sim" .
        "$tmp/sim" gate verify -file gate.json "${gomods[@]}"
    else
        echo "$kind: the libraries go.mod names" >&2
        GOWORK=off go build -o "$tmp/sim" .
    fi
    "$tmp/sim" run -populations G7 -learners 1 -set {{ quote(set) }} ${commits[@]+"${commits[@]}"} > /dev/null
    n="$(nproc)"
    for k in $(seq 0 $((n - 1))); do
        "$tmp/sim" run -populations all -learners "$learners" -set {{ quote(set) }} -rule {{ quote(rule) }} ${commits[@]+"${commits[@]}"} -part "$k/$n" -out "$tmp/part-$k.jsonl" &
        pids+=("$!")
    done
    failed=0
    for p in "${pids[@]}"; do wait "$p" || failed=1; done
    pids=()
    if [ "$failed" != 0 ]; then
        echo "$kind: a part of the run failed" >&2
        exit 1
    fi
    "$tmp/sim" merge -out "$out" "$tmp"/part-*.jsonl
    sha256sum "$out"

# The gate (docs/architecture.md, sections 3.7 and 3.9; R09, R26, R27 in docs/decisions.md). It
# builds the bench only on the commits gate.json names, through the workspace the recipe workspace
# writes in a temporary folder: each library at its commit, or the caller's own code in its place
# from GATE_LEARNER_MODEL or GATE_ENGINE, a folder or a commit by its full SHA. It points GOWORK at
# that workspace and builds sim once. The run goes as processes, one part of the learners each, as many as the
# processors: R09 keeps goroutines out of code. Then sim holds the build to gate.json — the versions
# the workspace chose, the go lines, the point sim's go.mod names, the packages — and the numbers to
# their bands, and writes the table. mode is bands, as the callers hold them; exact, the recorded
# numbers to the last bit, as sim's own CI holds them on amd64 and arm64; record, which writes the
# bands into gate.json anew, for the reason given as the one argument; or modules, the gate of the
# CI of modules, whose arguments are those of `sim gate changed` (README.md): it asks sim what
# changed of what the bench reads, before it fetches anything, runs the gate on the packages of the
# changed modules alone, and holds the numbers to thresholds, not bands. Relative paths among those
# arguments are of the caller's folder.
# Run the gate on the commits of gate.json: mode bands, exact, record or modules.
[positional-arguments]
gate mode="bands" *args:
    #!/usr/bin/env bash
    set -euo pipefail
    mode="$1"
    shift
    case "$mode" in
        bands|exact) [ "$#" -eq 0 ] || { echo "gate: mode $mode takes no argument" >&2; exit 2; } ;;
        record) [ "$#" -eq 1 ] || { echo "gate: mode record takes the reason, one argument" >&2; exit 2; } ;;
        modules) [ "$#" -gt 0 ] || { echo "gate: mode modules takes the arguments of sim gate changed" >&2; exit 2; } ;;
        *) echo "gate: mode $mode: bands, exact, record or modules" >&2; exit 2 ;;
    esac
    root="$PWD"
    tmp="$(mktemp -d)"
    pids=()
    # A part that fails stops the others before the folder they write into goes.
    trap 'for p in "${pids[@]}"; do kill "$p" 2>/dev/null || true; done; wait; rm -rf "$tmp"' EXIT
    GOWORK=off go build -o "$tmp/plain" .
    packages=()
    if [ "$mode" = modules ]; then
        changed=()
        for a in "$@"; do
            case "$a" in
                -*=/*|-*=) changed+=("$a") ;;
                -*=*) changed+=("${a%%=*}="{{ quote(invocation_directory()) }}/"${a#*=}") ;;
                -*|/*) changed+=("$a") ;;
                *) changed+=({{ quote(invocation_directory()) }}/"$a") ;;
            esac
        done
        run="$("$tmp/plain" gate changed "${changed[@]}")"
        if [ -z "$run" ]; then
            echo "gate: nothing the bench reads changed, so the gate of modules has nothing to run"
            exit 0
        fi
        # A line "along <folder>" is a package that did not change and runs only on the graphs it
        # shares with one that did (R31 in docs/decisions.md).
        while IFS= read -r line; do
            case "$line" in
                "along "*) packages+=(-along "${line#along }") ;;
                *) packages+=(-package "$line") ;;
            esac
        done <<< "$run"
    fi
    libraries="$({{ just_executable() }} -f {{ quote(justfile()) }} workspace "$tmp" GATE no {{ quote(invocation_directory()) }})"
    gomods=(-go "$root")
    commits=()
    while IFS=$'\t' read -r path dir head _; do
        gomods+=(-go "$dir")
        commits+=(-commit "$path=$head")
    done <<< "$libraries"
    export GOWORK="$tmp/go.work"
    go build -o "$tmp/sim" .
    n="$(nproc)"
    for k in $(seq 0 $((n - 1))); do
        "$tmp/sim" gate run -part "$k/$n" "${commits[@]}" ${packages[@]+"${packages[@]}"} -out "$tmp/part-$k.jsonl" &
        pids+=("$!")
    done
    failed=0
    for p in "${pids[@]}"; do wait "$p" || failed=1; done
    pids=()
    if [ "$failed" != 0 ]; then
        echo "gate: a part of the run failed" >&2
        exit 1
    fi
    case "$mode" in
        bands) "$tmp/sim" gate check -file gate.json "${gomods[@]}" "$tmp"/part-*.jsonl ;;
        exact) "$tmp/sim" gate check -exact -file gate.json "${gomods[@]}" "$tmp"/part-*.jsonl ;;
        record) "$tmp/sim" gate record -reason "$1" -file gate.json "${gomods[@]}" "$tmp"/part-*.jsonl ;;
        modules) "$tmp/sim" gate thresholds -file gate.json "${gomods[@]}" "$tmp"/part-*.jsonl ;;
    esac
