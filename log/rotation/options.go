package rotation

import (
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
)

// Constants for defaults.
const (
	DefaultMaxSize    = 10 * 1024 * 1024
	DefaultMaxBackups = 1
)

// DefaultPath returns the default path that will be used.
func DefaultPath() string {
	return cmdline.AppCmdName + ".log"
}

// Path specifies the file to write logs to. Backup log files will be retained
// in the same directory. Defaults to the value of DefaultPath().
func Path(path string) func(*Rotator) error {
	return func(r *Rotator) error {
		if path == "" {
			return errs.New("Must specify a path")
		}
		r.path = path
		return nil
	}
}

// MaxSize sets the maximum size of the log file before it gets rotated.
// Defaults to DefaultMaxSize.
func MaxSize(maxSize int64) func(*Rotator) error {
	return func(r *Rotator) error {
		r.maxSize = maxSize
		return nil
	}
}

// MaxBackups sets the maximum number of old log files to retain.  Defaults
// to DefaultMaxBackups.
func MaxBackups(maxBackups int) func(*Rotator) error {
	return func(r *Rotator) error {
		r.maxBackups = maxBackups
		return nil
	}
}
