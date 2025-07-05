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
	c := check.New(t)

	// Test explicit nil
	c.True(xreflect.IsNil(nil))

	// Test non-nilable types
	c.False(xreflect.IsNil(42))
	c.False(xreflect.IsNil(42.0))
	c.False(xreflect.IsNil("hello"))
	c.False(xreflect.IsNil(true))
	c.False(xreflect.IsNil(struct{}{}))
	c.False(xreflect.IsNil(complex(0, 0)))

	// Test nil pointers
	var p *int
	c.True(xreflect.IsNil(p))
	var sp *string
	c.True(xreflect.IsNil(sp))
	var stp *struct{}
	c.True(xreflect.IsNil(stp))
	var up unsafe.Pointer
	c.True(xreflect.IsNil(up))

	// Test non-nil pointer
	n := 42
	c.False(xreflect.IsNil(&n))
	str := "hi"
	c.False(xreflect.IsNil(&str))
	var strct struct{}
	c.False(xreflect.IsNil(&strct))
	c.False(xreflect.IsNil(unsafe.Pointer(&n)))

	// Test nil slice
	var s []int
	c.True(xreflect.IsNil(s))

	// Test non-nil slice
	c.False(xreflect.IsNil([]int{1, 2, 3}))
	c.False(xreflect.IsNil(make([]int, 0)))

	// Test nil map
	var m map[string]int
	c.True(xreflect.IsNil(m))

	// Test non-nil map
	c.False(xreflect.IsNil(make(map[string]int)))

	// Test nil channel
	var ch chan int
	c.True(xreflect.IsNil(ch))

	// Test non-nil channel
	c.False(xreflect.IsNil(make(chan int)))

	// Test nil function
	var f func()
	c.True(xreflect.IsNil(f))

	// Test non-nil function
	c.False(xreflect.IsNil(func() {}))

	// Test nil interface
	var err error
	c.True(xreflect.IsNil(err))

	// Test interface with nil pointer
	var nilPtr *int
	var iface any = nilPtr
	c.True(xreflect.IsNil(iface))

	// Test interface with non-nil value
	var nonNilIface any = 42
	c.False(xreflect.IsNil(nonNilIface))
}
