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

type int32Value int32

// NewInt32Option creates a new int32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt32Option(val *int32) *Option {
	return cl.NewOption((*int32Value)(val))
}

// Set implements the Value interface.
func (val *int32Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 32)
	*val = int32Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int32Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int32ArrayValue []int32

// NewInt32ArrayOption creates a new []int32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt32ArrayOption(val *[]int32) *Option {
	return cl.NewOption((*int32ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int32ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 32)
	*val = append(*val, int32(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int32ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
