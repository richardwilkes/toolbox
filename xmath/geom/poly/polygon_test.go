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
	"time"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/xmath"
)

type ptCheck struct {
	pt        Point
	in        bool
	inEvenOdd bool
}

type containmentTestCase struct {
	p      Polygon
	checks []ptCheck
}

func TestContains(t *testing.T) {
	tests := []containmentTestCase{
		{
			p: Polygon{
				{{200, 20}, {300, 20}, {300, 120}, {200, 120}},
				{{250, 50}, {280, 50}, {280, 80}, {250, 80}},
				{{260, 60}, {290, 60}, {290, 90}, {260, 90}},
				{{290, 110}, {320, 110}, {320, 140}, {290, 140}},
			},
			checks: []ptCheck{
				{pt: Point{X: 199, Y: 20}},
				{pt: Point{X: 200, Y: 20}, in: true, inEvenOdd: true},
				{pt: Point{X: 300, Y: 120}, in: true, inEvenOdd: true},
				{pt: Point{X: 250, Y: 50}, in: true},
				{pt: Point{X: 260, Y: 60}, in: true, inEvenOdd: true},
				{pt: Point{X: 290, Y: 110}, in: true},
				{pt: Point{X: 319, Y: 139}, in: true, inEvenOdd: true},
				{pt: Point{X: 321, Y: 140}},
				{pt: Point{X: 320, Y: 141}},
			},
		},
		{
			p: Polygon{{{0, 0}, {10, 0}, {0, 10}}},
			checks: []ptCheck{
				{pt: Point{X: 1, Y: 1}, in: true, inEvenOdd: true},
				{pt: Point{X: 2, Y: 0.1}, in: true, inEvenOdd: true},
				{pt: Point{X: 10, Y: 10}},
				{pt: Point{X: 11, Y: 0}},
				{pt: Point{X: 0, Y: 11}},
				{pt: Point{X: -1, Y: -1}},
			},
		},
		{
			p: Polygon{{{0, 0}, {0, 10}, {10, 0}}},
			checks: []ptCheck{
				{pt: Point{X: 1, Y: 1}, in: true, inEvenOdd: true},
				{pt: Point{X: 2, Y: 0.1}, in: true, inEvenOdd: true},
				{pt: Point{X: 10, Y: 10}},
				{pt: Point{X: 11, Y: 0}},
				{pt: Point{X: 0, Y: 11}},
				{pt: Point{X: -1, Y: -1}},
			},
		},
		{
			p: Polygon{{{55, 35}, {25, 35}, {25, 119}, {55, 119}}},
			checks: []ptCheck{
				{pt: Point{X: 54.95, Y: 77}, in: true, inEvenOdd: true},
				{pt: Point{X: 55.05, Y: 77}},
			},
		},
		{
			p: Polygon{{{145, 35}, {145, 77}, {105, 77}, {105, 119}, {55, 119}, {55, 35}}},
			checks: []ptCheck{
				{pt: Point{X: 54.95, Y: 77}},
				{pt: Point{X: 55.05, Y: 77}, in: true, inEvenOdd: true},
			},
		},
	}
	for i, test := range tests {
		for j, tc := range test.checks {
			check.Equal(t, tc.in, test.p.Contains(tc.pt), "test case %d:%d", i, j)
			check.Equal(t, tc.inEvenOdd, test.p.ContainsEvenOdd(tc.pt), "test case for evenodd %d:%d", i, j)
		}
	}
}

type segCases struct {
	subject  Polygon
	clipping Polygon
}

func TestNonReductiveSegmentDivisions(t *testing.T) {
	tests := []segCases{
		{
			subject: Polygon{{
				{1.427255375e+06, -2.3283064365386963e-10},
				{1.4271285e+06, 134.7111358642578},
				{1.427109e+06, 178.30108642578125},
			}},
			clipping: Polygon{{
				{1.416e+06, -12000},
				{1.428e+06, -12000},
				{1.428e+06, 0},
				{1.416e+06, 0},
				{1.416e+06, -12000},
			}},
		},
		{
			subject: Polygon{{
				{1.7714672107465276e+06, -102506.68254093888},
				{1.7713768917571804e+06, -102000.75485953009},
				{1.7717109214841307e+06, -101912.19625031832},
			}},
			clipping: Polygon{{
				{1.7714593229229522e+06, -102470.35230830211},
				{1.7714672107465276e+06, -102506.68254093867},
				{1.771439738086082e+06, -102512.92027456204},
			}},
		},
		{
			subject: Polygon{{
				{-1.8280000000000012e+06, -492999.99999999953},
				{-1.8289999999999995e+06, -494000.0000000006},
				{-1.828e+06, -493999.9999999991},
				{-1.8280000000000012e+06, -492999.99999999953},
			}},
			clipping: Polygon{{
				{-1.8280000000000005e+06, -495999.99999999977},
				{-1.8280000000000007e+06, -492000.0000000014},
				{-1.8240000000000007e+06, -492000.0000000014},
				{-1.8280000000000005e+06, -495999.99999999977},
			}},
		},
		{
			subject: Polygon{{
				{-2.0199999999999988e+06, -394999.99999999825},
				{-2.0199999999999988e+06, -392000.0000000009},
				{-2.0240000000000012e+06, -395999.9999999993},
				{-2.0199999999999988e+06, -394999.99999999825},
			}},
			clipping: Polygon{{
				{-2.0199999999999988e+06, -394999.99999999825},
				{-2.020000000000001e+06, -394000.0000000001},
				{-2.0190000000000005e+06, -394999.9999999997},
				{-2.0199999999999988e+06, -394999.99999999825},
			}},
		},
		{
			subject: Polygon{{
				{-47999.99999999992, -23999.999999998756},
				{0, -24000.00000000017},
				{0, 24000.00000000017},
				{-48000.00000000014, 24000.00000000017},
				{-47999.99999999992, -23999.999999998756},
			}},
			clipping: Polygon{{{-48000, -24000}, {0, -24000}, {0, 24000}, {-48000, 24000}, {-48000, -24000}}},
		},
		{
			subject: Polygon{{
				{-2.137000000000001e+06, -122000.00000000093},
				{-2.1360000000000005e+06, -121999.99999999907},
				{-2.1360000000000014e+06, -121000.00000000186},
			}},
			clipping: Polygon{{
				{-2.1120000000000005e+06, -120000},
				{-2.136000000000001e+06, -120000.00000000093},
				{-2.1360000000000005e+06, -144000},
			}},
		},
		{
			subject: Polygon{{
				{1.556e+06, -1.139999999999999e+06},
				{1.5600000000000002e+06, -1.140000000000001e+06},
				{1.56e+06, -1.136000000000001e+06},
			}},
			clipping: Polygon{{{1.56e+06, -1.127999999999999e+06}, {1.5600000000000002e+06, -1.151999999999999e+06}}},
		},
		{
			subject: Polygon{{
				{1.0958876176594219e+06, -567467.5197556159},
				{1.0956330600760083e+06, -567223.72588934},
				{1.0958876176594219e+06, -567467.5197556159},
			}},
			clipping: Polygon{{
				{1.0953516248896217e+06, -564135.1861293605},
				{1.0959085007300845e+06, -568241.1879245406},
				{1.0955136237022132e+06, -581389.3748769956},
			}},
		},
		{
			subject: Polygon{{
				{608000, -113151.36476426799},
				{608000, -114660.04962779157},
				{612000, -115414.39205955336},
				{1.616e+06, -300000},
				{1.608e+06, -303245.6575682382},
				{0, 0},
			}},
			clipping: Polygon{{{1.612e+06, -296000}}},
		},
		{
			subject: Polygon{{
				{1.1458356382266793e+06, -251939.4635597784},
				{1.1460824662209095e+06, -251687.86194535438},
				{1.1458356382266793e+06, -251939.4635597784},
			}},
			clipping: Polygon{{
				{1.1486683769211173e+06, -251759.06331944838},
				{1.1468807511323579e+06, -251379.90576799586},
				{1.1457914974731328e+06, -251816.31287551578},
			}},
		},
		{
			subject: Polygon{{
				{426694.6365274183, -668547.1611580737},
				{426714.57523030025, -668548.9238652373},
				{426745.39648089616, -668550.4651249861},
			}},
			clipping: Polygon{{
				{426714.5752302991, -668548.9238652373},
				{426744.63718662335, -668550.0591896093},
				{426745.3964821229, -668550.4652243527},
			}},
		},
		{
			subject: Polygon{{
				{99.67054939325573, 23.50752393246498},
				{99.88993946188153, 20.999883973365655},
				{100.01468418889, 20.53433031419374},
			}},
			clipping: Polygon{{
				{100.15374164547939, 20.015360821030836},
				{95.64222842284941, 36.85255738690467},
				{100.15374164547939, -14.714274712355238},
			}},
		},
	}

	// Test multiple rotations across 360° of each case to catch any orientation assumptions.
	const rotations = 360
	for i, one := range tests {
		for j := 0; j < rotations; j++ {
			angle := 2 * xmath.Pi * fpType(j) / rotations
			subject := rotate(one.subject, angle)
			clipping := rotate(one.clipping, angle)
			// Using require to force termination, since otherwise this could go on for quite some time
			check.True(t, doTimedFunc(func(ch chan Polygon) { ch <- subject.Union(clipping) }), "test case union %d:%d", i, j)
			check.True(t, doTimedFunc(func(ch chan Polygon) { ch <- subject.Intersect(clipping) }), "test case intersect %d:%d", i, j)
			check.True(t, doTimedFunc(func(ch chan Polygon) { ch <- subject.Subtract(clipping) }), "test case subtract %d:%d", i, j)
			check.True(t, doTimedFunc(func(ch chan Polygon) { ch <- subject.Xor(clipping) }), "test case xor %d:%d", i, j)
		}
	}
}

func doTimedFunc(f func(chan Polygon)) bool {
	ch := make(chan Polygon)
	go f(ch)
	select {
	case <-ch: // check that we get a result in finite time
		return true
	case <-time.After(1 * time.Second):
		return false
	}
}

func rotate(p Polygon, radians fpType) Polygon {
	result := p.Clone()
	for i, contour := range p {
		for j, point := range contour {
			result[i][j] = Point{
				X: point.X*xmath.Cos(radians) - point.Y*xmath.Sin(radians),
				Y: point.Y*xmath.Cos(radians) + point.X*xmath.Sin(radians),
			}
		}
	}
	return result
}
