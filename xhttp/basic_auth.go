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
	"crypto/subtle"
	"fmt"
	"net/http"
)

// BasicAuthWrap wraps the given handler, providing basic HTTP authentication. Populates the User field of Metadata and
// adjusts the logger within the Metadata to add a user attribute.
func BasicAuthWrap(next http.Handler, lookup func(user, realm string) ([]byte, bool), hasher func(input string) []byte, realm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pw, ok := r.BasicAuth(); ok {
			stored, found := lookup(user, realm)
			passwordMatch := subtle.ConstantTimeCompare(hasher(pw), stored) == 1
			if found && passwordMatch {
				if md := MetadataFromRequest(r); md != nil {
					md.User = user
					md.Logger = md.Logger.With("user", user)
				}
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q, charset="UTF-8"`, realm))
		ErrorStatus(w, http.StatusUnauthorized)
	})
}
