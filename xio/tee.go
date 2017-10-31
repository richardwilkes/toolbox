package xio

import "io"

// TeeWriter is a writer that writes to multiple other writers.
type TeeWriter struct {
	Writers []io.Writer
}

// Write to each of the underlying streams.
func (t *TeeWriter) Write(p []byte) (n int, err error) {
	var curErr error
	for _, w := range t.Writers {
		if n, curErr = w.Write(p); curErr != nil {
			err = curErr
		}
	}
	return
}
