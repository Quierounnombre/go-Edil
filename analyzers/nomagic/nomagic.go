package nomagic

import (
	"go/ast"
	"go/token"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nomagic",
	Doc:  "no magic numbers (except 0 and 1); no magic strings that are compared or repeated; use named constants",
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
		markCompared(file, compared)
		ast.Inspect(file, func(node ast.Node) bool {
			lit, ok := node.(*ast.BasicLit)
			if ok && !exempt[lit] && isInteresting(lit) {
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
