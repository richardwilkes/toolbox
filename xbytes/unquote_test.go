// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xbytes_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xbytes"
)

func TestUnquoteBytes(t *testing.T) {
	c := check.New(t)

	c.Equal([]byte("hello"), xbytes.Unquote([]byte(`"hello"`)))
	c.Equal([]byte("world"), xbytes.Unquote([]byte(`"world"`)))
	c.Equal([]byte(""), xbytes.Unquote([]byte(`""`)))

	c.Equal([]byte("hello"), xbytes.Unquote([]byte("'hello'")))
	c.Equal([]byte("world"), xbytes.Unquote([]byte("'world'")))
	c.Equal([]byte(""), xbytes.Unquote([]byte("''")))

	c.Equal([]byte("'hello\""), xbytes.Unquote([]byte("'hello\"")))
	c.Equal([]byte("\"hello'"), xbytes.Unquote([]byte("\"hello'")))

	c.Equal([]byte("hello"), xbytes.Unquote([]byte("hello")))
	c.Equal([]byte("world"), xbytes.Unquote([]byte("world")))

	c.Equal([]byte("a"), xbytes.Unquote([]byte("a")))
	c.Equal([]byte("\""), xbytes.Unquote([]byte("\"")))
	c.Equal([]byte("'"), xbytes.Unquote([]byte("'")))

	c.Equal([]byte(""), xbytes.Unquote([]byte("")))

	c.Equal([]byte("\"hello"), xbytes.Unquote([]byte("\"hello")))
	c.Equal([]byte("'hello"), xbytes.Unquote([]byte("'hello")))

	c.Equal([]byte("hello\""), xbytes.Unquote([]byte("hello\"")))
	c.Equal([]byte("hello'"), xbytes.Unquote([]byte("hello'")))

	c.Equal([]byte("he\"llo"), xbytes.Unquote([]byte("he\"llo")))
	c.Equal([]byte("he'llo"), xbytes.Unquote([]byte("he'llo")))

	c.Equal([]byte("he\\\"llo"), xbytes.Unquote([]byte("\"he\\\"llo\"")))
	c.Equal([]byte("he\\'llo"), xbytes.Unquote([]byte("'he\\'llo'")))

	c.Equal([]byte("'hello'"), xbytes.Unquote([]byte("\"'hello'\"")))
	c.Equal([]byte("\"hello\""), xbytes.Unquote([]byte("'\"hello\"'")))

	c.Equal([]byte("\"hello\""), xbytes.Unquote([]byte("\"\"hello\"\"")))
	c.Equal([]byte("'hello'"), xbytes.Unquote([]byte("''hello''")))

	c.Equal([]byte("café"), xbytes.Unquote([]byte("\"café\"")))
	c.Equal([]byte("北京"), xbytes.Unquote([]byte("'北京'")))
	c.Equal([]byte("🚀🎉"), xbytes.Unquote([]byte("\"🚀🎉\"")))

	c.Equal([]byte("hello\nworld"), xbytes.Unquote([]byte("\"hello\nworld\"")))
	c.Equal([]byte("hello\tworld"), xbytes.Unquote([]byte("'hello\tworld'")))
	c.Equal([]byte("hello\\world"), xbytes.Unquote([]byte("\"hello\\world\"")))

	c.Equal([]byte(" hello "), xbytes.Unquote([]byte("\" hello \"")))
	c.Equal([]byte("\thello\t"), xbytes.Unquote([]byte("'\thello\t'")))

	c.Equal([]byte("123"), xbytes.Unquote([]byte("\"123\"")))
	c.Equal([]byte("3.14"), xbytes.Unquote([]byte("'3.14'")))

	c.Equal([]byte("{\"key\":\"value\"}"), xbytes.Unquote([]byte("\"{\"key\":\"value\"}\"")))

	longContent := make([]byte, 1000)
	for i := range longContent {
		longContent[i] = byte('a' + (i % 26))
	}
	quotedLong := make([]byte, 0, len(longContent)+2)
	quotedLong = append(quotedLong, '"')
	quotedLong = append(quotedLong, longContent...)
	quotedLong = append(quotedLong, '"')
	c.Equal(longContent, xbytes.Unquote(quotedLong))

	original := []byte("\"hello\"")
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	result := xbytes.Unquote(original)

	c.Equal(originalCopy, original)
	c.Equal([]byte("hello"), result)

	var nilSlice []byte
	c.Equal(nilSlice, xbytes.Unquote(nilSlice))

	noQuotes := []byte("hello")
	resultNoQuotes := xbytes.Unquote(noQuotes)
	c.Equal(noQuotes, resultNoQuotes)

	unicodeQuotes := []byte("\u201chello\u201d") // Unicode left/right double quotes
	c.Equal(unicodeQuotes, xbytes.Unquote(unicodeQuotes))

	leftQuote := []byte("\u2018hello\u2019") // Unicode left/right single quotes
	c.Equal(leftQuote, xbytes.Unquote(leftQuote))

	invalidUTF8 := []byte{0xFF, 'h', 'e', 'l', 'l', 'o', 0xFF}
	c.Equal(invalidUTF8, xbytes.Unquote(invalidUTF8))

	backticks := []byte("`hello`")
	c.Equal(backticks, xbytes.Unquote(backticks))

	guillemets := []byte("«hello»")
	c.Equal(guillemets, xbytes.Unquote(guillemets))
}
