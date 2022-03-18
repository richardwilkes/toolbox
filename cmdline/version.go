// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

// ShortVersion returns the app version.
func ShortVersion() string {
	if VCSModified {
		return AppVersion + "~"
	}
	return AppVersion
}

// LongVersion returns a combination of the app version and the build number.
func LongVersion() string {
	version := AppVersion
	if BuildNumber != "" {
		version += "-" + BuildNumber
	}
	if VCSModified {
		return version + "~"
	}
	return version
}
