// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

type float32Value float32

// NewFloat32Option creates a new float32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat32Option(val *float32) *Option {
	return cl.NewOption((*float32Value)(val))
}

// Set implements the Value interface.
func (val *float32Value) Set(str string) error {
	v, err := strconv.ParseFloat(str, 32)
	*val = float32Value(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float32Value) String() string {
	return fmt.Sprintf("%v", *val)
}

type float32ArrayValue []float32

// NewFloat32ArrayOption creates a new []float32 Option and attaches it to this CmdLine.
func (cl *CmdLine) NewFloat32ArrayOption(val *[]float32) *Option {
	return cl.NewOption((*float32ArrayValue)(val))
}

// Set implements the Value interface.
func (val *float32ArrayValue) Set(str string) error {
	v, err := strconv.ParseFloat(str, 32)
	*val = append(*val, float32(v))
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *float32ArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}
