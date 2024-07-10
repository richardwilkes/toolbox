// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
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
	"crypto/subtle"
	"fmt"
	"net/http"
)

// BasicAuth provides basic HTTP authentication.
type BasicAuth struct {
	Realm string
	// Lookup provides a way to map a user in a realm to a password. The returned password should have already been
	// passed through the Hasher function.
	Lookup func(user, realm string) ([]byte, bool)
	Hasher func(input string) []byte
}

// Wrap an http.Handler, requiring Basic Authentication.
func (ba *BasicAuth) Wrap(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pw, ok := r.BasicAuth(); ok {
			stored, found := ba.Lookup(user, ba.Realm)
			passwordMatch := subtle.ConstantTimeCompare(ba.Hasher(pw), stored) == 1
			if found && passwordMatch {
				if md := MetadataFromRequest(r); md != nil {
					md.User = user
					md.Logger = md.Logger.With("user", user)
				}
				handler.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q, charset="UTF-8"`, ba.Realm))
		ErrorStatus(w, http.StatusUnauthorized)
	})
}
