// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package toolbox

import (
	"cmp"
	"os"
	"os/user"
)

// CurrentUserName returns the current user's name. This will attempt to retrieve the user's display name, but will fall
// back to the account name if it isn't available.
func CurrentUserName() string {
	if u, err := user.Current(); err == nil {
		return cmp.Or(u.Name, u.Username)
	}
	return os.Getenv("USER")
}
