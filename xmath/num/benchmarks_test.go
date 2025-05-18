// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package num_test

import (
	"math/big"
	"testing"

	"github.com/richardwilkes/toolbox/xmath/num"
)

const (
	leftStr       = "2307687492367321180488"
	rightStr      = "8022819849149681007238941328161"
	shortRightStr = "8491817238132811"
)

var (
	benchInt128  num.Int128
	benchUint128 num.Uint128
	benchBigInt  *big.Int
)

func BenchmarkInt128Add(b *testing.B) {
	left := num.Int128FromStringNoCheck(leftStr)
	right := num.Int128FromStringNoCheck(rightStr)
	var dest num.Int128
	for b.Loop() {
		dest = left.Add(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Add(b *testing.B) {
	left := num.Uint128FromStringNoCheck(leftStr)
	right := num.Uint128FromStringNoCheck(rightStr)
	var dest num.Uint128
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
	left := num.Int128FromStringNoCheck(leftStr)
	right := num.Int128FromStringNoCheck(rightStr)
	var dest num.Int128
	for b.Loop() {
		dest = left.Sub(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Sub(b *testing.B) {
	left := num.Uint128FromStringNoCheck(leftStr)
	right := num.Uint128FromStringNoCheck(rightStr)
	var dest num.Uint128
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
	left := num.Int128FromStringNoCheck(leftStr)
	right := num.Int128FromStringNoCheck(shortRightStr)
	var dest num.Int128
	for b.Loop() {
		dest = left.Mul(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Mul(b *testing.B) {
	left := num.Uint128FromStringNoCheck(leftStr)
	right := num.Uint128FromStringNoCheck(shortRightStr)
	var dest num.Uint128
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
	left := num.Int128FromStringNoCheck(leftStr)
	right := num.Int128FromStringNoCheck(shortRightStr)
	var dest num.Int128
	for b.Loop() {
		dest = left.Div(right)
	}
	benchInt128 = dest
}

func BenchmarkUint128Div(b *testing.B) {
	left := num.Uint128FromStringNoCheck(leftStr)
	right := num.Uint128FromStringNoCheck(shortRightStr)
	var dest num.Uint128
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
