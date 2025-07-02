// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package desktop provides desktop integration utilities.
package xos

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xruntime"
)

// OpenBrowser asks the system to open the provided path or URL.
func OpenBrowser(pathOrURL string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case xruntime.MacOS:
		cmd = "open"
		args = append(args, pathOrURL)
	case xruntime.LinuxOS:
		cmd = "xdg-open"
		args = append(args, pathOrURL)
	case xruntime.WindowsOS:
		if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
			cmd = "cmd"
			args = append(args, "/c", "start", pathOrURL)
		} else {
			cmd = "explorer"
			args = append(args, pathOrURL)
		}
	default:
		return errs.New("unsupported OS")
	}
	if err := exec.Command(cmd, args...).Start(); err != nil {
		return errs.NewWithCause("unable to open "+pathOrURL, err)
	}
	return nil
}
