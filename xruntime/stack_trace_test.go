// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xruntime_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xruntime"
)

func TestStackTrace(t *testing.T) {
	stack := xruntime.StackTrace(0)
	check.True(t, len(stack) > 1, "stack trace should have more than one entry")
	check.True(t, slices.ContainsFunc(stack, func(s string) bool { return strings.Contains(s, "TestStackTrace") }),
		"stack trace should contain TestStackTrace")
	check.True(t, slices.ContainsFunc(stack, func(s string) bool { return strings.Contains(s, "testing.tRunner") }),
		"stack trace should contain testing.tRunner")
}

func TestStackTraceSkip(t *testing.T) {
	stack0 := xruntime.StackTrace(0)
	stack1 := xruntime.StackTrace(1)
	stack2 := xruntime.StackTrace(2)
	check.True(t, len(stack0) > len(stack1))
	check.True(t, len(stack1) > len(stack2))
}

func TestStackTracePath(t *testing.T) {
	check.Equal(t, "xruntime/stack_trace.go",
		xruntime.StackTracePath("/Users/user/go/src/github.com/richardwilkes/toolbox/xruntime/stack_trace.go"),
		"full path with more than one directory")
	check.Equal(t, "xruntime/stack_trace.go", xruntime.StackTracePath("toolbox/xruntime/stack_trace.go"),
		"relative path with more than one directory")
	check.Equal(t, "home/main.go", xruntime.StackTracePath("/home/main.go"), "path with one directory")
	check.Equal(t, "main.go", xruntime.StackTracePath("main.go"), "path with just filename")
	check.Equal(t, "", xruntime.StackTracePath(""), "empty path")
	check.Equal(t, "dir/", xruntime.StackTracePath("/some/path/dir/"), "path ending with slash")
	check.Equal(t, "/file.go", xruntime.StackTracePath("/file.go"), "path with file at root")
}
