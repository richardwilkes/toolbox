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
	"log/slog"
	"net/http"
)

type ctxKey int

const metadataKey ctxKey = 1

// Metadata holds auxiliary information for a request.
type Metadata struct {
	// Logger holds the logger for the request.
	Logger *slog.Logger
	// User holds the user that made the request, if any. Populated by the basicauth middleware.
	User string
	// LogMsg will be used as the message in the final log call for the request if it isn't empty.
	LogMsg string
}

func metadataInContext(ctx context.Context, md *Metadata) context.Context {
	return context.WithValue(ctx, metadataKey, md)
}

// MetadataFromRequest returns the Metadata from the request.
func MetadataFromRequest(req *http.Request) *Metadata {
	if md, ok := req.Context().Value(metadataKey).(*Metadata); ok {
		return md
	}
	return nil
}

// SetMetadataLogMsg sets the LogMsg field of the Metadata in the request context.
// If there is no Metadata in the request, it does nothing.
func SetMetadataLogMsg(req *http.Request, msg string) {
	if md := MetadataFromRequest(req); md != nil {
		md.LogMsg = msg
	}
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
