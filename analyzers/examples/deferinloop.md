```go
// bad: files stay open until the function returns
for _, name := range names {
	file, _ := os.Open(name)
	defer file.Close()
}

// good
for _, name := range names {
	err := process(name)
	if err != nil {
		return err
	}
}
```
