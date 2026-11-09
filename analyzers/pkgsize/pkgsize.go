package pkgsize

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "pkgsize",
	Doc:  "a package must stay under MAX_TYPES types and MAX_LINES lines; split it",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	var decl ast.Decl
	var gen *ast.GenDecl
	var tokFile *token.File
	var ok bool
	var first token.Pos
	var typeCount int
	var lineCount int
	for _, file = range pass.Files {
		tokFile = pass.Fset.File(file.Pos())
		if ast.IsGenerated(file) || strings.HasSuffix(tokFile.Name(), "_test.go") {
			continue
		}
		if first == token.NoPos {
			first = file.Name.Pos()
		}
		lineCount += tokFile.LineCount()
		for _, decl = range file.Decls {
			gen, ok = decl.(*ast.GenDecl)
			if ok && gen.Tok == token.TYPE {
				typeCount += len(gen.Specs)
			}
		}
	}
	if typeCount > MAX_TYPES {
		pass.Reportf(first, "package declares %d types (max %d); split the package", typeCount, MAX_TYPES)
	}
	if lineCount > MAX_LINES {
		pass.Reportf(first, "package has %d lines (max %d); split the package", lineCount, MAX_LINES)
	}
	return nil, nil
}
