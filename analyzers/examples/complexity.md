```go
// bad: nested branches and loops pile up score
func handle(items []Item) {
	for _, item := range items {
		if item.Valid {
			for _, tag := range item.Tags {
				if tag != "" && tag != "skip" {
				}
			}
		}
	}
}

// good: one idea per function
func handle(items []Item) {
	for _, item := range items {
		handleItem(item)
	}
}
```
