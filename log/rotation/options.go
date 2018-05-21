package rotation

import (
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

// Path specifies the file to write logs to. Backup log files will be retained
// in the same directory. Defaults to <cmdline.AppCmdName>.log in the
// os.TempDir().
func Path(path string) func(*Rotator) error {
	return func(r *Rotator) error {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			path = path[:len(path)-4]
		}
		if path == "" {
			return errs.New("Must specify a path")
		}
		r.path = path
		return nil
	}
}

// MaxSize sets the maximum size of the log file before it gets rotated.
// Defaults to 10 MB.
func MaxSize(maxSize int64) func(*Rotator) error {
	return func(r *Rotator) error {
		r.maxSize = maxSize
		return nil
	}
}

// MaxBackups sets the maximum number of old log files to retain.  The default
// is to retain 1.
func MaxBackups(maxBackups int) func(*Rotator) error {
	return func(r *Rotator) error {
		r.maxBackups = maxBackups
		return nil
	}
}
