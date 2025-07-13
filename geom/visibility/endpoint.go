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
)

type endPoint struct {
	segmentIndex int
	angle        float32
	start        bool
}

func (ep *endPoint) pt(segments []Segment) geom.Point {
	if ep.start {
		return segments[ep.segmentIndex].Start
	}
	return segments[ep.segmentIndex].End
}
