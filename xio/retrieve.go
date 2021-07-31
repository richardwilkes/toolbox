// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/richardwilkes/toolbox/errs"
)

// RetrieveDataFromURL loads the bytes from the given URL. Only file, http,
// and https URLs are currently supported.
func RetrieveDataFromURL(urlStr string) ([]byte, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errs.NewWithCause(urlStr, err)
	}
	var data []byte
	switch u.Scheme {
	case "file":
		if data, err = ioutil.ReadFile(u.Path); err != nil {
			return nil, errs.NewWithCause(urlStr, err)
		}
		return data, nil
	case "http", "https":
		var rsp *http.Response
		if rsp, err = http.Get(urlStr); err != nil { //nolint:gosec
			return nil, errs.NewWithCause(urlStr, err)
		}
		defer CloseIgnoringErrors(rsp.Body)
		if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
			return nil, errs.NewWithCause(urlStr, errs.Newf("received status %d (%s)", rsp.StatusCode, rsp.Status))
		}
		data, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, errs.NewWithCause(urlStr, err)
		}
		return data, nil
	default:
		return nil, errs.Newf("invalid url: %s", urlStr)
	}
}
