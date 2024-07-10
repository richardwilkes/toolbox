// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package slice

import "slices"

// ZeroedDelete removes the elements s[i:j] from s, returning the modified slice. This function panics if s[i:j] is not
// a valid slice of s. This function modifies the contents of the slice s; it does not create a new slice. The elements
// that are removed are zeroed so that any references can be garbage collected. If you do not need this, use
// slices.Delete instead.
//
// Deprecated: As of Go 1.22, slices.Delete now zeroes out removed elements, so use it instead. This function was
// deprecated on March 29, 2024 and will be removed on or after January 1, 2025.
func ZeroedDelete[S ~[]E, E any](s S, i, j int) S {
	return slices.Delete(s, i, j)
}

// ZeroedDeleteFunc removes any elements from s for which del returns true, returning the modified slice. This function
// modifies the contents of the slice s; it does not create a new slice. The elements that are removed are zeroed so
// that any references can be garbage collected. If you do not need this, use slices.DeleteFunc instead.
//
// Deprecated: As of Go 1.22, slices.DeleteFunc now zeroes out removed elements, so use it instead. This function was
// deprecated on March 29, 2024 and will be removed on or after January 1, 2025.
func ZeroedDeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	return slices.DeleteFunc(s, del)
}
