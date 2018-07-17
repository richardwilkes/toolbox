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
	out, err = os.Create(dst)
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
