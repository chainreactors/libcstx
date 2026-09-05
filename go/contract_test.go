package cstx

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

// This gate prevents the SDK from quietly growing another public graph or
// repository model beside the generated protobuf package.
func TestNoDuplicatedPublicModelTypes(t *testing.T) {
	banned := map[string]bool{
		"Node": true, "Edge": true, "Relationship": true, "GraphStats": true,
		"Delta": true, "ChangeSet": true, "Commit": true, "GraphDiff": true,
		"History": true, "HistoryEntry": true, "RepositorySync": true,
		"MissingPlan": true, "NodeFilter": true, "EdgeFilter": true,
		"RelationshipFilter": true, "QueryOptions": true, "Config": true,
	}
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if filepath.Ext(file) != ".go" || filepath.Base(file) == "contract_test.go" {
			continue
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		for _, declaration := range parsed.Decls {
			general, ok := declaration.(*ast.GenDecl)
			if !ok || general.Tok != token.TYPE {
				continue
			}
			for _, spec := range general.Specs {
				name := spec.(*ast.TypeSpec).Name.Name
				if banned[name] {
					t.Fatalf("%s declares duplicated public model %s; use cstxproto.%s", file, name, name)
				}
			}
		}
	}
}
