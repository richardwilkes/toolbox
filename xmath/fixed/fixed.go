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
	// Max holds the maximum fixed-point value.
	Max = Fixed(math.MaxInt64)
	// Min holds the minimum fixed-point value.
	Min = Fixed(math.MinInt64)
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
//
// Deprecated: Use one of the F64d... types instead.
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
//
// Deprecated.
func SetDigitsAfterDecimal(digits int) {
	precision = digits
	multiplier = int64(math.Pow(10, float64(precision)))
}

// New creates a new fixed-point value.
//
// Deprecated: Use FromFloat64() instead.
func New(value float64) Fixed {
	return FromFloat64(value)
}

// FromFloat64 creates a new fixed-point value from a float64.
//
// Deprecated: Use one of the F64d... types instead.
func FromFloat64(value float64) Fixed {
	return Fixed(value * float64(multiplier))
}

// FromInt creates a new fixed-point value from an int.
//
// Deprecated: Use one of the F64d... types instead.
func FromInt(value int) Fixed {
	return Fixed(int64(value) * multiplier)
}

// Parse a string to extract a fixed-point value from it.
//
// Deprecated: Use one of the F64d... types instead.
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
func (fxd Fixed) MarshalText() ([]byte, error) {
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
