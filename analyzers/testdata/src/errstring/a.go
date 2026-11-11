package errstring

import (
	"errors"
	"fmt"
)

func bad() {
	_ = errors.New("Something failed")          // want `must not be capitalized`
	_ = errors.New("failed.")                   // want `must not end with punctuation`
	_ = fmt.Errorf("Oops: %w", errors.New("x")) // want `must not be capitalized`
}

func good() {
	_ = errors.New("something failed")
	_ = errors.New("HTTP request failed")
}
