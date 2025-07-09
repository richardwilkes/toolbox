package txt

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestUnquoteBytes(t *testing.T) {
	c := check.New(t)

	// Test with double quotes
	c.Equal([]byte("hello"), UnquoteBytes([]byte(`"hello"`)))
	c.Equal([]byte("world"), UnquoteBytes([]byte(`"world"`)))
	c.Equal([]byte(""), UnquoteBytes([]byte(`""`)))

	// Test with single quotes
	c.Equal([]byte("hello"), UnquoteBytes([]byte("'hello'")))
	c.Equal([]byte("world"), UnquoteBytes([]byte("'world'")))
	c.Equal([]byte(""), UnquoteBytes([]byte("''")))

	// Test with mixed quotes (should not unquote)
	c.Equal([]byte("'hello\""), UnquoteBytes([]byte("'hello\"")))
	c.Equal([]byte("\"hello'"), UnquoteBytes([]byte("\"hello'")))

	// Test with no quotes
	c.Equal([]byte("hello"), UnquoteBytes([]byte("hello")))
	c.Equal([]byte("world"), UnquoteBytes([]byte("world")))

	// Test with single character
	c.Equal([]byte("a"), UnquoteBytes([]byte("a")))
	c.Equal([]byte("\""), UnquoteBytes([]byte("\"")))
	c.Equal([]byte("'"), UnquoteBytes([]byte("'")))

	// Test with empty bytes
	c.Equal([]byte(""), UnquoteBytes([]byte("")))

	// Test with only opening quote
	c.Equal([]byte("\"hello"), UnquoteBytes([]byte("\"hello")))
	c.Equal([]byte("'hello"), UnquoteBytes([]byte("'hello")))

	// Test with only closing quote
	c.Equal([]byte("hello\""), UnquoteBytes([]byte("hello\"")))
	c.Equal([]byte("hello'"), UnquoteBytes([]byte("hello'")))

	// Test with quotes in the middle
	c.Equal([]byte("he\"llo"), UnquoteBytes([]byte("he\"llo")))
	c.Equal([]byte("he'llo"), UnquoteBytes([]byte("he'llo")))

	// Test with escaped quotes (should still unquote outer quotes)
	c.Equal([]byte("he\\\"llo"), UnquoteBytes([]byte("\"he\\\"llo\"")))
	c.Equal([]byte("he\\'llo"), UnquoteBytes([]byte("'he\\'llo'")))

	// Test with nested quotes
	c.Equal([]byte("'hello'"), UnquoteBytes([]byte("\"'hello'\"")))
	c.Equal([]byte("\"hello\""), UnquoteBytes([]byte("'\"hello\"'")))

	// Test with multiple quotes at start/end
	c.Equal([]byte("\"hello\""), UnquoteBytes([]byte("\"\"hello\"\"")))
	c.Equal([]byte("'hello'"), UnquoteBytes([]byte("''hello''")))

	// Test with Unicode content
	c.Equal([]byte("café"), UnquoteBytes([]byte("\"café\"")))
	c.Equal([]byte("北京"), UnquoteBytes([]byte("'北京'")))
	c.Equal([]byte("🚀🎉"), UnquoteBytes([]byte("\"🚀🎉\"")))

	// Test with special characters
	c.Equal([]byte("hello\nworld"), UnquoteBytes([]byte("\"hello\nworld\"")))
	c.Equal([]byte("hello\tworld"), UnquoteBytes([]byte("'hello\tworld'")))
	c.Equal([]byte("hello\\world"), UnquoteBytes([]byte("\"hello\\world\"")))

	// Test with whitespace
	c.Equal([]byte(" hello "), UnquoteBytes([]byte("\" hello \"")))
	c.Equal([]byte("\thello\t"), UnquoteBytes([]byte("'\thello\t'")))

	// Test with numeric content
	c.Equal([]byte("123"), UnquoteBytes([]byte("\"123\"")))
	c.Equal([]byte("3.14"), UnquoteBytes([]byte("'3.14'")))

	// Test with JSON-like content
	c.Equal([]byte("{\"key\":\"value\"}"), UnquoteBytes([]byte("\"{\"key\":\"value\"}\"")))

	// Test with very long content
	longContent := make([]byte, 1000)
	for i := range longContent {
		longContent[i] = byte('a' + (i % 26))
	}
	quotedLong := append([]byte("\""), append(longContent, '"')...)
	c.Equal(longContent, UnquoteBytes(quotedLong))
}

func TestUnquote(t *testing.T) {
	c := check.New(t)

	// Test with double quotes
	c.Equal("hello", Unquote(`"hello"`))
	c.Equal("world", Unquote(`"world"`))
	c.Equal("", Unquote(`""`))

	// Test with single quotes
	c.Equal("hello", Unquote("'hello'"))
	c.Equal("world", Unquote("'world'"))
	c.Equal("", Unquote("''"))

	// Test with mixed quotes (should not unquote)
	c.Equal("'hello\"", Unquote("'hello\""))
	c.Equal("\"hello'", Unquote("\"hello'"))

	// Test with no quotes
	c.Equal("hello", Unquote("hello"))
	c.Equal("world", Unquote("world"))

	// Test with single character
	c.Equal("a", Unquote("a"))
	c.Equal("\"", Unquote("\""))
	c.Equal("'", Unquote("'"))

	// Test with empty string
	c.Equal("", Unquote(""))

	// Test with only opening quote
	c.Equal("\"hello", Unquote("\"hello"))
	c.Equal("'hello", Unquote("'hello"))

	// Test with only closing quote
	c.Equal("hello\"", Unquote("hello\""))
	c.Equal("hello'", Unquote("hello'"))

	// Test with quotes in the middle
	c.Equal("he\"llo", Unquote("he\"llo"))
	c.Equal("he'llo", Unquote("he'llo"))

	// Test with escaped quotes (should still unquote outer quotes)
	c.Equal("he\\\"llo", Unquote("\"he\\\"llo\""))
	c.Equal("he\\'llo", Unquote("'he\\'llo'"))

	// Test with nested quotes
	c.Equal("'hello'", Unquote("\"'hello'\""))
	c.Equal("\"hello\"", Unquote("'\"hello\"'"))

	// Test with multiple quotes at start/end
	c.Equal("\"hello\"", Unquote("\"\"hello\"\""))
	c.Equal("'hello'", Unquote("''hello''"))

	// Test with Unicode content
	c.Equal("café", Unquote("\"café\""))
	c.Equal("北京", Unquote("'北京'"))
	c.Equal("🚀🎉", Unquote("\"🚀🎉\""))

	// Test with special characters
	c.Equal("hello\nworld", Unquote("\"hello\nworld\""))
	c.Equal("hello\tworld", Unquote("'hello\tworld'"))
	c.Equal("hello\\world", Unquote("\"hello\\world\""))

	// Test with whitespace
	c.Equal(" hello ", Unquote("\" hello \""))
	c.Equal("\thello\t", Unquote("'\thello\t'"))

	// Test with numeric content
	c.Equal("123", Unquote("\"123\""))
	c.Equal("3.14", Unquote("'3.14'"))

	// Test with JSON-like content
	c.Equal("{\"key\":\"value\"}", Unquote("\"{\"key\":\"value\"}\""))

	// Test with very long content
	longContent := ""
	for i := 0; i < 1000; i++ {
		longContent += string(rune('a' + (i % 26)))
	}
	quotedLong := "\"" + longContent + "\""
	c.Equal(longContent, Unquote(quotedLong))
}

func TestUnquoteBytesEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test that original slice is not modified
	original := []byte("\"hello\"")
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	result := UnquoteBytes(original)

	// Original should be unchanged
	c.Equal(originalCopy, original)
	// Result should be unquoted
	c.Equal([]byte("hello"), result)

	// Test with nil slice
	var nilSlice []byte
	c.Equal(nilSlice, UnquoteBytes(nilSlice))

	// Test with slice that shares memory with original when no unquoting
	noQuotes := []byte("hello")
	resultNoQuotes := UnquoteBytes(noQuotes)
	// Should return the same underlying data when no quotes to remove
	c.Equal(noQuotes, resultNoQuotes)

	// Test Unicode quote characters (should not be processed)
	unicodeQuotes := []byte("\u201chello\u201d") // Unicode left/right double quotes
	c.Equal(unicodeQuotes, UnquoteBytes(unicodeQuotes))

	leftQuote := []byte("\u2018hello\u2019") // Unicode left/right single quotes
	c.Equal(leftQuote, UnquoteBytes(leftQuote))

	// Test with byte sequences that look like quotes but aren't valid UTF-8
	invalidUTF8 := []byte{0xFF, 'h', 'e', 'l', 'l', 'o', 0xFF}
	c.Equal(invalidUTF8, UnquoteBytes(invalidUTF8))

	// Test with various quote-like characters
	backticks := []byte("`hello`")
	c.Equal(backticks, UnquoteBytes(backticks))

	guillemets := []byte("«hello»")
	c.Equal(guillemets, UnquoteBytes(guillemets))
}

func TestUnquoteEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test Unicode quote characters (should not be processed)
	unicodeQuotes := "\u201chello\u201d" // Unicode left/right double quotes
	c.Equal(unicodeQuotes, Unquote(unicodeQuotes))

	leftQuote := "\u2018hello\u2019" // Unicode left/right single quotes
	c.Equal(leftQuote, Unquote(leftQuote))

	// Test with various quote-like characters
	backticks := "`hello`"
	c.Equal(backticks, Unquote(backticks))

	guillemets := "«hello»"
	c.Equal(guillemets, Unquote(guillemets))

	// Test with mathematical symbols that might look like quotes
	primes := "′hello′"
	c.Equal(primes, Unquote(primes))

	doublePrimes := "″hello″"
	c.Equal(doublePrimes, Unquote(doublePrimes))

	// Test with control characters
	withNull := "\x00hello\x00"
	c.Equal(withNull, Unquote(withNull))

	// Test with very short strings
	oneChar := "a"
	c.Equal(oneChar, Unquote(oneChar))

	twoChars := "ab"
	c.Equal(twoChars, Unquote(twoChars))
}

func TestUnquoteVsUnquoteBytes(t *testing.T) {
	c := check.New(t)

	// Test that string and byte versions produce equivalent results
	testCases := []string{
		`"hello"`,
		`'world'`,
		`"café"`,
		`'🚀'`,
		`"hello\nworld"`,
		`'mixed"quotes'`,
		`"single'quote"`,
		`hello`,
		`"`,
		`'`,
		`""`,
		`''`,
		`"a"`,
		`'b'`,
	}

	for _, testCase := range testCases {
		stringResult := Unquote(testCase)
		bytesResult := string(UnquoteBytes([]byte(testCase)))
		c.Equal(stringResult, bytesResult)
	}
}
