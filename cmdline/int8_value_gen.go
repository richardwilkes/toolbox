// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type int8Value int8

// NewInt8Option creates a new int8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt8Option(val *int8) *Option {
	return cl.NewOption((*int8Value)(val))
}

// Set implements the Value interface.
func (val *int8Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 8)
	*val = int8Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int8Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int8ArrayValue []int8

// NewInt8ArrayOption creates a new []int8 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt8ArrayOption(val *[]int8) *Option {
	return cl.NewOption((*int8ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int8ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 8)
	*val = append(*val, int8(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int8ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
