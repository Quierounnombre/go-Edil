package noifinit

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noifinit",
	Doc:  "if and switch must not have an init statement; declare it on its own line",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.IfStmt:
				report(pass, typed.Init, "if")
			case *ast.SwitchStmt:
				report(pass, typed.Init, "switch")
			}
			return true
		})
	}
	return nil, nil
}

func report(pass *analysis.Pass, stmt ast.Stmt, keyword string) {
	if stmt != nil {
		kit.Reportf(pass, stmt.Pos(), "init statement in %s; declare it on its own line before the %s", keyword, keyword)
	}
}
