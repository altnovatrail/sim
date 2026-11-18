package measure

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// partEnv names the part a process of the test binary runs, for engineSumsOn: "<k>/<n> <file>".
const partEnv = "SIM_MEASURE_PART"

// engineSumsOn runs n learners of a population on its g-th graph, as recordsOn does, and gives the
// sums of the engine's measures of each, in the learners' order, and the graph. Where the test
// binary can start processes, the learners go to as many processes of it as the processors, the k-th
// of n those whose number gives k mod n, as the recipes split a run (R09, R67): each process runs
// the test alone, writes its learners' sums into a file, in their order, and skips, and the test
// reads them. A learner comes from its seeds alone, so its sums are the same in any process. In
// js/wasm and wasip1, which start no process, and when the run counts coverage, which a process of
// its own would not add to the profile, the learners run here, one after another.
func engineSumsOn(t *testing.T, name string, g, n int) ([]Sums, []string) {
	t.Helper()
	if v, ok := os.LookupEnv(partEnv); ok {
		// The file is all after the first space, so that a folder with a space in its name is whole.
		part, file, _ := strings.Cut(v, " ")
		var k, of int
		if _, err := fmt.Sscanf(part, "%d/%d", &k, &of); err != nil || of < 1 || k < 0 || k >= of || file == "" {
			t.Fatalf("%s is %q, not <k>/<n> <file>", partEnv, v)
		}
		var ks []int
		for i := k; i < n; i += of {
			ks = append(ks, i)
		}
		recs, _ := recordsOf(t, name, g, ks)
		sums := make([]Sums, len(recs))
		for i, rec := range recs {
			sums[i] = Engine(rec)
		}
		b, err := json.Marshal(sums)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, b, 0o600); err != nil {
			t.Fatal(err)
		}
		t.Skipf("part %d of %d of %s, written into %s", k, of, name, file)
	}
	procs := min(runtime.NumCPU(), n)
	if runtime.GOOS == "js" || runtime.GOOS == "wasip1" || testing.CoverMode() != "" || procs < 2 {
		recs, graph := recordsOn(t, name, g, n)
		sums := make([]Sums, len(recs))
		for i, rec := range recs {
			sums[i] = Engine(rec)
		}
		return sums, graph
	}
	// The pattern of -test.run that matches this test alone: each level of its name, anchored, which
	// holds only for a name with no character a regular expression reads.
	levels := strings.Split(t.Name(), "/")
	for i, l := range levels {
		if strings.ContainsAny(l, `\.+*?()|[]{}^$`) {
			t.Fatalf("the name %q is read as a pattern by -test.run", t.Name())
		}
		levels[i] = "^" + l + "$"
	}
	// The processes stop at the test binary's own time limit, so that none outlives it.
	timeout := "0"
	if f := flag.Lookup("test.timeout"); f != nil {
		timeout = f.Value.String()
	}
	dir := t.TempDir()
	file := func(k int, ext string) string { return filepath.Join(dir, fmt.Sprintf("part-%d.%s", k, ext)) }
	// The output of a process goes into a file, which it writes itself: a writer of another kind would
	// need a goroutine of os/exec to copy it, and the tests start none (R10).
	cmds := make([]*exec.Cmd, procs)
	for k := range procs {
		out, err := os.Create(file(k, "out"))
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(os.Args[0], "-test.run="+strings.Join(levels, "/"), "-test.count=1", "-test.timeout="+timeout)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%d/%d %s", partEnv, k, procs, file(k, "json")))
		cmd.Stdout = out
		cmd.Stderr = out
		err = cmd.Start()
		out.Close()
		if err != nil {
			for _, c := range cmds[:k] {
				_ = c.Process.Kill()
				_ = c.Wait()
			}
			t.Fatal(err)
		}
		cmds[k] = cmd
	}
	failed := false
	for k, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			out, _ := os.ReadFile(file(k, "out"))
			t.Errorf("part %d of %d of %s: %v\n%s", k, procs, name, err, out)
			failed = true
		}
	}
	if failed {
		t.FailNow()
	}
	sums := make([]Sums, n)
	for k := range procs {
		b, err := os.ReadFile(file(k, "json"))
		if err != nil {
			t.Fatal(err)
		}
		// The part gives its learners k, k + procs, … in that order.
		var part []Sums
		if err := json.Unmarshal(b, &part); err != nil {
			t.Fatal(err)
		}
		if want := (n - k + procs - 1) / procs; len(part) != want {
			t.Fatalf("part %d of %d of %s gives %d learners, not %d", k, procs, name, len(part), want)
		}
		for j, s := range part {
			sums[k+j*procs] = s
		}
	}
	_, graph := recordsOf(t, name, g, nil)
	return sums, graph
}
