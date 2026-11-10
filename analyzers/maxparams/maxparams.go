package maxparams

import (
	"go/ast"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "maxparams",
	Doc:  "functions take at most MAX_PARAMS (5) parameters, context.Context not counted; group the rest in a struct",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			total := countParams(pass, fn.Type.Params)
			if total > MAX_PARAMS {
				kit.Reportf(pass, fn.Name.Pos(), "function %q has %d parameters (max %d); group them in a struct", fn.Name.Name, total, MAX_PARAMS)
			}
		}
	}
	return nil, nil
}

func countParams(pass *analysis.Pass, params *ast.FieldList) int {
	total := 0
	if params == nil {
		return total
	}
	for _, field := range params.List {
		if isContext(pass.TypesInfo.TypeOf(field.Type)) {
			continue
		}
		total += max(len(field.Names), 1)
	}
	return total
}

func isContext(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok || named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == CONTEXT_PKG && named.Obj().Name() == CONTEXT_NAME
}
