package noelse

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noelse",
	Doc:  "no else after a block that returns, breaks or continues",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			stmt, ok := node.(*ast.IfStmt)
			if ok && stmt.Else != nil && exits(stmt.Body) {
				kit.Reportf(pass, stmt.Else.Pos(), "else after a block that exits; drop the else and unindent it")
			}
			return true
		})
	}
	return nil, nil
}

func exits(body *ast.BlockStmt) bool {
	if len(body.List) == 0 {
		return false
	}
	switch body.List[len(body.List)-1].(type) {
	case *ast.ReturnStmt, *ast.BranchStmt:
		return true
	}
	return false
}
