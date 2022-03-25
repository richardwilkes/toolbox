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

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
)

// Value represents a value that can be set for an Option.
type Value interface {
	// Set the contents of this value.
	Set(value string) error
	// String implements the fmt.Stringer interface.
	String() string
}

// Option represents an option available on the command line.
type Option struct {
	name   string
	usage  string
	arg    string
	def    string
	value  Value
	single rune
}

func (op *Option) isValid() (bool, error) {
	if op.value == nil {
		return false, errs.New(i18n.Text("Option must have a value"))
	}
	if op.single == 0 && op.name == "" {
		return false, errs.New(i18n.Text("Option must be named"))
	}
	return true, nil
}

func (op *Option) isBool() bool {
	if generalValue, ok := op.value.(*GeneralValue); ok {
		if _, ok = generalValue.Value.(*bool); ok {
			return true
		}
	}
	return false
}

// SetName sets the name for this option. Returns self for easy chaining.
func (op *Option) SetName(name string) *Option {
	if len(name) > 1 {
		op.name = name
	} else {
		fmt.Println("Name must be 2+ characters")
		atexit.Exit(1)
	}
	return op
}

// SetSingle sets the single character name for this option. Returns self for easy chaining.
func (op *Option) SetSingle(ch rune) *Option {
	op.single = ch
	return op
}

// SetArg sets the argument name for this option. Returns self for easy chaining.
func (op *Option) SetArg(name string) *Option {
	op.arg = name
	return op
}

// SetDefault sets the default value for this option. Returns self for easy chaining.
func (op *Option) SetDefault(def string) *Option {
	op.def = def
	return op
}

// SetUsage sets the usage message for this option. Returns self for easy chaining.
func (op *Option) SetUsage(usage string) *Option {
	op.usage = usage
	return op
}
