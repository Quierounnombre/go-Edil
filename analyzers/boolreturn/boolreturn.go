package boolreturn

import (
	"go/ast"
	"go/token"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "boolreturn",
	Doc:  "flags are bool: no 0/1 returns, no comparing against true or false",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		checkComparisons(pass, file)
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Body != nil && returnsInteger(pass, fn) && onlyZeroOne(fn.Body) {
				kit.Reportf(pass, fn.Name.Pos(), "function %q returns only 0 and 1; return bool", fn.Name.Name)
			}
		}
	}
	return nil, nil
}

func checkComparisons(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		bin, ok := node.(*ast.BinaryExpr)
		if !ok || (bin.Op != token.EQL && bin.Op != token.NEQ) {
			return true
		}
		if isBoolConst(pass, bin.X) || isBoolConst(pass, bin.Y) {
			kit.Reportf(pass, bin.Pos(), "comparison with true/false; use the bool directly")
		}
		return true
	})
}

func isBoolConst(pass *analysis.Pass, expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	obj := pass.TypesInfo.Uses[ident]
	return obj == types.Universe.Lookup(TRUE_NAME) || obj == types.Universe.Lookup(FALSE_NAME)
}

func returnsInteger(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	results := fn.Type.Results
	if results == nil || len(results.List) != 1 || len(results.List[0].Names) > 1 {
		return false
	}
	basic, ok := pass.TypesInfo.TypeOf(results.List[0].Type).(*types.Basic)
	return ok && basic.Info()&types.IsInteger != 0
}

func onlyZeroOne(body *ast.BlockStmt) bool {
	sawZero, sawOne, sawOther := false, false, false
	ast.Inspect(body, func(node ast.Node) bool {
		_, isLit := node.(*ast.FuncLit)
		if isLit {
			return false
		}
		ret, ok := node.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		value := literalValue(ret)
		sawZero = sawZero || value == ZERO_LITERAL
		sawOne = sawOne || value == ONE_LITERAL
		sawOther = sawOther || (value != ZERO_LITERAL && value != ONE_LITERAL)
		return true
	})
	return sawZero && sawOne && !sawOther
}

func literalValue(ret *ast.ReturnStmt) string {
	if len(ret.Results) != 1 {
		return ""
	}
	lit, ok := ret.Results[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return ""
	}
	return lit.Value
}
