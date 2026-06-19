// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xnet_test

import (
	"context"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xnet"
)

// TestExternalAddressesWithCancelledContext verifies that no network requests are attempted once the context is already
// canceled, and that nil is returned. This keeps the test hermetic (no outbound traffic).
func TestExternalAddressesWithCancelledContext(t *testing.T) {
	c := check.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Nil(xnet.ExternalIPAddress(ctx, time.Second))
	c.Nil(xnet.ExternalIPv4Address(ctx, time.Second))
	c.Nil(xnet.ExternalIPv6Address(ctx, time.Second))
}

func TestPrimaryIPAddress(t *testing.T) {
	c := check.New(t)
	// PrimaryIPAddress must always return a usable address, falling back to loopback when no route is available.
	ip := xnet.PrimaryIPAddress()
	c.NotNil(ip)
	c.True(ip.To4() != nil || ip.To16() != nil, "expected a valid IPv4 or IPv6 address, got %v", ip)
}
