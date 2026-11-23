```go
// bad
cfg := defaultConfig()
for _, path := range paths {
	cfg, err := merge(cfg, path) // a new cfg; the outer one never changes
	if err != nil {
		return err
	}
	log.Println(cfg.Name)
}

// good
cfg := defaultConfig()
for _, path := range paths {
	merged, err := merge(cfg, path)
	if err != nil {
		return err
	}
	cfg = merged
}
```
