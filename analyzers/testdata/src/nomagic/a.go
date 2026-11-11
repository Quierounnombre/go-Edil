package nomagic

import (
	"os"
	"strings"
)

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

func zeroOne(n int, items []int) bool {
	os.Exit(1) // want `magic number 1`
	_ = make([]int, 0)
	total := 0
	total++
	_ = len(items) == 0
	_ = n%LIMIT == 0
	_ = total
	return n == 0 // want `magic number 0`
}

func switched(n int) {
	switch n {
	case 1: // want `magic number 1`
	case LIMIT:
	}
}
