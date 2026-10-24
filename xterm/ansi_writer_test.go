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

// TestWrapTextWrapsAtVisibleWidth relies on a non-file writer being treated as 80 columns wide, which leaves a
// 79-column budget with an empty prefix.
func TestWrapTextWrapsAtVisibleWidth(t *testing.T) {
	c := check.New(t)
	const eraseLine = "\033[2K"
	cyrillic := strings.TrimSpace(strings.Repeat("абвг ", 9))
	first := strings.Repeat("w", 70)
	a := strings.Repeat("a", 40)
	b := strings.Repeat("b", 40)
	for _, tc := range []struct {
		prefix string
		text   string
		want   string
	}{
		// The budget resets for each input line: line 2 fits on its own even though line 1 nearly exhausted it.
		{text: first + "\naaaa bbbb", want: first + "\naaaa bbbb\n"},
		// Two 40-column words cannot share a 79-column line.
		{text: a + " " + b, want: a + "\n" + b + "\n"},
		// Nine 4-rune Cyrillic words are 44 visible columns (80 bytes), so they stay on one line.
		{text: cyrillic, want: cyrillic + "\n"},
		// "→ " is 2 columns (4 bytes), so continuation lines are indented by 2 spaces.
		{prefix: "→ ", text: a + " " + b, want: "→ " + a + "\n  " + b + "\n"},
		// 40 + 1 + 37 columns fit on one line only if the 4-byte erase-line escape is zero-width.
		{
			text: a + " " + eraseLine + strings.Repeat("b", 37),
			want: a + " " + eraseLine + strings.Repeat("b", 37) + "\n",
		},
		// A non-SGR escape in the prefix is likewise zero-width, so it does not widen the continuation indent.
		{prefix: eraseLine + "→ ", text: a + " " + b, want: eraseLine + "→ " + a + "\n  " + b + "\n"},
	} {
		var buf bytes.Buffer
		xterm.NewAnsiWriter(&buf).WrapText(tc.prefix, tc.text)
		c.Equal(tc.want, buf.String(), "WrapText(%q, %q)", tc.prefix, tc.text)
	}
}
