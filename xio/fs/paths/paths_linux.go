// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package paths

import "path/filepath"

func addPlatformAppLogSubDirs(base string) string {
	return filepath.Join(base, ".logs")
}

func addPlatformAppDataSubDirs(base string) string {
	return filepath.Join(base, ".appdata")
}

// FontDirs returns the standard font directories, in order of priority.
func FontDirs() []string {
	return []string{filepath.Join(HomeDir(), ".fonts"), "/usr/local/share/fonts", "/usr/share/fonts"}
}
