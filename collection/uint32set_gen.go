// Code generated - DO NOT EDIT.
package collection

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

// Uint32Set holds a set of uint32 values.
type Uint32Set map[uint32]bool

// NewUint32Set creates a new set from its input values.
func NewUint32Set(values ...uint32) Uint32Set {
	s := Uint32Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Uint32Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Uint32Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Uint32Set) Add(values ...uint32) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Uint32Set) Contains(value uint32) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Uint32Set) Clone() Uint32Set {
	if s == nil {
		return nil
	}
	clone := Uint32Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Uint32Set) Values() []uint32 {
	values := make([]uint32, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Uint32Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Uint32Set) UnmarshalJSON(data []byte) error {
	var values []uint32
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Uint32Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Uint32Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []uint32
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
