package softref

import (
	"runtime"
	"sync"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/log/logadapter"
)

// Pool is used to track soft references to resources.
type Pool struct {
	logger logadapter.WarnLogger
	lock   sync.Mutex
	refs   map[string]*softRef
}

// Resource is a resource that will be used with a pool.
type Resource interface {
	// Key returns a unique key for this resource. Must never change.
	Key() string
	// Release is called when the resource is no longer being referenced by
	// any remaining soft references.
	Release()
}

// SoftRef is a soft reference to a given resource.
type SoftRef struct {
	Key      string
	Resource Resource
}

type softRef struct {
	resource Resource
	count    int
}

// DefaultPool is a global default soft reference pool.
var DefaultPool = NewPool(&jot.Logger{})

// NewPool creates a new soft reference pool. 'logger' may be nil.
func NewPool(logger logadapter.WarnLogger) *Pool {
	if logger == nil {
		logger = &logadapter.Discarder{}
	}
	return &Pool{
		logger: logger,
		refs:   make(map[string]*softRef),
	}
}

// NewSoftRef returns a soft reference to the given resource, along with a
// flag indicating if a reference existed previously.
func (p *Pool) NewSoftRef(resource Resource) (ref *SoftRef, existedPreviously bool) {
	key := resource.Key()
	p.lock.Lock()
	defer p.lock.Unlock()
	r := p.refs[key]
	if r != nil {
		r.count++
	} else {
		r = &softRef{
			resource: resource,
			count:    1,
		}
		p.refs[key] = r
	}
	sr := &SoftRef{
		Key:      key,
		Resource: r.resource,
	}
	runtime.SetFinalizer(sr, p.finalizeSoftRef)
	return sr, r.count > 1
}

func (p *Pool) finalizeSoftRef(ref *SoftRef) {
	p.lock.Lock()
	if r, ok := p.refs[ref.Key]; ok {
		r.count--
		if r.count == 0 {
			delete(p.refs, ref.Key)
			r.resource.Release()
		} else if r.count < 0 {
			p.logger.Warnf("internalImageRef for %v finalized but count is now %d", ref.Key, r.count)
		}
	} else {
		p.logger.Warnf("internalImageRef for %v finalized but hash is not present", ref.Key)
	}
	p.lock.Unlock()
}
