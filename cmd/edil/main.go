package main

import (
	"slices"

	"golang.org/x/tools/go/analysis/multichecker"

	"edil/analyzers"
)

func main() {
	multichecker.Main(slices.Concat(analyzers.ANALYZERS, analyzers.UPSTREAM)...)
}
