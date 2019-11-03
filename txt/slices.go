package txt

// StringSliceToMap returns a map created from the strings of a slice.
func StringSliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool)
	for _, str := range slice {
		m[str] = true
	}
	return m
}

// MapToStringSlice returns a slice created from the keys of a map.
func MapToStringSlice(m map[string]bool) []string {
	s := make([]string, 0, len(m))
	for str := range m {
		s = append(s, str)
	}
	return s
}
