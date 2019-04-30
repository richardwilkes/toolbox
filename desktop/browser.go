// Package desktop provides desktop integration utilities.
package desktop

import (
	"os/exec"
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/errs"
)

// OpenBrowser opens 'url' with the user's preferred browser.
func OpenBrowser(url string) error {
	var cmd string
	switch runtime.GOOS {
	case toolbox.MacOS:
		cmd = "open"
	case toolbox.LinuxOS:
		cmd = "xdg-open"
	case toolbox.WindowsOS:
		cmd = "explorer"
	default:
		return errs.New("Unsupported platform")
	}
	if err := exec.Command(cmd, url).Start(); err != nil {
		return errs.NewWithCause("Unable to open "+url, err)
	}
	return nil
}
