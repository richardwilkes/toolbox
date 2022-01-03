// Code created from "set.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Int64Set holds a set of int64 values.
type Int64Set map[int64]bool

// NewInt64Set creates a new set from its input values.
func NewInt64Set(values ...int64) Int64Set {
	s := Int64Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Int64Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Int64Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Int64Set) Add(values ...int64) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Int64Set) Contains(value int64) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Int64Set) Clone() Int64Set {
	if s == nil {
		return nil
	}
	clone := Int64Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Int64Set) Values() []int64 {
	values := make([]int64, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Int64Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Int64Set) UnmarshalJSON(data []byte) error {
	var values []int64
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Int64Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Int64Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []int64
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
