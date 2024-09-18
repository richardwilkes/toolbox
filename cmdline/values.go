// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

var _ Value = &GeneralValue{}

// GeneralValue holds a general option value. Valid value types are: *bool, *int, *int8, *int16, *int32, *int64, *uint,
// *uint8, *uint16, *uint32, *uint64, *float32, *float64, *string, *time.Duration, *[]bool, *[]uint8, *[]uint16,
// *[]uint32, *[]uint64, *[]int8, *[]int16, *[]int32, *[]int64, *[]string, *[]time.Duration
type GeneralValue struct {
	Value any
}

// Set implements Value
func (v *GeneralValue) Set(str string) error {
	var err error
	var signedValue int64
	var unsignedValue uint64
	var floatValue float64
	switch value := v.Value.(type) {
	case *bool:
		if *value, err = strconv.ParseBool(str); err != nil {
			return errs.Wrap(err)
		}
	case *int:
		if signedValue, err = strconv.ParseInt(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
		*value = int(signedValue)
	case *int8:
		if signedValue, err = strconv.ParseInt(str, 0, 8); err != nil {
			return errs.Wrap(err)
		}
		*value = int8(signedValue)
	case *int16:
		if signedValue, err = strconv.ParseInt(str, 0, 16); err != nil {
			return errs.Wrap(err)
		}
		*value = int16(signedValue)
	case *int32:
		if signedValue, err = strconv.ParseInt(str, 0, 32); err != nil {
			return errs.Wrap(err)
		}
		*value = int32(signedValue)
	case *int64:
		if *value, err = strconv.ParseInt(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
	case *uint:
		if unsignedValue, err = strconv.ParseUint(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
		*value = uint(unsignedValue)
	case *uint8:
		if unsignedValue, err = strconv.ParseUint(str, 0, 8); err != nil {
			return errs.Wrap(err)
		}
		*value = uint8(unsignedValue)
	case *uint16:
		if unsignedValue, err = strconv.ParseUint(str, 0, 16); err != nil {
			return errs.Wrap(err)
		}
		*value = uint16(unsignedValue)
	case *uint32:
		if unsignedValue, err = strconv.ParseUint(str, 0, 32); err != nil {
			return errs.Wrap(err)
		}
		*value = uint32(unsignedValue)
	case *uint64:
		if *value, err = strconv.ParseUint(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
	case *float32:
		if floatValue, err = strconv.ParseFloat(str, 32); err != nil {
			return errs.Wrap(err)
		}
		*value = float32(floatValue)
	case *float64:
		if *value, err = strconv.ParseFloat(str, 64); err != nil {
			return errs.Wrap(err)
		}
	case *string:
		*value = str
	case *time.Duration:
		if *value, err = time.ParseDuration(str); err != nil {
			return errs.Wrap(err)
		}
	case *[]bool:
		var b bool
		if b, err = strconv.ParseBool(str); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, b)
	case *[]int:
		if signedValue, err = strconv.ParseInt(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, int(signedValue))
	case *[]int8:
		if signedValue, err = strconv.ParseInt(str, 0, 8); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, int8(signedValue))
	case *[]int16:
		if signedValue, err = strconv.ParseInt(str, 0, 16); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, int16(signedValue))
	case *[]int32:
		if signedValue, err = strconv.ParseInt(str, 0, 32); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, int32(signedValue))
	case *[]int64:
		if signedValue, err = strconv.ParseInt(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, signedValue)
	case *[]uint:
		if unsignedValue, err = strconv.ParseUint(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, uint(unsignedValue))
	case *[]uint8:
		if unsignedValue, err = strconv.ParseUint(str, 0, 8); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, uint8(unsignedValue))
	case *[]uint16:
		if unsignedValue, err = strconv.ParseUint(str, 0, 16); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, uint16(unsignedValue))
	case *[]uint32:
		if unsignedValue, err = strconv.ParseUint(str, 0, 32); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, uint32(unsignedValue))
	case *[]uint64:
		if unsignedValue, err = strconv.ParseUint(str, 0, 64); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, unsignedValue)
	case *[]string:
		*value = append(*value, str)
	case *[]time.Duration:
		var d time.Duration
		if d, err = time.ParseDuration(str); err != nil {
			return errs.Wrap(err)
		}
		*value = append(*value, d)
	default:
		return errs.Newf("<unhandled type: %v>", reflect.TypeOf(v.Value))
	}
	return nil
}

func (v *GeneralValue) String() string {
	k := reflect.TypeOf(v.Value).Kind()
	if k != reflect.Ptr {
		return fmt.Sprintf("<unhandled type: %v>", k)
	}
	e := reflect.ValueOf(v.Value).Elem()
	switch e.Kind() {
	case reflect.Slice:
		var buffer strings.Builder
		count := e.Len()
		for i := 0; i < count; i++ {
			if buffer.Len() != 0 {
				buffer.WriteString(", ")
			}
			fmt.Fprintf(&buffer, "%v", e.Index(i))
		}
		return buffer.String()
	case reflect.String:
		return fmt.Sprintf(`"%s"`, e) //nolint:gocritic // We don't want escape sequences, so can't use %q
	default:
		return fmt.Sprintf("%v", e)
	}
}
