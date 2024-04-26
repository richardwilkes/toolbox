// Copyright ©2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package slice

import "slices"

// ColumnSort sorts the slice in place using the provided comparison function. The resulting order will be as if the
// slice was divided into columns and each column was sorted independently. If the slice is not evenly divisible by
// the number of columns, the extra elements will be distributed across the columns from left to right.
func ColumnSort[S ~[]E, E any](s S, columns int, cmp func(a, b E) int) {
	slices.SortFunc[S, E](s, cmp)
	if columns > 1 && len(s) > columns {
		replacement := make([]E, len(s))
		step := len(s) / columns
		extra := len(s) - step*columns
		i := 0
		j := 0
		k := 1
		for i < len(s) {
			for c := 0; c < columns; c++ {
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
