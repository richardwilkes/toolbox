package fs_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	for i, one := range []struct {
		in  string
		out []string
	}{
		{
			in:  "/one/two.txt",
			out: []string{"/", "one", "two.txt"},
		},
		{
			in:  "/one",
			out: []string{"/", "one"},
		},
		{
			in:  "one",
			out: []string{".", "one"},
		},
		{
			in:  "/one////two.txt",
			out: []string{"/", "one", "two.txt"},
		},
		{
			in:  "/one//..//two.txt",
			out: []string{"/", "two.txt"},
		},
		{
			in:  "/one/../..//two.txt",
			out: []string{"/", "two.txt"},
		},
		{
			in:  "/one/../..//two.txt/",
			out: []string{"/", "two.txt"},
		},
		{
			in:  "/one/../..//two.txt/.",
			out: []string{"/", "two.txt"},
		},
	} {
		assert.Equal(t, one.out, fs.Split(one.in), "%d", i)
	}
}
