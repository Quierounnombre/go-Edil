package errwrap

import (
	"fmt"
	"go/ast"
	"go/token"
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
	for index, found := range verbs(format) {
		if index >= len(args) {
			return
		}
		if (found.verb == VERB_V || found.verb == VERB_S) && isErrorValue(pass, args[index]) {
			report(pass, call, found, args[index])
		}
	}
}

func report(pass *analysis.Pass, call *ast.CallExpr, found verbAt, arg ast.Expr) {
	message := fmt.Sprintf("error formatted with %%%c; use %%w to keep the error chain", found.verb)
	edits := fixEdits(pass, call, found, arg)
	if edits == nil {
		kit.Report(pass, arg.Pos(), message)
		return
	}
	kit.ReportFix(pass, arg.Pos(), message, "use %w", edits)
}

// fixEdits is nil unless the format is a plain interpreted literal, where
// source offsets equal string offsets.
func fixEdits(pass *analysis.Pass, call *ast.CallExpr, found verbAt, arg ast.Expr) []analysis.TextEdit {
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING || strings.ContainsRune(lit.Value, ESCAPE) || lit.Value[0] != QUOTE {
		return nil
	}
	at := lit.Pos() + token.Pos(len(QUOTE_STRING)+found.offset)
	edits := []analysis.TextEdit{kit.Edit(at, at+token.Pos(len(string(found.verb))), string(VERB_W))}
	method, isCall := arg.(*ast.CallExpr)
	if !isCall || !isErrorMethod(pass, method, types.Universe.Lookup(ERROR_NAME).Type()) {
		return edits
	}
	sel := method.Fun.(*ast.SelectorExpr)
	return append(edits, kit.Edit(sel.X.End(), method.End(), ""))
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

func verbs(format string) []verbAt {
	var found []verbAt
	for index := 0; index < len(format); index++ {
		if format[index] != PERCENT {
			continue
		}
		verb, at := readVerb(format, index+1)
		index = at
		if verb != PERCENT && verb != NO_VERB {
			found = append(found, verbAt{verb: verb, offset: at})
		}
	}
	return found
}

func readVerb(format string, start int) (rune, int) {
	index := start
	for index < len(format) && strings.ContainsRune(FLAG_CHARS, rune(format[index])) {
		index++
	}
	if index >= len(format) {
		return NO_VERB, index
	}
	return rune(format[index]), index
}
