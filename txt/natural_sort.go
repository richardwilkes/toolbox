// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import (
	"slices"
)

// NaturalLess compares two strings using natural ordering. This means that "a2" < "a12".
//
// Non-digit sequences and numbers are compared separately. The former are compared byte-wise, while the latter are
// compared numerically (except that the number of leading zeros is used as a tie-breaker, so "2" < "02").
//
// Limitations:
//   - only ASCII digits (0-9) are considered.
//
// Original algorithm: https://github.com/fvbommel/util/blob/master/sortorder/natsort.go
func NaturalLess(s1, s2 string, caseInsensitive bool) bool {
	return NaturalCmp(s1, s2, caseInsensitive) < 0
}

// NaturalCmp compares two strings using natural ordering. This means that "a2" < "a12".
//
// Non-digit sequences and numbers are compared separately. The former are compared byte-wise, while the latter are
// compared numerically (except that the number of leading zeros is used as a tie-breaker, so "2" < "02").
//
// Limitations:
//   - only ASCII digits (0-9) are considered.
//
// Original algorithm: https://github.com/fvbommel/util/blob/master/sortorder/natsort.go
func NaturalCmp(s1, s2 string, caseInsensitive bool) int {
	i1 := 0
	i2 := 0
	for i1 < len(s1) && i2 < len(s2) {
		c1 := s1[i1]
		c2 := s2[i2]
		d1 := c1 >= '0' && c1 <= '9'
		d2 := c2 >= '0' && c2 <= '9'
		switch {
		case d1 != d2: // Digits before other characters.
			if d1 { // True if LHS is a digit, false if the RHS is one.
				return -1
			}
			return 1
		case !d1: // && !d2, because d1 == d2
			// UTF-8 compares byte-wise-lexicographically, no need to decode code-points.
			if caseInsensitive {
				if c1 >= 'a' && c1 <= 'z' {
					c1 -= 'a' - 'A'
				}
				if c2 >= 'a' && c2 <= 'z' {
					c2 -= 'a' - 'A'
				}
			}
			if c1 != c2 {
				if c1 < c2 {
					return -1
				}
				return 1
			}
			i1++
			i2++
		default: // Digits
			// Eat zeros.
			for i1 < len(s1) && s1[i1] == '0' {
				i1++
			}
			for i1 < len(s1) && s1[i1] == '0' {
				i1++
			}
			for i2 < len(s2) && s2[i2] == '0' {
				i2++
			}
			// Eat all digits.
			nz1, nz2 := i1, i2
			for i1 < len(s1) && s1[i1] >= '0' && s1[i1] <= '9' {
				i1++
			}
			for i2 < len(s2) && s2[i2] >= '0' && s2[i2] <= '9' {
				i2++
			}
			// If lengths of numbers with non-zero prefix differ, the shorter one is less.
			if len1, len2 := i1-nz1, i2-nz2; len1 != len2 {
				if len1 < len2 {
					return -1
				}
				return 1
			}
			// If they're not equal, string comparison is correct.
			if nr1, nr2 := s1[nz1:i1], s2[nz2:i2]; nr1 != nr2 {
				if nr1 < nr2 {
					return -1
				}
				return 1
			}
			// Otherwise, the one with less zeros is less. Because everything up to the number is equal, comparing the
			// index after the zeros is sufficient.
			if nz1 != nz2 {
				if nz1 < nz2 {
					return -1
				}
				return 1
			}
		}
		// They're identical so far, so continue comparing.
	}
	// So far they are identical. At least one is ended. If the other continues, it sorts last. If the are the same
	// length and the caseInsensitive flag was set, compare again, but without the flag.
	switch {
	case len(s1) == len(s2):
		if caseInsensitive {
			return NaturalCmp(s1, s2, false)
		}
		return 0
	case len(s1) < len(s2):
		return -1
	default:
		return 1
	}
}

// SortStringsNaturalAscending sorts a slice of strings using NaturalLess in least to most order.
func SortStringsNaturalAscending(in []string) {
	slices.SortFunc(in, func(a, b string) int { return NaturalCmp(a, b, true) })
}

// SortStringsNaturalDescending sorts a slice of strings using NaturalLess in most to least order.
func SortStringsNaturalDescending(in []string) {
	slices.SortFunc(in, func(a, b string) int { return NaturalCmp(b, a, true) })
}
