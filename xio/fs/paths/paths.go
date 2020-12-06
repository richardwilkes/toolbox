// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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

// AppLogDir returns the application log directory.
func AppLogDir() string {
	var path string
	if u, err := user.Current(); err == nil {
		path = u.HomeDir
		switch runtime.GOOS {
		case toolbox.MacOS:
			path = filepath.Join(path, "Library", "Logs")
		case toolbox.WindowsOS:
			path = filepath.Join(path, "AppData")
		default:
			path = filepath.Join(path, ".logs")
		}
		if cmdline.AppIdentifier != "" {
			path = filepath.Join(path, cmdline.AppIdentifier)
		}
	}
	return path
}

// AppDataDir returns the application data directory.
func AppDataDir() string {
	var path string
	if u, err := user.Current(); err == nil {
		path = u.HomeDir
		switch runtime.GOOS {
		case toolbox.MacOS:
			path = filepath.Join(path, "Library", "Application Support")
		case toolbox.WindowsOS:
			path = filepath.Join(path, "AppData")
		default:
			path = filepath.Join(path, ".appdata")
		}
		if cmdline.AppIdentifier != "" {
			path = filepath.Join(path, cmdline.AppIdentifier)
		}
	}
	return path
}

// FontDirs returns the standard font directories, in order of priority.
func FontDirs() []string {
	var paths []string
	switch runtime.GOOS {
	case toolbox.MacOS:
		if u, err := user.Current(); err == nil {
			paths = append(paths, filepath.Join(u.HomeDir, "Library", "Fonts"))
		}
		paths = append(paths, "/Library/Fonts", "/System/Library/Fonts")
	case toolbox.WindowsOS:
		windir := os.Getenv("WINDIR")
		if windir == "" {
			windir = "C:\\Windows"
		}
		paths = append(paths, filepath.Join(windir, "Fonts"))
	case toolbox.LinuxOS:
		if u, err := user.Current(); err == nil {
			paths = append(paths, filepath.Join(u.HomeDir, ".fonts"))
		}
		paths = append(paths, "/usr/local/share/fonts", "/usr/share/fonts")
	default:
	}
	return paths
}
