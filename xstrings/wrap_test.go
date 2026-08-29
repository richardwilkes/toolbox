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
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestWrap(t *testing.T) {
	c := check.New(t)

	c.Equal("// short", xstrings.Wrap("// ", "short", 78))

	c.Equal("// some text\n// that is\n// longer", xstrings.Wrap("// ", "some text that is longer", 12))

	c.Equal("// some text\n// with embedded\n// line feeds", xstrings.Wrap("// ", "some text\nwith embedded line feeds", 16))

	c.Equal("some text\nthat is\nlonger", xstrings.Wrap("", "some text that is longer", 12))

	c.Equal("some\ntext\nthat\nis\nlonger", xstrings.Wrap("", "some text that is longer", 4))

	c.Equal("some\ntext\nthat\nis\nlonger,\nyep", xstrings.Wrap("", "some text that is longer, yep", 4))

	c.Equal("some text\nwith embedded\nline feeds", xstrings.Wrap("", "some text\nwith embedded line feeds", 16))

	c.Equal("", xstrings.Wrap("", "", 10))

	c.Equal("// ", xstrings.Wrap("// ", "", 10))

	c.Equal("verylongwordthatexceedsmaxcolumns", xstrings.Wrap("", "verylongwordthatexceedsmaxcolumns", 10))

	c.Equal("// verylongwordthatexceedsmaxcolumns", xstrings.Wrap("// ", "verylongwordthatexceedsmaxcolumns", 10))

	c.Equal("a", xstrings.Wrap("", "a", 10))

	c.Equal("// a", xstrings.Wrap("// ", "a", 10))

	c.Equal("", xstrings.Wrap("", "   ", 10))

	c.Equal("// ", xstrings.Wrap("// ", "   ", 10))

	c.Equal("// word", xstrings.Wrap("// ", "word", 2))

	c.Equal("// word", xstrings.Wrap("// ", "word", 1))

	c.Equal("word1 word2", xstrings.Wrap("", "word1    word2", 15))

	c.Equal("Hello,\nworld!", xstrings.Wrap("", "Hello, world!", 8))

	c.Equal("# Test 123\n# and 456", xstrings.Wrap("# ", "Test 123 and 456", 12))

	c.Equal("Café\nemotion\n🚀", xstrings.Wrap("", "Café emotion 🚀", 8))

	c.Equal("// Café\n// emotion\n// 🚀", xstrings.Wrap("// ", "Café emotion 🚀", 10))

	c.Equal("word1 word2", xstrings.Wrap("", "word1\tword2", 20))

	c.Equal("@user\n#hashtag\n$money", xstrings.Wrap("", "@user #hashtag $money", 8))

	c.Equal("\n", xstrings.Wrap("", "\n", 10))

	c.Equal("// \n// ", xstrings.Wrap("// ", "\n", 10))

	c.Equal("line1\n\nline3", xstrings.Wrap("", "line1\n\nline3", 10))

	c.Equal("// line1\n// \n// line3", xstrings.Wrap("// ", "line1\n\nline3", 15))

	c.Equal("\nfirst line", xstrings.Wrap("", "\nfirst line", 15))

	c.Equal("last line\n", xstrings.Wrap("", "last line\n", 15))

	c.Equal("\n\n", xstrings.Wrap("", "\n\n", 10))

	c.Equal("    word1\n    word2", xstrings.Wrap("    ", "word1 word2", 10))

	c.Equal("* item\n* one", xstrings.Wrap("* ", "item one", 8))

	c.Equal("1. first\n1. item", xstrings.Wrap("1. ", "first item", 10))

	c.Equal("PREFIX: word\nPREFIX: two", xstrings.Wrap("PREFIX: ", "word two", 15))

	c.Equal("🔸 hello\n🔸 world", xstrings.Wrap("🔸 ", "hello world", 10))

	c.Equal("hello world", xstrings.Wrap("", "hello world", 15))

	c.Equal("a\nb\nc", xstrings.Wrap("", "a b c", 1))

	c.Equal("a\nb\nc", xstrings.Wrap("", "a b c", 0))

	c.Equal("a\nb\nc", xstrings.Wrap("", "a b c", -5))

	c.Equal("hello", xstrings.Wrap("", "hello", 5))

	c.Equal("> hello", xstrings.Wrap("> ", "hello", 7))

	c.Equal("hi\nbye", xstrings.Wrap("", "hi bye", 3))

	c.Equal("short text here", xstrings.Wrap("", "short text here", 1000))
}

// TestWrapMeasuresRunesNotBytes verifies that wrapping measures rune width rather than byte length, so multibyte text
// is not wrapped prematurely and a multibyte prefix does not shrink the available width.
func TestWrapMeasuresRunesNotBytes(t *testing.T) {
	c := check.New(t)

	// 9 rune columns (4 + 1 + 4); byte counting would see 17 bytes and wrap.
	c.Equal("абвг абвг", xstrings.Wrap("", "абвг абвг", 9))
	// One column short forces the wrap.
	c.Equal("абвг\nабвг", xstrings.Wrap("", "абвг абвг", 8))

	// Three two-rune words fit in 8 columns (2 + 1 + 2 + 1 + 2); byte counting would see 20 bytes and wrap.
	c.Equal("日本 中文 한국", xstrings.Wrap("", "日本 中文 한국", 8))

	// The prefix "→ " is 2 columns (4 bytes): 13 columns fit "hello world" exactly, while 8 force a wrap.
	c.Equal("→ hello\n→ world", xstrings.Wrap("→ ", "hello world", 8))
	c.Equal("→ hello world", xstrings.Wrap("→ ", "hello world", 13))
}
