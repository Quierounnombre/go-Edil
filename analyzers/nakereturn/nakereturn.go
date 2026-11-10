package nakedreturn

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nakedreturn",
	Doc:  "functions with named results must not use a bare return",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.FuncDecl:
				if typed.Body != nil {
					checkFunc(pass, typed.Type, typed.Body)
				}
			case *ast.FuncLit:
				checkFunc(pass, typed.Type, typed.Body)
			}
			return true
		})
	}
	return nil, nil
}

func checkFunc(pass *analysis.Pass, fnType *ast.FuncType, body *ast.BlockStmt) {
	if !hasNamedResults(fnType) {
		return
	}
	ast.Inspect(body, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt:
			if len(typed.Results) == 0 {
				kit.Reportf(pass, typed.Pos(), "bare return with named results; return the values explicitly")
			}
		}
		return true
	})
}

func hasNamedResults(fnType *ast.FuncType) bool {
	if fnType.Results == nil {
		return false
	}
	for _, field := range fnType.Results.List {
		if len(field.Names) > 0 {
			return true
		}
	}
	return false
}
