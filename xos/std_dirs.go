// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// HomeDir returns the current user's home directory. The returned path uses platform separators.
func HomeDir() string {
	var home string
	switch runtime.GOOS {
	case WindowsOS:
		home = os.Getenv("USERPROFILE")
	default:
		home = os.Getenv("HOME")
	}
	if home == "" {
		home = os.TempDir()
	}
	return home
}

// AppDir returns the absolute path to the logical directory the application resides within. The returned path uses
// platform separators. For macOS, this means the directory where the .app bundle resides, not the binary that's tucked
// down inside it.
func AppDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", errs.Wrap(err)
	}
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return "", errs.Wrap(err)
	}
	if path, err = filepath.Abs(path); err != nil {
		return "", errs.Wrap(err)
	}
	path = filepath.Dir(path)
	if runtime.GOOS == MacOS {
		// Account for macOS bundles
		if i := strings.LastIndex(path, ".app/"); i != -1 {
			path = filepath.Dir(path[:i])
		}
	}
	return path, nil
}

// AppDataDir returns the path to use for user-specific data for the application. The returned path uses platform
// separators and may need to be created before use.
func AppDataDir(withAppIdentifier bool) string {
	var dir string
	switch runtime.GOOS {
	case MacOS:
		dir = filepath.Join(HomeDir(), "Library", "Application Support")
	case WindowsOS:
		if dir = os.Getenv("LOCALAPPDATA"); dir == "" {
			dir = filepath.Join(HomeDir(), "AppData", "Local")
		}
	default:
		if dir = os.Getenv("XDG_DATA_HOME"); dir == "" {
			dir = filepath.Join(HomeDir(), ".local", "share")
		}
	}
	if withAppIdentifier && AppIdentifier != "" {
		dir = filepath.Join(dir, AppIdentifier)
	}
	return dir
}

// AppLogDir returns the application log directory. The returned path uses platform separators and may need to be
// created before use.
func AppLogDir(withAppIdentifier bool) string {
	switch runtime.GOOS {
	case MacOS:
		dir := filepath.Join(HomeDir(), "Library", "Logs")
		if withAppIdentifier && AppIdentifier != "" {
			dir = filepath.Join(dir, AppIdentifier)
		}
		return dir
	default:
		return filepath.Join(AppDataDir(withAppIdentifier), "Logs")
	}
}

// FontDirs returns the standard font directories, in order of priority. The returned paths use platform separators and
// may not exist.
func FontDirs() []string {
	switch runtime.GOOS {
	case MacOS:
		return []string{
			filepath.Join(HomeDir(), "Library", "Fonts"),
			"/Library/Fonts",
			"/System/Library/Fonts",
		}
	case WindowsOS:
		windir := os.Getenv("WINDIR")
		if windir == "" {
			windir = `C:\Windows`
		}
		return []string{
			filepath.Join(AppDataDir(false), "Microsoft", "Windows", "Fonts"),
			filepath.Join(windir, "Fonts"),
		}
	case LinuxOS:
		return []string{
			filepath.Join(HomeDir(), ".fonts"),
			"/usr/local/share/fonts",
			"/usr/share/fonts",
		}
	default:
		return nil
	}
}
