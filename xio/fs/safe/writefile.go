// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package safe

import (
	"bufio"
	"io"
	"os"
)

// WriteFile uses writer to write data safely and atomically to a file.
func WriteFile(filename string, writer func(io.Writer) error) (err error) {
	return WriteFileWithMode(filename, writer, 0o644)
}

// WriteFileWithMode uses writer to write data safely and atomically to a file.
func WriteFileWithMode(filename string, writer func(io.Writer) error, mode os.FileMode) (err error) {
	var f *File
	f, err = CreateWithMode(filename, mode)
	if err != nil {
		return
	}
	w := bufio.NewWriterSize(f, 1<<16)
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if err = writer(w); err != nil {
		return
	}
	if err = w.Flush(); err != nil {
		return
	}
	if err = f.Commit(); err != nil {
		return
	}
	return
}
