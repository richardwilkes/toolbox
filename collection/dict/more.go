// Copyright (c) 2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package dict

// MapByKey returns a map of the values in 'data' keyed by the result of keyFunc. If there are duplicate keys, the last
// value in data with that key will be the one in the map.
func MapByKey[T any, K comparable](data []T, keyFunc func(T) K) map[K]T {
	m := make(map[K]T)
	for _, v := range data {
		m[keyFunc(v)] = v
	}
	return m
}

// MapOfSlicesByKey returns a map of the values in 'data' keyed by the result of keyFunc. Duplicate keys will have their
// values appended to a slice in the map.
func MapOfSlicesByKey[T any, K comparable](data []T, keyFunc func(T) K) map[K][]T {
	m := make(map[K][]T)
	for _, v := range data {
		k := keyFunc(v)
		m[k] = append(m[k], v)
	}
	return m
}
