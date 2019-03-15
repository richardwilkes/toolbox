// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package term

import (
	"syscall"
	"unsafe"
)

// Size returns the number of columns and rows comprising the terminal.
func Size() (columns, rows int) {
	var ws [4]uint16
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws[0]))); errno == 0 { //nolint:gosec
		return int(ws[1]), int(ws[0])
	}
	return defColumns, defRows
}
