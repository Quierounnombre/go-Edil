package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"edil/analyzers"
)

const (
	DOCS_PATH   = "docs/index.html"
	PLACEHOLDER = "<!-- RULES -->"
	REPO        = "github:quierounnombre/edil"
)

//go:embed template.html

var TEMPLATE string

func renderRows(indent string) string {
	var builder strings.Builder

	for _, analyzer := range analyzers.ANALYZERS {
		fmt.Fprintf(&builder, "%s<tr><td><code>%s</code></td><td>%s</td></tr>\n",
			indent, analyzer.Name, html.EscapeString(analyzer.Doc))
	}
	return builder.String()
}

func render() (string, error) {
	var position = strings.Index(TEMPLATE, PLACEHOLDER)
	var rendered string
	var lineStart int
	var indent string

	if position < 0 {
		return "", errors.New("placeholder not found in template.html")
	}
	lineStart = strings.LastIndex(TEMPLATE[:position], "\n") + 1
	indent = TEMPLATE[lineStart:position]
	rendered = TEMPLATE[:lineStart] + renderRows(indent) + TEMPLATE[position+len(PLACEHOLDER)+1:]
	return strings.ReplaceAll(rendered, "@REPO@", REPO), nil
}

func run() error {
	var check = flag.Bool("check", false, "fail if docs are out of date")
	var current []byte
	var updated string
	var err error

	flag.Parse()
	updated, err = render()
	if err != nil {
		return err
	}
	if *check {
		current, err = os.ReadFile(DOCS_PATH)
		if err != nil {
			return err
		}
		if string(current) != updated {
			return errors.New(DOCS_PATH + " is out of date, run: go run ./cmd/gendocs")
		}
		return nil
	}
	err = os.MkdirAll(filepath.Dir(DOCS_PATH), 0o755)
	if err != nil {
		return err
	}
	return os.WriteFile(DOCS_PATH, []byte(updated), 0o644)
}

func main() {
	var err = run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
