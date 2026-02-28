// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xruntime

import (
	"fmt"
	"path"
	"runtime"
	"slices"
	"strings"
)

// StackFuncPrefixesToFilter is a list of function prefixes to filter out of the stack trace.
//
// This variable is not used in a thread-safe manner, so any alterations should be done before any goroutines are
// started.
var StackFuncPrefixesToFilter = []string{
	"runtime.",
	"testing.",
	"github.com/richardwilkes/toolbox/v2/xos.PanicRecovery",
	"github.com/richardwilkes/toolbox/v2/errs.New",
	"github.com/richardwilkes/toolbox/v2/errs.Wrap",
}

// StackTrace returns a slice of strings, each of which is a text representation of a frame in the stack trace. 'skip'
// is the number of stack frames to skip before processing, with 0 identifying the caller of StackTrace.
func StackTrace(skip int) []string {
	var pcs [128]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	return PCsToStackTrace(pcs[:n])
}

// PCsToStackTrace converts a slice of program counters (PCs) into a slice of strings representing the stack trace.
func PCsToStackTrace(pcs []uintptr) []string {
	stack := make([]string, 0, len(pcs))
	frames := runtime.CallersFrames(pcs)
	more := true
	for more {
		var f runtime.Frame
		if f, more = frames.Next(); f.PC != 0 && (f.Function != "main.main" || f.File != "_testmain.go") &&
			!slices.ContainsFunc(StackFuncPrefixesToFilter, func(prefix string) bool {
				return strings.HasPrefix(f.Function, prefix)
			}) {
			stack = append(stack, fmt.Sprintf("[%s] %s:%d", f.Function, StackTracePath(f.Function, f.File), f.Line))
		}
	}
	return stack
}

// StackTracePath returns a shortened path for the function and file that should still be unique enough for an IDE to
// identify the location within the project.
func StackTracePath(function, file string) string {
	dirs := strings.Split(path.Dir(file), "/")
	functions := strings.Split(function, "/")
	if len(functions) > 0 {
		functions[len(functions)-1] = strings.TrimSuffix(
			strings.SplitN(functions[len(functions)-1], ".", 2)[0], "_test")
	}
	i := len(functions)
	j := len(dirs)
	for i > 0 && j > 0 && functions[i-1] == dirs[j-1] {
		i--
		j--
	}
	return strings.Join(append(dirs[j:], path.Base(file)), "/")
}
