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
	// AppCmdName holds the application's name as specified on the command line.
	AppCmdName string
	// AppName holds the name of the application. By default, this is the same as AppCmdName.
	AppName string
	// CopyrightStartYear holds the starting year to place in the copyright banner. If not set explicitly, will be set
	// to the year of the "vcs.time" build tag.
	CopyrightStartYear string
	// CopyrightEndYear holds the ending year to place in the copyright banner. If not set explicitly, will be set to
	// the year of the "vcs.time" build tag.
	CopyrightEndYear string
	// CopyrightHolder holds the name of the copyright holder.
	CopyrightHolder string
	// License holds the license the software is being distributed under. This is intended to be a simple one line
	// description, such as "Mozilla Public License 2.0" and not the full license itself.
	License string
	// AppVersion holds the application's version information. If not set explicitly, will be the version of the main
	// module or "0.0" if that isn't available. Unfortunately, this automatic setting only works for binaries created
	// using `go install <package>@<version>`.
	AppVersion string
	// VCSName holds the name of the version control system and is set to the value of the build tag "vcs".
	VCSName string
	// VCSVersion holds the vcs revision and is set to the value of the build tag "vcs.revision".
	VCSVersion string
	// VCSModified is true if the "vcs.modified" build tag is true.
	VCSModified bool
	// BuildNumber holds the build number. If not set explicitly, will be generated from the value of the build tag
	// "vcs.time".
	BuildNumber string
	// AppIdentifier holds the uniform type identifier (UTI) for the application. This should contain only alphanumeric
	// (A-Z,a-z,0-9), hyphen (-), and period (.) characters. The string should also be in reverse-DNS format. For
	// example, if your company’s domain is ajax.com and you create an application named Hello, you could assign the
	// string com.ajax.Hello as your AppIdentifier.
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

// ShortAppVersion returns the app version. If the code was built from a dirty version of a VCS checkout, a trailing '~'
// character will be added.
func ShortAppVersion() string {
	return markAppVersionModified(AppVersion)
}

// LongAppVersion returns a combination of the app version and the build number. If the code was built from a dirty
// version of a VCS checkout, a trailing '~' character will be added.
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
