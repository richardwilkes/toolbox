/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package xlsx_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/formats/xlsx"
)

func TestRef(t *testing.T) {
	for i, d := range []struct {
		Text string
		Col  int
		Row  int
	}{
		{"A1", 0, 0},
		{"Z9", 25, 8},
		{"AA1", 26, 0},
		{"AA99", 26, 98},
		{"ZZ100", 701, 99},
	} {
		ref := xlsx.ParseRef(d.Text)
		check.Equal(t, d.Col, ref.Col, "column for index %d: %s", i, d.Text)
		check.Equal(t, d.Row, ref.Row, "row for index %d: %s", i, d.Text)
		check.Equal(t, d.Text, ref.String(), "String() for index %d: %s", i, d.Text)
	}

	for r := 0; r < 100; r++ {
		for c := 0; c < 10000; c++ {
			in := xlsx.Ref{Row: r, Col: c}
			out := xlsx.ParseRef(in.String())
			check.Equal(t, in, out)
		}
	}
}
