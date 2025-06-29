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
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestSetupDev(t *testing.T) {
	if os.Getenv("XSLOG_DEV_TEST") == "1" {
		// This is the subprocess
		xslog.SetupStd(true)
		slog.Info("test message")
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestSetupDev")
	cmd.Env = append(os.Environ(), "XSLOG_DEV_TEST=1")
	output, err := cmd.CombinedOutput()
	check.NoError(t, err)
	check.Contains(t, string(output), "test message | xslog/setup_test.go:")
}

func TestSetupNotDev(t *testing.T) {
	if os.Getenv("XSLOG_NOTDEV_TEST") == "1" {
		// This is the subprocess
		xslog.SetupStd(false)
		slog.Info("test message")
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestSetupNotDev")
	cmd.Env = append(os.Environ(), "XSLOG_NOTDEV_TEST=1")
	output, err := cmd.CombinedOutput()
	check.NoError(t, err)
	check.Contains(t, string(output), `"msg":"test message"`)
}
