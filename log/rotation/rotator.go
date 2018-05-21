package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
)

const ext = ".log"

// Rotator holds the rotator data.
type Rotator struct {
	path       string
	maxSize    int64
	maxBackups int
	lock       sync.Mutex
	file       *os.File
	size       int64
}

// New creates a new Rotator with the specified options.
func New(options ...func(*Rotator) error) (*Rotator, error) {
	r := &Rotator{
		path:       filepath.Join(os.TempDir(), cmdline.AppCmdName),
		maxSize:    10 * 1024 * 1024,
		maxBackups: 1,
	}
	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Write implements io.Writer.
func (r *Rotator) Write(b []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.file == nil {
		path := r.path + ext
		file, err := os.Open(path)
		if os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return 0, errs.Wrap(err)
			}
			file, err = os.Create(path)
			if err != nil {
				return 0, errs.Wrap(err)
			}
			r.file = file
			r.size = 0
		} else if err != nil {
			return 0, errs.Wrap(err)
		} else {
			r.file = file
			fi, err := r.file.Stat()
			if err != nil {
				return 0, errs.Wrap(err)
			}
			r.size = fi.Size()
		}
	}
	writeSize := int64(len(b))
	if r.size+writeSize > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := r.file.Write(b)
	if err != nil {
		fmt.Println(err)
		fmt.Println(b)
		fmt.Println(r.file)
		err = errs.Wrap(err)
	}
	r.size += int64(n)
	return n, err
}

// Close implements io.Closer.
func (r *Rotator) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.file != nil {
		if err := r.file.Close(); err != nil {
			return errs.Wrap(err)
		}
		r.file = nil
	}
	return nil
}

func (r *Rotator) rotate() error {
	if r.file != nil {
		if err := r.file.Close(); err != nil {
			return errs.Wrap(err)
		}
		r.file = nil
	}
	path := r.path + ext
	if r.maxBackups < 1 {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return errs.Wrap(err)
		}
	} else {
		if err := os.Remove(fmt.Sprintf("%s.%d%s", r.path, r.maxBackups, ext)); err != nil && !os.IsNotExist(err) {
			return errs.Wrap(err)
		}
		for i := r.maxBackups; i > 0; i-- {
			var oldPath string
			if i != 1 {
				oldPath = fmt.Sprintf("%s.%d%s", r.path, i-1, ext)
			} else {
				oldPath = path
			}
			if err := os.Rename(oldPath, fmt.Sprintf("%s.%d%s", r.path, i, ext)); err != nil && !os.IsNotExist(err) {
				return errs.Wrap(err)
			}
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return errs.Wrap(err)
	}
	r.file = file
	r.size = 0
	return nil
}
