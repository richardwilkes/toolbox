// +build darwin dragonfly freebsd linux netbsd openbsd

package term

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// IsTerminal returns true if the writer's file descriptor is a terminal.
func IsTerminal(f io.Writer) bool {
	var termios syscall.Termios
	switch v := f.(type) {
	case *os.File:
		_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, v.Fd(), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0) //nolint:gosec
		return errno == 0
	default:
		return false
	}
}
