# Copilot Instructions for Toolbox

## Project Overview

This is a comprehensive Go utility library (v2) providing extended functionality for common patterns. The codebase follows a strict "x-prefix" convention where packages extending standard library functionality are prefixed with "x" (e.g., `xmath`, `xio`, `xhttp`).

## Key Architecture Patterns

### Generics-First Design

-   Heavy use of Go generics with type constraints (see `xmath.Numeric` interface)
-   Generic data structures like `Point[T xmath.Numeric]`, `Tree[K, V any]`, `Int[T fixed.Dx]`
-   Type conversion functions like `ConvertPoint[T, F xmath.Numeric](pt Point[F]) Point[T]`

### Error Handling Convention

-   Use `errs.Error` for structured errors with stack traces instead of basic Go errors
-   All errors include source location and can be chained: `errs.NewWithCause("message", causeErr)`
-   The `errs` package provides `ErrorWrapper` and `StackError` interfaces
-   Use `xos.PanicRecovery()` for panic handling in defer statements

### Testing Patterns

-   Use the internal `check` package instead of raw testing.T: `c := check.New(t)`
-   Testing methods: `c.Equal(expected, actual)`, `c.HasPrefix()`, `c.NoError()`, `c.HasError()`
-   Example tests follow `ExampleFunction` naming for documentation
-   All test files use `package_test` naming convention

### Fixed-Point Arithmetic

-   `fixed64.Int[T]` and `fixed128.Int[T]` for precise decimal calculations
-   Always use `.Mul()` and `.Div()` methods, never direct operators for multiplication/division
-   Different precision types via type parameters: `fixed64.Int[fixed.D2]` for 2 decimal places

### Command-Line Applications

-   Use `xflag.SetUsage()` for rich CLI help with version info and formatting
-   Set app metadata in main(): `xos.AppIdentifier`, `xos.CopyrightStartYear`, etc.
-   Use `i18n.Text()` for all user-facing strings for localization support

## Build and Development

### Build System

```bash
./build.sh --all    # Full build with linting and race-checked tests
./build.sh --lint   # Run golangci-lint (auto-installs latest version)
./build.sh --test   # Run tests
./build.sh --race   # Run tests with race detection
```

### Dependencies

-   Minimal external dependencies: only `golang.org/x/*` and `gopkg.in/yaml.v3`
-   Self-contained utilities - avoid adding new external deps without strong justification

## Package-Specific Notes

### Geometry (`geom/`)

-   All geometric types are generic: `Point[T]`, `Rect[T]`, `Size[T]` with `xmath.Numeric` constraint
-   Conversion between numeric types via `ConvertPoint`, `ConvertRect`, etc.

### Collections (`collection/`)

-   Red-black tree: `redblack.New[K, V any](compareFunc)` - requires comparison function
-   Bitset and quadtree implementations available

### Extended Standard Library

-   `xmath`: Generic math functions with proper type constraints
-   `xio`: BOM stripping, safe file operations
-   `xstrings`: Enhanced string utilities
-   `xhttp`: HTTP utilities with metadata and compression helpers

### 128-bit Numbers (`num128/`)

-   `Int` and `UInt` types for 128-bit arithmetic
-   Always check for overflow in mathematical operations
-   String conversion methods handle different bases

## File Organization

-   Each package should have comprehensive tests (`*_test.go`)
-   Example tests for public APIs to serve as documentation
-   License header on all files (Mozilla Public License 2.0)
-   Internal utilities in appropriate x-prefixed packages
