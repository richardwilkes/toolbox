package txt

// FirstN returns the first n runes of a string.
func FirstN(s string, n int) string {
	if n < 1 {
		return ""
	}
	r := []rune(s)
	if n > len(r) {
		return s
	}
	return string(r[:n])
}

// LastN returns the last n runes of a string.
func LastN(s string, n int) string {
	if n < 1 {
		return ""
	}
	r := []rune(s)
	if n > len(r) {
		return s
	}
	return string(r[len(r)-n:])
}
