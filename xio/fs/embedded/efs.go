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
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/log/jot"
)

type efs struct {
	files      map[string]*File
	dirModTime time.Time
}

func (f *efs) IsLive() bool {
	return false
}

func (f *efs) Open(path string) (http.File, error) {
	one, ok := f.files[ToEFSPath(path)]
	if !ok {
		return nil, os.ErrNotExist
	}
	if one.isDir {
		return one, nil
	}
	if err := one.uncompressData(); err != nil {
		return nil, err
	}
	return &File{
		Reader:  bytes.NewReader(one.data),
		name:    one.name,
		size:    one.size,
		modTime: one.modTime,
		data:    one.data,
	}, nil
}

func (f *efs) ContentAsBytes(path string) ([]byte, bool) {
	if one, ok := f.files[ToEFSPath(path)]; ok {
		if err := one.uncompressData(); err != nil {
			return nil, false
		}
		return one.data, true
	}
	return nil, false
}

func (f *efs) MustContentAsBytes(path string) []byte {
	if d, ok := f.ContentAsBytes(path); ok {
		return d
	}
	jot.Fatal(1, path+" does not exist")
	return nil
}

func (f *efs) ContentAsString(path string) (string, bool) {
	if d, ok := f.ContentAsBytes(path); ok {
		return string(d), true
	}
	return "", false
}

func (f *efs) MustContentAsString(path string) string {
	if s, ok := f.ContentAsString(path); ok {
		return s
	}
	jot.Fatal(1, path+" does not exist")
	return ""
}

// ToEFSPath converts a native file system path into one used by the EFS.
func ToEFSPath(path string) string {
	path = filepath.ToSlash(filepath.Clean(path))
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
