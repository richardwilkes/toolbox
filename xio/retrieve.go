// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// HasHTTPOrFileURLPrefix returns true if the provided URL has a http, https, or file scheme.
func HasHTTPOrFileURLPrefix(urlStr string) bool {
	return strings.HasPrefix(urlStr, "http://") ||
		strings.HasPrefix(urlStr, "https://") ||
		strings.HasPrefix(urlStr, "file://")
}

// RetrieveData loads the bytes from the given file path or URL with scheme file, http, or https. If client is nil and a
// network request is necessary, the http.DefaultClient will be used.
func RetrieveData(ctx context.Context, client *http.Client, filePathOrURL string) ([]byte, error) {
	r, err := StreamData(ctx, client, filePathOrURL)
	if err != nil {
		return nil, err
	}
	defer CloseIgnoringErrors(r)
	var data []byte
	data, err = io.ReadAll(r)
	if err != nil {
		return nil, errs.NewWithCause(filePathOrURL, err)
	}
	return data, nil
}

// StreamData returns an io.ReadCloser that streams the data from the given file path or URL with scheme file, http, or
// https. If client is nil and a network request is necessary, the http.DefaultClient will be used. The caller is
// responsible for closing the returned ReadCloser. Note that for network requests, the entire stream should be read to
// allow reuse of the underlying connection.
func StreamData(ctx context.Context, client *http.Client, filePathOrURL string) (io.ReadCloser, error) {
	if HasHTTPOrFileURLPrefix(filePathOrURL) {
		u, err := url.Parse(filePathOrURL)
		if err != nil {
			return nil, errs.NewWithCause(filePathOrURL, err)
		}
		switch u.Scheme {
		case "file":
			filePathOrURL = u.Path
		case "http", "https":
			var req *http.Request
			req, err = http.NewRequestWithContext(ctx, http.MethodGet, filePathOrURL, http.NoBody)
			if err != nil {
				return nil, errs.NewWithCause("unable to create request", err)
			}
			var rsp *http.Response
			if client == nil {
				client = http.DefaultClient
			}
			if rsp, err = client.Do(req); err != nil {
				return nil, errs.NewWithCause(filePathOrURL, err)
			}
			if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
				return nil, errs.NewWithCause(filePathOrURL, errs.Newf("received status %d (%s)", rsp.StatusCode, rsp.Status))
			}
			return rsp.Body, nil
		default:
			// Shouldn't be possible to reach this
			return nil, errs.Newf("invalid url: %s", filePathOrURL)
		}
	}
	r, err := os.Open(filePathOrURL)
	if err != nil {
		return nil, errs.NewWithCause(filePathOrURL, err)
	}
	return r, nil
}
