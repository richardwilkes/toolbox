// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio/term"
)

var (
	// AppCmdName holds the application's name as specified on the command line.
	AppCmdName string
	// AppName holds the name of the application. By default, this is the same as AppCmdName.
	AppName string
	// CopyrightYears holds the years to place in the copyright banner. Instead of setting this explicitly, consider
	// using CopyrightStartYear and CopyrightEndYear instead. For example, setting CopyrightStartYear early in your
	// main() method, and allowing the build system to populate CopyrightEndYear for you.
	CopyrightYears string
	// CopyrightStartYear holds the starting year to place in the copyright banner. Will not be used if CopyrightYears
	// is already set. If not set explicitly, will be set to the year of the "vcs.time" build tag, if available.
	CopyrightStartYear string
	// CopyrightEndYear holds the ending year to place in the copyright banner. Will not be used if CopyrightYears is
	// already set. If not set explicitly, will be set to the year of the "vcs.time" build tag, if available.
	CopyrightEndYear string
	// CopyrightHolder holds the name of the copyright holder.
	CopyrightHolder string
	// License holds the license the software is being distributed under. This is intended to be a simple one line
	// description, such as "Mozilla Public License 2.0" and not the full license itself.
	License string
	// AppVersion holds the application's version information. If not set explicitly, will be the version of the main
	// module.
	AppVersion string
	// GitVersion holds the vcs revision and clean/dirty status. If not set explicitly, will be generated from the value
	// of the build tags "vcs.revision" and "vcs.modified".
	GitVersion string
	// BuildNumber holds the build number. If not set explicitly, will be generated from the value of the build tag
	// "vcs.time".
	BuildNumber string
	// AppIdentifier holds the uniform type identifier (UTI) for the application. This should contain only alphanumeric
	// (A-Z,a-z,0-9), hyphen (-), and period (.) characters. The string should also be in reverse-DNS format. For
	// example, if your company’s domain is ajax.com and you create an application named Hello, you could assign the
	// string com.ajax.Hello as your AppIdentifier.
	AppIdentifier string
	vcs           = "git"
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
	var vcsRevision string
	var vcsTime time.Time
	var vcsModified bool
	if info, ok := debug.ReadBuildInfo(); ok {
		if AppVersion == "" {
			if info.Main.Version == "(devel)" {
				AppVersion = "0.0"
			} else {
				AppVersion = info.Main.Version
			}
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs":
				vcs = setting.Value
			case "vcs.revision":
				vcsRevision = setting.Value
			case "vcs.time":
				if t, err := time.Parse(time.RFC3339, setting.Value); err == nil {
					vcsTime = t
				}
			case "vcs.modified":
				if setting.Value == "true" {
					vcsModified = true
				}
			}
		}
	}
	if GitVersion == "" && vcsRevision != "" {
		GitVersion = vcsRevision
		if vcsModified {
			GitVersion += "-dirty"
		}
	}
	if !vcsModified && !vcsTime.IsZero() {
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
}

// ResolveCopyrightYears resolves the copyright years. If the CopyrightYears has been explicitly set, that will be
// returned unmodified. Otherwise, it will be generated based on the values of CopyrightStartYear and CopyrightEndYear.
func ResolveCopyrightYears() string {
	if CopyrightYears != "" {
		return CopyrightYears
	}
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

// Copyright returns the copyright notice.
func Copyright() string {
	var dot string
	if !strings.HasSuffix(CopyrightHolder, ".") {
		dot = "."
	}
	return fmt.Sprintf(i18n.Text("Copyright © %[1]s by %[2]s%[3]s All rights reserved."), ResolveCopyrightYears(),
		CopyrightHolder, dot)
}

// DisplayUsage displays the program usage information.
func (cl *CmdLine) DisplayUsage() {
	term.WrapText(cl, "", AppName)
	version := AppVersion
	if version == "" {
		version = "0.0"
	}
	buildInfo := fmt.Sprintf(i18n.Text("Version %s"), version)
	if BuildNumber != "" {
		buildInfo = fmt.Sprintf(i18n.Text("%s, Build %s"), buildInfo, BuildNumber)
	}
	term.WrapText(cl, "  ", buildInfo)
	if GitVersion != "" {
		term.WrapText(cl, "  ", vcs+": "+GitVersion)
	}
	term.WrapText(cl, "  ", Copyright())
	if License != "" {
		term.WrapText(cl, "  ", fmt.Sprintf(i18n.Text("License: %s"), License))
	}
	fmt.Fprintln(cl)
	if cl.Description != "" {
		term.WrapText(cl, "", cl.Description)
		fmt.Fprintln(cl)
	}
	usage := fmt.Sprintf(i18n.Text("%s [options]"), AppCmdName)
	opts := cl
	var stack []*CmdLine
	for opts != nil {
		stack = append(stack, opts)
		opts = opts.parent
	}
	for i := len(stack) - 1; i >= 0; i-- {
		one := stack[i]
		if one.cmd == nil {
			if i == 0 && len(cl.cmds) > 0 {
				usage += i18n.Text(" <command> [command options]")
			}
		} else {
			usage += fmt.Sprintf(i18n.Text(" %[1]s [%[1]s options]"), one.cmd.Name())
		}
	}
	if cl.UsageSuffix != "" {
		usage += " " + cl.UsageSuffix
	}
	term.WrapText(cl, i18n.Text("Usage: "), usage)
	for i := len(stack) - 1; i >= 0; i-- {
		one := stack[i]
		fmt.Fprintln(one)
		if one.cmd == nil {
			if i == 0 {
				usage += i18n.Text(" <command> [command options]")
			}
			fmt.Fprintln(one, i18n.Text("Options:"))
		} else {
			fmt.Fprintf(one, i18n.Text("%s options:\n"), one.cmd.Name())
		}
		fmt.Fprintln(one)
		one.displayOptions()
	}
	cl.displayCommands(2)
}

func (cl *CmdLine) displayOptions() {
	sort.Sort(cl.options)
	hasShort := false
	largest := 0
	for _, option := range cl.options {
		if option.usage == "" {
			continue
		}
		if option.single != 0 {
			hasShort = true
		}
		length := len([]rune(option.name))
		if length > 0 {
			length += 2
		}
		if !option.isBool() {
			if length > 0 {
				length++
			}
			length += 2 + len([]rune(option.arg))
		}
		if length > largest {
			largest = length
		}
	}
	largest += 2
	for _, option := range cl.options {
		if option.usage == "" {
			continue
		}
		var sn string
		if hasShort {
			if option.single != 0 {
				sn = "-" + string(option.single)
				if option.name != "" {
					sn += ", "
				} else {
					sn += "  "
				}
			} else {
				sn = "    "
			}
		}
		var ln string
		if option.name != "" {
			ln = "--" + option.name
		}
		if !option.isBool() {
			if ln != "" {
				ln += " "
			}
			ln += "<" + option.arg + ">"
		}
		prefix := "  " + sn + ln + strings.Repeat(" ", largest-len([]rune(ln)))
		usage := option.usage
		if !strings.HasSuffix(usage, ".") {
			usage += "."
		}
		if !option.isBool() && option.def != "" {
			usage += i18n.Text(" Default: ")
			usage += option.def
		}
		term.WrapText(cl, prefix, usage)
	}
}

func (cl *CmdLine) displayCommands(indent int) {
	if len(cl.cmds) > 0 {
		fmt.Fprintln(cl)
		term.WrapText(cl, "", i18n.Text("Available commands:"))
		fmt.Fprintln(cl)
		var all []string
		largest := 0
		for key := range cl.cmds {
			all = append(all, key)
			length := len(key)
			if length > largest {
				largest = length
			}
		}
		sort.Strings(all)
		format := fmt.Sprintf("%s%%-%ds  ", strings.Repeat(" ", indent), largest)
		for _, cmd := range all {
			term.WrapText(cl, fmt.Sprintf(format, cmd), cl.cmds[cmd].Usage())
		}
		fmt.Fprintln(cl)
		term.WrapText(cl, "", fmt.Sprintf(i18n.Text("Use '%s help <command>' to see command options"), AppCmdName))
	}
}
