package softref_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/softref"
	"github.com/stretchr/testify/assert"
)

type res struct {
	key      string
	released chan<- string
}

func newRes(key string, released chan<- string) *res {
	return &res{
		key:      key,
		released: released,
	}
}

func (r *res) Key() string {
	return r.key
}

func (r *res) Release() {
	r.released <- r.key
}

func TestSoftRef(t *testing.T) {
	p := softref.NewPool(&jot.Logger{})
	ch := make(chan string, 128)
	sr1, existed := p.NewSoftRef(newRes("1", ch))
	assert.False(t, existed)
	_, existed = p.NewSoftRef(newRes("2", ch))
	assert.False(t, existed)
	sr3, existed := p.NewSoftRef(newRes("3", ch))
	assert.False(t, existed)
	r4 := newRes("4", ch)
	sr4a, existed := p.NewSoftRef(r4)
	assert.False(t, existed)
	sr4b, existed := p.NewSoftRef(r4)
	assert.True(t, existed)
	lookfor(t, "2", ch)
	key := sr3.Resource.(*res).key
	lookfor(t, key, ch)
	key = sr1.Resource.(*res).key
	lookfor(t, key, ch)
	key = sr4a.Resource.(*res).key
	assert.Equal(t, key, sr4b.Resource.(*res).key)
	lookforExpectingTimeout(t, ch)
	assert.Equal(t, "4", sr4b.Key) // Keeps refs to r4 alive for the above call
	lookfor(t, key, ch)
}

func lookfor(t *testing.T, key string, ch <-chan string) {
	t.Helper()
	runtime.GC()
	select {
	case <-time.After(time.Second):
		assert.Failf(t, "timed out waiting for %s", key)
	case k := <-ch:
		assert.Equal(t, key, k)
	}
}

func lookforExpectingTimeout(t *testing.T, ch <-chan string) {
	t.Helper()
	runtime.GC()
	select {
	case <-time.After(time.Second):
	case k := <-ch:
		assert.Failf(t, "received key '%s' when none expected", k)
	}
}