// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ico_test

import (
	"bytes"
	"image"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/ximage/ico"
)

// TestEncodeDoesNotReorderCallerSlice ensures Encode sorts a copy of the variadic images rather than the caller's
// backing array. The images are supplied in ascending-width order; an in-place sort would permute them into
// descending-width order, so the slice must read back in its original order after the call.
func TestEncodeDoesNotReorderCallerSlice(t *testing.T) {
	c := check.New(t)
	images := []image.Image{
		image.NewNRGBA(image.Rect(0, 0, 16, 16)),
		image.NewNRGBA(image.Rect(0, 0, 32, 32)),
		image.NewNRGBA(image.Rect(0, 0, 64, 64)),
	}
	widthsBefore := widths(images)
	var buf bytes.Buffer
	c.NoError(ico.Encode(&buf, images...))
	c.Equal(widthsBefore, widths(images))
}

func widths(images []image.Image) []int {
	result := make([]int, len(images))
	for i, img := range images {
		result[i] = img.Bounds().Dx()
	}
	return result
}
