package txt

import (
	"unicode"
	"unicode/utf8"
)

// ToCamelCase converts a string to CamelCase.
func ToCamelCase(in string) string {
	runes := []rune(in)
	out := make([]rune, 0, len(runes))
	up := true
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '_' {
			up = true
		} else {
			if up {
				r = unicode.ToUpper(r)
				up = false
			}
			out = append(out, r)
		}
	}
	return string(out)
}

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(in string) string {
	runes := []rune(in)
	out := make([]rune, 0, 1+len(runes))
	for i := 0; i < len(runes); i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

// FirstToUpper converts the first character to upper case.
func FirstToUpper(in string) string {
	if in == "" {
		return in
	}
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError {
		return in
	}
	return string(unicode.ToUpper(r)) + in[size:]
}
