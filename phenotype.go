package goevo

// Forwarder is an interface for somthing that can take a set of inputs ([]float64) and return a set of outputs.
type Forwarder interface {
	Forward([]float64) []float64
}
