// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package term provides terminal utilities.
package term

import (
	"io"
	"os"

	"github.com/richardwilkes/toolbox/v2/errs"
	"golang.org/x/term"
)

// IsTerminal returns true if the writer is a terminal.
func IsTerminal(w io.Writer) bool {
	switch t := w.(type) {
	case *os.File:
		return term.IsTerminal(int(t.Fd()))
	case *AnsiWriter:
		return IsTerminal(t.w)
	default:
		return false
	}
}

// Size returns the number of columns and rows comprising the terminal.
func Size(w io.Writer) (columns, rows int) {
	switch t := w.(type) {
	case *os.File:
		var err error
		if columns, rows, err = term.GetSize(int(t.Fd())); err == nil {
			return columns, rows
		}
	case *AnsiWriter:
		return Size(t.w)
	}
	return 80, 24
}

// DetectKind returns the kind of support available in the terminal.
func DetectKind(w io.Writer) Kind {
	if IsTerminal(w) {
		envTerm := os.Getenv("TERM")
		if envTerm == "dumb" {
			return Dumb
		}
		if !enableColor() {
			return Dumb
		}
		return colorSupport(envTerm)
	}
	return Dumb
}

// RawRead reads a byte from the terminal without requiring the enter/return key to be pressed.
func RawRead(r io.Reader) (ch byte, ok bool) {
	var f *os.File
	if f, ok = r.(*os.File); ok {
		fd := int(f.Fd())
		if term.IsTerminal(fd) {
			if state, err := term.MakeRaw(fd); err == nil {
				defer func() {
					if err = term.Restore(fd, state); err != nil {
						errs.Log(errs.NewWithCause("unable to restore terminal state", err))
					}
				}()
			}
		}
	}
	data := make([]byte, 1)
	if numRead, err := r.Read(data); err != nil || numRead == 0 {
		return 0, false
	}
	return data[0], true
}
