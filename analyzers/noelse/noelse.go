package noelse

import (
	"go/ast"
	"go/token"

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
				report(pass, file, stmt)
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

func report(pass *analysis.Pass, file *ast.File, stmt *ast.IfStmt) {
	message := "else after a block that exits; drop the else and unindent it"
	if !flattenable(file, stmt) {
		kit.Report(pass, stmt.If, message)
		return
	}
	edits := []analysis.TextEdit{kit.Edit(stmt.Body.End(), stmt.Else.Pos()+token.Pos(len(token.LBRACE.String())), "")}
	block, isBlock := stmt.Else.(*ast.BlockStmt)
	if !isBlock {
		edits = []analysis.TextEdit{kit.Edit(stmt.Body.End(), stmt.Else.Pos(), NEWLINE)}
	} else {
		edits = append(edits, kit.Edit(block.List[len(block.List)-1].End(), block.End(), ""))
	}
	kit.ReportFix(pass, stmt.If, message, "remove the else", edits)
}

// flattenable is false when moving the else body out would leak its
// declarations into the parent scope or lose the scope of an init statement.
func flattenable(file *ast.File, stmt *ast.IfStmt) bool {
	if stmt.Init != nil {
		return false
	}
	block, isBlock := stmt.Else.(*ast.BlockStmt)
	if !isBlock {
		return true
	}
	if len(block.List) == 0 || hasComment(file, stmt.Else.Pos(), stmt.Else.End()) {
		return false
	}
	for _, inner := range block.List {
		if declares(inner) {
			return false
		}
	}
	return true
}

func declares(stmt ast.Stmt) bool {
	switch typed := stmt.(type) {
	case *ast.DeclStmt, *ast.LabeledStmt:
		return true
	case *ast.AssignStmt:
		return typed.Tok == token.DEFINE
	}
	return false
}

// hasComment guards the fix: deleting text must never delete a comment.
func hasComment(file *ast.File, from token.Pos, to token.Pos) bool {
	for _, group := range file.Comments {
		if group.End() > from && group.Pos() < to {
			return true
		}
	}
	return false
}
