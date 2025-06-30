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
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/cmdline"
	"github.com/richardwilkes/toolbox/v2/xos"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestSetupStd(t *testing.T) {
	if os.Getenv("XSLOG_SETUP_STD_TEST") == "1" {
		// This is the subprocess
		saved := os.Args      // Save the original args
		os.Args = os.Args[:1] // Reset args to just the program name
		logFile, _ := xslog.SetupStd(cmdline.New(false))
		os.Args = saved
		slog.Info("test message")
		if data, err := os.ReadFile(logFile); err != nil {
			slog.Error("Failed to read log file", "error", err)
		} else {
			fmt.Println(string(data))
		}
		xos.Exit(0)
	}
	if os.Getenv("XSLOG_SETUP_STD_TEST") == "2" {
		// This is the subprocess
		saved := os.Args      // Save the original args
		os.Args = os.Args[:1] // Reset args to just the program name
		xslog.SetupStd(cmdline.New(false))
		os.Args = saved
		slog.Info("test message")
		xos.Exit(0)
	}
	// Run the tests in subprocesses
	cmd := exec.Command(os.Args[0], "-test.run=TestSetupStd")
	cmd.Env = append(os.Environ(), "XSLOG_SETUP_STD_TEST=1")
	output, err := cmd.CombinedOutput()
	check.NoError(t, err)
	check.Contains(t, string(output), " | test message | xslog/setup_test.go:")

	cmd = exec.Command(os.Args[0], "-test.run=TestSetupStd")
	cmd.Env = append(os.Environ(), "XSLOG_SETUP_STD_TEST=2")
	output, err = cmd.CombinedOutput()
	check.NoError(t, err)
	check.Equal(t, "", string(output))
}

func TestSetupConsole(t *testing.T) {
	if os.Getenv("XSLOG_SETUP_STD_TO_CONSOLE_TEST") == "1" {
		// This is the subprocess
		saved := os.Args      // Save the original args
		os.Args = os.Args[:1] // Reset args to just the program name
		xslog.SetupStdToConsole(cmdline.New(false))
		os.Args = saved
		slog.Info("test message")
		xos.Exit(0)
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestSetupConsole")
	cmd.Env = append(os.Environ(), "XSLOG_SETUP_STD_TO_CONSOLE_TEST=1")
	output, err := cmd.CombinedOutput()
	check.NoError(t, err)
	check.Contains(t, string(output), ` | test message | xslog/setup_test.go:`)
}
