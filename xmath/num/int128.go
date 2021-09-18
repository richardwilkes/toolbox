// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package num

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
)

const (
	signBit        = 0x8000000000000000
	minInt128Float = float64(-170141183460469231731687303715884105728)
	maxInt128Float = float64(170141183460469231731687303715884105727)
)

var (
	// MaxInt128 is the maximum value representable by an Int128.
	MaxInt128 = Int128{hi: 0x7FFFFFFFFFFFFFFF, lo: 0xFFFFFFFFFFFFFFFF}
	// MinInt128 is the minimum value representable by an Int128.
	MinInt128 = Int128{hi: signBit, lo: 0}
)

var (
	minInt128AsAbsUint128 = Uint128{hi: signBit, lo: 0}
	maxInt128AsUint128    = Uint128{hi: 0x7FFFFFFFFFFFFFFF, lo: 0xFFFFFFFFFFFFFFFF}
	maxBigUint128, _      = new(big.Int).SetString("340282366920938463463374607431768211455", 10)
	big1                  = new(big.Int).SetInt64(1)
)

// Int128 represents a signed 128-bit integer.
type Int128 struct {
	hi uint64
	lo uint64
}

// Int128From64 creates an Int128 from an int64 value.
func Int128From64(v int64) Int128 {
	var hi uint64
	if v < 0 {
		hi = math.MaxUint64
	}
	return Int128{hi: hi, lo: uint64(v)}
}

// Int128FromUint64 creates an Int128 from a uint64 value.
func Int128FromUint64(v uint64) Int128 {
	return Int128{lo: v}
}

// Int128FromFloat64 creates an Int128 from a float64 value.
func Int128FromFloat64(f float64) Int128 {
	switch {
	case f == 0 || f != f: // 0 or NaN
		return Int128{}
	case f < 0:
		switch {
		case f >= -float64(math.MaxUint64)-1:
			return Int128{
				hi: math.MaxUint64,
				lo: uint64(f),
			}
		case f >= minInt128Float:
			f = -f
			lo := math.Mod(f, wrapUint64Float)
			return Int128{
				hi: ^uint64(f / wrapUint64Float),
				lo: ^uint64(lo),
			}
		default:
			return MinInt128
		}
	default:
		switch {
		case f <= float64(math.MaxUint64):
			return Int128{lo: uint64(f)}
		case f <= maxInt128Float:
			return Int128{
				hi: uint64(f / wrapUint64Float),
				lo: uint64(math.Mod(f, wrapUint64Float)),
			}
		default:
			return MaxInt128
		}
	}
}

// Int128FromBigInt creates an Int128 from a big.Int.
func Int128FromBigInt(v *big.Int) Int128 {
	var i Uint128
	words := v.Bits()
	switch len(words) {
	case 0:
	case 1:
		i.lo = uint64(words[0])
	case 2:
		if intSize == 64 {
			i.hi = uint64(words[1])
			i.lo = uint64(words[0])
		} else {
			i.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		}
	case 3:
		if intSize == 64 {
			i = MaxUint128
		} else {
			i.hi = uint64(words[2])
			i.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		}
	case 4:
		if intSize == 64 {
			i = MaxUint128
		} else {
			i.hi = (uint64(words[3]) << 32) | (uint64(words[2]))
			i.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		}
	default:
		i = MaxUint128
	}
	if v.Sign() >= 0 {
		if i.LessThan(maxInt128AsUint128) {
			return i.AsInt128()
		}
		return MaxInt128
	}
	if i.LessThan(minInt128AsAbsUint128) {
		return i.AsInt128().Neg()
	}
	return MinInt128
}

// Int128FromString creates an Int128 from a string.
func Int128FromString(s string) (Int128, error) {
	b, err := parseToBigInt(s)
	if err != nil {
		return Int128{}, err
	}
	return Int128FromBigInt(b), nil
}

// Int128FromStringNoCheck creates an Int128 from a string. Unlike Int128FromString, this allows any string as input.
func Int128FromStringNoCheck(s string) Int128 {
	i, _ := Int128FromString(s) //nolint:errcheck // Failure results in 0
	return i
}

// Int128FromComponents creates an Int128 from two uint64 values representing the high and low bits.
func Int128FromComponents(high, low uint64) Int128 {
	return Int128{hi: high, lo: low}
}

// Int128FromRand generates a signed 128-bit random integer.
func Int128FromRand(source RandomSource) Int128 {
	return Int128{hi: source.Uint64(), lo: source.Uint64()}
}

// Components returns the two uint64 values representing the high and low bits.
func (i Int128) Components() (high, low uint64) {
	return i.hi, i.lo
}

// IsZero returns true if the value is 0.
func (i Int128) IsZero() bool {
	return i.hi|i.lo == 0
}

// ToBigInt stores the Int128's value into the specified big.Int.
func (i Int128) ToBigInt(b *big.Int) {
	Uint128(i).ToBigInt(b)
	if !i.IsUint128() {
		b.Xor(b, maxBigUint128).Add(b, big1).Neg(b)
	}
}

// AsBigInt returns the Int128 as a big.Int.
func (i Int128) AsBigInt() *big.Int {
	var b big.Int
	i.ToBigInt(&b)
	return &b
}

// AsBigFloat returns the Int128 as a big.Float.
func (i Int128) AsBigFloat() (b *big.Float) {
	return new(big.Float).SetInt(i.AsBigInt())
}

// AsFloat64 returns the Int128 as a float64.
func (i Int128) AsFloat64() float64 {
	switch {
	case i.hi == 0:
		if i.lo == 0 {
			return 0
		}
		return float64(i.lo)
	case i.hi == math.MaxUint64:
		return -float64((^i.lo) + 1)
	case i.hi&signBit == 0:
		return (float64(i.hi) * maxUint64Float) + float64(i.lo)
	default:
		return (-float64(^i.hi) * maxUint64Float) + -float64(^i.lo)
	}
}

// IsUint128 returns true if this value can be represented as an Uint128 without any loss.
func (i Int128) IsUint128() bool {
	return i.hi&signBit == 0
}

// AsUint128 returns the Int128 as a Uint128.
func (i Int128) AsUint128() Uint128 {
	return Uint128(i)
}

// IsInt64 returns true if this value can be represented as an int64 without any loss.
func (i Int128) IsInt64() bool {
	if i.hi&signBit != 0 {
		return i.hi == math.MaxUint64 && i.lo >= signBit
	}
	return i.hi == 0 && i.lo <= math.MaxInt64
}

// AsInt64 returns the Int128 as an int64.
func (i Int128) AsInt64() int64 {
	if i.hi&signBit != 0 {
		return -int64(^(i.lo - 1))
	}
	return int64(i.lo)
}

// IsUint64 returns true if this value can be represented as a uint64 without any loss.
func (i Int128) IsUint64() bool {
	return i.hi == 0
}

// AsUint64 returns the Int128 as a uint64.
func (i Int128) AsUint64() uint64 {
	return i.lo
}

// Add returns i + n.
func (i Int128) Add(n Int128) Int128 {
	lo, carry := bits.Add64(i.lo, n.lo, 0)
	hi, _ := bits.Add64(i.hi, n.hi, carry)
	return Int128{
		hi: hi,
		lo: lo,
	}
}

// Add64 returns i + n.
func (i Int128) Add64(n int64) Int128 {
	lo, carry := bits.Add64(i.lo, uint64(n), 0)
	if n < 0 {
		carry += math.MaxUint64
	}
	return Int128{
		hi: i.hi + carry,
		lo: lo,
	}
}

// Sub returns i - n.
func (i Int128) Sub(n Int128) Int128 {
	lo, borrow := bits.Sub64(i.lo, n.lo, 0)
	hi, _ := bits.Sub64(i.hi, n.hi, borrow)
	return Int128{
		hi: hi,
		lo: lo,
	}
}

// Sub64 returns i - n.
func (i Int128) Sub64(n int64) Int128 {
	lo, borrow := bits.Sub64(i.lo, uint64(n), 0)
	hi := i.hi - borrow
	if n < 0 {
		hi -= math.MaxUint64
	}
	return Int128{
		hi: hi,
		lo: lo,
	}
}

// Inc returns i + 1.
func (i Int128) Inc() Int128 {
	return Int128(Uint128(i).Inc())
}

// Dec returns i - 1.
func (i Int128) Dec() Int128 {
	return Int128(Uint128(i).Dec())
}

// Sign returns 1 if i > 0, 0 if i == 0, and -1 if i < 0.
func (i Int128) Sign() int {
	switch {
	case i.hi|i.lo == 0:
		return 0
	case i.hi&signBit == 0:
		return 1
	default:
		return -1
	}
}

// Neg returns -i.
func (i Int128) Neg() Int128 {
	switch {
	case i.hi|i.lo == 0 || i == MinInt128:
		return i
	case i.hi&signBit != 0:
		hi := ^i.hi
		lo := ^(i.lo - 1)
		if lo == 0 {
			hi++
		}
		return Int128{hi: hi, lo: lo}
	default:
		hi := ^i.hi
		lo := (^i.lo) + 1
		if lo == 0 {
			hi++
		}
		return Int128{hi: hi, lo: lo}
	}
}

// Abs returns the absolute value of i as an Int128.
func (i Int128) Abs() Int128 {
	if i.hi&signBit != 0 {
		i.hi = ^i.hi
		i.lo = ^(i.lo - 1)
		if i.lo == 0 {
			i.hi++
		}
	}
	return i
}

// AbsUint128 returns the absolute value of i as a Uint128.
func (i Int128) AbsUint128() Uint128 {
	v := Uint128(i)
	if i == MinInt128 {
		return v
	}
	if i.hi&signBit != 0 {
		v.hi = ^i.hi
		v.lo = ^(i.lo - 1)
		if v.lo == 0 {
			v.hi++
		}
	}
	return v
}

// Cmp returns 1 if i > n, 0 if i == n, and -1 if i < n.
func (i Int128) Cmp(n Int128) int {
	switch {
	case i.hi == n.hi && i.lo == n.lo:
		return 0
	case i.hi&signBit == n.hi&signBit:
		if i.hi > n.hi || (i.hi == n.hi && i.lo > n.lo) {
			return 1
		}
	case i.hi&signBit == 0:
		return 1
	}
	return -1
}

// Cmp64 returns 1 if i > n, 0 if i == n, and -1 if i < n.
func (i Int128) Cmp64(n int64) int {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	switch {
	case i.hi == nhi && i.lo == nlo:
		return 0
	case i.hi&signBit == nhi&signBit:
		if i.hi > nhi || (i.hi == nhi && i.lo > nlo) {
			return 1
		}
	case i.hi&signBit == 0:
		return 1
	}
	return -1
}

// GreaterThan returns true if i > n.
func (i Int128) GreaterThan(n Int128) bool {
	switch {
	case i.hi&signBit == n.hi&signBit:
		return i.hi > n.hi || (i.hi == n.hi && i.lo > n.lo)
	case i.hi&signBit == 0:
		return true
	default:
		return false
	}
}

// GreaterThan64 returns true if i > n.
func (i Int128) GreaterThan64(n int64) bool {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	switch {
	case i.hi&signBit == nhi&signBit:
		return i.hi > nhi || (i.hi == nhi && i.lo > nlo)
	case i.hi&signBit == 0:
		return true
	default:
		return false
	}
}

// GreaterOrEqualTo returns true if i >= n.
func (i Int128) GreaterOrEqualTo(n Int128) bool {
	switch {
	case i.hi == n.hi && i.lo == n.lo:
		return true
	case i.hi&signBit == n.hi&signBit:
		return i.hi > n.hi || (i.hi == n.hi && i.lo > n.lo)
	case i.hi&signBit == 0:
		return true
	default:
		return false
	}
}

// GreaterOrEqualTo64 returns true if i >= n.
func (i Int128) GreaterOrEqualTo64(n int64) bool {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	switch {
	case i.hi == nhi && i.lo == nlo:
		return true
	case i.hi&signBit == nhi&signBit:
		return i.hi > nhi || (i.hi == nhi && i.lo > nlo)
	case i.hi&signBit == 0:
		return true
	default:
		return false
	}
}

// Equal returns true if i == n.
func (i Int128) Equal(n Int128) bool {
	return i.hi == n.hi && i.lo == n.lo
}

// Equal64 returns true if i == n.
func (i Int128) Equal64(n int64) bool {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	return i.hi == nhi && i.lo == nlo
}

// LessThan returns true if i < n.
func (i Int128) LessThan(n Int128) bool {
	switch {
	case i.hi&signBit == n.hi&signBit:
		return i.hi < n.hi || (i.hi == n.hi && i.lo < n.lo)
	case i.hi&signBit != 0:
		return true
	default:
		return false
	}
}

// LessThan64 returns true if i < n.
func (i Int128) LessThan64(n int64) bool {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	switch {
	case i.hi&signBit == nhi&signBit:
		return i.hi < nhi || (i.hi == nhi && i.lo < nlo)
	case i.hi&signBit != 0:
		return true
	default:
		return false
	}
}

// LessOrEqualTo returns true if i <= n.
func (i Int128) LessOrEqualTo(n Int128) bool {
	switch {
	case i.hi == n.hi && i.lo == n.lo:
		return true
	case i.hi&signBit == n.hi&signBit:
		return i.hi < n.hi || (i.hi == n.hi && i.lo < n.lo)
	case i.hi&signBit != 0:
		return true
	default:
		return false
	}
}

// LessOrEqualTo64 returns true if i <= n.
func (i Int128) LessOrEqualTo64(n int64) bool {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	switch {
	case i.hi == nhi && i.lo == nlo:
		return true
	case i.hi&signBit == nhi&signBit:
		return i.hi < nhi || (i.hi == nhi && i.lo < nlo)
	case i.hi&signBit != 0:
		return true
	default:
		return false
	}
}

// Mul returns i * n.
func (i Int128) Mul(n Int128) Int128 {
	hi, lo := bits.Mul64(i.lo, n.lo)
	return Int128{
		hi: hi + i.hi*n.lo + i.lo*n.hi,
		lo: lo,
	}
}

// Mul64 returns i * n.
func (i Int128) Mul64(n int64) Int128 {
	return i.Mul(Int128From64(n))
}

// Div returns i / n. If n == 0, a divide by zero panic will occur.
func (i Int128) Div(n Int128) Int128 {
	qSign := 1
	if i.LessThan(Int128{}) {
		qSign = -1
		//goland:noinspection GoAssignmentToReceiver
		i = i.Neg()
	}
	if n.LessThan(Int128{}) {
		qSign = -qSign
		n = n.Neg()
	}
	q := Int128(Uint128(i).Div(Uint128(n)))
	if qSign < 0 {
		q = q.Neg()
	}
	return q
}

// Div64 returns i / n. If n == 0, a divide by zero panic will occur.
func (i Int128) Div64(n int64) Int128 {
	qSign := 1
	if i.LessThan(Int128{}) {
		qSign = -1
		//goland:noinspection GoAssignmentToReceiver
		i = i.Neg()
	}
	if n < 0 {
		qSign = -qSign
		n = -n
	}
	q := Int128(Uint128(i).Div64(uint64(n)))
	if qSign < 0 {
		q = q.Neg()
	}
	return q
}

// DivMod returns both the result of i / n as well i % n. If n == 0, a divide by zero panic will occur.
func (i Int128) DivMod(n Int128) (q, r Int128) {
	qSign := 1
	rSign := 1
	if i.LessThan(Int128{}) {
		qSign = -1
		rSign = -1
		//goland:noinspection GoAssignmentToReceiver
		i = i.Neg()
	}
	if n.LessThan(Int128{}) {
		qSign = -qSign
		n = n.Neg()
	}
	qu, ru := Uint128(i).DivMod(Uint128(n))
	q = Int128(qu)
	r = Int128(ru)
	if qSign < 0 {
		q = q.Neg()
	}
	if rSign < 0 {
		r = r.Neg()
	}
	return q, r
}

// DivMod64 returns both the result of i / n as well i % n. If n == 0, a divide by zero panic will occur.
func (i Int128) DivMod64(n int64) (q, r Int128) {
	var hi uint64
	if n < 0 {
		hi = math.MaxUint64
	}
	return i.DivMod(Int128{hi: hi, lo: uint64(n)})
}

// Mod returns i % n. If n == 0, a divide by zero panic will occur.
func (i Int128) Mod(n Int128) (r Int128) {
	_, r = i.DivMod(n)
	return r
}

// Mod64 returns i % n. If n == 0, a divide by zero panic will occur.
func (i Int128) Mod64(n int64) (r Int128) {
	_, r = i.DivMod64(n)
	return r
}

// String implements fmt.Stringer.
func (i Int128) String() string {
	if i.hi == 0 {
		if i.lo == 0 {
			return "0"
		}
		return strconv.FormatUint(i.lo, 10)
	}
	return i.AsBigInt().String()
}

// Format implements fmt.Formatter.
func (i Int128) Format(s fmt.State, c rune) {
	i.AsBigInt().Format(s, c)
}

// Scan implements fmt.Scanner.
func (i *Int128) Scan(state fmt.ScanState, verb rune) error {
	t, err := state.Token(true, nil)
	if err != nil {
		return errs.Wrap(err)
	}
	v, err := Int128FromString(string(t))
	if err != nil {
		return errs.Wrap(err)
	}
	*i = v
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int128) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int128) UnmarshalText(text []byte) (err error) {
	v, err := Int128FromString(string(text))
	if err != nil {
		return err
	}
	*i = v
	return nil
}

// Float64 implements json.Number. Intentionally always returns an error, as we never want to emit floating point values
// into json for Int128.
func (i Int128) Float64() (float64, error) {
	return 0, errNoFloat64
}

// Int64 implements json.Number.
func (i Int128) Int64() (int64, error) {
	if !i.IsInt64() {
		return 0, errDoesNotFitInInt64
	}
	return i.AsInt64(), nil
}

// MarshalJSON implements json.Marshaler.
func (i Int128) MarshalJSON() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int128) UnmarshalJSON(in []byte) error {
	v, err := Int128FromString(string(in))
	if err != nil {
		return err
	}
	*i = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (i Int128) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (i *Int128) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := Int128FromString(str)
	if err != nil {
		return err
	}
	*i = v
	return nil
}
