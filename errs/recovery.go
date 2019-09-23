package errs

// RecoveryHandler defines the callback used when panics occur with Recovery.
type RecoveryHandler func(error)

// Recovery provides an easy way to run code that may panic. 'handler' will be
// called with the panic turned into an error. Pass in nil to silently ignore
// any panic.
//
// Typical usage:
//
// func runSomeCode(handler errs.RecoveryHandler) {
//     defer errs.Recovery(handler)
//     // ... run the code here ...
// }
func Recovery(handler RecoveryHandler) {
	if recovered := recover(); recovered != nil && handler != nil {
		defer Recovery(nil) // Guard against a bad handler implementation
		handler(Newf("recovered from panic: %+v", recovered))
	}
}
