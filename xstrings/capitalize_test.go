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

	c.Equal("", xstrings.CapitalizeWords(""))

	c.Equal("Hello", xstrings.CapitalizeWords("hello"))

	c.Equal("Hello", xstrings.CapitalizeWords("HELLO"))

	c.Equal("Hello", xstrings.CapitalizeWords("hELLo"))

	c.Equal("Hello World", xstrings.CapitalizeWords("hello world"))

	c.Equal("Hello World", xstrings.CapitalizeWords("HELLO WORLD"))

	c.Equal("Hello World", xstrings.CapitalizeWords("hELLo WoRLd"))

	c.Equal("Hello World", xstrings.CapitalizeWords("hello    world"))

	c.Equal("Hello World", xstrings.CapitalizeWords("  hello world  "))

	c.Equal("A B C", xstrings.CapitalizeWords("a b c"))

	c.Equal("Hello World123 Test", xstrings.CapitalizeWords("hello world123 test"))

	c.Equal("Hello-world Test_case", xstrings.CapitalizeWords("hello-world test_case"))
}
