// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog_test

import (
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/cmdline"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestLevelValueSet(t *testing.T) {
	var v xslog.LevelValue
	c := check.New(t)
	c.NoError(v.Set("error"))
	c.Equal(v.Level(), slog.LevelError)
}

func TestLevelValueCmdLineOpts(t *testing.T) {
	cl := cmdline.New(false)
	var v xslog.LevelValue
	v.AddStdCmdLineOptions(cl)
	if os.Getenv("LEVEL_VALUE_CMDLINE_TEST") == "1" {
		// This is the subprocess
		cl.DisplayUsage()
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestLevelValueCmdLineOpts")
	cmd.Env = append(os.Environ(), "LEVEL_VALUE_CMDLINE_TEST=1")
	output, err := cmd.CombinedOutput()
	c := check.New(t)
	c.NoError(err)
	c.Contains(string(output), "--log-level <value>")
}
