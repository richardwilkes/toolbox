package embedded

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type efs struct {
	files      map[string]File
	dirModTime time.Time
}

func (f *efs) IsLive() bool {
	return false
}

func (f *efs) Open(path string) (http.File, error) {
	path = filepath.Clean("/" + path)
	one, ok := f.files[path]
	if !ok {
		var files []os.FileInfo
		for k, v := range f.files {
			if strings.HasPrefix(k, path) {
				files = append(files, &v)
			}
		}
		if len(files) == 0 {
			return nil, os.ErrNotExist
		}
		return &File{
			name:    filepath.Base(path),
			modTime: f.dirModTime,
			isDir:   true,
			files:   files,
		}, nil
	}
	one.Reader = bytes.NewReader(one.data)
	return &one, nil
}

func (f *efs) ContentAsBytes(path string) ([]byte, bool) {
	if one, ok := f.files[filepath.Clean("/"+path)]; ok {
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
