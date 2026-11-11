package kit

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Reportf reports a diagnostic unless an //edil:ignore <rule> <reason>
// comment sits on the same line or the line above.
func Reportf(pass *analysis.Pass, pos token.Pos, format string, args ...any) {
	if ignored(pass, pos) {
		return
	}
	pass.Reportf(pos, format, args...)
}

// Skip reports whether a file is generated or a test file.
func Skip(pass *analysis.Pass, file *ast.File) bool {
	name := pass.Fset.File(file.Pos()).Name()
	return ast.IsGenerated(file) || strings.HasSuffix(name, TEST_SUFFIX)
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

func ignored(pass *analysis.Pass, pos token.Pos) bool {
	file := fileOf(pass, pos)
	if file == nil {
		return false
	}
	line := pass.Fset.Position(pos).Line
	for _, group := range file.Comments {
		for _, comment := range group.List {
			commentLine := pass.Fset.Position(comment.Pos()).Line
			near := commentLine == line || commentLine == line-1
			if near && matches(pass.Analyzer.Name, comment.Text) {
				return true
			}
		}
	}
	return false
}

func matches(rule string, text string) bool {
	parts := strings.Fields(text)
	return len(parts) > IGNORE_MIN_PARTS && parts[0] == IGNORE_PREFIX && parts[1] == rule
}

func fileOf(pass *analysis.Pass, pos token.Pos) *ast.File {
	for _, file := range pass.Files {
		if file.FileStart <= pos && pos <= file.FileEnd {
			return file
		}
	}
	return nil
}
