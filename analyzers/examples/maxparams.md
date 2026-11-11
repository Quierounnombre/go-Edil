```go
// bad
func connect(host string, port int, user string, pass string, db string, tls bool) {}

// good
func connect(opts Options) {}
```
