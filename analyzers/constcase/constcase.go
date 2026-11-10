package constcase

import (
	"go/ast"
	"go/types"

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
			kit.Reportf(pass, name.Pos(), "constant %q must be SCREAMING_SNAKE_CASE", name.Name)
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
