package rand

// Randomizer defines a source of random integer values.
type Randomizer interface {
	// Intn returns, as an int, a non-negative random number from 0 to n-1.
	// If n <= 0, the implementation should return 0.
	Intn(n int) int
}
