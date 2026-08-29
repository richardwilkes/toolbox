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
	"runtime"
	"testing"
	"unsafe"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xruntime"
)

// TestPtrFromUintptrHeapRoundTrip round-trips a pinned Go heap pointer through a uintptr, the pattern used when one of
// our own pinned Go objects is handed back to us by the OS. Under checkptr instrumentation (enabled by -race), a direct
// uintptr-to-unsafe.Pointer conversion of such a value aborts the test binary, so this also guards against regressing
// to an instrumented conversion.
func TestPtrFromUintptrHeapRoundTrip(t *testing.T) {
	c := check.New(t)
	var pin runtime.Pinner
	defer pin.Unpin()
	v := new(int)
	*v = 42
	pin.Pin(v)
	p := xruntime.PtrFromUintptr[int](uintptr(unsafe.Pointer(v)))
	c.True(p == v)
	c.Equal(42, *p)
	*p = 17
	c.Equal(17, *v)
}

// namedUintptr stands in for an OS handle alias, exercising the ~uintptr constraint.
type namedUintptr uintptr

func TestPtrFromUintptrNamedType(t *testing.T) {
	c := check.New(t)
	var pin runtime.Pinner
	defer pin.Unpin()
	v := new(uint32)
	*v = 0xDEADBEEF
	pin.Pin(v)
	c.Equal(uint32(0xDEADBEEF), *xruntime.PtrFromUintptr[uint32](namedUintptr(uintptr(unsafe.Pointer(v)))))
}
