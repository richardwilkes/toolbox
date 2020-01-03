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

	"github.com/richardwilkes/toolbox/errs"
)

type uint32Value uint32

// NewUint32Option creates a new uint32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint32Option(val *uint32) *Option {
	return cl.NewOption((*uint32Value)(val))
}

// Set implements the Value interface.
func (val *uint32Value) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 32)
	*val = uint32Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint32Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type uint32ArrayValue []uint32

// NewUint32ArrayOption creates a new []uint32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint32ArrayOption(val *[]uint32) *Option {
	return cl.NewOption((*uint32ArrayValue)(val))
}

// Set implements the Value interface.
func (val *uint32ArrayValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 32)
	*val = append(*val, uint32(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint32ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
