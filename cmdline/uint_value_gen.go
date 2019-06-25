// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

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
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
