package errreturn

import (
	"go/ast"
	"go/token"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "errreturn",
	Doc:  "errors must be returned, not assigned to _ or checked and swallowed",
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
			checkBlank(pass, fn.Body)
			if returnsError(pass, fn) {
				checkSwallowed(pass, fn.Body)
			}
		}
	}
	return nil, nil
}

func isError(typ types.Type) bool {
	return typ != nil && types.Identical(typ, types.Universe.Lookup(ERROR_NAME).Type())
}

func returnsError(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil {
		return false
	}
	for _, field := range fn.Type.Results.List {
		if isError(pass.TypesInfo.TypeOf(field.Type)) {
			return true
		}
	}
	return false
}

func checkBlank(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.DeferStmt:
			return false
		case *ast.AssignStmt:
			reportBlank(pass, typed)
		}
		return true
	})
}

func reportBlank(pass *analysis.Pass, assign *ast.AssignStmt) {
	for index, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if ok && ident.Name == BLANK_NAME && isError(assignedType(pass, assign, index)) {
			kit.Reportf(pass, ident.Pos(), "error assigned to _; handle it or return it")
		}
	}
}

func assignedType(pass *analysis.Pass, assign *ast.AssignStmt, index int) types.Type {
	if len(assign.Lhs) == len(assign.Rhs) {
		return pass.TypesInfo.TypeOf(assign.Rhs[index])
	}
	tuple, ok := pass.TypesInfo.TypeOf(assign.Rhs[0]).(*types.Tuple)
	if !ok || index >= tuple.Len() {
		return nil
	}
	return tuple.At(index).Type()
}

func checkSwallowed(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.FuncLit:
			return false
		case *ast.IfStmt:
			reportSwallowed(pass, typed)
		}
		return true
	})
}

func reportSwallowed(pass *analysis.Pass, stmt *ast.IfStmt) {
	obj := checkedError(pass, stmt.Cond)
	if obj == nil || exits(stmt.Body) || reassigns(pass, stmt.Body, obj) {
		return
	}
	kit.Reportf(pass, stmt.Pos(), "error %q is checked but neither returned nor propagated", obj.Name())
}

func checkedError(pass *analysis.Pass, cond ast.Expr) types.Object {
	binary, ok := cond.(*ast.BinaryExpr)
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
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	_, isNilObj := pass.TypesInfo.Uses[ident].(*types.Nil)
	return isNilObj
}

func errorVar(pass *analysis.Pass, expr ast.Expr) types.Object {
	ident, ok := expr.(*ast.Ident)
	if !ok || !isError(pass.TypesInfo.TypeOf(ident)) {
		return nil
	}
	return pass.TypesInfo.Uses[ident]
}

func exits(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt, *ast.BranchStmt:
			found = true
		}
		return !found
	})
	return found
}

func reassigns(pass *analysis.Pass, body *ast.BlockStmt, obj types.Object) bool {
	found := false
	ast.Inspect(body, func(node ast.Node) bool {
		assign, ok := node.(*ast.AssignStmt)
		if ok {
			found = found || targets(pass, assign, obj)
		}
		return !found
	})
	return found
}

func targets(pass *analysis.Pass, assign *ast.AssignStmt, obj types.Object) bool {
	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if ok && pass.TypesInfo.Uses[ident] == obj {
			return true
		}
	}
	return false
}
