package noifinit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "noifinit",
	Doc:  "if and switch must not have an init statement; declare it on its own line",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	for _, file = range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			var ifStmt *ast.IfStmt
			var switchStmt *ast.SwitchStmt
			var ok bool
			ifStmt, ok = node.(*ast.IfStmt)
			if ok && ifStmt.Init != nil {
				pass.Reportf(ifStmt.Init.Pos(), "init statement in if; declare it on its own line before the if")
			}
			switchStmt, ok = node.(*ast.SwitchStmt)
			if ok && switchStmt.Init != nil {
				pass.Reportf(switchStmt.Init.Pos(), "init statement in switch; declare it on its own line before the switch")
			}
			return true
		})
	}
	return nil, nil
}
