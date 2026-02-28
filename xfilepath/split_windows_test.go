// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xfilepath_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xfilepath"
)

func TestWindowsSplit(t *testing.T) {
	data := []splitData{
		{
			in:  `C:\one/two.txt`,
			out: []string{`C:\`, "one", "two.txt"},
		},
		{
			in:  `\\host\share\one\two.txt`,
			out: []string{`\\host\share\`, "one", "two.txt"},
		},
	}
	c := check.New(t)
	for i, one := range data {
		c.Equal(one.out, xfilepath.Split(one.in), "%d", i)
	}
}
