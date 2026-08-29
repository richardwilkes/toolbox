// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xfilepath_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xfilepath"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestSanitizeName(t *testing.T) {
	c := check.New(t)

	c.Equal("@0", xfilepath.SanitizeName(""))
	c.Equal("@1", xfilepath.SanitizeName("."))
	c.Equal("@2", xfilepath.SanitizeName(".."))

	c.Equal("@3", xfilepath.SanitizeName("@"))
	c.Equal("@4", xfilepath.SanitizeName("/"))
	c.Equal("@5", xfilepath.SanitizeName("\\"))
	c.Equal("@6", xfilepath.SanitizeName(":"))

	c.Equal("normal_file.txt", xfilepath.SanitizeName("normal_file.txt"))

	c.Equal("file@4with@5special@6chars@3here", xfilepath.SanitizeName("file/with\\special:chars@here"))

	c.Equal("my file (1).txt", xfilepath.SanitizeName("my file (1).txt"))
}

func TestUnsanitizeName(t *testing.T) {
	c := check.New(t)

	c.Equal("", xfilepath.UnsanitizeName("@0"))
	c.Equal(".", xfilepath.UnsanitizeName("@1"))
	c.Equal("..", xfilepath.UnsanitizeName("@2"))

	c.Equal("@", xfilepath.UnsanitizeName("@3"))
	c.Equal("/", xfilepath.UnsanitizeName("@4"))
	c.Equal("\\", xfilepath.UnsanitizeName("@5"))
	c.Equal(":", xfilepath.UnsanitizeName("@6"))

	c.Equal("normal_file.txt", xfilepath.UnsanitizeName("normal_file.txt"))

	c.Equal("file/with\\special:chars@here", xfilepath.UnsanitizeName("file@4with@5special@6chars@3here"))

	c.Equal("filename@", xfilepath.UnsanitizeName("filename@"))

	c.Equal("@x", xfilepath.UnsanitizeName("@x"))
	c.Equal("@9", xfilepath.UnsanitizeName("@9"))
}

func TestSanitizeUnsanitizeRoundTrip(t *testing.T) {
	c := check.New(t)

	testCases := []string{
		"",
		".",
		"..",
		"normal_file.txt",
		"file/with\\special:chars@here",
		"my file (1).txt",
		"@",
		"/",
		"\\",
		":",
		"complex@file/name\\with:many@special@chars",
		"file@end",
		"@start",
		"mid@dle",
	}

	for _, original := range testCases {
		sanitized := xfilepath.SanitizeName(original)
		restored := xfilepath.UnsanitizeName(sanitized)
		c.Equal(original, restored, "Round trip failed for: %q", original)
	}
}

func TestBaseName(t *testing.T) {
	c := check.New(t)

	c.Equal("file", xfilepath.BaseName("file"))
	c.Equal("file", xfilepath.BaseName("file.txt"))
	c.Equal("file", xfilepath.BaseName("/path/to/file.txt"))
	if runtime.GOOS == xos.WindowsOS {
		c.Equal("file", xfilepath.BaseName("C:\\path\\to\\file.txt"))
	}

	c.Equal("file.backup", xfilepath.BaseName("file.backup.txt"))
	c.Equal("archive.tar", xfilepath.BaseName("/home/user/archive.tar.gz"))

	c.Equal("README", xfilepath.BaseName("README"))
	c.Equal("Makefile", xfilepath.BaseName("/project/Makefile"))

	// A leading dot marks a hidden file, not an extension.
	c.Equal(".txt", xfilepath.BaseName(".txt"))
	c.Equal(".hidden", xfilepath.BaseName(".hidden"))
	c.Equal(".profile", xfilepath.BaseName("/etc/.profile"))
	c.Equal(".hidden", xfilepath.BaseName("dir/.hidden.txt"))
	c.Equal(".config", xfilepath.BaseName(".config.yaml"))

	c.Equal(".", xfilepath.BaseName("."))
	c.Equal(".", xfilepath.BaseName(".."))

	// filepath.Base("") is ".", which has no extension, so BaseName("") is ".".
	c.Equal(".", xfilepath.BaseName(""))
	c.Equal(string(filepath.Separator), xfilepath.BaseName("/"))
	if runtime.GOOS == xos.WindowsOS {
		c.Equal(string(filepath.Separator), xfilepath.BaseName("\\"))
	}
}

func TestTrimExtension(t *testing.T) {
	c := check.New(t)

	c.Equal("file", xfilepath.TrimExtension("file.txt"))
	c.Equal("/path/to/file", xfilepath.TrimExtension("/path/to/file.txt"))
	if runtime.GOOS == xos.WindowsOS {
		c.Equal("C:\\path\\to\\file", xfilepath.TrimExtension("C:\\path\\to\\file.txt"))
	}

	c.Equal("file.backup", xfilepath.TrimExtension("file.backup.txt"))
	c.Equal("archive.tar", xfilepath.TrimExtension("archive.tar.gz"))

	c.Equal("README", xfilepath.TrimExtension("README"))
	c.Equal("/project/Makefile", xfilepath.TrimExtension("/project/Makefile"))

	// A leading dot marks a hidden file, not an extension.
	c.Equal(".txt", xfilepath.TrimExtension(".txt"))
	c.Equal(".hidden", xfilepath.TrimExtension(".hidden"))
	c.Equal(".bashrc", xfilepath.TrimExtension(".bashrc"))
	c.Equal("/etc/.profile", xfilepath.TrimExtension("/etc/.profile"))
	c.Equal("dir/.hidden", xfilepath.TrimExtension("dir/.hidden.txt"))
	c.Equal(".config", xfilepath.TrimExtension(".config.yaml"))

	c.Equal(".", xfilepath.TrimExtension("."))
	c.Equal(".", xfilepath.TrimExtension(".."))

	c.Equal("", xfilepath.TrimExtension(""))

	c.Equal("file", xfilepath.TrimExtension("file."))
	c.Equal("path/file", xfilepath.TrimExtension("path/file."))
}
