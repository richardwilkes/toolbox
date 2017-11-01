package cmdline

import (
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
)

// Value represents a value that can be set for an Option.
type Value interface {
	// Set the contents of this value.
	Set(value string) error
	// String implements the fmt.Stringer interface.
	String() string
}

// Option represents an option available on the command line.
type Option struct {
	name   string
	usage  string
	arg    string
	def    string
	value  Value
	single rune
}

func (op *Option) isValid() (bool, error) {
	if op.value == nil {
		return false, errs.New(i18n.Text("Option must have a value"))
	}
	if op.single == 0 && op.name == "" {
		return false, errs.New(i18n.Text("Option must be named"))
	}
	return true, nil
}

func (op *Option) isBool() bool {
	_, ok := op.value.(*boolValue)
	return ok
}

func (op *Option) isString() bool {
	_, ok := op.value.(*stringValue)
	return ok
}

// SetName sets the name for this option. Returns self for easy chaining.
func (op *Option) SetName(name string) *Option {
	if len(name) > 1 {
		op.name = name
	} else {
		panic("Name must be 2+ characters")
	}
	return op
}

// SetSingle sets the single character name for this option. Returns self for
// easy chaining.
func (op *Option) SetSingle(ch rune) *Option {
	op.single = ch
	return op
}

// SetArg sets the argument name for this option. Returns self for easy
// chaining.
func (op *Option) SetArg(name string) *Option {
	op.arg = name
	return op
}

// SetDefault sets the default value for this option. Returns self for easy
// chaining.
func (op *Option) SetDefault(def string) *Option {
	op.def = def
	return op
}

// SetUsage sets the usage message for this option. Returns self for easy
// chaining.
func (op *Option) SetUsage(usage string) *Option {
	op.usage = usage
	return op
}
