// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"golang.org/x/exp/constraints"
)

type endPoint[T constraints.Float] struct {
	angle        T
	segmentIndex int
	start        bool
}

func (ep *endPoint[T]) pt(segments []Segment[T]) geom.Point[T] {
	if ep.start {
		return segments[ep.segmentIndex].Start
	}
	return segments[ep.segmentIndex].End
}
