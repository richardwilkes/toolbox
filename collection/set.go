// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package collection

// Set holds a set of values.
type Set[T comparable] map[T]struct{}

// NewSet creates a new set from its input values.
func NewSet[T comparable](values ...T) Set[T] {
	s := make(Set[T], len(values))
	s.Add(values...)
	return s
}

// Len returns the number of values in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// Empty returns true if there are no values in the set.
func (s Set[T]) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s *Set[T]) Clear() {
	*s = make(Set[T])
}

// Add values to the set.
func (s Set[T]) Add(values ...T) {
	for _, v := range values {
		s[v] = struct{}{}
	}
}

// Contains returns true if the value exists within the set.
func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Set[T]) Clone() Set[T] {
	if s == nil {
		return nil
	}
	clone := make(Set[T], len(s))
	for value := range s {
		clone[value] = struct{}{}
	}
	return clone
}

// Values returns all values in the set.
func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}
