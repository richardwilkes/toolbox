// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/v2/i18n"
)

var (
	// AppCmdName holds the base name of the executable, or "<unknown>" if it cannot be determined.
	AppCmdName string
	// AppName holds the name of the application. By default, this is the same as AppCmdName.
	AppName string
	// CopyrightStartYear holds the starting year to place in the copyright banner. If not set explicitly, it is the year
	// of the "vcs.time" build setting, or the current year if that is unavailable.
	CopyrightStartYear string
	// CopyrightEndYear holds the ending year to place in the copyright banner. If not set explicitly, it is the year of
	// the "vcs.time" build setting, or the current year if that is unavailable.
	CopyrightEndYear string
	// CopyrightHolder holds the name of the copyright holder.
	CopyrightHolder string
	// License holds a one-line description of the license the software is distributed under, such as "Mozilla Public
	// License 2.0", not the full license text.
	License string
	// AppVersion holds the application's version. If not set explicitly, it is the main module's version (without a
	// leading 'v'), or "0.0" if that is unavailable. The module version is only available in binaries built with
	// `go install <package>@<version>`.
	AppVersion string
	// VCSName holds the name of the version control system, from the "vcs" build setting.
	VCSName string
	// VCSVersion holds the VCS revision, from the "vcs.revision" build setting.
	VCSVersion string
	// VCSModified is true if the "vcs.modified" build setting is true.
	VCSModified bool
	// BuildNumber holds the build number. If not set explicitly, it is the "vcs.time" build setting (or the current time
	// if that is unavailable) formatted as YYYYMMDDhhmmss.
	BuildNumber string
	// AppIdentifier holds the application's uniform type identifier (UTI), in reverse-DNS form and containing only
	// alphanumeric (A-Z, a-z, 0-9), hyphen (-), and period (.) characters, e.g. "com.ajax.Hello".
	AppIdentifier string
)

func init() {
	if path, err := os.Executable(); err == nil {
		path = filepath.Base(path)
		if path != "." {
			AppCmdName = path
		}
	}
	if AppCmdName == "" {
		AppCmdName = "<unknown>"
	}
	if AppName == "" {
		AppName = AppCmdName
	}
	VCSName = ""
	VCSVersion = ""
	VCSModified = false
	var vcsTime time.Time
	if info, ok := debug.ReadBuildInfo(); ok {
		if AppVersion == "" && info.Main.Version != "(devel)" {
			AppVersion = strings.TrimLeft(info.Main.Version, "v")
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs":
				VCSName = setting.Value
			case "vcs.revision":
				VCSVersion = setting.Value
			case "vcs.time":
				if t, err := time.Parse(time.RFC3339, setting.Value); err == nil {
					vcsTime = t
				}
			case "vcs.modified":
				if setting.Value == "true" {
					VCSModified = true
				}
			}
		}
	}
	if AppVersion == "" {
		AppVersion = "0.0"
	}
	if vcsTime.IsZero() {
		vcsTime = time.Now()
	}
	if BuildNumber == "" {
		BuildNumber = vcsTime.Format("20060102150405")
	}
	year := strconv.Itoa(vcsTime.Year())
	if CopyrightStartYear == "" {
		CopyrightStartYear = year
	}
	if CopyrightEndYear == "" {
		CopyrightEndYear = year
	}
}

// ShortAppVersion returns AppVersion, with a trailing '~' if VCSModified is true.
func ShortAppVersion() string {
	return markAppVersionModified(AppVersion)
}

// LongAppVersion returns AppVersion and BuildNumber joined by '-' (or just AppVersion if BuildNumber is empty), with a
// trailing '~' if VCSModified is true.
func LongAppVersion() string {
	version := AppVersion
	if BuildNumber != "" {
		version += "-" + BuildNumber
	}
	return markAppVersionModified(version)
}

func markAppVersionModified(in string) string {
	if VCSModified {
		return in + "~"
	}
	return in
}

// CopyrightYears returns the copyright years, either as a single year or as a range of years, e.g. "2025" or
// "2016-2025".
func CopyrightYears() string {
	years := CopyrightStartYear
	if CopyrightEndYear != "" && CopyrightEndYear != CopyrightStartYear {
		if years == "" {
			years = CopyrightEndYear
		} else {
			years += "-" + CopyrightEndYear
		}
	}
	return years
}

// Copyright returns the copyright notice. If no copyright years have been set, an empty string will be returned.
func Copyright() string {
	var buf strings.Builder
	years := CopyrightYears()
	if years != "" {
		buf.WriteString(i18n.Text("Copyright ©"))
		buf.WriteString(years)
		if CopyrightHolder != "" {
			buf.WriteString(i18n.Text(" by "))
			buf.WriteString(CopyrightHolder)
		}
		buf.WriteString(i18n.Text(". All rights reserved."))
	}
	return buf.String()
}
