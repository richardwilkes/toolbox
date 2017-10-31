package txt

import "bytes"

// CollapseSpaces removes leading and trailing spaces and reduces any runs of
// two or more spaces to a single space.
func CollapseSpaces(in string) string {
	var buffer bytes.Buffer
	lastWasSpace := false
	for i, r := range in {
		if r == ' ' {
			if !lastWasSpace {
				if i != 0 {
					buffer.WriteByte(' ')
				}
				lastWasSpace = true
			}
		} else {
			buffer.WriteRune(r)
			lastWasSpace = false
		}
	}
	if lastWasSpace && buffer.Len() > 0 {
		buffer.Truncate(buffer.Len() - 1)
	}
	return buffer.String()
}
