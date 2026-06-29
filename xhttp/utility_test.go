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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xhttp"
)

func TestClientIP(t *testing.T) {
	for _, tc := range []struct {
		name          string
		xForwardedFor string
		forwarded     string
		remoteAddr    string
		want          string
	}{
		{name: "X-Forwarded-For wins", xForwardedFor: "203.0.113.7", forwarded: "for=192.0.2.60", remoteAddr: "10.0.0.1:80", want: "203.0.113.7"},
		{name: "X-Forwarded-For first of list", xForwardedFor: "203.0.113.7, 10.0.0.1", want: "203.0.113.7"},
		{name: "Forwarded simple", forwarded: "for=192.0.2.60", want: "192.0.2.60"},
		{name: "Forwarded uppercase param", forwarded: "For=192.0.2.60", want: "192.0.2.60"},
		{name: "Forwarded mixed-case param", forwarded: "FOR=192.0.2.60", want: "192.0.2.60"},
		{name: "Forwarded quoted IPv4 with port", forwarded: `for="192.0.2.60:4711"`, want: "192.0.2.60"},
		{name: "Forwarded quoted IPv6 with brackets and port", forwarded: `for="[2001:db8::1]:443"`, want: "2001:db8::1"},
		{name: "Forwarded quoted IPv6 with brackets no port", forwarded: `for="[2001:db8::1]"`, want: "2001:db8::1"},
		{name: "Forwarded bare IPv6", forwarded: "for=2001:db8::1", want: "2001:db8::1"},
		{name: "Forwarded with other params first", forwarded: "proto=https;for=192.0.2.60;by=203.0.113.43", want: "192.0.2.60"},
		{name: "Forwarded with whitespace", forwarded: "for = 192.0.2.60 ; proto=http", want: "192.0.2.60"},
		{name: "Forwarded multiple elements takes first", forwarded: "for=192.0.2.43, for=198.51.100.17", want: "192.0.2.43"},
		{name: "Forwarded obfuscated falls through to RemoteAddr", forwarded: "for=unknown", remoteAddr: "203.0.113.5:1234", want: "203.0.113.5"},
		{name: "RemoteAddr only", remoteAddr: "198.51.100.23:9999", want: "198.51.100.23"},
		{name: "RemoteAddr IPv6", remoteAddr: "[2001:db8::99]:443", want: "2001:db8::99"},
		{name: "nothing yields nil", want: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := check.New(t)
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.RemoteAddr = tc.remoteAddr
			if tc.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tc.xForwardedFor)
			}
			if tc.forwarded != "" {
				req.Header.Set("Forwarded", tc.forwarded)
			}
			got := xhttp.ClientIP(req)
			if tc.want == "" {
				c.True(got == nil, "%s: expected nil, got %v", tc.name, got)
				return
			}
			c.True(got != nil, "%s: expected %s, got nil", tc.name, tc.want)
			if got != nil {
				c.Equal(tc.want, got.String(), tc.name)
			}
		})
	}
}
