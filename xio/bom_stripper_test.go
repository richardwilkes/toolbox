// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xio"
)

func TestNewBOMStripper_NoBOM(t *testing.T) {
	c := check.New(t)
	input := "hello world"
	r := strings.NewReader(input)
	br, err := xio.NewBOMStripper(r)
	c.NoError(err)
	out, err := io.ReadAll(br)
	c.NoError(err)
	c.Equal([]byte(input), out)
}

func TestNewBOMStripper_WithBOM(t *testing.T) {
	c := check.New(t)
	// UTF-8 BOM is 0xEF,0xBB,0xBF
	input := string([]byte{0xEF, 0xBB, 0xBF}) + "hello world"
	r := strings.NewReader(input)
	br, err := xio.NewBOMStripper(r)
	c.NoError(err)
	out, err := io.ReadAll(br)
	c.NoError(err)
	c.Equal([]byte("hello world"), out)
}

func TestNewBOMStripper_AlreadyBufioReader(t *testing.T) {
	c := check.New(t)
	input := "hello world"
	buf := bytes.NewBufferString(input)
	br := bufio.NewReader(buf)
	br2, err := xio.NewBOMStripper(br)
	c.NoError(err)
	c.True(br == br2) // should reuse the same bufio.Reader
	out, err := io.ReadAll(br2)
	c.NoError(err)
	c.Equal([]byte(input), out)
}

func TestNewBOMStripper_ErrorOnRead(t *testing.T) {
	c := check.New(t)
	r := &errReader{err: errors.New("fail")}
	_, err := xio.NewBOMStripper(r)
	c.HasError(err)
}

type errReader struct{ err error }

func (e *errReader) Read(p []byte) (int, error) {
	return 0, e.err
}
