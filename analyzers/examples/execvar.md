```go
// bad: user input reaches a shell
exec.Command("sh", "-c", "ls "+dir)

// good: constant binary, arguments separate
exec.Command("ls", dir)
```
