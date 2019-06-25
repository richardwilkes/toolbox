// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type int32Value int32

// NewInt32Option creates a new int32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt32Option(val *int32) *Option {
	return cl.NewOption((*int32Value)(val))
}

// Set implements the Value interface.
func (val *int32Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 32)
	*val = int32Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int32Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int32ArrayValue []int32

// NewInt32ArrayOption creates a new []int32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt32ArrayOption(val *[]int32) *Option {
	return cl.NewOption((*int32ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int32ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 32)
	*val = append(*val, int32(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int32ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
