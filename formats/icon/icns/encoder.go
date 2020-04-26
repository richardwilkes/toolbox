// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package icns

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"image/png"
	"io"
	"sort"

	"github.com/richardwilkes/toolbox/errs"
)

// See https://en.wikipedia.org/wiki/Apple_Icon_Image_format for information
// on the .icns file format.

type header struct {
	Magic  [4]uint8
	Length uint32
}

type entryHeader struct {
	IconType [4]uint8
	Length   uint32
}

// Encode one or more images into an .icns. At least one image must be
// provided. macOS recommends providing 1024x1024, 512x512, 256x256, 128x128,
// 64x64, 32x32, and 16x16. Note that sizes other than these will not be
// considered valid.
func Encode(w io.Writer, images ...image.Image) error {
	if len(images) == 0 {
		return errs.New("must supply at least 1 image")
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].Bounds().Dx() > images[j].Bounds().Dx()
	})
	list := make([][]byte, 0, len(images))
	var totalBytes uint32
	for _, img := range images {
		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()
		if width != height || !validSize(width) || !validSize(height) {
			return errs.New("invalid image size")
		}
		if _, ok := img.(*image.RGBA); !ok {
			m := image.NewRGBA(bounds)
			draw.Draw(m, bounds, img, bounds.Min, draw.Src)
			img = m
		}
		var buffer bytes.Buffer
		if err := png.Encode(&buffer, img); err != nil {
			return errs.Wrap(err)
		}
		totalBytes += uint32(buffer.Len())
		list = append(list, buffer.Bytes())
	}
	if err := binary.Write(w, binary.BigEndian, header{
		Magic:  [4]uint8{'i', 'c', 'n', 's'},
		Length: 8 + 8*uint32(len(images)) + totalBytes,
	}); err != nil {
		return errs.Wrap(err)
	}
	for i, img := range images {
		if err := binary.Write(w, binary.BigEndian, entryHeader{
			IconType: iconTypeForSize(img.Bounds().Dx()),
			Length:   8 + uint32(len(list[i])),
		}); err != nil {
			return errs.Wrap(err)
		}
		if _, err := w.Write(list[i]); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func validSize(size int) bool {
	return size == 1024 || size == 512 || size == 256 || size == 128 || size == 64 || size == 32 || size == 16
}

func iconTypeForSize(size int) [4]uint8 {
	switch size {
	case 1024:
		return [4]uint8{'i', 'c', '1', '0'}
	case 512:
		return [4]uint8{'i', 'c', '0', '9'}
	case 256:
		return [4]uint8{'i', 'c', '0', '8'}
	case 128:
		return [4]uint8{'i', 'c', '0', '7'}
	case 64:
		return [4]uint8{'i', 'c', 'p', '6'}
	case 32:
		return [4]uint8{'i', 'c', 'p', '5'}
	case 16:
		return [4]uint8{'i', 'c', 'p', '4'}
	default:
		return [4]uint8{} // Can't actually happen
	}
}
