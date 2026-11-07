package constcase

import (
	"go/ast"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"
)

var SCREAMING_SNAKE = regexp.MustCompile(`^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$`)

var Analyzer = &analysis.Analyzer{
	Name: "constcase",
	Doc:  "constants must be SCREAMING_SNAKE_CASE",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				return true
			}
			for _, spec := range decl.Specs {
				for _, name := range spec.(*ast.ValueSpec).Names {
					if name.Name != "_" && !SCREAMING_SNAKE.MatchString(name.Name) {
						pass.Reportf(name.Pos(), "constant %q must be SCREAMING_SNAKE_CASE", name.Name)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

