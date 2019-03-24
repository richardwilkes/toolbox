package embedded

import (
	"archive/zip"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// NewFileSystemFromEmbeddedZip creates a new FileSystem from the contents of
// a zip file appended to the end of the executable. If no such data can be
// found, then 'fallbackLiveFSRoot' is used to return a FileSystem based upon
// the local disk.
func NewFileSystemFromEmbeddedZip(fallbackLiveFSRoot string) FileSystem {
	if efs, err := NewEFSFromEmbeddedZip(); err == nil {
		return efs.PrimaryFileSystem()
	}
	return NewLiveFS(fallbackLiveFSRoot)
}

// NewEFSFromEmbeddedZip creates a new EFS from the contents of a zip file
// appended to the end of the executable.
func NewEFSFromEmbeddedZip() (*EFS, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	exeFile, err := os.Open(exePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(exeFile)
	fi, err := exeFile.Stat()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	start := findZipStart(exeFile)
	if start == -1 {
		return nil, errs.New("unknown executable type")
	}
	section := io.NewSectionReader(exeFile, start, fi.Size()-start)
	r, err := zip.NewReader(section, section.Size())
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return NewEFSFromZip(r)
}

// NewEFSFromZip creates a new EFS from the contents of a zip file.
func NewEFSFromZip(zr *zip.Reader) (*EFS, error) {
	files := make(map[string]*File)
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		r, err := f.Open()
		if err != nil {
			return nil, errs.Wrap(err)
		}
		data, err := ioutil.ReadAll(r)
		xio.CloseIgnoringErrors(r)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		name := f.Name
		if !strings.HasPrefix(name, "/") {
			name = "/" + name
		}
		name = filepath.Clean(name)
		files[name] = NewFile(filepath.Base(name), f.Modified, int64(f.UncompressedSize64), false, data)
	}
	return NewEFS(files), nil
}

func findZipStart(r io.ReaderAt) int64 {
	if start := findZipStartForMacho(r); start != -1 {
		return start
	}
	if start := findZipStartForPE(r); start != -1 {
		return start
	}
	return findZipStartForElf(r)
}

func findZipStartForMacho(r io.ReaderAt) int64 {
	f, err := macho.NewFile(r)
	if err != nil {
		return -1
	}
	var max int64
	for _, load := range f.Loads {
		if segment, ok := load.(*macho.Segment); ok {
			end := int64(segment.Offset + segment.Filesz)
			if end > max {
				max = end
			}
		}
	}
	return max
}

func findZipStartForPE(r io.ReaderAt) int64 {
	f, err := pe.NewFile(r)
	if err != nil {
		return -1
	}
	var max int64
	for _, section := range f.Sections {
		end := int64(section.Offset + section.Size)
		if end > max {
			max = end
		}
	}
	return max
}

func findZipStartForElf(r io.ReaderAt) int64 {
	f, err := elf.NewFile(r)
	if err != nil {
		return -1
	}
	var max int64
	for _, section := range f.Sections {
		if section.Type == elf.SHT_NOBITS {
			continue
		}
		end := int64(section.Offset + section.Size)
		if end > max {
			max = end
		}
	}
	return max
}
