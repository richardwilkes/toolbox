package term

import (
	"fmt"
	"io"
	"strings"
)

const (
	defColumns = 80
	defRows    = 24
)

// WrapText prints the 'prefix' to 'out' and then wraps 'text' in the
// remaining space.
func WrapText(out io.Writer, prefix, text string) {
	fmt.Fprint(out, prefix)
	avail, _ := Size()
	avail -= 1 + len(prefix)
	if avail < 1 {
		avail = 1
	}
	remaining := avail
	indent := strings.Repeat(" ", len(prefix))
	for i, token := range strings.Fields(text) {
		length := len(token) + 1
		if i != 0 {
			if length > remaining {
				fmt.Fprintln(out)
				fmt.Fprint(out, indent)
				remaining = avail
			} else {
				fmt.Fprint(out, " ")
			}
		}
		fmt.Fprint(out, token)
		remaining -= length
	}
	fmt.Fprintln(out)
}
