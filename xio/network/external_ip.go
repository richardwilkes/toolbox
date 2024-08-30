// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package network

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/xio"
)

var sites = []string{
	// These seem to prefer ipv4 responses, if possible
	"http://whatismyip.akamai.com/",
	"https://myip.dnsomatic.com/",
	"http://api.ipify.org/",
	"http://checkip.amazonaws.com/",

	// These seem to prefer ipv6 responses, if possible
	"http://icanhazip.com/",
	"https://myexternalip.com/raw",
	"http://ifconfig.io/ip",
	"http://ident.me/",
}

// ExternalIP returns your IP address as seen by external sites. It does this by iterating through a list of websites
// that will return your IP address as they see it. The first response with a valid IP address will be returned. timeout
// sets the maximum amount of time for each attempt.
func ExternalIP(timeout time.Duration) string {
	client := &http.Client{Timeout: timeout}
	for _, site := range sites {
		if ip := externalIP(client, site); ip != "" {
			return ip
		}
	}
	return ""
}

func externalIP(client *http.Client, site string) string {
	if resp, err := client.Get(site); err == nil { //nolint:noctx // The timeout on the client provides the same effect
		defer xio.DiscardAndCloseIgnoringErrors(resp.Body)
		var body []byte
		if body, err = io.ReadAll(resp.Body); err == nil {
			if ip := net.ParseIP(strings.TrimSpace(string(body))); ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}
