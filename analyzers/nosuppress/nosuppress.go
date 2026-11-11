package nosuppress

import (
	"go/ast"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nosuppress",
	Doc:  "no lint suppression comments (nolint, nosec, lint:ignore, edil:ignore); Edil does not honor them, fix the code or the rule",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		checkFile(pass, file)
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File) {
	for _, group := range file.Comments {
		checkGroup(pass, group)
	}
}

func checkGroup(pass *analysis.Pass, group *ast.CommentGroup) {
	for _, comment := range group.List {
		if isSuppression(comment.Text) {
			kit.Reportf(pass, comment.Pos(), "suppression comments are not honored and not allowed; fix the code")
		}
	}
}

func directives() []string {
	return []string{DIRECTIVE_NOLINT, DIRECTIVE_NOSEC, DIRECTIVE_LINT_IGNORE, DIRECTIVE_LINT_FILE_IGNORE, DIRECTIVE_EDIL_IGNORE}
}

func isSuppression(text string) bool {
	body := strings.TrimSpace(strings.TrimPrefix(text, COMMENT_PREFIX))
	for _, directive := range directives() {
		if strings.HasPrefix(body, directive) {
			return true
		}
	}
	return false
}
