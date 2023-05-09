package goevo

// The zeroth dimension is the layer, the following dimensions are just used for positional info
type HyperNode struct {
	Position    []float64
	accumulator float64
}

type HyperNEATPhenotype struct {
	Substrate []*HyperNode
}
