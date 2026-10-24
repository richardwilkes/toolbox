// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xbytes_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xbytes"
)

func TestStripBOM(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		in   []byte
		want []byte
	}{
		{in: []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}, want: []byte{'h', 'e', 'l', 'l', 'o'}},
		{in: []byte{0xef, 0xbb, 0xbf}, want: []byte{}},
		{in: []byte{0xef, 0xbb, 0xbf, 'A'}, want: []byte{'A'}},
		{in: []byte{0xef, 0xbb, 0xbf, 0xc3, 0xa9}, want: []byte{0xc3, 0xa9}},                         // BOM + é
		{in: []byte{0xef, 0xbb, 0xbf, 0xf0, 0x9f, 0x98, 0x80}, want: []byte{0xf0, 0x9f, 0x98, 0x80}}, // BOM + 😀
		{in: []byte{'h', 'e', 'l', 'l', 'o'}, want: []byte{'h', 'e', 'l', 'l', 'o'}},
		{in: nil, want: nil},
		{in: []byte{}, want: []byte{}},
		{in: []byte{'A'}, want: []byte{'A'}},
		{in: []byte{'A', 'B'}, want: []byte{'A', 'B'}},
		{in: []byte{0xef}, want: []byte{0xef}},
		{in: []byte{0xef, 0xbb}, want: []byte{0xef, 0xbb}},
		{in: []byte{0xef, 'h', 'e', 'l', 'l', 'o'}, want: []byte{0xef, 'h', 'e', 'l', 'l', 'o'}},
		{in: []byte{0xef, 0xbb, 'h', 'e', 'l', 'l', 'o'}, want: []byte{0xef, 0xbb, 'h', 'e', 'l', 'l', 'o'}},
		{in: []byte{'h', 'e', 0xef, 0xbb, 0xbf, 'l', 'o'}, want: []byte{'h', 'e', 0xef, 0xbb, 0xbf, 'l', 'o'}},
		{in: []byte{'h', 'e', 'l', 'l', 'o', 0xef, 0xbb, 0xbf}, want: []byte{'h', 'e', 'l', 'l', 'o', 0xef, 0xbb, 0xbf}},
		{in: []byte{0xef, 0xbb, 0xbf, 0xef, 0xbb, 0xbf, 'h', 'i'}, want: []byte{0xef, 0xbb, 0xbf, 'h', 'i'}},
		{
			in:   []byte{0xef, 0xbb, 0xbe, 'h', 'e', 'l', 'l', 'o'}, // Last byte different
			want: []byte{0xef, 0xbb, 0xbe, 'h', 'e', 'l', 'l', 'o'},
		},
		{
			in:   []byte{0xef, 0xbc, 0xbf, 'h', 'e', 'l', 'l', 'o'}, // Middle byte different
			want: []byte{0xef, 0xbc, 0xbf, 'h', 'e', 'l', 'l', 'o'},
		},
		{
			in:   []byte{0xee, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}, // First byte different
			want: []byte{0xee, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'},
		},
		{in: []byte{0xef, 0xbb, 0xbf, 0x00, 'h', 'e', 'l', 'l', 'o'}, want: []byte{0x00, 'h', 'e', 'l', 'l', 'o'}},
		{in: []byte{0xef, 0xbb, 0xbf, 0x01, 0x02, 0x03, 0xff, 0xfe}, want: []byte{0x01, 0x02, 0x03, 0xff, 0xfe}},
		{in: append([]byte{0xef, 0xbb, 0xbf}, `{"key":"value"}`...), want: []byte(`{"key":"value"}`)},
		{in: append([]byte{0xef, 0xbb, 0xbf}, `<?xml version="1.0"?>`...), want: []byte(`<?xml version="1.0"?>`)},
	} {
		c.Equal(tc.want, xbytes.StripBOM(tc.in), "StripBOM(%v)", tc.in)
	}

	// StripBOM returns a sub-slice of its input rather than a copy, and leaves the input unmodified.
	withBOM := []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}
	result := xbytes.StripBOM(withBOM)
	c.True(&withBOM[3] == &result[0])
	c.Equal([]byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}, withBOM)
	noBOM := []byte{'h', 'e', 'l', 'l', 'o'}
	result = xbytes.StripBOM(noBOM)
	c.True(&noBOM[0] == &result[0])
}
