// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/i18n"
)

// Cmd represents a sub-command available on the command-line.
type Cmd interface {
	// Name should return the name of the command as it needs to be entered on
	// the command line.
	Name() string
	// Usage should return a description of what the command does.
	Usage() string
	// Run will be called to run the command. It will be passed a fresh
	// command line created from the command line that was parsed to determine
	// this command would be run, along with the arguments that have not yet
	// been consumed. The first argument, much like an application called from
	// main, will be the name of the command.
	Run(cmdLine *CmdLine, args []string) error
}

func (cl *CmdLine) newWithCmd(cmd Cmd) *CmdLine {
	cmdLine := New(false)
	cmdLine.out = cl.out
	cmdLine.parent = cl
	cmdLine.cmd = cmd
	return cmdLine
}

// AddCommand adds a command to the available commands.
func (cl *CmdLine) AddCommand(cmd Cmd) {
	if len(cl.cmds) == 0 {
		hc := &helpCmd{}
		cl.cmds[hc.Name()] = hc
	}
	cl.cmds[cmd.Name()] = cmd
}

// RunCommand attempts to run a command. Pass in the command line arguments,
// usually the result from calling Parse().
func (cl *CmdLine) RunCommand(args []string) error {
	if len(args) < 1 {
		cl.FatalMsg(i18n.Text("Must specify a command name"))
	}
	cmd := cl.cmds[args[0]]
	if cmd == nil {
		cl.FatalMsg(fmt.Sprintf(i18n.Text("'%[1]s' is not a valid %[2]s command"), args[0], AppCmdName))
	}
	return cmd.Run(cl.newWithCmd(cmd), args[1:])
}
