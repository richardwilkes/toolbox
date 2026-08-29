// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings

import (
	"slices"
)

// NaturalLess reports whether s1 sorts before s2 in natural order. See NaturalCmp.
func NaturalLess(s1, s2 string, caseInsensitive bool) bool {
	return NaturalCmp(s1, s2, caseInsensitive) < 0
}

// NaturalCmp compares two strings using natural ordering, so "a2" < "a12", and returns -1, 0, or 1. Non-digit
// sequences are compared byte-wise and digit sequences numerically, with the number of leading zeros as a
// tie-breaker, so "2" < "02". Only ASCII digits are recognized. When caseInsensitive is true, ASCII letters are first
// compared ignoring case; if the strings are equal that way, the comparison falls back to case-sensitive.
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
			// UTF-8 byte order matches code point order, so no decoding is needed.
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
			for i1 < len(s1) && s1[i1] == '0' {
				i1++
			}
			for i2 < len(s2) && s2[i2] == '0' {
				i2++
			}
			nz1, nz2 := i1, i2
			for i1 < len(s1) && s1[i1] >= '0' && s1[i1] <= '9' {
				i1++
			}
			for i2 < len(s2) && s2[i2] >= '0' && s2[i2] <= '9' {
				i2++
			}
			// After stripping leading zeros, the shorter number is less.
			if len1, len2 := i1-nz1, i2-nz2; len1 != len2 {
				if len1 < len2 {
					return -1
				}
				return 1
			}
			// Same length, so string comparison is numeric comparison.
			if nr1, nr2 := s1[nz1:i1], s2[nz2:i2]; nr1 != nr2 {
				if nr1 < nr2 {
					return -1
				}
				return 1
			}
			// Same number, so fewer leading zeros is less. Everything before the number is equal, so comparing the
			// indices after the zeros suffices.
			if nz1 != nz2 {
				if nz1 < nz2 {
					return -1
				}
				return 1
			}
		}
	}
	// Identical so far and at least one has ended, so the longer sorts last. If the same length and caseInsensitive,
	// compare again case-sensitively.
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

// SortStringsNaturalAscending sorts the slice in ascending order using NaturalCmp with caseInsensitive true.
func SortStringsNaturalAscending(in []string) {
	slices.SortFunc(in, func(a, b string) int { return NaturalCmp(a, b, true) })
}

// SortStringsNaturalDescending sorts the slice in descending order using NaturalCmp with caseInsensitive true.
func SortStringsNaturalDescending(in []string) {
	slices.SortFunc(in, func(a, b string) int { return NaturalCmp(b, a, true) })
}
