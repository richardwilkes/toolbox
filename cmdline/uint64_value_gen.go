// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type uint64Value uint64

// NewUint64Option creates a new uint64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint64Option(val *uint64) *Option {
	return cl.NewOption((*uint64Value)(val))
}

// Set implements the Value interface.
func (val *uint64Value) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 64)
	*val = uint64Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint64Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type uint64ArrayValue []uint64

// NewUint64ArrayOption creates a new []uint64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewUint64ArrayOption(val *[]uint64) *Option {
	return cl.NewOption((*uint64ArrayValue)(val))
}

// Set implements the Value interface.
func (val *uint64ArrayValue) Set(str string) error {
	v, err := strconv.ParseUint(str, 0, 64)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *uint64ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
