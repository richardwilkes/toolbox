// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package atexit provides functionality similar to the C standard library's
// atexit() call.
package atexit

import (
	"fmt"
	"log" //nolint:depguard
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/richardwilkes/toolbox/errs"
)

var (
	// RecoveryHandler will be used to capture any panics caused by functions
	// that have been installed when run during exit. It may be set to nil to
	// silently ignore them.
	RecoveryHandler errs.RecoveryHandler = func(err error) { log.Println(err) }
	lock            sync.Mutex
	funcs           []func()
)

// Register a function to be run at exit.
func Register(f func()) {
	lock.Lock()
	defer lock.Unlock()
	if len(funcs) == 0 {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			s := <-sigChan
			if s == syscall.SIGINT {
				fmt.Print("\b\b") // Removes the unsightly ^C in the terminal
			}
			Exit(1)
		}()
	}
	funcs = append(funcs, f)
}

// Exit runs any registered exit functions in the inverse order they were
// registered and then exits with the specified status.
func Exit(status int) {
	lock.Lock() // Intentionally don't unlock. Prevents secondary calls to Exit from causing early exits.
	all := make([]func(), len(funcs))
	copy(all, funcs)
	for i := len(all) - 1; i >= 0; i-- {
		run(all[i])
	}
	os.Exit(status)
}

func run(f func()) {
	defer errs.Recovery(RecoveryHandler)
	f()
}
