// Code generated - DO NOT EDIT.
package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// UintSet holds a set of uint values.
type UintSet map[uint]bool

// NewUintSet creates a new set from its input values.
func NewUintSet(values ...uint) UintSet {
	s := UintSet{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s UintSet) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s UintSet) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s UintSet) Add(values ...uint) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s UintSet) Contains(value uint) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s UintSet) Clone() UintSet {
	if s == nil {
		return nil
	}
	clone := UintSet{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s UintSet) Values() []uint {
	values := make([]uint, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s UintSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s UintSet) UnmarshalJSON(data []byte) error {
	var values []uint
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s UintSet) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s UintSet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []uint
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
