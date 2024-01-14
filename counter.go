package goevo

import "sync/atomic"

// An implementation of the Counter interface which is goroutine safe
type Counter struct {
	I int64
}

// Get the next innovation id
func (a *Counter) Next() int {
	v := atomic.AddInt64(&a.I, 1)
	return int(v)
}

// Create a new AtomicCounter which starts from id 0
func NewCounter() *Counter {
	return &Counter{
		I: -1,
	}
}
