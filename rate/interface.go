// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package rate

// Limiter provides a rate limiter.
type Limiter interface {
	// New returns a new limiter that is subordinate to this limiter, meaning that its cap rate is also capped by its
	// parent.
	New(capacity int) Limiter

	// Cap returns the capacity per time period.
	Cap(applyParentCaps bool) int

	// SetCap sets the capacity.
	SetCap(capacity int)

	// LastUsed returns the capacity used in the last time period.
	LastUsed() int

	// Use returns a channel that will return nil when the request is successful, or an error if the request cannot be
	// fulfilled.
	Use(amount int) <-chan error

	// Closed returns true if the limiter is closed.
	Closed() bool

	// Close this limiter and any children it may have.
	Close()
}
