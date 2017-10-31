package safe

import (
	"bufio"
	"io"
)

// WriteFile uses writer to write data safely and atomically to a file.
func WriteFile(filename string, writer func(io.Writer) error) (err error) {
	var f *File
	f, err = Create(filename)
	if err != nil {
		return
	}
	w := bufio.NewWriterSize(f, 1<<16)
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	if err = writer(w); err != nil {
		return
	}
	if err = w.Flush(); err != nil {
		return
	}
	if err = f.Commit(); err != nil {
		return
	}
	return
}
