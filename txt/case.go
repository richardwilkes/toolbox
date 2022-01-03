// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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

// ToCamelCaseWithExceptions converts a string to CamelCase, but forces certain words to all caps.
func ToCamelCaseWithExceptions(in string, exceptions *AllCaps) string {
	out := ToCamelCase(in)
	pos := 0
	runes := []rune(out)
	rr := RuneReader{}
	for {
		rr.Src = runes[pos:]
		rr.Pos = 0
		matches := exceptions.regex.FindReaderIndex(&rr)
		if len(matches) == 0 {
			break
		}
		for i := matches[0] + 1; i < matches[1]; i++ {
			runes[pos+i] = unicode.ToUpper(runes[pos+i])
		}
		pos += matches[0] + 1
	}
	return string(runes)
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

// FirstToLower converts the first character to lower case.
func FirstToLower(in string) string {
	if in == "" {
		return in
	}
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError {
		return in
	}
	return string(unicode.ToLower(r)) + in[size:]
}
