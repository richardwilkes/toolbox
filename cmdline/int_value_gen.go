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

type intValue int

// NewIntOption creates a new int Option and attaches it to this CmdLine.
func (cl *CmdLine) NewIntOption(val *int) *Option {
	return cl.NewOption((*intValue)(val))
}

// Set implements the Value interface.
func (val *intValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = intValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *intValue) String() string {
	return fmt.Sprintf("%v", *val)
}

type intArrayValue []int

// NewIntArrayOption creates a new []int Option and attaches it to this CmdLine.
func (cl *CmdLine) NewIntArrayOption(val *[]int) *Option {
	return cl.NewOption((*intArrayValue)(val))
}

// Set implements the Value interface.
func (val *intArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = append(*val, int(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *intArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%v", v)
	}
	return buffer.String()
}
