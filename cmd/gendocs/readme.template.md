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

There is no way to suppress a finding. Fix the code; if a rule is wrong for a real case, change the rule in Edil. Only `_test.go` files are exempt. The full rule list with examples is in [CONVENTIONS.md](CONVENTIONS.md).

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
2. Start the file loop with `kit.Skip(pass, file)` and report with `kit.Reportf(pass, pos, ...)`, so the test-file exemption behaves the same in every rule.
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

## Tools

* `edil -explain` lists every rule; `edil -explain RULE` prints its description and a bad/good example.
* `edil -fix ./...` applies the mechanical fixes (errwrap, varinit, noelse, constcase, tabindent). Run gofmt afterwards.
* `git config core.hooksPath cmd/fix/hooks` runs gofmt, Edil and the docs check before every commit.
* `cmd/fix/fixloop.sh` runs an agent until Edil is clean and logs each round to `~/.local/state/edil/metrics.tsv`; `cmd/fix/metrics.sh` summarizes it. Round 1 counts are the rules the agent breaks first: put those in its prompt.

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

Narrow it in its source, rebuild, and re-run on a real project. Suppression comments do not exist; the rule is the thing to fix.

* * *

Status: experimental. Rules change as they meet real code.
