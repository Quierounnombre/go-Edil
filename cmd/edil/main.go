package main

import (
	"fmt"
	"os"
	"slices"

	"golang.org/x/tools/go/analysis/multichecker"

	"edil/analyzers"
)

func main() {
	if len(os.Args) >= MIN_ARGS && os.Args[1] == EXPLAIN_FLAG {
		err := printExplain(os.Args[MIN_ARGS:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(EXIT_FAILURE)
		}
		return
	}
	multichecker.Main(slices.Concat(analyzers.ANALYZERS, analyzers.UPSTREAM)...)
}
