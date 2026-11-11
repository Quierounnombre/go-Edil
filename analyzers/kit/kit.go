package kit

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Reportf reports a diagnostic. There is deliberately no way to suppress one:
// a wrong rule gets fixed in Edil, never bypassed in the code under check.
func Reportf(pass *analysis.Pass, pos token.Pos, format string, args ...any) {
	pass.Reportf(pos, format, args...)
}

// Report reports an already formatted message.
func Report(pass *analysis.Pass, pos token.Pos, message string) {
	pass.Report(analysis.Diagnostic{Pos: pos, Message: message})
}

// ReportFix reports a diagnostic with one suggested fix. A fix is only offered
// when it is mechanical and cannot change behavior.
func ReportFix(pass *analysis.Pass, pos token.Pos, message string, fixMessage string, edits []analysis.TextEdit) {
	pass.Report(analysis.Diagnostic{
		Pos:     pos,
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message:   fixMessage,
			TextEdits: edits,
		}},
	})
}

// Edit builds a text edit replacing [start, end) with text.
func Edit(start token.Pos, end token.Pos, text string) analysis.TextEdit {
	return analysis.TextEdit{Pos: start, End: end, NewText: []byte(text)}
}

// Skip reports whether a file is a test file. It is the only exemption.
func Skip(pass *analysis.Pass, file *ast.File) bool {
	return strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), TEST_SUFFIX)
}

// IsMain reports whether fn is func main in package main.
func IsMain(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	return pass.Pkg.Name() == MAIN_NAME && fn.Recv == nil && fn.Name.Name == MAIN_NAME
}

// Callee returns the object a call expression invokes, or nil.
func Callee(pass *analysis.Pass, call *ast.CallExpr) types.Object {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return pass.TypesInfo.Uses[fun]
	case *ast.SelectorExpr:
		return pass.TypesInfo.Uses[fun.Sel]
	}
	return nil
}

// Qualified returns "pkgpath.Name" for package-level objects, "" for
// builtins and methods.
func Qualified(obj types.Object) string {
	if obj == nil || obj.Pkg() == nil {
		return ""
	}
	fn, isFunc := obj.(*types.Func)
	if isFunc && fn.Type().(*types.Signature).Recv() != nil {
		return ""
	}
	return obj.Pkg().Path() + QUALIFY_SEP + obj.Name()
}

// ConstString returns the value of a constant string expression.
func ConstString(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	value := pass.TypesInfo.Types[expr].Value
	if value == nil || value.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(value), true
}
