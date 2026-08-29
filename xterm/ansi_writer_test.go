// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xterm_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xterm"
)

func TestWrapTextMeasuresRunesNotBytes(t *testing.T) {
	c := check.New(t)

	// A non-file writer is treated as 80 columns wide. Nine 4-rune Cyrillic words are 44 visible columns (80 bytes),
	// so they must stay on one line.
	parts := make([]string, 9)
	for i := range parts {
		parts[i] = "абвг"
	}
	text := strings.Join(parts, " ")
	var buf bytes.Buffer
	xterm.NewAnsiWriter(&buf).WrapText("", text)
	c.Equal(text+"\n", buf.String())

	// "→ " is 2 columns (4 bytes), so continuation lines must be indented by 2 spaces.
	a := strings.Repeat("a", 40)
	b := strings.Repeat("b", 40)
	buf.Reset()
	xterm.NewAnsiWriter(&buf).WrapText("→ ", a+" "+b)
	c.Equal("→ "+a+"\n  "+b+"\n", buf.String())
}

func TestWrapTextResetsBudgetPerLine(t *testing.T) {
	c := check.New(t)

	// The width budget (79 columns with an empty prefix) must reset for each input line: line 2 fits on its own even
	// though line 1 nearly exhausted it.
	first := strings.Repeat("w", 70)
	var buf bytes.Buffer
	xterm.NewAnsiWriter(&buf).WrapText("", first+"\naaaa bbbb")
	c.Equal(first+"\naaaa bbbb\n", buf.String())

	// Two 40-column words still cannot share a 79-column line.
	a := strings.Repeat("a", 40)
	b := strings.Repeat("b", 40)
	buf.Reset()
	xterm.NewAnsiWriter(&buf).WrapText("", a+" "+b)
	c.Equal(a+"\n"+b+"\n", buf.String())
}

func TestWrapTextIgnoresNonSGREscapes(t *testing.T) {
	c := check.New(t)

	// word1 and word2 fit on one 79-column line (40 + 1 + 37 = 78) only if the 4-byte erase-line escape is treated as
	// zero-width.
	word1 := strings.Repeat("a", 40)
	word2 := strings.Repeat("b", 37)
	const eraseLine = "\033[2K"

	var plain bytes.Buffer
	xterm.NewAnsiWriter(&plain).WrapText("", word1+" "+word2)

	var withEscape bytes.Buffer
	xterm.NewAnsiWriter(&withEscape).WrapText("", word1+" "+eraseLine+word2)

	// Stripping the escape from the escaped output must yield the plain output.
	c.Equal(plain.String(), strings.ReplaceAll(withEscape.String(), eraseLine, ""))

	// Both stay on a single line, with the escape passed through untouched.
	c.Equal(word1+" "+word2+"\n", plain.String())
	c.Equal(word1+" "+eraseLine+word2+"\n", withEscape.String())

	// A non-SGR escape in the prefix is likewise zero-width, so it does not widen the continuation indent.
	x := strings.Repeat("x", 40)
	y := strings.Repeat("y", 40)
	var buf bytes.Buffer
	xterm.NewAnsiWriter(&buf).WrapText(eraseLine+"→ ", x+" "+y)
	c.Equal(eraseLine+"→ "+x+"\n  "+y+"\n", buf.String())
}
