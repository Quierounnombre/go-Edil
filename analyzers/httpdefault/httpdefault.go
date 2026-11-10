package httpdefault

import (
	"go/ast"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "httpdefault",
	Doc:  "no http.Get/Post/Head/PostForm, http.DefaultClient or http.Client{} without Timeout; they never time out",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.SelectorExpr:
				checkSelector(pass, typed)
			case *ast.CompositeLit:
				checkClient(pass, typed)
			}
			return true
		})
	}
	return nil, nil
}

func checkSelector(pass *analysis.Pass, sel *ast.SelectorExpr) {
	switch kit.Qualified(pass.TypesInfo.Uses[sel.Sel]) {
	case HTTP_GET, HTTP_HEAD, HTTP_POST, HTTP_POST_FORM, DEFAULT_CLIENT:
		kit.Reportf(pass, sel.Pos(), "%s has no timeout; use an http.Client with Timeout", sel.Sel.Name)
	}
}

func checkClient(pass *analysis.Pass, lit *ast.CompositeLit) {
	named, ok := pass.TypesInfo.TypeOf(lit).(*types.Named)
	if !ok || named.Obj().Pkg() == nil {
		return
	}
	if named.Obj().Pkg().Path() != HTTP_PKG || named.Obj().Name() != CLIENT_NAME {
		return
	}
	if !hasTimeout(lit) {
		kit.Reportf(pass, lit.Pos(), "http.Client without Timeout never times out")
	}
}

func hasTimeout(lit *ast.CompositeLit) bool {
	for _, element := range lit.Elts {
		pair, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, isIdent := pair.Key.(*ast.Ident)
		if isIdent && key.Name == TIMEOUT_NAME {
			return true
		}
	}
	return false
}
