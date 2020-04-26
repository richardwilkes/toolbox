// Code created from "values.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

type int8Value int8

// NewInt8Option creates a new int8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt8Option(val *int8) *Option {
	return cl.NewOption((*int8Value)(val))
}

// Set implements the Value interface.
func (val *int8Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 8)
	*val = int8Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int8Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int8ArrayValue []int8

// NewInt8ArrayOption creates a new []int8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt8ArrayOption(val *[]int8) *Option {
	return cl.NewOption((*int8ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int8ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 8)
	*val = append(*val, int8(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int8ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
