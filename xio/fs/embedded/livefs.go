package embedded

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type livefs struct {
	base string
}

func (f *livefs) IsLive() bool {
	return true
}

func (f *livefs) Open(path string) (http.File, error) {
	return os.Open(f.actualPath(path))
}

func (f *livefs) actualPath(path string) string {
	return filepath.Join(f.base, filepath.FromSlash(filepath.Clean("/"+path)))
}

func (f *livefs) ContentAsBytes(path string) ([]byte, bool) {
	if d, err := ioutil.ReadFile(f.actualPath(path)); err == nil {
		return d, true
	}
	return nil, false
}

func (f *livefs) MustContentAsBytes(path string) []byte {
	if d, ok := f.ContentAsBytes(path); ok {
		return d
	}
	panic(path + " does not exist") // @allow
}

func (f *livefs) ContentAsString(path string) (string, bool) {
	if d, ok := f.ContentAsBytes(path); ok {
		return string(d), true
	}
	return "", false
}

func (f *livefs) MustContentAsString(path string) string {
	if s, ok := f.ContentAsString(path); ok {
		return s
	}
	panic(path + " does not exist") // @allow
}
