package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"edil/analyzers"
)

func main() {
	multichecker.Main(analyzers.ANALYZERS...)
}
