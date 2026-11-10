package tlsinsecure

import (
	"go/ast"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "tlsinsecure",
	Doc:  "tls.Config.InsecureSkipVerify must not be set to true",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.KeyValueExpr:
				check(pass, typed.Key, typed.Value)
			case *ast.AssignStmt:
				checkAssign(pass, typed)
			}
			return true
		})
	}
	return nil, nil
}

func checkAssign(pass *analysis.Pass, assign *ast.AssignStmt) {
	if len(assign.Lhs) != len(assign.Rhs) {
		return
	}
	for index, lhs := range assign.Lhs {
		check(pass, lhs, assign.Rhs[index])
	}
}

func check(pass *analysis.Pass, target ast.Expr, value ast.Expr) {
	if !isInsecureField(pass, target) || !isTrue(pass, value) {
		return
	}
	kit.Reportf(pass, target.Pos(), "InsecureSkipVerify disables certificate checks; trust a CA instead")
}

func isInsecureField(pass *analysis.Pass, expr ast.Expr) bool {
	ident := fieldIdent(expr)
	if ident == nil {
		return false
	}
	field, ok := pass.TypesInfo.Uses[ident].(*types.Var)
	return ok && field.IsField() && field.Name() == INSECURE && field.Pkg() != nil && field.Pkg().Path() == TLS_PKG
}

func fieldIdent(expr ast.Expr) *ast.Ident {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed
	case *ast.SelectorExpr:
		return typed.Sel
	}
	return nil
}

func isTrue(pass *analysis.Pass, expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && pass.TypesInfo.Uses[ident] == types.Universe.Lookup(TRUE_NAME)
}
