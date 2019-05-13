// Package fixed provides simple fixed-point values that can be added,
// subtracted, multiplied and divided.
package fixed

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
)

const (
	// Max holds the maximum fixed-point value.
	Max = Fixed(1<<63 - 1)
	// Min holds the minimum fixed-point value.
	Min = Fixed(-1 << 63)
)

var (
	precision  int
	multiplier int64
)

// Fixed holds a fixed-point value that contains up to N decimal places, where
// N is the value passed to SetDigitsAfterDecimal (default is 4). Values are
// truncated, not rounded. Values can be added and subtracted directly. For
// multiplication and division, the provided Mul() and Div() methods should be
// used.
type Fixed int64

func init() {
	SetDigitsAfterDecimal(4)
}

// SetDigitsAfterDecimal controls the number of digits after the decimal place
// that are tracked. WARNING: This has a global effect on all fixed-point
// values and should only be set once prior to use of this package. Changes to
// this value invalidate any fixed-point values there were created prior to
// the call -- there is no enforcement of this, however, so use of a
// pre-existing value will quietly generate bad results.
func SetDigitsAfterDecimal(digits int) {
	precision = digits
	multiplier = int64(math.Pow(10, float64(precision)))
}

// New creates a new fixed-point value.
// Deprecated: Use FromFloat64() instead.
func New(value float64) Fixed {
	return FromFloat64(value)
}

// FromFloat64 creates a new fixed-point value from a float64.
func FromFloat64(value float64) Fixed {
	return Fixed(value * float64(multiplier))
}

// FromInt creates a new fixed-point value from an int.
func FromInt(value int) Fixed {
	return Fixed(int64(value) * multiplier)
}

// Parse a string to extract a fixed-point value from it.
func Parse(str string) (Fixed, error) {
	if str == "" {
		return 0, errs.New("Empty string is not valid")
	}
	if strings.ContainsRune(str, 'E') {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, err
		}
		return FromFloat64(f), nil
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
		value *= multiplier
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < precision+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > precision+1 {
			frac = frac[:precision+1]
		}
		if fraction, err = strconv.ParseInt(frac, 10, 64); err != nil {
			return 0, errs.Wrap(err)
		}
		value += fraction - multiplier
	}
	if neg {
		value = -value
	}
	return Fixed(value), nil
}

// Mul multiplies this value by the passed-in value, returning a new
// fixed-point value.
func (fxd Fixed) Mul(value Fixed) Fixed {
	return fxd * value / Fixed(multiplier)
}

// Div divides this value by the passed-in value, returning a new fixed-point
// value.
func (fxd Fixed) Div(value Fixed) Fixed {
	return fxd * Fixed(multiplier) / value
}

// Trunc returns a new value which has everything to the right of the decimal
// place truncated.
func (fxd Fixed) Trunc() Fixed {
	return fxd / Fixed(multiplier) * Fixed(multiplier)
}

// Int64 returns the truncated equivalent integer to this fixed-point value.
func (fxd Fixed) Int64() int64 {
	return int64(fxd / Fixed(multiplier))
}

// Float64 returns the floating-point equivalent to this fixed-point value.
func (fxd Fixed) Float64() float64 {
	return float64(fxd) / float64(multiplier)
}

// Comma returns the same as String(), but with commas for values of 1000 and
// greater.
func (fxd Fixed) Comma() string {
	integer := fxd / Fixed(multiplier)
	fraction := fxd % Fixed(multiplier)
	if fraction == 0 {
		return humanize.Comma(int64(integer))
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += Fixed(multiplier)
	fstr := fmt.Sprintf("%d", fraction)
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer == 0 && fxd < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%s.%s", neg, humanize.Comma(int64(integer)), fstr)
}

func (fxd Fixed) String() string {
	integer := fxd / Fixed(multiplier)
	fraction := fxd % Fixed(multiplier)
	if fraction == 0 {
		return fmt.Sprintf("%d", integer)
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += Fixed(multiplier)
	fstr := fmt.Sprintf("%d", fraction)
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer == 0 && fxd < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%d.%s", neg, integer, fstr)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (fxd *Fixed) MarshalText() ([]byte, error) {
	return []byte(fxd.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (fxd *Fixed) UnmarshalText(text []byte) error {
	f, err := Parse(string(text))
	if err != nil {
		return err
	}
	*fxd = f
	return nil
}

// MarshalJSON implements the json.Marshaler interface. Note that this
// intentionally generates a string where necessary to ensure the correct
// value is retained.
func (fxd *Fixed) MarshalJSON() ([]byte, error) {
	f := fxd.Float64()
	str := fxd.String()
	if FromFloat64(f) == *fxd && fmt.Sprint(f) == str {
		return json.Marshal(f)
	}
	return json.Marshal(str)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (fxd *Fixed) UnmarshalJSON(data []byte) error {
	var dv interface{}
	err := json.Unmarshal(data, &dv)
	if err != nil {
		return errs.Wrap(err)
	}
	var f Fixed
	switch v := dv.(type) {
	case string:
		f, err = Parse(v)
		if err != nil {
			return err
		}
	case float64:
		f = FromFloat64(v)
	default:
		return errs.New("Invalid type")
	}
	*fxd = f
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface. Note that this
// intentionally generates a string where necessary to ensure the correct
// value is retained.
func (fxd Fixed) MarshalYAML() (interface{}, error) {
	f := fxd.Float64()
	str := fxd.String()
	if FromFloat64(f) == fxd && fmt.Sprint(f) == str {
		return f, nil
	}
	return str, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (fxd *Fixed) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return errs.Wrap(err)
	}
	f, err := Parse(str)
	if err != nil {
		return err
	}
	*fxd = f
	return nil
}
