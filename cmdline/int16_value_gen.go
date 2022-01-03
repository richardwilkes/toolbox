// Code created from "values.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
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

type int16Value int16

// NewInt16Option creates a new int16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt16Option(val *int16) *Option {
	return cl.NewOption((*int16Value)(val))
}

// Set implements the Value interface.
func (val *int16Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 16)
	*val = int16Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int16Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int16ArrayValue []int16

// NewInt16ArrayOption creates a new []int16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt16ArrayOption(val *[]int16) *Option {
	return cl.NewOption((*int16ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int16ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 16)
	*val = append(*val, int16(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int16ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
