// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type int16Value int16

// NewInt16Option creates a new int16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt16Option(val *int16) *Option {
	return cl.NewOption((*int16Value)(val))
}

// Set implements the Value interface.
func (val *int16Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 16)
	*val = int16Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int16Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int16ArrayValue []int16

// NewInt16ArrayOption creates a new []int16 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt16ArrayOption(val *[]int16) *Option {
	return cl.NewOption((*int16ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int16ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 16)
	*val = append(*val, int16(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int16ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
