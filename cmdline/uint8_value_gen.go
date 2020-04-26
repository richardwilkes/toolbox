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

type uint8Value uint8

// NewUint8Option creates a new uint8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint8Option(val *uint8) *Option {
	return cl.NewOption((*uint8Value)(val))
}

// Set implements the Value interface.
func (val *uint8Value) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 8)
	*val = uint8Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint8Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type uint8ArrayValue []uint8

// NewUint8ArrayOption creates a new []uint8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint8ArrayOption(val *[]uint8) *Option {
	return cl.NewOption((*uint8ArrayValue)(val))
}

// Set implements the Value interface.
func (val *uint8ArrayValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 8)
	*val = append(*val, uint8(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint8ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
