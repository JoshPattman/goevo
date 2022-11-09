package goevo

import "sync/atomic"

type Counter interface {
	Next() int
}
type AtomicCounter struct {
	I int64
}

func (a *AtomicCounter) Next() int {
	v := atomic.AddInt64(&a.I, 1)
	return int(v)
}

func NewAtomicCounter() Counter {
	return &AtomicCounter{
		I: -1,
	}
}
