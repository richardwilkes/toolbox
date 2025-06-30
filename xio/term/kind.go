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
	"math"
	"strconv"
)

// Kind represents the kind of terminal in use.
type Kind int

// Possible terminal kinds.
const (
	InvalidKind Kind = iota // Invalid terminal kind
	Dumb                    // No color support
	Mono                    // Monochrome, but supports some ANSI codes
	Color4                  // 4-bit (16 colors)
	Color8                  // 8-bit (256 colors)
	Color24                 // 24-bit (16,777,216 colors)
)

const (
	start = "\033["
	end   = "m"
)

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

// Bell causes the bell to sound.
func (k Kind) Bell() string {
	return "\007"
}

// Reset colors and styles.
func (k Kind) Reset() string {
	if k == Dumb {
		return ""
	}
	return "\033[m"
}

// Up moves the cursor up 'count' rows. If this would put it beyond the top edge of the screen, it will instead go to
// the top edge of the screen.
func (k Kind) Up(count int) string {
	if k == Dumb {
		return ""
	}
	return start + strconv.Itoa(count) + "A"
}

// Down moves the cursor down 'count' rows. If this would put it beyond the bottom edge of the screen, it will instead
// go to the bottom edge of the screen.
func (k Kind) Down(count int) string {
	if k == Dumb {
		return ""
	}
	return start + strconv.Itoa(count) + "B"
}

// Right moves the cursor right 'count' columns. If this would put it beyond the right edge of the screen, it will
// instead go to the right edge of the screen.
func (k Kind) Right(count int) string {
	if k == Dumb {
		return ""
	}
	return start + strconv.Itoa(count) + "C"
}

// Left moves the cursor left 'count' columns. If this would put it beyond the left edge of the screen, it will instead
// go to the left edge of the screen.
func (k Kind) Left(count int) string {
	if k == Dumb {
		return ""
	}
	return start + strconv.Itoa(count) + "D"
}

// Position the cursor at 'row' and 'column'. Both values are 1-based.
func (k Kind) Position(row, column int) string {
	if k == Dumb {
		return ""
	}
	return start + strconv.Itoa(row) + ";" + strconv.Itoa(column) + "H"
}

// Clear the screen and position the cursor at row 1, column 1.
func (k Kind) Clear() string {
	if k == Dumb {
		return ""
	}
	return "\033[2J"
}

// ClearToStart clears the screen from the cursor to the beginning of the screen.
func (k Kind) ClearToStart() string {
	if k == Dumb {
		return ""
	}
	return "\033[1J"
}

// ClearToEnd clears the screen from the cursor to the end of the screen.
func (k Kind) ClearToEnd() string {
	if k == Dumb {
		return ""
	}
	return "\033[J"
}

// EraseLine clears the current row.
func (k Kind) EraseLine() string {
	if k == Dumb {
		return ""
	}
	return "\033[2K"
}

// EraseLineToStart clears from the cursor position to the start of the current row.
func (k Kind) EraseLineToStart() string {
	if k == Dumb {
		return ""
	}
	return "\033[1K"
}

// EraseLineToEnd clears from the cursor position to the end of the current row.
func (k Kind) EraseLineToEnd() string {
	if k == Dumb {
		return ""
	}
	return "\033[K"
}

// SavePosition saves the current cursor position.
func (k Kind) SavePosition() string {
	if k == Dumb {
		return ""
	}
	return "\033[s"
}

// RestorePosition restores the previously saved cursor position.
func (k Kind) RestorePosition() string {
	if k == Dumb {
		return ""
	}
	return "\033[u"
}

// HideCursor makes the cursor invisible.
func (k Kind) HideCursor() string {
	if k == Dumb {
		return ""
	}
	return "\033[?25l"
}

// ShowCursor makes the cursor visible.
func (k Kind) ShowCursor() string {
	if k == Dumb {
		return ""
	}
	return "\033[?25h"
}

// Bold turns on bold text formatting.
func (k Kind) Bold() string {
	if k == Dumb {
		return ""
	}
	return "\033[1m"
}

// NoBold turns off bold text formatting.
func (k Kind) NoBold() string {
	if k == Dumb {
		return ""
	}
	return "\033[22m"
}

// Dim turns on dim text formatting.
func (k Kind) Dim() string {
	if k == Dumb {
		return ""
	}
	return "\033[2m"
}

// NoDim turns off dim text formatting.
func (k Kind) NoDim() string {
	return k.NoBold()
}

// Italic turns on italic text formatting.
func (k Kind) Italic() string {
	if k == Dumb {
		return ""
	}
	return "\033[3m"
}

// NoItalic turns off italic text formatting.
func (k Kind) NoItalic() string {
	if k == Dumb {
		return ""
	}
	return "\033[23m"
}

// Underline turns on underline text formatting.
func (k Kind) Underline() string {
	if k == Dumb {
		return ""
	}
	return "\033[4m"
}

// NoUnderline turns off underline text formatting.
func (k Kind) NoUnderline() string {
	if k == Dumb {
		return ""
	}
	return "\033[24m"
}

// Blink turns on blink text formatting.
func (k Kind) Blink() string {
	if k == Dumb {
		return ""
	}
	return "\033[5m"
}

// NoBlink turns off blink text formatting.
func (k Kind) NoBlink() string {
	if k == Dumb {
		return ""
	}
	return "\033[25m"
}

// Inverse turns on inverse text formatting (foreground and background colors are swapped).
func (k Kind) Inverse() string {
	if k == Dumb {
		return ""
	}
	return "\033[7m"
}

// NoInverse turns off inverse text formatting.
func (k Kind) NoInverse() string {
	if k == Dumb {
		return ""
	}
	return "\033[27m"
}

// StrikeThrough turns on strikethrough text formatting.
func (k Kind) StrikeThrough() string {
	if k == Dumb {
		return ""
	}
	return "\033[9m"
}

// NoStrikeThrough turns off strikethrough text formatting.
func (k Kind) NoStrikeThrough() string {
	if k == Dumb {
		return ""
	}
	return "\033[29m"
}

// Overline turns on overline text formatting.
func (k Kind) Overline() string {
	if k == Dumb {
		return ""
	}
	return "\033[53m"
}

// NoOverline turns off overline text formatting.
func (k Kind) NoOverline() string {
	if k == Dumb {
		return ""
	}
	return "\033[55m"
}

// Black foreground.
func (k Kind) Black() string {
	if k < Color4 {
		return ""
	}
	return "\033[30m"
}

// Red foreground.
func (k Kind) Red() string {
	if k < Color4 {
		return ""
	}
	return "\033[31m"
}

// Green foreground.
func (k Kind) Green() string {
	if k < Color4 {
		return ""
	}
	return "\033[32m"
}

// Yellow foreground.
func (k Kind) Yellow() string {
	if k < Color4 {
		return ""
	}
	return "\033[33m"
}

// Blue foreground.
func (k Kind) Blue() string {
	if k < Color4 {
		return ""
	}
	return "\033[34m"
}

// Magenta foreground.
func (k Kind) Magenta() string {
	if k < Color4 {
		return ""
	}
	return "\033[35m"
}

// Cyan foreground.
func (k Kind) Cyan() string {
	if k < Color4 {
		return ""
	}
	return "\033[36m"
}

// White foreground.
func (k Kind) White() string {
	if k < Color4 {
		return ""
	}
	return "\033[37m"
}

// Grey foreground. Same as BrightBlack.
func (k Kind) Grey() string {
	return k.BrightBlack()
}

// Gray foreground. Same as BrightBlack.
func (k Kind) Gray() string {
	return k.BrightBlack()
}

// BrightBlack foreground.
func (k Kind) BrightBlack() string {
	if k < Color4 {
		return ""
	}
	return "\033[90m"
}

// BrightRed foreground.
func (k Kind) BrightRed() string {
	if k < Color4 {
		return ""
	}
	return "\033[91m"
}

// BrightGreen foreground.
func (k Kind) BrightGreen() string {
	if k < Color4 {
		return ""
	}
	return "\033[92m"
}

// BrightYellow foreground.
func (k Kind) BrightYellow() string {
	if k < Color4 {
		return ""
	}
	return "\033[93m"
}

// BrightBlue foreground.
func (k Kind) BrightBlue() string {
	if k < Color4 {
		return ""
	}
	return "\033[94m"
}

// BrightMagenta foreground.
func (k Kind) BrightMagenta() string {
	if k < Color4 {
		return ""
	}
	return "\033[95m"
}

// BrightCyan foreground.
func (k Kind) BrightCyan() string {
	if k < Color4 {
		return ""
	}
	return "\033[96m"
}

// BrightWhite foreground.
func (k Kind) BrightWhite() string {
	if k < Color4 {
		return ""
	}
	return "\033[97m"
}

// FgReset resets the foreground color to the default.
func (k Kind) FgReset() string {
	if k < Color4 {
		return ""
	}
	return "\033[39m"
}

// BgBlack background.
func (k Kind) BgBlack() string {
	if k < Color4 {
		return ""
	}
	return "\033[40m"
}

// BgRed background.
func (k Kind) BgRed() string {
	if k < Color4 {
		return ""
	}
	return "\033[41m"
}

// BgGreen background.
func (k Kind) BgGreen() string {
	if k < Color4 {
		return ""
	}
	return "\033[42m"
}

// BgYellow background.
func (k Kind) BgYellow() string {
	if k < Color4 {
		return ""
	}
	return "\033[43m"
}

// BgBlue background.
func (k Kind) BgBlue() string {
	if k < Color4 {
		return ""
	}
	return "\033[44m"
}

// BgMagenta background.
func (k Kind) BgMagenta() string {
	if k < Color4 {
		return ""
	}
	return "\033[45m"
}

// BgCyan background.
func (k Kind) BgCyan() string {
	if k < Color4 {
		return ""
	}
	return "\033[46m"
}

// BgWhite background.
func (k Kind) BgWhite() string {
	if k < Color4 {
		return ""
	}
	return "\033[47m"
}

// BgGrey background. Same as BgBrightBlack.
func (k Kind) BgGrey() string {
	return k.BgBrightBlack()
}

// BgGray background. Same as BgBrightBlack.
func (k Kind) BgGray() string {
	return k.BgBrightBlack()
}

// BgBrightBlack background.
func (k Kind) BgBrightBlack() string {
	if k < Color4 {
		return ""
	}
	return "\033[100m"
}

// BgBrightRed background.
func (k Kind) BgBrightRed() string {
	if k < Color4 {
		return ""
	}
	return "\033[101m"
}

// BgBrightGreen background.
func (k Kind) BgBrightGreen() string {
	if k < Color4 {
		return ""
	}
	return "\033[102m"
}

// BgBrightYellow background.
func (k Kind) BgBrightYellow() string {
	if k < Color4 {
		return ""
	}
	return "\033[103m"
}

// BgBrightBlue background.
func (k Kind) BgBrightBlue() string {
	if k < Color4 {
		return ""
	}
	return "\033[104m"
}

// BgBrightMagenta background.
func (k Kind) BgBrightMagenta() string {
	if k < Color4 {
		return ""
	}
	return "\033[105m"
}

// BgBrightCyan background.
func (k Kind) BgBrightCyan() string {
	if k < Color4 {
		return ""
	}
	return "\033[106m"
}

// BgBrightWhite background.
func (k Kind) BgBrightWhite() string {
	if k < Color4 {
		return ""
	}
	return "\033[107m"
}

// BgReset resets the background color to the default.
func (k Kind) BgReset() string {
	if k < Color4 {
		return ""
	}
	return "\033[49m"
}

// Color256 sets the foreground color using an 8-bit ANSI color code.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit.
func (k Kind) Color256(color uint8) string {
	switch k {
	case Color4:
		return start + strconv.Itoa(int(ansi256ToAnsi16[color])) + end
	case Color8, Color24:
		return "\033[38;5;" + strconv.Itoa(int(color)) + end
	default:
		return ""
	}
}

// BgColor256 sets the background color using an 8-bit ANSI color code.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit.
func (k Kind) BgColor256(color uint8) string {
	switch k {
	case Color4:
		return start + strconv.Itoa(int(ansi256ToAnsi16[color]+10)) + end
	case Color8, Color24:
		return "\033[48;5;" + strconv.Itoa(int(color)) + end
	default:
		return ""
	}
}

// RGB sets the foreground color using RGB values.
func (k Kind) RGB(r, g, b uint8) string {
	switch k {
	case Color4:
		return start + strconv.Itoa(int(ansi256ToAnsi16[rgbToAnsi256(r, g, b)])) + end
	case Color8:
		return k.Color256(rgbToAnsi256(r, g, b))
	case Color24:
		return "\033[38;2;" + strconv.Itoa(int(r)) + ";" + strconv.Itoa(int(g)) + ";" + strconv.Itoa(int(b)) + end
	default:
		return ""
	}
}

// BgRGB sets the background color using RGB values.
func (k Kind) BgRGB(r, g, b uint8) string {
	switch k {
	case Color4:
		return start + strconv.Itoa(int(ansi256ToAnsi16[rgbToAnsi256(r, g, b)]+10)) + end
	case Color8:
		return k.BgColor256(rgbToAnsi256(r, g, b))
	case Color24:
		return "\033[48;2;" + strconv.Itoa(int(r)) + ";" + strconv.Itoa(int(g)) + ";" + strconv.Itoa(int(b)) + end
	default:
		return ""
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
