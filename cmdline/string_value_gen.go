// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"github.com/richardwilkes/toolbox/errs"
)

type stringValue string

// NewStringOption creates a new string Option and attaches it to this CmdLine.
func (cl *CmdLine) NewStringOption(val *string) *Option {
	return cl.NewOption((*stringValue)(val))
}

// Set implements the Value interface.
func (val *stringValue) Set(str string) error {
	v, err := str, error(nil)
	*val = stringValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *stringValue) String() string {
	return string(*val)
}

type stringArrayValue []string

// NewStringArrayOption creates a new []string Option and attaches it to this CmdLine.
func (cl *CmdLine) NewStringArrayOption(val *[]string) *Option {
	return cl.NewOption((*stringArrayValue)(val))
}

// Set implements the Value interface.
func (val *stringArrayValue) Set(str string) error {
	v, err := str, error(nil)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *stringArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += v
	}
	return str
}
