// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/stretchr/testify/assert"
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
	if runtime.GOOS == toolbox.WindowsOS {
		data = append(data,
			splitData{
				in:  `C:\one/two.txt`,
				out: []string{`C:\`, "one", "two.txt"},
			},
			splitData{
				in:  `\\host\share\one\two.txt`,
				out: []string{`\\host\share\`, "one", "two.txt"},
			},
		)
	}
	for i, one := range data {
		assert.Equal(t, one.out, fs.Split(one.in), "%d", i)
	}
}
