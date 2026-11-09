# Edil

Aedile of the codebase: keeps the city in order. Custom Go analyzers for syntax form and readability, reproducible with Nix.

---

## Rules

Generated from the analyzers in code, so this table cannot go stale.

| Analyzer | Enforces |
|---|---|
<!-- RULES -->

---

## Tutorial

### 1. Get the binary

With Nix, from any machine:

```sh
nix build @REPO@
```

Or from a clone of the repo, without Nix:

```sh
go build -o edil ./cmd/edil
```

### 2. Run it on a project

Edil runs as a vet tool. From the root of the Go project you want to check:

```sh
go vet -vettool=$(readlink -f result/bin/edil) ./...
```

Use the path of your own binary if you built it with `go build`.

### 3. Read the output

Each finding is one line, in the usual `file:line:column: message` form:

```
store.go:12:7: constant "maxRetries" must be SCREAMING_SNAKE_CASE
```

Fix the code, or narrow the rule if the rule is wrong. Do not disable a whole analyzer.

### 4. Add it to the project gate

In the flake of the project you are checking:

```nix
inputs.edil.url = "@REPO@";
inputs.edil.inputs.nixpkgs.follows = "nixpkgs";

# add to your tools list
edil.packages.${system}.default

# inside checks.lint
go vet -vettool=$(command -v edil) ./...
```

> [!TIP]
> While developing both repos side by side, use `path:../edil` instead of the repo URL.

### 5. Add a rule to Edil

1. Create `analyzers/NAME/NAME.go` exporting `var Analyzer *analysis.Analyzer`.
2. Add it to the `ANALYZERS` list in `analyzers/all.go`.
3. Refresh this file `with go run ./cmd/gendocs/main.go`, then commit:

```sh
go run ./cmd/gendocs
```

> [!NOTE]
> If you skip the refresh, `nix flake check` fails and tells you the command to run.

---

## Troubleshooting

<details>
<summary><b>Nix says the build hash is wrong</b></summary>

Expected on first build. `vendorHash` starts as `lib.fakeHash`. Paste the real hash from the error, then rebuild. Repeat whenever `go.mod` changes.

</details>

<details>
<summary><b>Nix cannot find my files</b></summary>

Flakes only see tracked files:

```sh
git add -A
```

</details>

<details>
<summary><b>Sandbox fails to download Go modules</b></summary>

The Nix sandbox has no network. In the consuming project, run this and commit `vendor/`:

```sh
go mod vendor
```

</details>

<details>
<summary><b>A rule flags too much</b></summary>

Narrow it in its source, rebuild, and re-run on a real project.

</details>

---

<sub>Status: experimental. Rules change as they meet real code.</sub>
