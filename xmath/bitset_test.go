package xmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitSet(t *testing.T) {
	var bs BitSet
	assert.Equal(t, 0, bs.Count())
	bs.Set(0)
	assert.Equal(t, 1, bs.Count())
	bs.Set(7)
	assert.Equal(t, 2, bs.Count())
	bs.Set(dataBitsPerWord - 1)
	assert.Equal(t, 3, bs.Count())
	bs.Set(dataBitsPerWord)
	assert.Equal(t, 4, bs.Count())
	bs.Set(dataBitsPerWord + 1)
	assert.Equal(t, 5, bs.Count())
	bs.Set(0)
	assert.Equal(t, 5, bs.Count())
	bs.Clear(0)
	assert.Equal(t, 4, bs.Count())
	bs.Clear(1)
	assert.Equal(t, 4, bs.Count())
	bs.Clear(1000)
	assert.Equal(t, 4, bs.Count())
	assert.False(t, bs.State(0))
	assert.False(t, bs.State(1))
	assert.True(t, bs.State(7))
	assert.False(t, bs.State(77))
	assert.True(t, bs.State(dataBitsPerWord))
	bs.Flip(22)
	assert.True(t, bs.State(22))
	bs.Flip(22)
	assert.False(t, bs.State(22))
	assert.Equal(t, 7, bs.NextSet(0))
	assert.Equal(t, 7, bs.NextSet(7))
	assert.Equal(t, dataBitsPerWord-1, bs.NextSet(8))
	assert.Equal(t, dataBitsPerWord, bs.NextSet(dataBitsPerWord))
	bs.Set(1234)
	assert.Equal(t, 1234, bs.NextSet(dataBitsPerWord+2))
	assert.Equal(t, 0, bs.NextClear(0))
	assert.Equal(t, dataBitsPerWord+2, bs.NextClear(dataBitsPerWord-1))
	assert.Equal(t, 1235, bs.NextClear(1234))
	bs.Set(dataBitsPerWord*100 - 1)
	assert.Equal(t, dataBitsPerWord*100, bs.NextClear(dataBitsPerWord*100-1))
	assert.Equal(t, dataBitsPerWord*100-1, bs.PreviousSet(dataBitsPerWord*100))
	assert.Equal(t, 1234, bs.PreviousSet(dataBitsPerWord*100-2))
	assert.Equal(t, -1, bs.PreviousSet(0))
	assert.Equal(t, dataBitsPerWord*1000, bs.PreviousClear(dataBitsPerWord*1000))
	assert.Equal(t, dataBitsPerWord*100-2, bs.PreviousClear(dataBitsPerWord*100-1))
	assert.Equal(t, 0, bs.PreviousClear(0))
	bs.Set(0)
	assert.Equal(t, -1, bs.PreviousClear(0))

	bs.Reset()
	bs.Set(65)
	bs.SetRange(10, 300)
	assert.Equal(t, 291, bs.Count())
	for i := 10; i < 301; i++ {
		assert.True(t, bs.State(i))
	}
	assert.Equal(t, 301, bs.NextClear(10))
	assert.Equal(t, 9, bs.PreviousClear(300))
	assert.Equal(t, 10, bs.NextSet(0))
	assert.Equal(t, 300, bs.PreviousSet(1000))
	bs.ClearRange(15, 295)
	assert.Equal(t, 10, bs.Count())
	for i := 15; i < 296; i++ {
		assert.False(t, bs.State(i))
	}
	bs.FlipRange(10, 300)
	assert.Equal(t, 281, bs.Count())
}
