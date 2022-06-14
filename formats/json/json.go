// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package json provides manipulation of JSON data.
package json

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

// Data provides conveniences for working with JSON data.
type Data struct {
	obj any
}

// MustParse is the same as calling Parse, but without the error code on return.
func MustParse(data []byte) *Data {
	result, err := Parse(data)
	if err != nil {
		result = &Data{}
	}
	return result
}

// Parse JSON data from bytes. If the data can't be loaded, a valid, empty Data will still be returned, along with an
// error.
func Parse(data []byte) (*Data, error) {
	var obj any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&obj); err != nil {
		return &Data{}, err
	}
	return &Data{obj: obj}, nil
}

// MustParseStream is the same as calling ParseStream, but without the error code on return.
func MustParseStream(in io.Reader) *Data {
	result, err := ParseStream(in)
	if err != nil {
		result = &Data{}
	}
	return result
}

// ParseStream parses JSON data from the stream. If the data can't be loaded, a valid, empty Data will still be
// returned, along with an error.
func ParseStream(in io.Reader) (*Data, error) {
	var obj any
	decoder := json.NewDecoder(in)
	decoder.UseNumber()
	if err := decoder.Decode(&obj); err != nil {
		return &Data{}, err
	}
	return &Data{obj: obj}, nil
}

// Raw returns the underlying data.
func (j *Data) Raw() any {
	return j.obj
}

// Path searches the dot-separated path and returns the object at that point. If the search encounters an array and has
// not reached the end target, then it will iterate through the array for the target and return all results in a Data
// array.
func (j *Data) Path(path string) *Data {
	return j.path(strings.Split(path, ".")...)
}

func (j *Data) path(path ...string) *Data {
	if len(path) == 1 && path[0] == "" {
		path = nil
	}
	obj := j.obj
	for i := 0; i < len(path); i++ {
		if m, ok := obj.(map[string]any); ok {
			obj = m[path[i]]
		} else {
			var a []any
			if a, ok = obj.([]any); ok {
				t := make([]any, 0)
				for _, one := range a {
					tj := &Data{obj: one}
					if result := tj.path(path[i:]...).obj; result != nil {
						t = append(t, result)
					}
				}
				if len(a) != 0 {
					return &Data{obj: t}
				}
			}
			return &Data{}
		}
	}
	return &Data{obj}
}

// Exists returns true if the path exists in the data.
func (j *Data) Exists(path string) bool {
	return j.Path(path).obj != nil
}

// IsArray returns true if this is a Data array.
func (j *Data) IsArray() bool {
	_, ok := j.obj.([]any)
	return ok
}

// IsMap returns true if this is a Data map.
func (j *Data) IsMap() bool {
	_, ok := j.obj.(map[string]any)
	return ok
}

// Keys returns the keys of a map, or an empty slice if this is not a map.
func (j *Data) Keys() []string {
	if m, ok := j.obj.(map[string]any); ok {
		keys := make([]string, 0, len(m))
		for key := range m {
			keys = append(keys, key)
		}
		return keys
	}
	return make([]string, 0)
}

// Size returns the number of elements in an array or map, or 0 if this is neither type.
func (j *Data) Size() int {
	if m, ok := j.obj.(map[string]any); ok {
		return len(m)
	}
	if a, ok := j.obj.([]any); ok {
		return len(a)
	}
	return 0
}

// Index returns the object at the specified index within an array, or nil if this isn't an array or the index isn't
// valid.
func (j *Data) Index(index int) *Data {
	if a, ok := j.obj.([]any); ok {
		if index >= 0 && index < len(a) {
			return &Data{obj: a[index]}
		}
	}
	return &Data{}
}

// Bytes converts the data into a Data []byte.
func (j *Data) Bytes() []byte {
	if j.obj != nil {
		if data, err := json.Marshal(j.obj); err == nil {
			return data
		}
	}
	return []byte("{}")
}

// String converts the data into a Data string.
func (j *Data) String() string {
	return string(j.Bytes())
}

// Str extracts a string from the path. Returns the empty string if the path isn't present or isn't a string type.
func (j *Data) Str(path string) string {
	if str, ok := j.Path(path).obj.(string); ok {
		return str
	}
	return ""
}

// Bool extracts a bool from the path. Returns false if the path isn't present or isn't a boolean type.
func (j *Data) Bool(path string) bool {
	if b, ok := j.Path(path).obj.(bool); ok {
		return b
	}
	return false
}

// BoolRelaxed extracts a bool from the path. Returns false if the path isn't present or can't be converted to a boolean
// type.
func (j *Data) BoolRelaxed(path string) bool {
	if b, ok := j.Path(path).obj.(bool); ok {
		return b
	}
	return strings.EqualFold(j.Str(path), "true")
}

// Float64 extracts an float64 from the path. Returns 0 if the path isn't present or isn't a numeric type.
func (j *Data) Float64(path string) float64 {
	if n, ok := j.Path(path).obj.(json.Number); ok {
		if f, err := n.Float64(); err == nil {
			return f
		}
	}
	return 0
}

// Float64Relaxed extracts an float64 from the path. Returns 0 if the path isn't present or can't be converted to a
// numeric type.
func (j *Data) Float64Relaxed(path string) float64 {
	if n, ok := j.Path(path).obj.(json.Number); ok {
		if f, err := n.Float64(); err == nil {
			return f
		}
	} else {
		if f, err := strconv.ParseFloat(j.Str(path), 64); err == nil {
			return f
		}
	}
	return 0
}

// Int64 extracts an int64 from the path. Returns 0 if the path isn't present or isn't a numeric type.
func (j *Data) Int64(path string) int64 {
	if n, ok := j.Path(path).obj.(json.Number); ok {
		if i, err := n.Int64(); err == nil {
			return i
		}
	}
	return 0
}

// Int64Relaxed extracts an int64 from the path. Returns 0 if the path isn't present or can't be converted to a numeric
// type.
func (j *Data) Int64Relaxed(path string) int64 {
	if n, ok := j.Path(path).obj.(json.Number); ok {
		if i, err := n.Int64(); err == nil {
			return i
		}
	} else {
		if i, err := strconv.ParseInt(j.Str(path), 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// Unmarshal parses the data at the path and stores the result into value.
func (j *Data) Unmarshal(path string, value any) error {
	return json.Unmarshal(j.Path(path).Bytes(), value)
}

// NewMap creates a map at the specified path. Any parts of the path that do not exist will be created. Returns true if
// successful, or false if a collision occurs with a non-object type while traversing the path.
func (j *Data) NewMap(path string) bool {
	return j.set(path, make(map[string]any))
}

// NewArray creates an array at the specified path. Any parts of the path that do not exist will be created. Returns
// true if successful, or false if a collision occurs with a non-object type while traversing the path.
func (j *Data) NewArray(path string) bool {
	return j.set(path, make([]any, 0))
}

// SetStr a string at the specified path. Any parts of the path that do not exist will be created. Returns true if
// successful, or false if a collision occurs with a non-object type while traversing the path.
func (j *Data) SetStr(path, value string) bool {
	return j.set(path, value)
}

// SetBool a bool at the specified path. Any parts of the path that do not exist will be created. Returns true if
// successful, or false if a collision occurs with a non-object type while traversing the path.
func (j *Data) SetBool(path string, value bool) bool {
	return j.set(path, value)
}

// SetFloat64 a float64 at the specified path. Any parts of the path that do not exist will be created. Returns true if
// successful, or false if a collision occurs with a non-object type while traversing the path.
func (j *Data) SetFloat64(path string, value float64) bool {
	return j.set(path, value)
}

// SetInt64 an int64 at the specified path. Any parts of the path that do not exist will be created. Returns true if
// successful, or false if a collision occurs with a non-object type while traversing the path.
func (j *Data) SetInt64(path string, value int64) bool {
	return j.set(path, value)
}

// Set a Data value at the specified path. Any parts of the path that do not exist will be created. Returns true if
// successful, or false if a collision occurs with a non-object type while traversing the path.
func (j *Data) Set(path string, value *Data) bool {
	var v any
	if value != nil {
		v = value.obj
	}
	return j.set(path, v)
}

func (j *Data) set(path string, value any) bool {
	paths := strings.Split(path, ".")
	if len(paths) == 0 {
		j.obj = value
	} else {
		if j.obj == nil {
			j.obj = make(map[string]any)
		}
		obj := j.obj
		for i := 0; i < len(paths); i++ {
			if m, ok := obj.(map[string]any); ok {
				if i == len(paths)-1 {
					m[paths[i]] = value
				} else if m[paths[i]] == nil {
					m[paths[i]] = make(map[string]any)
				}
				obj = m[paths[i]]
			} else {
				return false
			}
		}
	}
	return true
}

// AppendMap appends a new map to an array at the specified path. The array must already exist. Returns true if
// successful.
func (j *Data) AppendMap(path string) bool {
	return j.append(path, make(map[string]any))
}

// AppendArray appends a new array to an array at the specified path. The array must already exist. Returns true if
// successful.
func (j *Data) AppendArray(path string) bool {
	return j.append(path, make([]any, 0))
}

// AppendStr appends a string to an array at the specified path. The array must already exist. Returns true if
// successful.
func (j *Data) AppendStr(path, value string) bool {
	return j.append(path, value)
}

// AppendBool appends a bool to an array at the specified path. The array must already exist. Returns true if
// successful.
func (j *Data) AppendBool(path string, value bool) bool {
	return j.append(path, value)
}

// AppendFloat64 appends a float64 to an array at the specified path. The array must already exist. Returns true if
// successful.
func (j *Data) AppendFloat64(path string, value float64) bool {
	return j.append(path, value)
}

// AppendInt64 appends an int64 to an array at the specified path. The array must already exist. Returns true if
// successful.
func (j *Data) AppendInt64(path string, value int64) bool {
	return j.append(path, value)
}

// Append a Data value to an array at the specified path. The array must already exist. Returns true if successful.
func (j *Data) Append(path string, value *Data) bool {
	var v any
	if value != nil {
		v = value.obj
	}
	return j.append(path, v)
}

func (j *Data) append(path string, value any) bool {
	if array, ok := j.Path(path).obj.([]any); ok {
		return j.set(path, append(array, value))
	}
	return false
}

// Delete a value at the specified path. Returns true if successful.
func (j *Data) Delete(path string) bool {
	if j.obj == nil {
		return false
	}
	paths := strings.Split(path, ".")
	if len(paths) == 0 {
		j.obj = nil
		return true
	}
	obj := j.obj
	for i := 0; i < len(paths); i++ {
		if m, ok := obj.(map[string]any); ok {
			if i == len(paths)-1 {
				if _, ok = m[paths[i]]; ok {
					delete(m, paths[i])
					return true
				}
			}
			obj = m[paths[i]]
		}
	}
	return false
}
