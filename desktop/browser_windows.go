// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package desktop provides desktop integration utilities.
package desktop

import "strings"

func cmdAndArgs(pathOrURL string) (cmdName string, args []string) {
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		return "cmd", []string{"/c", "start", pathOrURL}
	}
	return "explorer", []string{pathOrURL}
}
