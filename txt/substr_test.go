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

func TestFirstN(t *testing.T) {
	table := []struct {
		In  string
		Out string
		N   int
	}{
		{In: "abcd", N: 3, Out: "abc"},
		{In: "abcd", N: 5, Out: "abcd"},
		{In: "abcd", N: 0, Out: ""},
		{In: "abcd", N: -1, Out: ""},
		{In: "aécd", N: 3, Out: "aéc"},
		{In: "aécd", N: 5, Out: "aécd"},
		{In: "aécd", N: 0, Out: ""},
		{In: "aécd", N: -1, Out: ""},
	}
	c := check.New(t)
	for i, one := range table {
		c.Equal(one.Out, txt.FirstN(one.In, one.N), "#%d", i)
	}
}

func TestLastN(t *testing.T) {
	table := []struct {
		In  string
		Out string
		N   int
	}{
		{In: "abcd", N: 3, Out: "bcd"},
		{In: "abcd", N: 5, Out: "abcd"},
		{In: "abcd", N: 0, Out: ""},
		{In: "abcd", N: -1, Out: ""},
		{In: "aécd", N: 3, Out: "écd"},
		{In: "aécd", N: 5, Out: "aécd"},
		{In: "aécd", N: 0, Out: ""},
		{In: "aécd", N: -1, Out: ""},
	}
	c := check.New(t)
	for i, one := range table {
		c.Equal(one.Out, txt.LastN(one.In, one.N), "#%d", i)
	}
}
