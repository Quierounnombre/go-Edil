# Conventions

Generated from the analyzers, do not edit by hand: `go run ./cmd/gendocs`.

These are the rules the code is checked against. There is no way to suppress a finding: fix the code. Only `_test.go` files are exempt.

Code style: tabs for indentation, constants in `defines.go`, types in `type.go`.

## boolreturn

flags are bool: no 0/1 returns, no comparing against true or false

```go
// bad
func isReady() int {
	return 1
}

if ready == true {
}

// good
func isReady() bool {
	return true
}

if ready {
}
```

## complexity

cognitive complexity per function must stay at or below MAX_COMPLEXITY (15)

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

## constcase

constants must be SCREAMING_SNAKE_CASE

```go
// bad
const maxRetries = 3

// good
const MAX_RETRIES = 3
```

## deferinloop

no defer inside a loop; it runs when the function returns, not per iteration

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

## errreturn

errors must be returned, not assigned to _ or checked and swallowed

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

## errstring

error strings are lowercase and do not end with punctuation or a newline

```go
// bad
return errors.New("Failed to open file.")

// good
return errors.New("failed to open file")
```

## errwrap

fmt.Errorf must wrap errors with %w, not format them with %v or %s

```go
// bad: the chain is lost
return fmt.Errorf("load config: %v", err)

// good (edil -fix does this)
return fmt.Errorf("load config: %w", err)
```

## execvar

exec.Command runs a constant binary; a shell must never get a non-constant script

```go
// bad: user input reaches a shell
exec.Command("sh", "-c", "ls "+dir)

// good: constant binary, arguments separate
exec.Command("ls", dir)
```

## funclen

functions must not exceed MAX_FUNC_LINES (120) lines

```go
// bad: a 200 line function

// good: split it into named steps
func run() error {
	cfg := load()
	return serve(cfg)
}
```

## gostmt

no go statements outside main; use errgroup.Group.Go or sync.WaitGroup.Go so goroutines are waited on and report errors

```go
// bad
go worker()

// good: waited on, errors returned
var group errgroup.Group
group.Go(worker)
err := group.Wait()
```

## logsecret

do not pass password, token or secret values to log, slog or fmt print calls

```go
// bad
log.Printf("login %s", password)

// good
log.Printf("login %s", userID)
```

## hardcodedsecret

no secret-looking string literals assigned to password, token, secret or key names; load them from the environment

```go
// bad
password := "hunter2hunter2"

// good
password := os.Getenv("DB_PASSWORD")
```

## httpdefault

no http.Get/Post/Head/PostForm, http.DefaultClient or http.Client{} without Timeout; they never time out

```go
// bad: never times out
resp, err := http.Get(url)

// good
client := &http.Client{Timeout: HTTP_TIMEOUT}
resp, err := client.Get(url)
```

## layout

types must live in type.go, constants in defines.go

```go
// bad: type and constant next to the logic (store.go)
type Store struct{}
const LIMIT = 10

// good
// type.go:    type Store struct{}
// defines.go: const LIMIT = 10
```

## mathrand

math/rand is predictable; use crypto/rand

```go
// bad
import "math/rand"

// good
import "crypto/rand"
```

## maxdepth

control flow must not nest deeper than MAX_DEPTH (3) levels

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

## maxparams

functions take at most MAX_PARAMS (5) parameters, context.Context not counted; group the rest in a struct

```go
// bad
func connect(host string, port int, user string, pass string, db string, tls bool) {}

// good
func connect(opts Options) {}
```

## nakedreturn

functions with named results must not use a bare return

```go
// bad
func parse(s string) (n int, err error) {
	n, err = strconv.Atoi(s)
	return
}

// good
func parse(s string) (int, error) {
	return strconv.Atoi(s)
}
```

## nobackgroundctx

context.Background and context.TODO only in main; accept a ctx parameter everywhere else

```go
// bad
func fetch() error {
	ctx := context.Background()
	return call(ctx)
}

// good
func fetch(ctx context.Context) error {
	return call(ctx)
}
```

## nobuildtag

no build constraints, GOOS/GOARCH file names, excluded files or hidden directories: all code must be visible to the linter

```go
// bad: never checked on a CI that is not windows
//go:build windows

// bad: file names that select a platform
store_windows.go

// bad: hidden from ./...
_old/, .tmp/, vendor/, a nested go.mod

// good: one build, every file visible
package store
```

## noelse

no else after a block that returns, breaks or continues

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

## noglobals

no package-level var, except //go:embed, _ assertions and analyzers; use const or pass values in

```go
// bad
var cache = map[string]int{}

// good: pass state in
func newCache() map[string]int {
	return map[string]int{}
}
```

## noifinit

if and switch must not have an init statement; declare it on its own line

```go
// bad
if err := run(); err != nil {
	return err
}

// good
err := run()
if err != nil {
	return err
}
```

## nomagic

no magic numbers (0 and 1 are fine, except as call arguments or in ==/!=); no magic strings that are compared or repeated; use named constants

```go
// bad
time.Sleep(30 * time.Second)
if role == "admin" {
}
os.Exit(1)

// good
const POLL_INTERVAL = 30 * time.Second
const ROLE_ADMIN = "admin"
const EXIT_FAILURE = 1
```

## nopanic

forbid panic, log.Fatal*, log.Panic* and os.Exit outside main; return errors instead

```go
// bad
func load() {
	panic("missing config")
}

// good
func load() error {
	return errors.New("missing config")
}
```

## nosuppress

no lint suppression comments (nolint, nosec, lint:ignore, edil:ignore); Edil does not honor them, fix the code or the rule

```go
// bad
//nolint:errcheck
//edil:ignore nomagic

// good: fix the code, or ask for the rule to change
```

## pkgsize

a package must stay under MAX_TYPES (10) types and MAX_LINES (1500) lines; split it

```go
// bad: one package with 15 types and 3000 lines

// good: split by responsibility into packages
```

## sqlconcat

database/sql queries must be constant strings; pass values as placeholder arguments

```go
// bad
db.Query("SELECT * FROM users WHERE id = " + id)

// good
db.Query("SELECT * FROM users WHERE id = ?", id)
```

## tabindent

indentation must use tabs, not spaces

```go
// bad (edil -fix does this)
func run() {
    work()
}

// good
func run() {
	work()
}
```

## tlsinsecure

tls.Config.InsecureSkipVerify must not be set to true

```go
// bad
cfg := &tls.Config{InsecureSkipVerify: true}

// good
cfg := &tls.Config{MinVersion: tls.VersionTLS12}
```

## varinit

inside functions use :=; var only for zero values that are not assigned right after

```go
// bad (edil -fix does this)
var name = "x"
var total int
total = compute()

// good
name := "x"
total := compute()
```

## weakcrypto

no md5, sha1, des or rc4; use sha256+, AES-GCM or ChaCha20

```go
// bad
sum := md5.Sum(data)

// good
sum := sha256.Sum256(data)
```
