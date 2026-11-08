package seed

import (
	"encoding/json"
	"io/fs"
	"maps"
	"path"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/altnovatrail/spec"
	"github.com/altnovatrail/spec/types/formatv0"
	"go.yaml.in/yaml/v4"
)

// readVector reads a vector of spec's format/vectors into v: the YAML into plain values, and those
// through JSON into the types of formatv0, whose tags are JSON's. The vectors are embedded in
// spec's Go module, so no submodule is needed (R64 spec).
func readVector(t *testing.T, name string, v any) {
	t.Helper()
	data, err := fs.ReadFile(spec.Format, name)
	if err != nil {
		t.Fatal(err)
	}
	var plain any
	if err := yaml.Unmarshal(data, &plain); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	j, err := json.Marshal(plain)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	dec := json.NewDecoder(strings.NewReader(string(j)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
}

// vectors gives the names of the vectors of one kind, failing if there are none.
func vectors(t *testing.T, kind string) []string {
	t.Helper()
	names, err := fs.Glob(spec.Format, "vectors/"+kind+"/*.yaml")
	if err != nil || len(names) == 0 {
		t.Fatalf("no vectors/%s in spec: %v", kind, err)
	}
	return names
}

func unsigned(t *testing.T, u formatv0.Unsigned) uint64 {
	t.Helper()
	v, err := strconv.ParseUint(string(u), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func signed(t *testing.T, s formatv0.Signed) int64 {
	t.Helper()
	v, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// Every vector of the generator: the state a seed and names give, and each draw from that state.
func TestRandomVectors(t *testing.T) {
	for _, name := range vectors(t, "random") {
		t.Run(path.Base(name), func(t *testing.T) {
			var v formatv0.RandomVector
			readVector(t, name, &v)
			high, low := unsigned(t, v.State.High), unsigned(t, v.State.Low)
			if v.Seed != nil {
				s := unsigned(t, *v.Seed)
				got, err := High(s, v.Names...)
				if err != nil {
					t.Fatal(err)
				}
				if got != high || s != low {
					t.Fatalf("the place of %d %q has the state %d, %d; want %d, %d", s, v.Names, got, s, high, low)
				}
				p, err := Place(s, v.Names...)
				if err != nil {
					t.Fatal(err)
				}
				if p.Integer() != New(high, low).Integer() {
					t.Fatal("Place and New of the same state draw apart")
				}
			}
			st := New(high, low)
			for i, d := range v.Draws {
				switch {
				case d.Integer != nil:
					if got, want := st.Integer(), unsigned(t, d.Integer.Value); got != want {
						t.Fatalf("draw %d: integer %d, want %d", i, got, want)
					}
				case d.Below != nil:
					if got, want := st.Below(unsigned(t, d.Below.Bound)), unsigned(t, d.Below.Value); got != want {
						t.Fatalf("draw %d: below %d, want %d", i, got, want)
					}
				case d.Between != nil:
					got := st.Between(signed(t, d.Between.From), signed(t, d.Between.To))
					if want := signed(t, d.Between.Value); got != want {
						t.Fatalf("draw %d: between %d, want %d", i, got, want)
					}
				case d.Fraction != nil:
					if got, want := st.Fraction(), d.Fraction.Value; got != want {
						t.Fatalf("draw %d: fraction %v, want %v", i, got, want)
					}
				case d.Permutation != nil:
					got := st.Permutation(int(d.Permutation.Count))
					want := make([]int, len(d.Permutation.Value))
					for k, x := range d.Permutation.Value {
						want[k] = int(x)
					}
					if !slices.Equal(got, want) {
						t.Fatalf("draw %d: permutation %v, want %v", i, got, want)
					}
				default:
					t.Fatalf("draw %d: no function", i)
				}
			}
		})
	}
}

// Every vector of the templates: the instance an attempt's seed takes is the one whose numbers the
// vector gives. The instances are the combinations of the parameters in the order of an odometer —
// names by code units, the first slowest, each domain upwards — that the constraint keeps (the
// text of the format, section 23). sim reads the instances from the canonical view and evaluates
// no expression; this test evaluates only the constraint of the form "p != q" that the vectors
// have, and fails on any other, so a new vector is read before it is trusted.
func TestTemplateVectors(t *testing.T) {
	for _, name := range vectors(t, "templates") {
		t.Run(path.Base(name), func(t *testing.T) {
			var v formatv0.TemplateVector
			readVector(t, name, &v)
			tpl := v.Task.Template
			if tpl == nil {
				t.Fatal("the vector's task is no template")
			}
			instances := combinations(t, tpl)
			if int64(len(instances)) != v.Count {
				t.Fatalf("%d instances, the vector says %d", len(instances), v.Count)
			}
			for _, a := range v.Attempts {
				i, err := Instance(unsigned(t, a.Seed), string(v.Task.ID), len(instances))
				if err != nil {
					t.Fatal(err)
				}
				for p, x := range instances[i] {
					if y, ok := a.Numbers[p]; !ok || y != x {
						t.Fatalf("seed %s: instance %d has %s = %d, the vector %d", a.Seed, i, p, x, a.Numbers[p])
					}
				}
			}
		})
	}
}

// combinations gives the instances of a template as values of its parameters.
func combinations(t *testing.T, tpl *formatv0.Template) []map[formatv0.TemplateName]int64 {
	t.Helper()
	names := make([]formatv0.TemplateName, 0, len(tpl.Parameters))
	for n := range tpl.Parameters {
		names = append(names, n)
	}
	slices.Sort(names)
	var keep func(map[formatv0.TemplateName]int64) bool
	if tpl.Constraint == nil {
		keep = func(map[formatv0.TemplateName]int64) bool { return true }
	} else {
		sides := strings.Split(string(*tpl.Constraint), " != ")
		if len(sides) != 2 {
			t.Fatalf("the test reads no constraint %q", *tpl.Constraint)
		}
		p, q := formatv0.TemplateName(sides[0]), formatv0.TemplateName(sides[1])
		_, okP := tpl.Parameters[p]
		_, okQ := tpl.Parameters[q]
		if !okP || !okQ {
			t.Fatalf("the test reads no constraint %q: both sides must be parameters", *tpl.Constraint)
		}
		keep = func(c map[formatv0.TemplateName]int64) bool { return c[p] != c[q] }
	}
	var out []map[formatv0.TemplateName]int64
	var walk func(int, map[formatv0.TemplateName]int64)
	walk = func(k int, c map[formatv0.TemplateName]int64) {
		if k == len(names) {
			if keep(c) {
				out = append(out, maps.Clone(c))
			}
			return
		}
		d := tpl.Parameters[names[k]]
		for x := d.From; x <= d.To; x++ {
			c[names[k]] = x
			walk(k+1, c)
		}
	}
	walk(0, map[formatv0.TemplateName]int64{})
	return out
}

func TestInstanceRefuses(t *testing.T) {
	if _, err := Instance(1, "a-task", 0); err != ErrNoInstance {
		t.Errorf("Instance with no instances: %v, want ErrNoInstance", err)
	}
	if _, err := Instance(1, "a|task", 2); err != ErrName {
		t.Errorf("Instance of a task ID with |: %v, want ErrName", err)
	}
}
