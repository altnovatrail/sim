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
# another version than GOPLS_VERSION in the Dockerfile.
# Check formatting, go vet and the version of gopls.
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
    # image's build stops when the gopls it built does not name GOPLS_VERSION (R11). CI never
    # builds that stage, so the two are held to each other here: a pull request that raises one
    # alone, or a move to a new tag of spec that drops the step from the Dockerfile, fails here.
    lock="$(sed -n 's/^[[:space:]]*\(require[[:space:]]\{1,\}\)\{0,1\}golang\.org\/x\/tools\/gopls[[:space:]]\{1,\}v\([^[:space:]]\{1,\}\).*$/\2/p' \
        .devcontainer/gopls/go.mod)"
    arg="$(sed -n 's/^ARG GOPLS_VERSION=\(.*\)$/\1/p' .devcontainer/Dockerfile)"
    if [ -z "$lock" ] || [ "$arg" != "$lock" ]; then
        echo "ci-lint: .devcontainer/gopls/go.mod names gopls ${lock:-at no version} and" \
            ".devcontainer/Dockerfile has GOPLS_VERSION ${arg:-unset}; raise them together (R11)" >&2
        exit 1
    fi

# The coverage goes into coverage/go.out, one Go profile, which the job coverage of .github's shared
# go.yml sends to Codecov (R48 spec; R10). Go counts statements and leaves the tests themselves out;
# a package that no test runs is in it too, its statements at zero. Every test counts for every
# package of the module, -coverpkg=./...: a test of one part runs the code of others, and its lines
# count where they are.
# The checks of the same bits on every build (R09) are the shared go.yml's, with same-bits:
# "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules". "Go / arm64" runs this recipe on
# arm64 hardware, then the targets. The recipes test-targets, fma-diff and rules below run the same
# in the container, and same-bits runs all three; CI calls none of them.
# Run the tests, with the race detector, and write their coverage.
ci-test:
    mkdir -p coverage
    go test -race -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage/go.out ./...

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

# The gate (docs/architecture.md, sections 3.7 and 3.9; R09, R26, R27 in docs/decisions.md). It
# builds the bench only on the commits gate.json names: it fetches each library at its commit into
# a temporary folder — or takes the caller's own code from GATE_LEARNER_MODEL or GATE_ENGINE, a
# folder, in its place —, writes go.work there with the three modules and points GOWORK at it, and
# builds sim once. The run goes as processes, one part of the learners each, as many as the
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
        while IFS= read -r folder; do packages+=(-package "$folder"); done <<< "$run"
    fi
    libraries="$("$tmp/plain" gate libraries gate.json)"
    uses=("$root")
    gomods=(-go "$root")
    commits=()
    while read -r path repo commit version; do
        name="${path##*/}"
        var="GATE_$(printf '%s' "$name" | tr 'a-z-' 'A-Z_')"
        dir="${!var:-}"
        if [ -z "$dir" ]; then
            dir="$tmp/$name"
            git init -q "$dir"
            git -C "$dir" fetch -q --depth 1 "https://github.com/$repo" "$commit"
            git -C "$dir" -c advice.detachedHead=false checkout -q FETCH_HEAD
            # The tag gate.json names is on the commit it names: the job "Go" builds go.mod's tag,
            # and the gate the commit (R26). An annotated tag points at its commit through ^{}.
            if [ -n "${version:-}" ]; then
                tags="$(git ls-remote "https://github.com/$repo" "refs/tags/$version" "refs/tags/$version^{}")"
                tagged="$(printf '%s\n' "$tags" | awk -v ref="refs/tags/$version^{}" '$2 == ref { print $1 }')"
                if [ -z "$tagged" ]; then
                    tagged="$(printf '%s\n' "$tags" | awk '{ print $1 }')"
                fi
                if [ "$tagged" != "$commit" ]; then
                    echo "gate: $path $version is on ${tagged:-no commit}, and gate.json names $commit" >&2
                    exit 1
                fi
            fi
        fi
        dir="$(cd "$dir" && pwd)"
        head="$(git -C "$dir" rev-parse HEAD)"
        echo "gate: $path at $head${!var:+, from the caller}" >&2
        uses+=("$dir")
        gomods+=(-go "$dir")
        commits+=(-commit "$path=$head")
    done <<< "$libraries"
    {
        echo "go $(sed -n 's/^go //p' go.mod)"
        echo "use ("
        for u in "${uses[@]}"; do echo "	$u"; done
        echo ")"
    } > "$tmp/go.work"
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
