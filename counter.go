package goevo

import "sync/atomic"

// An interface that implements an innovation counter
type Counter interface {
	// Get the next innovation number
	Next() int
}

// An implementation of the Counter interface which is multigoroutine safe
type AtomicCounter struct {
	I int64
}

// Get the next innovation number
func (a *AtomicCounter) Next() int {
	v := atomic.AddInt64(&a.I, 1)
	return int(v)
}

// Create a new AtomicCounter which starts from 0
func NewAtomicCounter() Counter {
	return &AtomicCounter{
		I: -1,
	}
}
