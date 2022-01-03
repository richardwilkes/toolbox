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
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

type durationValue time.Duration

// NewDurationOption creates a new time.Duration Option and attaches it to this CmdLine.
func (cl *CmdLine) NewDurationOption(val *time.Duration) *Option {
	return cl.NewOption((*durationValue)(val))
}

// Set implements the Value interface.
func (val *durationValue) Set(str string) error {
	v, err := time.ParseDuration(str)
	*val = durationValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *durationValue) String() string {
	return time.Duration(*val).String()
}

type durationArrayValue []time.Duration

// NewDurationArrayOption creates a new []time.Duration Option and attaches it to this CmdLine.
func (cl *CmdLine) NewDurationArrayOption(val *[]time.Duration) *Option {
	return cl.NewOption((*durationArrayValue)(val))
}

// Set implements the Value interface.
func (val *durationArrayValue) Set(str string) error {
	v, err := time.ParseDuration(str)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *durationArrayValue) String() string {
	var buffer strings.Builder
	for _, v := range *val {
		if buffer.Len() != 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(v.String())
	}
	return buffer.String()
}
