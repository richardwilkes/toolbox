// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//nolint:goconst // I'd rather have the strings inline than extracted out into a constant for the tests.
package xstrings_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestFirstN(t *testing.T) {
	table := []struct {
		In  string
		Out string
		N   int
	}{
		{In: "abcd", N: 3, Out: "abc"},
		{In: "abcd", N: 5, Out: "abcd"},
		{In: "abcd", N: 0, Out: ""},
		{In: "abcd", N: -1, Out: ""},
		{In: "aécd", N: 3, Out: "aéc"},
		{In: "aécd", N: 5, Out: "aécd"},
		{In: "aécd", N: 0, Out: ""},
		{In: "aécd", N: -1, Out: ""},
	}
	c := check.New(t)
	for i, one := range table {
		c.Equal(one.Out, xstrings.FirstN(one.In, one.N), "#%d", i)
	}
}

func TestLastN(t *testing.T) {
	table := []struct {
		In  string
		Out string
		N   int
	}{
		{In: "abcd", N: 3, Out: "bcd"},
		{In: "abcd", N: 5, Out: "abcd"},
		{In: "abcd", N: 0, Out: ""},
		{In: "abcd", N: -1, Out: ""},
		{In: "aécd", N: 3, Out: "écd"},
		{In: "aécd", N: 5, Out: "aécd"},
		{In: "aécd", N: 0, Out: ""},
		{In: "aécd", N: -1, Out: ""},
	}
	c := check.New(t)
	for i, one := range table {
		c.Equal(one.Out, xstrings.LastN(one.In, one.N), "#%d", i)
	}
}

func TestTruncate(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", xstrings.Truncate("", 5, true))
	c.Equal("", xstrings.Truncate("", 5, false))
	c.Equal("", xstrings.Truncate("", 0, true))
	c.Equal("", xstrings.Truncate("", 0, false))

	// Test string shorter than count (no truncation needed)
	c.Equal("hello", xstrings.Truncate("hello", 10, true))
	c.Equal("hello", xstrings.Truncate("hello", 10, false))
	c.Equal("test", xstrings.Truncate("test", 4, true))
	c.Equal("test", xstrings.Truncate("test", 4, false))
	c.Equal("ab", xstrings.Truncate("ab", 5, true))
	c.Equal("ab", xstrings.Truncate("ab", 5, false))

	// Test string equal to count (no truncation needed)
	c.Equal("hello", xstrings.Truncate("hello", 5, true))
	c.Equal("hello", xstrings.Truncate("hello", 5, false))
	c.Equal("test", xstrings.Truncate("test", 4, true))
	c.Equal("test", xstrings.Truncate("test", 4, false))

	// Test keepFirst = true (truncate from end, ellipsis at end)
	c.Equal("hel…", xstrings.Truncate("hello", 3, true))
	c.Equal("hel…", xstrings.Truncate("hello world", 3, true))
	c.Equal("h…", xstrings.Truncate("hello", 1, true))
	c.Equal("hello wor…", xstrings.Truncate("hello world", 9, true))
	c.Equal("a…", xstrings.Truncate("abcdef", 1, true))
	c.Equal("ab…", xstrings.Truncate("abcdef", 2, true))

	// Test keepFirst = false (truncate from start, ellipsis at start)
	c.Equal("…llo", xstrings.Truncate("hello", 3, false))
	c.Equal("…rld", xstrings.Truncate("hello world", 3, false))
	c.Equal("…o", xstrings.Truncate("hello", 1, false))
	c.Equal("…llo world", xstrings.Truncate("hello world", 9, false))
	c.Equal("…f", xstrings.Truncate("abcdef", 1, false))
	c.Equal("…ef", xstrings.Truncate("abcdef", 2, false))

	// Test with count = 0
	c.Equal("…", xstrings.Truncate("hello", 0, true))
	c.Equal("…", xstrings.Truncate("hello", 0, false))
	c.Equal("…", xstrings.Truncate("hello world", 0, true))
	c.Equal("…", xstrings.Truncate("hello world", 0, false))

	// Test with negative count
	c.Equal("…", xstrings.Truncate("hello", -1, true))
	c.Equal("…", xstrings.Truncate("hello", -1, false))
	c.Equal("…", xstrings.Truncate("hello world", -5, true))
	c.Equal("…", xstrings.Truncate("hello world", -5, false))

	// Test with Unicode characters
	c.Equal("café", xstrings.Truncate("café", 4, true))
	c.Equal("café", xstrings.Truncate("café", 4, false))
	c.Equal("caf…", xstrings.Truncate("café", 3, true))
	c.Equal("…afé", xstrings.Truncate("café", 3, false))
	c.Equal("c…", xstrings.Truncate("café", 1, true))
	c.Equal("…é", xstrings.Truncate("café", 1, false))

	// Test with emoji and complex Unicode
	c.Equal("hello 🚀 world", xstrings.Truncate("hello 🚀 world", 13, true))
	c.Equal("hello 🚀 world", xstrings.Truncate("hello 🚀 world", 13, false))
	c.Equal("hello 🚀 …", xstrings.Truncate("hello 🚀 world", 8, true))
	c.Equal("… 🚀 world", xstrings.Truncate("hello 🚀 world", 8, false))
	c.Equal("hel…", xstrings.Truncate("hello 🚀 world", 3, true))
	c.Equal("…rld", xstrings.Truncate("hello 🚀 world", 3, false))

	// Test with single character strings
	c.Equal("a", xstrings.Truncate("a", 1, true))
	c.Equal("a", xstrings.Truncate("a", 1, false))
	c.Equal("a", xstrings.Truncate("a", 5, true))
	c.Equal("a", xstrings.Truncate("a", 5, false))
	c.Equal("…", xstrings.Truncate("a", 0, true))
	c.Equal("…", xstrings.Truncate("a", 0, false))

	// Test with special characters
	c.Equal("hello\nworld", xstrings.Truncate("hello\nworld", 11, true))
	c.Equal("hello\nworld", xstrings.Truncate("hello\nworld", 11, false))
	c.Equal("hello\n…", xstrings.Truncate("hello\nworld", 6, true))
	c.Equal("…\nworld", xstrings.Truncate("hello\nworld", 6, false))
	c.Equal("hell…", xstrings.Truncate("hello\tworld", 4, true))
	c.Equal("…orld", xstrings.Truncate("hello\tworld", 4, false))

	// Test with numbers and symbols
	c.Equal("123456789", xstrings.Truncate("123456789", 9, true))
	c.Equal("123456789", xstrings.Truncate("123456789", 9, false))
	c.Equal("12345…", xstrings.Truncate("123456789", 5, true))
	c.Equal("…56789", xstrings.Truncate("123456789", 5, false))
	c.Equal("test@exam…", xstrings.Truncate("test@example.com", 9, true))
	c.Equal("…ample.com", xstrings.Truncate("test@example.com", 9, false))

	// Test with whitespace
	c.Equal("  hello  ", xstrings.Truncate("  hello  ", 9, true))
	c.Equal("  hello  ", xstrings.Truncate("  hello  ", 9, false))
	c.Equal("  hel…", xstrings.Truncate("  hello  ", 5, true))
	c.Equal("…llo  ", xstrings.Truncate("  hello  ", 5, false))
	c.Equal(" …", xstrings.Truncate("  hello  ", 1, true))
	c.Equal("… ", xstrings.Truncate("  hello  ", 1, false))

	// Test edge case: exactly one more character than count
	c.Equal("hell…", xstrings.Truncate("hello", 4, true))
	c.Equal("…ello", xstrings.Truncate("hello", 4, false))
	c.Equal("hel…", xstrings.Truncate("hell", 3, true))
	c.Equal("…ell", xstrings.Truncate("hell", 3, false))

	// Test with very long strings
	longString := "This is a very long string that contains many words and should test how the truncate function handles extended content with proper ellipsis placement"
	c.Equal("This is a very…", xstrings.Truncate(longString, 14, true))
	c.Equal("…psis placement", xstrings.Truncate(longString, 14, false))
	c.Equal("This…", xstrings.Truncate(longString, 4, true))
	c.Equal("…ment", xstrings.Truncate(longString, 4, false))

	// Test with international characters
	c.Equal("Ñoño test", xstrings.Truncate("Ñoño test", 9, true))
	c.Equal("Ñoño test", xstrings.Truncate("Ñoño test", 9, false))
	c.Equal("Ñoño…", xstrings.Truncate("Ñoño test", 4, true))
	c.Equal("…test", xstrings.Truncate("Ñoño test", 4, false))
	c.Equal("北京上海…", xstrings.Truncate("北京上海广州", 4, true))
	c.Equal("…上海广州", xstrings.Truncate("北京上海广州", 4, false))
}
