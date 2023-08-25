// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package softref

import (
	"log/slog"
	"runtime"
	"sync"
)

// Pool is used to track soft references to resources.
type Pool struct {
	lock sync.Mutex
	refs map[string]*softRef
}

// Resource is a resource that will be used with a pool.
type Resource interface {
	// Key returns a unique key for this resource. Must never change.
	Key() string
	// Release is called when the resource is no longer being referenced by any remaining soft references.
	Release()
}

// SoftRef is a soft reference to a given resource.
type SoftRef struct {
	Key      string
	Resource Resource
}

type softRef struct {
	resource Resource
	count    int
}

// DefaultPool is a global default soft reference pool.
var DefaultPool = NewPool()

// NewPool creates a new soft reference pool.
func NewPool() *Pool {
	return &Pool{refs: make(map[string]*softRef)}
}

// NewSoftRef returns a SoftRef to the given resource, along with a flag indicating if a reference existed previously.
func (p *Pool) NewSoftRef(resource Resource) (ref *SoftRef, existedPreviously bool) {
	key := resource.Key()
	p.lock.Lock()
	defer p.lock.Unlock()
	r := p.refs[key]
	if r != nil {
		r.count++
	} else {
		r = &softRef{
			resource: resource,
			count:    1,
		}
		p.refs[key] = r
	}
	sr := &SoftRef{
		Key:      key,
		Resource: r.resource,
	}
	runtime.SetFinalizer(sr, p.finalizeSoftRef)
	return sr, r.count > 1
}

func (p *Pool) finalizeSoftRef(ref *SoftRef) {
	p.lock.Lock()
	if r, ok := p.refs[ref.Key]; ok {
		r.count--
		if r.count == 0 {
			delete(p.refs, ref.Key)
			r.resource.Release()
		} else if r.count < 0 {
			slog.Debug("SoftRef count is invalid", "key", ref.Key, "count", r.count)
		}
	} else {
		slog.Debug("SoftRef finalized for unknown key", "key", ref.Key)
	}
	p.lock.Unlock()
}
