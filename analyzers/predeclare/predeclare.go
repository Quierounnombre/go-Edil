package predeclare

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "predeclare",
	Doc:  "variables must be declared with var, not := mid-function",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		var exempt = map[ast.Node]bool{}
		ast.Inspect(file, func(node ast.Node) bool {
			switch stmt := node.(type) {
			case *ast.IfStmt:
				exempt[stmt.Init] = true
			case *ast.ForStmt:
				exempt[stmt.Init] = true
			case *ast.SwitchStmt:
				exempt[stmt.Init] = true
			case *ast.TypeSwitchStmt:
				exempt[stmt.Init] = true
				exempt[stmt.Assign] = true
			case *ast.CommClause:
				exempt[stmt.Comm] = true
			case *ast.AssignStmt:
				if stmt.Tok == token.DEFINE && !exempt[stmt] {
					pass.Reportf(stmt.Pos(), "declare variables with var, not :=")
				}
			}
			return true
		})
	}
	return nil, nil
}
