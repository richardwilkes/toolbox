// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type intValue int

// NewIntOption creates a new int Option and attaches it to this CmdLine.
func (cl *CmdLine) NewIntOption(val *int) *Option {
	return cl.NewOption((*intValue)(val))
}

// Set implements the Value interface.
func (val *intValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = intValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *intValue) String() string {
	return fmt.Sprintf("%v", *val)
}

type intArrayValue []int

// NewIntArrayOption creates a new []int Option and attaches it to this CmdLine.
func (cl *CmdLine) NewIntArrayOption(val *[]int) *Option {
	return cl.NewOption((*intArrayValue)(val))
}

// Set implements the Value interface.
func (val *intArrayValue) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, 64)
	*val = append(*val, int(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *intArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
