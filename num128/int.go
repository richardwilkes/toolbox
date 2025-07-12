// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package num128

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"

	"github.com/richardwilkes/toolbox/v2/errs"
)

const (
	signBit     = 0x8000000000000000
	minIntFloat = float64(-170141183460469231731687303715884105728)
	maxIntFloat = float64(170141183460469231731687303715884105727)
)

var (
	// MaxInt is the maximum value representable by an Int.
	MaxInt = Int{hi: 0x7FFFFFFFFFFFFFFF, lo: 0xFFFFFFFFFFFFFFFF}
	// MinInt is the minimum value representable by an Int.
	MinInt = Int{hi: signBit, lo: 0}
)

var (
	minIntAsAbsUint = Uint{hi: signBit, lo: 0}
	maxIntAsUint    = Uint{hi: 0x7FFFFFFFFFFFFFFF, lo: 0xFFFFFFFFFFFFFFFF}
	maxBigUint, _   = new(big.Int).SetString("340282366920938463463374607431768211455", 10)
	big1            = new(big.Int).SetInt64(1)
)

// Int represents a signed 128-bit integer.
type Int struct {
	hi uint64
	lo uint64
}

// IntFrom64 creates an Int from an int64 value.
func IntFrom64(v int64) Int {
	var hi uint64
	if v < 0 {
		hi = math.MaxUint64
	}
	return Int{hi: hi, lo: uint64(v)}
}

// IntFromUint64 creates an Int from a uint64 value.
func IntFromUint64(v uint64) Int {
	return Int{lo: v}
}

// IntFromFloat64 creates an Int from a float64 value.
func IntFromFloat64(f float64) Int {
	switch {
	case f == 0 || f != f: // 0 or NaN
		return Int{}
	case f < 0:
		switch {
		case f >= -float64(math.MaxUint64)-1:
			return Int{
				hi: math.MaxUint64,
				lo: uint64(f),
			}
		case f >= minIntFloat:
			f = -f
			lo := math.Mod(f, wrapUint64Float)
			return Int{
				hi: ^uint64(f / wrapUint64Float),
				lo: ^uint64(lo),
			}
		default:
			return MinInt
		}
	default:
		switch {
		case f <= float64(math.MaxUint64):
			return Int{lo: uint64(f)}
		case f <= maxIntFloat:
			return Int{
				hi: uint64(f / wrapUint64Float),
				lo: uint64(math.Mod(f, wrapUint64Float)),
			}
		default:
			return MaxInt
		}
	}
}

// IntFromBigInt creates an Int from a big.Int.
func IntFromBigInt(v *big.Int) Int {
	var i Uint
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
			i = MaxUint
		} else {
			i.hi = uint64(words[2])
			i.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		}
	case 4:
		if intSize == 64 {
			i = MaxUint
		} else {
			i.hi = (uint64(words[3]) << 32) | (uint64(words[2]))
			i.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		}
	default:
		i = MaxUint
	}
	if v.Sign() >= 0 {
		if i.LessThan(maxIntAsUint) {
			return i.AsInt()
		}
		return MaxInt
	}
	if i.LessThan(minIntAsAbsUint) {
		return i.AsInt().Neg()
	}
	return MinInt
}

// IntFromString creates an Int from a string.
func IntFromString(s string) (Int, error) {
	b, err := parseToBigInt(s)
	if err != nil {
		return Int{}, err
	}
	return IntFromBigInt(b), nil
}

// IntFromStringNoCheck creates an Int from a string. Unlike IntFromString, this allows any string as input.
func IntFromStringNoCheck(s string) Int {
	i, _ := IntFromString(s) //nolint:errcheck // Failure results in 0
	return i
}

// IntFromComponents creates an Int from two uint64 values representing the high and low bits.
func IntFromComponents(high, low uint64) Int {
	return Int{hi: high, lo: low}
}

// IntFromRand generates a signed 128-bit random integer.
func IntFromRand(source RandomSource) Int {
	return Int{hi: source.Uint64(), lo: source.Uint64()}
}

// Components returns the two uint64 values representing the high and low bits.
func (i Int) Components() (high, low uint64) {
	return i.hi, i.lo
}

// IsZero returns true if the value is 0.
func (i Int) IsZero() bool {
	return i.hi|i.lo == 0
}

// ToBigInt stores the Int's value into the specified big.Int.
func (i Int) ToBigInt(b *big.Int) {
	Uint(i).ToBigInt(b)
	if !i.IsUint() {
		b.Xor(b, maxBigUint).Add(b, big1).Neg(b)
	}
}

// AsBigInt returns the Int as a big.Int.
func (i Int) AsBigInt() *big.Int {
	var b big.Int
	i.ToBigInt(&b)
	return &b
}

// AsBigFloat returns the Int as a big.Float.
func (i Int) AsBigFloat() (b *big.Float) {
	return new(big.Float).SetInt(i.AsBigInt())
}

// AsFloat64 returns the Int as a float64.
func (i Int) AsFloat64() float64 {
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

// IsUint returns true if this value can be represented as an Uint without any loss.
func (i Int) IsUint() bool {
	return i.hi&signBit == 0
}

// AsUint returns the Int as a Uint.
func (i Int) AsUint() Uint {
	return Uint(i)
}

// IsInt64 returns true if this value can be represented as an int64 without any loss.
func (i Int) IsInt64() bool {
	if i.hi&signBit != 0 {
		return i.hi == math.MaxUint64 && i.lo >= signBit
	}
	return i.hi == 0 && i.lo <= math.MaxInt64
}

// AsInt64 returns the Int as an int64.
func (i Int) AsInt64() int64 {
	if i.hi&signBit != 0 {
		return -int64(^(i.lo - 1))
	}
	return int64(i.lo)
}

// IsUint64 returns true if this value can be represented as a uint64 without any loss.
func (i Int) IsUint64() bool {
	return i.hi == 0
}

// AsUint64 returns the Int as a uint64.
func (i Int) AsUint64() uint64 {
	return i.lo
}

// Add returns i + n.
func (i Int) Add(n Int) Int {
	lo, carry := bits.Add64(i.lo, n.lo, 0)
	hi, _ := bits.Add64(i.hi, n.hi, carry)
	return Int{
		hi: hi,
		lo: lo,
	}
}

// Add64 returns i + n.
func (i Int) Add64(n int64) Int {
	lo, carry := bits.Add64(i.lo, uint64(n), 0)
	if n < 0 {
		carry += math.MaxUint64
	}
	return Int{
		hi: i.hi + carry,
		lo: lo,
	}
}

// Sub returns i - n.
func (i Int) Sub(n Int) Int {
	lo, borrow := bits.Sub64(i.lo, n.lo, 0)
	hi, _ := bits.Sub64(i.hi, n.hi, borrow)
	return Int{
		hi: hi,
		lo: lo,
	}
}

// Sub64 returns i - n.
func (i Int) Sub64(n int64) Int {
	lo, borrow := bits.Sub64(i.lo, uint64(n), 0)
	hi := i.hi - borrow
	if n < 0 {
		hi -= math.MaxUint64
	}
	return Int{
		hi: hi,
		lo: lo,
	}
}

// Inc returns i + 1.
func (i Int) Inc() Int {
	return Int(Uint(i).Inc())
}

// Dec returns i - 1.
func (i Int) Dec() Int {
	return Int(Uint(i).Dec())
}

// Sign returns 1 if i > 0, 0 if i == 0, and -1 if i < 0.
func (i Int) Sign() int {
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
func (i Int) Neg() Int {
	switch {
	case i.hi|i.lo == 0 || i == MinInt:
		return i
	case i.hi&signBit != 0:
		hi := ^i.hi
		lo := ^(i.lo - 1)
		if lo == 0 {
			hi++
		}
		return Int{hi: hi, lo: lo}
	default:
		hi := ^i.hi
		lo := (^i.lo) + 1
		if lo == 0 {
			hi++
		}
		return Int{hi: hi, lo: lo}
	}
}

// Abs returns the absolute value of i as an Int.
func (i Int) Abs() Int {
	if i.hi&signBit != 0 {
		i.hi = ^i.hi
		i.lo = ^(i.lo - 1)
		if i.lo == 0 {
			i.hi++
		}
	}
	return i
}

// AbsUint returns the absolute value of i as a Uint.
func (i Int) AbsUint() Uint {
	v := Uint(i)
	if i == MinInt {
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
func (i Int) Cmp(n Int) int {
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
func (i Int) Cmp64(n int64) int {
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
func (i Int) GreaterThan(n Int) bool {
	switch i.hi & signBit {
	case n.hi & signBit:
		return i.hi > n.hi || (i.hi == n.hi && i.lo > n.lo)
	case 0:
		return true
	default:
		return false
	}
}

// GreaterThan64 returns true if i > n.
func (i Int) GreaterThan64(n int64) bool {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	switch i.hi & signBit {
	case nhi & signBit:
		return i.hi > nhi || (i.hi == nhi && i.lo > nlo)
	case 0:
		return true
	default:
		return false
	}
}

// GreaterThanOrEqual returns true if i >= n.
func (i Int) GreaterThanOrEqual(n Int) bool {
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

// GreaterThanOrEqual64 returns true if i >= n.
func (i Int) GreaterThanOrEqual64(n int64) bool {
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
func (i Int) Equal(n Int) bool {
	return i.hi == n.hi && i.lo == n.lo
}

// Equal64 returns true if i == n.
func (i Int) Equal64(n int64) bool {
	var nhi uint64
	nlo := uint64(n)
	if n < 0 {
		nhi = math.MaxUint64
	}
	return i.hi == nhi && i.lo == nlo
}

// LessThan returns true if i < n.
func (i Int) LessThan(n Int) bool {
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
func (i Int) LessThan64(n int64) bool {
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

// LessThanOrEqual returns true if i <= n.
func (i Int) LessThanOrEqual(n Int) bool {
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

// LessThanOrEqual64 returns true if i <= n.
func (i Int) LessThanOrEqual64(n int64) bool {
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
func (i Int) Mul(n Int) Int {
	hi, lo := bits.Mul64(i.lo, n.lo)
	return Int{
		hi: hi + i.hi*n.lo + i.lo*n.hi,
		lo: lo,
	}
}

// Mul64 returns i * n.
func (i Int) Mul64(n int64) Int {
	return i.Mul(IntFrom64(n))
}

// Div returns i / n. If n == 0, a divide by zero panic will occur.
func (i Int) Div(n Int) Int {
	qSign := 1
	if i.LessThan(Int{}) {
		qSign = -1
		//goland:noinspection GoAssignmentToReceiver
		i = i.Neg()
	}
	if n.LessThan(Int{}) {
		qSign = -qSign
		n = n.Neg()
	}
	q := Int(Uint(i).Div(Uint(n)))
	if qSign < 0 {
		q = q.Neg()
	}
	return q
}

// Div64 returns i / n. If n == 0, a divide by zero panic will occur.
func (i Int) Div64(n int64) Int {
	qSign := 1
	if i.LessThan(Int{}) {
		qSign = -1
		//goland:noinspection GoAssignmentToReceiver
		i = i.Neg()
	}
	if n < 0 {
		qSign = -qSign
		n = -n
	}
	q := Int(Uint(i).Div64(uint64(n)))
	if qSign < 0 {
		q = q.Neg()
	}
	return q
}

// DivMod returns both the result of i / n as well i % n. If n == 0, a divide by zero panic will occur.
func (i Int) DivMod(n Int) (q, r Int) {
	qSign := 1
	rSign := 1
	if i.LessThan(Int{}) {
		qSign = -1
		rSign = -1
		//goland:noinspection GoAssignmentToReceiver
		i = i.Neg()
	}
	if n.LessThan(Int{}) {
		qSign = -qSign
		n = n.Neg()
	}
	qu, ru := Uint(i).DivMod(Uint(n))
	q = Int(qu)
	r = Int(ru)
	if qSign < 0 {
		q = q.Neg()
	}
	if rSign < 0 {
		r = r.Neg()
	}
	return q, r
}

// DivMod64 returns both the result of i / n as well i % n. If n == 0, a divide by zero panic will occur.
func (i Int) DivMod64(n int64) (q, r Int) {
	var hi uint64
	if n < 0 {
		hi = math.MaxUint64
	}
	return i.DivMod(Int{hi: hi, lo: uint64(n)})
}

// Mod returns i % n. If n == 0, a divide by zero panic will occur.
func (i Int) Mod(n Int) (r Int) {
	_, r = i.DivMod(n)
	return r
}

// Mod64 returns i % n. If n == 0, a divide by zero panic will occur.
func (i Int) Mod64(n int64) (r Int) {
	_, r = i.DivMod64(n)
	return r
}

// String implements fmt.Stringer.
func (i Int) String() string {
	if i.hi == 0 {
		if i.lo == 0 {
			return "0"
		}
		return strconv.FormatUint(i.lo, 10)
	}
	return i.AsBigInt().String()
}

// Format implements fmt.Formatter.
func (i Int) Format(s fmt.State, c rune) {
	i.AsBigInt().Format(s, c)
}

// Scan implements fmt.Scanner.
func (i *Int) Scan(state fmt.ScanState, _ rune) error {
	t, err := state.Token(true, nil)
	if err != nil {
		return errs.Wrap(err)
	}
	var v Int
	if v, err = IntFromString(string(t)); err != nil {
		return errs.Wrap(err)
	}
	*i = v
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int) UnmarshalText(text []byte) error {
	v, err := IntFromString(string(text))
	if err != nil {
		return err
	}
	*i = v
	return nil
}

// Float64 implements json.Number. Intentionally always returns an error, as we never want to emit floating point values
// into json for Int.
func (i Int) Float64() (float64, error) {
	return 0, errNoFloat64
}

// Int64 implements json.Number.
func (i Int) Int64() (int64, error) {
	if !i.IsInt64() {
		return 0, errDoesNotFitInInt64
	}
	return i.AsInt64(), nil
}

// MarshalJSON implements json.Marshaler.
func (i Int) MarshalJSON() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int) UnmarshalJSON(in []byte) error {
	v, err := IntFromString(string(in))
	if err != nil {
		return err
	}
	*i = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (i Int) MarshalYAML() (any, error) {
	return i.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (i *Int) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := IntFromString(str)
	if err != nil {
		return err
	}
	*i = v
	return nil
}
