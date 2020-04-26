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

type float64Value float64

// NewFloat64Option creates a new float64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat64Option(val *float64) *Option {
	return cl.NewOption((*float64Value)(val))
}

// Set implements the Value interface.
func (val *float64Value) Set(str string) error {
	v, err := strconv.ParseFloat(str, 64)
	*val = float64Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float64Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type float64ArrayValue []float64

// NewFloat64ArrayOption creates a new []float64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat64ArrayOption(val *[]float64) *Option {
	return cl.NewOption((*float64ArrayValue)(val))
}

// Set implements the Value interface.
func (val *float64ArrayValue) Set(str string) error {
	v, err := strconv.ParseFloat(str, 64)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float64ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
