package layout

import (
	"go/ast"
	"go/token"
	"path/filepath"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "layout",
	Doc:  "types must live in type.go, constants in defines.go",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		name := filepath.Base(pass.Fset.File(file.Pos()).Name())
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if ok {
				check(pass, gen, name)
			}
		}
	}
	return nil, nil
}

func check(pass *analysis.Pass, gen *ast.GenDecl, name string) {
	if gen.Tok == token.TYPE && name != TYPE_FILE {
		reportTypes(pass, gen)
	}
	if gen.Tok == token.CONST && name != DEFINES_FILE {
		kit.Reportf(pass, gen.Pos(), "constants must live in defines.go")
	}
}

func reportTypes(pass *analysis.Pass, gen *ast.GenDecl) {
	for _, spec := range gen.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if ok {
			kit.Reportf(pass, typeSpec.Pos(), "type %q must live in type.go", typeSpec.Name.Name)
		}
	}
}
