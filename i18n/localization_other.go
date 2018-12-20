// +build !windows

package i18n

import "os"

// Locale returns the value of the LC_ALL environment variable, if set. If
// not, then it falls back to the value of the LANG environment variable. If
// that is also not set, then it returns "en_US.UTF-8".
func Locale() string {
	locale := os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LANG")
		if locale == "" {
			locale = "en_US.UTF-8"
		}
	}
	return locale
}
