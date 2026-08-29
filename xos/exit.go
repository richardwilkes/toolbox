// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	// ExitRecoveryHandler receives, as an error, any panic raised by an exit function run by Exit. If nil (the default),
	// the panic is logged.
	ExitRecoveryHandler func(err error)
	// ExitCodeForSIGINT is the exit code used when the program is terminated by a SIGINT (Ctrl+C). Defaults to 1.
	ExitCodeForSIGINT = 1
	// ExitCodeForSIGTERM is the exit code used when the program is terminated by a SIGTERM. Defaults to 1.
	ExitCodeForSIGTERM      = 1
	exitLock                sync.Mutex
	exitFuncs               []exitFunction
	lastExitID              int
	exiting                 bool
	signalHandlersInstalled bool
)

type exitFunction struct {
	f  func()
	id int
}

// EnsureAtSignalHandlersAreInstalled installs the SIGINT and SIGTERM handlers if they are not already installed.
// RunAtExit calls this automatically.
func EnsureAtSignalHandlersAreInstalled() {
	exitLock.Lock()
	defer exitLock.Unlock()
	if !signalHandlersInstalled {
		signalHandlersInstalled = true
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

// RunAtExit registers a function to be run when Exit is called and returns an ID that can be passed to
// CancelRunAtExit. Calls made after Exit has been called have no effect.
func RunAtExit(f func()) int {
	EnsureAtSignalHandlersAreInstalled()
	exitLock.Lock()
	defer exitLock.Unlock()
	lastExitID++
	exitFuncs = append(exitFuncs, exitFunction{id: lastExitID, f: f})
	return lastExitID
}

// CancelRunAtExit unregisters a function previously registered by RunAtExit. If the ID is not present, or Exit has
// already been called, nothing happens.
func CancelRunAtExit(id int) {
	exitLock.Lock()
	defer exitLock.Unlock()
	exitFuncs = slices.DeleteFunc(exitFuncs, func(p exitFunction) bool { return p.id == id })
}

// Exit runs the registered exit functions in reverse registration order, then exits the program with the given status.
// If an exit is already in progress, a call from another goroutine blocks forever, while a recursive call from within
// an exit function panics; that panic is caught and reported by the exit handling, which then proceeds as normal.
// Changes to the registered functions after Exit has been called have no effect.
func Exit(status int) {
	var f []func()
	exitLock.Lock()
	wasExiting := exiting
	if !wasExiting {
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
				// Panic so the exit function's SafeCall reports it and the exit proceeds.
				panic("recursive call of xos.Exit()")
			}
			if !more {
				break
			}
		}
		// Called from another goroutine; park it until the exit completes.
		select {}
	}
	for i := len(f) - 1; i >= 0; i-- {
		SafeCall(f[i], ExitRecoveryHandler)
	}
	os.Exit(status)
}

// ExitIfErr calls ExitWithErr(err) if err is not nil, as determined by xreflect.IsNil.
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

// ExitWithMsg writes the message to stderr and then exits with code 1.
func ExitWithMsg(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	Exit(1)
}
