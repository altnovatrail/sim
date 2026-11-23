package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/altnovatrail/sim/internal/bench"
	"github.com/altnovatrail/sim/internal/criterion"
	"github.com/altnovatrail/sim/internal/gate"
	"github.com/altnovatrail/sim/internal/measure"
	"github.com/altnovatrail/sim/internal/population"
	"github.com/altnovatrail/sim/internal/seed"
)

// fairRun writes a file of results that stands in for a decisive run of the base with the cells of
// fairness by the grade alone, on the libraries of the embedded gate file, and gives its SHA-256.
func fairRun(t *testing.T, path string) string {
	t.Helper()
	g, err := gate.Read(bytes.NewReader(gateJSON))
	if err != nil {
		t.Fatal(err)
	}
	f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Libraries: map[string]string{}, Learners: criterion.Decisive, Invariants: true}}
	// A library in no release is named by a pseudo-version of its commit, as go.mod names it then.
	for p, lib := range g.Libraries {
		f.Header.Libraries[p] = lib.Version
		if lib.Version == "" {
			f.Header.Libraries[p] = "v0.0.0-20261117000000-" + lib.Commit[:12]
		}
	}
	for _, p := range population.Fairness {
		f.Header.Populations = append(f.Header.Populations, p.Name)
		for n := range criterion.Decisive {
			tr := &measure.Trace{Levels: make([]float64, measure.FairLength), Bands: []int{1, 1, 1}}
			f.Lines = append(f.Lines, bench.Line{Population: p.Name, Graph: []string{population.Olympiad}, N: n, Sums: measure.Sums{"learners": 1}, Fair: tr})
		}
	}
	var b bytes.Buffer
	if err := bench.Write(f, &b); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(b.Bytes())
	return hex.EncodeToString(sum[:])
}

// sim evaluate writes the report of an evaluation whose run is the file of the SHA-256 it names, and
// refuses another file, a date not YYYY-MM-DD, a field it does not know, a name not its folder's, and
// a call without the one evaluation file.
func TestEvaluate(t *testing.T) {
	dir := t.TempDir()
	folder := filepath.Join(dir, "evaluations", "fairness-1")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	digest := fairRun(t, filepath.Join(dir, "base.jsonl"))
	write := func(name, body string) string {
		path := filepath.Join(folder, name)
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	file := func(name, sha, date string) string {
		e, _ := json.Marshal(evaluation{Evaluation: name, Run: "../../base.jsonl", SHA256: sha, Date: date})
		return string(e)
	}
	var stdout, stderr bytes.Buffer
	good := write("evaluation.json", file("fairness-1", digest, "2026-11-14"))
	if got := run([]string{"evaluate", good}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim evaluate = %d: %s", got, stderr.String())
	}
	for _, want := range []string{"# The fairness of the learner model by the grade: fairness-1", "SHA-256 " + digest, "The base holds the criterion.", "On 2026-11-14"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("the report does not say %q:\n%s", want, stdout.String())
		}
	}
	for name, body := range map[string]string{
		"another file":   file("fairness-1", strings.Repeat("0", 64), "2026-11-14"),
		"a date":         file("fairness-1", digest, "14.11.2026"),
		"a month":        file("fairness-1", digest, "2026-13-01"),
		"a day":          file("fairness-1", digest, "2026-02-29"),
		"a sign":         file("fairness-1", digest, "2026-+1-01"),
		"a short digest": file("fairness-1", "abc", "2026-11-14"),
		"another folder": file("fairness-2", digest, "2026-11-14"),
		"an extra field": `{"evaluation":"fairness-1","run":"../../base.jsonl","sha256":"` + digest + `","date":"2026-11-14","extra":1}`,
		"no run":         `{"evaluation":"fairness-1","run":"none.jsonl","sha256":"` + digest + `","date":"2026-11-14"}`,
	} {
		stderr.Reset()
		if got := run([]string{"evaluate", write("bad.json", body)}, &stdout, &stderr); got != 1 {
			t.Errorf("%s: sim evaluate = %d: %s", name, got, stderr.String())
		}
	}
	for date, want := range map[string]bool{"2028-02-29": true, "2000-02-29": true, "2100-02-29": false, "2026-04-31": false, "2026-12-31": true, "2026-11-00": false} {
		if isDate(date) != want {
			t.Errorf("isDate(%s) = %v", date, !want)
		}
	}
	if got := run([]string{"evaluate", filepath.Join(folder, "none.json")}, &stdout, &stderr); got != 1 {
		t.Errorf("no evaluation file: sim evaluate = %d", got)
	}
	if got := run([]string{"evaluate"}, &stdout, &stderr); got != 2 {
		t.Errorf("sim evaluate with no file = %d", got)
	}
}

// Every evaluation of evaluations/ has its evaluation file, which names its folder, a run by its
// SHA-256 and a date, and its report; fairness-1 is the evaluation of the base.
func TestEvaluations(t *testing.T) {
	dirs, err := filepath.Glob(filepath.Join("evaluations", "*"))
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range dirs {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			continue
		}
		e, err := readEvaluation(filepath.Join(dir, "evaluation.json"))
		if err != nil {
			t.Errorf("%s: %v", dir, err)
			continue
		}
		if !isSHA256(e.SHA256) || e.Run == "" {
			t.Errorf("%s names the run %q by %q", dir, e.Run, e.SHA256)
		}
		if _, err := os.Stat(filepath.Join(dir, "report.md")); err != nil {
			t.Errorf("%s: %v", dir, err)
		}
	}
}

// An evaluation file of the three rules names its criterion and the base's run beside its own, and
// one of fairness neither; sim evaluate reads both runs by their SHA-256 and holds the base to the
// paths'.
func TestEvaluateRules(t *testing.T) {
	dir := t.TempDir()
	folder := filepath.Join(dir, "evaluations", "rules-1")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	digest := fairRun(t, filepath.Join(dir, "run.jsonl"))
	write := func(e evaluation) string {
		body, _ := json.Marshal(e)
		path := filepath.Join(folder, "evaluation.json")
		if err := os.WriteFile(path, body, 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	rules := func(base, sha string) evaluation {
		return evaluation{Evaluation: "rules-1", Criterion: "rules", Run: "../../run.jsonl", SHA256: digest, Base: base, BaseSHA256: sha, Date: "2026-11-17"}
	}
	for name, e := range map[string]evaluation{
		"no base":              rules("", ""),
		"a short base digest":  rules("../../run.jsonl", "abc"),
		"another base digest":  rules("../../run.jsonl", strings.Repeat("0", 64)),
		"no base file":         rules("../../none.jsonl", digest),
		"the base not paths'":  rules("../../run.jsonl", digest),
		"another criterion":    {Evaluation: "rules-1", Criterion: "model", Run: "../../run.jsonl", SHA256: digest, Date: "2026-11-17"},
		"fairness with a base": {Evaluation: "rules-1", Run: "../../run.jsonl", SHA256: digest, Base: "../../run.jsonl", BaseSHA256: digest, Date: "2026-11-17"},
	} {
		var stdout, stderr bytes.Buffer
		if got := run([]string{"evaluate", write(e)}, &stdout, &stderr); got != 1 {
			t.Errorf("%s: sim evaluate = %d: %s", name, got, stderr.String())
		}
	}
}

// An evaluation file of the correction names its run alone, and sim evaluate writes its report.
func TestEvaluateCorrection(t *testing.T) {
	dir := t.TempDir()
	folder := filepath.Join(dir, "evaluations", "rules-2")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Learners: criterion.Decisive, Invariants: true, Populations: []string{"D-grade1", "C-grade1"},
		Libraries: map[string]string{"github.com/altnovatrail/engine": "v0.1.6", "github.com/altnovatrail/learner-model": "v0.1.13-0.20261116235821-" + criterion.RulesCandidate[:12]}}}
	for _, p := range f.Header.Populations {
		for n := range criterion.Decisive {
			tr := &measure.Trace{Levels: make([]float64, measure.FairLength), Truth: make([]float64, measure.FairLength), Bands: []int{1}}
			f.Lines = append(f.Lines, bench.Line{Population: p, Graph: []string{population.Olympiad}, N: n, Sums: measure.Sums{"learners": 1}, Fair: tr})
		}
	}
	var b bytes.Buffer
	if err := bench.Write(f, &b); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "run.jsonl"), b.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(b.Bytes())
	write := func(e evaluation) string {
		body, _ := json.Marshal(e)
		path := filepath.Join(folder, "evaluation.json")
		if err := os.WriteFile(path, body, 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	e := evaluation{Evaluation: "rules-2", Criterion: "correction", Run: "../../run.jsonl", SHA256: hex.EncodeToString(sum[:]), Date: "2026-11-17"}
	var stdout, stderr bytes.Buffer
	if got := run([]string{"evaluate", write(e)}, &stdout, &stderr); got != 0 {
		t.Fatalf("sim evaluate = %d: %s", got, stderr.String())
	}
	if !strings.Contains(stdout.String(), "# The adult's correction against the truth: rules-2") {
		t.Errorf("the report:\n%s", stdout.String())
	}
	e.Base, e.BaseSHA256 = "../../run.jsonl", e.SHA256
	if got := run([]string{"evaluate", write(e)}, &stdout, &stderr); got != 1 {
		t.Errorf("a correction with a base: sim evaluate = %d", got)
	}
	e.Base, e.BaseSHA256, e.Date = "", "", "2026-11-17"
	f.Lines = f.Lines[1:]
	b.Reset()
	if err := bench.Write(f, &b); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "run.jsonl"), b.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	sum = sha256.Sum256(b.Bytes())
	e.SHA256 = hex.EncodeToString(sum[:])
	if got := run([]string{"evaluate", write(e)}, &stdout, &stderr); got != 1 {
		t.Errorf("a run short of a learner: sim evaluate = %d", got)
	}
}
