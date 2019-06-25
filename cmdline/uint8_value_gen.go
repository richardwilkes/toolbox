// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type uint8Value uint8

// NewUint8Option creates a new uint8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint8Option(val *uint8) *Option {
	return cl.NewOption((*uint8Value)(val))
}

// Set implements the Value interface.
func (val *uint8Value) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 8)
	*val = uint8Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint8Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type uint8ArrayValue []uint8

// NewUint8ArrayOption creates a new []uint8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint8ArrayOption(val *[]uint8) *Option {
	return cl.NewOption((*uint8ArrayValue)(val))
}

// Set implements the Value interface.
func (val *uint8ArrayValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 8)
	*val = append(*val, uint8(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint8ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
