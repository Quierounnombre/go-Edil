package errwrap

import (
	"go/ast"
	"go/types"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "errwrap",
	Doc:  "fmt.Errorf must wrap errors with %w, not format them with %v or %s",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if ok && kit.Qualified(kit.Callee(pass, call)) == FMT_ERRORF {
				check(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

func check(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) < MIN_ARGS {
		return
	}
	format, ok := kit.ConstString(pass, call.Args[0])
	if !ok {
		return
	}
	args := call.Args[1:]
	for index, verb := range verbs(format) {
		if index >= len(args) {
			return
		}
		if (verb == VERB_V || verb == VERB_S) && isErrorValue(pass, args[index]) {
			kit.Reportf(pass, args[index].Pos(), "error formatted with %%%c; use %%w to keep the error chain", verb)
		}
	}
}

func isErrorValue(pass *analysis.Pass, expr ast.Expr) bool {
	errType := types.Universe.Lookup(ERROR_NAME).Type()
	typ := pass.TypesInfo.TypeOf(expr)
	if typ != nil && typ != types.Typ[types.UntypedNil] && types.AssignableTo(typ, errType) {
		return true
	}
	call, isCall := expr.(*ast.CallExpr)
	return isCall && isErrorMethod(pass, call, errType)
}

func isErrorMethod(pass *analysis.Pass, call *ast.CallExpr, errType types.Type) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != ERROR_METHOD || len(call.Args) != 0 {
		return false
	}
	recv := pass.TypesInfo.TypeOf(sel.X)
	return recv != nil && types.AssignableTo(recv, errType)
}

func verbs(format string) []rune {
	var found []rune
	runes := []rune(format)
	for index := 0; index < len(runes); index++ {
		if runes[index] != PERCENT {
			continue
		}
		verb, next := readVerb(runes, index+1)
		index = next
		if verb != PERCENT && verb != 0 {
			found = append(found, verb)
		}
	}
	return found
}

func readVerb(runes []rune, start int) (rune, int) {
	index := start
	for index < len(runes) && strings.ContainsRune(FLAG_CHARS, runes[index]) {
		index++
	}
	if index >= len(runes) {
		return 0, index
	}
	return runes[index], index
}
