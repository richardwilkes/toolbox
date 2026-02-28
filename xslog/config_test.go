// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xflag"
	"github.com/richardwilkes/toolbox/v2/xos"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestConfig(t *testing.T) {
	if testFlag := os.Getenv("XSLOG_FLAGS_TEST"); testFlag != "" {
		// This is the subprocess
		var c xslog.Config
		c.RotatorCfg.Path = os.Getenv("XSLOG_PATH_TEST")
		os.Args = os.Args[:1]
		switch testFlag {
		case "1":
			c.Console = true
		case "2":
			os.Args = append(os.Args, "-quiet")
			c.Console = true
		case "3":
			os.Args = append(os.Args, "-q")
			c.Console = true
		case "4":
		case "5":
			os.Args = append(os.Args, "-console")
		}
		c.AddFlags()
		xflag.Parse()
		slog.Info("test message")
		xos.Exit(0)
	}

	// Run test 1 (console enabled, no options on command line) in subprocesses
	c := check.New(t)
	tmpLogFile := filepath.Join(c.TempDir(), "xslog_test1")
	cmd := exec.Command(os.Args[0], "-test.run=TestConfig")
	cmd.Env = append(os.Environ(), "XSLOG_FLAGS_TEST=1", "XSLOG_PATH_TEST="+tmpLogFile)
	output, err := cmd.CombinedOutput()
	c.NoError(err)
	expectedMsgFragment := " | test message | xslog/config_test.go:"
	c.Contains(string(output), expectedMsgFragment)
	data, err := os.ReadFile(tmpLogFile + ".log")
	c.NoError(err)
	c.Contains(string(data), expectedMsgFragment)

	// Run test 2 (console enabled, -quiet on command line) in subprocesses
	tmpLogFile = filepath.Join(c.TempDir(), "xslog_test2")
	cmd = exec.Command(os.Args[0], "-test.run=TestConfig")
	cmd.Env = append(os.Environ(), "XSLOG_FLAGS_TEST=2", "XSLOG_PATH_TEST="+tmpLogFile)
	output, err = cmd.CombinedOutput()
	c.NoError(err)
	c.Equal(0, len(output))
	data, err = os.ReadFile(tmpLogFile + ".log")
	c.NoError(err)
	c.Contains(string(data), expectedMsgFragment)

	// Run test 3 (console enabled, -q on command line) in subprocesses
	tmpLogFile = filepath.Join(c.TempDir(), "xslog_test3")
	cmd = exec.Command(os.Args[0], "-test.run=TestConfig")
	cmd.Env = append(os.Environ(), "XSLOG_FLAGS_TEST=3", "XSLOG_PATH_TEST="+tmpLogFile)
	output, err = cmd.CombinedOutput()
	c.NoError(err)
	c.Equal(0, len(output))
	data, err = os.ReadFile(tmpLogFile + ".log")
	c.NoError(err)
	c.Contains(string(data), expectedMsgFragment)

	// Run test 4 (console not enabled, no options on command line) in subprocesses
	tmpLogFile = filepath.Join(c.TempDir(), "xslog_test4")
	cmd = exec.Command(os.Args[0], "-test.run=TestConfig")
	cmd.Env = append(os.Environ(), "XSLOG_FLAGS_TEST=4", "XSLOG_PATH_TEST="+tmpLogFile)
	output, err = cmd.CombinedOutput()
	c.NoError(err)
	c.Equal(0, len(output))
	data, err = os.ReadFile(tmpLogFile + ".log")
	c.NoError(err)
	c.Contains(string(data), expectedMsgFragment)

	// Run test 5 (console not enabled, -console on command line) in subprocesses
	tmpLogFile = filepath.Join(c.TempDir(), "xslog_test5")
	cmd = exec.Command(os.Args[0], "-test.run=TestConfig")
	cmd.Env = append(os.Environ(), "XSLOG_FLAGS_TEST=5", "XSLOG_PATH_TEST="+tmpLogFile)
	output, err = cmd.CombinedOutput()
	c.NoError(err)
	c.Contains(string(output), expectedMsgFragment)
	data, err = os.ReadFile(tmpLogFile + ".log")
	c.NoError(err)
	c.Contains(string(data), expectedMsgFragment)
}
