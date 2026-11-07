package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"edil/analyzers"
)

const (
	README_PATH = "README.md"
	PLACEHOLDER = "<!-- RULES -->\n"
	REPO        = "github:quierounnombre/edil"
)

//go:embed readme.template.md
var TEMPLATE string

func renderRows() string {
	var builder strings.Builder

	for _, analyzer := range analyzers.ANALYZERS {
		fmt.Fprintf(&builder, "| `%s` | %s |\n",
			analyzer.Name, strings.ReplaceAll(analyzer.Doc, "|", "\\|"))
	}
	return builder.String()
}

func render() (string, error) {
	var rendered string

	if !strings.Contains(TEMPLATE, PLACEHOLDER) {
		return "", errors.New("placeholder not found in readme.template.md")
	}
	rendered = strings.Replace(TEMPLATE, PLACEHOLDER, renderRows(), 1)
	return strings.ReplaceAll(rendered, "@REPO@", REPO), nil
}

func run() error {
	var check = flag.Bool("check", false, "fail if README.md is out of date")
	var current []byte
	var updated string
	var err error

	flag.Parse()
	updated, err = render()
	if err != nil {
		return err
	}
	if !*check {
		return os.WriteFile(README_PATH, []byte(updated), 0o644)
	}
	current, err = os.ReadFile(README_PATH)
	if err != nil {
		return err
	}
	if string(current) != updated {
		return errors.New(README_PATH + " is out of date, run: go run ./cmd/gendocs")
	}
	return nil
}

func main() {
	var err = run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
