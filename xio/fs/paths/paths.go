// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package paths provides platform-specific standard paths.
package paths

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/richardwilkes/toolbox/cmdline"
)

// HomeDir returns the home directory. If this cannot be determined for some reason, "." will be returned instead.
func HomeDir() string {
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	if dir, err := os.UserHomeDir(); err == nil {
		return dir
	}
	return "."
}

// AppLogDir returns the application log directory.
func AppLogDir() string {
	path := HomeDir()
	if path != "." {
		path = addPlatformAppLogSubDirs(path)
	} else {
		path = filepath.Join(path, "logs")
	}
	if cmdline.AppIdentifier != "" {
		path = filepath.Join(path, cmdline.AppIdentifier)
	}
	return path
}

// AppDataDir returns the application data directory.
func AppDataDir() string {
	path := HomeDir()
	if path != "." {
		path = addPlatformAppDataSubDirs(path)
	} else {
		path = filepath.Join(path, "app_data")
	}
	if cmdline.AppIdentifier != "" {
		path = filepath.Join(path, cmdline.AppIdentifier)
	}
	return path
}
