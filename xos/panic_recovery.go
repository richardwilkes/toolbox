// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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

// PanicRecovery provides an easy way to run code that may panic. An optional handler may be passed in and will be
// called with the panic turned into an error. Note that even though the handler is passed in as a variadic parameter,
// only the first one will be used. This was done to allow passing in no handler at all, which will result in the panic
// being logged as an error.
//
// Typical usage:
//
//	func RunSomeCode() {
//	    defer xos.PanicRecovery()
//	    // ... run the code here ...
//	}
func PanicRecovery(handler ...func(error)) {
	if recovered := recover(); recovered != nil {
		err, ok := recovered.(error)
		if !ok {
			err = fmt.Errorf("%+v", recovered)
		}
		err = errs.NewWithCause("recovered from panic", err)
		if len(handler) == 0 || handler[0] == nil {
			errs.Log(err)
		} else {
			defer PanicRecovery() // Guard against a bad handler implementation
			handler[0](err)
		}
	}
}
