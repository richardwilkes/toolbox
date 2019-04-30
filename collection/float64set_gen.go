// Code generated - DO NOT EDIT.

package collection

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

// Float64Set holds a set of float64 values.
type Float64Set map[float64]bool

// NewFloat64Set creates a new set from its input values.
func NewFloat64Set(values ...float64) Float64Set {
	s := Float64Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Float64Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Float64Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Float64Set) Add(values ...float64) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Float64Set) Contains(value float64) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Float64Set) Clone() Float64Set {
	if s == nil {
		return nil
	}
	clone := Float64Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Float64Set) Values() []float64 {
	values := make([]float64, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Float64Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Float64Set) UnmarshalJSON(data []byte) error {
	var values []float64
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Float64Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Float64Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []float64
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
