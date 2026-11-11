package analyzers

import (
	"embed"
	"strings"
)

// EXAMPLES holds one short bad/good example per rule: examples/<rule>.md.
//
//go:embed examples
var EXAMPLES embed.FS

// Example returns the bad/good example of a rule, and false when it has none.
func Example(name string) (string, bool) {
	text, err := EXAMPLES.ReadFile(EXAMPLES_DIR + name + EXAMPLE_EXT)
	return strings.TrimSpace(string(text)), err == nil
}
