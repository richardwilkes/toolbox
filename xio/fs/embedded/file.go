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
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

// File holds the data for an embedded file.
type File struct { //nolint:maligned
	*bytes.Reader
	name       string
	size       int64
	modTime    time.Time
	isDir      bool
	files      []os.FileInfo
	lock       sync.Mutex
	compressed bool
	data       []byte
}

// NewFile creates a new embedded file.
func NewFile(name string, modTime time.Time, size int64, compressed bool, data []byte) *File {
	return &File{
		name:       path.Base(ToEFSPath(name)),
		size:       size,
		modTime:    modTime,
		compressed: compressed,
		data:       data,
	}
}

func (f *File) uncompressData() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.compressed {
		r, err := gzip.NewReader(bytes.NewReader(f.data))
		if err != nil {
			return err
		}
		buffer := make([]byte, int(f.size))
		if _, err = io.ReadFull(r, buffer); err != nil {
			return err
		}
		if err = r.Close(); err != nil {
			return err
		}
		f.compressed = false
		f.data = buffer
	}
	return nil
}

// Close the file. Does nothing and always returns nil. Implements the
// io.Closer interface.
func (f *File) Close() error {
	return nil
}

// Readdir reads a directory and returns information about its contents.
// Implements the http.File interface.
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if f.isDir {
		return f.files, nil
	}
	return nil, os.ErrNotExist
}

// Stat returns information about the file. Implements the http.File
// interface.
func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

// Name returns the base name of the file. Implements the os.FileInfo
// interface.
func (f *File) Name() string {
	return f.name
}

// Size returns the size of the file in bytes. Implements the os.FileInfo
// interface.
func (f *File) Size() int64 {
	return f.size
}

// Mode returns the file mode bits. Implements the os.FileInfo interface.
func (f *File) Mode() os.FileMode {
	if f.isDir {
		return 0555
	}
	return 0444
}

// ModTime returns the file modification time. Implements the os.FileInfo
// interface.
func (f *File) ModTime() time.Time {
	return f.modTime
}

// IsDir returns true if this represents a directory. Implements the
// os.FileInfo interface.
func (f *File) IsDir() bool {
	return f.isDir
}

// Sys returns nil. Implements the os.FileInfo interface.
func (f *File) Sys() interface{} {
	return nil
}
