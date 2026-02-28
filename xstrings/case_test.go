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

func TestToCamelCase(t *testing.T) {
	c := check.New(t)

	// Test converting snake_case to CamelCase
	c.Equal("SnakeCase", xstrings.ToCamelCase("snake_case"))

	// Test handling multiple consecutive underscores
	c.Equal("SnakeCase", xstrings.ToCamelCase("snake__case"))

	// Test that already CamelCase strings are unchanged
	c.Equal("CamelCase", xstrings.ToCamelCase("CamelCase"))
}

func TestToCamelCaseWithExceptions(t *testing.T) {
	c := check.New(t)

	// Test single exception word converts to all caps
	c.Equal("ID", xstrings.ToCamelCaseWithExceptions("id", xstrings.StdAllCaps))

	// Test unicode characters with exception word
	c.Equal("世界ID", xstrings.ToCamelCaseWithExceptions("世界_id", xstrings.StdAllCaps))

	// Test exception word at the end of compound word
	c.Equal("OneID", xstrings.ToCamelCaseWithExceptions("one_id", xstrings.StdAllCaps))

	// Test exception word at the beginning of compound word
	c.Equal("IDOne", xstrings.ToCamelCaseWithExceptions("id_one", xstrings.StdAllCaps))

	// Test exception word in the middle of compound word
	c.Equal("OneIDTwo", xstrings.ToCamelCaseWithExceptions("one_id_two", xstrings.StdAllCaps))

	// Test multiple exception words in compound word
	c.Equal("OneIDTwoID", xstrings.ToCamelCaseWithExceptions("one_id_two_id", xstrings.StdAllCaps))

	// Test consecutive exception words
	c.Equal("OneIDID", xstrings.ToCamelCaseWithExceptions("one_id_id", xstrings.StdAllCaps))

	// Test word containing exception word but not matching (partial match)
	c.Equal("Orchid", xstrings.ToCamelCaseWithExceptions("orchid", xstrings.StdAllCaps))

	// Test different exception word (URL) in compound word
	c.Equal("OneURLTwo", xstrings.ToCamelCaseWithExceptions("one_url_two", xstrings.StdAllCaps))

	// Test consecutive different exception words
	c.Equal("URLID", xstrings.ToCamelCaseWithExceptions("url_id", xstrings.StdAllCaps))
}

func TestToSnakeCase(t *testing.T) {
	c := check.New(t)

	// Test that already snake_case strings are unchanged
	c.Equal("snake_case", xstrings.ToSnakeCase("snake_case"))

	// Test converting CamelCase to snake_case
	c.Equal("camel_case", xstrings.ToSnakeCase("CamelCase"))
}

func TestFirstToUpper(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", xstrings.FirstToUpper(""))

	// Test single lowercase character
	c.Equal("A", xstrings.FirstToUpper("a"))
	c.Equal("Z", xstrings.FirstToUpper("z"))

	// Test single uppercase character (should remain unchanged)
	c.Equal("A", xstrings.FirstToUpper("A"))
	c.Equal("Z", xstrings.FirstToUpper("Z"))

	// Test lowercase words
	c.Equal("Hello", xstrings.FirstToUpper("hello"))
	c.Equal("World", xstrings.FirstToUpper("world"))

	// Test uppercase words (should remain unchanged)
	c.Equal("Hello", xstrings.FirstToUpper("Hello"))
	c.Equal("HELLO", xstrings.FirstToUpper("HELLO"))

	// Test mixed case words
	c.Equal("HeLLo", xstrings.FirstToUpper("heLLo"))
	c.Equal("WoRLd", xstrings.FirstToUpper("woRLd"))

	// Test numbers (should remain unchanged)
	c.Equal("123", xstrings.FirstToUpper("123"))
	c.Equal("456abc", xstrings.FirstToUpper("456abc"))

	// Test special characters (should remain unchanged)
	c.Equal("!hello", xstrings.FirstToUpper("!hello"))
	c.Equal("@world", xstrings.FirstToUpper("@world"))
	c.Equal("#test", xstrings.FirstToUpper("#test"))

	// Test unicode characters
	c.Equal("Ñoño", xstrings.FirstToUpper("ñoño"))
	c.Equal("Über", xstrings.FirstToUpper("über"))
	c.Equal("Café", xstrings.FirstToUpper("café"))

	// Test already uppercase unicode
	c.Equal("Ñoño", xstrings.FirstToUpper("Ñoño"))
	c.Equal("Über", xstrings.FirstToUpper("Über"))

	// Test words starting with numbers or symbols
	c.Equal("123test", xstrings.FirstToUpper("123test"))
	c.Equal("_hello", xstrings.FirstToUpper("_hello"))
	c.Equal("-world", xstrings.FirstToUpper("-world"))

	// Test single character strings
	c.Equal("X", xstrings.FirstToUpper("x"))
	c.Equal("!", xstrings.FirstToUpper("!"))
	c.Equal("1", xstrings.FirstToUpper("1"))
}

func TestFirstToLower(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", xstrings.FirstToLower(""))

	// Test single uppercase character
	c.Equal("a", xstrings.FirstToLower("A"))
	c.Equal("z", xstrings.FirstToLower("Z"))

	// Test single lowercase character (should remain unchanged)
	c.Equal("a", xstrings.FirstToLower("a"))
	c.Equal("z", xstrings.FirstToLower("z"))

	// Test uppercase words
	c.Equal("hello", xstrings.FirstToLower("Hello"))
	c.Equal("world", xstrings.FirstToLower("World"))

	// Test lowercase words (should remain unchanged)
	c.Equal("hello", xstrings.FirstToLower("hello"))
	c.Equal("world", xstrings.FirstToLower("world"))

	// Test all uppercase words
	c.Equal("hELLO", xstrings.FirstToLower("HELLO"))
	c.Equal("wORLD", xstrings.FirstToLower("WORLD"))

	// Test mixed case words
	c.Equal("heLLo", xstrings.FirstToLower("HeLLo"))
	c.Equal("woRLd", xstrings.FirstToLower("WoRLd"))

	// Test numbers (should remain unchanged)
	c.Equal("123", xstrings.FirstToLower("123"))
	c.Equal("456ABC", xstrings.FirstToLower("456ABC"))

	// Test special characters (should remain unchanged)
	c.Equal("!Hello", xstrings.FirstToLower("!Hello"))
	c.Equal("@World", xstrings.FirstToLower("@World"))
	c.Equal("#Test", xstrings.FirstToLower("#Test"))

	// Test unicode characters
	c.Equal("ñoño", xstrings.FirstToLower("Ñoño"))
	c.Equal("über", xstrings.FirstToLower("Über"))
	c.Equal("café", xstrings.FirstToLower("Café"))

	// Test already lowercase unicode
	c.Equal("ñoño", xstrings.FirstToLower("ñoño"))
	c.Equal("über", xstrings.FirstToLower("über"))

	// Test words starting with numbers or symbols
	c.Equal("123Test", xstrings.FirstToLower("123Test"))
	c.Equal("_Hello", xstrings.FirstToLower("_Hello"))
	c.Equal("-World", xstrings.FirstToLower("-World"))

	// Test single character strings
	c.Equal("x", xstrings.FirstToLower("X"))
	c.Equal("!", xstrings.FirstToLower("!"))
	c.Equal("1", xstrings.FirstToLower("1"))
}
