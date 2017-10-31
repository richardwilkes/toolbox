package xmath

// AbsInt returns the absolute value of x.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MinInt returns the smaller of a or b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt returns the larger of a or b.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
