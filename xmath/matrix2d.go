package xmath

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
)

// Matrix2D provides a 2D matrix.
type Matrix2D struct {
	XX, YX, XY, YY, X0, Y0 float64
}

// NewMatrix2D creates a new 2D matrix.
func NewMatrix2D(xx, yx, xy, yy, x0, y0 float64) *Matrix2D {
	return &Matrix2D{XX: xx, YX: yx, XY: xy, YY: yy, X0: x0, Y0: y0}
}

// NewIdentityMatrix2D creates a new identity transformation 2D matrix.
func NewIdentityMatrix2D() *Matrix2D {
	return &Matrix2D{XX: 1, YY: 1}
}

// NewTranslationMatrix2D creates a new 2D matrix that translates by 'tx' and
// 'ty'.
func NewTranslationMatrix2D(tx, ty float64) *Matrix2D {
	return &Matrix2D{XX: 1, YY: 1, X0: tx, Y0: ty}
}

// NewScaleMatrix2D creates a new 2D matrix that scales by 'sx' and 'sy'.
func NewScaleMatrix2D(sx, sy float64) *Matrix2D {
	return &Matrix2D{XX: sx, YY: sy}
}

// NewRotationMatrix2D creates a new 2D matrix that rotates by 'radians'.
// Positive values are clockwise.
func NewRotationMatrix2D(radians float64) *Matrix2D {
	s := math.Sin(radians)
	c := math.Cos(radians)
	return &Matrix2D{XX: c, YX: s, XY: -s, YY: c}
}

// Translate this matrix by 'tx' and 'ty'.
func (m *Matrix2D) Translate(tx, ty float64) {
	m.X0 += tx
	m.Y0 += ty
}

// Scale this matrix by 'sx' and 'sy'.
func (m *Matrix2D) Scale(sx, sy float64) {
	m.XX *= sx
	m.YX *= sy
	m.XY *= sx
	m.YY *= sy
	m.X0 *= sx
	m.Y0 *= sy
}

// Rotate this matrix by 'radians'. Positive values are clockwise.
func (m *Matrix2D) Rotate(radians float64) {
	s := math.Sin(radians)
	c := math.Cos(radians)
	x := m.XX*c - s*m.YX
	m.YX = m.XX*s + m.YX*c
	m.XX = x
	x = m.XY*c - s*m.YY
	m.YY = m.XY*s + m.YY*c
	m.XY = x
	x = m.X0*c - s*m.Y0
	m.Y0 = m.X0*s + m.Y0*c
	m.X0 = x
}

// Multiply this matrix by 'other'.
func (m *Matrix2D) Multiply(other *Matrix2D) {
	x := m.XX*other.XX + m.YX*other.XY
	m.YX = m.XX*other.YX + m.YX*other.YY
	m.XX = x
	x = m.XY*other.XX + m.YY*other.XY
	m.YY = m.XY*other.YX + m.YY*other.YY
	m.XY = x
	x = m.X0*other.XX + m.Y0*other.XY + other.X0
	m.Y0 = m.X0*other.YX + m.Y0*other.YY + other.Y0
	m.X0 = x
}

// TransformDistance returns the result of transforming the distance vector
// (size.Width and size.Height) by this matrix. This is similar to
// TransformPoint(), except that the translation components of the
// transformation are ignored.
func (m *Matrix2D) TransformDistance(size geom.Size) geom.Size {
	x := m.XX*size.Width + m.XY*size.Height
	size.Height = m.YX*size.Width + m.YY*size.Height
	size.Width = x
	return size
}

// TransformPoint returns the result of transforming the point by this matrix.
func (m *Matrix2D) TransformPoint(where geom.Point) geom.Point {
	x := m.XX*where.X + m.XY*where.Y + m.X0
	where.Y = m.YX*where.X + m.YY*where.Y + m.Y0
	where.X = x
	return where
}
