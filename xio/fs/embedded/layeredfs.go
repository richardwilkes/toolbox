// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package embedded

import (
	"net/http"
	"sync"
)

var _ FileSystem = &LayeredFS{}

// LayeredFS holds the contents of a layered file system.
type LayeredFS struct {
	primaries map[string]FileSystem
	fallback  FileSystem
	lock      sync.RWMutex
	primary   string
}

// NewLayeredFS creates a composite FileSystem. Multiple file systems may be
// designated as the potential primary and are chosen between dynamically at
// the time of a request by the current value of primary. Should the primary
// file system be unable to fulfill the request, then the request is passed to
// the fallback file system.
func NewLayeredFS(primary string, primaries map[string]FileSystem, fallback FileSystem) *LayeredFS {
	return &LayeredFS{
		primaries: primaries,
		fallback:  fallback,
		primary:   primary,
	}
}

// SetPrimary sets the primary filesystem
func (fs *LayeredFS) SetPrimary(primary string) {
	fs.lock.Lock()
	fs.primary = primary
	fs.lock.Unlock()
}

// Open a file
func (fs *LayeredFS) Open(name string) (http.File, error) {
	fs.lock.RLock()
	primary := fs.primary
	fs.lock.RUnlock()
	if fileSystem, ok := fs.primaries[primary]; ok {
		if f, err := fileSystem.Open(name); err == nil {
			return f, nil
		}
	}
	return fs.fallback.Open(name)
}

// IsLive returns true if the underlying filesystem is considered to be "live"
func (fs *LayeredFS) IsLive() bool {
	fs.lock.RLock()
	primary := fs.primary
	fs.lock.RUnlock()
	if fileSystem, ok := fs.primaries[primary]; ok {
		return fileSystem.IsLive()
	}
	return fs.fallback.IsLive()
}

// ContentAsBytes returns the file contents as bytes
func (fs *LayeredFS) ContentAsBytes(path string) ([]byte, bool) {
	fs.lock.RLock()
	primary := fs.primary
	fs.lock.RUnlock()
	if fileSystem, ok := fs.primaries[primary]; ok {
		var b []byte
		if b, ok = fileSystem.ContentAsBytes(path); ok {
			return b, ok
		}
	}
	return fs.fallback.ContentAsBytes(path)
}

// MustContentAsBytes returns the file contents as bytes, exiting if unable to
func (fs *LayeredFS) MustContentAsBytes(path string) []byte {
	fs.lock.RLock()
	primary := fs.primary
	fs.lock.RUnlock()
	if fileSystem, ok := fs.primaries[primary]; ok {
		var b []byte
		if b, ok = fileSystem.ContentAsBytes(path); ok {
			return b
		}
	}
	return fs.fallback.MustContentAsBytes(path)
}

// ContentAsString returns the file contents a string
func (fs *LayeredFS) ContentAsString(path string) (string, bool) {
	fs.lock.RLock()
	primary := fs.primary
	fs.lock.RUnlock()
	if fileSystem, ok := fs.primaries[primary]; ok {
		var str string
		if str, ok = fileSystem.ContentAsString(path); ok {
			return str, ok
		}
	}
	return fs.fallback.ContentAsString(path)
}

// MustContentAsString returns the file contents as a string, exiting if
// unable to
func (fs *LayeredFS) MustContentAsString(path string) string {
	fs.lock.RLock()
	primary := fs.primary
	fs.lock.RUnlock()
	if fileSystem, ok := fs.primaries[primary]; ok {
		var str string
		if str, ok = fileSystem.ContentAsString(path); ok {
			return str
		}
	}
	return fs.fallback.MustContentAsString(path)
}
