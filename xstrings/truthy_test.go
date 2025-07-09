package xstrings_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestIsTruthy(t *testing.T) {
	c := check.New(t)

	// Test truthy values - lowercase
	c.True(xstrings.IsTruthy("1"))
	c.True(xstrings.IsTruthy("true"))
	c.True(xstrings.IsTruthy("yes"))
	c.True(xstrings.IsTruthy("on"))

	// Test truthy values - uppercase
	c.True(xstrings.IsTruthy("TRUE"))
	c.True(xstrings.IsTruthy("YES"))
	c.True(xstrings.IsTruthy("ON"))

	// Test truthy values - mixed case
	c.True(xstrings.IsTruthy("True"))
	c.True(xstrings.IsTruthy("TRUE"))
	c.True(xstrings.IsTruthy("tRuE"))
	c.True(xstrings.IsTruthy("Yes"))
	c.True(xstrings.IsTruthy("YES"))
	c.True(xstrings.IsTruthy("yEs"))
	c.True(xstrings.IsTruthy("On"))
	c.True(xstrings.IsTruthy("ON"))
	c.True(xstrings.IsTruthy("oN"))

	// Test falsy values - common false representations
	c.False(xstrings.IsTruthy("0"))
	c.False(xstrings.IsTruthy("false"))
	c.False(xstrings.IsTruthy("no"))
	c.False(xstrings.IsTruthy("off"))
	c.False(xstrings.IsTruthy("FALSE"))
	c.False(xstrings.IsTruthy("NO"))
	c.False(xstrings.IsTruthy("OFF"))

	// Test empty string
	c.False(xstrings.IsTruthy(""))

	// Test whitespace strings
	c.False(xstrings.IsTruthy(" "))
	c.False(xstrings.IsTruthy("  "))
	c.False(xstrings.IsTruthy("\t"))
	c.False(xstrings.IsTruthy("\n"))
	c.False(xstrings.IsTruthy("\r"))
	c.False(xstrings.IsTruthy(" \t\n\r "))

	// Test truthy values with leading/trailing whitespace (should be false)
	c.False(xstrings.IsTruthy(" 1"))
	c.False(xstrings.IsTruthy("1 "))
	c.False(xstrings.IsTruthy(" 1 "))
	c.False(xstrings.IsTruthy(" true"))
	c.False(xstrings.IsTruthy("true "))
	c.False(xstrings.IsTruthy(" true "))
	c.False(xstrings.IsTruthy(" yes"))
	c.False(xstrings.IsTruthy("yes "))
	c.False(xstrings.IsTruthy(" yes "))
	c.False(xstrings.IsTruthy(" on"))
	c.False(xstrings.IsTruthy("on "))
	c.False(xstrings.IsTruthy(" on "))

	// Test similar but non-truthy values
	c.False(xstrings.IsTruthy("11"))
	c.False(xstrings.IsTruthy("10"))
	c.False(xstrings.IsTruthy("01"))
	c.False(xstrings.IsTruthy("2"))
	c.False(xstrings.IsTruthy("-1"))
	c.False(xstrings.IsTruthy("1.0"))
	c.False(xstrings.IsTruthy("1.1"))

	// Test partial matches
	c.False(xstrings.IsTruthy("t"))
	c.False(xstrings.IsTruthy("tr"))
	c.False(xstrings.IsTruthy("tru"))
	c.False(xstrings.IsTruthy("truee"))
	c.False(xstrings.IsTruthy("trues"))
	c.False(xstrings.IsTruthy("y"))
	c.False(xstrings.IsTruthy("ye"))
	c.False(xstrings.IsTruthy("yess"))
	c.False(xstrings.IsTruthy("o"))
	c.False(xstrings.IsTruthy("onn"))

	// Test words containing truthy values
	c.False(xstrings.IsTruthy("true1"))
	c.False(xstrings.IsTruthy("1true"))
	c.False(xstrings.IsTruthy("atrue"))
	c.False(xstrings.IsTruthy("truea"))
	c.False(xstrings.IsTruthy("yes1"))
	c.False(xstrings.IsTruthy("1yes"))
	c.False(xstrings.IsTruthy("ayes"))
	c.False(xstrings.IsTruthy("yesa"))
	c.False(xstrings.IsTruthy("on1"))
	c.False(xstrings.IsTruthy("1on"))
	c.False(xstrings.IsTruthy("aon"))
	c.False(xstrings.IsTruthy("ona"))

	// Test random strings
	c.False(xstrings.IsTruthy("hello"))
	c.False(xstrings.IsTruthy("world"))
	c.False(xstrings.IsTruthy("test"))
	c.False(xstrings.IsTruthy("abc"))
	c.False(xstrings.IsTruthy("xyz"))
	c.False(xstrings.IsTruthy("123"))
	c.False(xstrings.IsTruthy("hello world"))

	// Test special characters
	c.False(xstrings.IsTruthy("@"))
	c.False(xstrings.IsTruthy("#"))
	c.False(xstrings.IsTruthy("$"))
	c.False(xstrings.IsTruthy("%"))
	c.False(xstrings.IsTruthy("&"))
	c.False(xstrings.IsTruthy("*"))
	c.False(xstrings.IsTruthy("!"))
	c.False(xstrings.IsTruthy("?"))
	c.False(xstrings.IsTruthy("."))
	c.False(xstrings.IsTruthy(","))
	c.False(xstrings.IsTruthy(";"))
	c.False(xstrings.IsTruthy(":"))

	// Test Unicode characters
	c.False(xstrings.IsTruthy("ñ"))
	c.False(xstrings.IsTruthy("ü"))
	c.False(xstrings.IsTruthy("café"))
	c.False(xstrings.IsTruthy("北京"))
	c.False(xstrings.IsTruthy("🚀"))
	c.False(xstrings.IsTruthy("🎉"))
	c.False(xstrings.IsTruthy("😀"))

	// Test mathematical symbols
	c.False(xstrings.IsTruthy("∑"))
	c.False(xstrings.IsTruthy("∞"))
	c.False(xstrings.IsTruthy("π"))
	c.False(xstrings.IsTruthy("α"))
	c.False(xstrings.IsTruthy("β"))

	// Test currency symbols
	c.False(xstrings.IsTruthy("$"))
	c.False(xstrings.IsTruthy("€"))
	c.False(xstrings.IsTruthy("£"))
	c.False(xstrings.IsTruthy("¥"))

	// Test programming-related strings
	c.False(xstrings.IsTruthy("null"))
	c.False(xstrings.IsTruthy("undefined"))
	c.False(xstrings.IsTruthy("nil"))
	c.False(xstrings.IsTruthy("void"))
	c.False(xstrings.IsTruthy("NaN"))

	// Test other common boolean-like strings that are not truthy
	c.False(xstrings.IsTruthy("enable"))
	c.False(xstrings.IsTruthy("enabled"))
	c.False(xstrings.IsTruthy("disable"))
	c.False(xstrings.IsTruthy("disabled"))
	c.False(xstrings.IsTruthy("active"))
	c.False(xstrings.IsTruthy("inactive"))
	c.False(xstrings.IsTruthy("ok"))
	c.False(xstrings.IsTruthy("OK"))
	c.False(xstrings.IsTruthy("fail"))
	c.False(xstrings.IsTruthy("pass"))

	// Test numeric strings that might be confused with "1"
	c.False(xstrings.IsTruthy("01"))
	c.False(xstrings.IsTruthy("10"))
	c.False(xstrings.IsTruthy("21"))
	c.False(xstrings.IsTruthy("12"))
	c.False(xstrings.IsTruthy("100"))
	c.False(xstrings.IsTruthy("001"))

	// Test strings that contain truthy substrings but are not exact matches
	c.False(xstrings.IsTruthy("trueish"))
	c.False(xstrings.IsTruthy("yesno"))
	c.False(xstrings.IsTruthy("onoff"))
	c.False(xstrings.IsTruthy("1and1"))
	c.False(xstrings.IsTruthy("nottrue"))
	c.False(xstrings.IsTruthy("notyes"))
	c.False(xstrings.IsTruthy("noton"))

	// Test mixed alphanumeric
	c.False(xstrings.IsTruthy("abc123"))
	c.False(xstrings.IsTruthy("123abc"))
	c.False(xstrings.IsTruthy("a1b2c3"))
	c.False(xstrings.IsTruthy("test123"))
	c.False(xstrings.IsTruthy("123test"))

	// Test case sensitivity is properly handled
	c.True(xstrings.IsTruthy("tRuE"))
	c.True(xstrings.IsTruthy("yEs"))
	c.True(xstrings.IsTruthy("oN"))
	c.True(xstrings.IsTruthy("TRUE"))
	c.True(xstrings.IsTruthy("YES"))
	c.True(xstrings.IsTruthy("ON"))

	// Test that function only accepts exact matches
	c.False(xstrings.IsTruthy("truthy"))
	c.False(xstrings.IsTruthy("truth"))
	c.False(xstrings.IsTruthy("truly"))
	c.False(xstrings.IsTruthy("yesyes"))
	c.False(xstrings.IsTruthy("onon"))
	c.False(xstrings.IsTruthy("11"))

	// Test that whitespace is not trimmed
	c.False(xstrings.IsTruthy(" true"))
	c.False(xstrings.IsTruthy("true "))
	c.False(xstrings.IsTruthy(" yes"))
	c.False(xstrings.IsTruthy("yes "))
	c.False(xstrings.IsTruthy(" on"))
	c.False(xstrings.IsTruthy("on "))
	c.False(xstrings.IsTruthy(" 1"))
	c.False(xstrings.IsTruthy("1 "))

	// Test tab and newline characters
	c.False(xstrings.IsTruthy("true\t"))
	c.False(xstrings.IsTruthy("\ttrue"))
	c.False(xstrings.IsTruthy("true\n"))
	c.False(xstrings.IsTruthy("\ntrue"))
	c.False(xstrings.IsTruthy("true\r"))
	c.False(xstrings.IsTruthy("\rtrue"))

	// Test empty-like strings
	c.False(xstrings.IsTruthy(""))
	c.False(xstrings.IsTruthy(" "))
	c.False(xstrings.IsTruthy("\t"))
	c.False(xstrings.IsTruthy("\n"))
	c.False(xstrings.IsTruthy("\r"))
	c.False(xstrings.IsTruthy("\r\n"))

	// Test strings that might look like boolean values but aren't supported
	c.False(xstrings.IsTruthy("T"))
	c.False(xstrings.IsTruthy("F"))
	c.False(xstrings.IsTruthy("Y"))
	c.False(xstrings.IsTruthy("N"))
	c.False(xstrings.IsTruthy("t"))
	c.False(xstrings.IsTruthy("f"))
	c.False(xstrings.IsTruthy("y"))
	c.False(xstrings.IsTruthy("n"))
}
