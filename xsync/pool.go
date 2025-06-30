// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xsync

import "sync"

// Pool provides a type-safe wrapper around sync.Pool.
type Pool[T any] struct {
	internalPool sync.Pool
}

// NewPool creates a new, empty, Pool.
func NewPool[T any](newFunc func() T) (p Pool[T]) {
	if newFunc == nil {
		panic("newFunc must not be nil")
	}
	p.internalPool.New = func() any { return newFunc() }
	return
}

// Put adds x to the pool.
func (t *Pool[T]) Put(x T) {
	t.internalPool.Put(x)
}

// Get returns a value from the pool.
func (t *Pool[T]) Get() T {
	return t.internalPool.Get().(T) //nolint:errcheck // We know the type is correct
}
