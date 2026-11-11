package errstring

import (
	"go/ast"
	"strings"
	"unicode"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "errstring",
	Doc:  "error strings are lowercase and do not end with punctuation or a newline",
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
	switch kit.Qualified(kit.Callee(pass, call)) {
	case ERRORS_NEW, FMT_ERRORF:
		if len(call.Args) == 0 {
			return
		}
		text, ok := kit.ConstString(pass, call.Args[0])
		problem := problemWith(text)
		if ok && problem != "" {
			kit.Reportf(pass, call.Args[0].Pos(), "error string %q %s", text, problem)
		}
	}
}

func problemWith(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}
	if len(runes) > 1 && unicode.IsUpper(runes[0]) && unicode.IsLower(runes[1]) {
		return "must not be capitalized"
	}
	if strings.ContainsRune(END_CHARS, runes[len(runes)-1]) {
		return "must not end with punctuation or a newline"
	}
	return ""
}
