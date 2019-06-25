// Code created from "set.go.tmpl" - don't edit by hand

package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Uint16Set holds a set of uint16 values.
type Uint16Set map[uint16]bool

// NewUint16Set creates a new set from its input values.
func NewUint16Set(values ...uint16) Uint16Set {
	s := Uint16Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Uint16Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Uint16Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Uint16Set) Add(values ...uint16) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Uint16Set) Contains(value uint16) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Uint16Set) Clone() Uint16Set {
	if s == nil {
		return nil
	}
	clone := Uint16Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Uint16Set) Values() []uint16 {
	values := make([]uint16, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Uint16Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Uint16Set) UnmarshalJSON(data []byte) error {
	var values []uint16
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Uint16Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Uint16Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []uint16
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
