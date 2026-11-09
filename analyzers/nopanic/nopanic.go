package nopanic

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "nopanic",
	Doc:  "forbid panic, log.Fatal*, log.Panic* and os.Exit outside main; return errors instead",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	for _, file = range pass.Files {
		if strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			continue
		}
		checkFile(pass, file)
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File) {
	var decl ast.Decl
	var fn *ast.FuncDecl
	var ok bool
	var inMain bool
	for _, decl = range file.Decls {
		fn, ok = decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		inMain = pass.Pkg.Name() == "main" && fn.Recv == nil && fn.Name.Name == "main"
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			var call *ast.CallExpr
			var isCall bool
			var reason string
			call, isCall = node.(*ast.CallExpr)
			if !isCall {
				return true
			}
			reason = forbidden(pass, call, inMain)
			if reason != "" {
				pass.Reportf(call.Pos(), "%s; return an error instead", reason)
			}
			return true
		})
	}
}

func forbidden(pass *analysis.Pass, call *ast.CallExpr, inMain bool) string {
	var obj types.Object
	var builtin *types.Builtin
	var fn *types.Func
	var ok bool
	obj = calleeObject(pass, call.Fun)
	builtin, ok = obj.(*types.Builtin)
	if ok && builtin.Name() == "panic" {
		return "panic is forbidden"
	}
	fn, ok = obj.(*types.Func)
	if !ok || fn.Pkg() == nil {
		return ""
	}
	switch fn.Pkg().Path() + "." + fn.Name() {
	case "log.Panic", "log.Panicf", "log.Panicln":
		return fn.Name() + " is forbidden"
	case "os.Exit", "log.Fatal", "log.Fatalf", "log.Fatalln":
		if !inMain {
			return fn.Name() + " is only allowed in main"
		}
	}
	return ""
}

func calleeObject(pass *analysis.Pass, expr ast.Expr) types.Object {
	var ident *ast.Ident
	var sel *ast.SelectorExpr
	var ok bool
	ident, ok = expr.(*ast.Ident)
	if ok {
		return pass.TypesInfo.Uses[ident]
	}
	sel, ok = expr.(*ast.SelectorExpr)
	if ok {
		return pass.TypesInfo.Uses[sel.Sel]
	}
	return nil
}
