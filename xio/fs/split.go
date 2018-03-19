package fs

import (
	"path/filepath"
)

// Split a path into its component parts.
func Split(path string) []string {
	var parts []string
	path = filepath.Clean(path)
	parts = append(parts, filepath.Base(path))
	sep := string(filepath.Separator)
	for {
		path = filepath.Dir(path)
		parts = append(parts, filepath.Base(path))
		if path == "." || path == sep {
			break
		}
	}
	result := make([]string, len(parts))
	for i := 0; i < len(parts); i++ {
		result[len(parts)-(i+1)] = parts[i]
	}
	return result
}
