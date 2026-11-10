package varinit

import (
	"go/ast"
	"go/token"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "varinit",
	Doc:  "inside functions use :=; var only for zero values that are not assigned right after",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.BlockStmt:
				checkStmts(pass, typed.List)
			case *ast.CaseClause:
				checkStmts(pass, typed.Body)
			case *ast.CommClause:
				checkStmts(pass, typed.Body)
			}
			return true
		})
	}
	return nil, nil
}

func checkStmts(pass *analysis.Pass, stmts []ast.Stmt) {
	for index, stmt := range stmts {
		decl, ok := stmt.(*ast.DeclStmt)
		if !ok {
			continue
		}
		gen, ok := decl.Decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			checkSpec(pass, spec.(*ast.ValueSpec), stmts[index+1:])
		}
	}
}

func checkSpec(pass *analysis.Pass, spec *ast.ValueSpec, rest []ast.Stmt) {
	if len(spec.Values) > 0 && spec.Type == nil {
		kit.Reportf(pass, spec.Pos(), "use := instead of var with an initial value")
		return
	}
	if len(spec.Values) > 0 {
		return
	}
	for _, name := range spec.Names {
		obj := pass.TypesInfo.Defs[name]
		if obj != nil && assignedFirst(pass, rest, obj) {
			kit.Reportf(pass, name.Pos(), "%q is assigned right after its declaration; use := at that assignment", name.Name)
		}
	}
}

func assignedFirst(pass *analysis.Pass, rest []ast.Stmt, obj types.Object) bool {
	for _, stmt := range rest {
		if !mentions(pass, stmt, obj) {
			continue
		}
		assign, ok := stmt.(*ast.AssignStmt)
		return ok && assign.Tok == token.ASSIGN && onlyTarget(pass, assign, obj)
	}
	return false
}

func onlyTarget(pass *analysis.Pass, assign *ast.AssignStmt, obj types.Object) bool {
	targeted := false
	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		targeted = targeted || (ok && pass.TypesInfo.Uses[ident] == obj)
	}
	for _, rhs := range assign.Rhs {
		if mentions(pass, rhs, obj) {
			return false
		}
	}
	return targeted
}

func mentions(pass *analysis.Pass, root ast.Node, obj types.Object) bool {
	found := false
	ast.Inspect(root, func(node ast.Node) bool {
		ident, ok := node.(*ast.Ident)
		found = found || (ok && pass.TypesInfo.Uses[ident] == obj)
		return !found
	})
	return found
}
