package txt

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestCapitalizeWords(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", CapitalizeWords(""))

	// Test single word lowercase
	c.Equal("Hello", CapitalizeWords("hello"))

	// Test single word uppercase
	c.Equal("Hello", CapitalizeWords("HELLO"))

	// Test single word mixed case
	c.Equal("Hello", CapitalizeWords("hELLo"))

	// Test multiple words lowercase
	c.Equal("Hello World", CapitalizeWords("hello world"))

	// Test multiple words uppercase
	c.Equal("Hello World", CapitalizeWords("HELLO WORLD"))

	// Test multiple words mixed case
	c.Equal("Hello World", CapitalizeWords("hELLo WoRLd"))

	// Test multiple spaces between words
	c.Equal("Hello World", CapitalizeWords("hello    world"))

	// Test leading and trailing spaces
	c.Equal("Hello World", CapitalizeWords("  hello world  "))

	// Test single character words
	c.Equal("A B C", CapitalizeWords("a b c"))

	// Test words with numbers
	c.Equal("Hello World123 Test", CapitalizeWords("hello world123 test"))

	// Test special characters
	c.Equal("Hello-world Test_case", CapitalizeWords("hello-world test_case"))
}
