// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestCommaInt(t *testing.T) {
	c := check.New(t)

	// Test integers
	c.Equal("0", xstrings.CommaInt(0))
	c.Equal("1", xstrings.CommaInt(1))
	c.Equal("12", xstrings.CommaInt(12))
	c.Equal("123", xstrings.CommaInt(123))
	c.Equal("1,234", xstrings.CommaInt(1234))
	c.Equal("12,345", xstrings.CommaInt(12345))
	c.Equal("123,456", xstrings.CommaInt(123456))
	c.Equal("1,234,567", xstrings.CommaInt(1234567))
	c.Equal("12,345,678", xstrings.CommaInt(12345678))
	c.Equal("123,456,789", xstrings.CommaInt(123456789))
	c.Equal("1,234,567,890", xstrings.CommaInt(1234567890))

	// Test negative integers
	c.Equal("-1", xstrings.CommaInt(-1))
	c.Equal("-12", xstrings.CommaInt(-12))
	c.Equal("-123", xstrings.CommaInt(-123))
	c.Equal("-1,234", xstrings.CommaInt(-1234))
	c.Equal("-12,345", xstrings.CommaInt(-12345))
	c.Equal("-123,456", xstrings.CommaInt(-123456))
	c.Equal("-1,234,567", xstrings.CommaInt(-1234567))

	// Test different numeric types
	c.Equal("42", xstrings.CommaInt(int8(42)))
	c.Equal("1,234", xstrings.CommaInt(int16(1234)))
	c.Equal("123,456", xstrings.CommaInt(int32(123456)))
	c.Equal("1,234,567", xstrings.CommaInt(int64(1234567)))
	c.Equal("255", xstrings.CommaInt(uint8(255)))
	c.Equal("65,535", xstrings.CommaInt(uint16(65535)))
	c.Equal("4,294,967,295", xstrings.CommaInt(uint32(4294967295)))
}

func TestCommaFloat(t *testing.T) {
	c := check.New(t)

	// Test integers
	c.Equal("0", xstrings.CommaFloat(0.0))
	c.Equal("1", xstrings.CommaFloat(1.0))
	c.Equal("12", xstrings.CommaFloat(12.0))
	c.Equal("123", xstrings.CommaFloat(123.0))
	c.Equal("1,234", xstrings.CommaFloat(1234.0))
	c.Equal("12,345", xstrings.CommaFloat(12345.0))
	c.Equal("123,456", xstrings.CommaFloat(123456.0))
	c.Equal("1,234,567", xstrings.CommaFloat(1234567.0))
	c.Equal("12,345,678", xstrings.CommaFloat(12345678.0))
	c.Equal("123,456,789", xstrings.CommaFloat(123456789.0))
	c.Equal("1,234,567,890", xstrings.CommaFloat(1234567890.0))

	// Test negative integers
	c.Equal("-1", xstrings.CommaFloat(-1.0))
	c.Equal("-12", xstrings.CommaFloat(-12.0))
	c.Equal("-123", xstrings.CommaFloat(-123.0))
	c.Equal("-1,234", xstrings.CommaFloat(-1234.0))
	c.Equal("-12,345", xstrings.CommaFloat(-12345.0))
	c.Equal("-123,456", xstrings.CommaFloat(-123456.0))
	c.Equal("-1,234,567", xstrings.CommaFloat(-1234567.0))

	// Test floating point numbers
	c.Equal("1.5", xstrings.CommaFloat(1.5))
	c.Equal("12.34", xstrings.CommaFloat(12.34))
	c.Equal("123.456", xstrings.CommaFloat(123.456))
	c.Equal("1,234.56", xstrings.CommaFloat(1234.56))
	c.Equal("12,345.678", xstrings.CommaFloat(12345.678))
	c.Equal("123,456.789", xstrings.CommaFloat(123456.789))
	c.Equal("1,234,567.89", xstrings.CommaFloat(1234567.89))

	// Test negative floating point numbers
	c.Equal("-1.5", xstrings.CommaFloat(-1.5))
	c.Equal("-12.34", xstrings.CommaFloat(-12.34))
	c.Equal("-123.456", xstrings.CommaFloat(-123.456))
	c.Equal("-1,234.56", xstrings.CommaFloat(-1234.56))
	c.Equal("-1,234.567", xstrings.CommaFloat(-1234.567))
	c.Equal("-12,345.678", xstrings.CommaFloat(-12345.678))

	// Test float32
	c.Equal("42.5", xstrings.CommaFloat(float32(42.5)))
}

func TestCommaFromStringNum(t *testing.T) {
	c := check.New(t)

	// Test empty string
	c.Equal("", xstrings.CommaFromStringNum(""))

	// Test single digits
	c.Equal("0", xstrings.CommaFromStringNum("0"))
	c.Equal("1", xstrings.CommaFromStringNum("1"))
	c.Equal("9", xstrings.CommaFromStringNum("9"))

	// Test two digits
	c.Equal("10", xstrings.CommaFromStringNum("10"))
	c.Equal("99", xstrings.CommaFromStringNum("99"))

	// Test three digits
	c.Equal("100", xstrings.CommaFromStringNum("100"))
	c.Equal("123", xstrings.CommaFromStringNum("123"))
	c.Equal("999", xstrings.CommaFromStringNum("999"))

	// Test four digits (first comma)
	c.Equal("1,000", xstrings.CommaFromStringNum("1000"))
	c.Equal("1,234", xstrings.CommaFromStringNum("1234"))
	c.Equal("9,999", xstrings.CommaFromStringNum("9999"))

	// Test five digits
	c.Equal("10,000", xstrings.CommaFromStringNum("10000"))
	c.Equal("12,345", xstrings.CommaFromStringNum("12345"))
	c.Equal("99,999", xstrings.CommaFromStringNum("99999"))

	// Test six digits
	c.Equal("100,000", xstrings.CommaFromStringNum("100000"))
	c.Equal("123,456", xstrings.CommaFromStringNum("123456"))
	c.Equal("999,999", xstrings.CommaFromStringNum("999999"))

	// Test seven digits (second comma)
	c.Equal("1,000,000", xstrings.CommaFromStringNum("1000000"))
	c.Equal("1,234,567", xstrings.CommaFromStringNum("1234567"))
	c.Equal("9,999,999", xstrings.CommaFromStringNum("9999999"))

	// Test larger numbers
	c.Equal("12,345,678", xstrings.CommaFromStringNum("12345678"))
	c.Equal("123,456,789", xstrings.CommaFromStringNum("123456789"))
	c.Equal("1,234,567,890", xstrings.CommaFromStringNum("1234567890"))
	c.Equal("12,345,678,901", xstrings.CommaFromStringNum("12345678901"))
	c.Equal("123,456,789,012", xstrings.CommaFromStringNum("123456789012"))

	// Test negative numbers
	c.Equal("-1", xstrings.CommaFromStringNum("-1"))
	c.Equal("-12", xstrings.CommaFromStringNum("-12"))
	c.Equal("-123", xstrings.CommaFromStringNum("-123"))
	c.Equal("-1,234", xstrings.CommaFromStringNum("-1234"))
	c.Equal("-12,345", xstrings.CommaFromStringNum("-12345"))
	c.Equal("-123,456", xstrings.CommaFromStringNum("-123456"))
	c.Equal("-1,234,567", xstrings.CommaFromStringNum("-1234567"))
	c.Equal("-12,345,678", xstrings.CommaFromStringNum("-12345678"))

	// Test positive numbers with explicit sign
	c.Equal("+1", xstrings.CommaFromStringNum("+1"))
	c.Equal("+12", xstrings.CommaFromStringNum("+12"))
	c.Equal("+123", xstrings.CommaFromStringNum("+123"))
	c.Equal("+1,234", xstrings.CommaFromStringNum("+1234"))
	c.Equal("+12,345", xstrings.CommaFromStringNum("+12345"))
	c.Equal("+123,456", xstrings.CommaFromStringNum("+123456"))
	c.Equal("+1,234,567", xstrings.CommaFromStringNum("+1234567"))

	// Test decimal numbers
	c.Equal("1.0", xstrings.CommaFromStringNum("1.0"))
	c.Equal("12.34", xstrings.CommaFromStringNum("12.34"))
	c.Equal("123.456", xstrings.CommaFromStringNum("123.456"))
	c.Equal("1,234.56", xstrings.CommaFromStringNum("1234.56"))
	c.Equal("12,345.678", xstrings.CommaFromStringNum("12345.678"))
	c.Equal("123,456.789", xstrings.CommaFromStringNum("123456.789"))
	c.Equal("1,234,567.89", xstrings.CommaFromStringNum("1234567.89"))

	// Test negative decimal numbers
	c.Equal("-1.0", xstrings.CommaFromStringNum("-1.0"))
	c.Equal("-12.34", xstrings.CommaFromStringNum("-12.34"))
	c.Equal("-123.456", xstrings.CommaFromStringNum("-123.456"))
	c.Equal("-1,234.56", xstrings.CommaFromStringNum("-1234.56"))
	c.Equal("-12,345.678", xstrings.CommaFromStringNum("-12345.678"))

	// Test positive decimal numbers with explicit sign
	c.Equal("+1.0", xstrings.CommaFromStringNum("+1.0"))
	c.Equal("+12.34", xstrings.CommaFromStringNum("+12.34"))
	c.Equal("+1,234.56", xstrings.CommaFromStringNum("+1234.56"))

	// Test decimal numbers with different decimal place counts
	c.Equal("1,234.1", xstrings.CommaFromStringNum("1234.1"))
	c.Equal("1,234.12", xstrings.CommaFromStringNum("1234.12"))
	c.Equal("1,234.123", xstrings.CommaFromStringNum("1234.123"))
	c.Equal("1,234.1234", xstrings.CommaFromStringNum("1234.1234"))
	c.Equal("1,234.12345", xstrings.CommaFromStringNum("1234.12345"))

	// Test numbers with leading zeros in decimal part
	c.Equal("1,234.01", xstrings.CommaFromStringNum("1234.01"))
	c.Equal("1,234.001", xstrings.CommaFromStringNum("1234.001"))
	c.Equal("1,234.0001", xstrings.CommaFromStringNum("1234.0001"))

	// Test numbers that are exactly divisible by 3 digits
	c.Equal("1,000", xstrings.CommaFromStringNum("1000"))
	c.Equal("10,000", xstrings.CommaFromStringNum("10000"))
	c.Equal("100,000", xstrings.CommaFromStringNum("100000"))
	c.Equal("1,000,000", xstrings.CommaFromStringNum("1000000"))

	// Test numbers with 1 digit before first comma
	c.Equal("1,000", xstrings.CommaFromStringNum("1000"))
	c.Equal("1,000,000", xstrings.CommaFromStringNum("1000000"))
	c.Equal("1,000,000,000", xstrings.CommaFromStringNum("1000000000"))

	// Test numbers with 2 digits before first comma
	c.Equal("10,000", xstrings.CommaFromStringNum("10000"))
	c.Equal("12,000", xstrings.CommaFromStringNum("12000"))
	c.Equal("99,000", xstrings.CommaFromStringNum("99000"))

	// Test very large numbers
	c.Equal("1,000,000,000,000", xstrings.CommaFromStringNum("1000000000000"))
	c.Equal("123,456,789,012,345", xstrings.CommaFromStringNum("123456789012345"))

	// Test various zero representations
	c.Equal("0", xstrings.CommaFromStringNum("0"))
	c.Equal("-0", xstrings.CommaFromStringNum("-0"))
	c.Equal("+0", xstrings.CommaFromStringNum("+0"))
	c.Equal("0.0", xstrings.CommaFromStringNum("0.0"))
	c.Equal("0.00", xstrings.CommaFromStringNum("0.00"))
	c.Equal("0.000", xstrings.CommaFromStringNum("0.000"))
	c.Equal("-0.0", xstrings.CommaFromStringNum("-0.0"))
	c.Equal("+0.0", xstrings.CommaFromStringNum("+0.0"))
}
