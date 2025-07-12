// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package num128_test

import (
	"math/big"
	"testing"

	"github.com/richardwilkes/toolbox/v2/num128"
)

const (
	leftStr       = "2307687492367321180488"
	rightStr      = "8022819849149681007238941328161"
	shortRightStr = "8491817238132811"
)

var (
	benchInt128  num128.Int
	benchUint128 num128.Uint
	benchBigInt  *big.Int
)

func BenchmarkInt128Add(b *testing.B) {
	left := num128.IntFromStringNoCheck(leftStr)
	right := num128.IntFromStringNoCheck(rightStr)
	var dest num128.Int
	for b.Loop() {
		dest = left.Add(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Add(b *testing.B) {
	left := num128.UintFromStringNoCheck(leftStr)
	right := num128.UintFromStringNoCheck(rightStr)
	var dest num128.Uint
	for b.Loop() {
		dest = left.Add(right)
	}
	benchUint128 = dest
}

func BenchmarkBigIntAdd(b *testing.B) {
	left, _ := new(big.Int).SetString(leftStr, 0)
	right, _ := new(big.Int).SetString(rightStr, 0)
	dest := new(big.Int)
	for b.Loop() {
		dest.Add(left, right)
	}
	benchBigInt = dest
}

func BenchmarkInt128Sub(b *testing.B) {
	left := num128.IntFromStringNoCheck(leftStr)
	right := num128.IntFromStringNoCheck(rightStr)
	var dest num128.Int
	for b.Loop() {
		dest = left.Sub(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Sub(b *testing.B) {
	left := num128.UintFromStringNoCheck(leftStr)
	right := num128.UintFromStringNoCheck(rightStr)
	var dest num128.Uint
	for b.Loop() {
		dest = left.Sub(right)
	}
	benchUint128 = dest
}

func BenchmarkBigIntSub(b *testing.B) {
	left, _ := new(big.Int).SetString(leftStr, 0)
	right, _ := new(big.Int).SetString(rightStr, 0)
	dest := new(big.Int)
	for b.Loop() {
		dest.Sub(left, right)
	}
	benchBigInt = dest
}

func BenchmarkInt128Mul(b *testing.B) {
	left := num128.IntFromStringNoCheck(leftStr)
	right := num128.IntFromStringNoCheck(shortRightStr)
	var dest num128.Int
	for b.Loop() {
		dest = left.Mul(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Mul(b *testing.B) {
	left := num128.UintFromStringNoCheck(leftStr)
	right := num128.UintFromStringNoCheck(shortRightStr)
	var dest num128.Uint
	for b.Loop() {
		dest = left.Mul(right)
	}
	benchUint128 = dest
}

func BenchmarkBigIntMul(b *testing.B) {
	left, _ := new(big.Int).SetString(leftStr, 0)
	right, _ := new(big.Int).SetString(shortRightStr, 0)
	dest := new(big.Int)
	for b.Loop() {
		dest.Mul(left, right)
	}
	benchBigInt = dest
}

func BenchmarkInt128Div(b *testing.B) {
	left := num128.IntFromStringNoCheck(leftStr)
	right := num128.IntFromStringNoCheck(shortRightStr)
	var dest num128.Int
	for b.Loop() {
		dest = left.Div(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Div(b *testing.B) {
	left := num128.UintFromStringNoCheck(leftStr)
	right := num128.UintFromStringNoCheck(shortRightStr)
	var dest num128.Uint
	for b.Loop() {
		dest = left.Div(right)
	}
	benchUint128 = dest
}

func BenchmarkBigIntDiv(b *testing.B) {
	left, _ := new(big.Int).SetString(leftStr, 0)
	right, _ := new(big.Int).SetString(shortRightStr, 0)
	dest := new(big.Int)
	for b.Loop() {
		dest.Div(left, right)
	}
	benchBigInt = dest
}
