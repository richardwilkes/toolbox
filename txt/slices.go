// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

// StringSliceToMap returns a map created from the strings of a slice.
func StringSliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool)
	for _, str := range slice {
		m[str] = true
	}
	return m
}

// MapToStringSlice returns a slice created from the keys of a map.
func MapToStringSlice(m map[string]bool) []string {
	s := make([]string, 0, len(m))
	for str := range m {
		s = append(s, str)
	}
	return s
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
