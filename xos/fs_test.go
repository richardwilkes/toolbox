// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestIsDir(t *testing.T) {
	c := check.New(t)

	// path is a directory
	tmpDir := t.TempDir()
	c.True(xos.IsDir(tmpDir))

	// file does not exist
	file := filepath.Join(tmpDir, "isdir-regular.txt")
	c.False(xos.IsDir(file))

	// file exists and is regular
	c.NoError(os.WriteFile(file, []byte("test content"), 0o644))
	c.False(xos.IsDir(file))
}

func TestFileExists(t *testing.T) {
	c := check.New(t)

	// path is a directory
	tmpDir := t.TempDir()
	c.False(xos.FileExists(tmpDir))

	// file does not exist
	file := filepath.Join(tmpDir, "fileexists-regular.txt")
	c.False(xos.FileExists(file))

	// file exists and is regular
	c.NoError(os.WriteFile(file, []byte("test content"), 0o644))
	c.True(xos.FileExists(file))
}

func TestFileIsReadable(t *testing.T) {
	c := check.New(t)

	// path is a directory
	tmpDir := t.TempDir()
	c.False(xos.FileIsReadable(tmpDir))

	// file does not exist
	file := filepath.Join(tmpDir, "readable-regular.txt")
	c.False(xos.FileIsReadable(file))

	// file exists and is readable
	c.NoError(os.WriteFile(file, []byte("test content"), 0o644))
	c.True(xos.FileIsReadable(file))

	// file exists but not readable
	if runtime.GOOS != xos.WindowsOS { // Windows seems to ignore the write-only permission and give it read access too
		noReadFile := filepath.Join(tmpDir, "not-readable.txt")
		c.NoError(os.WriteFile(noReadFile, []byte("test content"), 0o200))
		c.False(xos.FileIsReadable(noReadFile))
	}
}

func TestMoveFile(t *testing.T) {
	c := check.New(t)

	// Test moving a non-existent source file
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")
	c.HasError(xos.MoveFile(srcFile, dstFile))

	// Test moving a directory as source
	c.HasError(xos.MoveFile(tmpDir, dstFile))

	// Test moving to a directory as destination
	c.NoError(os.WriteFile(srcFile, []byte("test content"), 0o644))
	dstDir := filepath.Join(tmpDir, "destdir")
	c.NoError(os.MkdirAll(dstDir, 0o755))
	c.HasError(xos.MoveFile(srcFile, dstDir))

	// Test successful move
	dstFile = filepath.Join(tmpDir, "dest.txt")
	c.NoError(xos.MoveFile(srcFile, dstFile))
	c.False(xos.FileExists(srcFile))
	c.True(xos.FileExists(dstFile))
	content, err := os.ReadFile(dstFile)
	c.NoError(err)
	c.Equal("test content", string(content))

	// Test moving to same file
	srcFile = filepath.Join(tmpDir, "same.txt")
	c.NoError(os.WriteFile(srcFile, []byte("same content"), 0o644))
	c.NoError(xos.MoveFile(srcFile, srcFile))
	c.True(xos.FileExists(srcFile))
}

func TestCopySymlinkIntoMissingParent(t *testing.T) {
	if runtime.GOOS == xos.WindowsOS { // Windows doesn't support symlinks without special permissions enabled first
		t.Skip("symlinks not supported without special permissions on Windows")
	}
	c := check.New(t)

	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "target.txt")
	c.NoError(os.WriteFile(srcFile, []byte("test content"), 0o644))
	srcLink := filepath.Join(tmpDir, "link.txt")
	c.NoError(os.Symlink(srcFile, srcLink))

	// Destination's parent directory does not exist yet; linkCopy must create it.
	dstLink := filepath.Join(tmpDir, "sub", "nested", "copylink.txt")
	c.NoError(xos.Copy(srcLink, dstLink))

	linkTarget, err := os.Readlink(dstLink)
	c.NoError(err)
	c.Equal(srcFile, linkTarget)
}

func TestCopy(t *testing.T) {
	c := check.New(t)

	// Test copying regular file
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")
	dstFile := filepath.Join(tmpDir, "dst.txt")
	content := []byte("test content")
	c.NoError(os.WriteFile(srcFile, content, 0o644))
	c.NoError(xos.Copy(srcFile, dstFile))
	copiedContent, err := os.ReadFile(dstFile)
	c.NoError(err)
	c.Equal(string(content), string(copiedContent))

	// Test copying directory
	srcDir := filepath.Join(tmpDir, "srcdir")
	dstDir := filepath.Join(tmpDir, "dstdir")
	c.NoError(os.MkdirAll(srcDir, 0o755))
	c.NoError(os.WriteFile(filepath.Join(srcDir, "file.txt"), content, 0o644))
	c.NoError(xos.Copy(srcDir, dstDir))
	copiedContent, err = os.ReadFile(filepath.Join(dstDir, "file.txt"))
	c.NoError(err)
	c.Equal(string(content), string(copiedContent))

	// Test copying symlink
	if runtime.GOOS != xos.WindowsOS { // Windows doesn't support symlinks without special permissions enabled first
		srcLink := filepath.Join(tmpDir, "link.txt")
		dstLink := filepath.Join(tmpDir, "copylink.txt")
		c.NoError(os.Symlink(srcFile, srcLink))
		c.NoError(xos.Copy(srcLink, dstLink))
		var linkTarget, origTarget string
		linkTarget, err = os.Readlink(dstLink)
		c.NoError(err)
		origTarget, err = os.Readlink(srcLink)
		c.NoError(err)
		c.Equal(origTarget, linkTarget)
	}

	// Test copying non-existent file
	c.HasError(xos.Copy(filepath.Join(tmpDir, "nonexistent"), dstFile))

	// Test copying file that is not writable
	srcFile = filepath.Join(tmpDir, "src-no-write.txt")
	dstFile = filepath.Join(tmpDir, "dst-no-write.txt")
	c.NoError(os.WriteFile(srcFile, content, 0o444))
	c.NoError(xos.Copy(srcFile, dstFile))
	copiedContent, err = os.ReadFile(dstFile)
	c.NoError(err)
	c.Equal(string(content), string(copiedContent))
	var fi os.FileInfo
	fi, err = os.Stat(dstFile)
	c.NoError(err)
	c.Equal(fs.FileMode(0o444), fi.Mode())
}

// TestCopyWithMaskClearingOwnerWrite verifies that a mask which strips the owner's write bit does not abort a directory
// copy partway. Before the fix, destination directories were created without owner write, so creating files inside them
// failed with a permission error.
func TestCopyWithMaskClearingOwnerWrite(t *testing.T) {
	if runtime.GOOS == xos.WindowsOS {
		t.Skip("POSIX permission semantics not applicable on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks, so the masked-copy scenario cannot be exercised")
	}
	c := check.New(t)

	// Build a small source tree and normalize its modes so the assertions are deterministic regardless of umask.
	srcDir := filepath.Join(t.TempDir(), "src")
	c.NoError(os.MkdirAll(filepath.Join(srcDir, "sub"), 0o755))
	content := []byte("payload")
	c.NoError(os.WriteFile(filepath.Join(srcDir, "file.txt"), content, 0o644))
	c.NoError(os.WriteFile(filepath.Join(srcDir, "sub", "nested.txt"), content, 0o644))
	c.NoError(os.Chmod(filepath.Join(srcDir, "sub"), 0o755))
	c.NoError(os.Chmod(srcDir, 0o755))

	dstDir := filepath.Join(t.TempDir(), "dst")
	// Restore owner write on the copied directories so t.TempDir cleanup can remove the tree. Registered after
	// t.TempDir, so it runs before the temp dir removal (cleanups run in LIFO order).
	t.Cleanup(func() {
		c.NoError(os.Chmod(dstDir, 0o700))
		c.NoError(os.Chmod(filepath.Join(dstDir, "sub"), 0o700))
	})

	// Mask 0o577 clears the owner write bit (0o755 & 0o577 == 0o555 for dirs, 0o644 & 0o577 == 0o444 for files).
	c.NoError(xos.CopyWithMask(srcDir, dstDir, 0o577))

	// All content was copied.
	got, err := os.ReadFile(filepath.Join(dstDir, "file.txt"))
	c.NoError(err)
	c.Equal(string(content), string(got))
	got, err = os.ReadFile(filepath.Join(dstDir, "sub", "nested.txt"))
	c.NoError(err)
	c.Equal(string(content), string(got))

	// Directories carry exactly the masked mode.
	fi, err := os.Stat(dstDir)
	c.NoError(err)
	c.Equal(fs.FileMode(0o555), fi.Mode().Perm())
	fi, err = os.Stat(filepath.Join(dstDir, "sub"))
	c.NoError(err)
	c.Equal(fs.FileMode(0o555), fi.Mode().Perm())

	// Files carry exactly the masked mode.
	fi, err = os.Stat(filepath.Join(dstDir, "file.txt"))
	c.NoError(err)
	c.Equal(fs.FileMode(0o444), fi.Mode().Perm())
}

// TestCopySymlinkOverExisting verifies that copying a symlink onto an existing destination overwrites it, matching the
// overwrite behavior of the regular-file path and making tree copies containing symlinks idempotent. Before the fix,
// os.Symlink failed with EEXIST when the destination already existed.
func TestCopySymlinkOverExisting(t *testing.T) {
	if runtime.GOOS == xos.WindowsOS {
		t.Skip("symlinks not supported without special permissions on Windows")
	}
	c := check.New(t)

	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target.txt")
	c.NoError(os.WriteFile(target, []byte("content"), 0o644))
	srcLink := filepath.Join(tmpDir, "src-link")
	c.NoError(os.Symlink(target, srcLink))

	// Destination is an existing symlink pointing elsewhere (the re-copy / idempotency case).
	otherTarget := filepath.Join(tmpDir, "other.txt")
	c.NoError(os.WriteFile(otherTarget, []byte("other"), 0o644))
	dstLink := filepath.Join(tmpDir, "dst-link")
	c.NoError(os.Symlink(otherTarget, dstLink))
	c.NoError(xos.Copy(srcLink, dstLink))
	got, err := os.Readlink(dstLink)
	c.NoError(err)
	c.Equal(target, got)

	// Destination is an existing regular file.
	dstFile := filepath.Join(tmpDir, "dst-file")
	c.NoError(os.WriteFile(dstFile, []byte("existing"), 0o644))
	c.NoError(xos.Copy(srcLink, dstFile))
	got, err = os.Readlink(dstFile)
	c.NoError(err)
	c.Equal(target, got)
}
