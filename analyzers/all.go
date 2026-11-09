package analyzers

import (
	"golang.org/x/tools/go/analysis"

	"edil/analyzers/constcase"
	"edil/analyzers/errreturn"
	"edil/analyzers/layout"
	"edil/analyzers/nopanic"
	"edil/analyzers/tabindent"
)

var ANALYZERS = []*analysis.Analyzer{
	constcase.Analyzer,
	nopanic.Analyzer,
	layout.Analyzer,
	errreturn.Analyzer,
	tabindent.Analyzer,
}
