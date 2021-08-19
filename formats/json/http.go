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
	"net/http"
)

// GetRequest calls http.Get with the URL and returns the response body as a
// new Data object.
func GetRequest(url string) (statusCode int, body *Data, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err == nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		statusCode = resp.StatusCode
		body = MustParseStream(resp.Body)
	} else {
		body = &Data{}
	}
	return
}

// PostRequest calls http.Post with the URL and the contents of this Data
// object and returns the response body as a new Data object.
func (j *Data) PostRequest(url string) (statusCode int, body *Data, err error) {
	var resp *http.Response
	if resp, err = http.Post(url, "application/json", bytes.NewBuffer(j.Bytes())); err == nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		statusCode = resp.StatusCode
		body = MustParseStream(resp.Body)
	} else {
		body = &Data{}
	}
	return
}
