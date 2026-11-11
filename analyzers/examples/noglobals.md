```go
// bad
var cache = map[string]int{}

// good: pass state in
func newCache() map[string]int {
	return map[string]int{}
}
```
