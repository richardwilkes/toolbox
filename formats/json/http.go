// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package json

import (
	"bytes"
	"context"
	"net/http"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// GetRequest calls http.Get with the URL and returns the response body as a new Data object.
func GetRequest(url string) (statusCode int, body *Data, err error) {
	return GetRequestWithContext(context.Background(), url)
}

// GetRequestWithContext calls http.Get with the URL and returns the response body as a new Data object.
func GetRequestWithContext(ctx context.Context, url string) (statusCode int, body *Data, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, errs.NewWithCause("unable to create request", err)
	}
	return makeRequest(req)
}

// PostRequest calls http.Post with the URL and the contents of this Data object and returns the response body as a new
// Data object.
func (j *Data) PostRequest(url string) (statusCode int, body *Data, err error) {
	return j.PostRequestWithContext(context.Background(), url)
}

// PostRequestWithContext calls http.Post with the URL and the contents of this Data object and returns the response
// body as a new Data object.
func (j *Data) PostRequestWithContext(ctx context.Context, url string) (statusCode int, body *Data, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(j.Bytes()))
	if err != nil {
		return 0, nil, errs.NewWithCause("unable to create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return makeRequest(req)
}

func makeRequest(req *http.Request) (statusCode int, body *Data, err error) {
	var rsp *http.Response
	if rsp, err = http.DefaultClient.Do(req); err != nil {
		return 0, &Data{}, errs.NewWithCause("request failed", err)
	}
	defer xio.DiscardAndCloseIgnoringErrors(rsp.Body)
	statusCode = rsp.StatusCode
	body = MustParseStream(rsp.Body)
	return
}
