package nakedreturn

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "nakedreturn",
	Doc:  "functions with named results must not use a bare return",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	for _, file = range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			var decl *ast.FuncDecl
			var lit *ast.FuncLit
			var ok bool
			decl, ok = node.(*ast.FuncDecl)
			if ok && decl.Body != nil {
				checkFunc(pass, decl.Type, decl.Body)
			}
			lit, ok = node.(*ast.FuncLit)
			if ok {
				checkFunc(pass, lit.Type, lit.Body)
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
		var isLit bool
		var ret *ast.ReturnStmt
		var ok bool
		_, isLit = node.(*ast.FuncLit)
		if isLit {
			return false
		}
		ret, ok = node.(*ast.ReturnStmt)
		if ok && len(ret.Results) == 0 {
			pass.Reportf(ret.Pos(), "bare return with named results; return the values explicitly")
		}
		return true
	})
}

func hasNamedResults(fnType *ast.FuncType) bool {
	var field *ast.Field
	if fnType.Results == nil {
		return false
	}
	for _, field = range fnType.Results.List {
		if len(field.Names) > 0 {
			return true
		}
	}
	return false
}
