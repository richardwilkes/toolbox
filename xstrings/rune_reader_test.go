// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings_test

import (
	"io"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestRuneReader(t *testing.T) {
	c := check.New(t)

	// Test that xstrings.RuneReader implements io.xstrings.RuneReader interface
	var _ io.RuneReader = &xstrings.RuneReader{}

	// Test empty rune slice
	rr := &xstrings.RuneReader{Src: []rune{}, Pos: 0}
	r, size, err := rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)
	c.Equal(0, rr.Pos)

	// Test nil rune slice (should behave like empty)
	rr = &xstrings.RuneReader{Src: nil, Pos: 0}
	r, size, err = rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)
	c.Equal(0, rr.Pos)

	// Test single ASCII character
	rr = &xstrings.RuneReader{Src: []rune{'a'}, Pos: 0}
	r, size, err = rr.ReadRune()
	c.Equal('a', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(1, rr.Pos)

	// Try to read beyond end
	r, size, err = rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)
	c.Equal(1, rr.Pos)

	// Test single Unicode character
	rr = &xstrings.RuneReader{Src: []rune{'🚀'}, Pos: 0}
	r, size, err = rr.ReadRune()
	c.Equal('🚀', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(1, rr.Pos)

	// Test multiple ASCII characters
	runes := []rune{'h', 'e', 'l', 'l', 'o'}
	rr = &xstrings.RuneReader{Src: runes, Pos: 0}

	// Read first character
	r, size, err = rr.ReadRune()
	c.Equal('h', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(1, rr.Pos)

	// Read second character
	r, size, err = rr.ReadRune()
	c.Equal('e', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(2, rr.Pos)

	// Read third character
	r, size, err = rr.ReadRune()
	c.Equal('l', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(3, rr.Pos)

	// Read fourth character
	r, size, err = rr.ReadRune()
	c.Equal('l', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(4, rr.Pos)

	// Read fifth character
	r, size, err = rr.ReadRune()
	c.Equal('o', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(5, rr.Pos)

	// Try to read beyond end
	r, size, err = rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)
	c.Equal(5, rr.Pos)

	// Test mixed ASCII and Unicode characters
	runes = []rune{'H', 'e', 'l', 'l', 'o', ' ', '🌍', '!', ' ', '🚀', ' ', 'T', 'e', 's', 't'}
	rr = &xstrings.RuneReader{Src: runes, Pos: 0}

	// Read ASCII characters
	r, size, err = rr.ReadRune()
	c.Equal('H', r)
	c.Equal(1, size)
	c.NoError(err)

	// Skip to emoji
	rr.Pos = 6
	r, size, err = rr.ReadRune()
	c.Equal('🌍', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(7, rr.Pos)

	// Read exclamation
	r, size, err = rr.ReadRune()
	c.Equal('!', r)
	c.Equal(1, size)
	c.NoError(err)

	// Skip to rocket emoji
	rr.Pos = 9
	r, size, err = rr.ReadRune()
	c.Equal('🚀', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(10, rr.Pos)

	// Test special characters and symbols
	runes = []rune{'\n', '\t', '\r', ' ', '@', '#', '$', '%', '^', '&', '*'}
	rr = &xstrings.RuneReader{Src: runes, Pos: 0}

	// Read newline
	r, size, err = rr.ReadRune()
	c.Equal('\n', r)
	c.Equal(1, size)
	c.NoError(err)

	// Read tab
	r, size, err = rr.ReadRune()
	c.Equal('\t', r)
	c.Equal(1, size)
	c.NoError(err)

	// Read carriage return
	r, size, err = rr.ReadRune()
	c.Equal('\r', r)
	c.Equal(1, size)
	c.NoError(err)

	// Read space
	r, size, err = rr.ReadRune()
	c.Equal(' ', r)
	c.Equal(1, size)
	c.NoError(err)

	// Read symbols
	r, size, err = rr.ReadRune()
	c.Equal('@', r)
	c.Equal(1, size)
	c.NoError(err)

	// Test numeric characters
	runes = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	rr = &xstrings.RuneReader{Src: runes, Pos: 0}

	for i, expected := range runes {
		r, size, err = rr.ReadRune()
		c.Equal(expected, r)
		c.Equal(1, size)
		c.NoError(err)
		c.Equal(i+1, rr.Pos)
	}

	// Try to read beyond end
	r, size, err = rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)

	runes = []rune{'a', 'b', 'c', 'd', 'e'}
	rr = &xstrings.RuneReader{Src: runes, Pos: 0}

	// Read first character
	r, size, err = rr.ReadRune()
	c.Equal('a', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(1, rr.Pos)

	// Manually set position to middle
	rr.Pos = 2
	r, size, err = rr.ReadRune()
	c.Equal('c', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(3, rr.Pos)

	// Manually set position to end
	rr.Pos = 5
	r, size, err = rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)
	c.Equal(5, rr.Pos)

	// Manually set position beyond end
	rr.Pos = 10
	r, size, err = rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)
	c.Equal(10, rr.Pos)

	// Reset position to beginning
	rr.Pos = 0
	r, size, err = rr.ReadRune()
	c.Equal('a', r)
	c.Equal(1, size)
	c.NoError(err)
	c.Equal(1, rr.Pos)

	// Test various Unicode categories
	runes = []rune{
		'A', // ASCII uppercase
		'z', // ASCII lowercase
		'5', // ASCII digit
		'ñ', // Latin extended
		'ü', // Latin extended
		'α', // Greek
		'β', // Greek
		'中', // CJK
		'国', // CJK
		'🎉', // Emoji
		'🚀', // Emoji
		'📝', // Emoji
		'∑', // Mathematical symbol
		'∞', // Mathematical symbol
		'€', // Currency symbol
	}

	rr = &xstrings.RuneReader{Src: runes, Pos: 0}

	for i, expected := range runes {
		r, size, err = rr.ReadRune()
		c.Equal(expected, r)
		c.Equal(1, size)
		c.NoError(err)
		c.Equal(i+1, rr.Pos)
	}

	// Verify we're at the end
	r, size, err = rr.ReadRune()
	c.Equal(rune(-1), r)
	c.Equal(0, size)
	c.Equal(io.EOF, err)

	// Test reading the same xstrings.RuneReader multiple times after EOF
	rr = &xstrings.RuneReader{Src: []rune{'x'}, Pos: 0}

	// First read should succeed
	r, size, err = rr.ReadRune()
	c.Equal('x', r)
	c.Equal(1, size)
	c.NoError(err)

	// Multiple reads after EOF should all return EOF
	for range 5 {
		r, size, err = rr.ReadRune()
		c.Equal(rune(-1), r)
		c.Equal(0, size)
		c.Equal(io.EOF, err)
	}

	// Test that size is always 1 for successful reads, regardless of actual Unicode byte size
	runes = []rune{
		'A', // 1 byte in UTF-8
		'ü', // 2 bytes in UTF-8
		'中', // 3 bytes in UTF-8
		'🚀', // 4 bytes in UTF-8
	}

	rr = &xstrings.RuneReader{Src: runes, Pos: 0}

	for _, expected := range runes {
		r, size, err = rr.ReadRune()
		c.Equal(expected, r)
		c.Equal(1, size) // Size should always be 1 for rune count, not byte count
		c.NoError(err)
	}
}
