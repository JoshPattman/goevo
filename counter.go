package goevo

// Counter is a simple counter that can be used to generate unique IDs.
type Counter struct {
	n int
}

// NewCounter creates a new counter, starting at 0.
func NewCounter() *Counter {
	return &Counter{0}
}

// Next returns the next value of the counter
//
// TODO(make this thread-safe)
func (c *Counter) Next() int {
	c.n++
	return c.n
}
