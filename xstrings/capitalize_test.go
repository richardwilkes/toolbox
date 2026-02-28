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
