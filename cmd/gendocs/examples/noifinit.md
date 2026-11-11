```go
// bad
if err := run(); err != nil {
	return err
}

// good
err := run()
if err != nil {
	return err
}
```
