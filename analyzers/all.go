package analyzers

import (
	"edil/analyzers/boolreturn"
	"edil/analyzers/complexity"
	"edil/analyzers/constcase"
	"edil/analyzers/deferinloop"
	"edil/analyzers/errreturn"
	"edil/analyzers/errstring"
	"edil/analyzers/errwrap"
	"edil/analyzers/execvar"
	"edil/analyzers/funclen"
	"edil/analyzers/hardcodedsecret"
	"edil/analyzers/httpdefault"
	"edil/analyzers/layout"
	"edil/analyzers/mathrand"
	"edil/analyzers/maxdepth"
	"edil/analyzers/maxparams"
	"edil/analyzers/nakereturn"
	"edil/analyzers/nobackgroundctx"
	"edil/analyzers/noelse"
	"edil/analyzers/noglobals"
	"edil/analyzers/noifinit"
	"edil/analyzers/nomagic"
	"edil/analyzers/nopanic"
	"edil/analyzers/pkgsize"
	"edil/analyzers/sqlconcat"
	"edil/analyzers/tabindent"
	"edil/analyzers/tlsinsecure"
	"edil/analyzers/varinit"
	"edil/analyzers/weakcrypto"

	"golang.org/x/tools/go/analysis"
)

var ANALYZERS = []*analysis.Analyzer{
	boolreturn.Analyzer,
	complexity.Analyzer,
	constcase.Analyzer,
	deferinloop.Analyzer,
	errreturn.Analyzer,
	errstring.Analyzer,
	errwrap.Analyzer,
	execvar.Analyzer,
	funclen.Analyzer,
	hardcodedsecret.Analyzer,
	httpdefault.Analyzer,
	layout.Analyzer,
	mathrand.Analyzer,
	maxdepth.Analyzer,
	maxparams.Analyzer,
	nakedreturn.Analyzer,
	nobackgroundctx.Analyzer,
	noelse.Analyzer,
	noglobals.Analyzer,
	noifinit.Analyzer,
	nomagic.Analyzer,
	nopanic.Analyzer,
	pkgsize.Analyzer,
	sqlconcat.Analyzer,
	tabindent.Analyzer,
	tlsinsecure.Analyzer,
	varinit.Analyzer,
	weakcrypto.Analyzer,
}
