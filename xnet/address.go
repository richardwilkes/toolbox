// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/v2/xhttp"
	"github.com/richardwilkes/toolbox/v2/xio"
)

var (
	// These seem to prefer ipv4 responses, if possible
	v4Sites = []string{
		"https://whatismyip.akamai.com/",
		"https://myip.dnsomatic.com/",
		"https://api.ipify.org/",
		"https://checkip.amazonaws.com/",
		"https://4.ident.me/",
	}

	// These seem to prefer ipv6 responses, if possible
	v6Sites = []string{
		"https://icanhazip.com/",
		"https://myexternalip.com/raw",
		"https://ifconfig.io/ip",
		"https://6.ident.me/",
	}
)

// ExternalIPAddress returns your IP address as seen by external sites, by querying websites that report it. The IPv4
// sites are tried first (see ExternalIPv4Address), then the IPv6 sites (see ExternalIPv6Address). timeout bounds each
// individual site request. If no site returns a valid IP address, nil is returned.
func ExternalIPAddress(ctx context.Context, timeout time.Duration) net.IP {
	if v4 := ExternalIPv4Address(ctx, timeout); v4 != nil {
		return v4
	}
	return ExternalIPv6Address(ctx, timeout)
}

// ExternalIPv4Address returns your IPv4 address as seen by external sites. A list of websites that report it is queried
// concurrently and the first valid IPv4 address returned is used. timeout bounds each individual site request. If no
// site returns a valid IPv4 address, nil is returned.
func ExternalIPv4Address(ctx context.Context, timeout time.Duration) net.IP {
	return externalIPAddress(ctx, timeout, v4Sites, true)
}

// ExternalIPv6Address returns your IPv6 address as seen by external sites. A list of websites that report it is queried
// concurrently and the first valid IPv6 address returned is used. timeout bounds each individual site request. If no
// site returns a valid IPv6 address, nil is returned.
func ExternalIPv6Address(ctx context.Context, timeout time.Duration) net.IP {
	return externalIPAddress(ctx, timeout, v6Sites, false)
}

func externalIPAddress(ctx context.Context, timeout time.Duration, sites []string, v4 bool) net.IP {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	client := &http.Client{Timeout: timeout}
	results := make(chan net.IP, len(sites))
	for _, site := range sites {
		go func(addr string) {
			if ctx.Err() != nil {
				results <- nil
				return
			}
			data, err := xhttp.RetrieveData(ctx, client, addr)
			if err != nil {
				results <- nil
				return
			}
			ip := net.ParseIP(strings.TrimSpace(string(data)))
			if ip == nil {
				results <- nil
				return
			}
			switch {
			case v4:
				// Reject IPv6 responses (To4 returns nil for them).
				results <- ip.To4()
			case ip.To4() == nil:
				results <- ip
			default:
				// Reject IPv4 responses to an IPv6 query.
				results <- nil
			}
		}(site)
	}
	for range sites {
		var ip net.IP
		select {
		case ip = <-results:
			if ip != nil {
				return ip
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

// PrimaryIPAddress returns the local IP address used to reach external hosts, or 127.0.0.1 if it cannot be determined.
func PrimaryIPAddress() net.IP {
	// A UDP dial sends nothing; it just selects the local address that would be used to reach an external host.
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer xio.CloseIgnoringErrors(conn)
		if localAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			return localAddr.IP
		}
	}
	return net.IPv4(127, 0, 0, 1)
}
