// Code created from "set.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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

	"gopkg.in/yaml.v2"
)

// Float32Set holds a set of float32 values.
type Float32Set map[float32]bool

// NewFloat32Set creates a new set from its input values.
func NewFloat32Set(values ...float32) Float32Set {
	s := Float32Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Float32Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Float32Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Float32Set) Add(values ...float32) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Float32Set) Contains(value float32) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Float32Set) Clone() Float32Set {
	if s == nil {
		return nil
	}
	clone := Float32Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Float32Set) Values() []float32 {
	values := make([]float32, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Float32Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Float32Set) UnmarshalJSON(data []byte) error {
	var values []float32
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Float32Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Float32Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []float32
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
