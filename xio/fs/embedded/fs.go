package embedded

import (
	"net/http"
	"os"
	"time"
)

// FileSystem defines the methods available for a live or embedded filesystem.
type FileSystem interface {
	http.FileSystem
	IsLive() bool
	ContentAsBytes(path string) ([]byte, bool)
	MustContentAsBytes(path string) []byte
	ContentAsString(path string) (string, bool)
	MustContentAsString(path string) string
}

// EFS holds an embedded filesystem.
type EFS struct {
	efs FileSystem
}

// NewEFS creates a new embedded filesystem.
func NewEFS(files map[string]*File) *EFS {
	return &EFS{
		efs: &efs{
			files:      files,
			dirModTime: time.Now(),
		},
	}
}

// FileSystem returns either the embedded filesystem or a live filesystem
// rooted at localRoot if localRoot isn't an empty string and points to a
// directory.
func (efs *EFS) FileSystem(localRoot string) FileSystem {
	if localRoot != "" {
		if fi, err := os.Stat(localRoot); err == nil && fi.IsDir() {
			return &livefs{base: localRoot}
		}
	}
	return efs.efs
}
