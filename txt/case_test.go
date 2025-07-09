// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/txt"
)

func TestToCamelCase(t *testing.T) {
	c := check.New(t)

	// Test converting snake_case to CamelCase
	c.Equal("SnakeCase", txt.ToCamelCase("snake_case"))

	// Test handling multiple consecutive underscores
	c.Equal("SnakeCase", txt.ToCamelCase("snake__case"))

	// Test that already CamelCase strings are unchanged
	c.Equal("CamelCase", txt.ToCamelCase("CamelCase"))
}

func TestToCamelCaseWithExceptions(t *testing.T) {
	c := check.New(t)

	// Test single exception word converts to all caps
	c.Equal("ID", txt.ToCamelCaseWithExceptions("id", txt.StdAllCaps))

	// Test unicode characters with exception word
	c.Equal("世界ID", txt.ToCamelCaseWithExceptions("世界_id", txt.StdAllCaps))

	// Test exception word at the end of compound word
	c.Equal("OneID", txt.ToCamelCaseWithExceptions("one_id", txt.StdAllCaps))

	// Test exception word at the beginning of compound word
	c.Equal("IDOne", txt.ToCamelCaseWithExceptions("id_one", txt.StdAllCaps))

	// Test exception word in the middle of compound word
	c.Equal("OneIDTwo", txt.ToCamelCaseWithExceptions("one_id_two", txt.StdAllCaps))

	// Test multiple exception words in compound word
	c.Equal("OneIDTwoID", txt.ToCamelCaseWithExceptions("one_id_two_id", txt.StdAllCaps))

	// Test consecutive exception words
	c.Equal("OneIDID", txt.ToCamelCaseWithExceptions("one_id_id", txt.StdAllCaps))

	// Test word containing exception word but not matching (partial match)
	c.Equal("Orchid", txt.ToCamelCaseWithExceptions("orchid", txt.StdAllCaps))

	// Test different exception word (URL) in compound word
	c.Equal("OneURLTwo", txt.ToCamelCaseWithExceptions("one_url_two", txt.StdAllCaps))

	// Test consecutive different exception words
	c.Equal("URLID", txt.ToCamelCaseWithExceptions("url_id", txt.StdAllCaps))
}

func TestToSnakeCase(t *testing.T) {
	c := check.New(t)

	// Test that already snake_case strings are unchanged
	c.Equal("snake_case", txt.ToSnakeCase("snake_case"))

	// Test converting CamelCase to snake_case
	c.Equal("camel_case", txt.ToSnakeCase("CamelCase"))
}

func TestFirstToUpper(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", txt.FirstToUpper(""))

	// Test single lowercase character
	c.Equal("A", txt.FirstToUpper("a"))
	c.Equal("Z", txt.FirstToUpper("z"))

	// Test single uppercase character (should remain unchanged)
	c.Equal("A", txt.FirstToUpper("A"))
	c.Equal("Z", txt.FirstToUpper("Z"))

	// Test lowercase words
	c.Equal("Hello", txt.FirstToUpper("hello"))
	c.Equal("World", txt.FirstToUpper("world"))

	// Test uppercase words (should remain unchanged)
	c.Equal("Hello", txt.FirstToUpper("Hello"))
	c.Equal("HELLO", txt.FirstToUpper("HELLO"))

	// Test mixed case words
	c.Equal("HeLLo", txt.FirstToUpper("heLLo"))
	c.Equal("WoRLd", txt.FirstToUpper("woRLd"))

	// Test numbers (should remain unchanged)
	c.Equal("123", txt.FirstToUpper("123"))
	c.Equal("456abc", txt.FirstToUpper("456abc"))

	// Test special characters (should remain unchanged)
	c.Equal("!hello", txt.FirstToUpper("!hello"))
	c.Equal("@world", txt.FirstToUpper("@world"))
	c.Equal("#test", txt.FirstToUpper("#test"))

	// Test unicode characters
	c.Equal("Ñoño", txt.FirstToUpper("ñoño"))
	c.Equal("Über", txt.FirstToUpper("über"))
	c.Equal("Café", txt.FirstToUpper("café"))

	// Test already uppercase unicode
	c.Equal("Ñoño", txt.FirstToUpper("Ñoño"))
	c.Equal("Über", txt.FirstToUpper("Über"))

	// Test words starting with numbers or symbols
	c.Equal("123test", txt.FirstToUpper("123test"))
	c.Equal("_hello", txt.FirstToUpper("_hello"))
	c.Equal("-world", txt.FirstToUpper("-world"))

	// Test single character strings
	c.Equal("X", txt.FirstToUpper("x"))
	c.Equal("!", txt.FirstToUpper("!"))
	c.Equal("1", txt.FirstToUpper("1"))
}

func TestFirstToLower(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", txt.FirstToLower(""))

	// Test single uppercase character
	c.Equal("a", txt.FirstToLower("A"))
	c.Equal("z", txt.FirstToLower("Z"))

	// Test single lowercase character (should remain unchanged)
	c.Equal("a", txt.FirstToLower("a"))
	c.Equal("z", txt.FirstToLower("z"))

	// Test uppercase words
	c.Equal("hello", txt.FirstToLower("Hello"))
	c.Equal("world", txt.FirstToLower("World"))

	// Test lowercase words (should remain unchanged)
	c.Equal("hello", txt.FirstToLower("hello"))
	c.Equal("world", txt.FirstToLower("world"))

	// Test all uppercase words
	c.Equal("hELLO", txt.FirstToLower("HELLO"))
	c.Equal("wORLD", txt.FirstToLower("WORLD"))

	// Test mixed case words
	c.Equal("heLLo", txt.FirstToLower("HeLLo"))
	c.Equal("woRLd", txt.FirstToLower("WoRLd"))

	// Test numbers (should remain unchanged)
	c.Equal("123", txt.FirstToLower("123"))
	c.Equal("456ABC", txt.FirstToLower("456ABC"))

	// Test special characters (should remain unchanged)
	c.Equal("!Hello", txt.FirstToLower("!Hello"))
	c.Equal("@World", txt.FirstToLower("@World"))
	c.Equal("#Test", txt.FirstToLower("#Test"))

	// Test unicode characters
	c.Equal("ñoño", txt.FirstToLower("Ñoño"))
	c.Equal("über", txt.FirstToLower("Über"))
	c.Equal("café", txt.FirstToLower("Café"))

	// Test already lowercase unicode
	c.Equal("ñoño", txt.FirstToLower("ñoño"))
	c.Equal("über", txt.FirstToLower("über"))

	// Test words starting with numbers or symbols
	c.Equal("123Test", txt.FirstToLower("123Test"))
	c.Equal("_Hello", txt.FirstToLower("_Hello"))
	c.Equal("-World", txt.FirstToLower("-World"))

	// Test single character strings
	c.Equal("x", txt.FirstToLower("X"))
	c.Equal("!", txt.FirstToLower("!"))
	c.Equal("1", txt.FirstToLower("1"))
}
