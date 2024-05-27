// Copyright (c) 2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package paths

import (
	"path/filepath"

	"github.com/richardwilkes/toolbox/cmdline"
)

// AppDataDir returns the application data directory.
func AppDataDir() string {
	path := filepath.Join(HomeDir(), "Library", "Application Support")
	if cmdline.AppIdentifier != "" {
		path = filepath.Join(path, cmdline.AppIdentifier)
	}
	return path
}

// AppLogDir returns the application log directory.
func AppLogDir() string {
	path := filepath.Join(HomeDir(), "Library", "Logs")
	if cmdline.AppIdentifier != "" {
		path = filepath.Join(path, cmdline.AppIdentifier)
	}
	return path
}

// FontDirs returns the standard font directories, in order of priority.
func FontDirs() []string {
	return []string{filepath.Join(HomeDir(), "Library", "Fonts"), "/Library/Fonts", "/System/Library/Fonts"}
}
