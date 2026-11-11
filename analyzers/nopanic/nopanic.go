package nopanic

import (
	"go/ast"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nopanic",
	Doc:  "forbid panic, log.Fatal*, log.Panic* and os.Exit outside main; return errors instead",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Body != nil {
				checkBody(pass, fn.Body, kit.IsMain(pass, fn))
			}
		}
	}
	return nil, nil
}

func checkBody(pass *analysis.Pass, body *ast.BlockStmt, inMain bool) {
	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		reason := forbidden(kit.Callee(pass, call), inMain)
		if reason != "" {
			kit.Reportf(pass, call.Pos(), "%s; return an error instead", reason)
		}
		return true
	})
}

func forbidden(obj types.Object, inMain bool) string {
	_, isBuiltin := obj.(*types.Builtin)
	if isBuiltin && obj.Name() == PANIC_NAME {
		return "panic is forbidden"
	}
	switch kit.Qualified(obj) {
	case LOG_PANIC, LOG_PANICF, LOG_PANICLN:
		return obj.Name() + " is forbidden"
	case OS_EXIT, LOG_FATAL, LOG_FATALF, LOG_FATALLN:
		if !inMain {
			return obj.Name() + " is only allowed in main"
		}
	}
	return ""
}
