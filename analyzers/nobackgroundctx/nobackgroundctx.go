package nobackgroundctx

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nobackgroundctx",
	Doc:  "context.Background and context.TODO only in main; accept a ctx parameter everywhere else",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Body != nil && !isMain(pass, fn) {
				checkBody(pass, fn.Body)
			}
		}
	}
	return nil, nil
}

func isMain(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	return pass.Pkg.Name() == MAIN_PKG && fn.Recv == nil && fn.Name.Name == MAIN_FUNC
}

func checkBody(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch kit.Qualified(kit.Callee(pass, call)) {
		case CTX_BACKGROUND, CTX_TODO:
			kit.Reportf(pass, call.Pos(), "root context created outside main; accept a ctx parameter")
		}
		return true
	})
}
