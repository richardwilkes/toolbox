// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom_test

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/xmath/geom"
)

type Pt = geom.Point[float64]

//go:embed orient.txt
var orientData []byte

func TestOrient(t *testing.T) {
	check.True(t, geom.Orient(Pt{}, Pt{X: 1, Y: 1}, Pt{Y: 1}) < 0, "clockwise")
	check.True(t, geom.Orient(Pt{}, Pt{Y: 1}, Pt{X: 1, Y: 1}) > 0, "counterclockwise")
	check.True(t, geom.Orient(Pt{}, Pt{X: 0.5, Y: 0.5}, Pt{X: 1, Y: 1}) == 0, "collinear")

	line := 0
	s := bufio.NewScanner(bytes.NewBuffer(orientData))
	for s.Scan() {
		line++
		parts := strings.Split(s.Text(), " ")
		if len(parts) != 8 {
			fmt.Printf("skipped line #%d\n", line)
			continue
		}
		var err error
		var a, b, c Pt
		var sign int64
		a, err = parsePoint(parts[1], parts[2])
		check.NoError(t, err, "parsing point A for line #%d", line)
		b, err = parsePoint(parts[3], parts[4])
		check.NoError(t, err, "parsing point B for line #%d", line)
		c, err = parsePoint(parts[5], parts[6])
		check.NoError(t, err, "parsing point C for line #%d", line)
		sign, err = strconv.ParseInt(parts[7], 10, 64)
		sign = -sign // sign field in test file is inverted in original data source, so flip it to match expectations
		check.NoError(t, err, "parsing sign for line #%d", line)
		result := geom.Orient(a, b, c)
		check.Equal(t, sign < 0, result < 0, "checking line #%d: expected %v, got %v (result: %#v)", line, sign < 0, result < 0, result)
	}
	check.NoError(t, s.Err())
}

func parsePoint(xs, ys string) (pt Pt, err error) {
	if pt.X, err = strconv.ParseFloat(xs, 64); err != nil {
		return
	}
	pt.Y, err = strconv.ParseFloat(ys, 64)
	return
}
