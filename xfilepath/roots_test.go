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
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xfilepath"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestUniquePaths_DuplicatePaths(t *testing.T) {
	c := check.New(t)
	dir := t.TempDir()
	resolved := resolve(c, dir)
	target := filepath.Join(dir, "target")
	c.NoError(os.Mkdir(target, 0o755))
	symlink := filepath.Join(dir, "symlink")
	// Creating a symlink requires permissions that aren't available by default on Windows.
	canSymlink := runtime.GOOS != xos.WindowsOS
	if canSymlink {
		c.NoError(os.Symlink(target, symlink))
	}
	t.Chdir(filepath.Dir(dir))
	for _, tc := range []struct {
		paths   []string
		want    []string
		symlink bool
	}{
		{want: []string{}},
		{paths: []string{dir, dir, dir}, want: []string{resolved}},
		{paths: []string{filepath.Base(dir), dir}, want: []string{resolved}},
		{paths: []string{target, symlink}, want: []string{filepath.Join(resolved, "target")}, symlink: true},
	} {
		if tc.symlink && !canSymlink {
			continue
		}
		result, err := xfilepath.UniquePaths(tc.paths...)
		c.NoError(err)
		c.Equal(tc.want, result, "UniquePaths(%q)", tc.paths)
	}
}

func TestUniquePaths_NestedPaths(t *testing.T) {
	c := check.New(t)
	dir := t.TempDir()
	other := t.TempDir()
	subDir := filepath.Join(dir, "subdir")
	subSubDir := filepath.Join(subDir, "subsubdir")
	c.NoError(os.MkdirAll(subSubDir, 0o755))
	// Relative to dir this child is "..foo", which must not be mistaken for an ascending path.
	dotDotChild := filepath.Join(dir, "..foo")
	c.NoError(os.Mkdir(dotDotChild, 0o755))
	resolvedDir := resolve(c, dir)
	resolvedOther := resolve(c, other)
	for _, tc := range []struct {
		paths []string
		want  []string
	}{
		{paths: []string{dir, subDir, subSubDir}, want: []string{resolvedDir}},
		{paths: []string{dir, dotDotChild}, want: []string{resolvedDir}},
		{paths: []string{dotDotChild, dir}, want: []string{resolvedDir}},
		{paths: []string{dir, other}, want: []string{resolvedDir, resolvedOther}},
		{paths: []string{dir, other, subDir}, want: []string{resolvedDir, resolvedOther}},
	} {
		result, err := xfilepath.UniquePaths(tc.paths...)
		c.NoError(err)
		slices.Sort(result)
		slices.Sort(tc.want)
		c.Equal(tc.want, result, "UniquePaths(%q)", tc.paths)
	}
}

func TestUniquePaths_NonexistentPath(t *testing.T) {
	c := check.New(t)
	nonexistent := "/nonexistent/path/that/does/not/exist"
	_, err := xfilepath.UniquePaths(nonexistent)
	c.HasError(err)
}

func resolve(c check.Checker, path string) string {
	c.Helper()
	abs, err := filepath.Abs(path)
	c.NoError(err)
	resolved, err := filepath.EvalSymlinks(abs)
	c.NoError(err)
	return resolved
}
