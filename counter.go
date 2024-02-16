package goevo

type Counter struct {
	n int
}

func NewCounter() *Counter {
	return &Counter{0}
}

func (c *Counter) Next() int { // TODO: make this thread-safe
	c.n++
	return c.n
}
