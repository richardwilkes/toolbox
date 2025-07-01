// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package paths

import (
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/v2/xflag"
)

// AppDataDir returns the application data directory.
func AppDataDir() string {
	path := filepath.Join(HomeDir(), "AppData", "Local")
	if xflag.AppIdentifier != "" {
		path = filepath.Join(path, xflag.AppIdentifier)
	}
	return path
}

// AppLogDir returns the application log directory.
func AppLogDir() string {
	return filepath.Join(AppDataDir(), "Logs")
}

// FontDirs returns the standard font directories, in order of priority.
func FontDirs() []string {
	windir := os.Getenv("WINDIR")
	if windir == "" {
		windir = `C:\Windows`
	}
	return []string{filepath.Join(windir, "Fonts")}
}
