// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

type vertexType int

const (
	_ vertexType = iota
	externalMaximum
	externalLeftIntermediate
	_ // topEdge, not used
	externalRightIntermediate
	rightEdge
	internalMaximumAndMinimum
	internalMinimum
	externalMinimum
	externalMaximumAndMinimum
	leftEdge
	internalLeftIntermediate
	_ // bottomEdge, not used
	internalRightIntermediate
	internalMaximum
	_ // non-intersection, not used
)

func calcVertexType(tr, tl, br, bl bool) vertexType {
	var vt vertexType
	if tr {
		vt = externalMaximum
	}
	if tl {
		vt |= externalLeftIntermediate
	}
	if br {
		vt |= externalRightIntermediate
	}
	if bl {
		vt |= externalMinimum
	}
	return vt
}
