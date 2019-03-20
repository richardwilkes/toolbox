package txt

import (
	"strings"
)

// Wrap text to a certain length, giving it an optional prefix on each line.
// Words will not be broken, even if they exceed the maximum column size and
// instead will extend past the desired length.
func Wrap(prefix, text string, maxColumns int) string {
	var buffer strings.Builder
	buffer.WriteString(prefix)
	avail := maxColumns - len(prefix)
	for i, token := range strings.Fields(text) {
		if i != 0 {
			if 1+len(token) > avail {
				buffer.WriteByte('\n')
				buffer.WriteString(prefix)
				avail = maxColumns - len(prefix)
			} else {
				buffer.WriteByte(' ')
				avail--
			}
		}
		buffer.WriteString(token)
		avail -= len(token)
	}
	return buffer.String()
}
