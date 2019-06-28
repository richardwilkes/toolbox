// Code created from "fixed64.go.tmpl" - don't edit by hand

package fixed

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
)

const (
	// F64d2Max holds the maximum F64d2 value.
	F64d2Max = F64d2(1<<63 - 1)
	// F64d2Min holds the minimum F64d2 value.
	F64d2Min = F64d2(^(1<<63 - 1))
)

var multiplierF64d2 = int64(math.Pow(10, 2))

// F64d2 holds a fixed-point value that contains up to 2 decimal places.
// Values are truncated, not rounded. Values can be added and subtracted
// directly. For multiplication and division, the provided Mul() and Div()
// methods should be used.
type F64d2 int64

// F64d2FromFloat64 creates a new F64d2 value from a float64.
func F64d2FromFloat64(value float64) F64d2 {
	return F64d2(value * float64(multiplierF64d2))
}

// F64d2FromInt64 creates a new F64d2 value from an int64.
func F64d2FromInt64(value int64) F64d2 {
	return F64d2(value * multiplierF64d2)
}

// F64d2FromString creates a new F64d2 value from a string.
func F64d2FromString(str string) (F64d2, error) {
	if str == "" {
		return 0, errs.New("empty string is not valid")
	}
	if strings.ContainsRune(str, 'E') {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, err
		}
		return F64d2FromFloat64(f), nil
	}
	parts := strings.SplitN(str, ".", 2)
	var value, fraction int64
	var neg bool
	var err error
	switch parts[0] {
	case "":
	case "-", "-0":
		neg = true
	default:
		if value, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
			return 0, errs.Wrap(err)
		}
		if value < 0 {
			neg = true
			value = -value
		}
		value *= multiplierF64d2
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < 2+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > 2+1 {
			frac = frac[:2+1]
		}
		if fraction, err = strconv.ParseInt(frac, 10, 64); err != nil {
			return 0, errs.Wrap(err)
		}
		value += fraction - multiplierF64d2
	}
	if neg {
		value = -value
	}
	return F64d2(value), nil
}

// F64d2FromStringForced creates a new F64d2 value from a string.
func F64d2FromStringForced(str string) F64d2 {
	f, _ := F64d2FromString(str) //nolint:errcheck
	return f
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f F64d2) Mul(value F64d2) F64d2 {
	return f * value / F64d2(multiplierF64d2)
}

// Div divides this value by the passed-in value, returning a new value.
func (f F64d2) Div(value F64d2) F64d2 {
	return f * F64d2(multiplierF64d2) / value
}

// Trunc returns a new value which has everything to the right of the decimal
// place truncated.
func (f F64d2) Trunc() F64d2 {
	return f / F64d2(multiplierF64d2) * F64d2(multiplierF64d2)
}

// Int64 returns the truncated equivalent integer to this value.
func (f F64d2) Int64() int64 {
	return int64(f / F64d2(multiplierF64d2))
}

// Float64 returns the floating-point equivalent to this value.
func (f F64d2) Float64() float64 {
	return float64(f) / float64(multiplierF64d2)
}

// Comma returns the same as String(), but with commas for values of 1000 and
// greater.
func (f F64d2) Comma() string {
	integer := f / F64d2(multiplierF64d2)
	fraction := f % F64d2(multiplierF64d2)
	if fraction == 0 {
		return humanize.Comma(int64(integer))
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += F64d2(multiplierF64d2)
	fstr := strconv.FormatInt(int64(fraction), 10)
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer == 0 && f < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%s.%s", neg, humanize.Comma(int64(integer)), fstr)
}

func (f F64d2) String() string {
	integer := f / F64d2(multiplierF64d2)
	fraction := f % F64d2(multiplierF64d2)
	if fraction == 0 {
		return strconv.FormatInt(int64(integer), 10)
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += F64d2(multiplierF64d2)
	fstr := strconv.FormatInt(int64(fraction), 10)
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer == 0 && f < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%d.%s", neg, integer, fstr)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (f *F64d2) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *F64d2) UnmarshalText(text []byte) error {
	f1, err := F64d2FromString(string(text))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// MarshalJSON implements the json.Marshaler interface. Note that this
// intentionally generates a string to ensure the correct value is retained.
func (f *F64d2) MarshalJSON() ([]byte, error) {
	return []byte(`"` + f.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *F64d2) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		if data[len(data)-1] != '"' {
			return fmt.Errorf("invalid JSON %q", string(data))
		}
		data = data[1 : len(data)-1]
	}
	v, err := F64d2FromString(string(data))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface. Note that this
// intentionally generates a string to ensure the correct value is retained.
func (f F64d2) MarshalYAML() (interface{}, error) {
	return f.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (f *F64d2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return errs.Wrap(err)
	}
	f1, err := F64d2FromString(str)
	if err != nil {
		return err
	}
	*f = f1
	return nil
}
