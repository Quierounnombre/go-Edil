# Conventions

Generated from the analyzers, do not edit by hand: `go run ./cmd/gendocs`.

These are the rules the code is checked against. There is no way to suppress a finding: fix the code. Only `_test.go` files are exempt.

Code style: tabs for indentation, constants in `defines.go`, types in `type.go`.
{{range .Rules}}
## {{.Name}}

{{.Doc}}

{{.Example}}
{{end}}