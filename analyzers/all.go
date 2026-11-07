package analyzers

import (
	"golang.org/x/tools/go/analysis"

	"edil/analyzers/constcase"
	"edil/analyzers/predeclare"
)

var ANALYZERS = []*analysis.Analyzer{
	constcase.Analyzer,
	predeclare.Analyzer,
}
