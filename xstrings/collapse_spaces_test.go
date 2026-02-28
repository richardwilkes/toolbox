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

func TestCollapseSpaces(t *testing.T) {
	c := check.New(t)

	// Test string with no spaces
	c.Equal("123", xstrings.CollapseSpaces("123"))

	// Test string with leading space
	c.Equal("123", xstrings.CollapseSpaces(" 123"))

	// Test string with leading and trailing spaces
	c.Equal("123", xstrings.CollapseSpaces(" 123 "))

	// Test string with multiple leading and trailing spaces
	c.Equal("abc", xstrings.CollapseSpaces("    abc  "))

	// Test string with multiple spaces between words
	c.Equal("a b c d", xstrings.CollapseSpaces("  a b c   d"))

	// Test empty string
	c.Equal("", xstrings.CollapseSpaces(""))

	// Test string with only spaces
	c.Equal("", xstrings.CollapseSpaces(" "))
}
