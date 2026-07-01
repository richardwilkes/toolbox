// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xhttp_test

import (
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xhttp"
)

const validUserName = "alice"

func sha256Hasher(input string) []byte {
	sum := sha256.Sum256([]byte(input))
	return sum[:]
}

func newBasicAuthHandler(hasher func(string) []byte) http.Handler {
	stored := map[string][]byte{validUserName: sha256Hasher("secret")}
	lookup := func(user, _ string) ([]byte, bool) {
		h, ok := stored[user]
		return h, ok
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok")) //nolint:errcheck // For test purposes, we don't care about the error from Write
	})
	return xhttp.BasicAuthWrap(next, lookup, hasher, "test-realm")
}

func TestBasicAuthWrap(t *testing.T) {
	c := check.New(t)
	handler := newBasicAuthHandler(sha256Hasher)

	for _, tc := range []struct {
		name       string
		user       string
		pw         string
		wantStatus int
		setAuth    bool
	}{
		{name: "valid credentials", user: validUserName, pw: "secret", wantStatus: http.StatusOK, setAuth: true},
		{name: "wrong password", user: validUserName, pw: "nope", wantStatus: http.StatusUnauthorized, setAuth: true},
		{name: "unknown user", user: "bob", pw: "secret", wantStatus: http.StatusUnauthorized, setAuth: true},
		{name: "no auth header", wantStatus: http.StatusUnauthorized, setAuth: false},
	} {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		if tc.setAuth {
			req.SetBasicAuth(tc.user, tc.pw)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		c.Equal(tc.wantStatus, rec.Code, tc.name)
		if tc.wantStatus == http.StatusOK {
			c.Equal("ok", rec.Body.String(), tc.name)
		} else {
			// Every rejection advertises the realm so clients can re-authenticate.
			c.NotEqual("", rec.Header().Get("WWW-Authenticate"), tc.name)
		}
	}
}

// TestBasicAuthWrapConstantWork guards the timing-side-channel fix: the supplied password must be hashed and compared
// regardless of whether the user exists, and an unknown user must be indistinguishable from a known user with a wrong
// password. If the middleware short-circuited (e.g. skipped hashing or the comparison when the user was not found), it
// would leak user existence through timing and this test would catch the regression.
func TestBasicAuthWrapConstantWork(t *testing.T) {
	c := check.New(t)

	var hasherCalls int
	handler := newBasicAuthHandler(func(input string) []byte {
		hasherCalls++
		return sha256Hasher(input)
	})

	serve := func(user, pw string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		req.SetBasicAuth(user, pw)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}

	// Known user with a wrong password: the password is hashed once and the request is rejected.
	hasherCalls = 0
	knownWrong := serve(validUserName, "wrong")
	c.Equal(1, hasherCalls)
	c.Equal(http.StatusUnauthorized, knownWrong.Code)

	// Unknown user: the password must still be hashed exactly once (no early-out that would reveal user existence).
	hasherCalls = 0
	unknown := serve("bob", "wrong")
	c.Equal(1, hasherCalls)
	c.Equal(http.StatusUnauthorized, unknown.Code)

	// The two rejections must be byte-for-byte identical so the response itself carries no user-existence signal.
	c.Equal(knownWrong.Code, unknown.Code)
	c.Equal(knownWrong.Body.String(), unknown.Body.String())
	c.Equal(knownWrong.Header().Get("WWW-Authenticate"), unknown.Header().Get("WWW-Authenticate"))
}
