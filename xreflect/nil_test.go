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

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xreflect"
)

func TestIsNil(t *testing.T) {
	// Test explicit nil
	check.True(t, xreflect.IsNil(nil))

	// Test non-nilable types
	check.False(t, xreflect.IsNil(42))
	check.False(t, xreflect.IsNil(42.0))
	check.False(t, xreflect.IsNil("hello"))
	check.False(t, xreflect.IsNil(true))
	check.False(t, xreflect.IsNil(struct{}{}))
	check.False(t, xreflect.IsNil(complex(0, 0)))

	// Test nil pointers
	var p *int
	check.True(t, xreflect.IsNil(p))
	var sp *string
	check.True(t, xreflect.IsNil(sp))
	var stp *struct{}
	check.True(t, xreflect.IsNil(stp))
	var up unsafe.Pointer
	check.True(t, xreflect.IsNil(up))

	// Test non-nil pointer
	n := 42
	check.False(t, xreflect.IsNil(&n))
	str := "hi"
	check.False(t, xreflect.IsNil(&str))
	var strct struct{}
	check.False(t, xreflect.IsNil(&strct))
	check.False(t, xreflect.IsNil(unsafe.Pointer(&n)))

	// Test nil slice
	var s []int
	check.True(t, xreflect.IsNil(s))

	// Test non-nil slice
	check.False(t, xreflect.IsNil([]int{1, 2, 3}))
	check.False(t, xreflect.IsNil(make([]int, 0)))

	// Test nil map
	var m map[string]int
	check.True(t, xreflect.IsNil(m))

	// Test non-nil map
	check.False(t, xreflect.IsNil(make(map[string]int)))

	// Test nil channel
	var ch chan int
	check.True(t, xreflect.IsNil(ch))

	// Test non-nil channel
	check.False(t, xreflect.IsNil(make(chan int)))

	// Test nil function
	var f func()
	check.True(t, xreflect.IsNil(f))

	// Test non-nil function
	check.False(t, xreflect.IsNil(func() {}))

	// Test nil interface
	var err error
	check.True(t, xreflect.IsNil(err))

	// Test interface with nil pointer
	var nilPtr *int
	var iface any = nilPtr
	check.True(t, xreflect.IsNil(iface))

	// Test interface with non-nil value
	var nonNilIface any = 42
	check.False(t, xreflect.IsNil(nonNilIface))
}
