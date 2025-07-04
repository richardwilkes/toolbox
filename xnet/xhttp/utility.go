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

// LoggerForRequest returns a logger for use with the request.
func LoggerForRequest(r *http.Request) *slog.Logger {
	var logger *slog.Logger
	if md := MetadataFromRequest(r); md != nil {
		logger = md.Logger
	} else {
		logger = slog.Default()
	}
	return logger
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

// ClientIP implements a best effort algorithm to return the real client IP address from the
// request. It parses True-Client-Ip, X-Forwarded-For, and X-Real-IP in order to work properly with
// reverse proxies such as akamai, nginx, or haproxy.
func ClientIP(req *http.Request) string {
	ip := strings.TrimSpace(req.Header.Get("True-Client-Ip"))
	if ip != "" && net.ParseIP(ip) != nil {
		return ip
	}
	ip = req.Header.Get("X-Forwarded-For")
	if index := strings.IndexByte(ip, ','); index >= 0 {
		ip = ip[0:index]
	}
	if ip = strings.TrimSpace(ip); ip != "" && net.ParseIP(ip) != nil {
		return ip
	}
	if ip = strings.TrimSpace(req.Header.Get("X-Real-Ip")); ip != "" && net.ParseIP(ip) != nil {
		return ip
	}
	return req.RemoteAddr
}
