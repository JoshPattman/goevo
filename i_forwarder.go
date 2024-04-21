package goevo

// Forwarder is an interface for somthing that can take a set of inputs ([]float64) and return a set of outputs.
// It can be thought of as a function with a vector input and output.
type Forwarder interface {
	// Forward takes a set of inputs and returns a set of outputs.
	Forward([]float64) []float64
}
