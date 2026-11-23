package uniquename

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "uniquename",
	Doc:  "a name declared in a block cannot be declared again in any block inside it, closures included; it must not reuse a package, import or builtin name either",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.Ident:
				check(pass, pass.TypesInfo.Defs[typed], typed.Pos())
			case *ast.TypeSwitchStmt:
				checkTypeSwitch(pass, typed)
			}
			return true
		})
	}
	return nil, nil
}

// A type switch binding has no Defs entry: its object lives in each clause.
func checkTypeSwitch(pass *analysis.Pass, stmt *ast.TypeSwitchStmt) {
	assign, isAssign := stmt.Assign.(*ast.AssignStmt)
	if !isAssign {
		return
	}
	for _, clause := range stmt.Body.List {
		obj := pass.TypesInfo.Implicits[clause]
		if obj != nil {
			check(pass, obj, assign.Lhs[0].Pos())
			return
		}
	}
}

func check(pass *analysis.Pass, obj types.Object, pos token.Pos) {
	if !isLocal(pass, obj) {
		return
	}
	outer := enclosing(obj)
	if outer == nil {
		return
	}
	kit.Reportf(pass, pos, "%q is already declared outside this block (%s); use a new name, or assign with =", obj.Name(), describe(pass, outer))
}

func isLocal(pass *analysis.Pass, obj types.Object) bool {
	switch obj.(type) {
	case *types.Var, *types.Const, *types.TypeName:
	default:
		return false
	}
	scope := obj.Parent()
	return obj.Name() != BLANK && scope != nil && scope != pass.Pkg.Scope()
}

// enclosing finds the same name in any scope around obj's own, regardless of
// position: a later declaration in an outer block still collides.
func enclosing(obj types.Object) types.Object {
	for scope := obj.Parent().Parent(); scope != nil; scope = scope.Parent() {
		found := scope.Lookup(obj.Name())
		if found != nil {
			return found
		}
	}
	return nil
}

func describe(pass *analysis.Pass, outer types.Object) string {
	if outer.Pkg() == nil {
		return "builtin"
	}
	_, isImport := outer.(*types.PkgName)
	if isImport || outer.Parent() == pass.Pkg.Scope() {
		return "package level"
	}
	return fmt.Sprintf("line %d", pass.Fset.Position(outer.Pos()).Line)
}
