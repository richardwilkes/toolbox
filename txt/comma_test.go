package txt

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestCommaInt(t *testing.T) {
	c := check.New(t)

	// Test integers
	c.Equal("0", CommaInt(0))
	c.Equal("1", CommaInt(1))
	c.Equal("12", CommaInt(12))
	c.Equal("123", CommaInt(123))
	c.Equal("1,234", CommaInt(1234))
	c.Equal("12,345", CommaInt(12345))
	c.Equal("123,456", CommaInt(123456))
	c.Equal("1,234,567", CommaInt(1234567))
	c.Equal("12,345,678", CommaInt(12345678))
	c.Equal("123,456,789", CommaInt(123456789))
	c.Equal("1,234,567,890", CommaInt(1234567890))

	// Test negative integers
	c.Equal("-1", CommaInt(-1))
	c.Equal("-12", CommaInt(-12))
	c.Equal("-123", CommaInt(-123))
	c.Equal("-1,234", CommaInt(-1234))
	c.Equal("-12,345", CommaInt(-12345))
	c.Equal("-123,456", CommaInt(-123456))
	c.Equal("-1,234,567", CommaInt(-1234567))

	// Test different numeric types
	c.Equal("42", CommaInt(int8(42)))
	c.Equal("1,234", CommaInt(int16(1234)))
	c.Equal("123,456", CommaInt(int32(123456)))
	c.Equal("1,234,567", CommaInt(int64(1234567)))
	c.Equal("255", CommaInt(uint8(255)))
	c.Equal("65,535", CommaInt(uint16(65535)))
	c.Equal("4,294,967,295", CommaInt(uint32(4294967295)))
}

func TestCommaFloat(t *testing.T) {
	c := check.New(t)

	// Test integers
	c.Equal("0", CommaFloat(0.0))
	c.Equal("1", CommaFloat(1.0))
	c.Equal("12", CommaFloat(12.0))
	c.Equal("123", CommaFloat(123.0))
	c.Equal("1,234", CommaFloat(1234.0))
	c.Equal("12,345", CommaFloat(12345.0))
	c.Equal("123,456", CommaFloat(123456.0))
	c.Equal("1,234,567", CommaFloat(1234567.0))
	c.Equal("12,345,678", CommaFloat(12345678.0))
	c.Equal("123,456,789", CommaFloat(123456789.0))
	c.Equal("1,234,567,890", CommaFloat(1234567890.0))

	// Test negative integers
	c.Equal("-1", CommaFloat(-1.0))
	c.Equal("-12", CommaFloat(-12.0))
	c.Equal("-123", CommaFloat(-123.0))
	c.Equal("-1,234", CommaFloat(-1234.0))
	c.Equal("-12,345", CommaFloat(-12345.0))
	c.Equal("-123,456", CommaFloat(-123456.0))
	c.Equal("-1,234,567", CommaFloat(-1234567.0))

	// Test floating point numbers
	c.Equal("1.5", CommaFloat(1.5))
	c.Equal("12.34", CommaFloat(12.34))
	c.Equal("123.456", CommaFloat(123.456))
	c.Equal("1,234.56", CommaFloat(1234.56))
	c.Equal("12,345.678", CommaFloat(12345.678))
	c.Equal("123,456.789", CommaFloat(123456.789))
	c.Equal("1,234,567.89", CommaFloat(1234567.89))

	// Test negative floating point numbers
	c.Equal("-1.5", CommaFloat(-1.5))
	c.Equal("-12.34", CommaFloat(-12.34))
	c.Equal("-123.456", CommaFloat(-123.456))
	c.Equal("-1,234.56", CommaFloat(-1234.56))
	c.Equal("-1,234.567", CommaFloat(-1234.567))
	c.Equal("-12,345.678", CommaFloat(-12345.678))

	// Test float32
	c.Equal("42.5", CommaFloat(float32(42.5)))
}

func TestCommaFromStringNum(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", CommaFromStringNum(""))

	// Test single digits
	c.Equal("0", CommaFromStringNum("0"))
	c.Equal("1", CommaFromStringNum("1"))
	c.Equal("9", CommaFromStringNum("9"))

	// Test two digits
	c.Equal("10", CommaFromStringNum("10"))
	c.Equal("99", CommaFromStringNum("99"))

	// Test three digits
	c.Equal("100", CommaFromStringNum("100"))
	c.Equal("123", CommaFromStringNum("123"))
	c.Equal("999", CommaFromStringNum("999"))

	// Test four digits (first comma)
	c.Equal("1,000", CommaFromStringNum("1000"))
	c.Equal("1,234", CommaFromStringNum("1234"))
	c.Equal("9,999", CommaFromStringNum("9999"))

	// Test five digits
	c.Equal("10,000", CommaFromStringNum("10000"))
	c.Equal("12,345", CommaFromStringNum("12345"))
	c.Equal("99,999", CommaFromStringNum("99999"))

	// Test six digits
	c.Equal("100,000", CommaFromStringNum("100000"))
	c.Equal("123,456", CommaFromStringNum("123456"))
	c.Equal("999,999", CommaFromStringNum("999999"))

	// Test seven digits (second comma)
	c.Equal("1,000,000", CommaFromStringNum("1000000"))
	c.Equal("1,234,567", CommaFromStringNum("1234567"))
	c.Equal("9,999,999", CommaFromStringNum("9999999"))

	// Test larger numbers
	c.Equal("12,345,678", CommaFromStringNum("12345678"))
	c.Equal("123,456,789", CommaFromStringNum("123456789"))
	c.Equal("1,234,567,890", CommaFromStringNum("1234567890"))
	c.Equal("12,345,678,901", CommaFromStringNum("12345678901"))
	c.Equal("123,456,789,012", CommaFromStringNum("123456789012"))
}

func TestCommaFromStringNumWithNegativeSign(t *testing.T) {
	c := check.New(t)

	// Test negative numbers
	c.Equal("-1", CommaFromStringNum("-1"))
	c.Equal("-12", CommaFromStringNum("-12"))
	c.Equal("-123", CommaFromStringNum("-123"))
	c.Equal("-1,234", CommaFromStringNum("-1234"))
	c.Equal("-12,345", CommaFromStringNum("-12345"))
	c.Equal("-123,456", CommaFromStringNum("-123456"))
	c.Equal("-1,234,567", CommaFromStringNum("-1234567"))
	c.Equal("-12,345,678", CommaFromStringNum("-12345678"))
}

func TestCommaFromStringNumWithPositiveSign(t *testing.T) {
	c := check.New(t)

	// Test positive numbers with explicit sign
	c.Equal("+1", CommaFromStringNum("+1"))
	c.Equal("+12", CommaFromStringNum("+12"))
	c.Equal("+123", CommaFromStringNum("+123"))
	c.Equal("+1,234", CommaFromStringNum("+1234"))
	c.Equal("+12,345", CommaFromStringNum("+12345"))
	c.Equal("+123,456", CommaFromStringNum("+123456"))
	c.Equal("+1,234,567", CommaFromStringNum("+1234567"))
}

func TestCommaFromStringNumWithDecimals(t *testing.T) {
	c := check.New(t)

	// Test decimal numbers
	c.Equal("1.0", CommaFromStringNum("1.0"))
	c.Equal("12.34", CommaFromStringNum("12.34"))
	c.Equal("123.456", CommaFromStringNum("123.456"))
	c.Equal("1,234.56", CommaFromStringNum("1234.56"))
	c.Equal("12,345.678", CommaFromStringNum("12345.678"))
	c.Equal("123,456.789", CommaFromStringNum("123456.789"))
	c.Equal("1,234,567.89", CommaFromStringNum("1234567.89"))

	// Test negative decimal numbers
	c.Equal("-1.0", CommaFromStringNum("-1.0"))
	c.Equal("-12.34", CommaFromStringNum("-12.34"))
	c.Equal("-123.456", CommaFromStringNum("-123.456"))
	c.Equal("-1,234.56", CommaFromStringNum("-1234.56"))
	c.Equal("-12,345.678", CommaFromStringNum("-12345.678"))

	// Test positive decimal numbers with explicit sign
	c.Equal("+1.0", CommaFromStringNum("+1.0"))
	c.Equal("+12.34", CommaFromStringNum("+12.34"))
	c.Equal("+1,234.56", CommaFromStringNum("+1234.56"))

	// Test decimal numbers with different decimal place counts
	c.Equal("1,234.1", CommaFromStringNum("1234.1"))
	c.Equal("1,234.12", CommaFromStringNum("1234.12"))
	c.Equal("1,234.123", CommaFromStringNum("1234.123"))
	c.Equal("1,234.1234", CommaFromStringNum("1234.1234"))
	c.Equal("1,234.12345", CommaFromStringNum("1234.12345"))

	// Test numbers with leading zeros in decimal part
	c.Equal("1,234.01", CommaFromStringNum("1234.01"))
	c.Equal("1,234.001", CommaFromStringNum("1234.001"))
	c.Equal("1,234.0001", CommaFromStringNum("1234.0001"))
}

func TestCommaFromStringNumEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test numbers that are exactly divisible by 3 digits
	c.Equal("1,000", CommaFromStringNum("1000"))
	c.Equal("10,000", CommaFromStringNum("10000"))
	c.Equal("100,000", CommaFromStringNum("100000"))
	c.Equal("1,000,000", CommaFromStringNum("1000000"))

	// Test numbers with 1 digit before first comma
	c.Equal("1,000", CommaFromStringNum("1000"))
	c.Equal("1,000,000", CommaFromStringNum("1000000"))
	c.Equal("1,000,000,000", CommaFromStringNum("1000000000"))

	// Test numbers with 2 digits before first comma
	c.Equal("10,000", CommaFromStringNum("10000"))
	c.Equal("12,000", CommaFromStringNum("12000"))
	c.Equal("99,000", CommaFromStringNum("99000"))

	// Test very large numbers
	c.Equal("1,000,000,000,000", CommaFromStringNum("1000000000000"))
	c.Equal("123,456,789,012,345", CommaFromStringNum("123456789012345"))
}

func TestCommaFromStringNumZeroCases(t *testing.T) {
	c := check.New(t)

	// Test various zero representations
	c.Equal("0", CommaFromStringNum("0"))
	c.Equal("-0", CommaFromStringNum("-0"))
	c.Equal("+0", CommaFromStringNum("+0"))
	c.Equal("0.0", CommaFromStringNum("0.0"))
	c.Equal("0.00", CommaFromStringNum("0.00"))
	c.Equal("0.000", CommaFromStringNum("0.000"))
	c.Equal("-0.0", CommaFromStringNum("-0.0"))
	c.Equal("+0.0", CommaFromStringNum("+0.0"))
}
