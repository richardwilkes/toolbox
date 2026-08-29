// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslices

import "slices"

// ColumnSort sorts 's' in place using 'cmp', then rearranges it so that laying the result out row by row into
// 'columns' columns places the sorted order down each column in turn. If the slice is not evenly divisible by the
// number of columns, the extra elements are distributed across the columns from left to right.
func ColumnSort[S ~[]E, E any](s S, columns int, cmp func(a, b E) int) {
	slices.SortFunc(s, cmp)
	if columns > 1 && len(s) > columns {
		replacement := make([]E, len(s))
		step := len(s) / columns
		extra := len(s) - step*columns
		i := 0
		j := 0
		k := 1
		for i < len(s) {
			for c := range columns {
				replacement[i] = s[j]
				i++
				if i >= len(s) {
					break
				}
				j += step
				if extra > c {
					j++
				}
				if j >= len(s) {
					j = k
					k++
				}
			}
		}
		copy(s, replacement)
	}
}
