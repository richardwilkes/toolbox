// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xflag_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xflag"
)

func TestSplitCommandLine(t *testing.T) {
	c := check.New(t)
	splitChecker(c, "cmd", []string{"cmd"}, false)
	splitChecker(c, `cmd "world hello"`, []string{"cmd", "world hello"}, false)
	splitChecker(c, `'cmd with spaces' "hello world"`, []string{"cmd with spaces", "hello world"}, false)
	splitChecker(c, `cmd \"hello\ world\"`, []string{"cmd", `"hello world"`}, false)
	splitChecker(c, "cmd 世界", []string{"cmd", "世界"}, false)
	splitChecker(c, `spacey\ cmd`, []string{"spacey cmd"}, false)
	splitChecker(c, "", []string(nil), false)
	splitChecker(c, "   ", []string(nil), false)
	splitChecker(c, "cmd   arg1   arg2", []string{"cmd", "arg1", "arg2"}, false)
	splitChecker(c, `cmd\ with\ spaces --option=value`, []string{"cmd with spaces", "--option=value"}, false)
	splitChecker(c, "cmd -o value", []string{"cmd", "-o", "value"}, false)
	splitChecker(c, `cmd --long-option 'some value'`, []string{"cmd", "--long-option", "some value"}, false)
	splitChecker(c, `cmd \"quoted\ arg\"`, []string{"cmd", `"quoted arg"`}, false)
	splitChecker(c, `cmd 'unterminated`, []string(nil), true)
	splitChecker(c, `cmd "unterminated`, []string(nil), true)
	splitChecker(c, `cmd 'single \' quote'`, []string{"cmd", "single ' quote"}, false)
	splitChecker(c, `cmd "double \" quote"`, []string{"cmd", `double " quote`}, false)
	splitChecker(c, `cmd	arg1	arg2`, []string{"cmd", "arg1", "arg2"}, false)
	splitChecker(c, `cmd	arg1 arg2`, []string{"cmd", "arg1", "arg2"}, false)
	splitChecker(c, `cmd\ arg1`, []string{"cmd arg1"}, false)
	splitChecker(c, `cmd\ arg1\ arg2`, []string{"cmd arg1 arg2"}, false)
	splitChecker(c, `cmd\"arg\"`, []string{"cmd\"arg\""}, false)
	splitChecker(c, `cmd \"arg\"`, []string{"cmd", `"arg"`}, false)
	splitChecker(c, `cmd 'a "b" c'`, []string{"cmd", `a "b" c`}, false)
	splitChecker(c, `cmd \"a 'b' c\"`, []string{"cmd", `"a`, "b", `c"`}, false)
	splitChecker(c, "cmd\narg1", []string{"cmd\narg1"}, false)
	splitChecker(c, "cmd arg1\narg2", []string{"cmd", "arg1\narg2"}, false)
	splitChecker(c, `cmd\`, []string(nil), true)
	splitChecker(c, `cmd 'a b'c`, []string{"cmd", "a b", "c"}, false)
	splitChecker(c, `cmd ''`, []string{"cmd", ""}, false)
	splitChecker(c, `cmd \"\"`, []string{"cmd", `""`}, false)
}

func splitChecker(c check.Checker, in string, expected []string, shouldErr bool) {
	c.Helper()
	parts, err := xflag.SplitCommandLine(in)
	if shouldErr {
		c.HasError(err)
	} else {
		c.NoError(err)
		c.Equal(expected, parts)
		if !slices.Equal(expected, parts) {
			for i, one := range parts {
				fmt.Printf("%d: %q\n", i, one)
			}
		}
	}
}
