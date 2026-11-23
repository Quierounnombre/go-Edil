package main

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"

	"edil/analyzers"
)

// printExplain lists every rule, or explains each rule named in args.
func printExplain(args []string) error {
	rules := slices.Concat(analyzers.ANALYZERS, analyzers.UPSTREAM)
	if len(args) == NO_RULES {
		for _, rule := range rules {
			fmt.Printf(LIST_FORMAT, rule.Name, summary(rule))
		}
		return nil
	}
	for _, name := range args {
		err := explain(rules, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// explain prints a rule's full description and, for Edil rules, its example.
func explain(rules []*analysis.Analyzer, name string) error {
	index := slices.IndexFunc(rules, func(candidate *analysis.Analyzer) bool {
		return candidate.Name == name
	})
	if index == NOT_FOUND {
		return fmt.Errorf("unknown rule %q; run edil -explain to list them", name)
	}
	rule := rules[index]
	fmt.Printf(TITLE_FORMAT, rule.Name, rule.Doc)
	example, found := analyzers.Example(name)
	if found {
		fmt.Println(example)
	}
	return nil
}

func summary(rule *analysis.Analyzer) string {
	first, _, _ := strings.Cut(rule.Doc, NEWLINE)
	return first
}
