package goevo

import "sync/atomic"

// Counter is a type which counts up.
// It is used to generate unique species IDs and innovation IDs
type Counter struct {
	i int64
}

// Next gets the next ID, and increments it for the future.
// It is thread-safe.
func (a *Counter) Next() int {
	v := atomic.AddInt64(&a.i, 1)
	return int(v)
}

// NewCounter creates a new Counter, starting from id 0
func NewCounter() *Counter {
	return &Counter{
		i: -1,
	}
}
