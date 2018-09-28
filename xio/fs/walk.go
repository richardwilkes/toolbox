package fs

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/richardwilkes/toolbox/xio"
)

// Walk performs the same function as filepath.Walk() does, but works on
// http.FileSystem objects.
func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	info, err := stat(fs, root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walk(fs, root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}

func stat(fs http.FileSystem, path string) (os.FileInfo, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	xio.CloseIgnoringErrors(f)
	return info, err
}

func walk(fs http.FileSystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}
	names, err := readDirNames(fs, path)
	err1 := walkFn(path, info, err)
	if err != nil || err1 != nil {
		return err1
	}
	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := stat(fs, filename)
		if err != nil {
			if err = walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

func readDirNames(fs http.FileSystem, dirname string) ([]string, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	xio.CloseIgnoringErrors(f)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(list))
	for i := range list {
		names[i] = list[i].Name()
	}
	sort.Strings(names)
	return names, nil
}
