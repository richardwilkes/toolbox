package netutil

import (
	"os/exec"
	"runtime"

	"github.com/richardwilkes/gokit/errs"
)

// OpenBrowser opens 'url' with the user's preferred browser.
func OpenBrowser(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "explorer"
	default:
		return errs.New("Unsupported platform")
	}
	if err := exec.Command(cmd, url).Start(); err != nil {
		return errs.NewWithCause("Unable to open "+url, err)
	}
	return nil
}
