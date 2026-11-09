package layout

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "layout",
	Doc:  "types must live in type.go, constants in defines.go",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	var decl ast.Decl
	var gen *ast.GenDecl
	var ok bool
	var name string
	for _, file = range pass.Files {
		name = pass.Fset.File(file.Pos()).Name()
		if ast.IsGenerated(file) || strings.HasSuffix(name, "_test.go") {
			continue
		}
		name = filepath.Base(name)
		for _, decl = range file.Decls {
			gen, ok = decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if gen.Tok == token.TYPE && name != "type.go" {
				reportTypes(pass, gen)
			}
			if gen.Tok == token.CONST && name != "defines.go" {
				pass.Reportf(gen.Pos(), "constants must live in defines.go")
			}
		}
	}
	return nil, nil
}

func reportTypes(pass *analysis.Pass, gen *ast.GenDecl) {
	var spec ast.Spec
	var typeSpec *ast.TypeSpec
	var ok bool
	for _, spec = range gen.Specs {
		typeSpec, ok = spec.(*ast.TypeSpec)
		if ok {
			pass.Reportf(typeSpec.Pos(), "type %q must live in type.go", typeSpec.Name.Name)
		}
	}
}
