package nomagic

import "strings"

const LIMIT = 5

func numbers(n int) int {
	if n > 100 { // want `magic number 100`
		return n * 2 // want `magic number 2`
	}
	return n + 1 - 0
}

func compared(s string) bool {
	return s == "admin" // want `magic string "admin" in a comparison`
}

func repeated() string {
	_ = strings.ToUpper("dup")    // want `repeated 2 times`
	return strings.ToLower("dup") // want `repeated 2 times`
}

func fine(s string) string {
	return strings.TrimSpace(s) + ""
}
