// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package rotation provides file rotation when files hit a given size.
package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/richardwilkes/toolbox/errs"
)

// Rotator holds the rotator data.
type Rotator struct {
	path       string
	maxSize    int64
	maxBackups int
	mask       os.FileMode
	lock       sync.Mutex
	file       *os.File
	size       int64
}

// New creates a new Rotator with the specified options.
func New(options ...func(*Rotator) error) (*Rotator, error) {
	r := &Rotator{
		path:       DefaultPath(),
		maxSize:    DefaultMaxSize,
		maxBackups: DefaultMaxBackups,
		mask:       0777,
	}
	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Write implements io.Writer.
func (r *Rotator) Write(b []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.file == nil {
		fi, err := os.Stat(r.path)
		switch {
		case os.IsNotExist(err):
			if err = os.MkdirAll(filepath.Dir(r.path), 0755&r.mask); err != nil {
				return 0, errs.Wrap(err)
			}
			file, fErr := os.OpenFile(r.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644&r.mask)
			if fErr != nil {
				return 0, errs.Wrap(fErr)
			}
			r.file = file
			r.size = 0
		case err != nil:
			return 0, errs.Wrap(err)
		default:
			var file *os.File
			if file, err = os.OpenFile(r.path, os.O_WRONLY|os.O_APPEND, 0644&r.mask); err != nil {
				return 0, errs.Wrap(err)
			}
			r.file = file
			r.size = fi.Size()
		}
	}
	writeSize := int64(len(b))
	if r.size+writeSize > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := r.file.Write(b)
	if err != nil {
		err = errs.Wrap(err)
	}
	r.size += int64(n)
	return n, err
}

// Close implements io.Closer.
func (r *Rotator) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.file != nil {
		file := r.file
		r.file = nil
		if err := file.Close(); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (r *Rotator) rotate() error {
	if r.file != nil {
		err := r.file.Close()
		r.file = nil
		if err != nil {
			return errs.Wrap(err)
		}
	}
	if r.maxBackups < 1 {
		if err := os.Remove(r.path); err != nil && !os.IsNotExist(err) {
			return errs.Wrap(err)
		}
	} else {
		if err := os.Remove(fmt.Sprintf("%s-%d", r.path, r.maxBackups)); err != nil && !os.IsNotExist(err) {
			return errs.Wrap(err)
		}
		for i := r.maxBackups; i > 0; i-- {
			var oldPath string
			if i != 1 {
				oldPath = fmt.Sprintf("%s-%d", r.path, i-1)
			} else {
				oldPath = r.path
			}
			if err := os.Rename(oldPath, fmt.Sprintf("%s-%d", r.path, i)); err != nil && !os.IsNotExist(err) {
				return errs.Wrap(err)
			}
		}
	}
	file, err := os.OpenFile(r.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644&r.mask)
	if err != nil {
		return errs.Wrap(err)
	}
	r.file = file
	r.size = 0
	return nil
}
