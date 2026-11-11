package constcase

import (
	"fmt"
	"go/ast"
	"go/types"
	"unicode"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "constcase",
	Doc:  "constants must be SCREAMING_SNAKE_CASE",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			spec, ok := node.(*ast.ValueSpec)
			if ok {
				checkSpec(pass, spec)
			}
			return true
		})
	}
	return nil, nil
}

func checkSpec(pass *analysis.Pass, spec *ast.ValueSpec) {
	for _, name := range spec.Names {
		_, isConst := pass.TypesInfo.Defs[name].(*types.Const)
		if isConst && name.Name != BLANK_NAME && !isScreamingSnake(name.Name) {
			report(pass, name)
		}
	}
}

func isScreamingSnake(name string) bool {
	runes := []rune(name)
	if len(runes) == 0 || !isUpper(runes[0]) {
		return false
	}
	for index, char := range runes {
		if !isAllowed(char) || (char == UNDERSCORE && isEdgeOrDouble(runes, index)) {
			return false
		}
	}
	return true
}

func isEdgeOrDouble(runes []rune, index int) bool {
	return index == len(runes)-1 || runes[index+1] == UNDERSCORE
}

func isUpper(char rune) bool {
	return char >= 'A' && char <= 'Z'
}

func isAllowed(char rune) bool {
	return isUpper(char) || (char >= '0' && char <= '9') || char == UNDERSCORE
}

func report(pass *analysis.Pass, name *ast.Ident) {
	message := fmt.Sprintf("constant %q must be SCREAMING_SNAKE_CASE", name.Name)
	renamed := screaming(name.Name)
	obj := pass.TypesInfo.Defs[name]
	if !renameable(obj, renamed) {
		kit.Report(pass, name.Pos(), message)
		return
	}
	var edits []analysis.TextEdit
	for _, file := range pass.Files {
		edits = append(edits, renameEdits(pass, file, obj, renamed)...)
	}
	kit.ReportFix(pass, name.Pos(), message, "rename to "+renamed, edits)
}

// renameable allows only unexported names: exported ones are used by packages
// this analysis cannot see.
func renameable(obj types.Object, renamed string) bool {
	if obj == nil || obj.Exported() || !isScreamingSnake(renamed) {
		return false
	}
	_, existing := obj.Parent().LookupParent(renamed, obj.Pos())
	return existing == nil
}

func renameEdits(pass *analysis.Pass, file *ast.File, obj types.Object, renamed string) []analysis.TextEdit {
	var edits []analysis.TextEdit
	ast.Inspect(file, func(node ast.Node) bool {
		ident, ok := node.(*ast.Ident)
		if ok && (pass.TypesInfo.Defs[ident] == obj || pass.TypesInfo.Uses[ident] == obj) {
			edits = append(edits, kit.Edit(ident.Pos(), ident.End(), renamed))
		}
		return true
	})
	return edits
}

// screaming turns camelCase into CAMEL_CASE.
func screaming(name string) string {
	var out []rune
	runes := []rune(name)
	for index, char := range runes {
		if index > 0 && unicode.IsUpper(char) && startsWord(runes, index) {
			out = append(out, UNDERSCORE)
		}
		out = append(out, unicode.ToUpper(char))
	}
	return string(out)
}

func startsWord(runes []rune, index int) bool {
	previous := runes[index-1]
	if previous == UNDERSCORE {
		return false
	}
	nextLower := index+1 < len(runes) && unicode.IsLower(runes[index+1])
	return !unicode.IsUpper(previous) || nextLower
}
