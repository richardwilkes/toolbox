// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestExecutionOrderExit(t *testing.T) {
	if os.Getenv("EXIT_TEST") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { fmt.Print("1") })
		xos.RunAtExit(func() { fmt.Print("2") })
		xos.RunAtExit(func() { fmt.Print("3") })
		xos.Exit(22)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExecutionOrderExit")
	cmd.Env = append(os.Environ(), "EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	c := check.New(t)
	c.HasError(err)
	c.Equal("321", string(output))
	var exitError *exec.ExitError
	hasExitErr := errors.As(err, &exitError)
	c.True(hasExitErr)
	if hasExitErr {
		c.Equal(22, exitError.ExitCode())
	}
}

func TestExecutionOrderWithCancelExit(t *testing.T) {
	if os.Getenv("EXIT_TEST_WITH_CANCEL") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { fmt.Print("1") })
		id := xos.RunAtExit(func() { fmt.Print("2") })
		xos.RunAtExit(func() { fmt.Print("3") })
		xos.CancelRunAtExit(id)
		xos.Exit(22)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExecutionOrderWithCancelExit")
	cmd.Env = append(os.Environ(), "EXIT_TEST_WITH_CANCEL=1")
	output, err := cmd.CombinedOutput()
	c := check.New(t)
	c.HasError(err)
	c.Equal("31", string(output))
	var exitError *exec.ExitError
	hasExitErr := errors.As(err, &exitError)
	c.True(hasExitErr)
	if hasExitErr {
		c.Equal(22, exitError.ExitCode())
	}
}

func TestConcurrentCallsExit(t *testing.T) {
	if os.Getenv("CONCURRENT_EXIT_TEST") == "1" {
		// This is the subprocess
		var executionCount int64
		var wg sync.WaitGroup
		xos.RunAtExit(func() {
			count := atomic.AddInt64(&executionCount, 1)
			fmt.Print("executed")
			time.Sleep(10 * time.Millisecond) // Simulate work
			if count > 1 {
				fmt.Print("multiple")
			}
		})
		for range 5 {
			wg.Go(func() {
				xos.Exit(0)
			})
		}
		wg.Wait()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestConcurrentCallsExit")
	cmd.Env = append(os.Environ(), "CONCURRENT_EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	c := check.New(t)
	c.NoError(err)
	c.Equal("executed", string(output))
}

func TestPanicInExitFunction(t *testing.T) {
	if os.Getenv("PANIC_EXIT_TEST") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { fmt.Print("first") })
		xos.RunAtExit(func() { panic("test panic in exit function") })
		xos.RunAtExit(func() { fmt.Print("third") })
		xos.Exit(0)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPanicInExitFunction")
	cmd.Env = append(os.Environ(), "PANIC_EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	c := check.New(t)
	c.NoError(err)
	c.Contains(string(output), "first")
	c.Contains(string(output), "test panic in exit function")
	c.Contains(string(output), "third")
}

func TestRecursiveExit(t *testing.T) {
	if os.Getenv("RECURSIVE_EXIT_TEST") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { xos.Exit(2) })
		xos.RunAtExit(func() { fmt.Print("normal") })
		xos.Exit(0)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRecursiveExit")
	cmd.Env = append(os.Environ(), "RECURSIVE_EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	c := check.New(t)
	c.NoError(err)
	c.Contains(string(output), "recursive call of xos.Exit()")
	c.Contains(string(output), "normal")
}

func TestSignalsExitWithConfiguredCode(t *testing.T) {
	if runtime.GOOS == xos.WindowsOS {
		t.Skip("This test requires permissions that aren't available by default on Windows")
	}
	cases := []struct {
		exitCode *int
		name     string
		signal   syscall.Signal
		code     int
	}{
		{name: "SIGINT", signal: syscall.SIGINT, exitCode: &xos.ExitCodeForSIGINT, code: 99},
		{name: "SIGTERM", signal: syscall.SIGTERM, exitCode: &xos.ExitCodeForSIGTERM, code: 123},
	}
	if name := os.Getenv("SIGNAL_EXIT_TEST"); name != "" {
		// This is the subprocess
		for _, tc := range cases {
			if tc.name == name {
				*tc.exitCode = tc.code
				xos.EnsureAtSignalHandlersAreInstalled()
				select {}
			}
		}
		return
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestSignalsExitWithConfiguredCode")
			cmd.Env = append(os.Environ(), "SIGNAL_EXIT_TEST="+tc.name)
			c := check.New(t)
			c.NoError(cmd.Start())
			time.Sleep(100 * time.Millisecond) // Give the command time to start
			c.NoError(cmd.Process.Signal(tc.signal))
			err := cmd.Wait()
			c.HasError(err)
			var exitError *exec.ExitError
			hasExitErr := errors.As(err, &exitError)
			c.True(hasExitErr)
			if hasExitErr {
				c.Equal(tc.code, exitError.ExitCode())
			}
		})
	}
}

func TestEnsureAtSignalHandlersAreInstalledIsIdempotent(t *testing.T) {
	c := check.New(t)
	// Install the handlers and give the goroutine a moment to start and block on its channel.
	xos.EnsureAtSignalHandlersAreInstalled()
	time.Sleep(50 * time.Millisecond)
	before := runtime.NumGoroutine()
	// Repeated direct calls (with no intervening RunAtExit) must not install duplicate handlers or leak goroutines.
	for range 20 {
		xos.EnsureAtSignalHandlersAreInstalled()
	}
	time.Sleep(50 * time.Millisecond)
	after := runtime.NumGoroutine()
	// With the previous bug, each repeated call leaked a goroutine, so after would be before+20. Allow for a little
	// unrelated runtime goroutine churn, but anything near +20 indicates the leak has returned.
	c.True(after-before < 10)
}

func TestExitHelpersReportFailureAndExit(t *testing.T) {
	cases := []struct {
		fn       func()
		name     string
		output   string
		exitCode int
	}{
		{name: "ExitIfErr", fn: func() { xos.ExitIfErr(errors.New("bang")) }, output: "bang", exitCode: 1},
		{name: "ExitIfErr nil", fn: func() { xos.ExitIfErr(nil) }},
		{name: "ExitWithErr", fn: func() { xos.ExitWithErr(errors.New("boom")) }, output: "boom", exitCode: 1},
		{name: "ExitWithMsg", fn: func() { xos.ExitWithMsg("Outta here") }, output: "Outta here", exitCode: 1},
	}
	if name := os.Getenv("EXIT_HELPER_TEST"); name != "" {
		// This is the subprocess
		for _, tc := range cases {
			if tc.name == name {
				tc.fn()
			}
		}
		return
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestExitHelpersReportFailureAndExit")
			cmd.Env = append(os.Environ(), "EXIT_HELPER_TEST="+tc.name)
			output, err := cmd.CombinedOutput()
			c := check.New(t)
			if tc.exitCode == 0 {
				c.NoError(err)
				return
			}
			c.HasError(err)
			c.Contains(string(output), tc.output)
			var exitError *exec.ExitError
			hasExitErr := errors.As(err, &exitError)
			c.True(hasExitErr)
			if hasExitErr {
				c.Equal(tc.exitCode, exitError.ExitCode())
			}
		})
	}
}
