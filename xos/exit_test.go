// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestExit_ExecutionOrder(t *testing.T) {
	if os.Getenv("EXIT_TEST") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { fmt.Print("1") })
		xos.RunAtExit(func() { fmt.Print("2") })
		xos.RunAtExit(func() { fmt.Print("3") })
		xos.Exit(22)
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExit_ExecutionOrder")
	cmd.Env = append(os.Environ(), "EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	check.Error(t, err)
	check.Equal(t, "321", string(output))
	var exitError *exec.ExitError
	hasExitErr := errors.As(err, &exitError)
	check.True(t, hasExitErr)
	if hasExitErr {
		check.Equal(t, 22, exitError.ExitCode())
	}
}

func TestExit_ExecutionOrderWithCancel(t *testing.T) {
	if os.Getenv("EXIT_TEST_WITH_CANCEL") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { fmt.Print("1") })
		id := xos.RunAtExit(func() { fmt.Print("2") })
		xos.RunAtExit(func() { fmt.Print("3") })
		xos.CancelRunAtExit(id)
		xos.Exit(22)
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExit_ExecutionOrderWithCancel")
	cmd.Env = append(os.Environ(), "EXIT_TEST_WITH_CANCEL=1")
	output, err := cmd.CombinedOutput()
	check.Error(t, err)
	check.Equal(t, "31", string(output))
	var exitError *exec.ExitError
	hasExitErr := errors.As(err, &exitError)
	check.True(t, hasExitErr)
	if hasExitErr {
		check.Equal(t, 22, exitError.ExitCode())
	}
}

func TestExit_ConcurrentCalls(t *testing.T) {
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
			wg.Add(1)
			go func() {
				defer wg.Done()
				xos.Exit(0)
			}()
		}
		wg.Wait()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExit_ConcurrentCalls")
	cmd.Env = append(os.Environ(), "CONCURRENT_EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	check.NoError(t, err)
	check.Equal(t, "executed", string(output))
}

func TestExit_PanicInExitFunction(t *testing.T) {
	if os.Getenv("PANIC_EXIT_TEST") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { fmt.Print("first") })
		xos.RunAtExit(func() { panic("test panic in exit function") })
		xos.RunAtExit(func() { fmt.Print("third") })
		xos.Exit(0)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExit_PanicInExitFunction")
	cmd.Env = append(os.Environ(), "PANIC_EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	check.NoError(t, err)
	check.Contains(t, string(output), "first")
	check.Contains(t, string(output), "test panic in exit function")
	check.Contains(t, string(output), "third")
}

func TestExit_RecursiveExit(t *testing.T) {
	if os.Getenv("RECURSIVE_EXIT_TEST") == "1" {
		// This is the subprocess
		xos.RunAtExit(func() { xos.Exit(2) })
		xos.RunAtExit(func() { fmt.Print("normal") })
		xos.Exit(0)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExit_RecursiveExit")
	cmd.Env = append(os.Environ(), "RECURSIVE_EXIT_TEST=1")
	output, err := cmd.CombinedOutput()
	check.NoError(t, err)
	check.Contains(t, string(output), "recursive call of xos.Exit()")
	check.Contains(t, string(output), "normal")
}

func TestExit_SIGINT(t *testing.T) {
	if os.Getenv("SIGINT_EXIT_TEST") == "1" {
		// This is the subprocess
		xos.ExitCodeForSIGINT = 234
		xos.EnsureAtSignalHandlersAreInstalled()
		select {}
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExit_SIGINT")
	cmd.Env = append(os.Environ(), "SIGINT_EXIT_TEST=1")
	check.NoError(t, cmd.Start())
	time.Sleep(100 * time.Millisecond) // Give the command time to start
	check.NoError(t, cmd.Process.Signal(syscall.SIGINT))
	err := cmd.Wait()
	check.Error(t, err)
	var exitError *exec.ExitError
	hasExitErr := errors.As(err, &exitError)
	check.True(t, hasExitErr)
	if hasExitErr {
		check.Equal(t, 234, exitError.ExitCode())
	}
}

func TestExit_SIGTERM(t *testing.T) {
	if os.Getenv("SIGTERM_EXIT_TEST") == "1" {
		// This is the subprocess
		xos.ExitCodeForSIGTERM = 123
		xos.EnsureAtSignalHandlersAreInstalled()
		select {}
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExit_SIGTERM")
	cmd.Env = append(os.Environ(), "SIGTERM_EXIT_TEST=1")
	check.NoError(t, cmd.Start())
	time.Sleep(100 * time.Millisecond) // Give the command time to start
	check.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
	err := cmd.Wait()
	check.Error(t, err)
	var exitError *exec.ExitError
	hasExitErr := errors.As(err, &exitError)
	check.True(t, hasExitErr)
	if hasExitErr {
		check.Equal(t, 123, exitError.ExitCode())
	}
}
