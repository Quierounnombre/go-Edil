```go
// bad
cfg := &tls.Config{InsecureSkipVerify: true}

// good
cfg := &tls.Config{MinVersion: tls.VersionTLS12}
```
