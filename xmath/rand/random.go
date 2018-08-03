package rand

import (
	"crypto/rand"
	mrnd "math/rand"
)

var cryptoRandInstance = &cryptoRand{}

// NewCryptoRand returns a Randomizer based on the crypto/rand package. This
// method returns a shared singleton instance and does not allocate.
func NewCryptoRand() Randomizer {
	return cryptoRandInstance
}

type cryptoRand struct {
}

func (r *cryptoRand) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	var buffer [8]byte
	size := 8
	n64 := int64(n)
	for i := 1; i < 8; i++ {
		if n64 < int64(1)<<uint(i*8) {
			size = i
			break
		}
	}
	if _, err := rand.Read(buffer[:size]); err != nil {
		// Fallback to pseudo-random number generator if crypto/rand fails
		return mrnd.Intn(n)
	}
	var v int
	for i := size - 1; i >= 0; i-- {
		v |= int(buffer[i]) << uint(i*8)
	}
	if v < 0 {
		v = -v
	}
	return v % n
}
