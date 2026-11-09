package errreturn

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "errreturn",
	Doc:  "errors must be returned, not assigned to _ or checked and swallowed",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var file *ast.File
	var decl ast.Decl
	var fn *ast.FuncDecl
	var ok bool
	for _, file = range pass.Files {
		if ast.IsGenerated(file) || strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			continue
		}
		for _, decl = range file.Decls {
			fn, ok = decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			checkBlank(pass, fn.Body)
			if returnsError(pass, fn) {
				checkSwallowed(pass, fn.Body)
			}
		}
	}
	return nil, nil
}

func isError(typ types.Type) bool {
	return typ != nil && types.Identical(typ, types.Universe.Lookup("error").Type())
}

func returnsError(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	var field *ast.Field
	if fn.Type.Results == nil {
		return false
	}
	for _, field = range fn.Type.Results.List {
		if isError(pass.TypesInfo.TypeOf(field.Type)) {
			return true
		}
	}
	return false
}

func checkBlank(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		var assign *ast.AssignStmt
		var isDefer bool
		var ok bool
		_, isDefer = node.(*ast.DeferStmt)
		if isDefer {
			return false
		}
		assign, ok = node.(*ast.AssignStmt)
		if ok {
			reportBlank(pass, assign)
		}
		return true
	})
}

func reportBlank(pass *analysis.Pass, assign *ast.AssignStmt) {
	var index int
	var lhs ast.Expr
	var ident *ast.Ident
	var ok bool
	for index, lhs = range assign.Lhs {
		ident, ok = lhs.(*ast.Ident)
		if ok && ident.Name == "_" && isError(assignedType(pass, assign, index)) {
			pass.Reportf(ident.Pos(), "error assigned to _; handle it or return it")
		}
	}
}

func assignedType(pass *analysis.Pass, assign *ast.AssignStmt, index int) types.Type {
	var tuple *types.Tuple
	var ok bool
	if len(assign.Lhs) == len(assign.Rhs) {
		return pass.TypesInfo.TypeOf(assign.Rhs[index])
	}
	_, ok = assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return nil
	}
	tuple, ok = pass.TypesInfo.TypeOf(assign.Rhs[0]).(*types.Tuple)
	if !ok || index >= tuple.Len() {
		return nil
	}
	return tuple.At(index).Type()
}

func checkSwallowed(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		var isLit bool
		var stmt *ast.IfStmt
		var ok bool
		_, isLit = node.(*ast.FuncLit)
		if isLit {
			return false
		}
		stmt, ok = node.(*ast.IfStmt)
		if ok {
			reportSwallowed(pass, stmt)
		}
		return true
	})
}

func reportSwallowed(pass *analysis.Pass, stmt *ast.IfStmt) {
	var obj types.Object
	obj = checkedError(pass, stmt.Cond)
	if obj == nil || exits(stmt.Body) || reassigns(pass, stmt.Body, obj) {
		return
	}
	pass.Reportf(stmt.Pos(), "error %q is checked but neither returned nor propagated", obj.Name())
}

func checkedError(pass *analysis.Pass, cond ast.Expr) types.Object {
	var binary *ast.BinaryExpr
	var ok bool
	binary, ok = cond.(*ast.BinaryExpr)
	if !ok || binary.Op != token.NEQ {
		return nil
	}
	if isNil(pass, binary.Y) {
		return errorVar(pass, binary.X)
	}
	if isNil(pass, binary.X) {
		return errorVar(pass, binary.Y)
	}
	return nil
}

func isNil(pass *analysis.Pass, expr ast.Expr) bool {
	var ident *ast.Ident
	var ok bool
	ident, ok = expr.(*ast.Ident)
	if !ok {
		return false
	}
	_, ok = pass.TypesInfo.Uses[ident].(*types.Nil)
	return ok
}

func errorVar(pass *analysis.Pass, expr ast.Expr) types.Object {
	var ident *ast.Ident
	var ok bool
	ident, ok = expr.(*ast.Ident)
	if !ok || !isError(pass.TypesInfo.TypeOf(ident)) {
		return nil
	}
	return pass.TypesInfo.Uses[ident]
}

func exits(body *ast.BlockStmt) bool {
	var found bool
	ast.Inspect(body, func(node ast.Node) bool {
		var isReturn bool
		var isBranch bool
		var isLit bool
		_, isLit = node.(*ast.FuncLit)
		if isLit {
			return false
		}
		_, isReturn = node.(*ast.ReturnStmt)
		_, isBranch = node.(*ast.BranchStmt)
		if isReturn || isBranch {
			found = true
		}
		return !found
	})
	return found
}

func reassigns(pass *analysis.Pass, body *ast.BlockStmt, obj types.Object) bool {
	var found bool
	ast.Inspect(body, func(node ast.Node) bool {
		var assign *ast.AssignStmt
		var lhs ast.Expr
		var ident *ast.Ident
		var ok bool
		assign, ok = node.(*ast.AssignStmt)
		if !ok {
			return !found
		}
		for _, lhs = range assign.Lhs {
			ident, ok = lhs.(*ast.Ident)
			if ok && pass.TypesInfo.Uses[ident] == obj {
				found = true
			}
		}
		return !found
	})
	return found
}
