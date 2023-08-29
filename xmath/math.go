// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xmath

import (
	"math"
	"reflect"

	"golang.org/x/exp/constraints"
)

// Numeric is a constraint that permits any integer or float type.
type Numeric interface {
	constraints.Float | constraints.Integer
}

// Abs returns the absolute value of x.
//
// Special cases are:
//
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs[T Numeric](x T) T {
	switch reflect.TypeOf(x).Kind() {
	case reflect.Float32:
		return T(math.Float32frombits(math.Float32bits(float32(x)) &^ (1 << 31)))
	case reflect.Float64:
		return T(math.Abs(float64(x)))
	default:
		if x < 0 {
			return -x
		}
		return x
	}
}

// Acos returns the arccosine of x.
func Acos[T constraints.Float](x T) T {
	return T(math.Acos(float64(x)))
}

// Acosh returns the inverse hyperbolic cosine of x.
//
// Special cases are:
//
//	Acosh(+Inf) = +Inf
//	Acosh(x) = NaN if x < 1
//	Acosh(NaN) = NaN
func Acosh[T constraints.Float](x T) T {
	return T(math.Acosh(float64(x)))
}

// Asin returns the arcsine, in radians, of x.
//
// Special cases are:
//
//	Asin(±0) = ±0
//	Asin(x) = NaN if x < -1 or x > 1
func Asin[T constraints.Float](x T) T {
	return T(math.Asin(float64(x)))
}

// Asinh returns the inverse hyperbolic sine of x.
//
// Special cases are:
//
//	Asinh(±0) = ±0
//	Asinh(±Inf) = ±Inf
//	Asinh(NaN) = NaN
func Asinh[T constraints.Float](x T) T {
	return T(math.Asinh(float64(x)))
}

// Atan2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
//
// Special cases are (in order):
//
//	Atan2(y, NaN) = NaN
//	Atan2(NaN, x) = NaN
//	Atan2(+0, x>=0) = +0
//	Atan2(-0, x>=0) = -0
//	Atan2(+0, x<=-0) = +Pi
//	Atan2(-0, x<=-0) = -Pi
//	Atan2(y>0, 0) = +Pi/2
//	Atan2(y<0, 0) = -Pi/2
//	Atan2(+Inf, +Inf) = +Pi/4
//	Atan2(-Inf, +Inf) = -Pi/4
//	Atan2(+Inf, -Inf) = 3Pi/4
//	Atan2(-Inf, -Inf) = -3Pi/4
//	Atan2(y, +Inf) = 0
//	Atan2(y>0, -Inf) = +Pi
//	Atan2(y<0, -Inf) = -Pi
//	Atan2(+Inf, x) = +Pi/2
//	Atan2(-Inf, x) = -Pi/2
func Atan2[T constraints.Float](y, x T) T {
	return T(math.Atan2(float64(y), float64(x)))
}

// Atan returns the arctangent, in radians, of x.
//
// Special cases are:
//
//	Atan(±0) = ±0
//	Atan(±Inf) = ±Pi/2
func Atan[T constraints.Float](x T) T {
	return T(math.Atan(float64(x)))
}

// Atanh returns the inverse hyperbolic tangent of x.
//
// Special cases are:
//
//	Atanh(1) = +Inf
//	Atanh(±0) = ±0
//	Atanh(-1) = -Inf
//	Atanh(x) = NaN if x < -1 or x > 1
//	Atanh(NaN) = NaN
func Atanh[T constraints.Float](x T) T {
	return T(math.Atanh(float64(x)))
}

// Cbrt returns the cube root of x.
func Cbrt[T constraints.Float](x T) T {
	return T(math.Cbrt(float64(x)))
}

// Ceil returns the smallest integer value greater than or equal to x.
func Ceil[T constraints.Float](x T) T {
	return T(math.Ceil(float64(x)))
}

// Copysign returns a value with the magnitude of x and the sign of y.
func Copysign[T constraints.Float](x, y T) T {
	if reflect.TypeOf(x).Kind() == reflect.Float32 {
		const sign = 1 << 31
		return T(math.Float32frombits(math.Float32bits(float32(x))&^sign | math.Float32bits(float32(y))&sign))
	}
	return T(math.Copysign(float64(x), float64(y)))
}

// Cos returns the cosine of the radian argument x.
//
// Special cases are:
//
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func Cos[T constraints.Float](x T) T {
	return T(math.Cos(float64(x)))
}

// Cosh returns the hyperbolic cosine of x.
//
// Special cases are:
//
//	Cosh(±0) = 1
//	Cosh(±Inf) = +Inf
//	Cosh(NaN) = NaN
func Cosh[T constraints.Float](x T) T {
	return T(math.Cosh(float64(x)))
}

// Dim returns the maximum of x-y or 0.
//
// Special cases are:
//
//	Dim(+Inf, +Inf) = NaN
//	Dim(-Inf, -Inf) = NaN
//	Dim(x, NaN) = Dim(NaN, x) = NaN
func Dim[T constraints.Float](x, y T) T {
	if v := x - y; v > 0 {
		return v
	}
	return 0
}

// Erf returns the error function of x.
//
// Special cases are:
//
//	Erf(+Inf) = 1
//	Erf(-Inf) = -1
//	Erf(NaN) = NaN
func Erf[T constraints.Float](x T) T {
	return T(math.Erf(float64(x)))
}

// Erfc returns the complementary error function of x.
//
// Special cases are:
//
//	Erfc(+Inf) = 0
//	Erfc(-Inf) = 2
//	Erfc(NaN) = NaN
func Erfc[T constraints.Float](x T) T {
	return T(math.Erfc(float64(x)))
}

// Erfinv returns the inverse error function of x.
//
// Special cases are:
//
//	Erfinv(1) = +Inf
//	Erfinv(-1) = -Inf
//	Erfinv(x) = NaN if x < -1 or x > 1
//	Erfinv(NaN) = NaN
func Erfinv[T constraints.Float](x T) T {
	return T(math.Erfinv(float64(x)))
}

// Erfcinv returns the inverse of Erfc(x).
//
// Special cases are:
//
//	Erfcinv(0) = +Inf
//	Erfcinv(2) = -Inf
//	Erfcinv(x) = NaN if x < 0 or x > 2
//	Erfcinv(NaN) = NaN
func Erfcinv[T constraints.Float](x T) T {
	return Erfinv(1 - x)
}

// Exp returns e**x, the base-e exponential of x.
//
// Special cases are:
//
//	Exp(+Inf) = +Inf
//	Exp(NaN) = NaN
//
// Very large values overflow to 0 or +Inf.
// Very small values underflow to 1.
func Exp[T constraints.Float](x T) T {
	return T(math.Exp(float64(x)))
}

// Exp2 returns 2**x, the base-2 exponential of x.
//
// Special cases are the same as Exp.
func Exp2[T constraints.Float](x T) T {
	return T(math.Exp2(float64(x)))
}

// Expm1 returns e**x - 1, the base-e exponential of x minus 1.
// It is more accurate than Exp(x) - 1 when x is near zero.
//
// Special cases are:
//
//	Expm1(+Inf) = +Inf
//	Expm1(-Inf) = -1
//	Expm1(NaN) = NaN
//
// Very large values overflow to -1 or +Inf.
func Expm1[T constraints.Float](x T) T {
	return T(math.Expm1(float64(x)))
}

// Floor returns the greatest integer value less than or equal to x.
func Floor[T constraints.Float](x T) T {
	return T(math.Floor(float64(x)))
}

// FMA returns x * y + z, computed with only one rounding.
// (That is, FMA returns the fused multiply-add of x, y, and z.)
func FMA[T constraints.Float](x, y, z T) T {
	return T(math.FMA(float64(x), float64(y), float64(z)))
}

// Frexp breaks f into a normalized fraction
// and an integral power of two.
// It returns frac and exp satisfying f == frac × 2**exp,
// with the absolute value of frac in the interval [½, 1).
//
// Special cases are:
//
//	Frexp(±0) = ±0, 0
//	Frexp(±Inf) = ±Inf, 0
//	Frexp(NaN) = NaN, 0
func Frexp[T constraints.Float](f T) (frac T, exp int) {
	fr, e := math.Frexp(float64(f))
	return T(fr), e
}

// Gamma returns the Gamma function of x.
//
// Special cases are:
//
//	Gamma(+Inf) = +Inf
//	Gamma(+0) = +Inf
//	Gamma(-0) = -Inf
//	Gamma(x) = NaN for integer x < 0
//	Gamma(-Inf) = NaN
//	Gamma(NaN) = NaN
func Gamma[T constraints.Float](x T) T {
	return T(math.Gamma(float64(x)))
}

// Hypot returns Sqrt(p*p + q*q), taking care to avoid
// unnecessary overflow and underflow.
//
// Special cases are:
//
//	Hypot(±Inf, q) = +Inf
//	Hypot(p, ±Inf) = +Inf
//	Hypot(NaN, q) = NaN
//	Hypot(p, NaN) = NaN
func Hypot[T constraints.Float](p, q T) T {
	return T(math.Hypot(float64(p), float64(q)))
}

// Ilogb returns the binary exponent of x as an integer.
//
// Special cases are:
//
//	Ilogb(±Inf) = MaxInt32
//	Ilogb(0) = MinInt32
//	Ilogb(NaN) = MaxInt32
func Ilogb[T constraints.Float](x T) int {
	return math.Ilogb(float64(x))
}

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func Inf[T constraints.Float](sign int) T {
	var t T
	if reflect.TypeOf(t).Kind() == reflect.Float32 {
		var v uint32
		if sign >= 0 {
			v = 0x7FF00000
		} else {
			v = 0xFFF00000
		}
		return T(math.Float32frombits(v))
	}
	return T(math.Inf(sign))
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func IsInf[T constraints.Float](f T, sign int) bool {
	if reflect.TypeOf(f).Kind() == reflect.Float32 {
		return sign >= 0 && f > math.MaxFloat32 || sign <= 0 && f < -math.MaxFloat32
	}
	return math.IsInf(float64(f), sign)
}

// IsNaN reports whether f is a "not-a-number" value.
func IsNaN[T constraints.Float](f T) bool {
	// Only NaNs satisfy f != f.
	return f != f
}

// J0 returns the order-zero Bessel function of the first kind.
//
// Special cases are:
//
//	J0(±Inf) = 0
//	J0(0) = 1
//	J0(NaN) = NaN
func J0[T constraints.Float](x T) T {
	return T(math.J0(float64(x)))
}

// J1 returns the order-one Bessel function of the first kind.
//
// Special cases are:
//
//	J1(±Inf) = 0
//	J1(NaN) = NaN
func J1[T constraints.Float](x T) T {
	return T(math.J1(float64(x)))
}

// Jn returns the order-n Bessel function of the first kind.
//
// Special cases are:
//
//	Jn(n, ±Inf) = 0
//	Jn(n, NaN) = NaN
func Jn[T constraints.Float](n int, x T) T {
	return T(math.Jn(n, float64(x)))
}

// Ldexp is the inverse of Frexp.
// It returns frac × 2**exp.
//
// Special cases are:
//
//	Ldexp(±0, exp) = ±0
//	Ldexp(±Inf, exp) = ±Inf
//	Ldexp(NaN, exp) = NaN
func Ldexp[T constraints.Float](frac T, exp int) T {
	return T(math.Ldexp(float64(frac), exp))
}

// Lgamma returns the natural logarithm and sign (-1 or +1) of Gamma(x).
//
// Special cases are:
//
//	Lgamma(+Inf) = +Inf
//	Lgamma(0) = +Inf
//	Lgamma(-integer) = +Inf
//	Lgamma(-Inf) = -Inf
//	Lgamma(NaN) = NaN
func Lgamma[T constraints.Float](x T) (lgamma T, sign int) {
	f64, s := math.Lgamma(float64(x))
	return T(f64), s
}

// Log returns the natural logarithm of x.
//
// Special cases are:
//
//	Log(+Inf) = +Inf
//	Log(0) = -Inf
//	Log(x < 0) = NaN
//	Log(NaN) = NaN
func Log[T constraints.Float](x T) T {
	return T(math.Log(float64(x)))
}

// Log10 returns the decimal logarithm of x. The special cases are the same as for Log.
func Log10[T constraints.Float](x T) T {
	return T(math.Log10(float64(x)))
}

// Log1p returns the natural logarithm of 1 plus its argument x. It is more accurate than Log(1 + x) when x is near
// zero.
//
// Special cases are:
//
//	Log1p(+Inf) = +Inf
//	Log1p(±0) = ±0
//	Log1p(-1) = -Inf
//	Log1p(x < -1) = NaN
//	Log1p(NaN) = NaN
func Log1p[T constraints.Float](x T) T {
	return T(math.Log1p(float64(x)))
}

// Log2 returns the binary logarithm of x. The special cases are the same as for Log.
func Log2[T constraints.Float](x T) T {
	return T(math.Log2(float64(x)))
}

// Logb returns the binary exponent of x.
//
// Special cases are:
//
//	Logb(±Inf) = +Inf
//	Logb(0) = -Inf
//	Logb(NaN) = NaN
func Logb[T constraints.Float](x T) T {
	return T(math.Logb(float64(x)))
}

// Max returns the larger of x or y.
//
// Special cases are:
//
//	Max(x, +Inf) = Max(+Inf, x) = +Inf
//	Max(x, NaN) = Max(NaN, x) = NaN
//	Max(+0, ±0) = Max(±0, +0) = +0
//	Max(-0, -0) = -0
//
// Deprecated: Use the Go 1.21+ built-in max() instead. August 8, 2023
func Max[T Numeric](a, b T) T {
	return max(a, b)
}

// MaxValue returns the maximum value for the type.
func MaxValue[T Numeric]() T {
	var t T
	v := reflect.Indirect(reflect.ValueOf(&t))
	switch reflect.TypeOf(t).Kind() {
	case reflect.Int:
		v.SetInt(math.MaxInt)
	case reflect.Int8:
		v.SetInt(math.MaxInt8)
	case reflect.Int16:
		v.SetInt(math.MaxInt16)
	case reflect.Int32:
		v.SetInt(math.MaxInt32)
	case reflect.Int64:
		v.SetInt(math.MaxInt64)
	case reflect.Uint:
		v.SetUint(math.MaxUint)
	case reflect.Uint8:
		v.SetUint(math.MaxUint8)
	case reflect.Uint16:
		v.SetUint(math.MaxUint16)
	case reflect.Uint32:
		v.SetUint(math.MaxUint32)
	case reflect.Uint64:
		v.SetUint(math.MaxUint64)
	case reflect.Uintptr:
		v.SetUint(math.MaxUint64)
	case reflect.Float32:
		v.SetFloat(math.MaxFloat32)
	case reflect.Float64:
		v.SetFloat(math.MaxFloat64)
	default:
		panic("unhandled type")
	}
	return t
}

// Min returns the smaller of x or y.
//
// Special cases are:
//
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN) = Min(NaN, x) = NaN
//	Min(-0, ±0) = Min(±0, -0) = -0
//
// Deprecated: Use the Go 1.21+ built-in min() instead. August 8, 2023
func Min[T Numeric](a, b T) T {
	return min(a, b)
}

// MinValue returns the minimum value for the type.
func MinValue[T Numeric]() T {
	var t T
	v := reflect.Indirect(reflect.ValueOf(&t))
	switch reflect.TypeOf(t).Kind() {
	case reflect.Int:
		v.SetInt(MinInt)
	case reflect.Int8:
		v.SetInt(MinInt8)
	case reflect.Int16:
		v.SetInt(MinInt16)
	case reflect.Int32:
		v.SetInt(MinInt32)
	case reflect.Int64:
		v.SetInt(MinInt64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		// 0
	case reflect.Float32:
		v.SetFloat(-MaxFloat32)
	case reflect.Float64:
		v.SetFloat(-MaxFloat64)
	default:
		panic("unhandled type")
	}
	return t
}

// SmallestPositiveNonZeroValue returns the smallest, positive, non-zero value for the type.
func SmallestPositiveNonZeroValue[T Numeric]() T {
	var t T
	v := reflect.Indirect(reflect.ValueOf(&t))
	switch reflect.TypeOf(t).Kind() {
	case reflect.Float32:
		v.SetFloat(SmallestNonzeroFloat32)
	case reflect.Float64:
		v.SetFloat(SmallestNonzeroFloat64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v.SetInt(1)
	default:
		panic("unhandled type")
	}
	return t
}

// Mod returns the floating-point remainder of x/y. The magnitude of the result is less than y and its sign agrees with
// that of x.
//
// Special cases are:
//
//	Mod(±Inf, y) = NaN
//	Mod(NaN, y) = NaN
//	Mod(x, 0) = NaN
//	Mod(x, ±Inf) = x
//	Mod(x, NaN) = NaN
func Mod[T constraints.Float](x, y T) T {
	return T(math.Mod(float64(x), float64(y)))
}

// Modf returns integer and fractional floating-point numbers
// that sum to f. Both values have the same sign as f.
//
// Special cases are:
//
//	Modf(±Inf) = ±Inf, NaN
//	Modf(NaN) = NaN, NaN
func Modf[T constraints.Float](f T) (i, frac T) {
	i64, f64 := math.Modf(float64(f))
	return T(i64), T(f64)
}

// NaN returns the "not-a-number" value.
func NaN[T constraints.Float]() T {
	var t T
	if reflect.TypeOf(t).Kind() == reflect.Float32 {
		return T(math.Float32frombits(0x7FF80001))
	}
	return T(math.NaN())
}

// Nextafter returns the next representable float32 value after x towards y.
//
// Special cases are:
//
//	Nextafter(x, x)   = x
//	Nextafter(NaN, y) = NaN
//	Nextafter(x, NaN) = NaN
func Nextafter[T constraints.Float](x, y T) (r T) {
	if reflect.TypeOf(x).Kind() == reflect.Float32 {
		return T(math.Nextafter32(float32(x), float32(y)))
	}
	return T(math.Nextafter(float64(x), float64(y)))
}

// Pow returns x**y, the base-x exponential of y.
//
// Special cases are (in order):
//
//	Pow(x, ±0) = 1 for any x
//	Pow(1, y) = 1 for any y
//	Pow(x, 1) = x for any x
//	Pow(NaN, y) = NaN
//	Pow(x, NaN) = NaN
//	Pow(±0, y) = ±Inf for y an odd integer < 0
//	Pow(±0, -Inf) = +Inf
//	Pow(±0, +Inf) = +0
//	Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
//	Pow(±0, y) = ±0 for y an odd integer > 0
//	Pow(±0, y) = +0 for finite y > 0 and not an odd integer
//	Pow(-1, ±Inf) = 1
//	Pow(x, +Inf) = +Inf for |x| > 1
//	Pow(x, -Inf) = +0 for |x| > 1
//	Pow(x, +Inf) = +0 for |x| < 1
//	Pow(x, -Inf) = +Inf for |x| < 1
//	Pow(+Inf, y) = +Inf for y > 0
//	Pow(+Inf, y) = +0 for y < 0
//	Pow(-Inf, y) = Pow(-0, -y)
//	Pow(x, y) = NaN for finite x < 0 and finite non-integer y
func Pow[T constraints.Float](x, y T) T {
	return T(math.Pow(float64(x), float64(y)))
}

// Pow10 returns 10**n, the base-10 exponential of n.
//
// Special cases are:
//
//	Pow10(n) =    0 for n < -323
//	Pow10(n) = +Inf for n > 308
func Pow10[T constraints.Float](n int) T {
	return T(math.Pow10(n))
}

// Remainder returns the IEEE 754 floating-point remainder of x/y.
//
// Special cases are:
//
//	Remainder(±Inf, y) = NaN
//	Remainder(NaN, y) = NaN
//	Remainder(x, 0) = NaN
//	Remainder(x, ±Inf) = x
//	Remainder(x, NaN) = NaN
func Remainder[T constraints.Float](x, y T) T {
	return T(math.Remainder(float64(x), float64(y)))
}

// Round returns the nearest integer, rounding half away from zero.
//
// Special cases are:
//
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
func Round[T constraints.Float](x T) T {
	return T(math.Round(float64(x)))
}

// RoundToEven returns the nearest integer, rounding ties to even.
//
// Special cases are:
//
//	RoundToEven(±0) = ±0
//	RoundToEven(±Inf) = ±Inf
//	RoundToEven(NaN) = NaN
func RoundToEven[T constraints.Float](x T) T {
	return T(math.RoundToEven(float64(x)))
}

// Signbit reports whether x is negative or negative zero.
func Signbit[T constraints.Float](x T) bool {
	if reflect.TypeOf(x).Kind() == reflect.Float32 {
		return math.Float32bits(float32(x))&(1<<31) != 0
	}
	return math.Signbit(float64(x))
}

// Sin returns the sine of the radian argument x.
//
// Special cases are:
//
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func Sin[T constraints.Float](x T) T {
	return T(math.Sin(float64(x)))
}

// Sincos returns Sin(x), Cos(x).
//
// Special cases are:
//
//	Sincos(±0) = ±0, 1
//	Sincos(±Inf) = NaN, NaN
//	Sincos(NaN) = NaN, NaN
func Sincos[T constraints.Float](x T) (sin, cos T) {
	s, c := math.Sincos(float64(x))
	return T(s), T(c)
}

// Sinh returns the hyperbolic sine of x.
//
// Special cases are:
//
//	Sinh(±0) = ±0
//	Sinh(±Inf) = ±Inf
//	Sinh(NaN) = NaN
func Sinh[T constraints.Float](x T) T {
	return T(math.Sinh(float64(x)))
}

// Sqrt returns the square root of x.
//
// Special cases are:
//
//	Sqrt(+Inf) = +Inf
//	Sqrt(±0) = ±0
//	Sqrt(x < 0) = NaN
//	Sqrt(NaN) = NaN
func Sqrt[T constraints.Float](x T) T {
	return T(math.Sqrt(float64(x)))
}

// Tan returns the tangent of the radian argument x.
//
// Special cases are:
//
//	Tan(±0) = ±0
//	Tan(±Inf) = NaN
//	Tan(NaN) = NaN
func Tan[T constraints.Float](x T) T {
	return T(math.Tan(float64(x)))
}

// Tanh returns the hyperbolic tangent of x.
//
// Special cases are:
//
//	Tanh(±0) = ±0
//	Tanh(±Inf) = ±1
//	Tanh(NaN) = NaN
func Tanh[T constraints.Float](x T) T {
	return T(math.Tanh(float64(x)))
}

// Trunc returns the integer value of x.
//
// Special cases are:
//
//	Trunc(±0) = ±0
//	Trunc(±Inf) = ±Inf
//	Trunc(NaN) = NaN
func Trunc[T constraints.Float](x T) T {
	return T(math.Trunc(float64(x)))
}

// Y0 returns the order-zero Bessel function of the second kind.
//
// Special cases are:
//
//	Y0(+Inf) = 0
//	Y0(0) = -Inf
//	Y0(x < 0) = NaN
//	Y0(NaN) = NaN
func Y0[T constraints.Float](x T) T {
	return T(math.Y0(float64(x)))
}

// Y1 returns the order-one Bessel function of the second kind.
//
// Special cases are:
//
//	Y1(+Inf) = 0
//	Y1(0) = -Inf
//	Y1(x < 0) = NaN
//	Y1(NaN) = NaN
func Y1[T constraints.Float](x T) T {
	return T(math.Y1(float64(x)))
}

// Yn returns the order-n Bessel function of the second kind.
//
// Special cases are:
//
//	Yn(n, +Inf) = 0
//	Yn(n ≥ 0, 0) = -Inf
//	Yn(n < 0, 0) = +Inf if n is odd, -Inf if n is even
//	Yn(n, x < 0) = NaN
//	Yn(n, NaN) = NaN
func Yn[T constraints.Float](n int, x T) T {
	return T(math.Yn(n, float64(x)))
}

// EqualWithin returns true if a and b are within the given tolerance of each other.
func EqualWithin[T constraints.Float](a, b, tolerance T) bool {
	if a == b {
		return true
	}
	delta := Abs(a - b)
	if delta <= tolerance {
		return true
	}
	var mv T
	if reflect.TypeOf(mv).Kind() == reflect.Float32 {
		mv = 0x1p-126
	} else {
		mv = 0x1p-1022
	}
	if delta <= mv {
		return delta <= tolerance*mv
	}
	return delta/max(Abs(a), Abs(b)) <= tolerance
}
