// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package slice_test

import (
	"cmp"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/collection/slice"
)

func TestColumnSort(t *testing.T) {
	c := check.New(t)
	s := []int{0, 1, 2, 3, 4, 5, 6}
	slice.ColumnSort(s, 2, cmp.Compare)
	// 0 4
	// 1 5
	// 2 6
	// 3
	c.Equal([]int{0, 4, 1, 5, 2, 6, 3}, s)

	slice.ColumnSort(s, 3, cmp.Compare)
	// 0 3 5
	// 1 4 6
	// 2
	c.Equal([]int{0, 3, 5, 1, 4, 6, 2}, s)

	s = []int{0, 1, 2, 3, 4, 5}
	slice.ColumnSort(s, 2, cmp.Compare)
	// 0 3
	// 1 4
	// 2 5
	c.Equal([]int{0, 3, 1, 4, 2, 5}, s)

	slice.ColumnSort(s, 4, cmp.Compare)
	// 0 2 4 5
	// 1 3
	c.Equal([]int{0, 2, 4, 5, 1, 3}, s)

	slice.ColumnSort(s, 10, cmp.Compare)
	// 0 1 2 3 4 5
	c.Equal([]int{0, 1, 2, 3, 4, 5}, s)

	s = []int{0, 1, 2, 3, 4, 5, 6, 7}
	slice.ColumnSort(s, 3, cmp.Compare)
	// 0 3 6
	// 1 4 7
	// 2 5
	c.Equal([]int{0, 3, 6, 1, 4, 7, 2, 5}, s)
}
