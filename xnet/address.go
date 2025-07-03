// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xnet

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/v2/xio"
)

var (
	// These seem to prefer ipv4 responses, if possible
	v4Sites = []string{
		"http://whatismyip.akamai.com/",
		"https://myip.dnsomatic.com/",
		"http://api.ipify.org/",
		"http://checkip.amazonaws.com/",
		"http://4.ident.me/",
	}

	// These seem to prefer ipv6 responses, if possible
	v6Sites = []string{
		"http://icanhazip.com/",
		"https://myexternalip.com/raw",
		"http://ifconfig.io/ip",
		"http://6.ident.me/",
	}
)

// ExternalIPAddress returns your IP address as seen by external sites. It does this by iterating through a list of
// websites that will return your IP address as they see it. The first response with a valid IP address will be
// returned. timeout sets the maximum amount of time for each attempt to connect with a site. If no valid IP address is
// returned from any of the sites, nil is returned. Sites that usually return IPv4 addresses are checked first, then
// those that usually return IPv6 addresses.
func ExternalIPAddress(ctx context.Context, timeout time.Duration) net.IP {
	if v4 := ExternalIPv4Address(ctx, timeout); v4 != nil {
		return v4
	}
	return ExternalIPv6Address(ctx, timeout)
}

// ExternalIPv4Address returns your IPv4 address as seen by external sites. It does this by iterating through a list of
// websites that will return your IPv4 address as they see it. The first response with a valid IPv4 address will be
// returned. timeout sets the maximum amount of time for each attempt to connect with a site. If no valid IPv4 address
// is returned from any of the sites, nil is returned.
func ExternalIPv4Address(ctx context.Context, timeout time.Duration) net.IP {
	return externalIPAddress(ctx, timeout, v4Sites, true)
}

// ExternalIPv6Address returns your IPv6 address as seen by external sites. It does this by iterating through a list of
// websites that will return your IPv6 address as they see it. The first response with a valid IPv6 address will be
// returned. timeout sets the maximum amount of time for each attempt to connect with a site. If no valid IPv6 address
// is returned from any of the sites, nil is returned.
func ExternalIPv6Address(ctx context.Context, timeout time.Duration) net.IP {
	return externalIPAddress(ctx, timeout, v6Sites, false)
}

func externalIPAddress(ctx context.Context, timeout time.Duration, sites []string, v4 bool) net.IP {
	client := &http.Client{Timeout: timeout}
	for _, site := range sites {
		if ctx.Err() != nil {
			return nil
		}
		if data, err := xio.RetrieveData(ctx, client, site); err == nil {
			ip := net.ParseIP(strings.TrimSpace(string(data)))
			if v4 {
				ip = ip.To4()
			}
			if ip != nil {
				return ip
			}
		}
	}
	return nil
}

// PrimaryIPAddress returns the primary IP address.
func PrimaryIPAddress() net.IP {
	// Since we're using udp, no connection will actually be made. We just need an external IP address.
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		if localAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			return localAddr.IP
		}
	}
	return net.IPv4(127, 0, 0, 1)
}
