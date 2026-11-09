package noelse

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "noelse",
	Doc:  "no else after a block that returns, breaks or continues",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	for _, file = range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			var stmt *ast.IfStmt
			var ok bool
			stmt, ok = node.(*ast.IfStmt)
			if ok && stmt.Else != nil && exits(stmt.Body) {
				pass.Reportf(stmt.Else.Pos(), "else after a block that exits; drop the else and unindent it")
			}
			return true
		})
	}
	return nil, nil
}

func exits(body *ast.BlockStmt) bool {
	var last ast.Stmt
	var isReturn bool
	var isBranch bool
	if len(body.List) == 0 {
		return false
	}
	last = body.List[len(body.List)-1]
	_, isReturn = last.(*ast.ReturnStmt)
	_, isBranch = last.(*ast.BranchStmt)
	return isReturn || isBranch
}
