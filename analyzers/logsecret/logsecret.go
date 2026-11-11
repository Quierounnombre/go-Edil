package logsecret

import (
	"go/ast"
	"go/types"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "logsecret",
	Doc:  "do not pass password, token or secret values to log, slog or fmt print calls",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if ok {
				check(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

func check(pass *analysis.Pass, call *ast.CallExpr) {
	fn, ok := kit.Callee(pass, call).(*types.Func)
	if !ok || fn.Pkg() == nil || !isLogging(fn) {
		return
	}
	structured := fn.Pkg().Path() == SLOG_PKG
	for _, arg := range call.Args {
		name := secretName(pass, arg, structured)
		if name != "" {
			kit.Reportf(pass, arg.Pos(), "possible secret %q in a logging call; log an identifier, not the secret", name)
		}
	}
}

func isLogging(fn *types.Func) bool {
	switch fn.Pkg().Path() {
	case LOG_PKG, SLOG_PKG:
		return true
	case FMT_PKG:
		return isFmtPrint(fn.Name())
	}
	return false
}

func isFmtPrint(name string) bool {
	switch name {
	case FMT_PRINT, FMT_PRINTF, FMT_PRINTLN, FMT_FPRINT, FMT_FPRINTF, FMT_FPRINTL, FMT_ERRORF:
		return true
	}
	return false
}

func secretName(pass *analysis.Pass, arg ast.Expr, structured bool) string {
	key, isConst := kit.ConstString(pass, arg)
	if isConst {
		return keyName(key, structured)
	}
	name := exprName(pass, arg)
	if kit.IsSecretName(name) {
		return name
	}
	return ""
}

func keyName(key string, structured bool) string {
	if structured && !strings.Contains(key, SPACE) && kit.IsSecretName(key) {
		return key
	}
	return ""
}

func exprName(pass *analysis.Pass, expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.SelectorExpr:
		return typed.Sel.Name
	case *ast.StarExpr:
		return exprName(pass, typed.X)
	case *ast.UnaryExpr:
		return exprName(pass, typed.X)
	case *ast.CallExpr:
		return callName(pass, typed)
	}
	return ""
}

func callName(pass *analysis.Pass, call *ast.CallExpr) string {
	if pass.TypesInfo.Types[call.Fun].IsType() && len(call.Args) == 1 {
		return exprName(pass, call.Args[0])
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok && len(call.Args) == 0 {
		return sel.Sel.Name
	}
	return ""
}
