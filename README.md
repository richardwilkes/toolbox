# toolbox
Toolbox for Go.

To install this package and the tools it provides:
```
go get -u github.com/richardwilkes/toolbox/...
```

## atexit
Provides functionality similar to the C standard library's atexit() call. To
function properly, use `atexit.Exit(result)` rather than `os.Exit(result)`.

## cmdline
Command line handling. Provides the tool `genversion` for generating version
numbers with an embedded date.

## collection
Provides type-safe sets for the various primitive types.

## desktop
Desktop integration utilities.

## errs
Errors that contain stack traces with source locations, along with nested
causes, if any.

## formats/json
Manipulation of JSON data.

## formats/xlsx
Extract text from Excel spreadsheets.

## i18n
Internationalization support for applications. Provides the tool `go-i18n` for
generating a template for a localization file from source code.

## log/jot
Simple asynchronous logging.

### Sample usage:

```Go
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

## log/jotrotate
Provides a pre-canned way to add jot logging with file rotation, along with
command-line options for controlling it.

## log/logadapter
This package defines an API to use for logging, which actual logging
implementations can implement directly or provide an adapter to use.

It also provides an implementation that just discards data given to it as
well as an implementation that wraps another logger and prefixes all output.

## log/rotation
Provides file rotation when files hit a given size.

## rate
Rate limiting which supports a hierarchy of limiters, each capped by their
parent.

## taskqueue
Provides a simple asynchronous task queue.

## txt
Text utilities.

## vcs/git
git repository access

## xio
io utilities.

## xio/fs
Filesystem utilities.

## xio/fs/embedded
Provides an implementation of an embedded filesystem.

## xio/fs/embedded/htmltmpl
Provides convenience utilities for using html templates in an embedded filesystem.

## xio/fs/embedded/texttmpl
Provides convenience utilities for using text templates in an embedded filesystem.

## xio/fs/mkembeddedfs
Tool for generating the embedded filesystem. Note: this utility has been
deprecated in favor of using `embedded.NewFileSystemFromEmbeddedZip()` instead

## xio/fs/paths
Platform-specific standard paths.

## xio/fs/safe
Safe, atomic saving of files.

## xio/fs/zip
Simple zip extraction.

## xio/network
Network-related utilities.

## xio/network/natpmp
Implementation of NAT-PMP. See https://tools.ietf.org/html/rfc6886

## xio/network/xhttp
HTTP-related utilities.

## xio/network/xhttp/web
Web server with some standardized logging and handler wrapping.

## xio/term
Terminal utilities.

## xmath
Math utilities.

## xmath/fixed
Fixed-point types of varying sizes. More can be added by adjusting the types
created in the generator.

## xmath/num
128-bit int and uint types.

## xmath/geom
Geometry primitives.

## xmath/rand
Randomizer based upon the crypto/rand package.
