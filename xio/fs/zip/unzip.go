package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// ExtractArchive extracts the contents of a zip archive at 'src' into the
// 'dst' directory.
func ExtractArchive(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	return Extract(&r.Reader, dst)
}

// Extract the contents of a zip reader into the 'dst' directory.
func Extract(zr *zip.Reader, dst string) error {
	root, err := filepath.Abs(dst)
	if err != nil {
		return errs.Wrap(err)
	}
	rootWithTrailingSep := fmt.Sprintf("%s%c", root, filepath.Separator)
	for _, f := range zr.File {
		path := filepath.Join(root, f.Name)
		if !strings.HasPrefix(path, rootWithTrailingSep) {
			return errs.Newf("Path outside of root is not permitted: %s", f.Name)
		}
		fi := f.FileInfo()
		mode := fi.Mode()
		if mode&os.ModeSymlink != 0 {
			if err := extractSymLink(f, path); err != nil {
				return err
			}
		} else if fi.IsDir() {
			if err := os.MkdirAll(path, mode.Perm()); err != nil {
				return errs.Wrap(err)
			}
		} else {
			if err := extractFile(f, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractSymLink(f *zip.File, dst string) error {
	r, err := f.Open()
	if err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	var buffer []byte
	if buffer, err = ioutil.ReadAll(r); err != nil {
		return errs.Wrap(err)
	}
	if err = os.MkdirAll(filepath.Dir(dst), 0775); err != nil {
		return errs.Wrap(err)
	}
	if err = os.Symlink(string(buffer), dst); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func extractFile(f *zip.File, dst string) (err error) {
	var r io.ReadCloser
	if r, err = f.Open(); err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	if err = os.MkdirAll(filepath.Dir(dst), 0775); err != nil {
		return errs.Wrap(err)
	}
	var file *os.File
	if file, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.FileInfo().Mode().Perm()); err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = errs.Wrap(cerr)
		}
	}()
	if _, err = io.Copy(file, r); err != nil {
		err = errs.Wrap(err)
	}
	return
}
