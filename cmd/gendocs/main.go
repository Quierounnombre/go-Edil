package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"edil/analyzers"
)

//go:embed readme.template.md conventions.template.md
var FILES embed.FS

func buildPage() (Page, error) {
	rules := make([]Rule, 0, len(analyzers.ANALYZERS))
	for _, analyzer := range analyzers.ANALYZERS {
		summary, _, _ := strings.Cut(analyzer.Doc, NEWLINE)
		doc := strings.ReplaceAll(summary, PIPE, PIPE_ESCAPED)
		example, found := analyzers.Example(analyzer.Name)
		if !found {
			return Page{}, fmt.Errorf("rule %s has no example in analyzers/examples", analyzer.Name)
		}
		rules = append(rules, Rule{Name: analyzer.Name, Doc: doc, Example: example})
	}
	return Page{Repo: REPO, Rules: rules}, nil
}

func render(page Page, source string) (string, error) {
	text, err := FILES.ReadFile(source)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(source).Parse(string(text))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = tmpl.Execute(&out, page)
	return out.String(), err
}

func outputs() []output {
	return []output{
		{path: README_PATH, template: README_TEMPLATE},
		{path: CONVENTIONS_PATH, template: CONVENTIONS_TMPL},
	}
}

func apply(page Page, target output, check bool) error {
	updated, err := render(page, target.template)
	if err != nil {
		return err
	}
	if !check {
		return os.WriteFile(target.path, []byte(updated), FILE_MODE)
	}
	current, err := os.ReadFile(target.path)
	if err != nil {
		return err
	}
	if string(current) != updated {
		return errors.New(target.path + " is out of date, run: go run ./cmd/gendocs")
	}
	return nil
}

func run() error {
	check := flag.Bool("check", false, "fail if generated docs are out of date")
	flag.Parse()
	page, err := buildPage()
	if err != nil {
		return err
	}
	for _, target := range outputs() {
		err = apply(page, target, *check)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EXIT_FAILURE)
	}
}
