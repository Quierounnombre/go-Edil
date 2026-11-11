```go
// bad
time.Sleep(30 * time.Second)
if role == "admin" {
}
os.Exit(1)

// good
const POLL_INTERVAL = 30 * time.Second
const ROLE_ADMIN = "admin"
const EXIT_FAILURE = 1
```
