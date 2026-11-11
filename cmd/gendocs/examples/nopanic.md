```go
// bad
func load() {
	panic("missing config")
}

// good
func load() error {
	return errors.New("missing config")
}
```
