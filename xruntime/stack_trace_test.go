// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	c := check.New(t)
	for _, tc := range []struct {
		want []string
		omit []string
		skip int
	}{
		{want: []string{"nextLevelStackTrace", "TestStackTrace"}},
		{want: []string{"TestStackTrace"}, omit: []string{"nextLevelStackTrace"}, skip: 1},
		{omit: []string{"nextLevelStackTrace", "TestStackTrace"}, skip: 2},
	} {
		stack := nextLevelStackTrace(tc.skip)
		for _, name := range tc.want {
			c.True(slices.ContainsFunc(stack, func(s string) bool { return strings.Contains(s, name) }),
				"skip %d: stack trace should contain %s", tc.skip, name)
		}
		for _, name := range tc.omit {
			c.False(slices.ContainsFunc(stack, func(s string) bool { return strings.Contains(s, name) }),
				"skip %d: stack trace should not contain %s", tc.skip, name)
		}
	}
}

func nextLevelStackTrace(skip int) []string {
	return xruntime.StackTrace(skip)
}

func TestStackTracePath(t *testing.T) {
	c := check.New(t)
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
		c.Equal(one.expected, xruntime.StackTracePath(one.function, one.file))
	}
}
