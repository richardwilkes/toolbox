package jotrotate

import (
	"io"
	"os"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/log/rotation"
	"github.com/richardwilkes/toolbox/xio"
)

// ParseAndSetup adds command-line options for controlling logging, parses the
// command line, then instantiates a rotator and attaches it to jot. Returns
// the remaining arguments that weren't used for option content.
func ParseAndSetup(cl *cmdline.CmdLine) []string {
	logFile := rotation.DefaultPath()
	var maxSize int64 = rotation.DefaultMaxSize
	maxBackups := rotation.DefaultMaxBackups
	quiet := false
	cl.NewStringOption(&logFile).SetSingle('l').SetName("log-file").SetUsage("The file to write logs to")
	cl.NewInt64Option(&maxSize).SetName("log-file-size").SetUsage("The maximum number of bytes to write to a log file before rotating it")
	cl.NewIntOption(&maxBackups).SetName("log-file-backups").SetUsage("The maximum number of old logs files to retain")
	cl.NewBoolOption(&quiet).SetSingle('q').SetName("log-only").SetUsage("Do not copy the log output to the console")
	remainingArgs := cl.Parse(os.Args[1:])
	if rotator, err := rotation.New(rotation.Path(logFile), rotation.MaxSize(maxSize), rotation.MaxBackups(maxBackups)); err == nil {
		if quiet {
			jot.SetWriter(rotator)
		} else {
			jot.SetWriter(&xio.TeeWriter{Writers: []io.Writer{rotator, os.Stdout}})
		}
	} else {
		jot.Error(err)
	}
	return remainingArgs
}
