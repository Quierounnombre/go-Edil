```go
// bad: the chain is lost
return fmt.Errorf("load config: %v", err)

// good (edil -fix does this)
return fmt.Errorf("load config: %w", err)
```
