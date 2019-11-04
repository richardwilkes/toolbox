package txt

// StripBOM removes the BOM marker from UTF-8 data, if present.
func StripBOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return b[3:]
	}
	return b
}
