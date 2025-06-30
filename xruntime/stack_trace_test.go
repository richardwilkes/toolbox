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
	stack := nextLevelStackTrace(0)
	check.True(t, len(stack) > 1, "stack trace should have more than one entry")
	check.True(t, slices.ContainsFunc(stack, func(s string) bool { return strings.Contains(s, "TestStackTrace") }),
		"stack trace should contain TestStackTrace")
	check.True(t, slices.ContainsFunc(stack, func(s string) bool { return strings.Contains(s, "nextLevelStackTrace") }),
		"stack trace should contain nextLevelStackTrace")
}

func nextLevelStackTrace(skip int) []string {
	return xruntime.StackTrace(skip)
}

func TestStackTraceSkip(t *testing.T) {
	stack0 := nextLevelStackTrace(0)
	stack1 := nextLevelStackTrace(1)
	stack2 := nextLevelStackTrace(2)
	check.True(t, len(stack0) > len(stack1))
	check.True(t, len(stack1) > len(stack2))
}

func TestStackTracePath(t *testing.T) {
	for _, one := range []struct {
		function string
		file     string
		expected string
	}{
		{
			function: "github.com/user/play/internal/stuff.DoSomething",
			file:     "/Users/user/code/play/internal/stuff/stuff.go",
			expected: "play/internal/stuff/stuff.go",
		},
		{
			function: "main.main",
			file:     "/Users/user/code/play/main.go",
			expected: "main.go",
		},
		{
			function: "main.foo.bar",
			file:     "/Users/user/code/play/main.go",
			expected: "main.go",
		},
	} {
		check.Equal(t, one.expected, xruntime.StackTracePath(one.function, one.file))
	}
}
