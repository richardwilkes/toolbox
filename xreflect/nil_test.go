// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xreflect_test

import (
	"testing"
	"unsafe"

	"github.com/richardwilkes/toolbox/v2/xreflect"
)

func TestIsNil(t *testing.T) {
	// Test explicit nil
	checkTrue(t, xreflect.IsNil(nil))

	// Test non-nilable types
	checkFalse(t, xreflect.IsNil(42))
	checkFalse(t, xreflect.IsNil(42.0))
	checkFalse(t, xreflect.IsNil("hello"))
	checkFalse(t, xreflect.IsNil(true))
	checkFalse(t, xreflect.IsNil(struct{}{}))
	checkFalse(t, xreflect.IsNil(complex(0, 0)))

	// Test nil pointers
	var p *int
	checkTrue(t, xreflect.IsNil(p))
	var sp *string
	checkTrue(t, xreflect.IsNil(sp))
	var stp *struct{}
	checkTrue(t, xreflect.IsNil(stp))
	var up unsafe.Pointer
	checkTrue(t, xreflect.IsNil(up))

	// Test non-nil pointer
	n := 42
	checkFalse(t, xreflect.IsNil(&n))
	str := "hi"
	checkFalse(t, xreflect.IsNil(&str))
	var strct struct{}
	checkFalse(t, xreflect.IsNil(&strct))
	checkFalse(t, xreflect.IsNil(unsafe.Pointer(&n)))

	// Test nil slice
	var s []int
	checkTrue(t, xreflect.IsNil(s))

	// Test non-nil slice
	checkFalse(t, xreflect.IsNil([]int{1, 2, 3}))
	checkFalse(t, xreflect.IsNil(make([]int, 0)))

	// Test nil map
	var m map[string]int
	checkTrue(t, xreflect.IsNil(m))

	// Test non-nil map
	checkFalse(t, xreflect.IsNil(make(map[string]int)))

	// Test nil channel
	var ch chan int
	checkTrue(t, xreflect.IsNil(ch))

	// Test non-nil channel
	checkFalse(t, xreflect.IsNil(make(chan int)))

	// Test nil function
	var f func()
	checkTrue(t, xreflect.IsNil(f))

	// Test non-nil function
	checkFalse(t, xreflect.IsNil(func() {}))

	// Test nil interface
	var err error
	checkTrue(t, xreflect.IsNil(err))

	// Test interface with nil pointer
	var nilPtr *int
	var iface any = nilPtr
	checkTrue(t, xreflect.IsNil(iface))

	// Test interface with non-nil value
	var nonNilIface any = 42
	checkFalse(t, xreflect.IsNil(nonNilIface))
}

func checkTrue(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Error("Expected true")
	}
}

func checkFalse(t *testing.T, value bool) {
	t.Helper()
	if value {
		t.Error("Expected false")
	}
}
