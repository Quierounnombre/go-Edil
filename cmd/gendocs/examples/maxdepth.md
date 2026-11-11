```go
// bad: 4 levels
for _, a := range as {
	if a.ok {
		for _, b := range a.bs {
			if b.ok {
			}
		}
	}
}

// good: guard clauses and helpers
for _, a := range as {
	if !a.ok {
		continue
	}
	visit(a.bs)
}
```
