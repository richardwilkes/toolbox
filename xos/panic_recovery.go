// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

import (
	"fmt"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// PanicRecovery provides an easy way to run code that may panic. If the handler is nil, the panic will be logged as an
// error.
//
// Typical usage:
//
//	func RunSomeCode() {
//	    defer xos.PanicRecovery(nil /* or provide a handler function */)
//	    // ... run the code here ...
//	}
func PanicRecovery(handler func(error)) {
	if recovered := recover(); recovered != nil {
		err, ok := recovered.(error)
		if !ok {
			err = fmt.Errorf("%+v", recovered)
		}
		err = errs.NewWithCause("recovered from panic", err)
		if handler == nil {
			errs.Log(err)
		} else {
			defer PanicRecovery(nil) // Guard against a bad handler implementation
			handler(err)
		}
	}
}

// SafeCall calls the provided function, safely wrapped by xos.PanicRecovery().
func SafeCall(f func(), handler func(error)) {
	defer PanicRecovery(handler)
	f()
}
