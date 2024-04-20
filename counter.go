package goevo

import (
	"sync/atomic"
)

// Counter is a simple counter that can be used to generate unique IDs.
type Counter struct {
	n int64
}

// NewCounter creates a new counter, starting at 0.
func NewCounter() *Counter {
	return &Counter{0}
}

// Next returns the next value of the counter
func (c *Counter) Next() int {
	return int(atomic.AddInt64(&c.n, 1))
}
