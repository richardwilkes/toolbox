package txt

import "strings"

// IsTruthy returns true for "truthy" values, i.e. ones that should be
// interpreted as true.
func IsTruthy(in string) bool {
	in = strings.ToLower(in)
	return in == "1" || in == "true" || in == "yes" || in == "on"
}
