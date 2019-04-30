// Code generated - DO NOT EDIT.

package collection

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

// Int32Set holds a set of int32 values.
type Int32Set map[int32]bool

// NewInt32Set creates a new set from its input values.
func NewInt32Set(values ...int32) Int32Set {
	s := Int32Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Int32Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Int32Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Int32Set) Add(values ...int32) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Int32Set) Contains(value int32) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Int32Set) Clone() Int32Set {
	if s == nil {
		return nil
	}
	clone := Int32Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Int32Set) Values() []int32 {
	values := make([]int32, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Int32Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Int32Set) UnmarshalJSON(data []byte) error {
	var values []int32
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Int32Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Int32Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []int32
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
