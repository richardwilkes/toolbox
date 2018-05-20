package jot

// LoggerWriter provides a bridge between the standard log.Logger and the jot
// package. You can use it like this:
//
// log.New(&jot.LoggerWriter{}, "", 0)
//
// This will send all output for this logger to the jot.Error() call.
//
// You can also set the Filter function to direct the output to a particular
// jot logging method:
//
// log.New(&jot.LoggerWriter{Filter: jot.Info}), "", 0)
type LoggerWriter struct {
	Filter func(v ...interface{})
}

// Write implements the io.Writer interface required by log.Logger.
func (w *LoggerWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		filter := w.Filter
		if filter == nil {
			filter = Error
		}
		filter(string(p[:len(p)-1]))
		Flush() // To ensure the output is recorded.
	}
	return len(p), nil
}
