package xmath

// Round returns the closest integer.
func Round(x float64) float64 {
	return float64(int(x + 0.5))
}
