// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package toolbox

import (
	"testing"
	"unsafe"
)

func TestIsNil(t *testing.T) {
	// Test explicit nil
	checkTrue(t, IsNil(nil))

	// Test non-nilable types
	checkFalse(t, IsNil(42))
	checkFalse(t, IsNil(42.0))
	checkFalse(t, IsNil("hello"))
	checkFalse(t, IsNil(true))
	checkFalse(t, IsNil(struct{}{}))
	checkFalse(t, IsNil(complex(0, 0)))

	// Test nil pointers
	var p *int
	checkTrue(t, IsNil(p))
	var sp *string
	checkTrue(t, IsNil(sp))
	var stp *struct{}
	checkTrue(t, IsNil(stp))
	var up unsafe.Pointer
	checkTrue(t, IsNil(up))

	// Test non-nil pointer
	n := 42
	checkFalse(t, IsNil(&n))
	str := "hi"
	checkFalse(t, IsNil(&str))
	var strct struct{}
	checkFalse(t, IsNil(&strct))
	checkFalse(t, IsNil(unsafe.Pointer(&n)))

	// Test nil slice
	var s []int
	checkTrue(t, IsNil(s))

	// Test non-nil slice
	checkFalse(t, IsNil([]int{1, 2, 3}))
	checkFalse(t, IsNil(make([]int, 0)))

	// Test nil map
	var m map[string]int
	checkTrue(t, IsNil(m))

	// Test non-nil map
	checkFalse(t, IsNil(make(map[string]int)))

	// Test nil channel
	var ch chan int
	checkTrue(t, IsNil(ch))

	// Test non-nil channel
	checkFalse(t, IsNil(make(chan int)))

	// Test nil function
	var f func()
	checkTrue(t, IsNil(f))

	// Test non-nil function
	checkFalse(t, IsNil(func() {}))

	// Test nil interface
	var err error
	checkTrue(t, IsNil(err))

	// Test interface with nil pointer
	var nilPtr *int
	var iface any = nilPtr
	checkTrue(t, IsNil(iface))

	// Test interface with non-nil value
	var nonNilIface any = 42
	checkFalse(t, IsNil(nonNilIface))
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
