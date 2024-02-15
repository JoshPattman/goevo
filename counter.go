package goevo

type Counter struct {
	n int
}

func (c *Counter) Next() int { // TODO: make this thread-safe
	c.n++
	return c.n
}
