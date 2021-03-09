// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs

import (
	"io"
	"os"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// MoveFile moves a file in the file system or across volumes, using rename if
// possible, but falling back to copying the file if not. This will error if
// either src or dst are not regular files.
func MoveFile(src, dst string) (err error) {
	var srcInfo, dstInfo os.FileInfo
	srcInfo, err = os.Stat(src)
	if err != nil {
		return errs.Wrap(err)
	}
	if !srcInfo.Mode().IsRegular() {
		return errs.Newf("%s is not a regular file", src)
	}
	dstInfo, err = os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return errs.Wrap(err)
		}
	} else {
		if !dstInfo.Mode().IsRegular() {
			return errs.Newf("%s is not a regular file", dst)
		}
		if os.SameFile(srcInfo, dstInfo) {
			return nil
		}
	}
	if os.Rename(src, dst) == nil {
		return nil
	}
	var in, out *os.File
	out, err = os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if in, err = os.Open(src); err != nil {
		err = errs.Wrap(err)
		return
	}
	_, err = io.Copy(out, in)
	xio.CloseIgnoringErrors(in)
	if err != nil {
		err = errs.Wrap(err)
		return
	}
	if err = os.Remove(src); err != nil {
		err = errs.Wrap(err)
	}
	return
}
