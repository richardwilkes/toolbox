package xjson

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// JSON provides conveniences for working with JSON data.
type JSON struct {
	obj interface{}
}

// MustParseJSON is the same as calling ParseJSON, but without the error code
// on return.
func MustParseJSON(data []byte) *JSON {
	result, err := ParseJSON(data)
	if err != nil {
		result = &JSON{}
	}
	return result
}

// ParseJSON from the data. If the data can't be loaded, a valid, empty JSON
// will still be returned, along with an error.
func ParseJSON(data []byte) (*JSON, error) {
	var obj interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&obj); err != nil {
		return &JSON{}, err
	}
	return &JSON{obj: obj}, nil
}

// MustParseJSONStream is the same as calling ParseJSONStream, but without the
// error code on return.
func MustParseJSONStream(in io.Reader) *JSON {
	result, err := ParseJSONStream(in)
	if err != nil {
		result = &JSON{}
	}
	return result
}

// ParseJSONStream from the stream. If the data can't be loaded, a valid,
// empty JSON will still be returned, along with an error.
func ParseJSONStream(in io.Reader) (*JSON, error) {
	var obj interface{}
	decoder := json.NewDecoder(in)
	decoder.UseNumber()
	if err := decoder.Decode(&obj); err != nil {
		return &JSON{}, err
	}
	return &JSON{obj: obj}, nil
}

// Data returns the underlying data.
func (j *JSON) Data() interface{} {
	return j.obj
}

// Path searches the dot-separated path and returns the object at that point.
// If the search encounters an array and has not reached the end target, then
// it will iterate through the array for the target and return all results in
// a JSON array.
func (j *JSON) Path(path string) *JSON {
	return j.path(strings.Split(path, ".")...)
}

func (j *JSON) path(path ...string) *JSON {
	obj := j.obj
	for i := 0; i < len(path); i++ {
		if m, ok := obj.(map[string]interface{}); ok {
			obj = m[path[i]]
		} else if a, ok := obj.([]interface{}); ok {
			t := make([]interface{}, 0)
			for _, one := range a {
				tj := &JSON{obj: one}
				if result := tj.path(path[i:]...).obj; result != nil {
					t = append(t, result)
				}
			}
			if len(a) == 0 {
				return &JSON{}
			}
			return &JSON{obj: t}
		} else {
			return &JSON{}
		}
	}
	return &JSON{obj}
}

// Exists returns true if the path exists in the data.
func (j *JSON) Exists(path string) bool {
	return j.Path(path).obj != nil
}

// IsArray returns true if this is a JSON array.
func (j *JSON) IsArray() bool {
	_, ok := j.obj.([]interface{})
	return ok
}

// IsMap returns true if this is a JSON map.
func (j *JSON) IsMap() bool {
	_, ok := j.obj.(map[string]interface{})
	return ok
}

// Keys returns the keys of a map, or an empty slice if this is not a map.
func (j *JSON) Keys() []string {
	if m, ok := j.obj.(map[string]interface{}); ok {
		keys := make([]string, 0, len(m))
		for key := range m {
			keys = append(keys, key)
		}
		return keys
	}
	return make([]string, 0)
}

// Size returns the number of elements in an array or map, or 0 if this is
// neither type.
func (j *JSON) Size() int {
	if m, ok := j.obj.(map[string]interface{}); ok {
		return len(m)
	}
	if a, ok := j.obj.([]interface{}); ok {
		return len(a)
	}
	return 0
}

// Index returns the object at the specified index within an array, or nil if
// this isn't an array or the index isn't valid.
func (j *JSON) Index(index int) *JSON {
	if a, ok := j.obj.([]interface{}); ok {
		if index >= 0 && index < len(a) {
			return &JSON{obj: a[index]}
		}
	}
	return &JSON{}
}

// Bytes converts the data into a JSON []byte.
func (j *JSON) Bytes() []byte {
	if j.obj != nil {
		if data, err := json.Marshal(j.obj); err == nil {
			return data
		}
	}
	return []byte("{}")
}

// String converts the data into a JSON string.
func (j *JSON) String() string {
	return string(j.Bytes())
}

// Str extracts a string from the path. Returns the empty string if the path
// isn't present or isn't a string type.
func (j *JSON) Str(path string) string {
	if str, ok := j.Path(path).obj.(string); ok {
		return str
	}
	return ""
}

// Bool extracts a bool from the path. Returns false if the path isn't present
// or isn't a boolean type.
func (j *JSON) Bool(path string) bool {
	if b, ok := j.Path(path).obj.(bool); ok {
		return b
	}
	return false
}

// Float64 extracts an float64 from the path. Returns 0 if the path isn't
// present or isn't a numeric type.
func (j *JSON) Float64(path string) float64 {
	if n, ok := j.Path(path).obj.(json.Number); ok {
		if f, err := n.Float64(); err == nil {
			return f
		}
	}
	return 0
}

// Int64 extracts an int64 from the path. Returns 0 if the path isn't present
// or isn't a numeric type.
func (j *JSON) Int64(path string) int64 {
	if n, ok := j.Path(path).obj.(json.Number); ok {
		if i, err := n.Int64(); err == nil {
			return i
		}
	}
	return 0
}

// Unmarshal parses the data at the path and stores the result into value.
func (j *JSON) Unmarshal(path string, value interface{}) error {
	return json.Unmarshal(j.Path(path).Bytes(), value)
}

// NewMap creates a map at the specified path. Any parts of the path that do
// not exist will be created. Returns true if successful, or false if a
// collision occurs with a non-object type while traversing the path.
func (j *JSON) NewMap(path string) bool {
	return j.set(path, make(map[string]interface{}))
}

// NewArray creates an array at the specified path. Any parts of the path that
// do not exist will be created. Returns true if successful, or false if a
// collision occurs with a non-object type while traversing the path.
func (j *JSON) NewArray(path string) bool {
	return j.set(path, make([]interface{}, 0))
}

// SetStr a string at the specified path. Any parts of the path that do not
// exist will be created. Returns true if successful, or false if a collision
// occurs with a non-object type while traversing the path.
func (j *JSON) SetStr(path, value string) bool {
	return j.set(path, value)
}

// SetBool a bool at the specified path. Any parts of the path that do not
// exist will be created. Returns true if successful, or false if a collision
// occurs with a non-object type while traversing the path.
func (j *JSON) SetBool(path string, value bool) bool {
	return j.set(path, value)
}

// SetFloat64 a float64 at the specified path. Any parts of the path that do
// not exist will be created. Returns true if successful, or false if a
// collision occurs with a non-object type while traversing the path.
func (j *JSON) SetFloat64(path string, value float64) bool {
	return j.set(path, value)
}

// SetInt64 an int64 at the specified path. Any parts of the path that do not
// exist will be created. Returns true if successful, or false if a collision
// occurs with a non-object type while traversing the path.
func (j *JSON) SetInt64(path string, value int64) bool {
	return j.set(path, value)
}

// Set a JSON value at the specified path. Any parts of the path that do not
// exist will be created. Returns true if successful, or false if a collision
// occurs with a non-object type while traversing the path.
func (j *JSON) Set(path string, value *JSON) bool {
	var v interface{}
	if value != nil {
		v = value.obj
	}
	return j.set(path, v)
}

func (j *JSON) set(path string, value interface{}) bool {
	paths := strings.Split(path, ".")
	if len(paths) == 0 {
		j.obj = value
	} else {
		if j.obj == nil {
			j.obj = make(map[string]interface{})
		}
		obj := j.obj
		for i := 0; i < len(paths); i++ {
			if m, ok := obj.(map[string]interface{}); ok {
				if i == len(paths)-1 {
					m[paths[i]] = value
				} else if m[paths[i]] == nil {
					m[paths[i]] = make(map[string]interface{})
				}
				obj = m[paths[i]]
			} else {
				return false
			}
		}
	}
	return true
}

// AppendMap appends a new map to an array at the specified path. The array
// must already exist. Returns true if successful.
func (j *JSON) AppendMap(path string) bool {
	return j.append(path, make(map[string]interface{}))
}

// AppendArray appends a new array to an array at the specified path. The
// array must already exist. Returns true if successful.
func (j *JSON) AppendArray(path string) bool {
	return j.append(path, make([]interface{}, 0))
}

// AppendStr appends a string to an array at the specified path. The array
// must already exist. Returns true if successful.
func (j *JSON) AppendStr(path, value string) bool {
	return j.append(path, value)
}

// AppendBool appends a bool to an array at the specified path. The array must
// already exist. Returns true if successful.
func (j *JSON) AppendBool(path string, value bool) bool {
	return j.append(path, value)
}

// AppendFloat64 appends a float64 to an array at the specified path. The
// array must already exist. Returns true if successful.
func (j *JSON) AppendFloat64(path string, value float64) bool {
	return j.append(path, value)
}

// AppendInt64 appends an int64 to an array at the specified path. The array
// must already exist. Returns true if successful.
func (j *JSON) AppendInt64(path string, value int64) bool {
	return j.append(path, value)
}

// Append a JSON value to an array at the specified path. The array must
// already exist. Returns true if successful.
func (j *JSON) Append(path string, value *JSON) bool {
	var v interface{}
	if value != nil {
		v = value.obj
	}
	return j.append(path, v)
}

func (j *JSON) append(path string, value interface{}) bool {
	if array, ok := j.Path(path).obj.([]interface{}); ok {
		return j.set(path, append(array, value))
	}
	return false
}

// Delete a value at the specified path. Returns true if successful.
func (j *JSON) Delete(path string) bool {
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
		if m, ok := obj.(map[string]interface{}); ok {
			if i == len(paths)-1 {
				if _, ok := m[paths[i]]; ok {
					delete(m, paths[i])
					return true
				}
			}
			obj = m[paths[i]]
		}
	}
	return false
}
