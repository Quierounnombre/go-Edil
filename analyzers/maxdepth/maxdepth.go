package maxdepth

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "maxdepth",
	Doc:  "control flow must not nest deeper than MAX_DEPTH levels",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	var decl ast.Decl
	var fn *ast.FuncDecl
	var ok bool
	for _, file = range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		for _, decl = range file.Decls {
			fn, ok = decl.(*ast.FuncDecl)
			if ok && fn.Body != nil {
				walk(pass, fn.Body, 0, map[ast.Node]bool{})
			}
		}
	}
	return nil, nil
}

func walk(pass *analysis.Pass, root ast.Node, depth int, chained map[ast.Node]bool) {
	ast.Inspect(root, func(child ast.Node) bool {
		var lit *ast.FuncLit
		var ok bool
		markChained(child, chained)
		if child == root {
			return true
		}
		lit, ok = child.(*ast.FuncLit)
		if ok {
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
			pass.Reportf(child.Pos(), "nesting deeper than %d levels; return early or extract a function", MAX_DEPTH)
			return false
		}
		walk(pass, child, depth+1, chained)
		return false
	})
}

func markChained(node ast.Node, chained map[ast.Node]bool) {
	var stmt *ast.IfStmt
	var next *ast.IfStmt
	var ok bool
	stmt, ok = node.(*ast.IfStmt)
	if !ok || stmt.Else == nil {
		return
	}
	next, ok = stmt.Else.(*ast.IfStmt)
	if ok {
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
