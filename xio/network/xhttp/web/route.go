// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package web

import (
	"net/http"
	"strings"
)

type routeCtxKey int

var routeKey routeCtxKey = 1

type route struct {
	path string
	last string
}

func (r *route) shift() string {
	i := strings.Index(r.path[1:], "/") + 1
	if i <= 0 {
		r.last = r.path[1:]
		r.path = "/"
	} else {
		r.last = r.path[1:i]
		r.path = r.path[i:]
	}
	return r.last
}

func (r *route) remaining() string {
	return r.path[1:]
}

// PathHeadThenShift returns the head segment of the request's adjusted path, then shifts it left by one segment. This
// does not adjust the path stored in req.URL.Path.
func PathHeadThenShift(req *http.Request) string {
	r, ok := req.Context().Value(routeKey).(*route)
	if !ok {
		return req.URL.Path
	}
	return r.shift()
}

// LastPathHead returns the last result obtained from a call to PathHeadThenShift() for the request.
func LastPathHead(req *http.Request) string {
	r, ok := req.Context().Value(routeKey).(*route)
	if !ok {
		return req.URL.Path
	}
	return r.last
}

// HasMorePathSegments returns true if more path segments will be returned from future calls to PathHeadThenShift() for
// the request.
func HasMorePathSegments(req *http.Request) bool {
	r, ok := req.Context().Value(routeKey).(*route)
	if !ok {
		return false
	}
	return r.path != "/"
}

// RemainingPath returns the remaining path for a request.
func RemainingPath(req *http.Request) string {
	r, ok := req.Context().Value(routeKey).(*route)
	if !ok {
		return req.URL.Path
	}
	return r.remaining()
}
