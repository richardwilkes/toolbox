package xbytes_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xbytes"
)

func TestUnquoteBytes(t *testing.T) {
	c := check.New(t)

	// Test with double quotes
	c.Equal([]byte("hello"), xbytes.Unquote([]byte(`"hello"`)))
	c.Equal([]byte("world"), xbytes.Unquote([]byte(`"world"`)))
	c.Equal([]byte(""), xbytes.Unquote([]byte(`""`)))

	// Test with single quotes
	c.Equal([]byte("hello"), xbytes.Unquote([]byte("'hello'")))
	c.Equal([]byte("world"), xbytes.Unquote([]byte("'world'")))
	c.Equal([]byte(""), xbytes.Unquote([]byte("''")))

	// Test with mixed quotes (should not unquote)
	c.Equal([]byte("'hello\""), xbytes.Unquote([]byte("'hello\"")))
	c.Equal([]byte("\"hello'"), xbytes.Unquote([]byte("\"hello'")))

	// Test with no quotes
	c.Equal([]byte("hello"), xbytes.Unquote([]byte("hello")))
	c.Equal([]byte("world"), xbytes.Unquote([]byte("world")))

	// Test with single character
	c.Equal([]byte("a"), xbytes.Unquote([]byte("a")))
	c.Equal([]byte("\""), xbytes.Unquote([]byte("\"")))
	c.Equal([]byte("'"), xbytes.Unquote([]byte("'")))

	// Test with empty bytes
	c.Equal([]byte(""), xbytes.Unquote([]byte("")))

	// Test with only opening quote
	c.Equal([]byte("\"hello"), xbytes.Unquote([]byte("\"hello")))
	c.Equal([]byte("'hello"), xbytes.Unquote([]byte("'hello")))

	// Test with only closing quote
	c.Equal([]byte("hello\""), xbytes.Unquote([]byte("hello\"")))
	c.Equal([]byte("hello'"), xbytes.Unquote([]byte("hello'")))

	// Test with quotes in the middle
	c.Equal([]byte("he\"llo"), xbytes.Unquote([]byte("he\"llo")))
	c.Equal([]byte("he'llo"), xbytes.Unquote([]byte("he'llo")))

	// Test with escaped quotes (should still unquote outer quotes)
	c.Equal([]byte("he\\\"llo"), xbytes.Unquote([]byte("\"he\\\"llo\"")))
	c.Equal([]byte("he\\'llo"), xbytes.Unquote([]byte("'he\\'llo'")))

	// Test with nested quotes
	c.Equal([]byte("'hello'"), xbytes.Unquote([]byte("\"'hello'\"")))
	c.Equal([]byte("\"hello\""), xbytes.Unquote([]byte("'\"hello\"'")))

	// Test with multiple quotes at start/end
	c.Equal([]byte("\"hello\""), xbytes.Unquote([]byte("\"\"hello\"\"")))
	c.Equal([]byte("'hello'"), xbytes.Unquote([]byte("''hello''")))

	// Test with Unicode content
	c.Equal([]byte("café"), xbytes.Unquote([]byte("\"café\"")))
	c.Equal([]byte("北京"), xbytes.Unquote([]byte("'北京'")))
	c.Equal([]byte("🚀🎉"), xbytes.Unquote([]byte("\"🚀🎉\"")))

	// Test with special characters
	c.Equal([]byte("hello\nworld"), xbytes.Unquote([]byte("\"hello\nworld\"")))
	c.Equal([]byte("hello\tworld"), xbytes.Unquote([]byte("'hello\tworld'")))
	c.Equal([]byte("hello\\world"), xbytes.Unquote([]byte("\"hello\\world\"")))

	// Test with whitespace
	c.Equal([]byte(" hello "), xbytes.Unquote([]byte("\" hello \"")))
	c.Equal([]byte("\thello\t"), xbytes.Unquote([]byte("'\thello\t'")))

	// Test with numeric content
	c.Equal([]byte("123"), xbytes.Unquote([]byte("\"123\"")))
	c.Equal([]byte("3.14"), xbytes.Unquote([]byte("'3.14'")))

	// Test with JSON-like content
	c.Equal([]byte("{\"key\":\"value\"}"), xbytes.Unquote([]byte("\"{\"key\":\"value\"}\"")))

	// Test with very long content
	longContent := make([]byte, 1000)
	for i := range longContent {
		longContent[i] = byte('a' + (i % 26))
	}
	quotedLong := append([]byte("\""), append(longContent, '"')...)
	c.Equal(longContent, xbytes.Unquote(quotedLong))

	// Test that original slice is not modified
	original := []byte("\"hello\"")
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	result := xbytes.Unquote(original)

	// Original should be unchanged
	c.Equal(originalCopy, original)
	// Result should be unquoted
	c.Equal([]byte("hello"), result)

	// Test with nil slice
	var nilSlice []byte
	c.Equal(nilSlice, xbytes.Unquote(nilSlice))

	// Test with slice that shares memory with original when no unquoting
	noQuotes := []byte("hello")
	resultNoQuotes := xbytes.Unquote(noQuotes)
	// Should return the same underlying data when no quotes to remove
	c.Equal(noQuotes, resultNoQuotes)

	// Test Unicode quote characters (should not be processed)
	unicodeQuotes := []byte("\u201chello\u201d") // Unicode left/right double quotes
	c.Equal(unicodeQuotes, xbytes.Unquote(unicodeQuotes))

	leftQuote := []byte("\u2018hello\u2019") // Unicode left/right single quotes
	c.Equal(leftQuote, xbytes.Unquote(leftQuote))

	// Test with byte sequences that look like quotes but aren't valid UTF-8
	invalidUTF8 := []byte{0xFF, 'h', 'e', 'l', 'l', 'o', 0xFF}
	c.Equal(invalidUTF8, xbytes.Unquote(invalidUTF8))

	// Test with various quote-like characters
	backticks := []byte("`hello`")
	c.Equal(backticks, xbytes.Unquote(backticks))

	guillemets := []byte("«hello»")
	c.Equal(guillemets, xbytes.Unquote(guillemets))
}
