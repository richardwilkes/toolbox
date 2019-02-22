package paths

import (
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/cmdline"
)

// AppLogDir returns the application log directory.
func AppLogDir() string {
	var path string
	if u, err := user.Current(); err == nil {
		path = u.HomeDir
		switch runtime.GOOS {
		case toolbox.MacOS:
			path = filepath.Join(path, "Library", "Logs")
		case toolbox.WindowsOS:
			path = filepath.Join(path, "AppData")
		default:
			path = filepath.Join(path, ".logs")
		}
		if cmdline.AppIdentifier != "" {
			path = filepath.Join(path, cmdline.AppIdentifier)
		}
	}
	return path
}

// AppDataDir returns the application data directory.
func AppDataDir() string {
	var path string
	if u, err := user.Current(); err == nil {
		path = u.HomeDir
		switch runtime.GOOS {
		case toolbox.MacOS:
			path = filepath.Join(path, "Library", "Application Support")
		case toolbox.WindowsOS:
			path = filepath.Join(path, "AppData")
		default:
			path = filepath.Join(path, ".appdata")
		}
		if cmdline.AppIdentifier != "" {
			path = filepath.Join(path, cmdline.AppIdentifier)
		}
	}
	return path
}
