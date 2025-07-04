// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/richardwilkes/toolbox/v2/xio"
	"github.com/richardwilkes/toolbox/v2/xnet/xhttp"
)

var _ http.ResponseWriter = &gzipResponseWriter{}

type gzipResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

// Wrap the given handler and provide automatic gzip compression when requests advertise that they accept the gzip
// encoding.
func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gw := gzip.NewWriter(w)
			defer func() { xio.CloseLoggingAnyErrorTo(xhttp.LoggerForRequest(req), gw) }()
			w = &gzipResponseWriter{
				w:  w,
				gw: gw,
			}
		}
		next.ServeHTTP(w, req)
	})
}

// Header implements http.ResponseWriter.
func (w *gzipResponseWriter) Header() http.Header {
	return w.w.Header()
}

// Write implements http.ResponseWriter.
func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.gw.Write(data)
}

// WriteHeader implements http.ResponseWriter.
func (w *gzipResponseWriter) WriteHeader(status int) {
	w.w.WriteHeader(status)
}
