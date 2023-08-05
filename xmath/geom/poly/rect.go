// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

// Rect creates a new polygon in the shape of a rectangle.
func Rect[T constraints.Float](bounds geom.Rect[T]) Polygon[T] {
	return Polygon[T]{Contour[T]{
		bounds.Point,
		geom.Point[T]{X: bounds.X, Y: bounds.Bottom() - 1},
		geom.Point[T]{X: bounds.Right() - 1, Y: bounds.Bottom() - 1},
		geom.Point[T]{X: bounds.Right() - 1, Y: bounds.Y},
	}}
}
