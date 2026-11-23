package uniquename

import "strings"

const LIMIT = 3

func load() error { return nil }

func nested() error {
	err := load()
	for range LIMIT {
		err := load() // want `"err" is already declared outside this block \(line 10\)`
		_ = err
	}
	return err
}

func closure() {
	err := load()
	_ = err
	func() {
		err := load() // want `"err" is already declared outside this block`
		_ = err
	}()
	_ = func(err error) error { // want `"err" is already declared outside this block`
		return err
	}
}

func param(name string) {
	{
		name := "x" // want `"name" is already declared outside this block`
		_ = name
	}
}

func later() {
	{
		count := 1 // want `"count" is already declared outside this block`
		_ = count
	}
	count := 2
	_ = count
}

func outside() {
	strings := 1 // want `"strings" is already declared outside this block \(package level\)`
	len := 2     // want `"len" is already declared outside this block \(builtin\)`
	load := 3    // want `"load" is already declared outside this block \(package level\)`
	_, _, _ = strings, len, load
}

func typeswitch(value any) {
	switch value := value.(type) { // want `"value" is already declared outside this block`
	case int:
		_ = value
	}
}

func siblings(flag bool) {
	if flag {
		err := load()
		_ = err
	}
	if !flag {
		err := load()
		_ = err
	}
	for i := range LIMIT {
		_ = i
	}
	for i := range LIMIT {
		_ = i
	}
}

func reuse() (string, error) {
	first, err := "a", load()
	_ = first
	second, err := "b", load()
	_ = second
	return "", err
}

func good() {
	_ = strings.ToUpper("x")
	_ = func(input string) string {
		return input
	}
}
