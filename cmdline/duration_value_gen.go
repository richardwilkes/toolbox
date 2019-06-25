// Code created from "values.go.tmpl" - don't edit by hand

package cmdline

import (
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

type durationValue time.Duration

// NewDurationOption creates a new time.Duration Option and attaches it to this CmdLine.
func (cl *CmdLine) NewDurationOption(val *time.Duration) *Option {
	return cl.NewOption((*durationValue)(val))
}

// Set implements the Value interface.
func (val *durationValue) Set(str string) error {
	v, err := time.ParseDuration(str)
	*val = durationValue(v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *durationValue) String() string {
	return time.Duration(*val).String()
}

type durationArrayValue []time.Duration

// NewDurationArrayOption creates a new []time.Duration Option and attaches it to this CmdLine.
func (cl *CmdLine) NewDurationArrayOption(val *[]time.Duration) *Option {
	return cl.NewOption((*durationArrayValue)(val))
}

// Set implements the Value interface.
func (val *durationArrayValue) Set(str string) error {
	v, err := time.ParseDuration(str)
	*val = append(*val, v)
	return errs.Wrap(err)
}

// String implements the Value interface.
func (val *durationArrayValue) String() string {
	var str string
	for _, v := range *val {
		if str == "" {
			str += ", "
		}
		str += v.String()
	}
	return str
}
