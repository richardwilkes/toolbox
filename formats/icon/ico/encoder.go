// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ico

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

// See https://en.wikipedia.org/wiki/ICO_(file_format) for information on the
// .ico file format.

type header struct {
	Reserved  uint16 //nolint:unused,structcheck // necessary for the format on disk
	ImageType uint16
	Count     uint16
}

type entry struct {
	Width       uint8
	Height      uint8
	Colors      uint8 //nolint:unused,structcheck // necessary for the format on disk
	Reserved    uint8 //nolint:unused,structcheck // necessary for the format on disk
	Planes      uint16
	BitPerPixel uint16
	Size        uint32
	Offset      uint32
}

// Encode one or more images into an .ico. At least one image must be provided
// and no image may have a width or height greater than 256 pixels. Windows
// recommends providing 256x256, 48x48, 32x32, and 16x16 icons.
func Encode(w io.Writer, images ...image.Image) error {
	if len(images) == 0 {
		return errs.New("must supply at least 1 image")
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].Bounds().Dx() > images[j].Bounds().Dx()
	})
	list := make([][]byte, 0, len(images))
	for _, img := range images {
		bounds := img.Bounds()
		if bounds.Dx() > 256 || bounds.Dy() > 256 {
			return errs.New("image too large - .ico has a 256x256 size limit")
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
		list = append(list, buffer.Bytes())
	}
	if err := binary.Write(w, binary.LittleEndian, header{
		ImageType: 1,
		Count:     uint16(len(images)),
	}); err != nil {
		return errs.Wrap(err)
	}
	offset := 6 + uint32(len(images))*16
	for i, img := range images {
		bounds := img.Bounds()
		width := bounds.Dx()
		if width > 255 {
			width = 0
		}
		height := bounds.Dy()
		if height > 255 {
			height = 0
		}
		e := entry{
			Width:       uint8(width),
			Height:      uint8(height),
			Planes:      1,
			BitPerPixel: 32,
			Size:        uint32(len(list[i])),
			Offset:      offset,
		}
		if err := binary.Write(w, binary.LittleEndian, e); err != nil {
			return errs.Wrap(err)
		}
		offset += e.Size
	}
	for _, data := range list {
		if _, err := w.Write(data); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}
