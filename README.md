# toolbox

[![Go Reference](https://pkg.go.dev/badge/github.com/richardwilkes/toolbox/v2.svg)](https://pkg.go.dev/github.com/richardwilkes/toolbox/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/richardwilkes/toolbox/v2)](https://goreportcard.com/report/github.com/richardwilkes/toolbox/v2)

Toolbox for Go.

Contains a wide variety of code I've found useful in my own projects over the years. For cases where code exists to help
use standard library code, the package has been named the same as the standard library one, but with a preceeding "x"
(for extended). This allows both to be used in the same file without having to do import renaming.

## Package Overview

### Core Utilities

- **`check`** - Enhanced testing utilities with a fluent API that provides more informative error messages and better test assertions than the standard testing package.

- **`errs`** - Structured error handling with stack traces, error chaining, and detailed error objects that provide source locations and nested causes for better debugging.

- **`i18n`** - Internationalization support for applications, providing localization capabilities for user-facing text and messages.

- **`notifier`** - Event notification system for implementing the observer pattern, allowing objects to register for and receive notifications about events.

- **`tid`** - Thread-safe unique identifier generation using cryptographically secure random values encoded in base64.

### Mathematical and Numerical

- **`xmath`** - Extended mathematical functions with generic type constraints (e.g., `Numeric` interface) that work across integer and floating-point types.

- **`num128`** - 128-bit integer arithmetic with signed (`Int`) and unsigned (`UInt`) types for high-precision calculations that exceed native Go integer limits.

- **`fixed`** - Fixed-point decimal arithmetic for precise financial and monetary calculations, with separate packages for 64-bit (`fixed64`) and 128-bit (`fixed128`) precision.

- **`eval`** - Expression evaluator that can parse and evaluate mathematical expressions with variables, supporting both fixed-point and floating-point arithmetic.

### Geometry and Graphics

- **`geom`** - Generic geometric primitives including `Point[T]`, `Rect[T]`, `Size[T]`, `Line[T]`, `Matrix[T]`, and `Insets[T]` with type conversion utilities.

- **`geom/poly`** - Polygon boolean operations (⚠️ experimental - not recommended for production use).

- **`geom/visibility`** - Visibility calculations for geometric objects (⚠️ experimental - not recommended for production use).

### Data Structures and Collections

- **`collection/bitset`** - Efficient bitset implementation for working with sets of bits.

- **`collection/quadtree`** - Spatial data structure for efficient 2D range queries and collision detection.

- **`collection/redblack`** - Generic red-black tree implementation providing balanced binary search tree functionality.

### Extended Standard Library

- **`xbytes`** - Extended byte manipulation utilities including BOM (Byte Order Mark) handling and buffer utilities.

- **`xcrc64`** - Extended CRC64 checksum utilities for data integrity verification.

- **`xcrypto`** - Cryptographic utilities including stream encryption/decryption helpers.

- **`xfilepath`** - Extended file path utilities with filename manipulation, root detection, and cross-platform path splitting.

- **`xflag`** - Enhanced command-line flag parsing with rich usage formatting, automatic version flags, and post-parse function handling.

- **`xhash`** - Extended hashing utilities for creating consistent hash values across different data types.

- **`xhttp`** - HTTP utilities including basic authentication, gzip handling, metadata management, file retrieval, and server helpers.

- **`ximage`** - Image processing utilities with support for various image formats.

- **`xio`** - Extended I/O utilities including BOM stripping, safe file closing, and specialized readers/writers.

- **`xjson`** - Enhanced JSON handling utilities for parsing and marshaling with additional features.

- **`xnet`** - Network utilities for address manipulation and network-related operations.

- **`xos`** - Operating system utilities including application information, browser launching, filesystem operations, panic recovery, safe file handling, task queues, and user information.

- **`xrand`** - Extended random number generation utilities for cryptographically secure randomness.

- **`xreflect`** - Reflection utilities for working with Go's reflection API more effectively.

- **`xruntime`** - Runtime utilities including detailed stack trace generation with source location information.

- **`xslices`** - Enhanced slice manipulation utilities including column-based sorting and other advanced slice operations.

- **`xslog`** - Enhanced structured logging with a "pretty" formatter that provides colorful output, stack trace formatting, and improved readability.

- **`xstrings`** - String manipulation utilities including case conversion, text wrapping, natural sorting, emoji handling, capitalization, space collapsing, and various string processing functions.

- **`xsync`** - Synchronization utilities extending Go's sync package with additional concurrent programming tools.

- **`xterm`** - Terminal utilities for ANSI color codes, terminal detection, and formatted output with cross-platform compatibility.

- **`xtime`** - Time manipulation utilities extending Go's time package with additional date/time functionality.

- **`xunicode`** - Unicode utilities for advanced text processing and character manipulation.

- **`xyaml`** - YAML processing utilities for parsing and marshaling YAML data with enhanced features.

### Specialized Utilities

- **`rate`** - Rate limiting with hierarchical limiters, where each limiter can be capped by its parent, useful for implementing tiered rate limiting.

- **`softref`** - Soft reference implementation for memory management, allowing resources to be garbage collected when memory pressure occurs.

### Command Line Tools

- **`cmd/i18n`** - Command-line tool for extracting and managing internationalization strings from Go source code.

All packages follow consistent patterns with use of Go generics for type safety where appropriate, comprehensive error handling with the `errs` package, and thorough testing with the `check` package. The "x" prefix convention allows seamless use alongside standard library packages without import conflicts.
