// Code created from "set.go.tmpl" - don't edit by hand

package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Uint8Set holds a set of uint8 values.
type Uint8Set map[uint8]bool

// NewUint8Set creates a new set from its input values.
func NewUint8Set(values ...uint8) Uint8Set {
	s := Uint8Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Uint8Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Uint8Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Uint8Set) Add(values ...uint8) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Uint8Set) Contains(value uint8) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Uint8Set) Clone() Uint8Set {
	if s == nil {
		return nil
	}
	clone := Uint8Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Uint8Set) Values() []uint8 {
	values := make([]uint8, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Uint8Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Uint8Set) UnmarshalJSON(data []byte) error {
	var values []uint8
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Uint8Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Uint8Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []uint8
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
