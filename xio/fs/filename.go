// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs

import (
	"path/filepath"
	"strings"
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
	var buffer strings.Builder
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
	var buffer strings.Builder
	found := false
	for _, r := range name {
		switch {
		case found:
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
		case r == '@':
			found = true
		default:
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
	return TrimExtension(filepath.Base(path))
}

// TrimExtension trims any extension from the path.
func TrimExtension(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}
