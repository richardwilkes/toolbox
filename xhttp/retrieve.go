// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xhttp

import (
	"context"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xio"
	"github.com/richardwilkes/toolbox/v2/xos"
)

// HasHTTPOrFileURLPrefix returns true if the provided URL has a http, https, or file scheme. The scheme comparison is
// case-insensitive, since URI schemes are not case-sensitive (RFC 3986).
func HasHTTPOrFileURLPrefix(urlStr string) bool {
	return hasSchemePrefix(urlStr, "http://") ||
		hasSchemePrefix(urlStr, "https://") ||
		hasSchemePrefix(urlStr, "file://")
}

func hasSchemePrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

// RetrieveData loads the bytes from the given file path or URL with scheme file, http, or https. If client is nil and a
// network request is necessary, the http.DefaultClient will be used. No limit is placed on the amount of data that will
// be read; use RetrieveDataWithLimit if the source is untrusted.
func RetrieveData(ctx context.Context, client *http.Client, filePathOrURL string) ([]byte, error) {
	return RetrieveDataWithLimit(ctx, client, filePathOrURL, 0)
}

// RetrieveDataWithLimit behaves like RetrieveData, but returns an error if more than maxBytes would be read. A maxBytes
// of zero or less means no limit.
func RetrieveDataWithLimit(ctx context.Context, client *http.Client, filePathOrURL string, maxBytes int64) ([]byte, error) {
	r, err := StreamData(ctx, client, filePathOrURL)
	if err != nil {
		return nil, err
	}
	defer xio.CloseIgnoringErrors(r)
	var reader io.Reader = r
	enforceLimit := maxBytes > 0 && maxBytes != math.MaxInt64
	if enforceLimit {
		reader = io.LimitReader(r, maxBytes+1) // +1 so reading exactly maxBytes can be distinguished from exceeding it
	}
	var data []byte
	if data, err = io.ReadAll(reader); err != nil {
		return nil, errs.NewWithCause(filePathOrURL, err)
	}
	if enforceLimit && int64(len(data)) > maxBytes {
		return nil, errs.Newf("data from %s exceeds maximum of %d bytes", filePathOrURL, maxBytes)
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
		case ProtocolFile:
			if filePathOrURL, err = fileURLToPath(u); err != nil {
				return nil, err
			}
		case ProtocolHTTP, ProtocolHTTPS:
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
				xio.DiscardAndCloseIgnoringErrors(rsp.Body)
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

// fileURLToPath converts a parsed file URL into a local filesystem path. Only local file URLs are supported, so a
// non-empty host other than "localhost" is rejected, since os.Open cannot reach it.
func fileURLToPath(u *url.URL) (string, error) {
	host := u.Host
	p := u.Path
	if runtime.GOOS == xos.WindowsOS && len(host) == 2 && host[1] == ':' && isASCIILetter(host[0]) {
		// A file URL incorrectly written as file://C:/path (two slashes rather than three) puts the drive letter in the
		// host. Fold it back into the path so it is treated as a local drive path rather than an unreachable remote
		// host.
		p = host + p
		host = ""
	}
	if host != "" && !strings.EqualFold(host, "localhost") {
		return "", errs.Newf("unsupported host in file URL: %s", u.String())
	}
	if runtime.GOOS == xos.WindowsOS {
		// Windows file URLs encode drive paths as /C:/path; strip the leading slash so os.Open sees C:/path.
		if len(p) >= 3 && p[0] == '/' && p[2] == ':' && isASCIILetter(p[1]) {
			p = p[1:]
		}
		p = filepath.FromSlash(p)
	}
	return p, nil
}

func isASCIILetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
