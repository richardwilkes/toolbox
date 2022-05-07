// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed

// Dx is the constraint for allowed decimal place types. These methods must return a static value and be usable with the
// zero value of the type.
type Dx interface {
	Places() int
	Multiplier() int64
}

// D1 is used for fixed-point types with 1 decimal place.
type D1 int64

// Places implements Dx.
func (d D1) Places() int {
	return 1
}

// Multiplier implements Dx.
func (d D1) Multiplier() int64 {
	return 10
}

// D2 is used for fixed-point types with 2 decimal places.
type D2 int64

// Places implements Dx.
func (d D2) Places() int {
	return 2
}

// Multiplier implements Dx.
func (d D2) Multiplier() int64 {
	return 100
}

// D3 is used for fixed-point types with 3 decimal places.
type D3 int64

// Places implements Dx.
func (d D3) Places() int {
	return 3
}

// Multiplier implements Dx.
func (d D3) Multiplier() int64 {
	return 1_000
}

// D4 is used for fixed-point types with 4 decimal places.
type D4 int64

// Places implements Dx.
func (d D4) Places() int {
	return 4
}

// Multiplier implements Dx.
func (d D4) Multiplier() int64 {
	return 10_000
}

// D5 is used for fixed-point types with 5 decimal places.
type D5 int64

// Places implements Dx.
func (d D5) Places() int {
	return 5
}

// Multiplier implements Dx.
func (d D5) Multiplier() int64 {
	return 100_000
}

// D6 is used for fixed-point types with 6 decimal places.
type D6 int64

// Places implements Dx.
func (d D6) Places() int {
	return 6
}

// Multiplier implements Dx.
func (d D6) Multiplier() int64 {
	return 1_000_000
}

// D7 is used for fixed-point types with 7 decimal places.
type D7 int64

// Places implements Dx.
func (d D7) Places() int {
	return 7
}

// Multiplier implements Dx.
func (d D7) Multiplier() int64 {
	return 10_000_000
}

// D8 is used for fixed-point types with 8 decimal places.
type D8 int64

// Places implements Dx.
func (d D8) Places() int {
	return 8
}

// Multiplier implements Dx.
func (d D8) Multiplier() int64 {
	return 100_000_000
}

// D9 is used for fixed-point types with 9 decimal places.
type D9 int64

// Places implements Dx.
func (d D9) Places() int {
	return 9
}

// Multiplier implements Dx.
func (d D9) Multiplier() int64 {
	return 1_000_000_000
}

// D10 is used for fixed-point types with 10 decimal places.
type D10 int64

// Places implements Dx.
func (d D10) Places() int {
	return 10
}

// Multiplier implements Dx.
func (d D10) Multiplier() int64 {
	return 10_000_000_000
}

// D11 is used for fixed-point types with 11 decimal places.
type D11 int64

// Places implements Dx.
func (d D11) Places() int {
	return 11
}

// Multiplier implements Dx.
func (d D11) Multiplier() int64 {
	return 100_000_000_000
}

// D12 is used for fixed-point types with 12 decimal places.
type D12 int64

// Places implements Dx.
func (d D12) Places() int {
	return 12
}

// Multiplier implements Dx.
func (d D12) Multiplier() int64 {
	return 1_000_000_000_000
}

// D13 is used for fixed-point types with 13 decimal places.
type D13 int64

// Places implements Dx.
func (d D13) Places() int {
	return 13
}

// Multiplier implements Dx.
func (d D13) Multiplier() int64 {
	return 10_000_000_000_000
}

// D14 is used for fixed-point types with 14 decimal places.
type D14 int64

// Places implements Dx.
func (d D14) Places() int {
	return 14
}

// Multiplier implements Dx.
func (d D14) Multiplier() int64 {
	return 100_000_000_000_000
}

// D15 is used for fixed-point types with 15 decimal places.
type D15 int64

// Places implements Dx.
func (d D15) Places() int {
	return 15
}

// Multiplier implements Dx.
func (d D15) Multiplier() int64 {
	return 1_000_000_000_000_000
}

// D16 is used for fixed-point types with 16 decimal places.
type D16 int64

// Places implements Dx.
func (d D16) Places() int {
	return 16
}

// Multiplier implements Dx.
func (d D16) Multiplier() int64 {
	return 10_000_000_000_000_000
}
