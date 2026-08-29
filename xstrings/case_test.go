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

	c.Equal("SnakeCase", xstrings.ToCamelCase("snake_case"))

	c.Equal("SnakeCase", xstrings.ToCamelCase("snake__case"))

	c.Equal("CamelCase", xstrings.ToCamelCase("CamelCase"))
}

func TestToCamelCaseWithExceptions(t *testing.T) {
	c := check.New(t)

	c.Equal("ID", xstrings.ToCamelCaseWithExceptions("id", xstrings.StdAllCaps))

	c.Equal("世界ID", xstrings.ToCamelCaseWithExceptions("世界_id", xstrings.StdAllCaps))

	c.Equal("OneID", xstrings.ToCamelCaseWithExceptions("one_id", xstrings.StdAllCaps))

	c.Equal("IDOne", xstrings.ToCamelCaseWithExceptions("id_one", xstrings.StdAllCaps))

	c.Equal("OneIDTwo", xstrings.ToCamelCaseWithExceptions("one_id_two", xstrings.StdAllCaps))

	c.Equal("OneIDTwoID", xstrings.ToCamelCaseWithExceptions("one_id_two_id", xstrings.StdAllCaps))

	c.Equal("OneIDID", xstrings.ToCamelCaseWithExceptions("one_id_id", xstrings.StdAllCaps))

	// "id" inside "orchid" must not be treated as an exception
	c.Equal("Orchid", xstrings.ToCamelCaseWithExceptions("orchid", xstrings.StdAllCaps))

	c.Equal("OneURLTwo", xstrings.ToCamelCaseWithExceptions("one_url_two", xstrings.StdAllCaps))

	c.Equal("URLID", xstrings.ToCamelCaseWithExceptions("url_id", xstrings.StdAllCaps))
}

func TestToSnakeCase(t *testing.T) {
	c := check.New(t)

	c.Equal("snake_case", xstrings.ToSnakeCase("snake_case"))

	c.Equal("camel_case", xstrings.ToSnakeCase("CamelCase"))
}

func TestFirstToUpper(t *testing.T) {
	c := check.New(t)

	c.Equal("", xstrings.FirstToUpper(""))

	c.Equal("A", xstrings.FirstToUpper("a"))
	c.Equal("Z", xstrings.FirstToUpper("z"))

	c.Equal("A", xstrings.FirstToUpper("A"))
	c.Equal("Z", xstrings.FirstToUpper("Z"))

	c.Equal("Hello", xstrings.FirstToUpper("hello"))
	c.Equal("World", xstrings.FirstToUpper("world"))

	c.Equal("Hello", xstrings.FirstToUpper("Hello"))
	c.Equal("HELLO", xstrings.FirstToUpper("HELLO"))

	c.Equal("HeLLo", xstrings.FirstToUpper("heLLo"))
	c.Equal("WoRLd", xstrings.FirstToUpper("woRLd"))

	c.Equal("123", xstrings.FirstToUpper("123"))
	c.Equal("456abc", xstrings.FirstToUpper("456abc"))

	c.Equal("!hello", xstrings.FirstToUpper("!hello"))
	c.Equal("@world", xstrings.FirstToUpper("@world"))
	c.Equal("#test", xstrings.FirstToUpper("#test"))

	c.Equal("Ñoño", xstrings.FirstToUpper("ñoño"))
	c.Equal("Über", xstrings.FirstToUpper("über"))
	c.Equal("Café", xstrings.FirstToUpper("café"))

	c.Equal("Ñoño", xstrings.FirstToUpper("Ñoño"))
	c.Equal("Über", xstrings.FirstToUpper("Über"))

	c.Equal("123test", xstrings.FirstToUpper("123test"))
	c.Equal("_hello", xstrings.FirstToUpper("_hello"))
	c.Equal("-world", xstrings.FirstToUpper("-world"))

	c.Equal("X", xstrings.FirstToUpper("x"))
	c.Equal("!", xstrings.FirstToUpper("!"))
	c.Equal("1", xstrings.FirstToUpper("1"))
}

func TestFirstToLower(t *testing.T) {
	c := check.New(t)

	c.Equal("", xstrings.FirstToLower(""))

	c.Equal("a", xstrings.FirstToLower("A"))
	c.Equal("z", xstrings.FirstToLower("Z"))

	c.Equal("a", xstrings.FirstToLower("a"))
	c.Equal("z", xstrings.FirstToLower("z"))

	c.Equal("hello", xstrings.FirstToLower("Hello"))
	c.Equal("world", xstrings.FirstToLower("World"))

	c.Equal("hello", xstrings.FirstToLower("hello"))
	c.Equal("world", xstrings.FirstToLower("world"))

	c.Equal("hELLO", xstrings.FirstToLower("HELLO"))
	c.Equal("wORLD", xstrings.FirstToLower("WORLD"))

	c.Equal("heLLo", xstrings.FirstToLower("HeLLo"))
	c.Equal("woRLd", xstrings.FirstToLower("WoRLd"))

	c.Equal("123", xstrings.FirstToLower("123"))
	c.Equal("456ABC", xstrings.FirstToLower("456ABC"))

	c.Equal("!Hello", xstrings.FirstToLower("!Hello"))
	c.Equal("@World", xstrings.FirstToLower("@World"))
	c.Equal("#Test", xstrings.FirstToLower("#Test"))

	c.Equal("ñoño", xstrings.FirstToLower("Ñoño"))
	c.Equal("über", xstrings.FirstToLower("Über"))
	c.Equal("café", xstrings.FirstToLower("Café"))

	c.Equal("ñoño", xstrings.FirstToLower("ñoño"))
	c.Equal("über", xstrings.FirstToLower("über"))

	c.Equal("123Test", xstrings.FirstToLower("123Test"))
	c.Equal("_Hello", xstrings.FirstToLower("_Hello"))
	c.Equal("-World", xstrings.FirstToLower("-World"))

	c.Equal("x", xstrings.FirstToLower("X"))
	c.Equal("!", xstrings.FirstToLower("!"))
	c.Equal("1", xstrings.FirstToLower("1"))
}
