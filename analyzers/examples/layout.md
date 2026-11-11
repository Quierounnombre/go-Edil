```go
// bad: type and constant next to the logic (store.go)
type Store struct{}
const LIMIT = 10

// good
// type.go:    type Store struct{}
// defines.go: const LIMIT = 10
```
