// Code generated - DO NOT EDIT.
package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Complex64Set holds a set of complex64 values.
type Complex64Set map[complex64]bool

// NewComplex64Set creates a new set from its input values.
func NewComplex64Set(values ...complex64) Complex64Set {
	s := Complex64Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Complex64Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Complex64Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Complex64Set) Add(values ...complex64) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Complex64Set) Contains(value complex64) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Complex64Set) Clone() Complex64Set {
	if s == nil {
		return nil
	}
	clone := Complex64Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Complex64Set) Values() []complex64 {
	values := make([]complex64, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Complex64Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Complex64Set) UnmarshalJSON(data []byte) error {
	var values []complex64
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Complex64Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Complex64Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []complex64
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
