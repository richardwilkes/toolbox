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

	"github.com/stretchr/testify/assert"
)

type segmentTestCase struct {
	name  string
	s1    Segment
	s2    Segment
	count int
	ip1   Point
}

func TestFindIntersection(t *testing.T) {
	tests := []segmentTestCase{
		{
			name: "Almost (but not) parallel lines 1",
			s1:   Segment{End: Point{X: 100, Y: 0.0001}},
			s2:   Segment{Start: Point{X: 1}, End: Point{X: 100}},
		},
		{
			name: "Almost (but not) parallel lines 2",
			s1:   Segment{End: Point{X: 100, Y: 0.0000001}},
			s2:   Segment{Start: Point{X: 1}, End: Point{X: 100}},
		},
		{
			name:  "Cross",
			s1:    Segment{Start: Point{X: 1}, End: Point{X: 1, Y: 3}},
			s2:    Segment{Start: Point{Y: 1}, End: Point{X: 3, Y: 1}},
			count: 1,
			ip1:   Point{X: 1, Y: 1},
		},
		{
			name:  "Rays",
			s1:    Segment{Start: Point{Y: 1}, End: Point{X: 1, Y: 3}},
			s2:    Segment{Start: Point{Y: 1}, End: Point{X: 3, Y: 1}},
			count: 1,
			ip1:   Point{Y: 1},
		},
		{
			name:  "Colinear rays",
			s1:    Segment{Start: Point{X: 2, Y: 1}, End: Point{Y: 1}},
			s2:    Segment{Start: Point{X: 2, Y: 1}, End: Point{X: 1, Y: 1}},
			count: 2,
			ip1:   Point{X: 2, Y: 1},
		},
		{
			name:  "Colinear rays",
			s1:    Segment{Start: Point{Y: 3}, End: Point{Y: 1}},
			s2:    Segment{Start: Point{Y: 3}, End: Point{Y: 2}},
			count: 2,
			ip1:   Point{Y: 3},
		},
		{
			name:  "Overlapping segments 1",
			s1:    Segment{Start: Point{Y: 1}, End: Point{X: 3, Y: 1}},
			s2:    Segment{Start: Point{X: 1, Y: 1}, End: Point{X: 2, Y: 1}},
			count: 2,
			ip1:   Point{X: 1, Y: 1},
		},
		{
			name:  "Overlapping segments 2",
			s1:    Segment{Start: Point{Y: 1}, End: Point{Y: 4}},
			s2:    Segment{Start: Point{Y: 2}, End: Point{Y: 3}},
			count: 2,
			ip1:   Point{Y: 2},
		},
		{
			name:  "Overlapping segments 3",
			s1:    Segment{Start: Point{X: 43.2635182233307, Y: 170.15192246987792}, End: Point{X: 41.57979856674331, Y: 170.60307379214092}},
			s2:    Segment{Start: Point{X: 43.2635182233307, Y: 170.15192246987792}, End: Point{X: 42.78116786015871, Y: 170.28116786015872}},
			count: 2,
			ip1:   Point{X: 43.2635182233307, Y: 170.15192246987792},
		},
		{
			name:  "Overlapping segments 4",
			s1:    Segment{Start: Point{X: 41.57979856674331, Y: 170.60307379214092}, End: Point{X: 43.2635182233307, Y: 170.15192246987792}},
			s2:    Segment{Start: Point{X: 42.78116786015871, Y: 170.28116786015872}, End: Point{X: 43.2635182233307, Y: 170.15192246987792}},
			count: 2,
			ip1:   Point{X: 42.78116786015871, Y: 170.28116786015872},
		},
		{
			name:  "Overlapping segments 5",
			s1:    Segment{Start: Point{X: 43.2635182233307, Y: 170.15192246987792}, End: Point{X: 41.57979856674331, Y: 170.60307379214092}},
			s2:    Segment{Start: Point{X: 42.78116786015871, Y: 170.28116786015872}, End: Point{X: 43.2635182233307, Y: 170.15192246987792}},
			count: 2,
			ip1:   Point{X: 43.2635182233307, Y: 170.15192246987792},
		},
		{
			name:  "Overlapping segments 2",
			s1:    Segment{Start: Point{X: 41.57979856674331, Y: 170.60307379214092}, End: Point{X: 43.2635182233307, Y: 170.15192246987792}},
			s2:    Segment{Start: Point{X: 43.2635182233307, Y: 170.15192246987792}, End: Point{X: 42.78116786015871, Y: 170.28116786015872}},
			count: 2,
			ip1:   Point{X: 43.2635182233307, Y: 170.15192246987792},
		},
		{
			name:  "Identical segments",
			s1:    Segment{Start: Point{X: 66, Y: 160}, End: Point{X: 67.1242262770966, Y: 147.15003485264717}},
			s2:    Segment{Start: Point{X: 66, Y: 160}, End: Point{X: 67.1242262770966, Y: 147.15003485264717}},
			count: 2,
			ip1:   Point{X: 66, Y: 160},
		},
	}
	for i, test := range tests {
		num, ip1, _ := test.s1.FindIntersection(test.s2, true)
		assert.Equal(t, test.count, num, "test case %d (%s)", i, test.name)
		assert.Equal(t, test.ip1, ip1, "test case %d (%s)", i, test.name)
	}
}
