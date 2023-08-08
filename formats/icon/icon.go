// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package icon

import (
	"image"

	"golang.org/x/image/draw"
)

// Stack a set of images on top of each other, producing a new image. The first image in the series will be on the
// bottom and the last will be on the top. If the images are of different sizes, the resulting image will be the size of
// the largest image and all other images will be centered within that area.
func Stack(images ...image.Image) image.Image {
	var width, height int
	for _, img := range images {
		bounds := img.Bounds()
		if width < bounds.Dx() {
			width = bounds.Dx()
		}
		if height < bounds.Dy() {
			height = bounds.Dy()
		}
	}
	base := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, img := range images {
		bounds := img.Bounds()
		w := bounds.Dx()
		h := bounds.Dy()
		draw.Copy(base, image.Pt((width-w)/2, (height-h)/2), img, bounds, draw.Over, nil)
	}
	return base
}

// ImageAt provides storage for an image and an origin point.
type ImageAt struct {
	Image  image.Image
	Origin image.Point
}

// StackAt stacks a set of images on top of each other, producing a new image. The first image in the series will be on
// the bottom and the last will be on the top. The resulting image will be the size of the largest area covered based on
// each image's size plus origin. Note that if an origin has a negative value, it will be normalized such that the
// largest negative will become the new origin for the resulting image.
func StackAt(images ...*ImageAt) image.Image {
	var x, y, width, height int
	for _, img := range images {
		bounds := img.Image.Bounds()
		if x > img.Origin.X {
			x = img.Origin.X
		}
		if y > img.Origin.Y {
			y = img.Origin.Y
		}
		w := bounds.Dx() + img.Origin.X
		if width < w {
			width = w
		}
		h := bounds.Dy() + img.Origin.Y
		if height < h {
			height = h
		}
	}
	base := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, img := range images {
		bounds := img.Image.Bounds()
		draw.Copy(base, image.Pt(img.Origin.X-x, img.Origin.Y-y), img.Image, bounds, draw.Over, nil)
	}
	return base
}

// Scale an image.
func Scale(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w == width && h == height {
		return img
	}
	scaled := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), img, bounds, draw.Over, nil)
	return scaled
}

// ScaleTo scales the image to the desired sizes. If an image cannot be scaled exactly to the desired size, it will be
// scaled proportionally and then centered within the available space.
func ScaleTo(img image.Image, sizes []image.Point) []image.Image {
	list := make([]image.Image, 0, len(sizes))
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	for _, size := range sizes {
		w, h := ScaleProportionally(width, height, size.X, size.Y)
		scaled := Scale(img, w, h)
		if w != size.X || h != size.Y {
			scaled = StackAt(
				&ImageAt{
					Image:  image.NewRGBA(image.Rect(0, 0, size.X, size.Y)),
					Origin: image.Point{},
				},
				&ImageAt{
					Image:  scaled,
					Origin: image.Pt((size.X-w)/2, (size.Y-h)/2),
				},
			)
		}
		list = append(list, scaled)
	}
	return list
}

// ScaleProportionally returns the width and height that are closest to the desired values without distorting the size.
func ScaleProportionally(currentWidth, currentHeight, desiredWidth, desiredHeight int) (width, height int) {
	if desiredWidth != currentWidth || desiredHeight != currentHeight {
		scaleX := float64(desiredWidth) / float64(currentWidth)
		scaleY := float64(desiredHeight) / float64(currentHeight)
		scale := min(scaleX, scaleY)
		return int(float64(currentWidth) * scale), int(float64(currentHeight) * scale)
	}
	return currentWidth, currentHeight
}

// ScaleUpProportionally returns the width and height that are closest to the desired values without distorting the
// size, but won't decrease the current values.
func ScaleUpProportionally(currentWidth, currentHeight, desiredWidth, desiredHeight int) (width, height int) {
	if desiredWidth != currentWidth || desiredHeight != currentHeight {
		scaleX := float64(desiredWidth) / float64(currentWidth)
		scaleY := float64(desiredHeight) / float64(currentHeight)
		if scale := min(scaleX, scaleY); scale > 1 {
			return int(float64(currentWidth) * scale), int(float64(currentHeight) * scale)
		}
	}
	return currentWidth, currentHeight
}

// ScaleDownProportionally returns the width and height that are closest to the desired values without distorting the
// size, but won't increase the current values.
func ScaleDownProportionally(currentWidth, currentHeight, desiredWidth, desiredHeight int) (width, height int) {
	if desiredWidth != currentWidth || desiredHeight != currentHeight {
		scaleX := float64(desiredWidth) / float64(currentWidth)
		scaleY := float64(desiredHeight) / float64(currentHeight)
		if scale := min(scaleX, scaleY); scale < 1 {
			return int(float64(currentWidth) * scale), int(float64(currentHeight) * scale)
		}
	}
	return currentWidth, currentHeight
}
