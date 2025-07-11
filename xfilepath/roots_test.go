// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xfilepath_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xfilepath"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestUniquePaths_Empty(t *testing.T) {
	c := check.New(t)
	result, err := xfilepath.UniquePaths()
	c.NoError(err)
	c.Equal([]string{}, result)
}

func TestUniquePaths_SinglePath(t *testing.T) {
	c := check.New(t)
	tempDir := t.TempDir()
	result, err := xfilepath.UniquePaths(tempDir)
	c.NoError(err)
	c.Equal(1, len(result))
	// The result should be the absolute path
	abs, err := filepath.Abs(tempDir)
	c.NoError(err)
	noSymLinks, err := filepath.EvalSymlinks(abs)
	c.NoError(err)
	c.Equal(noSymLinks, result[0])
}

func TestUniquePaths_DuplicatePaths(t *testing.T) {
	c := check.New(t)
	tempDir := t.TempDir()
	result, err := xfilepath.UniquePaths(tempDir, tempDir, tempDir)
	c.NoError(err)
	c.Equal(1, len(result))
	abs, err := filepath.Abs(tempDir)
	c.NoError(err)
	noSymLinks, err := filepath.EvalSymlinks(abs)
	c.NoError(err)
	c.Equal(noSymLinks, result[0])
}

func TestUniquePaths_RelativeAndAbsolute(t *testing.T) {
	c := check.New(t)
	tempDir := t.TempDir()
	abs, err := filepath.Abs(tempDir)
	c.NoError(err)

	// Change to a different directory to test relative path
	originalWd, err := os.Getwd()
	c.NoError(err)
	defer func() {
		c.NoError(os.Chdir(originalWd))
	}()

	parentDir := filepath.Dir(tempDir)
	c.NoError(os.Chdir(parentDir))
	relPath := filepath.Base(tempDir)

	result, err := xfilepath.UniquePaths(relPath, abs)
	c.NoError(err)
	c.Equal(1, len(result))
	noSymLinks, err := filepath.EvalSymlinks(abs)
	c.NoError(err)
	c.Equal(noSymLinks, result[0])
}

func TestUniquePaths_NestedPaths(t *testing.T) {
	c := check.New(t)
	tempDir := t.TempDir()

	// Create nested directories
	subDir := filepath.Join(tempDir, "subdir")
	c.NoError(os.MkdirAll(subDir, 0o755))

	subSubDir := filepath.Join(subDir, "subsubdir")
	c.NoError(os.MkdirAll(subSubDir, 0o755))

	result, err := xfilepath.UniquePaths(tempDir, subDir, subSubDir)
	c.NoError(err)
	c.Equal(1, len(result))

	// Only the root directory should remain
	abs, err := filepath.Abs(tempDir)
	c.NoError(err)
	noSymLinks, err := filepath.EvalSymlinks(abs)
	c.NoError(err)
	c.Equal(noSymLinks, result[0])
}

func TestUniquePaths_SeparateDirectories(t *testing.T) {
	c := check.New(t)
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()
	c.NotEqual(tempDir1, tempDir2)

	result, err := xfilepath.UniquePaths(tempDir1, tempDir2)
	c.NoError(err)
	c.Equal(2, len(result))

	abs1, err := filepath.Abs(tempDir1)
	c.NoError(err)
	noSymLinks1, err := filepath.EvalSymlinks(abs1)
	c.NoError(err)
	abs2, err := filepath.Abs(tempDir2)
	c.NoError(err)
	noSymLinks2, err := filepath.EvalSymlinks(abs2)
	c.NoError(err)

	// Both directories should be in the result
	found1, found2 := false, false
	for _, path := range result {
		if path == noSymLinks1 {
			found1 = true
		}
		if path == noSymLinks2 {
			found2 = true
		}
	}
	c.True(found1)
	c.True(found2)
}

func TestUniquePaths_MixedNestedAndSeparate(t *testing.T) {
	c := check.New(t)
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()
	c.NotEqual(tempDir1, tempDir2)

	// Create a subdirectory in tempDir1
	subDir := filepath.Join(tempDir1, "subdir")
	c.NoError(os.MkdirAll(subDir, 0o755))

	result, err := xfilepath.UniquePaths(tempDir1, tempDir2, subDir)
	c.NoError(err)
	c.Equal(2, len(result))

	abs1, err := filepath.Abs(tempDir1)
	c.NoError(err)
	noSymLinks1, err := filepath.EvalSymlinks(abs1)
	c.NoError(err)
	abs2, err := filepath.Abs(tempDir2)
	c.NoError(err)
	noSymLinks2, err := filepath.EvalSymlinks(abs2)
	c.NoError(err)

	// Should contain tempDir1 and tempDir2, but not subDir (subset of tempDir1)
	found1, found2 := false, false
	for _, path := range result {
		if path == noSymLinks1 {
			found1 = true
		}
		if path == noSymLinks2 {
			found2 = true
		}
		// subDir should not be in the result
		c.NotEqual(filepath.Join(noSymLinks1, "subdir"), path)
	}
	c.True(found1)
	c.True(found2)
}

func TestUniquePaths_WithSymlinks(t *testing.T) {
	if runtime.GOOS == xos.WindowsOS {
		t.Skip("This test requires permissions that aren't available by default on Windows")
	}
	c := check.New(t)
	tempDir := t.TempDir()

	// Create a target directory
	targetDir := filepath.Join(tempDir, "target")
	c.NoError(os.MkdirAll(targetDir, 0o755))

	// Create a symlink
	symlinkPath := filepath.Join(tempDir, "symlink")
	c.NoError(os.Symlink(targetDir, symlinkPath))

	result, err := xfilepath.UniquePaths(targetDir, symlinkPath)
	c.NoError(err)
	c.Equal(1, len(result))

	// Both should resolve to the same target
	absTarget, err := filepath.Abs(targetDir)
	c.NoError(err)
	noSymLinks, err := filepath.EvalSymlinks(absTarget)
	c.NoError(err)
	c.Equal(noSymLinks, result[0])
}

func TestUniquePaths_NonexistentPath(t *testing.T) {
	c := check.New(t)
	nonexistent := "/nonexistent/path/that/does/not/exist"
	_, err := xfilepath.UniquePaths(nonexistent)
	c.HasError(err)
}

func TestUniquePaths_CurrentDirectory(t *testing.T) {
	c := check.New(t)
	result, err := xfilepath.UniquePaths(".")
	c.NoError(err)
	c.Equal(1, len(result))

	// Should resolve to absolute path of current directory
	cwd, err := os.Getwd()
	c.NoError(err)
	// Case-insensitive comparison for Windows
	c.True(strings.EqualFold(cwd, result[0]), "Expected current directory (%q) to match result, got: %q", cwd, result[0])
}

func TestUniquePaths_ParentDirectory(t *testing.T) {
	c := check.New(t)
	tempDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	c.NoError(err)
	defer func() {
		c.NoError(os.Chdir(originalWd))
	}()
	c.NoError(os.Chdir(tempDir))

	result, err := xfilepath.UniquePaths("..")
	c.NoError(err)
	c.Equal(1, len(result))

	// Should resolve to absolute path of parent directory
	parentDir := filepath.Dir(tempDir)
	noSymLinks, err := filepath.EvalSymlinks(parentDir)
	c.NoError(err)
	c.Equal(noSymLinks, result[0])
}

func TestUniquePaths_ComplexNesting(t *testing.T) {
	c := check.New(t)
	tempDir := t.TempDir()

	// Create a complex directory structure
	// root/
	//   ├── a/
	//   │   └── b/
	//   ├── c/
	//   └── d/
	//       └── e/
	//           └── f/

	pathA := filepath.Join(tempDir, "a")
	pathB := filepath.Join(pathA, "b")
	pathC := filepath.Join(tempDir, "c")
	pathD := filepath.Join(tempDir, "d")
	pathE := filepath.Join(pathD, "e")
	pathF := filepath.Join(pathE, "f")

	c.NoError(os.MkdirAll(pathB, 0o755))
	c.NoError(os.MkdirAll(pathC, 0o755))
	c.NoError(os.MkdirAll(pathF, 0o755))

	result, err := xfilepath.UniquePaths(tempDir, pathA, pathB, pathC, pathD, pathE, pathF)
	c.NoError(err)
	c.Equal(1, len(result))

	// Only the root should remain
	abs, err := filepath.Abs(tempDir)
	c.NoError(err)
	noSymLinks, err := filepath.EvalSymlinks(abs)
	c.NoError(err)
	c.Equal(noSymLinks, result[0])
}
