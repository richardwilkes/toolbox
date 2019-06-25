// Code created from "set.go.tmpl" - don't edit by hand

package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Complex128Set holds a set of complex128 values.
type Complex128Set map[complex128]bool

// NewComplex128Set creates a new set from its input values.
func NewComplex128Set(values ...complex128) Complex128Set {
	s := Complex128Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Complex128Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Complex128Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Complex128Set) Add(values ...complex128) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Complex128Set) Contains(value complex128) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Complex128Set) Clone() Complex128Set {
	if s == nil {
		return nil
	}
	clone := Complex128Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Complex128Set) Values() []complex128 {
	values := make([]complex128, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Complex128Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Complex128Set) UnmarshalJSON(data []byte) error {
	var values []complex128
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Complex128Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Complex128Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []complex128
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
