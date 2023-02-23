package goevo

import "sync/atomic"

// An interface that describes a type that can generate unique ids
type Counter interface {
	// Get the next innovation id
	Next() int
}

// An implementation of the Counter interface which is goroutine safe
type AtomicCounter struct {
	I int64
}

// Get the next innovation id
func (a *AtomicCounter) Next() int {
	v := atomic.AddInt64(&a.I, 1)
	return int(v)
}

// Create a new AtomicCounter which starts from id 0
func NewAtomicCounter() Counter {
	return &AtomicCounter{
		I: -1,
	}
}
