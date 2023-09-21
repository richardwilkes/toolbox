// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom

import (
	"math"
)

const (
	orientEpsilon  = 1.1102230246251565e-16
	ccwErrBoundA   = (3 + 16*orientEpsilon) * orientEpsilon
	ccwErrBoundB   = (2 + 12*orientEpsilon) * orientEpsilon
	ccwErrBoundC   = (9 + 64*orientEpsilon) * orientEpsilon * orientEpsilon
	resultErrBound = (3 + 8*orientEpsilon) * orientEpsilon
	splitter       = 134217729
)

// OrientFast is a simple, non-robust, version of Orient().
//
// Returns 1 if the points 'a', 'b', and 'c' occur in counterclockwise order ('c' lies to the left of the directed line
// defined by points 'a' and 'b'). Returns -1 if they occur in clockwise order ('c' lies to the right of the directed
// line defined by points 'a' and 'b'). Returns 0 if they are collinear.
//
// Due to precision issues that arise for float32, only the float64 version of this function is provided. Use the helper
// functions Pt32to64() and Pt64to32() to convert as needed.
//
// Based on the Javascript code found here: https://github.com/mourner/robust-predicates
func OrientFast(a, b, c Point[float64]) int {
	return normalizeOrientResult((a.Y-c.Y)*(b.X-c.X) - (a.X-c.X)*(b.Y-c.Y))
}

// Orient returns 1 if the points 'a', 'b', and 'c' occur in counterclockwise order ('c' lies to the left of the
// directed line defined by points 'a' and 'b'). Returns -1 if they occur in clockwise order ('c' lies to the right of
// the directed line defined by points 'a' and 'b'). Returns 0 if they are collinear.
//
// Due to precision issues that arise for float32, only the float64 version of this function is provided. Use the helper
// functions Pt32to64() and Pt64to32() to convert as needed.
//
// Based on the Javascript code found here: https://github.com/mourner/robust-predicates
func Orient(a, b, c Point[float64]) int {
	detLeft := (a.Y - c.Y) * (b.X - c.X)
	detRight := (a.X - c.X) * (b.Y - c.Y)
	det := detLeft - detRight
	detSum := math.Abs(detLeft + detRight)
	if math.Abs(det) >= ccwErrBoundA*detSum {
		return normalizeOrientResult(det)
	}
	acx := a.X - c.X
	bcx := b.X - c.X
	acy := a.Y - c.Y
	bcy := b.Y - c.Y

	s1 := acx * bcy
	c1 := splitter * acx
	ahi := c1 - (c1 - acx)
	alo := acx - ahi
	c1 = splitter * bcy
	bhi := c1 - (c1 - bcy)
	blo := bcy - bhi
	s0 := alo*blo - (s1 - ahi*bhi - alo*bhi - ahi*blo)
	t1 := acy * bcx
	c1 = splitter * acy
	ahi = c1 - (c1 - acy)
	alo = acy - ahi
	c1 = splitter * bcx
	bhi = c1 - (c1 - bcx)
	blo = bcx - bhi
	t0 := alo*blo - (t1 - ahi*bhi - alo*bhi - ahi*blo)
	i := s0 - t0
	bvirt := s0 - i
	var B [4]float64
	B[0] = s0 - (i + bvirt) + (bvirt - t0)
	j := s1 + i
	bvirt = j - s1
	z := s1 - (j - bvirt) + (i - bvirt)
	i = z - t1
	bvirt = z - i
	B[1] = z - (i + bvirt) + (bvirt - t1)
	u3 := j + i
	bvirt = u3 - j
	B[2] = j - (u3 - bvirt) + (i - bvirt)
	B[3] = u3

	det = B[0] + B[1] + B[2] + B[3]
	errbound := ccwErrBoundB * detSum
	if det >= errbound || -det >= errbound {
		return normalizeOrientResult(-det)
	}

	bvirt = a.X - acx
	acxtail := a.X - (acx + bvirt) + (bvirt - c.X)
	bvirt = b.X - bcx
	bcxtail := b.X - (bcx + bvirt) + (bvirt - c.X)
	bvirt = a.Y - acy
	acytail := a.Y - (acy + bvirt) + (bvirt - c.Y)
	bvirt = b.Y - bcy
	bcytail := b.Y - (bcy + bvirt) + (bvirt - c.Y)

	if acxtail == 0 && acytail == 0 && bcxtail == 0 && bcytail == 0 {
		return normalizeOrientResult(-det)
	}

	errbound = ccwErrBoundC*detSum + resultErrBound*math.Abs(det)
	det += (acx*bcytail + bcy*acxtail) - (acy*bcxtail + bcx*acytail)
	if det >= errbound || -det >= errbound {
		return normalizeOrientResult(-det)
	}

	s1 = acxtail * bcy
	c1 = splitter * acxtail
	ahi = c1 - (c1 - acxtail)
	alo = acxtail - ahi
	c1 = splitter * bcy
	bhi = c1 - (c1 - bcy)
	blo = bcy - bhi
	s0 = alo*blo - (s1 - ahi*bhi - alo*bhi - ahi*blo)
	t1 = acytail * bcx
	c1 = splitter * acytail
	ahi = c1 - (c1 - acytail)
	alo = acytail - ahi
	c1 = splitter * bcx
	bhi = c1 - (c1 - bcx)
	blo = bcx - bhi
	t0 = alo*blo - (t1 - ahi*bhi - alo*bhi - ahi*blo)
	i = s0 - t0
	bvirt = s0 - i
	var u [4]float64
	u[0] = s0 - (i + bvirt) + (bvirt - t0)
	j = s1 + i
	bvirt = j - s1
	z = s1 - (j - bvirt) + (i - bvirt)
	i = z - t1
	bvirt = z - i
	u[1] = z - (i + bvirt) + (bvirt - t1)
	u3 = j + i
	bvirt = u3 - j
	u[2] = j - (u3 - bvirt) + (i - bvirt)
	u[3] = u3
	var C1 [8]float64
	C1len := orientSum(4, B[:], 4, u[:], C1[:])

	s1 = acx * bcytail
	c1 = splitter * acx
	ahi = c1 - (c1 - acx)
	alo = acx - ahi
	c1 = splitter * bcytail
	bhi = c1 - (c1 - bcytail)
	blo = bcytail - bhi
	s0 = alo*blo - (s1 - ahi*bhi - alo*bhi - ahi*blo)
	t1 = acy * bcxtail
	c1 = splitter * acy
	ahi = c1 - (c1 - acy)
	alo = acy - ahi
	c1 = splitter * bcxtail
	bhi = c1 - (c1 - bcxtail)
	blo = bcxtail - bhi
	t0 = alo*blo - (t1 - ahi*bhi - alo*bhi - ahi*blo)
	i = s0 - t0
	bvirt = s0 - i
	u[0] = s0 - (i + bvirt) + (bvirt - t0)
	j = s1 + i
	bvirt = j - s1
	z = s1 - (j - bvirt) + (i - bvirt)
	i = z - t1
	bvirt = z - i
	u[1] = z - (i + bvirt) + (bvirt - t1)
	u3 = j + i
	bvirt = u3 - j
	u[2] = j - (u3 - bvirt) + (i - bvirt)
	u[3] = u3
	var C2 [12]float64
	C2len := orientSum(C1len, C1[:], 4, u[:], C2[:])

	s1 = acxtail * bcytail
	c1 = splitter * acxtail
	ahi = c1 - (c1 - acxtail)
	alo = acxtail - ahi
	c1 = splitter * bcytail
	bhi = c1 - (c1 - bcytail)
	blo = bcytail - bhi
	s0 = alo*blo - (s1 - ahi*bhi - alo*bhi - ahi*blo)
	t1 = acytail * bcxtail
	c1 = splitter * acytail
	ahi = c1 - (c1 - acytail)
	alo = acytail - ahi
	c1 = splitter * bcxtail
	bhi = c1 - (c1 - bcxtail)
	blo = bcxtail - bhi
	t0 = alo*blo - (t1 - ahi*bhi - alo*bhi - ahi*blo)
	i = s0 - t0
	bvirt = s0 - i
	u[0] = s0 - (i + bvirt) + (bvirt - t0)
	j = s1 + i
	bvirt = j - s1
	z = s1 - (j - bvirt) + (i - bvirt)
	i = z - t1
	bvirt = z - i
	u[1] = z - (i + bvirt) + (bvirt - t1)
	u3 = j + i
	bvirt = u3 - j
	u[2] = j - (u3 - bvirt) + (i - bvirt)
	u[3] = u3
	var D [16]float64
	Dlen := orientSum(C2len, C2[:], 4, u[:], D[:])

	return normalizeOrientResult(-D[Dlen-1])
}

func orientSum(elen int, e []float64, flen int, f, h []float64) int {
	enow := e[0]
	fnow := f[0]
	eindex := 0
	findex := 0
	var Q float64
	if (fnow > enow) == (fnow > -enow) {
		Q = enow
		eindex++
		enow = e[eindex]
	} else {
		Q = fnow
		findex++
		fnow = f[findex]
	}
	hindex := 0
	var Qnew float64
	var hh float64
	if eindex < elen && findex < flen {
		if (fnow > enow) == (fnow > -enow) {
			Qnew = enow + Q
			hh = Q - (Qnew - enow)
			enow = e[eindex]
			eindex++
		} else {
			Qnew = fnow + Q
			hh = Q - (Qnew - fnow)
			fnow = f[findex]
			findex++
		}
		Q = Qnew
		if hh != 0 {
			h[hindex] = hh
			hindex++
		}
		for eindex < elen && findex < flen {
			if (fnow > enow) == (fnow > -enow) {
				Qnew = Q + enow
				bvirt := Qnew - Q
				hh = Q - (Qnew - bvirt) + (enow - bvirt)
				enow = e[eindex]
				eindex++
			} else {
				Qnew = Q + fnow
				bvirt := Qnew - Q
				hh = Q - (Qnew - bvirt) + (fnow - bvirt)
				fnow = f[findex]
				findex++
			}
			Q = Qnew
			if hh != 0 {
				h[hindex] = hh
				hindex++
			}
		}
	}
	for eindex < elen {
		Qnew = Q + enow
		bvirt := Qnew - Q
		hh = Q - (Qnew - bvirt) + (enow - bvirt)
		enow = e[eindex]
		eindex++
		Q = Qnew
		if hh != 0 {
			h[hindex] = hh
			hindex++
		}
	}
	for findex < flen {
		Qnew = Q + fnow
		bvirt := Qnew - Q
		hh = Q - (Qnew - bvirt) + (fnow - bvirt)
		fnow = f[findex]
		findex++
		Q = Qnew
		if hh != 0 {
			h[hindex] = hh
			hindex++
		}
	}
	if Q != 0 || hindex == 0 {
		h[hindex] = Q
		hindex++
	}
	return hindex
}

func normalizeOrientResult(result float64) int {
	if result == 0 {
		return 0
	}
	if result < 0 {
		return -1
	}
	return 1
}
