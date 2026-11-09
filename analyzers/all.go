package analyzers

import (
	"golang.org/x/tools/go/analysis"

	"edil/analyzers/constcase"
	"edil/analyzers/errreturn"
	"edil/analyzers/layout"
	"edil/analyzers/nopanic"
	"edil/analyzers/tabindent"
	"edil/analyzers/maxdepth"
	"edil/analyzers/nakereturn"
	"edil/analyzers/noelse"
	"edil/analyzers/noifinit"
	"edil/analyzers/pkgsize"
)

var ANALYZERS = []*analysis.Analyzer{
	constcase.Analyzer,
	nopanic.Analyzer,
	layout.Analyzer,
	errreturn.Analyzer,
	tabindent.Analyzer,
	maxdepth.Analyzer,
	nakedreturn.Analyzer,
	noelse.Analyzer,
	noifinit.Analyzer,
	pkgsize.Analyzer,
}
