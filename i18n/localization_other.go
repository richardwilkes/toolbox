// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//go:build !windows

package i18n

import "os"

// Locale returns the value of the LC_ALL environment variable, if set. If not, then it falls back to the value of the
// LANG environment variable. If that is also not set, then it returns "en_US.UTF-8".
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
