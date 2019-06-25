// Code created from "set.go.tmpl" - don't edit by hand

package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Uint64Set holds a set of uint64 values.
type Uint64Set map[uint64]bool

// NewUint64Set creates a new set from its input values.
func NewUint64Set(values ...uint64) Uint64Set {
	s := Uint64Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Uint64Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Uint64Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Uint64Set) Add(values ...uint64) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Uint64Set) Contains(value uint64) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Uint64Set) Clone() Uint64Set {
	if s == nil {
		return nil
	}
	clone := Uint64Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Uint64Set) Values() []uint64 {
	values := make([]uint64, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Uint64Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Uint64Set) UnmarshalJSON(data []byte) error {
	var values []uint64
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Uint64Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Uint64Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []uint64
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
