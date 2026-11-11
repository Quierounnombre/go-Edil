package varinit

import (
	"fmt"
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
			checkSpec(pass, decl, spec.(*ast.ValueSpec), stmts[index+1:])
		}
	}
}

func checkSpec(pass *analysis.Pass, decl *ast.DeclStmt, spec *ast.ValueSpec, rest []ast.Stmt) {
	if len(spec.Values) > 0 && spec.Type == nil {
		reportInitialValue(pass, decl, spec)
		return
	}
	if len(spec.Values) > 0 {
		return
	}
	for _, name := range spec.Names {
		obj := pass.TypesInfo.Defs[name]
		assign := firstAssign(pass, rest, obj)
		if obj == nil || assign == nil {
			continue
		}
		message := fmt.Sprintf("%q is assigned right after its declaration; use := at that assignment", name.Name)
		if canMerge(pass, decl, spec, assign, obj) {
			kit.ReportFix(pass, name.Pos(), message, "declare at the assignment", mergeEdits(decl, rest, assign))
			continue
		}
		kit.Report(pass, name.Pos(), message)
	}
}

func firstAssign(pass *analysis.Pass, rest []ast.Stmt, obj types.Object) *ast.AssignStmt {
	for _, stmt := range rest {
		if !mentions(pass, stmt, obj) {
			continue
		}
		assign, ok := stmt.(*ast.AssignStmt)
		if ok && assign.Tok == token.ASSIGN && onlyTarget(pass, assign, obj) {
			return assign
		}
		return nil
	}
	return nil
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

func reportInitialValue(pass *analysis.Pass, decl *ast.DeclStmt, spec *ast.ValueSpec) {
	message := "use := instead of var with an initial value"
	gen := decl.Decl.(*ast.GenDecl)
	if gen.Lparen.IsValid() || len(spec.Names) != len(spec.Values) {
		kit.Report(pass, spec.Pos(), message)
		return
	}
	last := spec.Names[len(spec.Names)-1]
	edits := []analysis.TextEdit{
		kit.Edit(gen.Pos(), spec.Names[0].Pos(), ""),
		kit.Edit(last.End(), spec.Values[0].Pos(), " := "),
	}
	kit.ReportFix(pass, spec.Pos(), message, "use :=", edits)
}

// canMerge is true when "var x T; ...; x = e" can become "x := e" without
// changing the type of x or its scope.
func canMerge(pass *analysis.Pass, decl *ast.DeclStmt, spec *ast.ValueSpec, assign *ast.AssignStmt, obj types.Object) bool {
	gen := decl.Decl.(*ast.GenDecl)
	if gen.Lparen.IsValid() || len(spec.Names) != SINGLE || len(assign.Lhs) != SINGLE || len(assign.Rhs) != SINGLE {
		return false
	}
	value := pass.TypesInfo.Types[assign.Rhs[0]]
	return value.Value == nil && !value.IsNil() && types.Identical(obj.Type(), value.Type)
}

func mergeEdits(decl *ast.DeclStmt, rest []ast.Stmt, assign *ast.AssignStmt) []analysis.TextEdit {
	return []analysis.TextEdit{
		kit.Edit(decl.Pos(), rest[0].Pos(), ""),
		kit.Edit(assign.TokPos, assign.TokPos+token.Pos(len(token.ASSIGN.String())), token.DEFINE.String()),
	}
}
