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
	nextID          = 1
	pairs           []pair
)

type pair struct {
	id int
	f  func()
}

// Register a function to be run at exit. Returns an ID that can be used to
// remove the function later, if desired.
func Register(f func()) int {
	lock.Lock()
	defer lock.Unlock()
	if nextID == 1 {
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
	pairs = append(pairs, pair{id: nextID, f: f})
	nextID++
	return nextID - 1
}

// Unregister a function that was previously registered to be run at exit. If
// the ID is no longer present, nothing happens.
func Unregister(id int) {
	lock.Lock()
	defer lock.Unlock()
	for i := range pairs {
		if pairs[i].id == id {
			if i < len(pairs)-1 {
				copy(pairs[i:], pairs[i+1:])
			}
			pairs = pairs[:len(pairs)-1]
		}
	}
}

// Exit runs any registered exit functions in the inverse order they were
// registered and then exits with the specified status.
func Exit(status int) {
	lock.Lock() // Intentionally don't unlock. Prevents secondary calls to Exit from causing early exits.
	for i := len(pairs) - 1; i >= 0; i-- {
		run(pairs[i].f)
	}
	os.Exit(status)
}

func run(f func()) {
	defer errs.Recovery(RecoveryHandler)
	f()
}
