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

type int64Value int64

// NewInt64Option creates a new int64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt64Option(val *int64) *Option {
	return cl.NewOption((*int64Value)(val))
}

// Set implements the Value interface.
func (val *int64Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = int64Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int64Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int64ArrayValue []int64

// NewInt64ArrayOption creates a new []int64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt64ArrayOption(val *[]int64) *Option {
	return cl.NewOption((*int64ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int64ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int64ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
