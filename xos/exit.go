// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"slices"
	"sync"
	"syscall"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xreflect"
)

var (
	// ExitRecoveryHandler will be used to capture any panics caused by functions that have been installed when run
	// during exit. By default, this will be nil, which means the panics will be logged.
	ExitRecoveryHandler func(err error)
	// ExitCodeForSIGINT is the exit code used when the program is terminated by a SIGINT (Ctrl+C). Defaults to 1.
	ExitCodeForSIGINT = 1
	// ExitCodeForSIGTERM is the exit code used when the program is terminated by a SIGTERM. Defaults to 1.
	ExitCodeForSIGTERM = 1
	exitLock           sync.Mutex
	exitFuncs          []exitFunction
	lastExitID         int
	exiting            bool
)

type exitFunction struct {
	f  func()
	id int
}

// EnsureAtSignalHandlersAreInstalled ensures that the signal handlers for SIGINT and SIGTERM are installed. If they are
// already installed, this function does nothing. This will be called automatically by RunAtExit() if the signal
// handlers have not yet been installed.
func EnsureAtSignalHandlersAreInstalled() {
	exitLock.Lock()
	defer exitLock.Unlock()
	if lastExitID == 0 {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			if <-sigChan == syscall.SIGINT {
				fmt.Print("\b\b") // Removes the unsightly ^C in the terminal
				Exit(ExitCodeForSIGINT)
			}
			Exit(ExitCodeForSIGTERM)
		}()
	}
}

// RunAtExit registers a function to be run when xos.Exit() is called. Returns an ID that can be used to remove the
// function later, if needed. Calling xos.RunAtExit() after xos.Exit() has been called will have no effect.
func RunAtExit(f func()) int {
	EnsureAtSignalHandlersAreInstalled()
	exitLock.Lock()
	defer exitLock.Unlock()
	lastExitID++
	exitFuncs = append(exitFuncs, exitFunction{id: lastExitID, f: f})
	return lastExitID
}

// CancelRunAtExit unregisters a function that was previously registered by xos.RunAtExit(). If the ID is no longer
// present, nothing happens. Calling xos.CancelRunAtExit() after xos.Exit() has been called will have no effect.
func CancelRunAtExit(id int) {
	exitLock.Lock()
	defer exitLock.Unlock()
	exitFuncs = slices.DeleteFunc(exitFuncs, func(p exitFunction) bool { return p.id == id })
}

// Exit runs any registered exit functions in the inverse order they were registered and then exits the progream with
// the specified status. If a previous call to xos.Exit() is already being handled, this method does nothing and does
// not return. Recursive calls to xos.Exit() will trigger a panic, which the exit handling will catch and report, but
// will then proceed with exit as normal. Note that once xos.Exit() is called, no subsequent changes to the registered
// list of functions will have an effect.
func Exit(status int) {
	var f []func()
	exitLock.Lock()
	wasExiting := exiting
	if !wasExiting {
		// We weren't already exiting, so mark us as exiting and make a copy of the exit functions
		exiting = true
		f = make([]func(), len(exitFuncs))
		for i, one := range exitFuncs {
			f[i] = one.f
		}
	}
	exitLock.Unlock()
	if wasExiting {
		// Check for recursive calls
		var pcs [512]uintptr
		n := runtime.Callers(2, pcs[:])
		frames := runtime.CallersFrames(pcs[:n])
		for {
			frame, more := frames.Next()
			if frame.Function == "github.com/richardwilkes/toolbox/v2/xos.Exit" {
				// We're in a recursive call, so we need to panic to trigger the recovery mechanism
				panic("recursive call of xos.Exit()")
			}
			if !more {
				break
			}
		}
		// We're being called from another goroutine, so we need to park it and allow the exit to complete
		select {}
	}
	// Run the exit functions in reverse order from how they were registered, then exit
	for i := len(f) - 1; i >= 0; i-- {
		SafeCall(f[i], ExitRecoveryHandler)
	}
	os.Exit(status)
}

// ExitIfErr checks the error and if it isn't nil, calls xos.ExitWithErr(err).
func ExitIfErr(err error) {
	if !xreflect.IsNil(err) {
		ExitWithErr(err)
	}
}

// ExitWithErr logs the error and then exits with code 1.
func ExitWithErr(err error) {
	errs.Log(err)
	Exit(1)
}
