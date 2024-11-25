package goevo

import (
	"sync"
)

// Implements a memory pool, similar to sync.Pool.
// However, this pool is able to create and manage objects with different shapes.
// This means that, for example, we can use the pool to create both arrays of length x and also length y.
// The length x arrays will only ever be re-used to create new other length x arrays.
type Pool[T any, U comparable] struct {
	// The maximum combined size of all objects in the pool
	maxSize int
	// The old object queues
	oldObjects map[U]chan T
	// The lock
	oldObjectLock *sync.Mutex
	// The constructor
	constructor func(U) T
	// The current number of stored objects
	stored int
	// The number of new allocations
	allocated int
}

func NewPool[T any, U comparable](maxSize int, constructor func(U) T) *Pool[T, U] {
	if maxSize < 1 {
		panic("must have max size larger than 0")
	}
	if constructor == nil {
		panic("cannot have nil constructor")
	}
	return &Pool[T, U]{
		maxSize:       maxSize,
		oldObjects:    make(map[U]chan T),
		oldObjectLock: &sync.Mutex{},
		constructor:   constructor,
	}
}

// Gets an object from the pool with the specified create params
func (p *Pool[T, U]) Get(params U) T {
	var obj T
	reused := false

	p.oldObjectLock.Lock()
	p.ensureOldObjects(params)
	if len(p.oldObjects[params]) > 0 {
		obj = <-p.oldObjects[params]
		reused = true
		p.stored--
	} else {
		// We do this here to be thread safe, knowing that we will allocate in a second
		p.allocated++
	}
	p.oldObjectLock.Unlock()

	if !reused {
		// We do the constructor out of the lock because it might take a long time
		obj = p.constructor(params)
	}
	return obj
}

// Puts an old object (that was created with params) on the queue. You must stop using the object after this!
func (p *Pool[T, U]) Put(params U, old T) {
	p.oldObjectLock.Lock()
	// Make sure we have a target q
	p.ensureOldObjects(params)
	// Remove items until we are at the max size-1
	for p.stored >= p.maxSize {
		var maxShape U
		maxCount := -1
		for shape, q := range p.oldObjects {
			lq := len(q)
			if lq > maxCount {
				maxCount = lq
				maxShape = shape
			}
		}
		<-p.oldObjects[maxShape]
	}
	// Put new item on q
	p.oldObjects[params] <- old
	p.stored++
	p.oldObjectLock.Unlock()
}

func (p *Pool[T, U]) MaxSize() int {
	return p.maxSize
}

func (p *Pool[T, U]) CurrentSize() int {
	return p.stored
}

// Create a queue for  old objects of that constructor.
// This is NOT thread-safe and should be wrapped within a lock
func (p *Pool[T, U]) ensureOldObjects(params U) {
	if _, ok := p.oldObjects[params]; !ok {
		// This will never fill up to it's max size if there are more than one type of param in the pool
		p.oldObjects[params] = make(chan T, p.maxSize)
	}
}
