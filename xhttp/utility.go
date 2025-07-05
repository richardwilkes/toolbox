// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xio"
)

// ErrorStatus sends an HTTP response header with 'statusCode' and follows it with the standard text for that code as
// the body.
func ErrorStatus(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// RequestError logs an error to the logger associated with the request.
func RequestError(req *http.Request, err error) {
	errs.LogTo(LoggerForRequest(req), err)
}

// RequestWarning logs an error to the logger associated with the request as a warning.
func RequestWarning(req *http.Request, err error) {
	errs.LogWithLevel(context.Background(), slog.LevelWarn, LoggerForRequest(req), err)
}

// DisableCaching disables caching for the given response writer. To be effective, should be called before any data is
// written.
func DisableCaching(w http.ResponseWriter) {
	header := w.Header()
	header.Set("Cache-Control", "no-store")
	header.Set("Pragma", "no-cache")
}

// JSONResponse writes a JSON response with a status code.
func JSONResponse(w http.ResponseWriter, req *http.Request, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		RequestError(req, err)
	}
}

// ExtractJSONBody extracts the body and tries to marshal it from JSON into the data parameter.
func ExtractJSONBody(req *http.Request, data any) error {
	defer xio.DiscardAndCloseIgnoringErrors(req.Body)
	decoder := json.NewDecoder(req.Body)
	decoder.UseNumber()
	return decoder.Decode(data)
}

// ClientIP looks at the X-Forwarded-For, Forwarded, and RemoteAddr headers (in that order) to determine the client's
// actual IP address.
func ClientIP(req *http.Request) net.IP {
	if xForwardedFor := req.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For can contain multiple values, we take the first one.
		// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For for more information.
		if ip := net.ParseIP(strings.TrimSpace(strings.SplitN(xForwardedFor, ",", 2)[0])); ip != nil {
			return ip
		}
	}
	if forwarded := req.Header.Get("Forwarded"); forwarded != "" {
		// Forwarded can contain multiple values, we take the first one.
		// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Forwarded for more information.
		for _, forwarded = range strings.Split(strings.SplitN(forwarded, ",", 2)[0], ";") {
			if strings.HasPrefix(forwarded, "for=") {
				forwarded = strings.TrimPrefix(forwarded, "for=")
				forwarded = strings.TrimPrefix(forwarded, `"`)
				forwarded = strings.TrimSuffix(forwarded, `"`)
				if ip := net.ParseIP(forwarded); ip != nil {
					return ip
				}
			}
		}
	}
	if req.RemoteAddr != "" {
		// RemoteAddr is in the form "IP:port", we take the IP part.
		if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			if ip := net.ParseIP(host); ip != nil {
				return ip
			}
		}
	}
	return nil
}
