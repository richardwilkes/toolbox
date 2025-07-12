package xstrings_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestNormalizeLineEndings(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", xstrings.NormalizeLineEndings(""))

	// Test string with no line endings
	c.Equal("hello world", xstrings.NormalizeLineEndings("hello world"))

	// Test string with only LF (Unix) - should remain unchanged
	c.Equal("line1\nline2", xstrings.NormalizeLineEndings("line1\nline2"))
	c.Equal("line1\nline2\nline3", xstrings.NormalizeLineEndings("line1\nline2\nline3"))

	// Test string with CRLF (Windows) - should convert to LF
	c.Equal("line1\nline2", xstrings.NormalizeLineEndings("line1\r\nline2"))
	c.Equal("line1\nline2\nline3", xstrings.NormalizeLineEndings("line1\r\nline2\r\nline3"))

	// Test string with only CR (old Mac) - should convert to LF
	c.Equal("line1\nline2", xstrings.NormalizeLineEndings("line1\rline2"))
	c.Equal("line1\nline2\nline3", xstrings.NormalizeLineEndings("line1\rline2\rline3"))

	// Test mixed line endings - should all convert to LF
	c.Equal("line1\nline2\nline3\nline4", xstrings.NormalizeLineEndings("line1\r\nline2\rline3\nline4"))
	c.Equal("line1\nline2\nline3\nline4\nline5", xstrings.NormalizeLineEndings("line1\nline2\r\nline3\rline4\nline5"))

	// Test multiple consecutive line endings
	c.Equal("line1\n\nline2", xstrings.NormalizeLineEndings("line1\r\n\r\nline2"))
	c.Equal("line1\n\nline2", xstrings.NormalizeLineEndings("line1\r\rline2"))
	c.Equal("line1\n\n\nline2", xstrings.NormalizeLineEndings("line1\r\n\r\n\r\nline2"))

	// Test string starting with line endings
	c.Equal("\nline1", xstrings.NormalizeLineEndings("\r\nline1"))
	c.Equal("\nline1", xstrings.NormalizeLineEndings("\rline1"))
	c.Equal("\nline1", xstrings.NormalizeLineEndings("\nline1"))

	// Test string ending with line endings
	c.Equal("line1\n", xstrings.NormalizeLineEndings("line1\r\n"))
	c.Equal("line1\n", xstrings.NormalizeLineEndings("line1\r"))
	c.Equal("line1\n", xstrings.NormalizeLineEndings("line1\n"))

	// Test string with only line endings
	c.Equal("\n", xstrings.NormalizeLineEndings("\r\n"))
	c.Equal("\n", xstrings.NormalizeLineEndings("\r"))
	c.Equal("\n", xstrings.NormalizeLineEndings("\n"))
	c.Equal("\n\n", xstrings.NormalizeLineEndings("\r\n\r\n"))
	c.Equal("\n\n", xstrings.NormalizeLineEndings("\r\r"))
	c.Equal("\n\n", xstrings.NormalizeLineEndings("\n\n"))

	// Test complex mixed scenarios
	c.Equal("\nstart\nmiddle\nend\n", xstrings.NormalizeLineEndings("\r\nstart\rmiddle\r\nend\n"))
	c.Equal("a\nb\nc\nd\ne", xstrings.NormalizeLineEndings("a\r\nb\rc\nd\r\ne"))

	// Test text with special characters and line endings
	c.Equal("Special: @#$%\nMore text", xstrings.NormalizeLineEndings("Special: @#$%\r\nMore text"))
	c.Equal("Unicode: 🚀\nEmoji text", xstrings.NormalizeLineEndings("Unicode: 🚀\rEmoji text"))

	// Test very long strings with line endings
	longText := "This is a very long line of text that continues for a while to test how the function handles longer content"
	c.Equal(longText+"\nSecond line", xstrings.NormalizeLineEndings(longText+"\r\nSecond line"))
	c.Equal(longText+"\nSecond line", xstrings.NormalizeLineEndings(longText+"\rSecond line"))

	// Test multiple different line ending patterns in sequence
	c.Equal("a\n\n\nb", xstrings.NormalizeLineEndings("a\r\n\r\rb"))
	c.Equal("a\n\n\nb", xstrings.NormalizeLineEndings("a\r\r\n\rb"))
	c.Equal("a\n\n\nb", xstrings.NormalizeLineEndings("a\n\r\n\rb"))

	// Test text with numbers and line endings
	c.Equal("123\n456\n789", xstrings.NormalizeLineEndings("123\r\n456\r789"))
	c.Equal("1.0\n2.5\n3.14", xstrings.NormalizeLineEndings("1.0\r2.5\r\n3.14"))

	// Test whitespace around line endings
	c.Equal("line1 \n line2", xstrings.NormalizeLineEndings("line1 \r\n line2"))
	c.Equal(" line1\nline2 ", xstrings.NormalizeLineEndings(" line1\rline2 "))
	c.Equal("line1\t\nline2", xstrings.NormalizeLineEndings("line1\t\r\nline2"))

	// Test tabs and other whitespace with line endings
	c.Equal("col1\tcol2\nrow2", xstrings.NormalizeLineEndings("col1\tcol2\r\nrow2"))
	c.Equal("  spaced  \n  text  ", xstrings.NormalizeLineEndings("  spaced  \r  text  "))
}
