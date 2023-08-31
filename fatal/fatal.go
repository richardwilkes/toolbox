// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fatal

import (
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/errs"
)

// IfErr checks the error and if it isn't nil, calls fatal.WithErr(err).
func IfErr(err error) {
	if !toolbox.IsNil(err) {
		WithErr(err)
	}
}

// WithErr logs the error and then exits with code 1.
func WithErr(err error) {
	errs.Log(err)
	atexit.Exit(1)
}
