// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package term

import (
	"io"
	"regexp"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
)

var (
	_ io.Writer       = &AnsiWriter{}
	_ io.StringWriter = &AnsiWriter{}
	_ io.ByteWriter   = &AnsiWriter{}
)

var colorSequenceMatcher = regexp.MustCompile(`\033\[(?:\d+(?:;\d+)*)*m`)

// AnsiWriter provides support for ANSI terminal escape sequences.
type AnsiWriter struct {
	w    io.Writer
	kind Kind
}

// NewAnsiWriter creates a new writer capable of emitting ANSI escape sequences for terminal control as well as color
// for those terminals that support it.
func NewAnsiWriter(w io.Writer) *AnsiWriter {
	return &AnsiWriter{w: w, kind: DetectKind(w)}
}

// Kind returns the term.Kind that this writer has been configured for.
func (a *AnsiWriter) Kind() Kind {
	return a.kind
}

// SetKind overrides the automatically detected terminal kind with the specified one.
func (a *AnsiWriter) SetKind(kind Kind) {
	if kind < Dumb || kind > Color24 {
		errs.Log(errs.Newf("invalid terminal kind %d", kind))
		return
	}
	a.kind = kind
}

// Bell causes the bell to sound.
func (a *AnsiWriter) Bell() {
	a.writeString(a.kind.Bell())
}

// Reset colors and styles.
func (a *AnsiWriter) Reset() {
	a.writeString(a.kind.Reset())
}

// Up moves the cursor up 'count' rows. If this would put it beyond the top edge of the screen, it will instead go to
// the top edge of the screen.
func (a *AnsiWriter) Up(count int) {
	a.writeString(a.kind.Up(count))
}

// Down moves the cursor down 'count' rows. If this would put it beyond the bottom edge of the screen, it will instead
// go to the bottom edge of the screen.
func (a *AnsiWriter) Down(count int) {
	a.writeString(a.kind.Down(count))
}

// Right moves the cursor right 'count' columns. If this would put it beyond the right edge of the screen, it will
// instead go to the right edge of the screen.
func (a *AnsiWriter) Right(count int) {
	a.writeString(a.kind.Right(count))
}

// Left moves the cursor left 'count' columns. If this would put it beyond the left edge of the screen, it will instead
// go to the left edge of the screen.
func (a *AnsiWriter) Left(count int) {
	a.writeString(a.kind.Left(count))
}

// Position the cursor at 'row' and 'column'. Both values are 1-based.
func (a *AnsiWriter) Position(row, column int) {
	a.writeString(a.kind.Position(row, column))
}

// Clear the screen and position the cursor at row 1, column 1.
func (a *AnsiWriter) Clear() {
	a.writeString(a.kind.Clear())
}

// ClearToStart clears the screen from the cursor to the beginning of the screen.
func (a *AnsiWriter) ClearToStart() {
	a.writeString(a.kind.ClearToStart())
}

// ClearToEnd clears the screen from the cursor to the end of the screen.
func (a *AnsiWriter) ClearToEnd() {
	a.writeString(a.kind.ClearToEnd())
}

// EraseLine clears the current row.
func (a *AnsiWriter) EraseLine() {
	a.writeString(a.kind.EraseLine())
}

// EraseLineToStart clears from the cursor position to the start of the current row.
func (a *AnsiWriter) EraseLineToStart() {
	a.writeString(a.kind.EraseLineToStart())
}

// EraseLineToEnd clears from the cursor position to the end of the current row.
func (a *AnsiWriter) EraseLineToEnd() {
	a.writeString(a.kind.EraseLineToEnd())
}

// SavePosition saves the current cursor position.
func (a *AnsiWriter) SavePosition() {
	a.writeString(a.kind.SavePosition())
}

// RestorePosition restores the previously saved cursor position.
func (a *AnsiWriter) RestorePosition() {
	a.writeString(a.kind.RestorePosition())
}

// HideCursor makes the cursor invisible.
func (a *AnsiWriter) HideCursor() {
	a.writeString(a.kind.HideCursor())
}

// ShowCursor makes the cursor visible.
func (a *AnsiWriter) ShowCursor() {
	a.writeString(a.kind.ShowCursor())
}

// Bold turns on bold text formatting.
func (a *AnsiWriter) Bold() {
	a.writeString(a.kind.Bold())
}

// NoBold turns off bold text formatting.
func (a *AnsiWriter) NoBold() {
	a.writeString(a.kind.NoBold())
}

// Dim turns on dim text formatting.
func (a *AnsiWriter) Dim() {
	a.writeString(a.kind.Dim())
}

// NoDim turns off dim text formatting.
func (a *AnsiWriter) NoDim() {
	a.writeString(a.kind.NoDim())
}

// Italic turns on italic text formatting.
func (a *AnsiWriter) Italic() {
	a.writeString(a.kind.Italic())
}

// NoItalic turns off italic text formatting.
func (a *AnsiWriter) NoItalic() {
	a.writeString(a.kind.NoItalic())
}

// Underline turns on underline text formatting.
func (a *AnsiWriter) Underline() {
	a.writeString(a.kind.Underline())
}

// NoUnderline turns off underline text formatting.
func (a *AnsiWriter) NoUnderline() {
	a.writeString(a.kind.NoUnderline())
}

// Blink turns on blink text formatting.
func (a *AnsiWriter) Blink() {
	a.writeString(a.kind.Blink())
}

// NoBlink turns off blink text formatting.
func (a *AnsiWriter) NoBlink() {
	a.writeString(a.kind.NoBlink())
}

// Inverse turns on inverse text formatting (foreground and background colors are swapped).
func (a *AnsiWriter) Inverse() {
	a.writeString(a.kind.Inverse())
}

// NoInverse turns off inverse text formatting.
func (a *AnsiWriter) NoInverse() {
	a.writeString(a.kind.NoInverse())
}

// StrikeThrough turns on strikethrough text formatting.
func (a *AnsiWriter) StrikeThrough() {
	a.writeString(a.kind.StrikeThrough())
}

// NoStrikeThrough turns off strikethrough text formatting.
func (a *AnsiWriter) NoStrikeThrough() {
	a.writeString(a.kind.NoStrikeThrough())
}

// Overline turns on overline text formatting.
func (a *AnsiWriter) Overline() {
	a.writeString(a.kind.Overline())
}

// NoOverline turns off overline text formatting.
func (a *AnsiWriter) NoOverline() {
	a.writeString(a.kind.NoOverline())
}

// Black foreground.
func (a *AnsiWriter) Black() {
	a.writeString(a.kind.Black())
}

// Red foreground.
func (a *AnsiWriter) Red() {
	a.writeString(a.kind.Red())
}

// Green foreground.
func (a *AnsiWriter) Green() {
	a.writeString(a.kind.Green())
}

// Yellow foreground.
func (a *AnsiWriter) Yellow() {
	a.writeString(a.kind.Yellow())
}

// Blue foreground.
func (a *AnsiWriter) Blue() {
	a.writeString(a.kind.Blue())
}

// Magenta foreground.
func (a *AnsiWriter) Magenta() {
	a.writeString(a.kind.Magenta())
}

// Cyan foreground.
func (a *AnsiWriter) Cyan() {
	a.writeString(a.kind.Cyan())
}

// White foreground.
func (a *AnsiWriter) White() {
	a.writeString(a.kind.White())
}

// Grey foreground. Same as BrightBlack.
func (a *AnsiWriter) Grey() {
	a.writeString(a.kind.Grey())
}

// Gray foreground. Same as BrightBlack.
func (a *AnsiWriter) Gray() {
	a.writeString(a.kind.Gray())
}

// BrightBlack foreground.
func (a *AnsiWriter) BrightBlack() {
	a.writeString(a.kind.BrightBlack())
}

// BrightRed foreground.
func (a *AnsiWriter) BrightRed() {
	a.writeString(a.kind.BrightRed())
}

// BrightGreen foreground.
func (a *AnsiWriter) BrightGreen() {
	a.writeString(a.kind.BrightGreen())
}

// BrightYellow foreground.
func (a *AnsiWriter) BrightYellow() {
	a.writeString(a.kind.BrightYellow())
}

// BrightBlue foreground.
func (a *AnsiWriter) BrightBlue() {
	a.writeString(a.kind.BrightBlue())
}

// BrightMagenta foreground.
func (a *AnsiWriter) BrightMagenta() {
	a.writeString(a.kind.BrightMagenta())
}

// BrightCyan foreground.
func (a *AnsiWriter) BrightCyan() {
	a.writeString(a.kind.BrightCyan())
}

// BrightWhite foreground.
func (a *AnsiWriter) BrightWhite() {
	a.writeString(a.kind.BrightWhite())
}

// FgReset resets the foreground color to the default.
func (a *AnsiWriter) FgReset() {
	a.writeString(a.kind.FgReset())
}

// BgBlack background.
func (a *AnsiWriter) BgBlack() {
	a.writeString(a.kind.BgBlack())
}

// BgRed background.
func (a *AnsiWriter) BgRed() {
	a.writeString(a.kind.BgRed())
}

// BgGreen background.
func (a *AnsiWriter) BgGreen() {
	a.writeString(a.kind.BgGreen())
}

// BgYellow background.
func (a *AnsiWriter) BgYellow() {
	a.writeString(a.kind.BgYellow())
}

// BgBlue background.
func (a *AnsiWriter) BgBlue() {
	a.writeString(a.kind.BgBlue())
}

// BgMagenta background.
func (a *AnsiWriter) BgMagenta() {
	a.writeString(a.kind.BgMagenta())
}

// BgCyan background.
func (a *AnsiWriter) BgCyan() {
	a.writeString(a.kind.BgCyan())
}

// BgWhite background.
func (a *AnsiWriter) BgWhite() {
	a.writeString(a.kind.BgWhite())
}

// BgGrey background. Same as BgBrightBlack.
func (a *AnsiWriter) BgGrey() {
	a.writeString(a.kind.BgGrey())
}

// BgGray background. Same as BgBrightBlack.
func (a *AnsiWriter) BgGray() {
	a.writeString(a.kind.BgGray())
}

// BgBrightBlack background.
func (a *AnsiWriter) BgBrightBlack() {
	a.writeString(a.kind.BgBrightBlack())
}

// BgBrightRed background.
func (a *AnsiWriter) BgBrightRed() {
	a.writeString(a.kind.BgBrightRed())
}

// BgBrightGreen background.
func (a *AnsiWriter) BgBrightGreen() {
	a.writeString(a.kind.BgBrightGreen())
}

// BgBrightYellow background.
func (a *AnsiWriter) BgBrightYellow() {
	a.writeString(a.kind.BgBrightYellow())
}

// BgBrightBlue background.
func (a *AnsiWriter) BgBrightBlue() {
	a.writeString(a.kind.BgBrightBlue())
}

// BgBrightMagenta background.
func (a *AnsiWriter) BgBrightMagenta() {
	a.writeString(a.kind.BgBrightMagenta())
}

// BgBrightCyan background.
func (a *AnsiWriter) BgBrightCyan() {
	a.writeString(a.kind.BgBrightCyan())
}

// BgBrightWhite background.
func (a *AnsiWriter) BgBrightWhite() {
	a.writeString(a.kind.BgBrightWhite())
}

// BgReset resets the background color to the default.
func (a *AnsiWriter) BgReset() {
	a.writeString(a.kind.BgReset())
}

// Color256 sets the foreground color using an 8-bit ANSI color code.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit.
func (a *AnsiWriter) Color256(color uint8) {
	a.writeString(a.kind.Color256(color))
}

// BgColor256 sets the background color using an 8-bit ANSI color code.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit.
func (a *AnsiWriter) BgColor256(color uint8) {
	a.writeString(a.kind.BgColor256(color))
}

// RGB sets the foreground color using RGB values.
func (a *AnsiWriter) RGB(r, g, b uint8) {
	a.writeString(a.kind.RGB(r, g, b))
}

// BgRGB sets the background color using RGB values.
func (a *AnsiWriter) BgRGB(r, g, b uint8) {
	a.writeString(a.kind.BgRGB(r, g, b))
}

// Write implements the io.Writer interface.
func (a *AnsiWriter) Write(p []byte) (n int, err error) {
	return a.w.Write(p)
}

// WriteString implements the io.StringWriter interface.
func (a *AnsiWriter) WriteString(s string) (n int, err error) {
	return io.WriteString(a.w, s)
}

func (a *AnsiWriter) writeString(s string) {
	_, _ = io.WriteString(a.w, s) //nolint:errcheck // We don't care about the error here.
}

// WriteByte implements the io.ByteWriter interface.
func (a *AnsiWriter) WriteByte(c byte) error {
	_, err := a.w.Write([]byte{c})
	return err
}

func (a *AnsiWriter) writeByte(c byte) {
	_, _ = a.w.Write([]byte{c}) //nolint:errcheck // We don't care about the error here.
}

// WrapText prints the 'prefix' to 'out' and then wraps 'text' in the remaining space.
func (a *AnsiWriter) WrapText(prefix, text string) {
	a.writeString(prefix)
	avail, _ := Size(a.w)
	prefixLength := len(colorSequenceMatcher.ReplaceAllString(prefix, ""))
	avail -= 1 + prefixLength
	if avail < 1 {
		avail = 1
	}
	remaining := avail
	indent := strings.Repeat(" ", prefixLength)
	for line := range strings.SplitSeq(text, "\n") {
		for _, ch := range line {
			if ch == ' ' {
				a.writeByte(' ')
				remaining--
			} else {
				break
			}
		}
		for i, token := range strings.Fields(line) {
			length := len(colorSequenceMatcher.ReplaceAllString(token, "")) + 1
			if i != 0 {
				if length > remaining {
					a.writeByte('\n')
					a.writeString(indent)
					remaining = avail
				} else {
					a.writeByte(' ')
				}
			}
			a.writeString(token)
			remaining -= length
		}
		a.writeByte('\n')
	}
}
