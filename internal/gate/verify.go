package gate

import (
	"bufio"
	"bytes"
	"fmt"
	"go/version"
	"maps"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/altnovatrail/sim/internal/published"
)

// Module is a module of the gate's workspace as its go.mod says: its path, its go line, and the
// versions it requires, by path.
type Module struct {
	Path     string
	Go       string
	Requires map[string]string
}

// ParseGoMod reads the module path, the go line and the requirements of a go.mod.
func ParseGoMod(data []byte) (Module, error) {
	m := Module{Requires: map[string]string{}}
	block := false
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := sc.Text()
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[:i]
		}
		f := fields(line)
		switch {
		case len(f) == 0:
		case block && f[0] == ")":
			block = false
		case block && len(f) >= 2:
			m.Requires[f[0]] = f[1]
		case f[0] == "module" && len(f) == 2:
			m.Path = f[1]
		case f[0] == "go" && len(f) == 2:
			m.Go = f[1]
		case f[0] == "require" && len(f) == 2 && f[1] == "(":
			block = true
		case f[0] == "require" && len(f) >= 3:
			m.Requires[f[1]] = f[2]
		}
	}
	if m.Path == "" || m.Go == "" {
		return m, fmt.Errorf("gate: a go.mod with no module path or no go line")
	}
	return m, nil
}

// fields splits a line of go.mod at its spaces and tabs, by bytes: go.mod is ASCII where it
// matters here.
func fields(line string) []string {
	var out []string
	for f := range strings.FieldsFuncSeq(line, func(r rune) bool { return r == 32 || r == 9 || r == 13 }) {
		out = append(out, f)
	}
	return out
}

// Verify gives every way the build of the gate departs from the gate file (docs/architecture.md,
// section 3.9): a library not built from the workspace, a shared dependency the workspace chose at
// another version, a module whose go line is another than the recorded or newer than the Go that
// built the gate, sim's go.mod requiring a library at another point than the file names — the tag
// of its release, or the pseudo-version of its commit while it is in none —, and packages other than
// those the bench embeds.
func Verify(f *File, info *debug.BuildInfo, goVersion string, modules []Module) []string {
	var out []string
	dep := func(path string) *debug.Module {
		if info == nil {
			return nil
		}
		i := slices.IndexFunc(info.Deps, func(d *debug.Module) bool { return d.Path == path })
		if i < 0 {
			return nil
		}
		if d := info.Deps[i]; d.Replace != nil {
			return d.Replace
		}
		return info.Deps[i]
	}
	for _, path := range slices.Sorted(maps.Keys(f.Libraries)) {
		if d := dep(path); d != nil && d.Version != "(devel)" {
			out = append(out, fmt.Sprintf("%s is built at %s, not from the workspace the gate file names", path, d.Version))
		}
	}
	for _, path := range slices.Sorted(maps.Keys(f.Dependencies)) {
		want := f.Dependencies[path]
		if d := dep(path); d == nil || d.Version != want {
			got := "nothing"
			if d != nil {
				got = d.Version
			}
			out = append(out, fmt.Sprintf("the workspace chose %s at %s, and the gate file records %s", path, got, want))
		}
	}
	byPath := map[string]Module{}
	for _, m := range modules {
		byPath[m.Path] = m
	}
	for _, path := range slices.Sorted(maps.Keys(f.Go)) {
		m, ok := byPath[path]
		switch {
		case !ok:
			out = append(out, fmt.Sprintf("the gate was given no go.mod of %s", path))
		case m.Go != f.Go[path]:
			out = append(out, fmt.Sprintf("%s has the go line %s, and the gate file records %s", path, m.Go, f.Go[path]))
		case version.Compare("go"+m.Go, goVersion) > 0:
			out = append(out, fmt.Sprintf("%s has the go line %s, newer than the Go of the gate, %s", path, m.Go, goVersion))
		}
	}
	if sim, ok := byPath["github.com/altnovatrail/sim"]; ok {
		for _, path := range slices.Sorted(maps.Keys(f.Libraries)) {
			lib, got := f.Libraries[path], sim.Requires[path]
			want := lib.Version
			if want == "" && len(lib.Commit) >= 12 && strings.HasSuffix(got, "-"+lib.Commit[:12]) {
				want = got // the pseudo-version of the commit
			}
			if got != want || got == "" {
				out = append(out, fmt.Sprintf("sim's go.mod requires %s at %q, and the gate file names %s", path, got, point(lib)))
			}
		}
	} else {
		out = append(out, "the gate was given no go.mod of github.com/altnovatrail/sim")
	}
	out = append(out, pinned(f.Modules)...)
	return out
}

// point names the point a gate file names a library at.
func point(l Library) string {
	if l.Version != "" {
		return l.Version
	}
	return "the pseudo-version of commit " + l.Commit
}

// pinned gives every way the packages the gate file names depart from those the bench embeds.
func pinned(m Modules) []string {
	var out []string
	if m.Release != published.Release || m.Index != "sha256:"+published.IndexDigest {
		out = append(out, fmt.Sprintf("the gate file names modules %s with index %s, and the bench embeds %s with sha256:%s", m.Release, m.Index, published.Release, published.IndexDigest))
	}
	ps, err := published.Pinned()
	if err != nil {
		return append(out, err.Error())
	}
	if !slices.Equal(ps, m.Packages) {
		out = append(out, fmt.Sprintf("the gate file names the packages %v, and the bench embeds %v", m.Packages, ps))
	}
	return out
}
