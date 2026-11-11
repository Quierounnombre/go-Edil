package errwrap

import "fmt"

func bad(err error) error {
	_ = fmt.Errorf("a: %v", err)         // want `use %w`
	_ = fmt.Errorf("b: %s", err.Error()) // want `use %w`
	return fmt.Errorf("c: %s", err)      // want `use %w`
}

func good(err error) error {
	_ = fmt.Errorf("a: %w", err)
	_ = fmt.Errorf("n: %d %s", 3, "x")
	return fmt.Errorf("100%% %w", err)
}
