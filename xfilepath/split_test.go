// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xfilepath_test

import (
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xfilepath"
)

type splitData struct {
	in  string
	out []string
}

func TestSplit(t *testing.T) {
	full := string([]rune{filepath.Separator})
	data := []splitData{
		{
			in:  "/one/two.txt",
			out: []string{full, "one", "two.txt"},
		},
		{
			in:  "/one",
			out: []string{full, "one"},
		},
		{
			in:  "one",
			out: []string{".", "one"},
		},
		{
			in:  "/one////two.txt",
			out: []string{full, "one", "two.txt"},
		},
		{
			in:  "/one//..//two.txt",
			out: []string{full, "two.txt"},
		},
		{
			in:  "/one/../..//two.txt",
			out: []string{full, "two.txt"},
		},
		{
			in:  "/one/../..//two.txt/",
			out: []string{full, "two.txt"},
		},
		{
			in:  "/one/../..//two.txt/.",
			out: []string{full, "two.txt"},
		},
	}
	c := check.New(t)
	for i, one := range data {
		c.Equal(one.out, xfilepath.Split(one.in), "%d", i)
	}
}
