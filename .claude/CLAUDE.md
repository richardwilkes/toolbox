# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A general-purpose Go utility library (module `github.com/richardwilkes/toolbox/v2`, Go 1.26+). It collects code the author reuses across projects. Packages that extend a standard library package are named the same with an `x` prefix (`xio`, `xhttp`, `xmath`, ...) so both can be imported in one file without renaming. Dependencies are intentionally minimal — only `golang.org/x/*` and `gopkg.in/yaml.v3`. Avoid adding new external dependencies without strong justification.

## Build & test

```bash
./build.sh            # go build ./... only
./build.sh --all      # build + lint + race-checked tests (the full check)
./build.sh --lint     # golangci-lint (auto-installs the pinned version into $GOPATH/bin)
./build.sh --test     # go test ./...
./build.sh --race     # go test -race ./...
```

Run a single test directly with the toolchain:

```bash
go test ./xstrings/ -run TestName -v
```

`build.sh` also `go install`s the one executable, `./cmd/i18n` (extracts `i18n.Text(...)` calls into localization template files).

## Conventions to follow

These are enforced by lint config (`.golangci.yml`) and repository style — match them in new code.

- **License header**: every `.go` file starts with the MPL-2.0 header block (copy it from any existing file, e.g. `check/check.go`). New files must include it.
- **Tests use the `check` package, not raw `testing.T`**: `c := check.New(t)` then `c.Equal(expected, actual)`, `c.NoError(err)`, `c.HasError(err)`, `c.HasPrefix(...)`, etc. Test files use the `package foo_test` (external) naming. Add `ExampleXxx` functions for public APIs — they double as documentation.
- **Errors**: use the `errs` package for structured errors with stack traces — `errs.New(...)`, `errs.NewWithCause(msg, cause)`, `errs.Wrap(err)` — not bare `fmt.Errorf`. Use `errs.As` / the `ErrorWrapper` / `StackError` interfaces for inspection. Use `xos.PanicRecovery()` in `defer` for panic handling.
- **errcheck is strict**: `check-blank` and `check-type-assertions` are on. Don't discard errors with `_`, and guard type assertions with the comma-ok form.

## Architecture notes

- **Generics-first.** Numeric and geometric types are generic over `xmath.Numeric`: `Point[T]`, `Rect[T]`, `Size[T]`, `Matrix[T]`, with conversion helpers like `ConvertPoint[T, F]`. Collections are generic too (`redblack.New[K, V]` takes a compare func). When extending these, keep the type-parameter style consistent rather than adding concrete-typed variants.
- **Fixed-point arithmetic** (`fixed`, `fixed64`, `fixed128`): precision is a type parameter — `fixed64.Int[fixed.D2]` is 2 decimal places. Use the `.Mul()` / `.Div()` methods for multiply/divide; the bare `*` and `/` operators give wrong results for fixed-point values.
- **128-bit integers** (`num128`): `Int` (signed) and `UInt` (unsigned) for arithmetic beyond native 64-bit. Watch for overflow.
- **CLI apps**: use `xflag.SetUsage()` for help/version formatting, set app metadata (`xos.AppIdentifier`, `xos.CopyrightStartYear`, ...) in `main()`, and wrap every user-facing string in `i18n.Text(...)` so `cmd/i18n` can extract it.

For a per-package summary, see `README.md`; deeper architecture and pattern notes are in `.github/copilot-instructions.md`.
