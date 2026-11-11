```go
// bad: never checked on a CI that is not windows
//go:build windows

// bad: file names that select a platform
store_windows.go

// bad: hidden from ./...
_old/, .tmp/, vendor/, a nested go.mod

// good: one build, every file visible
package store
```
