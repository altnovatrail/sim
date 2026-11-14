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
// fairness by the grade alone, on the libraries of the embedded gate file, and gives its SHA-256: no
// population of fairness runs on a published package before T64 (R55).
func fairRun(t *testing.T, path string) string {
	t.Helper()
	g, err := gate.Read(bytes.NewReader(gateJSON))
	if err != nil {
		t.Fatal(err)
	}
	f := &bench.File{Header: bench.Header{Set: seed.Working, Rule: "engine", Libraries: map[string]string{}, Learners: criterion.Decisive, Invariants: true}}
	for p, lib := range g.Libraries {
		f.Header.Libraries[p] = lib.Version
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
// SHA-256 and a date, and its report (R55). None is there before T64, and then nothing is held.
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
