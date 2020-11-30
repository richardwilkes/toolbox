// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio

import (
	"context"
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

// ContextSleep sleeps for the specified time, or until the context is done.
// You can check the return error to see if the context deadline was
// exceeded by using errors.Is(err, context.DeadlineExceeded).
func ContextSleep(ctx context.Context, waitTime time.Duration) error {
	timer := time.NewTimer(waitTime)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		err := ctx.Err()
		return errs.NewWithCause(err.Error(), err)
	case <-timer.C:
		return nil
	}
}
