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
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xio"
)

func TestHasHttpOrFileURLPrefix(t *testing.T) {
	c := check.New(t)
	c.True(xio.HasHTTPOrFileURLPrefix("http://example.com"))
	c.True(xio.HasHTTPOrFileURLPrefix("https://example.com"))
	c.True(xio.HasHTTPOrFileURLPrefix("file:///tmp/foo"))
	c.False(xio.HasHTTPOrFileURLPrefix("ftp://example.com"))
	c.False(xio.HasHTTPOrFileURLPrefix("/tmp/foo"))
	c.False(xio.HasHTTPOrFileURLPrefix("c:/tmp/foo"))
	c.False(xio.HasHTTPOrFileURLPrefix("c:\\tmp\\foo"))
}

func TestRetrieveData_File(t *testing.T) {
	c := check.New(t)
	file := filepath.Join(t.TempDir(), "retrieve_test_1.txt")
	content := []byte("hello file")
	c.NoError(os.WriteFile(file, content, 0o600))
	data, err := xio.RetrieveData(context.Background(), nil, file)
	c.NoError(err)
	c.Equal(content, data)
}

func TestRetrieveData_FileURL(t *testing.T) {
	c := check.New(t)
	file := filepath.Join(t.TempDir(), "retrieve_test_2.txt")
	content := []byte("hello fileurl")
	c.NoError(os.WriteFile(file, content, 0o600))
	fileURL := "file://" + file
	if runtime.GOOS == "windows" {
		fileURL = "file:///" + filepath.ToSlash(file)
	}
	data, err := xio.RetrieveData(context.Background(), nil, fileURL)
	c.NoError(err)
	c.Equal(content, data)
}

func TestRetrieveData_HTTP(t *testing.T) {
	c := check.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("hello http")) //nolint:errcheck // Can't fail for test, so no need to check
	}))
	defer server.Close()
	data, err := xio.RetrieveData(context.Background(), nil, server.URL)
	c.NoError(err)
	c.Equal([]byte("hello http"), data)
}

func TestRetrieveData_HTTPS(t *testing.T) {
	c := check.New(t)
	// Use httptest.Server (http, not https), but test the code path
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("hello https")) //nolint:errcheck // Can't fail for test, so no need to check
	}))
	defer server.Close()
	client := server.Client()
	data, err := xio.RetrieveData(context.Background(), client, server.URL)
	c.NoError(err)
	c.Equal([]byte("hello https"), data)
}

func TestRetrieveData_HTTPErrorStatus(t *testing.T) {
	c := check.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("not found")) //nolint:errcheck // Can't fail for test, so no need to check
	}))
	defer server.Close()
	_, err := xio.RetrieveData(context.Background(), nil, server.URL)
	c.HasError(err)
}

func TestRetrieveData_FileNotFound(t *testing.T) {
	c := check.New(t)
	_, err := xio.RetrieveData(context.Background(), nil, "nonexistent_file_123456789.txt")
	c.HasError(err)
}

func TestStreamData_InvalidURL(t *testing.T) {
	c := check.New(t)
	_, err := xio.StreamData(context.Background(), nil, "http://%41:8080/")
	c.HasError(err)
}

func TestStreamData_UnsupportedScheme(t *testing.T) {
	c := check.New(t)
	_, err := xio.StreamData(context.Background(), nil, "ftp://example.com")
	c.HasError(err)
}

func TestStreamData_File(t *testing.T) {
	c := check.New(t)
	file := filepath.Join(t.TempDir(), "retrieve_test_3.txt")
	content := []byte("stream file")
	c.NoError(os.WriteFile(file, content, 0o600))
	r, err := xio.StreamData(context.Background(), nil, file)
	c.NoError(err)
	defer xio.CloseIgnoringErrors(r)
	data, err := io.ReadAll(r)
	c.NoError(err)
	c.Equal(content, data)
}

func TestStreamData_HTTP(t *testing.T) {
	c := check.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("stream http")) //nolint:errcheck // Can't fail for test, so no need to check
	}))
	defer server.Close()
	r, err := xio.StreamData(context.Background(), nil, server.URL)
	c.NoError(err)
	defer xio.CloseIgnoringErrors(r)
	data, err := io.ReadAll(r)
	c.NoError(err)
	c.Equal([]byte("stream http"), data)
}

func TestStreamData_HTTPErrorStatus(t *testing.T) {
	c := check.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("fail")) //nolint:errcheck // Can't fail for test, so no need to check
	}))
	defer server.Close()
	_, err := xio.StreamData(context.Background(), nil, server.URL)
	c.HasError(err)
}

func TestStreamData_FileNotFound(t *testing.T) {
	c := check.New(t)
	_, err := xio.StreamData(context.Background(), nil, "nonexistent_file_987654321.txt")
	c.HasError(err)
}
