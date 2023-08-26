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

func TestClipLine(t *testing.T) {
	tests := []testCase{
		{
			name:    "clip line",
			subject: Polygon{{{0, 1}, {1.25, 1}, {1.5, 1.1}, {1.75, 1}, {5, 1}, {5, 2}, {0, 2}}},
			clipping: Polygon{
				{{1, 0}, {4, 0}, {4, 3}, {1, 3}},
				{{2, 0.5}, {3, 0.5}, {3, 2.5}, {2, 2.5}},
			},
			expected: Polygon{
				{{2, 1}, {1.75, 1}, {1.5, 1.1}, {1.25, 1}, {1, 1}},
				{{2, 2}, {1, 2}},
				{{4, 1}, {3, 1}},
				{{4, 2}, {3, 2}},
			},
		},
		{
			name:     "Clip line within 1",
			subject:  Polygon{{{-3999, -3999}, {-3500, -3500}}},
			clipping: Polygon{{{-4000, -4000}, {0, -4000}, {0, 0}, {-4000, 0}, {-4000, -4000}}},
			expected: Polygon{{{-3500, -3500}, {-3999, -3999}}},
		},
		{
			name:     "Clip line within 2",
			subject:  Polygon{{{X: 1.893757843025658e+06, Y: 358279.0127257189}, {X: 1.893986642180132e+06, Y: 359465.8124818327}, {X: 1.893983849777607e+06, Y: 359429.8946016282}}},
			clipping: Polygon{{{X: 1.89e+06, Y: 340000}, {X: 1.91e+06, Y: 340000}, {X: 1.91e+06, Y: 360000}, {X: 1.89e+06, Y: 360000}, {X: 1.89e+06, Y: 340000}}},
			expected: Polygon{{{X: 1.893757843025658e+06, Y: 358279.0127257189}, {X: 1.893986642180132e+06, Y: 359465.8124818327}, {X: 1.893983849777607e+06, Y: 359429.8946016282}}},
		},
	}
	for i, test := range tests {
		check.Equal(t, test.expected, test.subject.ClipLine(test.clipping), "test case %d (%s)", i, test.name)
	}
}
