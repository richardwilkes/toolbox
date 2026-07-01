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

// TestWrapTextMeasuresRunesNotBytes verifies that WrapText measures column width by visible runes rather than bytes. A
// non-file writer reports a fixed 80-column width, so these cases are deterministic.
func TestWrapTextMeasuresRunesNotBytes(t *testing.T) {
	c := check.New(t)

	// Nine 4-rune Cyrillic words total 44 visible columns, well within the 80-column width, so they stay on one line.
	// With byte counting each word is 8 bytes (total 80), which exceeded the available width and forced a wrap.
	parts := make([]string, 9)
	for i := range parts {
		parts[i] = "абвг"
	}
	text := strings.Join(parts, " ")
	var buf bytes.Buffer
	xterm.NewAnsiWriter(&buf).WrapText("", text)
	c.Equal(text+"\n", buf.String())

	// A multibyte prefix must contribute only its visible width to the continuation indent. "→ " is 2 columns (4
	// bytes); byte counting indented continuation lines with 4 spaces and shrank the usable width.
	a := strings.Repeat("a", 40)
	b := strings.Repeat("b", 40)
	buf.Reset()
	xterm.NewAnsiWriter(&buf).WrapText("→ ", a+" "+b)
	c.Equal("→ "+a+"\n  "+b+"\n", buf.String())
}

// TestWrapTextResetsBudgetPerLine verifies that the remaining-width budget is reset for each input line, so a line that
// individually fits is not wrongly wrapped just because an earlier line consumed most of the width. A non-file writer
// reports a fixed 80-column width, so with an empty prefix each line has 79 columns available.
func TestWrapTextResetsBudgetPerLine(t *testing.T) {
	c := check.New(t)

	// Line 1 nearly fills the 79-column budget; line 2 ("aaaa bbbb", 9 columns) easily fits. Before the fix, line 2
	// inherited line 1's nearly-exhausted budget and was split into "aaaa\nbbbb".
	first := strings.Repeat("w", 70)
	var buf bytes.Buffer
	xterm.NewAnsiWriter(&buf).WrapText("", first+"\naaaa bbbb")
	c.Equal(first+"\naaaa bbbb\n", buf.String())

	// Sanity check the reset did not disable wrapping: two 40-column words cannot share a single 79-column line, so the
	// second must wrap onto its own line.
	a := strings.Repeat("a", 40)
	b := strings.Repeat("b", 40)
	buf.Reset()
	xterm.NewAnsiWriter(&buf).WrapText("", a+" "+b)
	c.Equal(a+"\n"+b+"\n", buf.String())
}
