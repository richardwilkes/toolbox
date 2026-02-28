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

func TestUnquote(t *testing.T) {
	c := check.New(t)

	// Test with double quotes
	c.Equal("hello", xstrings.Unquote(`"hello"`))
	c.Equal("world", xstrings.Unquote(`"world"`))
	c.Equal("", xstrings.Unquote(`""`))

	// Test with single quotes
	c.Equal("hello", xstrings.Unquote("'hello'"))
	c.Equal("world", xstrings.Unquote("'world'"))
	c.Equal("", xstrings.Unquote("''"))

	// Test with mixed quotes (should not unquote)
	c.Equal("'hello\"", xstrings.Unquote("'hello\""))
	c.Equal("\"hello'", xstrings.Unquote("\"hello'"))

	// Test with no quotes
	c.Equal("hello", xstrings.Unquote("hello"))
	c.Equal("world", xstrings.Unquote("world"))

	// Test with single character
	c.Equal("a", xstrings.Unquote("a"))
	c.Equal("\"", xstrings.Unquote("\""))
	c.Equal("'", xstrings.Unquote("'"))

	// Test with empty string
	c.Equal("", xstrings.Unquote(""))

	// Test with only opening quote
	c.Equal("\"hello", xstrings.Unquote("\"hello"))
	c.Equal("'hello", xstrings.Unquote("'hello"))

	// Test with only closing quote
	c.Equal("hello\"", xstrings.Unquote("hello\""))
	c.Equal("hello'", xstrings.Unquote("hello'"))

	// Test with quotes in the middle
	c.Equal("he\"llo", xstrings.Unquote("he\"llo"))
	c.Equal("he'llo", xstrings.Unquote("he'llo"))

	// Test with escaped quotes (should still unquote outer quotes)
	c.Equal("he\\\"llo", xstrings.Unquote("\"he\\\"llo\""))
	c.Equal("he\\'llo", xstrings.Unquote("'he\\'llo'"))

	// Test with nested quotes
	c.Equal("'hello'", xstrings.Unquote("\"'hello'\""))
	c.Equal("\"hello\"", xstrings.Unquote("'\"hello\"'"))

	// Test with multiple quotes at start/end
	c.Equal("\"hello\"", xstrings.Unquote("\"\"hello\"\""))
	c.Equal("'hello'", xstrings.Unquote("''hello''"))

	// Test with Unicode content
	c.Equal("café", xstrings.Unquote("\"café\""))
	c.Equal("北京", xstrings.Unquote("'北京'"))
	c.Equal("🚀🎉", xstrings.Unquote("\"🚀🎉\""))

	// Test with special characters
	c.Equal("hello\nworld", xstrings.Unquote("\"hello\nworld\""))
	c.Equal("hello\tworld", xstrings.Unquote("'hello\tworld'"))
	c.Equal("hello\\world", xstrings.Unquote("\"hello\\world\""))

	// Test with whitespace
	c.Equal(" hello ", xstrings.Unquote("\" hello \""))
	c.Equal("\thello\t", xstrings.Unquote("'\thello\t'"))

	// Test with numeric content
	c.Equal("123", xstrings.Unquote("\"123\""))
	c.Equal("3.14", xstrings.Unquote("'3.14'"))

	// Test with JSON-like content
	c.Equal("{\"key\":\"value\"}", xstrings.Unquote("\"{\"key\":\"value\"}\""))

	// Test with very long content
	longContent := ""
	for i := range 1000 {
		longContent += string(rune('a' + (i % 26)))
	}
	quotedLong := "\"" + longContent + "\""
	c.Equal(longContent, xstrings.Unquote(quotedLong))

	// Test Unicode quote characters (should not be processed)
	unicodeQuotes := "\u201chello\u201d" // Unicode left/right double quotes
	c.Equal(unicodeQuotes, xstrings.Unquote(unicodeQuotes))

	leftQuote := "\u2018hello\u2019" // Unicode left/right single quotes
	c.Equal(leftQuote, xstrings.Unquote(leftQuote))

	// Test with various quote-like characters
	backticks := "`hello`"
	c.Equal(backticks, xstrings.Unquote(backticks))

	guillemets := "«hello»"
	c.Equal(guillemets, xstrings.Unquote(guillemets))

	// Test with mathematical symbols that might look like quotes
	primes := "′hello′"
	c.Equal(primes, xstrings.Unquote(primes))

	doublePrimes := "″hello″"
	c.Equal(doublePrimes, xstrings.Unquote(doublePrimes))

	// Test with control characters
	withNull := "\x00hello\x00"
	c.Equal(withNull, xstrings.Unquote(withNull))

	// Test with very short strings
	oneChar := "a"
	c.Equal(oneChar, xstrings.Unquote(oneChar))

	twoChars := "ab"
	c.Equal(twoChars, xstrings.Unquote(twoChars))
}
