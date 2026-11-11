```go
// bad (edil -fix does this)
if err != nil {
	return err
} else {
	use(value)
}

// good
if err != nil {
	return err
}
use(value)
```
