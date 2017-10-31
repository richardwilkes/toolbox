// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd

package term

import (
	"io"
)

// IsTerminal returns true if the writer's file descriptor is a terminal.
func IsTerminal(f io.Writer) bool {
	return false
}
