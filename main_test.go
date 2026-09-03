package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGeneratorSupportsArbitraryReferenceKinds(t *testing.T) {
	types := parseTestTypes(t, `package fixture
type Item struct {
	UserID string `+"`ref:\"user:UserName\"`"+`
	UserName string
	TenantId string `+"`ref:\"tenant:TenantName\"`"+`
	TenantName string
}
type Page struct {
	List []Item `+"`ref:\"inline\"`"+`
}`)

	generator := &generator{types: types}
	if err := generator.walkStruct(types["Page"], "value", map[string]bool{"Page": true}); err != nil {
		t.Fatal(err)
	}
	collector := generator.collector.String()
	applier := generator.applier.String()
	for _, fragment := range []string{
		`kind == "user"`, `kind == "tenant"`, `.UserID`, `.TenantId`,
	} {
		if !strings.Contains(collector, fragment) {
			t.Fatalf("collector does not contain %q:\n%s", fragment, collector)
		}
	}
	for _, fragment := range []string{`.UserName =`, `.TenantName =`} {
		if !strings.Contains(applier, fragment) {
			t.Fatalf("applier does not contain %q:\n%s", fragment, applier)
		}
	}
}

func TestGeneratorRejectsReferenceWithoutKind(t *testing.T) {
	types := parseTestTypes(t, `package fixture
type Item struct {
	UserID string `+"`ref:\"UserName\"`"+`
	UserName string
}`)
	generator := &generator{types: types}
	err := generator.walkStruct(types["Item"], "value", map[string]bool{"Item": true})
	if err == nil || !strings.Contains(err.Error(), "kind:Destination") {
		t.Fatalf("walkStruct() error = %v", err)
	}
}

func parseTestTypes(t *testing.T, source string) map[string]*ast.StructType {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "fixture.go", source, 0)
	if err != nil {
		t.Fatal(err)
	}
	types := make(map[string]*ast.StructType)
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.TYPE {
			continue
		}
		for _, specification := range general.Specs {
			typeSpec := specification.(*ast.TypeSpec)
			if structure, ok := typeSpec.Type.(*ast.StructType); ok {
				types[typeSpec.Name.Name] = structure
			}
		}
	}
	return types
}
