// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package slice_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/collection/slice"
)

func TestZeroedDelete(t *testing.T) {
	type data struct {
		a int
	}
	d1 := &data{a: 1}
	d2 := &data{a: 2}
	d3 := &data{a: 3}
	for _, test := range []struct {
		s    []*data
		i, j int
		want []*data
	}{
		{
			[]*data{d1, d2, d3},
			0,
			0,
			[]*data{d1, d2, d3},
		},
		{
			[]*data{d1, d2, d3},
			0,
			1,
			[]*data{d2, d3},
		},
		{
			[]*data{d1, d2, d3},
			3,
			3,
			[]*data{d1, d2, d3},
		},
		{
			[]*data{d1, d2, d3},
			0,
			2,
			[]*data{d3},
		},
		{
			[]*data{d1, d2, d3},
			0,
			3,
			[]*data{},
		},
	} {
		theCopy := append([]*data{}, test.s...)
		result := slice.ZeroedDelete(theCopy, test.i, test.j)
		check.Equal(t, result, test.want, "ZeroedDelete(%v, %d, %d) = %v, want %v", test.s, test.i, test.j, result, test.want)
		for i := len(result); i < len(theCopy); i++ {
			check.Nil(t, theCopy[i], "residual element %d should have been nil, was %v", i, theCopy[i])
		}
	}
}
