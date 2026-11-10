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
	"edil/analyzers/boolreturn"
	"edil/analyzers/complexity"
	"edil/analyzers/funclen"
	"edil/analyzers/httpdefault"
	"edil/analyzers/mathrand"
	"edil/analyzers/maxparams"
	"edil/analyzers/nobackgroundctx"
	"edil/analyzers/noglobals"
	"edil/analyzers/nomagic"
	"edil/analyzers/tlsinsecure"
	"edil/analyzers/varinit"
	"edil/analyzers/weakcrypto"
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
	boolreturn.Analyzer,
	complexity.Analyzer,
	funclen.Analyzer,
	httpdefault.Analyzer,
	mathrand.Analyzer,
	maxparams.Analyzer,
	nobackgroundctx.Analyzer,
	noglobals.Analyzer,
	nomagic.Analyzer,
	tlsinsecure.Analyzer,
	varinit.Analyzer,
	weakcrypto.Analyzer,
}
