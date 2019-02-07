package txt

import "io"

// RuneReader implements io.RuneReader
type RuneReader struct {
	Src []rune
	Pos int
}

// ReadRune returns the next rune and its size in bytes.
func (rr *RuneReader) ReadRune() (r rune, size int, err error) {
	if rr.Pos >= len(rr.Src) {
		return -1, 0, io.EOF
	}
	nextRune := rr.Src[rr.Pos]
	rr.Pos++
	return nextRune, 1, nil
}
