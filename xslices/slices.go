// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslices

// Set returns a map using the values of 'data' as its keys.
func Set[K comparable](data []K) map[K]struct{} {
	m := make(map[K]struct{})
	for _, v := range data {
		m[v] = struct{}{}
	}
	return m
}

// MapFromData returns a map of the values in 'data', keyed by 'keyFunc'.
func MapFromData[K comparable, V any](data []V, keyFunc func(V) K) map[K]V {
	m := make(map[K]V)
	for _, v := range data {
		m[keyFunc(v)] = v
	}
	return m
}

// MapFromKeys returns a map of 'keys' with each value supplied by 'dataFunc'.
func MapFromKeys[K comparable, V any](keys []K, dataFunc func(K) V) map[K]V {
	m := make(map[K]V)
	for _, k := range keys {
		m[k] = dataFunc(k)
	}
	return m
}
