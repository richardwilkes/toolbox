package fs

import (
	"bytes"
	"path/filepath"
)

// SanitizeName sanitizes a file name by replacing invalid characters.
func SanitizeName(name string) string {
	if name == "" {
		return "@0"
	}
	if name == "." {
		return "@1"
	}
	if name == ".." {
		return "@2"
	}
	var buffer bytes.Buffer
	for _, r := range name {
		switch r {
		case '@':
			buffer.WriteString("@3")
		case '/':
			buffer.WriteString("@4")
		case '\\':
			buffer.WriteString("@5")
		case ':':
			buffer.WriteString("@6")
		default:
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

// UnsanitizeName reverses the effects of a call to SanitizeName.
func UnsanitizeName(name string) string {
	if name == "@0" {
		return ""
	}
	if name == "@1" {
		return "."
	}
	if name == "@2" {
		return ".."
	}
	var buffer bytes.Buffer
	found := false
	for _, r := range name {
		if found {
			switch r {
			case '3':
				buffer.WriteByte('@')
			case '4':
				buffer.WriteByte('/')
			case '5':
				buffer.WriteByte('\\')
			case '6':
				buffer.WriteByte(':')
			default:
				buffer.WriteByte('@')
				buffer.WriteRune(r)
			}
			found = false
		} else if r == '@' {
			found = true
		} else {
			buffer.WriteRune(r)
		}
	}
	if found {
		buffer.WriteByte('@')
	}
	return buffer.String()
}

// BaseName returns the file name without the directory or extension.
func BaseName(path string) string {
	path = filepath.Base(path)
	return path[:len(path)-len(filepath.Ext(path))]
}