package rate_test

import (
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/rate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCap(t *testing.T) {
	rl := rate.New(50*1024, time.Second)
	assert.Equal(t, 50*1024, rl.Cap(true))
	sub := rl.New(100 * 1024)
	assert.Equal(t, 100*1024, sub.Cap(false))
	assert.Equal(t, 50*1024, sub.Cap(true))
	sub.SetCap(1024)
	assert.Equal(t, 1024, sub.Cap(true))
	rl.Close()
	assert.True(t, sub.Closed())
	assert.True(t, rl.Closed())
}

func TestUse(t *testing.T) {
	rl := rate.New(100, 100*time.Millisecond)
	endAfter := time.Now().Add(250 * time.Millisecond)
	for endAfter.After(time.Now()) {
		err := <-rl.Use(1)
		require.NoError(t, err)
	}
	assert.Equal(t, 100, rl.LastUsed())
	rl.Close()
	assert.True(t, rl.Closed())
}
