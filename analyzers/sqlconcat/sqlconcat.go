package sqlconcat

import (
	"go/ast"
	"go/types"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "sqlconcat",
	Doc:  "database/sql queries must be constant strings; pass values as placeholder arguments",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if ok {
				check(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

func check(pass *analysis.Pass, call *ast.CallExpr) {
	fn, ok := kit.Callee(pass, call).(*types.Func)
	if !ok {
		return
	}
	index, isQuery := queryIndex(fn)
	if !isQuery || index >= len(call.Args) {
		return
	}
	query := call.Args[index]
	if pass.TypesInfo.Types[query].Value == nil {
		kit.Reportf(pass, query.Pos(), "SQL query must be a constant string; pass values as placeholder arguments")
	}
}

func queryIndex(fn *types.Func) (int, bool) {
	recv := fn.Type().(*types.Signature).Recv()
	if recv == nil || !isSQLReceiver(recv.Type()) {
		return 0, false
	}
	base := strings.TrimSuffix(fn.Name(), SUFFIX_CONTEXT)
	switch base {
	case QUERY, QUERY_ROW, EXEC, PREPARE:
		if base != fn.Name() {
			return 1, true
		}
		return 0, true
	}
	return 0, false
}

func isSQLReceiver(typ types.Type) bool {
	pointer, isPointer := typ.(*types.Pointer)
	if isPointer {
		typ = pointer.Elem()
	}
	named, ok := typ.(*types.Named)
	if !ok || named.Obj().Pkg() == nil || named.Obj().Pkg().Path() != SQL_PKG {
		return false
	}
	switch named.Obj().Name() {
	case DB_NAME, TX_NAME, CONN_NAME:
		return true
	}
	return false
}
