package xstrings_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestCapitalizeWords(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", xstrings.CapitalizeWords(""))

	// Test single word lowercase
	c.Equal("Hello", xstrings.CapitalizeWords("hello"))

	// Test single word uppercase
	c.Equal("Hello", xstrings.CapitalizeWords("HELLO"))

	// Test single word mixed case
	c.Equal("Hello", xstrings.CapitalizeWords("hELLo"))

	// Test multiple words lowercase
	c.Equal("Hello World", xstrings.CapitalizeWords("hello world"))

	// Test multiple words uppercase
	c.Equal("Hello World", xstrings.CapitalizeWords("HELLO WORLD"))

	// Test multiple words mixed case
	c.Equal("Hello World", xstrings.CapitalizeWords("hELLo WoRLd"))

	// Test multiple spaces between words
	c.Equal("Hello World", xstrings.CapitalizeWords("hello    world"))

	// Test leading and trailing spaces
	c.Equal("Hello World", xstrings.CapitalizeWords("  hello world  "))

	// Test single character words
	c.Equal("A B C", xstrings.CapitalizeWords("a b c"))

	// Test words with numbers
	c.Equal("Hello World123 Test", xstrings.CapitalizeWords("hello world123 test"))

	// Test special characters
	c.Equal("Hello-world Test_case", xstrings.CapitalizeWords("hello-world test_case"))
}
