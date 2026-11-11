package mathrand

import (
	"strconv"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mathrand",
	Doc:  "math/rand is predictable; use crypto/rand",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, spec := range file.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return nil, err
			}
			if path == MATH_RAND || path == MATH_RAND_V2 {
				kit.Reportf(pass, spec.Pos(), "%s is predictable; use crypto/rand for tokens, keys and ids", path)
			}
		}
	}
	return nil, nil
}
