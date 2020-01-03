// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	for i, one := range []struct {
		in  string
		out []string
	}{
		{
			in:  "/one/two.txt",
			out: []string{"/", "one", "two.txt"},
		},
		{
			in:  "/one",
			out: []string{"/", "one"},
		},
		{
			in:  "one",
			out: []string{".", "one"},
		},
		{
			in:  "/one////two.txt",
			out: []string{"/", "one", "two.txt"},
		},
		{
			in:  "/one//..//two.txt",
			out: []string{"/", "two.txt"},
		},
		{
			in:  "/one/../..//two.txt",
			out: []string{"/", "two.txt"},
		},
		{
			in:  "/one/../..//two.txt/",
			out: []string{"/", "two.txt"},
		},
		{
			in:  "/one/../..//two.txt/.",
			out: []string{"/", "two.txt"},
		},
	} {
		assert.Equal(t, one.out, fs.Split(one.in), "%d", i)
	}
}
