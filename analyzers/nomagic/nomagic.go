package nomagic

import (
	"go/ast"
	"go/token"
	"go/types"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nomagic",
	Doc:  "no magic numbers (0 and 1 are fine, except as call arguments or in ==/!=); no magic strings that are compared or repeated; use named constants",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var literals []*ast.BasicLit
	compared := map[*ast.BasicLit]bool{}
	counts := map[string]int{}
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		exempt := exemptLiterals(file)
		zeroOne := zeroOneLiterals(pass, file)
		markCompared(file, compared)
		ast.Inspect(file, func(node ast.Node) bool {
			lit, ok := node.(*ast.BasicLit)
			if ok && !exempt[lit] && (isInteresting(lit) || zeroOne[lit]) {
				literals = append(literals, lit)
				counts[lit.Value]++
			}
			return true
		})
	}
	for _, lit := range literals {
		report(pass, lit, compared[lit], counts[lit.Value])
	}
	return nil, nil
}

func isInteresting(lit *ast.BasicLit) bool {
	switch lit.Kind {
	case token.INT, token.FLOAT, token.IMAG:
		return lit.Value != ZERO_LITERAL && lit.Value != ONE_LITERAL
	case token.STRING:
		return lit.Value != EMPTY_STRING
	}
	return false
}

func report(pass *analysis.Pass, lit *ast.BasicLit, compared bool, repeats int) {
	if lit.Kind != token.STRING {
		kit.Reportf(pass, lit.Pos(), "magic number %s; define a constant", lit.Value)
		return
	}
	if compared {
		kit.Reportf(pass, lit.Pos(), "magic string %s in a comparison; define a constant", lit.Value)
		return
	}
	if repeats >= MIN_REPEATS {
		kit.Reportf(pass, lit.Pos(), "magic string %s repeated %d times; define a constant", lit.Value, repeats)
	}
}

func exemptLiterals(file *ast.File) map[*ast.BasicLit]bool {
	exempt := map[*ast.BasicLit]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.GenDecl:
			if typed.Tok == token.CONST {
				markAll(typed, exempt)
			}
		case *ast.ImportSpec:
			exempt[typed.Path] = true
		case *ast.Field:
			if typed.Tag != nil {
				exempt[typed.Tag] = true
			}
		}
		return true
	})
	return exempt
}

func markAll(root ast.Node, set map[*ast.BasicLit]bool) {
	ast.Inspect(root, func(node ast.Node) bool {
		lit, ok := node.(*ast.BasicLit)
		if ok {
			set[lit] = true
		}
		return true
	})
}

func markCompared(file *ast.File, compared map[*ast.BasicLit]bool) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.BinaryExpr:
			if typed.Op == token.EQL || typed.Op == token.NEQ {
				markLiteral(typed.X, compared)
				markLiteral(typed.Y, compared)
			}
		case *ast.CaseClause:
			for _, expr := range typed.List {
				markLiteral(expr, compared)
			}
		}
		return true
	})
}

func markLiteral(expr ast.Expr, set map[*ast.BasicLit]bool) {
	lit, ok := expr.(*ast.BasicLit)
	if ok {
		set[lit] = true
	}
}

// zeroOneLiterals finds 0 and 1 used as call arguments, in == / != comparisons
// and in case clauses. len/cap comparisons and modulo checks are idiomatic.
func zeroOneLiterals(pass *analysis.Pass, file *ast.File) map[*ast.BasicLit]bool {
	flagged := map[*ast.BasicLit]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.CallExpr:
			markArgs(pass, typed, flagged)
		case *ast.BinaryExpr:
			markComparison(pass, typed, flagged)
		case *ast.CaseClause:
			for _, expr := range typed.List {
				markZeroOne(expr, flagged)
			}
		}
		return true
	})
	return flagged
}

func markArgs(pass *analysis.Pass, call *ast.CallExpr, flagged map[*ast.BasicLit]bool) {
	_, isFunc := kit.Callee(pass, call).(*types.Func)
	if !isFunc {
		return
	}
	for _, arg := range call.Args {
		markZeroOne(arg, flagged)
	}
}

func markComparison(pass *analysis.Pass, bin *ast.BinaryExpr, flagged map[*ast.BasicLit]bool) {
	if bin.Op != token.EQL && bin.Op != token.NEQ {
		return
	}
	if isSizeOrModulo(pass, bin.X) || isSizeOrModulo(pass, bin.Y) {
		return
	}
	markZeroOne(bin.X, flagged)
	markZeroOne(bin.Y, flagged)
}

func isSizeOrModulo(pass *analysis.Pass, expr ast.Expr) bool {
	switch typed := expr.(type) {
	case *ast.BinaryExpr:
		return typed.Op == token.REM || typed.Op == token.AND
	case *ast.CallExpr:
		builtin, ok := kit.Callee(pass, typed).(*types.Builtin)
		return ok && (builtin.Name() == LEN_NAME || builtin.Name() == CAP_NAME)
	}
	return false
}

func markZeroOne(expr ast.Expr, flagged map[*ast.BasicLit]bool) {
	lit, ok := expr.(*ast.BasicLit)
	if ok && lit.Kind == token.INT && (lit.Value == ZERO_LITERAL || lit.Value == ONE_LITERAL) {
		flagged[lit] = true
	}
}
