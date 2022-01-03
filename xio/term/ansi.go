// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package term

import (
	"fmt"
	"io"
)

// ANSI color constants
const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// ANSI style constants. Multiple styles may be or'd together.
const (
	Bold Style = 1 << iota
	Underline
	Blink
	Normal = 0
)

// Color represents an ANSI terminal color
type Color int

// Style represents an ANSI terminal style. Multiple styles may be or'd together.
type Style int

// ANSI provides support for ANSI terminal escape sequences.
type ANSI struct {
	out io.Writer
	ok  bool
}

// NewANSI creates a new ANSI terminal and attaches it to 'out'.
func NewANSI(out io.Writer) *ANSI {
	return &ANSI{out: out, ok: IsTerminal(out)}
}

// Bell causes the bell to sound.
func (a *ANSI) Bell() {
	fmt.Fprint(a, "\007")
}

// Reset colors and styles.
func (a *ANSI) Reset() {
	if a.ok {
		fmt.Fprint(a, "\033[m")
	}
}

// Up moves the cursor up 'count' rows. If this would put it beyond the top edge of the screen, it will instead go to
// the top edge of the screen.
func (a *ANSI) Up(count int) {
	if a.ok {
		fmt.Fprintf(a, "\033[%dA", count)
	}
}

// Down moves the cursor down 'count' rows. If this would put it beyond the bottom edge of the screen, it will instead
// go to the bottom edge of the screen.
func (a *ANSI) Down(count int) {
	if a.ok {
		fmt.Fprintf(a, "\033[%dB", count)
	}
}

// Left moves the cursor left 'count' columns. If this would put it beyond the left edge of the screen, it will instead
// go to the left edge of the screen.
func (a *ANSI) Left(count int) {
	if a.ok {
		fmt.Fprintf(a, "\033[%dD", count)
	}
}

// Right moves the cursor right 'count' columns. If this would put it beyond the right edge of the screen, it will
// instead go to the right edge of the screen.
func (a *ANSI) Right(count int) {
	if a.ok {
		fmt.Fprintf(a, "\033[%dC", count)
	}
}

// Position the cursor at 'row' and 'column'. Both values are 1-based.
func (a *ANSI) Position(row, column int) {
	if a.ok {
		fmt.Fprintf(a, "\033[%d;%dH", row, column)
	}
}

// Clear the screen and position the cursor at row 1, column 1.
func (a *ANSI) Clear() {
	if a.ok {
		fmt.Fprint(a, "\033[2J")
	}
}

// ClearToStart clears the screen from the cursor to the beginning of the screen.
func (a *ANSI) ClearToStart() {
	if a.ok {
		fmt.Fprint(a, "\033[1J")
	}
}

// ClearToEnd clears the screen from the cursor to the end of the screen.
func (a *ANSI) ClearToEnd() {
	if a.ok {
		fmt.Fprint(a, "\033[J")
	}
}

// EraseLine clears the current row.
func (a *ANSI) EraseLine() {
	if a.ok {
		fmt.Fprint(a, "\033[2K")
	}
}

// EraseLineToStart clears from the cursor position to the start of the current row.
func (a *ANSI) EraseLineToStart() {
	if a.ok {
		fmt.Fprint(a, "\033[1K")
	}
}

// EraseLineToEnd clears from the cursor position to the end of the current row.
func (a *ANSI) EraseLineToEnd() {
	if a.ok {
		fmt.Fprint(a, "\033[K")
	}
}

// SavePosition saves the current cursor position.
func (a *ANSI) SavePosition() {
	if a.ok {
		fmt.Fprint(a, "\033[s")
	}
}

// RestorePosition restores the previously saved cursor position.
func (a *ANSI) RestorePosition() {
	if a.ok {
		fmt.Fprint(a, "\033[u")
	}
}

// HideCursor makes the cursor invisible.
func (a *ANSI) HideCursor() {
	if a.ok {
		fmt.Fprint(a, "\033[?25l")
	}
}

// ShowCursor makes the cursor visible.
func (a *ANSI) ShowCursor() {
	if a.ok {
		fmt.Fprint(a, "\033[?25h")
	}
}

// Foreground sets the foreground color and style for subsequent output.
func (a *ANSI) Foreground(color Color, style Style) {
	if a.ok {
		fmt.Fprint(a, "\033[0;")
		if style&Bold == Bold {
			fmt.Fprint(a, "1;")
		}
		if style&Underline == Underline {
			fmt.Fprint(a, "4;")
		}
		if style&Blink == Blink {
			fmt.Fprint(a, "5;")
		}
		fmt.Fprintf(a, "%dm", 30+color)
	}
}

// Background sets the background color for subsequent output.
func (a *ANSI) Background(color Color) {
	if a.ok {
		fmt.Fprintf(a, "\033[%dm", 40+color)
	}
}

// Write implements the io.Writer interface.
func (a *ANSI) Write(p []byte) (n int, err error) {
	return a.out.Write(p)
}
