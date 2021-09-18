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
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

const (
	divBinaryShiftThreshold = 16
	divByZero               = "divide by zero"
	bit32                   = uint64(1) << 32
)

// MaxUint128 is the maximum value representable by a Uint128.
var MaxUint128 = Uint128{hi: math.MaxUint64, lo: math.MaxUint64}

var (
	intSize                      = 32 << (^uint(0) >> 63)
	maxUint64Float               = float64(math.MaxUint64)
	maxRepresentableUint64Float  = math.Nextafter(maxUint64Float, 0)
	maxRepresentableUint128Float = math.Nextafter(float64(340282366920938463463374607431768211455), 0)
	wrapUint64Float              = float64(math.MaxUint64) + 1
	errNoFloat64                 = errors.New("no float64 conversion for json/yaml")
	errDoesNotFitInInt64         = errors.New("does not fit in int64")
)

// RandomSource defines the method required of a source of random bits. This is a subset of the rand.Source64 interface.
type RandomSource interface {
	Uint64() uint64
}

// Uint128 represents an unsigned 128-bit integer.
type Uint128 struct {
	hi uint64
	lo uint64
}

// Uint128From64 creates a Uint128 from a uint64 value.
func Uint128From64(v uint64) Uint128 {
	return Uint128{lo: v}
}

// Uint128FromFloat64 creates a Uint128 from a float64 value.
func Uint128FromFloat64(f float64) Uint128 {
	switch {
	case f <= 0 || f != f: // <= 0 or NaN
		return Uint128{}
	case f <= maxRepresentableUint64Float:
		return Uint128{lo: uint64(f)}
	case f <= maxRepresentableUint128Float:
		return Uint128{
			hi: uint64(f / wrapUint64Float),
			lo: uint64(math.Mod(f, wrapUint64Float)),
		}
	default:
		return MaxUint128
	}
}

// Uint128FromBigInt creates a Uint128 from a big.Int.
func Uint128FromBigInt(v *big.Int) Uint128 {
	if v.Sign() < 0 {
		return Uint128{}
	}
	words := v.Bits()
	switch len(words) {
	case 0:
		return Uint128{}
	case 1:
		return Uint128{lo: uint64(words[0])}
	case 2:
		if intSize == 64 {
			return Uint128{
				hi: uint64(words[1]),
				lo: uint64(words[0]),
			}
		}
		return Uint128{lo: (uint64(words[1]) << 32) | (uint64(words[0]))}
	case 3:
		if intSize == 64 {
			return MaxUint128
		}
		return Uint128{
			hi: uint64(words[2]),
			lo: (uint64(words[1]) << 32) | (uint64(words[0])),
		}
	case 4:
		if intSize == 64 {
			return MaxUint128
		}
		return Uint128{
			hi: (uint64(words[3]) << 32) | (uint64(words[2])),
			lo: (uint64(words[1]) << 32) | (uint64(words[0])),
		}
	default:
		return MaxUint128
	}
}

// Uint128FromString creates a Uint128 from a string.
func Uint128FromString(s string) (Uint128, error) {
	b, err := parseToBigInt(s)
	if err != nil {
		return Uint128{}, err
	}
	return Uint128FromBigInt(b), nil
}

func parseToBigInt(s string) (*big.Int, error) {
	var b *big.Int
	var ok bool
	if strings.ContainsAny(s, "Ee") {
		// Given a floating-point value with an exponent, which technically isn't valid input, but we'll try to convert
		// it anyway.
		var f *big.Float
		f, ok = new(big.Float).SetString(s)
		if ok && !f.IsInt() {
			ok = false
		}
		if ok {
			b, _ = f.Int(nil)
		}
	} else {
		b, ok = new(big.Int).SetString(s, 0)
	}
	if !ok {
		return nil, errs.Newf("invalid input: %s", s)
	}
	return b, nil
}

// Uint128FromStringNoCheck creates a Uint128 from a string. Unlike Uint128FromString, this allows any string as input.
func Uint128FromStringNoCheck(s string) Uint128 {
	out, _ := Uint128FromString(s) //nolint:errcheck // Failure results in 0
	return out
}

// Uint128FromComponents creates a Uint128 from two uint64 values representing the high and low bits.
func Uint128FromComponents(high, low uint64) Uint128 {
	return Uint128{hi: high, lo: low}
}

// Uint128FromRand generates an unsigned 128-bit random integer.
func Uint128FromRand(source RandomSource) Uint128 {
	return Uint128{hi: source.Uint64(), lo: source.Uint64()}
}

// Components returns the two uint64 values representing the high and low bits.
func (u Uint128) Components() (high, low uint64) {
	return u.hi, u.lo
}

// IsZero returns true if the value is 0.
func (u Uint128) IsZero() bool {
	return u.hi|u.lo == 0
}

// ToBigInt stores the Uint128's value into the specified big.Int.
func (u Uint128) ToBigInt(b *big.Int) {
	words := b.Bits()
	if intSize == 64 {
		if len(words) < 2 {
			words = append(words, make([]big.Word, 2-len(words))...)
		}
		words = words[:2]
		words[0] = big.Word(u.lo)
		words[1] = big.Word(u.hi)
	} else {
		if len(words) < 4 {
			words = append(words, make([]big.Word, 4-len(words))...)
		}
		words = words[:4]
		words[0] = big.Word(u.lo & 0xFFFFFFFF)
		words[1] = big.Word(u.lo >> 32)
		words[2] = big.Word(u.hi & 0xFFFFFFFF)
		words[3] = big.Word(u.hi >> 32)
	}
	b.SetBits(words)
}

// AsBigInt returns the Uint128 as a big.Int.
func (u Uint128) AsBigInt() *big.Int {
	var b big.Int
	u.ToBigInt(&b)
	return &b
}

// AsBigFloat returns the Uint128 as a big.Float.
func (u Uint128) AsBigFloat() *big.Float {
	return new(big.Float).SetInt(u.AsBigInt())
}

// AsFloat64 returns the Uint128 as a float64.
func (u Uint128) AsFloat64() float64 {
	if u.hi == 0 {
		if u.lo == 0 {
			return 0
		}
		return float64(u.lo)
	}
	return (float64(u.hi) * wrapUint64Float) + float64(u.lo)
}

// IsInt128 returns true if this value can be represented as an Int128 without any loss.
func (u Uint128) IsInt128() bool {
	return u.hi&signBit == 0
}

// AsInt128 returns the Uint128 as an Int128.
func (u Uint128) AsInt128() Int128 {
	return Int128(u)
}

// IsUint64 returns true if this value can be represented as a uint64 without any loss.
func (u Uint128) IsUint64() bool {
	return u.hi == 0
}

// AsUint64 returns the Uint128 as a uint64.
func (u Uint128) AsUint64() uint64 {
	return u.lo
}

// Add returns u + n.
func (u Uint128) Add(n Uint128) Uint128 {
	lo, carry := bits.Add64(u.lo, n.lo, 0)
	hi, _ := bits.Add64(u.hi, n.hi, carry)
	return Uint128{
		hi: hi,
		lo: lo,
	}
}

// Add64 returns u + n.
func (u Uint128) Add64(n uint64) Uint128 {
	lo, carry := bits.Add64(u.lo, n, 0)
	return Uint128{
		hi: u.hi + carry,
		lo: lo,
	}
}

// Sub returns u - n.
func (u Uint128) Sub(n Uint128) Uint128 {
	lo, borrow := bits.Sub64(u.lo, n.lo, 0)
	hi, _ := bits.Sub64(u.hi, n.hi, borrow)
	return Uint128{
		hi: hi,
		lo: lo,
	}
}

// Sub64 returns u - n.
func (u Uint128) Sub64(n uint64) Uint128 {
	lo, borrow := bits.Sub64(u.lo, n, 0)
	return Uint128{
		hi: u.hi - borrow,
		lo: lo,
	}
}

// Inc returns u + 1.
func (u Uint128) Inc() Uint128 {
	lo, carry := bits.Add64(u.lo, 1, 0)
	return Uint128{
		hi: u.hi + carry,
		lo: lo,
	}
}

// Dec returns u - 1.
func (u Uint128) Dec() Uint128 {
	lo, borrow := bits.Sub64(u.lo, 1, 0)
	return Uint128{
		hi: u.hi - borrow,
		lo: lo,
	}
}

// Cmp returns 1 if u > n, 0 if u == n, and -1 if u < n.
func (u Uint128) Cmp(n Uint128) int {
	switch {
	case u.hi == n.hi:
		if u.lo > n.lo {
			return 1
		} else if u.lo < n.lo {
			return -1
		}
	case u.hi > n.hi:
		return 1
	case u.hi < n.hi:
		return -1
	}
	return 0
}

// Cmp64 returns 1 if u > n, 0 if u == n, and -1 if u < n.
func (u Uint128) Cmp64(n uint64) int {
	switch {
	case u.hi > 0 || u.lo > n:
		return 1
	case u.lo < n:
		return -1
	default:
		return 0
	}
}

// GreaterThan returns true if u > n.
func (u Uint128) GreaterThan(n Uint128) bool {
	return u.hi > n.hi || (u.hi == n.hi && u.lo > n.lo)
}

// GreaterThan64 returns true if u > n.
func (u Uint128) GreaterThan64(n uint64) bool {
	return u.hi > 0 || u.lo > n
}

// GreaterOrEqualTo returns true if u >= n.
func (u Uint128) GreaterOrEqualTo(n Uint128) bool {
	return u.hi > n.hi || (u.hi == n.hi && u.lo >= n.lo)
}

// GreaterOrEqualTo64 returns true if u >= n.
func (u Uint128) GreaterOrEqualTo64(n uint64) bool {
	return u.hi > 0 || u.lo >= n
}

// Equal returns true if u == n.
func (u Uint128) Equal(n Uint128) bool {
	return u.hi == n.hi && u.lo == n.lo
}

// Equal64 returns true if u == n.
func (u Uint128) Equal64(n uint64) bool {
	return u.hi == 0 && u.lo == n
}

// LessThan returns true if u < n.
func (u Uint128) LessThan(n Uint128) bool {
	return u.hi < n.hi || (u.hi == n.hi && u.lo < n.lo)
}

// LessThan64 returns true if u < n.
func (u Uint128) LessThan64(n uint64) bool {
	return u.hi == 0 && u.lo < n
}

// LessOrEqualTo returns true if u <= n.
func (u Uint128) LessOrEqualTo(n Uint128) bool {
	return u.hi < n.hi || (u.hi == n.hi && u.lo <= n.lo)
}

// LessOrEqualTo64 returns true if u <= n.
func (u Uint128) LessOrEqualTo64(n uint64) bool {
	return u.hi == 0 && u.lo <= n
}

// BitLen returns the length of the absolute value of u in bits. The bit length of 0 is 0.
func (u Uint128) BitLen() int {
	if u.hi != 0 {
		return bits.Len64(u.hi) + 64
	}
	return bits.Len64(u.lo)
}

// OnesCount returns the number of one bits ("population count") in u.
func (u Uint128) OnesCount() int {
	if u.hi != 0 {
		return bits.OnesCount64(u.hi) + 64
	}
	return bits.OnesCount64(u.lo)
}

// Bit returns the value of the i'th bit of x. That is, it returns (x>>i)&1. If the bit index is less than 0 or greater
// than 127, zero will be returned.
func (u Uint128) Bit(i int) uint {
	switch {
	case i < 0 || i > 127:
		return 0
	case i < 64:
		return uint((u.lo >> uint(i)) & 1)
	default:
		return uint((u.hi >> uint(i-64)) & 1)
	}
}

// SetBit returns a Uint128 with u's i'th bit set to b (0 or 1). Values of b that are not 0 will be treated as 1. If the
// bit index is less than 0 or greater than 127, nothing will happen.
func (u Uint128) SetBit(i int, b uint) Uint128 {
	if i < 0 || i > 127 {
		return u
	}
	if b == 0 {
		if i >= 64 {
			u.hi &^= 1 << uint(i-64)
		} else {
			u.lo &^= 1 << uint(i)
		}
	} else {
		if i >= 64 {
			u.hi |= 1 << uint(i-64)
		} else {
			u.lo |= 1 << uint(i)
		}
	}
	return u
}

// Not returns ^u.
func (u Uint128) Not() Uint128 {
	return Uint128{
		hi: ^u.hi,
		lo: ^u.lo,
	}
}

// And returns u & n.
func (u Uint128) And(n Uint128) Uint128 {
	return Uint128{
		hi: u.hi & n.hi,
		lo: u.lo & n.lo,
	}
}

// And64 returns u & n.
func (u Uint128) And64(n uint64) Uint128 {
	return Uint128{lo: u.lo & n}
}

// AndNot returns u &^ n.
func (u Uint128) AndNot(n Uint128) Uint128 {
	return Uint128{
		hi: u.hi &^ n.hi,
		lo: u.lo &^ n.lo,
	}
}

// AndNot64 returns u &^ n.
func (u Uint128) AndNot64(n Uint128) Uint128 {
	return Uint128{
		hi: u.hi,
		lo: u.lo &^ n.lo,
	}
}

// Or returns u | n.
func (u Uint128) Or(n Uint128) Uint128 {
	return Uint128{
		hi: u.hi | n.hi,
		lo: u.lo | n.lo,
	}
}

// Or64 returns u | n.
func (u Uint128) Or64(n uint64) Uint128 {
	return Uint128{
		hi: u.hi,
		lo: u.lo | n,
	}
}

// Xor returns u ^ n.
func (u Uint128) Xor(n Uint128) Uint128 {
	return Uint128{
		hi: u.hi ^ n.hi,
		lo: u.lo ^ n.lo,
	}
}

// Xor64 returns u ^ n.
func (u Uint128) Xor64(n uint64) Uint128 {
	return Uint128{
		hi: u.hi,
		lo: u.lo ^ n,
	}
}

// LeadingZeros returns the number of leading bits set to 0.
func (u Uint128) LeadingZeros() uint {
	if u.hi == 0 {
		return uint(bits.LeadingZeros64(u.lo)) + 64
	}
	return uint(bits.LeadingZeros64(u.hi))
}

// TrailingZeros returns the number of trailing bits set to 0.
func (u Uint128) TrailingZeros() uint {
	if u.lo == 0 {
		return uint(bits.TrailingZeros64(u.hi)) + 64
	}
	return uint(bits.TrailingZeros64(u.lo))
}

// LeftShift returns u << n.
func (u Uint128) LeftShift(n uint) Uint128 {
	switch {
	case n == 0:
	case n > 64:
		u.hi = u.lo << (n - 64)
		u.lo = 0
	case n < 64:
		u.hi = (u.hi << n) | (u.lo >> (64 - n))
		u.lo <<= n
	default:
		u.hi = u.lo
		u.lo = 0
	}
	return u
}

// RightShift returns u >> n.
func (u Uint128) RightShift(n uint) Uint128 {
	switch {
	case n == 0:
	case n > 64:
		u.lo = u.hi >> (n - 64)
		u.hi = 0
	case n < 64:
		u.lo = (u.lo >> n) | (u.hi << (64 - n))
		u.hi >>= n
	default:
		u.lo = u.hi
		u.hi = 0
	}
	return u
}

// Mul returns u * n.
func (u Uint128) Mul(n Uint128) Uint128 {
	hi, lo := bits.Mul64(u.lo, n.lo)
	return Uint128{
		hi: hi + u.hi*n.lo + u.lo*n.hi,
		lo: lo,
	}
}

// Mul64 returns u * n.
func (u Uint128) Mul64(n uint64) (dest Uint128) {
	x0 := u.lo & 0xFFFFFFFF
	x1 := u.lo >> 32
	y0 := n & 0xFFFFFFFF
	y1 := n >> 32
	t := x1*y0 + (x0*y0)>>32
	return Uint128{
		hi: (x1 * y1) + (t >> 32) + (((t & 0xFFFFFFFF) + (x0 * y1)) >> 32) + u.hi*n,
		lo: u.lo * n,
	}
}

// Div returns u / n. If n == 0, a divide by zero panic will occur.
func (u Uint128) Div(n Uint128) Uint128 {
	var nLoLeading0, nHiLeading0, nLeading0 uint
	if n.hi == 0 {
		if n.lo == 0 {
			panic(divByZero)
		}
		if n.lo == 1 { // divide by 1
			return u
		}
		if u.hi == 0 { // 64-bit division only
			u.lo /= n.lo
			return u
		}
		nLoLeading0 = uint(bits.LeadingZeros64(n.lo))
		nHiLeading0 = 64
		nLeading0 = nLoLeading0 + 64
	} else {
		nHiLeading0 = uint(bits.LeadingZeros64(n.hi))
		nLeading0 = nHiLeading0
	}
	nTrailing0 := n.TrailingZeros()
	if (nLeading0 + nTrailing0) == 127 { // Only one bit set in divisor, so use right shift
		return u.RightShift(nTrailing0)
	}
	if cmp := u.Cmp(n); cmp < 0 {
		return Uint128{} // nothing but remainder
	} else if cmp == 0 { // division by same value
		return Uint128{lo: 1}
	}
	uLeading0 := u.LeadingZeros()
	if nLeading0-uLeading0 > divBinaryShiftThreshold {
		q, _ := u.divmod128by128(n, nHiLeading0, nLoLeading0)
		return q
	}
	q, _ := u.divmod128bin(n, uLeading0, nLeading0)
	return q
}

// Div64 returns u / n. If n == 0, a divide by zero panic will occur.
func (u Uint128) Div64(n uint64) Uint128 {
	if n == 0 {
		panic(divByZero)
	}
	if n == 1 {
		return u
	}
	if u.hi == 0 { // 64-bit division only
		u.lo /= n
		return u
	}
	nLoLeading0 := uint(bits.LeadingZeros64(n))
	nLeading0 := nLoLeading0 + 64
	nTrailing0 := uint(bits.TrailingZeros64(n))
	if nLeading0+nTrailing0 == 127 { // Only one bit set in divisor, so use right shift
		return u.RightShift(nTrailing0)
	}
	if cmp := u.Cmp64(n); cmp < 0 {
		return Uint128{} // nothing but remainder
	} else if cmp == 0 { // division by same value
		return Uint128{lo: 1}
	}
	uLeading0 := u.LeadingZeros()
	if nLeading0-uLeading0 > divBinaryShiftThreshold {
		if u.hi < n {
			u.lo, _ = u.divmod128by64(n, nLoLeading0)
			u.hi = 0
		} else {
			hi := u.hi / n
			u.hi %= n
			u.lo, _ = u.divmod128by64(n, nLoLeading0)
			u.hi = hi
		}
		return u
	}
	q, _ := u.divmod128bin(Uint128{lo: n}, uLeading0, nLeading0)
	return q
}

// DivMod returns both the result of u / n as well u % n. If n == 0, a divide by zero panic will occur.
func (u Uint128) DivMod(n Uint128) (q, r Uint128) {
	var nLoLeading0, nHiLeading0, nLeading0 uint
	if n.hi == 0 {
		if n.lo == 0 {
			panic(divByZero)
		}
		if n.lo == 1 { // divide by 1
			return u, r
		}
		if u.hi == 0 { // 64-bit division only
			q.lo = u.lo / n.lo
			r.lo = u.lo % n.lo
			return q, r
		}
		nLoLeading0 = uint(bits.LeadingZeros64(n.lo))
		nHiLeading0 = 64
		nLeading0 = nLoLeading0 + 64
	} else {
		nHiLeading0 = uint(bits.LeadingZeros64(n.hi))
		nLeading0 = nHiLeading0
	}
	nTrailing0 := n.TrailingZeros()
	if (nLeading0 + nTrailing0) == 127 { // Only one bit set in divisor, so use right shift
		q = u.RightShift(nTrailing0)
		r = n.Dec().And(u)
		return q, r
	}
	if cmp := u.Cmp(n); cmp < 0 {
		return q, u // nothing but remainder
	} else if cmp == 0 { // division by same value
		q.lo = 1
		return q, r
	}
	uLeading0 := u.LeadingZeros()
	if nLeading0-uLeading0 > divBinaryShiftThreshold {
		return u.divmod128by128(n, nHiLeading0, nLoLeading0)
	}
	return u.divmod128bin(n, uLeading0, nLeading0)
}

// DivMod64 returns both the result of u / n as well u % n. If n == 0, a divide by zero panic will occur.
func (u Uint128) DivMod64(n uint64) (q, r Uint128) {
	if n == 0 {
		panic(divByZero)
	}
	if n == 1 {
		return u, r
	}
	if u.hi == 0 { // 64-bit division only
		q.lo = u.lo / n
		r.lo = u.lo % n
		return q, r
	}
	nLoLeading0 := uint(bits.LeadingZeros64(n))
	nLeading0 := nLoLeading0 + 64
	nTrailing0 := uint(bits.TrailingZeros64(n))
	if nLeading0+nTrailing0 == 127 { // Only one bit set in divisor, so use right shift
		q = u.RightShift(nTrailing0)
		r = u.And64(n - 1)
		return q, r
	}
	if cmp := u.Cmp64(n); cmp < 0 {
		return q, u // nothing but remainder
	} else if cmp == 0 { // division by same value
		q.lo = 1
		return q, r
	}
	uLeading0 := u.LeadingZeros()
	if nLeading0-uLeading0 > divBinaryShiftThreshold {
		if u.hi < n {
			q.lo, r.lo = u.divmod128by64(n, nLoLeading0)
		} else {
			q.hi = u.hi / n
			u.hi %= n
			q.lo, r.lo = u.divmod128by64(n, nLoLeading0)
		}
		return q, r
	}
	return u.divmod128bin(Uint128{lo: n}, uLeading0, nLeading0)
}

// Mod returns u % n. If n == 0, a divide by zero panic will occur.
func (u Uint128) Mod(n Uint128) Uint128 {
	var nLoLeading0, nHiLeading0, nLeading0 uint
	if n.hi == 0 {
		if n.lo == 0 {
			panic(divByZero)
		}
		if n.lo == 1 { // divide by 1
			return Uint128{}
		}
		if u.hi == 0 { // 64-bit division only
			u.lo %= n.lo
			return u
		}
		nLoLeading0 = uint(bits.LeadingZeros64(n.lo))
		nHiLeading0 = 64
		nLeading0 = nLoLeading0 + 64
	} else {
		nHiLeading0 = uint(bits.LeadingZeros64(n.hi))
		nLeading0 = nHiLeading0
	}
	nTrailing0 := n.TrailingZeros()
	if (nLeading0 + nTrailing0) == 127 { // Only one bit set in divisor, so use right shift
		return n.Dec().And(u)
	}
	if cmp := u.Cmp(n); cmp < 0 {
		return u // nothing but remainder
	} else if cmp == 0 { // division by same value
		return Uint128{}
	}
	uLeading0 := u.LeadingZeros()
	if nLeading0-uLeading0 > divBinaryShiftThreshold {
		_, r := u.divmod128by128(n, nHiLeading0, nLoLeading0)
		return r
	}
	_, r := u.divmod128bin(n, uLeading0, nLeading0)
	return r
}

// Mod64 returns u % n. If n == 0, a divide by zero panic will occur.
func (u Uint128) Mod64(n uint64) Uint128 {
	if n == 0 {
		panic(divByZero)
	}
	if n == 1 {
		return Uint128{}
	}
	if u.hi == 0 { // 64-bit division only
		u.lo %= n
		return u
	}
	nLoLeading0 := uint(bits.LeadingZeros64(n))
	nLeading0 := nLoLeading0 + 64
	nTrailing0 := uint(bits.TrailingZeros64(n))
	if nLeading0+nTrailing0 == 127 { // Only one bit set in divisor, so use right shift
		return u.And64(n - 1)
	}
	if cmp := u.Cmp64(n); cmp < 0 {
		return u // nothing but remainder
	} else if cmp == 0 { // division by same value
		return Uint128{}
	}
	uLeading0 := u.LeadingZeros()
	if nLeading0-uLeading0 > divBinaryShiftThreshold {
		if u.hi >= n {
			u.hi %= n
		}
		_, r := u.divmod128by64(n, nLoLeading0)
		return Uint128{lo: r}
	}
	_, r := u.divmod128bin(Uint128{lo: n}, uLeading0, nLeading0)
	return r
}

// divmod128by64 was adapted from https://www.codeproject.com/Tips/785014/UInt-Division-Modulus
func (u Uint128) divmod128by64(n uint64, nLeading0 uint) (q, r uint64) {
	n <<= nLeading0
	vn1 := n >> 32
	vn0 := n & 0xffffffff
	if nLeading0 > 0 {
		u.hi = (u.hi << nLeading0) | (u.lo >> (64 - nLeading0))
		u.lo <<= nLeading0
	}
	un1 := u.lo >> 32
	un0 := u.lo & 0xffffffff
	q1 := u.hi / vn1
	rhat := u.hi % vn1
	left := q1 * vn0
	right := (rhat << 32) + un1
loop1:
	if (q1 >= bit32) || (left > right) {
		q1--
		rhat += vn1
		if rhat < bit32 {
			left -= vn0
			right = (rhat << 32) | un1
			goto loop1
		}
	}
	un21 := (u.hi << 32) + (un1 - (q1 * n))
	q0 := un21 / vn1
	rhat = un21 % vn1
	left = q0 * vn0
	right = (rhat << 32) | un0
loop2:
	if (q0 >= bit32) || (left > right) {
		q0--
		rhat += vn1
		if rhat < bit32 {
			left -= vn0
			right = (rhat << 32) | un0
			goto loop2
		}
	}
	return (q1 << 32) | q0, ((un21 << 32) + (un0 - (q0 * n))) >> nLeading0
}

// divmod128by128 was adapted from https://www.codeproject.com/Tips/785014/UInt-Division-Modulus
func (u Uint128) divmod128by128(n Uint128, nHiLeading0, nLoLeading0 uint) (q, r Uint128) {
	if n.hi == 0 {
		if u.hi < n.lo {
			q.lo, r.lo = u.divmod128by64(n.lo, nLoLeading0)
			return q, r
		}
		q.hi = u.hi / n.lo
		u.hi %= n.lo
		q.lo, r.lo = u.divmod128by64(n.lo, nLoLeading0)
		r.hi = 0
		return q, r
	}
	q.lo, _ = u.RightShift(1).divmod128by64(n.LeftShift(nHiLeading0).hi, nLoLeading0)
	q.lo >>= 63 - nHiLeading0
	if q.lo != 0 {
		q.lo--
	}
	r = u.Sub(q.Mul(n))
	if r.Cmp(n) >= 0 {
		q = q.Inc()
		r = r.Sub(n)
	}
	return q, r
}

// divmod128bin was adapted from https://www.codeproject.com/Tips/785014/UInt-Division-Modulus
func (u Uint128) divmod128bin(n Uint128, uLeading0, byLeading0 uint) (q, r Uint128) {
	shift := int(byLeading0 - uLeading0)
	n = n.LeftShift(uint(shift))
	for {
		if u.GreaterOrEqualTo(n) {
			//goland:noinspection GoAssignmentToReceiver
			u = u.Sub(n)
			q.lo |= 1
		}
		if shift <= 0 {
			break
		}
		n = n.RightShift(1)
		q = q.LeftShift(1)
		shift--
	}
	return q, u
}

// String implements fmt.Stringer.
func (u Uint128) String() string {
	if u.hi == 0 {
		if u.lo == 0 {
			return "0"
		}
		return strconv.FormatUint(u.lo, 10)
	}
	return u.AsBigInt().String()
}

// Format implements fmt.Formatter.
func (u Uint128) Format(s fmt.State, c rune) {
	u.AsBigInt().Format(s, c)
}

// Scan implements fmt.Scanner.
func (u *Uint128) Scan(state fmt.ScanState, verb rune) error {
	t, err := state.Token(true, nil)
	if err != nil {
		return errs.Wrap(err)
	}
	v, err := Uint128FromString(string(t))
	if err != nil {
		return errs.Wrap(err)
	}
	*u = v
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint128) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Uint128) UnmarshalText(text []byte) (err error) {
	v, err := Uint128FromString(string(text))
	if err != nil {
		return err
	}
	*u = v
	return nil
}

// Float64 implements json.Number. Intentionally always returns an error, as we never want to emit floating point values
// into json for Uint128.
func (u Uint128) Float64() (float64, error) {
	return 0, errNoFloat64
}

// Int64 implements json.Number.
func (u Uint128) Int64() (int64, error) {
	if u.IsInt128() {
		i128 := Int128(u)
		if i128.IsInt64() {
			return i128.AsInt64(), nil
		}
	}
	return 0, errDoesNotFitInInt64
}

// MarshalJSON implements json.Marshaler.
func (u Uint128) MarshalJSON() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Uint128) UnmarshalJSON(in []byte) error {
	v, err := Uint128FromString(string(in))
	if err != nil {
		return err
	}
	*u = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (u Uint128) MarshalYAML() (interface{}, error) {
	return u.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (u *Uint128) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := Uint128FromString(str)
	if err != nil {
		return err
	}
	*u = v
	return nil
}
