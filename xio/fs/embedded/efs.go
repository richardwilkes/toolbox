package embedded

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type efs struct {
	files      map[string]*File
	dirModTime time.Time
}

func (f *efs) IsLive() bool {
	return false
}

func (f *efs) Open(path string) (http.File, error) {
	path = filepath.Clean("/" + path)
	one, ok := f.files[path]
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
	if one, ok := f.files[filepath.Clean("/"+path)]; ok {
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
	panic(path + " does not exist") // @allow
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
	panic(path + " does not exist") // @allow
}
