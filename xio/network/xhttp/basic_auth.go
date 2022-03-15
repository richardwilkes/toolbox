// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package xhttp provides HTTP-related utilities.
package xhttp

import (
	"fmt"
	"net/http"
)

// PasswordLookup provides a way to map a user in a realm to a password
type PasswordLookup func(user, realm string) string

// BasicAuth provides basic HTTP authentication.
type BasicAuth struct {
	realm  string
	lookup PasswordLookup
}

// NewBasicAuth creates a new BasicAuth.
func NewBasicAuth(realm string, lookup PasswordLookup) *BasicAuth {
	return &BasicAuth{realm: realm, lookup: lookup}
}

// Wrap an http.Handler.
func (auth *BasicAuth) Wrap(handler http.Handler) http.Handler {
	return &wrapper{auth: auth, handler: handler}
}

type wrapper struct {
	auth    *BasicAuth
	handler http.Handler
}

func (hw *wrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if user, pw, ok := req.BasicAuth(); ok {
		if pw == hw.auth.lookup(user, hw.auth.realm) {
			hw.handler.ServeHTTP(w, req)
			return
		}
	}
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q`, hw.auth.realm))
	WriteHTTPStatus(w, http.StatusUnauthorized)
}
