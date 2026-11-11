```go
// bad
func parse(s string) (n int, err error) {
	n, err = strconv.Atoi(s)
	return
}

// good
func parse(s string) (int, error) {
	return strconv.Atoi(s)
}
```
