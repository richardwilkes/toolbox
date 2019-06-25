// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type float64Value float64

// NewFloat64Option creates a new float64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat64Option(val *float64) *Option {
	return cl.NewOption((*float64Value)(val))
}

// Set implements the Value interface.
func (val *float64Value) Set(str string) error {
	v, err := strconv.ParseFloat(str, 64)
	*val = float64Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float64Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type float64ArrayValue []float64

// NewFloat64ArrayOption creates a new []float64 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat64ArrayOption(val *[]float64) *Option {
	return cl.NewOption((*float64ArrayValue)(val))
}

// Set implements the Value interface.
func (val *float64ArrayValue) Set(str string) error {
	v, err := strconv.ParseFloat(str, 64)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float64ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
