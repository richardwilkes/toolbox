// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package embedded provides an implementation of an embedded filesystem.
package embedded

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/richardwilkes/toolbox/collection"
	"github.com/richardwilkes/toolbox/txt"
)

// FileSystem defines the methods available for a live or embedded filesystem.
//
// Deprecated: use Go 1.16's embedded support instead
type FileSystem interface {
	http.FileSystem
	IsLive() bool
	ContentAsBytes(path string) ([]byte, bool)
	MustContentAsBytes(path string) []byte
	ContentAsString(path string) (string, bool)
	MustContentAsString(path string) string
}

// EFS holds an embedded filesystem.
//
// Deprecated: use Go 1.16's embedded support instead
type EFS struct {
	efs FileSystem
}

// NewEFS creates a new embedded filesystem.
//
// Deprecated: use Go 1.16's embedded support instead
func NewEFS(files map[string]*File) *EFS {
	// Generate immediate directories for files
	now := time.Now()
	all := make(map[string]*File)
	for k, v := range files {
		all[k] = v
	}
	type dInfo struct {
		f *File
		m collection.StringSet
	}
	dirs := make(map[string]*dInfo)
	for k, v := range files {
		dir, _ := filepath.Split(k)
		dir = ToEFSPath(dir)
		di, ok := dirs[dir]
		if !ok {
			di = &dInfo{
				f: &File{
					name:    path.Base(dir),
					modTime: now,
					isDir:   true,
				},
				m: collection.NewStringSet(),
			}
			dirs[dir] = di
		}
		di.f.files = append(di.f.files, v)
		// Ensure parents are present
		parent := dir
		for {
			if dir, _ = path.Split(dir); dir == "" || parent == dir {
				break
			}
			dir = path.Clean(dir)
			var p *dInfo
			if p, ok = dirs[dir]; !ok {
				p = &dInfo{
					f: &File{
						name:    path.Base(dir),
						modTime: now,
						isDir:   true,
					},
					m: collection.NewStringSet(),
				}
				dirs[dir] = p
			}
			if p.m.Contains(parent) {
				break
			}
			p.m.Add(parent)
			p.f.files = append(p.f.files, di.f)
			di = p
			parent = dir
		}
	}
	// For each dir, sort its file list and add it to our "all" list
	for k, v := range dirs {
		sort.Slice(v.f.files, func(i, j int) bool {
			return txt.NaturalLess(v.f.files[i].Name(), v.f.files[j].Name(), true) //nolint:scopelint
		})
		all[k] = v.f
	}
	return &EFS{
		efs: &efs{
			files:      all,
			dirModTime: time.Now(),
		},
	}
}

// PrimaryFileSystem returns the primary filesystem this EFS represents.
func (efs *EFS) PrimaryFileSystem() FileSystem {
	return efs.efs
}

// FileSystem returns either the embedded filesystem or a live filesystem
// rooted at localRoot if localRoot isn't an empty string and points to a
// directory.
func (efs *EFS) FileSystem(localRoot string) FileSystem {
	if localRoot != "" {
		if fi, err := os.Stat(localRoot); err == nil && fi.IsDir() {
			return NewLiveFS(localRoot)
		}
	}
	return efs.PrimaryFileSystem()
}
