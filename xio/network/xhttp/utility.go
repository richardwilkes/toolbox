// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xhttp

import "net/http"

// ErrorStatus sends an HTTP response header with 'statusCode' and follows it with the standard text for that code as
// the body.
func ErrorStatus(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// DisableCaching disables caching for the given response writer. To be effective, should be called before any data is
// written.
func DisableCaching(w http.ResponseWriter) {
	header := w.Header()
	header.Set("Cache-Control", "no-store")
	header.Set("Pragma", "no-cache")
}
