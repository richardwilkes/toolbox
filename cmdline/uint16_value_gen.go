// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type uint16Value uint16

// NewUint16Option creates a new uint16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint16Option(val *uint16) *Option {
	return cl.NewOption((*uint16Value)(val))
}

// Set implements the Value interface.
func (val *uint16Value) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 16)
	*val = uint16Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint16Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type uint16ArrayValue []uint16

// NewUint16ArrayOption creates a new []uint16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint16ArrayOption(val *[]uint16) *Option {
	return cl.NewOption((*uint16ArrayValue)(val))
}

// Set implements the Value interface.
func (val *uint16ArrayValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 16)
	*val = append(*val, uint16(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint16ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
