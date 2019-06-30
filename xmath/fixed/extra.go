package fixed

import "errors"

var (
	errDoesNotFitInFloat64 = errors.New("does not fit in float64")
	errDoesNotFitInInt64   = errors.New("does not fit in int64")
)
