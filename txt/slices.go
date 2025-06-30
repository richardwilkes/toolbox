// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import (
	"strings"

	"github.com/richardwilkes/toolbox/v2/collection/dict"
)

// StringSliceToMap returns a map created from the strings of a slice.
func StringSliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool, len(slice))
	for _, str := range slice {
		m[str] = true
	}
	return m
}

// MapToStringSlice returns a slice created from the keys of a map.
//
// Deprecated: Use dict.Keys instead. This function was deprecated on May 3, 2024 and will be removed on or after
// January 1, 2025.
func MapToStringSlice(m map[string]bool) []string {
	return dict.Keys(m)
}

// CloneStringSlice returns a copy of the slice of strings.
func CloneStringSlice(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

// RunesEqual returns true if the two slices of runes are equal.
func RunesEqual(left, right []rune) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

// CaselessSliceContains returns true if the target is within the slice, regardless of case.
func CaselessSliceContains(slice []string, target string) bool {
	for _, one := range slice {
		if strings.EqualFold(one, target) {
			return true
		}
	}
	return false
}
