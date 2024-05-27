// Copyright (c) 2016-2022 by Richard A. Wilkes. All rights reserved.
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

// See https://en.wikipedia.org/wiki/Apple_Icon_Image_format for information on the .icns file format.

type header struct {
	Magic  [4]uint8
	Length uint32
}

type entryHeader struct {
	IconType [4]uint8
	Length   uint32
}

type iconInfo struct {
	iconType [4]byte
	buffer   []byte
}

// Encode one or more images into an .icns. At least one image must be provided. macOS recommends providing 1024x1024,
// 512x512, 256x256, 128x128, 64x64, 32x32, and 16x16. Note that sizes other than these will not be considered valid.
func Encode(w io.Writer, images ...image.Image) error {
	if len(images) == 0 {
		return errs.New("must supply at least 1 image")
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].Bounds().Dx() > images[j].Bounds().Dx()
	})
	var info []*iconInfo
	for _, img := range images {
		width := img.Bounds().Dx()
		if width != img.Bounds().Dy() {
			return errs.New("image must be square")
		}
		var ii *iconInfo
		var err error
		switch width {
		case 1024:
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '1', '0'}); err != nil {
				return err
			}
			info = append(info, ii)
		case 512:
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '0', '9'}); err != nil {
				return err
			}
			info = append(info, ii)
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '1', '4'}); err != nil {
				return err
			}
			info = append(info, ii)
		case 256:
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '0', '8'}); err != nil {
				return err
			}
			info = append(info, ii)
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '1', '3'}); err != nil {
				return err
			}
			info = append(info, ii)
		case 128:
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '0', '7'}); err != nil {
				return err
			}
			info = append(info, ii)
		case 64:
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '1', '2'}); err != nil {
				return err
			}
			info = append(info, ii)
		case 32:
			if ii, err = createPNGData(img, [4]byte{'i', 'c', '1', '1'}); err != nil {
				return err
			}
			info = append(info, ii)
			if ii, err = createARGBData(img, [4]byte{'i', 'c', '0', '5'}); err != nil {
				return err
			}
			info = append(info, ii)
		case 16:
			if ii, err = createARGBData(img, [4]byte{'i', 'c', '0', '4'}); err != nil {
				return err
			}
			info = append(info, ii)
		default:
			return errs.New("invalid image size")
		}
	}
	totalBytes := 8 + 8*len(info)
	for _, one := range info {
		totalBytes += len(one.buffer)
	}
	if err := binary.Write(w, binary.BigEndian, header{
		Magic:  [4]uint8{'i', 'c', 'n', 's'},
		Length: uint32(totalBytes),
	}); err != nil {
		return errs.Wrap(err)
	}
	for _, one := range info {
		if err := binary.Write(w, binary.BigEndian, entryHeader{
			IconType: one.iconType,
			Length:   8 + uint32(len(one.buffer)),
		}); err != nil {
			return errs.Wrap(err)
		}
		if _, err := w.Write(one.buffer); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func createPNGData(img image.Image, iconType [4]byte) (*iconInfo, error) {
	var buffer bytes.Buffer
	if _, ok := img.(*image.RGBA); !ok {
		bounds := img.Bounds()
		m := image.NewRGBA(bounds)
		draw.Draw(m, bounds, img, bounds.Min, draw.Src)
		img = m
	}
	if err := png.Encode(&buffer, img); err != nil {
		return nil, errs.Wrap(err)
	}
	return &iconInfo{
		iconType: iconType,
		buffer:   buffer.Bytes(),
	}, nil
}

func createARGBData(img image.Image, iconType [4]byte) (*iconInfo, error) {
	var buffer bytes.Buffer
	buffer.Write([]byte{'A', 'R', 'G', 'B'})
	var nrgba *image.NRGBA
	nrgba, ok := img.(*image.NRGBA)
	if !ok {
		bounds := img.Bounds()
		nrgba = image.NewNRGBA(bounds)
		draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	}
	size := len(nrgba.Pix)
	a := make([]byte, size/4)
	r := make([]byte, size/4)
	g := make([]byte, size/4)
	b := make([]byte, size/4)
	for i := 0; i < size; i += 4 {
		j := i / 4
		r[j] = nrgba.Pix[i]
		g[j] = nrgba.Pix[i+1]
		b[j] = nrgba.Pix[i+2]
		a[j] = nrgba.Pix[i+3]
	}
	if err := writeChannel(&buffer, a); err != nil {
		return nil, err
	}
	if err := writeChannel(&buffer, r); err != nil {
		return nil, err
	}
	if err := writeChannel(&buffer, g); err != nil {
		return nil, err
	}
	if err := writeChannel(&buffer, b); err != nil {
		return nil, err
	}
	return &iconInfo{
		iconType: iconType,
		buffer:   buffer.Bytes(),
	}, nil
}

func writeChannel(buffer *bytes.Buffer, data []byte) error {
	size := len(data)
	for i := 0; i < size; i += 128 {
		count := size - i
		if count > 128 {
			count = 128
		}
		if err := buffer.WriteByte(byte(count - 1)); err != nil {
			return errs.Wrap(err)
		}
		if _, err := buffer.Write(data[i : i+count]); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}
