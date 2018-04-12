package embedded

import (
	"bytes"
	"os"
	"path/filepath"
	"time"
)

// File holds the data for an embedded file.
type File struct {
	*bytes.Reader
	name    string
	size    int64
	modTime time.Time
	isDir   bool
	files   []os.FileInfo
	data    []byte
}

// NewFile creates a new embedded file.
func NewFile(name string, modTime time.Time, data []byte) File {
	return File{
		name:    filepath.Base(name),
		size:    int64(len(data)),
		modTime: modTime,
		data:    data,
	}
}

func (f *File) Close() error {
	return nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if f.isDir {
		return f.files, nil
	}
	return nil, os.ErrNotExist
}

func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int64 {
	return f.size
}

func (f *File) Mode() os.FileMode {
	if f.isDir {
		return 0555
	}
	return 0444
}

func (f *File) ModTime() time.Time {
	return f.modTime
}

func (f *File) IsDir() bool {
	return f.isDir
}

func (f *File) Sys() interface{} {
	return nil
}
