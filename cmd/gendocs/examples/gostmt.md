```go
// bad
go worker()

// good: waited on, errors returned
var group errgroup.Group
group.Go(worker)
err := group.Wait()
```
