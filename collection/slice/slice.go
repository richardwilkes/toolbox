// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package slice

// ZeroedDelete removes the elements s[i:j] from s, returning the modified slice. This function panics if s[i:j] is not
// a valid slice of s. This function modifies the contents of the slice s; it does not create a new slice. The elements
// that are removed are zeroed so that any references can be garbage collected. If you do not need this, use
// slices.Delete instead.
func ZeroedDelete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j] // bounds check
	copy(s[i:], s[j:])
	var e E
	end := len(s) - j + i
	for k, n := end, len(s); k < n; k++ {
		s[k] = e
	}
	return s[:end]
}

// ZeroedDeleteFunc removes any elements from s for which del returns true, returning the modified slice. This function
// modifies the contents of the slice s; it does not create a new slice. The elements that are removed are zeroed so
// that any references can be garbage collected. If you do not need this, use slices.DeleteFunc instead.
func ZeroedDeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	// Don't start copying elements until we find one to delete.
	for i, v := range s {
		if del(v) {
			j := i
			for i++; i < len(s); i++ {
				v = s[i]
				if !del(v) {
					s[j] = v
					j++
				}
			}
			var e E
			for k, n := j, len(s); k < n; k++ {
				s[k] = e
			}
			return s[:j]
		}
	}
	return s
}
