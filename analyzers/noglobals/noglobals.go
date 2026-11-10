package noglobals

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noglobals",
	Doc:  "no package-level var, except //go:embed, _ assertions and analyzers; use const or pass values in",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}
			checkDecl(pass, gen)
		}
	}
	return nil, nil
}

func checkDecl(pass *analysis.Pass, gen *ast.GenDecl) {
	for _, spec := range gen.Specs {
		value := spec.(*ast.ValueSpec)
		if hasEmbed(gen.Doc) || hasEmbed(value.Doc) {
			continue
		}
		for _, name := range value.Names {
			if !allowed(pass, name) {
				kit.Reportf(pass, name.Pos(), "package-level var %q; use a const, or pass the value in", name.Name)
			}
		}
	}
}

func hasEmbed(group *ast.CommentGroup) bool {
	if group == nil {
		return false
	}
	for _, comment := range group.List {
		if strings.HasPrefix(comment.Text, EMBED_PREFIX) {
			return true
		}
	}
	return false
}

func allowed(pass *analysis.Pass, name *ast.Ident) bool {
	if name.Name == BLANK_NAME {
		return true
	}
	obj := pass.TypesInfo.Defs[name]
	if obj == nil {
		return false
	}
	switch types.TypeString(obj.Type(), nil) {
	case ANALYZER_TYPE, ANALYZER_LIST_TYPE:
		return true
	}
	return false
}
