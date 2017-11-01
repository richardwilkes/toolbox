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

## errs
Errors that contain stack traces with source locations, along with nested
causes, if any.

## desktop
Desktop integration utilities.

## formats/json
Easier manipulation of JSON data.

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

You can easily add rolling log files by using Nate Finch's lumberjack package:

```Go
package main

import (
    "github.com/richardwilkes/toolbox/log/jot"
    "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
    defer jot.Flush()
    jot.SetWriter(&lumberjack.Logger{
        Filename:   "test.log",
        MaxSize:    1, // megabytes
        MaxBackups: 4,
        MaxAge:     7, // days
        LocalTime:  true,
    })
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

## log/logadapter
This package defines an API to use for logging, which actual logging
implementations can implement directly or provide an adapter to use.

It also provides an implementation that just discards data given to it as
well as an implementation that wraps another logger and prefixes all output.

## rate
Rate limiting which supports a hierarchy of limiters, each capped by their
parent.

## txt
Text utilities.

## xio
io utilities.

## xio/fs
Filesystem utilities.

## xio/fs/safe
Safe, atomic saving of files.

## xio/network
Network-related utilities.

## xio/network/xhttp
HTTP-related utilities.

## xio/term
Terminal utilities.

## xmath
Math utilities.

## xmath/fixed
Simple fixed-point values that can be added, subtracted, multiplied and
divided. The values implement these interfaces for convenient
encoding/decoding:

- encoding.TextMarshaler
- encoding.TextUnmarshaler
- json.Marshaler
- json.Unmarshaler
- yaml.Marshaler
- yaml.Unmarshaler

## xmath/geom
Geometry primitives.
