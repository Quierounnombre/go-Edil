package tabindent

import (
	"go/ast"
	"go/token"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "tabindent",
	Doc:  "indentation must use tabs, not spaces",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		err := checkFile(pass, file)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File) error {
	tokFile := pass.Fset.File(file.Pos())
	src, err := pass.ReadFile(tokFile.Name())
	if err != nil {
		return err
	}
	skip := verbatimLines(pass, file)
	count, first := 0, 0
	for index, line := range strings.Split(string(src), NEWLINE) {
		if skip[index+1] || !hasSpaceIndent(line) {
			continue
		}
		if count == 0 {
			first = index + 1
		}
		count++
	}
	if count > 0 {
		kit.Reportf(pass, tokFile.LineStart(first), "indent with tabs, not spaces (%d lines, first here)", count)
	}
	return nil
}

func hasSpaceIndent(line string) bool {
	trimmed := strings.TrimLeft(line, INDENT_CHARS)
	if trimmed == "" {
		return false
	}
	indent := line[:len(line)-len(trimmed)]
	return strings.Contains(indent, SPACE)
}

func verbatimLines(pass *analysis.Pass, file *ast.File) map[int]bool {
	skip := map[int]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		lit, ok := node.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING && lit.Value[0] == '`' {
			markLines(pass, skip, lit.Pos(), lit.End())
		}
		return true
	})
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if strings.HasPrefix(comment.Text, BLOCK_OPEN) {
				markLines(pass, skip, comment.Pos(), comment.End())
			}
		}
	}
	return skip
}

func markLines(pass *analysis.Pass, skip map[int]bool, start token.Pos, end token.Pos) {
	last := pass.Fset.Position(end).Line
	for line := pass.Fset.Position(start).Line + 1; line <= last; line++ {
		skip[line] = true
	}
}
