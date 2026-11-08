package synthetic

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"path"
	"testing"

	"github.com/altnovatrail/spec"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// canonicalView compiles the schema of the canonical view of spec's module, with every schema of
// the format beside it by its $id, as module-kit checks a package.
func canonicalView(t *testing.T) *jsonschema.Schema {
	t.Helper()
	c := jsonschema.NewCompiler()
	names, err := fs.Glob(spec.Format, "schemas/*.schema.json")
	if err != nil || len(names) == 0 {
		t.Fatalf("no schemas in spec: %v", err)
	}
	for _, name := range names {
		data, err := fs.ReadFile(spec.Format, name)
		if err != nil {
			t.Fatal(err)
		}
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if err := c.AddResource("https://altnovatrail.github.io/spec/format/0.3/"+path.Base(name), doc); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	s, err := c.Compile("https://altnovatrail.github.io/spec/format/0.3/canonical-view.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// Every synthetic view is a valid canonical view of the format 0.3: the schema spec gives it holds
// what the Go types cannot — required members, enums, the shared taxonomy (the text of the format,
// section 25).
func TestSchema(t *testing.T) {
	s := canonicalView(t)
	for id, view := range views {
		data, err := json.Marshal(view())
		if err != nil {
			t.Fatal(err)
		}
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		if err := s.Validate(doc); err != nil {
			t.Errorf("%s is no valid canonical view: %v", id, err)
		}
	}
}
