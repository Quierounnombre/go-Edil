# Edil

Aedile of the codebase: keeps the city in order. Custom Go analyzers for code style, readability and a few security checks, reproducible with Nix.

* * *

## Rules

Generated from the analyzers in code, so this table cannot go stale.

| Analyzer | Enforces |
|----------|----------|
{{range .Rules}}| `{{.Name}}` | {{.Doc}} |
{{end}}
Limits (function length, parameters, complexity, nesting) live in each analyzer's `defines.go`.

* * *

## Tutorial

### 1. Get the binary

With Nix, from any machine:

```
nix build {{.Repo}}
```

Or from a clone of the repo, without Nix:

```
go build -o edil ./cmd/edil
```

### 2. Run it on a project

Edil runs as a vet tool. From the root of the Go project you want to check:

```
go vet -vettool=$(readlink -f result/bin/edil) ./...
```

Use the path of your own binary if you built it with `go build`.

### 3. Read the output

Each finding is one line, in the usual `file:line:column: message` form:

```
store.go:12:7: constant "maxRetries" must be SCREAMING_SNAKE_CASE
```

Fix the code, or narrow the rule if the rule is wrong. Do not disable a whole analyzer.

When one line is a deliberate exception, say so on that line or the line above, with a reason:

```
//edil:ignore mathrand jitter only, not security
```

The rule name and the reason are both required. Without a reason the finding stays.

### 4. Add it to the project gate

In the flake of the project you are checking:

```
inputs.edil.url = "{{.Repo}}";
inputs.edil.inputs.nixpkgs.follows = "nixpkgs";

# add to your tools list
edil.packages.${system}.default

# inside checks.lint
go vet -vettool=$(command -v edil) ./...
```

> [!TIP]
> While developing both repos side by side, use `path:../edil` instead of the repo URL.

### 5. Add a rule to Edil

1. Create `analyzers/NAME/NAME.go` exporting `var Analyzer *analysis.Analyzer`. Constants go in `analyzers/NAME/defines.go`.
2. Start the file loop with `kit.Skip(pass, file)` and report with `kit.Reportf(pass, pos, ...)`, so tests, generated files and `//edil:ignore` behave the same in every rule.
3. Add it to the `ANALYZERS` list in `analyzers/all.go`.
4. Refresh this file, then commit:

```
go run ./cmd/gendocs
```

> [!NOTE]
> If you skip the refresh, `nix flake check` fails and tells you the command to run.

Edil must pass its own rules. Run it on this repo before committing:

```
go vet -vettool=$(readlink -f result/bin/edil) ./...
```

* * *

## Troubleshooting

**Nix says the build hash is wrong**

Expected on first build. `vendorHash` starts as `lib.fakeHash`. Paste the real hash from the error, then rebuild. Repeat whenever `go.mod` changes.

**Nix cannot find my files**

Flakes only see tracked files:

```
git add -A
```

**Sandbox fails to download Go modules**

The Nix sandbox has no network. In the consuming project, run this and commit `vendor/`:

```
go mod vendor
```

**A rule flags too much**

Narrow it in its source, rebuild, and re-run on a real project. For a single justified line, use `//edil:ignore` with a reason.

* * *

Status: experimental. Rules change as they meet real code.
