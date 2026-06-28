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
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xhttp"
)

// TestStatusWriterUnwrapReachesDeadlines verifies that http.ResponseController can reach the underlying writer's
// deadline methods through the StatusWriter, which it can only do if StatusWriter implements Unwrap. Since the server
// wraps every response in a StatusWriter, this is what keeps per-request deadline control working for handlers.
func TestStatusWriterUnwrapReachesDeadlines(t *testing.T) {
	c := check.New(t)
	underlying := &deadlineWriter{ResponseWriter: httptest.NewRecorder()}
	sw := xhttp.NewStatusWriter(underlying, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	rc := http.NewResponseController(sw)
	deadline := time.Now().Add(time.Minute)
	c.NoError(rc.SetReadDeadline(deadline))
	c.NoError(rc.SetWriteDeadline(deadline))
	c.Equal(deadline, underlying.readDeadline)
	c.Equal(deadline, underlying.writeDeadline)
}
