package fs

import "os"

// FileExists returns true if the path points to a regular file.
func FileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		mode := fi.Mode()
		return !mode.IsDir() && mode.IsRegular()
	}
	return false
}
