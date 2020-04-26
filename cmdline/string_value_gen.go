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
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

type stringValue string

// NewStringOption creates a new string Option and attaches it to this CmdLine.
func (cl *CmdLine) NewStringOption(val *string) *Option {
	return cl.NewOption((*stringValue)(val))
}

// Set implements the Value interface.
func (val *stringValue) Set(str string) error {
	v, err := str, error(nil)
	*val = stringValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *stringValue) String() string {
	return string(*val)
}

type stringArrayValue []string

// NewStringArrayOption creates a new []string Option and attaches it to this CmdLine.
func (cl *CmdLine) NewStringArrayOption(val *[]string) *Option {
	return cl.NewOption((*stringArrayValue)(val))
}

// Set implements the Value interface.
func (val *stringArrayValue) Set(str string) error {
	v, err := str, error(nil)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *stringArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(v)
	}
	return buffer.String()
}
