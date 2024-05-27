/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package xio

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

// RetrieveData loads the bytes from the given file path or URL of type file, http, or https.
func RetrieveData(filePathOrURL string) ([]byte, error) {
	return RetrieveDataWithContext(context.Background(), filePathOrURL)
}

// RetrieveDataWithContext loads the bytes from the given file path or URL of type file, http, or https.
func RetrieveDataWithContext(ctx context.Context, filePathOrURL string) ([]byte, error) {
	if strings.HasPrefix(filePathOrURL, "http://") ||
		strings.HasPrefix(filePathOrURL, "https://") ||
		strings.HasPrefix(filePathOrURL, "file://") {
		return RetrieveDataFromURLWithContext(ctx, filePathOrURL)
	}
	data, err := os.ReadFile(filePathOrURL)
	if err != nil {
		return nil, errs.NewWithCause(filePathOrURL, err)
	}
	return data, nil
}

// RetrieveDataFromURL loads the bytes from the given URL of type file, http, or https.
func RetrieveDataFromURL(urlStr string) ([]byte, error) {
	return RetrieveDataFromURLWithContext(context.Background(), urlStr)
}

// RetrieveDataFromURLWithContext loads the bytes from the given URL of type file, http, or https.
func RetrieveDataFromURLWithContext(ctx context.Context, urlStr string) ([]byte, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errs.NewWithCause(urlStr, err)
	}
	var data []byte
	switch u.Scheme {
	case "file":
		if data, err = os.ReadFile(u.Path); err != nil {
			return nil, errs.NewWithCause(urlStr, err)
		}
		return data, nil
	case "http", "https":
		var req *http.Request
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
		if err != nil {
			return nil, errs.NewWithCause("unable to create request", err)
		}
		var rsp *http.Response
		if rsp, err = http.DefaultClient.Do(req); err != nil {
			return nil, errs.NewWithCause(urlStr, err)
		}
		defer DiscardAndCloseIgnoringErrors(rsp.Body)
		if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
			return nil, errs.NewWithCause(urlStr, errs.Newf("received status %d (%s)", rsp.StatusCode, rsp.Status))
		}
		data, err = io.ReadAll(rsp.Body)
		if err != nil {
			return nil, errs.NewWithCause(urlStr, err)
		}
		return data, nil
	default:
		return nil, errs.Newf("invalid url: %s", urlStr)
	}
}
