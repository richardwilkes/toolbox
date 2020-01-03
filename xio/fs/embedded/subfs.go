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
	"path/filepath"
)

type subfs struct {
	parent FileSystem
	base   string
}

// NewSubFileSystem creates a new FileSystem rooted at 'base' within an
// existing FileSystem.
func NewSubFileSystem(parent FileSystem, base string) FileSystem {
	base = filepath.Clean(base)
	if !filepath.IsAbs(base) {
		base = "/" + base
	}
	return &subfs{
		parent: parent,
		base:   base,
	}
}

func (f *subfs) IsLive() bool {
	return f.parent.IsLive()
}

func (f *subfs) adjustPath(path string) string {
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		path = "/" + path
	}
	return filepath.Join(f.base, path)
}

func (f *subfs) Open(path string) (http.File, error) {
	return f.parent.Open(f.adjustPath(path))
}

func (f *subfs) ContentAsBytes(path string) ([]byte, bool) {
	return f.parent.ContentAsBytes(f.adjustPath(path))
}

func (f *subfs) MustContentAsBytes(path string) []byte {
	return f.parent.MustContentAsBytes(f.adjustPath(path))
}

func (f *subfs) ContentAsString(path string) (string, bool) {
	return f.parent.ContentAsString(f.adjustPath(path))
}

func (f *subfs) MustContentAsString(path string) string {
	return f.parent.MustContentAsString(f.adjustPath(path))
}
