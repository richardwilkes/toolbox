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

type uintValue uint

// NewUintOption creates a new uint Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUintOption(val *uint) *Option {
	return cl.NewOption((*uintValue)(val))
}

// Set implements the Value interface.
func (val *uintValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 64)
	*val = uintValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uintValue) String() string {
	return fmt.Sprintf("%v", *val)
}

type uintArrayValue []uint

// NewUintArrayOption creates a new []uint Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUintArrayOption(val *[]uint) *Option {
	return cl.NewOption((*uintArrayValue)(val))
}

// Set implements the Value interface.
func (val *uintArrayValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 64)
	*val = append(*val, uint(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uintArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
