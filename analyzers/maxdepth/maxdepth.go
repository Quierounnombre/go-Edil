package maxdepth

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "maxdepth",
	Doc:  "control flow must not nest deeper than MAX_DEPTH (3) levels",
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
				walk(pass, fn.Body, 0, map[ast.Node]bool{})
			}
		}
	}
	return nil, nil
}

func walk(pass *analysis.Pass, root ast.Node, depth int, chained map[ast.Node]bool) {
	ast.Inspect(root, func(child ast.Node) bool {
		markChained(child, chained)
		if child == nil || child == root {
			return true
		}
		return visit(pass, child, depth, chained)
	})
}

func visit(pass *analysis.Pass, child ast.Node, depth int, chained map[ast.Node]bool) bool {
	lit, isLit := child.(*ast.FuncLit)
	if isLit {
		walk(pass, lit.Body, 0, chained)
		return false
	}
	if !isNesting(child) {
		return true
	}
	if chained[child] {
		walk(pass, child, depth, chained)
		return false
	}
	if depth+1 > MAX_DEPTH {
		kit.Reportf(pass, child.Pos(), "nesting deeper than %d levels; return early or extract a function", MAX_DEPTH)
		return false
	}
	walk(pass, child, depth+1, chained)
	return false
}

func markChained(node ast.Node, chained map[ast.Node]bool) {
	stmt, ok := node.(*ast.IfStmt)
	if !ok || stmt.Else == nil {
		return
	}
	next, isIf := stmt.Else.(*ast.IfStmt)
	if isIf {
		chained[next] = true
	}
}

func isNesting(node ast.Node) bool {
	switch node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		return true
	}
	return false
}
