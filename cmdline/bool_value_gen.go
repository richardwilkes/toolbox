// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

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
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
