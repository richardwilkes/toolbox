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
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestUnquote(t *testing.T) {
	c := check.New(t)

	c.Equal("hello", xstrings.Unquote(`"hello"`))
	c.Equal("world", xstrings.Unquote(`"world"`))
	c.Equal("", xstrings.Unquote(`""`))

	c.Equal("hello", xstrings.Unquote("'hello'"))
	c.Equal("world", xstrings.Unquote("'world'"))
	c.Equal("", xstrings.Unquote("''"))

	c.Equal("'hello\"", xstrings.Unquote("'hello\""))
	c.Equal("\"hello'", xstrings.Unquote("\"hello'"))

	c.Equal("hello", xstrings.Unquote("hello"))
	c.Equal("world", xstrings.Unquote("world"))

	c.Equal("a", xstrings.Unquote("a"))
	c.Equal("\"", xstrings.Unquote("\""))
	c.Equal("'", xstrings.Unquote("'"))

	c.Equal("", xstrings.Unquote(""))

	c.Equal("\"hello", xstrings.Unquote("\"hello"))
	c.Equal("'hello", xstrings.Unquote("'hello"))

	c.Equal("hello\"", xstrings.Unquote("hello\""))
	c.Equal("hello'", xstrings.Unquote("hello'"))

	c.Equal("he\"llo", xstrings.Unquote("he\"llo"))
	c.Equal("he'llo", xstrings.Unquote("he'llo"))

	c.Equal("he\\\"llo", xstrings.Unquote("\"he\\\"llo\""))
	c.Equal("he\\'llo", xstrings.Unquote("'he\\'llo'"))

	c.Equal("'hello'", xstrings.Unquote("\"'hello'\""))
	c.Equal("\"hello\"", xstrings.Unquote("'\"hello\"'"))

	c.Equal("\"hello\"", xstrings.Unquote("\"\"hello\"\""))
	c.Equal("'hello'", xstrings.Unquote("''hello''"))

	c.Equal("café", xstrings.Unquote("\"café\""))
	c.Equal("北京", xstrings.Unquote("'北京'"))
	c.Equal("🚀🎉", xstrings.Unquote("\"🚀🎉\""))

	c.Equal("hello\nworld", xstrings.Unquote("\"hello\nworld\""))
	c.Equal("hello\tworld", xstrings.Unquote("'hello\tworld'"))
	c.Equal("hello\\world", xstrings.Unquote("\"hello\\world\""))

	c.Equal(" hello ", xstrings.Unquote("\" hello \""))
	c.Equal("\thello\t", xstrings.Unquote("'\thello\t'"))

	c.Equal("123", xstrings.Unquote("\"123\""))
	c.Equal("3.14", xstrings.Unquote("'3.14'"))

	c.Equal("{\"key\":\"value\"}", xstrings.Unquote("\"{\"key\":\"value\"}\""))

	var longContent strings.Builder
	for i := range 1000 {
		longContent.WriteString(string(rune('a' + (i % 26))))
	}
	c.Equal(longContent.String(), xstrings.Unquote("\""+longContent.String()+"\""))

	unicodeQuotes := "\u201chello\u201d" // Unicode left/right double quotes
	c.Equal(unicodeQuotes, xstrings.Unquote(unicodeQuotes))

	leftQuote := "\u2018hello\u2019" // Unicode left/right single quotes
	c.Equal(leftQuote, xstrings.Unquote(leftQuote))

	backticks := "`hello`"
	c.Equal(backticks, xstrings.Unquote(backticks))

	guillemets := "«hello»"
	c.Equal(guillemets, xstrings.Unquote(guillemets))

	primes := "′hello′"
	c.Equal(primes, xstrings.Unquote(primes))

	doublePrimes := "″hello″"
	c.Equal(doublePrimes, xstrings.Unquote(doublePrimes))

	withNull := "\x00hello\x00"
	c.Equal(withNull, xstrings.Unquote(withNull))

	oneChar := "a"
	c.Equal(oneChar, xstrings.Unquote(oneChar))

	twoChars := "ab"
	c.Equal(twoChars, xstrings.Unquote(twoChars))
}
