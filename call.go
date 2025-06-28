// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package toolbox

import "github.com/richardwilkes/toolbox/v2/xos"

// TODO: Revisit and probably consolidate into xos

// Call the provided function, safely wrapped in a xos.PanicRecovery() handler that logs any errors via slog.
func Call(f func()) {
	CallWithHandler(f, nil)
}

// CallWithHandler calls the provided function, safely wrapped in a xos.PanicRecovery() handler.
func CallWithHandler(f func(), errHandler func(err error)) {
	defer xos.PanicRecovery(errHandler)
	f()
}
