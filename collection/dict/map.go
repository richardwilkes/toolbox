// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package dict

// Functions in this pcakge are only here because the Go team decided not to bring them over in the Go 1.21 maps package
// when they migrated the existing code from golang.org/x/exp/maps. Why, I'm not sure, since these can be useful.
//
// I chose not to use the package name "maps" to avoid collisions with the standard library. Unlike with the "slices"
// package, though, I couldn't use "map", since that's a keyword.

// Keys returns the keys of the map m. The keys will be in an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Values returns the values of the map m. The values will be in an indeterminate order.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
