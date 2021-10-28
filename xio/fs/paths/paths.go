// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
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
	"runtime"

	"github.com/richardwilkes/toolbox"
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
		switch runtime.GOOS {
		case toolbox.MacOS:
			path = filepath.Join(path, "Library", "Logs")
		case toolbox.WindowsOS:
			path = filepath.Join(path, "AppData", "Logs")
		default:
			path = filepath.Join(path, ".logs")
		}
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
		switch runtime.GOOS {
		case toolbox.MacOS:
			path = filepath.Join(path, "Library", "Application Support")
		case toolbox.WindowsOS:
			path = filepath.Join(path, "AppData")
		default:
			path = filepath.Join(path, ".appdata")
		}
	} else {
		path = filepath.Join(path, "app_data")
	}
	if cmdline.AppIdentifier != "" {
		path = filepath.Join(path, cmdline.AppIdentifier)
	}
	return path
}

// FontDirs returns the standard font directories, in order of priority.
func FontDirs() []string {
	switch runtime.GOOS {
	case toolbox.MacOS:
		return []string{filepath.Join(HomeDir(), "Library", "Fonts"), "/Library/Fonts", "/System/Library/Fonts"}
	case toolbox.WindowsOS:
		windir := os.Getenv("WINDIR")
		if windir == "" {
			windir = "C:\\Windows"
		}
		return []string{filepath.Join(windir, "Fonts")}
	case toolbox.LinuxOS:
		return []string{filepath.Join(HomeDir(), ".fonts"), "/usr/local/share/fonts", "/usr/share/fonts"}
	default:
		return nil
	}
}
