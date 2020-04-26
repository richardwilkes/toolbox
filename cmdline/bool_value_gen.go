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

type boolValue bool

// NewBoolOption creates a new bool Option and attaches it to this CmdLine.
func (cl *CmdLine) NewBoolOption(val *bool) *Option {
	return cl.NewOption((*boolValue)(val))
}

// Set implements the Value interface.
func (val *boolValue) Set(str string) error {
	v, err := strconv.ParseBool(str)
	*val = boolValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *boolValue) String() string {
	return fmt.Sprintf("%v", *val)
}

type boolArrayValue []bool

// NewBoolArrayOption creates a new []bool Option and attaches it to this CmdLine.
func (cl *CmdLine) NewBoolArrayOption(val *[]bool) *Option {
	return cl.NewOption((*boolArrayValue)(val))
}

// Set implements the Value interface.
func (val *boolArrayValue) Set(str string) error {
	v, err := strconv.ParseBool(str)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *boolArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
