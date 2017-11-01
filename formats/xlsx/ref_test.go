package xlsx_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/formats/xlsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRef(t *testing.T) {
	for i, d := range []struct {
		Text string
		Col  int
		Row  int
	}{
		{"A1", 0, 0},
		{"Z9", 25, 8},
		{"AA1", 26, 0},
		{"AA99", 26, 98},
		{"ZZ100", 701, 99},
	} {
		ref := xlsx.ParseRef(d.Text)
		assert.Equal(t, d.Col, ref.Col, "column for index %d: %s", i, d.Text)
		assert.Equal(t, d.Row, ref.Row, "row for index %d: %s", i, d.Text)
		assert.Equal(t, d.Text, ref.String(), "String() for index %d: %s", i, d.Text)
	}

	for r := 0; r < 100; r++ {
		for c := 0; c < 10000; c++ {
			in := xlsx.Ref{Row: r, Col: c}
			out := xlsx.ParseRef(in.String())
			require.Equal(t, in, out)
		}
	}
}
