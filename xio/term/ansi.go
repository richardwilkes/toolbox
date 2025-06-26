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
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"golang.org/x/term"
)

// Possible terminal kinds.
const (
	Dumb    Kind = iota // No color support
	Color4              // 4-bit (16 colors)
	Color8              // 8-bit (256 colors)
	Color24             // 24-bit (16,777,216 colors)
)

const (
	start = "\033["
	end   = "m"
)

var colorSequenceMatcher = regexp.MustCompile(`\033\[(?:\d+(?:;\d+)*)*m`)

// Kind represents the kind of terminal in use.
type Kind int

// ANSI provides support for ANSI terminal escape sequences.
type ANSI struct {
	out  io.Writer
	kind Kind
}

// IsTerminal returns true if the writer is a terminal.
func IsTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

// ColorSupport returns the kind of color support available in the terminal.
func ColorSupport(w io.Writer) Kind {
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

// Size returns the number of columns and rows comprising the terminal.
func Size(w io.Writer) (columns, rows int) {
	if f, ok := w.(*os.File); ok {
		var err error
		if columns, rows, err = term.GetSize(int(f.Fd())); err == nil {
			return columns, rows
		}
	}
	return 80, 24
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

// NewANSI creates a new ANSI terminal and attaches it to 'out'.
func NewANSI(out io.Writer) *ANSI {
	return &ANSI{out: out, kind: ColorSupport(out)}
}

// SetKind overrides the automatically detected terminal kind with the specified one.
func (a *ANSI) SetKind(kind Kind) {
	if kind < Dumb || kind > Color24 {
		errs.Log(errs.Newf("invalid terminal kind %d", kind))
		return
	}
	a.kind = kind
}

// Bell causes the bell to sound.
func (a *ANSI) Bell() {
	fmt.Fprint(a, "\007")
}

// Reset colors and styles.
func (a *ANSI) Reset() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[m")
	}
}

// Up moves the cursor up 'count' rows. If this would put it beyond the top edge of the screen, it will instead go to
// the top edge of the screen.
func (a *ANSI) Up(count int) {
	if a.kind != Dumb {
		fmt.Fprint(a, start+strconv.Itoa(count)+"A")
	}
}

// Down moves the cursor down 'count' rows. If this would put it beyond the bottom edge of the screen, it will instead
// go to the bottom edge of the screen.
func (a *ANSI) Down(count int) {
	if a.kind != Dumb {
		fmt.Fprint(a, start+strconv.Itoa(count)+"B")
	}
}

// Right moves the cursor right 'count' columns. If this would put it beyond the right edge of the screen, it will
// instead go to the right edge of the screen.
func (a *ANSI) Right(count int) {
	if a.kind != Dumb {
		fmt.Fprint(a, start+strconv.Itoa(count)+"C")
	}
}

// Left moves the cursor left 'count' columns. If this would put it beyond the left edge of the screen, it will instead
// go to the left edge of the screen.
func (a *ANSI) Left(count int) {
	if a.kind != Dumb {
		fmt.Fprint(a, start+strconv.Itoa(count)+"D")
	}
}

// Position the cursor at 'row' and 'column'. Both values are 1-based.
func (a *ANSI) Position(row, column int) {
	if a.kind != Dumb {
		fmt.Fprint(a, start+strconv.Itoa(row)+";"+strconv.Itoa(column)+"H")
	}
}

// Clear the screen and position the cursor at row 1, column 1.
func (a *ANSI) Clear() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[2J")
	}
}

// ClearToStart clears the screen from the cursor to the beginning of the screen.
func (a *ANSI) ClearToStart() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[1J")
	}
}

// ClearToEnd clears the screen from the cursor to the end of the screen.
func (a *ANSI) ClearToEnd() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[J")
	}
}

// EraseLine clears the current row.
func (a *ANSI) EraseLine() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[2K")
	}
}

// EraseLineToStart clears from the cursor position to the start of the current row.
func (a *ANSI) EraseLineToStart() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[1K")
	}
}

// EraseLineToEnd clears from the cursor position to the end of the current row.
func (a *ANSI) EraseLineToEnd() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[K")
	}
}

// SavePosition saves the current cursor position.
func (a *ANSI) SavePosition() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[s")
	}
}

// RestorePosition restores the previously saved cursor position.
func (a *ANSI) RestorePosition() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[u")
	}
}

// HideCursor makes the cursor invisible.
func (a *ANSI) HideCursor() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[?25l")
	}
}

// ShowCursor makes the cursor visible.
func (a *ANSI) ShowCursor() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[?25h")
	}
}

// Bold turns on bold text formatting.
func (a *ANSI) Bold() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[1m")
	}
}

// NoBold turns off bold text formatting.
func (a *ANSI) NoBold() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[22m")
	}
}

// Dim turns on dim text formatting.
func (a *ANSI) Dim() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[2m")
	}
}

// NoDim turns off dim text formatting.
func (a *ANSI) NoDim() {
	a.NoBold()
}

// Italic turns on italic text formatting.
func (a *ANSI) Italic() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[3m")
	}
}

// NoItalic turns off italic text formatting.
func (a *ANSI) NoItalic() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[23m")
	}
}

// Underline turns on underline text formatting.
func (a *ANSI) Underline() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[4m")
	}
}

// NoUnderline turns off underline text formatting.
func (a *ANSI) NoUnderline() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[24m")
	}
}

// Blink turns on blink text formatting.
func (a *ANSI) Blink() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[5m")
	}
}

// NoBlink turns off blink text formatting.
func (a *ANSI) NoBlink() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[25m")
	}
}

// Inverse turns on inverse text formatting (foreground and background colors are swapped).
func (a *ANSI) Inverse() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[7m")
	}
}

// NoInverse turns off inverse text formatting.
func (a *ANSI) NoInverse() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[27m")
	}
}

// StrikeThrough turns on strikethrough text formatting.
func (a *ANSI) StrikeThrough() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[9m")
	}
}

// NoStrikeThrough turns off strikethrough text formatting.
func (a *ANSI) NoStrikeThrough() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[29m")
	}
}

// Overline turns on overline text formatting.
func (a *ANSI) Overline() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[53m")
	}
}

// NoOverline turns off overline text formatting.
func (a *ANSI) NoOverline() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[55m")
	}
}

// Black foreground.
func (a *ANSI) Black() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[30m")
	}
}

// Red foreground.
func (a *ANSI) Red() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[31m")
	}
}

// Green foreground.
func (a *ANSI) Green() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[32m")
	}
}

// Yellow foreground.
func (a *ANSI) Yellow() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[33m")
	}
}

// Blue foreground.
func (a *ANSI) Blue() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[34m")
	}
}

// Magenta foreground.
func (a *ANSI) Magenta() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[35m")
	}
}

// Cyan foreground.
func (a *ANSI) Cyan() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[36m")
	}
}

// White foreground.
func (a *ANSI) White() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[37m")
	}
}

// Grey foreground. Same as BrightBlack.
func (a *ANSI) Grey() {
	a.BrightBlack()
}

// Gray foreground. Same as BrightBlack.
func (a *ANSI) Gray() {
	a.BrightBlack()
}

// BrightBlack foreground.
func (a *ANSI) BrightBlack() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[90m")
	}
}

// BrightRed foreground.
func (a *ANSI) BrightRed() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[91m")
	}
}

// BrightGreen foreground.
func (a *ANSI) BrightGreen() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[92m")
	}
}

// BrightYellow foreground.
func (a *ANSI) BrightYellow() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[93m")
	}
}

// BrightBlue foreground.
func (a *ANSI) BrightBlue() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[94m")
	}
}

// BrightMagenta foreground.
func (a *ANSI) BrightMagenta() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[95m")
	}
}

// BrightCyan foreground.
func (a *ANSI) BrightCyan() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[96m")
	}
}

// BrightWhite foreground.
func (a *ANSI) BrightWhite() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[97m")
	}
}

// FgReset resets the foreground color to the default.
func (a *ANSI) FgReset() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[39m")
	}
}

// BgBlack background.
func (a *ANSI) BgBlack() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[40m")
	}
}

// BgRed background.
func (a *ANSI) BgRed() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[41m")
	}
}

// BgGreen background.
func (a *ANSI) BgGreen() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[42m")
	}
}

// BgYellow background.
func (a *ANSI) BgYellow() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[43m")
	}
}

// BgBlue background.
func (a *ANSI) BgBlue() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[44m")
	}
}

// BgMagenta background.
func (a *ANSI) BgMagenta() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[45m")
	}
}

// BgCyan background.
func (a *ANSI) BgCyan() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[46m")
	}
}

// BgWhite background.
func (a *ANSI) BgWhite() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[47m")
	}
}

// BgGrey background. Same as BgBrightBlack.
func (a *ANSI) BgGrey() {
	a.BgBrightBlack()
}

// BgGray background. Same as BgBrightBlack.
func (a *ANSI) BgGray() {
	a.BgBrightBlack()
}

// BgBrightBlack background.
func (a *ANSI) BgBrightBlack() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[100m")
	}
}

// BgBrightRed background.
func (a *ANSI) BgBrightRed() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[101m")
	}
}

// BgBrightGreen background.
func (a *ANSI) BgBrightGreen() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[102m")
	}
}

// BgBrightYellow background.
func (a *ANSI) BgBrightYellow() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[103m")
	}
}

// BgBrightBlue background.
func (a *ANSI) BgBrightBlue() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[104m")
	}
}

// BgBrightMagenta background.
func (a *ANSI) BgBrightMagenta() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[105m")
	}
}

// BgBrightCyan background.
func (a *ANSI) BgBrightCyan() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[106m")
	}
}

// BgBrightWhite background.
func (a *ANSI) BgBrightWhite() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[107m")
	}
}

// BgReset resets the background color to the default.
func (a *ANSI) BgReset() {
	if a.kind != Dumb {
		fmt.Fprint(a, "\033[49m")
	}
}

// Ansi256 sets the foreground color using an 8-bit ANSI color code.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit.
func (a *ANSI) Ansi256(color uint8) {
	switch a.kind {
	case Color4:
		fmt.Fprint(a, start+strconv.Itoa(int(ansi256ToAnsi16[color]))+end)
	case Color8, Color24:
		fmt.Fprint(a, "\033[38;5;"+strconv.Itoa(int(color))+end)
	}
}

// BgAnsi256 sets the background color using an 8-bit ANSI color code.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit.
func (a *ANSI) BgAnsi256(color uint8) {
	switch a.kind {
	case Color4:
		fmt.Fprint(a, start+strconv.Itoa(int(ansi256ToAnsi16[color]+10))+end)
	case Color8, Color24:
		fmt.Fprint(a, "\033[48;5;"+strconv.Itoa(int(color))+end)
	}
}

// RGB sets the foreground color using RGB values.
func (a *ANSI) RGB(r, g, b uint8) {
	switch a.kind {
	case Color4:
		fmt.Fprint(a, start+strconv.Itoa(int(ansi256ToAnsi16[rgbToAnsi256(r, g, b)]))+end)
	case Color8:
		a.Ansi256(rgbToAnsi256(r, g, b))
	case Color24:
		fmt.Fprint(a, "\033[38;2;"+strconv.Itoa(int(r))+";"+strconv.Itoa(int(g))+";"+strconv.Itoa(int(b))+end)
	}
}

// BgRGB sets the background color using RGB values.
func (a *ANSI) BgRGB(r, g, b uint8) {
	switch a.kind {
	case Color4:
		fmt.Fprint(a, start+strconv.Itoa(int(ansi256ToAnsi16[rgbToAnsi256(r, g, b)]+10))+end)
	case Color8:
		a.BgAnsi256(rgbToAnsi256(r, g, b))
	case Color24:
		fmt.Fprint(a, "\033[48;2;"+strconv.Itoa(int(r))+";"+strconv.Itoa(int(g))+";"+strconv.Itoa(int(b))+end)
	}
}

func rgbToAnsi256(red, green, blue uint8) uint8 {
	if red == green && green == blue {
		if red < 8 {
			return 16
		}
		if red > 248 {
			return 231
		}
		return uint8(math.Round(((float64(red)-8)/247)*24)) + 232
	}
	return 16 + uint8((36*rgbTo256s(red))+(6*rgbTo256s(green))+rgbTo256s(blue))
}

func rgbTo256s(value uint8) float64 {
	return math.Round(float64(value) / 255 * 5)
}

var ansi256ToAnsi16 = []uint8{
	// Standard colors
	30, 31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96, 97,
	// Colors
	30, 30, 30, 34, 34, 34, 30, 30, 34, 34, 34, 34, 32, 32, 90, 34, 34, 34,
	32, 32, 36, 36, 36, 36, 32, 32, 36, 36, 36, 36, 32, 32, 92, 36, 36, 36,
	30, 30, 30, 34, 34, 34, 30, 30, 90, 34, 34, 34, 32, 90, 90, 90, 94, 94,
	32, 32, 90, 36, 36, 94, 32, 32, 92, 36, 36, 96, 32, 92, 92, 92, 96, 96,
	30, 30, 90, 90, 34, 94, 31, 90, 90, 90, 94, 94, 90, 90, 90, 90, 94, 94,
	33, 90, 90, 90, 94, 94, 33, 92, 92, 92, 96, 96, 92, 92, 92, 92, 96, 96,
	31, 31, 90, 35, 35, 35, 31, 31, 90, 35, 35, 35, 31, 90, 90, 90, 94, 94,
	33, 33, 90, 37, 37, 94, 33, 33, 92, 37, 37, 37, 33, 92, 92, 92, 37, 96,
	31, 31, 31, 35, 35, 35, 31, 31, 91, 35, 35, 35, 31, 91, 91, 35, 35, 95,
	33, 33, 91, 37, 37, 95, 33, 33, 93, 37, 37, 37, 33, 93, 93, 93, 37, 97,
	31, 31, 91, 35, 35, 35, 31, 91, 91, 35, 35, 95, 31, 91, 91, 91, 95, 95,
	33, 91, 91, 91, 95, 95, 33, 93, 93, 93, 37, 97, 33, 93, 93, 93, 97, 97,
	// Greyscale
	30, 30, 30, 30, 30, 90, 90, 90, 90, 90, 90, 90,
	90, 90, 90, 37, 37, 37, 37, 37, 37, 37, 97, 97,
}

// Write implements the io.Writer interface.
func (a *ANSI) Write(p []byte) (n int, err error) {
	return a.out.Write(p)
}

// WrapText prints the 'prefix' to 'out' and then wraps 'text' in the remaining space.
func (a *ANSI) WrapText(prefix, text string) {
	fmt.Fprint(a, prefix)
	avail, _ := Size(a.out)
	avail -= 1 + len(prefix)
	if avail < 1 {
		avail = 1
	}
	remaining := avail
	indent := strings.Repeat(" ", len(prefix))
	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' {
				fmt.Fprint(a, " ")
				remaining--
			} else {
				break
			}
		}
		for i, token := range strings.Fields(line) {
			length := len(colorSequenceMatcher.ReplaceAllString(token, "")) + 1
			if i != 0 {
				if length > remaining {
					fmt.Fprintln(a)
					fmt.Fprint(a, indent)
					remaining = avail
				} else {
					fmt.Fprint(a, " ")
				}
			}
			fmt.Fprint(a, token)
			remaining -= length
		}
		fmt.Fprintln(a)
	}
}
