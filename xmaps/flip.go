// Copyright (c) 2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package xmaps provides map utilities.
package xmaps

// Flip returns a new map with the keys and values of m swapped. If multiple keys in m share the same value, only one
// of them will be present in the result and which one is unspecified, since map iteration order is not defined.
func Flip[MI ~map[T1]T2, T1 comparable, T2 comparable](m MI) map[T2]T1 {
	flipped := make(map[T2]T1, len(m))
	for k, v := range m {
		flipped[v] = k
	}
	return flipped
}

// FlipAlt is the same as Flip, but allows the type of the returned map to be specified, which is useful when a named
// map type is desired rather than the plain map[T2]T1 that Flip returns.
func FlipAlt[MI ~map[T1]T2, MO ~map[T2]T1, T1 comparable, T2 comparable](m MI) MO {
	flipped := make(MO, len(m))
	for k, v := range m {
		flipped[v] = k
	}
	return flipped
}
