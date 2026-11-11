```go
// bad
password := "hunter2hunter2"

// good
password := os.Getenv("DB_PASSWORD")
```
