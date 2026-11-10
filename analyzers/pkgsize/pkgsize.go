package pkgsize

import (
	"go/ast"
	"go/token"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "pkgsize",
	Doc:  "a package must stay under MAX_TYPES (10) types and MAX_LINES (1500) lines; split it",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	first := token.NoPos
	typeCount, lineCount := 0, 0
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		if first == token.NoPos {
			first = file.Name.Pos()
		}
		lineCount += pass.Fset.File(file.Pos()).LineCount()
		typeCount += countTypes(file)
	}
	if typeCount > MAX_TYPES {
		kit.Reportf(pass, first, "package declares %d types (max %d); split the package", typeCount, MAX_TYPES)
	}
	if lineCount > MAX_LINES {
		kit.Reportf(pass, first, "package has %d lines (max %d); split the package", lineCount, MAX_LINES)
	}
	return nil, nil
}

func countTypes(file *ast.File) int {
	total := 0
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.TYPE {
			total += len(gen.Specs)
		}
	}
	return total
}
