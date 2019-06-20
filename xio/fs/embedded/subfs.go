package embedded

import (
	"net/http"
	"path/filepath"
	"strings"
)

type subfs struct {
	parent FileSystem
	base   string
}

// NewSubFileSystem creates a new FileSystem rooted at 'base' within an
// existing FileSystem.
func NewSubFileSystem(parent FileSystem, base string) FileSystem {
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	return &subfs{
		parent: parent,
		base:   filepath.Clean(base),
	}
}

func (f *subfs) IsLive() bool {
	return f.parent.IsLive()
}

func (f *subfs) adjustPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return filepath.Join(f.base, filepath.Clean(path))
}

func (f *subfs) Open(path string) (http.File, error) {
	return f.parent.Open(f.adjustPath(path))
}

func (f *subfs) ContentAsBytes(path string) ([]byte, bool) {
	return f.parent.ContentAsBytes(f.adjustPath(path))
}

func (f *subfs) MustContentAsBytes(path string) []byte {
	return f.parent.MustContentAsBytes(f.adjustPath(path))
}

func (f *subfs) ContentAsString(path string) (string, bool) {
	return f.parent.ContentAsString(f.adjustPath(path))
}

func (f *subfs) MustContentAsString(path string) string {
	return f.parent.MustContentAsString(f.adjustPath(path))
}
