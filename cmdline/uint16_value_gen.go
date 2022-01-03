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

type uint16Value uint16

// NewUint16Option creates a new uint16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint16Option(val *uint16) *Option {
	return cl.NewOption((*uint16Value)(val))
}

// Set implements the Value interface.
func (val *uint16Value) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 16)
	*val = uint16Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint16Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type uint16ArrayValue []uint16

// NewUint16ArrayOption creates a new []uint16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint16ArrayOption(val *[]uint16) *Option {
	return cl.NewOption((*uint16ArrayValue)(val))
}

// Set implements the Value interface.
func (val *uint16ArrayValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 16)
	*val = append(*val, uint16(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint16ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
