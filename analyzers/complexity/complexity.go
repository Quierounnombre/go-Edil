package complexity

import (
	"go/ast"
	"go/token"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "complexity",
	Doc:  "cognitive complexity per function must stay at or below MAX_COMPLEXITY (15)",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			total := scoreChildren(fn.Body, 0)
			if total > MAX_COMPLEXITY {
				kit.Reportf(pass, fn.Name.Pos(), "function %q has cognitive complexity %d (max %d)", fn.Name.Name, total, MAX_COMPLEXITY)
			}
		}
	}
	return nil, nil
}

func scoreChildren(parent ast.Node, nesting int) int {
	total := 0
	ast.Inspect(parent, func(child ast.Node) bool {
		if child == nil || child == parent {
			return true
		}
		added, descend := scoreNode(child, nesting)
		total += added
		return descend
	})
	return total
}

func scoreNode(node ast.Node, nesting int) (int, bool) {
	switch typed := node.(type) {
	case *ast.IfStmt:
		return scoreIf(typed, nesting, false), false
	case *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		return 1 + nesting + scoreChildren(node, nesting+1), false
	case *ast.FuncLit:
		return scoreChildren(typed.Body, nesting+1), false
	case *ast.BranchStmt:
		return branchScore(typed), false
	case *ast.BinaryExpr:
		return logicalScore(typed), !isLogical(typed)
	}
	return 0, true
}

func scoreIf(stmt *ast.IfStmt, nesting int, chained bool) int {
	total := logicalScore(stmt.Cond) + 1
	if !chained {
		total += nesting
	}
	total += scoreChildren(stmt.Body, nesting+1)
	switch tail := stmt.Else.(type) {
	case *ast.IfStmt:
		total += scoreIf(tail, nesting, true)
	case *ast.BlockStmt:
		total += 1 + scoreChildren(tail, nesting+1)
	}
	return total
}

func branchScore(stmt *ast.BranchStmt) int {
	score := 0
	if stmt.Label != nil {
		score++
	}
	return score
}

func isLogical(expr *ast.BinaryExpr) bool {
	return expr.Op == token.LAND || expr.Op == token.LOR
}

func logicalScore(expr ast.Expr) int {
	total := 0
	last := token.ILLEGAL
	for _, op := range flatten(expr, nil) {
		if op != last {
			total++
		}
		last = op
	}
	return total
}

func flatten(expr ast.Expr, ops []token.Token) []token.Token {
	paren, isParen := expr.(*ast.ParenExpr)
	if isParen {
		return flatten(paren.X, ops)
	}
	bin, isBin := expr.(*ast.BinaryExpr)
	if !isBin || !isLogical(bin) {
		return ops
	}
	ops = flatten(bin.X, ops)
	ops = append(ops, bin.Op)
	return flatten(bin.Y, ops)
}
