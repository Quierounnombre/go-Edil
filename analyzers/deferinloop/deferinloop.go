package deferinloop

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "deferinloop",
	Doc:  "no defer inside a loop; it runs when the function returns, not per iteration",
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
				scan(pass, fn.Body, false)
			}
		}
	}
	return nil, nil
}

func scan(pass *analysis.Pass, root ast.Node, inLoop bool) {
	ast.Inspect(root, func(node ast.Node) bool {
		if node == nil || node == root {
			return true
		}
		switch typed := node.(type) {
		case *ast.FuncLit:
			scan(pass, typed.Body, false)
			return false
		case *ast.ForStmt:
			scan(pass, typed.Body, true)
			return false
		case *ast.RangeStmt:
			scan(pass, typed.Body, true)
			return false
		case *ast.DeferStmt:
			if inLoop {
				kit.Reportf(pass, typed.Pos(), "defer inside a loop runs at function exit; move the body into its own function")
			}
		}
		return true
	})
}
