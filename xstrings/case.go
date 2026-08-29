// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xos"
)

// StdAllCaps is the list of initialisms golint expects to be all caps, from 'commonInitialisms' in
// https://github.com/golang/lint/blob/master/lint.go#L771-L808
var StdAllCaps = MustNewAllCaps(
	"acl",
	"api",
	"ascii",
	"cpu",
	"css",
	"dns",
	"eof",
	"guid",
	"html",
	"http",
	"https",
	"id",
	"ip",
	"json",
	"lhs",
	"qps",
	"ram",
	"rhs",
	"rpc",
	"sla",
	"smtp",
	"sql",
	"ssh",
	"tcp",
	"tls",
	"ttl",
	"udp",
	"ui",
	"uid",
	"uuid",
	"uri",
	"url",
	"utf8",
	"vm",
	"xml",
	"xmpp",
	"xsrf",
	"xss",
)

// AllCaps holds the words that ToCamelCaseWithExceptions forces to all caps.
type AllCaps struct {
	regex *regexp.Regexp
}

// NewAllCaps creates an AllCaps from words that should be all uppercase when part of a camel-cased string.
func NewAllCaps(in ...string) (*AllCaps, error) {
	var buffer strings.Builder
	for _, str := range in {
		if buffer.Len() > 0 {
			buffer.WriteByte('|')
		}
		buffer.WriteString(FirstToUpper(strings.ToLower(str)))
	}
	r, err := regexp.Compile(fmt.Sprintf("(%s)(?:$|[A-Z])", buffer.String()))
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &AllCaps{regex: r}, nil
}

// MustNewAllCaps is like NewAllCaps, but panics if the AllCaps cannot be created.
func MustNewAllCaps(in ...string) *AllCaps {
	return xos.Must(NewAllCaps(in...))
}

// ToCamelCase converts a string to CamelCase.
func ToCamelCase(in string) string {
	runes := []rune(in)
	out := make([]rune, 0, len(runes))
	up := true
	for _, r := range runes {
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

// ToCamelCaseWithExceptions converts a string to CamelCase, forcing the words in exceptions to all caps.
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
	for i := range runes {
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
