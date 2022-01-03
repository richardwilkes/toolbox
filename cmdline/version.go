// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

import "strings"

// ShortVersion returns the app version. If AppVersion has not been set, then "0.0" will be returned instead.
func ShortVersion() string {
	if AppVersion == "" {
		return "0.0"
	}
	return AppVersion
}

// LongVersion returns a combination of the app version and the build number. If AppVersion has not been set, then "0.0"
// will be used instead.
func LongVersion() string {
	version := ShortVersion()
	if BuildNumber != "" {
		if !strings.HasSuffix(version, "~") {
			version += "-"
		}
		version += BuildNumber
	}
	return version
}
