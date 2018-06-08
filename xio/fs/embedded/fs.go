package embedded

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/richardwilkes/toolbox/collection"
	"github.com/richardwilkes/toolbox/txt"
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
	// Generate immediate directories for files
	now := time.Now()
	all := make(map[string]*File)
	for k, v := range files {
		all[k] = v
	}
	type dinfo struct {
		f *File
		m collection.StringSet
	}
	dirs := make(map[string]*dinfo)
	for k, v := range files {
		dir, _ := filepath.Split(k)
		dir = filepath.Clean(dir)
		di, ok := dirs[dir]
		if !ok {
			di = &dinfo{
				f: &File{
					name:    filepath.Base(dir),
					modTime: now,
					isDir:   true,
				},
				m: collection.NewStringSet(),
			}
			dirs[dir] = di
		}
		di.f.files = append(di.f.files, v)
		// Ensure parents are present
		path := dir
		for {
			if dir, _ = filepath.Split(dir); dir == "" || path == dir {
				break
			}
			dir = filepath.Clean(dir)
			p, ok := dirs[dir]
			if !ok {
				p = &dinfo{
					f: &File{
						name:    filepath.Base(dir),
						modTime: now,
						isDir:   true,
					},
					m: collection.NewStringSet(),
				}
				dirs[dir] = p
			}
			if p.m.Contains(path) {
				break
			}
			p.m.Add(path)
			p.f.files = append(p.f.files, di.f)
			di = p
			path = dir
		}
	}
	// For each dir, sort its file list and add it to our "all" list
	for k, v := range dirs {
		sort.Slice(v.f.files, func(i, j int) bool {
			return txt.NaturalLess(v.f.files[i].Name(), v.f.files[j].Name(), true)
		})
		all[k] = v.f
	}
	return &EFS{
		efs: &efs{
			files:      all,
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
