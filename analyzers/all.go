package analyzers

import (
	"golang.org/x/tools/go/analysis"

	"edil/analyzers/constcase"
	"edil/analyzers/errreturn"
	"edil/analyzers/layout"
	"edil/analyzers/nopanic"
	"edil/analyzers/predeclare"
	"edil/analyzers/tabindent"
)

var ANALYZERS = []*analysis.Analyzer{
	constcase.Analyzer,
	predeclare.Analyzer,
	nopanic.Analyzer,
	layout.Analyzer,
	errreturn.Analyzer,
	tabindent.Analyzer,
}
