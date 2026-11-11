```go
// bad
db.Query("SELECT * FROM users WHERE id = " + id)

// good
db.Query("SELECT * FROM users WHERE id = ?", id)
```
