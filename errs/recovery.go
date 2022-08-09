// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package errs

// RecoveryHandler defines the callback used when panics occur with Recovery.
type RecoveryHandler func(error)

// Recovery provides an easy way to run code that may panic. 'handler' will be called with the panic turned into an
// error. Pass in nil to silently ignore any panic.
//
// Typical usage:
//
//	func runSomeCode(handler errs.RecoveryHandler) {
//	    defer errs.Recovery(handler)
//	    // ... run the code here ...
//	}
func Recovery(handler RecoveryHandler) {
	if recovered := recover(); recovered != nil && handler != nil {
		err, ok := recovered.(error)
		if !ok {
			err = Newf("%+v", recovered)
		}
		defer Recovery(nil) // Guard against a bad handler implementation
		handler(NewWithCause("recovered from panic", err))
	}
}
