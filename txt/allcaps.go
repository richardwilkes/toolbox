// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package txt provides various text utilities.
package txt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
)

// StdAllCaps provides the standard list of words that golint expects to be
// capitalized, found in the variable 'commonInitialisms' in
// https://github.com/golang/lint/blob/master/lint.go#L771-L808
var StdAllCaps = MustNewAllCaps("acl", "api", "ascii", "cpu", "css", "dns",
	"eof", "guid", "html", "http", "https", "id", "ip", "json", "lhs", "qps",
	"ram", "rhs", "rpc", "sla", "smtp", "sql", "ssh", "tcp", "tls", "ttl",
	"udp", "ui", "uid", "uuid", "uri", "url", "utf8", "vm", "xml", "xmpp",
	"xsrf", "xss")

// AllCaps holds information for transforming text with
// ToCamelCaseWithExceptions.
type AllCaps struct {
	regex *regexp.Regexp
}

// NewAllCaps takes a list of words that should be all uppercase when part of
// a camel-cased string.
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

// MustNewAllCaps takes a list of words that should be all uppercase when part
// of a camel-cased string. Failure to create the AllCaps object causes the
// program to exit.
func MustNewAllCaps(in ...string) *AllCaps {
	result, err := NewAllCaps(in...)
	if err != nil {
		jot.Fatal(1, err)
	}
	return result
}
