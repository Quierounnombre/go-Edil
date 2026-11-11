```go
// bad
file, _ := os.Open(path)

if err != nil {
	log.Println(err)
}

// good
file, err := os.Open(path)
if err != nil {
	return err
}
```
