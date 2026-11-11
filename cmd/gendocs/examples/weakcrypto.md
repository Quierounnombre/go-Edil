```go
// bad
sum := md5.Sum(data)

// good
sum := sha256.Sum256(data)
```
