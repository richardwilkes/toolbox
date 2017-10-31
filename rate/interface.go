package rate

// Limiter provides a rate limiter.
type Limiter interface {
	// New returns a new limiter that is subordinate to this limiter, meaning
	// that its cap rate is also capped by its parent.
	New(cap int) Limiter

	// Cap returns the capacity per time period.
	Cap(applyParentCaps bool) int

	// SetCap sets the capacity.
	SetCap(cap int)

	// LastUsed returns the capacity used in the last time period.
	LastUsed() int

	// Use returns a channel that will return nil when the request is
	// successful, or an error if the request cannot be fulfilled.
	Use(amount int) <-chan error

	// Closed returns true if the limiter is closed.
	Closed() bool

	// Closes this limiter and any children it may have.
	Close()
}
