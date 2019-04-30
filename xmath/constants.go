// Package xmath provides math-related utilities.
package xmath

import (
	"math"
)

const (
	// DegreesToRadians converts a value in degrees to radians when multiplied
	// with the value.
	DegreesToRadians = math.Pi / 180
	// RadiansToDegrees converts a value in radians to degrees when multiplied
	// with the value.
	RadiansToDegrees = 180 / math.Pi
)
