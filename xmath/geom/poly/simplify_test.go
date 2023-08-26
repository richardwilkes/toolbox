// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
)

type testCaseSimplify struct {
	name   string
	poly   Polygon
	result Polygon
}

func TestSimplify(t *testing.T) {
	tests := []testCaseSimplify{
		{
			name: "Self-intersecting polygon",
			poly: Polygon{{{0, 0}, {1, 1}, {1, 0}, {0, 1}}},
			result: Polygon{
				{{0, 1}, {0, 0}, {0.5, 0.5}},
				{{1, 1}, {1, 0}, {0.5, 0.5}},
			},
		},
		{
			name:   "Polygon with repeated vertical edge",
			poly:   Polygon{{{0, 0}, {1, 0}, {1, 1}, {2, 1}, {2, 0}, {1, 0}, {1, 1}, {0, 1}}},
			result: Polygon{{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {2, 0}, {2, 1}}},
		},
		{
			name:   "Polygon with repeated horizontal edge",
			poly:   Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 2}, {1, 2}, {1, 1}, {0, 1}}},
			result: Polygon{{{0, 2}, {0, 1}, {0, 0}, {1, 0}, {1, 1}, {1, 2}}},
		},
		{
			name:   "Polygon with partially repeated edge",
			poly:   Polygon{{{0, 0}, {1, 0}, {1, 0.75}, {2, 0.75}, {2, 0.25}, {1, 0.25}, {1, 1}, {0, 1}}},
			result: Polygon{{{1, 0.75}, {1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 0.25}, {2, 0.25}, {2, 0.75}}},
		},
		{
			name: "Polygon with repeated edge in opposite direction",
			poly: Polygon{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
				{{1, 1}, {1, 0}, {2, 0}, {2, 1}},
			},
			result: Polygon{{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {2, 0}, {2, 1}}},
		},
		{
			name: "Polygon with partially repeated edge in opposite direction",
			poly: Polygon{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
				{{1, 0.25}, {2, 0.25}, {2, 0.75}, {1, 0.75}},
			},
			result: Polygon{{{1, 0.75}, {1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 0.25}, {2, 0.25}, {2, 0.75}}},
		},
		{
			name:   "Polygon with repeated dangling edge 1",
			poly:   Polygon{{{0, 0}, {1, 0}, {0, 0}, {1, 0}, {1, 1}, {0, 1}}},
			result: Polygon{{{0, 1}, {0, 0}, {1, 0}, {1, 1}}},
		},
		{
			name:   "Polygon with repeated dangling edge 2",
			poly:   Polygon{{{0, 0}, {1, 0}, {1, 1}, {1, 0}, {1, 1}, {0, 1}}},
			result: Polygon{{{0, 1}, {0, 0}, {1, 0}, {1, 1}}},
		},
		{
			name:   "Polygon with repeated dangling edge 3",
			poly:   Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {1, 1}, {0, 1}}},
			result: Polygon{{{0, 1}, {0, 0}, {1, 0}, {1, 1}}},
		},
		{
			name:   "Polygon with repeated dangling edge 4",
			poly:   Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}, {0, 1}}},
			result: Polygon{{{0, 1}, {0, 0}, {1, 0}, {1, 1}}},
		},
		{
			name: "Completely degenerate",
			poly: Polygon{{{1, 2}, {2, 2}, {2, 3}, {1, 2}, {2, 2}, {2, 3}}},
		},
	}
	for i, test := range tests {
		check.Equal(t, test.result, test.poly.Simplify(), "test case %d (%s)", i, test.name)
	}
}
