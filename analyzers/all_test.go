package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

const (
	TABINDENT_NAME = "tabindent"
	TABINDENT_FILE = "tabindent/a.go"
	// Space-indented on purpose: on disk gofmt would rewrite it to tabs.
	TABINDENT_GOLDEN_FILE = "tabindent/a.go.golden"
	TABINDENT_GOLDEN      = "package tabindent\n\nfunc good() {\n\t_ = \x60\n    raw string content is verbatim\n\x60\n}\n\nfunc bad() {\n\t_ = 1 // want \"indent with tabs\"\n}\n"
	TABINDENT_SOURCE      = "package tabindent\n\nfunc good() {\n\t_ = \x60\n    raw string content is verbatim\n\x60\n}\n\nfunc bad() {\n    _ = 1 // want \"indent with tabs\"\n}\n"
)

// Every Edil analyzer needs a fixture at testdata/src/<name>. Lines that must
// be reported carry a "// want" comment; every other line must stay silent.
// A.go.golden next to a fixture is the expected result of applying -fix.
func TestAnalyzers(t *testing.T) {
	for _, analyzer := range ANALYZERS {
		t.Run(analyzer.Name, func(t *testing.T) {
			analysistest.RunWithSuggestedFixes(t, fixtureDir(t, analyzer.Name), analyzer, analyzer.Name)
		})
	}
}

func fixtureDir(t *testing.T, name string) string {
	if name != TABINDENT_NAME {
		return analysistest.TestData()
	}
	dir, cleanup, err := analysistest.WriteFiles(map[string]string{TABINDENT_FILE: TABINDENT_SOURCE, TABINDENT_GOLDEN_FILE: TABINDENT_GOLDEN})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)
	return dir
}
