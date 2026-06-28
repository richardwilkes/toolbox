// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package icns_test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/ximage/icns"
)

// TestEncodeNRGBASubImageMatchesStandalone ensures the ARGB encoding path reads pixels via the image's bounds rather
// than walking the raw Pix slice. A sub-image has a stride wider than its width and a Pix slice that starts at a
// non-zero offset, so identical pixel content must still produce identical output regardless of the backing layout.
func TestEncodeNRGBASubImageMatchesStandalone(t *testing.T) {
	c := check.New(t)
	for _, size := range []int{16, 32} { // 16 -> ic04, 32 -> ic05, both exercise createARGBData
		fill := func(img *image.NRGBA) {
			b := img.Bounds()
			for y := range size {
				for x := range size {
					img.SetNRGBA(b.Min.X+x, b.Min.Y+y, color.NRGBA{
						R: uint8(x*16 + 1),
						G: uint8(y*16 + 2),
						B: uint8(x + y + 3),
						A: uint8(255 - x),
					})
				}
			}
		}

		standalone := image.NewNRGBA(image.Rect(0, 0, size, size))
		fill(standalone)

		// The same pixels as a sub-image of a larger image, giving a non-tight stride and a non-zero Pix offset.
		parent := image.NewNRGBA(image.Rect(0, 0, size*3, size*3))
		sub, ok := parent.SubImage(image.Rect(size, size, size*2, size*2)).(*image.NRGBA)
		c.True(ok)
		fill(sub)

		// Guard the test itself: the sub-image must actually have a non-tight layout, otherwise it wouldn't exercise
		// the bug being regression-tested.
		c.NotEqual(sub.Bounds().Dx()*4, sub.Stride)

		var standaloneBuf, subBuf bytes.Buffer
		c.NoError(icns.Encode(&standaloneBuf, standalone))
		c.NoError(icns.Encode(&subBuf, sub))

		c.Equal(standaloneBuf.Bytes(), subBuf.Bytes())
	}
}
