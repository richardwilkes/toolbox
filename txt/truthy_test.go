package txt

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestIsTruthy(t *testing.T) {
	c := check.New(t)

	// Test truthy values - lowercase
	c.True(IsTruthy("1"))
	c.True(IsTruthy("true"))
	c.True(IsTruthy("yes"))
	c.True(IsTruthy("on"))

	// Test truthy values - uppercase
	c.True(IsTruthy("TRUE"))
	c.True(IsTruthy("YES"))
	c.True(IsTruthy("ON"))

	// Test truthy values - mixed case
	c.True(IsTruthy("True"))
	c.True(IsTruthy("TRUE"))
	c.True(IsTruthy("tRuE"))
	c.True(IsTruthy("Yes"))
	c.True(IsTruthy("YES"))
	c.True(IsTruthy("yEs"))
	c.True(IsTruthy("On"))
	c.True(IsTruthy("ON"))
	c.True(IsTruthy("oN"))

	// Test falsy values - common false representations
	c.False(IsTruthy("0"))
	c.False(IsTruthy("false"))
	c.False(IsTruthy("no"))
	c.False(IsTruthy("off"))
	c.False(IsTruthy("FALSE"))
	c.False(IsTruthy("NO"))
	c.False(IsTruthy("OFF"))

	// Test empty string
	c.False(IsTruthy(""))

	// Test whitespace strings
	c.False(IsTruthy(" "))
	c.False(IsTruthy("  "))
	c.False(IsTruthy("\t"))
	c.False(IsTruthy("\n"))
	c.False(IsTruthy("\r"))
	c.False(IsTruthy(" \t\n\r "))

	// Test truthy values with leading/trailing whitespace (should be false)
	c.False(IsTruthy(" 1"))
	c.False(IsTruthy("1 "))
	c.False(IsTruthy(" 1 "))
	c.False(IsTruthy(" true"))
	c.False(IsTruthy("true "))
	c.False(IsTruthy(" true "))
	c.False(IsTruthy(" yes"))
	c.False(IsTruthy("yes "))
	c.False(IsTruthy(" yes "))
	c.False(IsTruthy(" on"))
	c.False(IsTruthy("on "))
	c.False(IsTruthy(" on "))

	// Test similar but non-truthy values
	c.False(IsTruthy("11"))
	c.False(IsTruthy("10"))
	c.False(IsTruthy("01"))
	c.False(IsTruthy("2"))
	c.False(IsTruthy("-1"))
	c.False(IsTruthy("1.0"))
	c.False(IsTruthy("1.1"))

	// Test partial matches
	c.False(IsTruthy("t"))
	c.False(IsTruthy("tr"))
	c.False(IsTruthy("tru"))
	c.False(IsTruthy("truee"))
	c.False(IsTruthy("trues"))
	c.False(IsTruthy("y"))
	c.False(IsTruthy("ye"))
	c.False(IsTruthy("yess"))
	c.False(IsTruthy("o"))
	c.False(IsTruthy("onn"))

	// Test words containing truthy values
	c.False(IsTruthy("true1"))
	c.False(IsTruthy("1true"))
	c.False(IsTruthy("atrue"))
	c.False(IsTruthy("truea"))
	c.False(IsTruthy("yes1"))
	c.False(IsTruthy("1yes"))
	c.False(IsTruthy("ayes"))
	c.False(IsTruthy("yesa"))
	c.False(IsTruthy("on1"))
	c.False(IsTruthy("1on"))
	c.False(IsTruthy("aon"))
	c.False(IsTruthy("ona"))

	// Test random strings
	c.False(IsTruthy("hello"))
	c.False(IsTruthy("world"))
	c.False(IsTruthy("test"))
	c.False(IsTruthy("abc"))
	c.False(IsTruthy("xyz"))
	c.False(IsTruthy("123"))
	c.False(IsTruthy("hello world"))

	// Test special characters
	c.False(IsTruthy("@"))
	c.False(IsTruthy("#"))
	c.False(IsTruthy("$"))
	c.False(IsTruthy("%"))
	c.False(IsTruthy("&"))
	c.False(IsTruthy("*"))
	c.False(IsTruthy("!"))
	c.False(IsTruthy("?"))
	c.False(IsTruthy("."))
	c.False(IsTruthy(","))
	c.False(IsTruthy(";"))
	c.False(IsTruthy(":"))

	// Test Unicode characters
	c.False(IsTruthy("ñ"))
	c.False(IsTruthy("ü"))
	c.False(IsTruthy("café"))
	c.False(IsTruthy("北京"))
	c.False(IsTruthy("🚀"))
	c.False(IsTruthy("🎉"))
	c.False(IsTruthy("😀"))

	// Test mathematical symbols
	c.False(IsTruthy("∑"))
	c.False(IsTruthy("∞"))
	c.False(IsTruthy("π"))
	c.False(IsTruthy("α"))
	c.False(IsTruthy("β"))

	// Test currency symbols
	c.False(IsTruthy("$"))
	c.False(IsTruthy("€"))
	c.False(IsTruthy("£"))
	c.False(IsTruthy("¥"))

	// Test programming-related strings
	c.False(IsTruthy("null"))
	c.False(IsTruthy("undefined"))
	c.False(IsTruthy("nil"))
	c.False(IsTruthy("void"))
	c.False(IsTruthy("NaN"))

	// Test other common boolean-like strings that are not truthy
	c.False(IsTruthy("enable"))
	c.False(IsTruthy("enabled"))
	c.False(IsTruthy("disable"))
	c.False(IsTruthy("disabled"))
	c.False(IsTruthy("active"))
	c.False(IsTruthy("inactive"))
	c.False(IsTruthy("ok"))
	c.False(IsTruthy("OK"))
	c.False(IsTruthy("fail"))
	c.False(IsTruthy("pass"))

	// Test numeric strings that might be confused with "1"
	c.False(IsTruthy("01"))
	c.False(IsTruthy("10"))
	c.False(IsTruthy("21"))
	c.False(IsTruthy("12"))
	c.False(IsTruthy("100"))
	c.False(IsTruthy("001"))

	// Test strings that contain truthy substrings but are not exact matches
	c.False(IsTruthy("trueish"))
	c.False(IsTruthy("yesno"))
	c.False(IsTruthy("onoff"))
	c.False(IsTruthy("1and1"))
	c.False(IsTruthy("nottrue"))
	c.False(IsTruthy("notyes"))
	c.False(IsTruthy("noton"))

	// Test mixed alphanumeric
	c.False(IsTruthy("abc123"))
	c.False(IsTruthy("123abc"))
	c.False(IsTruthy("a1b2c3"))
	c.False(IsTruthy("test123"))
	c.False(IsTruthy("123test"))
}

func TestIsTruthyEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test case sensitivity is properly handled
	c.True(IsTruthy("tRuE"))
	c.True(IsTruthy("yEs"))
	c.True(IsTruthy("oN"))
	c.True(IsTruthy("TRUE"))
	c.True(IsTruthy("YES"))
	c.True(IsTruthy("ON"))

	// Test that function only accepts exact matches
	c.False(IsTruthy("truthy"))
	c.False(IsTruthy("truth"))
	c.False(IsTruthy("truly"))
	c.False(IsTruthy("yesyes"))
	c.False(IsTruthy("onon"))
	c.False(IsTruthy("11"))

	// Test that whitespace is not trimmed
	c.False(IsTruthy(" true"))
	c.False(IsTruthy("true "))
	c.False(IsTruthy(" yes"))
	c.False(IsTruthy("yes "))
	c.False(IsTruthy(" on"))
	c.False(IsTruthy("on "))
	c.False(IsTruthy(" 1"))
	c.False(IsTruthy("1 "))

	// Test tab and newline characters
	c.False(IsTruthy("true\t"))
	c.False(IsTruthy("\ttrue"))
	c.False(IsTruthy("true\n"))
	c.False(IsTruthy("\ntrue"))
	c.False(IsTruthy("true\r"))
	c.False(IsTruthy("\rtrue"))

	// Test empty-like strings
	c.False(IsTruthy(""))
	c.False(IsTruthy(" "))
	c.False(IsTruthy("\t"))
	c.False(IsTruthy("\n"))
	c.False(IsTruthy("\r"))
	c.False(IsTruthy("\r\n"))

	// Test strings that might look like boolean values but aren't supported
	c.False(IsTruthy("T"))
	c.False(IsTruthy("F"))
	c.False(IsTruthy("Y"))
	c.False(IsTruthy("N"))
	c.False(IsTruthy("t"))
	c.False(IsTruthy("f"))
	c.False(IsTruthy("y"))
	c.False(IsTruthy("n"))
}
