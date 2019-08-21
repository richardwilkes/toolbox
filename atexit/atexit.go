// Package atexit provides functionality similar to the C standard library's
// atexit() call.
package atexit

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	lock  sync.Mutex
	funcs []func()
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic during cleanup: %+v\n", err)
		}
	}()
	f()
}
