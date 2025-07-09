// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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

	// Test short text with prefix - should fit on one line
	c.Equal("// short", xstrings.Wrap("// ", "short", 78))

	// Test longer text with prefix that needs wrapping
	c.Equal("// some text\n// that is\n// longer", xstrings.Wrap("// ", "some text that is longer", 12))

	// Test text with embedded line feeds and prefix
	c.Equal("// some text\n// with embedded\n// line feeds", xstrings.Wrap("// ", "some text\nwith embedded line feeds", 16))

	// Test longer text without prefix that needs wrapping
	c.Equal("some text\nthat is\nlonger", xstrings.Wrap("", "some text that is longer", 12))

	// Test text without prefix with very short max columns
	c.Equal("some\ntext\nthat\nis\nlonger", xstrings.Wrap("", "some text that is longer", 4))

	// Test text without prefix with punctuation and very short max columns
	c.Equal("some\ntext\nthat\nis\nlonger,\nyep", xstrings.Wrap("", "some text that is longer, yep", 4))

	// Test text with embedded line feeds and no prefix
	c.Equal("some text\nwith embedded\nline feeds", xstrings.Wrap("", "some text\nwith embedded line feeds", 16))

	// Test empty text
	c.Equal("", xstrings.Wrap("", "", 10))

	// Test empty text with prefix
	c.Equal("// ", xstrings.Wrap("// ", "", 10))

	// Test single word that exceeds max columns (should not break)
	c.Equal("verylongwordthatexceedsmaxcolumns", xstrings.Wrap("", "verylongwordthatexceedsmaxcolumns", 10))

	// Test single word with prefix that exceeds max columns
	c.Equal("// verylongwordthatexceedsmaxcolumns", xstrings.Wrap("// ", "verylongwordthatexceedsmaxcolumns", 10))

	// Test single character
	c.Equal("a", xstrings.Wrap("", "a", 10))

	// Test single character with prefix
	c.Equal("// a", xstrings.Wrap("// ", "a", 10))

	// Test text with only spaces
	c.Equal("", xstrings.Wrap("", "   ", 10))

	// Test text with only spaces and prefix
	c.Equal("// ", xstrings.Wrap("// ", "   ", 10))

	// Test max columns equal to prefix length
	c.Equal("// word", xstrings.Wrap("// ", "word", 2))

	// Test max columns less than prefix length
	c.Equal("// word", xstrings.Wrap("// ", "word", 1))

	// Test multiple consecutive spaces (extra spaces are collapsed by strings.Fields)
	c.Equal("word1 word2", xstrings.Wrap("", "word1    word2", 15))

	// Test text with punctuation
	c.Equal("Hello,\nworld!", xstrings.Wrap("", "Hello, world!", 8))

	// Test text with numbers
	c.Equal("# Test 123\n# and 456", xstrings.Wrap("# ", "Test 123 and 456", 12))

	// Test text with Unicode characters (emotion doesn't have accent)
	c.Equal("Café\nemotion\n🚀", xstrings.Wrap("", "Café emotion 🚀", 8))

	// Test text with Unicode and prefix (emotion doesn't have accent)
	c.Equal("// Café\n// emotion\n// 🚀", xstrings.Wrap("// ", "Café emotion 🚀", 10))

	// Test text with tabs (tabs are treated as whitespace and collapsed)
	c.Equal("word1 word2", xstrings.Wrap("", "word1\tword2", 20))

	// Test text with special symbols
	c.Equal("@user\n#hashtag\n$money", xstrings.Wrap("", "@user #hashtag $money", 8))

	// Test multiple empty lines
	c.Equal("\n", xstrings.Wrap("", "\n", 10))

	// Test multiple empty lines with prefix
	c.Equal("// \n// ", xstrings.Wrap("// ", "\n", 10))

	// Test text with multiple line breaks
	c.Equal("line1\n\nline3", xstrings.Wrap("", "line1\n\nline3", 10))

	// Test text with multiple line breaks and prefix
	c.Equal("// line1\n// \n// line3", xstrings.Wrap("// ", "line1\n\nline3", 15))

	// Test text starting with newline
	c.Equal("\nfirst line", xstrings.Wrap("", "\nfirst line", 15))

	// Test text ending with newline
	c.Equal("last line\n", xstrings.Wrap("", "last line\n", 15))

	// Test text with only newlines
	c.Equal("\n\n", xstrings.Wrap("", "\n\n", 10))

	// Test with indentation prefix
	c.Equal("    word1\n    word2", xstrings.Wrap("    ", "word1 word2", 10))

	// Test with bullet point prefix
	c.Equal("* item\n* one", xstrings.Wrap("* ", "item one", 8))

	// Test with numbered prefix
	c.Equal("1. first\n1. item", xstrings.Wrap("1. ", "first item", 10))

	// Test with long prefix
	c.Equal("PREFIX: word\nPREFIX: two", xstrings.Wrap("PREFIX: ", "word two", 15))

	// Test with Unicode prefix
	c.Equal("🔸 hello\n🔸 world", xstrings.Wrap("🔸 ", "hello world", 10))

	// Test with empty prefix (equivalent to no prefix)
	c.Equal("hello world", xstrings.Wrap("", "hello world", 15))

	// Test with max columns = 1 (each word gets its own line)
	c.Equal("a\nb\nc", xstrings.Wrap("", "a b c", 1))

	// Test with max columns = 0 (edge case - each word gets its own line)
	c.Equal("a\nb\nc", xstrings.Wrap("", "a b c", 0))

	// Test with negative max columns (edge case - each word gets its own line)
	c.Equal("a\nb\nc", xstrings.Wrap("", "a b c", -5))

	// Test word exactly fitting the line
	c.Equal("hello", xstrings.Wrap("", "hello", 5))

	// Test word exactly fitting with prefix
	c.Equal("> hello", xstrings.Wrap("> ", "hello", 7))

	// Test multiple words exactly fitting
	c.Equal("hi\nbye", xstrings.Wrap("", "hi bye", 3))

	// Test very large max columns
	c.Equal("short text here", xstrings.Wrap("", "short text here", 1000))
}
