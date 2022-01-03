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

type float32Value float32

// NewFloat32Option creates a new float32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat32Option(val *float32) *Option {
	return cl.NewOption((*float32Value)(val))
}

// Set implements the Value interface.
func (val *float32Value) Set(str string) error {
	v, err := strconv.ParseFloat(str, 32)
	*val = float32Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float32Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type float32ArrayValue []float32

// NewFloat32ArrayOption creates a new []float32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat32ArrayOption(val *[]float32) *Option {
	return cl.NewOption((*float32ArrayValue)(val))
}

// Set implements the Value interface.
func (val *float32ArrayValue) Set(str string) error {
	v, err := strconv.ParseFloat(str, 32)
	*val = append(*val, float32(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float32ArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
