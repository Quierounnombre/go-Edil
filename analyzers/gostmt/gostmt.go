package gostmt

import (
	"go/ast"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "gostmt",
	Doc:  "no go statements outside main; use errgroup.Group.Go or sync.WaitGroup.Go so goroutines are waited on and report errors",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.FuncDecl:
				return !kit.IsMain(pass, typed)
			case *ast.GoStmt:
				kit.Reportf(pass, typed.Pos(), "go statement; use errgroup.Group.Go or sync.WaitGroup.Go so the goroutine is waited on and its error returned")
			}
			return true
		})
	}
	return nil, nil
}
