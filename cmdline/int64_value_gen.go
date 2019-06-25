// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type int64Value int64

// NewInt64Option creates a new int64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt64Option(val *int64) *Option {
	return cl.NewOption((*int64Value)(val))
}

// Set implements the Value interface.
func (val *int64Value) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = int64Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int64Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type int64ArrayValue []int64

// NewInt64ArrayOption creates a new []int64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewInt64ArrayOption(val *[]int64) *Option {
	return cl.NewOption((*int64ArrayValue)(val))
}

// Set implements the Value interface.
func (val *int64ArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *int64ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
