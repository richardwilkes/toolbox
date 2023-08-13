[![Go Reference](https://pkg.go.dev/badge/github.com/richardwilkes/toolbox.svg)](https://pkg.go.dev/github.com/richardwilkes/toolbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/richardwilkes/toolbox)](https://goreportcard.com/report/github.com/richardwilkes/toolbox)

# toolbox
Toolbox for Go.

To install this package and the tools it provides:
```
./build.sh
```

> NOTE: This library already had a v1.x.y when Go modules were first introduced. Due to this, it doesn't follow the
>       normal convention and instead treats its releases as if they are of the v0.x.y variety (i.e. it could introduce
>       breaking API changes). Keep this in mind when deciding whether or not to use it.

## Package summaries

### toolbox
Utilities that didn't have a home elsewhere.

### toolbox/atexit
Provides functionality similar to the C standard library's atexit() call. To function properly, use
`atexit.Exit(result)` rather than `os.Exit(result)`.

### toolbox/cmdline
Command line handling. Provides the tool `genversion` for generating version numbers with an embedded date.

### toolbox/collection
Provides type-safe sets for the various primitive types.

### toolbox/collection/dict
Provides some useful `map` functions that were inexplicably left out of the `maps` package introduced in Go 1.21.

### toolbox/collection/quadtree
Provides an implementation of a [Quadtree](https://en.wikipedia.org/wiki/Quadtree).

### toolbox/collection/slice
Provides some useful `slice` functions that were inexplicably left out of the `slices` package introduced in Go 1.21.

### toolbox/desktop
Desktop integration utilities.

### toolbox/errs
Errors that contain stack traces with source locations, along with nested causes, if any.

### toolbox/eval
Dynamically evaluate expressions.

### toolbox/formats/icon
Provides image scaling and stacking utilities.

### toolbox/formats/icon/icns
Provides an encoder for [Apple Icon Image](https://en.wikipedia.org/wiki/Apple_Icon_Image_format) files.

### toolbox/formats/icon/ico
Provides an encoder for [Windows Icon Image](https://en.wikipedia.org/wiki/ICO_(file_format)) files.

### toolbox/formats/json
Manipulation of JSON data.

### toolbox/formats/xlsx
Extract text from Excel spreadsheets.

### toolbox/i18n
Internationalization support for applications. Provides the tool `go-i18n` for generating a template for a localization
file from source code.

### toolbox/log/jot
Simple asynchronous logging.

#### Sample usage:
```go
package main

import "github.com/richardwilkes/toolbox/log/jot"

func main() {
    defer jot.Flush()
    jot.Debug("Debug level")
    jot.Debugf("Debug level with %s", "args")
    jot.Info("Info level")
    jot.Infof("Info level with %s", "args")
    jot.Warn("Warning level")
    jot.Warnf("Warning level with %s", "args")
    jot.Error("Error level")
    jot.Errorf("Error level with %s", "args")
    jot.Fatal(1, "Fatal level")
    jot.Fatalf(1, "Fatal level with %s", "args")    // Will never be reached due to previous line
}
```

### toolbox/log/jotrotate
Provides a pre-canned way to add jot logging with file rotation, along with command-line options for controlling it.

### toolbox/log/logadapter
This package defines an API to use for logging, which actual logging implementations can implement directly or provide
an adapter to use.

It also provides an implementation that just discards data given to it as well as an implementation that wraps another
logger and prefixes all output.

### toolbox/log/rotation
Provides file rotation when files hit a given size.

### toolbox/notifier
Provides a mechanism for tracking targets of notifications and methods for notifying them.

### toolbox/rate
Rate limiting which supports a hierarchy of limiters, each capped by their parent.

### toolbox/softref
Soft references.

### toolbox/taskqueue
Provides a simple asynchronous task queue.

### toolbox/txt
Text utilities.

### toolbox/vcs/git
git repository access.

### toolbox/xcrypto
Provides convenience utilities for encrypting and decrypting streams of data with public & private keys.

### toolbox/xio
io utilities.

### toolbox/xio/fs
Filesystem utilities.

### toolbox/xio/fs/paths
Platform-specific standard paths.

### toolbox/xio/fs/safe
Safe, atomic saving of files.

### toolbox/xio/fs/tar
Provides tar file extraction utilities.

### toolbox/xio/fs/zip
Provides zip file extraction utilities.

### toolbox/xio/network
Network-related utilities.

### toolbox/xio/network/natpmp
Implementation of [NAT-PMP](https://tools.ietf.org/html/rfc6886).

### toolbox/xio/network/xhttp
HTTP-related utilities.

### toolbox/xio/network/xhttp/web
Web server with some standardized logging and handler wrapping.

### toolbox/xio/term
Terminal utilities.

### toolbox/xmath
Math utilities.

### toolbox/xmath/crc
Provides CRC helpers.

### toolbox/xmath/fixed
Fixed-point types with a configurable number of decimal places. These types implement the marshal/unmarshal interfaces
for JSON and YAML.

### toolbox/xmath/geom
Geometry primitives.

### toolbox/xmath/geom/poly
Provides polygon boolean operations. These are not as robust as I'd like.

### toolbox/xmath/geom/visibility
Calculates a visibility polygon from a given point in the presence of a set of obstructions, also known as an [Isovist](https://en.wikipedia.org/wiki/Isovist).

### toolbox/xmath/num
128-bit int and uint types. These types implement the marshal/unmarshal interfaces for JSON and YAML.

### toolbox/xmath/rand
Randomizer based upon the crypto/rand package.
