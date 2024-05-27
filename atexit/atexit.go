/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

// Package atexit provides functionality similar to the C standard library's atexit() call.
package atexit

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"slices"
	"sync"
	"syscall"

	"github.com/richardwilkes/toolbox/errs"
)

var (
	// RecoveryHandler will be used to capture any panics caused by functions that have been installed when run during
	// exit. It may be set to nil to silently ignore them.
	RecoveryHandler errs.RecoveryHandler = func(err error) { errs.Log(err) }
	lock            sync.Mutex
	nextID          = 1
	pairs           []pair
	exiting         bool
)

type pair struct {
	f  func()
	id int
}

// Register a function to be run at exit. Returns an ID that can be used to remove the function later, if desired.
// Registering a function after Exit() has been called (i.e. in a function that was registered) will have no effect.
func Register(f func()) int {
	lock.Lock()
	defer lock.Unlock()
	if nextID == 1 {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go waitForSigInt(sigChan)
	}
	pairs = append(pairs, pair{id: nextID, f: f})
	nextID++
	return nextID - 1
}

func waitForSigInt(sigChan <-chan os.Signal) {
	s := <-sigChan
	if s == syscall.SIGINT {
		fmt.Print("\b\b") // Removes the unsightly ^C in the terminal
	}
	Exit(1)
}

// Unregister a function that was previously registered to be run at exit. If the ID is no longer present, nothing
// happens. Unregistering a function after Exit() has been called (i.e. in a function that was registered) will have no
// effect.
func Unregister(id int) {
	lock.Lock()
	defer lock.Unlock()
	pairs = slices.DeleteFunc(pairs, func(p pair) bool { return p.id == id })
}

// Exit runs any registered exit functions in the inverse order they were registered and then exits the progream with
// the specified status. If a previous call to Exit() is already being handled, this method does nothing but does not
// return. Recursive calls to Exit() will trigger a panic, which the exit handling will catch and report, but will then
// proceed with exit as normal. Note that once Exit() is called, no subsequent changes to the registered list of
// functions will have an effect (i.e. you cannot Unregister() a function inside an exit handler to prevent its
// execution, nor can you Register() a new function).
func Exit(status int) {
	var pcs [512]uintptr
	recursive := false
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if frame.Function == "github.com/richardwilkes/toolbox/atexit.Exit" {
			recursive = true
			break
		}
		if !more {
			break
		}
	}
	var f []func()
	lock.Lock()
	wasExiting := exiting
	if !wasExiting {
		exiting = true
		f = make([]func(), len(pairs))
		for i, p := range pairs {
			f[i] = p.f
		}
	}
	lock.Unlock()
	if wasExiting {
		if recursive {
			panic("recursive call of atexit.Exit()") // force the recovery mechanism to deal with it
		}
		select {} // halt progress so that we don't return
	} else {
		for i := len(f) - 1; i >= 0; i-- {
			run(f[i])
		}
		os.Exit(status)
	}
}

func run(f func()) {
	defer errs.Recovery(RecoveryHandler)
	f()
}
