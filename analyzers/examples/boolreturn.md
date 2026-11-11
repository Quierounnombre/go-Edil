```go
// bad
func isReady() int {
	return 1
}

if ready == true {
}

// good
func isReady() bool {
	return true
}

if ready {
}
```
