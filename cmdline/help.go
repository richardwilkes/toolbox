package cmdline

import (
	"fmt"

	"github.com/richardwilkes/gokit/atexit"
	"github.com/richardwilkes/gokit/i18n"
)

type helpCmd struct {
}

// Name implements the Cmd interface.
func (c *helpCmd) Name() string {
	return "help"
}

// Usage implements the Cmd interface.
func (c *helpCmd) Usage() string {
	return i18n.Text("Display help information for a command and exit")
}

// Run implements the Cmd interface.
func (c *helpCmd) Run(cmdLine *CmdLine, args []string) error {
	cmdLine = cmdLine.parent
	if len(args) > 0 {
		if "help" != args[0] {
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
