package funclen

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "funclen",
	Doc:  "functions must not exceed MAX_FUNC_LINES (120) lines",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			first := pass.Fset.Position(fn.Pos()).Line
			last := pass.Fset.Position(fn.End()).Line
			total := last - first + 1
			if total > MAX_FUNC_LINES {
				kit.Reportf(pass, fn.Name.Pos(), "function %q is %d lines (max %d)", fn.Name.Name, total, MAX_FUNC_LINES)
			}
		}
	}
	return nil, nil
}
