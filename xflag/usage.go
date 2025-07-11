// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xflag

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xos"
	"github.com/richardwilkes/toolbox/v2/xterm"
)

// SetUsage replaces any existing Usage function on the given flagSet with one that provides more information. You may
// pass in nil for the flagSet to use the default flag set (flag.CommandLine).
func SetUsage(flagSet *flag.FlagSet, description, argsUsage string) {
	if flagSet == nil {
		flagSet = flag.CommandLine
	}
	flagSet.Usage = func() {
		var w *xterm.AnsiWriter
		switch out := flagSet.Output().(type) {
		case *xterm.AnsiWriter:
			w = out
		default:
			w = xterm.NewAnsiWriter(out)
		}
		kind := w.Kind()
		w.WriteString(kind.Reset() + "\n")
		w.WrapText("", kind.Green()+xos.AppName+kind.Dim())
		buildInfo := fmt.Sprintf(i18n.Text("Version %s"), xos.ShortAppVersion())
		if xos.BuildNumber != "" {
			buildInfo = fmt.Sprintf(i18n.Text("%s, Build %s"), buildInfo, xos.BuildNumber)
		}
		w.WrapText("", buildInfo)
		if xos.VCSName != "" && xos.VCSVersion != "" {
			str := xos.VCSName + ": " + xos.VCSVersion
			if xos.VCSModified {
				str += "-modified"
			}
			w.WrapText("", str)
		}
		copyright := xos.Copyright()
		if copyright != "" {
			w.WrapText("", copyright)
		}
		if xos.License != "" {
			w.WrapText("", fmt.Sprintf(i18n.Text("License: %s"), xos.License))
		}
		w.WriteString(kind.Reset() + "\n")
		if description != "" {
			w.WrapText("", description)
			w.WriteByte('\n')
		}
		usage := kind.Yellow() + xos.AppCmdName
		type state struct {
			flag  *flag.Flag
			name  string
			usage string
			size  int
			zero  bool
		}
		var flags []state
		largest := 0
		flagSet.VisitAll(func(f *flag.Flag) {
			argName, revisedUsage := flag.UnquoteUsage(f)
			size := 3 + len([]rune(f.Name))
			if argName != "" {
				size += 3 + len([]rune(argName))
			}
			if largest < size {
				largest = size
			}
			t := reflect.TypeOf(f.Value)
			var z reflect.Value
			if t.Kind() == reflect.Pointer {
				z = reflect.New(t.Elem())
			} else {
				z = reflect.Zero(t)
			}
			v, ok := z.Interface().(flag.Value)
			if !ok {
				panic("unable to call .String() on flag: -" + f.Name)
			}
			flags = append(flags, state{
				flag:  f,
				name:  argName,
				usage: revisedUsage,
				size:  size,
				zero:  f.DefValue == v.String(),
			})
		})
		if len(flags) != 0 {
			usage += i18n.Text(" [options]")
		}
		if argsUsage != "" {
			usage += " " + argsUsage
		}
		usage += kind.Reset()
		w.WrapText(i18n.Text("Usage: "), usage)
		w.WriteByte('\n')
		if len(flags) != 0 {
			fmt.Fprintln(w, i18n.Text("Options:"))
			w.WriteByte('\n')
			for _, f := range flags {
				prefix := kind.Yellow() + "  -" + f.flag.Name
				if f.name != "" {
					prefix += kind.Dim() + " <" + f.name + ">"
				}
				prefix += kind.Reset()
				full := f.usage
				if !f.zero {
					full += fmt.Sprintf(i18n.Text(" (default: %s%s%s)"), kind.Blue(), f.flag.DefValue, kind.Reset())
				}
				w.WrapText(prefix+strings.Repeat(" ", 1+largest-f.size), full)
			}
		}
		w.WriteByte('\n')
	}
}
