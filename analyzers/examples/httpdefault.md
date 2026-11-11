```go
// bad: never times out
resp, err := http.Get(url)

// good
client := &http.Client{Timeout: HTTP_TIMEOUT}
resp, err := client.Get(url)
```
