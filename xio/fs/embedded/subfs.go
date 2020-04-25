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
	"path"
)

type subFS struct {
	parent FileSystem
	base   string
}

// NewSubFileSystem creates a new FileSystem rooted at 'base' within an
// existing FileSystem.
func NewSubFileSystem(parent FileSystem, base string) FileSystem {
	return &subFS{
		parent: parent,
		base:   ToEFSPath(base),
	}
}

func (f *subFS) IsLive() bool {
	return f.parent.IsLive()
}

func (f *subFS) adjustPath(p string) string {
	return path.Join(f.base, ToEFSPath(p))
}

func (f *subFS) Open(p string) (http.File, error) {
	return f.parent.Open(f.adjustPath(p))
}

func (f *subFS) ContentAsBytes(p string) ([]byte, bool) {
	return f.parent.ContentAsBytes(f.adjustPath(p))
}

func (f *subFS) MustContentAsBytes(p string) []byte {
	return f.parent.MustContentAsBytes(f.adjustPath(p))
}

func (f *subFS) ContentAsString(p string) (string, bool) {
	return f.parent.ContentAsString(f.adjustPath(p))
}

func (f *subFS) MustContentAsString(p string) string {
	return f.parent.MustContentAsString(f.adjustPath(p))
}
