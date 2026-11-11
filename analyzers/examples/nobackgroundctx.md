```go
// bad
func fetch() error {
	ctx := context.Background()
	return call(ctx)
}

// good
func fetch(ctx context.Context) error {
	return call(ctx)
}
```
