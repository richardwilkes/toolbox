// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

import (
	"fmt"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/i18n"
)

type helpCmd struct{}

// Name implements the Cmd interface.
func (c *helpCmd) Name() string {
	return "help"
}

// Usage implements the Cmd interface.
func (c *helpCmd) Usage() string {
	return i18n.Text("Display help information for a command and exit.")
}

// Run implements the Cmd interface.
func (c *helpCmd) Run(cmdLine *CmdLine, args []string) error {
	cmdLine = cmdLine.parent
	if len(args) > 0 {
		if args[0] != "help" {
			const helpFlag = "-h"
			if helpFlag != args[0] {
				for name, cmd := range cmdLine.cmds {
					if name == args[0] {
						return cmd.Run(cmdLine.newWithCmd(cmd), []string{helpFlag})
					}
				}
			}
			fmt.Fprintf(cmdLine, i18n.Text("'%[1]s' is not a valid %[2]s command\n"), args[0], AppCmdName)
		}
	}
	cmdLine.DisplayUsage()
	atexit.Exit(1)
	return nil
}
