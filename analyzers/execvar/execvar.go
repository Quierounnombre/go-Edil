package execvar

import (
	"go/ast"
	"go/token"
	"path"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "execvar",
	Doc:  "exec.Command runs a constant binary; a shell must never get a non-constant script",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if ok {
				check(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

func commandStart(qualified string) (int, bool) {
	switch qualified {
	case EXEC_COMMAND:
		return 0, true
	case EXEC_COMMAND_CTX:
		return 1, true
	}
	return 0, false
}

func check(pass *analysis.Pass, call *ast.CallExpr) {
	start, isExec := commandStart(kit.Qualified(kit.Callee(pass, call)))
	if !isExec || start >= len(call.Args) {
		return
	}
	name, isConst := kit.ConstString(pass, call.Args[start])
	if !isConst {
		kit.Reportf(pass, call.Args[start].Pos(), "command name must be a constant; run a fixed binary with variable arguments")
		return
	}
	if isShell(name) && hasDynamicArgs(pass, call, start+1) {
		kit.Reportf(pass, call.Pos(), "shell %q runs a non-constant script; never build shell strings", name)
	}
}

func isShell(name string) bool {
	switch strings.ToLower(path.Base(name)) {
	case SH, BASH, ZSH, DASH, CMD, CMD_EXE, POWERSHELL, PWSH:
		return true
	}
	return false
}

func hasDynamicArgs(pass *analysis.Pass, call *ast.CallExpr, from int) bool {
	if call.Ellipsis != token.NoPos {
		return true
	}
	for index := from; index < len(call.Args); index++ {
		_, isConst := kit.ConstString(pass, call.Args[index])
		if !isConst {
			return true
		}
	}
	return false
}
