package managed_test

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This snapshot was extracted from the pinned upstream commit de6914c,
// independently of the Qoder implementation. Every upstream public field must
// remain available, including deeply nested union variants and path parameters.
func TestQoderFieldAlignment(t *testing.T) {
	b, e := os.ReadFile("testdata/upstream-fields.json")
	if e != nil {
		t.Fatal(e)
	}
	var expected map[string][]string
	if e = json.Unmarshal(b, &expected); e != nil {
		t.Fatal(e)
	}
	actual := map[string]map[string]bool{}
	files, e := filepath.Glob("*.go")
	if e != nil {
		t.Fatal(e)
	}
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, e := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if e != nil {
			t.Fatal(e)
		}
		for _, d := range f.Decls {
			g, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, s := range g.Specs {
				ts, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if strings.HasSuffix(ts.Name.Name, "Service") {
					if _, ok := ts.Type.(*ast.StructType); !ok {
						t.Errorf("service %s must be a concrete struct", ts.Name.Name)
					}
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				names := map[string]bool{}
				for _, f := range st.Fields.List {
					for _, n := range f.Names {
						names[n.Name] = true
					}
				}
				actual[ts.Name.Name] = names
			}
		}
	}
	for name, fields := range expected {
		for _, f := range fields {
			if !actual[name][f] {
				t.Errorf("upstream field missing: %s.%s", name, f)
			}
		}
	}
}
