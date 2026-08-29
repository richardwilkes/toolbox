// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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

	c.True(xreflect.IsNil(nil))

	c.False(xreflect.IsNil(42))
	c.False(xreflect.IsNil(42.0))
	c.False(xreflect.IsNil("hello"))
	c.False(xreflect.IsNil(true))
	c.False(xreflect.IsNil(struct{}{}))
	c.False(xreflect.IsNil(complex(0, 0)))

	var p *int
	c.True(xreflect.IsNil(p))
	var sp *string
	c.True(xreflect.IsNil(sp))
	var stp *struct{}
	c.True(xreflect.IsNil(stp))
	var up unsafe.Pointer
	c.True(xreflect.IsNil(up))

	n := 42
	c.False(xreflect.IsNil(&n))
	str := "hi"
	c.False(xreflect.IsNil(&str))
	var strct struct{}
	c.False(xreflect.IsNil(&strct))
	c.False(xreflect.IsNil(unsafe.Pointer(&n)))

	var s []int
	c.True(xreflect.IsNil(s))

	c.False(xreflect.IsNil([]int{1, 2, 3}))
	c.False(xreflect.IsNil(make([]int, 0)))

	var m map[string]int
	c.True(xreflect.IsNil(m))

	c.False(xreflect.IsNil(make(map[string]int)))

	var ch chan int
	c.True(xreflect.IsNil(ch))

	c.False(xreflect.IsNil(make(chan int)))

	var f func()
	c.True(xreflect.IsNil(f))

	c.False(xreflect.IsNil(func() {}))

	var err error
	c.True(xreflect.IsNil(err))

	var nilPtr *int
	var iface any = nilPtr
	c.True(xreflect.IsNil(iface))

	var nonNilIface any = 42
	c.False(xreflect.IsNil(nonNilIface))
}
