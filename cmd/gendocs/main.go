package main

import (
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"edil/analyzers"
)

//go:embed readme.template.md
var TEMPLATE string

func buildPage() Page {
	rules := make([]Rule, 0, len(analyzers.ANALYZERS))
	for _, analyzer := range analyzers.ANALYZERS {
		summary, _, _ := strings.Cut(analyzer.Doc, NEWLINE)
		doc := strings.ReplaceAll(summary, PIPE, PIPE_ESCAPED)
		rules = append(rules, Rule{Name: analyzer.Name, Doc: doc})
	}
	return Page{Repo: REPO, Rules: rules}
}

func render() (string, error) {
	tmpl, err := template.New(TEMPLATE_NAME).Parse(TEMPLATE)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = tmpl.Execute(&out, buildPage())
	return out.String(), err
}

func run() error {
	check := flag.Bool("check", false, "fail if README.md is out of date")
	flag.Parse()
	updated, err := render()
	if err != nil {
		return err
	}
	if !*check {
		return os.WriteFile(README_PATH, []byte(updated), FILE_MODE)
	}
	current, err := os.ReadFile(README_PATH)
	if err != nil {
		return err
	}
	if string(current) != updated {
		return errors.New(README_PATH + " is out of date, run: go run ./cmd/gendocs")
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
