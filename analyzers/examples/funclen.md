```go
// bad: a 200 line function

// good: split it into named steps
func run() error {
	cfg := load()
	return serve(cfg)
}
```
