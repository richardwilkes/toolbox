// genversion provides a simple way to generate a version with build date and
// time of 'now'.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
)

func main() {
	var major, minor, patch int
	cmdline.CopyrightYears = "2016-2017"
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.License = "Mozilla Public License 2.0"
	cl := cmdline.New(true)
	cl.NewIntOption(&major).SetName("major").SetUsage(i18n.Text("The major version number"))
	cl.NewIntOption(&minor).SetName("minor").SetUsage(i18n.Text("The minor version number"))
	cl.NewIntOption(&patch).SetName("patch").SetUsage(i18n.Text("The patch number"))
	cl.Description = i18n.Text("Generates a version number with a build time of now.")
	cl.Parse(os.Args[1:])
	fmt.Println(cmdline.NewVersion(major, minor, patch, time.Now()).Format(false, true))
}
